package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave26C_EndToEndLoweringIdempotencyAndRace(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w26ctest\n\ngo 1.24\n")

	// 1. Minimal shims/nodestore with Wave 26A JoinDot + JoinDot3
	mustWrite("shims/nodestore/nodestore.go", `package nodestore

import (
	"sync"
	"sync/atomic"
)

type Node uint32

type NodeRec struct {
	Token    uint8
	Payload  string
	OrigName string
	First    Node
	Next     Node
	Parent   Node
}

var (
	mu            sync.RWMutex
	recs          = make([]NodeRec, 1, 1024)
	JoinDotCount  atomic.Int64
	JoinDot3Count atomic.Int64
	JoinDot4Count atomic.Int64
)

func AllocRaw(rec NodeRec) Node {
	mu.Lock()
	defer mu.Unlock()
	id := Node(len(recs))
	recs = append(recs, rec)
	return id
}

func GetRec(n Node) NodeRec {
	mu.RLock()
	defer mu.RUnlock()
	return recs[n]
}

func MutateRec(n Node, fn func(*NodeRec)) {
	mu.Lock()
	defer mu.Unlock()
	fn(&recs[n])
}

func JoinDot(a, b string) string {
	JoinDotCount.Add(1)
	return a + "." + b
}

func JoinDot3(a, b, c string) string {
	JoinDot3Count.Add(1)
	return a + "." + b + "." + c
}

func JoinDot4(a, b, c, d string) string {
	JoinDot4Count.Add(1)
	return a + "." + b + "." + c + "." + d
}
`)

	// 2. Minimal gen/j2gort with Wave 26B ReleaseOrderedMap, AcquireOrderedMapPtr/Cap, TrimStaleSlabPools,
	// OrderedMap, CustomOrderedMap, OrderedMapGetAny, ArrayList, LinkedHashMap, Stream, etc.
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
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
func CheckNotNullIface[T any](v T) T {
	if isJavaNil(v) {
		ThrowNPE()
	}
	return v
}
func CheckNotNullString(s string) string {
	if s == NullString {
		ThrowNPE()
	}
	return s
}

func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}
func StrVal(v any) string {
	if isJavaNil(v) {
		return "null"
	}
	switch x := v.(type) {
	case string:
		return StrValString(x)
	case *BoxInt:
		return fmt.Sprintf("%d", x.Val)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

type BoxInt struct{ Val int32 }
func (b *BoxInt) IsJavaNil_() bool { return b == nil }
func (b *BoxInt) String() string {
	if b == nil {
		return "null"
	}
	return fmt.Sprintf("%d", b.Val)
}

var smallInts [256]BoxInt
func init() {
	for i := 0; i < 256; i++ {
		smallInts[i].Val = int32(i - 128)
	}
}
func ValueOfInt(v int32) *BoxInt {
	if v >= -128 && v <= 127 {
		return &smallInts[v+128]
	}
	return &BoxInt{Val: v}
}
func UnboxInt(b *BoxInt) int32 {
	if b == nil {
		ThrowNPE()
	}
	return b.Val
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
func BoxNullableVal[T JavaNilChecker](v T) any {
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
func ToAny[T any](v T) any {
	if isJavaNil(v) {
		return nil
	}
	if s, ok := any(v).(string); ok && s == NullString {
		return nil
	}
	return v
}
func CastPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}
func Cast[T any](v any) T {
	if v == nil {
		var z T
		return z
	}
	return v.(T)
}
func AsIface[T any](v any) T {
	if isJavaNil(v) {
		var z T
		return z
	}
	return v.(T)
}
func RefEquals(a, b any) bool {
	return a == b
}

// Lifecycle event recording for LIFO defer verification
var (
	EventMu                 sync.Mutex
	Events                  []string
	SlabAllocCount          atomic.Int64
	SlabReleaseCount        atomic.Int64
	OrderedMapPtrAcquireCnt atomic.Int64
	OrderedMapPtrReleaseCnt atomic.Int64
	OrderedMapValReleaseCnt atomic.Int64
	TrimCalls               []int
)

func RecordEvent(ev string) {
	EventMu.Lock()
	Events = append(Events, ev)
	EventMu.Unlock()
}

func ResetEvents() {
	EventMu.Lock()
	Events = nil
	TrimCalls = nil
	EventMu.Unlock()
}

func SnapshotEvents() ([]string, []int) {
	EventMu.Lock()
	defer EventMu.Unlock()
	ev := append([]string(nil), Events...)
	tc := append([]int(nil), TrimCalls...)
	return ev, tc
}

func TrimStaleSlabPools(keepRecent int) {
	EventMu.Lock()
	TrimCalls = append(TrimCalls, keepRecent)
	Events = append(Events, fmt.Sprintf("trim:%d", keepRecent))
	EventMu.Unlock()
}

type Ext_java_lang_Iterable interface {
	IterItems_() []any
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	Stream() *JStream
}
type Ext_java_util_Set interface {
	Ext_java_util_Collection
	Iterator() Ext_java_util_Iterator
}
type Ext_java_util_List interface {
	Ext_java_util_Collection
}
type Ext_java_util_Iterator interface {
	HasNext() bool
	Next() any
}
type Ext_java_util_Map_Entry interface {
	GetKey() any
	GetValue() any
}
type Ext_java_util_function_Function interface {
	Apply(any) any
}
type Ext_java_util_function_Predicate interface {
	Test(any) bool
}

type fnLambda struct{ fn func(any) any }
func (f fnLambda) Apply(x any) any { return f.fn(x) }
func NewLambda_Ext_java_util_function_Function(fn func(any) any) Ext_java_util_function_Function {
	return fnLambda{fn: fn}
}

type predLambda struct{ fn func(any) bool }
func (p predLambda) Test(x any) bool { return p.fn(x) }
func NewLambda_Ext_java_util_function_Predicate(fn func(any) bool) Ext_java_util_function_Predicate {
	return predLambda{fn: fn}
}

type JStream struct{ items []any }
func (s *JStream) Filter(p Ext_java_util_function_Predicate) *JStream {
	var out []any
	for _, it := range s.items {
		if p.Test(it) {
			out = append(out, it)
		}
	}
	return &JStream{items: out}
}
func (s *JStream) Map(f Ext_java_util_function_Function) *JStream {
	out := make([]any, len(s.items))
	for i, it := range s.items {
		out[i] = f.Apply(it)
	}
	return &JStream{items: out}
}
func (s *JStream) CollectToSlice() []any {
	return append([]any(nil), s.items...)
}

type sliceIter struct {
	items []any
	idx   int
}
func (it *sliceIter) HasNext() bool { return it.idx < len(it.items) }
func (it *sliceIter) Next() any {
	v := it.items[it.idx]
	it.idx++
	return v
}

type IterSlice struct{ Items []any }
func (s *IterSlice) IterItems_() []any            { return s.Items }
func (s *IterSlice) Size() int32                  { return int32(len(s.Items)) }
func (s *IterSlice) Stream() *JStream             { return &JStream{items: append([]any(nil), s.Items...)} }
func (s *IterSlice) Iterator() Ext_java_util_Iterator { return &sliceIter{items: s.Items} }

type JIterCursor struct {
	items []any
	idx   int
}
func NewIterCursor(it Ext_java_lang_Iterable) JIterCursor {
	if it == nil {
		ThrowNPE()
	}
	return JIterCursor{items: it.IterItems_()}
}
func (c *JIterCursor) HasNext() bool { return c.idx < len(c.items) }
func (c *JIterCursor) Next() any {
	v := c.items[c.idx]
	c.idx++
	return v
}

// ArrayList (for Patch 69 []JSType lowering test)
type Ext_java_util_ArrayList struct {
	items []any
}
func NewExt_java_util_ArrayList() *Ext_java_util_ArrayList { return &Ext_java_util_ArrayList{} }
func (l *Ext_java_util_ArrayList) Add_1(x any) bool {
	l.items = append(l.items, x)
	return true
}
func (l *Ext_java_util_ArrayList) Set_2(idx int32, x any) any {
	old := l.items[idx]
	l.items[idx] = x
	return old
}
func (l *Ext_java_util_ArrayList) Get(idx int32) any { return l.items[idx] }
func (l *Ext_java_util_ArrayList) Size() int32       { return int32(len(l.items)) }
func (l *Ext_java_util_ArrayList) IterItems_() []any { return l.items }
func (l *Ext_java_util_ArrayList) Stream() *JStream  { return &JStream{items: append([]any(nil), l.items...)} }

// Ext_java_util_Map & Ext_java_util_LinkedHashMap
type mapEntry struct {
	k, v any
}
func (e mapEntry) GetKey() any   { return e.k }
func (e mapEntry) GetValue() any { return e.v }

type Ext_java_util_Map interface {
	IsJavaNil_() bool
	Size() int32
	IsEmpty() bool
	Clear()
	ContainsKey(any) bool
	Get(any) any
	Put_2(any, any) any
	PutIfAbsent(any, any) any
	PutAll(Ext_java_util_Map)
	ComputeIfAbsent(any, Ext_java_util_function_Function) any
	KeySet_0() Ext_java_util_Set
	Values() Ext_java_util_Collection
	EntrySet() Ext_java_util_Set
	String() string
}

type Ext_java_util_LinkedHashMap struct {
	entries []mapEntry
}
func NewExt_java_util_LinkedHashMap() *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{}
}
func (m *Ext_java_util_LinkedHashMap) IsJavaNil_() bool { return m == nil }
func (m *Ext_java_util_LinkedHashMap) Size() int32      { return int32(len(m.entries)) }
func (m *Ext_java_util_LinkedHashMap) IsEmpty() bool    { return len(m.entries) == 0 }
func (m *Ext_java_util_LinkedHashMap) Clear()           { m.entries = m.entries[:0] }
func (m *Ext_java_util_LinkedHashMap) findIdx(k any) int {
	for i, e := range m.entries {
		if e.k == k {
			return i
		}
	}
	return -1
}
func (m *Ext_java_util_LinkedHashMap) ContainsKey(k any) bool { return m.findIdx(k) >= 0 }
func (m *Ext_java_util_LinkedHashMap) Get(k any) any {
	if i := m.findIdx(k); i >= 0 {
		return m.entries[i].v
	}
	return nil
}
func (m *Ext_java_util_LinkedHashMap) Put_2(k, v any) any {
	if i := m.findIdx(k); i >= 0 {
		old := m.entries[i].v
		m.entries[i].v = v
		return old
	}
	m.entries = append(m.entries, mapEntry{k: k, v: v})
	return nil
}
func (m *Ext_java_util_LinkedHashMap) PutIfAbsent(k, v any) any {
	if i := m.findIdx(k); i >= 0 && m.entries[i].v != nil {
		return m.entries[i].v
	} else if i >= 0 {
		m.entries[i].v = v
		return nil
	}
	m.entries = append(m.entries, mapEntry{k: k, v: v})
	return nil
}
func (m *Ext_java_util_LinkedHashMap) PutAll(other Ext_java_util_Map) {
	for _, raw := range other.EntrySet().IterItems_() {
		e := raw.(Ext_java_util_Map_Entry)
		m.Put_2(e.GetKey(), e.GetValue())
	}
}
func (m *Ext_java_util_LinkedHashMap) ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any {
	if cur := m.Get(k); cur != nil {
		return cur
	}
	nv := fn.Apply(k)
	if nv != nil {
		m.Put_2(k, nv)
	}
	return nv
}
func (m *Ext_java_util_LinkedHashMap) KeySet_0() Ext_java_util_Set {
	items := make([]any, len(m.entries))
	for i, e := range m.entries {
		items[i] = e.k
	}
	return &IterSlice{Items: items}
}
func (m *Ext_java_util_LinkedHashMap) Values() Ext_java_util_Collection {
	items := make([]any, len(m.entries))
	for i, e := range m.entries {
		items[i] = e.v
	}
	return &IterSlice{Items: items}
}
func (m *Ext_java_util_LinkedHashMap) EntrySet() Ext_java_util_Set {
	items := make([]any, len(m.entries))
	for i, e := range m.entries {
		items[i] = e
	}
	return &IterSlice{Items: items}
}
func (m *Ext_java_util_LinkedHashMap) String() string {
	if m == nil {
		ThrowNPE()
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range m.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(e.k))
		sb.WriteByte('=')
		sb.WriteString(StrVal(e.v))
	}
	sb.WriteByte('}')
	return sb.String()
}

func MapContainsKeyStr(m Ext_java_util_Map, k string) bool { return m.ContainsKey(BoxString(k)) }
func MapGetStr(m Ext_java_util_Map, k string) any          { return m.Get(BoxString(k)) }
func MapPutStr(m Ext_java_util_Map, k string, v any) any   { return m.Put_2(BoxString(k), v) }
func MapTryGetComputeAny(m Ext_java_util_Map, k any) (any, bool) {
	v := m.Get(k)
	return v, v != nil
}
func NewMapKeyCursor(m Ext_java_util_Map) JIterCursor { return NewIterCursor(m.KeySet_0()) }
func NewMapValCursor(m Ext_java_util_Map) JIterCursor { return NewIterCursor(m.Values()) }

// OrderedMap[K, V] + Wave 26B pointer pool helpers
type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func NewOrderedMapCap[K comparable, V any](_ int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}
func NewOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func AcquireOrderedMapCap[K comparable, V any](_ int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	OrderedMapPtrAcquireCnt.Add(1)
	return &OrderedMap[K, V]{}
}
func AcquireOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] {
	OrderedMapPtrAcquireCnt.Add(1)
	return &OrderedMap[K, V]{}
}
func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m == nil {
		return
	}
	OrderedMapPtrReleaseCnt.Add(1)
	m.Clear()
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	OrderedMapValReleaseCnt.Add(1)
	m.Clear()
}
func (m *OrderedMap[K, V]) IsNil() bool      { return m == nil }
func (m *OrderedMap[K, V]) IsJavaNil_() bool { return m == nil }
func (m *OrderedMap[K, V]) Size() int32 {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.keys))
}
func (m *OrderedMap[K, V]) IsEmpty() bool { return m.Size() == 0 }
func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		ThrowNPE()
	}
	clear(m.keys)
	clear(m.vals)
	m.keys = m.keys[:0]
	m.vals = m.vals[:0]
	clear(m.idx)
}
func (m *OrderedMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	if m.idx != nil {
		if i, ok := m.idx[k]; ok {
			return m.vals[i], true
		}
	}
	var z V
	return z, false
}
func (m *OrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}
func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	if m == nil {
		ThrowNPE()
	}
	if m.idx == nil {
		m.idx = make(map[K]int)
	}
	if i, ok := m.idx[k]; ok {
		old := m.vals[i]
		m.vals[i] = v
		return old
	}
	m.idx[k] = len(m.keys)
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	var z V
	return z
}
func (m *OrderedMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, ok := m.GetOk(k)
	if !ok || isJavaNil(cur) {
		m.Put(k, v)
		var z V
		return z
	}
	return cur
}

type OrderedMapKeyCursor[K comparable, V any] struct {
	m   *OrderedMap[K, V]
	pos int
}
func (m *OrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	return OrderedMapKeyCursor[K, V]{m: m}
}
func (c *OrderedMapKeyCursor[K, V]) HasNext() bool { return c.pos < len(c.m.keys) }
func (c *OrderedMapKeyCursor[K, V]) Next() K {
	k := c.m.keys[c.pos]
	c.pos++
	return k
}

type OrderedMapValCursor[K comparable, V any] struct {
	m   *OrderedMap[K, V]
	pos int
}
func (m *OrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	return OrderedMapValCursor[K, V]{m: m}
}
func (c *OrderedMapValCursor[K, V]) HasNext() bool { return c.pos < len(c.m.vals) }
func (c *OrderedMapValCursor[K, V]) Next() V {
	v := c.m.vals[c.pos]
	c.pos++
	return v
}

func (m *OrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
	if m == nil {
		ThrowNPE()
	}
	items := make([]any, len(m.keys))
	for i, k := range m.keys {
		items[i] = ToAny(k)
	}
	return &IterSlice{Items: items}
}
func (m *OrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
	if m == nil {
		ThrowNPE()
	}
	items := make([]any, len(m.vals))
	for i, v := range m.vals {
		items[i] = ToAny(v)
	}
	return &IterSlice{Items: items}
}
func (m *OrderedMap[K, V]) ToJavaLinkedHashMap() *Ext_java_util_LinkedHashMap {
	if m == nil {
		return nil
	}
	out := NewExt_java_util_LinkedHashMap()
	for i := range m.keys {
		out.Put_2(ToAny(m.keys[i]), ToAny(m.vals[i]))
	}
	return out
}
func (m *OrderedMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i := range m.keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(ToAny(m.keys[i])))
		sb.WriteByte('=')
		sb.WriteString(StrVal(ToAny(m.vals[i])))
	}
	sb.WriteByte('}')
	return sb.String()
}

func OrderedMapGetAny[K comparable, V any](m *OrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	if rawKey == nil {
		var zero K
		v, ok := m.GetOk(zero)
		if !ok {
			return nil
		}
		return ToAny(v)
	}
	k, ok := rawKey.(K)
	if !ok {
		return nil
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return ToAny(v)
}
func OrderedMapGetNodeAny[K ~uint32, V any](m *OrderedMap[K, V], rawKey any) any {
	return OrderedMapGetAny(m, rawKey)
}

// FastCustomKey + CustomOrderedMap[K, V]
type FastCustomKey interface {
	IsJavaNil_() bool
	HashCode() int32
	Equals(any) bool
}
type FastCustomKeyComparable interface {
	comparable
	FastCustomKey
}

type CustomOrderedMap[K FastCustomKeyComparable, V any] struct {
	keys []K
	vals []V
}
func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func NewCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
	return &CustomOrderedMap[K, V]{}
}
func (m *CustomOrderedMap[K, V]) findIdx(k K) int {
	for i, ek := range m.keys {
		if ek == k || (!ek.IsJavaNil_() && ek.Equals(k)) {
			return i
		}
	}
	return -1
}
func (m *CustomOrderedMap[K, V]) Size() int32 {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.keys))
}
func (m *CustomOrderedMap[K, V]) Clear() {
	if m == nil {
		ThrowNPE()
	}
	m.keys = m.keys[:0]
	m.vals = m.vals[:0]
}
func (m *CustomOrderedMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	if i := m.findIdx(k); i >= 0 {
		return m.vals[i], true
	}
	var z V
	return z, false
}
func (m *CustomOrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}
func (m *CustomOrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}
func (m *CustomOrderedMap[K, V]) Put(k K, v V) V {
	if m == nil {
		ThrowNPE()
	}
	if i := m.findIdx(k); i >= 0 {
		old := m.vals[i]
		m.vals[i] = v
		return old
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	var z V
	return z
}
func (m *CustomOrderedMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, ok := m.GetOk(k)
	if !ok || isJavaNil(cur) {
		m.Put(k, v)
		var z V
		return z
	}
	return cur
}
func (m *CustomOrderedMap[K, V]) PutAllJavaMap(src any) {
	if m == nil {
		ThrowNPE()
	}
	mp := src.(Ext_java_util_Map)
	for _, raw := range mp.EntrySet().IterItems_() {
		e := raw.(Ext_java_util_Map_Entry)
		m.Put(e.GetKey().(K), e.GetValue().(V))
	}
}
func CustomOrderedMapComputeIfAbsent[K FastCustomKeyComparable, V comparable](m *CustomOrderedMap[K, V], k K, nullV V, fn func(K) V) V {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	nv := fn(k)
	if nv != nullV {
		m.Put(k, nv)
	}
	return nv
}
func CustomOrderedMapGetAny[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	k, ok := rawKey.(K)
	if !ok {
		return nil
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return ToAny(v)
}

// OrderedSet + Deque + Slabs
type OrderedSet[T comparable] struct{ items []T }
func NewOrderedSet[T comparable]() OrderedSet[T] { return OrderedSet[T]{} }
func (s *OrderedSet[T]) Add(v T) bool {
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	return true
}
func (s *OrderedSet[T]) Clear() { s.items = s.items[:0] }

type Deque[T any] struct{ items []T }
func (d *Deque[T]) Push(v T)   { d.items = append(d.items, v) }
func (d *Deque[T]) Pop() T     { v := d.items[len(d.items)-1]; d.items = d.items[:len(d.items)-1]; return v }
func (d *Deque[T]) Peek() T    { return d.items[len(d.items)-1] }
func (d *Deque[T]) Clear()     { d.items = d.items[:0] }
func (d *Deque[T]) IsEmpty() bool { return len(d.items) == 0 }

type Slab[T any] struct{ items []*T }
func (s *Slab[T]) Alloc() *T {
	SlabAllocCount.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

type PooledSlab[T any] struct{ items []*T }
func (s *PooledSlab[T]) Alloc() *T {
	SlabAllocCount.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}
func (s *PooledSlab[T]) Reset() {
	for _, it := range s.items {
		var z T
		*it = z
	}
	s.items = s.items[:0]
}
func (s *PooledSlab[T]) Release() {
	SlabReleaseCount.Add(1)
	RecordEvent("slab_release")
	s.Reset()
}
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for _, it := range s.items {
		fn(it)
	}
	s.Release()
}
`)

	// 3. Minimal gen/colors (for ColorId + Color in ColorPool_Builder & ColorGraphNode)
	mustWrite("gen/colors/colors.go", `package colors

import "example.com/w26ctest/gen/j2gort"

type ColorId struct {
	id int32
}

func NewColorId(id int32) *ColorId    { return &ColorId{id: id} }
func (c *ColorId) IsJavaNil_() bool   { return c == nil }
func (c *ColorId) HashCode() int32    { if c == nil { return 0 }; return c.id * 31 }
func (c *ColorId) Equals(o any) bool {
	other, ok := o.(*ColorId)
	return ok && c != nil && other != nil && c.id == other.id
}
func (c *ColorId) String() string { return "cid" }

type Color interface {
	Implements_colors_Color()
	IsJavaNil_() bool
	GetId() *ColorId
}

type SingleColor struct {
	id *ColorId
}
func (*SingleColor) Implements_colors_Color() {}
func (s *SingleColor) IsJavaNil_() bool       { return s == nil }
func (s *SingleColor) GetId() *ColorId        { return s.id }

func NewSingleColor(id *ColorId) Color { return &SingleColor{id: id} }

var axiomatic = func() j2gort.Ext_java_util_Map {
	m := j2gort.NewExt_java_util_LinkedHashMap()
	id0 := NewColorId(1)
	m.Put_2(id0, NewSingleColor(id0))
	return m
}()

func StaticGet_StandardColors_AXIOMATIC_COLORS() j2gort.Ext_java_util_Map {
	return axiomatic
}
`)

	// 4. Minimal gen/graph/standard_union_find.go (Task 3.2: StandardUnionFind[E comparable].elmap)
	mustWrite("gen/graph/standard_union_find.go", `package graph

import "example.com/w26ctest/gen/j2gort"

type StandardUnionFind[E comparable] struct {
	elmap                       j2gort.Ext_java_util_Map
	slab_StandardUnionFind_Node j2gort.Slab[StandardUnionFind_Node[E]]
}

func (this_ *StandardUnionFind[E]) IsJavaNil_() bool    { return this_ == nil }
func (*StandardUnionFind[E]) HasIdentityEquals_()       {}

func (this_ *StandardUnionFind[E]) initFields_StandardUnionFind() {
	this_.elmap = j2gort.NewExt_java_util_LinkedHashMap()
}

func NewStandardUnionFind[E comparable]() *StandardUnionFind[E] {
	this_ := &StandardUnionFind[E]{}
	InitStandardUnionFind[E](this_)
	return this_
}

func InitStandardUnionFind[E comparable](this_ *StandardUnionFind[E]) {
	this_.initFields_StandardUnionFind()
}

func (this_ *StandardUnionFind[E]) Add(e E) {
	this_.Union(e, e)
}

func (this_ *StandardUnionFind[E]) Union(a E, b E) E {
	nodeA := this_.priv_StandardUnionFind_findRootOrCreateNode(a)
	nodeB := this_.priv_StandardUnionFind_findRootOrCreateNode(b)
	if nodeA == nodeB {
		return nodeA.element
	}
	nodeB.parent = nodeA
	nodeA.size += nodeB.size
	return nodeA.element
}

func (this_ *StandardUnionFind[E]) Find(e E) E {
	if !(this_.elmap.ContainsKey(j2gort.ToAny(e))) {
		panic("missing")
	}
	return this_.priv_StandardUnionFind_findRoot(j2gort.CastPtr[StandardUnionFind_Node[E]](this_.elmap.Get(j2gort.ToAny(e)))).element
}

func (this_ *StandardUnionFind[E]) Elements() j2gort.Ext_java_util_Set {
	return this_.elmap.KeySet_0()
}

func (this_ *StandardUnionFind[E]) AllRepresentatives() []any {
	return this_.elmap.Values().Stream().Filter(j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp any) bool {
		n := j2gort.CastPtr[StandardUnionFind_Node[E]](lp)
		return n == n.parent
	})).Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp any) any {
		n := j2gort.CastPtr[StandardUnionFind_Node[E]](lp)
		return j2gort.ToAny(n.element)
	})).CollectToSlice()
}

func (this_ *StandardUnionFind[E]) CountLiveNodes() int {
	cnt := 0
	cur := j2gort.NewMapValCursor(this_.elmap)
	for cur.HasNext() {
		raw := cur.Next()
		elem := j2gort.CastPtr[StandardUnionFind_Node[E]](raw)
		if elem != nil {
			cnt++
		}
	}
	return cnt
}

func (this_ *StandardUnionFind[E]) priv_StandardUnionFind_findRootOrCreateNode(e E) *StandardUnionFind_Node[E] {
	var node *StandardUnionFind_Node[E] = j2gort.CastPtr[StandardUnionFind_Node[E]](this_.elmap.Get(j2gort.ToAny(e)))
	if node != nil {
		return this_.priv_StandardUnionFind_findRoot(node)
	}
	node = this_.slab_StandardUnionFind_Node.Alloc()
	node.parent = node
	node.element = e
	node.size = 1
	this_.elmap.Put_2(j2gort.ToAny(e), j2gort.PtrToAny(node))
	return node
}

func (this_ *StandardUnionFind[E]) priv_StandardUnionFind_findRoot(node *StandardUnionFind_Node[E]) *StandardUnionFind_Node[E] {
	if node.parent != node {
		node.parent = this_.priv_StandardUnionFind_findRoot(node.parent)
	}
	return node.parent
}

func (this_ *StandardUnionFind[E]) SameRootPredicate(value E) *StandardUnionFind__Anon_122[E] {
	return NewStandardUnionFind__Anon_122[E](this_, value)
}

type StandardUnionFind_Node[E comparable] struct {
	parent  *StandardUnionFind_Node[E]
	element E
	size    int32
}

type StandardUnionFind__Anon_122[E comparable] struct {
	outer_       *StandardUnionFind[E]
	cap_0_value  E
	nodeForValue *StandardUnionFind_Node[E]
}

func NewStandardUnionFind__Anon_122[E comparable](outer_ *StandardUnionFind[E], v E) *StandardUnionFind__Anon_122[E] {
	a := &StandardUnionFind__Anon_122[E]{outer_: outer_, cap_0_value: v}
	a.nodeForValue = j2gort.CastPtr[StandardUnionFind_Node[E]](a.outer_.elmap.Get(j2gort.ToAny(a.cap_0_value)))
	return a
}

func (this_ *StandardUnionFind__Anon_122[E]) Apply(arg_21 any) bool {
	b := arg_21
	var nodeForB *StandardUnionFind_Node[E] = j2gort.CastPtr[StandardUnionFind_Node[E]](this_.outer_.elmap.Get(b))
	if nodeForB == nil {
		return false
	}
	this_.nodeForValue = this_.outer_.priv_StandardUnionFind_findRoot(this_.nodeForValue)
	return this_.outer_.priv_StandardUnionFind_findRoot(nodeForB) == this_.nodeForValue
}
`)

	// 5. Minimal gen/rhino/node.go (Task 4: GetQualifiedName / GetOriginalQualifiedName + Patch 69 FunctionType->[]JSType)
	mustWrite("gen/rhino/node.go", `package rhino

import (
	"example.com/w26ctest/gen/j2gort"
	"example.com/w26ctest/shims/nodestore"
)

type Token uint8

const (
	TokSCRIPT   Token = 1
	TokFUNCTION Token = 2
	TokBLOCK    Token = 3
	TokNAME     Token = 33
	TokGETPROP  Token = 27
	TokTHIS     Token = 38
)

type Node nodestore.Node

func (n Node) IsJavaNil_() bool        { return n == 0 }
func (n Node) GetToken() Token         { return Token(nodestore.GetRec(nodestore.Node(n)).Token + 1) }
func (n Node) TokenSwitchOrd_() int32  { return int32(nodestore.GetRec(nodestore.Node(n)).Token) }
func (n Node) TokenOrd_() int32        { return int32(nodestore.GetRec(nodestore.Node(n)).Token) }
func (n Node) IsScript() bool          { return n != 0 && n.GetToken() == TokSCRIPT }
func (n Node) IsFunction() bool        { return n != 0 && n.GetToken() == TokFUNCTION }
func (n Node) IsBlock() bool           { return n != 0 && n.GetToken() == TokBLOCK }
func (n Node) IsName() bool            { return n != 0 && n.GetToken() == TokNAME }
func (n Node) IsGetProp() bool         { return n != 0 && n.GetToken() == TokGETPROP }
func (n Node) GetFirstChild() Node     { return Node(nodestore.GetRec(nodestore.Node(n)).First) }
func (n Node) nodeGet_first() Node     { return n.GetFirstChild() }
func (n Node) GetNext() Node           { return Node(nodestore.GetRec(nodestore.Node(n)).Next) }
func (n Node) GetParent() Node         { return Node(nodestore.GetRec(nodestore.Node(n)).Parent) }
func (n Node) GetString() string       { return nodestore.GetRec(nodestore.Node(n)).Payload }
func (n Node) GetOriginalName() string {
	orig := nodestore.GetRec(nodestore.Node(n)).OrigName
	if orig == "" {
		return j2gort.NullString
	}
	return orig
}

func (n Node) AddChildToBack(c Node) {
	nodestore.MutateRec(nodestore.Node(c), func(r *nodestore.NodeRec) { r.Parent = nodestore.Node(n) })
	first := n.GetFirstChild()
	if first == 0 {
		nodestore.MutateRec(nodestore.Node(n), func(r *nodestore.NodeRec) { r.First = nodestore.Node(c) })
		return
	}
	cur := first
	for cur.GetNext() != 0 {
		cur = cur.GetNext()
	}
	nodestore.MutateRec(nodestore.Node(cur), func(r *nodestore.NodeRec) { r.Next = nodestore.Node(c) })
}

func NewNode_1(tok Token) Node {
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1}))
}
func NewNode_StringNode_2(tok Token, s string) Node {
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1, Payload: s}))
}
func NewGetProp(recv Node, prop string, orig string) Node {
	gp := Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(TokGETPROP) - 1, Payload: prop, OrigName: orig}))
	gp.AddChildToBack(recv)
	return gp
}

func (this_ Node) GetQualifiedName() string {
	switch this_.TokenSwitchOrd_() {
	case int32(32):
		var name string = this_.GetString()
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		}
		return name
	case int32(26):
		var p1 string = this_.GetString()
		var n1 Node = this_.nodeGet_first()
		if !(n1.IsGetProp()) {
			var left string = n1.GetQualifiedName()
			if left == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left) + "." + j2gort.StrValString(p1))
		}
		var p2 string = n1.GetString()
		var n2 Node = n1.nodeGet_first()
		if !(n2.IsGetProp()) {
			var left_2 string = n2.GetQualifiedName()
			if left_2 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_2) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p3 string = n2.GetString()
		var n3 Node = n2.nodeGet_first()
		var left_3 string = n3.GetQualifiedName()
		if left_3 == j2gort.NullString {
			return j2gort.NullString
		}
		return (j2gort.StrValString(left_3) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
	case int32(37):
		return "this"
	default:
		return j2gort.NullString
	}
}

func (this_ Node) GetOriginalQualifiedName() string {
	if this_.TokenOrd_() == int32(32) {
		var name string = this_.GetOriginalName()
		if name == j2gort.NullString {
			name = this_.GetString()
		}
		return name
	} else if this_.TokenOrd_() == int32(26) {
		var left string = this_.GetFirstChild().GetOriginalQualifiedName()
		if left == j2gort.NullString {
			return j2gort.NullString
		}
		var right string = this_.GetOriginalName()
		if right == j2gort.NullString {
			right = this_.GetString()
		}
		return (j2gort.StrValString(left) + "." + j2gort.StrValString(right))
	}
	return j2gort.NullString
}

// Patch 69 check: final class FunctionType (*FunctionType concrete struct) written into ArrayList<JSType> -> []JSType
type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	TypeLabel() string
}

type FunctionType struct {
	label string
}
func (*FunctionType) Implements_rhino_JSType() {}
func (f *FunctionType) IsJavaNil_() bool       { return f == nil }
func (f *FunctionType) TypeLabel() string      { return f.label }
func NewFunctionType(label string) *FunctionType { return &FunctionType{label: label} }

func BuildTypeAlternates(base JSType, maybeFn *FunctionType) string {
	var alts *j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	alts.Add_1(base)
	alts.Add_1(j2gort.PtrToAny(maybeFn))
	first := j2gort.AsIface[JSType](alts.Get(0))
	second := j2gort.AsIface[JSType](alts.Get(1))
	if second == nil {
		return first.TypeLabel() + ":nil"
	}
	return first.TypeLabel() + ":" + second.TypeLabel()
}
`)

	// 6. Minimal gen/jscomp/scope.go + node_traversal.go + wave26c_passes.go
	mustWrite("gen/jscomp/scope.go", `package jscomp

import (
	"example.com/w26ctest/gen/j2gort"
	pkg_rhino "example.com/w26ctest/gen/rhino"
)

type AbstractCompiler interface {
	Implements_jscomp_AbstractCompiler()
	IsJavaNil_() bool
	ReportCodeChange()
}

type Compiler struct {
	passCount int
}
func (*Compiler) Implements_jscomp_AbstractCompiler() {}
func (c *Compiler) IsJavaNil_() bool                  { return c == nil }
func (c *Compiler) ReportCodeChange()                 { c.passCount++ }
func (this_ *Compiler) AfterPass(name string) {
	j2gort.RecordEvent("afterPass:" + name)
}

func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	defer func() {
		this_.AfterPass("parseInputs")
	}()
	j2gort.RecordEvent("parseBody")
	return pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
}

type AbstractVar interface {
	Implements_jscomp_AbstractVar()
	IsJavaNil_() bool
	GetName() string
	GetScope() AbstractScope
}

type Var struct {
	name  string
	node  pkg_rhino.Node
	scope *Scope
}

func (*Var) Implements_jscomp_AbstractVar() {}
func (v *Var) IsJavaNil_() bool             { return v == nil }
func (v *Var) HashCode() int32              { if v == nil { return 0 }; return int32(len(v.name)) }
func (v *Var) Equals(o any) bool            { ov, ok := o.(*Var); return ok && v == ov }
func (v *Var) GetName() string              { return v.name }
func (v *Var) GetNode() pkg_rhino.Node      { return v.node }
func (v *Var) GetNameNode() pkg_rhino.Node  { return v.node }
func (v *Var) GetDeclaration() *Var         { return v }
func (v *Var) GetScope() AbstractScope      { return v.scope }

func NewVar(name string, node pkg_rhino.Node, scope *Scope) *Var {
	v := &Var{}
	InitVar(v, name, node, scope)
	return v
}
func InitVar(v *Var, name string, node pkg_rhino.Node, scope *Scope) {
	v.name = name
	v.node = node
	v.scope = scope
}

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetRootNode() pkg_rhino.Node
	Untyped() *Scope
}

type AbstractScopeBase struct {
	rootNode pkg_rhino.Node
}
func (*AbstractScopeBase) Implements_jscomp_AbstractScope() {}
func (s *AbstractScopeBase) GetRootNode() pkg_rhino.Node    { return s.rootNode }
func (s *AbstractScopeBase) GetRootNode_0() pkg_rhino.Node  { return s.rootNode }

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
}

func (s *Scope) IsJavaNil_() bool                 { return s == nil }
func (s *Scope) Untyped() *Scope                  { return s }
func (s *Scope) GetParent() AbstractScope         { if s.parent == nil { return nil }; return s.parent }
func (s *Scope) GetClosestContainerScope() *Scope { return s }
func (s *Scope) GetClosestHoistScope() *Scope     { return s }
func (s *Scope) IsGlobal() bool                   { return s.parent == nil }

func NewScope(parent *Scope, root pkg_rhino.Node) *Scope {
	s := &Scope{}
	InitScope(s, parent, root)
	return s
}
func InitScope(s *Scope, parent *Scope, root pkg_rhino.Node) {
	s.parent = parent
	s.rootNode = root
}

func (this_ *Scope) Declare(name string, n pkg_rhino.Node, v *Var) *Var {
	if this_.vars != nil {
		this_.vars.Put(name, v)
		return v
	}
	switch {
	case this_.var0 == nil:
		this_.var0 = v
	case this_.var1 == nil:
		this_.var1 = v
	case this_.var2 == nil:
		this_.var2 = v
	case this_.var3 == nil:
		this_.var3 = v
	case this_.var4 == nil:
		this_.var4 = v
	case this_.var5 == nil:
		this_.var5 = v
	default:
		this_.vars = j2gort.NewOrderedMapPtrCap[string, *Var](int32(16))
		for _, old := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
			this_.vars.Put(old.name, old)
		}
		this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5 = nil, nil, nil, nil, nil, nil
		this_.vars.Put(name, v)
	}
	return v
}

func (this_ *Scope) GetVar(name string) *Var {
	if this_.vars != nil {
		return this_.vars.Get(name)
	}
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil && v.name == name {
			return v
		}
	}
	if this_.parent != nil {
		return this_.parent.GetVar(name)
	}
	return nil
}

func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if this_.vars != nil {
		return this_.vars.LiveValues()
	}
	var items []any
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil {
			items = append(items, v)
		}
	}
	return &j2gort.IterSlice{Items: items}
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope(pkg_rhino.Node, AbstractScope) AbstractScope
}

type SyntacticScopeCreator_RedeclarationHandler interface{}
type defaultRedeclHandler struct{}
func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	return defaultRedeclHandler{}
}

type SyntacticScopeCreator_ScopeScanner struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
}
func (s *SyntacticScopeCreator_ScopeScanner) Reset() {}

type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	reusableScanner               *SyntacticScopeCreator_ScopeScanner
	treatProvidesAsRedeclarations bool
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
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{compiler: compiler, redeclarationHandler: this_.redeclarationHandler}
}

func (this_ *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	var p *Scope
	if parent != nil && !parent.IsJavaNil_() {
		p = parent.Untyped()
	}
	sc := NewScope(p, n)
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		if c.IsName() {
			sc.Declare(c.GetString(), c, NewVar(c.GetString(), c, sc))
		}
	}
	return sc
}

type MemoizedScopeCreator struct {
	delegate          ScopeCreator
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope]
}
func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }
func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	m := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator(m, delegate)
	return m
}
func InitMemoizedScopeCreator(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	this_.delegate = delegate
	this_.scopesByScopeRoot = j2gort.NewOrderedMap[pkg_rhino.Node, AbstractScope]()
}
func (this_ *MemoizedScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	if existing := this_.scopesByScopeRoot.Get(n); existing != nil {
		return existing
	}
	sc := this_.delegate.CreateScope(n, parent)
	this_.scopesByScopeRoot.Put(n, sc)
	return sc
}
`)

	mustWrite("gen/jscomp/node_traversal.go", `package jscomp

import (
	pkg_rhino "example.com/w26ctest/gen/rhino"
)

type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal_AbstractPostOrderCallbackBase struct{}
func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (*NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}

type NodeTraversal_AbstractScopedCallbackBase struct{}
func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractScopedCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (*NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (*NodeTraversal_AbstractScopedCallbackBase) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}

func clinit_NodeTraversal() {}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder { return &NodeTraversal_Builder{} }
func (this_ *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	this_.compiler = c
	return this_
}
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	this_.callback = cb
	return this_
}
func (this_ *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	this_.scopeCreator = sc
	return this_
}
func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	nt := &NodeTraversal{}
	InitNodeTraversal(nt, this_)
	return nt
}
func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_1(root)
}
func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverseRoots_2(externs, root)
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	curNode                                        pkg_rhino.Node
	scopeNodes                                     []pkg_rhino.Node
	scopeObjs                                      []AbstractScope
	scopeDepth                                     int32
	cfgs                                           any
	scopeCreator                                   ScopeCreator
	scopedCallback                                 NodeTraversal_ScopedCallback
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func InitNodeTraversal(this_ *NodeTraversal, b *NodeTraversal_Builder) {
	this_.compiler = b.compiler
	this_.callback = b.callback
	this_.obeyDestructuringAndDefaultValueExecutionOrder = b.obeyDestructuringAndDefaultValueExecutionOrder
	if sc, ok := b.callback.(NodeTraversal_ScopedCallback); ok {
		this_.scopedCallback = sc
	}
	if b.scopeCreator == nil || b.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = b.scopeCreator
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_1() {
	if int(this_.scopeDepth) == len(this_.scopeNodes) {
		n := len(this_.scopeNodes)*2 + 4
		ns := make([]pkg_rhino.Node, n)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		os := make([]AbstractScope, n)
		copy(os, this_.scopeObjs)
		this_.scopeObjs = os
	}
}

func (this_ *NodeTraversal) pushScopeNode(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_1()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if this_.scopedCallback != nil {
		this_.scopedCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) popScope() {
	if this_.scopedCallback != nil {
		this_.scopedCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) GetScope() *Scope {
	var sc AbstractScope
	for i := int(this_.scopeDepth) - 1; i >= 0; i-- {
		if this_.scopeObjs[i] != nil {
			sc = this_.scopeObjs[i]
			for j := i + 1; j < int(this_.scopeDepth); j++ {
				sc = this_.scopeCreator.CreateScope(this_.scopeNodes[j], sc)
				this_.scopeObjs[j] = sc
			}
			return sc.Untyped()
		}
	}
	for j := 0; j < int(this_.scopeDepth); j++ {
		sc = this_.scopeCreator.CreateScope(this_.scopeNodes[j], sc)
		this_.scopeObjs[j] = sc
	}
	return sc.Untyped()
}

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node { return this_.scopeNodes[this_.scopeDepth-1] }
func (this_ *NodeTraversal) ReportCodeChange_0()          { this_.compiler.ReportCodeChange() }

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_1(root pkg_rhino.Node) {
	this_.pushScopeNode(root)
	this_.priv_NodeTraversal_traverseBranch_2(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.pushScopeNode(root)
	if externs != 0 {
		this_.priv_NodeTraversal_traverseBranch_2(externs, 0)
	}
	this_.priv_NodeTraversal_traverseBranch_2(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_2(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.curNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	createsScope := n.IsBlock()
	if createsScope {
		this_.pushScopeNode(n)
	}
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.priv_NodeTraversal_traverseBranch_2(c, n)
	}
	this_.curNode = n
	this_.callback.Visit_3(this_, n, parent)
	if createsScope {
		this_.popScope()
	}
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}
`)

	mustWrite("gen/jscomp/wave26c_passes.go", `package jscomp

import (
	pkg_colors "example.com/w26ctest/gen/colors"
	pkg_graph "example.com/w26ctest/gen/graph"
	"example.com/w26ctest/gen/j2gort"
	pkg_rhino "example.com/w26ctest/gen/rhino"
)

// --- Task 2a: VariableReferenceCheck (LIFO defer TrimStaleSlabPools(2) AFTER ReleaseSlabs_()) + Patch 69 4-slot ReferenceCollection ---
type Reference struct {
	node pkg_rhino.Node
}
func NewReference(n pkg_rhino.Node) *Reference {
	r := &Reference{}
	InitReference(r, n)
	return r
}
func InitReference(r *Reference, n pkg_rhino.Node) { r.node = n }

type ReferenceCollection struct {
	ref0 *Reference
	ref1 *Reference
	ref2 *Reference
	ref3 *Reference
	rest []*Reference
}
func NewReferenceCollection() *ReferenceCollection {
	rc := &ReferenceCollection{}
	InitReferenceCollection(rc)
	return rc
}
func InitReferenceCollection(_ *ReferenceCollection) {}
func (rc *ReferenceCollection) Add(r *Reference) {
	switch {
	case rc.ref0 == nil:
		rc.ref0 = r
	case rc.ref1 == nil:
		rc.ref1 = r
	case rc.ref2 == nil:
		rc.ref2 = r
	case rc.ref3 == nil:
		rc.ref3 = r
	default:
		rc.rest = append(rc.rest, r)
	}
}

type ReferenceCollector struct {
	slab_Reference           j2gort.PooledSlab[Reference]
	slab_ReferenceCollection j2gort.PooledSlab[ReferenceCollection]
}
func (rc *ReferenceCollector) Collect(n pkg_rhino.Node) {
	col := rc.slab_ReferenceCollection.Alloc()
	for i := 0; i < 6; i++ {
		ref := rc.slab_Reference.Alloc()
		ref.node = n
		col.Add(ref)
	}
}
func (rc *ReferenceCollector) ReleaseSlabs_() {
	rc.slab_ReferenceCollection.ReleaseWith(func(c *ReferenceCollection) { c.rest = nil })
	rc.slab_Reference.Release()
}

type VariableReferenceCheck struct {
	compiler AbstractCompiler
}
func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	{
		var val_1 ReferenceCollector
		defer func() {
			(&val_1).ReleaseSlabs_()
		}()
		(&val_1).Collect(root)
	}
}

// --- Task 2b: Normalize_NormalizeStatements delegating to MakeDeclaredNamesUnique + RewriteLogicalAssignmentOperatorsHelper ---
type MakeDeclaredNamesUnique struct {
	NodeTraversal_AbstractScopedCallbackBase
	seenNames []string
}
func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur.HasNext() {
		raw := cur.Next()
		v := j2gort.CastPtr[Var](raw)
		this_.seenNames = append(this_.seenNames, v.GetName())
	}
}
func (this_ *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	_ = t.GetScopeRoot()
}
func (this_ *MakeDeclaredNamesUnique) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}
func (this_ *MakeDeclaredNamesUnique) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := t.GetScope().GetVar(n.GetString()); v != nil {
			this_.seenNames = append(this_.seenNames, v.GetName())
		}
	}
}

type RewriteLogicalAssignmentOperatorsHelper struct{}
func (h *RewriteLogicalAssignmentOperatorsHelper) VisitLogicalAssignmentOperator(t *NodeTraversal, n pkg_rhino.Node) {
	_ = n
	_ = t.GetScopeRoot()
	t.ReportCodeChange_0()
}

type Normalize_NormalizeStatements struct {
	compiler                                AbstractCompiler
	makeDeclaredNamesUnique                 *MakeDeclaredNamesUnique
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper
}
func (*Normalize_NormalizeStatements) Implements_jscomp_NodeTraversal_Callback()       {}
func (*Normalize_NormalizeStatements) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (n *Normalize_NormalizeStatements) IsJavaNil_() bool                              { return n == nil }
func (this_ *Normalize_NormalizeStatements) EnterScope(t *NodeTraversal) {
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.EnterScope(t)
	}
}
func (this_ *Normalize_NormalizeStatements) ExitScope(t *NodeTraversal) {
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.ExitScope(t)
	}
}
func (this_ *Normalize_NormalizeStatements) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	if this_.makeDeclaredNamesUnique != nil {
		var unused bool = this_.makeDeclaredNamesUnique.ShouldTraverse_3(t, n, parent)
		_ = unused
	}
	return true
}
func (this_ *Normalize_NormalizeStatements) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.Visit_3(t, n, parent)
	}
	this_.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, n)
}

// --- Task 2b: CheckMissingRequires using v.GetDeclaration() ---
type CheckMissingRequires struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	declNodes []pkg_rhino.Node
}
func (this_ *CheckMissingRequires) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		v := t.GetScope().GetVar(n.GetString())
		if v != nil && v.GetDeclaration() != nil {
			this_.declNodes = append(this_.declNodes, v.GetDeclaration().GetNode())
		}
	}
}

// --- Negative soundness check: EscapingVarCallback MUST NOT get CanReleaseSyntacticScopes_() ---
type EscapingVarCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	escaped []*Var
}
func (this_ *EscapingVarCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := t.GetScope().GetVar(n.GetString()); v != nil {
			this_.escaped = append(this_.escaped, v)
		}
	}
}

// --- Task 2b Pass-Boundary Releases: RenameVars, ContextualRenameInverter, GatherModuleMetadata ---
type RenameVars struct {
	compiler                       AbstractCompiler
	localBleedingFunctions         j2gort.OrderedSet[*Var]
	localBleedingFunctionsPerScope j2gort.OrderedMap[*Scope, *Var]
	renamedCount                   int
}
func (this_ *RenameVars) GetVariableMap() int { return this_.renamedCount }
func (this_ *RenameVars) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_Traverse(this_.compiler, root, NewRenameVars_ProcessVars(this_))
}

type RenameVars_ProcessVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RenameVars
}
func NewRenameVars_ProcessVars(outer *RenameVars) *RenameVars_ProcessVars {
	return &RenameVars_ProcessVars{outer_: outer}
}
func (this_ *RenameVars_ProcessVars) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		sc := t.GetScope()
		if v := sc.GetVar(n.GetString()); v != nil {
			this_.outer_.localBleedingFunctions.Add(v)
			this_.outer_.localBleedingFunctionsPerScope.Put(sc, v)
			this_.outer_.renamedCount++
		}
	}
}

type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	NodeTraversal_AbstractScopedCallbackBase
	compiler               AbstractCompiler
	scopeContextStack      j2gort.Deque[*Scope]
	variablesStackedByName j2gort.OrderedMap[string, *Var]
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	this_.scopeContextStack.Push(t.GetScope())
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(_ *NodeTraversal) {
	_ = this_.scopeContextStack.Pop()
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := t.GetScope().GetVar(n.GetString()); v != nil {
			this_.variablesStackedByName.Put(v.GetName(), v)
		}
	}
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) Process_2(_ pkg_rhino.Node, js pkg_rhino.Node) {
	NodeTraversal_Traverse(this_.compiler, js, this_)
}

type GatherModuleMetadata struct {
	compiler AbstractCompiler
	visited  int
}
func (this_ *GatherModuleMetadata) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_Traverse(this_.compiler, externs, NewGatherModuleMetadata_Finder(this_))
	NodeTraversal_Traverse(this_.compiler, root, NewGatherModuleMetadata_Finder(this_))
}
type GatherModuleMetadata_Finder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *GatherModuleMetadata
	seenScope *Scope
}
func NewGatherModuleMetadata_Finder(o *GatherModuleMetadata) *GatherModuleMetadata_Finder {
	return &GatherModuleMetadata_Finder{outer_: o}
}
func (this_ *GatherModuleMetadata_Finder) Visit_3(t *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {
	this_.seenScope = t.GetScope()
	this_.outer_.visited++
}

// --- Task 3.1a: MakeDeclaredNamesUnique_ContextualRenamer_NameUsage.counts ---
type MakeDeclaredNamesUnique_ContextualRenamer_NameUsage struct {
	counts *j2gort.Ext_java_util_LinkedHashMap
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) initFields_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage() {
	this_.counts = j2gort.NewExt_java_util_LinkedHashMap()
}
func NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage() *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage {
	u := &MakeDeclaredNamesUnique_ContextualRenamer_NameUsage{}
	InitMakeDeclaredNamesUnique_ContextualRenamer_NameUsage(u)
	return u
}
func InitMakeDeclaredNamesUnique_ContextualRenamer_NameUsage(this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) {
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage()
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) Contains(name string) bool {
	return j2gort.MapContainsKeyStr(this_.counts, name)
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) Reserve(name string) {
	this_.counts.PutIfAbsent(j2gort.BoxString(name), j2gort.ValueOfInt(int32(1)))
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) SetZeroForTesting(name string) {
	j2gort.MapPutStr(this_.counts, name, j2gort.ValueOfInt(int32(0)))
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) AddOne(name string) int32 {
	var prev *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](j2gort.MapGetStr(this_.counts, name))
	var old int32
	if prev != nil {
		old = j2gort.UnboxInt(prev)
	} else {
		old = int32(0)
	}
	j2gort.MapPutStr(this_.counts, name, j2gort.ValueOfInt(old+int32(1)))
	return old
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) String() string {
	return this_.counts.String()
}

// --- Task 3.1b: ColorGraphNode.associatedProps + ClusterPropagator ---
type PropertyClustering struct {
	name        string
	invalidated bool
	clusters    *pkg_graph.StandardUnionFind[*ColorGraphNode]
}
func (*PropertyClustering) HasIdentityEquals_()   {}
func (p *PropertyClustering) IsJavaNil_() bool    { return p == nil }
func (p *PropertyClustering) GetName() string     { return p.name }
func (p *PropertyClustering) IsInvalidated() bool { return p.invalidated }
func NewPropertyClustering(name string) *PropertyClustering {
	return &PropertyClustering{
		name:     name,
		clusters: pkg_graph.NewStandardUnionFind[*ColorGraphNode](),
	}
}

type ColorGraphNode_PropAssociation uint8
const (
	PropAssocAST         ColorGraphNode_PropAssociation = 1
	PropAssocTYPE_SYSTEM ColorGraphNode_PropAssociation = 2
	PropAssocSUPERTYPE   ColorGraphNode_PropAssociation = 3
)
func (e ColorGraphNode_PropAssociation) IsJavaNil_() bool                 { return e == 0 }
func (ColorGraphNode_PropAssociation) Implements_j2gort_Ext_java_lang_Enum() {}

type ColorGraphNode struct {
	associatedProps *j2gort.Ext_java_util_LinkedHashMap
	index           int32
}
func (*ColorGraphNode) HasIdentityEquals_() {}
func (c *ColorGraphNode) IsJavaNil_() bool  { return c == nil }
func (this_ *ColorGraphNode) initFields_ColorGraphNode() {
	this_.associatedProps = j2gort.NewExt_java_util_LinkedHashMap()
}
func NewColorGraphNode(idx int32) *ColorGraphNode {
	c := &ColorGraphNode{index: idx}
	InitColorGraphNode(c)
	return c
}
func InitColorGraphNode(this_ *ColorGraphNode) {
	this_.initFields_ColorGraphNode()
}
func (this_ *ColorGraphNode) GetAssociatedProps() *j2gort.Ext_java_util_LinkedHashMap {
	return this_.associatedProps
}

func PropagateClusterProps(src, dest *ColorGraphNode) bool {
	var start int32 = dest.GetAssociatedProps().Size()
	cur := j2gort.NewMapKeyCursor(src.GetAssociatedProps())
	for cur.HasNext() {
		raw := cur.Next()
		prop := j2gort.CastPtr[PropertyClustering](raw)
		if !prop.IsInvalidated() {
			dest.GetAssociatedProps().PutIfAbsent(j2gort.PtrToAny(prop), j2gort.BoxNullableVal(PropAssocSUPERTYPE))
		}
	}
	return start < dest.GetAssociatedProps().Size()
}

func DumpAssociatedPropsNames(n *ColorGraphNode) []any {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](n.GetAssociatedProps().EntrySet()).Stream().Map(
		j2gort.NewLambda_Ext_java_util_function_Function(func(lp any) any {
			e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp)
			return j2gort.CastPtr[PropertyClustering](e.GetKey()).GetName()
		}),
	).CollectToSlice()
}

// --- Task 3.1c: ColorPool_Builder.idToColor -> NewColorPool / InitColorPool -> ColorPool.idToColor ---
var pendingColor pkg_colors.Color = pkg_colors.NewSingleColor(pkg_colors.NewColorId(-1))

type ColorPool struct {
	idToColor *j2gort.Ext_java_util_LinkedHashMap
}
func (this_ *ColorPool) initFields_ColorPool() {}
func NewColorPool(idToColor *j2gort.Ext_java_util_LinkedHashMap) *ColorPool {
	c := &ColorPool{}
	InitColorPool(c, idToColor)
	return c
}
func InitColorPool(this_ *ColorPool, idToColor *j2gort.Ext_java_util_LinkedHashMap) {
	_ = idToColor
	this_.initFields_ColorPool()
	this_.idToColor = idToColor
}
func (this_ *ColorPool) GetColor(id *pkg_colors.ColorId) pkg_colors.Color {
	return j2gort.AsIface[pkg_colors.Color](this_.idToColor.Get(j2gort.PtrToAny(id)))
}

type ColorPool_Builder struct {
	idToColor *j2gort.Ext_java_util_LinkedHashMap
}
func (this_ *ColorPool_Builder) initFields_ColorPool_Builder() {
	this_.idToColor = j2gort.NewExt_java_util_LinkedHashMap()
}
func NewColorPool_Builder() *ColorPool_Builder {
	b := &ColorPool_Builder{}
	InitColorPool_Builder(b)
	return b
}
func InitColorPool_Builder(this_ *ColorPool_Builder) {
	this_.initFields_ColorPool_Builder()
	this_.idToColor.PutAll(pkg_colors.StaticGet_StandardColors_AXIOMATIC_COLORS())
}
func (this_ *ColorPool_Builder) Build(extraIds []*pkg_colors.ColorId) *ColorPool {
	for _, cid := range extraIds {
		c := func() pkg_colors.Color {
			m_19 := this_.idToColor
			k_20 := j2gort.PtrToAny(cid)
			if v_21, ok_22 := j2gort.MapTryGetComputeAny(m_19, k_20); ok_22 {
				return j2gort.AsIface[pkg_colors.Color](v_21)
			}
			return j2gort.AsIface[pkg_colors.Color](m_19.ComputeIfAbsent(k_20, j2gort.NewLambda_Ext_java_util_function_Function(func(lp any) any {
				_ = j2gort.CastPtr[pkg_colors.ColorId](lp)
				return pkg_colors.NewSingleColor(cid)
			})))
		}()
		_ = c
	}
	reconcileId := pkg_colors.NewColorId(99)
	existing := j2gort.AsIface[pkg_colors.Color](this_.idToColor.PutIfAbsent(j2gort.PtrToAny(reconcileId), pendingColor))
	if existing == nil || existing.IsJavaNil_() {
		this_.idToColor.Put_2(j2gort.PtrToAny(reconcileId), pkg_colors.NewSingleColor(reconcileId))
	}
	var val_28 ColorPool
	InitColorPool(&val_28, this_.idToColor)
	return &val_28
}
`)

	allPasses := passes.AllPasses()
	for _, relPkg := range []string{"gen/graph", "gen/rhino", "gen/jscomp"} {
		pkgDir := filepath.Join(tmpDir, relPkg)
		res1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", relPkg, err)
		}
		if res1.TotalSites == 0 {
			t.Fatalf("expected >0 rewrite sites on %s Run 1, got %+v", relPkg, res1)
		}
		// Run 2 MUST produce 0 sites and 0 modified files (100% idempotency)!
		res2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 on %s failed: %v", relPkg, err)
		}
		if res2.TotalSites != 0 || res2.FilesModified != 0 {
			t.Fatalf("expected 0 sites and 0 modified files on %s Run 2, got %+v", relPkg, res2)
		}
	}

	readLowered := func(rel string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(tmpDir, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	// Verify Task 4 + Patch 69 in gen/rhino/node.go
	rhinoCode := readLowered("gen/rhino/node.go")
	for _, want := range []string{
		"nodestore.JoinDot(",
		"nodestore.JoinDot3(",
		"nodestore.JoinDot4(",
		"var alts []JSType",
	} {
		if !strings.Contains(rhinoCode, want) {
			t.Fatalf("expected gen/rhino/node.go to contain %q:\n%s", want, rhinoCode)
		}
	}
	if strings.Count(rhinoCode, "nodestore.JoinDot(") != 2 ||
		strings.Count(rhinoCode, "nodestore.JoinDot3(") != 1 ||
		strings.Count(rhinoCode, "nodestore.JoinDot4(") != 1 {
		t.Fatalf("expected exact 2 JoinDot, 1 JoinDot3, and 1 JoinDot4 calls, got:\n%s", rhinoCode)
	}

	// Verify Task 3.2 in gen/graph/standard_union_find.go
	normWS := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}
	graphCode := readLowered("gen/graph/standard_union_find.go")
	for _, want := range []string{
		"elmap j2gort.OrderedMap[E, *StandardUnionFind_Node[E]]",
		"j2gort.OrderedMapGetAny(&this_.outer_.elmap, b)",
		"this_.elmap.IterValues()",
	} {
		if !strings.Contains(normWS(graphCode), want) {
			t.Fatalf("expected gen/graph/standard_union_find.go to contain %q:\n%s", want, graphCode)
		}
	}

	// Verify Task 1 + Task 2a in gen/jscomp/scope.go (no mid-compilation TrimStaleSlabPools injection)
	scopeCode := readLowered("gen/jscomp/scope.go")
	for _, want := range []string{
		"j2gort.AcquireOrderedMapPtrCap[string, *Var](int32(16))",
		"j2gort.ReleaseOrderedMap(s_.vars)",
		"this_.scopesByScopeRoot.ReleaseToPool()",
	} {
		if !strings.Contains(normWS(scopeCode), want) {
			t.Fatalf("expected gen/jscomp/scope.go to contain %q:\n%s", want, scopeCode)
		}
	}
	if strings.Contains(scopeCode, "TrimStaleSlabPools(") || strings.Contains(scopeCode, "PurgeAllRuntimePools(") {
		t.Fatalf("expected gen/jscomp/scope.go NOT to contain mid-compilation pool trimming:\n%s", scopeCode)
	}

	// Verify Task 2 + Task 3.1 in gen/jscomp/wave26c_passes.go
	jscompCode := readLowered("gen/jscomp/wave26c_passes.go")
	for _, want := range []string{
		"func (*Normalize_NormalizeStatements) CanReleaseSyntacticScopes_()",
		"func (*CheckMissingRequires) CanReleaseSyntacticScopes_()",
		"this_.localBleedingFunctions.Clear()",
		"this_.localBleedingFunctionsPerScope.Clear()",
		"this_.scopeContextStack.Clear()",
		"this_.variablesStackedByName.Clear()",
		"counts j2gort.OrderedMap[string, *j2gort.BoxInt]",
		"this_.counts.String()",
		"associatedProps j2gort.OrderedMap[*PropertyClustering, ColorGraphNode_PropAssociation]",
		"GetAssociatedProps() *j2gort.OrderedMap[*PropertyClustering, ColorGraphNode_PropAssociation]",
		"ToJavaLinkedHashMap().EntrySet()",
		"idToColor *j2gort.CustomOrderedMap[*pkg_colors.ColorId, pkg_colors.Color]",
		"idToColor j2gort.CustomOrderedMap[*pkg_colors.ColorId, pkg_colors.Color]",
		"j2gort.CustomOrderedMapComputeIfAbsent(",
	} {
		if !strings.Contains(normWS(jscompCode), want) {
			t.Fatalf("expected gen/jscomp/wave26c_passes.go to contain %q:\n%s", want, jscompCode)
		}
	}
	if strings.Contains(jscompCode, "func (*EscapingVarCallback) CanReleaseSyntacticScopes_()") {
		t.Fatalf("EscapingVarCallback MUST NOT receive CanReleaseSyntacticScopes_()!")
	}

	// Write runtime semantic + LIFO defer order + -race test across lowered packages!
	mustWrite("gen/jscomp/wave26c_rt_test.go", `package jscomp

import (
	"reflect"
	"sync"
	"testing"

	pkg_colors "example.com/w26ctest/gen/colors"
	"example.com/w26ctest/gen/j2gort"
	pkg_rhino "example.com/w26ctest/gen/rhino"
	"example.com/w26ctest/shims/nodestore"
)

func buildDeepScopeAST() (pkg_rhino.Node, pkg_rhino.Node) {
	ext := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	ext.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "ext1"))

	root := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	blk := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
	// >6 vars in blk so Scope.vars spills to AcquireOrderedMapPtrCap!
	for _, nm := range []string{"v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8"} {
		blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, nm))
	}
	root.AddChildToBack(blk)
	return ext, root
}

func TestWave26CRuntimeSemanticsAndLifoOrder(t *testing.T) {
	// 1. Task 2a LIFO defer order in Compiler.ParseInputs & VariableReferenceCheck.Process_2 (without mid-compilation pool destruction)
	j2gort.ResetEvents()
	comp := &Compiler{}
	_ = comp.ParseInputs()
	ev1, tc1 := j2gort.SnapshotEvents()
	wantEv1 := []string{"parseBody", "afterPass:parseInputs"}
	if !reflect.DeepEqual(ev1, wantEv1) || len(tc1) != 0 {
		t.Fatalf("ParseInputs event mismatch: got events=%v trims=%v, want %v", ev1, tc1, wantEv1)
	}

	j2gort.ResetEvents()
	ext, root := buildDeepScopeAST()
	vrc := &VariableReferenceCheck{compiler: comp}
	vrc.Process_2(ext, root)
	ev2, tc2 := j2gort.SnapshotEvents()
	wantEv2 := []string{"slab_release", "slab_release"}
	if !reflect.DeepEqual(ev2, wantEv2) || len(tc2) != 0 {
		t.Fatalf("VariableReferenceCheck.Process_2 event mismatch: got events=%v trims=%v, want %v", ev2, tc2, wantEv2)
	}

	// 2. Task 1 + Task 2b: Scope.vars pointer acquire + ReleaseOrderedMap on Normalize_NormalizeStatements!
	j2gort.OrderedMapPtrAcquireCnt.Store(0)
	j2gort.OrderedMapPtrReleaseCnt.Store(0)
	mdnu := &MakeDeclaredNamesUnique{}
	norm := &Normalize_NormalizeStatements{
		compiler:                                comp,
		makeDeclaredNamesUnique:                 mdnu,
		rewriteLogicalAssignmentOperatorsHelper: &RewriteLogicalAssignmentOperatorsHelper{},
	}
	NodeTraversal_Traverse(comp, root, norm)
	if j2gort.OrderedMapPtrAcquireCnt.Load() == 0 || j2gort.OrderedMapPtrReleaseCnt.Load() != j2gort.OrderedMapPtrAcquireCnt.Load() {
		t.Fatalf("expected Scope.vars pointer pool acquire == release > 0, got acq=%d rel=%d",
			j2gort.OrderedMapPtrAcquireCnt.Load(), j2gort.OrderedMapPtrReleaseCnt.Load())
	}

	// CheckMissingRequires + EscapingVarCallback + Pass-boundary releases
	cmr := &CheckMissingRequires{}
	NodeTraversal_Traverse(comp, root, cmr)
	if len(cmr.declNodes) != 9 {
		t.Fatalf("expected 9 declNodes in CheckMissingRequires, got %d", len(cmr.declNodes))
	}
	esc := &EscapingVarCallback{}
	NodeTraversal_Traverse(comp, root, esc)
	if len(esc.escaped) != 9 || esc.escaped[8].GetName() != "v8" {
		t.Fatalf("expected EscapingVarCallback to retain live *Var pointers after Traverse, got %+v", esc.escaped)
	}

	rv := &RenameVars{compiler: comp}
	rv.Process_2(ext, root)
	if rv.GetVariableMap() != 9 {
		t.Fatalf("expected RenameVars count 9, got %d", rv.GetVariableMap())
	}
	cri := &MakeDeclaredNamesUnique_ContextualRenameInverter{compiler: comp}
	cri.Process_2(ext, root)
	gmm := &GatherModuleMetadata{compiler: comp}
	gmm.Process_2(ext, root)

	// 3. Task 3.1a: NameUsage.counts including Integer(0) vs null semantics + String()
	nu := NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage()
	nu.Reserve("alpha")
	nu.Reserve("alpha")
	if got := nu.AddOne("alpha"); got != 1 {
		t.Fatalf("expected AddOne(alpha)=1 after Reserve, got %d", got)
	}
	nu.SetZeroForTesting("zeroKey")
	nu.Reserve("zeroKey") // Integer(0) is NOT null, so PutIfAbsent MUST NOT overwrite 0 with 1!
	if got := nu.AddOne("zeroKey"); got != 0 {
		t.Fatalf("Java Integer(0) vs null violated: expected 0, got %d", got)
	}
	if s := nu.String(); s != "{alpha=2, zeroKey=1}" {
		t.Fatalf("unexpected NameUsage.String(): %q", s)
	}

	// 4. Task 3.1b + 3.2: ColorGraphNode.associatedProps + StandardUnionFind[*ColorGraphNode]
	nA := NewColorGraphNode(0)
	nB := NewColorGraphNode(1)
	nC := NewColorGraphNode(2)
	propX := NewPropertyClustering("propX")
	nA.GetAssociatedProps().Put(propX, PropAssocAST)
	if !PropagateClusterProps(nA, nB) || PropagateClusterProps(nA, nB) {
		t.Fatalf("ClusterPropagator PutIfAbsent failed")
	}
	if names := DumpAssociatedPropsNames(nB); len(names) != 1 || names[0] != "propX" {
		t.Fatalf("EntrySet().Stream() mismatch: %+v", names)
	}
	propX.clusters.Union(nA, nB)
	propX.clusters.Add(nC)
	pred := propX.clusters.SameRootPredicate(nA)
	if !pred.Apply(nB) || pred.Apply(nC) || pred.Apply("wrongType") || pred.Apply(nil) {
		t.Fatalf("StandardUnionFind SameRootPredicate failed")
	}
	if propX.clusters.CountLiveNodes() != 3 || len(propX.clusters.AllRepresentatives()) != 2 {
		t.Fatalf("StandardUnionFind counts mismatch")
	}

	// 5. Task 3.1c: ColorPool_Builder -> ColorPool
	cpb := NewColorPool_Builder()
	cid42 := pkg_colors.NewColorId(42)
	cp := cpb.Build([]*pkg_colors.ColorId{cid42})
	if cp.GetColor(pkg_colors.NewColorId(42)) == nil || cp.GetColor(pkg_colors.NewColorId(99)) == nil || cp.GetColor(pkg_colors.NewColorId(1)) == nil {
		t.Fatalf("ColorPool lookup failed across custom ColorId keys")
	}

	// 6. Task 4 + Patch 69: Node.GetQualifiedName / GetOriginalQualifiedName + *FunctionType in []JSType
	nodestore.JoinDotCount.Store(0)
	nodestore.JoinDot3Count.Store(0)
	nodestore.JoinDot4Count.Store(0)
	base := pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "goog")
	gp1 := pkg_rhino.NewGetProp(base, "dom", "origDom")
	gp2 := pkg_rhino.NewGetProp(gp1, "TagName", "")
	gp3 := pkg_rhino.NewGetProp(gp2, "DIV", "")
	if q1 := gp1.GetQualifiedName(); q1 != "goog.dom" {
		t.Fatalf("got q1=%q", q1)
	}
	if q2 := gp2.GetQualifiedName(); q2 != "goog.dom.TagName" {
		t.Fatalf("got q2=%q", q2)
	}
	if q3 := gp3.GetQualifiedName(); q3 != "goog.dom.TagName.DIV" {
		t.Fatalf("got q3=%q", q3)
	}
	if oq := gp1.GetOriginalQualifiedName(); oq != "goog.origDom" {
		t.Fatalf("got oq=%q", oq)
	}
	if nodestore.JoinDotCount.Load() != 2 || nodestore.JoinDot3Count.Load() != 1 || nodestore.JoinDot4Count.Load() != 1 {
		t.Fatalf("expected JoinDot=2 JoinDot3=1 JoinDot4=1, got %d, %d, %d",
			nodestore.JoinDotCount.Load(), nodestore.JoinDot3Count.Load(), nodestore.JoinDot4Count.Load())
	}
	fnT := pkg_rhino.NewFunctionType("fn()")
	if got := pkg_rhino.BuildTypeAlternates(fnT, nil); got != "fn():nil" {
		t.Fatalf("typed-nil hazard in []JSType: got %q", got)
	}
	if got := pkg_rhino.BuildTypeAlternates(fnT, pkg_rhino.NewFunctionType("cb()")); got != "fn():cb()" {
		t.Fatalf("got %q", got)
	}

	// 7. Concurrent -race stress test on pooled Scope.vars + NodeTraversal + JoinDot
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localComp := &Compiler{}
			e, r := buildDeepScopeAST()
			_ = e
			cb := &Normalize_NormalizeStatements{
				compiler:                                localComp,
				makeDeclaredNamesUnique:                 &MakeDeclaredNamesUnique{},
				rewriteLogicalAssignmentOperatorsHelper: &RewriteLogicalAssignmentOperatorsHelper{},
			}
			for i := 0; i < 20; i++ {
				NodeTraversal_Traverse(localComp, r, cb)
				_ = gp2.GetQualifiedName()
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-race", "-v", "./...")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lowered runtime test with -race failed: %v\n%s", err, string(out))
	}
}

func TestWave26C_UpgradeExistingSyntacticScopeReleaseIdempotent(t *testing.T) {
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
	mustWrite("go.mod", "module example.com/w26cupg\n\ngo 1.24\n")
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

type OrderedMap[K comparable, V any] struct{ m map[K]V }
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V]            { return &OrderedMap[K, V]{} }
func AcquireOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func ReleaseOrderedMap[K comparable, V any](_ *OrderedMap[K, V]) {}
func (m *OrderedMap[K, V]) ReleaseToPool() {}

type Slab[T any] struct{}
func (s *Slab[T]) Alloc() *T { return new(T) }

type PooledSlab[T any] struct{}
func (s *PooledSlab[T]) Alloc() *T            { return new(T) }
func (s *PooledSlab[T]) Reset()               {}
func (s *PooledSlab[T]) Release()             {}
func (s *PooledSlab[T]) ReleaseWith(func(*T)) {}
`)
	mustWrite("gen/jscomp/syntactic_scope_creator.go", `package jscomp

import "example.com/w26cupg/gen/j2gort"

type Var struct{ name string }
type Scope struct {
	vars *j2gort.OrderedMap[string, *Var]
}
func (this_ *Scope) Spill() {
	this_.vars = &j2gort.OrderedMap[string, *Var]{}
}

type SyntacticScopeCreator struct {
	slab_Scope j2gort.PooledSlab[Scope]
	slab_Var   j2gort.PooledSlab[Var]
}
func (this_ *SyntacticScopeCreator) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_Scope.ReleaseWith(func(s_ *Scope) {
		if s_.vars != nil {
			s_.vars.ReleaseToPool()
			s_.vars = nil
		}
	})
	this_.slab_Var.Release()
}
`)

	pkgDir := filepath.Join(tmpDir, "gen", "jscomp")
	slabPass := []passes.Pass{passes.NewSlabAllocLoweringPass()}
	res1, err := passes.RunOnPackageDir(pkgDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.TotalSites != 2 {
		t.Fatalf("expected 2 sites (&OrderedMap{} -> AcquireOrderedMapPtr + ReleaseToPool -> ReleaseOrderedMap), got %+v", res1)
	}
	res2, err := passes.RunOnPackageDir(pkgDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", res2)
	}
}

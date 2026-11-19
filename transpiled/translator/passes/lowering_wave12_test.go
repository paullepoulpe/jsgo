package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

const wave12ExtraJ2gortSrc = `package j2gort

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
)

type Ext_java_util_Iterator interface {
	IsJavaNil_() bool
	HasNext() bool
	Next() any
}

func Ext_java_util_Collections_EmptyIterator() Ext_java_util_Iterator {
	return &sliceIter{}
}

func Ext_java_util_Collections_EmptyList() Ext_java_util_List {
	return &arrayListImpl{}
}

func Ext_java_util_Collections_SingletonList(v any) Ext_java_util_List {
	return &arrayListImpl{items: []any{ToAny(v)}}
}

func Ext_java_util_Collections_UnmodifiableList(l Ext_java_util_List) Ext_java_util_List {
	CheckNotNull(l)
	return l
}

func Ext_com_google_common_collect_Iterables_GetLast_1(it Ext_java_lang_Iterable) any {
	CheckNotNull(it)
	s := it.(interface{ IterSlice_() []any }).IterSlice_()
	if len(s) == 0 {
		panic("NoSuchElementException")
	}
	return s[len(s)-1]
}

type Ext_com_google_common_collect_ImmutableList interface {
	Ext_java_util_List
	IterSlice_() []any
}

func Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(c Ext_java_util_Collection) Ext_com_google_common_collect_ImmutableList {
	CheckNotNull(c)
	src := c.(interface{ IterSlice_() []any }).IterSlice_()
	cp := make([]any, len(src))
	for i, v := range src {
		cp[i] = CheckNotNull(ToAny(v))
	}
	return &arrayListImpl{items: cp}
}

type Ext_com_google_common_base_MoreObjects_ToStringHelper struct {
	parts []string
}

func Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(_ any) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
	return &Ext_com_google_common_base_MoreObjects_ToStringHelper{}
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) Add_2_string_any(name string, v any) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
	if si, ok := v.(interface{ IterSlice_() []any }); ok && !isNilVal(v) {
		var elemStrs []string
		for _, e := range si.IterSlice_() {
			elemStrs = append(elemStrs, StrVal(e))
		}
		h.parts = append(h.parts, name+"=["+strings.Join(elemStrs, ", ")+"]")
		return h
	}
	h.parts = append(h.parts, name+"="+StrVal(v))
	return h
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) String() string {
	return "{" + strings.Join(h.parts, "; ") + "}"
}

func boxJavaElem(v any) any {
	switch x := v.(type) {
	case int32:
		return ValueOfInt(x)
	case bool:
		return ValueOfBoolean(x)
	default:
		return ToAny(v)
	}
}

func unboxJavaElem[T any](v any) T {
	var zero T
	switch any(zero).(type) {
	case string:
		return any(CastString(v)).(T)
	case int32:
		return any(UnboxInt(CastPtr[BoxInt](v))).(T)
	case bool:
		return any(UnboxBoolean(CastPtr[BoxBoolean](v))).(T)
	default:
		return Cast[T](v)
	}
}

func MakeNonNilSlice[T any](cap int32) []T {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	if cap == 0 {
		return []T{}
	}
	return make([]T, 0, cap)
}

func ImmListFromSlice[T any](s []T) Ext_com_google_common_collect_ImmutableList {
	if s == nil {
		panic("NPE")
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = CheckNotNull(boxJavaElem(v))
	}
	return &arrayListImpl{items: items}
}

func SliceUnmodifiableList[T any](s []T) Ext_java_util_List {
	if s == nil {
		panic("NPE")
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = boxJavaElem(v)
	}
	return &arrayListImpl{items: items}
}

func SliceToJavaList[T any](s []T) Ext_java_util_List {
	if s == nil {
		return nil
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = boxJavaElem(v)
	}
	return &arrayListImpl{items: items}
}

func SliceIterator[T any](s []T) Ext_java_util_Iterator {
	if s == nil {
		panic("NPE")
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = boxJavaElem(v)
	}
	return &sliceIter{items: items}
}

func SliceGetFirst[T any](s []T) T {
	if s == nil {
		panic("NPE")
	}
	if len(s) == 0 {
		panic("NoSuchElementException")
	}
	return s[0]
}

func SliceGetLast[T any](s []T) T {
	if s == nil {
		panic("NPE")
	}
	if len(s) == 0 {
		panic("NoSuchElementException")
	}
	return s[len(s)-1]
}

func SliceRemoveAt[T any](s []T, idx int32) []T {
	if s == nil {
		panic("NPE")
	}
	if idx < 0 || int(idx) >= len(s) {
		panic("IndexOutOfBoundsException")
	}
	i := int(idx)
	var zero T
	copy(s[i:], s[i+1:])
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func AppendStrCol(dst []string, col any) []string {
	if dst == nil {
		panic("NPE")
	}
	CheckNotNull(col)
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, cur.NextStr())
	}
	return dst
}

func AppendNodeCol[N ~uint32](dst []N, col any) []N {
	if dst == nil {
		panic("NPE")
	}
	CheckNotNull(col)
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, N(cur.NextNode()))
	}
	return dst
}

func AppendCollection[T any](dst []T, col any) []T {
	if dst == nil {
		panic("NPE")
	}
	CheckNotNull(col)
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, unboxJavaElem[T](cur.Next()))
	}
	return dst
}

// --- Mock Set / LinkedHashSet / ImmutableSet / ArrayDeque + OrderedSet ---

type Ext_java_util_Set interface {
	Ext_java_util_Collection
	Contains(v any) bool
	Remove_1_any(v any) bool
	AddAll_1(c Ext_java_util_Collection) bool
	GetFirst() any
	IterSlice_() []any
}

type Ext_com_google_common_collect_ImmutableSet interface {
	Ext_java_util_Set
}

type Ext_java_util_LinkedHashSet struct {
	items []any
}

func (s *Ext_java_util_LinkedHashSet) IsJavaNil_() bool  { return s == nil }
func (s *Ext_java_util_LinkedHashSet) Size() int32       { return int32(len(s.items)) }
func (s *Ext_java_util_LinkedHashSet) IsEmpty() bool     { return len(s.items) == 0 }
func (s *Ext_java_util_LinkedHashSet) Clear()            { clear(s.items); s.items = s.items[:0] }
func (s *Ext_java_util_LinkedHashSet) IterSlice_() []any { return s.items }
func (s *Ext_java_util_LinkedHashSet) Iterator() any     { return &sliceIter{items: s.items} }
func (s *Ext_java_util_LinkedHashSet) GetFirst() any {
	if len(s.items) == 0 {
		panic("NoSuchElementException")
	}
	return s.items[0]
}
func javaElemEq(a, b any) bool {
	if eq, ok := a.(interface{ Equals(any) bool }); ok && !isNilVal(a) {
		return eq.Equals(b)
	}
	return a == b
}
func (s *Ext_java_util_LinkedHashSet) Contains(v any) bool {
	v = ToAny(v)
	for _, x := range s.items {
		if javaElemEq(x, v) {
			return true
		}
	}
	return false
}
func (s *Ext_java_util_LinkedHashSet) Add_1(v any) bool {
	v = ToAny(v)
	if s.Contains(v) {
		return false
	}
	s.items = append(s.items, v)
	return true
}
func (s *Ext_java_util_LinkedHashSet) Remove_1_any(v any) bool {
	v = ToAny(v)
	for i, x := range s.items {
		if javaElemEq(x, v) {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}
func (s *Ext_java_util_LinkedHashSet) AddAll_1(c Ext_java_util_Collection) bool {
	CheckNotNull(c)
	changed := false
	for cur := NewIterCursor(c); cur.HasNext(); {
		if s.Add_1(cur.Next()) {
			changed = true
		}
	}
	return changed
}

func NewExt_java_util_LinkedHashSet() *Ext_java_util_LinkedHashSet {
	return &Ext_java_util_LinkedHashSet{}
}
func NewExt_java_util_LinkedHashSet_1_int32(cap int32) *Ext_java_util_LinkedHashSet {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	return &Ext_java_util_LinkedHashSet{items: make([]any, 0, cap)}
}
func NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(c Ext_java_util_Collection) *Ext_java_util_LinkedHashSet {
	s := &Ext_java_util_LinkedHashSet{}
	s.AddAll_1(c)
	return s
}

func SetAddStr(s Ext_java_util_Set, k string) bool      { return CheckNotNull(s).Add_1(BoxString(k)) }
func SetContainsStr(s Ext_java_util_Set, k string) bool { return CheckNotNull(s).Contains(BoxString(k)) }
func SetRemoveStr(s Ext_java_util_Set, k string) bool   { return CheckNotNull(s).Remove_1_any(BoxString(k)) }
func SetAddNode[N ~uint32](s Ext_java_util_Set, n N) bool {
	return CheckNotNull(s).Add_1(BoxNode(n))
}
func SetContainsNode[N ~uint32](s Ext_java_util_Set, n N) bool {
	return CheckNotNull(s).Contains(BoxNode(n))
}
func SetRemoveNode[N ~uint32](s Ext_java_util_Set, n N) bool {
	return CheckNotNull(s).Remove_1_any(BoxNode(n))
}

func Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(c Ext_java_util_Collection) Ext_com_google_common_collect_ImmutableSet {
	CheckNotNull(c)
	out := &Ext_java_util_LinkedHashSet{}
	for cur := NewIterCursor(c); cur.HasNext(); {
		out.Add_1(CheckNotNull(cur.Next()))
	}
	return out
}

type Ext_java_util_Queue interface {
	Ext_java_util_Collection
	Remove_0() any
}

type Ext_java_util_ArrayDeque struct {
	items []any
}

func (q *Ext_java_util_ArrayDeque) IsJavaNil_() bool  { return q == nil }
func (q *Ext_java_util_ArrayDeque) Size() int32       { return int32(len(q.items)) }
func (q *Ext_java_util_ArrayDeque) IsEmpty() bool     { return len(q.items) == 0 }
func (q *Ext_java_util_ArrayDeque) Clear()            { q.items = nil }
func (q *Ext_java_util_ArrayDeque) IterSlice_() []any { return q.items }
func (q *Ext_java_util_ArrayDeque) Iterator() any     { return &sliceIter{items: q.items} }
func (q *Ext_java_util_ArrayDeque) AddAll_1(c Ext_java_util_Collection) bool {
	CheckNotNull(c)
	changed := false
	for cur := NewIterCursor(c); cur.HasNext(); {
		q.Add_1(cur.Next())
		changed = true
	}
	return changed
}
func (q *Ext_java_util_ArrayDeque) Add_1(v any) bool {
	q.items = append(q.items, CheckNotNull(ToAny(v)))
	return true
}
func (q *Ext_java_util_ArrayDeque) Remove_0() any {
	v := q.items[0]
	q.items = q.items[1:]
	return v
}

func NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection(c Ext_java_util_Collection) *Ext_java_util_ArrayDeque {
	CheckNotNull(c)
	q := &Ext_java_util_ArrayDeque{}
	for cur := NewIterCursor(c); cur.HasNext(); {
		q.Add_1(cur.Next())
	}
	return q
}

type OrderedSet[T comparable] struct {
	inline [4]T
	items  []T
	idx    map[T]int32
	live   int32
	lenInl uint8
	isInit bool
}

type OrderedSetCursor[T comparable] struct {
	s   *OrderedSet[T]
	pos int32
}

func NewOrderedSet[T comparable]() OrderedSet[T] { return OrderedSet[T]{isInit: true} }
func NewOrderedSetCap[T comparable](cap int32) OrderedSet[T] {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	s := OrderedSet[T]{isInit: true}
	if cap > 4 {
		s.items = make([]T, 0, cap)
	}
	return s
}
func OrderedSetFromSet[T comparable](src OrderedSet[T]) OrderedSet[T] {
	if !src.isInit {
		panic("NPE")
	}
	dst := OrderedSet[T]{isInit: true}
	dst.AddAllSet(src)
	return dst
}
func OrderedSetFromStrCol(col any) OrderedSet[string] {
	s := OrderedSet[string]{isInit: true}
	OrderedSetAddAllStr(&s, col)
	return s
}
func OrderedSetFromNodeCol[N ~uint32](col any) OrderedSet[N] {
	s := OrderedSet[N]{isInit: true}
	OrderedSetAddAllNode(&s, col)
	return s
}
func OrderedSetFromCol[T comparable](col any) OrderedSet[T] {
	s := OrderedSet[T]{isInit: true}
	OrderedSetAddAllCol(&s, col)
	return s
}
func (s OrderedSet[T]) IsNil() bool { return !s.isInit }
func (s OrderedSet[T]) Size() int32 {
	if !s.isInit {
		panic("NPE")
	}
	return s.live
}
func (s OrderedSet[T]) IsEmpty() bool {
	if !s.isInit {
		panic("NPE")
	}
	return s.live == 0
}
func (s OrderedSet[T]) Contains(v T) bool {
	if !s.isInit {
		panic("NPE")
	}
	if s.live == 0 {
		return false
	}
	if s.idx != nil {
		return s.idx[v] != 0
	}
	if s.items != nil {
		for _, x := range s.items {
			if x == v {
				return true
			}
		}
		return false
	}
	for i := uint8(0); i < s.lenInl; i++ {
		if s.inline[i] == v {
			return true
		}
	}
	return false
}
func (s OrderedSet[T]) GetFirst() T {
	if !s.isInit {
		panic("NPE")
	}
	if s.live == 0 {
		panic("NoSuchElementException")
	}
	if s.items == nil {
		return s.inline[0]
	}
	if s.idx != nil && int(s.live) < len(s.items) {
		for i, v := range s.items {
			if s.idx[v] == int32(i+1) {
				return v
			}
		}
	}
	return s.items[0]
}
func (s *OrderedSet[T]) Add(v T) bool {
	if !s.isInit {
		panic("NPE")
	}
	if s.items == nil {
		for i := uint8(0); i < s.lenInl; i++ {
			if s.inline[i] == v {
				return false
			}
		}
		if s.lenInl < 4 {
			s.inline[s.lenInl] = v
			s.lenInl++
			s.live++
			return true
		}
		items := make([]T, 5, 8)
		copy(items, s.inline[:4])
		items[4] = v
		var z [4]T
		s.inline = z
		s.lenInl = 0
		s.items = items
		s.live = 5
		return true
	}
	if s.idx != nil {
		if s.idx[v] != 0 {
			return false
		}
		s.items = append(s.items, v)
		s.idx[v] = int32(len(s.items))
		s.live++
		return true
	}
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	s.live++
	if len(s.items) > 8 {
		m := make(map[T]int32, 16)
		for i, x := range s.items {
			m[x] = int32(i + 1)
		}
		s.idx = m
	}
	return true
}
func (s *OrderedSet[T]) Remove(v T) bool {
	if !s.isInit {
		panic("NPE")
	}
	if s.live == 0 {
		return false
	}
	if s.items == nil {
		for i := uint8(0); i < s.lenInl; i++ {
			if s.inline[i] == v {
				var zero T
				copy(s.inline[i:], s.inline[i+1:s.lenInl])
				s.lenInl--
				s.inline[s.lenInl] = zero
				s.live--
				return true
			}
		}
		return false
	}
	if s.idx == nil {
		for i, x := range s.items {
			if x == v {
				var zero T
				copy(s.items[i:], s.items[i+1:])
				last := len(s.items) - 1
				s.items[last] = zero
				s.items = s.items[:last]
				s.live--
				return true
			}
		}
		return false
	}
	pos1 := s.idx[v]
	if pos1 == 0 {
		return false
	}
	delete(s.idx, v)
	var zero T
	s.items[pos1-1] = zero
	s.live--
	if s.live == 0 {
		s.items = s.items[:0]
		clear(s.idx)
	} else if len(s.items) > 16 && int(s.live)*2 < len(s.items) {
		w := 0
		for i, x := range s.items {
			if s.idx[x] == int32(i+1) {
				s.items[w] = x
				s.idx[x] = int32(w + 1)
				w++
			}
		}
		clear(s.items[w:])
		s.items = s.items[:w]
	}
	return true
}
func (s *OrderedSet[T]) Clear() {
	if !s.isInit {
		panic("NPE")
	}
	var z [4]T
	s.inline = z
	s.lenInl = 0
	if s.items != nil {
		clear(s.items)
		s.items = s.items[:0]
	}
	if s.idx != nil {
		clear(s.idx)
	}
	s.live = 0
}
func (s *OrderedSet[T]) AddAllSet(other OrderedSet[T]) bool {
	if !s.isInit || !other.isInit {
		panic("NPE")
	}
	if other.live == 0 {
		return false
	}
	changed := false
	if other.items == nil {
		for i := uint8(0); i < other.lenInl; i++ {
			if s.Add(other.inline[i]) {
				changed = true
			}
		}
		return changed
	}
	if other.idx != nil && int(other.live) < len(other.items) {
		for i, v := range other.items {
			if other.idx[v] == int32(i+1) && s.Add(v) {
				changed = true
			}
		}
		return changed
	}
	for _, v := range other.items {
		if s.Add(v) {
			changed = true
		}
	}
	return changed
}
func OrderedSetAddAllStr(s *OrderedSet[string], col any) bool {
	if s == nil || !s.isInit {
		panic("NPE")
	}
	CheckNotNull(col)
	changed := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(cur.NextStr()) {
			changed = true
		}
	}
	return changed
}
func OrderedSetAddAllNode[N ~uint32](s *OrderedSet[N], col any) bool {
	if s == nil || !s.isInit {
		panic("NPE")
	}
	CheckNotNull(col)
	changed := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(N(cur.NextNode())) {
			changed = true
		}
	}
	return changed
}
func OrderedSetAddAllCol[T comparable](s *OrderedSet[T], col any) bool {
	if s == nil || !s.isInit {
		panic("NPE")
	}
	CheckNotNull(col)
	changed := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(unboxJavaElem[T](cur.Next())) {
			changed = true
		}
	}
	return changed
}
func (s *OrderedSet[T]) Iter() OrderedSetCursor[T] {
	if s == nil || !s.isInit {
		panic("NPE")
	}
	return OrderedSetCursor[T]{s: s}
}
func (c *OrderedSetCursor[T]) HasNext() bool {
	s := c.s
	if s.live == 0 {
		return false
	}
	if s.items == nil {
		return c.pos < int32(s.lenInl)
	}
	if s.idx != nil && int(s.live) < len(s.items) {
		for i := int(c.pos); i < len(s.items); i++ {
			if s.idx[s.items[i]] == int32(i+1) {
				c.pos = int32(i)
				return true
			}
		}
		c.pos = int32(len(s.items))
		return false
	}
	return int(c.pos) < len(s.items)
}
func (c *OrderedSetCursor[T]) Next() T {
	s := c.s
	i := c.pos
	c.pos = i + 1
	if s.items == nil {
		return s.inline[i]
	}
	return s.items[i]
}
func (s OrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	if !s.isInit {
		panic("NPE")
	}
	out := &Ext_java_util_LinkedHashSet{}
	for cur := s.Iter(); cur.HasNext(); {
		out.Add_1(CheckNotNull(boxJavaElem(cur.Next())))
	}
	return out
}
func (s OrderedSet[T]) ToJavaSet() *Ext_java_util_LinkedHashSet {
	if !s.isInit {
		return nil
	}
	out := &Ext_java_util_LinkedHashSet{}
	for cur := s.Iter(); cur.HasNext(); {
		out.Add_1(boxJavaElem(cur.Next()))
	}
	return out
}
func ArrayDequeFromOrderedSet[T comparable](s OrderedSet[T]) *Ext_java_util_ArrayDeque {
	if !s.isInit {
		panic("NPE")
	}
	q := &Ext_java_util_ArrayDeque{}
	for cur := s.Iter(); cur.HasNext(); {
		q.Add_1(boxJavaElem(cur.Next()))
	}
	return q
}

// --- Mock Ext_java_lang_StringBuilder + LocalStringBuilder ---

type Ext_java_lang_StringBuilder struct {
	b LocalStringBuilder
}

func NewExt_java_lang_StringBuilder() *Ext_java_lang_StringBuilder {
	return &Ext_java_lang_StringBuilder{}
}
func NewExt_java_lang_StringBuilder_1_int32(cap int32) *Ext_java_lang_StringBuilder {
	sb := &Ext_java_lang_StringBuilder{}
	sb.b.InitCap(cap)
	return sb
}
func NewExt_java_lang_StringBuilder_1_string(s string) *Ext_java_lang_StringBuilder {
	sb := &Ext_java_lang_StringBuilder{}
	sb.b.InitStr(s)
	return sb
}
func (s *Ext_java_lang_StringBuilder) IsJavaNil_() bool { return s == nil }
func (s *Ext_java_lang_StringBuilder) Append_1_string(v string) any {
	s.b.AppendStr(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_uint16(c uint16) any {
	s.b.AppendChar(c)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_int32(v int32) any {
	s.b.AppendInt(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_int64(v int64) any {
	s.b.AppendLong(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_bool(v bool) any {
	s.b.AppendBool(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_float64(v float64) any {
	s.b.AppendDouble(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) Append_1_any(v any) any {
	s.b.AppendAny(v)
	return s
}
func (s *Ext_java_lang_StringBuilder) SetLength(n int32) {
	if n == 0 {
		s.b.Reset()
		return
	}
	str := s.b.String()
	s.b.Reset()
	s.b.AppendStr(str[:n])
}
func (s *Ext_java_lang_StringBuilder) CharAt(idx int32) uint16 {
	return uint16(s.b.String()[idx])
}
func (s *Ext_java_lang_StringBuilder) IsEmpty() bool  { return s.b.IsEmpty() }
func (s *Ext_java_lang_StringBuilder) String() string { return s.b.String() }

type LocalStringBuilder struct {
	buf         strings.Builder
	pendingHigh uint16
}

func (b *LocalStringBuilder) InitCap(cap int32) *LocalStringBuilder {
	if cap < 0 {
		panic("NegativeArraySizeException")
	}
	if cap > 0 && cap <= (1<<24) {
		b.buf.Grow(int(cap))
	}
	return b
}
func (b *LocalStringBuilder) InitStr(s string) *LocalStringBuilder {
	CheckNotNullString(s)
	return b.AppendStr(s)
}
func (b *LocalStringBuilder) flushPendingHigh() {
	if b.pendingHigh != 0 {
		ch := b.pendingHigh
		b.pendingHigh = 0
		b.buf.WriteByte(byte(0xE0 | (ch >> 12)))
		b.buf.WriteByte(byte(0x80 | ((ch >> 6) & 0x3F)))
		b.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
	}
}
func (b *LocalStringBuilder) AppendStr(raw string) *LocalStringBuilder {
	s := raw
	if s == NullString {
		s = "null"
	}
	if len(s) == 0 {
		return b
	}
	if b.pendingHigh == 0 && s[len(s)-1] < 0x80 {
		b.buf.WriteString(s)
		return b
	}
	if b.pendingHigh != 0 {
		if len(s) >= 3 && s[0] == 0xED && s[1] >= 0xB0 && s[1] <= 0xBF && (s[2]&0xC0) == 0x80 {
			low := (uint16(s[0]&0x0F) << 12) | (uint16(s[1]&0x3F) << 6) | uint16(s[2]&0x3F)
			r := utf16.DecodeRune(rune(b.pendingHigh), rune(low))
			b.pendingHigh = 0
			b.buf.WriteRune(r)
			s = s[3:]
			if len(s) == 0 {
				return b
			}
		} else {
			b.flushPendingHigh()
		}
	}
	n := len(s)
	if n >= 3 && s[n-3] == 0xED && s[n-2] >= 0xA0 && s[n-2] <= 0xAF && (s[n-1]&0xC0) == 0x80 {
		high := (uint16(s[n-3]&0x0F) << 12) | (uint16(s[n-2]&0x3F) << 6) | uint16(s[n-1]&0x3F)
		b.buf.WriteString(s[:n-3])
		b.pendingHigh = high
		return b
	}
	b.buf.WriteString(s)
	return b
}
func (b *LocalStringBuilder) AppendChar(ch uint16) *LocalStringBuilder {
	if b.pendingHigh == 0 && ch < 0x80 {
		b.buf.WriteByte(byte(ch))
		return b
	}
	if b.pendingHigh != 0 {
		if ch >= 0xDC00 && ch <= 0xDFFF {
			r := utf16.DecodeRune(rune(b.pendingHigh), rune(ch))
			b.pendingHigh = 0
			b.buf.WriteRune(r)
			return b
		}
		b.flushPendingHigh()
	}
	if ch >= 0xD800 && ch <= 0xDBFF {
		b.pendingHigh = ch
		return b
	}
	if ch < 0x80 {
		b.buf.WriteByte(byte(ch))
	} else if ch < 0x800 {
		b.buf.WriteByte(byte(0xC0 | (ch >> 6)))
		b.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
	} else {
		b.buf.WriteByte(byte(0xE0 | (ch >> 12)))
		b.buf.WriteByte(byte(0x80 | ((ch >> 6) & 0x3F)))
		b.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
	}
	return b
}
func (b *LocalStringBuilder) AppendInt(v int32) *LocalStringBuilder {
	return b.AppendLong(int64(v))
}
func (b *LocalStringBuilder) AppendLong(v int64) *LocalStringBuilder {
	b.flushPendingHigh()
	var tmp [20]byte
	b.buf.Write(strconv.AppendInt(tmp[:0], v, 10))
	return b
}
func (b *LocalStringBuilder) AppendBool(v bool) *LocalStringBuilder {
	b.flushPendingHigh()
	if v {
		b.buf.WriteString("true")
	} else {
		b.buf.WriteString("false")
	}
	return b
}
func (b *LocalStringBuilder) AppendDouble(v float64) *LocalStringBuilder {
	b.flushPendingHigh()
	b.buf.WriteString(fmt.Sprint(v))
	return b
}
func (b *LocalStringBuilder) AppendFloat(v float32) *LocalStringBuilder {
	b.flushPendingHigh()
	b.buf.WriteString(fmt.Sprint(v))
	return b
}
func (b *LocalStringBuilder) AppendAny(v any) *LocalStringBuilder {
	return b.AppendStr(StrVal(v))
}
func (b *LocalStringBuilder) IsEmpty() bool { return b.buf.Len() == 0 && b.pendingHigh == 0 }
func (b *LocalStringBuilder) Reset()        { b.buf.Reset(); b.pendingHigh = 0 }
func (b *LocalStringBuilder) String() string {
	if b.pendingHigh == 0 {
		return b.buf.String()
	}
	s := b.buf.String()
	ch := b.pendingHigh
	out := make([]byte, len(s)+3)
	copy(out, s)
	out[len(s)] = byte(0xE0 | (ch >> 12))
	out[len(s)+1] = byte(0x80 | ((ch >> 6) & 0x3F))
	out[len(s)+2] = byte(0x80 | (ch & 0x3F))
	return string(out)
}
`

func writeWave12Workspace(t *testing.T, pkgFiles map[string]string) (string, string) {
	t.Helper()
	tmpDir, appDir := writeTempJ2gortWorkspace(t, pkgFiles)
	if err := os.WriteFile(filepath.Join(tmpDir, "j2gort", "wave12_ext.go"), []byte(wave12ExtraJ2gortSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpDir, appDir
}

func TestWave12SliceLowering_LazyReferenceCollectionAndBuilderHandoffs(t *testing.T) {
	refCollSrc := `package app

import "example.com/j2gotest/j2gort"

type Reference struct{ Name string }

func (*Reference) HasIdentityEquals_() {}
func (r *Reference) IsJavaNil_() bool  { return r == nil }
func (r *Reference) String() string    { return r.Name }

type ReferenceCollection struct {
	ref0 *Reference
	rest j2gort.Ext_java_util_ArrayList
}

func (this_ *ReferenceCollection) initFields_ReferenceCollection() {
	_ = this_
}

func NewReferenceCollection() *ReferenceCollection {
	c := &ReferenceCollection{}
	c.initFields_ReferenceCollection()
	return c
}

func (this_ *ReferenceCollection) IsEmpty() bool {
	return (this_.ref0 == nil) && (this_.rest == nil)
}

func (this_ *ReferenceCollection) Size() int32 {
	if this_.rest != nil {
		return this_.rest.Size()
	}
	if this_.ref0 != nil {
		return 1
	}
	return 0
}

func (this_ *ReferenceCollection) Get(index int32) *Reference {
	if this_.rest != nil {
		return j2gort.CastPtr[Reference](this_.rest.Get(index))
	}
	if index == 0 && this_.ref0 != nil {
		return this_.ref0
	}
	panic("IndexOutOfBoundsException")
}

func (this_ *ReferenceCollection) Add(reference *Reference) {
	if this_.rest != nil {
		this_.rest.Add_1(j2gort.PtrToAny(reference))
	} else if this_.ref0 == nil {
		this_.ref0 = reference
	} else {
		var list j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
		_ = list
		list.Add_1(j2gort.PtrToAny(this_.ref0))
		list.Add_1(j2gort.PtrToAny(reference))
		this_.rest = list
		this_.ref0 = nil
	}
}

func (this_ *ReferenceCollection) Iterator() any {
	if this_.rest != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.rest.Iterator())
	}
	if this_.ref0 == nil {
		return j2gort.Ext_java_util_Collections_EmptyIterator()
	}
	return j2gort.Ext_java_util_Collections_SingletonList(j2gort.PtrToAny(this_.ref0)).Iterator()
}

func (this_ *ReferenceCollection) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("references", func() any {
		if this_.rest != nil {
			return this_.rest
		}
		if this_.ref0 != nil {
			return j2gort.Ext_java_util_Collections_SingletonList(j2gort.PtrToAny(this_.ref0))
		}
		return j2gort.Ext_java_util_Collections_EmptyList()
	}()).String()
}

// UnionType_Builder + FunctionParamBuilder + RemovableBuilder patterns:
type JSType interface {
	IsJavaNil_() bool
	TypeName() string
}
type FunctionType interface {
	JSType
	Arity() int32
}

type ParamBuilder struct {
	parameters j2gort.Ext_java_util_ArrayList
}

func NewParamBuilder(cap int32) *ParamBuilder {
	p := &ParamBuilder{}
	p.parameters = j2gort.NewExt_java_util_ArrayList_1_int32(cap)
	return p
}

func (this_ *ParamBuilder) AddType(t JSType, sub FunctionType) j2gort.Ext_com_google_common_collect_ImmutableList {
	this_.parameters.Add_1(t)
	if sub != nil && !sub.IsJavaNil_() {
		this_.parameters.Set(int32(0), sub)
		this_.parameters.Add_1(t)
		this_.parameters.Remove_1_int32(int32(1))
	}
	_ = j2gort.AsIface[JSType](j2gort.Ext_com_google_common_collect_Iterables_GetLast_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.parameters)))
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.parameters))
}

type Cont struct{ ID int32 }
type RemovableBuilder struct {
	continuations j2gort.Ext_java_util_List
}
type Removable struct {
	continuations j2gort.Ext_java_util_List
}

func NewRemovableBuilder() *RemovableBuilder {
	b := &RemovableBuilder{}
	b.continuations = j2gort.NewExt_java_util_ArrayList()
	return b
}
func (b *RemovableBuilder) AddContinuation(c *Cont) *RemovableBuilder {
	b.continuations.Add_1(j2gort.PtrToAny(c))
	return b
}
func (b *RemovableBuilder) Build() *Removable {
	r := &Removable{}
	InitRemovable(r, b)
	return r
}
func InitRemovable(this_ *Removable, builder *RemovableBuilder) {
	this_.continuations = builder.continuations
}
func (this_ *Removable) Drain() int32 {
	var sum int32
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.continuations))
	for cur.HasNext() {
		raw := cur.Next()
		c := j2gort.CastPtr[Cont](raw)
		sum += c.ID
	}
	this_.continuations.Clear()
	return sum
}
`
	_, appDir := writeWave12Workspace(t, map[string]string{"refcoll.go": refCollSrc})

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir failed: %v", err)
	}
	if stats1.ByPass["native_slice_lowering"] < 20 {
		t.Fatalf("expected >= 20 native_slice_lowering sites, got %+v", stats1.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(appDir, "refcoll.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"rest []*Reference",
		"var list []*Reference = make([]*Reference, 0, int32(4))",
		"this_.rest = list",
		"return j2gort.SliceIterator(this_.rest)",
		"return j2gort.SliceToJavaList(this_.rest)",
		"parameters []JSType",
		"p.parameters = j2gort.MakeNonNilSlice[JSType](cap)",
		"this_.parameters[int32(0)] = sub",
		"this_.parameters = j2gort.SliceRemoveAt(this_.parameters, int32(1))",
		"j2gort.SliceGetLast(this_.parameters)",
		"return j2gort.ImmListFromSlice(this_.parameters)",
		"continuations []*Cont",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected snippet %q in refcoll.go:\n%s", want, got)
		}
	}

	stats2, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	// Differential runtime parity check between unlowered orig and lowered app!
	origDir := filepath.Join(filepath.Dir(appDir), "orig")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origDir, "refcoll.go"), []byte(strings.Replace(refCollSrc, "package app", "package orig", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	sliceDiff := `package refcoll_diff_test

import (
	"fmt"
	"testing"

	"example.com/j2gotest/app"
	"example.com/j2gotest/j2gort"
	"example.com/j2gotest/orig"
)

type origFnType struct{ name string }
func (f *origFnType) IsJavaNil_() bool { return f == nil }
func (f *origFnType) TypeName() string { return f.name }
func (f *origFnType) Arity() int32     { return 1 }

type appFnType struct{ name string }
func (f *appFnType) IsJavaNil_() bool { return f == nil }
func (f *appFnType) TypeName() string { return f.name }
func (f *appFnType) Arity() int32     { return 1 }

func TestRefCollAndBuildersParity(t *testing.T) {
	for _, n := range []int{0, 1, 2, 5} {
		oc := orig.NewReferenceCollection()
		ac := app.NewReferenceCollection()
		for i := 0; i < n; i++ {
			lbl := fmt.Sprintf("r%d", i)
			oc.Add(&orig.Reference{Name: lbl})
			ac.Add(&app.Reference{Name: lbl})
		}
		if oc.IsEmpty() != ac.IsEmpty() || oc.Size() != ac.Size() || oc.String() != ac.String() {
			t.Fatalf("n=%d mismatch: orig=(%v,%d,%q) vs app=(%v,%d,%q)", n, oc.IsEmpty(), oc.Size(), oc.String(), ac.IsEmpty(), ac.Size(), ac.String())
		}
		oit := oc.Iterator().(j2gort.Ext_java_util_Iterator)
		ait := ac.Iterator().(j2gort.Ext_java_util_Iterator)
		for oit.HasNext() || ait.HasNext() {
			if oit.HasNext() != ait.HasNext() || j2gort.StrVal(oit.Next()) != j2gort.StrVal(ait.Next()) {
				t.Fatalf("n=%d iterator mismatch", n)
			}
		}
	}

	opb := orig.NewParamBuilder(0)
	apb := app.NewParamBuilder(0)
	ol := opb.AddType(&origFnType{"t0"}, &origFnType{"fn1"})
	al := apb.AddType(&appFnType{"t0"}, &appFnType{"fn1"})
	if ol.Size() != al.Size() || ol.Get(0).(orig.JSType).TypeName() != al.Get(0).(app.JSType).TypeName() {
		t.Fatalf("ParamBuilder mismatch: %d vs %d", ol.Size(), al.Size())
	}

	orb := orig.NewRemovableBuilder().AddContinuation(&orig.Cont{ID: 10}).AddContinuation(&orig.Cont{ID: 25}).Build()
	arb := app.NewRemovableBuilder().AddContinuation(&app.Cont{ID: 10}).AddContinuation(&app.Cont{ID: 25}).Build()
	if orb.Drain() != arb.Drain() || orb.Drain() != arb.Drain() {
		t.Fatalf("RemovableBuilder Drain mismatch")
	}
}
`
	if err := os.WriteFile(filepath.Join(filepath.Dir(appDir), "refcoll_diff_test.go"), []byte(sliceDiff), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = filepath.Dir(appDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("refcoll differential runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave12SetAndStringBuilderLowering_DifferentialParityAndIdempotence(t *testing.T) {
	workloadSrc := `package orig

import "example.com/j2gotest/j2gort"

type Node uint32

type Tree struct {
	Name     string
	LValue   bool
	Children []*Tree
}

// Mirrors FunctionArgumentInjector_FindModifiedParameters recursive Set return across call stack!
func FindModifiedParameters(t *Tree, params j2gort.Ext_com_google_common_collect_ImmutableSet) j2gort.Ext_com_google_common_collect_ImmutableSet {
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(
		j2gort.AsIface[j2gort.Ext_java_util_Collection](FunctionArgumentInjector_FindModifiedParameters(t, params)),
	)
}

func FunctionArgumentInjector_FindModifiedParameters(t *Tree, names j2gort.Ext_com_google_common_collect_ImmutableSet) j2gort.Ext_java_util_Set {
	var unsafe *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = unsafe
	if t == nil {
		return unsafe
	}
	if t.Name != "" && j2gort.SetContainsStr(names, t.Name) && t.LValue {
		j2gort.SetAddStr(unsafe, t.Name)
	}
	for _, ch := range t.Children {
		unsafe.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](FunctionArgumentInjector_FindModifiedParameters(ch, names)))
	}
	return unsafe
}

// Mirrors DeadPropertyAssignmentElimination_Property.children + ArrayDeque BFS!
type Prop struct {
	Label    string
	children *j2gort.Ext_java_util_LinkedHashSet
}

func (*Prop) HasIdentityEquals_() {}

func NewProp(label string) *Prop {
	p := &Prop{Label: label}
	p.initFields_Prop()
	return p
}

func (this_ *Prop) initFields_Prop() {
	this_.children = j2gort.NewExt_java_util_LinkedHashSet()
}

func (this_ *Prop) AddChild(c *Prop) {
	this_.children.Add_1(j2gort.PtrToAny(c))
}

func (this_ *Prop) BFSLabels() string {
	var propertiesSet j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.children))
	var q j2gort.Ext_java_util_Queue = j2gort.NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection(propertiesSet)
	propertiesSet.Add_1(j2gort.PtrToAny(this_))

	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](
		j2gort.NewExt_java_lang_StringBuilder_1_string("bfs:").Append_1_string(this_.Label),
	)
	_ = sb
	for !(q.IsEmpty()) {
		child := j2gort.CastPtr[Prop](q.Remove_0())
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(','))).Append_1_string(child.Label)
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(child.children))
		for cur.HasNext() {
			raw := cur.Next()
			gc := j2gort.CastPtr[Prop](raw)
			if propertiesSet.Add_1(j2gort.PtrToAny(gc)) {
				q.Add_1(j2gort.PtrToAny(gc))
			}
		}
	}
	return sb.String()
}

// Exercises OrderedSet spill (> 8 elements), Remove + re-Add insertion order, NullString, and StringBuilder surrogates!
func ExerciseSetSpillAndStringBuilder(items []string, maybeNull string, capHint int32) string {
	var s j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(2))
	for _, it := range items {
		j2gort.SetAddStr(s, it)
	}
	if len(items) > 0 {
		j2gort.SetRemoveStr(s, items[0])
		j2gort.SetAddStr(s, items[0])
	}
	j2gort.SetAddStr(s, maybeNull)

	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_int32(capHint)
	_ = sb
	sb.Append_1_string("junk")
	sb.SetLength(int32(0))
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_bool(s.IsEmpty())).Append_1_uint16(uint16('|')),
	).Append_1_int32(s.Size())
	// Surrogate pair D83D + DE00 -> U+1F600, plus lone high surrogate D800 before int!
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(0xD83D))).Append_1_uint16(uint16(0xDE00))
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(0xD800))).Append_1_int32(int32(7))
	sb.Append_1_string(maybeNull)
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(s))
		for cur.HasNext() {
			k := cur.NextStr()
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(':'))).Append_1_any(j2gort.BoxString(k))
		}
	}
	// Trailing lone high surrogate read via String(), then paired!
	var tail *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	tail.Append_1_uint16(uint16(0xD83D))
	snap1 := tail.String()
	tail.Append_1_uint16(uint16(0xDE00))
	return j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(snap1)).Append_1_string(tail.String()),
	).String()
}

// Negative guards: custom-Equals interface Set & indexed StringBuilder MUST NOT be lowered!
type CustomKey interface {
	IsJavaNil_() bool
	Equals(any) bool
}

func KeepUnsafePatterns(k1, k2 CustomKey) (int32, uint16) {
	var customSet j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	customSet.Add_1(k1)
	customSet.Add_1(k2)

	var indexedSB *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string("hello")
	ch := indexedSB.CharAt(int32(1))
	return customSet.Size(), ch
}
`

	tmpDir, _ := writeWave12Workspace(t, nil)
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

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(opt) failed: %v", err)
	}
	if stats1.ByPass["native_set_lowering"] < 15 {
		t.Fatalf("expected >= 15 native_set_lowering sites, got %+v", stats1.ByPass)
	}
	if stats1.ByPass["string_builder_lowering"] < 10 {
		t.Fatalf("expected >= 10 string_builder_lowering sites, got %+v", stats1.ByPass)
	}

	// Verify AST shape of lowered opt/workload.go.
	optBytes, err := os.ReadFile(filepath.Join(optDir, "workload.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotOpt := string(optBytes)
	for _, want := range []string{
		"func ordSet_FunctionArgumentInjector_FindModifiedParameters(t *Tree, names j2gort.Ext_com_google_common_collect_ImmutableSet) j2gort.OrderedSet[string]",
		"return ordSet_FunctionArgumentInjector_FindModifiedParameters(t, params).ToImmutableSet()",
		"return ordSet_FunctionArgumentInjector_FindModifiedParameters(t, names).ToJavaSet()",
		"unsafe.AddAllSet(ordSet_FunctionArgumentInjector_FindModifiedParameters(ch, names))",
		"children j2gort.OrderedSet[*Prop]",
		"j2gort.OrderedSetFromSet(this_.children)",
		"j2gort.ArrayDequeFromOrderedSet(propertiesSet)",
		"var sb j2gort.LocalStringBuilder",
		`sb.InitStr("bfs:").AppendStr(this_.Label)`,
		"sb.Reset()",
		// Negative guards preserved:
		"var customSet j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()",
		`var indexedSB *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string("hello")`,
	} {
		if !strings.Contains(gotOpt, want) {
			t.Errorf("missing expected snippet %q in opt/workload.go:\n%s", want, gotOpt)
		}
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2, got %+v", stats2)
	}

	// Differential runtime test (including >8 OrderedSet spill, removals, NullString, negative cap panics, surrogates).
	diffTest := `package wave12_diff_test

import (
	"fmt"
	"testing"

	"example.com/j2gotest/j2gort"
	"example.com/j2gotest/opt"
	"example.com/j2gotest/orig"
)

func runSafe(fn func() string) (out string, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	return fn(), false
}

func TestWave12Parity(t *testing.T) {
	// 1. Recursive FindModifiedParameters parity
	origNames := j2gort.NewExt_java_util_LinkedHashSet()
	origNames.Add_1("a")
	origNames.Add_1("b")
	origNames.Add_1("c")
	paramSet := j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(origNames)

	origTree := &orig.Tree{
		Children: []*orig.Tree{
			{Name: "b", LValue: true},
			{Name: "x", LValue: true},
			{Name: "a", LValue: false, Children: []*orig.Tree{{Name: "a", LValue: true}, {Name: "b", LValue: true}}},
		},
	}
	optTree := &opt.Tree{
		Children: []*opt.Tree{
			{Name: "b", LValue: true},
			{Name: "x", LValue: true},
			{Name: "a", LValue: false, Children: []*opt.Tree{{Name: "a", LValue: true}, {Name: "b", LValue: true}}},
		},
	}
	resOrig := orig.FindModifiedParameters(origTree, paramSet).IterSlice_()
	resOpt := opt.FindModifiedParameters(optTree, paramSet).IterSlice_()
	if fmt.Sprint(resOrig) != fmt.Sprint(resOpt) || fmt.Sprint(resOpt) != "[b a]" {
		t.Fatalf("FindModifiedParameters mismatch: orig=%v vs opt=%v", resOrig, resOpt)
	}
	// Also test exported wrapper parity!
	wrapOrig := orig.FunctionArgumentInjector_FindModifiedParameters(origTree, paramSet).IterSlice_()
	wrapOpt := opt.FunctionArgumentInjector_FindModifiedParameters(optTree, paramSet).IterSlice_()
	if fmt.Sprint(wrapOrig) != fmt.Sprint(wrapOpt) {
		t.Fatalf("Exported wrapper mismatch: orig=%v vs opt=%v", wrapOrig, wrapOpt)
	}

	// 2. Prop graph BFS parity (with cycle!)
	oa, ob, oc := orig.NewProp("A"), orig.NewProp("B"), orig.NewProp("C")
	oa.AddChild(ob)
	ob.AddChild(oc)
	oc.AddChild(oa)

	xa, xb, xc := opt.NewProp("A"), opt.NewProp("B"), opt.NewProp("C")
	xa.AddChild(xb)
	xb.AddChild(xc)
	xc.AddChild(xa)

	if oa.BFSLabels() != xa.BFSLabels() {
		t.Fatalf("BFSLabels mismatch: orig=%q vs opt=%q", oa.BFSLabels(), xa.BFSLabels())
	}

	// 3. OrderedSet inline/slice/map spill + compaction + LocalStringBuilder surrogates & NullString & negative cap panic
	var twenty []string
	for i := 0; i < 20; i++ {
		twenty = append(twenty, fmt.Sprintf("k%d", i))
	}
	cases := []struct {
		items []string
		s     string
		cap   int32
	}{
		{nil, "ok", 0},
		{[]string{"a", "b", "c"}, j2gort.NullString, 16},
		{[]string{"1", "2", "3", "4", "5", "6"}, "", 4},
		{twenty, j2gort.NullString, 8},
		{[]string{"neg"}, "x", -1},
	}
	for i, tc := range cases {
		goOrig, panOrig := runSafe(func() string { return orig.ExerciseSetSpillAndStringBuilder(tc.items, tc.s, tc.cap) })
		goOpt, panOpt := runSafe(func() string { return opt.ExerciseSetSpillAndStringBuilder(tc.items, tc.s, tc.cap) })
		if panOrig != panOpt || goOrig != goOpt {
			t.Fatalf("case %d mismatch: orig=(%q, panic=%v) vs opt=(%q, panic=%v)", i, goOrig, panOrig, goOpt, panOpt)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "wave12_diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime parity test failed: %v\n%s", err, string(out))
	}
}

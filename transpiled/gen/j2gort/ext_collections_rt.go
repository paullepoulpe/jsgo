// Auto-generated clean-room deterministic collections, Guava, and Stream runtime for j2gort.
package j2gort

import (
	"fmt"
	"math/bits"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

var _ = javautil.NewBitSet
var _ = nodestore.Node(0)
var _ = bits.TrailingZeros64
var _ = reflect.TypeFor[int]
var _ = (*atomic.Bool)(nil)
var _ = strconv.AppendInt
var _ = utf16.DecodeRune
var _ = slices.Sort[[]string]

var (
	compPoolResetMu    sync.Mutex
	compPoolResets     []func()
	runtimePoolPurgers []func()
)

func registerCompPoolReset(fn func()) {
	compPoolResetMu.Lock()
	compPoolResets = append(compPoolResets, fn)
	compPoolResetMu.Unlock()
}

func registerRuntimePoolPurger(fn func()) {
	compPoolResetMu.Lock()
	runtimePoolPurgers = append(runtimePoolPurgers, fn)
	compPoolResetMu.Unlock()
}

func PurgeAllRuntimePools() {
	TrimStaleSlabPools(0)
	compPoolResetMu.Lock()
	fns := runtimePoolPurgers
	compPoolResetMu.Unlock()
	for _, fn := range fns {
		fn()
	}
}

func FlushCompilationPools() {
	compPoolResetMu.Lock()
	fns := compPoolResets
	compPoolResetMu.Unlock()
	for _, fn := range fns {
		fn()
	}
	slabGlobalPooledBytes.Store(0)
}

func ResetCompilationPools() {
	FlushCompilationPools()
}

func init() {
	nodestore.RegisterStoreResetHook(FlushCompilationPools)
	RegisterTypePredicate("com.google.common.collect.Multimap", "com.google.common.collect.Multimap", func(v any) bool {
		_, ok := v.(Ext_com_google_common_collect_Multimap)
		return ok
	})
	RegisterTypePredicate("java.util.Collection", "java.util.Collection", func(v any) bool {
		_, ok := v.(Ext_java_util_Collection)
		return ok
	})
	RegisterTypePredicate("java.util.List", "java.util.List", func(v any) bool {
		_, ok := v.(Ext_java_util_List)
		return ok
	})
	RegisterTypePredicate("java.util.Set", "java.util.Set", func(v any) bool {
		_, ok := v.(Ext_java_util_Set)
		return ok
	})
	RegisterTypePredicate("java.util.Map", "java.util.Map", func(v any) bool {
		_, ok := v.(Ext_java_util_Map)
		return ok
	})
}

// --- External Interface Default Methods ---

func ExtDefault_Iterable_ForEach(this_ Ext_java_lang_Iterable, action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	for it := NewIterCursor(this_); it.HasNext(); {
		action.Accept(it.Next())
	}
}

func ExtDefault_Iterator_ForEachRemaining(this_ Ext_java_util_Iterator, action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	for this_.HasNext() {
		action.Accept(this_.Next())
	}
}

type jNatOrderCore struct{}

type jComparatorFunc struct {
	fn  func(a, b any) int32
	nat bool
}

var jNaturalOrderCmp Ext_java_util_Comparator = &jComparatorFunc{fn: CompareValues, nat: true}
var jReverseOrderCmp Ext_java_util_Comparator = &jComparatorFunc{fn: func(a, b any) int32 { return CompareValues(b, a) }}

func newComparatorFromFunc(fn func(a, b any) int32) Ext_java_util_Comparator {
	return &jComparatorFunc{fn: fn}
}

func (c *jComparatorFunc) isNatOrder_() bool                         { return c != nil && c.nat }
func (*jComparatorFunc) Implements_j2gort_Ext_java_util_Comparator() {}
func (c *jComparatorFunc) Compare(a, b any) int32                    { return c.fn(a, b) }
func (c *jComparatorFunc) Reversed() Ext_java_util_Comparator {
	return ExtDefault_Comparator_Reversed(c)
}
func (c *jComparatorFunc) ThenComparing_1_Ext_java_util_Comparator(other Ext_java_util_Comparator) Ext_java_util_Comparator {
	return ExtDefault_Comparator_ThenComparing(c, other)
}
func (c *jComparatorFunc) ThenComparing_1_Ext_java_util_function_Function(ke Ext_java_util_function_Function) Ext_java_util_Comparator {
	return ExtDefault_Comparator_ThenComparingKey(c, ke)
}
func (c *jComparatorFunc) ThenComparingInt(ke Ext_java_util_function_ToIntFunction) Ext_java_util_Comparator {
	return ExtDefault_Comparator_ThenComparingInt(c, ke)
}
func (c *jComparatorFunc) ThenComparingLong(ke Ext_java_util_function_ToLongFunction) Ext_java_util_Comparator {
	return ExtDefault_Comparator_ThenComparingLong(c, ke)
}

func arraysSortStrSlice(sub []string) {
	ascii := true
	for _, s := range sub {
		if !ascii {
			if s == NullString {
				ThrowNPE()
			}
			continue
		}
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				if s == NullString {
					ThrowNPE()
				}
				ascii = false
				break
			}
		}
	}
	n := len(sub)
	if n < 2 {
		return
	}
	if ascii {
		if n <= 16 {
			for i := 1; i < n; i++ {
				v := sub[i]
				j := i - 1
				for j >= 0 && sub[j] > v {
					sub[j+1] = sub[j]
					j--
				}
				sub[j+1] = v
			}
			return
		}
		slices.Sort(sub)
		return
	}
	slices.SortFunc(sub, javalang.StringCompareTo)
}

func ArraysSortStrings(a []string) {
	if a == nil {
		ThrowNPE()
	}
	arraysSortStrSlice(a)
}

func ArraysSortStringsRange(a []string, fromIndex, toIndex int32) {
	if a == nil {
		ThrowNPE()
	}
	if fromIndex > toIndex {
		panic(javalang.NewIllegalArgumentException("fromIndex > toIndex"))
	}
	if fromIndex < 0 || int(toIndex) > len(a) {
		panic(javalang.NewArrayIndexOutOfBoundsException(""))
	}
	arraysSortStrSlice(a[fromIndex:toIndex])
}

func ComparatorComparingStr(keyFn func(any) string) Ext_java_util_Comparator {
	if keyFn == nil {
		ThrowNPE()
	}
	return newComparatorFromFunc(func(a, b any) int32 {
		ka := keyFn(a)
		kb := keyFn(b)
		if ka == NullString || kb == NullString {
			ThrowNPE()
		}
		return int32(javalang.StringCompareTo(ka, kb))
	})
}

func ComparatorThenComparingStr(base Ext_java_util_Comparator, keyFn func(any) string) Ext_java_util_Comparator {
	CheckNotNull(base)
	if keyFn == nil {
		ThrowNPE()
	}
	if cf, ok := base.(*jComparatorFunc); ok {
		baseFn := cf.fn
		return newComparatorFromFunc(func(a, b any) int32 {
			c := baseFn(a, b)
			if c != 0 {
				return c
			}
			ka := keyFn(a)
			kb := keyFn(b)
			if ka == NullString || kb == NullString {
				ThrowNPE()
			}
			return int32(javalang.StringCompareTo(ka, kb))
		})
	}
	return newComparatorFromFunc(func(a, b any) int32 {
		c := base.Compare(a, b)
		if c != 0 {
			return c
		}
		ka := keyFn(a)
		kb := keyFn(b)
		if ka == NullString || kb == NullString {
			ThrowNPE()
		}
		return int32(javalang.StringCompareTo(ka, kb))
	})
}

func isNaturalOrderComparator(cmp Ext_java_util_Comparator) bool {
	if isJavaNilAny(cmp) {
		return true
	}
	if m, ok := any(cmp).(interface{ isNatOrder_() bool }); ok && m.isNatOrder_() {
		return true
	}
	return false
}

func initSortedCore(mc *JMapCore, cmp func(a, b any) int32) *JMapCore {
	if cmp == nil {
		mc.mode = jMapModeSortedNat
		mc.ext = nil
	} else {
		mc.mode = jMapModeSortedCmp
		mc.ext = &jSortedMapState{cmp: cmp}
	}
	return mc
}

func initSortedCoreCmp(mc *JMapCore, cmp Ext_java_util_Comparator) *JMapCore {
	if isNaturalOrderComparator(cmp) {
		return initSortedCore(mc, nil)
	}
	return initSortedCore(mc, comparatorFunc(cmp))
}

func initSetCapLoad(sc *JSetCore, initialCap int32, loadFactor float32) {
	if initialCap < 0 {
		panic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(initialCap)))
	}
	if loadFactor <= 0 || loadFactor != loadFactor {
		panic(javalang.NewIllegalArgumentException("Illegal load factor: " + fmt.Sprint(loadFactor)))
	}
	sc.InitExpectedSize(initialCap)
}

func newSetFromMap(m Ext_java_util_Map) Ext_java_util_Set {
	CheckNotNull(m)
	if m.Size() != 0 {
		panic(javalang.NewIllegalArgumentException("Map is non-empty"))
	}
	obj := &Ext_java_util_HashSetBase{}
	obj.BindExtSelf_(obj)
	obj.Core_ = m
	return obj
}

func newSortedCoreFromCmp(cmp Ext_java_util_Comparator) *JMapCore {
	if isNaturalOrderComparator(cmp) {
		return NewSortedJMapCore(nil)
	}
	return NewSortedJMapCore(comparatorFunc(cmp))
}

func sortedCoreNav(mc *JMapCore, k any, op int) (any, any, bool) {
	k = ToAny(k)
	if mc == nil || mc.Live == 0 || mc.ext == nil {
		if mc != nil && mc.mode == jMapModeSortedNat && k == nil {
			ThrowNPE()
		}
		return nil, nil, false
	}
	st := mc.ext.(*jSortedMapState)
	if st.natOrder && k == nil {
		ThrowNPE()
	}
	best := int32(0)
	cur := st.root
	for cur != 0 {
		c := st.avlCmpAny(k, cur)
		if c == 0 {
			if op == -1 || op == 1 {
				return st.treeNodeKey(cur), st.tree[cur-1].val, true
			}
			if op < 0 {
				cur = st.tree[cur-1].left
			} else {
				cur = st.tree[cur-1].right
			}
		} else if c < 0 {
			if op > 0 {
				best = cur
			}
			cur = st.tree[cur-1].left
		} else {
			if op < 0 {
				best = cur
			}
			cur = st.tree[cur-1].right
		}
	}
	if best == 0 {
		return nil, nil, false
	}
	return st.treeNodeKey(best), st.tree[best-1].val, true
}

func sortedCoreNavKey(mc *JMapCore, k any, op int) any {
	rk, _, ok := sortedCoreNav(mc, k, op)
	if !ok {
		return nil
	}
	return rk
}

func sortedCoreNavEntry(mc *JMapCore, k any, op int) any {
	rk, rv, ok := sortedCoreNav(mc, k, op)
	if !ok {
		return nil
	}
	return &jSimpleEntry{k: rk, v: rv, imm: true}
}

func mapCoreFirstEntry(mc *JMapCore) any {
	if mc == nil || mc.Live == 0 {
		return nil
	}
	k, v := mc.FirstEntry()
	return &jSimpleEntry{k: k, v: v, imm: true}
}

func mapCoreLastEntry(mc *JMapCore) any {
	if mc == nil || mc.Live == 0 {
		return nil
	}
	k := mc.LastKey()
	return &jSimpleEntry{k: k, v: mc.Get(k), imm: true}
}

func mapCorePollFirstEntry(mc *JMapCore) any {
	if mc == nil || mc.Live == 0 {
		return nil
	}
	k, v := mc.FirstEntry()
	mc.Remove(k)
	return &jSimpleEntry{k: k, v: v, imm: true}
}

func mapCorePollLastEntry(mc *JMapCore) any {
	if mc == nil || mc.Live == 0 {
		return nil
	}
	k := mc.LastKey()
	v := mc.Remove(k)
	return &jSimpleEntry{k: k, v: v, imm: true}
}

func comparatorFunc(cmp Ext_java_util_Comparator) func(a, b any) int32 {
	if isNaturalOrderComparator(cmp) {
		return CompareValues
	}
	return func(a, b any) int32 { return cmp.Compare(a, b) }
}

func comparatorFuncOrNatural(cmp Ext_java_util_Comparator) func(a, b any) int32 {
	return comparatorFunc(cmp)
}

func ExtDefault_Comparator_Reversed(this_ Ext_java_util_Comparator) Ext_java_util_Comparator {
	return newComparatorFromFunc(func(a, b any) int32 { return this_.Compare(b, a) })
}

func ExtDefault_Comparator_ThenComparing(this_ Ext_java_util_Comparator, other Ext_java_util_Comparator) Ext_java_util_Comparator {
	CheckNotNull(other)
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		return other.Compare(a, b)
	})
}

func ExtDefault_Comparator_ThenComparingKeyCmp(this_ Ext_java_util_Comparator, ke Ext_java_util_function_Function, kc Ext_java_util_Comparator) Ext_java_util_Comparator {
	CheckNotNull(ke)
	CheckNotNull(kc)
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		return kc.Compare(ke.Apply(a), ke.Apply(b))
	})
}

func ExtDefault_Comparator_ThenComparingKey(this_ Ext_java_util_Comparator, ke Ext_java_util_function_Function) Ext_java_util_Comparator {
	CheckNotNull(ke)
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		return CompareValues(ke.Apply(a), ke.Apply(b))
	})
}

func ExtDefault_Comparator_ThenComparingInt(this_ Ext_java_util_Comparator, ke Ext_java_util_function_ToIntFunction) Ext_java_util_Comparator {
	CheckNotNull(ke)
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		x, y := ke.ApplyAsInt(a), ke.ApplyAsInt(b)
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
}

func ExtDefault_Comparator_ThenComparingLong(this_ Ext_java_util_Comparator, ke Ext_java_util_function_ToLongFunction) Ext_java_util_Comparator {
	CheckNotNull(ke)
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		x, y := ke.ApplyAsLong(a), ke.ApplyAsLong(b)
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
}

func ExtDefault_Comparator_ThenComparingDouble(this_ Ext_java_util_Comparator, ke any) Ext_java_util_Comparator {
	CheckNotNull(ke)
	df := ke.(interface{ ApplyAsDouble(any) float64 })
	return newComparatorFromFunc(func(a, b any) int32 {
		r := this_.Compare(a, b)
		if r != 0 {
			return r
		}
		return int32(javalang.DoubleCompare(df.ApplyAsDouble(a), df.ApplyAsDouble(b)))
	})
}

type jPredFunc struct {
	fn     func(any) bool
	nodeFn func(uint32) bool
}

func (*jPredFunc) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*jPredFunc) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*jPredFunc) HasIdentityEquals_()                                     {}
func (p *jPredFunc) Test(x any) bool                                       { return p.fn(x) }
func (p *jPredFunc) Apply(x any) bool                                      { return p.fn(x) }
func (p *jPredFunc) TestNode_(n uint32) bool {
	if p == nil {
		ThrowNPE()
	}
	if p.nodeFn != nil {
		return p.nodeFn(n)
	}
	return p.fn(boxNodeById(n))
}
func (p *jPredFunc) Equals(o any) bool { return RefEquals(p, o) }

func ExtDefault_Predicate_Negate(this_ Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return &jPredFunc{
		fn:     func(x any) bool { return !this_.Test(x) },
		nodeFn: func(n uint32) bool { return !PredTestNode(this_, n) },
	}
}

func ExtDefault_Predicate_And(this_ Ext_java_util_function_Predicate, other Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNull(other)
	return &jPredFunc{
		fn:     func(x any) bool { return this_.Test(x) && other.Test(x) },
		nodeFn: func(n uint32) bool { return PredTestNode(this_, n) && PredTestNode(other, n) },
	}
}

func ExtDefault_Predicate_Or(this_ Ext_java_util_function_Predicate, other Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNull(other)
	return &jPredFunc{
		fn:     func(x any) bool { return this_.Test(x) || other.Test(x) },
		nodeFn: func(n uint32) bool { return PredTestNode(this_, n) || PredTestNode(other, n) },
	}
}

type jFuncIdentity struct{}

func (*jFuncIdentity) Implements_j2gort_Ext_java_util_function_Function()     {}
func (*jFuncIdentity) Implements_j2gort_Ext_com_google_common_base_Function() {}
func (*jFuncIdentity) Apply(x any) any                                        { return x }
func (f *jFuncIdentity) Equals(o any) bool                                    { _, ok := o.(*jFuncIdentity); return ok }

type jFuncWrap struct {
	fn func(any) any
}

func (*jFuncWrap) Implements_j2gort_Ext_java_util_function_Function()     {}
func (*jFuncWrap) Implements_j2gort_Ext_com_google_common_base_Function() {}
func (f *jFuncWrap) Apply(x any) any                                      { return f.fn(x) }
func (f *jFuncWrap) Equals(o any) bool                                    { return RefEquals(f, o) }

func ExtDefault_Function_AndThen(this_ Ext_java_util_function_Function, after Ext_java_util_function_Function) Ext_java_util_function_Function {
	CheckNotNull(after)
	return &jFuncWrap{fn: func(x any) any { return after.Apply(this_.Apply(x)) }}
}

func ExtDefault_Function_Compose(this_ Ext_java_util_function_Function, before Ext_java_util_function_Function) Ext_java_util_function_Function {
	CheckNotNull(before)
	return &jFuncWrap{fn: func(x any) any { return this_.Apply(before.Apply(x)) }}
}

type jConsumerWrap struct {
	fn func(any)
}

func (*jConsumerWrap) Implements_j2gort_Ext_java_util_function_Consumer() {}
func (c *jConsumerWrap) Accept(x any)                                     { c.fn(x) }

func ExtDefault_Consumer_AndThen(this_ Ext_java_util_function_Consumer, after Ext_java_util_function_Consumer) Ext_java_util_function_Consumer {
	CheckNotNull(after)
	return &jConsumerWrap{fn: func(x any) { this_.Accept(x); after.Accept(x) }}
}

// --- Core Collection Equality, Hash, ToArray, String &Iterators ---

func sliceToArray1(items []any, dst []any) []any {
	CheckNotNullSlice(dst)
	n := len(items)
	if len(dst) < n {
		out := make([]any, n)
		copy(out, items)
		return out
	}
	copy(dst, items)
	if len(dst) > n {
		dst[n] = nil
	}
	return dst
}

func ColToStringArray(col any, dst []string) []string {
	if col == nil {
		ThrowNPE()
	}
	if ssi, ok := col.(StrSliceIterable); ok {
		if nc, ok2 := col.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			ThrowNPE()
		}
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			CheckNotNullSlice(dst)
			n := len(ss)
			if len(dst) < n {
				out := make([]string, n)
				copy(out, ss)
				return out
			}
			copy(dst, ss)
			if len(dst) > n {
				dst[n] = NullString
			}
			return dst
		}
	}
	if jsp, ok := col.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			CheckNotNullSlice(dst)
			n := int(sc.Size())
			var out []string
			if len(dst) < n {
				out = make([]string, n)
			} else {
				out = dst
				if len(dst) > n {
					dst[n] = NullString
				}
			}
			if n == 0 {
				return out
			}
			if sc.mode == jSetModeStr {
				strs := sc.ext.(*jSetStrState).strs
				w := 0
				for i := int(sc.firstIdx); i < len(strs); i++ {
					if s := strs[i]; s != jTombstoneStr {
						out[w] = s
						w++
					}
				}
				return out
			}
			w := 0
			sc.ForEach(func(k any) bool {
				out[w] = CastString(k)
				w++
				return true
			})
			return out
		}
	}
	if isJavaNilAny(col) {
		ThrowNPE()
	}
	CheckNotNullSlice(dst)
	if sp, ok := col.(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil && mc.mode == jMapModeStr {
			n := int(mc.Size())
			var out []string
			if len(dst) < n {
				out = make([]string, n)
			} else {
				out = dst
				if len(dst) > n {
					dst[n] = NullString
				}
			}
			w := 0
			for i := int(mc.firstIdx); i < len(mc.strEntries); i++ {
				if s := mc.strEntries[i].key; s != jTombstoneStr {
					out[w] = s
					w++
				}
			}
			return out
		}
	}
	capHint := 0
	if c, ok := col.(Ext_java_util_Collection); ok {
		capHint = int(c.Size())
	}
	cur := NewIterCursor(col)
	if !cur.HasNext() {
		if len(dst) > 0 {
			dst[0] = NullString
		}
		return dst
	}
	var tmp []string
	if capHint > 0 {
		tmp = make([]string, 0, capHint)
	}
	for cur.HasNext() {
		tmp = append(tmp, cur.NextStr())
	}
	if len(dst) >= len(tmp) {
		copy(dst, tmp)
		if len(dst) > len(tmp) {
			dst[len(tmp)] = NullString
		}
		return dst
	}
	return tmp
}

func ColToNewStringArray(col any) []string {
	return ColToStringArray(col, emptyStrSlice_)
}

func ColToPtrArray[T any](col any, dst []*T) []*T {
	if col == nil {
		ThrowNPE()
	}
	if si, ok := col.(SliceIterable); ok {
		if nc, ok2 := col.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			ThrowNPE()
		}
		items := si.IterSlice_()
		CheckNotNullSlice(dst)
		n := len(items)
		var out []*T
		if len(dst) < n {
			out = make([]*T, n)
		} else {
			out = dst
			if len(dst) > n {
				dst[n] = nil
			}
		}
		for i, x := range items {
			out[i] = CastPtr[T](x)
		}
		return out
	}
	if jsp, ok := col.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			CheckNotNullSlice(dst)
			n := int(sc.Size())
			var out []*T
			if len(dst) < n {
				out = make([]*T, n)
			} else {
				out = dst
				if len(dst) > n {
					dst[n] = nil
				}
			}
			if n == 0 {
				return out
			}
			if sc.mode == jSetModeAny {
				items := sc.ext.(*jSetAnyState).items
				w := 0
				for i := int(sc.firstIdx); i < len(items); i++ {
					if it := items[i]; it != jTombstoneObj {
						out[w] = CastPtr[T](it)
						w++
					}
				}
				return out
			}
			w := 0
			sc.ForEach(func(k any) bool {
				out[w] = CastPtr[T](k)
				w++
				return true
			})
			return out
		}
	}
	if isJavaNilAny(col) {
		ThrowNPE()
	}
	CheckNotNullSlice(dst)
	capHint := 0
	if c, ok := col.(Ext_java_util_Collection); ok {
		capHint = int(c.Size())
	}
	cur := NewIterCursor(col)
	if !cur.HasNext() {
		if len(dst) > 0 {
			dst[0] = nil
		}
		return dst
	}
	var tmp []*T
	if capHint > 0 {
		tmp = make([]*T, 0, capHint)
	}
	for cur.HasNext() {
		tmp = append(tmp, CastPtr[T](cur.Next()))
	}
	if len(dst) >= len(tmp) {
		copy(dst, tmp)
		if len(dst) > len(tmp) {
			dst[len(tmp)] = nil
		}
		return dst
	}
	return tmp
}

func ColToNewPtrArray[T any](col any) []*T {
	return ColToPtrArray(col, []*T{})
}

type jFastSetCoreProvider interface {
	fastSetCore_() *JMapCore
}

type jFastJSetCoreProvider interface {
	fastSetJCore_() *JSetCore
}

type jFastMapCoreProvider interface {
	fastMapCore_() *JMapCore
}

type jFastEnumSetProvider interface {
	fastEnumSet_() *jEnumSetCore
}

type jFastEnumMapProvider interface {
	fastEnumMap_() *jEnumMapCore
}

func immSliceContains(items []any, target any) bool {
	if len(items) == 0 {
		return false
	}
	t := ToAny(target)
	if t == nil {
		return false
	}
	for _, x := range items {
		if elemJavaEquals(x, t, false) {
			return true
		}
	}
	return false
}

func immSliceIndexOf(items []any, target any) int32 {
	if len(items) == 0 {
		return -1
	}
	t := ToAny(target)
	if t == nil {
		return -1
	}
	for i, x := range items {
		if elemJavaEquals(x, t, false) {
			return int32(i)
		}
	}
	return -1
}

func immSliceLastIndexOf(items []any, target any) int32 {
	if len(items) == 0 {
		return -1
	}
	t := ToAny(target)
	if t == nil {
		return -1
	}
	for i := len(items) - 1; i >= 0; i-- {
		if elemJavaEquals(items[i], t, false) {
			return int32(i)
		}
	}
	return -1
}

func colContainsAll(containsFn func(any) bool, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if tv, ok := any(other).(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			for i := 0; i < n; i++ {
				if !containsFn(tv.typedViewAt_(i)) {
					return false
				}
			}
			return true
		}
	}
	if si, ok := any(other).(SliceIterable); ok {
		for _, x := range si.IterSlice_() {
			if !containsFn(x) {
				return false
			}
		}
		return true
	}
	if jsp, ok := any(other).(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			all := true
			sc.ForEach(func(k any) bool {
				if !containsFn(k) {
					all = false
					return false
				}
				return true
			})
			return all
		}
	}
	if sp, ok := any(other).(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil {
			all := true
			mc.ForEachEntry(func(k, _ any) bool {
				if !containsFn(k) {
					all = false
					return false
				}
				return true
			})
			return all
		}
	}
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			return es.ForEach(func(k any) bool { return containsFn(k) })
		}
	}
	for it := NewIterCursor(other); it.HasNext(); {
		if !containsFn(it.Next()) {
			return false
		}
	}
	return true
}

func listCoreAddAllIterable(lc *JListCore, src any) bool {
	CheckNotNull(src)
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return false
			}
			lc.InitCap(lc.Size() + int32(n))
			for i := 0; i < n; i++ {
				lc.Add(tv.typedViewAt_(i))
			}
			return true
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return false
		}
		lc.InitCap(lc.Size() + int32(len(items)))
		for _, x := range items {
			lc.Add(x)
		}
		return true
	}
	if mv, ok := src.(*jLiveMapValues); ok {
		if mv == nil {
			ThrowNPE()
		}
		if mv.mc.Size() == 0 {
			return false
		}
		lc.InitCap(lc.Size() + mv.mc.Size())
		mv.mc.ForEachValue(func(v any) bool {
			lc.Add(v)
			return true
		})
		return true
	}
	if mk, ok := src.(*jLiveMapKeySet); ok {
		if mk == nil {
			ThrowNPE()
		}
		if mk.mc.Size() == 0 {
			return false
		}
		lc.InitCap(lc.Size() + mk.mc.Size())
		mk.mc.ForEachEntry(func(k, _ any) bool {
			lc.Add(k)
			return true
		})
		return true
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			if sc.Size() == 0 {
				return false
			}
			lc.InitCap(lc.Size() + sc.Size())
			sc.ForEach(func(k any) bool {
				lc.Add(k)
				return true
			})
			return true
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil {
			if mc.Size() == 0 {
				return false
			}
			lc.InitCap(lc.Size() + mc.Size())
			mc.ForEachEntry(func(k, _ any) bool {
				lc.Add(k)
				return true
			})
			return true
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return false
			}
			lc.InitCap(lc.Size() + es.size)
			es.ForEach(func(k any) bool {
				lc.Add(k)
				return true
			})
			return true
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		n := c.Size()
		if n == 0 {
			return false
		}
		lc.InitCap(lc.Size() + n)
	}
	added := false
	for it := NewIterCursor(src); it.HasNext(); {
		lc.Add(it.Next())
		added = true
	}
	return added
}

func listCoreAddAllNonNull(lc *JListCore, src any) {
	CheckNotNull(src)
	if ssi, ok := src.(StrSliceIterable); ok {
		if strs, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(strs) == 0 {
				return
			}
			lc.InitCap(lc.Size() + int32(len(strs)))
			for _, s := range strs {
				if s == NullString {
					ThrowNPE()
				}
				lc.Add(s)
			}
			return
		}
	}
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return
			}
			lc.InitCap(lc.Size() + int32(n))
			for i := 0; i < n; i++ {
				lc.Add(CheckNotNull(ToAny(tv.typedViewAt_(i))))
			}
			return
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return
		}
		lc.InitCap(lc.Size() + int32(len(items)))
		for _, x := range items {
			lc.Add(CheckNotNull(ToAny(x)))
		}
		return
	}
	if mv, ok := src.(*jLiveMapValues); ok {
		if mv == nil {
			ThrowNPE()
		}
		if mv.mc.Size() == 0 {
			return
		}
		lc.InitCap(lc.Size() + mv.mc.Size())
		mv.mc.ForEachValue(func(v any) bool {
			lc.Add(CheckNotNull(v))
			return true
		})
		return
	}
	if mk, ok := src.(*jLiveMapKeySet); ok {
		if mk == nil {
			ThrowNPE()
		}
		if mk.mc.Size() == 0 {
			return
		}
		lc.InitCap(lc.Size() + mk.mc.Size())
		mk.mc.ForEachEntry(func(k, _ any) bool {
			lc.Add(CheckNotNull(k))
			return true
		})
		return
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			if sc.Size() == 0 {
				return
			}
			lc.InitCap(lc.Size() + sc.Size())
			sc.ForEach(func(k any) bool {
				lc.Add(CheckNotNull(k))
				return true
			})
			return
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil {
			if mc.Size() == 0 {
				return
			}
			lc.InitCap(lc.Size() + mc.Size())
			mc.ForEachEntry(func(k, _ any) bool {
				lc.Add(CheckNotNull(k))
				return true
			})
			return
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return
			}
			lc.InitCap(lc.Size() + es.size)
			es.ForEach(func(k any) bool {
				lc.Add(k)
				return true
			})
			return
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		lc.InitCap(lc.Size() + c.Size())
	}
	for it := NewIterCursor(src); it.HasNext(); {
		lc.Add(CheckNotNull(it.Next()))
	}
}

func listCoreRemoveAll(lc *JListCore, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	return listCoreRemoveIf(lc, func(x any) bool { return other.Contains(x) })
}

func listCoreRemoveIf(lc *JListCore, pred func(any) bool) bool {
	changed := false
	w := 0
	for _, x := range lc.Items {
		if pred(x) {
			changed = true
		} else {
			lc.Items[w] = x
			w++
		}
	}
	for i := w; i < len(lc.Items); i++ {
		lc.Items[i] = nil
	}
	lc.Items = lc.Items[:w]
	return changed
}

func listCoreRemoveLastVal(lc *JListCore, v any) bool {
	idx := lc.LastIndexOf(v)
	if idx < 0 {
		return false
	}
	lc.RemoveAt(idx)
	return true
}

func listCoreSort(lc *JListCore, cmp Ext_java_util_Comparator) {
	fn := comparatorFuncOrNatural(cmp)
	sort.SliceStable(lc.Items, func(i, j int) bool {
		return fn(lc.Items[i], lc.Items[j]) < 0
	})
}

func listEqualsSlice(items []any, other any) bool {
	lst, ok := other.(Ext_java_util_List)
	if !ok || lst == nil || javalang.IsNil(lst) {
		return false
	}
	if int32(len(items)) != lst.Size() {
		return false
	}
	it := Iterate(lst)
	for _, x := range items {
		if !it.HasNext() || !javalang.ObjectsEquals(ToAny(x), ToAny(it.Next())) {
			return false
		}
	}
	return true
}

func listHashSlice(items []any) int32 {
	h := int32(1)
	for _, x := range items {
		h = 31*h + javalang.ObjectHashCode(ToAny(x))
	}
	return h
}

func sliceJavaString(items []any) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, x := range items {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(x))
	}
	sb.WriteByte(']')
	return sb.String()
}

func copySetCoreInto(dst *JSetCore, src *JSetCore) *JSetCore {
	dst.identity = src.identity
	if src.Size() == 0 {
		return dst
	}
	dst.InitExpectedSize(src.Size())
	if src.mode == jSetModeStr {
		strs := src.ext.(*jSetStrState).strs
		for i := int(src.firstIdx); i < len(strs); i++ {
			if s := strs[i]; s != jTombstoneStr {
				dst.AddStr(s)
			}
		}
		return dst
	}
	if src.mode == jSetModeNode {
		for i := int(src.firstIdx); i < len(src.nodes); i++ {
			if u := src.nodes[i]; u != jTombstoneNode {
				dst.AddNode(u)
			}
		}
		return dst
	}
	src.ForEach(func(k any) bool {
		dst.Add(k)
		return true
	})
	return dst
}

func setCoreAddAllIterable(sc *JSetCore, src any) bool {
	CheckNotNull(src)
	if ssi, ok := src.(StrSliceIterable); ok {
		if strs, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(strs) == 0 {
				return false
			}
			sc.InitExpectedSize(sc.Size() + int32(len(strs)))
			changed := false
			for _, s := range strs {
				if s != jTombstoneStr && sc.AddStr(s) {
					changed = true
				}
			}
			return changed
		}
	}
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return false
			}
			sc.InitExpectedSize(sc.Size() + int32(n))
			changed := false
			for i := 0; i < n; i++ {
				if sc.Add(tv.typedViewAt_(i)) {
					changed = true
				}
			}
			return changed
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return false
		}
		sc.InitExpectedSize(sc.Size() + int32(len(items)))
		changed := false
		for _, x := range items {
			if sc.Add(x) {
				changed = true
			}
		}
		return changed
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
			if osc.Size() == 0 {
				return false
			}
			sc.InitExpectedSize(sc.Size() + osc.Size())
			changed := false
			if osc.mode == jSetModeStr {
				strs := osc.ext.(*jSetStrState).strs
				for i := int(osc.firstIdx); i < len(strs); i++ {
					if s := strs[i]; s != jTombstoneStr && sc.AddStr(s) {
						changed = true
					}
				}
				return changed
			}
			if osc.mode == jSetModeNode {
				for i := int(osc.firstIdx); i < len(osc.nodes); i++ {
					if u := osc.nodes[i]; u != jTombstoneNode && sc.AddNode(u) {
						changed = true
					}
				}
				return changed
			}
			osc.ForEach(func(k any) bool {
				if sc.Add(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if smc := sp.fastSetCore_(); smc != nil {
			if smc.Size() == 0 {
				return false
			}
			sc.InitExpectedSize(sc.Size() + smc.Size())
			changed := false
			if smc.mode == jMapModeStr && !smc.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					if k := smc.strEntries[i].key; k != jTombstoneStr {
						if sc.AddStr(k) {
							changed = true
						}
					}
				}
				return changed
			}
			smc.ForEachEntry(func(k, _ any) bool {
				if sc.Add(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return false
			}
			sc.InitExpectedSize(sc.Size() + es.size)
			changed := false
			es.ForEach(func(k any) bool {
				if sc.Add(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		n := c.Size()
		if n == 0 {
			return false
		}
		sc.InitExpectedSize(sc.Size() + n)
	}
	changed := false
	for it := NewIterCursor(src); it.HasNext(); {
		if sc.Add(it.Next()) {
			changed = true
		}
	}
	return changed
}

func setCoreAddAllNonNull(sc *JSetCore, src any) {
	CheckNotNull(src)
	if ssi, ok := src.(StrSliceIterable); ok {
		if strs, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(strs) == 0 {
				return
			}
			sc.InitExpectedSize(sc.Size() + int32(len(strs)))
			for _, s := range strs {
				if s == NullString {
					ThrowNPE()
				}
				sc.AddStr(s)
			}
			return
		}
	}
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return
			}
			sc.InitExpectedSize(sc.Size() + int32(n))
			for i := 0; i < n; i++ {
				sc.Add(CheckNotNull(ToAny(tv.typedViewAt_(i))))
			}
			return
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return
		}
		sc.InitExpectedSize(sc.Size() + int32(len(items)))
		for _, x := range items {
			sc.Add(CheckNotNull(ToAny(x)))
		}
		return
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
			if osc.Size() == 0 {
				return
			}
			sc.InitExpectedSize(sc.Size() + osc.Size())
			osc.ForEach(func(k any) bool {
				sc.Add(CheckNotNull(k))
				return true
			})
			return
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if smc := sp.fastSetCore_(); smc != nil {
			if smc.Size() == 0 {
				return
			}
			sc.InitExpectedSize(sc.Size() + smc.Size())
			if smc.mode == jMapModeStr && !smc.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					k := smc.strEntries[i].key
					if k == jTombstoneStr {
						continue
					}
					if k == NullString {
						ThrowNPE()
					}
					sc.AddStr(k)
				}
				return
			}
			smc.ForEachEntry(func(k, _ any) bool {
				sc.Add(CheckNotNull(k))
				return true
			})
			return
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return
			}
			sc.InitExpectedSize(sc.Size() + es.size)
			es.ForEach(func(k any) bool {
				sc.Add(k)
				return true
			})
			return
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		sc.InitExpectedSize(sc.Size() + c.Size())
	}
	for it := NewIterCursor(src); it.HasNext(); {
		sc.Add(CheckNotNull(it.Next()))
	}
}

func hashSetAddAll(hs *Ext_java_util_HashSetBase, src any) bool {
	if mb := hs.mapBack(); mb != nil {
		CheckNotNull(src)
		changed := false
		for it := Iterate(src); it.HasNext(); {
			if mb.Put_2(it.Next(), ValueOfBoolean(true)) == nil {
				changed = true
			}
		}
		return changed
	}
	return setCoreAddAllIterable(&hs.inlineSet, src)
}

func hashSetContainsAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
	if mb := hs.mapBack(); mb != nil {
		return colContainsAll(mb.ContainsKey, other)
	}
	return colContainsAll(hs.inlineSet.Contains, other)
}

func hashSetRemoveAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if mb := hs.mapBack(); mb != nil {
		changed := false
		for it := Iterate(other); it.HasNext(); {
			if mb.Remove_1(it.Next()) != nil {
				changed = true
			}
		}
		return changed
	}
	sc := &hs.inlineSet
	if sc.Size() == 0 {
		return false
	}
	if tv, ok := any(other).(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			changed := false
			for i := 0; i < n; i++ {
				if sc.Remove(tv.typedViewAt_(i)) {
					changed = true
				}
			}
			return changed
		}
	}
	if si, ok := any(other).(SliceIterable); ok {
		changed := false
		for _, x := range si.IterSlice_() {
			if sc.Remove(x) {
				changed = true
			}
		}
		return changed
	}
	if jsp, ok := any(other).(jFastJSetCoreProvider); ok {
		if osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
			changed := false
			osc.ForEach(func(k any) bool {
				if sc.Remove(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	changed := false
	for it := Iterate(other); it.HasNext(); {
		if sc.Remove(it.Next()) {
			changed = true
		}
	}
	return changed
}

func hashSetRetainAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if mb := hs.mapBack(); mb != nil {
		return abstractColRetainAll(Cast[Ext_java_util_Set](mb.KeySet_0()), other)
	}
	sc := &hs.inlineSet
	changed := false
	for _, k := range sc.LiveKeys() {
		if !other.Contains(k) && sc.Remove(k) {
			changed = true
		}
	}
	return changed
}

func hashSetRemoveIf(hs *Ext_java_util_HashSetBase, pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	if mb := hs.mapBack(); mb != nil {
		return abstractColRemoveIf(Cast[Ext_java_util_Set](mb.KeySet_0()), pred)
	}
	sc := &hs.inlineSet
	changed := false
	for _, k := range sc.LiveKeys() {
		if pred.Test(k) && sc.Remove(k) {
			changed = true
		}
	}
	return changed
}

func hashSetToArray(hs *Ext_java_util_HashSetBase) []any {
	if mb := hs.mapBack(); mb != nil {
		return colToSlice(Cast[Ext_java_util_Set](mb.KeySet_0()))
	}
	return hs.inlineSet.LiveKeys()
}

func hashSetForEach(hs *Ext_java_util_HashSetBase, action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	if mb := hs.mapBack(); mb != nil {
		for it := Iterate(mb.KeySet_0()); it.HasNext(); {
			action.Accept(it.Next())
		}
		return
	}
	hs.inlineSet.ForEach(func(k any) bool {
		action.Accept(k)
		return true
	})
}

func hashSetStream(hs *Ext_java_util_HashSetBase) Ext_java_util_stream_Stream {
	if mb := hs.mapBack(); mb != nil {
		return streamFromIterable(Cast[Ext_java_util_Set](mb.KeySet_0()))
	}
	sc := &hs.inlineSet
	if sc.IsEmpty() {
		return emptyStream_
	}
	var cur JIterCursor
	inited := false
	return &jStream{next: func() (any, bool) {
		if !inited {
			start := int(sc.firstIdx)
			cur = JIterCursor{sc: sc, idx: start, kind: 4}
			cur.advanceJSet()
			inited = true
		}
		if !cur.HasNext() {
			return nil, false
		}
		return ToAny(cur.Next()), true
	}}
}

func hashSetString(hs *Ext_java_util_HashSetBase) string {
	if mb := hs.mapBack(); mb != nil {
		return sliceJavaString(colToSlice(Cast[Ext_java_util_Set](mb.KeySet_0())))
	}
	sc := &hs.inlineSet
	if sc.IsEmpty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	if sc.mode == jSetModeStr {
		strs := sc.ext.(*jSetStrState).strs
		for i := int(sc.firstIdx); i < len(strs); i++ {
			s := strs[i]
			if s == jTombstoneStr {
				continue
			}
			if !first {
				sb.WriteString(", ")
			}
			first = false
			if s == NullString {
				sb.WriteString("null")
			} else {
				sb.WriteString(s)
			}
		}
		sb.WriteByte(']')
		return sb.String()
	}
	sc.ForEach(func(k any) bool {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(k))
		return true
	})
	sb.WriteByte(']')
	return sb.String()
}

func hashSetEquals(hs *Ext_java_util_HashSetBase, other any) bool {
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	if mb := hs.mapBack(); mb != nil {
		if mb.Size() != st.Size() {
			return false
		}
		return colContainsAll(mb.ContainsKey, st)
	}
	sc := &hs.inlineSet
	if sc.Size() != st.Size() {
		return false
	}
	if sc.Size() == 0 {
		return true
	}
	if jsp, ok := other.(jFastJSetCoreProvider); ok {
		if osc := jsp.fastSetJCore_(); osc != nil {
			if sc == osc {
				return true
			}
			all := true
			osc.ForEach(func(k any) bool {
				if !sc.Contains(k) {
					all = false
					return false
				}
				return true
			})
			return all
		}
	}
	return colContainsAll(sc.Contains, st)
}

func hashSetHash(hs *Ext_java_util_HashSetBase) int32 {
	if mb := hs.mapBack(); mb != nil {
		return javalang.ObjectHashCode(mb.KeySet_0())
	}
	return hs.inlineSet.SetHash()
}

type jSetIter struct {
	sc      *JSetCore
	idx     int
	lastVal any
	hasLast bool
}

func newHashSetIter(hs *Ext_java_util_HashSetBase) Ext_java_util_Iterator {
	if mb := hs.mapBack(); mb != nil {
		return Cast[Ext_java_util_Iterator](Cast[Ext_java_util_Set](mb.KeySet_0()).Iterator())
	}
	sc := &hs.inlineSet
	if sc.Size() == 0 {
		return emptyMapIter_
	}
	start := int(sc.firstIdx)
	it := &jSetIter{sc: sc, idx: start}
	it.advance()
	if it.idx != start && !sc.imm {
		sc.firstIdx = int32(it.idx)
	}
	return it
}

func (*jSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jSetIter) sliceLen() int {
	switch it.sc.mode {
	case jSetModeNode:
		return len(it.sc.nodes)
	case jSetModeStr:
		return len(it.sc.ext.(*jSetStrState).strs)
	case jSetModeAny:
		return len(it.sc.ext.(*jSetAnyState).items)
	default:
		return 0
	}
}
func (it *jSetIter) advance() {
	switch it.sc.mode {
	case jSetModeNode:
		for it.idx < len(it.sc.nodes) && it.sc.nodes[it.idx] == jTombstoneNode {
			it.idx++
		}
	case jSetModeStr:
		strs := it.sc.ext.(*jSetStrState).strs
		for it.idx < len(strs) && strs[it.idx] == jTombstoneStr {
			it.idx++
		}
	case jSetModeAny:
		items := it.sc.ext.(*jSetAnyState).items
		for it.idx < len(items) && items[it.idx] == jTombstoneObj {
			it.idx++
		}
	}
}
func (it *jSetIter) HasNext() bool {
	it.advance()
	return it.idx < it.sliceLen()
}
func (it *jSetIter) Next() any {
	it.advance()
	if it.idx >= it.sliceLen() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	var v any
	switch it.sc.mode {
	case jSetModeNode:
		v = boxNodeById(it.sc.nodes[it.idx])
	case jSetModeStr:
		s := it.sc.ext.(*jSetStrState).strs[it.idx]
		if s != NullString {
			v = s
		}
	case jSetModeAny:
		v = it.sc.ext.(*jSetAnyState).items[it.idx]
	}
	it.idx++
	it.lastVal = v
	it.hasLast = true
	it.advance()
	return v
}
func (it *jSetIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	if it.sc.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	it.sc.Remove(it.lastVal)
	it.hasLast = false
}

const jImmEnumMaxBits = 256

type jEnumOrdAccessor interface {
	EnumOrdinal_() int32
	JavaClassInfo_() (string, string)
}

type immSetIdx struct {
	strIdx      map[string]struct{}
	identIdx    map[any]struct{}
	nodeIdx     map[uint32]struct{}
	big         *JSetCore
	enumClass   string
	enumBits    [4]uint64
	hasEnumBits bool
}

func tryInitImmEnumSingle(obj *Ext_com_google_common_collect_ImmutableSetBase, e0 any) {
	if e0 == nil {
		return
	}
	rt0 := reflect.TypeOf(e0)
	if rt0.Kind() == reflect.Uint8 {
		u := uint8(reflect.ValueOf(e0).Uint())
		if u != 0 {
			tab := new([256]bool)
			tab[u] = true
			obj.u8Tab_ = tab
			obj.u8PtrNil_ = reflect.Zero(reflect.PointerTo(rt0)).Interface()
			obj.hasU8Bits_ = true
		}
	}
}

func tryInitImmEnumMulti(obj *Ext_com_google_common_collect_ImmutableSetBase, elems []any) bool {
	if len(elems) == 0 || elems[0] == nil {
		return false
	}
	rt0 := reflect.TypeOf(elems[0])
	if rt0.Kind() == reflect.Uint8 {
		u0 := uint8(reflect.ValueOf(elems[0]).Uint())
		if u0 == 0 {
			return false
		}
		tab := new([256]bool)
		tab[u0] = true
		for i := 1; i < len(elems); i++ {
			ei := elems[i]
			if ei == nil || reflect.TypeOf(ei) != rt0 {
				return false
			}
			u := uint8(reflect.ValueOf(ei).Uint())
			if u == 0 {
				return false
			}
			tab[u] = true
		}
		obj.u8Tab_ = tab
		obj.u8PtrNil_ = reflect.Zero(reflect.PointerTo(rt0)).Interface()
		obj.hasU8Bits_ = true
		return true
	}
	if len(elems) <= 4 {
		return false
	}
	e0, ok := elems[0].(jEnumOrdAccessor)
	if !ok {
		return false
	}
	ord0 := e0.EnumOrdinal_()
	if ord0 < 0 || ord0 >= jImmEnumMaxBits {
		return false
	}
	cls0, _ := e0.JavaClassInfo_()
	var maskBits [4]uint64
	maskBits[uint(ord0)>>6] |= uint64(1) << (uint(ord0) & 63)
	for i := 1; i < len(elems); i++ {
		ei := elems[i]
		if ei == nil {
			return false
		}
		eai, ok2 := ei.(jEnumOrdAccessor)
		if !ok2 {
			return false
		}
		clsi, _ := eai.JavaClassInfo_()
		if clsi != cls0 {
			return false
		}
		ord := eai.EnumOrdinal_()
		if ord < 0 || ord >= jImmEnumMaxBits {
			return false
		}
		maskBits[uint(ord)>>6] |= uint64(1) << (uint(ord) & 63)
	}
	obj.idx_.Store(&immSetIdx{
		enumClass:   cls0,
		enumBits:    maskBits,
		hasEnumBits: true,
	})
	return true
}

func isAllFastImmSetElems(elems []any) bool {
	allStr := true
	for _, e := range elems {
		str, ok := e.(string)
		if !ok || str == NullString {
			allStr = false
			break
		}
	}
	if allStr {
		return true
	}
	allNode := true
	for _, e := range elems {
		if e == nil {
			allNode = false
			break
		}
		if rn, ok := e.(jRawNodeProvider); ok {
			if rn.RawNodeId_() == 0 {
				allNode = false
				break
			}
		} else if nd, ok := e.(nodestore.Node); ok {
			if nd == 0 {
				allNode = false
				break
			}
		} else {
			allNode = false
			break
		}
	}
	if allNode {
		return true
	}
	for _, e := range elems {
		if !isFastIdentElem(e) {
			return false
		}
	}
	return true
}

func finishImmSetMulti(elems []any, scOwned *JSetCore) *Ext_com_google_common_collect_ImmutableSetBase {
	obj := &Ext_com_google_common_collect_ImmutableSetBase{elems_: elems[:len(elems):len(elems)], u8Tab_: &emptyU8Tab_}
	obj.BindExtSelf_(obj)
	if tryInitImmEnumMulti(obj, elems) {
		return obj
	}
	if len(elems) > 8 && scOwned != nil && !isAllFastImmSetElems(elems) {
		scOwned.sealForImmutable()
		obj.idx_.Store(&immSetIdx{big: scOwned})
	}
	return obj
}

func singletonElemMatchesNode(e any, n uint32) bool {
	if e == nil || n == 0 {
		return false
	}
	if rn, ok := e.(jRawNodeProvider); ok {
		return rn.RawNodeId_() == n
	}
	if nd, ok := e.(nodestore.Node); ok {
		return uint32(nd) == n
	}
	return false
}

func colSetEqualsSingleton(elem any, other any) bool {
	if other == nil {
		return false
	}
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	if st.Size() != 1 {
		return false
	}
	if si, ok := other.(SliceIterable); ok {
		s := si.IterSlice_()
		return len(s) == 1 && jSetElemEquals(elem, s[0])
	}
	it := Iterate(st)
	return it.HasNext() && jSetElemEquals(elem, it.Next())
}

func strsToAnySlice(strs []string) []any {
	if len(strs) == 0 {
		return []any{}
	}
	out := make([]any, len(strs))
	for i, s := range strs {
		if s != NullString {
			out[i] = s
		}
	}
	return out
}

func immStrSliceContains(strs []string, arg0 any) bool {
	s, ok := arg0.(string)
	if !ok || s == NullString {
		return false
	}
	for _, x := range strs {
		if x == s {
			return true
		}
	}
	return false
}

func immStrSliceIndexOf(strs []string, arg0 any) int32 {
	s, ok := arg0.(string)
	if !ok || s == NullString {
		return -1
	}
	for i, x := range strs {
		if x == s {
			return int32(i)
		}
	}
	return -1
}

func immStrSliceLastIndexOf(strs []string, arg0 any) int32 {
	s, ok := arg0.(string)
	if !ok || s == NullString {
		return -1
	}
	for i := len(strs) - 1; i >= 0; i-- {
		if strs[i] == s {
			return int32(i)
		}
	}
	return -1
}

func listEqualsStrSlice(strs []string, other any) bool {
	if other == nil {
		return false
	}
	lst, ok := other.(Ext_java_util_List)
	if !ok || lst == nil || javalang.IsNil(lst) {
		return false
	}
	n := len(strs)
	if int32(n) != lst.Size() {
		return false
	}
	if n == 0 {
		return true
	}
	if ssi, ok2 := other.(StrSliceIterable); ok2 {
		if os, ok3 := ssi.IterStrSlice_(); ok3 {
			if len(os) != n {
				return false
			}
			for i := 0; i < n; i++ {
				if strs[i] != os[i] {
					return false
				}
			}
			return true
		}
	}
	if si, ok2 := other.(SliceIterable); ok2 {
		os := si.IterSlice_()
		if len(os) != n {
			return false
		}
		for i := 0; i < n; i++ {
			es, ok3 := os[i].(string)
			if !ok3 || es != strs[i] {
				return false
			}
		}
		return true
	}
	it := Iterate(lst)
	for i := 0; i < n; i++ {
		if !it.HasNext() {
			return false
		}
		es, ok2 := it.Next().(string)
		if !ok2 || es != strs[i] {
			return false
		}
	}
	return !it.HasNext()
}

func listHashStrSlice(strs []string) int32 {
	h := int32(1)
	for _, s := range strs {
		if s == NullString {
			h = 31 * h
		} else {
			h = 31*h + javalang.StringHashCode(s)
		}
	}
	return h
}

func strSliceJavaString(strs []string) string {
	if len(strs) == 0 {
		return "[]"
	}
	if len(strs) == 1 {
		if strs[0] == NullString {
			return "[null]"
		}
		return "[" + strs[0] + "]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, s := range strs {
		if i > 0 {
			sb.WriteString(", ")
		}
		if s == NullString {
			sb.WriteString("null")
		} else {
			sb.WriteString(s)
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

func immStrSetHash(strs []string) int32 {
	var h int32
	for _, s := range strs {
		if s != NullString {
			h += javalang.StringHashCode(s)
		}
	}
	return h
}

func (s *Ext_com_google_common_collect_ImmutableSetBase) getOrBuildIdx() *immSetIdx {
	if idx := s.idx_.Load(); idx != nil {
		return idx
	}
	return s.buildIdxSlow()
}

//go:noinline
func (s *Ext_com_google_common_collect_ImmutableSetBase) buildIdxSlow() *immSetIdx {
	built := &immSetIdx{}
	if len(s.strs_) > 4 {
		idx := make(map[string]struct{}, len(s.strs_))
		for _, str := range s.strs_ {
			idx[str] = struct{}{}
		}
		built.strIdx = idx
	} else if len(s.elems_) > 4 {
		allStr := true
		for _, e := range s.elems_ {
			str, ok := e.(string)
			if !ok || str == NullString {
				allStr = false
				break
			}
		}
		if allStr {
			idx := make(map[string]struct{}, len(s.elems_))
			for _, e := range s.elems_ {
				idx[e.(string)] = struct{}{}
			}
			built.strIdx = idx
		} else {
			allNode := true
			for _, e := range s.elems_ {
				if e == nil {
					allNode = false
					break
				}
				if rn, ok := e.(jRawNodeProvider); ok {
					if rn.RawNodeId_() == 0 {
						allNode = false
						break
					}
				} else if nd, ok := e.(nodestore.Node); ok {
					if nd == 0 {
						allNode = false
						break
					}
				} else {
					allNode = false
					break
				}
			}
			if allNode {
				idx := make(map[uint32]struct{}, len(s.elems_))
				for _, e := range s.elems_ {
					if rn, ok := e.(jRawNodeProvider); ok {
						idx[rn.RawNodeId_()] = struct{}{}
					} else {
						idx[uint32(e.(nodestore.Node))] = struct{}{}
					}
				}
				built.nodeIdx = idx
			} else {
				allIdent := true
				for _, e := range s.elems_ {
					if !isFastIdentElem(e) {
						allIdent = false
						break
					}
				}
				if allIdent {
					idx := make(map[any]struct{}, len(s.elems_))
					for _, e := range s.elems_ {
						idx[e] = struct{}{}
					}
					built.identIdx = idx
				} else if len(s.elems_) > 8 {
					big := &JSetCore{}
					big.InitExpectedSize(int32(len(s.elems_)))
					for _, e := range s.elems_ {
						big.Add(e)
					}
					big.sealForImmutable()
					built.big = big
				}
			}
		}
	}
	if s.idx_.CompareAndSwap(nil, built) {
		return built
	}
	return s.idx_.Load()
}

func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsStr(k string) bool {
	if k == NullString || s.hasU8Bits_ {
		return false
	}
	if s.strs_ != nil {
		if len(s.strs_) == 0 {
			return false
		}
		if len(s.strs_) > 4 {
			idx := s.getOrBuildIdx()
			_, hit := idx.strIdx[k]
			return hit
		}
		for _, es := range s.strs_ {
			if es == k {
				return true
			}
		}
		return false
	}
	elems := s.elems_
	if len(elems) == 0 {
		return false
	}
	if len(elems) <= 4 {
		for _, e := range elems {
			if es, ok := e.(string); ok && es == k {
				return true
			}
		}
		return false
	}
	idx := s.getOrBuildIdx()
	if idx.hasEnumBits {
		return false
	}
	if idx.strIdx != nil {
		_, hit := idx.strIdx[k]
		return hit
	}
	if idx.nodeIdx != nil || idx.identIdx != nil {
		return false
	}
	if idx.big != nil {
		return idx.big.ContainsStr(k)
	}
	for _, e := range elems {
		if es, ok := e.(string); ok && es == k {
			return true
		}
	}
	return false
}

func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsNode(n uint32) bool {
	if n == 0 || s.strs_ != nil || s.hasU8Bits_ {
		return false
	}
	elems := s.elems_
	if len(elems) == 0 {
		return false
	}
	if len(elems) <= 4 {
		for _, e := range elems {
			if singletonElemMatchesNode(e, n) {
				return true
			}
		}
		return false
	}
	idx := s.getOrBuildIdx()
	if idx.hasEnumBits {
		return false
	}
	if idx.nodeIdx != nil {
		_, hit := idx.nodeIdx[n]
		return hit
	}
	if idx.strIdx != nil || idx.identIdx != nil {
		return false
	}
	if idx.big != nil {
		return idx.big.ContainsNode(n)
	}
	for _, e := range elems {
		if singletonElemMatchesNode(e, n) {
			return true
		}
	}
	return false
}

func immSmallSetContains(elems []any, v any) bool {
	switch x := v.(type) {
	case string:
		if x == NullString {
			return false
		}
		for _, e := range elems {
			if es, ok := e.(string); ok && es == x {
				return true
			}
		}
		return false
	case jIdentityKey:
		if x.IsJavaNil_() {
			return false
		}
		for _, e := range elems {
			if e == v {
				return true
			}
		}
		return false
	case jEnumOrdAccessor:
		if x.EnumOrdinal_() < 0 {
			return false
		}
		for _, e := range elems {
			if e == v {
				return true
			}
		}
		return false
	case *BoxInt:
		if x == nil {
			return false
		}
		xv := x.Val
		for _, e := range elems {
			if eb, ok := e.(*BoxInt); ok {
				if eb != nil && eb.Val == xv {
					return true
				}
			} else if ei, ok2 := e.(int32); ok2 && ei == xv {
				return true
			}
		}
		return false
	case int32:
		for _, e := range elems {
			if ei, ok := e.(int32); ok {
				if ei == x {
					return true
				}
			} else if eb, ok2 := e.(*BoxInt); ok2 && eb != nil && eb.Val == x {
				return true
			}
		}
		return false
	case jRawNodeProvider:
		uid := x.RawNodeId_()
		if uid == 0 {
			return false
		}
		for _, e := range elems {
			if singletonElemMatchesNode(e, uid) {
				return true
			}
		}
		return false
	case nodestore.Node:
		if x == 0 {
			return false
		}
		uid := uint32(x)
		for _, e := range elems {
			if singletonElemMatchesNode(e, uid) {
				return true
			}
		}
		return false
	case FastCustomKey:
		if x.IsJavaNil_() {
			return false
		}
		for _, e := range elems {
			if x.Equals(e) {
				return true
			}
		}
		return false
	case IdentityComparable:
		if nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
			return false
		}
		for _, e := range elems {
			if e == v {
				return true
			}
		}
		return false
	case *javalang.Class:
		if x == nil {
			return false
		}
		for _, e := range elems {
			if e == v {
				return true
			}
		}
		return false
	}
	v = ToAny(v)
	if v == nil {
		return false
	}
	for _, e := range elems {
		if jSetElemEquals(e, v) {
			return true
		}
	}
	return false
}

func (b *Ext_com_google_common_collect_ImmutableSetBase) ContainsU8_(u uint8) bool {
	if b != nil && b.u8Tab_ != nil {
		return b.u8Tab_[u]
	}
	return b.immSetContainsU8Fallback(u)
}

//go:noinline
func (b *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsU8Fallback(u uint8) bool {
	if b == nil {
		ThrowNPE()
	}
	if b.hasU8Bits_ && b.u8Tab_ != nil {
		return b.u8Tab_[u]
	}
	return false
}

func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContains(v any) bool {
	if s.hasU8Bits_ {
		if s.u8PtrNil_ == nil || v == nil {
			return false
		}
		if ea, ok := v.(jEnumOrdAccessor); ok {
			u := ea.EnumOrdinal_() + 1
			if u <= 0 || u >= 256 || !s.u8Tab_[u] {
				return false
			}
			return reflect.TypeOf(v) == reflect.TypeOf(s.elems_[0])
		}
		return false
	}
	if s.strs_ != nil {
		if len(s.strs_) == 0 || v == nil {
			return false
		}
		str, ok := v.(string)
		if !ok || str == NullString {
			return false
		}
		if len(s.strs_) > 4 {
			idx := s.getOrBuildIdx()
			_, hit := idx.strIdx[str]
			return hit
		}
		for _, es := range s.strs_ {
			if es == str {
				return true
			}
		}
		return false
	}
	elems := s.elems_
	if len(elems) == 0 || v == nil {
		return false
	}
	if len(elems) <= 4 {
		return immSmallSetContains(elems, v)
	}
	idx := s.idx_.Load()
	if idx == nil {
		idx = s.buildIdxSlow()
	}
	if idx.hasEnumBits {
		ea, ok := v.(jEnumOrdAccessor)
		if !ok {
			return false
		}
		ord := ea.EnumOrdinal_()
		if ord < 0 || ord >= jImmEnumMaxBits || (idx.enumBits[uint(ord)>>6]&(uint64(1)<<(uint(ord)&63))) == 0 {
			return false
		}
		cls, _ := ea.JavaClassInfo_()
		return cls == idx.enumClass
	}
	if idx.strIdx != nil {
		str, ok := v.(string)
		if !ok || str == NullString {
			return false
		}
		_, hit := idx.strIdx[str]
		return hit
	}
	if idx.nodeIdx != nil {
		if rn, ok := v.(jRawNodeProvider); ok {
			u := rn.RawNodeId_()
			if u == 0 {
				return false
			}
			_, hit := idx.nodeIdx[u]
			return hit
		}
		if nd, ok := v.(nodestore.Node); ok {
			if nd == 0 {
				return false
			}
			_, hit := idx.nodeIdx[uint32(nd)]
			return hit
		}
		return false
	}
	if idx.identIdx != nil {
		if !isFastIdentElem(v) {
			return false
		}
		_, hit := idx.identIdx[v]
		return hit
	}
	if idx.big != nil {
		return idx.big.Contains(v)
	}
	return immSmallSetContains(elems, v)
}

func findSetCoreEqVal[K FastEqValKey](sc *JSetCore, k K, h int32, nonNull, isIdent bool) bool {
	if sc.live == 0 {
		return false
	}
	if !nonNull {
		switch sc.mode {
		case jSetModeNode:
			return sc.ContainsNode(0)
		case jSetModeStr:
			return sc.ContainsStr(NullString)
		case jSetModeAny:
			return sc.containsAnyNormalized(nil)
		default:
			return false
		}
	}
	if sc.mode != jSetModeAny {
		return false
	}
	st := sc.ext.(*jSetAnyState)
	if st.allIdent {
		if !isIdent {
			return false
		}
		for _, it := range st.items {
			if ek, ok := it.(K); ok && ek == k {
				return true
			}
		}
		return false
	}
	if g := st.gen; g != nil && len(g.buckets) > 0 {
		slot := smearHash(h) & uint32(len(g.buckets)-1)
		for cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
			i := cur - 1
			if g.hashes[i] == h {
				if ek, ok := st.items[i].(K); ok && ek == k {
					return true
				}
			}
		}
		return false
	}
	for _, it := range st.items {
		if ek, ok := it.(K); ok && ek == k {
			return true
		}
	}
	return false
}

func fastSetContainsEqVal[K FastEqValKey](s any, k K) (bool, bool) {
	if isb, ok := s.(*Ext_com_google_common_collect_ImmutableSetBase); ok {
		if isb == nil {
			ThrowNPE()
		}
		if isb.Self_ != nil && isb.Self_ != Ext_java_util_AbstractCollection(isb) {
			return false, false
		}
		if isb.strs_ != nil || isb.hasU8Bits_ || len(isb.elems_) == 0 {
			return false, true
		}
		h, nonNull, isIdent := k.FastValHash_()
		if !nonNull {
			return false, true
		}
		elems := isb.elems_
		if len(elems) <= 8 {
			for _, e := range elems {
				if ek, ok2 := e.(K); ok2 && ek == k {
					return true, true
				}
			}
			return false, true
		}
		idx := isb.getOrBuildIdx()
		if idx.hasEnumBits || idx.strIdx != nil || idx.nodeIdx != nil {
			return false, true
		}
		if idx.big != nil {
			return findSetCoreEqVal(idx.big, k, h, true, isIdent), true
		}
		for _, e := range elems {
			if ek, ok2 := e.(K); ok2 && ek == k {
				return true, true
			}
		}
		return false, true
	}
	if jsp, ok := s.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			h, nonNull, isIdent := k.FastValHash_()
			return findSetCoreEqVal(sc, k, h, nonNull, isIdent), true
		}
	}
	if sp, ok := s.(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil {
			_, found := findMapCoreEqVal(mc, k)
			return found, true
		}
	}
	return false, false
}

func immSetEquals(s *Ext_com_google_common_collect_ImmutableSetBase, other any) bool {
	if other == nil {
		return false
	}
	if os, ok := other.(*Ext_com_google_common_collect_ImmutableSetBase); ok && os != nil {
		if os == s {
			return true
		}
		if s.hasU8Bits_ && os.hasU8Bits_ {
			return s.u8PtrNil_ == os.u8PtrNil_ && len(s.elems_) == len(os.elems_) && *s.u8Tab_ == *os.u8Tab_
		}
	}
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	n := int32(len(s.elems_))
	if s.strs_ != nil {
		n = int32(len(s.strs_))
	}
	if n != st.Size() {
		return false
	}
	if n == 0 {
		return true
	}
	if s.strs_ != nil {
		if ssi, ok2 := other.(StrSliceIterable); ok2 {
			if os, ok3 := ssi.IterStrSlice_(); ok3 {
				for _, x := range os {
					if !s.immSetContainsStr(x) {
						return false
					}
				}
				return true
			}
		}
	}
	return colContainsAll(s.immSetContains, st)
}

func immSetHash(elems []any) int32 {
	var h int32
	for _, e := range elems {
		h += elemJavaHashCode(e, false)
	}
	return h
}

func setCoreEquals(mc *JMapCore, other any) bool {
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	if mc.Size() != st.Size() {
		return false
	}
	if mc.Size() == 0 {
		return true
	}
	return colContainsAll(mc.ContainsKey, st)
}

func setCoreHash(mc *JMapCore) int32 {
	return mc.SetHash()
}

func mapCoreAddAllIterable(mc *JMapCore, src any) bool {
	CheckNotNull(src)
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return false
			}
			mc.InitExpectedSize(mc.Size() + int32(n))
			changed := false
			for i := 0; i < n; i++ {
				if mc.AddKey(tv.typedViewAt_(i)) {
					changed = true
				}
			}
			return changed
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return false
		}
		mc.InitExpectedSize(mc.Size() + int32(len(items)))
		changed := false
		for _, x := range items {
			if mc.AddKey(x) {
				changed = true
			}
		}
		return changed
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			if sc.Size() == 0 {
				return false
			}
			mc.InitExpectedSize(mc.Size() + sc.Size())
			changed := false
			sc.ForEach(func(k any) bool {
				if mc.AddKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if smc := sp.fastSetCore_(); smc != nil {
			if smc.Size() == 0 {
				return false
			}
			mc.InitExpectedSize(mc.Size() + smc.Size())
			changed := false
			if smc.mode == jMapModeStr && !smc.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					if k := smc.strEntries[i].key; k != jTombstoneStr {
						if mc.AddKeyStr(k) {
							changed = true
						}
					}
				}
				return changed
			}
			smc.ForEachEntry(func(k, _ any) bool {
				if mc.AddKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return false
			}
			mc.InitExpectedSize(mc.Size() + es.size)
			changed := false
			es.ForEach(func(k any) bool {
				if mc.AddKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		n := c.Size()
		if n == 0 {
			return false
		}
		mc.InitExpectedSize(mc.Size() + n)
	}
	changed := false
	for it := NewIterCursor(src); it.HasNext(); {
		if mc.AddKey(it.Next()) {
			changed = true
		}
	}
	return changed
}

func mapCoreAddAllNonNull(mc *JMapCore, src any) {
	CheckNotNull(src)
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return
			}
			mc.InitExpectedSize(mc.Size() + int32(n))
			for i := 0; i < n; i++ {
				mc.Put(CheckNotNull(ToAny(tv.typedViewAt_(i))), true)
			}
			return
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return
		}
		mc.InitExpectedSize(mc.Size() + int32(len(items)))
		for _, x := range items {
			mc.Put(CheckNotNull(ToAny(x)), true)
		}
		return
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			if sc.Size() == 0 {
				return
			}
			mc.InitExpectedSize(mc.Size() + sc.Size())
			sc.ForEach(func(k any) bool {
				mc.Put(CheckNotNull(k), true)
				return true
			})
			return
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if smc := sp.fastSetCore_(); smc != nil {
			if smc.Size() == 0 {
				return
			}
			mc.InitExpectedSize(mc.Size() + smc.Size())
			if smc.mode == jMapModeStr && !smc.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					k := smc.strEntries[i].key
					if k == jTombstoneStr {
						continue
					}
					if k == NullString {
						ThrowNPE()
					}
					mc.PutStr(k, true)
				}
				return
			}
			smc.ForEachEntry(func(k, _ any) bool {
				mc.Put(CheckNotNull(k), true)
				return true
			})
			return
		}
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			if es.size == 0 {
				return
			}
			mc.InitExpectedSize(mc.Size() + es.size)
			es.ForEach(func(k any) bool {
				mc.Put(k, true)
				return true
			})
			return
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		mc.InitExpectedSize(mc.Size() + c.Size())
	}
	for it := NewIterCursor(src); it.HasNext(); {
		mc.Put(CheckNotNull(it.Next()), true)
	}
}

func mapCoreRemoveAllKeys(mc *JMapCore, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if mc.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if mc.Size() == 0 {
		return false
	}
	if tv, ok := any(other).(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			changed := false
			for i := 0; i < n; i++ {
				if mc.RemoveKey(tv.typedViewAt_(i)) {
					changed = true
				}
			}
			return changed
		}
	}
	if si, ok := any(other).(SliceIterable); ok {
		changed := false
		for _, x := range si.IterSlice_() {
			if mc.RemoveKey(x) {
				changed = true
			}
		}
		return changed
	}
	if jsp, ok := any(other).(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			changed := false
			sc.ForEach(func(k any) bool {
				if mc.RemoveKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if sp, ok := any(other).(jFastSetCoreProvider); ok {
		if omc := sp.fastSetCore_(); omc != nil && omc != mc {
			changed := false
			if omc.mode == jMapModeStr && !omc.Identity && !mc.Identity {
				for i := int(omc.firstIdx); i < len(omc.strEntries); i++ {
					if k := omc.strEntries[i].key; k != jTombstoneStr {
						if mc.RemoveKeyStr(k) {
							changed = true
						}
					}
				}
				return changed
			}
			omc.ForEachEntry(func(k, _ any) bool {
				if mc.RemoveKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			changed := false
			es.ForEach(func(k any) bool {
				if mc.RemoveKey(k) {
					changed = true
				}
				return true
			})
			return changed
		}
	}
	changed := false
	for it := Iterate(other); it.HasNext(); {
		if mc.RemoveKey(it.Next()) {
			changed = true
		}
	}
	return changed
}

func mapCoreRetainAllKeys(mc *JMapCore, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if mc.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	changed := false
	for _, k := range mc.LiveKeys() {
		if !other.Contains(k) && mc.RemoveKey(k) {
			changed = true
		}
	}
	return changed
}

func mapCoreRemoveIfKeys(mc *JMapCore, pred func(any) bool) bool {
	if mc.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	changed := false
	for _, k := range mc.LiveKeys() {
		if pred(k) && mc.RemoveKey(k) {
			changed = true
		}
	}
	return changed
}

func copyHashCoreInto(dst *JMapCore, src *JMapCore) *JMapCore {
	dst.Identity = src.Identity
	dst.genNonCustom = src.Identity
	dst.accessOrd = src.accessOrd
	if dst.CopyFromSameMode(src) {
		return dst
	}
	dst.InitExpectedSize(src.Size())
	if src.mode == jMapModeStr && !src.Identity && !dst.Identity {
		for i := int(src.firstIdx); i < len(src.strEntries); i++ {
			e := src.strEntries[i]
			if e.key != jTombstoneStr {
				dst.PutStr(e.key, e.val)
			}
		}
		return dst
	}
	if src.IsSorted() && src.ext != nil {
		st := src.ext.(*jSortedMapState)
		st.avlCopyInto(st.root, dst, false)
		return dst
	}
	src.ForEachEntry(func(k, v any) bool {
		dst.Put(k, v)
		return true
	})
	return dst
}

func copySortedCoreInto(dst *JMapCore, src *JMapCore) *JMapCore {
	if src.IsSorted() {
		dst.CopySortedFrom(src)
		return dst
	}
	initSortedCore(dst, nil)
	if src.mode == jMapModeStr && !src.Identity {
		for i := int(src.firstIdx); i < len(src.strEntries); i++ {
			e := src.strEntries[i]
			if e.key != jTombstoneStr {
				dst.PutStr(e.key, e.val)
			}
		}
		return dst
	}
	src.ForEachEntry(func(k, v any) bool {
		dst.Put(k, v)
		return true
	})
	return dst
}

func cloneJMapCore(src *JMapCore) *JMapCore {
	if src.IsSorted() {
		dst := &JMapCore{}
		dst.CopySortedFrom(src)
		return dst
	}
	dst := NewJMapCore(src.Identity)
	dst.accessOrd = src.accessOrd
	if dst.CopyFromSameMode(src) {
		return dst
	}
	dst.InitExpectedSize(src.Size())
	if src.mode == jMapModeStr && !src.Identity {
		for i := int(src.firstIdx); i < len(src.strEntries); i++ {
			e := src.strEntries[i]
			if e.key != jTombstoneStr {
				dst.PutStr(e.key, e.val)
			}
		}
		return dst
	}
	src.ForEachEntry(func(k, v any) bool {
		dst.Put(k, v)
		return true
	})
	return dst
}

func copyExtMapInto(dst *JMapCore, src Ext_java_util_Map) {
	CheckNotNull(src)
	if mp, ok := any(src).(jFastMapCoreProvider); ok {
		if smc := mp.fastMapCore_(); smc != nil {
			if smc == dst || smc.Size() == 0 {
				return
			}
			if dst.Size() == 0 && dst.CopyFromSameMode(smc) {
				return
			}
			dst.InitExpectedSize(dst.Size() + smc.Size())
			if smc.mode == jMapModeStr && !smc.Identity && !dst.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					e := smc.strEntries[i]
					if e.key != jTombstoneStr {
						dst.PutStr(e.key, e.val)
					}
				}
				return
			}
			if smc.IsSorted() && smc.ext != nil && !dst.Identity {
				st := smc.ext.(*jSortedMapState)
				st.avlCopyInto(st.root, dst, false)
				return
			}
			smc.ForEachEntry(func(k, v any) bool {
				dst.Put(k, v)
				return true
			})
			return
		}
	}
	if ep, ok := any(src).(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			if em.size == 0 {
				return
			}
			dst.InitExpectedSize(dst.Size() + em.size)
			em.ForEachEntry(func(k, v any) bool {
				dst.Put(k, v)
				return true
			})
			return
		}
	}
	if es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
		for it := Iterate(es); it.HasNext(); {
			e := Cast[Ext_java_util_Map_Entry](it.Next())
			dst.Put(e.GetKey(), e.GetValue())
		}
	}
}

func copyExtMapIntoNonNull(dst *JMapCore, src Ext_java_util_Map) {
	CheckNotNull(src)
	if mp, ok := any(src).(jFastMapCoreProvider); ok {
		if smc := mp.fastMapCore_(); smc != nil {
			if smc.Size() == 0 {
				return
			}
			if smc != dst && dst.Size() == 0 && dst.CopyFromSameModeNonNull(smc) {
				return
			}
			dst.InitExpectedSize(dst.Size() + smc.Size())
			if smc != dst && smc.mode == jMapModeStr && !smc.Identity && !dst.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					e := smc.strEntries[i]
					if e.key == jTombstoneStr {
						continue
					}
					if e.key == NullString {
						ThrowNPE()
					}
					dst.PutStr(e.key, CheckNotNull(e.val))
				}
				return
			}
			if smc != dst && smc.IsSorted() && smc.ext != nil && !dst.Identity {
				st := smc.ext.(*jSortedMapState)
				st.avlCopyInto(st.root, dst, true)
				return
			}
			smc.ForEachEntry(func(k, v any) bool {
				dst.Put(CheckNotNull(ToAny(k)), CheckNotNull(ToAny(v)))
				return true
			})
			return
		}
	}
	if ep, ok := any(src).(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			if em.size == 0 {
				return
			}
			dst.InitExpectedSize(dst.Size() + em.size)
			em.ForEachEntry(func(k, v any) bool {
				dst.Put(CheckNotNull(ToAny(k)), CheckNotNull(ToAny(v)))
				return true
			})
			return
		}
	}
	dst.InitExpectedSize(dst.Size() + src.Size())
	if es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
		for it := Iterate(es); it.HasNext(); {
			e := Cast[Ext_java_util_Map_Entry](it.Next())
			dst.Put(CheckNotNull(ToAny(e.GetKey())), CheckNotNull(ToAny(e.GetValue())))
		}
	}
}

func populateImmSortedSetElems(obj *Ext_com_google_common_collect_ImmutableSortedSetBase, mc *JMapCore) {
	n := mc.Size()
	if n == 0 {
		obj.elems_ = emptyAnySlice_
	} else if n == 1 {
		obj.singleArr_[0] = mc.FirstKey()
		obj.elems_ = obj.singleArr_[:1:1]
	} else {
		elems := mc.LiveKeys()
		obj.elems_ = elems[:len(elems):len(elems)]
	}
}

func setCoreDescending(mc *JMapCore, imm bool) any {
	origCmp := mc.SortedCmp()
	if origCmp == nil {
		origCmp = CompareValues
	}
	revCmp := func(a, b any) int32 { return origCmp(b, a) }
	keys := mc.LiveKeys()
	if imm {
		out := NewSortedJMapCore(revCmp)
		for i := len(keys) - 1; i >= 0; i-- {
			out.Put(keys[i], true)
		}
		out.sealImmutable()
		obj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: out}
		obj.BindExtSelf_(obj)
		populateImmSortedSetElems(obj, out)
		return obj
	}
	ts := &Ext_java_util_TreeSet{}
	ts.BindExtSelf_(ts)
	out := initSortedCore(&ts.inlineCore, revCmp)
	for i := len(keys) - 1; i >= 0; i-- {
		out.Put(keys[i], true)
	}
	return ts
}

func setCoreTailSet(mc *JMapCore, fromElem any, inclusive bool, imm bool) any {
	cmp := mc.SortedCmp()
	nat := mc.IsNaturalOrder()
	if cmp == nil {
		cmp = CompareValues
	}
	f := ToAny(fromElem)
	if imm {
		var out *JMapCore
		if nat {
			out = NewSortedJMapCore(nil)
		} else {
			out = NewSortedJMapCore(cmp)
		}
		for _, k := range mc.LiveKeys() {
			c := cmp(k, f)
			if c > 0 || (inclusive && c == 0) {
				out.Put(k, true)
			}
		}
		return wrapImmutableSortedSet(out)
	}
	ts := &Ext_java_util_TreeSet{}
	ts.BindExtSelf_(ts)
	var out *JMapCore
	if nat {
		out = initSortedCore(&ts.inlineCore, nil)
	} else {
		out = initSortedCore(&ts.inlineCore, cmp)
	}
	for _, k := range mc.LiveKeys() {
		c := cmp(k, f)
		if c > 0 || (inclusive && c == 0) {
			out.Put(k, true)
		}
	}
	return ts
}

// --- Iterators (ListIterator, Descending, MapKeyIter, SubList, Arrays.asList) ---

type jFullListIter struct {
	lc      *JListCore
	items   []any
	imm     bool
	cursor  int32
	lastRet int32
}

func newFullListIter(lc *JListCore, index int32) Ext_java_util_ListIterator {
	CheckPositionIndex(index, lc.Size())
	return &jFullListIter{lc: lc, cursor: index, lastRet: -1}
}

func newImmListIter(items []any, index int32) Ext_java_util_ListIterator {
	CheckPositionIndex(index, int32(len(items)))
	return &jFullListIter{items: items, imm: true, cursor: index, lastRet: -1}
}

func (*jFullListIter) Implements_j2gort_Ext_java_util_Iterator()                                 {}
func (*jFullListIter) Implements_j2gort_Ext_java_util_ListIterator()                             {}
func (*jFullListIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator()     {}
func (*jFullListIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableListIterator() {}
func (*jFullListIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (*jFullListIter) GetExt_com_google_common_collect_UnmodifiableListIteratorBase() *Ext_com_google_common_collect_UnmodifiableListIteratorBase {
	return nil
}

func (it *jFullListIter) listSize() int32 {
	if it.imm {
		return int32(len(it.items))
	}
	return it.lc.Size()
}
func (it *jFullListIter) HasNext() bool { return it.cursor < it.listSize() }
func (it *jFullListIter) Next() any {
	if it.cursor >= it.listSize() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	var v any
	if it.imm {
		v = it.items[it.cursor]
	} else {
		v = it.lc.Get(it.cursor)
	}
	it.lastRet = it.cursor
	it.cursor++
	return v
}
func (it *jFullListIter) HasPrevious() bool { return it.cursor > 0 }
func (it *jFullListIter) Previous() any {
	if it.cursor <= 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	it.cursor--
	it.lastRet = it.cursor
	if it.imm {
		return it.items[it.cursor]
	}
	return it.lc.Get(it.cursor)
}
func (it *jFullListIter) NextIndex() int32     { return it.cursor }
func (it *jFullListIter) PreviousIndex() int32 { return it.cursor - 1 }
func (it *jFullListIter) Remove() {
	if it.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.lc.RemoveAt(it.lastRet)
	if it.lastRet < it.cursor {
		it.cursor--
	}
	it.lastRet = -1
}
func (it *jFullListIter) Set(x any) {
	if it.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("set before next"))
	}
	it.lc.Set(it.lastRet, x)
}
func (it *jFullListIter) Add(x any) {
	if it.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	it.lc.AddAt(it.cursor, x)
	it.cursor++
	it.lastRet = -1
}

type jDescListIter struct {
	lc      *JListCore
	cursor  int32
	lastRet int32
}

func newDescendingListIter(lc *JListCore) Ext_java_util_Iterator {
	return &jDescListIter{lc: lc, cursor: lc.Size() - 1, lastRet: -1}
}

func (*jDescListIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jDescListIter) HasNext() bool                          { return it.cursor >= 0 }
func (it *jDescListIter) Next() any {
	if it.cursor < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.lc.Get(it.cursor)
	it.lastRet = it.cursor
	it.cursor--
	return v
}
func (it *jDescListIter) Remove() {
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.lc.RemoveAt(it.lastRet)
	it.lastRet = -1
}

type jDequeIter struct {
	dc      *JDequeCore
	off     int32
	lastOff int32
}

func newDequeIter(dc *JDequeCore) Ext_java_util_Iterator {
	if dc == nil || dc.size == 0 {
		return emptyMapIter_
	}
	return &jDequeIter{dc: dc, off: 0, lastOff: -1}
}

func (*jDequeIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jDequeIter) HasNext() bool                          { return it.off < it.dc.size }
func (it *jDequeIter) Next() any {
	if it.off >= it.dc.size {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.dc.elemAt((it.dc.head + it.off) & it.dc.mask)
	it.lastOff = it.off
	it.off++
	return v
}
func (it *jDequeIter) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.dc.deleteAtOffset(it.lastOff)
	it.off = it.lastOff
	it.lastOff = -1
}

type jDequeDescIter struct {
	dc      *JDequeCore
	off     int32
	lastOff int32
}

func newDequeDescIter(dc *JDequeCore) Ext_java_util_Iterator {
	if dc == nil || dc.size == 0 {
		return emptyMapIter_
	}
	return &jDequeDescIter{dc: dc, off: dc.size - 1, lastOff: -1}
}

func (*jDequeDescIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jDequeDescIter) HasNext() bool                          { return it.off >= 0 && it.off < it.dc.size }
func (it *jDequeDescIter) Next() any {
	if it.off < 0 || it.off >= it.dc.size {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.dc.elemAt((it.dc.head + it.off) & it.dc.mask)
	it.lastOff = it.off
	it.off--
	return v
}
func (it *jDequeDescIter) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.dc.deleteAtOffset(it.lastOff)
	it.lastOff = -1
}

type jEmptyMapIter struct{}

var emptyMapIter_ = &jEmptyMapIter{}

func (*jEmptyMapIter) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*jEmptyMapIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*jEmptyMapIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (*jEmptyMapIter) HasNext() bool { return false }
func (*jEmptyMapIter) Next() any {
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}
func (*jEmptyMapIter) Remove() {
	panic(javalang.NewIllegalStateException("remove before next"))
}

type jMapCursor struct {
	mc       *JMapCore
	mode     uint8
	idx      int
	curPos   int32
	lastSlot int32
	lastKey  any
	hasLast  bool
	keys     []any
	strKeys  []string
	vals     []any
}

func initMapCursor(mc *JMapCore, needVals bool) (jMapCursor, bool) {
	if mc == nil || mc.Live == 0 {
		return jMapCursor{}, false
	}
	switch mc.mode {
	case jMapModeNullOnly, jMapModeNode1:
		return jMapCursor{mc: mc, mode: mc.mode, idx: 0, lastSlot: -1}, true
	case jMapModeStr, jMapModeNode, jMapModePtr:
		start := int(mc.firstIdx)
		c := jMapCursor{mc: mc, mode: mc.mode, idx: start, lastSlot: -1}
		c.advanceSlice()
		if c.idx != start && !mc.immutable {
			mc.firstIdx = int32(c.idx)
		}
		return c, c.hasSliceNext()
	case jMapModeGeneral:
		g := mc.ext.(*jGenMapState)
		return jMapCursor{mc: mc, mode: jMapModeGeneral, curPos: g.head, lastSlot: -1}, g.head != 0
	case jMapModeSortedNat, jMapModeSortedCmp:
		var vals []any
		if needVals {
			vals = make([]any, 0, mc.Live)
		}
		if mc.ext != nil {
			st := mc.ext.(*jSortedMapState)
			if st.avlAllStr(st.root) {
				strKeys := make([]string, 0, mc.Live)
				strKeys, vals = st.avlAppendStrEntries(st.root, strKeys, vals, needVals)
				return jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, strKeys: strKeys, vals: vals}, len(strKeys) > 0
			}
		}
		keys := make([]any, 0, mc.Live)
		mc.ForEachEntry(func(k, v any) bool {
			keys = append(keys, k)
			if needVals {
				vals = append(vals, v)
			}
			return true
		})
		return jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, keys: keys, vals: vals}, len(keys) > 0
	default:
		return jMapCursor{}, false
	}
}

func initMapValCursor(mc *JMapCore) (jMapCursor, bool) {
	if mc == nil || mc.Live == 0 {
		return jMapCursor{}, false
	}
	if mc.mode == jMapModeSortedNat || mc.mode == jMapModeSortedCmp {
		vals := make([]any, 0, mc.Live)
		mc.ForEachValue(func(v any) bool {
			vals = append(vals, v)
			return true
		})
		return jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, vals: vals}, len(vals) > 0
	}
	return initMapCursor(mc, false)
}

func (c *jMapCursor) hasSliceNext() bool {
	switch c.mode {
	case jMapModeStr:
		return c.idx < len(c.mc.strEntries)
	case jMapModeNode:
		return c.idx < len(c.mc.ext.(*jNodeMapState).entries)
	case jMapModePtr:
		return c.idx < len(c.mc.ext.(*jPtrMapState).entries)
	default:
		return false
	}
}

func (c *jMapCursor) advanceSlice() {
	switch c.mode {
	case jMapModeStr:
		for c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
			c.idx++
		}
	case jMapModeNode:
		es := c.mc.ext.(*jNodeMapState).entries
		for c.idx < len(es) && !es[c.idx].active {
			c.idx++
		}
	case jMapModePtr:
		es := c.mc.ext.(*jPtrMapState).entries
		for c.idx < len(es) && es[c.idx].key == jTombstoneObj {
			c.idx++
		}
	}
}

func (c *jMapCursor) hasNext() bool {
	switch c.mode {
	case jMapModeNullOnly:
		return c.idx == 0 && c.mc.mode == jMapModeNullOnly
	case jMapModeNode1:
		return c.idx == 0 && c.mc.mode == jMapModeNode1
	case jMapModeGeneral:
		return c.curPos != 0
	case jMapModeSortedNat, jMapModeSortedCmp:
		return c.idx < len(c.keys) || c.idx < len(c.strKeys)
	default:
		c.advanceSlice()
		return c.hasSliceNext()
	}
}

func (c *jMapCursor) hasValNext() bool {
	switch c.mode {
	case jMapModeNullOnly:
		return c.idx == 0 && c.mc.mode == jMapModeNullOnly
	case jMapModeNode1:
		return c.idx == 0 && c.mc.mode == jMapModeNode1
	case jMapModeGeneral:
		return c.curPos != 0
	case jMapModeSortedNat, jMapModeSortedCmp:
		return c.idx < len(c.vals)
	default:
		c.advanceSlice()
		return c.hasSliceNext()
	}
}

func (c *jMapCursor) nextVal() any {
	switch c.mode {
	case jMapModeNullOnly:
		if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return c.mc.ext
	case jMapModeNode1:
		if c.idx > 0 || c.mc.mode != jMapModeNode1 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return c.mc.ext
	case jMapModeStr:
		c.advanceSlice()
		if c.idx >= len(c.mc.strEntries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := c.mc.strEntries[c.idx].val
		c.idx++
		c.advanceSlice()
		return v
	case jMapModeNode:
		ns := c.mc.ext.(*jNodeMapState)
		c.advanceSlice()
		if c.idx >= len(ns.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := ns.entries[c.idx].val
		c.idx++
		c.advanceSlice()
		return v
	case jMapModePtr:
		ps := c.mc.ext.(*jPtrMapState)
		c.advanceSlice()
		if c.idx >= len(ps.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := ps.entries[c.idx].val
		c.idx++
		c.advanceSlice()
		return v
	case jMapModeGeneral:
		if c.curPos == 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		g := c.mc.ext.(*jGenMapState)
		pos := c.curPos
		v := g.entries[pos-1].Val
		c.curPos = g.entries[pos-1].next()
		return v
	default:
		if c.idx >= len(c.vals) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := c.vals[c.idx]
		c.idx++
		return v
	}
}

func (c *jMapCursor) nextKey() any {
	switch c.mode {
	case jMapModeNullOnly:
		if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return nil
	case jMapModeNode1:
		if c.idx > 0 || c.mc.mode != jMapModeNode1 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return boxNodeById(uint32(c.mc.firstIdx))
	case jMapModeStr:
		c.advanceSlice()
		if c.idx >= len(c.mc.strEntries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ek := c.mc.strEntries[c.idx].key
		c.idx++
		c.advanceSlice()
		return BoxString(ek)
	case jMapModeNode:
		ns := c.mc.ext.(*jNodeMapState)
		c.advanceSlice()
		if c.idx >= len(ns.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ek := ns.entries[c.idx].key
		c.idx++
		c.advanceSlice()
		return boxNodeById(ek)
	case jMapModePtr:
		ps := c.mc.ext.(*jPtrMapState)
		c.advanceSlice()
		if c.idx >= len(ps.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ek := ps.entries[c.idx].key
		c.idx++
		c.advanceSlice()
		return ek
	case jMapModeGeneral:
		if c.curPos == 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		g := c.mc.ext.(*jGenMapState)
		pos := c.curPos
		k := g.entryKey(pos)
		c.curPos = g.entries[pos-1].next()
		return k
	default:
		if c.strKeys != nil {
			if c.idx >= len(c.strKeys) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			s := c.strKeys[c.idx]
			c.idx++
			return BoxString(s)
		}
		if c.idx >= len(c.keys) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		k := c.keys[c.idx]
		c.idx++
		return k
	}
}

func avlAppendStrKeys(s *jSortedMapState, cur int32, out []string) []string {
	if cur == 0 {
		return out
	}
	n := &s.tree[cur-1]
	if n.left != 0 {
		out = avlAppendStrKeys(s, n.left, out)
	}
	if n.isStr {
		out = append(out, n.strKey)
	} else {
		out = append(out, CastString(n.key))
	}
	if n.right != 0 {
		out = avlAppendStrKeys(s, n.right, out)
	}
	return out
}

func initMapKeyStrCursor(mc *JMapCore) (jMapCursor, bool) {
	if mc == nil || mc.Live == 0 {
		return jMapCursor{}, false
	}
	if mc.mode == jMapModeSortedNat || mc.mode == jMapModeSortedCmp {
		st := mc.ext.(*jSortedMapState)
		strKeys := avlAppendStrKeys(st, st.root, make([]string, 0, mc.Live))
		return jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, strKeys: strKeys}, len(strKeys) > 0
	}
	return initMapCursor(mc, false)
}

func (c *jMapCursor) nextKeyStr() string {
	switch c.mode {
	case jMapModeNullOnly:
		if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return NullString
	case jMapModeStr:
		c.advanceSlice()
		if c.idx >= len(c.mc.strEntries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ek := c.mc.strEntries[c.idx].key
		c.idx++
		c.advanceSlice()
		return ek
	case jMapModeGeneral:
		if c.curPos == 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		g := c.mc.ext.(*jGenMapState)
		pos := c.curPos
		e := &g.entries[pos-1]
		c.curPos = e.next()
		if e.kind() == jKeyNull {
			return NullString
		}
		if s, ok := e.Key.(string); ok {
			return s
		}
		return CastString(g.entryKey(pos))
	default:
		if c.strKeys != nil {
			if c.idx >= len(c.strKeys) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			s := c.strKeys[c.idx]
			c.idx++
			return s
		}
		return CastString(c.nextKey())
	}
}

func (c *jMapCursor) nextKeyNode() uint32 {
	switch c.mode {
	case jMapModeNullOnly:
		if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return 0
	case jMapModeNode1:
		if c.idx > 0 || c.mc.mode != jMapModeNode1 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return uint32(c.mc.firstIdx)
	case jMapModeNode:
		ns := c.mc.ext.(*jNodeMapState)
		c.advanceSlice()
		if c.idx >= len(ns.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		u := ns.entries[c.idx].key
		c.idx++
		c.advanceSlice()
		return u
	case jMapModeGeneral:
		if c.curPos == 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		g := c.mc.ext.(*jGenMapState)
		pos := c.curPos
		e := &g.entries[pos-1]
		c.curPos = e.next()
		if e.kind() == jKeyNull {
			return 0
		}
		if e.kind() == jKeyNode {
			return e.u32OrHash
		}
		return CastNode[uint32](g.entryKey(pos))
	default:
		return CastNode[uint32](c.nextKey())
	}
}

func (c *jMapCursor) nextKV() (k any, v any, slot int32) {
	switch c.mode {
	case jMapModeNullOnly:
		if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		k = nil
		v = c.mc.ext
		slot = 0
		c.idx = 1
		c.lastSlot = 0
		c.lastKey = nil
		c.hasLast = true
		return
	case jMapModeNode1:
		if c.idx > 0 || c.mc.mode != jMapModeNode1 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		k = boxNodeById(uint32(c.mc.firstIdx))
		v = c.mc.ext
		slot = 0
		c.idx = 1
		c.lastSlot = 0
		c.lastKey = k
		c.hasLast = true
		return
	case jMapModeStr:
		c.advanceSlice()
		if c.idx >= len(c.mc.strEntries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		e := c.mc.strEntries[c.idx]
		if e.key != NullString {
			k = e.key
		}
		v = e.val
		slot = int32(c.idx)
		c.idx++
		c.lastSlot = slot
		c.lastKey = k
		c.hasLast = true
		c.advanceSlice()
		return
	case jMapModeNode:
		ns := c.mc.ext.(*jNodeMapState)
		c.advanceSlice()
		if c.idx >= len(ns.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		e := ns.entries[c.idx]
		k = boxNodeById(e.key)
		v = e.val
		slot = int32(c.idx)
		c.idx++
		c.lastSlot = slot
		c.lastKey = k
		c.hasLast = true
		c.advanceSlice()
		return
	case jMapModePtr:
		ps := c.mc.ext.(*jPtrMapState)
		c.advanceSlice()
		if c.idx >= len(ps.entries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		e := ps.entries[c.idx]
		k = e.key
		v = e.val
		slot = int32(c.idx)
		c.idx++
		c.lastSlot = slot
		c.lastKey = k
		c.hasLast = true
		c.advanceSlice()
		return
	case jMapModeGeneral:
		if c.curPos == 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		g := c.mc.ext.(*jGenMapState)
		pos := c.curPos
		c.lastSlot = pos
		c.hasLast = true
		k = g.entryKey(pos)
		v = g.entries[pos-1].Val
		slot = pos
		c.curPos = g.entries[pos-1].next()
		return
	default:
		if c.strKeys != nil {
			if c.idx >= len(c.strKeys) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			k = BoxString(c.strKeys[c.idx])
			if c.vals != nil {
				v = c.vals[c.idx]
			}
			slot = -1
			c.idx++
			c.lastSlot = -1
			c.lastKey = k
			c.hasLast = true
			return
		}
		if c.idx >= len(c.keys) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		k = c.keys[c.idx]
		if c.vals != nil {
			v = c.vals[c.idx]
		}
		slot = -1
		c.idx++
		c.lastSlot = -1
		c.lastKey = k
		c.hasLast = true
		return
	}
}

func (c *jMapCursor) remove() {
	if !c.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	if c.mc.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	c.hasLast = false
	if c.mode == jMapModeGeneral && c.mc.mode == jMapModeGeneral && c.lastSlot > 0 {
		c.mc.removeGenPos(c.lastSlot)
		c.lastSlot = -1
		return
	}
	c.mc.Remove(c.lastKey)
	c.lastSlot = -1
}

type jMapKeyIter struct {
	cur jMapCursor
}

func newMapKeyIter(mc *JMapCore) Ext_java_util_Iterator {
	c, ok := initMapCursor(mc, false)
	if !ok {
		return emptyMapIter_
	}
	return &jMapKeyIter{cur: c}
}

func (*jMapKeyIter) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*jMapKeyIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*jMapKeyIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}

func (it *jMapKeyIter) HasNext() bool { return it.cur.hasNext() }
func (it *jMapKeyIter) Next() any {
	k, _, _ := it.cur.nextKV()
	return k
}
func (it *jMapKeyIter) Remove() { it.cur.remove() }

// Live SubList view
type jSubListView struct {
	Ext_java_util_AbstractListBase
	parent *JListCore
	from   int32
	to     int32
}

func newSubListView(parent *JListCore, from, to int32) Ext_java_util_List {
	CheckPositionIndexes(from, to, parent.Size())
	v := &jSubListView{parent: parent, from: from, to: to}
	v.BindExtSelf_(v)
	return v
}

func (v *jSubListView) Size() int32 { return v.to - v.from }
func (v *jSubListView) Get(i int32) any {
	CheckElementIndex(i, v.Size())
	return v.parent.Get(v.from + i)
}
func (v *jSubListView) Set(i int32, x any) any {
	CheckElementIndex(i, v.Size())
	return v.parent.Set(v.from+i, x)
}
func (v *jSubListView) Add_2(i int32, x any) {
	CheckPositionIndex(i, v.Size())
	v.parent.AddAt(v.from+i, x)
	v.to++
}
func (v *jSubListView) Remove_1_int32(i int32) any {
	CheckElementIndex(i, v.Size())
	old := v.parent.RemoveAt(v.from + i)
	v.to--
	return old
}
func (v *jSubListView) Clear() {
	if v.from == v.to {
		return
	}
	CheckPositionIndexes(v.from, v.to, v.parent.Size())
	s := v.parent.Items
	newLen := len(s) - int(v.to-v.from)
	copy(s[v.from:], s[v.to:])
	clear(s[newLen:])
	v.parent.Items = s[:newLen]
	v.to = v.from
}

// Backed Arrays.asList view
type jBackedArrayList struct {
	Ext_java_util_AbstractListBase
	arr []any
}

func newBackedArrayList(arr []any) Ext_java_util_List {
	CheckNotNullSlice(arr)
	v := &jBackedArrayList{arr: arr}
	v.BindExtSelf_(v)
	return v
}

func (v *jBackedArrayList) Size() int32 { return int32(len(v.arr)) }
func (v *jBackedArrayList) Get(i int32) any {
	CheckElementIndex(i, int32(len(v.arr)))
	return ToAny(v.arr[i])
}
func (v *jBackedArrayList) Set(i int32, x any) any {
	CheckElementIndex(i, int32(len(v.arr)))
	old := ToAny(v.arr[i])
	v.arr[i] = ToAny(x)
	return old
}

// Lists.reverse & Lists.transform live views
type jReversedListView struct {
	Ext_java_util_AbstractListBase
	base Ext_java_util_List
}

func newReversedListView(base Ext_java_util_List) Ext_java_util_List {
	CheckNotNull(base)
	v := &jReversedListView{base: base}
	v.BindExtSelf_(v)
	return v
}

func (v *jReversedListView) Size() int32 { return v.base.Size() }
func (v *jReversedListView) Get(i int32) any {
	n := v.base.Size()
	CheckElementIndex(i, n)
	return v.base.Get(n - 1 - i)
}
func (v *jReversedListView) Set(i int32, x any) any {
	n := v.base.Size()
	CheckElementIndex(i, n)
	return v.base.Set(n-1-i, x)
}
func (v *jReversedListView) Add_2(i int32, x any) {
	n := v.base.Size()
	CheckPositionIndex(i, n)
	v.base.Add_2(n-i, x)
}
func (v *jReversedListView) Remove_1_int32(i int32) any {
	n := v.base.Size()
	CheckElementIndex(i, n)
	return v.base.Remove_1_int32(n - 1 - i)
}

type jTransformedListView struct {
	Ext_java_util_AbstractListBase
	base Ext_java_util_List
	fn   Ext_com_google_common_base_Function
}

func newTransformedListView(base Ext_java_util_List, fn Ext_com_google_common_base_Function) Ext_java_util_List {
	CheckNotNull(base)
	CheckNotNull(fn)
	v := &jTransformedListView{base: base, fn: fn}
	v.BindExtSelf_(v)
	return v
}

func (v *jTransformedListView) Size() int32     { return v.base.Size() }
func (v *jTransformedListView) Get(i int32) any { return v.fn.Apply(v.base.Get(i)) }
func (v *jTransformedListView) Remove_1_int32(i int32) any {
	return v.fn.Apply(v.base.Remove_1_int32(i))
}
func (v *jTransformedListView) Clear() { v.base.Clear() }

// --- Map Operations & Live Views (keySet, values, entrySet) ---

func mapCoreRemove2(mc *JMapCore, k, v any) bool {
	if !mc.ContainsKey(k) || !javalang.ObjectsEquals(mc.Get(k), ToAny(v)) {
		return false
	}
	mc.Remove(k)
	return true
}

func mapCoreReplace2(mc *JMapCore, k, newVal any) any {
	if !mc.ContainsKey(k) {
		return nil
	}
	return mc.Put(k, newVal)
}

func mapCoreReplace3(mc *JMapCore, k, oldVal, newVal any) bool {
	if !mc.ContainsKey(k) || !javalang.ObjectsEquals(mc.Get(k), ToAny(oldVal)) {
		return false
	}
	mc.Put(k, newVal)
	return true
}

func mapCoreReplaceAll(mc *JMapCore, fn Ext_java_util_function_BiFunction) {
	CheckNotNull(fn)
	mc.TransformValues(func(k, v any) any {
		return fn.Apply(k, v)
	})
}

func mapCoreForEach(mc *JMapCore, action Ext_java_util_function_BiConsumer) {
	CheckNotNull(action)
	mc.ForEachEntry(func(k, v any) bool {
		action.Accept(k, v)
		return true
	})
}

func mapCoreComputeIfAbsent(mc *JMapCore, k any, fn Ext_java_util_function_Function) any {
	CheckNotNull(fn)
	if v := mc.Get(k); v != nil {
		return v
	}
	newVal := ToAny(fn.Apply(ToAny(k)))
	if newVal != nil {
		mc.Put(k, newVal)
	}
	return newVal
}

func mapCoreComputeIfAbsentStr(mc *JMapCore, k string, fn Ext_java_util_function_Function) any {
	CheckNotNull(fn)
	if k == NullString && mc.IsNaturalOrder() {
		ThrowNPE()
	}
	if v := mc.GetStr(k); v != nil {
		return v
	}
	newVal := ToAny(fn.Apply(BoxString(k)))
	if newVal != nil {
		mc.PutStr(k, newVal)
	}
	return newVal
}

func mapCoreComputeIfAbsentNode(mc *JMapCore, u uint32, fn Ext_java_util_function_Function) any {
	CheckNotNull(fn)
	if u == 0 && mc.IsNaturalOrder() {
		ThrowNPE()
	}
	if v := mc.GetNode(u); v != nil {
		return v
	}
	newVal := ToAny(fn.Apply(boxNodeById(u)))
	if newVal != nil {
		mc.PutNode(u, newVal)
	}
	return newVal
}

func mapCoreTryGetComputeStr(mc *JMapCore, k string) (any, bool) {
	if k == NullString && mc.IsNaturalOrder() {
		return nil, false
	}
	if v := mc.GetStr(k); v != nil {
		return v, true
	}
	return nil, false
}

func mapCoreTryGetComputeNode(mc *JMapCore, u uint32) (any, bool) {
	if u == 0 && mc.IsNaturalOrder() {
		return nil, false
	}
	if v := mc.GetNode(u); v != nil {
		return v, true
	}
	return nil, false
}

func mapCoreTryGetComputeAny(mc *JMapCore, k any) (any, bool) {
	if (k == nil || isJavaNilAny(k)) && mc.IsNaturalOrder() {
		return nil, false
	}
	if v := mc.Get(k); v != nil {
		return v, true
	}
	return nil, false
}

func mapCoreComputeIfPresent(mc *JMapCore, k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	oldVal := mc.Get(k)
	if oldVal == nil {
		return nil
	}
	newVal := ToAny(fn.Apply(ToAny(k), oldVal))
	if newVal == nil {
		mc.Remove(k)
		return nil
	}
	mc.Put(k, newVal)
	return newVal
}

func mapCoreCompute(mc *JMapCore, k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	oldVal := mc.Get(k)
	haveOld := mc.ContainsKey(k)
	newVal := ToAny(fn.Apply(ToAny(k), oldVal))
	if newVal == nil {
		if haveOld {
			mc.Remove(k)
		}
		return nil
	}
	mc.Put(k, newVal)
	return newVal
}

func mapCoreMerge(mc *JMapCore, k, val any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	v := CheckNotNull(ToAny(val))
	oldVal := mc.Get(k)
	var newVal any
	if oldVal == nil {
		newVal = v
	} else {
		newVal = ToAny(fn.Apply(oldVal, v))
	}
	if newVal == nil {
		mc.Remove(k)
	} else {
		mc.Put(k, newVal)
	}
	return newVal
}

func mapCoreEquals(mc *JMapCore, other any) bool {
	mp, ok := other.(Ext_java_util_Map)
	if !ok || mp == nil || isJavaNilAny(mp) {
		return false
	}
	if mc.Size() != mp.Size() {
		return false
	}
	if mc.Size() == 0 {
		return true
	}
	if fmp, ok2 := other.(jFastMapCoreProvider); ok2 {
		if omc := fmp.fastMapCore_(); omc != nil {
			if mc == omc {
				return true
			}
			if mc.mode == jMapModeStr && omc.mode == jMapModeStr {
				for i := int(mc.firstIdx); i < len(mc.strEntries); i++ {
					e := mc.strEntries[i]
					if e.key == jTombstoneStr {
						continue
					}
					v2, ok3 := omc.getStr(e.key, false)
					if !ok3 || !javalang.ObjectsEquals(e.val, v2) {
						return false
					}
				}
				return true
			}
			if mc.mode == jMapModeNode1 {
				v2, ok3 := omc.getNode(uint32(mc.firstIdx), false)
				return ok3 && javalang.ObjectsEquals(mc.ext, v2)
			}
			if mc.mode == jMapModeNode && omc.mode == jMapModeNode {
				ns := mc.ext.(*jNodeMapState)
				for i := int(mc.firstIdx); i < len(ns.entries); i++ {
					e := ns.entries[i]
					if !e.active {
						continue
					}
					v2, ok3 := omc.getNode(e.key, false)
					if !ok3 || !javalang.ObjectsEquals(e.val, v2) {
						return false
					}
				}
				return true
			}
			eq := true
			mc.ForEachEntry(func(k, v any) bool {
				v2, ok3 := omc.GetEntry(k)
				if !ok3 || !javalang.ObjectsEquals(v, v2) {
					eq = false
					return false
				}
				return true
			})
			return eq
		}
	}
	eq := true
	mc.ForEachEntry(func(k, v any) bool {
		if v == nil {
			if mp.Get(k) != nil || !mp.ContainsKey(k) {
				eq = false
				return false
			}
		} else if !javalang.ObjectsEquals(v, ToAny(mp.Get(k))) {
			eq = false
			return false
		}
		return true
	})
	return eq
}

func mapCoreHash(mc *JMapCore) int32 {
	h := int32(0)
	if mc != nil && mc.mode == jMapModeStr && !mc.Identity {
		for i := int(mc.firstIdx); i < len(mc.strEntries); i++ {
			e := mc.strEntries[i]
			if e.key == jTombstoneStr {
				continue
			}
			var kh int32
			if e.key != NullString {
				kh = javalang.StringHashCode(e.key)
			}
			h += kh ^ javalang.ObjectHashCode(e.val)
		}
		return h
	}
	mc.ForEachEntry(func(k, v any) bool {
		if mc.Identity {
			h += javalang.IdentityHashCode(k) ^ javalang.IdentityHashCode(v)
		} else {
			h += javalang.ObjectHashCode(k) ^ javalang.ObjectHashCode(v)
		}
		return true
	})
	return h
}

func mapCoreString(mc *JMapCore) string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	if mc != nil && mc.mode == jMapModeStr {
		for i := int(mc.firstIdx); i < len(mc.strEntries); i++ {
			e := mc.strEntries[i]
			if e.key == jTombstoneStr {
				continue
			}
			if !first {
				sb.WriteString(", ")
			}
			first = false
			if e.key == NullString {
				sb.WriteString("null")
			} else {
				sb.WriteString(e.key)
			}
			sb.WriteByte('=')
			sb.WriteString(StrVal(e.val))
		}
		sb.WriteByte('}')
		return sb.String()
	}
	mc.ForEachEntry(func(k, v any) bool {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(k))
		sb.WriteByte('=')
		sb.WriteString(StrVal(v))
		return true
	})
	sb.WriteByte('}')
	return sb.String()
}

type jSimpleEntry struct {
	k   any
	v   any
	mc  *JMapCore
	imm bool
}

func (*jSimpleEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
func (e *jSimpleEntry) GetKey() any                              { return e.k }
func (e *jSimpleEntry) GetValue() any                            { return e.v }
func (e *jSimpleEntry) SetValue(nv any) any {
	if e.imm {
		panic(javalang.NewUnsupportedOperationException("immutable entry"))
	}
	old := e.v
	e.v = ToAny(nv)
	if e.mc != nil {
		e.mc.Put(e.k, e.v)
	}
	return old
}
func (e *jSimpleEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
func (e *jSimpleEntry) HashCode() int32   { return mapEntryHash(e) }
func (e *jSimpleEntry) String() string    { return StrVal(e.k) + "=" + StrVal(e.v) }

func newSimpleImmEntry(k, v any) Ext_java_util_AbstractMap_SimpleImmutableEntry {
	e := &Ext_java_util_AbstractMap_SimpleImmutableEntryBase{}
	e.BindExtSelf_(e)
	e.Core_ = &jSimpleEntry{k: ToAny(k), v: ToAny(v), imm: true}
	return e
}

func mapEntryEquals(self Ext_java_util_Map_Entry, o any) bool {
	other, ok := o.(Ext_java_util_Map_Entry)
	if !ok || other == nil || javalang.IsNil(other) {
		return false
	}
	return javalang.ObjectsEquals(ToAny(self.GetKey()), ToAny(other.GetKey())) &&
		javalang.ObjectsEquals(ToAny(self.GetValue()), ToAny(other.GetValue()))
}

func mapEntryHash(self Ext_java_util_Map_Entry) int32 {
	return javalang.ObjectHashCode(ToAny(self.GetKey())) ^ javalang.ObjectHashCode(ToAny(self.GetValue()))
}

type jLiveMapKeySet struct {
	Ext_java_util_AbstractSetBase
	mc *JMapCore
}

func newLiveMapKeySet(mc *JMapCore) Ext_java_util_Set {
	s := &jLiveMapKeySet{mc: mc}
	s.BindExtSelf_(s)
	return s
}
func (s *jLiveMapKeySet) Size() int32             { return s.mc.Size() }
func (s *jLiveMapKeySet) Contains(k any) bool     { return s.mc.ContainsKey(k) }
func (s *jLiveMapKeySet) Iterator() any           { return newMapKeyIter(s.mc) }
func (s *jLiveMapKeySet) Remove_1_any(k any) bool { return s.mc.RemoveKey(k) }
func (s *jLiveMapKeySet) Clear()                  { s.mc.Clear() }
func (s *jLiveMapKeySet) SetContainsStr_(k string) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.mc.ContainsStr(k)
}
func (s *jLiveMapKeySet) SetAddStr_(_ string) bool {
	if s == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jLiveMapKeySet) SetRemoveStr_(k string) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.mc.RemoveKeyStr(k)
}
func (s *jLiveMapKeySet) SetContainsNode_(n uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.mc.ContainsNode(n)
}
func (s *jLiveMapKeySet) SetAddNode_(_ uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jLiveMapKeySet) SetRemoveNode_(n uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.mc.RemoveKeyNode(n)
}

type jLiveMapValues struct {
	Ext_java_util_AbstractCollectionBase
	mc *JMapCore
}

func newLiveMapValues(mc *JMapCore) Ext_java_util_Collection {
	v := &jLiveMapValues{mc: mc}
	v.BindExtSelf_(v)
	return v
}
func (v *jLiveMapValues) Size() int32         { return v.mc.Size() }
func (v *jLiveMapValues) Contains(x any) bool { return v.mc.ContainsValue(x) }
func (v *jLiveMapValues) Clear()              { v.mc.Clear() }
func (v *jLiveMapValues) Iterator() any {
	c, ok := initMapCursor(v.mc, true)
	if !ok {
		return emptyMapIter_
	}
	return &jMapValIter{cur: c}
}

type jMapValIter struct {
	cur jMapCursor
}

func (*jMapValIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jMapValIter) HasNext() bool                          { return it.cur.hasNext() }
func (it *jMapValIter) Next() any {
	_, v, _ := it.cur.nextKV()
	return v
}
func (it *jMapValIter) Remove() { it.cur.remove() }

type jLiveMapEntrySet struct {
	Ext_java_util_AbstractSetBase
	mc *JMapCore
}

func newLiveMapEntrySet(mc *JMapCore) Ext_java_util_Set {
	s := &jLiveMapEntrySet{mc: mc}
	s.BindExtSelf_(s)
	return s
}
func (s *jLiveMapEntrySet) Size() int32 { return s.mc.Size() }
func (s *jLiveMapEntrySet) Clear()      { s.mc.Clear() }
func (s *jLiveMapEntrySet) Contains(o any) bool {
	e, ok := o.(Ext_java_util_Map_Entry)
	if !ok || e == nil || javalang.IsNil(e) {
		return false
	}
	k := e.GetKey()
	return s.mc.ContainsKey(k) && javalang.ObjectsEquals(s.mc.Get(k), ToAny(e.GetValue()))
}
func (s *jLiveMapEntrySet) Remove_1_any(o any) bool {
	e, ok := o.(Ext_java_util_Map_Entry)
	if !ok || e == nil || javalang.IsNil(e) {
		return false
	}
	k := e.GetKey()
	if !s.mc.ContainsKey(k) || !javalang.ObjectsEquals(s.mc.Get(k), ToAny(e.GetValue())) {
		return false
	}
	s.mc.Remove(k)
	return true
}
func (s *jLiveMapEntrySet) Iterator() any {
	c, ok := initMapCursor(s.mc, true)
	if !ok {
		return emptyMapIter_
	}
	return &jMapEntryIter{cur: c}
}

type jLiveMapEntry struct {
	mc   *JMapCore
	mode uint8
	slot int32
	k    any
	v    any
}

func (*jLiveMapEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
func (e *jLiveMapEntry) GetKey() any                              { return e.k }
func (e *jLiveMapEntry) GetValue() any                            { return e.v }
func (e *jLiveMapEntry) SetValue(nv any) any {
	if e.mc.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable entry"))
	}
	nv = ToAny(nv)
	old := e.v
	e.v = nv
	switch e.mode {
	case jMapModeNullOnly:
		if e.mc.mode == jMapModeNullOnly && e.slot == 0 {
			e.mc.ext = nv
			return old
		}
	case jMapModeNode1:
		if e.mc.mode == jMapModeNode1 && e.slot == 0 {
			e.mc.ext = nv
			return old
		}
	case jMapModeStr:
		if e.mc.mode == jMapModeStr && e.slot >= 0 && int(e.slot) < len(e.mc.strEntries) && e.mc.strEntries[e.slot].key != jTombstoneStr {
			e.mc.strEntries[e.slot].val = nv
			return old
		}
	case jMapModeNode:
		if e.mc.mode == jMapModeNode && e.slot >= 0 {
			ns := e.mc.ext.(*jNodeMapState)
			if int(e.slot) < len(ns.entries) && ns.entries[e.slot].active {
				ns.entries[e.slot].val = nv
				return old
			}
		}
	case jMapModePtr:
		if e.mc.mode == jMapModePtr && e.slot >= 0 {
			ps := e.mc.ext.(*jPtrMapState)
			if int(e.slot) < len(ps.entries) && ps.entries[e.slot].key != jTombstoneObj {
				ps.entries[e.slot].val = nv
				return old
			}
		}
	case jMapModeGeneral:
		if e.mc.mode == jMapModeGeneral && e.slot > 0 {
			g := e.mc.ext.(*jGenMapState)
			if int(e.slot) <= len(g.entries) {
				g.entries[e.slot-1].Val = nv
				return old
			}
		}
	}
	e.mc.Put(e.k, nv)
	return old
}
func (e *jLiveMapEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
func (e *jLiveMapEntry) HashCode() int32   { return mapEntryHash(e) }
func (e *jLiveMapEntry) String() string    { return StrVal(e.k) + "=" + StrVal(e.v) }

type jMapEntryIter struct {
	cur jMapCursor
}

func (*jMapEntryIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jMapEntryIter) HasNext() bool                          { return it.cur.hasNext() }
func (it *jMapEntryIter) Next() any {
	k, v, slot := it.cur.nextKV()
	if slot < 0 {
		return &jSimpleEntry{k: k, v: v, mc: it.cur.mc, imm: it.cur.mc.immutable}
	}
	return &jLiveMapEntry{mc: it.cur.mc, mode: it.cur.mode, slot: slot, k: k, v: v}
}
func (it *jMapEntryIter) Remove() { it.cur.remove() }

type JMapValCursor struct {
	cur          jMapCursor
	em           *jEnumMapCore
	fallbackIter SimpleIter
}

func NewMapValCursor(m any) JMapValCursor {
	if m == nil {
		ThrowNPE()
	}
	if mp, ok := m.(jFastMapCoreProvider); ok {
		if mc := mp.fastMapCore_(); mc != nil {
			c, _ := initMapValCursor(mc)
			return JMapValCursor{cur: c}
		}
	}
	if ep, ok := m.(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			if em.size == 0 {
				return JMapValCursor{}
			}
			return JMapValCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}}
		}
	}
	if isJavaNilAny(m) {
		ThrowNPE()
	}
	return JMapValCursor{fallbackIter: Iterate(m.(Ext_java_util_Map).Values())}
}

func (c *JMapValCursor) HasNext() bool {
	if c.cur.mc != nil {
		return c.cur.hasValNext()
	}
	if c.em != nil {
		return c.cur.idx >= 0
	}
	if c.fallbackIter != nil {
		return c.fallbackIter.HasNext()
	}
	return false
}

func (c *JMapValCursor) Next() any {
	if c.cur.mc != nil {
		return c.cur.nextVal()
	}
	if c.em != nil {
		if c.cur.idx < 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ord := c.cur.idx
		c.cur.idx = c.em.nextSetBit(ord + 1)
		return c.em.vals[ord]
	}
	if c.fallbackIter != nil {
		return c.fallbackIter.Next()
	}
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}

type JMapKeyCursor struct {
	cur          jMapCursor
	em           *jEnumMapCore
	fallbackIter SimpleIter
}

func NewMapKeyCursor(m any) JMapKeyCursor {
	if m == nil {
		ThrowNPE()
	}
	if mp, ok := m.(jFastMapCoreProvider); ok {
		if mc := mp.fastMapCore_(); mc != nil {
			c, _ := initMapCursor(mc, false)
			return JMapKeyCursor{cur: c}
		}
	}
	if ep, ok := m.(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			if em.size == 0 {
				return JMapKeyCursor{}
			}
			return JMapKeyCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}}
		}
	}
	if isJavaNilAny(m) {
		ThrowNPE()
	}
	return JMapKeyCursor{fallbackIter: Iterate(m.(Ext_java_util_Map).KeySet_0())}
}

func NewMapKeyStrCursor(m any) JMapKeyCursor {
	if m == nil {
		ThrowNPE()
	}
	if mp, ok := m.(jFastMapCoreProvider); ok {
		if mc := mp.fastMapCore_(); mc != nil {
			c, _ := initMapKeyStrCursor(mc)
			return JMapKeyCursor{cur: c}
		}
	}
	return NewMapKeyCursor(m)
}

func (c *JMapKeyCursor) HasNext() bool {
	if c.cur.mc != nil {
		return c.cur.hasNext()
	}
	if c.em != nil {
		return c.cur.idx >= 0
	}
	if c.fallbackIter != nil {
		return c.fallbackIter.HasNext()
	}
	return false
}

func (c *JMapKeyCursor) Next() any {
	if c.cur.mc != nil {
		return c.cur.nextKey()
	}
	if c.em != nil {
		if c.cur.idx < 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ord := c.cur.idx
		c.cur.idx = c.em.nextSetBit(ord + 1)
		return c.em.universe[ord]
	}
	if c.fallbackIter != nil {
		return c.fallbackIter.Next()
	}
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}

func (c *JMapKeyCursor) NextStr() string {
	if c.cur.mc != nil {
		return c.cur.nextKeyStr()
	}
	return CastString(c.Next())
}

func (c *JMapKeyCursor) NextNode() uint32 {
	if c.cur.mc != nil {
		return c.cur.nextKeyNode()
	}
	return CastNode[uint32](c.Next())
}

type JMapEntryCursor struct {
	cur          jMapCursor
	trm          *jTableRowMap
	em           *jEnumMapCore
	fallbackIter SimpleIter
	curEntry     Ext_java_util_Map_Entry
	curVal       any
	curKeyAny    any
	curKeyStr    string
	curNode      uint32
	curSlot      int32
	curKind      uint8 // 0=any(curKeyAny), 1=str(curKeyStr), 2=node(curNode), 3=int32(curNode)
	imm          bool
	singleReady  bool
}

func NewMapEntryCursor(m any) JMapEntryCursor {
	if m == nil {
		ThrowNPE()
	}
	if mp, ok := m.(jFastMapCoreProvider); ok {
		if mc := mp.fastMapCore_(); mc != nil {
			c, _ := initMapCursor(mc, true)
			_, isImmBi := m.(*Ext_com_google_common_collect_ImmutableBiMapBase)
			return JMapEntryCursor{cur: c, curSlot: -1, imm: mc.immutable || isImmBi}
		}
	}
	if trm, ok := m.(*jTableRowMap); ok {
		if trm == nil {
			ThrowNPE()
		}
		e, _ := trm.resolveEntry()
		if e == nil {
			return JMapEntryCursor{}
		}
		if e.hasSingle {
			return JMapEntryCursor{trm: trm, singleReady: true, curSlot: -1}
		}
		c, _ := initMapCursor(e.multi, true)
		return JMapEntryCursor{cur: c, curSlot: -1}
	}
	if ep, ok := m.(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			if em.size == 0 {
				return JMapEntryCursor{}
			}
			return JMapEntryCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}, curSlot: -1}
		}
	}
	if isJavaNilAny(m) {
		ThrowNPE()
	}
	return JMapEntryCursor{fallbackIter: Iterate(m.(Ext_java_util_Map).EntrySet()), curSlot: -1}
}

func (c *JMapEntryCursor) NextEntry() bool {
	if c.singleReady {
		c.singleReady = false
		e, _ := c.trm.resolveEntry()
		if e == nil {
			return false
		}
		if e.hasSingle {
			c.curKind = 0
			c.curKeyAny = e.col0
			c.curVal = e.val0
			c.curSlot = 0
			return true
		}
		mc, _ := initMapCursor(e.multi, true)
		c.cur = mc
	}
	if c.cur.mc != nil {
		switch c.cur.mode {
		case jMapModeNullOnly:
			if c.cur.idx > 0 || c.cur.mc.mode != jMapModeNullOnly {
				return false
			}
			c.cur.idx = 1
			c.curSlot = 0
			c.curKind = 0
			c.curKeyAny = nil
			c.curVal = c.cur.mc.ext
			return true
		case jMapModeNode1:
			if c.cur.idx > 0 || c.cur.mc.mode != jMapModeNode1 {
				return false
			}
			c.cur.idx = 1
			c.curSlot = 0
			c.curKind = 2
			c.curNode = uint32(c.cur.mc.firstIdx)
			c.curVal = c.cur.mc.ext
			return true
		case jMapModeStr:
			c.cur.advanceSlice()
			if c.cur.idx >= len(c.cur.mc.strEntries) {
				return false
			}
			slot := c.cur.idx
			e := c.cur.mc.strEntries[slot]
			c.cur.idx = slot + 1
			c.cur.advanceSlice()
			c.curSlot = int32(slot)
			c.curKind = 1
			c.curKeyStr = e.key
			c.curVal = e.val
			return true
		case jMapModeNode:
			ns := c.cur.mc.ext.(*jNodeMapState)
			c.cur.advanceSlice()
			if c.cur.idx >= len(ns.entries) {
				return false
			}
			slot := c.cur.idx
			e := ns.entries[slot]
			c.cur.idx = slot + 1
			c.cur.advanceSlice()
			c.curSlot = int32(slot)
			c.curKind = 2
			c.curNode = e.key
			c.curVal = e.val
			return true
		case jMapModePtr:
			ps := c.cur.mc.ext.(*jPtrMapState)
			c.cur.advanceSlice()
			if c.cur.idx >= len(ps.entries) {
				return false
			}
			slot := c.cur.idx
			e := ps.entries[slot]
			c.cur.idx = slot + 1
			c.cur.advanceSlice()
			c.curSlot = int32(slot)
			c.curKind = 0
			c.curKeyAny = e.key
			c.curVal = e.val
			return true
		case jMapModeGeneral:
			if c.cur.curPos == 0 {
				return false
			}
			g := c.cur.mc.ext.(*jGenMapState)
			pos := c.cur.curPos
			ge := &g.entries[pos-1]
			c.cur.curPos = ge.next()
			c.curSlot = pos
			c.curVal = ge.Val
			switch ge.kind() {
			case jKeyNull:
				c.curKind = 0
				c.curKeyAny = nil
			case jKeyStr:
				c.curKind = 1
				c.curKeyStr = ge.Key.(string)
				c.curKeyAny = ge.Key
			case jKeyNode:
				c.curKind = 2
				c.curNode = ge.u32OrHash
			case jKeyInt:
				c.curKind = 3
				c.curNode = ge.u32OrHash
				c.curKeyAny = ge.Key
			default:
				c.curKind = 0
				c.curKeyAny = ge.Key
			}
			return true
		default:
			if c.cur.strKeys != nil {
				if c.cur.idx >= len(c.cur.strKeys) {
					return false
				}
				i := c.cur.idx
				c.cur.idx = i + 1
				c.curSlot = -1
				c.curKind = 1
				c.curKeyStr = c.cur.strKeys[i]
				c.curVal = c.cur.vals[i]
				return true
			}
			if c.cur.idx >= len(c.cur.keys) {
				return false
			}
			i := c.cur.idx
			c.cur.idx = i + 1
			c.curSlot = -1
			c.curKind = 0
			c.curKeyAny = c.cur.keys[i]
			c.curVal = c.cur.vals[i]
			return true
		}
	}
	if c.em != nil {
		if c.cur.idx < 0 {
			return false
		}
		ord := c.cur.idx
		c.cur.idx = c.em.nextSetBit(ord + 1)
		c.curSlot = int32(ord)
		c.curKind = 0
		c.curKeyAny = c.em.universe[ord]
		c.curVal = c.em.vals[ord]
		return true
	}
	if c.fallbackIter != nil {
		if !c.fallbackIter.HasNext() {
			return false
		}
		c.curEntry = Cast[Ext_java_util_Map_Entry](c.fallbackIter.Next())
		return true
	}
	return false
}

func (c *JMapEntryCursor) GetKey() any {
	if c.fallbackIter != nil {
		return c.curEntry.GetKey()
	}
	switch c.curKind {
	case 1:
		return BoxString(c.curKeyStr)
	case 2:
		return boxNodeById(c.curNode)
	case 3:
		if c.curKeyAny != nil {
			return c.curKeyAny
		}
		return ValueOfInt(int32(c.curNode))
	default:
		return c.curKeyAny
	}
}

func (c *JMapEntryCursor) KeyStr() string {
	if c.fallbackIter != nil {
		return CastString(c.curEntry.GetKey())
	}
	if c.curKind == 1 {
		return c.curKeyStr
	}
	if c.curKind == 0 && c.curKeyAny == nil {
		return NullString
	}
	return CastString(c.GetKey())
}

func (c *JMapEntryCursor) KeyNode() uint32 {
	if c.fallbackIter != nil {
		return CastNode[uint32](c.curEntry.GetKey())
	}
	if c.curKind == 2 {
		return c.curNode
	}
	if c.curKind == 0 && c.curKeyAny == nil {
		return 0
	}
	return CastNode[uint32](c.GetKey())
}

func (c *JMapEntryCursor) GetValue() any {
	if c.fallbackIter != nil {
		return c.curEntry.GetValue()
	}
	return c.curVal
}

func (c *JMapEntryCursor) SetValue(nv any) any {
	if c.fallbackIter != nil {
		return c.curEntry.SetValue(ToAny(nv))
	}
	if c.imm || (c.cur.mc != nil && c.cur.mc.immutable) {
		panic(javalang.NewUnsupportedOperationException("immutable entry"))
	}
	if c.trm != nil && c.cur.mc == nil {
		nv = CheckNotNull(ToAny(nv))
		old := c.trm.putImpl(c.curKeyAny, nv)
		c.curVal = nv
		return old
	}
	nv = ToAny(nv)
	old := c.curVal
	c.curVal = nv
	if c.em != nil {
		if c.curSlot >= 0 && c.em.hasBit(int(c.curSlot)) {
			c.em.vals[c.curSlot] = nv
			return old
		}
		c.em.Put(c.curKeyAny, nv)
		return old
	}
	switch c.cur.mode {
	case jMapModeNullOnly:
		if c.cur.mc.mode == jMapModeNullOnly && c.curSlot == 0 {
			c.cur.mc.ext = nv
			return old
		}
	case jMapModeNode1:
		if c.cur.mc.mode == jMapModeNode1 && c.curSlot == 0 {
			c.cur.mc.ext = nv
			return old
		}
	case jMapModeStr:
		if c.cur.mc.mode == jMapModeStr && c.curSlot >= 0 && int(c.curSlot) < len(c.cur.mc.strEntries) && c.cur.mc.strEntries[c.curSlot].key != jTombstoneStr {
			c.cur.mc.strEntries[c.curSlot].val = nv
			return old
		}
	case jMapModeNode:
		if c.cur.mc.mode == jMapModeNode && c.curSlot >= 0 {
			ns := c.cur.mc.ext.(*jNodeMapState)
			if int(c.curSlot) < len(ns.entries) && ns.entries[c.curSlot].active {
				ns.entries[c.curSlot].val = nv
				return old
			}
		}
	case jMapModePtr:
		if c.cur.mc.mode == jMapModePtr && c.curSlot >= 0 {
			ps := c.cur.mc.ext.(*jPtrMapState)
			if int(c.curSlot) < len(ps.entries) && ps.entries[c.curSlot].key != jTombstoneObj {
				ps.entries[c.curSlot].val = nv
				return old
			}
		}
	case jMapModeGeneral:
		if c.cur.mc.mode == jMapModeGeneral && c.curSlot > 0 {
			g := c.cur.mc.ext.(*jGenMapState)
			if int(c.curSlot) <= len(g.entries) && g.entries[c.curSlot-1].kind() != jKeyFree {
				g.entries[c.curSlot-1].Val = nv
				return old
			}
		}
	}
	if c.curKind == 1 && !c.cur.mc.Identity {
		c.cur.mc.PutStr(c.curKeyStr, nv)
		return old
	}
	c.cur.mc.Put(c.GetKey(), nv)
	return old
}

type inlineStepIterable interface {
	iterInit_() (pos int32, has bool)
	iterNext_(pos int32) (elem any, nextPos int32)
}

type JIterCursor struct {
	slice        []any
	strs         []string
	tv           jTypedViewCursorSource
	step         inlineStepIterable
	mc           *JMapCore
	sc           *JSetCore
	dc           *JDequeCore
	em           *jEnumMapCore
	fallbackIter SimpleIter
	idx          int
	curPos       int32
	stepPos      int32
	kind         uint8 // 0=empty, 1=slice, 2=mapVal, 3=mapKey, 4=jset, 5=deque, 6=enumVal, 7=enumKey, 8=fallback, 9=strSlice, 10=typedView, 11=inlineStep
	mapMode      uint8
}

func initIterMapCursor(mc *JMapCore, vals bool) JIterCursor {
	if mc == nil || mc.Live == 0 {
		return JIterCursor{}
	}
	k := uint8(3)
	if vals {
		k = 2
	}
	switch mc.mode {
	case jMapModeNullOnly, jMapModeNode1:
		return JIterCursor{mc: mc, kind: k, mapMode: mc.mode}
	case jMapModeStr, jMapModeNode, jMapModePtr:
		start := int(mc.firstIdx)
		c := JIterCursor{mc: mc, idx: start, kind: k, mapMode: mc.mode}
		c.advanceMapSlice()
		if c.idx != start && !mc.immutable {
			mc.firstIdx = int32(c.idx)
		}
		return c
	case jMapModeGeneral:
		g := mc.ext.(*jGenMapState)
		if g.head == 0 {
			return JIterCursor{}
		}
		return JIterCursor{mc: mc, curPos: g.head, kind: k, mapMode: jMapModeGeneral}
	case jMapModeSortedNat, jMapModeSortedCmp:
		var s []any
		if vals {
			s = mc.LiveValues()
		} else {
			s = mc.LiveKeys()
		}
		if len(s) == 0 {
			return JIterCursor{}
		}
		return JIterCursor{slice: s, kind: 1}
	default:
		return JIterCursor{}
	}
}

func NewIterCursor(coll any) JIterCursor {
	if coll == nil {
		ThrowNPE()
	}
	if ssi, ok := coll.(StrSliceIterable); ok {
		if nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			ThrowNPE()
		}
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(ss) == 0 {
				return JIterCursor{}
			}
			return JIterCursor{strs: ss, kind: 9}
		}
	}
	if tv, ok := coll.(jTypedViewCursorSource); ok {
		if nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			ThrowNPE()
		}
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return JIterCursor{}
			}
			return JIterCursor{tv: tv, kind: 10}
		}
	}
	if si, ok := coll.(SliceIterable); ok {
		if nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			ThrowNPE()
		}
		s := si.IterSlice_()
		if len(s) == 0 {
			return JIterCursor{}
		}
		return JIterCursor{slice: s, kind: 1}
	}
	if st, ok := coll.(inlineStepIterable); ok {
		pos, has := st.iterInit_()
		if !has {
			return JIterCursor{}
		}
		return JIterCursor{step: st, stepPos: pos, kind: 11}
	}
	if mv, ok := coll.(*jLiveMapValues); ok {
		if mv == nil {
			ThrowNPE()
		}
		return initIterMapCursor(mv.mc, true)
	}
	if rv, ok := coll.(jTableRowValues); ok {
		if rv.m == nil {
			ThrowNPE()
		}
		e, _ := rv.m.resolveEntry()
		if e == nil {
			return JIterCursor{}
		}
		return initIterMapCursor(e.multi, true)
	}
	if mk, ok := coll.(*jLiveMapKeySet); ok {
		if mk == nil {
			ThrowNPE()
		}
		return initIterMapCursor(mk.mc, false)
	}
	if ms, ok := coll.(*jMultimapLiveSet); ok {
		if ms == nil {
			ThrowNPE()
		}
		b := ms.getBucket()
		if b == nil {
			return JIterCursor{}
		}
		return initIterMapCursor(b.(*JMapCore), false)
	}
	if mmk, ok := coll.(*jMultimapLiveKeySet); ok {
		if mmk == nil {
			ThrowNPE()
		}
		return initIterMapCursor(mmk.mm.keys, false)
	}
	if ev, ok := coll.(*jLiveEnumMapValues); ok {
		if ev == nil {
			ThrowNPE()
		}
		if ev.em == nil || ev.em.size == 0 {
			return JIterCursor{}
		}
		return JIterCursor{em: ev.em, idx: ev.em.nextSetBit(0), kind: 6}
	}
	if ek, ok := coll.(*jLiveEnumMapKeySet); ok {
		if ek == nil {
			ThrowNPE()
		}
		if ek.em == nil || ek.em.size == 0 {
			return JIterCursor{}
		}
		return JIterCursor{em: ek.em, idx: ek.em.nextSetBit(0), kind: 7}
	}
	if jsp, ok := coll.(jFastJSetCoreProvider); ok {
		if sc := jsp.fastSetJCore_(); sc != nil {
			if sc.Size() == 0 {
				return JIterCursor{}
			}
			start := int(sc.firstIdx)
			c := JIterCursor{sc: sc, idx: start, kind: 4}
			c.advanceJSet()
			if c.idx != start && !sc.imm {
				sc.firstIdx = int32(c.idx)
			}
			return c
		}
	}
	if sp, ok := coll.(jFastSetCoreProvider); ok {
		if mc := sp.fastSetCore_(); mc != nil {
			return initIterMapCursor(mc, false)
		}
	}
	if dp, ok := coll.(jFastDequeCoreProvider); ok {
		if dc := dp.fastDequeCore_(); dc != nil {
			if dc.size == 0 {
				return JIterCursor{}
			}
			return JIterCursor{dc: dc, kind: 5}
		}
	}
	if isJavaNilAny(coll) {
		ThrowNPE()
	}
	it := Iterate(coll)
	if ssi, ok := it.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(ss) == 0 {
				return JIterCursor{}
			}
			return JIterCursor{strs: ss, kind: 9}
		}
	}
	if tv, ok := it.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return JIterCursor{}
			}
			return JIterCursor{tv: tv, kind: 10}
		}
	}
	return JIterCursor{fallbackIter: it, kind: 8}
}

func (c *JIterCursor) advanceMapSlice() {
	switch c.mapMode {
	case jMapModeStr:
		for c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
			c.idx++
		}
	case jMapModeNode:
		es := c.mc.ext.(*jNodeMapState).entries
		for c.idx < len(es) && !es[c.idx].active {
			c.idx++
		}
	case jMapModePtr:
		es := c.mc.ext.(*jPtrMapState).entries
		for c.idx < len(es) && es[c.idx].key == jTombstoneObj {
			c.idx++
		}
	}
}

func (c *JIterCursor) jsetLen() int {
	switch c.sc.mode {
	case jSetModeNode:
		return len(c.sc.nodes)
	case jSetModeStr:
		return len(c.sc.ext.(*jSetStrState).strs)
	case jSetModeAny:
		return len(c.sc.ext.(*jSetAnyState).items)
	default:
		return 0
	}
}

func (c *JIterCursor) advanceJSet() {
	switch c.sc.mode {
	case jSetModeNode:
		for c.idx < len(c.sc.nodes) && c.sc.nodes[c.idx] == jTombstoneNode {
			c.idx++
		}
	case jSetModeStr:
		strs := c.sc.ext.(*jSetStrState).strs
		for c.idx < len(strs) && strs[c.idx] == jTombstoneStr {
			c.idx++
		}
	case jSetModeAny:
		items := c.sc.ext.(*jSetAnyState).items
		for c.idx < len(items) && items[c.idx] == jTombstoneObj {
			c.idx++
		}
	}
}

func (c *JIterCursor) HasNext() bool {
	switch c.kind {
	case 0:
		return false
	case 1:
		return c.idx < len(c.slice)
	case 9:
		for c.idx < len(c.strs) && c.strs[c.idx] == jTombstoneStr {
			c.idx++
		}
		return c.idx < len(c.strs)
	case 10:
		n, _ := c.tv.typedViewLen_()
		return c.idx < n
	case 2, 3:
		switch c.mapMode {
		case jMapModeNullOnly:
			return c.idx == 0 && c.mc.mode == jMapModeNullOnly
		case jMapModeNode1:
			return c.idx == 0 && c.mc.mode == jMapModeNode1
		case jMapModeGeneral:
			return c.curPos != 0
		case jMapModeStr:
			c.advanceMapSlice()
			return c.idx < len(c.mc.strEntries)
		case jMapModeNode:
			c.advanceMapSlice()
			return c.idx < len(c.mc.ext.(*jNodeMapState).entries)
		case jMapModePtr:
			c.advanceMapSlice()
			return c.idx < len(c.mc.ext.(*jPtrMapState).entries)
		default:
			return false
		}
	case 4:
		c.advanceJSet()
		return c.idx < c.jsetLen()
	case 5:
		return int32(c.idx) < c.dc.size
	case 6, 7:
		return c.idx >= 0
	case 8:
		return c.fallbackIter.HasNext()
	case 11:
		return c.stepPos >= 0
	default:
		return false
	}
}

func (c *JIterCursor) Next() any {
	switch c.kind {
	case 1:
		if c.idx >= len(c.slice) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := c.slice[c.idx]
		c.idx++
		return v
	case 2:
		switch c.mapMode {
		case jMapModeNullOnly:
			if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			c.idx = 1
			return c.mc.ext
		case jMapModeNode1:
			if c.idx > 0 || c.mc.mode != jMapModeNode1 {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			c.idx = 1
			return c.mc.ext
		case jMapModeStr:
			c.advanceMapSlice()
			if c.idx >= len(c.mc.strEntries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			v := c.mc.strEntries[c.idx].val
			c.idx++
			return v
		case jMapModeNode:
			ns := c.mc.ext.(*jNodeMapState)
			c.advanceMapSlice()
			if c.idx >= len(ns.entries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			v := ns.entries[c.idx].val
			c.idx++
			return v
		case jMapModePtr:
			ps := c.mc.ext.(*jPtrMapState)
			c.advanceMapSlice()
			if c.idx >= len(ps.entries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			v := ps.entries[c.idx].val
			c.idx++
			return v
		case jMapModeGeneral:
			if c.curPos == 0 {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			g := c.mc.ext.(*jGenMapState)
			pos := c.curPos
			v := g.entries[pos-1].Val
			c.curPos = g.entries[pos-1].next()
			return v
		}
	case 3:
		switch c.mapMode {
		case jMapModeNullOnly:
			if c.idx > 0 || c.mc.mode != jMapModeNullOnly {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			c.idx = 1
			return nil
		case jMapModeNode1:
			if c.idx > 0 || c.mc.mode != jMapModeNode1 {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			c.idx = 1
			return boxNodeById(uint32(c.mc.firstIdx))
		case jMapModeStr:
			c.advanceMapSlice()
			if c.idx >= len(c.mc.strEntries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			ek := c.mc.strEntries[c.idx].key
			c.idx++
			return BoxString(ek)
		case jMapModeNode:
			ns := c.mc.ext.(*jNodeMapState)
			c.advanceMapSlice()
			if c.idx >= len(ns.entries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			ek := ns.entries[c.idx].key
			c.idx++
			return boxNodeById(ek)
		case jMapModePtr:
			ps := c.mc.ext.(*jPtrMapState)
			c.advanceMapSlice()
			if c.idx >= len(ps.entries) {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			ek := ps.entries[c.idx].key
			c.idx++
			return ek
		case jMapModeGeneral:
			if c.curPos == 0 {
				panic(javalang.NewNoSuchElementException("iterator exhausted"))
			}
			g := c.mc.ext.(*jGenMapState)
			pos := c.curPos
			k := g.entryKey(pos)
			c.curPos = g.entries[pos-1].next()
			return k
		}
	case 4:
		c.advanceJSet()
		if c.idx >= c.jsetLen() {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		var v any
		switch c.sc.mode {
		case jSetModeNode:
			v = boxNodeById(c.sc.nodes[c.idx])
		case jSetModeStr:
			v = BoxString(c.sc.ext.(*jSetStrState).strs[c.idx])
		case jSetModeAny:
			v = c.sc.ext.(*jSetAnyState).items[c.idx]
		}
		c.idx++
		return v
	case 5:
		if int32(c.idx) >= c.dc.size {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := c.dc.elemAt((c.dc.head + int32(c.idx)) & c.dc.mask)
		c.idx++
		return v
	case 6:
		if c.idx < 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ord := c.idx
		c.idx = c.em.nextSetBit(ord + 1)
		return c.em.vals[ord]
	case 7:
		if c.idx < 0 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		ord := c.idx
		c.idx = c.em.nextSetBit(ord + 1)
		return c.em.universe[ord]
	case 8:
		return c.fallbackIter.Next()
	case 9:
		for c.idx < len(c.strs) && c.strs[c.idx] == jTombstoneStr {
			c.idx++
		}
		if c.idx >= len(c.strs) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		s := c.strs[c.idx]
		c.idx++
		return BoxString(s)
	case 10:
		n, _ := c.tv.typedViewLen_()
		if c.idx >= n {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		v := c.tv.typedViewAt_(c.idx)
		c.idx++
		return v
	case 11:
		v, np := c.step.iterNext_(c.stepPos)
		c.stepPos = np
		return v
	}
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}

func (c *JIterCursor) NextStr() string {
	if c.kind == 9 {
		for c.idx < len(c.strs) && c.strs[c.idx] == jTombstoneStr {
			c.idx++
		}
		if c.idx >= len(c.strs) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		s := c.strs[c.idx]
		c.idx++
		return s
	}
	if c.kind == 4 && c.sc.mode == jSetModeStr {
		strs := c.sc.ext.(*jSetStrState).strs
		for c.idx < len(strs) && strs[c.idx] == jTombstoneStr {
			c.idx++
		}
		if c.idx >= len(strs) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		s := strs[c.idx]
		c.idx++
		return s
	}
	if c.kind == 3 && c.mapMode == jMapModeStr {
		for c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
			c.idx++
		}
		if c.idx >= len(c.mc.strEntries) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		s := c.mc.strEntries[c.idx].key
		c.idx++
		return s
	}
	return CastString(c.Next())
}

func (c *JIterCursor) NextNode() uint32 {
	if c.kind == 10 {
		n, _ := c.tv.typedViewLen_()
		if c.idx >= n {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		if u, ok := c.tv.typedViewNodeAt_(c.idx); ok {
			c.idx++
			return u
		}
	}
	if c.kind == 4 && c.sc.mode == jSetModeNode {
		nodes := c.sc.nodes
		for c.idx < len(nodes) && nodes[c.idx] == jTombstoneNode {
			c.idx++
		}
		if c.idx >= len(nodes) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		u := nodes[c.idx]
		c.idx++
		return u
	}
	if c.kind == 3 && c.mapMode == jMapModeNode1 {
		if c.idx > 0 || c.mc.mode != jMapModeNode1 {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		c.idx = 1
		return uint32(c.mc.firstIdx)
	}
	if c.kind == 3 && c.mapMode == jMapModeNode {
		es := c.mc.ext.(*jNodeMapState).entries
		for c.idx < len(es) && !es[c.idx].active {
			c.idx++
		}
		if c.idx >= len(es) {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		u := es[c.idx].key
		c.idx++
		return u
	}
	if c.kind == 5 && c.dc.nodes != nil {
		if int32(c.idx) >= c.dc.size {
			panic(javalang.NewNoSuchElementException("iterator exhausted"))
		}
		u := c.dc.nodes[(c.dc.head+int32(c.idx))&c.dc.mask]
		c.idx++
		return u
	}
	return CastNode[uint32](c.Next())
}

func newJavaSet2(e0, e1 any) Ext_com_google_common_collect_ImmutableSet {
	if s0, ok0 := e0.(string); ok0 {
		if s1, ok1 := e1.(string); ok1 {
			if s0 == NullString || s1 == NullString {
				ThrowNPE()
			}
			if s0 == s1 {
				panic(javalang.NewIllegalArgumentException("duplicate element"))
			}
			return ImmSetOfStr2(s0, s1)
		}
	}
	e0 = CheckNotNull(ToAny(e0))
	e1 = CheckNotNull(ToAny(e1))
	if elemJavaEquals(e0, e1, false) {
		panic(javalang.NewIllegalArgumentException("duplicate element"))
	}
	c := &immAnySet2_{}
	c.buf[0] = e0
	c.buf[1] = e1
	c.set.elems_ = c.buf[:2:2]
	c.set.u8Tab_ = &emptyU8Tab_
	c.set.BindExtSelf_(&c.set)
	tryInitImmEnumMulti(&c.set, c.set.elems_)
	return &c.set
}

func immSetFromCheckedSlice(dupThrows bool, items []any) Ext_com_google_common_collect_ImmutableSet {
	if len(items) == 0 {
		return emptyImmutableSet_
	}
	if len(items) == 1 {
		return newImmSet1(items[0])
	}
	if len(items) <= 8 {
		var buf [8]any
		n := 0
		allStr := true
		for _, raw := range items {
			v := CheckNotNull(ToAny(raw))
			if _, ok := v.(string); !ok {
				allStr = false
			}
			dup := false
			for j := 0; j < n; j++ {
				if elemJavaEquals(buf[j], v, false) {
					if dupThrows {
						panic(javalang.NewIllegalArgumentException("duplicate element"))
					}
					dup = true
					break
				}
			}
			if !dup {
				buf[n] = v
				n++
			}
		}
		if n == 1 {
			return newImmSet1(buf[0])
		}
		if allStr {
			var sbuf [8]string
			for i := 0; i < n; i++ {
				sbuf[i] = buf[i].(string)
			}
			return finishImmStrSetSmall(sbuf[:n])
		}
		return finishImmSetSmall(buf[:n])
	}
	sc := &JSetCore{}
	sc.InitExpectedSize(int32(len(items)))
	for _, raw := range items {
		v := CheckNotNull(ToAny(raw))
		if !sc.Add(v) && dupThrows {
			panic(javalang.NewIllegalArgumentException("duplicate element"))
		}
	}
	return wrapImmutableSetFromSetCoreOwn(sc)
}

func wrapImmutableMapEntrySet(mc *JMapCore) Ext_com_google_common_collect_ImmutableSet {
	if mc == nil || mc.Size() == 0 {
		return emptyImmutableSet_
	}
	if mc.Size() == 1 {
		var e0 any
		mc.ForEachEntry(func(k, v any) bool {
			e0 = &jSimpleEntry{k: k, v: v, imm: true}
			return false
		})
		return newImmSet1(e0)
	}
	elems := make([]any, 0, mc.Size())
	mc.ForEachEntry(func(k, v any) bool {
		elems = append(elems, &jSimpleEntry{k: k, v: v, imm: true})
		return true
	})
	return finishImmSetMulti(elems, nil)
}

var emptyImmutableSortedSet_ Ext_com_google_common_collect_ImmutableSortedSet
var emptyImmutableSortedMap_ *Ext_com_google_common_collect_ImmutableSortedMap

func init() {
	{
		mc := NewSortedJMapCore(nil)
		mc.sealImmutable()
		obj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: mc}
		obj.elems_ = emptyAnySlice_
		obj.BindExtSelf_(obj)
		emptyImmutableSortedSet_ = obj
	}
	{
		obj := &Ext_com_google_common_collect_ImmutableSortedMap{}
		obj.BindExtSelf_(obj)
		obj.inlineCore.mode = jMapModeSortedNat
		obj.inlineCore.sealImmutable()
		emptyImmutableSortedMap_ = obj
	}
}

func wrapImmutableSortedSet(mc *JMapCore) Ext_com_google_common_collect_ImmutableSortedSet {
	if mc == nil || (mc.Size() == 0 && mc.IsNaturalOrder()) {
		return emptyImmutableSortedSet_
	}
	mc.sealImmutable()
	obj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: mc}
	obj.BindExtSelf_(obj)
	populateImmSortedSetElems(obj, mc)
	return obj
}

func wrapImmutableSortedMap(mc *JMapCore) *Ext_com_google_common_collect_ImmutableSortedMap {
	if mc == nil || (mc.Size() == 0 && mc.IsNaturalOrder()) {
		return emptyImmutableSortedMap_
	}
	mc.sealImmutable()
	obj := &Ext_com_google_common_collect_ImmutableSortedMap{}
	obj.BindExtSelf_(obj)
	obj.Core_ = mc
	return obj
}

func immSetFromStrSlice(dupThrows bool, src []string) Ext_com_google_common_collect_ImmutableSet {
	if len(src) == 0 {
		return emptyImmutableSet_
	}
	if len(src) == 1 {
		return ImmSetOfStr1(src[0])
	}
	if len(src) <= 8 {
		var buf [8]string
		n := 0
		for _, s := range src {
			if s == NullString {
				ThrowNPE()
			}
			dup := false
			for j := 0; j < n; j++ {
				if buf[j] == s {
					if dupThrows {
						panic(javalang.NewIllegalArgumentException("duplicate element"))
					}
					dup = true
					break
				}
			}
			if !dup {
				buf[n] = s
				n++
			}
		}
		if n == 1 {
			return ImmSetOfStr1(buf[0])
		}
		return finishImmStrSetSmall(buf[:n])
	}
	sc := &JSetCore{}
	sc.InitExpectedSize(int32(len(src)))
	for _, s := range src {
		if s == NullString {
			ThrowNPE()
		}
		if !sc.AddStr(s) && dupThrows {
			panic(javalang.NewIllegalArgumentException("duplicate element"))
		}
	}
	return wrapImmutableSetFromSetCoreOwn(sc)
}

func immListCopyOfSlice(src []any) Ext_com_google_common_collect_ImmutableList {
	CheckNotNullSlice(src)
	switch len(src) {
	case 0:
		return emptyImmutableList_
	case 1:
		return newImmList1(src[0])
	case 2:
		if s0, ok0 := src[0].(string); ok0 {
			if s1, ok1 := src[1].(string); ok1 {
				return ImmListOfStr2(s0, s1)
			}
		}
		return newImmList2(src[0], src[1])
	case 3:
		if s0, ok0 := src[0].(string); ok0 {
			if s1, ok1 := src[1].(string); ok1 {
				if s2, ok2 := src[2].(string); ok2 {
					return ImmListOfStr3(s0, s1, s2)
				}
			}
		}
		return newImmList3(src[0], src[1], src[2])
	case 4:
		if s0, ok0 := src[0].(string); ok0 {
			if s1, ok1 := src[1].(string); ok1 {
				if s2, ok2 := src[2].(string); ok2 {
					if s3, ok3 := src[3].(string); ok3 {
						return ImmListOfStr4(s0, s1, s2, s3)
					}
				}
			}
		}
		return newImmList4(src[0], src[1], src[2], src[3])
	case 5:
		return newImmList5(src[0], src[1], src[2], src[3], src[4])
	case 6:
		return newImmList6(src[0], src[1], src[2], src[3], src[4], src[5])
	case 7:
		return newImmList7(src[0], src[1], src[2], src[3], src[4], src[5], src[6])
	case 8:
		return newImmList8(src[0], src[1], src[2], src[3], src[4], src[5], src[6], src[7])
	default:
		if _, ok0 := src[0].(string); ok0 {
			allStr := true
			for _, v := range src {
				s, ok := v.(string)
				if !ok {
					allStr = false
					break
				}
				if s == NullString {
					ThrowNPE()
				}
			}
			if allStr {
				strs := make([]string, len(src))
				for i, v := range src {
					strs[i] = v.(string)
				}
				return &jImmStrList{strs: strs}
			}
		}
		items := make([]any, len(src))
		for i, v := range src {
			items[i] = CheckNotNull(ToAny(v))
		}
		return wrapImmutableListOwn(items)
	}
}

func immListCopyOfIterable(src any) Ext_com_google_common_collect_ImmutableList {
	CheckNotNull(src)
	if il, ok := src.(Ext_com_google_common_collect_ImmutableList); ok && !isJavaNilAny(il) {
		return il
	}
	if ssi, ok := src.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			return wrapImmStrSliceCopy(ss)
		}
	}
	if si, ok := src.(SliceIterable); ok {
		s := si.IterSlice_()
		if len(s) == 0 {
			return emptyImmutableList_
		}
		return immListCopyOfSlice(s)
	}
	if c, ok := src.(Ext_java_util_Collection); ok && c.IsEmpty() {
		return emptyImmutableList_
	}
	var b Ext_com_google_common_collect_ImmutableList_Builder
	b.AddAllVoid_(src)
	return b.BuildImmList_()
}

func immSetCopyOfSlice(src []any) Ext_com_google_common_collect_ImmutableSet {
	CheckNotNullSlice(src)
	return immSetFromCheckedSlice(false, src)
}

func immSetCopyOfIterable(src any) Ext_com_google_common_collect_ImmutableSet {
	CheckNotNull(src)
	if is, ok := src.(Ext_com_google_common_collect_ImmutableSet); ok && !isJavaNilAny(is) {
		return is
	}
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			return immSetFromEnumSet(es)
		}
	}
	if ssi, ok := src.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			return immSetFromStrSlice(false, ss)
		}
	}
	if si, ok := src.(SliceIterable); ok {
		return immSetFromCheckedSlice(false, si.IterSlice_())
	}
	if jsp, ok := src.(jFastJSetCoreProvider); ok {
		if osc := jsp.fastSetJCore_(); osc != nil && !osc.identity {
			return immSetFromSetCoreCopyChecked(osc)
		}
	}
	if sp, ok := src.(jFastSetCoreProvider); ok {
		if smc := sp.fastSetCore_(); smc != nil && !smc.Identity {
			if smc.Size() == 0 {
				return emptyImmutableSet_
			}
			if smc.mode == jMapModeStr {
				n := int(smc.Size())
				if n == 1 {
					for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
						s := smc.strEntries[i].key
						if s == jTombstoneStr {
							continue
						}
						if s == NullString {
							ThrowNPE()
						}
						return ImmSetOfStr1(s)
					}
				}
				if n <= 8 {
					var sbuf [8]string
					w := 0
					for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
						s := smc.strEntries[i].key
						if s == jTombstoneStr {
							continue
						}
						if s == NullString {
							ThrowNPE()
						}
						sbuf[w] = s
						w++
					}
					return finishImmStrSetSmall(sbuf[:w])
				}
				strs := make([]string, 0, n)
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					s := smc.strEntries[i].key
					if s == jTombstoneStr {
						continue
					}
					if s == NullString {
						ThrowNPE()
					}
					strs = append(strs, s)
				}
				return finishImmStrSetOwn(strs, nil)
			}
			if smc.Size() == 1 {
				return newImmSet1(smc.FirstKey())
			}
			return immSetFromCheckedSlice(false, smc.LiveKeys())
		}
	}
	if c, ok := src.(Ext_java_util_Collection); ok && c.IsEmpty() {
		return emptyImmutableSet_
	}
	sc := &JSetCore{}
	setCoreAddAllNonNull(sc, src)
	return wrapImmutableSetFromSetCoreOwn(sc)
}

func immSetFromSetCoreCopyChecked(src *JSetCore) Ext_com_google_common_collect_ImmutableSet {
	if src == nil || src.Size() == 0 {
		return emptyImmutableSet_
	}
	n := int(src.Size())
	if src.mode == jSetModeStr {
		st := src.ext.(*jSetStrState)
		if n == 1 {
			for i := int(src.firstIdx); i < len(st.strs); i++ {
				s := st.strs[i]
				if s == jTombstoneStr {
					continue
				}
				if s == NullString {
					ThrowNPE()
				}
				return ImmSetOfStr1(s)
			}
			return emptyImmutableSet_
		}
		if n <= 8 {
			var sbuf [8]string
			w := 0
			for i := int(src.firstIdx); i < len(st.strs); i++ {
				s := st.strs[i]
				if s == jTombstoneStr {
					continue
				}
				if s == NullString {
					ThrowNPE()
				}
				sbuf[w] = s
				w++
			}
			if w == 1 {
				return ImmSetOfStr1(sbuf[0])
			}
			return finishImmStrSetSmall(sbuf[:w])
		}
		strs := make([]string, 0, n)
		for i := int(src.firstIdx); i < len(st.strs); i++ {
			s := st.strs[i]
			if s == jTombstoneStr {
				continue
			}
			if s == NullString {
				ThrowNPE()
			}
			strs = append(strs, s)
		}
		if len(strs) == 1 {
			return ImmSetOfStr1(strs[0])
		}
		return finishImmStrSetOwn(strs, nil)
	}
	if n == 1 {
		return newImmSet1(src.FirstKey())
	}
	if n <= 8 {
		var buf [8]any
		w := 0
		src.ForEach(func(k any) bool {
			buf[w] = CheckNotNull(k)
			w++
			return true
		})
		if w == 1 {
			return newImmSet1(buf[0])
		}
		return finishImmSetSmall(buf[:w])
	}
	elems := make([]any, 0, n)
	src.ForEach(func(k any) bool {
		elems = append(elems, CheckNotNull(k))
		return true
	})
	return finishImmSetMulti(elems, nil)
}

func immSetFromSetCoreCopy(src *JSetCore) Ext_com_google_common_collect_ImmutableSet {
	return immSetFromSetCoreCopyChecked(src)
}

func immSortedSetFromCoreCopy(src *JMapCore) Ext_com_google_common_collect_ImmutableSortedSet {
	if src == nil || (src.Size() == 0 && src.IsNaturalOrder()) {
		return emptyImmutableSortedSet_
	}
	dst := cloneJMapCore(src)
	return wrapImmutableSortedSet(dst)
}

func immMapCopyOf(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableMap {
	CheckNotNull(src)
	if im, ok := any(src).(Ext_com_google_common_collect_ImmutableMap); ok && !isJavaNilAny(im) {
		return im
	}
	if src.IsEmpty() {
		return emptyImmutableMap_
	}
	obj := newImmutableMapShell()
	mc := &obj.inlineCore
	copyExtMapIntoNonNull(mc, src)
	mc.sealImmutable()
	return obj
}

func immMapCopyOfEntries(entries any) Ext_com_google_common_collect_ImmutableMap {
	CheckNotNull(entries)
	obj := newImmutableMapShell()
	mc := &obj.inlineCore
	for it := Iterate(entries); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](CheckNotNull(it.Next()))
		k := CheckNotNull(ToAny(e.GetKey()))
		v := CheckNotNull(ToAny(e.GetValue()))
		if mc.ContainsKey(k) {
			panic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
		}
		mc.Put(k, v)
	}
	if mc.IsEmpty() {
		return emptyImmutableMap_
	}
	mc.sealImmutable()
	return obj
}

func immSortedMapCopyOf(src Ext_java_util_Map) *Ext_com_google_common_collect_ImmutableSortedMap {
	CheckNotNull(src)
	if ism, ok := any(src).(*Ext_com_google_common_collect_ImmutableSortedMap); ok && ism != nil && ism.mc().IsNaturalOrder() {
		return ism
	}
	if src.IsEmpty() {
		return emptyImmutableSortedMap_
	}
	obj := &Ext_com_google_common_collect_ImmutableSortedMap{}
	obj.BindExtSelf_(obj)
	mc := initSortedCore(&obj.inlineCore, nil)
	copyExtMapIntoNonNull(mc, src)
	mc.sealImmutable()
	return obj
}

type jImmMapBuilderCore struct {
	mc     *JMapCore
	hasDup bool
	dupKey any
	sorted bool
}

func immMapBuilderPut(b *jImmMapBuilderCore, k, v any) {
	k = CheckNotNull(ToAny(k))
	v = CheckNotNull(ToAny(v))
	if b.mc.ContainsKey(k) {
		b.hasDup = true
		b.dupKey = k
	}
	b.mc.Put(k, v)
}

func immMapBuilderPutStr(b *jImmMapBuilderCore, k string, v any) {
	if k == NullString {
		ThrowNPE()
	}
	v = CheckNotNull(ToAny(v))
	if b.mc.ContainsStr(k) {
		b.hasDup = true
		b.dupKey = k
	}
	b.mc.PutStr(k, v)
}

func avlImmBuilderPutAll(s *jSortedMapState, cur int32, b *jImmMapBuilderCore) {
	if cur == 0 {
		return
	}
	n := &s.tree[cur-1]
	if n.left != 0 {
		avlImmBuilderPutAll(s, n.left, b)
	}
	if n.isStr {
		immMapBuilderPutStr(b, n.strKey, n.val)
	} else {
		immMapBuilderPut(b, n.key, n.val)
	}
	if n.right != 0 {
		avlImmBuilderPutAll(s, n.right, b)
	}
}

func immMapBuilderPutAll(b *jImmMapBuilderCore, m Ext_java_util_Map) {
	CheckNotNull(m)
	if mp, ok := any(m).(jFastMapCoreProvider); ok {
		if smc := mp.fastMapCore_(); smc != nil {
			if smc.Size() == 0 {
				return
			}
			b.mc.InitExpectedSize(b.mc.Size() + smc.Size())
			if smc.mode == jMapModeStr && !smc.Identity && !b.mc.Identity {
				for i := int(smc.firstIdx); i < len(smc.strEntries); i++ {
					e := smc.strEntries[i]
					if e.key != jTombstoneStr {
						immMapBuilderPutStr(b, e.key, e.val)
					}
				}
				return
			}
			if smc.IsSorted() && smc.ext != nil && !b.mc.Identity {
				st := smc.ext.(*jSortedMapState)
				avlImmBuilderPutAll(st, st.root, b)
				return
			}
			smc.ForEachEntry(func(k, v any) bool {
				immMapBuilderPut(b, k, v)
				return true
			})
			return
		}
	}
	if ep, ok := any(m).(jFastEnumMapProvider); ok {
		if em := ep.fastEnumMap_(); em != nil {
			b.mc.InitExpectedSize(b.mc.Size() + em.size)
			em.ForEachEntry(func(k, v any) bool {
				immMapBuilderPut(b, k, v)
				return true
			})
			return
		}
	}
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		immMapBuilderPut(b, e.GetKey(), e.GetValue())
	}
}

func immMapBuilderBuild(b *jImmMapBuilderCore, throwOnDup bool) any {
	if throwOnDup && b.hasDup {
		panic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(b.dupKey)))
	}
	if b.sorted {
		if b.mc.Size() == 0 && b.mc.IsNaturalOrder() {
			return emptyImmutableSortedMap_
		}
		obj := &Ext_com_google_common_collect_ImmutableSortedMap{}
		obj.BindExtSelf_(obj)
		dst := copySortedCoreInto(&obj.inlineCore, b.mc)
		dst.sealImmutable()
		return obj
	}
	if b.mc.Size() == 0 {
		return emptyImmutableMap_
	}
	obj := newImmutableMapShell()
	dst := copyHashCoreInto(&obj.inlineCore, b.mc)
	dst.sealImmutable()
	return obj
}

func immBiMapBuilderBuild(b *jImmMapBuilderCore, throwOnDup bool) any {
	if throwOnDup && b.hasDup {
		panic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(b.dupKey)))
	}
	return newImmutableBiMapFrom(wrapLinkedHashMap(cloneJMapCore(b.mc)))
}

func immListSortedCopyOf(cmp Ext_java_util_Comparator, it Ext_java_lang_Iterable) Ext_com_google_common_collect_ImmutableList {
	CheckNotNull(it)
	var items []any
	for iter := NewIterCursor(it); iter.HasNext(); {
		items = append(items, CheckNotNull(iter.Next()))
	}
	if len(items) == 0 {
		return emptyImmutableList_
	}
	fn := comparatorFuncOrNatural(cmp)
	sort.SliceStable(items, func(i, j int) bool { return fn(items[i], items[j]) < 0 })
	return wrapImmutableListOwn(items)
}

func immSortedSetCopyOf(cmp Ext_java_util_Comparator, col any) Ext_com_google_common_collect_ImmutableSortedSet {
	CheckNotNull(col)
	nat := isNaturalOrderComparator(cmp)
	if nat {
		if iss, ok := col.(Ext_com_google_common_collect_ImmutableSortedSet); ok && !isJavaNilAny(iss) {
			if sp, ok2 := col.(jFastSetCoreProvider); ok2 {
				if smc := sp.fastSetCore_(); smc != nil && smc.IsNaturalOrder() {
					return iss
				}
			}
		}
		if c, ok := col.(Ext_java_util_Collection); ok && c.IsEmpty() {
			return emptyImmutableSortedSet_
		}
	}
	mc := newSortedCoreFromCmp(cmp)
	mapCoreAddAllNonNull(mc, col)
	return wrapImmutableSortedSet(mc)
}

// --- Bitmask EnumSet & Slice/Inline EnumMap ---

func resolveEnumUniverseFromElem(v any) []any {
	if cip, ok := v.(interface{ JavaClassInfo_() (string, string) }); ok {
		binName, _ := cip.JavaClassInfo_()
		if u := LookupEnumConstantsByName(binName); u != nil {
			return u
		}
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", v)))
}

type jEnumSetCore struct {
	universe []any
	bits0    uint64
	bits1    uint64
	more     []uint64
	size     int32
}

func (es *jEnumSetCore) Size() int32 { return es.size }

func (es *jEnumSetCore) hasBit(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		return (es.bits0 & mask) != 0
	case 1:
		return (es.bits1 & mask) != 0
	default:
		mw := w - 2
		return mw < len(es.more) && (es.more[mw]&mask) != 0
	}
}

func (es *jEnumSetCore) setOrd(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		if (es.bits0 & mask) != 0 {
			return false
		}
		es.bits0 |= mask
		es.size++
		return true
	case 1:
		if (es.bits1 & mask) != 0 {
			return false
		}
		es.bits1 |= mask
		es.size++
		return true
	default:
		mw := w - 2
		if mw >= len(es.more) {
			need := ((len(es.universe) + 63) >> 6) - 2
			if mw+1 > need {
				need = mw + 1
			}
			next := make([]uint64, need)
			copy(next, es.more)
			es.more = next
		}
		if (es.more[mw] & mask) != 0 {
			return false
		}
		es.more[mw] |= mask
		es.size++
		return true
	}
}

func (es *jEnumSetCore) clearOrd(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		if (es.bits0 & mask) == 0 {
			return false
		}
		es.bits0 &^= mask
		es.size--
		return true
	case 1:
		if (es.bits1 & mask) == 0 {
			return false
		}
		es.bits1 &^= mask
		es.size--
		return true
	default:
		mw := w - 2
		if mw >= len(es.more) || (es.more[mw]&mask) == 0 {
			return false
		}
		es.more[mw] &^= mask
		es.size--
		return true
	}
}

func (es *jEnumSetCore) nextSetBit(fromOrd int) int {
	if fromOrd < 0 {
		fromOrd = 0
	}
	w := fromOrd >> 6
	shift := uint(fromOrd & 63)
	if w == 0 {
		word := es.bits0 & (^uint64(0) << shift)
		if word != 0 {
			return bits.TrailingZeros64(word)
		}
		w = 1
		shift = 0
	}
	if w == 1 {
		word := es.bits1 & (^uint64(0) << shift)
		if word != 0 {
			return 64 + bits.TrailingZeros64(word)
		}
		w = 2
		shift = 0
	}
	for mw := w - 2; mw < len(es.more); mw++ {
		word := es.more[mw] & (^uint64(0) << shift)
		shift = 0
		if word != 0 {
			return ((mw + 2) << 6) + bits.TrailingZeros64(word)
		}
	}
	return -1
}

func (es *jEnumSetCore) recomputeSize() {
	n := bits.OnesCount64(es.bits0) + bits.OnesCount64(es.bits1)
	for _, w := range es.more {
		n += bits.OnesCount64(w)
	}
	es.size = int32(n)
}

func (es *jEnumSetCore) ordIfValid(v any) (int, bool) {
	v = ToAny(v)
	if v == nil {
		return -1, false
	}
	ea, ok := v.(EnumAccessor)
	if !ok {
		return -1, false
	}
	eb := ea.GetEnumBase()
	if eb == nil {
		return -1, false
	}
	if es.universe == nil {
		return -1, false
	}
	ord := int(eb.Ordinal_)
	if ord >= 0 && ord < len(es.universe) && es.universe[ord] == v {
		return ord, true
	}
	return -1, false
}

func (es *jEnumSetCore) mustOrd(v any) int {
	v = CheckNotNull(ToAny(v))
	ea, ok := v.(EnumAccessor)
	if !ok {
		panic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", v)))
	}
	eb := ea.GetEnumBase()
	if eb == nil {
		ThrowNPE()
	}
	if es.universe == nil {
		es.universe = resolveEnumUniverseFromElem(v)
	}
	ord := int(eb.Ordinal_)
	if ord >= 0 && ord < len(es.universe) && es.universe[ord] == v {
		return ord
	}
	panic(javalang.NewClassCastException("enum type mismatch in EnumSet"))
}

func (es *jEnumSetCore) Contains(v any) bool {
	if es.size == 0 {
		return false
	}
	ord, ok := es.ordIfValid(v)
	return ok && es.hasBit(ord)
}

func (es *jEnumSetCore) Add(v any) bool {
	return es.setOrd(es.mustOrd(v))
}

func (es *jEnumSetCore) Remove(v any) bool {
	if es.size == 0 {
		return false
	}
	ord, ok := es.ordIfValid(v)
	if !ok {
		return false
	}
	return es.clearOrd(ord)
}

func (es *jEnumSetCore) Clear() {
	es.bits0 = 0
	es.bits1 = 0
	if len(es.more) > 0 {
		clear(es.more)
	}
	es.size = 0
}

func (es *jEnumSetCore) sameUniverse(other *jEnumSetCore) bool {
	if len(es.universe) != len(other.universe) {
		return false
	}
	if len(es.universe) == 0 {
		return es.universe != nil && other.universe != nil
	}
	return es.universe[0] == other.universe[0]
}

func (es *jEnumSetCore) initCopy(src *jEnumSetCore) {
	es.universe = src.universe
	es.bits0 = src.bits0
	es.bits1 = src.bits1
	if len(src.more) > 0 {
		es.more = make([]uint64, len(src.more))
		copy(es.more, src.more)
	}
	es.size = src.size
}

func (es *jEnumSetCore) initAll(u []any) {
	es.universe = u
	n := len(u)
	es.size = int32(n)
	if n == 0 {
		return
	}
	if n <= 64 {
		if n == 64 {
			es.bits0 = ^uint64(0)
		} else {
			es.bits0 = (uint64(1) << uint(n)) - 1
		}
		return
	}
	es.bits0 = ^uint64(0)
	if n <= 128 {
		if n == 128 {
			es.bits1 = ^uint64(0)
		} else {
			es.bits1 = (uint64(1) << uint(n-64)) - 1
		}
		return
	}
	es.bits1 = ^uint64(0)
	mWords := ((n + 63) >> 6) - 2
	es.more = make([]uint64, mWords)
	for i := 0; i < mWords-1; i++ {
		es.more[i] = ^uint64(0)
	}
	rem := n & 63
	if rem == 0 {
		es.more[mWords-1] = ^uint64(0)
	} else {
		es.more[mWords-1] = (uint64(1) << uint(rem)) - 1
	}
}

func (es *jEnumSetCore) ForEach(fn func(any) bool) bool {
	if es.size == 0 {
		return true
	}
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		if !fn(es.universe[ord]) {
			return false
		}
	}
	return true
}

func (es *jEnumSetCore) ForEachConsumer(action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		action.Accept(es.universe[ord])
	}
}

func (es *jEnumSetCore) ToArray() []any {
	if es.size == 0 {
		return []any{}
	}
	out := make([]any, 0, es.size)
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		out = append(out, es.universe[ord])
	}
	return out
}

func (es *jEnumSetCore) ContainsAll(other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			if oes.size == 0 {
				return true
			}
			if es.size == 0 || !es.sameUniverse(oes) {
				return false
			}
			if (oes.bits0&^es.bits0) != 0 || (oes.bits1&^es.bits1) != 0 {
				return false
			}
			for i, ow := range oes.more {
				var ew uint64
				if i < len(es.more) {
					ew = es.more[i]
				}
				if (ow &^ ew) != 0 {
					return false
				}
			}
			return true
		}
	}
	return colContainsAll(es.Contains, other)
}

func (es *jEnumSetCore) AddAll(other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			if oes.size == 0 {
				return false
			}
			if es.universe == nil {
				es.universe = oes.universe
			} else if !es.sameUniverse(oes) {
				panic(javalang.NewClassCastException("enum type mismatch in EnumSet.addAll"))
			}
			oldSize := es.size
			es.bits0 |= oes.bits0
			es.bits1 |= oes.bits1
			if len(oes.more) > 0 {
				if len(es.more) < len(oes.more) {
					next := make([]uint64, len(oes.more))
					copy(next, es.more)
					es.more = next
				}
				for i, ow := range oes.more {
					es.more[i] |= ow
				}
			}
			es.recomputeSize()
			return es.size != oldSize
		}
	}
	changed := false
	for it := Iterate(other); it.HasNext(); {
		if es.Add(it.Next()) {
			changed = true
		}
	}
	return changed
}

func (es *jEnumSetCore) RemoveAll(other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if es.size == 0 {
		return false
	}
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			if oes.size == 0 || !es.sameUniverse(oes) {
				return false
			}
			oldSize := es.size
			es.bits0 &^= oes.bits0
			es.bits1 &^= oes.bits1
			for i := 0; i < len(es.more) && i < len(oes.more); i++ {
				es.more[i] &^= oes.more[i]
			}
			es.recomputeSize()
			return es.size != oldSize
		}
	}
	changed := false
	for it := Iterate(other); it.HasNext(); {
		if es.Remove(it.Next()) {
			changed = true
		}
	}
	return changed
}

func (es *jEnumSetCore) RetainAll(other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	if es.size == 0 {
		return false
	}
	if ep, ok := any(other).(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			if oes.size == 0 || !es.sameUniverse(oes) {
				es.Clear()
				return true
			}
			oldSize := es.size
			es.bits0 &= oes.bits0
			es.bits1 &= oes.bits1
			for i := range es.more {
				if i < len(oes.more) {
					es.more[i] &= oes.more[i]
				} else {
					es.more[i] = 0
				}
			}
			es.recomputeSize()
			return es.size != oldSize
		}
	}
	changed := false
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		if !other.Contains(es.universe[ord]) {
			es.clearOrd(ord)
			changed = true
		}
	}
	return changed
}

func (es *jEnumSetCore) RemoveIf(pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	if es.size == 0 {
		return false
	}
	changed := false
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		if pred.Test(es.universe[ord]) {
			es.clearOrd(ord)
			changed = true
		}
	}
	return changed
}

func (es *jEnumSetCore) Equals(other any) bool {
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	if ep, ok := other.(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			if es == oes {
				return true
			}
			if es.size != oes.size {
				return false
			}
			if es.size == 0 {
				return true
			}
			if !es.sameUniverse(oes) {
				return false
			}
			if es.bits0 != oes.bits0 || es.bits1 != oes.bits1 {
				return false
			}
			maxM := max(len(es.more), len(oes.more))
			for i := 0; i < maxM; i++ {
				var a, b uint64
				if i < len(es.more) {
					a = es.more[i]
				}
				if i < len(oes.more) {
					b = oes.more[i]
				}
				if a != b {
					return false
				}
			}
			return true
		}
	}
	if es.size != st.Size() {
		return false
	}
	if es.size == 0 {
		return true
	}
	return es.ContainsAll(st)
}

func (es *jEnumSetCore) HashCode() int32 {
	if es.size == 0 {
		return 0
	}
	h := int32(0)
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		h += es.universe[ord].(EnumAccessor).GetEnumBase().HashCode()
	}
	return h
}

func (es *jEnumSetCore) String() string {
	if es.size == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(es.universe[ord]))
	}
	sb.WriteByte(']')
	return sb.String()
}

type jEnumSetIter struct {
	es      *jEnumSetCore
	nextOrd int
	lastOrd int
}

func (*jEnumSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jEnumSetIter) HasNext() bool                          { return it.nextOrd >= 0 }
func (it *jEnumSetIter) Next() any {
	if it.nextOrd < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	ord := it.nextOrd
	it.lastOrd = ord
	it.nextOrd = it.es.nextSetBit(ord + 1)
	return it.es.universe[ord]
}
func (it *jEnumSetIter) Remove() {
	if it.lastOrd < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.es.clearOrd(it.lastOrd)
	it.lastOrd = -1
}

func (es *jEnumSetCore) Iterator() Ext_java_util_Iterator {
	if es.size == 0 {
		return emptyMapIter_
	}
	return &jEnumSetIter{es: es, nextOrd: es.nextSetBit(0), lastOrd: -1}
}

func enumSetNoneOf(cls *javalang.Class) Ext_java_util_EnumSet {
	CheckNotNullPtr(cls)
	u := LookupEnumConstants(cls)
	obj := newEnumSetShell()
	obj.inlineEnumSet.universe = u
	return obj
}

func enumSetAllOf(cls *javalang.Class) Ext_java_util_EnumSet {
	CheckNotNullPtr(cls)
	u := LookupEnumConstants(cls)
	if u == nil {
		panic(javalang.NewClassCastException(cls.GetName() + " not an enum"))
	}
	obj := newEnumSetShell()
	obj.inlineEnumSet.initAll(u)
	return obj
}

func enumSetClone(src *jEnumSetCore) Ext_java_util_EnumSet {
	obj := newEnumSetShell()
	obj.inlineEnumSet.initCopy(src)
	return obj
}

func enumSetCopyOf(src any) Ext_java_util_EnumSet {
	CheckNotNull(src)
	if ep, ok := src.(jFastEnumSetProvider); ok {
		if oes := ep.fastEnumSet_(); oes != nil {
			return enumSetClone(oes)
		}
	}
	obj := newEnumSetShell()
	es := &obj.inlineEnumSet
	addedAny := false
	for it := Iterate(src); it.HasNext(); {
		es.Add(CheckNotNull(it.Next()))
		addedAny = true
	}
	if !addedAny && es.universe == nil {
		panic(javalang.NewIllegalArgumentException("Collection is empty"))
	}
	return obj
}

func enumSetComplementOf(src Ext_java_util_EnumSet) Ext_java_util_EnumSet {
	CheckNotNull(src)
	obj := newEnumSetShell()
	dst := &obj.inlineEnumSet
	if ep, ok := any(src).(jFastEnumSetProvider); ok {
		if ses := ep.fastEnumSet_(); ses != nil && ses.universe != nil {
			dst.initAll(ses.universe)
			dst.bits0 &^= ses.bits0
			dst.bits1 &^= ses.bits1
			for i := 0; i < len(dst.more) && i < len(ses.more); i++ {
				dst.more[i] &^= ses.more[i]
			}
			dst.size -= ses.size
			return obj
		}
	}
	return obj
}

func enumSetRange(from, to any) Ext_java_util_EnumSet {
	CheckNotNull(from)
	CheckNotNull(to)
	obj := newEnumSetShell()
	es := &obj.inlineEnumSet
	fOrd := es.mustOrd(from)
	tOrd := es.mustOrd(to)
	if fOrd > tOrd {
		panic(javalang.NewIllegalArgumentException(StrVal(from) + " > " + StrVal(to)))
	}
	for ord := fOrd; ord <= tOrd; ord++ {
		es.setOrd(ord)
	}
	return obj
}

func immSetFromEnumSet(es *jEnumSetCore) Ext_com_google_common_collect_ImmutableSet {
	if es == nil || es.size == 0 {
		return emptyImmutableSet_
	}
	if es.size == 1 {
		ord := es.nextSetBit(0)
		return newImmSet1(es.universe[ord])
	}
	elems := make([]any, 0, es.size)
	for ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
		elems = append(elems, es.universe[ord])
	}
	return finishImmSetMulti(elems, nil)
}

func setsImmutableEnumSetIterable(src Ext_java_lang_Iterable) Ext_com_google_common_collect_ImmutableSet {
	CheckNotNull(src)
	if ep, ok := any(src).(jFastEnumSetProvider); ok {
		if es := ep.fastEnumSet_(); es != nil {
			return immSetFromEnumSet(es)
		}
	}
	var es jEnumSetCore
	for it := Iterate(src); it.HasNext(); {
		es.Add(CheckNotNull(it.Next()))
	}
	return immSetFromEnumSet(&es)
}

type jEnumMapCore struct {
	universe   []any
	vals       []any
	inlineVals [4]any
	bits0      uint64
	bits1      uint64
	more       []uint64
	size       int32
}

func (em *jEnumMapCore) Size() int32 { return em.size }

func (em *jEnumMapCore) initClass(cls *javalang.Class) {
	CheckNotNullPtr(cls)
	em.universe = LookupEnumConstants(cls)
}

func (em *jEnumMapCore) ensureVals() {
	if em.vals != nil {
		return
	}
	n := len(em.universe)
	if n <= len(em.inlineVals) {
		em.vals = em.inlineVals[:n]
	} else {
		em.vals = make([]any, n)
	}
}

func (em *jEnumMapCore) initCopyCore(src *jEnumMapCore) {
	em.universe = src.universe
	em.bits0 = src.bits0
	em.bits1 = src.bits1
	if len(src.more) > 0 {
		em.more = make([]uint64, len(src.more))
		copy(em.more, src.more)
	}
	em.size = src.size
	if src.size > 0 && src.vals != nil {
		em.ensureVals()
		copy(em.vals, src.vals)
	}
}

func (em *jEnumMapCore) initCopyEnumMap(src any) {
	CheckNotNull(src)
	if ep, ok := src.(jFastEnumMapProvider); ok {
		if sem := ep.fastEnumMap_(); sem != nil {
			em.initCopyCore(sem)
			return
		}
	}
	if m, ok := src.(Ext_java_util_Map); ok {
		em.initCopyMap(m)
	}
}

func (em *jEnumMapCore) initCopyMap(src Ext_java_util_Map) {
	CheckNotNull(src)
	if ep, ok := any(src).(jFastEnumMapProvider); ok {
		if sem := ep.fastEnumMap_(); sem != nil {
			em.initCopyCore(sem)
			return
		}
	}
	if src.IsEmpty() {
		panic(javalang.NewIllegalArgumentException("Specified map is empty"))
	}
	em.PutAll(src)
}

func (em *jEnumMapCore) hasBit(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		return (em.bits0 & mask) != 0
	case 1:
		return (em.bits1 & mask) != 0
	default:
		mw := w - 2
		return mw < len(em.more) && (em.more[mw]&mask) != 0
	}
}

func (em *jEnumMapCore) setBit(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		if (em.bits0 & mask) != 0 {
			return false
		}
		em.bits0 |= mask
		em.size++
		return true
	case 1:
		if (em.bits1 & mask) != 0 {
			return false
		}
		em.bits1 |= mask
		em.size++
		return true
	default:
		mw := w - 2
		if mw >= len(em.more) {
			need := ((len(em.universe) + 63) >> 6) - 2
			if mw+1 > need {
				need = mw + 1
			}
			next := make([]uint64, need)
			copy(next, em.more)
			em.more = next
		}
		if (em.more[mw] & mask) != 0 {
			return false
		}
		em.more[mw] |= mask
		em.size++
		return true
	}
}

func (em *jEnumMapCore) clearBit(ord int) bool {
	w := ord >> 6
	mask := uint64(1) << uint(ord&63)
	switch w {
	case 0:
		if (em.bits0 & mask) == 0 {
			return false
		}
		em.bits0 &^= mask
		em.size--
		if ord < len(em.vals) {
			em.vals[ord] = nil
		}
		return true
	case 1:
		if (em.bits1 & mask) == 0 {
			return false
		}
		em.bits1 &^= mask
		em.size--
		if ord < len(em.vals) {
			em.vals[ord] = nil
		}
		return true
	default:
		mw := w - 2
		if mw >= len(em.more) || (em.more[mw]&mask) == 0 {
			return false
		}
		em.more[mw] &^= mask
		em.size--
		if ord < len(em.vals) {
			em.vals[ord] = nil
		}
		return true
	}
}

func (em *jEnumMapCore) nextSetBit(fromOrd int) int {
	if fromOrd < 0 {
		fromOrd = 0
	}
	w := fromOrd >> 6
	shift := uint(fromOrd & 63)
	if w == 0 {
		word := em.bits0 & (^uint64(0) << shift)
		if word != 0 {
			return bits.TrailingZeros64(word)
		}
		w = 1
		shift = 0
	}
	if w == 1 {
		word := em.bits1 & (^uint64(0) << shift)
		if word != 0 {
			return 64 + bits.TrailingZeros64(word)
		}
		w = 2
		shift = 0
	}
	for mw := w - 2; mw < len(em.more); mw++ {
		word := em.more[mw] & (^uint64(0) << shift)
		shift = 0
		if word != 0 {
			return ((mw + 2) << 6) + bits.TrailingZeros64(word)
		}
	}
	return -1
}

func (em *jEnumMapCore) ordIfValid(k any) (int, bool) {
	k = ToAny(k)
	if k == nil || em.universe == nil {
		return -1, false
	}
	ea, ok := k.(EnumAccessor)
	if !ok {
		return -1, false
	}
	eb := ea.GetEnumBase()
	if eb == nil {
		return -1, false
	}
	ord := int(eb.Ordinal_)
	if ord >= 0 && ord < len(em.universe) && em.universe[ord] == k {
		return ord, true
	}
	return -1, false
}

func (em *jEnumMapCore) mustOrd(k any) int {
	k = CheckNotNull(ToAny(k))
	ea, ok := k.(EnumAccessor)
	if !ok {
		panic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", k)))
	}
	eb := ea.GetEnumBase()
	if eb == nil {
		ThrowNPE()
	}
	if em.universe == nil {
		em.universe = resolveEnumUniverseFromElem(k)
	}
	ord := int(eb.Ordinal_)
	if ord >= 0 && ord < len(em.universe) && em.universe[ord] == k {
		return ord
	}
	panic(javalang.NewClassCastException("enum key type mismatch in EnumMap"))
}

func (em *jEnumMapCore) ContainsKey(k any) bool {
	if em.size == 0 {
		return false
	}
	ord, ok := em.ordIfValid(k)
	return ok && em.hasBit(ord)
}

func (em *jEnumMapCore) ContainsValue(v any) bool {
	if em.size == 0 {
		return false
	}
	v = ToAny(v)
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		if javalang.ObjectsEquals(em.vals[ord], v) {
			return true
		}
	}
	return false
}

func (em *jEnumMapCore) Get(k any) any {
	if em.size == 0 {
		return nil
	}
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) {
		return nil
	}
	return em.vals[ord]
}

func (em *jEnumMapCore) GetOrDefault(k, def any) any {
	if em.size > 0 {
		if ord, ok := em.ordIfValid(k); ok && em.hasBit(ord) {
			return em.vals[ord]
		}
	}
	return ToAny(def)
}

func (em *jEnumMapCore) Put(k, v any) any {
	ord := em.mustOrd(k)
	em.ensureVals()
	v = ToAny(v)
	if em.hasBit(ord) {
		old := em.vals[ord]
		em.vals[ord] = v
		return old
	}
	em.setBit(ord)
	em.vals[ord] = v
	return nil
}

func (em *jEnumMapCore) PutIfAbsent(k, v any) any {
	ord := em.mustOrd(k)
	em.ensureVals()
	if em.hasBit(ord) {
		if cur := em.vals[ord]; cur != nil {
			return cur
		}
		em.vals[ord] = ToAny(v)
		return nil
	}
	em.setBit(ord)
	em.vals[ord] = ToAny(v)
	return nil
}

func (em *jEnumMapCore) Remove(k any) any {
	if em.size == 0 {
		return nil
	}
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) {
		return nil
	}
	old := em.vals[ord]
	em.clearBit(ord)
	return old
}

func (em *jEnumMapCore) Remove2(k, v any) bool {
	if em.size == 0 {
		return false
	}
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) || !javalang.ObjectsEquals(em.vals[ord], ToAny(v)) {
		return false
	}
	em.clearBit(ord)
	return true
}

func (em *jEnumMapCore) PutAll(src Ext_java_util_Map) {
	CheckNotNull(src)
	if ep, ok := any(src).(jFastEnumMapProvider); ok {
		if sem := ep.fastEnumMap_(); sem != nil {
			if sem.size == 0 {
				return
			}
			sem.ForEachEntry(func(k, v any) bool {
				em.Put(k, v)
				return true
			})
			return
		}
	}
	if mp, ok := any(src).(jFastMapCoreProvider); ok {
		if smc := mp.fastMapCore_(); smc != nil {
			smc.ForEachEntry(func(k, v any) bool {
				em.Put(k, v)
				return true
			})
			return
		}
	}
	if es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
		for it := Iterate(es); it.HasNext(); {
			e := Cast[Ext_java_util_Map_Entry](it.Next())
			em.Put(e.GetKey(), e.GetValue())
		}
	}
}

func (em *jEnumMapCore) Clear() {
	em.bits0 = 0
	em.bits1 = 0
	if len(em.more) > 0 {
		clear(em.more)
	}
	if len(em.vals) > 0 {
		clear(em.vals)
	}
	em.size = 0
}

func (em *jEnumMapCore) Replace2(k, newVal any) any {
	if em.size == 0 {
		return nil
	}
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) {
		return nil
	}
	old := em.vals[ord]
	em.vals[ord] = ToAny(newVal)
	return old
}

func (em *jEnumMapCore) Replace3(k, oldVal, newVal any) bool {
	if em.size == 0 {
		return false
	}
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) || !javalang.ObjectsEquals(em.vals[ord], ToAny(oldVal)) {
		return false
	}
	em.vals[ord] = ToAny(newVal)
	return true
}

func (em *jEnumMapCore) ReplaceAll(fn Ext_java_util_function_BiFunction) {
	CheckNotNull(fn)
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		em.vals[ord] = ToAny(fn.Apply(em.universe[ord], em.vals[ord]))
	}
}

func (em *jEnumMapCore) tryGetNonNull(k any) (any, bool) {
	ord, ok := em.ordIfValid(k)
	if ok && em.hasBit(ord) {
		if cur := em.vals[ord]; cur != nil {
			return cur, true
		}
	}
	return nil, false
}

func (em *jEnumMapCore) ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any {
	CheckNotNull(fn)
	ord := em.mustOrd(k)
	if em.hasBit(ord) {
		if cur := em.vals[ord]; cur != nil {
			return cur
		}
	}
	newVal := ToAny(fn.Apply(ToAny(k)))
	if newVal != nil {
		em.ensureVals()
		em.setBit(ord)
		em.vals[ord] = newVal
	}
	return newVal
}

func (em *jEnumMapCore) ComputeIfPresent(k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	ord, ok := em.ordIfValid(k)
	if !ok || !em.hasBit(ord) || em.vals[ord] == nil {
		return nil
	}
	newVal := ToAny(fn.Apply(ToAny(k), em.vals[ord]))
	if newVal == nil {
		em.clearBit(ord)
		return nil
	}
	em.vals[ord] = newVal
	return newVal
}

func (em *jEnumMapCore) Compute(k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	ord := em.mustOrd(k)
	haveOld := em.hasBit(ord)
	var oldVal any
	if haveOld {
		oldVal = em.vals[ord]
	}
	newVal := ToAny(fn.Apply(ToAny(k), oldVal))
	if newVal == nil {
		if haveOld {
			em.clearBit(ord)
		}
		return nil
	}
	em.ensureVals()
	em.setBit(ord)
	em.vals[ord] = newVal
	return newVal
}

func (em *jEnumMapCore) Merge(k, val any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	v := CheckNotNull(ToAny(val))
	ord := em.mustOrd(k)
	haveOld := em.hasBit(ord)
	var oldVal any
	if haveOld {
		oldVal = em.vals[ord]
	}
	var newVal any
	if oldVal == nil {
		newVal = v
	} else {
		newVal = ToAny(fn.Apply(oldVal, v))
	}
	if newVal == nil {
		if haveOld {
			em.clearBit(ord)
		}
		return nil
	}
	em.ensureVals()
	em.setBit(ord)
	em.vals[ord] = newVal
	return newVal
}

func (em *jEnumMapCore) ForEachEntry(fn func(k, v any) bool) {
	if em.size == 0 {
		return
	}
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		if !fn(em.universe[ord], em.vals[ord]) {
			return
		}
	}
}

func (em *jEnumMapCore) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
	CheckNotNull(action)
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		action.Accept(em.universe[ord], em.vals[ord])
	}
}

func (em *jEnumMapCore) Equals(other any) bool {
	mp, ok := other.(Ext_java_util_Map)
	if !ok || mp == nil || isJavaNilAny(mp) {
		return false
	}
	if em.size != mp.Size() {
		return false
	}
	if em.size == 0 {
		return true
	}
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		k := em.universe[ord]
		v := em.vals[ord]
		if v == nil {
			if mp.Get(k) != nil || !mp.ContainsKey(k) {
				return false
			}
		} else if !javalang.ObjectsEquals(v, ToAny(mp.Get(k))) {
			return false
		}
	}
	return true
}

func (em *jEnumMapCore) HashCode() int32 {
	if em.size == 0 {
		return 0
	}
	h := int32(0)
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		h += em.universe[ord].(EnumAccessor).GetEnumBase().HashCode() ^ javalang.ObjectHashCode(em.vals[ord])
	}
	return h
}

func (em *jEnumMapCore) String() string {
	if em.size == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(em.universe[ord]))
		sb.WriteByte('=')
		sb.WriteString(StrVal(em.vals[ord]))
	}
	sb.WriteByte('}')
	return sb.String()
}

type jLiveEnumMapKeySet struct {
	Ext_java_util_AbstractSetBase
	em *jEnumMapCore
}

func newLiveEnumMapKeySet(em *jEnumMapCore) Ext_java_util_Set {
	s := &jLiveEnumMapKeySet{em: em}
	s.BindExtSelf_(s)
	return s
}
func (s *jLiveEnumMapKeySet) Size() int32         { return s.em.size }
func (s *jLiveEnumMapKeySet) Contains(k any) bool { return s.em.ContainsKey(k) }
func (s *jLiveEnumMapKeySet) Remove_1_any(k any) bool {
	if s.em.size == 0 {
		return false
	}
	ord, ok := s.em.ordIfValid(k)
	return ok && s.em.clearBit(ord)
}
func (s *jLiveEnumMapKeySet) Clear() { s.em.Clear() }
func (s *jLiveEnumMapKeySet) Iterator() any {
	if s.em == nil || s.em.size == 0 {
		return emptyMapIter_
	}
	return &jEnumMapKeyIter{em: s.em, nextOrd: s.em.nextSetBit(0), lastOrd: -1}
}
func (s *jLiveEnumMapKeySet) SetContainsStr_(_ string) bool {
	if s == nil {
		ThrowNPE()
	}
	return false
}
func (s *jLiveEnumMapKeySet) SetAddStr_(_ string) bool {
	if s == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jLiveEnumMapKeySet) SetRemoveStr_(_ string) bool {
	if s == nil {
		ThrowNPE()
	}
	return false
}
func (s *jLiveEnumMapKeySet) SetContainsNode_(_ uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	return false
}
func (s *jLiveEnumMapKeySet) SetAddNode_(_ uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jLiveEnumMapKeySet) SetRemoveNode_(_ uint32) bool {
	if s == nil {
		ThrowNPE()
	}
	return false
}

type jEnumMapKeyIter struct {
	em      *jEnumMapCore
	nextOrd int
	lastOrd int
}

func (*jEnumMapKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jEnumMapKeyIter) HasNext() bool                          { return it.nextOrd >= 0 }
func (it *jEnumMapKeyIter) Next() any {
	if it.nextOrd < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	ord := it.nextOrd
	it.lastOrd = ord
	it.nextOrd = it.em.nextSetBit(ord + 1)
	return it.em.universe[ord]
}
func (it *jEnumMapKeyIter) Remove() {
	if it.lastOrd < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.em.clearBit(it.lastOrd)
	it.lastOrd = -1
}

type jLiveEnumMapValues struct {
	Ext_java_util_AbstractCollectionBase
	em *jEnumMapCore
}

func newLiveEnumMapValues(em *jEnumMapCore) Ext_java_util_Collection {
	v := &jLiveEnumMapValues{em: em}
	v.BindExtSelf_(v)
	return v
}
func (v *jLiveEnumMapValues) Size() int32         { return v.em.size }
func (v *jLiveEnumMapValues) Contains(x any) bool { return v.em.ContainsValue(x) }
func (v *jLiveEnumMapValues) Clear()              { v.em.Clear() }
func (v *jLiveEnumMapValues) Iterator() any {
	if v.em == nil || v.em.size == 0 {
		return emptyMapIter_
	}
	return &jEnumMapValIter{em: v.em, nextOrd: v.em.nextSetBit(0), lastOrd: -1}
}

type jEnumMapValIter struct {
	em      *jEnumMapCore
	nextOrd int
	lastOrd int
}

func (*jEnumMapValIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jEnumMapValIter) HasNext() bool                          { return it.nextOrd >= 0 }
func (it *jEnumMapValIter) Next() any {
	if it.nextOrd < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	ord := it.nextOrd
	it.lastOrd = ord
	it.nextOrd = it.em.nextSetBit(ord + 1)
	return it.em.vals[ord]
}
func (it *jEnumMapValIter) Remove() {
	if it.lastOrd < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.em.clearBit(it.lastOrd)
	it.lastOrd = -1
}

type jLiveEnumMapEntrySet struct {
	Ext_java_util_AbstractSetBase
	em *jEnumMapCore
}

func newLiveEnumMapEntrySet(em *jEnumMapCore) Ext_java_util_Set {
	s := &jLiveEnumMapEntrySet{em: em}
	s.BindExtSelf_(s)
	return s
}
func (s *jLiveEnumMapEntrySet) Size() int32 { return s.em.size }
func (s *jLiveEnumMapEntrySet) Clear()      { s.em.Clear() }
func (s *jLiveEnumMapEntrySet) Contains(o any) bool {
	e, ok := o.(Ext_java_util_Map_Entry)
	if !ok || e == nil || javalang.IsNil(e) {
		return false
	}
	ord, ok := s.em.ordIfValid(e.GetKey())
	return ok && s.em.hasBit(ord) && javalang.ObjectsEquals(s.em.vals[ord], ToAny(e.GetValue()))
}
func (s *jLiveEnumMapEntrySet) Remove_1_any(o any) bool {
	e, ok := o.(Ext_java_util_Map_Entry)
	if !ok || e == nil || javalang.IsNil(e) {
		return false
	}
	ord, ok := s.em.ordIfValid(e.GetKey())
	if !ok || !s.em.hasBit(ord) || !javalang.ObjectsEquals(s.em.vals[ord], ToAny(e.GetValue())) {
		return false
	}
	s.em.clearBit(ord)
	return true
}
func (s *jLiveEnumMapEntrySet) Iterator() any {
	if s.em == nil || s.em.size == 0 {
		return emptyMapIter_
	}
	return &jEnumMapEntryIter{em: s.em, nextOrd: s.em.nextSetBit(0), lastOrd: -1}
}

type jEnumMapEntry struct {
	em  *jEnumMapCore
	ord int
	k   any
	v   any
}

func (*jEnumMapEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
func (e *jEnumMapEntry) GetKey() any                              { return e.k }
func (e *jEnumMapEntry) GetValue() any {
	if e.em.hasBit(e.ord) {
		return e.em.vals[e.ord]
	}
	return e.v
}
func (e *jEnumMapEntry) SetValue(nv any) any {
	nv = ToAny(nv)
	if e.em.hasBit(e.ord) {
		old := e.em.vals[e.ord]
		e.em.vals[e.ord] = nv
		e.v = nv
		return old
	}
	old := e.v
	e.v = nv
	return old
}
func (e *jEnumMapEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
func (e *jEnumMapEntry) HashCode() int32   { return mapEntryHash(e) }
func (e *jEnumMapEntry) String() string    { return StrVal(e.GetKey()) + "=" + StrVal(e.GetValue()) }

type jEnumMapEntryIter struct {
	em      *jEnumMapCore
	nextOrd int
	lastOrd int
}

func (*jEnumMapEntryIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jEnumMapEntryIter) HasNext() bool                          { return it.nextOrd >= 0 }
func (it *jEnumMapEntryIter) Next() any {
	if it.nextOrd < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	ord := it.nextOrd
	it.lastOrd = ord
	it.nextOrd = it.em.nextSetBit(ord + 1)
	return &jEnumMapEntry{em: it.em, ord: ord, k: it.em.universe[ord], v: it.em.vals[ord]}
}
func (it *jEnumMapEntryIter) Remove() {
	if it.lastOrd < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.em.clearBit(it.lastOrd)
	it.lastOrd = -1
}

func mapsImmutableEnumMap(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableMap {
	CheckNotNull(src)
	if src.IsEmpty() {
		return emptyImmutableMap_
	}
	if ep, ok := any(src).(jFastEnumMapProvider); ok {
		if sem := ep.fastEnumMap_(); sem != nil {
			obj := newImmutableMapShell()
			mc := &obj.inlineCore
			mc.InitExpectedSize(sem.size)
			sem.ForEachEntry(func(k, v any) bool {
				mc.Put(CheckNotNull(k), CheckNotNull(v))
				return true
			})
			mc.sealImmutable()
			return obj
		}
	}
	var tmp jEnumMapCore
	if es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
		for it := Iterate(es); it.HasNext(); {
			e := Cast[Ext_java_util_Map_Entry](it.Next())
			tmp.Put(CheckNotNull(ToAny(e.GetKey())), CheckNotNull(ToAny(e.GetValue())))
		}
	}
	obj := newImmutableMapShell()
	mc := &obj.inlineCore
	mc.InitExpectedSize(tmp.size)
	tmp.ForEachEntry(func(k, v any) bool {
		mc.Put(k, v)
		return true
	})
	mc.sealImmutable()
	return obj
}

// --- Abstract Delegation Helpers ---

func iterToSliceDirect(col Ext_java_util_Collection) []any {
	if col == nil || javalang.IsNil(col) {
		ThrowNPE()
	}
	out := make([]any, 0, col.Size())
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		out = append(out, ToAny(it.Next()))
	}
	return out
}

func colToSlice(col Ext_java_util_Collection) []any {
	if tv, ok := any(col).(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			out := make([]any, n)
			for i := 0; i < n; i++ {
				out[i] = ToAny(tv.typedViewAt_(i))
			}
			return out
		}
	}
	if si, ok := any(col).(SliceIterable); ok {
		src := si.IterSlice_()
		out := make([]any, len(src))
		for i, x := range src {
			out[i] = ToAny(x)
		}
		return out
	}
	out := make([]any, 0, col.Size())
	for it := NewIterCursor(col); it.HasNext(); {
		out = append(out, ToAny(it.Next()))
	}
	return out
}

func abstractColContains(col Ext_java_util_Collection, target any) bool {
	t := ToAny(target)
	for it := Iterate(col); it.HasNext(); {
		if javalang.ObjectsEquals(t, ToAny(it.Next())) {
			return true
		}
	}
	return false
}

func abstractColRemove(col Ext_java_util_Collection, target any) bool {
	t := ToAny(target)
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		if javalang.ObjectsEquals(t, ToAny(it.Next())) {
			it.Remove()
			return true
		}
	}
	return false
}

func abstractColAddAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	modified := false
	for it := Iterate(other); it.HasNext(); {
		if col.Add_1(it.Next()) {
			modified = true
		}
	}
	return modified
}

func abstractColRemoveAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	modified := false
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		if other.Contains(it.Next()) {
			it.Remove()
			modified = true
		}
	}
	return modified
}

func abstractColRetainAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
	CheckNotNull(other)
	modified := false
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		if !other.Contains(it.Next()) {
			it.Remove()
			modified = true
		}
	}
	return modified
}

func abstractColRemoveIf(col Ext_java_util_Collection, pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	modified := false
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		if pred.Test(it.Next()) {
			it.Remove()
			modified = true
		}
	}
	return modified
}

func abstractColClear(col Ext_java_util_Collection) {
	it := Cast[Ext_java_util_Iterator](col.Iterator())
	for it.HasNext() {
		_ = it.Next()
		it.Remove()
	}
}

func abstractSetEquals(self Ext_java_util_Collection, other any) bool {
	if RefEquals(self, other) {
		return true
	}
	st, ok := other.(Ext_java_util_Set)
	if !ok || st == nil || javalang.IsNil(st) {
		return false
	}
	if self.Size() != st.Size() {
		return false
	}
	return self.ContainsAll(st)
}

func abstractSetHash(self Ext_java_util_Collection) int32 {
	h := int32(0)
	for it := Iterate(self); it.HasNext(); {
		h += javalang.ObjectHashCode(ToAny(it.Next()))
	}
	return h
}

// AbstractList delegation
type jAbstractListIter struct {
	lst     Ext_java_util_List
	cursor  int32
	lastRet int32
}

func (*jAbstractListIter) Implements_j2gort_Ext_java_util_Iterator()     {}
func (*jAbstractListIter) Implements_j2gort_Ext_java_util_ListIterator() {}
func (it *jAbstractListIter) HasNext() bool                              { return it.cursor < it.lst.Size() }
func (it *jAbstractListIter) Next() any {
	if it.cursor >= it.lst.Size() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.lst.Get(it.cursor)
	it.lastRet = it.cursor
	it.cursor++
	return v
}
func (it *jAbstractListIter) Remove() {
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.lst.Remove_1_int32(it.lastRet)
	if it.lastRet < it.cursor {
		it.cursor--
	}
	it.lastRet = -1
}

func abstractListOp_size_0(l Ext_java_util_List) int32 {
	panic(javalang.NewUnsupportedOperationException("AbstractList.size"))
}
func abstractListOp_get_1(l Ext_java_util_List, i int32) any {
	panic(javalang.NewUnsupportedOperationException("AbstractList.get"))
}
func abstractListOp_set_2(l Ext_java_util_List, i int32, v any) any {
	panic(javalang.NewUnsupportedOperationException("AbstractList.set"))
}
func abstractListOp_add_2(l Ext_java_util_List, i int32, v any) {
	panic(javalang.NewUnsupportedOperationException("AbstractList.add"))
}
func abstractListOp_remove_1(l Ext_java_util_List, arg any) any {
	if idx, ok := arg.(int32); ok {
		_ = idx
		panic(javalang.NewUnsupportedOperationException("AbstractList.remove(int)"))
	}
	return abstractColRemove(l, arg)
}
func abstractListOp_add_1(l Ext_java_util_List, v any) bool {
	l.Add_2(l.Size(), v)
	return true
}
func abstractListOp_isEmpty_0(l Ext_java_util_List) bool         { return l.Size() == 0 }
func abstractListOp_contains_1(l Ext_java_util_List, v any) bool { return l.IndexOf_1(v) >= 0 }
func abstractListOp_indexOf_1(l Ext_java_util_List, v any) int32 {
	n := l.Size()
	t := ToAny(v)
	for i := int32(0); i < n; i++ {
		if javalang.ObjectsEquals(t, ToAny(l.Get(i))) {
			return i
		}
	}
	return -1
}
func abstractListOp_lastIndexOf_1(l Ext_java_util_List, v any) int32 {
	t := ToAny(v)
	for i := l.Size() - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(t, ToAny(l.Get(i))) {
			return i
		}
	}
	return -1
}
func abstractListOp_clear_0(l Ext_java_util_List) { abstractColClear(l) }
func abstractListOp_iterator_0(l Ext_java_util_List) any {
	return &jAbstractListIter{lst: l, lastRet: -1}
}
func abstractListOp_listIterator_0(l Ext_java_util_List) any {
	return abstractListOp_listIterator_1(l, 0)
}
func abstractListOp_listIterator_1(l Ext_java_util_List, i int32) any {
	CheckPositionIndex(i, l.Size())
	return &jAbstractListIter{lst: l, cursor: i, lastRet: -1}
}
func abstractListOp_subList_2(l Ext_java_util_List, from, to int32) any {
	CheckPositionIndexes(from, to, l.Size())
	items := make([]any, to-from)
	for i := from; i < to; i++ {
		items[i-from] = l.Get(i)
	}
	return wrapArrayList(items)
}
func abstractListOp_toArray_0(l Ext_java_util_List) []any { return colToSlice(l) }
func abstractListOp_toArray_1(l Ext_java_util_List, dst []any) []any {
	return sliceToArray1(colToSlice(l), dst)
}
func abstractListOp_containsAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
	return colContainsAll(l.Contains, c)
}
func abstractListOp_addAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
	return abstractColAddAll(l, c)
}
func abstractListOp_removeAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
	return abstractColRemoveAll(l, c)
}
func abstractListOp_retainAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
	return abstractColRetainAll(l, c)
}
func abstractListOp_removeIf_1(l Ext_java_util_List, p Ext_java_util_function_Predicate) bool {
	return abstractColRemoveIf(l, p)
}
func abstractListOp_equals_1(l Ext_java_util_List, o any) bool {
	return listEqualsSlice(colToSlice(l), o)
}
func abstractListOp_hashCode_0(l Ext_java_util_List) int32  { return listHashSlice(colToSlice(l)) }
func abstractListOp_toString_0(l Ext_java_util_List) string { return sliceJavaString(colToSlice(l)) }
func abstractListOp_sort_1(l Ext_java_util_List, cmp Ext_java_util_Comparator) {
	s := colToSlice(l)
	fn := comparatorFuncOrNatural(cmp)
	sort.SliceStable(s, func(i, j int) bool { return fn(s[i], s[j]) < 0 })
	for i, v := range s {
		l.Set(int32(i), v)
	}
}
func abstractListOp_addFirst_1(l Ext_java_util_List, v any) { l.Add_2(0, v) }
func abstractListOp_addLast_1(l Ext_java_util_List, v any)  { l.Add_2(l.Size(), v) }
func abstractListOp_getFirst_0(l Ext_java_util_List) any {
	if l.IsEmpty() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Get(0)
}
func abstractListOp_getLast_0(l Ext_java_util_List) any {
	if l.IsEmpty() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Get(l.Size() - 1)
}
func abstractListOp_removeFirst_0(l Ext_java_util_List) any {
	if l.IsEmpty() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Remove_1_int32(0)
}
func abstractListOp_removeLast_0(l Ext_java_util_List) any {
	if l.IsEmpty() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Remove_1_int32(l.Size() - 1)
}
func abstractListOp_stream_0(l Ext_java_util_List) Ext_java_util_stream_Stream {
	return streamFromIterable(l)
}
func abstractListOp_forEach_1(l Ext_java_util_List, c Ext_java_util_function_Consumer) {
	ExtDefault_Iterable_ForEach(l, c)
}

// AbstractQueue delegation
func abstractQueueOp_size_0(q Ext_java_util_Queue) int32 {
	panic(javalang.NewUnsupportedOperationException("Queue.size"))
}
func abstractQueueOp_iterator_0(q Ext_java_util_Queue) any {
	panic(javalang.NewUnsupportedOperationException("Queue.iterator"))
}
func abstractQueueOp_poll_0(q Ext_java_util_Queue) any {
	panic(javalang.NewUnsupportedOperationException("Queue.poll"))
}
func abstractQueueOp_peek_0(q Ext_java_util_Queue) any {
	panic(javalang.NewUnsupportedOperationException("Queue.peek"))
}
func abstractQueueOp_add_1(q Ext_java_util_Queue, v any) bool {
	panic(javalang.NewUnsupportedOperationException("Queue.add"))
}
func abstractQueueOp_remove_0(q Ext_java_util_Queue) any {
	x := q.Poll_0()
	if x != nil {
		return x
	}
	panic(javalang.NewNoSuchElementException("empty queue"))
}
func abstractQueueOp_element_0(q Ext_java_util_Queue) any {
	x := q.Peek()
	if x != nil {
		return x
	}
	panic(javalang.NewNoSuchElementException("empty queue"))
}
func abstractQueueOp_remove_1(q Ext_java_util_Queue, v any) bool { return abstractColRemove(q, v) }
func abstractQueueOp_retainAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
	return abstractColRetainAll(q, c)
}
func abstractQueueOp_removeIf_1(q Ext_java_util_Queue, p Ext_java_util_function_Predicate) bool {
	return abstractColRemoveIf(q, p)
}
func abstractQueueOp_clear_0(q Ext_java_util_Queue) {
	for q.Poll_0() != nil {
	}
}
func abstractQueueOp_addAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
	return abstractColAddAll(q, c)
}
func abstractQueueOp_isEmpty_0(q Ext_java_util_Queue) bool         { return q.Size() == 0 }
func abstractQueueOp_contains_1(q Ext_java_util_Queue, v any) bool { return abstractColContains(q, v) }
func abstractQueueOp_toArray_0(q Ext_java_util_Queue) []any        { return colToSlice(q) }
func abstractQueueOp_toArray_1(q Ext_java_util_Queue, d []any) []any {
	return sliceToArray1(colToSlice(q), d)
}
func abstractQueueOp_containsAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
	return colContainsAll(q.Contains, c)
}
func abstractQueueOp_removeAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
	return abstractColRemoveAll(q, c)
}
func abstractQueueOp_toString_0(q Ext_java_util_Queue) string    { return sliceJavaString(colToSlice(q)) }
func abstractQueueOp_equals_1(q Ext_java_util_Queue, o any) bool { return RefEquals(q, o) }
func abstractQueueOp_hashCode_0(q Ext_java_util_Queue) int32     { return javalang.IdentityHashCode(q) }
func abstractQueueOp_stream_0(q Ext_java_util_Queue) Ext_java_util_stream_Stream {
	return streamFromIterable(q)
}
func abstractQueueOp_forEach_1(q Ext_java_util_Queue, c Ext_java_util_function_Consumer) {
	ExtDefault_Iterable_ForEach(q, c)
}

// AbstractMap delegation (via EntrySet())
func abstractMapOp_entrySet_0(m Ext_java_util_Map) any {
	panic(javalang.NewUnsupportedOperationException("AbstractMap.entrySet"))
}
func abstractMapOp_size_0(m Ext_java_util_Map) int32 {
	return Cast[Ext_java_util_Set](m.EntrySet()).Size()
}
func abstractMapOp_isEmpty_0(m Ext_java_util_Map) bool { return m.Size() == 0 }
func abstractMapOp_containsKey_1(m Ext_java_util_Map, k any) bool {
	t := ToAny(k)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		if javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
			return true
		}
	}
	return false
}
func abstractMapOp_containsValue_1(m Ext_java_util_Map, v any) bool {
	t := ToAny(v)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		if javalang.ObjectsEquals(t, ToAny(e.GetValue())) {
			return true
		}
	}
	return false
}
func abstractMapOp_get_1(m Ext_java_util_Map, k any) any {
	t := ToAny(k)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		if javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
			return e.GetValue()
		}
	}
	return nil
}
func abstractMapOp_put_2(m Ext_java_util_Map, k, v any) any {
	panic(javalang.NewUnsupportedOperationException("AbstractMap.put"))
}
func abstractMapOp_remove_1(m Ext_java_util_Map, k any) any {
	t := ToAny(k)
	it := Cast[Ext_java_util_Iterator](Cast[Ext_java_util_Set](m.EntrySet()).Iterator())
	for it.HasNext() {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		if javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
			old := e.GetValue()
			it.Remove()
			return old
		}
	}
	return nil
}
func abstractMapOp_putAll_1(m Ext_java_util_Map, other Ext_java_util_Map) {
	CheckNotNull(other)
	for it := Iterate(other.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		m.Put_2(e.GetKey(), e.GetValue())
	}
}
func abstractMapOp_clear_0(m Ext_java_util_Map) {
	Cast[Ext_java_util_Set](m.EntrySet()).Clear()
}
func abstractMapOp_keySet_0(m Ext_java_util_Map) any {
	var keys []any
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		keys = append(keys, Cast[Ext_java_util_Map_Entry](it.Next()).GetKey())
	}
	return wrapLinkedHashSet(keys)
}
func abstractMapOp_values_0(m Ext_java_util_Map) any {
	var vals []any
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		vals = append(vals, Cast[Ext_java_util_Map_Entry](it.Next()).GetValue())
	}
	return wrapArrayList(vals)
}
func abstractMapOp_equals_1(m Ext_java_util_Map, o any) bool {
	mp, ok := o.(Ext_java_util_Map)
	if !ok || mp == nil || javalang.IsNil(mp) {
		return false
	}
	if m.Size() != mp.Size() {
		return false
	}
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		k, v := e.GetKey(), e.GetValue()
		if v == nil {
			if mp.Get(k) != nil || !mp.ContainsKey(k) {
				return false
			}
		} else if !javalang.ObjectsEquals(ToAny(v), ToAny(mp.Get(k))) {
			return false
		}
	}
	return true
}
func abstractMapOp_hashCode_0(m Ext_java_util_Map) int32 {
	h := int32(0)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		h += mapEntryHash(Cast[Ext_java_util_Map_Entry](it.Next()))
	}
	return h
}
func abstractMapOp_toString_0(m Ext_java_util_Map) string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		k, v := e.GetKey(), e.GetValue()
		if k == any(m) {
			sb.WriteString("(this Map)")
		} else {
			sb.WriteString(StrVal(k))
		}
		sb.WriteByte('=')
		if v == any(m) {
			sb.WriteString("(this Map)")
		} else {
			sb.WriteString(StrVal(v))
		}
	}
	sb.WriteByte('}')
	return sb.String()
}
func abstractMapOp_getOrDefault_2(m Ext_java_util_Map, k, def any) any {
	if v := m.Get(k); v != nil || m.ContainsKey(k) {
		return v
	}
	return ToAny(def)
}
func abstractMapOp_forEach_1(m Ext_java_util_Map, fn Ext_java_util_function_BiConsumer) {
	CheckNotNull(fn)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		fn.Accept(e.GetKey(), e.GetValue())
	}
}
func abstractMapOp_replaceAll_1(m Ext_java_util_Map, fn Ext_java_util_function_BiFunction) {
	CheckNotNull(fn)
	for it := Iterate(m.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		e.SetValue(fn.Apply(e.GetKey(), e.GetValue()))
	}
}
func abstractMapOp_putIfAbsent_2(m Ext_java_util_Map, k, v any) any {
	if cur := m.Get(k); cur != nil {
		return cur
	}
	return m.Put_2(k, v)
}
func abstractMapOp_replace_3(m Ext_java_util_Map, k, oldV, newV any) bool {
	if !m.ContainsKey(k) || !javalang.ObjectsEquals(ToAny(m.Get(k)), ToAny(oldV)) {
		return false
	}
	m.Put_2(k, newV)
	return true
}
func abstractMapOp_computeIfAbsent_2(m Ext_java_util_Map, k any, fn Ext_java_util_function_Function) any {
	CheckNotNull(fn)
	if v := m.Get(k); v != nil {
		return v
	}
	nv := ToAny(fn.Apply(ToAny(k)))
	if nv != nil {
		m.Put_2(k, nv)
	}
	return nv
}
func abstractMapOp_compute_2(m Ext_java_util_Map, k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	oldV := m.Get(k)
	nv := ToAny(fn.Apply(ToAny(k), oldV))
	if nv == nil {
		if m.ContainsKey(k) {
			m.Remove_1(k)
		}
		return nil
	}
	m.Put_2(k, nv)
	return nv
}
func abstractMapOp_merge_3(m Ext_java_util_Map, k, v any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	v = CheckNotNull(ToAny(v))
	oldV := m.Get(k)
	var nv any
	if oldV == nil {
		nv = v
	} else {
		nv = ToAny(fn.Apply(oldV, v))
	}
	if nv == nil {
		m.Remove_1(k)
	} else {
		m.Put_2(k, nv)
	}
	return nv
}

// --- Sets.SetView ---

type jSetViewCore struct {
	s1   Ext_java_util_Set
	s2   Ext_java_util_Set
	mode int // 0=union, 1=intersection, 2=difference
}

func newSetView(s1, s2 Ext_java_util_Set, mode int) Ext_com_google_common_collect_Sets_SetView {
	CheckNotNull(s1)
	CheckNotNull(s2)
	v := &Ext_com_google_common_collect_Sets_SetViewBase{}
	v.BindExtSelf_(v)
	v.Core_ = &jSetViewCore{s1: s1, s2: s2, mode: mode}
	return v
}

func (c *jSetViewCore) contains(x any) bool {
	switch c.mode {
	case 0:
		return c.s1.Contains(x) || c.s2.Contains(x)
	case 1:
		return c.s1.Contains(x) && c.s2.Contains(x)
	default:
		return c.s1.Contains(x) && !c.s2.Contains(x)
	}
}

func (c *jSetViewCore) compute() []any {
	var sc JSetCore
	switch c.mode {
	case 0:
		for it := Iterate(c.s1); it.HasNext(); {
			sc.Add(it.Next())
		}
		for it := Iterate(c.s2); it.HasNext(); {
			sc.Add(it.Next())
		}
	case 1:
		for it := Iterate(c.s1); it.HasNext(); {
			x := it.Next()
			if c.s2.Contains(x) {
				sc.Add(x)
			}
		}
	default:
		for it := Iterate(c.s1); it.HasNext(); {
			x := it.Next()
			if !c.s2.Contains(x) {
				sc.Add(x)
			}
		}
	}
	return sc.LiveKeys()
}

// --- ConcurrentHashMap ---

type jSyncMapCore struct {
	mu sync.Mutex
	mc *JMapCore
}

func syncMapCore(m *Ext_java_util_concurrent_ConcurrentHashMap) *jSyncMapCore {
	CheckNotNullPtr(m)
	if m.Core_ == nil {
		m.Core_ = &jSyncMapCore{mc: NewJMapCore(false)}
	}
	return m.Core_.(*jSyncMapCore)
}

func syncMapOp_size_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) int32 {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.Size()
}
func syncMapOp_isEmpty_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) bool {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.IsEmpty()
}
func syncMapOp_get_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) any {
	CheckNotNull(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.Get(k)
}
func syncMapOp_getOrDefault_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, def any) any {
	CheckNotNull(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.GetOrDefault(k, def)
}
func syncMapOp_containsKey_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) bool {
	CheckNotNull(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.ContainsKey(k)
}
func syncMapOp_containsValue_1(m *Ext_java_util_concurrent_ConcurrentHashMap, v any) bool {
	CheckNotNull(v)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.ContainsValue(v)
}
func syncMapOp_put_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any) any {
	CheckNotNull(k)
	CheckNotNull(v)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.Put(k, v)
}
func syncMapOp_putIfAbsent_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any) any {
	CheckNotNull(k)
	CheckNotNull(v)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.PutIfAbsent(k, v)
}
func syncMapOp_putAll_1(m *Ext_java_util_concurrent_ConcurrentHashMap, src Ext_java_util_Map) {
	CheckNotNull(src)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	copyExtMapInto(c.mc, src)
}
func syncMapOp_remove_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) any {
	CheckNotNull(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.Remove(k)
}
func syncMapOp_clear_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mc.Clear()
}
func syncMapOp_keySet_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return newLiveMapKeySet(c.mc)
}
func syncMapOp_values_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return newLiveMapValues(c.mc)
}
func syncMapOp_entrySet_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return newLiveMapEntrySet(c.mc)
}
func syncMapOp_hashCode_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) int32 {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreHash(c.mc)
}
func syncMapOp_toString_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) string {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreString(c.mc)
}
func syncMapOp_equals_1(m *Ext_java_util_concurrent_ConcurrentHashMap, o any) bool {
	c := syncMapCore(m)
	c.mu.Lock()
	snapshot := cloneJMapCore(c.mc)
	c.mu.Unlock()
	return mapCoreEquals(snapshot, o)
}
func syncMapOp_replace_3(m *Ext_java_util_concurrent_ConcurrentHashMap, k, ov, nv any) bool {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreReplace3(c.mc, k, ov, nv)
}
func syncMapOp_forEach_1(m *Ext_java_util_concurrent_ConcurrentHashMap, fn Ext_java_util_function_BiConsumer) {
	c := syncMapCore(m)
	c.mu.Lock()
	snapshot := cloneJMapCore(c.mc)
	c.mu.Unlock()
	mapCoreForEach(snapshot, fn)
}
func syncMapOp_replaceAll_1(m *Ext_java_util_concurrent_ConcurrentHashMap, fn Ext_java_util_function_BiFunction) {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	mapCoreReplaceAll(c.mc, fn)
}
func syncMapOp_computeIfAbsent_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k any, fn Ext_java_util_function_Function) any {
	CheckNotNull(k)
	CheckNotNull(fn)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreComputeIfAbsent(c.mc, k, fn)
}
func syncMapComputeIfAbsentStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, fn Ext_java_util_function_Function) any {
	CheckNotNullString(k)
	CheckNotNull(fn)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	if v := c.mc.GetStr(k); v != nil {
		return v
	}
	nv := ToAny(fn.Apply(BoxString(k)))
	if nv != nil {
		c.mc.PutStr(k, nv)
	}
	return nv
}
func syncMapComputeIfAbsentNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, fn Ext_java_util_function_Function) any {
	if n == 0 {
		ThrowNPE()
	}
	CheckNotNull(fn)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	if v := c.mc.GetNode(n); v != nil {
		return v
	}
	nv := ToAny(fn.Apply(boxNodeById(n)))
	if nv != nil {
		c.mc.PutNode(n, nv)
	}
	return nv
}
func syncMapTryGetComputeStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) (any, bool) {
	if m == nil {
		ThrowNPE()
	}
	if k == NullString {
		return nil, false
	}
	c := syncMapCore(m)
	c.mu.Lock()
	v := c.mc.GetStr(k)
	c.mu.Unlock()
	if v != nil {
		return v, true
	}
	return nil, false
}
func syncMapTryGetComputeNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) (any, bool) {
	if m == nil {
		ThrowNPE()
	}
	if n == 0 {
		return nil, false
	}
	c := syncMapCore(m)
	c.mu.Lock()
	v := c.mc.GetNode(n)
	c.mu.Unlock()
	if v != nil {
		return v, true
	}
	return nil, false
}
func syncMapTryGetComputeAny(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) (any, bool) {
	if m == nil {
		ThrowNPE()
	}
	if k == nil || isJavaNilAny(k) {
		return nil, false
	}
	c := syncMapCore(m)
	c.mu.Lock()
	v := c.mc.Get(k)
	c.mu.Unlock()
	if v != nil {
		return v, true
	}
	return nil, false
}
func syncMapOp_compute_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k any, fn Ext_java_util_function_BiFunction) any {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreCompute(c.mc, k, fn)
}
func syncMapOp_merge_3(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any, fn Ext_java_util_function_BiFunction) any {
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return mapCoreMerge(c.mc, k, v, fn)
}
func syncMapGetStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) any {
	CheckNotNullString(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.GetStr(k)
}
func syncMapGetOrDefaultStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, def any) any {
	CheckNotNullString(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.GetOrDefaultStr(k, def)
}
func syncMapContainsStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) bool {
	CheckNotNullString(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.ContainsStr(k)
}
func syncMapPutStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, v any) any {
	CheckNotNullString(k)
	CheckNotNull(v)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.PutStr(k, v)
}
func syncMapRemoveStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) any {
	CheckNotNullString(k)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.RemoveStr(k)
}
func syncMapGetNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) any {
	if n == 0 {
		ThrowNPE()
	}
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.GetNode(n)
}
func syncMapGetOrDefaultNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, def any) any {
	if n == 0 {
		ThrowNPE()
	}
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.GetOrDefaultNode(n, def)
}
func syncMapContainsNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) bool {
	if n == 0 {
		ThrowNPE()
	}
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.ContainsNode(n)
}
func syncMapPutNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, v any) any {
	if n == 0 {
		ThrowNPE()
	}
	CheckNotNull(v)
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.PutNode(n, v)
}
func syncMapRemoveNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) any {
	if n == 0 {
		ThrowNPE()
	}
	c := syncMapCore(m)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mc.RemoveNode(n)
}

// --- BiMap (HashBiMap & ImmutableBiMap) ---

type jBiMapCore struct {
	fwd           *JMapCore
	rev           *JMapCore
	inverse       any
	immutable     bool
	keySetCache   atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
	valuesCache   atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
	entrySetCache atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
}

func newHashBiMap() *Ext_com_google_common_collect_HashBiMap {
	fwd := NewJMapCore(false)
	rev := NewJMapCore(false)
	m1 := &Ext_com_google_common_collect_HashBiMap{}
	m2 := &Ext_com_google_common_collect_HashBiMap{}
	c1 := &jBiMapCore{fwd: fwd, rev: rev, inverse: m2}
	c2 := &jBiMapCore{fwd: rev, rev: fwd, inverse: m1}
	m1.Core_ = c1
	m2.Core_ = c2
	return m1
}

func newImmutableBiMapFrom(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableBiMap {
	fwd := NewJMapCore(false)
	rev := NewJMapCore(false)
	if src != nil && !javalang.IsNil(src) {
		for it := Iterate(src.EntrySet()); it.HasNext(); {
			e := Cast[Ext_java_util_Map_Entry](it.Next())
			k := CheckNotNull(ToAny(e.GetKey()))
			v := CheckNotNull(ToAny(e.GetValue()))
			if fwd.ContainsKey(k) || rev.ContainsKey(v) {
				panic(javalang.NewIllegalArgumentException("duplicate key or value in ImmutableBiMap"))
			}
			fwd.Put(k, v)
			rev.Put(v, k)
		}
	}
	m1 := &Ext_com_google_common_collect_ImmutableBiMapBase{}
	m1.BindExtSelf_(m1)
	m2 := &Ext_com_google_common_collect_ImmutableBiMapBase{}
	m2.BindExtSelf_(m2)
	m1.Core_ = &jBiMapCore{fwd: fwd, rev: rev, inverse: m2, immutable: true}
	m2.Core_ = &jBiMapCore{fwd: rev, rev: fwd, inverse: m1, immutable: true}
	return m1
}

func biMapOp_size_0(b *jBiMapCore) int32                 { return b.fwd.Size() }
func biMapOp_toString_0(b *jBiMapCore) string            { return mapCoreString(b.fwd) }
func biMapOp_isEmpty_0(b *jBiMapCore) bool               { return b.fwd.IsEmpty() }
func biMapOp_containsKey_1(b *jBiMapCore, k any) bool    { return b.fwd.ContainsKey(k) }
func biMapOp_containsValue_1(b *jBiMapCore, v any) bool  { return b.rev.ContainsKey(v) }
func biMapOp_get_1(b *jBiMapCore, k any) any             { return b.fwd.Get(k) }
func biMapOp_getOrDefault_2(b *jBiMapCore, k, d any) any { return b.fwd.GetOrDefault(k, d) }
func biMapOp_inverse_0(b *jBiMapCore) any                { return b.inverse }
func biMapOp_put_2(b *jBiMapCore, k, v any) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	k, v = ToAny(k), ToAny(v)
	if b.fwd.ContainsKey(k) && javalang.ObjectsEquals(b.fwd.Get(k), v) {
		return v
	}
	if b.rev.ContainsKey(v) {
		panic(javalang.NewIllegalArgumentException("value already present: " + StrVal(v)))
	}
	if b.fwd.ContainsKey(k) {
		oldV := b.fwd.Get(k)
		b.rev.Remove(oldV)
	}
	old := b.fwd.Put(k, v)
	b.rev.Put(v, k)
	return old
}
func biMapOp_putIfAbsent_2(b *jBiMapCore, k, v any) any {
	if b.fwd.ContainsKey(k) && b.fwd.Get(k) != nil {
		return b.fwd.Get(k)
	}
	return biMapOp_put_2(b, k, v)
}
func biMapPutStr(b *jBiMapCore, k string, v any) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	v = ToAny(v)
	if b.fwd.ContainsStr(k) && javalang.ObjectsEquals(b.fwd.GetStr(k), v) {
		return v
	}
	if b.rev.ContainsKey(v) {
		panic(javalang.NewIllegalArgumentException("value already present: " + StrVal(v)))
	}
	if b.fwd.ContainsStr(k) {
		oldV := b.fwd.GetStr(k)
		if os, ok := oldV.(string); ok {
			b.rev.RemoveStr(os)
		} else {
			b.rev.Remove(oldV)
		}
	}
	old := b.fwd.PutStr(k, v)
	bk := BoxString(k)
	if vs, ok := v.(string); ok {
		b.rev.PutStr(vs, bk)
	} else {
		b.rev.Put(v, bk)
	}
	return old
}
func biMapRemoveStr(b *jBiMapCore, k string) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if !b.fwd.ContainsStr(k) {
		return nil
	}
	oldV := b.fwd.RemoveStr(k)
	if os, ok := oldV.(string); ok {
		b.rev.RemoveStr(os)
	} else {
		b.rev.Remove(oldV)
	}
	return oldV
}
func biMapOp_remove_1(b *jBiMapCore, k any) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if !b.fwd.ContainsKey(k) {
		return nil
	}
	oldV := b.fwd.Remove(k)
	b.rev.Remove(oldV)
	return oldV
}
func biMapOp_clear_0(b *jBiMapCore) {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b.fwd.Clear()
	b.rev.Clear()
}
func biMapOp_putAll_1(b *jBiMapCore, src Ext_java_util_Map) {
	CheckNotNull(src)
	for it := Iterate(src.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		biMapOp_put_2(b, e.GetKey(), e.GetValue())
	}
}
func biMapOp_keySet_0(b *jBiMapCore) any {
	if b.immutable {
		if b.fwd.IsEmpty() {
			return emptyImmutableSet_
		}
		if s := b.keySetCache.Load(); s != nil {
			return s
		}
		s := wrapImmutableSetCore(b.fwd).(*Ext_com_google_common_collect_ImmutableSetBase)
		if b.keySetCache.CompareAndSwap(nil, s) {
			return s
		}
		return b.keySetCache.Load()
	}
	return newLiveMapKeySet(b.fwd)
}
func biMapOp_values_0(b *jBiMapCore) any {
	if b.immutable {
		if b.fwd.IsEmpty() {
			return emptyImmutableSet_
		}
		if s := b.valuesCache.Load(); s != nil {
			return s
		}
		s := wrapImmutableSet(b.fwd.LiveValues()).(*Ext_com_google_common_collect_ImmutableSetBase)
		if b.valuesCache.CompareAndSwap(nil, s) {
			return s
		}
		return b.valuesCache.Load()
	}
	return newLiveMapKeySet(b.rev)
}
func biMapOp_entrySet_0(b *jBiMapCore) any {
	if b.immutable {
		if b.fwd.IsEmpty() {
			return emptyImmutableSet_
		}
		if s := b.entrySetCache.Load(); s != nil {
			return s
		}
		s := wrapImmutableMapEntrySet(b.fwd).(*Ext_com_google_common_collect_ImmutableSetBase)
		if b.entrySetCache.CompareAndSwap(nil, s) {
			return s
		}
		return b.entrySetCache.Load()
	}
	return newLiveMapEntrySet(b.fwd)
}
func biMapOp_equals_1(b *jBiMapCore, o any) bool { return mapCoreEquals(b.fwd, o) }
func biMapOp_hashCode_0(b *jBiMapCore) int32     { return mapCoreHash(b.fwd) }
func biMapOp_forEach_1(b *jBiMapCore, fn Ext_java_util_function_BiConsumer) {
	mapCoreForEach(b.fwd, fn)
}
func biMapOp_replaceAll_1(b *jBiMapCore, fn Ext_java_util_function_BiFunction) {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	for _, k := range b.fwd.LiveKeys() {
		oldV := b.fwd.Get(k)
		biMapOp_put_2(b, k, fn.Apply(k, oldV))
	}
}
func biMapOp_replace_3(b *jBiMapCore, k, ov, nv any) bool {
	if !b.fwd.ContainsKey(k) || !javalang.ObjectsEquals(b.fwd.Get(k), ToAny(ov)) {
		return false
	}
	biMapOp_put_2(b, k, nv)
	return true
}
func biMapOp_computeIfAbsent_2(b *jBiMapCore, k any, fn Ext_java_util_function_Function) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(fn)
	if v := b.fwd.Get(k); v != nil {
		return v
	}
	nv := ToAny(fn.Apply(ToAny(k)))
	if nv != nil {
		biMapOp_put_2(b, k, nv)
	}
	return nv
}
func biMapComputeIfAbsentStr(b *jBiMapCore, k string, fn Ext_java_util_function_Function) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(fn)
	if v := b.fwd.GetStr(k); v != nil {
		return v
	}
	bk := BoxString(k)
	nv := ToAny(fn.Apply(bk))
	if nv != nil {
		biMapOp_put_2(b, bk, nv)
	}
	return nv
}
func biMapComputeIfAbsentNode(b *jBiMapCore, n uint32, fn Ext_java_util_function_Function) any {
	if b.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(fn)
	if v := b.fwd.GetNode(n); v != nil {
		return v
	}
	bn := boxNodeById(n)
	nv := ToAny(fn.Apply(bn))
	if nv != nil {
		biMapOp_put_2(b, bn, nv)
	}
	return nv
}
func biMapOp_compute_2(b *jBiMapCore, k any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	oldV := b.fwd.Get(k)
	nv := ToAny(fn.Apply(ToAny(k), oldV))
	if nv == nil {
		biMapOp_remove_1(b, k)
		return nil
	}
	biMapOp_put_2(b, k, nv)
	return nv
}
func biMapOp_merge_3(b *jBiMapCore, k, v any, fn Ext_java_util_function_BiFunction) any {
	CheckNotNull(fn)
	v = CheckNotNull(ToAny(v))
	oldV := b.fwd.Get(k)
	var nv any
	if oldV == nil {
		nv = v
	} else {
		nv = ToAny(fn.Apply(oldV, v))
	}
	if nv == nil {
		biMapOp_remove_1(b, k)
	} else {
		biMapOp_put_2(b, k, nv)
	}
	return nv
}

// --- Multimap (ListMultimap, SetMultimap, SortedSetMultimap, Immutable*Multimap) ---

type jMultimapCore struct {
	isSet     bool
	immutable bool
	natVal    bool
	keyCmp    func(a, b any) int32
	valCmp    func(a, b any) int32
	keys      *JMapCore // key -> *JListCore (if !isSet) or *JMapCore (if isSet)
	total     int32
}

func newMultimapCore(isSet, immutable bool, keyCmp, valCmp func(a, b any) int32) *jMultimapCore {
	var km *JMapCore
	if keyCmp != nil {
		km = NewSortedJMapCore(keyCmp)
	} else {
		km = NewJMapCore(false)
	}
	return &jMultimapCore{isSet: isSet, immutable: immutable, keyCmp: keyCmp, valCmp: valCmp, keys: km}
}

func newTreeMultimapCore() *jMultimapCore {
	return &jMultimapCore{
		isSet:  true,
		natVal: true,
		keyCmp: CompareValues,
		valCmp: CompareValues,
		keys:   NewSortedJMapCore(nil),
	}
}

func (mm *jMultimapCore) newBucket() any {
	if mm.isSet {
		if mm.natVal {
			return NewSortedJMapCore(nil)
		}
		if mm.valCmp != nil {
			return NewSortedJMapCore(mm.valCmp)
		}
		return NewJMapCore(false)
	}
	return NewJListCore()
}

type jMultimapNilElem struct{}

var jMultimapNilElem_ any = jMultimapNilElem{}

func mmWrapSingle(v any) any {
	if v == nil {
		return jMultimapNilElem_
	}
	return v
}

func mmUnwrapSingle(b any) any {
	if _, isNil := b.(jMultimapNilElem); isNil {
		return nil
	}
	return b
}

func (mm *jMultimapCore) sealImmutableBuckets() {
	mm.keys.TransformValues(func(_, b any) any {
		if mm.isSet {
			return wrapImmutableSetCore(b.(*JMapCore))
		}
		if lc, ok := b.(*JListCore); ok {
			return wrapImmutableListOwn(lc.Items)
		}
		return newImmList1(mmUnwrapSingle(b))
	})
	mm.keys.sealImmutable()
	mm.immutable = true
}

func multimapOp_size_0(mm *jMultimapCore) int32              { return mm.total }
func multimapOp_isEmpty_0(mm *jMultimapCore) bool            { return mm.total == 0 }
func multimapOp_containsKey_1(mm *jMultimapCore, k any) bool { return mm.keys.ContainsKey(k) }
func (mm *jMultimapCore) containsKeyStr(k string) bool       { return mm.keys.ContainsStr(k) }
func (mm *jMultimapCore) containsKeyNode(u uint32) bool      { return mm.keys.ContainsNode(u) }
func multimapBucketContains(mm *jMultimapCore, b any, v any) bool {
	if b == nil {
		return false
	}
	if mm.immutable {
		if mm.isSet {
			return b.(Ext_com_google_common_collect_ImmutableSet).Contains(v)
		}
		return b.(Ext_com_google_common_collect_ImmutableList).Contains(v)
	}
	if mm.isSet {
		return b.(*JMapCore).ContainsKey(v)
	}
	if lc, ok := b.(*JListCore); ok {
		return lc.Contains(v)
	}
	return jSetElemEquals(mmUnwrapSingle(b), v)
}
func multimapOp_containsEntry_2(mm *jMultimapCore, k, v any) bool {
	return multimapBucketContains(mm, mm.keys.Get(k), v)
}
func (mm *jMultimapCore) containsEntryStr(k string, v any) bool {
	return multimapBucketContains(mm, mm.keys.GetStr(k), v)
}
func (mm *jMultimapCore) containsEntryNode(u uint32, v any) bool {
	return multimapBucketContains(mm, mm.keys.GetNode(u), v)
}
func multimapOp_put_2(mm *jMultimapCore, k, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.putInternal(k, v)
}
func (mm *jMultimapCore) putSetBucketVal(b any, v any) bool {
	sm := b.(*JMapCore)
	if sm.ContainsKey(v) {
		return false
	}
	sm.Put(v, true)
	mm.total++
	return true
}
func (mm *jMultimapCore) putInternal(k, v any) bool {
	if s, ok := k.(string); ok {
		return mm.putStrInternal(s, v)
	}
	k, v = ToAny(k), ToAny(v)
	b := mm.keys.Get(k)
	if mm.isSet {
		if b == nil {
			b = mm.newBucket()
			mm.keys.Put(k, b)
		}
		return mm.putSetBucketVal(b, v)
	}
	if b == nil {
		mm.keys.Put(k, mmWrapSingle(v))
		mm.total++
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		lc.Add(v)
		mm.total++
		return true
	}
	lc := &JListCore{Items: make([]any, 2, 4)}
	lc.Items[0] = mmUnwrapSingle(b)
	lc.Items[1] = v
	mm.keys.Put(k, lc)
	mm.total++
	return true
}
func (mm *jMultimapCore) putStrInternal(k string, v any) bool {
	v = ToAny(v)
	b := mm.keys.GetStr(k)
	if mm.isSet {
		if b == nil {
			b = mm.newBucket()
			mm.keys.PutStr(k, b)
		}
		return mm.putSetBucketVal(b, v)
	}
	if b == nil {
		mm.keys.PutStr(k, mmWrapSingle(v))
		mm.total++
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		lc.Add(v)
		mm.total++
		return true
	}
	lc := &JListCore{Items: make([]any, 2, 4)}
	lc.Items[0] = mmUnwrapSingle(b)
	lc.Items[1] = v
	mm.keys.PutStr(k, lc)
	mm.total++
	return true
}
func (mm *jMultimapCore) putNodeInternal(u uint32, v any) bool {
	v = ToAny(v)
	b := mm.keys.GetNode(u)
	if mm.isSet {
		if b == nil {
			b = mm.newBucket()
			mm.keys.PutNode(u, b)
		}
		return mm.putSetBucketVal(b, v)
	}
	if b == nil {
		mm.keys.PutNode(u, mmWrapSingle(v))
		mm.total++
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		lc.Add(v)
		mm.total++
		return true
	}
	lc := &JListCore{Items: make([]any, 2, 4)}
	lc.Items[0] = mmUnwrapSingle(b)
	lc.Items[1] = v
	mm.keys.PutNode(u, lc)
	mm.total++
	return true
}
func (mm *jMultimapCore) putStr(k string, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.putStrInternal(k, v)
}
func (mm *jMultimapCore) putNode(u uint32, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.putNodeInternal(u, v)
}
func multimapOp_remove_2(mm *jMultimapCore, k, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b := mm.keys.Get(k)
	if b == nil {
		return false
	}
	if mm.isSet {
		sm := b.(*JMapCore)
		if !sm.ContainsKey(v) {
			return false
		}
		sm.Remove(v)
		mm.total--
		if sm.IsEmpty() {
			mm.keys.Remove(k)
		}
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		if !lc.RemoveVal(v) {
			return false
		}
		mm.total--
		if lc.IsEmpty() {
			mm.keys.Remove(k)
		}
		return true
	}
	if !jSetElemEquals(mmUnwrapSingle(b), v) {
		return false
	}
	mm.keys.Remove(k)
	mm.total--
	return true
}
func (mm *jMultimapCore) removeStr(k string, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b := mm.keys.GetStr(k)
	if b == nil {
		return false
	}
	if mm.isSet {
		sm := b.(*JMapCore)
		if !sm.ContainsKey(v) {
			return false
		}
		sm.Remove(v)
		mm.total--
		if sm.IsEmpty() {
			mm.keys.RemoveStr(k)
		}
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		if !lc.RemoveVal(v) {
			return false
		}
		mm.total--
		if lc.IsEmpty() {
			mm.keys.RemoveStr(k)
		}
		return true
	}
	if !jSetElemEquals(mmUnwrapSingle(b), v) {
		return false
	}
	mm.keys.RemoveStr(k)
	mm.total--
	return true
}
func (mm *jMultimapCore) removeNode(u uint32, v any) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b := mm.keys.GetNode(u)
	if b == nil {
		return false
	}
	if mm.isSet {
		sm := b.(*JMapCore)
		if !sm.ContainsKey(v) {
			return false
		}
		sm.Remove(v)
		mm.total--
		if sm.IsEmpty() {
			mm.keys.RemoveNode(u)
		}
		return true
	}
	if lc, ok := b.(*JListCore); ok {
		if !lc.RemoveVal(v) {
			return false
		}
		mm.total--
		if lc.IsEmpty() {
			mm.keys.RemoveNode(u)
		}
		return true
	}
	if !jSetElemEquals(mmUnwrapSingle(b), v) {
		return false
	}
	mm.keys.RemoveNode(u)
	mm.total--
	return true
}
func multimapOp_putAll_2(mm *jMultimapCore, k any, vals Ext_java_lang_Iterable) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(vals)
	changed := false
	for it := NewIterCursor(vals); it.HasNext(); {
		if mm.putInternal(k, it.Next()) {
			changed = true
		}
	}
	return changed
}
func (mm *jMultimapCore) putAllStr(k string, vals Ext_java_lang_Iterable) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(vals)
	changed := false
	for it := NewIterCursor(vals); it.HasNext(); {
		if mm.putStrInternal(k, it.Next()) {
			changed = true
		}
	}
	return changed
}
func (mm *jMultimapCore) putAllNode(u uint32, vals Ext_java_lang_Iterable) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(vals)
	changed := false
	for it := NewIterCursor(vals); it.HasNext(); {
		if mm.putNodeInternal(u, it.Next()) {
			changed = true
		}
	}
	return changed
}
func multimapOp_putAll_1(mm *jMultimapCore, other Ext_com_google_common_collect_Multimap) bool {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(other)
	changed := false
	asMap := Cast[Ext_java_util_Map](other.AsMap())
	for it := Iterate(asMap.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		col := Cast[Ext_java_lang_Iterable](e.GetValue())
		if multimapOp_putAll_2(mm, e.GetKey(), col) {
			changed = true
		}
	}
	return changed
}
func (mm *jMultimapCore) finishRemoveAllBucket(b any) any {
	if b == nil {
		if mm.isSet {
			if mm.natVal {
				ts := &Ext_java_util_TreeSet{}
				ts.BindExtSelf_(ts)
				ts.Core_ = NewSortedJMapCore(nil)
				return ts
			}
			if mm.valCmp != nil {
				ts := &Ext_java_util_TreeSet{}
				ts.BindExtSelf_(ts)
				ts.Core_ = NewSortedJMapCore(mm.valCmp)
				return ts
			}
			return wrapLinkedHashSet(nil)
		}
		return wrapArrayList(nil)
	}
	if mm.isSet {
		sm := b.(*JMapCore)
		mm.total -= sm.Size()
		if mm.valCmp != nil {
			ts := &Ext_java_util_TreeSet{}
			ts.BindExtSelf_(ts)
			ts.Core_ = sm
			return ts
		}
		obj := &Ext_java_util_LinkedHashSet{}
		obj.BindExtSelf_(obj)
		obj.inlineSet.InitExpectedSize(sm.Size())
		if sm.mode == jMapModeStr {
			for i := int(sm.firstIdx); i < len(sm.strEntries); i++ {
				if k := sm.strEntries[i].key; k != jTombstoneStr {
					obj.inlineSet.AddStr(k)
				}
			}
		} else {
			sm.ForEachEntry(func(k, _ any) bool {
				obj.inlineSet.Add(k)
				return true
			})
		}
		return obj
	}
	if lc, ok := b.(*JListCore); ok {
		mm.total -= lc.Size()
		return wrapArrayList(lc.Items)
	}
	mm.total--
	return wrapArrayListOwn([]any{mmUnwrapSingle(b)})
}
func multimapOp_removeAll_1(mm *jMultimapCore, k any) any {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.finishRemoveAllBucket(mm.keys.Remove(k))
}
func (mm *jMultimapCore) removeAllStr(k string) any {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.finishRemoveAllBucket(mm.keys.RemoveStr(k))
}
func (mm *jMultimapCore) removeAllNode(u uint32) any {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	return mm.finishRemoveAllBucket(mm.keys.RemoveNode(u))
}
func multimapOp_clear_0(mm *jMultimapCore) {
	if mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	mm.keys.Clear()
	mm.total = 0
}
func (mm *jMultimapCore) finishImmGet(b any) any {
	if b != nil {
		return b
	}
	if mm.isSet {
		return emptyImmutableSet_
	}
	return emptyImmutableList_
}
func multimapOp_get_1(mm *jMultimapCore, k any) any {
	if s, ok := k.(string); ok {
		return mm.getStr(s)
	}
	k = ToAny(k)
	if mm.immutable {
		return mm.finishImmGet(mm.keys.Get(k))
	}
	if mm.isSet {
		v := &jMultimapLiveSet{mm: mm, key: k}
		v.BindExtSelf_(v)
		return v
	}
	v := &jMultimapLiveList{mm: mm, key: k}
	v.BindExtSelf_(v)
	return v
}
func (mm *jMultimapCore) getStr(k string) any {
	if mm.immutable {
		return mm.finishImmGet(mm.keys.GetStr(k))
	}
	if mm.isSet {
		v := &jMultimapLiveSet{mm: mm, strKey: k, keyKind: 1}
		v.BindExtSelf_(v)
		return v
	}
	v := &jMultimapLiveList{mm: mm, strKey: k, keyKind: 1}
	v.BindExtSelf_(v)
	return v
}
func (mm *jMultimapCore) getNode(u uint32) any {
	if mm.immutable {
		return mm.finishImmGet(mm.keys.GetNode(u))
	}
	if mm.isSet {
		v := &jMultimapLiveSet{mm: mm, nodeKey: u, keyKind: 2}
		v.BindExtSelf_(v)
		return v
	}
	v := &jMultimapLiveList{mm: mm, nodeKey: u, keyKind: 2}
	v.BindExtSelf_(v)
	return v
}
func multimapOp_keySet_0(mm *jMultimapCore) any {
	if mm.immutable {
		return wrapImmutableSetCore(mm.keys)
	}
	if mm.keyCmp != nil {
		ts := &Ext_java_util_TreeSet{}
		ts.Core_ = cloneJMapCore(mm.keys)
		return ts
	}
	v := &jMultimapLiveKeySet{mm: mm}
	v.BindExtSelf_(v)
	return v
}
func multimapOp_values_0(mm *jMultimapCore) any {
	out := make([]any, 0, mm.total)
	if mm.immutable {
		mm.keys.ForEachValue(func(b any) bool {
			for it := NewIterCursor(b); it.HasNext(); {
				out = append(out, it.Next())
			}
			return true
		})
		return wrapImmutableListOwn(out)
	}
	mm.keys.ForEachValue(func(b any) bool {
		if mm.isSet {
			b.(*JMapCore).ForEachEntry(func(v, _ any) bool {
				out = append(out, v)
				return true
			})
		} else if lc, ok := b.(*JListCore); ok {
			out = append(out, lc.Items...)
		} else {
			out = append(out, mmUnwrapSingle(b))
		}
		return true
	})
	return wrapArrayList(out)
}
func multimapOp_asMap_0(mm *jMultimapCore) any {
	if mm.immutable {
		return wrapImmutableMap(mm.keys)
	}
	v := &jMultimapAsMap{mm: mm}
	v.BindExtSelf_(v)
	return v
}

func multimapInverse(mm *jMultimapCore, isSet bool) any {
	inv := newMultimapCore(isSet, false, nil, nil)
	if mm.immutable {
		mm.keys.ForEachEntry(func(k, b any) bool {
			for it := NewIterCursor(b); it.HasNext(); {
				inv.putInternal(it.Next(), k)
			}
			return true
		})
	} else {
		mm.keys.ForEachEntry(func(k, b any) bool {
			if mm.isSet {
				b.(*JMapCore).ForEachEntry(func(v, _ any) bool {
					inv.putInternal(v, k)
					return true
				})
			} else if lc, ok := b.(*JListCore); ok {
				for _, v := range lc.Items {
					inv.putInternal(v, k)
				}
			} else {
				inv.putInternal(mmUnwrapSingle(b), k)
			}
			return true
		})
	}
	inv.sealImmutableBuckets()
	if isSet {
		obj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
		obj.Core_ = inv
		return obj
	}
	obj := &Ext_com_google_common_collect_ImmutableListMultimap{}
	obj.Core_ = inv
	return obj
}

type jMultimapLiveList struct {
	Ext_java_util_AbstractListBase
	mm      *jMultimapCore
	key     any
	strKey  string
	nodeKey uint32
	keyKind uint8
}

func (l *jMultimapLiveList) getBucket() any {
	switch l.keyKind {
	case 1:
		return l.mm.keys.GetStr(l.strKey)
	case 2:
		return l.mm.keys.GetNode(l.nodeKey)
	default:
		return l.mm.keys.Get(l.key)
	}
}
func (l *jMultimapLiveList) putBucket(b any) {
	switch l.keyKind {
	case 1:
		l.mm.keys.PutStr(l.strKey, b)
	case 2:
		l.mm.keys.PutNode(l.nodeKey, b)
	default:
		l.mm.keys.Put(l.key, b)
	}
}
func (l *jMultimapLiveList) removeBucket() {
	switch l.keyKind {
	case 1:
		l.mm.keys.RemoveStr(l.strKey)
	case 2:
		l.mm.keys.RemoveNode(l.nodeKey)
	default:
		l.mm.keys.Remove(l.key)
	}
}

func (l *jMultimapLiveList) typedViewLen_() (int, bool) {
	if l == nil {
		ThrowNPE()
	}
	return int(l.Size()), true
}
func (l *jMultimapLiveList) typedViewAt_(i int) any            { return l.Get(int32(i)) }
func (*jMultimapLiveList) typedViewNodeAt_(int) (uint32, bool) { return 0, false }
func (l *jMultimapLiveList) Size() int32 {
	b := l.getBucket()
	if b == nil {
		return 0
	}
	if lc, ok := b.(*JListCore); ok {
		return lc.Size()
	}
	return 1
}
func (l *jMultimapLiveList) Get(i int32) any {
	b := l.getBucket()
	if b == nil {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0))
	}
	if lc, ok := b.(*JListCore); ok {
		return lc.Get(i)
	}
	CheckElementIndex(i, 1)
	return mmUnwrapSingle(b)
}
func (l *jMultimapLiveList) Set(i int32, v any) any {
	b := l.getBucket()
	if b == nil {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0))
	}
	if lc, ok := b.(*JListCore); ok {
		return lc.Set(i, v)
	}
	CheckElementIndex(i, 1)
	old := mmUnwrapSingle(b)
	l.putBucket(mmWrapSingle(ToAny(v)))
	return old
}
func (l *jMultimapLiveList) Add_2(i int32, v any) {
	b := l.getBucket()
	if b == nil {
		if i != 0 {
			panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0))
		}
		l.putBucket(mmWrapSingle(ToAny(v)))
		l.mm.total++
		return
	}
	if lc, ok := b.(*JListCore); ok {
		lc.AddAt(i, v)
		l.mm.total++
		return
	}
	CheckPositionIndex(i, 1)
	v0 := mmUnwrapSingle(b)
	v1 := ToAny(v)
	lc := &JListCore{Items: make([]any, 2, 4)}
	if i == 0 {
		lc.Items[0], lc.Items[1] = v1, v0
	} else {
		lc.Items[0], lc.Items[1] = v0, v1
	}
	l.putBucket(lc)
	l.mm.total++
}
func (l *jMultimapLiveList) Remove_1_int32(i int32) any {
	b := l.getBucket()
	if b == nil {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0))
	}
	if lc, ok := b.(*JListCore); ok {
		old := lc.RemoveAt(i)
		l.mm.total--
		if lc.IsEmpty() {
			l.removeBucket()
		}
		return old
	}
	CheckElementIndex(i, 1)
	old := mmUnwrapSingle(b)
	l.removeBucket()
	l.mm.total--
	return old
}

type jMultimapLiveSet struct {
	Ext_java_util_AbstractSetBase
	mm      *jMultimapCore
	key     any
	strKey  string
	nodeKey uint32
	keyKind uint8
}

func (s *jMultimapLiveSet) getBucket() any {
	switch s.keyKind {
	case 1:
		return s.mm.keys.GetStr(s.strKey)
	case 2:
		return s.mm.keys.GetNode(s.nodeKey)
	default:
		return s.mm.keys.Get(s.key)
	}
}
func (s *jMultimapLiveSet) putVal(v any) bool {
	switch s.keyKind {
	case 1:
		return s.mm.putStr(s.strKey, v)
	case 2:
		return s.mm.putNode(s.nodeKey, v)
	default:
		return multimapOp_put_2(s.mm, s.key, v)
	}
}
func (s *jMultimapLiveSet) remVal(v any) bool {
	switch s.keyKind {
	case 1:
		return s.mm.removeStr(s.strKey, v)
	case 2:
		return s.mm.removeNode(s.nodeKey, v)
	default:
		return multimapOp_remove_2(s.mm, s.key, v)
	}
}

func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SequencedCollection() {}
func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SequencedSet()        {}
func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SortedSet()           {}
func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_NavigableSet()        {}
func (s *jMultimapLiveSet) Size() int32 {
	b := s.getBucket()
	if b == nil {
		return 0
	}
	return b.(*JMapCore).Size()
}
func (s *jMultimapLiveSet) Contains(v any) bool {
	b := s.getBucket()
	return b != nil && b.(*JMapCore).ContainsKey(v)
}
func (s *jMultimapLiveSet) SetContainsStr_(k string) bool {
	b := s.getBucket()
	return b != nil && b.(*JMapCore).ContainsStr(k)
}
func (s *jMultimapLiveSet) SetAddStr_(k string) bool {
	if s.mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b := s.getBucket()
	if b == nil {
		sm := s.mm.newBucket().(*JMapCore)
		sm.PutStr(k, true)
		switch s.keyKind {
		case 1:
			s.mm.keys.PutStr(s.strKey, sm)
		case 2:
			s.mm.keys.PutNode(s.nodeKey, sm)
		default:
			s.mm.keys.Put(s.key, sm)
		}
		s.mm.total++
		return true
	}
	sm := b.(*JMapCore)
	if sm.ContainsStr(k) {
		return false
	}
	sm.PutStr(k, true)
	s.mm.total++
	return true
}
func (s *jMultimapLiveSet) SetRemoveStr_(k string) bool {
	if s.mm.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	b := s.getBucket()
	if b == nil {
		return false
	}
	sm := b.(*JMapCore)
	if !sm.ContainsStr(k) {
		return false
	}
	sm.RemoveStr(k)
	s.mm.total--
	if sm.IsEmpty() {
		switch s.keyKind {
		case 1:
			s.mm.keys.RemoveStr(s.strKey)
		case 2:
			s.mm.keys.RemoveNode(s.nodeKey)
		default:
			s.mm.keys.Remove(s.key)
		}
	}
	return true
}
func (s *jMultimapLiveSet) Add_1(v any) bool        { return s.putVal(v) }
func (s *jMultimapLiveSet) Remove_1_any(v any) bool { return s.remVal(v) }
func (s *jMultimapLiveSet) Clear() {
	switch s.keyKind {
	case 1:
		_ = s.mm.removeAllStr(s.strKey)
	case 2:
		_ = s.mm.removeAllNode(s.nodeKey)
	default:
		_ = multimapOp_removeAll_1(s.mm, s.key)
	}
}
func (s *jMultimapLiveSet) AddFirst(v any) { _ = s.putVal(v) }
func (s *jMultimapLiveSet) AddLast(v any)  { _ = s.putVal(v) }
func (s *jMultimapLiveSet) GetFirst() any {
	b := s.getBucket()
	if b == nil {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return b.(*JMapCore).FirstKey()
}
func (s *jMultimapLiveSet) GetLast() any {
	b := s.getBucket()
	if b == nil {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return b.(*JMapCore).LastKey()
}
func (s *jMultimapLiveSet) RemoveFirst() any {
	x := s.GetFirst()
	s.Remove_1_any(x)
	return x
}
func (s *jMultimapLiveSet) RemoveLast() any {
	x := s.GetLast()
	s.Remove_1_any(x)
	return x
}
func (s *jMultimapLiveSet) DescendingSet() any {
	b := s.getBucket()
	if b == nil {
		return setCoreDescending(NewSortedJMapCore(nil), false)
	}
	return setCoreDescending(b.(*JMapCore), false)
}
func (s *jMultimapLiveSet) TailSet_2(from any, incl bool) any {
	b := s.getBucket()
	if b == nil {
		return setCoreTailSet(NewSortedJMapCore(nil), from, incl, false)
	}
	return setCoreTailSet(b.(*JMapCore), from, incl, false)
}
func (s *jMultimapLiveSet) fastSetCore_() *JMapCore {
	if s == nil {
		return nil
	}
	b := s.getBucket()
	if b == nil {
		return nil
	}
	return b.(*JMapCore)
}
func (s *jMultimapLiveSet) Iterator() any {
	b := s.getBucket()
	if b == nil {
		return emptyMapIter_
	}
	mc := b.(*JMapCore)
	if mc.Live == 0 {
		return emptyMapIter_
	}
	if strs, ok := mc.LiveStrKeys(); ok {
		return &jMultimapSetIter{mm: s.mm, key: s.key, strKey: s.strKey, nodeKey: s.nodeKey, keyKind: s.keyKind, strs: strs, isStr: true}
	}
	return &jMultimapSetIter{mm: s.mm, key: s.key, strKey: s.strKey, nodeKey: s.nodeKey, keyKind: s.keyKind, vals: mc.LiveKeys()}
}

type jMultimapSetIter struct {
	mm      *jMultimapCore
	key     any
	strKey  string
	nodeKey uint32
	vals    []any
	strs    []string
	idx     int
	lastVal any
	lastStr string
	keyKind uint8
	isStr   bool
	hasLast bool
}

func (*jMultimapSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jMultimapSetIter) IterStrSlice_() ([]string, bool) {
	if it != nil && it.isStr && it.idx == 0 {
		return it.strs, true
	}
	return nil, false
}
func (it *jMultimapSetIter) HasNext() bool {
	if it.isStr {
		return it.idx < len(it.strs)
	}
	return it.idx < len(it.vals)
}
func (it *jMultimapSetIter) Next() any {
	if it.isStr {
		if it.idx >= len(it.strs) {
			panic(javalang.NewNoSuchElementException("exhausted"))
		}
		s := it.strs[it.idx]
		it.idx++
		it.lastStr = s
		it.hasLast = true
		return BoxString(s)
	}
	if it.idx >= len(it.vals) {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	v := it.vals[it.idx]
	it.idx++
	it.lastVal = v
	it.hasLast = true
	return v
}
func (it *jMultimapSetIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	lv := it.lastVal
	if it.isStr {
		lv = BoxString(it.lastStr)
	}
	switch it.keyKind {
	case 1:
		it.mm.removeStr(it.strKey, lv)
	case 2:
		it.mm.removeNode(it.nodeKey, lv)
	default:
		multimapOp_remove_2(it.mm, it.key, lv)
	}
	it.hasLast = false
}

type jMultimapLiveKeySet struct {
	Ext_java_util_AbstractSetBase
	mm *jMultimapCore
}

func (s *jMultimapLiveKeySet) fastSetCore_() *JMapCore {
	if s == nil {
		return nil
	}
	return s.mm.keys
}
func (s *jMultimapLiveKeySet) Size() int32                   { return s.mm.keys.Size() }
func (s *jMultimapLiveKeySet) Contains(k any) bool           { return s.mm.keys.ContainsKey(k) }
func (s *jMultimapLiveKeySet) SetContainsStr_(k string) bool { return s.mm.keys.ContainsStr(k) }
func (s *jMultimapLiveKeySet) SetAddStr_(_ string) bool {
	panic(javalang.NewUnsupportedOperationException("keySet.add"))
}
func (s *jMultimapLiveKeySet) SetRemoveStr_(k string) bool {
	if !s.mm.keys.ContainsStr(k) {
		return false
	}
	_ = s.mm.removeAllStr(k)
	return true
}
func (s *jMultimapLiveKeySet) Clear() { multimapOp_clear_0(s.mm) }
func (s *jMultimapLiveKeySet) Remove_1_any(k any) bool {
	if !s.mm.keys.ContainsKey(k) {
		return false
	}
	_ = multimapOp_removeAll_1(s.mm, k)
	return true
}
func (s *jMultimapLiveKeySet) Iterator() any {
	if s.mm.keys.Live == 0 {
		return emptyMapIter_
	}
	if strs, ok := s.mm.keys.LiveStrKeys(); ok {
		return &jMultimapKeyIter{mm: s.mm, strs: strs, isStr: true}
	}
	return &jMultimapKeyIter{mm: s.mm, keys: s.mm.keys.LiveKeys()}
}

type jMultimapKeyIter struct {
	mm      *jMultimapCore
	keys    []any
	strs    []string
	idx     int
	lastKey any
	lastStr string
	isStr   bool
	hasLast bool
}

func (*jMultimapKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jMultimapKeyIter) IterStrSlice_() ([]string, bool) {
	if it != nil && it.isStr && it.idx == 0 {
		return it.strs, true
	}
	return nil, false
}
func (it *jMultimapKeyIter) HasNext() bool {
	if it.isStr {
		return it.idx < len(it.strs)
	}
	return it.idx < len(it.keys)
}
func (it *jMultimapKeyIter) Next() any {
	if it.isStr {
		if it.idx >= len(it.strs) {
			panic(javalang.NewNoSuchElementException("exhausted"))
		}
		s := it.strs[it.idx]
		it.idx++
		it.lastStr = s
		it.hasLast = true
		return BoxString(s)
	}
	if it.idx >= len(it.keys) {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	k := it.keys[it.idx]
	it.idx++
	it.lastKey = k
	it.hasLast = true
	return k
}
func (it *jMultimapKeyIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	if it.isStr {
		_ = it.mm.removeAllStr(it.lastStr)
		it.hasLast = false
		return
	}
	_ = multimapOp_removeAll_1(it.mm, it.lastKey)
	it.hasLast = false
}

type jMultimapAsMap struct {
	Ext_java_util_AbstractMapBase
	mm *jMultimapCore
}

func (m *jMultimapAsMap) Size() int32            { return m.mm.keys.Size() }
func (m *jMultimapAsMap) ContainsKey(k any) bool { return m.mm.keys.ContainsKey(k) }
func (m *jMultimapAsMap) MapGetStr_(k string) any {
	if !m.mm.keys.ContainsStr(k) {
		return nil
	}
	return m.mm.getStr(k)
}
func (m *jMultimapAsMap) MapGetOrDefaultStr_(k string, def any) any {
	if !m.mm.keys.ContainsStr(k) {
		return ToAny(def)
	}
	return m.mm.getStr(k)
}
func (m *jMultimapAsMap) MapContainsStr_(k string) bool { return m.mm.keys.ContainsStr(k) }
func (m *jMultimapAsMap) MapPutStr_(_ string, _ any) any {
	panic(javalang.NewUnsupportedOperationException("asMap.put"))
}
func (m *jMultimapAsMap) MapRemoveStr_(k string) any {
	if !m.mm.keys.ContainsStr(k) {
		return nil
	}
	return m.mm.removeAllStr(k)
}
func (m *jMultimapAsMap) Get(k any) any {
	if !m.mm.keys.ContainsKey(k) {
		return nil
	}
	return multimapOp_get_1(m.mm, k)
}
func (m *jMultimapAsMap) Remove_1(k any) any {
	if !m.mm.keys.ContainsKey(k) {
		return nil
	}
	return multimapOp_removeAll_1(m.mm, k)
}
func (m *jMultimapAsMap) Clear() { multimapOp_clear_0(m.mm) }
func (m *jMultimapAsMap) EntrySet() any {
	s := &jMultimapAsMapEntries{mm: m.mm}
	s.BindExtSelf_(s)
	return s
}

type jMultimapAsMapEntries struct {
	Ext_java_util_AbstractSetBase
	mm *jMultimapCore
}

func (s *jMultimapAsMapEntries) Size() int32 { return s.mm.keys.Size() }
func (s *jMultimapAsMapEntries) Clear()      { multimapOp_clear_0(s.mm) }
func (s *jMultimapAsMapEntries) Iterator() any {
	if s.mm.keys.Live == 0 {
		return emptyMapIter_
	}
	items := make([]any, 0, s.mm.keys.Live)
	s.mm.keys.ForEachEntry(func(k, _ any) bool {
		items = append(items, &jSimpleEntry{k: k, v: multimapOp_get_1(s.mm, k), imm: true})
		return true
	})
	return &jMultimapAsMapIter{mm: s.mm, entries: items}
}

type jMultimapAsMapIter struct {
	mm      *jMultimapCore
	entries []any
	idx     int
	lastKey any
	hasLast bool
}

func (*jMultimapAsMapIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jMultimapAsMapIter) HasNext() bool                          { return it.idx < len(it.entries) }
func (it *jMultimapAsMapIter) Next() any {
	if it.idx >= len(it.entries) {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	e := it.entries[it.idx].(*jSimpleEntry)
	it.idx++
	it.lastKey = e.k
	it.hasLast = true
	return e
}
func (it *jMultimapAsMapIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	_ = multimapOp_removeAll_1(it.mm, it.lastKey)
	it.hasLast = false
}

type jMultimapBuilderEntry struct {
	strKey string
	key    any
	val    any
	isStr  bool
}

type jMultimapBuilderCore struct {
	isSet   bool
	entries []jMultimapBuilderEntry
	keyCmp  func(a, b any) int32
}

func newMultimapBuilderCore(isSet bool) *jMultimapBuilderCore {
	return &jMultimapBuilderCore{isSet: isSet}
}

func multimapBuilderPut(b *jMultimapBuilderCore, k, v any) {
	if s, ok := k.(string); ok {
		if s == NullString {
			ThrowNPE()
		}
		b.entries = append(b.entries, jMultimapBuilderEntry{strKey: s, val: CheckNotNull(ToAny(v)), isStr: true})
		return
	}
	b.entries = append(b.entries, jMultimapBuilderEntry{key: CheckNotNull(ToAny(k)), val: CheckNotNull(ToAny(v))})
}

func MultimapBuilderPutStr(b any, k string, v any) any {
	if b == nil || k == NullString {
		ThrowNPE()
	}
	v = CheckNotNull(ToAny(v))
	if cp, ok := b.(interface{ getCore() any }); ok {
		if isJavaNilAny(b) {
			ThrowNPE()
		}
		bc := cp.getCore().(*jMultimapBuilderCore)
		bc.entries = append(bc.entries, jMultimapBuilderEntry{strKey: k, val: v, isStr: true})
		return b
	}
	if isJavaNilAny(b) {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException(fmt.Sprintf("MultimapBuilderPutStr unsupported on %T", b)))
}

func multimapBuilderPutAll(b *jMultimapBuilderCore, k any, vals Ext_java_lang_Iterable) {
	if s, ok := k.(string); ok {
		if s == NullString {
			ThrowNPE()
		}
		CheckNotNull(vals)
		for it := NewIterCursor(vals); it.HasNext(); {
			b.entries = append(b.entries, jMultimapBuilderEntry{strKey: s, val: CheckNotNull(it.Next()), isStr: true})
		}
		return
	}
	k = CheckNotNull(ToAny(k))
	CheckNotNull(vals)
	for it := NewIterCursor(vals); it.HasNext(); {
		b.entries = append(b.entries, jMultimapBuilderEntry{key: k, val: CheckNotNull(it.Next())})
	}
}

func multimapBuilderBuild(b *jMultimapBuilderCore) any {
	mm := newMultimapCore(b.isSet, false, b.keyCmp, nil)
	for _, e := range b.entries {
		if e.isStr {
			mm.putStrInternal(e.strKey, e.val)
		} else {
			mm.putInternal(e.key, e.val)
		}
	}
	mm.sealImmutableBuckets()
	if b.isSet {
		obj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
		obj.Core_ = mm
		return obj
	}
	obj := &Ext_com_google_common_collect_ImmutableListMultimap{}
	obj.Core_ = mm
	return obj
}

func immSetMultimapCopyOf(src Ext_com_google_common_collect_Multimap) *Ext_com_google_common_collect_ImmutableSetMultimap {
	CheckNotNull(src)
	if ism, ok := any(src).(*Ext_com_google_common_collect_ImmutableSetMultimap); ok && ism != nil {
		return ism
	}
	mm := newMultimapCore(true, false, nil, nil)
	asMap := Cast[Ext_java_util_Map](src.AsMap())
	for it := Iterate(asMap.EntrySet()); it.HasNext(); {
		e := Cast[Ext_java_util_Map_Entry](it.Next())
		k := CheckNotNull(e.GetKey())
		for vit := NewIterCursor(e.GetValue()); vit.HasNext(); {
			mm.putInternal(k, CheckNotNull(vit.Next()))
		}
	}
	mm.sealImmutableBuckets()
	obj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
	obj.Core_ = mm
	return obj
}

func multimapsIndex(vals Ext_java_lang_Iterable, keyFn Ext_com_google_common_base_Function) *Ext_com_google_common_collect_ImmutableListMultimap {
	CheckNotNull(vals)
	CheckNotNull(keyFn)
	mm := newMultimapCore(false, false, nil, nil)
	for it := NewIterCursor(vals); it.HasNext(); {
		v := CheckNotNull(it.Next())
		k := CheckNotNull(ToAny(keyFn.Apply(v)))
		mm.putInternal(k, v)
	}
	mm.sealImmutableBuckets()
	obj := &Ext_com_google_common_collect_ImmutableListMultimap{}
	obj.Core_ = mm
	return obj
}

// --- Multiset ---

type jMultisetCore struct {
	counts    *JMapCore // elem -> int32
	total     int32
	immutable bool
	cmp       func(a, b any) int32
}

func newMultisetCore(immutable bool, cmp func(a, b any) int32) *jMultisetCore {
	var cm *JMapCore
	if cmp != nil {
		cm = NewSortedJMapCore(cmp)
	} else {
		cm = NewJMapCore(false)
	}
	return &jMultisetCore{counts: cm, immutable: immutable, cmp: cmp}
}

func newTreeMultisetCore() *jMultisetCore {
	return &jMultisetCore{counts: NewSortedJMapCore(nil), cmp: CompareValues}
}

func cloneMultisetCore(src *jMultisetCore, immutable bool) *jMultisetCore {
	dst := newMultisetCore(immutable, src.cmp)
	dst.counts = cloneJMapCore(src.counts)
	dst.total = src.total
	return dst
}

func wrapImmutableMultiset(c *jMultisetCore) Ext_com_google_common_collect_ImmutableMultiset {
	obj := &Ext_com_google_common_collect_ImmutableMultisetBase{}
	obj.BindExtSelf_(obj)
	obj.Core_ = c
	return obj
}

type jMultisetEntry struct {
	elem  any
	count int32
}

func (*jMultisetEntry) Implements_j2gort_Ext_com_google_common_collect_Multiset_Entry() {}
func (e *jMultisetEntry) GetElement() any                                               { return e.elem }
func (e *jMultisetEntry) GetCount() int32                                               { return e.count }
func (e *jMultisetEntry) Equals(o any) bool {
	other, ok := o.(Ext_com_google_common_collect_Multiset_Entry)
	return ok && other != nil && !javalang.IsNil(other) &&
		e.count == other.GetCount() && javalang.ObjectsEquals(e.elem, ToAny(other.GetElement()))
}
func (e *jMultisetEntry) HashCode() int32 {
	return javalang.ObjectHashCode(e.elem) ^ e.count
}
func (e *jMultisetEntry) String() string {
	if e.count == 1 {
		return StrVal(e.elem)
	}
	return fmt.Sprintf("%s x %d", StrVal(e.elem), e.count)
}

func multisetCountStr(c *jMultisetCore, s string) int32 {
	v := c.counts.GetStr(s)
	if v == nil {
		return 0
	}
	return v.(int32)
}
func multisetAddStr2(c *jMultisetCore, s string, n int32) int32 {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckArgument(n >= 0)
	old := multisetCountStr(c, s)
	if n > 0 {
		c.counts.PutStr(s, old+n)
		c.total += n
	}
	return old
}
func multisetRemoveStr1(c *jMultisetCore, s string) bool {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	old := multisetCountStr(c, s)
	if old == 0 {
		return false
	}
	if old == 1 {
		c.counts.RemoveStr(s)
	} else {
		c.counts.PutStr(s, old-1)
	}
	c.total--
	return true
}
func multisetSetCountStr3(c *jMultisetCore, s string, oldCount, newCount int32) bool {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckArgument(oldCount >= 0 && newCount >= 0)
	cur := multisetCountStr(c, s)
	if cur != oldCount {
		return false
	}
	if newCount == 0 {
		if cur > 0 {
			c.counts.RemoveStr(s)
		}
	} else {
		c.counts.PutStr(s, newCount)
	}
	c.total += newCount - cur
	return true
}
func multisetOp_size_0(c *jMultisetCore) int32   { return c.total }
func multisetOp_isEmpty_0(c *jMultisetCore) bool { return c.total == 0 }
func multisetOp_count_1(c *jMultisetCore, e any) int32 {
	v := c.counts.Get(e)
	if v == nil {
		return 0
	}
	return v.(int32)
}
func multisetOp_contains_1(c *jMultisetCore, e any) bool { return multisetOp_count_1(c, e) > 0 }
func multisetOp_add_2(c *jMultisetCore, e any, n int32) int32 {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckArgument(n >= 0)
	old := multisetOp_count_1(c, e)
	if n > 0 {
		c.counts.Put(e, old+n)
		c.total += n
	}
	return old
}
func multisetOp_add_1(c *jMultisetCore, e any) bool {
	multisetOp_add_2(c, e, 1)
	return true
}
func multisetOp_remove_1(c *jMultisetCore, e any) bool {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	old := multisetOp_count_1(c, e)
	if old == 0 {
		return false
	}
	if old == 1 {
		c.counts.Remove(e)
	} else {
		c.counts.Put(e, old-1)
	}
	c.total--
	return true
}
func multisetOp_setCount_3(c *jMultisetCore, e any, oldCount, newCount int32) bool {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckArgument(oldCount >= 0 && newCount >= 0)
	cur := multisetOp_count_1(c, e)
	if cur != oldCount {
		return false
	}
	if newCount == 0 {
		if cur > 0 {
			c.counts.Remove(e)
		}
	} else {
		c.counts.Put(e, newCount)
	}
	c.total += newCount - cur
	return true
}
func multisetOp_addAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
	CheckNotNull(col)
	changed := false
	for it := Iterate(col); it.HasNext(); {
		multisetOp_add_1(c, it.Next())
		changed = true
	}
	return changed
}
func multisetOp_removeAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	CheckNotNull(col)
	changed := false
	for it := Iterate(col); it.HasNext(); {
		x := it.Next()
		if cnt := multisetOp_count_1(c, x); cnt > 0 {
			c.counts.Remove(x)
			c.total -= cnt
			changed = true
		}
	}
	return changed
}
func multisetOp_clear_0(c *jMultisetCore) {
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	c.counts.Clear()
	c.total = 0
}
func multisetOp_retainAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
	CheckNotNull(col)
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	changed := false
	for _, x := range c.counts.LiveKeys() {
		if !col.Contains(x) {
			cnt := multisetOp_count_1(c, x)
			c.counts.Remove(x)
			c.total -= cnt
			changed = true
		}
	}
	return changed
}
func multisetOp_removeIf_1(c *jMultisetCore, pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	if c.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	changed := false
	for _, x := range c.counts.LiveKeys() {
		if pred.Test(x) {
			cnt := multisetOp_count_1(c, x)
			c.counts.Remove(x)
			c.total -= cnt
			changed = true
		}
	}
	return changed
}
func multisetOp_elementSet_0(c *jMultisetCore) any {
	if c.immutable {
		return wrapImmutableSet(c.counts.LiveKeys())
	}
	if c.cmp != nil {
		ts := &Ext_java_util_TreeSet{}
		ts.BindExtSelf_(ts)
		ts.Core_ = cloneJMapCore(c.counts)
		return ts
	}
	return newLiveMapKeySet(c.counts)
}
func multisetOp_entrySet_0(c *jMultisetCore) any {
	entries := make([]any, 0, c.counts.Size())
	c.counts.ForEachEntry(func(k, v any) bool {
		entries = append(entries, &jMultisetEntry{elem: k, count: v.(int32)})
		return true
	})
	if c.immutable {
		return wrapImmutableSet(entries)
	}
	return wrapLinkedHashSet(entries)
}
func multisetOp_toSlice(c *jMultisetCore) []any {
	out := make([]any, 0, c.total)
	c.counts.ForEachEntry(func(e, v any) bool {
		n := v.(int32)
		for j := int32(0); j < n; j++ {
			out = append(out, e)
		}
		return true
	})
	return out
}
func multisetOp_iterator_0(c *jMultisetCore) any {
	return &jMultisetIter{c: c, items: multisetOp_toSlice(c)}
}
func multisetOp_asList_0(c *jMultisetCore) Ext_com_google_common_collect_ImmutableList {
	return wrapImmutableList(multisetOp_toSlice(c))
}
func multisetOp_toArray_0(c *jMultisetCore) []any { return multisetOp_toSlice(c) }
func multisetOp_toArray_1(c *jMultisetCore, d []any) []any {
	return sliceToArray1(multisetOp_toSlice(c), d)
}
func multisetOp_containsAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
	return colContainsAll(func(x any) bool { return multisetOp_contains_1(c, x) }, col)
}
func multisetOp_forEach_1(c *jMultisetCore, fn Ext_java_util_function_Consumer) {
	CheckNotNull(fn)
	for _, x := range multisetOp_toSlice(c) {
		fn.Accept(x)
	}
}
func multisetOp_stream_0(c *jMultisetCore) Ext_java_util_stream_Stream {
	return streamFromSlice(multisetOp_toSlice(c))
}
func multisetOp_equals_1(c *jMultisetCore, o any) bool {
	ms, ok := o.(Ext_com_google_common_collect_Multiset)
	if !ok || ms == nil || javalang.IsNil(ms) {
		return false
	}
	if c.total != ms.Size() {
		return false
	}
	es := Cast[Ext_java_util_Set](ms.EntrySet())
	if c.counts.Size() != es.Size() {
		return false
	}
	for it := Iterate(es); it.HasNext(); {
		e := Cast[Ext_com_google_common_collect_Multiset_Entry](it.Next())
		if multisetOp_count_1(c, e.GetElement()) != e.GetCount() {
			return false
		}
	}
	return true
}
func multisetOp_hashCode_0(c *jMultisetCore) int32 {
	h := int32(0)
	c.counts.ForEachEntry(func(k, v any) bool {
		h += javalang.ObjectHashCode(k) ^ v.(int32)
		return true
	})
	return h
}
func multisetOp_toString_0(c *jMultisetCore) string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	c.counts.ForEachEntry(func(k, v any) bool {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		e := &jMultisetEntry{elem: k, count: v.(int32)}
		sb.WriteString(e.String())
		return true
	})
	sb.WriteByte(']')
	return sb.String()
}

type jMultisetIter struct {
	c       *jMultisetCore
	items   []any
	idx     int
	lastVal any
	hasLast bool
}

func (*jMultisetIter) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*jMultisetIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*jMultisetIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (it *jMultisetIter) HasNext() bool { return it.idx < len(it.items) }
func (it *jMultisetIter) Next() any {
	if it.idx >= len(it.items) {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	v := it.items[it.idx]
	it.idx++
	it.lastVal = v
	it.hasLast = true
	return v
}
func (it *jMultisetIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	multisetOp_remove_1(it.c, it.lastVal)
	it.hasLast = false
}

// --- Table (HashBasedTable, StandardTable, Tables.newCustomTable) ---

type jFastTableCoreProvider interface {
	fastTableCore_() (*jTableCore, bool)
}

type jTableRowEntry struct {
	rowKey    any
	rowKeyI64 int64
	col0      any
	val0      any
	multi     *JMapCore
	hashNext  int32
	hash      uint32
	hasSingle bool
	deleted   bool
	isColorId bool
}

func (e *jTableRowEntry) getRowKey() any {
	if e.isColorId {
		if e.rowKey == nil {
			e.rowKey = boxColorId(e.rowKeyI64)
		}
		return e.rowKey
	}
	return e.rowKey
}

func (e *jTableRowEntry) matchesAnyRowKey(r any, fk FastCustomKey, h uint32) bool {
	if e.hash != h {
		return false
	}
	if e.isColorId {
		if rc, ok := r.(jRawColorIdProvider); ok {
			raw, valid := rc.RawColorId_()
			return valid && raw == e.rowKeyI64
		}
		return false
	}
	if fk != nil {
		return fk.Equals(e.rowKey)
	}
	return javalang.ObjectsEquals(e.rowKey, r)
}

func (e *jTableRowEntry) isEmptyRow() bool {
	return e == nil || e.deleted || (!e.hasSingle && (e.multi == nil || e.multi.IsEmpty()))
}

type jTableCore struct {
	entries    []jTableRowEntry
	buckets    []int32
	strIdx     map[string]int32
	rowFactory func() *JMapCore
	liveRows   int32
	strLive    int32
	nonStrLive int32
	modVer     uint32
}

func newTableCore(rowFactory func() *JMapCore) *jTableCore {
	return &jTableCore{rowFactory: rowFactory, modVer: 1}
}

func newCustomTable(backingMap Ext_java_util_Map, factory Ext_com_google_common_base_Supplier) Ext_com_google_common_collect_Table {
	CheckNotNull(backingMap)
	CheckNotNull(factory)
	rf := func() *JMapCore {
		m := factory.Get()
		if mh, ok := m.(interface{ mc() *JMapCore }); ok {
			return mh.mc()
		}
		return NewJMapCore(false)
	}
	c := newTableCore(rf)
	if !backingMap.IsEmpty() {
		for cur := NewMapEntryCursor(backingMap); cur.NextEntry(); {
			r := cur.GetKey()
			if rm, ok := cur.GetValue().(Ext_java_util_Map); ok && rm != nil && !javalang.IsNil(rm) {
				for rc := NewMapEntryCursor(rm); rc.NextEntry(); {
					tableOp_put_3(c, r, rc.GetKey(), rc.GetValue())
				}
			}
		}
	}
	st := &Ext_com_google_common_collect_StandardTableBase{}
	st.BindExtSelf_(st)
	st.Core_ = c
	return st
}

func (t *jTableCore) newMultiRow() *JMapCore {
	if t.rowFactory != nil {
		return t.rowFactory()
	}
	return NewJMapCore(false)
}

func (t *jTableCore) rebuildBuckets(capB int) {
	nb := make([]int32, capB)
	mask := uint32(capB - 1)
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if !e.isColorId {
			if _, isStr := e.rowKey.(string); isStr {
				continue
			}
		}
		slot := smearHash(int32(e.hash)) & mask
		e.hashNext = nb[slot]
		nb[slot] = int32(i + 1)
	}
	t.buckets = nb
}

func (t *jTableCore) rebuildStrIdx() {
	m := make(map[string]int32, max(16, int(t.strLive)*2))
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted || e.isColorId {
			continue
		}
		if s, isStr := e.rowKey.(string); isStr {
			m[s] = int32(i)
		}
	}
	t.strIdx = m
}

func (t *jTableCore) findRowIdxColorId(raw int64) int {
	if len(t.entries) == 0 {
		return -1
	}
	h := uint32(int32(raw))
	if len(t.buckets) > 0 {
		slot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
		for pos := t.buckets[slot]; pos != 0; {
			e := &t.entries[pos-1]
			if e.hash == h {
				if e.isColorId {
					if e.rowKeyI64 == raw {
						return int(pos - 1)
					}
				} else if rc, ok := e.rowKey.(jRawColorIdProvider); ok {
					if eraw, valid := rc.RawColorId_(); valid && eraw == raw {
						return int(pos - 1)
					}
				}
			}
			pos = e.hashNext
		}
		return -1
	}
	for i := range t.entries {
		e := &t.entries[i]
		if !e.deleted && e.hash == h {
			if e.isColorId {
				if e.rowKeyI64 == raw {
					return i
				}
			} else if rc, ok := e.rowKey.(jRawColorIdProvider); ok {
				if eraw, valid := rc.RawColorId_(); valid && eraw == raw {
					return i
				}
			}
		}
	}
	return -1
}

func (t *jTableCore) findRowIdx(r any) int {
	if len(t.entries) == 0 || r == nil {
		return -1
	}
	if s, isStr := r.(string); isStr {
		if s == NullString {
			return -1
		}
		if t.strIdx != nil {
			if idx, ok := t.strIdx[s]; ok {
				return int(idx)
			}
			return -1
		}
		for i := range t.entries {
			e := &t.entries[i]
			if !e.deleted && !e.isColorId {
				if es, ok := e.rowKey.(string); ok && es == s {
					return i
				}
			}
		}
		return -1
	}
	if rc, ok := r.(jRawColorIdProvider); ok {
		raw, valid := rc.RawColorId_()
		if !valid {
			return -1
		}
		return t.findRowIdxColorId(raw)
	}
	if fk, ok := r.(FastCustomKey); ok {
		if fk.IsJavaNil_() {
			return -1
		}
		h := uint32(fk.HashCode())
		if len(t.buckets) > 0 {
			slot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
			for pos := t.buckets[slot]; pos != 0; {
				e := &t.entries[pos-1]
				if e.matchesAnyRowKey(r, fk, h) {
					return int(pos - 1)
				}
				pos = e.hashNext
			}
			return -1
		}
		for i := range t.entries {
			e := &t.entries[i]
			if !e.deleted && e.matchesAnyRowKey(r, fk, h) {
				return i
			}
		}
		return -1
	}
	if javalang.IsNil(r) {
		return -1
	}
	h := uint32(elemJavaHashCode(r, false))
	if len(t.buckets) > 0 {
		slot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
		for pos := t.buckets[slot]; pos != 0; {
			e := &t.entries[pos-1]
			if e.matchesAnyRowKey(r, nil, h) {
				return int(pos - 1)
			}
			pos = e.hashNext
		}
		return -1
	}
	for i := range t.entries {
		e := &t.entries[i]
		if !e.deleted && e.matchesAnyRowKey(r, nil, h) {
			return i
		}
	}
	return -1
}

func (t *jTableCore) insertNewRowColorId(raw int64, c, v any) {
	if t.entries == nil {
		t.entries = make([]jTableRowEntry, 0, 4)
	}
	idx := len(t.entries)
	h := uint32(int32(raw))
	if t.rowFactory != nil {
		mc := t.rowFactory()
		mc.Put(c, v)
		t.entries = append(t.entries, jTableRowEntry{rowKeyI64: raw, isColorId: true, multi: mc, hash: h})
	} else {
		t.entries = append(t.entries, jTableRowEntry{rowKeyI64: raw, isColorId: true, col0: c, val0: v, hash: h, hasSingle: true})
	}
	t.liveRows++
	t.nonStrLive++
	if len(t.buckets) == 0 {
		if t.nonStrLive >= 8 {
			t.rebuildBuckets(16)
		}
		return
	}
	if int(t.nonStrLive)*4 >= len(t.buckets)*3 {
		t.rebuildBuckets(len(t.buckets) * 2)
		return
	}
	slot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
	t.entries[idx].hashNext = t.buckets[slot]
	t.buckets[slot] = int32(idx + 1)
}

func (t *jTableCore) insertNewRow(r, c, v any) {
	if rc, ok := r.(jRawColorIdProvider); ok {
		if raw, valid := rc.RawColorId_(); valid {
			t.insertNewRowColorId(raw, c, v)
			return
		}
	}
	if t.entries == nil {
		t.entries = make([]jTableRowEntry, 0, 4)
	}
	idx := len(t.entries)
	if s, isStr := r.(string); isStr {
		if t.rowFactory != nil {
			mc := t.rowFactory()
			mc.Put(c, v)
			t.entries = append(t.entries, jTableRowEntry{rowKey: r, multi: mc})
		} else {
			t.entries = append(t.entries, jTableRowEntry{rowKey: r, col0: c, val0: v, hasSingle: true})
		}
		t.liveRows++
		t.strLive++
		if t.strIdx != nil {
			t.strIdx[s] = int32(idx)
		} else if t.strLive >= 8 {
			t.rebuildStrIdx()
		}
		return
	}
	var h uint32
	if fk, ok := r.(FastCustomKey); ok {
		h = uint32(fk.HashCode())
	} else {
		h = uint32(elemJavaHashCode(r, false))
	}
	if t.rowFactory != nil {
		mc := t.rowFactory()
		mc.Put(c, v)
		t.entries = append(t.entries, jTableRowEntry{rowKey: r, multi: mc, hash: h})
	} else {
		t.entries = append(t.entries, jTableRowEntry{rowKey: r, col0: c, val0: v, hash: h, hasSingle: true})
	}
	t.liveRows++
	t.nonStrLive++
	if len(t.buckets) == 0 {
		if t.nonStrLive >= 8 {
			t.rebuildBuckets(16)
		}
		return
	}
	if int(t.nonStrLive)*4 >= len(t.buckets)*3 {
		t.rebuildBuckets(len(t.buckets) * 2)
		return
	}
	slot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
	t.entries[idx].hashNext = t.buckets[slot]
	t.buckets[slot] = int32(idx + 1)
}

func (t *jTableCore) deleteRowAt(idx int) {
	if idx < 0 || idx >= len(t.entries) {
		return
	}
	e := &t.entries[idx]
	if e.deleted {
		return
	}
	if s, isStr := e.rowKey.(string); isStr {
		if t.strIdx != nil {
			delete(t.strIdx, s)
		}
		t.strLive--
	} else {
		if len(t.buckets) > 0 {
			pos := int32(idx + 1)
			slot := smearHash(int32(e.hash)) & uint32(len(t.buckets)-1)
			var prev int32
			for cur := t.buckets[slot]; cur != 0; {
				if cur == pos {
					if prev == 0 {
						t.buckets[slot] = e.hashNext
					} else {
						t.entries[prev-1].hashNext = e.hashNext
					}
					break
				}
				prev = cur
				cur = t.entries[cur-1].hashNext
			}
		}
		t.nonStrLive--
	}
	*e = jTableRowEntry{deleted: true}
	t.liveRows--
	if t.liveRows == 0 {
		t.entries = t.entries[:0]
		if len(t.buckets) > 0 {
			clear(t.buckets)
		}
		if t.strIdx != nil {
			clear(t.strIdx)
		}
		t.modVer++
		return
	}
	if len(t.entries) >= 32 && int(t.liveRows)*2 < len(t.entries) {
		w := 0
		for i := range t.entries {
			if !t.entries[i].deleted {
				if w != i {
					t.entries[w] = t.entries[i]
				}
				w++
			}
		}
		clear(t.entries[w:])
		t.entries = t.entries[:w]
		t.modVer++
		if t.strIdx != nil {
			t.rebuildStrIdx()
		}
		if len(t.buckets) > 0 {
			t.rebuildBuckets(len(t.buckets))
		}
	}
}

func tableOp_size_0(t *jTableCore) int32 {
	var n int32
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if e.hasSingle {
			n++
		} else if e.multi != nil {
			n += e.multi.Size()
		}
	}
	return n
}

func tableOp_isEmpty_0(t *jTableCore) bool {
	if t.liveRows == 0 {
		return true
	}
	for i := range t.entries {
		e := &t.entries[i]
		if !e.isEmptyRow() {
			return false
		}
	}
	return true
}

func tableOp_containsRow_1(t *jTableCore, r any) bool {
	idx := t.findRowIdx(r)
	return idx >= 0 && !t.entries[idx].isEmptyRow()
}

func tableOp_containsColumn_1(t *jTableCore, c any) bool {
	if c == nil || javalang.IsNil(c) {
		return false
	}
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if e.hasSingle {
			if javalang.ObjectsEquals(e.col0, c) {
				return true
			}
		} else if e.multi != nil && e.multi.ContainsKey(c) {
			return true
		}
	}
	return false
}

func tableOp_containsValue_1(t *jTableCore, v any) bool {
	if v == nil || javalang.IsNil(v) {
		return false
	}
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if e.hasSingle {
			if javalang.ObjectsEquals(e.val0, v) {
				return true
			}
		} else if e.multi != nil && e.multi.ContainsValue(v) {
			return true
		}
	}
	return false
}

func tableOp_contains_2(t *jTableCore, r, c any) bool {
	if c == nil || javalang.IsNil(c) {
		return false
	}
	idx := t.findRowIdx(r)
	if idx < 0 {
		return false
	}
	e := &t.entries[idx]
	if e.hasSingle {
		return javalang.ObjectsEquals(e.col0, c)
	}
	return e.multi != nil && e.multi.ContainsKey(c)
}

func tableOp_get_2(t *jTableCore, r, c any) any {
	if c == nil || javalang.IsNil(c) {
		return nil
	}
	idx := t.findRowIdx(r)
	if idx < 0 {
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if javalang.ObjectsEquals(e.col0, c) {
			return e.val0
		}
		return nil
	}
	if e.multi == nil {
		return nil
	}
	return e.multi.Get(c)
}

func tableOp_put_3(t *jTableCore, r, c, v any) any {
	r = CheckNotNull(ToAny(r))
	c = CheckNotNull(ToAny(c))
	v = CheckNotNull(ToAny(v))
	idx := t.findRowIdx(r)
	if idx < 0 {
		t.insertNewRow(r, c, v)
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if javalang.ObjectsEquals(e.col0, c) {
			old := e.val0
			e.val0 = v
			return old
		}
		mc := t.newMultiRow()
		mc.Put(e.col0, e.val0)
		old := mc.Put(c, v)
		e.col0 = nil
		e.val0 = nil
		e.hasSingle = false
		e.multi = mc
		return old
	}
	if e.multi == nil {
		if t.rowFactory == nil {
			e.col0 = c
			e.val0 = v
			e.hasSingle = true
			return nil
		}
		e.multi = t.newMultiRow()
	}
	return e.multi.Put(c, v)
}

func tableOp_remove_2(t *jTableCore, r, c any) any {
	if c == nil || javalang.IsNil(c) {
		return nil
	}
	idx := t.findRowIdx(r)
	if idx < 0 {
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if !javalang.ObjectsEquals(e.col0, c) {
			return nil
		}
		old := e.val0
		t.deleteRowAt(idx)
		return old
	}
	if e.multi == nil {
		return nil
	}
	old := e.multi.Remove(c)
	if e.multi.IsEmpty() {
		t.deleteRowAt(idx)
	}
	return old
}

func tableOp_clear_0(t *jTableCore) {
	if len(t.entries) == 0 {
		return
	}
	clear(t.entries)
	t.entries = t.entries[:0]
	if len(t.buckets) > 0 {
		clear(t.buckets)
	}
	if t.strIdx != nil {
		clear(t.strIdx)
	}
	t.liveRows = 0
	t.strLive = 0
	t.nonStrLive = 0
	t.modVer++
}

func tableOp_rowKeySet_0(t *jTableCore) Ext_java_util_Set {
	s := &jTableRowKeySet{t: t}
	s.BindExtSelf_(s)
	return s
}

func tableOp_columnKeySet_0(t *jTableCore) Ext_java_util_Set {
	var sc JSetCore
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if e.hasSingle {
			sc.Add(e.col0)
		} else if e.multi != nil {
			if e.multi.mode == jMapModeStr && !e.multi.Identity {
				for j := int(e.multi.firstIdx); j < len(e.multi.strEntries); j++ {
					if ck := e.multi.strEntries[j].key; ck != jTombstoneStr {
						sc.AddStr(ck)
					}
				}
			} else {
				e.multi.ForEachEntry(func(ck, _ any) bool {
					sc.Add(ck)
					return true
				})
			}
		}
	}
	obj := &Ext_java_util_LinkedHashSet{Ext_java_util_HashSetBase: Ext_java_util_HashSetBase{inlineSet: sc}}
	obj.BindExtSelf_(obj)
	return obj
}

func tableOp_values_0(t *jTableCore) Ext_java_util_Collection {
	var vals []any
	for i := range t.entries {
		e := &t.entries[i]
		if e.deleted {
			continue
		}
		if e.hasSingle {
			vals = append(vals, e.val0)
		} else if e.multi != nil {
			e.multi.ForEachValue(func(cv any) bool {
				vals = append(vals, cv)
				return true
			})
		}
	}
	return wrapArrayListOwn(vals)
}

func tableOp_row_1(t *jTableCore, r any) Ext_java_util_Map {
	r = CheckNotNull(ToAny(r))
	if rc, ok := r.(jRawColorIdProvider); ok {
		if raw, valid := rc.RawColorId_(); valid {
			return t.rowColorId(raw)
		}
	}
	idx := t.findRowIdx(r)
	return &jTableRowMap{t: t, r: r, hintIdx: int32(idx), hintVer: t.modVer}
}

func (t *jTableCore) containsRowColorId(raw int64) bool {
	idx := t.findRowIdxColorId(raw)
	return idx >= 0 && !t.entries[idx].isEmptyRow()
}

func (t *jTableCore) containsColorId(raw int64, c any) bool {
	if c == nil || javalang.IsNil(c) {
		return false
	}
	idx := t.findRowIdxColorId(raw)
	if idx < 0 {
		return false
	}
	e := &t.entries[idx]
	if e.hasSingle {
		return javalang.ObjectsEquals(e.col0, c)
	}
	return e.multi != nil && e.multi.ContainsKey(c)
}

func (t *jTableCore) getColorId(raw int64, c any) any {
	if c == nil || javalang.IsNil(c) {
		return nil
	}
	idx := t.findRowIdxColorId(raw)
	if idx < 0 {
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if javalang.ObjectsEquals(e.col0, c) {
			return e.val0
		}
		return nil
	}
	if e.multi == nil {
		return nil
	}
	return e.multi.Get(c)
}

func (t *jTableCore) putColorId(raw int64, c, v any) any {
	c = CheckNotNull(ToAny(c))
	v = CheckNotNull(ToAny(v))
	idx := t.findRowIdxColorId(raw)
	if idx < 0 {
		t.insertNewRowColorId(raw, c, v)
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if javalang.ObjectsEquals(e.col0, c) {
			old := e.val0
			e.val0 = v
			return old
		}
		mc := t.newMultiRow()
		mc.Put(e.col0, e.val0)
		old := mc.Put(c, v)
		e.col0 = nil
		e.val0 = nil
		e.hasSingle = false
		e.multi = mc
		return old
	}
	if e.multi == nil {
		if t.rowFactory == nil {
			e.col0 = c
			e.val0 = v
			e.hasSingle = true
			return nil
		}
		e.multi = t.newMultiRow()
	}
	return e.multi.Put(c, v)
}

func (t *jTableCore) removeColorId(raw int64, c any) any {
	if c == nil || javalang.IsNil(c) {
		return nil
	}
	idx := t.findRowIdxColorId(raw)
	if idx < 0 {
		return nil
	}
	e := &t.entries[idx]
	if e.hasSingle {
		if !javalang.ObjectsEquals(e.col0, c) {
			return nil
		}
		old := e.val0
		t.deleteRowAt(idx)
		return old
	}
	if e.multi == nil {
		return nil
	}
	old := e.multi.Remove(c)
	if e.multi.IsEmpty() {
		t.deleteRowAt(idx)
	}
	return old
}

func (t *jTableCore) rowColorId(raw int64) Ext_java_util_Map {
	idx := t.findRowIdxColorId(raw)
	return &jTableRowMap{t: t, rI64: raw, isColorId: true, hintIdx: int32(idx), hintVer: t.modVer}
}

func tableOp_rowMap_0(t *jTableCore) Ext_java_util_Map {
	m := &jTableRowMapOuter{t: t}
	m.BindExtSelf_(m)
	return m
}

type jTableRowKeySet struct {
	Ext_java_util_AbstractSetBase
	t *jTableCore
}

func (s *jTableRowKeySet) nextLiveIdx(from int) int {
	entries := s.t.entries
	for i := from; i < len(entries); i++ {
		if !entries[i].isEmptyRow() {
			return i
		}
	}
	return -1
}

func (s *jTableRowKeySet) iterInit_() (int32, bool) {
	if s == nil {
		ThrowNPE()
	}
	i := s.nextLiveIdx(0)
	if i < 0 {
		return -1, false
	}
	return int32(i), true
}

func (s *jTableRowKeySet) iterNext_(pos int32) (any, int32) {
	r := s.t.entries[pos].getRowKey()
	next := s.nextLiveIdx(int(pos) + 1)
	return r, int32(next)
}

func (s *jTableRowKeySet) Size() int32 {
	var n int32
	for i := range s.t.entries {
		if !s.t.entries[i].isEmptyRow() {
			n++
		}
	}
	return n
}

func (s *jTableRowKeySet) IsEmpty() bool       { return tableOp_isEmpty_0(s.t) }
func (s *jTableRowKeySet) Contains(r any) bool { return tableOp_containsRow_1(s.t, r) }
func (s *jTableRowKeySet) Clear()              { tableOp_clear_0(s.t) }
func (s *jTableRowKeySet) Remove_1_any(r any) bool {
	idx := s.t.findRowIdx(r)
	if idx < 0 || s.t.entries[idx].isEmptyRow() {
		return false
	}
	s.t.deleteRowAt(idx)
	return true
}
func (s *jTableRowKeySet) Iterator() any {
	var keys []any
	for i := range s.t.entries {
		if !s.t.entries[i].isEmptyRow() {
			keys = append(keys, s.t.entries[i].getRowKey())
		}
	}
	return &jTableRowKeyIter{t: s.t, keys: keys}
}

type jTableRowKeyIter struct {
	t       *jTableCore
	keys    []any
	idx     int
	lastKey any
	hasLast bool
}

func (*jTableRowKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jTableRowKeyIter) HasNext() bool                          { return it.idx < len(it.keys) }
func (it *jTableRowKeyIter) Next() any {
	if it.idx >= len(it.keys) {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	k := it.keys[it.idx]
	it.idx++
	it.lastKey = k
	it.hasLast = true
	return k
}
func (it *jTableRowKeyIter) Remove() {
	if !it.hasLast {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	if idx := it.t.findRowIdx(it.lastKey); idx >= 0 {
		it.t.deleteRowAt(idx)
	}
	it.hasLast = false
}

type jTableRowMapOuter struct {
	Ext_java_util_AbstractMapBase
	t *jTableCore
}

func (m *jTableRowMapOuter) Size() int32 {
	var n int32
	for i := range m.t.entries {
		if !m.t.entries[i].isEmptyRow() {
			n++
		}
	}
	return n
}
func (m *jTableRowMapOuter) ContainsKey(r any) bool { return tableOp_containsRow_1(m.t, r) }
func (m *jTableRowMapOuter) Get(r any) any {
	if !tableOp_containsRow_1(m.t, r) {
		return nil
	}
	return tableOp_row_1(m.t, r)
}
func (m *jTableRowMapOuter) Remove_1(r any) any {
	idx := m.t.findRowIdx(r)
	if idx < 0 || m.t.entries[idx].isEmptyRow() {
		return nil
	}
	e := &m.t.entries[idx]
	copyCore := NewJMapCore(false)
	if e.hasSingle {
		copyCore.Put(e.col0, e.val0)
	} else if e.multi != nil {
		if !copyCore.CopyFromSameMode(e.multi) {
			copyHashCoreInto(copyCore, e.multi)
		}
	}
	m.t.deleteRowAt(idx)
	lhm := &Ext_java_util_LinkedHashMap{Ext_java_util_HashMapBase: Ext_java_util_HashMapBase{inlineCore: *copyCore}}
	lhm.BindExtSelf_(lhm)
	return lhm
}
func (m *jTableRowMapOuter) Clear() { tableOp_clear_0(m.t) }
func (m *jTableRowMapOuter) EntrySet() any {
	var ents []any
	for i := range m.t.entries {
		e := &m.t.entries[i]
		if !e.isEmptyRow() {
			r := e.getRowKey()
			ents = append(ents, &jSimpleEntry{k: r, v: tableOp_row_1(m.t, r)})
		}
	}
	return wrapLinkedHashSet(ents)
}

type jTableRowMap struct {
	t         *jTableCore
	r         any
	rI64      int64
	hintIdx   int32
	hintVer   uint32
	isColorId bool
}

func (m *jTableRowMap) resolveEntry() (*jTableRowEntry, int) {
	t := m.t
	if m.hintIdx >= 0 && m.hintVer == t.modVer && int(m.hintIdx) < len(t.entries) {
		e := &t.entries[m.hintIdx]
		if !e.deleted {
			return e, int(m.hintIdx)
		}
		m.hintIdx = -1
	}
	var idx int
	if m.isColorId {
		idx = t.findRowIdxColorId(m.rI64)
	} else {
		idx = t.findRowIdx(m.r)
	}
	if idx < 0 {
		return nil, -1
	}
	m.hintIdx = int32(idx)
	m.hintVer = t.modVer
	return &t.entries[idx], idx
}

func (m *jTableRowMap) putImpl(c, v any) any {
	if m.isColorId {
		return m.t.putColorId(m.rI64, c, v)
	}
	return tableOp_put_3(m.t, m.r, c, v)
}

func (m *jTableRowMap) removeImpl(c any) any {
	if m.isColorId {
		return m.t.removeColorId(m.rI64, c)
	}
	return tableOp_remove_2(m.t, m.r, c)
}

func (m *jTableRowMap) sizeImpl() int32 {
	e, _ := m.resolveEntry()
	if e == nil {
		return 0
	}
	if e.hasSingle {
		return 1
	}
	if e.multi != nil {
		return e.multi.Size()
	}
	return 0
}

func (m *jTableRowMap) isEmptyImpl() bool {
	e, _ := m.resolveEntry()
	return e.isEmptyRow()
}

func (m *jTableRowMap) containsKeyImpl(c any) bool {
	if c == nil || javalang.IsNil(c) {
		return false
	}
	e, _ := m.resolveEntry()
	if e == nil {
		return false
	}
	if e.hasSingle {
		return javalang.ObjectsEquals(e.col0, c)
	}
	return e.multi != nil && e.multi.ContainsKey(c)
}
func (m *jTableRowMap) MapContainsStr_(c string) bool {
	if m == nil {
		ThrowNPE()
	}
	if c == NullString {
		return false
	}
	e, _ := m.resolveEntry()
	if e == nil {
		return false
	}
	if e.hasSingle {
		cs, ok := e.col0.(string)
		return ok && cs == c
	}
	return e.multi != nil && e.multi.ContainsStr(c)
}
func (m *jTableRowMap) MapGetStr_(c string) any {
	if m == nil {
		ThrowNPE()
	}
	if c == NullString {
		return nil
	}
	e, _ := m.resolveEntry()
	if e == nil {
		return nil
	}
	if e.hasSingle {
		if cs, ok := e.col0.(string); ok && cs == c {
			return e.val0
		}
		return nil
	}
	if e.multi != nil {
		return e.multi.GetStr(c)
	}
	return nil
}
func (m *jTableRowMap) MapGetOrDefaultStr_(c string, def any) any {
	if v := m.MapGetStr_(c); v != nil {
		return v
	}
	return ToAny(def)
}
func (m *jTableRowMap) MapPutStr_(c string, v any) any {
	if m == nil {
		ThrowNPE()
	}
	return m.putImpl(BoxString(c), v)
}
func (m *jTableRowMap) MapRemoveStr_(c string) any {
	if m == nil {
		ThrowNPE()
	}
	if c == NullString {
		return nil
	}
	if !m.MapContainsStr_(c) {
		return nil
	}
	return m.removeImpl(BoxString(c))
}

func (m *jTableRowMap) containsValueImpl(v any) bool {
	if v == nil || javalang.IsNil(v) {
		return false
	}
	e, _ := m.resolveEntry()
	if e == nil {
		return false
	}
	if e.hasSingle {
		return javalang.ObjectsEquals(e.val0, v)
	}
	return e.multi != nil && e.multi.ContainsValue(v)
}

func (m *jTableRowMap) getImpl(c any) any {
	if c == nil || javalang.IsNil(c) {
		return nil
	}
	e, _ := m.resolveEntry()
	if e == nil {
		return nil
	}
	if e.hasSingle {
		if javalang.ObjectsEquals(e.col0, c) {
			return e.val0
		}
		return nil
	}
	if e.multi == nil {
		return nil
	}
	return e.multi.Get(c)
}

func (m *jTableRowMap) getOrDefaultImpl(c, def any) any {
	v := m.getImpl(c)
	if v != nil {
		return v
	}
	return ToAny(def)
}

func (m *jTableRowMap) clearImpl() {
	_, idx := m.resolveEntry()
	if idx >= 0 {
		m.t.deleteRowAt(idx)
		m.hintIdx = -1
	}
}

func (m *jTableRowMap) forEachImpl(action Ext_java_util_function_BiConsumer) {
	CheckNotNull(action)
	e, _ := m.resolveEntry()
	if e == nil {
		return
	}
	if e.hasSingle {
		action.Accept(e.col0, e.val0)
		return
	}
	if e.multi != nil {
		mapCoreForEach(e.multi, action)
	}
}

func (m *jTableRowMap) keySetImpl() Ext_java_util_Set {
	s := &jTableRowKeyView{m: m}
	s.BindExtSelf_(s)
	return s
}

func (m *jTableRowMap) entrySetImpl() Ext_java_util_Set {
	s := &jTableRowEntrySet{m: m}
	s.BindExtSelf_(s)
	return s
}

func (m *jTableRowMap) equalsImpl(other any) bool {
	if other == m {
		return true
	}
	om, ok := other.(Ext_java_util_Map)
	if !ok || om == nil || javalang.IsNil(om) {
		return false
	}
	e, _ := m.resolveEntry()
	if e.isEmptyRow() {
		return om.IsEmpty()
	}
	if e.hasSingle {
		return om.Size() == 1 && om.ContainsKey(e.col0) && javalang.ObjectsEquals(e.val0, om.Get(e.col0))
	}
	return mapCoreEquals(e.multi, other)
}

func (m *jTableRowMap) hashCodeImpl() int32 {
	e, _ := m.resolveEntry()
	if e.isEmptyRow() {
		return 0
	}
	if e.hasSingle {
		return elemJavaHashCode(e.col0, false) ^ elemJavaHashCode(e.val0, false)
	}
	return mapCoreHash(e.multi)
}

func (m *jTableRowMap) stringImpl() string {
	e, _ := m.resolveEntry()
	if e.isEmptyRow() {
		return "{}"
	}
	if e.hasSingle {
		return "{" + StrVal(e.col0) + "=" + StrVal(e.val0) + "}"
	}
	return mapCoreString(e.multi)
}

type jTableRowKeyView struct {
	Ext_java_util_AbstractSetBase
	m *jTableRowMap
}

func (s *jTableRowKeyView) Size() int32             { return s.m.sizeImpl() }
func (s *jTableRowKeyView) IsEmpty() bool           { return s.m.isEmptyImpl() }
func (s *jTableRowKeyView) Contains(c any) bool     { return s.m.containsKeyImpl(c) }
func (s *jTableRowKeyView) Clear()                  { s.m.clearImpl() }
func (s *jTableRowKeyView) Remove_1_any(c any) bool { return s.m.removeImpl(c) != nil }
func (s *jTableRowKeyView) Iterator() any {
	e, _ := s.m.resolveEntry()
	if e.isEmptyRow() {
		return emptyMapIter_
	}
	if e.hasSingle {
		return &jTableRowSingleIter{m: s.m, col: e.col0, val: e.val0, mode: 0}
	}
	return newMapKeyIter(e.multi)
}

type jTableRowEntrySet struct {
	Ext_java_util_AbstractSetBase
	m *jTableRowMap
}

func (s *jTableRowEntrySet) Size() int32   { return s.m.sizeImpl() }
func (s *jTableRowEntrySet) IsEmpty() bool { return s.m.isEmptyImpl() }
func (s *jTableRowEntrySet) Clear()        { s.m.clearImpl() }
func (s *jTableRowEntrySet) Contains(o any) bool {
	ent, ok := o.(Ext_java_util_Map_Entry)
	if !ok || ent == nil || javalang.IsNil(ent) {
		return false
	}
	k := ent.GetKey()
	v := s.m.getImpl(k)
	return v != nil && javalang.ObjectsEquals(v, ToAny(ent.GetValue()))
}
func (s *jTableRowEntrySet) Remove_1_any(o any) bool {
	ent, ok := o.(Ext_java_util_Map_Entry)
	if !ok || ent == nil || javalang.IsNil(ent) {
		return false
	}
	k := ent.GetKey()
	v := s.m.getImpl(k)
	if v == nil || !javalang.ObjectsEquals(v, ToAny(ent.GetValue())) {
		return false
	}
	s.m.removeImpl(k)
	return true
}
func (s *jTableRowEntrySet) Iterator() any {
	e, _ := s.m.resolveEntry()
	if e.isEmptyRow() {
		return emptyMapIter_
	}
	if e.hasSingle {
		return &jTableRowSingleIter{m: s.m, col: e.col0, val: e.val0, mode: 2}
	}
	c, ok := initMapCursor(e.multi, true)
	if !ok {
		return emptyMapIter_
	}
	return &jMapEntryIter{cur: c}
}

type jTableRowSingleIter struct {
	m       *jTableRowMap
	col     any
	val     any
	mode    uint8 // 0=key, 1=val, 2=entry
	done    bool
	removed bool
}

func (*jTableRowSingleIter) Implements_j2gort_Ext_java_util_Iterator()  {}
func (*jTableRowSingleIter) Implements_j2gort_Ext_java_util_Map_Entry() {}
func (it *jTableRowSingleIter) GetKey() any                             { return it.col }
func (it *jTableRowSingleIter) GetValue() any                           { return it.val }
func (it *jTableRowSingleIter) SetValue(nv any) any {
	nv = CheckNotNull(ToAny(nv))
	old := it.m.putImpl(it.col, nv)
	it.val = nv
	return old
}
func (it *jTableRowSingleIter) Equals(o any) bool { return mapEntryEquals(it, o) }
func (it *jTableRowSingleIter) HashCode() int32   { return mapEntryHash(it) }
func (it *jTableRowSingleIter) String() string    { return StrVal(it.col) + "=" + StrVal(it.val) }
func (it *jTableRowSingleIter) HasNext() bool     { return !it.done }
func (it *jTableRowSingleIter) Next() any {
	if it.done {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	it.done = true
	switch it.mode {
	case 0:
		return it.col
	case 1:
		return it.val
	default:
		return it
	}
}
func (it *jTableRowSingleIter) Remove() {
	if !it.done || it.removed {
		panic(javalang.NewIllegalStateException("remove"))
	}
	it.removed = true
	it.m.removeImpl(it.col)
}

type jTableRowValues struct {
	m *jTableRowMap
}

func (v jTableRowValues) IsJavaNil_() bool { return v.m == nil }
func (v jTableRowValues) typedViewLen_() (int, bool) {
	if v.m == nil {
		ThrowNPE()
	}
	e, _ := v.m.resolveEntry()
	if e == nil || e.deleted {
		return 0, true
	}
	if e.hasSingle {
		return 1, true
	}
	if e.multi == nil || e.multi.IsEmpty() {
		return 0, true
	}
	return 0, false
}
func (v jTableRowValues) typedViewAt_(int) any {
	e, _ := v.m.resolveEntry()
	if e == nil || !e.hasSingle {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return e.val0
}
func (jTableRowValues) typedViewNodeAt_(int) (uint32, bool) { return 0, false }
func (v jTableRowValues) iteratorImpl() any {
	e, _ := v.m.resolveEntry()
	if e.isEmptyRow() {
		return emptyMapIter_
	}
	if e.hasSingle {
		return &jTableRowSingleIter{m: v.m, col: e.col0, val: e.val0, mode: 1}
	}
	c, ok := initMapCursor(e.multi, true)
	if !ok {
		return emptyMapIter_
	}
	return &jMapValIter{cur: c}
}
func (v jTableRowValues) toArrayImpl() []any {
	e, _ := v.m.resolveEntry()
	if e.isEmptyRow() {
		return []any{}
	}
	if e.hasSingle {
		return []any{e.val0}
	}
	var out []any
	if e.multi != nil {
		out = e.multi.LiveValues()
	}
	return out
}
func (v jTableRowValues) forEachImpl(action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	e, _ := v.m.resolveEntry()
	if e.isEmptyRow() {
		return
	}
	if e.hasSingle {
		action.Accept(e.val0)
		return
	}
	if e.multi != nil {
		e.multi.ForEachValue(func(cv any) bool {
			action.Accept(cv)
			return true
		})
	}
}

// --- PriorityQueue ---

type jPriorityQueue struct {
	items []any
	cmp   func(a, b any) int32
}

func newPriorityQueue(cmp func(a, b any) int32) *jPriorityQueue {
	if cmp == nil {
		cmp = CompareValues
	}
	return &jPriorityQueue{cmp: cmp}
}

func (pq *jPriorityQueue) offer(v any) bool {
	v = CheckNotNull(ToAny(v))
	pq.items = append(pq.items, v)
	pq.siftUp(len(pq.items) - 1)
	return true
}

func (pq *jPriorityQueue) siftUp(k int) {
	x := pq.items[k]
	for k > 0 {
		p := (k - 1) >> 1
		if pq.cmp(x, pq.items[p]) >= 0 {
			break
		}
		pq.items[k] = pq.items[p]
		k = p
	}
	pq.items[k] = x
}

func (pq *jPriorityQueue) siftDown(k int) {
	n := len(pq.items)
	x := pq.items[k]
	half := n >> 1
	for k < half {
		c := (k << 1) + 1
		right := c + 1
		if right < n && pq.cmp(pq.items[c], pq.items[right]) > 0 {
			c = right
		}
		if pq.cmp(x, pq.items[c]) <= 0 {
			break
		}
		pq.items[k] = pq.items[c]
		k = c
	}
	pq.items[k] = x
}

func (pq *jPriorityQueue) poll() any {
	n := len(pq.items)
	if n == 0 {
		return nil
	}
	res := pq.items[0]
	last := pq.items[n-1]
	pq.items[n-1] = nil
	pq.items = pq.items[:n-1]
	if n > 1 {
		pq.items[0] = last
		pq.siftDown(0)
	}
	return res
}

func pqOp_size_0(pq *jPriorityQueue) int32      { return int32(len(pq.items)) }
func pqOp_isEmpty_0(pq *jPriorityQueue) bool    { return len(pq.items) == 0 }
func pqOp_add_1(pq *jPriorityQueue, v any) bool { return pq.offer(v) }
func pqOp_peek_0(pq *jPriorityQueue) any {
	if len(pq.items) == 0 {
		return nil
	}
	return pq.items[0]
}
func pqOp_poll_0(pq *jPriorityQueue) any { return pq.poll() }
func pqOp_remove_0(pq *jPriorityQueue) any {
	if len(pq.items) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return pq.poll()
}
func pqOp_element_0(pq *jPriorityQueue) any {
	if len(pq.items) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return pq.items[0]
}
func pqOp_retainAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
	CheckNotNull(c)
	changed := false
	for i := len(pq.items) - 1; i >= 0; i-- {
		if !c.Contains(pq.items[i]) {
			pq.removeAt(i)
			changed = true
		}
	}
	return changed
}
func pqOp_removeIf_1(pq *jPriorityQueue, pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	changed := false
	for i := len(pq.items) - 1; i >= 0; i-- {
		if pred.Test(pq.items[i]) {
			pq.removeAt(i)
			changed = true
		}
	}
	return changed
}
func pqOp_remove_1(pq *jPriorityQueue, v any) bool {
	v = ToAny(v)
	for i, x := range pq.items {
		if javalang.ObjectsEquals(x, v) {
			pq.removeAt(i)
			return true
		}
	}
	return false
}
func (pq *jPriorityQueue) removeAt(i int) {
	n := len(pq.items) - 1
	if i == n {
		pq.items[n] = nil
		pq.items = pq.items[:n]
		return
	}
	moved := pq.items[n]
	pq.items[n] = nil
	pq.items = pq.items[:n]
	pq.items[i] = moved
	pq.siftDown(i)
	if pq.items[i] == moved {
		pq.siftUp(i)
	}
}
func pqOp_contains_1(pq *jPriorityQueue, v any) bool {
	v = ToAny(v)
	for _, x := range pq.items {
		if javalang.ObjectsEquals(x, v) {
			return true
		}
	}
	return false
}
func pqOp_clear_0(pq *jPriorityQueue)                  { pq.items = nil }
func pqOp_toArray_0(pq *jPriorityQueue) []any          { return CloneSlice(pq.items) }
func pqOp_toArray_1(pq *jPriorityQueue, d []any) []any { return sliceToArray1(pq.items, d) }
func pqOp_iterator_0(pq *jPriorityQueue) any           { return newSliceObjIter(CloneSlice(pq.items)) }
func pqOp_addAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
	CheckNotNull(c)
	added := false
	for it := Iterate(c); it.HasNext(); {
		pq.offer(it.Next())
		added = true
	}
	return added
}
func pqOp_containsAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
	return colContainsAll(func(x any) bool { return pqOp_contains_1(pq, x) }, c)
}
func pqOp_removeAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
	CheckNotNull(c)
	var kept []any
	changed := false
	for _, x := range pq.items {
		if c.Contains(x) {
			changed = true
		} else {
			kept = append(kept, x)
		}
	}
	if changed {
		pq.items = nil
		for _, x := range kept {
			pq.offer(x)
		}
	}
	return changed
}
func pqOp_forEach_1(pq *jPriorityQueue, fn Ext_java_util_function_Consumer) {
	CheckNotNull(fn)
	for _, x := range CloneSlice(pq.items) {
		fn.Accept(x)
	}
}
func pqOp_stream_0(pq *jPriorityQueue) Ext_java_util_stream_Stream {
	return streamFromSlice(CloneSlice(pq.items))
}
func pqOp_toString_0(pq *jPriorityQueue) string    { return sliceJavaString(pq.items) }
func pqOp_equals_1(pq *jPriorityQueue, o any) bool { return RefEquals(pq, o) }
func pqOp_hashCode_0(pq *jPriorityQueue) int32     { return javalang.IdentityHashCode(pq) }

// --- BitSet (inline [16]uint64 = 1024 dense bits + [8]int32 = 8 sparse high bits inside Ext_java_util_BitSet) ---

const (
	bsInlineWords = 16
	bsInlineBits  = bsInlineWords * 64
	bsSparseCap   = 8
)

//go:noinline
func bitSetThrowNegIdx(bitIndex int, name string) {
	panic(javalang.NewIndexOutOfBoundsException(name + " < 0: " + strconv.Itoa(bitIndex)))
}

//go:noinline
func bitSetThrowNegSize(nbits int32) {
	panic(javalang.NewNegativeArraySizeException("nbits < 0: " + strconv.Itoa(int(nbits))))
}

//go:noinline
func bitSetThrowBadRange(fromIndex, toIndex int) {
	panic(javalang.NewIndexOutOfBoundsException("fromIndex: " + strconv.Itoa(fromIndex) + " > toIndex: " + strconv.Itoa(toIndex)))
}

func bitSetCheckIndex(bitIndex int, name string) {
	if bitIndex < 0 {
		bitSetThrowNegIdx(bitIndex, name)
	}
}

func bitSetCheckRange(fromIndex, toIndex int) {
	if fromIndex < 0 {
		bitSetThrowNegIdx(fromIndex, "fromIndex")
	}
	if toIndex < 0 {
		bitSetThrowNegIdx(toIndex, "toIndex")
	}
	if fromIndex > toIndex {
		bitSetThrowBadRange(fromIndex, toIndex)
	}
}

func InitExt_java_util_BitSet_0(b *Ext_java_util_BitSet) {
	*b = Ext_java_util_BitSet{initW: 1}
}

func InitExt_java_util_BitSet_1(b *Ext_java_util_BitSet, nbits int32) {
	if nbits < 0 {
		bitSetThrowNegSize(nbits)
	}
	*b = Ext_java_util_BitSet{initW: int32((int64(nbits) + 63) >> 6)}
}

func bitSetInitCap(b *Ext_java_util_BitSet, nbits int32) {
	InitExt_java_util_BitSet_1(b, nbits)
}

func (b *Ext_java_util_BitSet) bsHiContains(idx int32) bool {
	n := int(b.hiLen)
	for i := 0; i < n; i++ {
		v := b.hi[i]
		if v >= idx {
			return v == idx
		}
	}
	return false
}

func (b *Ext_java_util_BitSet) bsHiInsert(idx int32) bool {
	n := int(b.hiLen)
	pos := 0
	for pos < n && b.hi[pos] < idx {
		pos++
	}
	if pos < n && b.hi[pos] == idx {
		return true
	}
	if n >= bsSparseCap {
		return false
	}
	copy(b.hi[pos+1:n+1], b.hi[pos:n])
	b.hi[pos] = idx
	b.hiLen = int32(n + 1)
	return true
}

func (b *Ext_java_util_BitSet) bsHiRemove(idx int32) bool {
	n := int(b.hiLen)
	for i := 0; i < n; i++ {
		v := b.hi[i]
		if v == idx {
			copy(b.hi[i:n-1], b.hi[i+1:n])
			b.hi[n-1] = 0
			b.hiLen = int32(n - 1)
			return true
		}
		if v > idx {
			return false
		}
	}
	return false
}

func (b *Ext_java_util_BitSet) bsWordAt(i int) uint64 {
	if b.words == nil {
		if uint(i) < bsInlineWords {
			return b.buf[i]
		}
		return 0
	}
	if uint(i) < uint(len(b.words)) {
		return b.words[i]
	}
	return 0
}

//go:noinline
func (b *Ext_java_util_BitSet) bsPromoteToHeap(minWords int) {
	if b.hiLen > 0 {
		hw := int(b.hi[b.hiLen-1]>>6) + 1
		if hw > minWords {
			minWords = hw
		}
	}
	if minWords < bsInlineWords {
		minWords = bsInlineWords
	}
	newCap := max(minWords, max(16, int(b.initW)))
	w := make([]uint64, minWords, newCap)
	copy(w[:bsInlineWords], b.buf[:])
	b.buf = [bsInlineWords]uint64{}
	for i := int32(0); i < b.hiLen; i++ {
		hb := int(b.hi[i])
		b.hi[i] = 0
		w[hb>>6] |= uint64(1) << (uint(hb) & 63)
	}
	b.hiLen = 0
	b.words = w
}

func (b *Ext_java_util_BitSet) bsMaterialize() {
	if b.hiLen > 0 && b.words == nil {
		b.bsPromoteToHeap(int(b.hi[b.hiLen-1]>>6) + 1)
	}
}

//go:noinline
func (b *Ext_java_util_BitSet) bsEnsure(wordIdx int) {
	if b.words == nil {
		b.bsPromoteToHeap(wordIdx + 1)
		return
	}
	if wordIdx < len(b.words) {
		return
	}
	if wordIdx < cap(b.words) {
		prev := len(b.words)
		b.words = b.words[:wordIdx+1]
		clear(b.words[prev:])
		return
	}
	newCap := max(wordIdx+1, max(2*cap(b.words), int(b.initW)))
	newWords := make([]uint64, wordIdx+1, newCap)
	copy(newWords, b.words)
	b.words = newWords
}

func (b *Ext_java_util_BitSet) bsUsedWords() int {
	if b.words == nil {
		for n := bsInlineWords; n > 0; n-- {
			if b.buf[n-1] != 0 {
				return n
			}
		}
		return 0
	}
	n := len(b.words)
	for n > 0 && b.words[n-1] == 0 {
		n--
	}
	return n
}

func (b *Ext_java_util_BitSet) bsTrim() {
	if b.words != nil {
		b.words = b.words[:b.bsUsedWords()]
	}
}

func (b *Ext_java_util_BitSet) bsLength() int {
	if b.hiLen > 0 {
		return int(b.hi[b.hiLen-1]) + 1
	}
	n := b.bsUsedWords()
	if n == 0 {
		return 0
	}
	last := b.bsWordAt(n - 1)
	return 64*(n-1) + (64 - bits.LeadingZeros64(last))
}

func (b *Ext_java_util_BitSet) bsSize() int {
	w := max(b.bsUsedWords(), int(b.initW))
	if cap(b.words) > bsInlineWords {
		w = max(w, cap(b.words))
	}
	if b.hiLen > 0 {
		w = max(w, int(b.hi[b.hiLen-1]>>6)+1)
	}
	return w * 64
}

func (b *Ext_java_util_BitSet) bsSet(bitIndex int) {
	if bitIndex < 0 {
		bitSetThrowNegIdx(bitIndex, "bitIndex")
	}
	w := bitIndex >> 6
	mask := uint64(1) << (uint(bitIndex) & 63)
	if b.words == nil {
		if w < bsInlineWords {
			b.buf[w] |= mask
			return
		}
		if b.bsHiInsert(int32(bitIndex)) {
			return
		}
		b.bsPromoteToHeap(w + 1)
		b.words[w] |= mask
		return
	}
	if w >= len(b.words) {
		b.bsEnsure(w)
	}
	b.words[w] |= mask
}

func (b *Ext_java_util_BitSet) bsClearBit(bitIndex int) {
	if bitIndex < 0 {
		bitSetThrowNegIdx(bitIndex, "bitIndex")
	}
	w := bitIndex >> 6
	if b.words == nil {
		if w < bsInlineWords {
			b.buf[w] &^= uint64(1) << (uint(bitIndex) & 63)
		} else if b.hiLen > 0 {
			b.bsHiRemove(int32(bitIndex))
		}
		return
	}
	if w < len(b.words) {
		b.words[w] &^= uint64(1) << (uint(bitIndex) & 63)
		b.bsTrim()
	}
}

func (b *Ext_java_util_BitSet) bsSetRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	if toIndex-1 == fromIndex {
		b.bsSet(fromIndex)
		return
	}
	if b.words == nil && fromIndex >= bsInlineBits && toIndex-fromIndex <= bsSparseCap {
		for i := fromIndex; i < toIndex; i++ {
			b.bsSet(i)
		}
		return
	}
	startWord := fromIndex >> 6
	endWord := (toIndex - 1) >> 6
	firstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
	lastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
	if b.words == nil && endWord < bsInlineWords {
		if startWord == endWord {
			b.buf[startWord] |= firstWordMask & lastWordMask
		} else {
			b.buf[startWord] |= firstWordMask
			for i := startWord + 1; i < endWord; i++ {
				b.buf[i] = ^uint64(0)
			}
			b.buf[endWord] |= lastWordMask
		}
		return
	}
	b.bsEnsure(endWord)
	if startWord == endWord {
		b.words[startWord] |= firstWordMask & lastWordMask
	} else {
		b.words[startWord] |= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			b.words[i] = ^uint64(0)
		}
		b.words[endWord] |= lastWordMask
	}
}

func (b *Ext_java_util_BitSet) bsClearRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	startWord := fromIndex >> 6
	endWord := (toIndex - 1) >> 6
	firstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
	lastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
	if b.words == nil {
		if b.hiLen > 0 && toIndex > bsInlineBits {
			wIdx := 0
			n := int(b.hiLen)
			for i := 0; i < n; i++ {
				hb := b.hi[i]
				if int(hb) < fromIndex || int(hb) >= toIndex {
					b.hi[wIdx] = hb
					wIdx++
				}
			}
			for i := wIdx; i < n; i++ {
				b.hi[i] = 0
			}
			b.hiLen = int32(wIdx)
		}
		if startWord >= bsInlineWords {
			return
		}
		if endWord >= bsInlineWords {
			endWord = bsInlineWords - 1
			lastWordMask = ^uint64(0)
		}
		if startWord == endWord {
			b.buf[startWord] &^= firstWordMask & lastWordMask
		} else {
			b.buf[startWord] &^= firstWordMask
			for i := startWord + 1; i < endWord; i++ {
				b.buf[i] = 0
			}
			b.buf[endWord] &^= lastWordMask
		}
		return
	}
	if startWord >= len(b.words) {
		return
	}
	if endWord >= len(b.words) {
		endWord = len(b.words) - 1
		lastWordMask = ^uint64(0)
	}
	if startWord == endWord {
		b.words[startWord] &^= firstWordMask & lastWordMask
	} else {
		b.words[startWord] &^= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			b.words[i] = 0
		}
		b.words[endWord] &^= lastWordMask
	}
	b.bsTrim()
}

func (b *Ext_java_util_BitSet) bsNextSetBit(fromIndex int) int {
	if fromIndex < 0 {
		bitSetThrowNegIdx(fromIndex, "fromIndex")
	}
	u := fromIndex >> 6
	if b.words == nil {
		if u < bsInlineWords {
			word := b.buf[u] & (^uint64(0) << (uint(fromIndex) & 63))
			for {
				if word != 0 {
					return u*64 + bits.TrailingZeros64(word)
				}
				u++
				if u >= bsInlineWords {
					break
				}
				word = b.buf[u]
			}
		}
		n := int(b.hiLen)
		for i := 0; i < n; i++ {
			if hb := int(b.hi[i]); hb >= fromIndex {
				return hb
			}
		}
		return -1
	}
	if u >= len(b.words) {
		return -1
	}
	word := b.words[u] & (^uint64(0) << (uint(fromIndex) & 63))
	for {
		if word != 0 {
			return u*64 + bits.TrailingZeros64(word)
		}
		u++
		if u >= len(b.words) {
			return -1
		}
		word = b.words[u]
	}
}

func bitSetOp_ToByteArray(b *Ext_java_util_BitSet) []int8 {
	if b == nil {
		ThrowNPE()
	}
	length := b.bsLength()
	if length == 0 {
		return []int8{}
	}
	nBytes := (length + 7) / 8
	out := make([]int8, nBytes)
	for i := b.bsNextSetBit(0); i >= 0; i = b.bsNextSetBit(i + 1) {
		out[i/8] |= 1 << uint(i%8)
	}
	return out
}
func bitSetOp_ToLongArray(b *Ext_java_util_BitSet) []int64 {
	if b == nil {
		ThrowNPE()
	}
	if b.hiLen > 0 {
		hw := int(b.hi[b.hiLen-1]) >> 6
		out := make([]int64, hw+1)
		n := b.bsUsedWords()
		for i := 0; i < n; i++ {
			out[i] = int64(b.bsWordAt(i))
		}
		for i := int32(0); i < b.hiLen; i++ {
			hb := int(b.hi[i])
			out[hb>>6] |= int64(uint64(1) << (uint(hb) & 63))
		}
		return out
	}
	n := b.bsUsedWords()
	if n == 0 {
		return []int64{}
	}
	out := make([]int64, n)
	for i := 0; i < n; i++ {
		out[i] = int64(b.bsWordAt(i))
	}
	return out
}
func bitSetOp_Flip_1(b *Ext_java_util_BitSet, bitIndex int32) {
	if b == nil {
		ThrowNPE()
	}
	i := int(bitIndex)
	if i < 0 {
		bitSetThrowNegIdx(i, "bitIndex")
	}
	w := i >> 6
	mask := uint64(1) << (uint(i) & 63)
	if b.words == nil {
		if w < bsInlineWords {
			b.buf[w] ^= mask
			return
		}
		if b.bsHiRemove(bitIndex) {
			return
		}
		if b.bsHiInsert(bitIndex) {
			return
		}
	}
	b.bsEnsure(w)
	b.words[w] ^= mask
	b.bsTrim()
}
func bitSetOp_Flip_2(b *Ext_java_util_BitSet, from, to int32) {
	if b == nil {
		ThrowNPE()
	}
	fromIndex, toIndex := int(from), int(to)
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	if toIndex-1 == fromIndex {
		bitSetOp_Flip_1(b, from)
		return
	}
	if b.words == nil && fromIndex >= bsInlineBits && toIndex-fromIndex <= bsSparseCap {
		for idx := from; idx < to; idx++ {
			bitSetOp_Flip_1(b, idx)
		}
		return
	}
	startWord := fromIndex >> 6
	endWord := (toIndex - 1) >> 6
	firstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
	lastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
	if b.words == nil && endWord < bsInlineWords {
		if startWord == endWord {
			b.buf[startWord] ^= firstWordMask & lastWordMask
		} else {
			b.buf[startWord] ^= firstWordMask
			for i := startWord + 1; i < endWord; i++ {
				b.buf[i] = ^b.buf[i]
			}
			b.buf[endWord] ^= lastWordMask
		}
		return
	}
	b.bsEnsure(endWord)
	if startWord == endWord {
		b.words[startWord] ^= firstWordMask & lastWordMask
	} else {
		b.words[startWord] ^= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			b.words[i] = ^b.words[i]
		}
		b.words[endWord] ^= lastWordMask
	}
	b.bsTrim()
}
func bitSetOp_Set_1(b *Ext_java_util_BitSet, i int32) {
	if b == nil {
		ThrowNPE()
	}
	b.bsSet(int(i))
}
func bitSetOp_Set_2_int32_bool(b *Ext_java_util_BitSet, i int32, v bool) {
	if b == nil {
		ThrowNPE()
	}
	if v {
		b.bsSet(int(i))
	} else {
		b.bsClearBit(int(i))
	}
}
func bitSetOp_Set_2_int32_int32(b *Ext_java_util_BitSet, f, t int32) {
	if b == nil {
		ThrowNPE()
	}
	b.bsSetRange(int(f), int(t))
}
func bitSetOp_Set_3(b *Ext_java_util_BitSet, f, t int32, v bool) {
	if b == nil {
		ThrowNPE()
	}
	if v {
		b.bsSetRange(int(f), int(t))
	} else {
		b.bsClearRange(int(f), int(t))
	}
}
func bitSetOp_Clear(b *Ext_java_util_BitSet) {
	if b == nil {
		ThrowNPE()
	}
	if b.words == nil {
		b.buf = [bsInlineWords]uint64{}
	} else {
		b.words = b.words[:0]
	}
	if b.hiLen > 0 {
		clear(b.hi[:b.hiLen])
		b.hiLen = 0
	}
}
func bitSetOp_Clear_0(b *Ext_java_util_BitSet) { bitSetOp_Clear(b) }
func bitSetOp_Clear_1(b *Ext_java_util_BitSet, i int32) {
	if b == nil {
		ThrowNPE()
	}
	b.bsClearBit(int(i))
}
func bitSetOp_Clear_2(b *Ext_java_util_BitSet, f, t int32) {
	if b == nil {
		ThrowNPE()
	}
	b.bsClearRange(int(f), int(t))
}
func bitSetOp_Get_1(b *Ext_java_util_BitSet, i int32) bool {
	if b == nil {
		ThrowNPE()
	}
	idx := int(i)
	if idx < 0 {
		bitSetThrowNegIdx(idx, "bitIndex")
	}
	w := idx >> 6
	if b.words == nil {
		if w < bsInlineWords {
			return (b.buf[w] & (uint64(1) << (uint(idx) & 63))) != 0
		}
		return b.hiLen > 0 && b.bsHiContains(i)
	}
	return w < len(b.words) && (b.words[w]&(uint64(1)<<(uint(idx)&63))) != 0
}
func bitSetOp_Get(b *Ext_java_util_BitSet, i int32) bool { return bitSetOp_Get_1(b, i) }
func bitSetOp_Get_2(b *Ext_java_util_BitSet, f, t int32) *Ext_java_util_BitSet {
	if b == nil {
		ThrowNPE()
	}
	fromIndex, toIndex := int(f), int(t)
	bitSetCheckRange(fromIndex, toIndex)
	dst := &Ext_java_util_BitSet{initW: 1}
	dst.BindExtSelf_(dst)
	length := b.bsLength()
	if length <= fromIndex || fromIndex == toIndex {
		return dst
	}
	if toIndex > length {
		toIndex = length
	}
	for i := b.bsNextSetBit(fromIndex); i >= 0 && i < toIndex; i = b.bsNextSetBit(i + 1) {
		dst.bsSet(i - fromIndex)
	}
	return dst
}
func bitSetOp_NextSetBit(b *Ext_java_util_BitSet, f int32) int32 {
	if b == nil {
		ThrowNPE()
	}
	return int32(b.bsNextSetBit(int(f)))
}
func bitSetOp_NextClearBit(b *Ext_java_util_BitSet, f int32) int32 {
	if b == nil {
		ThrowNPE()
	}
	fromIndex := int(f)
	if fromIndex < 0 {
		bitSetThrowNegIdx(fromIndex, "fromIndex")
	}
	u := fromIndex >> 6
	nLim := len(b.words)
	if b.words == nil {
		nLim = bsInlineWords
	}
	var cand int
	if u >= nLim {
		cand = fromIndex
	} else {
		word := ^b.bsWordAt(u) & (^uint64(0) << (uint(fromIndex) & 63))
		for {
			if word != 0 {
				cand = u*64 + bits.TrailingZeros64(word)
				break
			}
			u++
			if u >= nLim {
				cand = nLim * 64
				break
			}
			word = ^b.bsWordAt(u)
		}
	}
	if b.hiLen > 0 && cand >= bsInlineBits {
		n := int(b.hiLen)
		for i := 0; i < n; i++ {
			hb := int(b.hi[i])
			if hb == cand {
				cand++
			} else if hb > cand {
				break
			}
		}
	}
	return int32(cand)
}
func bitSetOp_PreviousSetBit(b *Ext_java_util_BitSet, f int32) int32 {
	if b == nil {
		ThrowNPE()
	}
	fromIndex := int(f)
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		bitSetThrowNegIdx(fromIndex, "fromIndex")
	}
	if b.hiLen > 0 {
		for i := int(b.hiLen) - 1; i >= 0; i-- {
			if hb := int(b.hi[i]); hb <= fromIndex {
				return int32(hb)
			}
		}
	}
	nUsed := b.bsUsedWords()
	if nUsed == 0 {
		return -1
	}
	w := fromIndex >> 6
	if w >= nUsed {
		last := b.bsWordAt(nUsed - 1)
		return int32(64*(nUsed-1) + (64 - bits.LeadingZeros64(last)) - 1)
	}
	word := b.bsWordAt(w) & (^uint64(0) >> (63 - (uint(fromIndex) & 63)))
	for {
		if word != 0 {
			return int32((w+1)*64 - 1 - bits.LeadingZeros64(word))
		}
		if w == 0 {
			return -1
		}
		w--
		word = b.bsWordAt(w)
	}
}
func bitSetOp_PreviousClearBit(b *Ext_java_util_BitSet, f int32) int32 {
	if b == nil {
		ThrowNPE()
	}
	fromIndex := int(f)
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		bitSetThrowNegIdx(fromIndex, "fromIndex")
	}
	if b.hiLen > 0 && fromIndex >= bsInlineBits {
		for i := int(b.hiLen) - 1; i >= 0; i-- {
			hb := int(b.hi[i])
			if hb == fromIndex {
				fromIndex--
			} else if hb < fromIndex {
				break
			}
		}
		if fromIndex >= bsInlineBits {
			return int32(fromIndex)
		}
	}
	nUsed := b.bsUsedWords()
	w := fromIndex >> 6
	if w >= nUsed {
		return int32(fromIndex)
	}
	word := ^b.bsWordAt(w) & (^uint64(0) >> (63 - (uint(fromIndex) & 63)))
	for {
		if word != 0 {
			return int32((w+1)*64 - 1 - bits.LeadingZeros64(word))
		}
		if w == 0 {
			return -1
		}
		w--
		word = ^b.bsWordAt(w)
	}
}
func bitSetOp_Length(b *Ext_java_util_BitSet) int32 {
	if b == nil {
		ThrowNPE()
	}
	return int32(b.bsLength())
}
func bitSetOp_IsEmpty(b *Ext_java_util_BitSet) bool {
	if b == nil {
		ThrowNPE()
	}
	if b.hiLen > 0 {
		return false
	}
	return b.bsUsedWords() == 0
}
func bitSetOp_Intersects(b, o *Ext_java_util_BitSet) bool {
	if b == nil || o == nil {
		ThrowNPE()
	}
	n := min(b.bsUsedWords(), o.bsUsedWords())
	for i := 0; i < n; i++ {
		if b.bsWordAt(i)&o.bsWordAt(i) != 0 {
			return true
		}
	}
	for i := int32(0); i < b.hiLen; i++ {
		if bitSetOp_Get_1(o, b.hi[i]) {
			return true
		}
	}
	if b.words != nil {
		for i := int32(0); i < o.hiLen; i++ {
			if bitSetOp_Get_1(b, o.hi[i]) {
				return true
			}
		}
	}
	return false
}
func bitSetOp_Cardinality(b *Ext_java_util_BitSet) int32 {
	if b == nil {
		ThrowNPE()
	}
	cnt := int(b.hiLen)
	if b.words == nil {
		for i := 0; i < bsInlineWords; i++ {
			cnt += bits.OnesCount64(b.buf[i])
		}
	} else {
		for _, w := range b.words {
			cnt += bits.OnesCount64(w)
		}
	}
	return int32(cnt)
}
func bitSetOp_And(b, o *Ext_java_util_BitSet) {
	if b == nil || o == nil {
		ThrowNPE()
	}
	if b == o {
		return
	}
	if b.words == nil {
		for i := 0; i < bsInlineWords; i++ {
			b.buf[i] &= o.bsWordAt(i)
		}
		if b.hiLen > 0 {
			wIdx := 0
			n := int(b.hiLen)
			for i := 0; i < n; i++ {
				hb := b.hi[i]
				if bitSetOp_Get_1(o, hb) {
					b.hi[wIdx] = hb
					wIdx++
				}
			}
			for i := wIdx; i < n; i++ {
				b.hi[i] = 0
			}
			b.hiLen = int32(wIdx)
		}
		return
	}
	var keepHi [bsSparseCap]int32
	keepCnt := 0
	for i := int32(0); i < o.hiLen; i++ {
		hb := o.hi[i]
		if bitSetOp_Get_1(b, hb) {
			keepHi[keepCnt] = hb
			keepCnt++
		}
	}
	n := min(len(b.words), o.bsUsedWords())
	for i := 0; i < n; i++ {
		b.words[i] &= o.bsWordAt(i)
	}
	b.words = b.words[:n]
	b.bsTrim()
	for i := 0; i < keepCnt; i++ {
		b.bsSet(int(keepHi[i]))
	}
}
func bitSetOp_Or(b, o *Ext_java_util_BitSet) {
	if b == nil || o == nil {
		ThrowNPE()
	}
	if b == o {
		return
	}
	nO := o.bsUsedWords()
	if nO > 0 {
		if b.words == nil && nO <= bsInlineWords {
			for i := 0; i < nO; i++ {
				b.buf[i] |= o.bsWordAt(i)
			}
		} else {
			b.bsEnsure(nO - 1)
			for i := 0; i < nO; i++ {
				b.words[i] |= o.bsWordAt(i)
			}
		}
	}
	for i := int32(0); i < o.hiLen; i++ {
		b.bsSet(int(o.hi[i]))
	}
}
func bitSetOp_Xor(b, o *Ext_java_util_BitSet) {
	if b == nil || o == nil {
		ThrowNPE()
	}
	if b == o {
		bitSetOp_Clear(b)
		return
	}
	nO := o.bsUsedWords()
	if nO > 0 {
		if b.words == nil && nO <= bsInlineWords {
			for i := 0; i < nO; i++ {
				b.buf[i] ^= o.bsWordAt(i)
			}
		} else {
			b.bsEnsure(nO - 1)
			for i := 0; i < nO; i++ {
				b.words[i] ^= o.bsWordAt(i)
			}
			b.bsTrim()
		}
	}
	for i := int32(0); i < o.hiLen; i++ {
		bitSetOp_Flip_1(b, o.hi[i])
	}
}
func bitSetOp_AndNot(b, o *Ext_java_util_BitSet) {
	if b == nil || o == nil {
		ThrowNPE()
	}
	if b == o {
		bitSetOp_Clear(b)
		return
	}
	if b.words == nil {
		for i := 0; i < bsInlineWords; i++ {
			b.buf[i] &^= o.bsWordAt(i)
		}
		if b.hiLen > 0 {
			wIdx := 0
			n := int(b.hiLen)
			for i := 0; i < n; i++ {
				hb := b.hi[i]
				if !bitSetOp_Get_1(o, hb) {
					b.hi[wIdx] = hb
					wIdx++
				}
			}
			for i := wIdx; i < n; i++ {
				b.hi[i] = 0
			}
			b.hiLen = int32(wIdx)
		}
		return
	}
	n := min(len(b.words), o.bsUsedWords())
	for i := 0; i < n; i++ {
		b.words[i] &^= o.bsWordAt(i)
	}
	for i := int32(0); i < o.hiLen; i++ {
		b.bsClearBit(int(o.hi[i]))
	}
	b.bsTrim()
}
func bitSetOp_HashCode(b *Ext_java_util_BitSet) int32 {
	if b == nil {
		ThrowNPE()
	}
	h := int64(1234)
	for i := int(b.hiLen) - 1; i >= 0; {
		hw := int(b.hi[i]) >> 6
		var hword uint64
		for i >= 0 && (int(b.hi[i])>>6) == hw {
			hword |= uint64(1) << (uint(b.hi[i]) & 63)
			i--
		}
		h ^= int64(hword) * int64(hw+1)
	}
	n := b.bsUsedWords()
	for i := n - 1; i >= 0; i-- {
		h ^= int64(b.bsWordAt(i)) * int64(i+1)
	}
	return int32((h >> 32) ^ h)
}
func bitSetOp_Size(b *Ext_java_util_BitSet) int32 {
	if b == nil {
		ThrowNPE()
	}
	return int32(b.bsSize())
}
func bitSetEqualsInlineVsHeap(inl, hp *Ext_java_util_BitSet) bool {
	nHp := hp.bsUsedWords()
	maxHW := int(inl.hi[inl.hiLen-1]) >> 6
	if nHp != maxHW+1 {
		return false
	}
	for i := 0; i < bsInlineWords; i++ {
		if inl.buf[i] != hp.bsWordAt(i) {
			return false
		}
	}
	curW := bsInlineWords
	n := int(inl.hiLen)
	for idx := 0; idx < n; {
		hw := int(inl.hi[idx]) >> 6
		for curW < hw {
			if hp.words[curW] != 0 {
				return false
			}
			curW++
		}
		var hword uint64
		for idx < n && (int(inl.hi[idx])>>6) == hw {
			hword |= uint64(1) << (uint(inl.hi[idx]) & 63)
			idx++
		}
		if hp.words[hw] != hword {
			return false
		}
		curW = hw + 1
	}
	return true
}
func bitSetEqualsPtr(b, ob *Ext_java_util_BitSet) bool {
	if b.words == nil && ob.words == nil {
		return b.buf == ob.buf && b.hiLen == ob.hiLen && b.hi == ob.hi
	}
	if b.hiLen == 0 && ob.hiLen == 0 {
		nB := b.bsUsedWords()
		nO := ob.bsUsedWords()
		if nB != nO {
			return false
		}
		for i := 0; i < nB; i++ {
			if b.bsWordAt(i) != ob.bsWordAt(i) {
				return false
			}
		}
		return true
	}
	if b.hiLen > 0 {
		return bitSetEqualsInlineVsHeap(b, ob)
	}
	return bitSetEqualsInlineVsHeap(ob, b)
}
func bitSetOp_Equals(b *Ext_java_util_BitSet, o any) bool {
	if b == nil {
		ThrowNPE()
	}
	ob, ok := o.(*Ext_java_util_BitSet)
	if !ok || ob == nil {
		return false
	}
	return b == ob || bitSetEqualsPtr(b, ob)
}
func bitSetOp_Clone(b *Ext_java_util_BitSet) any {
	if b == nil {
		ThrowNPE()
	}
	n := b.bsUsedWords()
	c := &Ext_java_util_BitSet{
		initW: b.initW,
		hiLen: b.hiLen,
		hi:    b.hi,
	}
	c.BindExtSelf_(c)
	if n <= bsInlineWords {
		if b.words == nil {
			c.buf = b.buf
		} else {
			copy(c.buf[:], b.words[:n])
		}
	} else {
		c.words = make([]uint64, n)
		copy(c.words, b.words[:n])
	}
	return c
}
func bitSetOp_String(b *Ext_java_util_BitSet) string {
	if b == nil {
		ThrowNPE()
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for i := b.bsNextSetBit(0); i >= 0; i = b.bsNextSetBit(i + 1) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(strconv.Itoa(i))
	}
	sb.WriteByte('}')
	return sb.String()
}

// --- Iterables, Iterators, AbstractIterator, PeekingIterator, Ordering, ComparisonChain, Range ---

func iterablesAddAll(dst Ext_java_util_Collection, src Ext_java_lang_Iterable) bool {
	CheckNotNull(dst)
	CheckNotNull(src)
	added := false
	for it := Iterate(src); it.HasNext(); {
		if dst.Add_1(it.Next()) {
			added = true
		}
	}
	return added
}

type jFuncIterable struct {
	makeIter func() Ext_java_util_Iterator
}

func (*jFuncIterable) Implements_j2gort_Ext_java_lang_Iterable() {}
func (f *jFuncIterable) Iterator() any                           { return f.makeIter() }
func (f *jFuncIterable) ForEach(c Ext_java_util_function_Consumer) {
	ExtDefault_Iterable_ForEach(f, c)
}
func (f *jFuncIterable) String() string {
	var items []any
	for it := f.makeIter(); it.HasNext(); {
		items = append(items, it.Next())
	}
	return sliceJavaString(items)
}

func iterablesConcatList(list []Ext_java_lang_Iterable) Ext_java_lang_Iterable {
	for _, x := range list {
		CheckNotNull(x)
	}
	return &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
		return &jConcatIter{outer: newSliceObjIter(toAnySliceIterables(list))}
	}}
}

func toAnySliceIterables(list []Ext_java_lang_Iterable) []any {
	out := make([]any, len(list))
	for i, v := range list {
		out[i] = v
	}
	return out
}

func iterablesConcatMany(outer Ext_java_lang_Iterable) Ext_java_lang_Iterable {
	CheckNotNull(outer)
	return &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
		return &jConcatIter{outer: Iterate(outer)}
	}}
}

type jConcatIter struct {
	outer SimpleIter
	cur   SimpleIter
}

func (*jConcatIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jConcatIter) HasNext() bool {
	for it.cur == nil || !it.cur.HasNext() {
		if !it.outer.HasNext() {
			return false
		}
		nextIt := CheckNotNull(it.outer.Next())
		it.cur = Iterate(nextIt)
	}
	return true
}
func (it *jConcatIter) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	return it.cur.Next()
}
func (it *jConcatIter) Remove() {
	if r, ok := it.cur.(Ext_java_util_Iterator); ok {
		r.Remove()
		return
	}
	panic(javalang.NewUnsupportedOperationException("remove"))
}

func iterablesContains(it Ext_java_lang_Iterable, target any) bool {
	CheckNotNull(it)
	if c, ok := any(it).(Ext_java_util_Collection); ok {
		return c.Contains(target)
	}
	t := ToAny(target)
	for iter := Iterate(it); iter.HasNext(); {
		if javalang.ObjectsEquals(t, ToAny(iter.Next())) {
			return true
		}
	}
	return false
}

func iterablesFilter(it Ext_java_lang_Iterable, pred Ext_com_google_common_base_Predicate) Ext_java_lang_Iterable {
	CheckNotNull(it)
	CheckNotNull(pred)
	return &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
		return iteratorsFilterSimple(Iterate(it), pred)
	}}
}

func iteratorsFilter(it Ext_java_util_Iterator, pred Ext_com_google_common_base_Predicate) Ext_com_google_common_collect_UnmodifiableIterator {
	CheckNotNull(it)
	CheckNotNull(pred)
	return iteratorsFilterSimple(it, pred)
}

type jFilterIter struct {
	src      SimpleIter
	pred     Ext_com_google_common_base_Predicate
	haveNext bool
	nextVal  any
}

func iteratorsFilterSimple(src SimpleIter, pred Ext_com_google_common_base_Predicate) *jFilterIter {
	return &jFilterIter{src: src, pred: pred}
}

func (*jFilterIter) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*jFilterIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*jFilterIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (it *jFilterIter) HasNext() bool {
	if it.haveNext {
		return true
	}
	for it.src.HasNext() {
		v := it.src.Next()
		if it.pred.Apply(v) {
			it.nextVal = v
			it.haveNext = true
			return true
		}
	}
	return false
}
func (it *jFilterIter) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	it.haveNext = false
	v := it.nextVal
	it.nextVal = nil
	return v
}
func (*jFilterIter) Remove() { panic(javalang.NewUnsupportedOperationException("remove")) }

func iterablesGet(it Ext_java_lang_Iterable, idx int32) any {
	CheckNotNull(it)
	if idx < 0 {
		panic(javalang.NewIndexOutOfBoundsException("negative index"))
	}
	if l, ok := any(it).(Ext_java_util_List); ok {
		return l.Get(idx)
	}
	i := int32(0)
	for iter := Iterate(it); iter.HasNext(); i++ {
		v := iter.Next()
		if i == idx {
			return v
		}
	}
	panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(idx), int(i)))
}

func iterablesGetFirst(it Ext_java_lang_Iterable, def any) any {
	CheckNotNull(it)
	iter := NewIterCursor(it)
	if iter.HasNext() {
		return iter.Next()
	}
	return ToAny(def)
}

func iterablesGetLast(it Ext_java_lang_Iterable) any {
	CheckNotNull(it)
	if l, ok := any(it).(Ext_java_util_List); ok {
		if l.IsEmpty() {
			panic(javalang.NewNoSuchElementException("empty"))
		}
		return l.Get(l.Size() - 1)
	}
	iter := NewIterCursor(it)
	if !iter.HasNext() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	last := iter.Next()
	for iter.HasNext() {
		last = iter.Next()
	}
	return last
}

func iterablesGetOnlyElement(it Ext_java_lang_Iterable) any {
	CheckNotNull(it)
	iter := NewIterCursor(it)
	if !iter.HasNext() {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	first := iter.Next()
	if iter.HasNext() {
		panic(javalang.NewIllegalArgumentException("expected one element but was multiple"))
	}
	return first
}

func iterablesSkip(it Ext_java_lang_Iterable, n int32) Ext_java_lang_Iterable {
	CheckNotNull(it)
	CheckArgument(n >= 0)
	return &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
		iter := Iterate(it)
		for i := int32(0); i < n && iter.HasNext(); i++ {
			_ = iter.Next()
		}
		return &jSimpleIterAdapter{it: iter}
	}}
}

func iterablesTransform(it Ext_java_lang_Iterable, fn Ext_com_google_common_base_Function) Ext_java_lang_Iterable {
	CheckNotNull(it)
	CheckNotNull(fn)
	return &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
		return &jTransformIter{src: Iterate(it), fn: fn}
	}}
}

type jSimpleIterAdapter struct {
	it SimpleIter
}

func (*jSimpleIterAdapter) Implements_j2gort_Ext_java_util_Iterator() {}
func (a *jSimpleIterAdapter) HasNext() bool                           { return a.it.HasNext() }
func (a *jSimpleIterAdapter) Next() any                               { return a.it.Next() }
func (a *jSimpleIterAdapter) Remove() {
	if r, ok := a.it.(Ext_java_util_Iterator); ok {
		r.Remove()
		return
	}
	panic(javalang.NewUnsupportedOperationException("remove"))
}

type jTransformIter struct {
	src SimpleIter
	fn  Ext_com_google_common_base_Function
}

func (*jTransformIter) Implements_j2gort_Ext_java_util_Iterator() {}
func (t *jTransformIter) HasNext() bool                           { return t.src.HasNext() }
func (t *jTransformIter) Next() any                               { return t.fn.Apply(t.src.Next()) }
func (t *jTransformIter) Remove() {
	if r, ok := t.src.(Ext_java_util_Iterator); ok {
		r.Remove()
		return
	}
	panic(javalang.NewUnsupportedOperationException("remove"))
}

type jPeekingIter struct {
	src       Ext_java_util_Iterator
	hasPeeked bool
	peeked    any
}

func newPeekingIter(src Ext_java_util_Iterator) Ext_com_google_common_collect_PeekingIterator {
	CheckNotNull(src)
	if p, ok := src.(Ext_com_google_common_collect_PeekingIterator); ok {
		return p
	}
	return &jPeekingIter{src: src}
}

func (*jPeekingIter) Implements_j2gort_Ext_java_util_Iterator()                        {}
func (*jPeekingIter) Implements_j2gort_Ext_com_google_common_collect_PeekingIterator() {}
func (p *jPeekingIter) HasNext() bool                                                  { return p.hasPeeked || p.src.HasNext() }
func (p *jPeekingIter) Next() any {
	if !p.hasPeeked {
		return p.src.Next()
	}
	v := p.peeked
	p.hasPeeked = false
	p.peeked = nil
	return v
}
func (p *jPeekingIter) Peek() any {
	if !p.hasPeeked {
		p.peeked = p.src.Next()
		p.hasPeeked = true
	}
	return p.peeked
}
func (p *jPeekingIter) Remove() {
	if p.hasPeeked {
		panic(javalang.NewIllegalStateException("cannot remove after peek"))
	}
	p.src.Remove()
}

type jAbstractIterState struct {
	state   int // 0=NOT_READY, 1=READY, 2=DONE, 3=FAILED
	nextVal any
}

func abstractIterHasNext(b *Ext_com_google_common_collect_AbstractIteratorBase) bool {
	if b.Core_ == nil {
		b.Core_ = &jAbstractIterState{}
	}
	st := b.Core_.(*jAbstractIterState)
	if st.state == 3 {
		panic(javalang.NewIllegalStateException("failed iterator"))
	}
	if st.state == 2 {
		return false
	}
	if st.state == 1 {
		return true
	}
	st.state = 3
	st.nextVal = b.Self_.ComputeNext()
	if st.state != 2 {
		st.state = 1
		return true
	}
	return false
}

func abstractIterNext(b *Ext_com_google_common_collect_AbstractIteratorBase) any {
	if !b.Self_.HasNext() {
		panic(javalang.NewNoSuchElementException("exhausted"))
	}
	st := b.Core_.(*jAbstractIterState)
	st.state = 0
	v := st.nextVal
	st.nextVal = nil
	return v
}

func newOrdering(cmp func(a, b any) int32) Ext_com_google_common_collect_Ordering {
	o := &Ext_com_google_common_collect_OrderingBase{}
	o.BindExtSelf_(o)
	o.Core_ = cmp
	return o
}

type jComparisonChain struct {
	res int32
}

var (
	comparisonChainActive  = newComparisonChainObj(0)
	comparisonChainLess    = newComparisonChainObj(-1)
	comparisonChainGreater = newComparisonChainObj(1)
)

func newComparisonChainObj(res int32) Ext_com_google_common_collect_ComparisonChain {
	c := &Ext_com_google_common_collect_ComparisonChainBase{}
	c.BindExtSelf_(c)
	c.Core_ = &jComparisonChain{res: res}
	return c
}

func (c *jComparisonChain) cmpInt(a, b int32) Ext_com_google_common_collect_ComparisonChain {
	if c.res != 0 {
		if c.res < 0 {
			return comparisonChainLess
		}
		return comparisonChainGreater
	}
	if a < b {
		return comparisonChainLess
	}
	if a > b {
		return comparisonChainGreater
	}
	return comparisonChainActive
}

func (c *jComparisonChain) cmpAny(a, b any) Ext_com_google_common_collect_ComparisonChain {
	if c.res != 0 {
		if c.res < 0 {
			return comparisonChainLess
		}
		return comparisonChainGreater
	}
	r := CompareValues(a, b)
	if r < 0 {
		return comparisonChainLess
	}
	if r > 0 {
		return comparisonChainGreater
	}
	return comparisonChainActive
}

type jRangeCore struct {
	lower any
	upper any
}

type jRangeSetCore struct {
	ranges []*jRangeCore
}

func collectionsShuffle(lst Ext_java_util_List, rnd *Ext_java_util_Random) {
	CheckNotNull(lst)
	CheckNotNullPtr(rnd)
	if rnd.Core_ == nil {
		rnd.Core_ = javautil.NewRandom()
	}
	r := rnd.Core_.(*javautil.Random)
	for i := lst.Size(); i > 1; i-- {
		j := r.NextIntBound(i)
		tmp := lst.Get(i - 1)
		lst.Set(i-1, lst.Get(j))
		lst.Set(j, tmp)
	}
}

// --- Streams, IntStream, Collectors ---

type jStream struct {
	next func() (any, bool)
}

var emptyStream_ = &jStream{next: func() (any, bool) { return nil, false }}

func streamFromSlice(items []any) Ext_java_util_stream_Stream {
	if len(items) == 0 {
		return emptyStream_
	}
	idx := 0
	return &jStream{next: func() (any, bool) {
		if idx >= len(items) {
			return nil, false
		}
		v := items[idx]
		idx++
		return ToAny(v), true
	}}
}

func streamFromIterable(it Ext_java_lang_Iterable) Ext_java_util_stream_Stream {
	CheckNotNull(it)
	if si, ok := any(it).(SliceIterable); ok {
		return streamFromSlice(si.IterSlice_())
	}
	var cur JIterCursor
	inited := false
	return &jStream{next: func() (any, bool) {
		if !inited {
			cur = NewIterCursor(it)
			inited = true
		}
		if !cur.HasNext() {
			return nil, false
		}
		return ToAny(cur.Next()), true
	}}
}

func streamConcat(a, b Ext_java_util_stream_Stream) Ext_java_util_stream_Stream {
	CheckNotNull(a)
	CheckNotNull(b)
	s1 := a.(*jStream)
	s2 := b.(*jStream)
	done1 := false
	return &jStream{next: func() (any, bool) {
		if !done1 {
			if v, ok := s1.next(); ok {
				return v, true
			}
			done1 = true
		}
		return s2.next()
	}}
}

func streamMapWithIndex(s Ext_java_util_stream_Stream, fn Ext_com_google_common_collect_Streams_FunctionWithIndex) Ext_java_util_stream_Stream {
	CheckNotNull(s)
	CheckNotNull(fn)
	src := s.(*jStream)
	var idx int64
	return &jStream{next: func() (any, bool) {
		v, ok := src.next()
		if !ok {
			return nil, false
		}
		i := idx
		idx++
		return ToAny(fn.Apply(v, i)), true
	}}
}

func (*jStream) Implements_j2gort_Ext_java_lang_AutoCloseable()     {}
func (*jStream) Implements_j2gort_Ext_java_util_stream_BaseStream() {}
func (*jStream) Implements_j2gort_Ext_java_util_stream_Stream()     {}
func (*jStream) Close_0()                                           {}

func (s *jStream) drain() []any {
	var out []any
	for {
		v, ok := s.next()
		if !ok {
			break
		}
		out = append(out, v)
	}
	return out
}

func (s *jStream) Filter(pred Ext_java_util_function_Predicate) Ext_java_util_stream_Stream {
	CheckNotNull(pred)
	return &jStream{next: func() (any, bool) {
		for {
			v, ok := s.next()
			if !ok {
				return nil, false
			}
			if pred.Test(v) {
				return v, true
			}
		}
	}}
}

func (s *jStream) Map(fn Ext_java_util_function_Function) Ext_java_util_stream_Stream {
	CheckNotNull(fn)
	return &jStream{next: func() (any, bool) {
		v, ok := s.next()
		if !ok {
			return nil, false
		}
		return ToAny(fn.Apply(v)), true
	}}
}

func (s *jStream) FlatMap(fn Ext_java_util_function_Function) Ext_java_util_stream_Stream {
	CheckNotNull(fn)
	var cur *jStream
	return &jStream{next: func() (any, bool) {
		for {
			if cur != nil {
				if v, ok := cur.next(); ok {
					return v, true
				}
				cur = nil
			}
			outer, ok := s.next()
			if !ok {
				return nil, false
			}
			sub := fn.Apply(outer)
			if sub != nil && !javalang.IsNil(sub) {
				cur = sub.(*jStream)
			}
		}
	}}
}

func (s *jStream) Distinct() Ext_java_util_stream_Stream {
	seen := NewJMapCore(false)
	return &jStream{next: func() (any, bool) {
		for {
			v, ok := s.next()
			if !ok {
				return nil, false
			}
			if !seen.ContainsKey(v) {
				seen.Put(v, true)
				return v, true
			}
		}
	}}
}

func (s *jStream) Sorted_0() Ext_java_util_stream_Stream {
	return s.Sorted_1(nil)
}

func (s *jStream) Sorted_1(cmp Ext_java_util_Comparator) Ext_java_util_stream_Stream {
	fn := comparatorFuncOrNatural(cmp)
	var buf []any
	sorted := false
	idx := 0
	return &jStream{next: func() (any, bool) {
		if !sorted {
			buf = s.drain()
			sort.SliceStable(buf, func(i, j int) bool { return fn(buf[i], buf[j]) < 0 })
			sorted = true
		}
		if idx >= len(buf) {
			return nil, false
		}
		v := buf[idx]
		idx++
		return v, true
	}}
}

func (s *jStream) Peek(action Ext_java_util_function_Consumer) Ext_java_util_stream_Stream {
	CheckNotNull(action)
	return &jStream{next: func() (any, bool) {
		v, ok := s.next()
		if ok {
			action.Accept(v)
		}
		return v, ok
	}}
}

func (s *jStream) Limit(maxSize int64) Ext_java_util_stream_Stream {
	CheckArgument(maxSize >= 0)
	var count int64
	return &jStream{next: func() (any, bool) {
		if count >= maxSize {
			return nil, false
		}
		v, ok := s.next()
		if !ok {
			return nil, false
		}
		count++
		return v, true
	}}
}

func (s *jStream) Skip(n int64) Ext_java_util_stream_Stream {
	CheckArgument(n >= 0)
	skipped := false
	return &jStream{next: func() (any, bool) {
		if !skipped {
			for i := int64(0); i < n; i++ {
				if _, ok := s.next(); !ok {
					skipped = true
					return nil, false
				}
			}
			skipped = true
		}
		return s.next()
	}}
}

func (s *jStream) ForEach(action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	for {
		v, ok := s.next()
		if !ok {
			break
		}
		action.Accept(v)
	}
}

func (s *jStream) ForEachOrdered(action Ext_java_util_function_Consumer) {
	s.ForEach(action)
}

func (s *jStream) ToArray_1(gen Ext_java_util_function_IntFunction) []any {
	CheckNotNull(gen)
	items := s.drain()
	arr := gen.Apply(int32(len(items)))
	if out, ok := arr.([]any); ok {
		copy(out, items)
		return out
	}
	return items
}

func (s *jStream) Reduce_2(identity any, op Ext_java_util_function_BinaryOperator) any {
	CheckNotNull(op)
	acc := ToAny(identity)
	for {
		v, ok := s.next()
		if !ok {
			break
		}
		acc = ToAny(op.Apply(acc, v))
	}
	return acc
}

func (s *jStream) Collect_1(collector Ext_java_util_stream_Collector) any {
	CheckNotNull(collector)
	return collector.(*jCollector).fn(s.drain())
}

func (s *jStream) AnyMatch(pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	for {
		v, ok := s.next()
		if !ok {
			return false
		}
		if pred.Test(v) {
			return true
		}
	}
}

func (s *jStream) AllMatch(pred Ext_java_util_function_Predicate) bool {
	CheckNotNull(pred)
	for {
		v, ok := s.next()
		if !ok {
			return true
		}
		if !pred.Test(v) {
			return false
		}
	}
}

func (s *jStream) FindFirst() *Ext_java_util_Optional {
	v, ok := s.next()
	opt := &Ext_java_util_Optional{}
	if ok {
		opt.Core_ = CheckNotNull(v)
	}
	return opt
}

type jStreamBuilder struct {
	items []any
}

func (*jStreamBuilder) Implements_j2gort_Ext_java_util_function_Consumer()     {}
func (*jStreamBuilder) Implements_j2gort_Ext_java_util_stream_Stream_Builder() {}
func (b *jStreamBuilder) Accept(v any)                                         { b.items = append(b.items, ToAny(v)) }
func (b *jStreamBuilder) Add(v any) Ext_java_util_stream_Stream_Builder {
	b.Accept(v)
	return b
}
func (b *jStreamBuilder) Build() Ext_java_util_stream_Stream {
	return streamFromSlice(b.items)
}

type jIntStream struct {
	start int32
	end   int32
}

func (*jIntStream) Implements_j2gort_Ext_java_lang_AutoCloseable()     {}
func (*jIntStream) Implements_j2gort_Ext_java_util_stream_BaseStream() {}
func (*jIntStream) Implements_j2gort_Ext_java_util_stream_IntStream()  {}
func (*jIntStream) Close_0()                                           {}
func (is *jIntStream) MapToObj(fn Ext_java_util_function_IntFunction) Ext_java_util_stream_Stream {
	CheckNotNull(fn)
	cur := is.start
	return &jStream{next: func() (any, bool) {
		if cur >= is.end {
			return nil, false
		}
		i := cur
		cur++
		return ToAny(fn.Apply(i)), true
	}}
}

type jCollector struct {
	fn func(items []any) any
}

func (*jCollector) Implements_j2gort_Ext_java_util_stream_Collector() {}

func collectorToList() Ext_java_util_stream_Collector {
	return &jCollector{fn: func(items []any) any { return wrapArrayListOwn(items) }}
}

func collectorToImmutableList() Ext_java_util_stream_Collector {
	return &jCollector{fn: func(items []any) any {
		for _, x := range items {
			CheckNotNull(x)
		}
		return wrapImmutableListOwn(items)
	}}
}

func collectorToImmutableSet() Ext_java_util_stream_Collector {
	return &jCollector{fn: func(items []any) any {
		for _, x := range items {
			CheckNotNull(x)
		}
		return wrapImmutableSet(items)
	}}
}

func collectorToImmutableSortedSet(cmp Ext_java_util_Comparator) Ext_java_util_stream_Collector {
	return &jCollector{fn: func(items []any) any {
		if len(items) == 0 && isNaturalOrderComparator(cmp) {
			return emptyImmutableSortedSet_
		}
		mc := newSortedCoreFromCmp(cmp)
		for _, x := range items {
			mc.Put(CheckNotNull(x), true)
		}
		return wrapImmutableSortedSet(mc)
	}}
}

func collectorToImmutableMultiset() Ext_java_util_stream_Collector {
	return &jCollector{fn: func(items []any) any {
		c := newMultisetCore(false, nil)
		for _, x := range items {
			multisetOp_add_2(c, CheckNotNull(x), 1)
		}
		c.immutable = true
		return wrapImmutableMultiset(c)
	}}
}

func collectorToImmutableMap(kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
	CheckNotNull(kFn)
	CheckNotNull(vFn)
	return &jCollector{fn: func(items []any) any {
		if len(items) == 0 {
			return emptyImmutableMap_
		}
		obj := newImmutableMapShell()
		mc := &obj.inlineCore
		mc.InitExpectedSize(int32(len(items)))
		for _, x := range items {
			k := CheckNotNull(ToAny(kFn.Apply(x)))
			v := CheckNotNull(ToAny(vFn.Apply(x)))
			if mc.ContainsKey(k) {
				panic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
			}
			mc.Put(k, v)
		}
		mc.sealImmutable()
		return obj
	}}
}

func collectorToImmutableSortedMap(cmp Ext_java_util_Comparator, kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
	CheckNotNull(cmp)
	CheckNotNull(kFn)
	CheckNotNull(vFn)
	return &jCollector{fn: func(items []any) any {
		if len(items) == 0 && isNaturalOrderComparator(cmp) {
			return emptyImmutableSortedMap_
		}
		obj := &Ext_com_google_common_collect_ImmutableSortedMap{}
		obj.BindExtSelf_(obj)
		mc := initSortedCoreCmp(&obj.inlineCore, cmp)
		for _, x := range items {
			k := CheckNotNull(ToAny(kFn.Apply(x)))
			v := CheckNotNull(ToAny(vFn.Apply(x)))
			if mc.ContainsKey(k) {
				panic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
			}
			mc.Put(k, v)
		}
		mc.sealImmutable()
		return obj
	}}
}

func collectorToCollection(sup Ext_java_util_function_Supplier) Ext_java_util_stream_Collector {
	CheckNotNull(sup)
	return &jCollector{fn: func(items []any) any {
		col := Cast[Ext_java_util_Collection](sup.Get())
		for _, x := range items {
			col.Add_1(x)
		}
		return col
	}}
}

func collectorJoining(delim any) Ext_java_util_stream_Collector {
	d := StrVal(CheckNotNull(delim))
	return &jCollector{fn: func(items []any) any {
		var sb strings.Builder
		for i, x := range items {
			if i > 0 {
				sb.WriteString(d)
			}
			sb.WriteString(StrVal(x))
		}
		return sb.String()
	}}
}

func collectorGroupingBy(classifier Ext_java_util_function_Function) Ext_java_util_stream_Collector {
	CheckNotNull(classifier)
	return &jCollector{fn: func(items []any) any {
		mc := NewJMapCore(false)
		for _, x := range items {
			k := CheckNotNull(ToAny(classifier.Apply(x)))
			lst := mc.Get(k)
			if lst == nil {
				lst = wrapArrayList(nil)
				mc.Put(k, lst)
			}
			lst.(Ext_java_util_List).Add_1(x)
		}
		return wrapLinkedHashMap(mc)
	}}
}

func collectorToMap(kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
	CheckNotNull(kFn)
	CheckNotNull(vFn)
	return &jCollector{fn: func(items []any) any {
		mc := NewJMapCore(false)
		for _, x := range items {
			k := ToAny(kFn.Apply(x))
			v := CheckNotNull(ToAny(vFn.Apply(x)))
			if mc.ContainsKey(k) {
				panic(javalang.NewIllegalStateException("Duplicate key " + StrVal(k)))
			}
			mc.Put(k, v)
		}
		return wrapLinkedHashMap(mc)
	}}
}

// --- Cache & Executor / Future / AtomicIntegerArray ---

type jCacheBuilderCore struct {
	maxSize  int64
	listener Ext_com_google_common_cache_RemovalListener
}

type jFuncCacheLoader struct {
	Ext_com_google_common_cache_CacheLoaderBase
	fn Ext_com_google_common_base_Function
}

func (l *jFuncCacheLoader) Load(k any) any { return l.fn.Apply(k) }

type jLoadingCache struct {
	mu       sync.Mutex
	mc       *JMapCore
	maxSize  int64
	listener Ext_com_google_common_cache_RemovalListener
	loader   Ext_com_google_common_cache_CacheLoader
}

func newLoadingCache(cfg *jCacheBuilderCore, loader Ext_com_google_common_cache_CacheLoader) Ext_com_google_common_cache_LoadingCache {
	return &jLoadingCache{
		mc:       NewJMapCore(false),
		maxSize:  cfg.maxSize,
		listener: cfg.listener,
		loader:   loader,
	}
}

func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_cache_Cache()        {}
func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_cache_LoadingCache() {}
func (*jLoadingCache) Implements_j2gort_Ext_java_util_function_Function()          {}
func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_base_Function()      {}
func (c *jLoadingCache) Equals(o any) bool                                         { return RefEquals(c, o) }
func (c *jLoadingCache) Apply(k any) any                                           { return c.GetUnchecked(k) }

func (c *jLoadingCache) GetUnchecked(k any) any {
	k = CheckNotNull(ToAny(k))
	c.mu.Lock()
	if c.mc.ContainsKey(k) {
		v := c.mc.Get(k)
		c.mu.Unlock()
		return v
	}
	c.mu.Unlock()
	if c.maxSize == 0 {
		v := CheckNotNull(ToAny(c.loader.Load(k)))
		if c.listener != nil {
			n := &Ext_com_google_common_cache_RemovalNotification{}
			n.BindExtSelf_(n)
			n.Core_ = &jSimpleEntry{k: k, v: v, imm: true}
			c.listener.OnRemoval(n)
		}
		return v
	}
	v := CheckNotNull(ToAny(c.loader.Load(k)))
	var evicted *jSimpleEntry
	c.mu.Lock()
	if c.maxSize > 0 && int64(c.mc.Size()) >= c.maxSize && !c.mc.ContainsKey(k) {
		ek := c.mc.FirstKey()
		ev := c.mc.Remove(ek)
		if c.listener != nil {
			evicted = &jSimpleEntry{k: ek, v: ev, imm: true}
		}
	}
	c.mc.Put(k, v)
	c.mu.Unlock()
	if evicted != nil {
		n := &Ext_com_google_common_cache_RemovalNotification{}
		n.BindExtSelf_(n)
		n.Core_ = evicted
		c.listener.OnRemoval(n)
	}
	return v
}

type jDirectExecutor struct{}

func (*jDirectExecutor) Implements_j2gort_Ext_java_util_concurrent_Executor()        {}
func (*jDirectExecutor) Implements_j2gort_Ext_java_lang_AutoCloseable()              {}
func (*jDirectExecutor) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
func (*jDirectExecutor) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {
}
func (*jDirectExecutor) Shutdown() {}
func (*jDirectExecutor) Close_0()  {}
func (*jDirectExecutor) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
	return execSubmitCallable(task)
}

type jThreadPoolExecutor struct {
	tasks     chan func()
	wg        sync.WaitGroup
	closeOnce sync.Once
}

func newThreadPoolExecutor(threads int32) Ext_java_util_concurrent_ExecutorService {
	if threads <= 1 {
		return &jDirectExecutor{}
	}
	n := int(threads)
	buf := n * 16
	if buf < 64 {
		buf = 64
	}
	p := &jThreadPoolExecutor{tasks: make(chan func(), buf)}
	p.wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer p.wg.Done()
			for fn := range p.tasks {
				fn()
			}
		}()
	}
	return p
}

func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_util_concurrent_Executor()        {}
func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_lang_AutoCloseable()              {}
func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
func (*jThreadPoolExecutor) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {
}
func (p *jThreadPoolExecutor) Shutdown() {
	p.closeOnce.Do(func() {
		close(p.tasks)
	})
}
func (p *jThreadPoolExecutor) Close_0() {
	p.Shutdown()
	p.wg.Wait()
}
func (p *jThreadPoolExecutor) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
	CheckNotNull(task)
	f := &jAsyncFuture{done: make(chan struct{})}
	nodestore.SuspendSingleOwnerAlloc()
	p.tasks <- func() {
		defer func() {
			if r := recover(); r != nil {
				f.err = javalang.RecoverException(r)
			}
			nodestore.ResumeSingleOwnerAlloc()
			close(f.done)
		}()
		f.val = ToAny(task.Call())
	}
	return f
}

type jAsyncFuture struct {
	done chan struct{}
	val  any
	err  javalang.Throwable
}

func (*jAsyncFuture) Implements_j2gort_Ext_java_util_concurrent_Future()                        {}
func (*jAsyncFuture) Implements_j2gort_Ext_com_google_common_util_concurrent_ListenableFuture() {}
func (f *jAsyncFuture) Get_0() any {
	<-f.done
	if f.err != nil {
		ex := &Ext_java_util_concurrent_ExecutionExceptionBase{}
		ex.BindExtSelf_(ex)
		InitUserThrowable(ex, "java.util.concurrent.ExecutionException", f.err)
		panic(ex)
	}
	return f.val
}
func (f *jAsyncFuture) Get_2(_ int64, _ *Ext_java_util_concurrent_TimeUnit) any {
	return f.Get_0()
}

type jListeningExecWrapper struct {
	exec Ext_java_util_concurrent_ExecutorService
}

func (*jListeningExecWrapper) Implements_j2gort_Ext_java_util_concurrent_Executor()        {}
func (*jListeningExecWrapper) Implements_j2gort_Ext_java_lang_AutoCloseable()              {}
func (*jListeningExecWrapper) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
func (*jListeningExecWrapper) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {
}
func (w *jListeningExecWrapper) Shutdown() { w.exec.Shutdown() }
func (w *jListeningExecWrapper) Close_0()  { w.exec.Close_0() }
func (w *jListeningExecWrapper) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
	return w.exec.Submit_1_Ext_java_util_concurrent_Callable(task)
}

type jCompletedFuture struct {
	val any
	err javalang.Throwable
}

func (*jCompletedFuture) Implements_j2gort_Ext_java_util_concurrent_Future()                        {}
func (*jCompletedFuture) Implements_j2gort_Ext_com_google_common_util_concurrent_ListenableFuture() {}
func (f *jCompletedFuture) Get_0() any {
	if f.err != nil {
		ex := &Ext_java_util_concurrent_ExecutionExceptionBase{}
		ex.BindExtSelf_(ex)
		InitUserThrowable(ex, "java.util.concurrent.ExecutionException", f.err)
		panic(ex)
	}
	return f.val
}
func (f *jCompletedFuture) Get_2(_ int64, _ *Ext_java_util_concurrent_TimeUnit) any {
	return f.Get_0()
}

func execSubmitCallable(task Ext_java_util_concurrent_Callable) Ext_com_google_common_util_concurrent_ListenableFuture {
	CheckNotNull(task)
	var res any
	var th javalang.Throwable
	func() {
		if nodestore.TryEnterSingleOwnerAlloc() {
			defer nodestore.ExitSingleOwnerAlloc()
		}
		defer func() {
			if r := recover(); r != nil {
				th = javalang.RecoverException(r)
			}
		}()
		res = ToAny(task.Call())
	}()
	return &jCompletedFuture{val: res, err: th}
}

func futuresAllAsList(futures Ext_java_lang_Iterable) Ext_com_google_common_util_concurrent_ListenableFuture {
	CheckNotNull(futures)
	var results []any
	for it := Iterate(futures); it.HasNext(); {
		f := Cast[Ext_java_util_concurrent_Future](it.Next())
		var v any
		var th javalang.Throwable
		func() {
			defer func() {
				if r := recover(); r != nil {
					th = javalang.RecoverException(r)
				}
			}()
			v = f.Get_0()
		}()
		if th != nil {
			if ee, ok := th.(Ext_java_util_concurrent_ExecutionException); ok {
				return &jCompletedFuture{err: ee.GetCause()}
			}
			return &jCompletedFuture{err: th}
		}
		results = append(results, v)
	}
	return &jCompletedFuture{val: wrapArrayList(results)}
}

type jAtomicIntArray struct {
	mu   sync.Mutex
	vals []int32
}

func (a *jAtomicIntArray) get(idx int32) int32 {
	a.mu.Lock()
	defer a.mu.Unlock()
	CheckElementIndex(idx, int32(len(a.vals)))
	return a.vals[idx]
}

func (a *jAtomicIntArray) updateAndGet(idx int32, op Ext_java_util_function_IntUnaryOperator) int32 {
	CheckNotNull(op)
	a.mu.Lock()
	defer a.mu.Unlock()
	CheckElementIndex(idx, int32(len(a.vals)))
	nv := op.ApplyAsInt(a.vals[idx])
	a.vals[idx] = nv
	return nv
}

func (x *JListCore) IsJavaNil_() bool             { return x == nil }
func (x *JDequeCore) IsJavaNil_() bool            { return x == nil }
func (x *JSetCore) IsJavaNil_() bool              { return x == nil }
func (x *JMapCore) IsJavaNil_() bool              { return x == nil }
func (x *jSetIter) IsJavaNil_() bool              { return x == nil }
func (x *jComparatorFunc) IsJavaNil_() bool       { return x == nil }
func (x *jPredFunc) IsJavaNil_() bool             { return x == nil }
func (x *jFuncIdentity) IsJavaNil_() bool         { return x == nil }
func (x *jFuncWrap) IsJavaNil_() bool             { return x == nil }
func (x *jConsumerWrap) IsJavaNil_() bool         { return x == nil }
func (x *jFullListIter) IsJavaNil_() bool         { return x == nil }
func (x *jDescListIter) IsJavaNil_() bool         { return x == nil }
func (x *jDequeIter) IsJavaNil_() bool            { return x == nil }
func (x *jDequeDescIter) IsJavaNil_() bool        { return x == nil }
func (x *jMapKeyIter) IsJavaNil_() bool           { return x == nil }
func (x *jSubListView) IsJavaNil_() bool          { return x == nil }
func (x *jBackedArrayList) IsJavaNil_() bool      { return x == nil }
func (x *jReversedListView) IsJavaNil_() bool     { return x == nil }
func (x *jTransformedListView) IsJavaNil_() bool  { return x == nil }
func (x *jSimpleEntry) IsJavaNil_() bool          { return x == nil }
func (x *jLiveMapKeySet) IsJavaNil_() bool        { return x == nil }
func (x *jLiveMapValues) IsJavaNil_() bool        { return x == nil }
func (x *jMapValIter) IsJavaNil_() bool           { return x == nil }
func (x *jLiveMapEntrySet) IsJavaNil_() bool      { return x == nil }
func (x *jMapEntryIter) IsJavaNil_() bool         { return x == nil }
func (x *jImmMapBuilderCore) IsJavaNil_() bool    { return x == nil }
func (x *jAbstractListIter) IsJavaNil_() bool     { return x == nil }
func (x *jSetViewCore) IsJavaNil_() bool          { return x == nil }
func (x *jSyncMapCore) IsJavaNil_() bool          { return x == nil }
func (x *jBiMapCore) IsJavaNil_() bool            { return x == nil }
func (x *jMultimapCore) IsJavaNil_() bool         { return x == nil }
func (x *jMultimapLiveList) IsJavaNil_() bool     { return x == nil }
func (x *jMultimapLiveSet) IsJavaNil_() bool      { return x == nil }
func (x *jMultimapSetIter) IsJavaNil_() bool      { return x == nil }
func (x *jMultimapLiveKeySet) IsJavaNil_() bool   { return x == nil }
func (x *jMultimapKeyIter) IsJavaNil_() bool      { return x == nil }
func (x *jMultimapAsMap) IsJavaNil_() bool        { return x == nil }
func (x *jMultimapAsMapEntries) IsJavaNil_() bool { return x == nil }
func (x *jMultimapAsMapIter) IsJavaNil_() bool    { return x == nil }
func (x *jMultimapBuilderCore) IsJavaNil_() bool  { return x == nil }
func (x *jMultisetCore) IsJavaNil_() bool         { return x == nil }
func (x *jMultisetEntry) IsJavaNil_() bool        { return x == nil }
func (x *jMultisetIter) IsJavaNil_() bool         { return x == nil }
func (x *jTableCore) IsJavaNil_() bool            { return x == nil }
func (x *jTableRowMap) IsJavaNil_() bool          { return x == nil }
func (x *jTableRowKeySet) IsJavaNil_() bool       { return x == nil }
func (x *jTableRowKeyIter) IsJavaNil_() bool      { return x == nil }
func (x *jTableRowMapOuter) IsJavaNil_() bool     { return x == nil }
func (x *jTableRowKeyView) IsJavaNil_() bool      { return x == nil }
func (x *jTableRowEntrySet) IsJavaNil_() bool     { return x == nil }
func (x *jTableRowSingleIter) IsJavaNil_() bool   { return x == nil }
func (x *jPriorityQueue) IsJavaNil_() bool        { return x == nil }
func (x *jFuncIterable) IsJavaNil_() bool         { return x == nil }
func (x *jConcatIter) IsJavaNil_() bool           { return x == nil }
func (x *jFilterIter) IsJavaNil_() bool           { return x == nil }
func (x *jSimpleIterAdapter) IsJavaNil_() bool    { return x == nil }
func (x *jTransformIter) IsJavaNil_() bool        { return x == nil }
func (x *jPeekingIter) IsJavaNil_() bool          { return x == nil }
func (x *jAbstractIterState) IsJavaNil_() bool    { return x == nil }
func (x *jComparisonChain) IsJavaNil_() bool      { return x == nil }
func (x *jRangeCore) IsJavaNil_() bool            { return x == nil }
func (x *jRangeSetCore) IsJavaNil_() bool         { return x == nil }
func (x *jStream) IsJavaNil_() bool               { return x == nil }
func (x *jStreamBuilder) IsJavaNil_() bool        { return x == nil }
func (x *jIntStream) IsJavaNil_() bool            { return x == nil }
func (x *jCollector) IsJavaNil_() bool            { return x == nil }
func (x *jCacheBuilderCore) IsJavaNil_() bool     { return x == nil }
func (x *jFuncCacheLoader) IsJavaNil_() bool      { return x == nil }
func (x *jLoadingCache) IsJavaNil_() bool         { return x == nil }
func (x *jEmptyMapIter) IsJavaNil_() bool         { return x == nil }
func (x *jLiveMapEntry) IsJavaNil_() bool         { return x == nil }
func (x *jDirectExecutor) IsJavaNil_() bool       { return x == nil }
func (x *jThreadPoolExecutor) IsJavaNil_() bool   { return x == nil }
func (x *jListeningExecWrapper) IsJavaNil_() bool { return x == nil }
func (x *jCompletedFuture) IsJavaNil_() bool      { return x == nil }
func (x *jAsyncFuture) IsJavaNil_() bool          { return x == nil }
func (x *jAtomicIntArray) IsJavaNil_() bool       { return x == nil }
func (x *jEnumSetCore) IsJavaNil_() bool          { return x == nil }
func (x *jEnumSetIter) IsJavaNil_() bool          { return x == nil }
func (x *jEnumMapCore) IsJavaNil_() bool          { return x == nil }
func (x *jLiveEnumMapKeySet) IsJavaNil_() bool    { return x == nil }
func (x *jEnumMapKeyIter) IsJavaNil_() bool       { return x == nil }
func (x *jLiveEnumMapValues) IsJavaNil_() bool    { return x == nil }
func (x *jEnumMapValIter) IsJavaNil_() bool       { return x == nil }
func (x *jLiveEnumMapEntrySet) IsJavaNil_() bool  { return x == nil }
func (x *jEnumMapEntry) IsJavaNil_() bool         { return x == nil }
func (x *jEnumMapEntryIter) IsJavaNil_() bool     { return x == nil }

// --- Post-Pass Native Slice, OrderedSet, and LocalStringBuilder Runtime Helpers ---

func boxJavaElem(v any) any {
	switch x := v.(type) {
	case int32:
		return ValueOfInt(x)
	case int64:
		return ValueOfLong(x)
	case bool:
		return ValueOfBoolean(x)
	case float64:
		return ValueOfDouble(x)
	case float32:
		return ValueOfFloat(x)
	case uint16:
		return ValueOfChar(x)
	case int8:
		return ValueOfByte(x)
	case int16:
		return ValueOfShort(x)
	case uint32:
		return boxNodeById(x)
	default:
		if rn, ok := v.(jRawNodeProvider); ok {
			if rn.RawNodeId_() == 0 {
				return nil
			}
			return v
		}
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
	case int64:
		return any(UnboxLong(CastPtr[BoxLong](v))).(T)
	case bool:
		return any(UnboxBoolean(CastPtr[BoxBoolean](v))).(T)
	case float64:
		return any(UnboxDouble(CastPtr[BoxDouble](v))).(T)
	case float32:
		return any(UnboxFloat(CastPtr[BoxFloat](v))).(T)
	case uint16:
		return any(UnboxChar(CastPtr[BoxChar](v))).(T)
	case int8:
		return any(UnboxByte(CastPtr[BoxByte](v))).(T)
	case int16:
		return any(UnboxShort(CastPtr[BoxShort](v))).(T)
	case uint32:
		return any(CastNode[uint32](v)).(T)
	default:
		return Cast[T](v)
	}
}

func MakeNonNilSlice[T any](cap int32) []T {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal Capacity: %d", cap)))
	}
	if cap == 0 {
		return []T{}
	}
	return make([]T, 0, cap)
}

func CheckNotNullElem[T any](v T) T {
	if any(v) == nil {
		ThrowNPE()
	}
	switch x := any(v).(type) {
	case string:
		if x == NullString {
			ThrowNPE()
		}
		return v
	case int32, int64, bool, float64, float32, uint16, int8, int16, uint8, int, uint64:
		return v
	default:
		switch nullValKindOf[T]() {
		case kindNeverNull:
			return v
		case kindZeroIsNull:
			var zero T
			if any(v) == any(zero) {
				ThrowNPE()
			}
			return v
		case kindNullStringType:
			if any(v).(string) == NullString {
				ThrowNPE()
			}
			return v
		default:
			if isJavaNullMapValSlow(v) {
				ThrowNPE()
			}
			return v
		}
	}
}

func ImmSliceCloneNonNull[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return []T{}
	}
	out := make([]T, len(s))
	switch nullValKindOf[T]() {
	case kindNeverNull:
		copy(out, s)
		return out
	case kindZeroIsNull:
		var zero T
		for i, v := range s {
			if any(v) == any(zero) {
				ThrowNPE()
			}
			out[i] = v
		}
		return out
	case kindNullStringType:
		for _, sv := range any(s).([]string) {
			if sv == NullString {
				ThrowNPE()
			}
		}
		copy(out, s)
		return out
	default:
		for i, v := range s {
			if isJavaNullMapValSlow(v) {
				ThrowNPE()
			}
			out[i] = v
		}
		return out
	}
}

type immTypedSliceView[T any] struct {
	s []T
}

var _ Ext_com_google_common_collect_ImmutableList = (*immTypedSliceView[int32])(nil)

func (l *immTypedSliceView[T]) IsJavaNil_() bool { return l == nil }
func (*immTypedSliceView[T]) JavaClassInfo_() (string, string) {
	return "com.google.common.collect.ImmutableList", "ImmutableList"
}
func (l *immTypedSliceView[T]) typedViewLen_() (int, bool) {
	if l == nil {
		ThrowNPE()
	}
	return len(l.s), true
}
func (l *immTypedSliceView[T]) typedViewAt_(i int) any {
	return boxJavaElem(l.s[i])
}
func (l *immTypedSliceView[T]) typedViewNodeAt_(i int) (uint32, bool) {
	return typedElemRawNode(l.s[i])
}
func (l *immTypedSliceView[T]) IterSlice_() []any {
	if l == nil {
		ThrowNPE()
	}
	if len(l.s) == 0 {
		return emptyAnySlice_
	}
	out := make([]any, len(l.s))
	for i, v := range l.s {
		out[i] = boxJavaElem(v)
	}
	return out
}
func (l *immTypedSliceView[T]) ListGetStr_(i int32) string {
	if l == nil {
		ThrowNPE()
	}
	CheckElementIndex(i, int32(len(l.s)))
	return CastString(boxJavaElem(l.s[i]))
}
func (l *immTypedSliceView[T]) SetContainsStr_(k string) bool {
	if l == nil {
		ThrowNPE()
	}
	if k == NullString {
		return false
	}
	return SliceContains(l.s, k)
}
func (l *immTypedSliceView[T]) SetAddStr_(string) bool {
	if l == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) SetRemoveStr_(string) bool {
	if l == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) SetContainsNode_(u uint32) bool {
	if l == nil {
		ThrowNPE()
	}
	if u == 0 {
		return false
	}
	for _, v := range l.s {
		if singletonElemMatchesNode(boxJavaElem(v), u) {
			return true
		}
	}
	return false
}
func (l *immTypedSliceView[T]) SetAddNode_(uint32) bool {
	if l == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) SetRemoveNode_(uint32) bool {
	if l == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Equals(arg0 any) bool {
	if l == nil {
		ThrowNPE()
	}
	if any(l) == arg0 {
		return true
	}
	return ImmSliceEquals(l.s, arg0)
}
func (l *immTypedSliceView[T]) HashCode() int32 {
	if l == nil {
		ThrowNPE()
	}
	return ImmSliceHashCode(l.s)
}
func (l *immTypedSliceView[T]) String() string {
	if l == nil {
		return "null"
	}
	if len(l.s) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range l.s {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(boxJavaElem(v)))
	}
	sb.WriteByte(']')
	return sb.String()
}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_com_google_common_collect_ImmutableList()       {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_com_google_common_collect_ImmutableCollection() {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_List()                                {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_RandomAccess()                        {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_AbstractCollection()                  {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_SequencedCollection()                 {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_Collection()                          {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_lang_Iterable()                            {}
func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_io_Serializable()                          {}
func (*immTypedSliceView[T]) GetExt_com_google_common_collect_ImmutableListBase() *Ext_com_google_common_collect_ImmutableListBase {
	return nil
}
func (*immTypedSliceView[T]) GetExt_com_google_common_collect_ImmutableCollectionBase() *Ext_com_google_common_collect_ImmutableCollectionBase {
	return nil
}
func (*immTypedSliceView[T]) GetExt_java_util_AbstractCollectionBase() *Ext_java_util_AbstractCollectionBase {
	return nil
}
func (l *immTypedSliceView[T]) Iterator() any {
	if l == nil {
		ThrowNPE()
	}
	return SliceIterator(l.s)
}
func (l *immTypedSliceView[T]) ListIterator_1(arg0 int32) any {
	if l == nil {
		ThrowNPE()
	}
	return newImmListIter(l.IterSlice_(), arg0)
}
func (l *immTypedSliceView[T]) ForEach(arg0 Ext_java_util_function_Consumer) {
	if l == nil {
		ThrowNPE()
	}
	CheckNotNull(arg0)
	for _, v := range l.s {
		arg0.Accept(boxJavaElem(v))
	}
}
func (l *immTypedSliceView[T]) IndexOf_1(arg0 any) int32     { return SliceIndexOf(l.s, arg0) }
func (l *immTypedSliceView[T]) LastIndexOf_1(arg0 any) int32 { return SliceLastIndexOf(l.s, arg0) }
func (l *immTypedSliceView[T]) Contains(arg0 any) bool       { return SliceContains(l.s, arg0) }
func (l *immTypedSliceView[T]) SubList(arg0 int32, arg1 int32) any {
	if l == nil {
		ThrowNPE()
	}
	n := int32(len(l.s))
	CheckPositionIndexes(arg0, arg1, n)
	if arg0 == 0 && arg1 == n {
		return l
	}
	return ImmListFromTypedSlice(l.s[arg0:arg1:arg1])
}
func (l *immTypedSliceView[T]) Set(int32, any) any {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Add_2(int32, any) {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Remove_1_int32(int32) any {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Sort(Ext_java_util_Comparator) {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) AsList() Ext_com_google_common_collect_ImmutableList { return l }
func (l *immTypedSliceView[T]) Reverse() Ext_com_google_common_collect_ImmutableList {
	if len(l.s) <= 1 {
		return l
	}
	return ImmListFromTypedSlice(ImmSliceReverse(l.s))
}
func (l *immTypedSliceView[T]) ToArray_1_slice_any(arg0 []any) []any {
	return sliceToArray1(l.IterSlice_(), arg0)
}
func (l *immTypedSliceView[T]) Add_1(any) bool {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Remove_1_any(any) bool {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) AddAll_1(Ext_java_util_Collection) bool {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) RemoveAll(Ext_java_util_Collection) bool {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Clear() {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Size() int32   { return int32(len(l.s)) }
func (l *immTypedSliceView[T]) IsEmpty() bool { return len(l.s) == 0 }
func (l *immTypedSliceView[T]) ContainsAll(arg0 Ext_java_util_Collection) bool {
	return colContainsAll(func(x any) bool { return SliceContains(l.s, x) }, arg0)
}
func (l *immTypedSliceView[T]) Get(arg0 int32) any {
	CheckElementIndex(arg0, int32(len(l.s)))
	return boxJavaElem(l.s[arg0])
}
func (l *immTypedSliceView[T]) AddFirst(any) {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) AddLast(any) {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) GetFirst() any {
	if len(l.s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return boxJavaElem(l.s[0])
}
func (l *immTypedSliceView[T]) GetLast() any {
	if len(l.s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return boxJavaElem(l.s[len(l.s)-1])
}
func (l *immTypedSliceView[T]) RemoveFirst() any {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) RemoveLast() any {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (l *immTypedSliceView[T]) Stream() Ext_java_util_stream_Stream {
	return streamFromSlice(l.IterSlice_())
}

func emptyNonNilSlice[T any]() []T {
	return make([]T, 0)
}

func ImmListFromTypedSlice[T any](s []T) Ext_com_google_common_collect_ImmutableList {
	if s == nil {
		return nil
	}
	if len(s) == 0 {
		return emptyImmutableList_
	}
	if ss, ok := any(s).([]string); ok {
		return &jImmStrList{strs: ss[:len(ss):len(ss)]}
	}
	return &immTypedSliceView[T]{s: s[:len(s):len(s)]}
}

func ImmListFromCompSlice[T comparable](s []T) Ext_com_google_common_collect_ImmutableList {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return emptyImmutableList_
	}
	var zero T
	cs := make([]T, len(s))
	for i, v := range s {
		if v == zero {
			ThrowNPE()
		}
		cs[i] = v
	}
	return &immTypedSliceView[T]{s: cs[:len(cs):len(cs)]}
}

func ImmListFromSlice[T any](s []T) Ext_com_google_common_collect_ImmutableList {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return emptyImmutableList_
	}
	if ss, ok := any(s).([]string); ok {
		return wrapImmStrSliceCopy(ss)
	}
	cs := ImmSliceCloneNonNull(s)
	return &immTypedSliceView[T]{s: cs[:len(cs):len(cs)]}
}

func unboxNonNullImmElem[T any](v any) T {
	if v == nil {
		ThrowNPE()
	}
	if res, ok := v.(T); ok {
		return res
	}
	return unboxJavaElem[T](CheckNotNull(v))
}

func ImmListCopyToSlice[T any](src any) []T {
	CheckNotNull(src)
	if tv, ok := src.(*immTypedSliceView[T]); ok {
		if len(tv.s) == 0 {
			return emptyNonNilSlice[T]()
		}
		return tv.s[:len(tv.s):len(tv.s)]
	}
	if al, ok := src.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
		if len(al.s) == 0 {
			return emptyNonNilSlice[T]()
		}
		return ImmSliceCloneNonNull(al.s)
	}
	if src == emptyImmutableList_ {
		return emptyNonNilSlice[T]()
	}
	if se, ok := src.(jSingleNonNullElem); ok {
		return []T{unboxNonNullImmElem[T](se.singleNonNullElem_())}
	}
	switch l := src.(type) {
	case *jImmList1:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxNonNullImmElem[T](l.elem[0])}
	case *jImmList2:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxNonNullImmElem[T](l.elem[0]), unboxNonNullImmElem[T](l.elem[1])}
	case *jSingletonList:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxJavaElem[T](CheckNotNull(l.elem[0]))}
	}
	if ssi, ok := src.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(ss) == 0 {
				return emptyNonNilSlice[T]()
			}
			if _, isStr := any((*T)(nil)).(*string); isStr {
				if l, ok3 := src.(*jImmStrList); ok3 && l != nil {
					return any(l.strs[:len(l.strs):len(l.strs)]).([]T)
				}
				out := make([]string, len(ss))
				for i, str := range ss {
					if str == NullString {
						ThrowNPE()
					}
					out[i] = str
				}
				return any(out).([]T)
			}
			out := make([]T, len(ss))
			for i, str := range ss {
				if str == NullString {
					ThrowNPE()
				}
				out[i] = unboxJavaElem[T](str)
			}
			return out
		}
	}
	if tv, ok := src.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return emptyNonNilSlice[T]()
			}
			out := make([]T, n)
			for i := 0; i < n; i++ {
				out[i] = unboxJavaElem[T](CheckNotNull(tv.typedViewAt_(i)))
			}
			return out
		}
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return emptyNonNilSlice[T]()
		}
		out := make([]T, len(items))
		for i, raw := range items {
			out[i] = unboxJavaElem[T](CheckNotNull(raw))
		}
		return out
	}
	if c, ok := src.(Ext_java_util_Collection); ok {
		sz := int(c.Size())
		if sz == 0 {
			return emptyNonNilSlice[T]()
		}
		out := make([]T, 0, sz)
		for cur := NewIterCursor(src); cur.HasNext(); {
			out = append(out, unboxJavaElem[T](CheckNotNull(cur.Next())))
		}
		return out
	}
	out := emptyNonNilSlice[T]()
	for cur := NewIterCursor(src); cur.HasNext(); {
		out = append(out, unboxJavaElem[T](CheckNotNull(cur.Next())))
	}
	return out
}

func ImmListCopyArrayToSlice[T any](arr []any) []T {
	CheckNotNullSlice(arr)
	if len(arr) == 0 {
		return emptyNonNilSlice[T]()
	}
	out := make([]T, len(arr))
	for i, raw := range arr {
		out[i] = unboxJavaElem[T](CheckNotNull(raw))
	}
	return out
}

func ImmListToSliceNullable[T any](imm Ext_com_google_common_collect_ImmutableList) []T {
	if imm == nil {
		return nil
	}
	if tv, ok := any(imm).(*immTypedSliceView[T]); ok {
		if tv == nil {
			return nil
		}
		if len(tv.s) == 0 {
			return emptyNonNilSlice[T]()
		}
		return tv.s[:len(tv.s):len(tv.s)]
	}
	if imm == emptyImmutableList_ {
		return emptyNonNilSlice[T]()
	}
	if isJavaNilAny(imm) {
		return nil
	}
	if imm.IsEmpty() {
		return emptyNonNilSlice[T]()
	}
	return ImmListCopyToSlice[T](imm)
}

func ImmListOfVarargsToSlice[T any](prefix []T, rest []any) []T {
	CheckNotNullSlice(rest)
	n := len(prefix) + len(rest)
	if n == 0 {
		return []T{}
	}
	out := make([]T, n)
	if len(prefix) > 0 {
		switch nullValKindOf[T]() {
		case kindNeverNull:
			copy(out, prefix)
		case kindZeroIsNull:
			var zero T
			for i, v := range prefix {
				if any(v) == any(zero) {
					ThrowNPE()
				}
				out[i] = v
			}
		case kindNullStringType:
			for _, sv := range any(prefix).([]string) {
				if sv == NullString {
					ThrowNPE()
				}
			}
			copy(out, prefix)
		default:
			for i, v := range prefix {
				if isJavaNullMapValSlow(v) {
					ThrowNPE()
				}
				out[i] = v
			}
		}
	}
	for i, raw := range rest {
		out[len(prefix)+i] = unboxJavaElem[T](CheckNotNull(raw))
	}
	return out
}

func SliceBuilderAdd[T any](b *[]T, v T) *[]T {
	if b == nil || *b == nil {
		ThrowNPE()
	}
	if isJavaNullMapVal(v) {
		ThrowNPE()
	}
	*b = append(*b, v)
	return b
}

func AppendImmCollection[T any](dst []T, src any) []T {
	if dst == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	if tv, ok := src.(*immTypedSliceView[T]); ok {
		return append(dst, tv.s...)
	}
	if al, ok := src.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
		if len(al.s) == 0 {
			return dst
		}
		return append(dst, ImmSliceCloneNonNull(al.s)...)
	}
	if ssi, ok := src.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if dstStr, ok3 := any(&dst).(*[]string); ok3 {
				for _, str := range ss {
					if str == NullString {
						ThrowNPE()
					}
				}
				*dstStr = append(*dstStr, ss...)
				return dst
			}
			for _, str := range ss {
				if str == NullString {
					ThrowNPE()
				}
				dst = append(dst, unboxJavaElem[T](str))
			}
			return dst
		}
	}
	if si, ok := src.(SliceIterable); ok {
		for _, raw := range si.IterSlice_() {
			dst = append(dst, unboxJavaElem[T](CheckNotNull(raw)))
		}
		return dst
	}
	for cur := NewIterCursor(src); cur.HasNext(); {
		dst = append(dst, unboxJavaElem[T](CheckNotNull(cur.Next())))
	}
	return dst
}

func SliceBuilderAddAll[T any](b *[]T, src any) *[]T {
	if b == nil {
		ThrowNPE()
	}
	*b = AppendImmCollection(*b, src)
	return b
}

func SliceBuilderAddAllSlice[T any](b *[]T, src []T) *[]T {
	if b == nil || *b == nil || src == nil {
		ThrowNPE()
	}
	*b = append(*b, src...)
	return b
}

func ImmListCopyNodeToSlice[N ~uint32](src any) []N {
	CheckNotNull(src)
	if tv, ok := src.(*immTypedSliceView[N]); ok {
		if len(tv.s) == 0 {
			return emptyNonNilSlice[N]()
		}
		return tv.s[:len(tv.s):len(tv.s)]
	}
	if src == emptyImmutableList_ {
		return emptyNonNilSlice[N]()
	}
	if si, ok := src.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return emptyNonNilSlice[N]()
		}
		out := make([]N, len(items))
		for i, raw := range items {
			out[i] = CastNode[N](CheckNotNull(raw))
		}
		return out
	}
	out := emptyNonNilSlice[N]()
	for cur := NewIterCursor(src); cur.HasNext(); {
		out = append(out, N(cur.NextNode()))
	}
	return out
}

func ImmListToNodeSliceNullable[N ~uint32](imm Ext_com_google_common_collect_ImmutableList) []N {
	if imm == nil {
		return nil
	}
	if tv, ok := any(imm).(*immTypedSliceView[N]); ok {
		if tv == nil {
			return nil
		}
		if len(tv.s) == 0 {
			return emptyNonNilSlice[N]()
		}
		return tv.s[:len(tv.s):len(tv.s)]
	}
	if imm == emptyImmutableList_ {
		return emptyNonNilSlice[N]()
	}
	if isJavaNilAny(imm) {
		return nil
	}
	if imm.IsEmpty() {
		return emptyNonNilSlice[N]()
	}
	return ImmListCopyNodeToSlice[N](imm)
}

func AppendImmNodeCollection[N ~uint32](dst []N, src any) []N {
	if dst == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	if tv, ok := src.(*immTypedSliceView[N]); ok {
		return append(dst, tv.s...)
	}
	for cur := NewIterCursor(src); cur.HasNext(); {
		dst = append(dst, N(cur.NextNode()))
	}
	return dst
}

func SliceBuilderAddAllNode[N ~uint32](b *[]N, src any) *[]N {
	if b == nil {
		ThrowNPE()
	}
	*b = AppendImmNodeCollection(*b, src)
	return b
}

func ImmSliceSize[T any](s []T) int32 {
	if s == nil {
		ThrowNPE()
	}
	return int32(len(s))
}

func ImmSliceIsEmpty[T any](s []T) bool {
	if s == nil {
		ThrowNPE()
	}
	return len(s) == 0
}

func ImmSliceGet[T any](s []T, idx int32) T {
	if s == nil {
		ThrowNPE()
	}
	CheckElementIndex(idx, int32(len(s)))
	return s[idx]
}

func ListGetVal[T any](l any, i int32) T {
	if l == nil {
		ThrowNPE()
	}
	if tv, ok := l.(*immTypedSliceView[T]); ok {
		if tv == nil {
			ThrowNPE()
		}
		CheckElementIndex(i, int32(len(tv.s)))
		return tv.s[i]
	}
	if al, ok := l.(*jTypedArrayList[T]); ok {
		if al == nil {
			ThrowNPE()
		}
		if al.lc != nil {
			return unboxJavaElem[T](al.lc.Get(i))
		}
		CheckElementIndex(i, int32(len(al.s)))
		return al.s[i]
	}
	if isJavaNilAny(l) {
		ThrowNPE()
	}
	return unboxJavaElem[T](l.(Ext_java_util_List).Get(i))
}

type TypedIterCursor[T any] struct {
	s    []T
	idx  int
	rest JIterCursor
	raw  bool
}

func NewTypedIterCursor[T any](coll any) TypedIterCursor[T] {
	if coll == nil {
		ThrowNPE()
	}
	if tv, ok := coll.(*immTypedSliceView[T]); ok {
		if tv == nil {
			ThrowNPE()
		}
		return TypedIterCursor[T]{s: tv.s}
	}
	if al, ok := coll.(*jTypedArrayList[T]); ok {
		if al == nil {
			ThrowNPE()
		}
		if al.lc == nil {
			return TypedIterCursor[T]{s: al.s}
		}
	}
	cur := NewIterCursor(coll)
	if !cur.HasNext() {
		return TypedIterCursor[T]{}
	}
	return TypedIterCursor[T]{rest: cur, raw: true}
}

func (c *TypedIterCursor[T]) HasNext() bool {
	if !c.raw {
		return c.idx < len(c.s)
	}
	return c.rest.HasNext()
}

func (c *TypedIterCursor[T]) Next() T {
	if !c.raw {
		v := c.s[c.idx]
		c.idx++
		return v
	}
	return unboxJavaElem[T](c.rest.Next())
}

func SliceContains[T any](s []T, target any) bool {
	if s == nil {
		ThrowNPE()
	}
	if target == nil || isJavaNilAny(target) {
		return false
	}
	t := ToAny(target)
	if ss, ok := any(s).([]string); ok {
		ts, ok2 := t.(string)
		if !ok2 || ts == NullString {
			return false
		}
		for _, x := range ss {
			if x == ts {
				return true
			}
		}
		return false
	}
	for _, v := range s {
		if javalang.ObjectsEquals(boxJavaElem(v), t) {
			return true
		}
	}
	return false
}

func SliceIndexOf[T any](s []T, target any) int32 {
	if s == nil {
		ThrowNPE()
	}
	if target == nil || isJavaNilAny(target) {
		return -1
	}
	t := ToAny(target)
	for i, v := range s {
		if javalang.ObjectsEquals(boxJavaElem(v), t) {
			return int32(i)
		}
	}
	return -1
}

func SliceLastIndexOf[T any](s []T, target any) int32 {
	if s == nil {
		ThrowNPE()
	}
	if target == nil || isJavaNilAny(target) {
		return -1
	}
	t := ToAny(target)
	for i := len(s) - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(boxJavaElem(s[i]), t) {
			return int32(i)
		}
	}
	return -1
}

func ImmSliceSubList[T any](s []T, fromIdx, toIdx int32) []T {
	if s == nil {
		ThrowNPE()
	}
	CheckPositionIndexes(fromIdx, toIdx, int32(len(s)))
	if fromIdx == toIdx {
		return []T{}
	}
	return s[fromIdx:toIdx]
}

func ImmSliceReverse[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return []T{}
	}
	if len(s) == 1 {
		return s
	}
	rev := make([]T, len(s))
	for i, v := range s {
		rev[len(s)-1-i] = v
	}
	return rev
}

func ImmSliceEqualsTyped[T any](a, b []T) bool {
	if a == nil {
		ThrowNPE()
	}
	if b == nil || len(a) != len(b) {
		return false
	}
	for i := range a {
		if !javaMapValsEqual(a[i], b[i], false) {
			return false
		}
	}
	return true
}

func ImmSliceEquals[T any](s []T, other any) bool {
	if s == nil {
		ThrowNPE()
	}
	if other == nil || isJavaNilAny(other) {
		return false
	}
	if tv, ok := other.(*immTypedSliceView[T]); ok {
		return ImmSliceEqualsTyped(s, tv.s)
	}
	if al, ok := other.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
		return ImmSliceEqualsTyped(s, al.s)
	}
	lst, ok := other.(Ext_java_util_List)
	if !ok {
		return false
	}
	if int32(len(s)) != lst.Size() {
		return false
	}
	if _, isRA := other.(Ext_java_util_RandomAccess); isRA {
		for i, v := range s {
			if !javalang.ObjectsEquals(boxJavaElem(v), lst.Get(int32(i))) {
				return false
			}
		}
		return true
	}
	it := Cast[Ext_java_util_Iterator](lst.Iterator())
	for _, v := range s {
		if !it.HasNext() || !javalang.ObjectsEquals(boxJavaElem(v), it.Next()) {
			return false
		}
	}
	return !it.HasNext()
}

func ImmSliceHashCode[T any](s []T) int32 {
	if s == nil {
		ThrowNPE()
	}
	var h int32 = 1
	for _, v := range s {
		h = 31*h + javalang.ObjectHashCode(boxJavaElem(v))
	}
	return h
}

func SliceUnmodifiableList[T any](s []T) Ext_java_util_List {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return emptyImmutableList_
	}
	cp := make([]T, len(s))
	copy(cp, s)
	hasNull := false
	switch nullValKindOf[T]() {
	case kindNeverNull:
	case kindNullStringType:
		for _, v := range any(cp).([]string) {
			if v == NullString {
				hasNull = true
				break
			}
		}
	case kindZeroIsNull:
		var zero T
		for _, v := range cp {
			if any(v) == any(zero) {
				hasNull = true
				break
			}
		}
	default:
		for _, v := range cp {
			if isJavaNullMapValSlow(v) {
				hasNull = true
				break
			}
		}
	}
	if hasNull {
		v := &unmodSlicePtrView[T]{p: &cp}
		v.BindExtSelf_(v)
		return v
	}
	return &immTypedSliceView[T]{s: cp[:len(cp):len(cp)]}
}

func SliceToJavaList[T any](s []T) Ext_java_util_List {
	if s == nil {
		return nil
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = boxJavaElem(v)
	}
	return wrapArrayListOwn(items)
}

type jEmptyImmSliceIter struct{}

func (*jEmptyImmSliceIter) IsJavaNil_() bool                                                      { return false }
func (*jEmptyImmSliceIter) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*jEmptyImmSliceIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*jEmptyImmSliceIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (*jEmptyImmSliceIter) HasNext() bool { return false }
func (*jEmptyImmSliceIter) Next() any {
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}
func (*jEmptyImmSliceIter) Remove() {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}
func (*jEmptyImmSliceIter) ForEachRemaining(action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
}
func (*jEmptyImmSliceIter) typedViewLen_() (int, bool)          { return 0, true }
func (*jEmptyImmSliceIter) typedViewAt_(int) any                { return nil }
func (*jEmptyImmSliceIter) typedViewNodeAt_(int) (uint32, bool) { return 0, false }

var emptyImmSliceIter_ Ext_com_google_common_collect_UnmodifiableIterator = &jEmptyImmSliceIter{}

type immTypedSliceIter[T any] struct {
	s   []T
	pos int
}

var _ Ext_com_google_common_collect_UnmodifiableIterator = (*immTypedSliceIter[int32])(nil)

func (it *immTypedSliceIter[T]) IsJavaNil_() bool                                                   { return it == nil }
func (*immTypedSliceIter[T]) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*immTypedSliceIter[T]) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*immTypedSliceIter[T]) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (it *immTypedSliceIter[T]) HasNext() bool {
	return it.pos < len(it.s)
}
func (it *immTypedSliceIter[T]) Next() any {
	if it.pos >= len(it.s) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.s[it.pos]
	it.pos++
	return boxJavaElem(v)
}
func (it *immTypedSliceIter[T]) NextTyped() T {
	if it.pos >= len(it.s) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.s[it.pos]
	it.pos++
	return v
}
func (it *immTypedSliceIter[T]) Remove() {
	panic(javalang.NewUnsupportedOperationException("immutable"))
}

func SliceIterator[T any](s []T) Ext_com_google_common_collect_UnmodifiableIterator {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return emptyImmSliceIter_
	}
	return &immTypedSliceIter[T]{s: s}
}

type TypedSliceIter[E any] struct {
	s   []E
	idx int
}

var _ Ext_java_util_Iterator = (*TypedSliceIter[int32])(nil)

func (it *TypedSliceIter[E]) IsJavaNil_() bool                                                   { return it == nil }
func (*TypedSliceIter[E]) Implements_j2gort_Ext_java_util_Iterator()                             {}
func (*TypedSliceIter[E]) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*TypedSliceIter[E]) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
	return nil
}
func (it *TypedSliceIter[E]) HasNext() bool {
	return it != nil && it.idx < len(it.s)
}
func (it *TypedSliceIter[E]) Next() any {
	if it == nil || it.idx >= len(it.s) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.s[it.idx]
	it.idx++
	return boxJavaElem(v)
}
func (it *TypedSliceIter[E]) NextTyped() E {
	if it == nil || it.idx >= len(it.s) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.s[it.idx]
	it.idx++
	return v
}
func (it *TypedSliceIter[E]) Remove() {
	panic(javalang.NewUnsupportedOperationException("remove"))
}
func (it *TypedSliceIter[E]) ForEachRemaining(action Ext_java_util_function_Consumer) {
	CheckNotNull(action)
	if it == nil {
		return
	}
	for it.idx < len(it.s) {
		v := it.s[it.idx]
		it.idx++
		action.Accept(boxJavaElem(v))
	}
}
func (it *TypedSliceIter[E]) typedViewLen_() (int, bool) {
	if it == nil {
		return 0, true
	}
	return len(it.s) - it.idx, true
}
func (it *TypedSliceIter[E]) typedViewAt_(i int) any {
	return boxJavaElem(it.s[it.idx+i])
}
func (it *TypedSliceIter[E]) typedViewNodeAt_(i int) (uint32, bool) {
	return typedElemRawNode(it.s[it.idx+i])
}
func (it *TypedSliceIter[E]) IterStrSlice_() ([]string, bool) {
	if it == nil {
		return nil, false
	}
	if ss, ok := any(it.s).([]string); ok {
		return ss[it.idx:], true
	}
	return nil, false
}

func NewTypedSliceIter[E any](s []E) Ext_java_util_Iterator {
	if len(s) == 0 {
		return emptyImmSliceIter_
	}
	return &TypedSliceIter[E]{s: s}
}

func SliceGetFirst[T any](s []T) T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return s[0]
}

func SliceGetLast[T any](s []T) T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return s[len(s)-1]
}

func SliceRemoveAt[T any](s []T, idx int32) []T {
	if s == nil {
		ThrowNPE()
	}
	CheckElementIndex(idx, int32(len(s)))
	i := int(idx)
	var zero T
	copy(s[i:], s[i+1:])
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func SliceRemoveAtElem[T any](sp *[]T, idx int32) T {
	s := *sp
	if s == nil {
		ThrowNPE()
	}
	CheckElementIndex(idx, int32(len(s)))
	i := int(idx)
	removed := s[i]
	var zero T
	if i+1 < len(s) {
		copy(s[i:], s[i+1:])
	}
	s[len(s)-1] = zero
	*sp = s[:len(s)-1]
	return removed
}

func AppendStrCol(dst []string, col any) []string {
	if dst == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			return append(dst, ss...)
		}
	}
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, cur.NextStr())
	}
	return dst
}

func AppendNodeCol[N ~uint32](dst []N, col any) []N {
	if dst == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, N(cur.NextNode()))
	}
	return dst
}

func AppendCollection[T any](dst []T, col any) []T {
	if dst == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	for cur := NewIterCursor(col); cur.HasNext(); {
		dst = append(dst, unboxJavaElem[T](cur.Next()))
	}
	return dst
}

func SliceInsertAt[T any](s []T, idx int32, elem T) []T {
	if s == nil {
		ThrowNPE()
	}
	CheckPositionIndex(idx, int32(len(s)))
	i := int(idx)
	if i == len(s) {
		return append(s, elem)
	}
	var zero T
	s = append(s, zero)
	copy(s[i+1:], s[i:len(s)-1])
	s[i] = elem
	return s
}

func SliceReverseInPlace[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func CloneNonNilSlice[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return []T{}
	}
	out := make([]T, len(s))
	copy(out, s)
	return out
}

func CopyStrColToSlice(col any) []string {
	CheckNotNull(col)
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(ss) == 0 {
				return []string{}
			}
			out := make([]string, len(ss))
			copy(out, ss)
			return out
		}
	}
	if tv, ok := col.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return []string{}
			}
			out := make([]string, n)
			for i := 0; i < n; i++ {
				out[i] = CastString(tv.typedViewAt_(i))
			}
			return out
		}
	}
	if si, ok := col.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return []string{}
		}
		out := make([]string, len(items))
		for i, raw := range items {
			out[i] = CastString(raw)
		}
		return out
	}
	if c, ok := col.(Ext_java_util_Collection); ok {
		sz := int(c.Size())
		if sz == 0 {
			return []string{}
		}
		out := make([]string, 0, sz)
		for cur := NewIterCursor(col); cur.HasNext(); {
			out = append(out, cur.NextStr())
		}
		return out
	}
	out := []string{}
	for cur := NewIterCursor(col); cur.HasNext(); {
		out = append(out, cur.NextStr())
	}
	return out
}

func CopyNodeColToSlice[N ~uint32](col any) []N {
	CheckNotNull(col)
	if tv, ok := col.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return []N{}
			}
			out := make([]N, n)
			for i := 0; i < n; i++ {
				if u, ok3 := tv.typedViewNodeAt_(i); ok3 {
					out[i] = N(u)
				} else {
					out[i] = CastNode[N](tv.typedViewAt_(i))
				}
			}
			return out
		}
	}
	if si, ok := col.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return []N{}
		}
		out := make([]N, len(items))
		for i, raw := range items {
			out[i] = CastNode[N](raw)
		}
		return out
	}
	if c, ok := col.(Ext_java_util_Collection); ok {
		sz := int(c.Size())
		if sz == 0 {
			return []N{}
		}
		out := make([]N, 0, sz)
		for cur := NewIterCursor(col); cur.HasNext(); {
			out = append(out, N(cur.NextNode()))
		}
		return out
	}
	out := []N{}
	for cur := NewIterCursor(col); cur.HasNext(); {
		out = append(out, N(cur.NextNode()))
	}
	return out
}

func CopyCollectionToSlice[T any](col any) []T {
	CheckNotNull(col)
	if tv, ok := col.(*immTypedSliceView[T]); ok {
		return CloneNonNilSlice(tv.s)
	}
	if al, ok := col.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
		return CloneNonNilSlice(al.s)
	}
	if se, ok := col.(jSingleNonNullElem); ok {
		return []T{unboxNonNullImmElem[T](se.singleNonNullElem_())}
	}
	switch l := col.(type) {
	case *jImmList1:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxNonNullImmElem[T](l.elem[0])}
	case *jImmList2:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxNonNullImmElem[T](l.elem[0]), unboxNonNullImmElem[T](l.elem[1])}
	case *jSingletonList:
		if l == nil {
			ThrowNPE()
		}
		return []T{unboxJavaElem[T](l.elem[0])}
	}
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			if len(ss) == 0 {
				return []T{}
			}
			out := make([]T, len(ss))
			if outStr, ok3 := any(&out).(*[]string); ok3 {
				copy(*outStr, ss)
				return out
			}
			for i, str := range ss {
				out[i] = unboxJavaElem[T](str)
			}
			return out
		}
	}
	if tv, ok := col.(jTypedViewCursorSource); ok {
		if n, ok2 := tv.typedViewLen_(); ok2 {
			if n == 0 {
				return []T{}
			}
			out := make([]T, n)
			for i := 0; i < n; i++ {
				out[i] = unboxJavaElem[T](tv.typedViewAt_(i))
			}
			return out
		}
	}
	if si, ok := col.(SliceIterable); ok {
		items := si.IterSlice_()
		if len(items) == 0 {
			return []T{}
		}
		out := make([]T, len(items))
		for i, raw := range items {
			out[i] = unboxJavaElem[T](raw)
		}
		return out
	}
	if c, ok := col.(Ext_java_util_Collection); ok {
		sz := int(c.Size())
		if sz == 0 {
			return []T{}
		}
		out := make([]T, 0, sz)
		for cur := NewIterCursor(col); cur.HasNext(); {
			out = append(out, unboxJavaElem[T](cur.Next()))
		}
		return out
	}
	out := []T{}
	for cur := NewIterCursor(col); cur.HasNext(); {
		out = append(out, unboxJavaElem[T](cur.Next()))
	}
	return out
}

func CopyIteratorToSlice[T any](it Ext_java_util_Iterator) []T {
	CheckNotNull(it)
	out := []T{}
	for it.HasNext() {
		out = append(out, unboxJavaElem[T](it.Next()))
	}
	return out
}

func SliceFormatJavaString[T any](s []T) string {
	if s == nil {
		return "null"
	}
	if len(s) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range s {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(boxJavaElem(v)))
	}
	sb.WriteByte(']')
	return sb.String()
}

func sliceTimSortBy[T any](s []T, less func(a, b T) bool) {
	// Stable insertion + merge sort matching Java TimSort stability.
	n := len(s)
	if n < 2 {
		return
	}
	if n <= 24 {
		for i := 1; i < n; i++ {
			v := s[i]
			j := i
			for j > 0 && less(v, s[j-1]) {
				s[j] = s[j-1]
				j--
			}
			s[j] = v
		}
		return
	}
	buf := make([]T, n)
	var sortRange func(lo, hi int)
	sortRange = func(lo, hi int) {
		if hi-lo <= 16 {
			for i := lo + 1; i < hi; i++ {
				v := s[i]
				j := i
				for j > lo && less(v, s[j-1]) {
					s[j] = s[j-1]
					j--
				}
				s[j] = v
			}
			return
		}
		mid := lo + (hi-lo)/2
		sortRange(lo, mid)
		sortRange(mid, hi)
		if !less(s[mid], s[mid-1]) {
			return
		}
		copy(buf[lo:hi], s[lo:hi])
		i, j, k := lo, mid, lo
		for i < mid && j < hi {
			if less(buf[j], buf[i]) {
				s[k] = buf[j]
				j++
			} else {
				s[k] = buf[i]
				i++
			}
			k++
		}
		if i < mid {
			copy(s[k:hi], buf[i:mid])
		}
	}
	sortRange(0, n)
}

func SliceSortNatural[T any](s []T) {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return
	}
	var zero T
	if _, isStr := any(zero).(string); isStr {
		ss := any(s).([]string)
		for _, v := range ss {
			if v == NullString {
				ThrowNPE()
			}
		}
		sliceTimSortBy(ss, func(a, b string) bool { return JStrCompareTo(a, b) < 0 })
		return
	}
	if nullValKindOf[T]() != kindNeverNull {
		for _, v := range s {
			if isJavaNullMapVal(v) {
				ThrowNPE()
			}
		}
	}
	if len(s) == 1 {
		return
	}
	sliceTimSortBy(s, func(a, b T) bool { return CompareValues(boxJavaElem(a), boxJavaElem(b)) < 0 })
}

func SliceSortComparator[T any](s []T, cmp Ext_java_util_Comparator) {
	if s == nil {
		ThrowNPE()
	}
	if isJavaNullMapVal(cmp) {
		SliceSortNatural(s)
		return
	}
	if len(s) < 2 {
		return
	}
	sliceTimSortBy(s, func(a, b T) bool { return cmp.Compare(boxJavaElem(a), boxJavaElem(b)) < 0 })
}

type jSliceIterableView[T any] struct {
	s []T
}

func (*jSliceIterableView[T]) Implements_j2gort_Ext_java_lang_Iterable() {}
func (v *jSliceIterableView[T]) IsJavaNil_() bool                        { return v == nil }
func (v *jSliceIterableView[T]) IterStrSlice_() ([]string, bool) {
	if ss, ok := any(v.s).([]string); ok {
		return ss, true
	}
	return nil, false
}
func (v *jSliceIterableView[T]) IterSlice_() []any {
	if len(v.s) == 0 {
		return emptyAnySlice_
	}
	out := make([]any, len(v.s))
	for i, x := range v.s {
		out[i] = boxJavaElem(x)
	}
	return out
}
func (v *jSliceIterableView[T]) Iterator() any { return SliceIterator(v.s) }
func (v *jSliceIterableView[T]) ForEach(c Ext_java_util_function_Consumer) {
	CheckNotNull(c)
	for _, x := range v.s {
		c.Accept(boxJavaElem(x))
	}
}
func (v *jSliceIterableView[T]) String() string { return SliceFormatJavaString(v.s) }

func SliceReadOnlyIterable[T any](s []T) Ext_java_lang_Iterable {
	if s == nil {
		ThrowNPE()
	}
	return &jSliceIterableView[T]{s: s}
}

type unmodSlicePtrView[T any] struct {
	Ext_java_util_AbstractListBase
	p *[]T
}

type unmodSlicePtrIter[T any] struct {
	p   *[]T
	pos int
}

func (*unmodSlicePtrIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *unmodSlicePtrIter[T]) IsJavaNil_() bool                       { return it == nil }
func (it *unmodSlicePtrIter[T]) HasNext() bool                          { return it.pos < len(*it.p) }
func (it *unmodSlicePtrIter[T]) Next() any {
	s := *it.p
	if it.pos >= len(s) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := s[it.pos]
	it.pos++
	return boxJavaElem(v)
}
func (it *unmodSlicePtrIter[T]) Remove() {
	panic(javalang.NewUnsupportedOperationException("unmodifiable"))
}

func (*unmodSlicePtrView[T]) Implements_j2gort_Ext_java_util_RandomAccess() {}
func (v *unmodSlicePtrView[T]) IsJavaNil_() bool                            { return v == nil }
func (v *unmodSlicePtrView[T]) Size() int32                                 { return int32(len(*v.p)) }
func (v *unmodSlicePtrView[T]) IsEmpty() bool                               { return len(*v.p) == 0 }
func (v *unmodSlicePtrView[T]) Get(i int32) any {
	s := *v.p
	CheckElementIndex(i, int32(len(s)))
	return boxJavaElem(s[i])
}
func (v *unmodSlicePtrView[T]) Contains(x any) bool {
	t := ToAny(x)
	for _, e := range *v.p {
		if javalang.ObjectsEquals(boxJavaElem(e), t) {
			return true
		}
	}
	return false
}
func (v *unmodSlicePtrView[T]) IndexOf_1(x any) int32 {
	t := ToAny(x)
	for i, e := range *v.p {
		if javalang.ObjectsEquals(boxJavaElem(e), t) {
			return int32(i)
		}
	}
	return -1
}
func (v *unmodSlicePtrView[T]) LastIndexOf_1(x any) int32 {
	t := ToAny(x)
	s := *v.p
	for i := len(s) - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(boxJavaElem(s[i]), t) {
			return int32(i)
		}
	}
	return -1
}
func (v *unmodSlicePtrView[T]) IterSlice_() []any {
	s := *v.p
	if len(s) == 0 {
		return emptyAnySlice_
	}
	out := make([]any, len(s))
	for i, e := range s {
		out[i] = boxJavaElem(e)
	}
	return out
}
func (v *unmodSlicePtrView[T]) Iterator() any {
	if len(*v.p) == 0 {
		return emptyImmSliceIter_
	}
	return &unmodSlicePtrIter[T]{p: v.p}
}
func (v *unmodSlicePtrView[T]) ForEach(c Ext_java_util_function_Consumer) {
	CheckNotNull(c)
	for _, e := range *v.p {
		c.Accept(boxJavaElem(e))
	}
}
func (v *unmodSlicePtrView[T]) Equals(o any) bool { return ImmSliceEquals(*v.p, o) }
func (v *unmodSlicePtrView[T]) HashCode() int32   { return ImmSliceHashCode(*v.p) }
func (v *unmodSlicePtrView[T]) String() string    { return SliceFormatJavaString(*v.p) }
func (v *unmodSlicePtrView[T]) Clear() {
	panic(javalang.NewUnsupportedOperationException("unmodifiable"))
}
func (v *unmodSlicePtrView[T]) Add_1(any) bool {
	panic(javalang.NewUnsupportedOperationException("unmodifiable"))
}
func (v *unmodSlicePtrView[T]) Remove_1_any(any) bool {
	panic(javalang.NewUnsupportedOperationException("unmodifiable"))
}

func SliceUnmodifiableListLive[T any](p *[]T) Ext_java_util_List {
	if p == nil || *p == nil {
		ThrowNPE()
	}
	v := &unmodSlicePtrView[T]{p: p}
	v.BindExtSelf_(v)
	return v
}

const orderedMapIndexThreshold = 12
const orderedSetIndexThreshold = 8

// OrderedSet is a slim value struct preserving Java LinkedHashSet insertion order
// with zero inline array bloat and lazy backing-slice allocation on first Add.
type OrderedSet[T comparable] struct {
	order    []T
	idx      map[T]int32
	holes    int32
	firstIdx int32
}

type OrderedSetCursor[T comparable] struct {
	s   *OrderedSet[T]
	pos int32
}

func NewOrderedSet[T comparable]() OrderedSet[T] {
	return OrderedSet[T]{}
}

func NewOrderedSetPtr[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{}
}

func NewOrderedSetCap[T comparable](cap int32) OrderedSet[T] {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	var s OrderedSet[T]
	if cap > 4 {
		s.order = make([]T, 0, cap)
	}
	return s
}

func NewOrderedSetPtrCap[T comparable](cap int32) *OrderedSet[T] {
	s := NewOrderedSetCap[T](cap)
	return &s
}

func OrderedSetFromSet[T comparable](src OrderedSet[T]) OrderedSet[T] {
	n := src.Size()
	if n == 0 {
		return OrderedSet[T]{}
	}
	var dst OrderedSet[T]
	dst.order = make([]T, 0, max(4, int(n)))
	dst.AddAllSet(src)
	return dst
}

func OrderedSetPtrFromSet[T comparable](src *OrderedSet[T]) *OrderedSet[T] {
	if src == nil {
		ThrowNPE()
	}
	dst := OrderedSetFromSet(*src)
	return &dst
}

func OrderedSetFromStrCol(col any) OrderedSet[string] {
	var s OrderedSet[string]
	OrderedSetAddAllStr(&s, col)
	return s
}

func OrderedSetPtrFromStrCol(col any) *OrderedSet[string] {
	s := OrderedSetFromStrCol(col)
	return &s
}

func OrderedSetFromNodeCol[N ~uint32](col any) OrderedSet[N] {
	var s OrderedSet[N]
	OrderedSetAddAllNode(&s, col)
	return s
}

func OrderedSetPtrFromNodeCol[N ~uint32](col any) *OrderedSet[N] {
	s := OrderedSetFromNodeCol[N](col)
	return &s
}

func OrderedSetFromCol[T comparable](col any) OrderedSet[T] {
	var s OrderedSet[T]
	OrderedSetAddAllCol(&s, col)
	return s
}

func OrderedSetPtrFromCol[T comparable](col any) *OrderedSet[T] {
	s := OrderedSetFromCol[T](col)
	return &s
}

func (s *OrderedSet[T]) Init() {}

func (s *OrderedSet[T]) InitCap(cap int32) {
	if s == nil {
		ThrowNPE()
	}
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	if cap > 4 && s.order == nil {
		s.order = make([]T, 0, cap)
	}
}

func (s *OrderedSet[T]) IsNil() bool      { return s == nil }
func (s *OrderedSet[T]) IsJavaNil_() bool { return s == nil }

func (s OrderedSet[T]) Size() int32 {
	return int32(len(s.order)) - s.holes
}

func (s OrderedSet[T]) IsEmpty() bool {
	return int32(len(s.order)) == s.holes
}

func (s OrderedSet[T]) Contains(v T) bool {
	if int32(len(s.order)) == s.holes {
		return false
	}
	if s.idx != nil {
		return s.idx[v] != 0
	}
	for _, x := range s.order {
		if x == v {
			return true
		}
	}
	return false
}

func (s OrderedSet[T]) GetFirst() T {
	if int32(len(s.order)) == s.holes {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return s.order[s.firstIdx]
}

func (s *OrderedSet[T]) Add(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	if s.order == nil {
		s.order = make([]T, 1, 4)
		s.order[0] = v
		return true
	}
	if s.idx != nil {
		if s.idx[v] != 0 {
			return false
		}
		s.order = append(s.order, v)
		s.idx[v] = int32(len(s.order))
		return true
	}
	for _, x := range s.order {
		if x == v {
			return false
		}
	}
	s.order = append(s.order, v)
	if len(s.order) > orderedSetIndexThreshold {
		m := make(map[T]int32, 16)
		for i, x := range s.order {
			m[x] = int32(i + 1)
		}
		s.idx = m
	}
	return true
}

func (s *OrderedSet[T]) Remove(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	if int32(len(s.order)) == s.holes {
		return false
	}
	if s.idx == nil {
		for i, x := range s.order {
			if x == v {
				var zero T
				copy(s.order[i:], s.order[i+1:])
				last := len(s.order) - 1
				s.order[last] = zero
				s.order = s.order[:last]
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
	slot := pos1 - 1
	s.order[slot] = zero
	s.holes++
	live := int32(len(s.order)) - s.holes
	if live == 0 {
		s.order = s.order[:0]
		s.idx = nil
		s.holes = 0
		s.firstIdx = 0
	} else if len(s.order) > 16 && int(live)*2 < len(s.order) {
		s.compact()
	} else if slot == s.firstIdx {
		for s.firstIdx++; int(s.firstIdx) < len(s.order); s.firstIdx++ {
			if s.idx[s.order[s.firstIdx]] == s.firstIdx+1 {
				break
			}
		}
	}
	return true
}

func (s *OrderedSet[T]) compact() {
	w := 0
	for i := int(s.firstIdx); i < len(s.order); i++ {
		v := s.order[i]
		if s.idx[v] == int32(i+1) {
			s.order[w] = v
			s.idx[v] = int32(w + 1)
			w++
		}
	}
	clear(s.order[w:])
	s.order = s.order[:w]
	s.holes = 0
	s.firstIdx = 0
}

func (s *OrderedSet[T]) Clear() {
	if s == nil {
		ThrowNPE()
	}
	if s.order != nil {
		clear(s.order)
		s.order = s.order[:0]
	}
	s.idx = nil
	s.holes = 0
	s.firstIdx = 0
}

func (s OrderedSet[T]) EqualsSet(other OrderedSet[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	if s.IsEmpty() {
		return true
	}
	if s.holes > 0 {
		for i := int(s.firstIdx); i < len(s.order); i++ {
			v := s.order[i]
			if s.idx[v] == int32(i+1) && !other.Contains(v) {
				return false
			}
		}
		return true
	}
	for _, v := range s.order {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s *OrderedSet[T]) EqualsSetPtr(other *OrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	return s.EqualsSet(*other)
}

func (s OrderedSet[T]) HashCode() int32 {
	if s.IsEmpty() {
		return 0
	}
	var h int32
	if s.holes > 0 {
		for i := int(s.firstIdx); i < len(s.order); i++ {
			v := s.order[i]
			if s.idx[v] == int32(i+1) {
				h += javalang.ObjectHashCode(boxJavaElem(v))
			}
		}
		return h
	}
	for _, v := range s.order {
		h += javalang.ObjectHashCode(boxJavaElem(v))
	}
	return h
}

func (s *OrderedSet[T]) AddAllSet(other OrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	if other.IsEmpty() {
		return false
	}
	changed := false
	if other.holes > 0 {
		for i := int(other.firstIdx); i < len(other.order); i++ {
			v := other.order[i]
			if other.idx[v] == int32(i+1) && s.Add(v) {
				changed = true
			}
		}
		return changed
	}
	for _, v := range other.order {
		if s.Add(v) {
			changed = true
		}
	}
	return changed
}

func (s *OrderedSet[T]) AddAllSetPtr(other *OrderedSet[T]) bool {
	if s == nil || other == nil {
		ThrowNPE()
	}
	return s.AddAllSet(*other)
}

func OrderedSetAddAllStr(s *OrderedSet[string], col any) bool {
	if s == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	changed := false
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			for _, k := range ss {
				if s.Add(k) {
					changed = true
				}
			}
			return changed
		}
	}
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(cur.NextStr()) {
			changed = true
		}
	}
	return changed
}

func OrderedSetAddAllNode[N ~uint32](s *OrderedSet[N], col any) bool {
	if s == nil {
		ThrowNPE()
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
	if s == nil {
		ThrowNPE()
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
	if s == nil {
		ThrowNPE()
	}
	return OrderedSetCursor[T]{s: s, pos: s.firstIdx}
}

func (c *OrderedSetCursor[T]) HasNext() bool {
	s := c.s
	if int32(len(s.order)) == s.holes {
		return false
	}
	if s.holes > 0 {
		return c.advanceDead()
	}
	return int(c.pos) < len(s.order)
}

func (c *OrderedSetCursor[T]) advanceDead() bool {
	s := c.s
	items := s.order
	idx := s.idx
	for i := int(c.pos); i < len(items); i++ {
		if idx[items[i]] == int32(i+1) {
			c.pos = int32(i)
			return true
		}
	}
	c.pos = int32(len(items))
	return false
}

func (c *OrderedSetCursor[T]) Next() T {
	i := c.pos
	c.pos = i + 1
	return c.s.order[i]
}

func (s *OrderedSet[T]) String() string {
	if s == nil {
		return "null"
	}
	if s.IsEmpty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for cur := s.Iter(); cur.HasNext(); {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(boxJavaElem(cur.Next())))
	}
	sb.WriteByte(']')
	return sb.String()
}

func (s OrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	live := s.Size()
	if live == 0 {
		return emptyImmutableSet_
	}
	var zero T
	if _, isStr := any(zero).(string); isStr {
		if live == 1 {
			return ImmSetOfStr1(any(s.GetFirst()).(string))
		}
		strs := make([]string, 0, live)
		if s.holes > 0 {
			for i, v := range s.order {
				if s.idx[v] == int32(i+1) {
					str := any(v).(string)
					if str == NullString {
						ThrowNPE()
					}
					strs = append(strs, str)
				}
			}
		} else {
			for _, v := range s.order {
				str := any(v).(string)
				if str == NullString {
					ThrowNPE()
				}
				strs = append(strs, str)
			}
		}
		return finishImmStrSetOwn(strs, nil)
	}
	if live == 1 {
		return newImmSet1(boxJavaElem(s.GetFirst()))
	}
	if live <= 8 {
		elems := make([]any, 0, live)
		if s.holes > 0 {
			for i, v := range s.order {
				if s.idx[v] == int32(i+1) {
					elems = append(elems, CheckNotNull(boxJavaElem(v)))
				}
			}
		} else {
			for _, v := range s.order {
				elems = append(elems, CheckNotNull(boxJavaElem(v)))
			}
		}
		return finishImmSetMulti(elems, nil)
	}
	sc := &JSetCore{}
	sc.InitExpectedSize(live)
	if s.holes > 0 {
		for i, v := range s.order {
			if s.idx[v] == int32(i+1) {
				sc.Add(CheckNotNull(boxJavaElem(v)))
			}
		}
	} else {
		for _, v := range s.order {
			sc.Add(CheckNotNull(boxJavaElem(v)))
		}
	}
	return wrapImmutableSetFromSetCoreOwn(sc)
}

func (s OrderedSet[T]) ToJavaSet() *Ext_java_util_LinkedHashSet {
	live := s.Size()
	obj := &Ext_java_util_LinkedHashSet{}
	obj.BindExtSelf_(obj)
	obj.inlineSet.InitExpectedSize(live)
	if live == 0 {
		return obj
	}
	var zero T
	if _, isStr := any(zero).(string); isStr {
		if s.holes > 0 {
			for i, v := range s.order {
				if s.idx[v] == int32(i+1) {
					obj.inlineSet.AddStr(any(v).(string))
				}
			}
		} else {
			for _, v := range s.order {
				obj.inlineSet.AddStr(any(v).(string))
			}
		}
		return obj
	}
	if s.holes > 0 {
		for i, v := range s.order {
			if s.idx[v] == int32(i+1) {
				obj.inlineSet.Add(boxJavaElem(v))
			}
		}
	} else {
		for _, v := range s.order {
			obj.inlineSet.Add(boxJavaElem(v))
		}
	}
	return obj
}

func ArrayDequeFromOrderedSet[T comparable](s OrderedSet[T]) *Ext_java_util_ArrayDeque {
	live := s.Size()
	obj := &Ext_java_util_ArrayDeque{}
	obj.BindExtSelf_(obj)
	dc := &obj.inlineCore
	dc.InitCap(live)
	if live == 0 {
		return obj
	}
	if s.holes > 0 {
		for i, v := range s.order {
			if s.idx[v] == int32(i+1) {
				dc.AddLast(boxJavaElem(v))
			}
		}
	} else {
		for _, v := range s.order {
			dc.AddLast(boxJavaElem(v))
		}
	}
	return obj
}

func ArrayDequeFromOrderedSetPtr[T comparable](s *OrderedSet[T]) *Ext_java_util_ArrayDeque {
	if s == nil {
		ThrowNPE()
	}
	return ArrayDequeFromOrderedSet(*s)
}

// LocalOrderedSet is a stack-only OrderedSet with [4]T inline storage so standalone stack
// local sets with <= 4 items allocate zero heap memory while struct fields remain slim.
type LocalOrderedSet[T comparable] struct {
	inline   [4]T
	items    []T
	idx      map[T]int32
	live     int32
	firstIdx int32
	lenInl   uint8
}

type LocalOrderedSetCursor[T comparable] struct {
	s   *LocalOrderedSet[T]
	pos int32
}

func NewLocalOrderedSet[T comparable]() LocalOrderedSet[T] {
	return LocalOrderedSet[T]{}
}

func NewLocalOrderedSetCap[T comparable](cap int32) LocalOrderedSet[T] {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	var s LocalOrderedSet[T]
	if cap > 4 {
		s.items = make([]T, 0, cap)
	}
	return s
}

func LocalOrderedSetFromStrCol(col any) LocalOrderedSet[string] {
	var s LocalOrderedSet[string]
	LocalOrderedSetAddAllStr(&s, col)
	return s
}

func LocalOrderedSetFromNodeCol[N ~uint32](col any) LocalOrderedSet[N] {
	var s LocalOrderedSet[N]
	LocalOrderedSetAddAllNode(&s, col)
	return s
}

func LocalOrderedSetFromCol[T comparable](col any) LocalOrderedSet[T] {
	var s LocalOrderedSet[T]
	LocalOrderedSetAddAllCol(&s, col)
	return s
}

func (s LocalOrderedSet[T]) Size() int32   { return s.live }
func (s LocalOrderedSet[T]) IsEmpty() bool { return s.live == 0 }
func (s LocalOrderedSet[T]) Contains(v T) bool {
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

func (s LocalOrderedSet[T]) GetFirst() T {
	if s.live == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if s.items == nil {
		return s.inline[0]
	}
	return s.items[s.firstIdx]
}

func (s *LocalOrderedSet[T]) Add(v T) bool {
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
	if len(s.items) > orderedSetIndexThreshold {
		m := make(map[T]int32, 16)
		for i, x := range s.items {
			m[x] = int32(i + 1)
		}
		s.idx = m
	}
	return true
}

func (s *LocalOrderedSet[T]) Remove(v T) bool {
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
	slot := pos1 - 1
	s.items[slot] = zero
	s.live--
	if s.live == 0 {
		s.items = s.items[:0]
		s.idx = nil
		s.firstIdx = 0
	} else if len(s.items) > 16 && int(s.live)*2 < len(s.items) {
		w := 0
		for i := int(s.firstIdx); i < len(s.items); i++ {
			x := s.items[i]
			if s.idx[x] == int32(i+1) {
				s.items[w] = x
				s.idx[x] = int32(w + 1)
				w++
			}
		}
		clear(s.items[w:])
		s.items = s.items[:w]
		s.firstIdx = 0
	} else if slot == s.firstIdx {
		for s.firstIdx++; int(s.firstIdx) < len(s.items); s.firstIdx++ {
			if s.idx[s.items[s.firstIdx]] == s.firstIdx+1 {
				break
			}
		}
	}
	return true
}

func (s *LocalOrderedSet[T]) Clear() {
	var z [4]T
	s.inline = z
	s.lenInl = 0
	if s.items != nil {
		clear(s.items)
		s.items = s.items[:0]
	}
	s.idx = nil
	s.live = 0
	s.firstIdx = 0
}

func (s *LocalOrderedSet[T]) HashCode() int32 {
	if s == nil {
		ThrowNPE()
	}
	return s.toOrderedSetView().HashCode()
}

func (s *LocalOrderedSet[T]) String() string {
	if s == nil {
		return "null"
	}
	v := s.toOrderedSetView()
	return (&v).String()
}

func LocalOrderedSetAddAllStr(s *LocalOrderedSet[string], col any) bool {
	if s == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	changed := false
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			for _, k := range ss {
				if s.Add(k) {
					changed = true
				}
			}
			return changed
		}
	}
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(cur.NextStr()) {
			changed = true
		}
	}
	return changed
}

func LocalOrderedSetAddAllNode[N ~uint32](s *LocalOrderedSet[N], col any) bool {
	if s == nil {
		ThrowNPE()
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

func LocalOrderedSetAddAllCol[T comparable](s *LocalOrderedSet[T], col any) bool {
	if s == nil {
		ThrowNPE()
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

func (s *LocalOrderedSet[T]) Iter() LocalOrderedSetCursor[T] {
	if s == nil {
		ThrowNPE()
	}
	return LocalOrderedSetCursor[T]{s: s, pos: s.firstIdx}
}

func (c *LocalOrderedSetCursor[T]) HasNext() bool {
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

func (c *LocalOrderedSetCursor[T]) Next() T {
	s := c.s
	i := c.pos
	c.pos = i + 1
	if s.items == nil {
		return s.inline[i]
	}
	return s.items[i]
}

func (s *LocalOrderedSet[T]) toOrderedSetView() OrderedSet[T] {
	if s.live == 0 {
		return OrderedSet[T]{}
	}
	if s.items == nil {
		return OrderedSet[T]{order: s.inline[:s.lenInl]}
	}
	return OrderedSet[T]{order: s.items, idx: s.idx, holes: int32(len(s.items)) - s.live, firstIdx: s.firstIdx}
}

func (s *LocalOrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	if s == nil {
		ThrowNPE()
	}
	return s.toOrderedSetView().ToImmutableSet()
}

func ArrayDequeFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) *Ext_java_util_ArrayDeque {
	if s == nil {
		ThrowNPE()
	}
	return ArrayDequeFromOrderedSet(s.toOrderedSetView())
}

type jIsNiler interface {
	IsNil() bool
}

const (
	kindNullFallback   uint8 = 0
	kindNeverNull      uint8 = 1
	kindZeroIsNull     uint8 = 2
	kindNullStringType uint8 = 3
)

var nullValKindMap sync.Map

func computeNullValKind[V any](rt reflect.Type) uint8 {
	if rt == nil || rt.Kind() == reflect.Interface {
		return kindNullFallback
	}
	if rt == reflect.TypeFor[string]() {
		return kindNullStringType
	}
	var zero V
	if nc, ok := any(zero).(jJavaNilChecker); ok {
		if !nc.IsJavaNil_() {
			return kindNeverNull
		}
		if rt.Comparable() {
			return kindZeroIsNull
		}
		return kindNullFallback
	}
	if rn, ok := any(zero).(jRawNodeProvider); ok {
		if rn.RawNodeId_() != 0 {
			return kindNeverNull
		}
		if rt.Comparable() {
			return kindZeroIsNull
		}
		return kindNullFallback
	}
	if rt.Implements(reflect.TypeFor[jIsNiler]()) {
		return kindNullFallback
	}
	switch rt.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Struct, reflect.Array:
		return kindNeverNull
	case reflect.Uint32, reflect.Pointer, reflect.UnsafePointer:
		return kindZeroIsNull
	}
	return kindNullFallback
}

func nullValKindOf[V any]() uint8 {
	rt := reflect.TypeFor[V]()
	if k, ok := nullValKindMap.Load(rt); ok {
		return k.(uint8)
	}
	k := computeNullValKind[V](rt)
	nullValKindMap.Store(rt, k)
	return k
}

func isJavaNullMapVal[V any](v V) bool {
	if any(v) == nil {
		return true
	}
	switch x := any(v).(type) {
	case string:
		return x == NullString
	case int32, int64, bool, float64, float32, uint16, int8, int16, uint8, int, uint64:
		return false
	default:
		switch nullValKindOf[V]() {
		case kindNeverNull:
			return false
		case kindZeroIsNull:
			var zero V
			return any(v) == any(zero)
		case kindNullStringType:
			return any(v).(string) == NullString
		default:
			return isJavaNullMapValSlow(v)
		}
	}
}

//go:noinline
func isJavaNullMapValSlow[V any](v V) bool {
	var zero V
	switch any(zero).(type) {
	case nil:
		x := any(v)
		if x == nil {
			return true
		}
		if nc, ok := x.(jJavaNilChecker); ok {
			return nc.IsJavaNil_()
		}
		if rn, ok := x.(jRawNodeProvider); ok {
			return rn.RawNodeId_() == 0
		}
		if nd, ok := x.(nodestore.Node); ok {
			return uint32(nd) == 0
		}
		return isJavaNilAny(x)
	case string:
		return any(v).(string) == NullString
	case int32, int64, bool, float64, float32, uint16, int8, int16:
		return false
	case uint32:
		return any(v).(uint32) == 0
	case nodestore.Node:
		return uint32(any(v).(nodestore.Node)) == 0
	case jJavaNilChecker:
		return any(v).(jJavaNilChecker).IsJavaNil_()
	case jRawNodeProvider:
		return any(v).(jRawNodeProvider).RawNodeId_() == 0
	default:
		return isJavaNilAny(any(v))
	}
}

func javaMapValsEqual[V any](a, b V, identity bool) bool {
	ba, bb := boxJavaElem(a), boxJavaElem(b)
	if identity {
		return ba == bb
	}
	return javalang.ObjectsEquals(ba, bb)
}

func javaMapEntryHash[K, V any](k K, v V, identity bool) int32 {
	bk, bv := boxJavaElem(k), boxJavaElem(v)
	if identity {
		return IdentityHash(bk) ^ IdentityHash(bv)
	}
	return ObjectHash(bk) ^ ObjectHash(bv)
}

// OrderedMap is a slim 40-byte struct preserving Java LinkedHashMap/HashMap/IdentityHashMap order
// for Go-comparable identity/primitive/string/Node keys with zero any-boxing and lazy backing allocation.
type OrderedMapEntry[K comparable, V any] struct {
	Key K
	Val V
}

type OrderedMap[K comparable, V any] struct {
	entries  []OrderedMapEntry[K, V]
	idx      map[K]int32
	holes    int32
	firstIdx int32
}

const (
	ordMapPoolSmallMax = 8192
	ordMapPoolMedMax   = 1024
	ordMapPoolLargeMax = 16
	ordMapPoolHugeMax  = 4
	ordMapPtrSmallMax  = 16384
	ordMapPtrLargeMax  = 64
	ordMapPtrHugeMax   = 4
)

type orderedMapBuf[K comparable, V any] struct {
	entries []OrderedMapEntry[K, V]
	idx     map[K]int32
}

type orderedMapTypePool[K comparable, V any] struct {
	mu    sync.Mutex
	small []orderedMapBuf[K, V]
	med   []orderedMapBuf[K, V]
	large []orderedMapBuf[K, V]
	huge  []orderedMapBuf[K, V]
}

type orderedMapPtrTypePool[K comparable, V any] struct {
	mu        sync.Mutex
	smallFree []*OrderedMap[K, V]
	largeFree []*OrderedMap[K, V]
	hugeFree  []*OrderedMap[K, V]
}

var (
	orderedMapPools    sync.Map
	orderedMapPtrPools sync.Map
)

func getOrderedMapPool[K comparable, V any]() *orderedMapTypePool[K, V] {
	rt := reflect.TypeFor[orderedMapBuf[K, V]]()
	if v, ok := orderedMapPools.Load(rt); ok {
		return v.(*orderedMapTypePool[K, V])
	}
	p := &orderedMapTypePool[K, V]{}
	actual, loaded := orderedMapPools.LoadOrStore(rt, p)
	if !loaded {
		purge := func() {
			p.mu.Lock()
			p.small = nil
			p.med = nil
			p.large = nil
			p.huge = nil
			p.mu.Unlock()
		}
		registerCompPoolReset(purge)
		registerRuntimePoolPurger(purge)
	}
	return actual.(*orderedMapTypePool[K, V])
}

func getOrderedMapPtrPool[K comparable, V any]() *orderedMapPtrTypePool[K, V] {
	rt := reflect.TypeFor[OrderedMap[K, V]]()
	if v, ok := orderedMapPtrPools.Load(rt); ok {
		return v.(*orderedMapPtrTypePool[K, V])
	}
	p := &orderedMapPtrTypePool[K, V]{}
	actual, loaded := orderedMapPtrPools.LoadOrStore(rt, p)
	if !loaded {
		purge := func() {
			p.mu.Lock()
			p.smallFree = nil
			p.largeFree = nil
			p.hugeFree = nil
			p.mu.Unlock()
		}
		registerCompPoolReset(purge)
		registerRuntimePoolPurger(purge)
	}
	return actual.(*orderedMapPtrTypePool[K, V])
}

func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	p := getOrderedMapPool[K, V]()
	p.mu.Lock()
	if n := len(p.small); n > 0 {
		b := p.small[n-1]
		p.small[n-1] = orderedMapBuf[K, V]{}
		p.small = p.small[:n-1]
		p.mu.Unlock()
		return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
	}
	p.mu.Unlock()
	return OrderedMap[K, V]{}
}

func AcquireOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint <= 64 {
		m := AcquireOrderedMap[K, V]()
		if cap(m.entries) == 0 && capHint > 8 {
			if capHint >= 64 {
				m.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(capHint))
			} else {
				m.entries = make([]OrderedMapEntry[K, V], 0, capHint)
			}
		}
		if m.idx == nil && int(capHint) > orderedMapIndexThreshold {
			m.idx = make(map[K]int32, capHint)
		}
		return m
	}
	p := getOrderedMapPool[K, V]()
	if capHint < 1024 {
		p.mu.Lock()
		if n := len(p.med); n > 0 {
			b := p.med[n-1]
			p.med[n-1] = orderedMapBuf[K, V]{}
			p.med = p.med[:n-1]
			p.mu.Unlock()
			return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
		}
		if n := len(p.large); n > 0 {
			b := p.large[n-1]
			p.large[n-1] = orderedMapBuf[K, V]{}
			p.large = p.large[:n-1]
			p.mu.Unlock()
			return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
		}
		p.mu.Unlock()
		return NewOrderedMapCap[K, V](capHint)
	}
	p.mu.Lock()
	if n := len(p.huge); n > 0 {
		b := p.huge[n-1]
		p.huge[n-1] = orderedMapBuf[K, V]{}
		p.huge = p.huge[:n-1]
		p.mu.Unlock()
		return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
	}
	if capHint <= 65536 {
		if n := len(p.large); n > 0 {
			b := p.large[n-1]
			p.large[n-1] = orderedMapBuf[K, V]{}
			p.large = p.large[:n-1]
			p.mu.Unlock()
			return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
		}
		if capHint <= 2048 {
			if n := len(p.med); n > 0 {
				b := p.med[n-1]
				p.med[n-1] = orderedMapBuf[K, V]{}
				p.med = p.med[:n-1]
				p.mu.Unlock()
				return OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
			}
		}
	}
	p.mu.Unlock()
	return NewOrderedMapCap[K, V](capHint)
}

func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	p := getOrderedMapPtrPool[K, V]()
	p.mu.Lock()
	if n := len(p.smallFree); n > 0 {
		m := p.smallFree[n-1]
		p.smallFree[n-1] = nil
		p.smallFree = p.smallFree[:n-1]
		p.mu.Unlock()
		return m
	}
	p.mu.Unlock()
	return &OrderedMap[K, V]{}
}

func AcquireOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	p := getOrderedMapPtrPool[K, V]()
	var m *OrderedMap[K, V]
	p.mu.Lock()
	if capHint >= 1024 {
		if n := len(p.hugeFree); n > 0 {
			m = p.hugeFree[n-1]
			p.hugeFree[n-1] = nil
			p.hugeFree = p.hugeFree[:n-1]
		} else if n := len(p.largeFree); n > 0 {
			m = p.largeFree[n-1]
			p.largeFree[n-1] = nil
			p.largeFree = p.largeFree[:n-1]
		} else if n := len(p.smallFree); n > 0 {
			m = p.smallFree[n-1]
			p.smallFree[n-1] = nil
			p.smallFree = p.smallFree[:n-1]
		}
	} else if capHint > 256 {
		if n := len(p.largeFree); n > 0 {
			m = p.largeFree[n-1]
			p.largeFree[n-1] = nil
			p.largeFree = p.largeFree[:n-1]
		} else if n := len(p.hugeFree); n > 0 {
			m = p.hugeFree[n-1]
			p.hugeFree[n-1] = nil
			p.hugeFree = p.hugeFree[:n-1]
		}
	} else {
		if n := len(p.smallFree); n > 0 {
			m = p.smallFree[n-1]
			p.smallFree[n-1] = nil
			p.smallFree = p.smallFree[:n-1]
		}
	}
	p.mu.Unlock()
	if m == nil {
		m = &OrderedMap[K, V]{}
	}
	if capHint >= 64 && cap(m.entries) < 64 {
		m.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(capHint))
	} else if cap(m.entries) == 0 && capHint > 0 {
		m.entries = make([]OrderedMapEntry[K, V], 0, capHint)
	}
	if m.idx == nil && int(capHint) > orderedMapIndexThreshold {
		m.idx = make(map[K]int32, capHint)
	}
	return m
}

func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m == nil {
		return
	}
	c := cap(m.entries)
	if c > 4194304 {
		es := m.entries
		clear(es)
		m.entries = nil
		m.idx = nil
		m.holes = 0
		m.firstIdx = 0
		if c >= 64 && c <= pooledSliceMaxCap {
			ReleasePooledSlice(es)
		}
		return
	}
	m.Clear()
	p := getOrderedMapPtrPool[K, V]()
	kept := false
	if c > 65536 {
		p.mu.Lock()
		if len(p.hugeFree) < ordMapPtrHugeMax {
			p.hugeFree = append(p.hugeFree, m)
			kept = true
		}
		p.mu.Unlock()
	} else if c > 256 {
		p.mu.Lock()
		if len(p.largeFree) < ordMapPtrLargeMax {
			p.largeFree = append(p.largeFree, m)
			kept = true
		}
		p.mu.Unlock()
	} else {
		p.mu.Lock()
		if len(p.smallFree) < ordMapPtrSmallMax {
			p.smallFree = append(p.smallFree, m)
			kept = true
		}
		p.mu.Unlock()
	}
	if !kept && c >= 64 && c <= pooledSliceMaxCap {
		es := m.entries
		m.entries = nil
		ReleasePooledSlice(es)
	}
}

func ReleaseOrderedMapPtr[K comparable, V any](m *OrderedMap[K, V]) {
	ReleaseOrderedMap(m)
}

func (m *OrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	es := m.entries
	idx := m.idx
	holes := m.holes
	firstIdx := m.firstIdx
	*m = OrderedMap[K, V]{}
	c := cap(es)
	if c == 0 && idx == nil {
		return
	}
	if c > 4194304 {
		if c >= 64 && c <= pooledSliceMaxCap {
			ReleasePooledSlice(es)
		}
		return
	}
	n := len(es)
	if idx != nil && len(idx) != 0 {
		live := n - int(holes)
		if live > 0 && live <= 32 && n-int(firstIdx) <= 64 && c >= 128 {
			for i := int(firstIdx); i < n; i++ {
				delete(idx, es[i].Key)
			}
			if len(idx) != 0 {
				clear(idx)
			}
		} else {
			clear(idx)
		}
	}
	if n > 0 {
		clear(es)
		es = es[:0]
	}
	p := getOrderedMapPool[K, V]()
	kept := false
	if c <= 64 {
		p.mu.Lock()
		if len(p.small) < ordMapPoolSmallMax {
			p.small = append(p.small, orderedMapBuf[K, V]{entries: es, idx: idx})
			kept = true
		}
		p.mu.Unlock()
	} else if c <= 2048 {
		p.mu.Lock()
		if len(p.med) < ordMapPoolMedMax {
			p.med = append(p.med, orderedMapBuf[K, V]{entries: es, idx: idx})
			kept = true
		}
		p.mu.Unlock()
	} else if c <= 65536 {
		p.mu.Lock()
		if len(p.large) < ordMapPoolLargeMax {
			p.large = append(p.large, orderedMapBuf[K, V]{entries: es, idx: idx})
			kept = true
		}
		p.mu.Unlock()
	} else {
		p.mu.Lock()
		if len(p.huge) < ordMapPoolHugeMax {
			p.huge = append(p.huge, orderedMapBuf[K, V]{entries: es, idx: idx})
			kept = true
		}
		p.mu.Unlock()
	}
	if !kept && c >= 64 && c <= pooledSliceMaxCap {
		ReleasePooledSlice(es)
	}
}

func (m *OrderedMap[K, V]) isEntryActive(i int) bool {
	if m.holes == 0 {
		return true
	}
	k := m.entries[i].Key
	return m.idx[k] == int32(i+1) || k != k
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

func NewOrderedMapCap[K comparable, V any](cap int32) OrderedMap[K, V] {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	var m OrderedMap[K, V]
	if cap >= 64 {
		m.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(cap))
	} else if cap > 8 {
		m.entries = make([]OrderedMapEntry[K, V], 0, cap)
	}
	if int(cap) > orderedMapIndexThreshold {
		m.idx = make(map[K]int32, cap)
	}
	return m
}

func NewOrderedMapPtrCap[K comparable, V any](cap int32) *OrderedMap[K, V] {
	m := NewOrderedMapCap[K, V](cap)
	return &m
}

func OrderedMapFromPtr[K comparable, V any](src *OrderedMap[K, V]) OrderedMap[K, V] {
	if src == nil {
		ThrowNPE()
	}
	live := int32(len(src.entries)) - src.holes
	if live == 0 {
		return OrderedMap[K, V]{}
	}
	capHint := int(live)
	if capHint < 4 {
		capHint = 4
	}
	var es []OrderedMapEntry[K, V]
	if capHint >= 64 {
		es = AcquirePooledSlice[OrderedMapEntry[K, V]](capHint)[:live]
	} else {
		es = make([]OrderedMapEntry[K, V], live, capHint)
	}
	if src.holes == 0 {
		copy(es, src.entries)
	} else {
		w := 0
		for i := int(src.firstIdx); i < len(src.entries); i++ {
			if src.isEntryActive(i) {
				es[w] = src.entries[i]
				w++
			}
		}
	}
	dst := OrderedMap[K, V]{entries: es}
	if len(es) > orderedMapIndexThreshold {
		idx := make(map[K]int32, max(16, len(es)*2))
		for i := range es {
			idx[es[i].Key] = int32(i + 1)
		}
		dst.idx = idx
	}
	return dst
}

func OrderedMapFrom[K comparable, V any](src OrderedMap[K, V]) OrderedMap[K, V] {
	return OrderedMapFromPtr(&src)
}

func OrderedMapPtrFrom[K comparable, V any](src *OrderedMap[K, V]) *OrderedMap[K, V] {
	dst := OrderedMapFromPtr(src)
	return &dst
}

func OrderedMapFromJavaMap[K comparable, V any](src any) OrderedMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewOrderedMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.Put(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func OrderedMapPtrFromJavaMap[K comparable, V any](src any) *OrderedMap[K, V] {
	m := OrderedMapFromJavaMap[K, V](src)
	return &m
}

func OrderedMapFromNodeKeyJavaMap[K ~uint32, V any](src any) OrderedMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewOrderedMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.Put(K(cur.KeyNode()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func OrderedMapPtrFromNodeKeyJavaMap[K ~uint32, V any](src any) *OrderedMap[K, V] {
	m := OrderedMapFromNodeKeyJavaMap[K, V](src)
	return &m
}

func (m *OrderedMap[K, V]) IsNil() bool      { return m == nil }
func (m *OrderedMap[K, V]) IsJavaNil_() bool { return m == nil }

func (m *OrderedMap[K, V]) Size() int32 {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.entries)) - m.holes
}

func (m *OrderedMap[K, V]) IsEmpty() bool {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.entries)) == m.holes
}

func (m *OrderedMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zero V
	if int32(len(m.entries)) == m.holes {
		return zero, false
	}
	if m.idx != nil {
		pos1 := m.idx[k]
		if pos1 == 0 {
			return zero, false
		}
		return m.entries[pos1-1].Val, true
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return m.entries[i].Val, true
		}
	}
	return zero, false
}

func (m *OrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}

func (m *OrderedMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return def
}

func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}

func (m *OrderedMap[K, V]) appendEntry(e OrderedMapEntry[K, V]) {
	es := m.entries
	if len(es) < cap(es) {
		m.entries = append(es, e)
		return
	}
	m.appendEntryGrow(e)
}

func (m *OrderedMap[K, V]) appendEntryGrow(e OrderedMapEntry[K, V]) {
	old := m.entries
	oldLen := len(old)
	oldCap := cap(old)
	newCap := oldCap * 2
	if newCap < 4 {
		newCap = 4
	}
	if newCap <= oldLen {
		newCap = oldLen + 1
	}
	if newCap < 64 {
		next := make([]OrderedMapEntry[K, V], oldLen+1, newCap)
		copy(next, old)
		next[oldLen] = e
		m.entries = next
		return
	}
	next := AcquirePooledSlice[OrderedMapEntry[K, V]](newCap)
	next = next[:oldLen+1]
	copy(next, old)
	next[oldLen] = e
	if oldCap >= 64 {
		ReleasePooledSlice(old)
	}
	m.entries = next
}

func (m *OrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zero V
	if m.idx != nil {
		if pos1 := m.idx[k]; pos1 != 0 {
			old := m.entries[pos1-1].Val
			m.entries[pos1-1].Val = v
			return old, true
		}
		m.appendEntry(OrderedMapEntry[K, V]{Key: k, Val: v})
		m.idx[k] = int32(len(m.entries))
		return zero, false
	}
	if m.entries == nil {
		es := make([]OrderedMapEntry[K, V], 1, 4)
		es[0] = OrderedMapEntry[K, V]{Key: k, Val: v}
		m.entries = es
		return zero, false
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			old := m.entries[i].Val
			m.entries[i].Val = v
			return old, true
		}
	}
	m.appendEntry(OrderedMapEntry[K, V]{Key: k, Val: v})
	if len(m.entries) > orderedMapIndexThreshold {
		idx := make(map[K]int32, max(24, cap(m.entries)))
		for i := range m.entries {
			idx[m.entries[i].Key] = int32(i + 1)
		}
		m.idx = idx
	}
	return zero, false
}

func (m *OrderedMap[K, V]) Put(k K, v V) V {
	old, _ := m.PutOk(k, v)
	return old
}

func (m *OrderedMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, ok := m.GetOk(k)
	if !ok || isJavaNullMapVal(cur) {
		m.PutOk(k, v)
	}
	return cur
}

func (m *OrderedMap[K, V]) RemoveOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zero V
	if int32(len(m.entries)) == m.holes {
		return zero, false
	}
	if m.idx == nil {
		for i := range m.entries {
			if m.entries[i].Key == k {
				old := m.entries[i].Val
				copy(m.entries[i:], m.entries[i+1:])
				last := len(m.entries) - 1
				m.entries[last] = OrderedMapEntry[K, V]{}
				m.entries = m.entries[:last]
				return old, true
			}
		}
		return zero, false
	}
	pos1 := m.idx[k]
	if pos1 == 0 {
		return zero, false
	}
	delete(m.idx, k)
	slot := pos1 - 1
	old := m.entries[slot].Val
	m.entries[slot] = OrderedMapEntry[K, V]{}
	m.holes++
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		m.entries = m.entries[:0]
		if m.idx != nil {
			clear(m.idx)
		}
		m.holes = 0
		m.firstIdx = 0
	} else {
		if slot == m.firstIdx {
			for int(m.firstIdx) < len(m.entries) && !m.isEntryActive(int(m.firstIdx)) {
				m.firstIdx++
			}
		}
		if len(m.entries) > 16 && int(live)*2 < len(m.entries) {
			m.compact()
		}
	}
	return old, true
}

func (m *OrderedMap[K, V]) Remove(k K) V {
	old, _ := m.RemoveOk(k)
	return old
}

func (m *OrderedMap[K, V]) compact() {
	w := 0
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		k := m.entries[i].Key
		if m.idx[k] == int32(i+1) || k != k {
			m.entries[w] = m.entries[i]
			if k == k {
				m.idx[k] = int32(w + 1)
			}
			w++
		}
	}
	clear(m.entries[w:])
	m.entries = m.entries[:w]
	m.holes = 0
	m.firstIdx = 0
}

func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		ThrowNPE()
	}
	n := len(m.entries)
	if n == 0 {
		if len(m.idx) != 0 {
			clear(m.idx)
		}
		m.holes = 0
		m.firstIdx = 0
		return
	}
	if m.idx != nil && len(m.idx) != 0 {
		live := n - int(m.holes)
		if live > 0 && live <= 32 && n-int(m.firstIdx) <= 64 && cap(m.entries) >= 128 {
			for i := int(m.firstIdx); i < n; i++ {
				delete(m.idx, m.entries[i].Key)
			}
			if len(m.idx) != 0 {
				clear(m.idx)
			}
		} else {
			clear(m.idx)
		}
	}
	clear(m.entries)
	m.entries = m.entries[:0]
	m.holes = 0
	m.firstIdx = 0
}

func (m *OrderedMap[K, V]) PutAllMapPtr(other *OrderedMap[K, V]) {
	if m == nil || other == nil {
		ThrowNPE()
	}
	if int32(len(other.entries)) == other.holes {
		return
	}
	if cap(m.entries) == 0 && m.idx == nil {
		*m = OrderedMapFromPtr(other)
		return
	}
	for i := int(other.firstIdx); i < len(other.entries); i++ {
		if other.isEntryActive(i) {
			m.PutOk(other.entries[i].Key, other.entries[i].Val)
		}
	}
}

func (m *OrderedMap[K, V]) PutAllMap(other OrderedMap[K, V]) {
	m.PutAllMapPtr(&other)
}

func (m *OrderedMap[K, V]) PutAllJavaMap(src any) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
}

func OrderedMapPtrGetStr[K comparable](m *OrderedMap[K, string], k K) string {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return NullString
}

func OrderedMapGetStr[K comparable](m OrderedMap[K, string], k K) string {
	return OrderedMapPtrGetStr(&m, k)
}

func OrderedMapPutStr[K comparable](m *OrderedMap[K, string], k K, v string) string {
	if old, ok := m.PutOk(k, v); ok {
		return old
	}
	return NullString
}

func OrderedMapRemoveStr[K comparable](m *OrderedMap[K, string], k K) string {
	if old, ok := m.RemoveOk(k); ok {
		return old
	}
	return NullString
}

func OrderedMapPutIfAbsentStr[K comparable](m *OrderedMap[K, string], k K, v string) string {
	cur, ok := m.GetOk(k)
	if !ok || cur == NullString {
		m.PutOk(k, v)
		return NullString
	}
	return cur
}

func OrderedMapComputeIfAbsent[K comparable, V comparable](m *OrderedMap[K, V], k K, nullV V, fn func(K) V) V {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	nv := fn(k)
	if nv != nullV {
		m.PutOk(k, nv)
	}
	return nv
}

func OrderedMapComputeIfAbsentFn[K comparable, V comparable](m *OrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(fn)
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	raw := fn.Apply(boxJavaElem(k))
	if raw == nil || isJavaNilAny(raw) {
		return nullV
	}
	nv := unboxJavaElem[V](raw)
	if nv != nullV {
		m.PutOk(k, nv)
	}
	return nv
}

type OrderedMapKeyCursor[K comparable, V any] struct {
	m   *OrderedMap[K, V]
	pos int32
}

func (m *OrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return OrderedMapKeyCursor[K, V]{m: m}
	}
	return OrderedMapKeyCursor[K, V]{m: m, pos: m.firstIdx}
}

func (c *OrderedMapKeyCursor[K, V]) HasNext() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.isEntryActive(int(c.pos)) {
			return true
		}
		c.pos++
	}
	return false
}

func (c *OrderedMapKeyCursor[K, V]) Next() K {
	m := c.m
	i := c.pos
	c.pos = i + 1
	return m.entries[i].Key
}

type OrderedMapValCursor[K comparable, V any] struct {
	m   *OrderedMap[K, V]
	pos int32
}

func (m *OrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return OrderedMapValCursor[K, V]{m: m}
	}
	return OrderedMapValCursor[K, V]{m: m, pos: m.firstIdx}
}

func (c *OrderedMapValCursor[K, V]) HasNext() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.isEntryActive(int(c.pos)) {
			return true
		}
		c.pos++
	}
	return false
}

func (c *OrderedMapValCursor[K, V]) Next() V {
	m := c.m
	i := c.pos
	c.pos = i + 1
	return m.entries[i].Val
}

type OrderedMapEntryCursor[K comparable, V any] struct {
	m        *OrderedMap[K, V]
	pos      int32
	lastSlot int32
}

func (m *OrderedMap[K, V]) IterEntries() OrderedMapEntryCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return OrderedMapEntryCursor[K, V]{m: m, lastSlot: -1}
	}
	return OrderedMapEntryCursor[K, V]{m: m, pos: m.firstIdx, lastSlot: -1}
}

func (c *OrderedMapEntryCursor[K, V]) NextEntry() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.isEntryActive(int(c.pos)) {
			c.lastSlot = c.pos
			c.pos++
			return true
		}
		c.pos++
	}
	return false
}

func (c *OrderedMapEntryCursor[K, V]) Key() K {
	return c.m.entries[c.lastSlot].Key
}

func (c *OrderedMapEntryCursor[K, V]) Value() V {
	return c.m.entries[c.lastSlot].Val
}

func (c *OrderedMapEntryCursor[K, V]) SetValue(nv V) V {
	old := c.m.entries[c.lastSlot].Val
	c.m.entries[c.lastSlot].Val = nv
	return old
}

func (m *OrderedMap[K, V]) ForEach(fn func(K, V)) {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.isEntryActive(i) {
			fn(m.entries[i].Key, m.entries[i].Val)
		}
	}
}

func (m *OrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(action)
	m.ForEach(func(k K, v V) {
		action.Accept(boxJavaElem(k), boxJavaElem(v))
	})
}

func OrderedMapPtrEquals[K comparable, V comparable](a, b *OrderedMap[K, V], identityVal bool) bool {
	if a == nil {
		ThrowNPE()
	}
	if b == nil {
		return false
	}
	aLive := int32(len(a.entries)) - a.holes
	bLive := int32(len(b.entries)) - b.holes
	if aLive != bLive {
		return false
	}
	if aLive == 0 {
		return true
	}
	for i := int(a.firstIdx); i < len(a.entries); i++ {
		if a.isEntryActive(i) {
			v2, ok := b.GetOk(a.entries[i].Key)
			if !ok || (a.entries[i].Val != v2 && (identityVal || !javaMapValsEqual(a.entries[i].Val, v2, false))) {
				return false
			}
		}
	}
	return true
}

func OrderedMapEquals[K comparable, V comparable](a, b OrderedMap[K, V], identityVal bool) bool {
	return OrderedMapPtrEquals(&a, &b, identityVal)
}

func (m *OrderedMap[K, V]) EqualsMapPtr(other *OrderedMap[K, V], identityVal bool) bool {
	if m == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	mLive := int32(len(m.entries)) - m.holes
	oLive := int32(len(other.entries)) - other.holes
	if mLive != oLive {
		return false
	}
	if mLive == 0 {
		return true
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.isEntryActive(i) {
			v2, ok := other.GetOk(m.entries[i].Key)
			if !ok || !javaMapValsEqual(m.entries[i].Val, v2, identityVal) {
				return false
			}
		}
	}
	return true
}

func (m *OrderedMap[K, V]) EqualsMap(other OrderedMap[K, V], identityVal bool) bool {
	return m.EqualsMapPtr(&other, identityVal)
}

func (m *OrderedMap[K, V]) HashCodeMap(identity bool) int32 {
	if m == nil {
		ThrowNPE()
	}
	var h int32
	m.ForEach(func(k K, v V) {
		h += javaMapEntryHash(k, v, identity)
	})
	return h
}

func (m *OrderedMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	if int32(len(m.entries)) == m.holes {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	m.ForEach(func(k K, v V) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(boxJavaElem(k)))
		sb.WriteByte('=')
		sb.WriteString(StrVal(boxJavaElem(v)))
	})
	sb.WriteByte('}')
	return sb.String()
}

func (m *OrderedMap[K, V]) ToImmutableMap() Ext_com_google_common_collect_ImmutableMap {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableMap_
	}
	obj := newImmutableMapShell()
	mc := &obj.inlineCore
	mc.InitExpectedSize(live)
	m.ForEach(func(k K, v V) {
		mc.Put(CheckNotNull(boxJavaElem(k)), CheckNotNull(boxJavaElem(v)))
	})
	mc.sealImmutable()
	return obj
}

func (m *OrderedMap[K, V]) KeysToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableSet_
	}
	s := NewOrderedSetCap[K](live)
	m.ForEach(func(k K, _ V) {
		s.Add(k)
	})
	return s.ToImmutableSet()
}

func (m *OrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableList_
	}
	items := make([]V, 0, live)
	m.ForEach(func(_ K, v V) {
		items = append(items, v)
	})
	return ImmListFromSlice(items)
}

func (m *OrderedMap[K, V]) ToJavaLinkedHashMap() *Ext_java_util_LinkedHashMap {
	if m == nil {
		return nil
	}
	live := int32(len(m.entries)) - m.holes
	obj := &Ext_java_util_LinkedHashMap{}
	obj.BindExtSelf_(obj)
	mc := &obj.inlineCore
	mc.InitExpectedSize(live)
	m.ForEach(func(k K, v V) {
		mc.Put(boxJavaElem(k), boxJavaElem(v))
	})
	return obj
}

func (m *OrderedMap[K, V]) ToJavaMap(identity bool) Ext_java_util_Map {
	if m == nil {
		return nil
	}
	if identity {
		live := int32(len(m.entries)) - m.holes
		obj := NewExt_java_util_IdentityHashMap()
		mc := &obj.inlineCore
		mc.InitExpectedSize(live)
		m.ForEach(func(k K, v V) {
			mc.Put(boxJavaElem(k), boxJavaElem(v))
		})
		return obj
	}
	return m.ToJavaLinkedHashMap()
}

type jOrdMapLiveValues[K comparable, V any] struct {
	Ext_java_util_AbstractCollectionBase
	m *OrderedMap[K, V]
}

func (v *jOrdMapLiveValues[K, V]) IsJavaNil_() bool { return v == nil }
func (v *jOrdMapLiveValues[K, V]) Size() int32 {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	return v.m.Size()
}
func (v *jOrdMapLiveValues[K, V]) IsEmpty() bool {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	return v.m.IsEmpty()
}
func (v *jOrdMapLiveValues[K, V]) Clear() {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	v.m.Clear()
}
func (v *jOrdMapLiveValues[K, V]) Contains(x any) bool {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	if int32(len(v.m.entries)) == v.m.holes {
		return false
	}
	target := ToAny(x)
	for cur := v.m.IterValues(); cur.HasNext(); {
		if javalang.ObjectsEquals(boxJavaElem(cur.Next()), target) {
			return true
		}
	}
	return false
}
func ordMapNextActive[K comparable, V any](m *OrderedMap[K, V], from int32) int32 {
	n := int32(len(m.entries))
	if n == m.holes || from < 0 || from >= n {
		return -1
	}
	if m.holes == 0 {
		return from
	}
	for i := int(from); i < int(n); i++ {
		if m.isEntryActive(i) {
			return int32(i)
		}
	}
	return -1
}

func (v *jOrdMapLiveValues[K, V]) iterInit_() (int32, bool) {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	p := ordMapNextActive(v.m, v.m.firstIdx)
	return p, p >= 0
}
func (v *jOrdMapLiveValues[K, V]) iterNext_(pos int32) (any, int32) {
	p := ordMapNextActive(v.m, pos)
	if p < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	return boxJavaElem(v.m.entries[p].Val), ordMapNextActive(v.m, p+1)
}
func (v *jOrdMapLiveValues[K, V]) Iterator() any {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	if int32(len(v.m.entries)) == v.m.holes {
		return emptyMapIter_
	}
	return &jOrdMapValIter[K, V]{m: v.m, pos: v.m.firstIdx, lastKeyValid: false}
}

type jOrdMapValIter[K comparable, V any] struct {
	m            *OrderedMap[K, V]
	pos          int32
	lastKey      K
	lastKeyValid bool
}

func (*jOrdMapValIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jOrdMapValIter[K, V]) IsJavaNil_() bool                       { return it == nil }
func (it *jOrdMapValIter[K, V]) HasNext() bool {
	if int32(len(it.m.entries)) == it.m.holes {
		return false
	}
	for int(it.pos) < len(it.m.entries) {
		if it.m.isEntryActive(int(it.pos)) {
			return true
		}
		it.pos++
	}
	return false
}
func (it *jOrdMapValIter[K, V]) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	i := it.pos
	it.pos = i + 1
	it.lastKey = it.m.entries[i].Key
	it.lastKeyValid = true
	return boxJavaElem(it.m.entries[i].Val)
}
func (it *jOrdMapValIter[K, V]) Remove() {
	if !it.lastKeyValid {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	wasShifted := it.m.idx == nil
	it.m.RemoveOk(it.lastKey)
	if wasShifted && it.pos > 0 {
		it.pos--
	}
	it.lastKeyValid = false
}

func (m *OrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
	if m == nil {
		ThrowNPE()
	}
	v := &jOrdMapLiveValues[K, V]{m: m}
	v.BindExtSelf_(v)
	return v
}

type jOrdMapLiveKeySet[K comparable, V any] struct {
	Ext_java_util_AbstractSetBase
	m *OrderedMap[K, V]
}

func (s *jOrdMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
func (s *jOrdMapLiveKeySet[K, V]) Size() int32 {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.Size()
}
func (s *jOrdMapLiveKeySet[K, V]) IsEmpty() bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.IsEmpty()
}
func (s *jOrdMapLiveKeySet[K, V]) Clear() {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	s.m.Clear()
}
func (s *jOrdMapLiveKeySet[K, V]) Contains(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	k, ok := tryUnboxOrderedMapKey[K](x)
	if !ok {
		return false
	}
	return s.m.ContainsKey(k)
}
func (s *jOrdMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	k, ok := tryUnboxOrderedMapKey[K](x)
	if !ok {
		return false
	}
	_, removed := s.m.RemoveOk(k)
	return removed
}
func (s *jOrdMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if k == NullString {
			return s.m.ContainsKey(zero)
		}
		uk, ok := any(k).(K)
		return ok && s.m.ContainsKey(uk)
	}
	if _, isStr := any(zero).(string); isStr {
		return s.m.ContainsKey(any(k).(K))
	}
	if k == NullString {
		if uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
			return s.m.ContainsKey(uk)
		}
	}
	return false
}
func (s *jOrdMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jOrdMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if k == NullString {
			_, removed := s.m.RemoveOk(zero)
			return removed
		}
		uk, ok := any(k).(K)
		if !ok {
			return false
		}
		_, removed := s.m.RemoveOk(uk)
		return removed
	}
	if _, isStr := any(zero).(string); isStr {
		_, removed := s.m.RemoveOk(any(k).(K))
		return removed
	}
	if k == NullString {
		if uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
			_, removed := s.m.RemoveOk(uk)
			return removed
		}
	}
	return false
}
func (s *jOrdMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	if u == 0 {
		uk, ok := tryUnboxOrderedMapKey[K](nil)
		return ok && s.m.ContainsKey(uk)
	}
	var zero K
	if any(zero) == nil {
		uk, ok := boxNodeById(u).(K)
		return ok && s.m.ContainsKey(uk)
	}
	switch any(zero).(type) {
	case uint32:
		return s.m.ContainsKey(any(u).(K))
	case nodestore.Node:
		return s.m.ContainsKey(any(nodestore.Node(u)).(K))
	case jRawNodeProvider:
		if globalNodeBoxer != nil {
			if uk, ok := globalNodeBoxer(u).(K); ok {
				return s.m.ContainsKey(uk)
			}
		}
		uk, ok := tryUnboxOrderedMapKey[K](u)
		return ok && s.m.ContainsKey(uk)
	default:
		return false
	}
}
func (s *jOrdMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jOrdMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	if u == 0 {
		uk, ok := tryUnboxOrderedMapKey[K](nil)
		if !ok {
			return false
		}
		_, removed := s.m.RemoveOk(uk)
		return removed
	}
	var zero K
	if any(zero) == nil {
		uk, ok := boxNodeById(u).(K)
		if !ok {
			return false
		}
		_, removed := s.m.RemoveOk(uk)
		return removed
	}
	switch any(zero).(type) {
	case uint32:
		_, removed := s.m.RemoveOk(any(u).(K))
		return removed
	case nodestore.Node:
		_, removed := s.m.RemoveOk(any(nodestore.Node(u)).(K))
		return removed
	case jRawNodeProvider:
		if globalNodeBoxer != nil {
			if uk, ok := globalNodeBoxer(u).(K); ok {
				_, removed := s.m.RemoveOk(uk)
				return removed
			}
		}
		uk, ok := tryUnboxOrderedMapKey[K](u)
		if !ok {
			return false
		}
		_, removed := s.m.RemoveOk(uk)
		return removed
	default:
		return false
	}
}
func (s *jOrdMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	p := ordMapNextActive(s.m, s.m.firstIdx)
	return p, p >= 0
}
func (s *jOrdMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
	p := ordMapNextActive(s.m, pos)
	if p < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	return boxJavaElem(s.m.entries[p].Key), ordMapNextActive(s.m, p+1)
}
func (s *jOrdMapLiveKeySet[K, V]) Iterator() any {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return emptyMapIter_
	}
	return &jOrdMapKeyIter[K, V]{m: s.m, pos: s.m.firstIdx}
}

type jOrdMapKeyIter[K comparable, V any] struct {
	m            *OrderedMap[K, V]
	pos          int32
	lastKey      K
	lastKeyValid bool
}

func (*jOrdMapKeyIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jOrdMapKeyIter[K, V]) IsJavaNil_() bool                       { return it == nil }
func (it *jOrdMapKeyIter[K, V]) HasNext() bool {
	if int32(len(it.m.entries)) == it.m.holes {
		return false
	}
	for int(it.pos) < len(it.m.entries) {
		if it.m.isEntryActive(int(it.pos)) {
			return true
		}
		it.pos++
	}
	return false
}
func (it *jOrdMapKeyIter[K, V]) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	i := it.pos
	it.pos = i + 1
	it.lastKey = it.m.entries[i].Key
	it.lastKeyValid = true
	return boxJavaElem(it.lastKey)
}
func (it *jOrdMapKeyIter[K, V]) Remove() {
	if !it.lastKeyValid {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	wasShifted := it.m.idx == nil
	it.m.RemoveOk(it.lastKey)
	if wasShifted && it.pos > 0 {
		it.pos--
	}
	it.lastKeyValid = false
}

func (m *OrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
	if m == nil {
		ThrowNPE()
	}
	s := &jOrdMapLiveKeySet[K, V]{m: m}
	s.BindExtSelf_(s)
	return s
}

func (m *OrderedMap[K, V]) FirstValue() V {
	if m == nil {
		ThrowNPE()
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.isEntryActive(i) {
			return m.entries[i].Val
		}
	}
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}

func OrderedMapGetNodeAny[K ~uint32, V any](m *OrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	var k K
	if rawKey != nil && !isJavaNilAny(rawKey) {
		if nk, ok := rawKey.(K); ok {
			k = nk
		} else if rn, ok := rawKey.(jRawNodeProvider); ok {
			k = K(rn.RawNodeId_())
		} else if nd, ok := rawKey.(nodestore.Node); ok {
			k = K(nd)
		} else if u, ok := rawKey.(uint32); ok {
			k = K(u)
		} else {
			return nil
		}
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

var boxedNullStringAny_ any = NullString

func tryUnboxOrderedMapKey[K comparable](raw any) (K, bool) {
	var zero K
	if any(zero) == nil {
		if raw == nil || isJavaNilAny(raw) {
			return zero, true
		}
		if !reflect.TypeOf(raw).Comparable() {
			return zero, false
		}
		k, ok := raw.(K)
		return k, ok
	}
	if raw != nil {
		if k, ok := raw.(K); ok {
			return k, true
		}
	}
	if raw == nil || isJavaNilAny(raw) {
		switch any(zero).(type) {
		case string:
			return boxedNullStringAny_.(K), true
		case int32, int64, bool, float64, float32, uint16, int8, int16:
			return zero, false
		default:
			return zero, true
		}
	}
	switch any(zero).(type) {
	case string:
		return zero, false
	case int32:
		if b, ok := raw.(*BoxInt); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case int64:
		if b, ok := raw.(*BoxLong); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case bool:
		if b, ok := raw.(*BoxBoolean); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case float64:
		if b, ok := raw.(*BoxDouble); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case float32:
		if b, ok := raw.(*BoxFloat); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case uint16:
		if b, ok := raw.(*BoxChar); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case int8:
		if b, ok := raw.(*BoxByte); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case int16:
		if b, ok := raw.(*BoxShort); ok && b != nil {
			return any(b.Val).(K), true
		}
		return zero, false
	case uint32:
		if rn, ok := raw.(jRawNodeProvider); ok {
			return any(rn.RawNodeId_()).(K), true
		}
		if nd, ok := raw.(nodestore.Node); ok {
			return any(uint32(nd)).(K), true
		}
		return zero, false
	case nodestore.Node:
		if rn, ok := raw.(jRawNodeProvider); ok {
			return any(nodestore.Node(rn.RawNodeId_())).(K), true
		}
		if u, ok := raw.(uint32); ok {
			return any(nodestore.Node(u)).(K), true
		}
		return zero, false
	case jRawNodeProvider:
		var u uint32
		var haveU bool
		if rn, ok := raw.(jRawNodeProvider); ok {
			u = rn.RawNodeId_()
			haveU = true
		} else if nd, ok := raw.(nodestore.Node); ok {
			u = uint32(nd)
			haveU = true
		} else if rawU, ok := raw.(uint32); ok {
			u = rawU
			haveU = true
		}
		if !haveU {
			return zero, false
		}
		if u == 0 {
			return zero, true
		}
		if globalNodeBoxer != nil {
			if b, ok := globalNodeBoxer(u).(K); ok {
				return b, true
			}
		}
		rt := reflect.TypeFor[K]()
		if rt.Kind() == reflect.Uint32 {
			return reflect.ValueOf(u).Convert(rt).Interface().(K), true
		}
		return zero, false
	default:
		return zero, false
	}
}

func OrderedMapGetAny[K comparable, V any](m *OrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	k, ok := tryUnboxOrderedMapKey[K](rawKey)
	if !ok {
		return nil
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

// PlainMap is an un-ordered value struct backed solely by a native Go map[K]V (zero entries slice)
// for non-iterated HashMap lowering.
type PlainMap[K comparable, V any] struct {
	m map[K]V
}

func NewPlainMap[K comparable, V any](capHint int) *PlainMap[K, V] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 {
		return &PlainMap[K, V]{m: make(map[K]V, capHint)}
	}
	return &PlainMap[K, V]{}
}

func NewPlainMapVal[K comparable, V any]() PlainMap[K, V] {
	return PlainMap[K, V]{}
}

func NewPlainMapCap[K comparable, V any](capHint int32) PlainMap[K, V] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 {
		return PlainMap[K, V]{m: make(map[K]V, capHint)}
	}
	return PlainMap[K, V]{}
}

func NewPlainMapPtr[K comparable, V any]() *PlainMap[K, V] {
	return &PlainMap[K, V]{}
}

func NewPlainMapPtrCap[K comparable, V any](capHint int32) *PlainMap[K, V] {
	return NewPlainMap[K, V](int(capHint))
}

func (m *PlainMap[K, V]) Init() {}

func (m *PlainMap[K, V]) InitCap(capHint int) {
	if m == nil {
		ThrowNPE()
	}
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 && m.m == nil {
		m.m = make(map[K]V, capHint)
	}
}

func (m *PlainMap[K, V]) IsNil() bool      { return m == nil }
func (m *PlainMap[K, V]) IsJavaNil_() bool { return m == nil }

func (m *PlainMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	return int32(len(m.m))
}

func (m *PlainMap[K, V]) IsEmpty() bool {
	return m == nil || len(m.m) == 0
}

func (m *PlainMap[K, V]) ContainsKey(k K) bool {
	if m == nil || len(m.m) == 0 {
		return false
	}
	_, ok := m.m[k]
	return ok
}

func (m *PlainMap[K, V]) Get(k K) V {
	if m == nil || len(m.m) == 0 {
		var zero V
		return zero
	}
	return m.m[k]
}

func (m *PlainMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil || len(m.m) == 0 {
		var zero V
		return zero, false
	}
	v, ok := m.m[k]
	return v, ok
}

func (m *PlainMap[K, V]) GetOrDefault(k K, def V) V {
	if m == nil || len(m.m) == 0 {
		return def
	}
	if v, ok := m.m[k]; ok {
		return v
	}
	return def
}

func (m *PlainMap[K, V]) Put(k K, v V) V {
	if m == nil {
		ThrowNPE()
	}
	if m.m == nil {
		m.m = make(map[K]V)
		m.m[k] = v
		var zero V
		return zero
	}
	old := m.m[k]
	m.m[k] = v
	return old
}

func (m *PlainMap[K, V]) PutOk(k K, v V) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	if m.m == nil {
		m.m = make(map[K]V)
		m.m[k] = v
		var zero V
		return zero, false
	}
	old, existed := m.m[k]
	m.m[k] = v
	return old, existed
}

func (m *PlainMap[K, V]) PutVoid(k K, v V) {
	if m == nil {
		ThrowNPE()
	}
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}

func (m *PlainMap[K, V]) PutIfAbsentOk(k K, v V) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	if m.m == nil {
		m.m = make(map[K]V)
		m.m[k] = v
		var zero V
		return zero, false
	}
	cur, ok := m.m[k]
	if !ok || isJavaNullMapVal(cur) {
		m.m[k] = v
		return cur, false
	}
	return cur, true
}

func (m *PlainMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, _ := m.PutIfAbsentOk(k, v)
	return cur
}

func (m *PlainMap[K, V]) RemoveOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	if len(m.m) == 0 {
		var zero V
		return zero, false
	}
	old, existed := m.m[k]
	if existed {
		delete(m.m, k)
	}
	return old, existed
}

func (m *PlainMap[K, V]) Remove(k K) V {
	old, _ := m.RemoveOk(k)
	return old
}

func (m *PlainMap[K, V]) Clear() {
	if m != nil && len(m.m) > 0 {
		clear(m.m)
	}
}

func PlainMapComputeIfAbsentFn[K comparable, V any](m *PlainMap[K, V], k K, fn func(K) V) V {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if cur, ok := m.GetOk(k); ok && !isJavaNullMapVal(cur) {
		return cur
	}
	nv := fn(k)
	if !isJavaNullMapVal(nv) {
		m.PutVoid(k, nv)
	}
	return nv
}

func PlainMapComputeIfAbsent[K comparable, V comparable](m *PlainMap[K, V], k K, nullV V, fn func(K) V) V {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	nv := fn(k)
	if nv != nullV {
		m.PutVoid(k, nv)
	}
	return nv
}

func PlainMapComputeIfAbsentExtFn[K comparable, V comparable](m *PlainMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(fn)
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	raw := fn.Apply(boxJavaElem(k))
	if raw == nil || isJavaNilAny(raw) {
		return nullV
	}
	nv := unboxJavaElem[V](raw)
	if nv != nullV {
		m.PutVoid(k, nv)
	}
	return nv
}

func PlainMapPtrGetStr[K comparable](m *PlainMap[K, string], k K) string {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return NullString
}

func PlainMapGetStr[K comparable](m PlainMap[K, string], k K) string {
	return PlainMapPtrGetStr(&m, k)
}

func PlainMapPutStr[K comparable](m *PlainMap[K, string], k K, v string) string {
	if old, ok := m.PutOk(k, v); ok {
		return old
	}
	return NullString
}

func PlainMapRemoveStr[K comparable](m *PlainMap[K, string], k K) string {
	if old, ok := m.RemoveOk(k); ok {
		return old
	}
	return NullString
}

func PlainMapPutIfAbsentStr[K comparable](m *PlainMap[K, string], k K, v string) string {
	cur, ok := m.GetOk(k)
	if !ok || cur == NullString {
		m.PutVoid(k, v)
		return NullString
	}
	return cur
}

func PlainMapFromJavaMap[K comparable, V any](src any) PlainMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewPlainMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutVoid(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func PlainMapPtrFromJavaMap[K comparable, V any](src any) *PlainMap[K, V] {
	m := PlainMapFromJavaMap[K, V](src)
	return &m
}

func PlainMapFromNodeKeyJavaMap[K ~uint32, V any](src any) PlainMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewPlainMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutVoid(K(cur.KeyNode()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func PlainMapPtrFromNodeKeyJavaMap[K ~uint32, V any](src any) *PlainMap[K, V] {
	m := PlainMapFromNodeKeyJavaMap[K, V](src)
	return &m
}

func (m *PlainMap[K, V]) PutAllJavaMap(src any) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutVoid(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
}

func PlainMapGetNodeAny[K ~uint32, V any](m *PlainMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	var k K
	if rawKey != nil && !isJavaNilAny(rawKey) {
		if nk, ok := rawKey.(K); ok {
			k = nk
		} else if rn, ok := rawKey.(jRawNodeProvider); ok {
			k = K(rn.RawNodeId_())
		} else if nd, ok := rawKey.(nodestore.Node); ok {
			k = K(nd)
		} else if u, ok := rawKey.(uint32); ok {
			k = K(u)
		} else {
			return nil
		}
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

func PlainMapGetAny[K comparable, V any](m *PlainMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	k, ok := tryUnboxOrderedMapKey[K](rawKey)
	if !ok {
		return nil
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

// PlainSet is an un-ordered value struct backed solely by a native Go map[T]struct{}
// for non-iterated HashSet lowering.
type PlainSet[T comparable] struct {
	m map[T]struct{}
}

func NewPlainSet[T comparable](capHint int) *PlainSet[T] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 {
		return &PlainSet[T]{m: make(map[T]struct{}, capHint)}
	}
	return &PlainSet[T]{}
}

func NewPlainSetVal[T comparable]() PlainSet[T] {
	return PlainSet[T]{}
}

func NewPlainSetCap[T comparable](capHint int32) PlainSet[T] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 {
		return PlainSet[T]{m: make(map[T]struct{}, capHint)}
	}
	return PlainSet[T]{}
}

func NewPlainSetPtr[T comparable]() *PlainSet[T] {
	return &PlainSet[T]{}
}

func NewPlainSetPtrCap[T comparable](capHint int32) *PlainSet[T] {
	return NewPlainSet[T](int(capHint))
}

func (s *PlainSet[T]) Init() {}

func (s *PlainSet[T]) InitCap(capHint int) {
	if s == nil {
		ThrowNPE()
	}
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint > 0 && s.m == nil {
		s.m = make(map[T]struct{}, capHint)
	}
}

func (s *PlainSet[T]) IsNil() bool      { return s == nil }
func (s *PlainSet[T]) IsJavaNil_() bool { return s == nil }

func (s *PlainSet[T]) Size() int32 {
	if s == nil {
		return 0
	}
	return int32(len(s.m))
}

func (s *PlainSet[T]) IsEmpty() bool {
	return s == nil || len(s.m) == 0
}

func (s *PlainSet[T]) Contains(v T) bool {
	if s == nil || len(s.m) == 0 {
		return false
	}
	_, ok := s.m[v]
	return ok
}

func (s *PlainSet[T]) Add(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	if s.m == nil {
		s.m = make(map[T]struct{})
		s.m[v] = struct{}{}
		return true
	}
	if _, ok := s.m[v]; ok {
		return false
	}
	s.m[v] = struct{}{}
	return true
}

func (s *PlainSet[T]) AddVoid(v T) {
	if s == nil {
		ThrowNPE()
	}
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[v] = struct{}{}
}

func (s *PlainSet[T]) Remove(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	if len(s.m) == 0 {
		return false
	}
	if _, ok := s.m[v]; ok {
		delete(s.m, v)
		return true
	}
	return false
}

func (s *PlainSet[T]) Clear() {
	if s != nil && len(s.m) > 0 {
		clear(s.m)
	}
}

func PlainSetFromStrCol(col any) PlainSet[string] {
	var s PlainSet[string]
	PlainSetAddAllStr(&s, col)
	return s
}

func PlainSetPtrFromStrCol(col any) *PlainSet[string] {
	s := PlainSetFromStrCol(col)
	return &s
}

func PlainSetFromNodeCol[N ~uint32](col any) PlainSet[N] {
	var s PlainSet[N]
	PlainSetAddAllNode(&s, col)
	return s
}

func PlainSetPtrFromNodeCol[N ~uint32](col any) *PlainSet[N] {
	s := PlainSetFromNodeCol[N](col)
	return &s
}

func PlainSetFromCol[T comparable](col any) PlainSet[T] {
	var s PlainSet[T]
	PlainSetAddAllCol(&s, col)
	return s
}

func PlainSetPtrFromCol[T comparable](col any) *PlainSet[T] {
	s := PlainSetFromCol[T](col)
	return &s
}

func PlainSetAddAllStr(s *PlainSet[string], col any) bool {
	if s == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	changed := false
	if ssi, ok := col.(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			for i := 0; i < len(ss); i++ {
				if s.Add(ss[i]) {
					changed = true
				}
			}
			return changed
		}
	}
	for cur := NewIterCursor(col); cur.HasNext(); {
		if s.Add(cur.NextStr()) {
			changed = true
		}
	}
	return changed
}

func PlainSetAddAllNode[N ~uint32](s *PlainSet[N], col any) bool {
	if s == nil {
		ThrowNPE()
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

func PlainSetAddAllCol[T comparable](s *PlainSet[T], col any) bool {
	if s == nil {
		ThrowNPE()
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

// CustomOrderedMap is a slim 64-byte struct preserving Java LinkedHashMap/HashMap order
// for FastCustomKey types (such as *Var, AbstractVar, ScopedName) with lazy slice and bucket allocation.
type FastCustomKeyComparable interface {
	comparable
	FastCustomKey
}

type CustomMapEntry[K FastCustomKeyComparable, V any] struct {
	Key      K
	Val      V
	Hash     int32
	HashNext int32
}

type CustomOrderedMap[K FastCustomKeyComparable, V any] struct {
	entries  []CustomMapEntry[K, V]
	buckets  []int32
	nullSlot int32
	holes    int32
	firstIdx int32
}

func isCustomKeyNull[K FastCustomKeyComparable](k K) bool {
	var zero K
	return k == zero || k.IsJavaNil_()
}

const (
	custMapPoolSmallMax = 2048
	custMapPoolMedMax   = 256
	custMapPoolLargeMax = 2
	custMapPoolHugeMax  = 2
)

type customOrderedMapBuf[K FastCustomKeyComparable, V any] struct {
	entries []CustomMapEntry[K, V]
	buckets []int32
}

type customOrderedMapTypePool[K FastCustomKeyComparable, V any] struct {
	mu    sync.Mutex
	small []customOrderedMapBuf[K, V]
	med   []customOrderedMapBuf[K, V]
	large []customOrderedMapBuf[K, V]
	huge  []customOrderedMapBuf[K, V]
}

var customOrderedMapPools sync.Map

func getCustomOrderedMapPool[K FastCustomKeyComparable, V any]() *customOrderedMapTypePool[K, V] {
	rt := reflect.TypeFor[customOrderedMapBuf[K, V]]()
	if v, ok := customOrderedMapPools.Load(rt); ok {
		return v.(*customOrderedMapTypePool[K, V])
	}
	p := &customOrderedMapTypePool[K, V]{}
	actual, loaded := customOrderedMapPools.LoadOrStore(rt, p)
	if !loaded {
		purge := func() {
			p.mu.Lock()
			p.small = nil
			p.med = nil
			p.large = nil
			p.huge = nil
			p.mu.Unlock()
		}
		registerCompPoolReset(purge)
		registerRuntimePoolPurger(purge)
	}
	return actual.(*customOrderedMapTypePool[K, V])
}

func AcquireCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	p := getCustomOrderedMapPool[K, V]()
	p.mu.Lock()
	if n := len(p.small); n > 0 {
		b := p.small[n-1]
		p.small[n-1] = customOrderedMapBuf[K, V]{}
		p.small = p.small[:n-1]
		p.mu.Unlock()
		return CustomOrderedMap[K, V]{entries: b.entries[:0]}
	}
	p.mu.Unlock()
	return CustomOrderedMap[K, V]{}
}

func AcquireCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	if capHint < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
	}
	if capHint <= 16 {
		m := AcquireCustomOrderedMap[K, V]()
		if cap(m.entries) == 0 && capHint > 8 {
			m.entries = make([]CustomMapEntry[K, V], 0, capHint)
		}
		return m
	}
	p := getCustomOrderedMapPool[K, V]()
	if capHint <= 512 {
		p.mu.Lock()
		if n := len(p.med); n > 0 {
			b := p.med[n-1]
			p.med[n-1] = customOrderedMapBuf[K, V]{}
			p.med = p.med[:n-1]
			p.mu.Unlock()
			return CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
		}
		p.mu.Unlock()
		return NewCustomOrderedMapCap[K, V](capHint)
	}
	p.mu.Lock()
	if capHint >= 1024 {
		if n := len(p.huge); n > 0 {
			b := p.huge[n-1]
			p.huge[n-1] = customOrderedMapBuf[K, V]{}
			p.huge = p.huge[:n-1]
			p.mu.Unlock()
			return CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
		}
	}
	if n := len(p.large); n > 0 {
		b := p.large[n-1]
		p.large[n-1] = customOrderedMapBuf[K, V]{}
		p.large = p.large[:n-1]
		p.mu.Unlock()
		return CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
	}
	p.mu.Unlock()
	return NewCustomOrderedMapCap[K, V](capHint)
}

func AcquireCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
	m := AcquireCustomOrderedMap[K, V]()
	return &m
}

func AcquireCustomOrderedMapPtrCap[K FastCustomKeyComparable, V any](capHint int32) *CustomOrderedMap[K, V] {
	m := AcquireCustomOrderedMapCap[K, V](capHint)
	return &m
}

func ReleaseCustomOrderedMap[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V]) {
	if m != nil {
		m.ReleaseToPool()
	}
}

func ReleaseCustomOrderedMapPtr[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V]) {
	if m != nil {
		m.ReleaseToPool()
	}
}

func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	es := m.entries
	bk := m.buckets
	nullSlot := m.nullSlot
	firstIdx := m.firstIdx
	*m = CustomOrderedMap[K, V]{}
	c := cap(es)
	bl := len(bk)
	bc := cap(bk)
	if c == 0 && bc == 0 {
		return
	}
	if c > 4194304 || bl > 8388608 {
		if c >= 64 && c <= pooledSliceMaxCap {
			ReleasePooledSlice(es)
		}
		if bc >= 64 && bc <= pooledSliceMaxCap {
			ReleasePooledSlice(bk)
		}
		return
	}
	n := len(es)
	if bk != nil && n > 0 {
		if n*4 < bl {
			mask := uint32(bl - 1)
			for i := int(firstIdx); i < n; i++ {
				e := &es[i]
				if e.HashNext >= 0 && int32(i+1) != nullSlot {
					bk[smearHash(e.Hash)&mask] = 0
				}
			}
		} else {
			clear(bk)
		}
	}
	if n > 0 {
		clear(es)
		es = es[:0]
	}
	p := getCustomOrderedMapPool[K, V]()
	kept := false
	if bk == nil && c <= 32 {
		p.mu.Lock()
		if len(p.small) < custMapPoolSmallMax {
			p.small = append(p.small, customOrderedMapBuf[K, V]{entries: es})
			kept = true
		}
		p.mu.Unlock()
	} else if c > 262144 || bl > 524288 {
		p.mu.Lock()
		if len(p.huge) < custMapPoolHugeMax {
			p.huge = append(p.huge, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
			kept = true
		}
		p.mu.Unlock()
	} else if c <= 512 && bl <= 1024 {
		p.mu.Lock()
		if len(p.med) < custMapPoolMedMax {
			p.med = append(p.med, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
			kept = true
		}
		p.mu.Unlock()
	} else {
		p.mu.Lock()
		if len(p.large) < custMapPoolLargeMax {
			p.large = append(p.large, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
			kept = true
		}
		p.mu.Unlock()
	}
	if !kept {
		if c >= 64 && c <= pooledSliceMaxCap {
			ReleasePooledSlice(es)
		}
		if bc >= 64 && bc <= pooledSliceMaxCap {
			ReleasePooledSlice(bk)
		}
	}
}

func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}

func NewCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
	return &CustomOrderedMap[K, V]{}
}

func NewCustomOrderedMapCap[K FastCustomKeyComparable, V any](cap int32) CustomOrderedMap[K, V] {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	var m CustomOrderedMap[K, V]
	if cap >= 64 {
		m.entries = AcquirePooledSlice[CustomMapEntry[K, V]](int(cap))
	} else if cap > 8 {
		m.entries = make([]CustomMapEntry[K, V], 0, cap)
	}
	return m
}

func NewCustomOrderedMapPtrCap[K FastCustomKeyComparable, V any](cap int32) *CustomOrderedMap[K, V] {
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	m := &CustomOrderedMap[K, V]{}
	if cap >= 64 {
		m.entries = AcquirePooledSlice[CustomMapEntry[K, V]](int(cap))
	} else if cap > 8 {
		m.entries = make([]CustomMapEntry[K, V], 0, cap)
	}
	return m
}

func CustomOrderedMapFromPtr[K FastCustomKeyComparable, V any](src *CustomOrderedMap[K, V]) CustomOrderedMap[K, V] {
	if src == nil {
		ThrowNPE()
	}
	live := int32(len(src.entries)) - src.holes
	if live == 0 {
		return CustomOrderedMap[K, V]{}
	}
	capHint := int(live)
	if capHint < 4 {
		capHint = 4
	}
	var es []CustomMapEntry[K, V]
	if capHint >= 64 {
		es = AcquirePooledSlice[CustomMapEntry[K, V]](capHint)[:live]
	} else {
		es = make([]CustomMapEntry[K, V], live, capHint)
	}
	if src.holes == 0 {
		copy(es, src.entries)
		dst := CustomOrderedMap[K, V]{entries: es, nullSlot: src.nullSlot}
		if src.buckets != nil && len(es) > orderedMapIndexThreshold {
			var b []int32
			if len(src.buckets) >= 64 {
				b = AcquirePooledSlice[int32](len(src.buckets))[:len(src.buckets)]
			} else {
				b = make([]int32, len(src.buckets))
			}
			copy(b, src.buckets)
			dst.buckets = b
		} else if len(es) <= orderedMapIndexThreshold {
			for i := range es {
				es[i].HashNext = 0
			}
		}
		return dst
	}
	dst := CustomOrderedMap[K, V]{entries: es}
	w := 0
	for i := int(src.firstIdx); i < len(src.entries); i++ {
		e := src.entries[i]
		if e.HashNext >= 0 {
			e.HashNext = 0
			if int32(i+1) == src.nullSlot {
				dst.nullSlot = int32(w + 1)
			}
			es[w] = e
			w++
		}
	}
	if len(es) > orderedMapIndexThreshold {
		bCap := 16
		for (bCap*3)/4 < len(es) {
			bCap *= 2
		}
		dst.rehash(bCap)
	}
	return dst
}

func CustomOrderedMapFrom[K FastCustomKeyComparable, V any](src CustomOrderedMap[K, V]) CustomOrderedMap[K, V] {
	return CustomOrderedMapFromPtr(&src)
}

func CustomOrderedMapPtrFrom[K FastCustomKeyComparable, V any](src *CustomOrderedMap[K, V]) *CustomOrderedMap[K, V] {
	dst := CustomOrderedMapFromPtr(src)
	return &dst
}

func CustomOrderedMapFromJavaMap[K FastCustomKeyComparable, V any](src any) CustomOrderedMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewCustomOrderedMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func CustomOrderedMapPtrFromJavaMap[K FastCustomKeyComparable, V any](src any) *CustomOrderedMap[K, V] {
	m := CustomOrderedMapFromJavaMap[K, V](src)
	return &m
}

func (m *CustomOrderedMap[K, V]) IsNil() bool      { return m == nil }
func (m *CustomOrderedMap[K, V]) IsJavaNil_() bool { return m == nil }

func (m *CustomOrderedMap[K, V]) Size() int32 {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.entries)) - m.holes
}

func (m *CustomOrderedMap[K, V]) IsEmpty() bool {
	if m == nil {
		ThrowNPE()
	}
	return int32(len(m.entries)) == m.holes
}

func (m *CustomOrderedMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zero V
	if int32(len(m.entries)) == m.holes {
		return zero, false
	}
	if isCustomKeyNull(k) {
		if m.nullSlot == 0 {
			return zero, false
		}
		return m.entries[m.nullSlot-1].Val, true
	}
	if m.buckets != nil {
		h := k.HashCode()
		slot := smearHash(h) & uint32(len(m.buckets)-1)
		for pos := m.buckets[slot]; pos != 0; {
			e := &m.entries[pos-1]
			if e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
				return e.Val, true
			}
			pos = e.HashNext
		}
		return zero, false
	}
	for i := range m.entries {
		if m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
			return m.entries[i].Val, true
		}
	}
	for i := range m.entries {
		e := &m.entries[i]
		if e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
			return e.Val, true
		}
	}
	return zero, false
}

func (m *CustomOrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}

func (m *CustomOrderedMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return def
}

func (m *CustomOrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}

func (m *CustomOrderedMap[K, V]) rehash(newCap int) {
	oldBk := m.buckets
	needHash := oldBk == nil
	var b []int32
	if newCap >= 64 {
		b = AcquirePooledSlice[int32](newCap)[:newCap]
	} else {
		b = make([]int32, newCap)
	}
	mask := uint32(newCap - 1)
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		e := &m.entries[i]
		if e.HashNext < 0 {
			continue
		}
		if int32(i+1) == m.nullSlot {
			e.HashNext = 0
			continue
		}
		if needHash {
			e.Hash = e.Key.HashCode()
		}
		slot := smearHash(e.Hash) & mask
		e.HashNext = b[slot]
		b[slot] = int32(i + 1)
	}
	m.buckets = b
	if cap(oldBk) >= 64 {
		ReleasePooledSlice(oldBk)
	}
}

func (m *CustomOrderedMap[K, V]) appendEntry(e CustomMapEntry[K, V]) {
	es := m.entries
	if len(es) < cap(es) {
		m.entries = append(es, e)
		return
	}
	m.appendEntryGrow(e)
}

func (m *CustomOrderedMap[K, V]) appendEntryGrow(e CustomMapEntry[K, V]) {
	old := m.entries
	oldLen := len(old)
	oldCap := cap(old)
	newCap := oldCap * 2
	if newCap < 4 {
		newCap = 4
	}
	if newCap <= oldLen {
		newCap = oldLen + 1
	}
	if newCap < 64 {
		next := make([]CustomMapEntry[K, V], oldLen+1, newCap)
		copy(next, old)
		next[oldLen] = e
		m.entries = next
		return
	}
	next := AcquirePooledSlice[CustomMapEntry[K, V]](newCap)
	next = next[:oldLen+1]
	copy(next, old)
	next[oldLen] = e
	if oldCap >= 64 {
		ReleasePooledSlice(old)
	}
	m.entries = next
}

func (m *CustomOrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zeroK K
	var zeroV V
	if isCustomKeyNull(k) {
		if m.nullSlot != 0 {
			old := m.entries[m.nullSlot-1].Val
			m.entries[m.nullSlot-1].Val = v
			return old, true
		}
		ne := CustomMapEntry[K, V]{Key: zeroK, Val: v}
		if m.entries == nil {
			es := make([]CustomMapEntry[K, V], 1, 4)
			es[0] = ne
			m.entries = es
			m.nullSlot = 1
			return zeroV, false
		}
		m.appendEntry(ne)
		m.nullSlot = int32(len(m.entries))
		if m.buckets == nil && len(m.entries) > orderedMapIndexThreshold {
			m.rehash(32)
		}
		return zeroV, false
	}
	if m.buckets != nil {
		h := k.HashCode()
		slot := smearHash(h) & uint32(len(m.buckets)-1)
		for pos := m.buckets[slot]; pos != 0; {
			e := &m.entries[pos-1]
			if e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
				old := e.Val
				e.Val = v
				return old, true
			}
			pos = e.HashNext
		}
		m.appendEntry(CustomMapEntry[K, V]{Key: k, Val: v, Hash: h})
		live := int32(len(m.entries)) - m.holes
		if int(live) > (len(m.buckets)*3)/4 {
			m.rehash(len(m.buckets) * 2)
		} else {
			pos := int32(len(m.entries))
			m.entries[pos-1].HashNext = m.buckets[slot]
			m.buckets[slot] = pos
		}
		return zeroV, false
	}
	if m.entries == nil {
		es := make([]CustomMapEntry[K, V], 1, 4)
		es[0] = CustomMapEntry[K, V]{Key: k, Val: v}
		m.entries = es
		return zeroV, false
	}
	for i := range m.entries {
		if m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
			old := m.entries[i].Val
			m.entries[i].Val = v
			return old, true
		}
	}
	for i := range m.entries {
		e := &m.entries[i]
		if e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
			old := e.Val
			e.Val = v
			return old, true
		}
	}
	m.appendEntry(CustomMapEntry[K, V]{Key: k, Val: v})
	if len(m.entries) > orderedMapIndexThreshold {
		m.rehash(32)
	}
	return zeroV, false
}

func (m *CustomOrderedMap[K, V]) Put(k K, v V) V {
	old, _ := m.PutOk(k, v)
	return old
}

func (m *CustomOrderedMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, ok := m.GetOk(k)
	if !ok || isJavaNullMapVal(cur) {
		m.PutOk(k, v)
	}
	return cur
}

func (m *CustomOrderedMap[K, V]) finishBucketedRemove(slotIdx int32) {
	m.holes++
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		m.entries = m.entries[:0]
		if len(m.buckets) > 8388608 {
			m.buckets = nil
		}
		m.nullSlot = 0
		m.holes = 0
		m.firstIdx = 0
		return
	}
	if slotIdx == m.firstIdx {
		for int(m.firstIdx) < len(m.entries) && m.entries[m.firstIdx].HashNext < 0 {
			m.firstIdx++
		}
	}
	if len(m.entries) > 16 && int(live)*2 < len(m.entries) {
		m.compact()
	}
}

func (m *CustomOrderedMap[K, V]) RemoveOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	var zeroV V
	if int32(len(m.entries)) == m.holes {
		return zeroV, false
	}
	if isCustomKeyNull(k) {
		if m.nullSlot == 0 {
			return zeroV, false
		}
		pos := m.nullSlot
		m.nullSlot = 0
		if m.buckets == nil {
			i := int(pos - 1)
			old := m.entries[i].Val
			copy(m.entries[i:], m.entries[i+1:])
			last := len(m.entries) - 1
			m.entries[last] = CustomMapEntry[K, V]{}
			m.entries = m.entries[:last]
			return old, true
		}
		old := m.entries[pos-1].Val
		m.entries[pos-1] = CustomMapEntry[K, V]{HashNext: -1}
		m.finishBucketedRemove(pos - 1)
		return old, true
	}
	if m.buckets == nil {
		for i := range m.entries {
			if m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
				old := m.entries[i].Val
				copy(m.entries[i:], m.entries[i+1:])
				last := len(m.entries) - 1
				m.entries[last] = CustomMapEntry[K, V]{}
				m.entries = m.entries[:last]
				if m.nullSlot > int32(i+1) {
					m.nullSlot--
				}
				return old, true
			}
		}
		for i := range m.entries {
			e := &m.entries[i]
			if e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
				old := e.Val
				copy(m.entries[i:], m.entries[i+1:])
				last := len(m.entries) - 1
				m.entries[last] = CustomMapEntry[K, V]{}
				m.entries = m.entries[:last]
				if m.nullSlot > int32(i+1) {
					m.nullSlot--
				}
				return old, true
			}
		}
		return zeroV, false
	}
	h := k.HashCode()
	slot := smearHash(h) & uint32(len(m.buckets)-1)
	var prev int32
	for pos := m.buckets[slot]; pos != 0; {
		e := &m.entries[pos-1]
		next := e.HashNext
		if e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
			if prev == 0 {
				m.buckets[slot] = next
			} else {
				m.entries[prev-1].HashNext = next
			}
			old := e.Val
			m.entries[pos-1] = CustomMapEntry[K, V]{HashNext: -1}
			m.finishBucketedRemove(pos - 1)
			return old, true
		}
		prev = pos
		pos = next
	}
	return zeroV, false
}

func (m *CustomOrderedMap[K, V]) Remove(k K) V {
	old, _ := m.RemoveOk(k)
	return old
}

func (m *CustomOrderedMap[K, V]) compact() {
	if m.buckets != nil {
		clear(m.buckets)
	}
	mask := uint32(len(m.buckets) - 1)
	oldNull := m.nullSlot
	var newNull int32
	w := 0
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		e := m.entries[i]
		if e.HashNext < 0 {
			continue
		}
		e.HashNext = 0
		if int32(i+1) == oldNull {
			newNull = int32(w + 1)
		} else if m.buckets != nil {
			slot := smearHash(e.Hash) & mask
			e.HashNext = m.buckets[slot]
			m.buckets[slot] = int32(w + 1)
		}
		m.entries[w] = e
		w++
	}
	clear(m.entries[w:])
	m.entries = m.entries[:w]
	m.nullSlot = newNull
	m.holes = 0
	m.firstIdx = 0
}

func (m *CustomOrderedMap[K, V]) Clear() {
	if m == nil {
		ThrowNPE()
	}
	n := len(m.entries)
	if n == 0 && m.nullSlot == 0 {
		m.holes = 0
		m.firstIdx = 0
		return
	}
	if m.buckets != nil {
		if len(m.buckets) > 8388608 {
			m.buckets = nil
		} else if n*4 < len(m.buckets) {
			mask := uint32(len(m.buckets) - 1)
			for i := int(m.firstIdx); i < n; i++ {
				e := &m.entries[i]
				if e.HashNext >= 0 && int32(i+1) != m.nullSlot {
					m.buckets[smearHash(e.Hash)&mask] = 0
				}
			}
		} else {
			clear(m.buckets)
		}
	}
	if n > 0 {
		clear(m.entries)
		m.entries = m.entries[:0]
	}
	m.nullSlot = 0
	m.holes = 0
	m.firstIdx = 0
}

func (m *CustomOrderedMap[K, V]) PutAllMapPtr(other *CustomOrderedMap[K, V]) {
	if m == nil || other == nil {
		ThrowNPE()
	}
	if int32(len(other.entries)) == other.holes {
		return
	}
	if len(m.entries) == 0 && m.buckets == nil {
		*m = CustomOrderedMapFromPtr(other)
		return
	}
	for i := int(other.firstIdx); i < len(other.entries); i++ {
		if other.entries[i].HashNext >= 0 {
			m.PutOk(other.entries[i].Key, other.entries[i].Val)
		}
	}
}

func (m *CustomOrderedMap[K, V]) PutAllMap(other CustomOrderedMap[K, V]) {
	m.PutAllMapPtr(&other)
}

func (m *CustomOrderedMap[K, V]) PutAllJavaMap(src any) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
}

func CustomOrderedMapPtrGetStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K) string {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return NullString
}

func CustomOrderedMapGetStr[K FastCustomKeyComparable](m CustomOrderedMap[K, string], k K) string {
	return CustomOrderedMapPtrGetStr(&m, k)
}

func CustomOrderedMapPutStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K, v string) string {
	if old, ok := m.PutOk(k, v); ok {
		return old
	}
	return NullString
}

func CustomOrderedMapRemoveStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K) string {
	if old, ok := m.RemoveOk(k); ok {
		return old
	}
	return NullString
}

func CustomOrderedMapPutIfAbsentStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K, v string) string {
	cur, ok := m.GetOk(k)
	if !ok || cur == NullString {
		m.PutOk(k, v)
		return NullString
	}
	return cur
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
		m.PutOk(k, nv)
	}
	return nv
}

func CustomOrderedMapComputeIfAbsentFn[K FastCustomKeyComparable, V comparable](m *CustomOrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(fn)
	if cur, ok := m.GetOk(k); ok && cur != nullV {
		return cur
	}
	raw := fn.Apply(boxJavaElem(k))
	if raw == nil || isJavaNilAny(raw) {
		return nullV
	}
	nv := unboxJavaElem[V](raw)
	if nv != nullV {
		m.PutOk(k, nv)
	}
	return nv
}

type CustomOrderedMapKeyCursor[K FastCustomKeyComparable, V any] struct {
	m   *CustomOrderedMap[K, V]
	pos int32
}

func (m *CustomOrderedMap[K, V]) IterKeys() CustomOrderedMapKeyCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return CustomOrderedMapKeyCursor[K, V]{m: m}
	}
	return CustomOrderedMapKeyCursor[K, V]{m: m, pos: m.firstIdx}
}

func (c *CustomOrderedMapKeyCursor[K, V]) HasNext() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.entries[c.pos].HashNext >= 0 {
			return true
		}
		c.pos++
	}
	return false
}

func (c *CustomOrderedMapKeyCursor[K, V]) Next() K {
	m := c.m
	i := c.pos
	c.pos = i + 1
	return m.entries[i].Key
}

type CustomOrderedMapValCursor[K FastCustomKeyComparable, V any] struct {
	m   *CustomOrderedMap[K, V]
	pos int32
}

func (m *CustomOrderedMap[K, V]) IterValues() CustomOrderedMapValCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return CustomOrderedMapValCursor[K, V]{m: m}
	}
	return CustomOrderedMapValCursor[K, V]{m: m, pos: m.firstIdx}
}

func (c *CustomOrderedMapValCursor[K, V]) HasNext() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.entries[c.pos].HashNext >= 0 {
			return true
		}
		c.pos++
	}
	return false
}

func (c *CustomOrderedMapValCursor[K, V]) Next() V {
	m := c.m
	i := c.pos
	c.pos = i + 1
	return m.entries[i].Val
}

type CustomOrderedMapEntryCursor[K FastCustomKeyComparable, V any] struct {
	m        *CustomOrderedMap[K, V]
	pos      int32
	lastSlot int32
}

func (m *CustomOrderedMap[K, V]) IterEntries() CustomOrderedMapEntryCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return CustomOrderedMapEntryCursor[K, V]{m: m, lastSlot: -1}
	}
	return CustomOrderedMapEntryCursor[K, V]{m: m, pos: m.firstIdx, lastSlot: -1}
}

func (c *CustomOrderedMapEntryCursor[K, V]) NextEntry() bool {
	m := c.m
	if int32(len(m.entries)) == m.holes {
		return false
	}
	for int(c.pos) < len(m.entries) {
		if m.entries[c.pos].HashNext >= 0 {
			c.lastSlot = c.pos
			c.pos++
			return true
		}
		c.pos++
	}
	return false
}

func (c *CustomOrderedMapEntryCursor[K, V]) Key() K {
	return c.m.entries[c.lastSlot].Key
}

func (c *CustomOrderedMapEntryCursor[K, V]) Value() V {
	return c.m.entries[c.lastSlot].Val
}

func (c *CustomOrderedMapEntryCursor[K, V]) SetValue(nv V) V {
	old := c.m.entries[c.lastSlot].Val
	c.m.entries[c.lastSlot].Val = nv
	return old
}

func (m *CustomOrderedMap[K, V]) ForEach(fn func(K, V)) {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.entries[i].HashNext >= 0 {
			fn(m.entries[i].Key, m.entries[i].Val)
		}
	}
}

func (m *CustomOrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(action)
	m.ForEach(func(k K, v V) {
		action.Accept(boxJavaElem(k), boxJavaElem(v))
	})
}

func CustomOrderedMapPtrEquals[K FastCustomKeyComparable, V comparable](a, b *CustomOrderedMap[K, V]) bool {
	if a == nil {
		ThrowNPE()
	}
	if b == nil {
		return false
	}
	aLive := int32(len(a.entries)) - a.holes
	bLive := int32(len(b.entries)) - b.holes
	if aLive != bLive {
		return false
	}
	if aLive == 0 {
		return true
	}
	for i := int(a.firstIdx); i < len(a.entries); i++ {
		if a.entries[i].HashNext >= 0 {
			v2, ok := b.GetOk(a.entries[i].Key)
			if !ok || (a.entries[i].Val != v2 && !javaMapValsEqual(a.entries[i].Val, v2, false)) {
				return false
			}
		}
	}
	return true
}

func CustomOrderedMapEquals[K FastCustomKeyComparable, V comparable](a, b CustomOrderedMap[K, V]) bool {
	return CustomOrderedMapPtrEquals(&a, &b)
}

func (m *CustomOrderedMap[K, V]) EqualsMapPtr(other *CustomOrderedMap[K, V]) bool {
	if m == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	mLive := int32(len(m.entries)) - m.holes
	oLive := int32(len(other.entries)) - other.holes
	if mLive != oLive {
		return false
	}
	if mLive == 0 {
		return true
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.entries[i].HashNext >= 0 {
			v2, ok := other.GetOk(m.entries[i].Key)
			if !ok || !javaMapValsEqual(m.entries[i].Val, v2, false) {
				return false
			}
		}
	}
	return true
}

func (m *CustomOrderedMap[K, V]) EqualsMap(other CustomOrderedMap[K, V]) bool {
	return m.EqualsMapPtr(&other)
}

func (m *CustomOrderedMap[K, V]) HashCodeMap(_ bool) int32 {
	return m.HashCode()
}

func (m *CustomOrderedMap[K, V]) HashCode() int32 {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return 0
	}
	var h int32
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		e := &m.entries[i]
		if e.HashNext >= 0 {
			kh := e.Hash
			if m.buckets == nil && int32(i+1) != m.nullSlot {
				kh = e.Key.HashCode()
			}
			h += kh ^ javalang.ObjectHashCode(boxJavaElem(e.Val))
		}
	}
	return h
}

func (m *CustomOrderedMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	if int32(len(m.entries)) == m.holes {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	m.ForEach(func(k K, v V) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(StrVal(boxJavaElem(k)))
		sb.WriteByte('=')
		sb.WriteString(StrVal(boxJavaElem(v)))
	})
	sb.WriteByte('}')
	return sb.String()
}

func (m *CustomOrderedMap[K, V]) ToImmutableMap() Ext_com_google_common_collect_ImmutableMap {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableMap_
	}
	obj := newImmutableMapShell()
	mc := &obj.inlineCore
	mc.InitExpectedSize(live)
	m.ForEach(func(k K, v V) {
		mc.Put(CheckNotNull(boxJavaElem(k)), CheckNotNull(boxJavaElem(v)))
	})
	mc.sealImmutable()
	return obj
}

func (m *CustomOrderedMap[K, V]) KeysToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableSet_
	}
	sc := &JSetCore{}
	sc.InitExpectedSize(live)
	m.ForEach(func(k K, _ V) {
		sc.Add(CheckNotNull(boxJavaElem(k)))
	})
	return wrapImmutableSetFromSetCoreOwn(sc)
}

func (m *CustomOrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
	if m == nil {
		ThrowNPE()
	}
	live := int32(len(m.entries)) - m.holes
	if live == 0 {
		return emptyImmutableList_
	}
	items := make([]V, 0, live)
	m.ForEach(func(_ K, v V) {
		items = append(items, v)
	})
	return ImmListFromSlice(items)
}

func (m *CustomOrderedMap[K, V]) ToJavaLinkedHashMap() *Ext_java_util_LinkedHashMap {
	if m == nil {
		return nil
	}
	live := int32(len(m.entries)) - m.holes
	obj := &Ext_java_util_LinkedHashMap{}
	obj.BindExtSelf_(obj)
	mc := &obj.inlineCore
	mc.InitExpectedSize(live)
	m.ForEach(func(k K, v V) {
		mc.Put(boxJavaElem(k), boxJavaElem(v))
	})
	return obj
}

func (m *CustomOrderedMap[K, V]) ToJavaMap(_ bool) Ext_java_util_Map {
	return m.ToJavaLinkedHashMap()
}

type jCustMapLiveValues[K FastCustomKeyComparable, V any] struct {
	Ext_java_util_AbstractCollectionBase
	m *CustomOrderedMap[K, V]
}

func (v *jCustMapLiveValues[K, V]) IsJavaNil_() bool { return v == nil }
func (v *jCustMapLiveValues[K, V]) Size() int32 {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	return v.m.Size()
}
func (v *jCustMapLiveValues[K, V]) IsEmpty() bool {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	return v.m.IsEmpty()
}
func (v *jCustMapLiveValues[K, V]) Clear() {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	v.m.Clear()
}
func (v *jCustMapLiveValues[K, V]) Contains(x any) bool {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	if int32(len(v.m.entries)) == v.m.holes {
		return false
	}
	target := ToAny(x)
	for cur := v.m.IterValues(); cur.HasNext(); {
		if javalang.ObjectsEquals(boxJavaElem(cur.Next()), target) {
			return true
		}
	}
	return false
}
func custMapNextActive[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V], from int32) int32 {
	n := int32(len(m.entries))
	if n == m.holes || from < 0 || from >= n {
		return -1
	}
	if m.holes == 0 {
		return from
	}
	for i := int(from); i < int(n); i++ {
		if m.entries[i].HashNext >= 0 {
			return int32(i)
		}
	}
	return -1
}

func (v *jCustMapLiveValues[K, V]) iterInit_() (int32, bool) {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	p := custMapNextActive(v.m, v.m.firstIdx)
	return p, p >= 0
}
func (v *jCustMapLiveValues[K, V]) iterNext_(pos int32) (any, int32) {
	p := custMapNextActive(v.m, pos)
	if p < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	return boxJavaElem(v.m.entries[p].Val), custMapNextActive(v.m, p+1)
}
func (v *jCustMapLiveValues[K, V]) Iterator() any {
	if v == nil || v.m == nil {
		ThrowNPE()
	}
	if int32(len(v.m.entries)) == v.m.holes {
		return emptyMapIter_
	}
	return &jCustMapValIter[K, V]{m: v.m, pos: v.m.firstIdx}
}

type jCustMapValIter[K FastCustomKeyComparable, V any] struct {
	m            *CustomOrderedMap[K, V]
	pos          int32
	lastKey      K
	lastKeyValid bool
}

func (*jCustMapValIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jCustMapValIter[K, V]) IsJavaNil_() bool                       { return it == nil }
func (it *jCustMapValIter[K, V]) HasNext() bool {
	if int32(len(it.m.entries)) == it.m.holes {
		return false
	}
	for int(it.pos) < len(it.m.entries) {
		if it.m.entries[it.pos].HashNext >= 0 {
			return true
		}
		it.pos++
	}
	return false
}
func (it *jCustMapValIter[K, V]) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	i := it.pos
	it.pos = i + 1
	it.lastKey = it.m.entries[i].Key
	it.lastKeyValid = true
	return boxJavaElem(it.m.entries[i].Val)
}
func (it *jCustMapValIter[K, V]) Remove() {
	if !it.lastKeyValid {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	wasShifted := it.m.buckets == nil
	it.m.RemoveOk(it.lastKey)
	if wasShifted && it.pos > 0 {
		it.pos--
	}
	it.lastKeyValid = false
}

func (m *CustomOrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
	if m == nil {
		ThrowNPE()
	}
	v := &jCustMapLiveValues[K, V]{m: m}
	v.BindExtSelf_(v)
	return v
}

type jCustMapLiveKeySet[K FastCustomKeyComparable, V any] struct {
	Ext_java_util_AbstractSetBase
	m *CustomOrderedMap[K, V]
}

func tryUnboxCustomMapKey[K FastCustomKeyComparable](raw any) (K, bool) {
	var zero K
	if raw == nil {
		return zero, true
	}
	if any(zero) == nil {
		if isJavaNilAny(raw) {
			return zero, true
		}
		if !reflect.TypeOf(raw).Comparable() {
			return zero, false
		}
		k, ok := raw.(K)
		return k, ok
	}
	if k, ok := raw.(K); ok {
		if isCustomKeyNull(k) {
			return zero, true
		}
		return k, true
	}
	if isJavaNilAny(raw) {
		return zero, true
	}
	return zero, false
}

func (s *jCustMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
func (s *jCustMapLiveKeySet[K, V]) Size() int32 {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.Size()
}
func (s *jCustMapLiveKeySet[K, V]) IsEmpty() bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.IsEmpty()
}
func (s *jCustMapLiveKeySet[K, V]) Clear() {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	s.m.Clear()
}
func (s *jCustMapLiveKeySet[K, V]) Contains(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	k, ok := tryUnboxCustomMapKey[K](x)
	if !ok {
		return false
	}
	return s.m.ContainsKey(k)
}
func (s *jCustMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	k, ok := tryUnboxCustomMapKey[K](x)
	if !ok {
		return false
	}
	_, removed := s.m.RemoveOk(k)
	return removed
}
func (s *jCustMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if k == NullString {
		return s.m.nullSlot != 0
	}
	return false
}
func (s *jCustMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jCustMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if k == NullString {
		var zero K
		_, removed := s.m.RemoveOk(zero)
		return removed
	}
	return false
}
func (s *jCustMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if u == 0 {
		return s.m.nullSlot != 0
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if uk, ok := tryUnboxCustomMapKey[K](boxNodeById(u)); ok {
			return s.m.ContainsKey(uk)
		}
	}
	return false
}
func (s *jCustMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jCustMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if u == 0 {
		var zero K
		_, removed := s.m.RemoveOk(zero)
		return removed
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if uk, ok := tryUnboxCustomMapKey[K](boxNodeById(u)); ok {
			_, removed := s.m.RemoveOk(uk)
			return removed
		}
	}
	return false
}
func (s *jCustMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	p := custMapNextActive(s.m, s.m.firstIdx)
	return p, p >= 0
}
func (s *jCustMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
	p := custMapNextActive(s.m, pos)
	if p < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	return boxJavaElem(s.m.entries[p].Key), custMapNextActive(s.m, p+1)
}
func (s *jCustMapLiveKeySet[K, V]) Iterator() any {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.entries)) == s.m.holes {
		return emptyMapIter_
	}
	return &jCustMapKeyIter[K, V]{m: s.m, pos: s.m.firstIdx}
}

type jCustMapKeyIter[K FastCustomKeyComparable, V any] struct {
	m            *CustomOrderedMap[K, V]
	pos          int32
	lastKey      K
	lastKeyValid bool
}

func (*jCustMapKeyIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jCustMapKeyIter[K, V]) IsJavaNil_() bool                       { return it == nil }
func (it *jCustMapKeyIter[K, V]) HasNext() bool {
	if int32(len(it.m.entries)) == it.m.holes {
		return false
	}
	for int(it.pos) < len(it.m.entries) {
		if it.m.entries[it.pos].HashNext >= 0 {
			return true
		}
		it.pos++
	}
	return false
}
func (it *jCustMapKeyIter[K, V]) Next() any {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	i := it.pos
	it.pos = i + 1
	it.lastKey = it.m.entries[i].Key
	it.lastKeyValid = true
	return boxJavaElem(it.lastKey)
}
func (it *jCustMapKeyIter[K, V]) Remove() {
	if !it.lastKeyValid {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	wasShifted := it.m.buckets == nil
	it.m.RemoveOk(it.lastKey)
	if wasShifted && it.pos > 0 {
		it.pos--
	}
	it.lastKeyValid = false
}

func (m *CustomOrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
	if m == nil {
		ThrowNPE()
	}
	s := &jCustMapLiveKeySet[K, V]{m: m}
	s.BindExtSelf_(s)
	return s
}

func (m *CustomOrderedMap[K, V]) FirstValue() V {
	if m == nil {
		ThrowNPE()
	}
	for i := int(m.firstIdx); i < len(m.entries); i++ {
		if m.entries[i].HashNext >= 0 {
			return m.entries[i].Val
		}
	}
	panic(javalang.NewNoSuchElementException("iterator exhausted"))
}

func CustomOrderedMapGetAny[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	if int32(len(m.entries)) == m.holes {
		return nil
	}
	k, ok := tryUnboxCustomMapKey[K](rawKey)
	if !ok {
		return nil
	}
	v, found := m.GetOk(k)
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

// IdentityOrderedMap is a slim 40-byte insertion-ordered map always using Go == reference identity
// with lazy slice allocation and orderedMapIndexThreshold = 12 before allocating the Go map index.
type IdentityOrderedMap[K comparable, V any] struct {
	raw OrderedMap[K, V]
}

func normIdentKey[K comparable](k K) K {
	var z K
	if any(z) == nil {
		if any(k) == nil || isJavaNilAny(k) {
			return z
		}
	}
	return k
}

func NewIdentityOrderedMap[K comparable, V any]() IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{}
}

func AcquireIdentityOrderedMap[K comparable, V any]() IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{raw: AcquireOrderedMap[K, V]()}
}

func AcquireIdentityOrderedMapCap[K comparable, V any](capHint int32) IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{raw: AcquireOrderedMapCap[K, V](capHint)}
}

func AcquireIdentityOrderedMapPtr[K comparable, V any]() *IdentityOrderedMap[K, V] {
	m := AcquireIdentityOrderedMap[K, V]()
	return &m
}

func AcquireIdentityOrderedMapPtrCap[K comparable, V any](capHint int32) *IdentityOrderedMap[K, V] {
	m := AcquireIdentityOrderedMapCap[K, V](capHint)
	return &m
}

func (m *IdentityOrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	m.raw.ReleaseToPool()
}

func NewIdentityOrderedMapPtr[K comparable, V any]() *IdentityOrderedMap[K, V] {
	return &IdentityOrderedMap[K, V]{}
}

func NewIdentityOrderedMapCap[K comparable, V any](cap int32) IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{raw: NewOrderedMapCap[K, V](cap)}
}

func NewIdentityOrderedMapPtrCap[K comparable, V any](cap int32) *IdentityOrderedMap[K, V] {
	m := NewIdentityOrderedMapCap[K, V](cap)
	return &m
}

func IdentityOrderedMapFromPtr[K comparable, V any](src *IdentityOrderedMap[K, V]) IdentityOrderedMap[K, V] {
	if src == nil {
		ThrowNPE()
	}
	return IdentityOrderedMap[K, V]{raw: OrderedMapFromPtr(&src.raw)}
}

func IdentityOrderedMapFrom[K comparable, V any](src IdentityOrderedMap[K, V]) IdentityOrderedMap[K, V] {
	return IdentityOrderedMapFromPtr(&src)
}

func IdentityOrderedMapPtrFrom[K comparable, V any](src *IdentityOrderedMap[K, V]) *IdentityOrderedMap[K, V] {
	dst := IdentityOrderedMapFromPtr(src)
	return &dst
}

func IdentityOrderedMapFromJavaMap[K comparable, V any](src any) IdentityOrderedMap[K, V] {
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	m := NewIdentityOrderedMapCap[K, V](mp.Size())
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
	return m
}

func IdentityOrderedMapPtrFromJavaMap[K comparable, V any](src any) *IdentityOrderedMap[K, V] {
	m := IdentityOrderedMapFromJavaMap[K, V](src)
	return &m
}

func (m *IdentityOrderedMap[K, V]) IsNil() bool      { return m == nil }
func (m *IdentityOrderedMap[K, V]) IsJavaNil_() bool { return m == nil }
func (m *IdentityOrderedMap[K, V]) Size() int32 {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.Size()
}
func (m *IdentityOrderedMap[K, V]) IsEmpty() bool {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.IsEmpty()
}
func (m *IdentityOrderedMap[K, V]) GetOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.GetOk(normIdentKey(k))
}
func (m *IdentityOrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}
func (m *IdentityOrderedMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return def
}
func (m *IdentityOrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}
func (m *IdentityOrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.PutOk(normIdentKey(k), v)
}
func (m *IdentityOrderedMap[K, V]) Put(k K, v V) V {
	old, _ := m.PutOk(k, v)
	return old
}
func (m *IdentityOrderedMap[K, V]) PutIfAbsent(k K, v V) V {
	cur, ok := m.GetOk(k)
	if !ok || isJavaNullMapVal(cur) {
		m.PutOk(k, v)
	}
	return cur
}
func (m *IdentityOrderedMap[K, V]) RemoveOk(k K) (V, bool) {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.RemoveOk(normIdentKey(k))
}
func (m *IdentityOrderedMap[K, V]) Remove(k K) V {
	old, _ := m.RemoveOk(k)
	return old
}
func (m *IdentityOrderedMap[K, V]) Clear() {
	if m == nil {
		ThrowNPE()
	}
	m.raw.Clear()
}
func (m *IdentityOrderedMap[K, V]) PutAllMapPtr(other *IdentityOrderedMap[K, V]) {
	if m == nil || other == nil {
		ThrowNPE()
	}
	m.raw.PutAllMapPtr(&other.raw)
}
func (m *IdentityOrderedMap[K, V]) PutAllMap(other IdentityOrderedMap[K, V]) {
	m.PutAllMapPtr(&other)
}
func (m *IdentityOrderedMap[K, V]) PutAllJavaMap(src any) {
	if m == nil {
		ThrowNPE()
	}
	CheckNotNull(src)
	mp := src.(Ext_java_util_Map)
	for cur := NewMapEntryCursor(mp); cur.NextEntry(); {
		m.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
	}
}

func IdentityOrderedMapPtrGetStr[K comparable](m *IdentityOrderedMap[K, string], k K) string {
	if v, ok := m.GetOk(k); ok {
		return v
	}
	return NullString
}
func IdentityOrderedMapGetStr[K comparable](m IdentityOrderedMap[K, string], k K) string {
	return IdentityOrderedMapPtrGetStr(&m, k)
}
func IdentityOrderedMapPutStr[K comparable](m *IdentityOrderedMap[K, string], k K, v string) string {
	if old, ok := m.PutOk(k, v); ok {
		return old
	}
	return NullString
}
func IdentityOrderedMapRemoveStr[K comparable](m *IdentityOrderedMap[K, string], k K) string {
	if old, ok := m.RemoveOk(k); ok {
		return old
	}
	return NullString
}
func IdentityOrderedMapPutIfAbsentStr[K comparable](m *IdentityOrderedMap[K, string], k K, v string) string {
	cur, ok := m.GetOk(k)
	if !ok || cur == NullString {
		m.PutOk(k, v)
		return NullString
	}
	return cur
}
func IdentityOrderedMapComputeIfAbsent[K comparable, V comparable](m *IdentityOrderedMap[K, V], k K, nullV V, fn func(K) V) V {
	if m == nil || fn == nil {
		ThrowNPE()
	}
	nk := normIdentKey(k)
	return OrderedMapComputeIfAbsent(&m.raw, nk, nullV, fn)
}
func IdentityOrderedMapComputeIfAbsentFn[K comparable, V comparable](m *IdentityOrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
	if m == nil {
		ThrowNPE()
	}
	nk := normIdentKey(k)
	return OrderedMapComputeIfAbsentFn(&m.raw, nk, nullV, fn)
}

func (m *IdentityOrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.IterKeys()
}
func (m *IdentityOrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.IterValues()
}
func (m *IdentityOrderedMap[K, V]) IterEntries() OrderedMapEntryCursor[K, V] {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.IterEntries()
}
func (m *IdentityOrderedMap[K, V]) ForEach(fn func(K, V)) {
	if m == nil {
		ThrowNPE()
	}
	m.raw.ForEach(fn)
}
func (m *IdentityOrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
	if m == nil {
		ThrowNPE()
	}
	m.raw.ForEachBiConsumer(action)
}
func IdentityOrderedMapPtrEquals[K comparable, V comparable](a, b *IdentityOrderedMap[K, V]) bool {
	if a == nil {
		ThrowNPE()
	}
	if b == nil {
		return false
	}
	return OrderedMapPtrEquals(&a.raw, &b.raw, true)
}
func IdentityOrderedMapEquals[K comparable, V comparable](a, b IdentityOrderedMap[K, V]) bool {
	return IdentityOrderedMapPtrEquals(&a, &b)
}
func (m *IdentityOrderedMap[K, V]) EqualsMapPtr(other *IdentityOrderedMap[K, V]) bool {
	if m == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	return m.raw.EqualsMapPtr(&other.raw, true)
}
func (m *IdentityOrderedMap[K, V]) EqualsMap(other IdentityOrderedMap[K, V]) bool {
	return m.EqualsMapPtr(&other)
}
func (m *IdentityOrderedMap[K, V]) HashCode() int32 {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.HashCodeMap(true)
}
func (m *IdentityOrderedMap[K, V]) HashCodeMap(_ bool) int32 {
	return m.HashCode()
}
func (m *IdentityOrderedMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	return m.raw.String()
}
func (m *IdentityOrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.ValuesToImmutableList()
}
func (m *IdentityOrderedMap[K, V]) ToJavaIdentityHashMap() *Ext_java_util_IdentityHashMap {
	if m == nil {
		return nil
	}
	return m.raw.ToJavaMap(true).(*Ext_java_util_IdentityHashMap)
}
func (m *IdentityOrderedMap[K, V]) ToJavaMap(_ bool) Ext_java_util_Map {
	if m == nil {
		return nil
	}
	return m.raw.ToJavaMap(true)
}
func (m *IdentityOrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.LiveValues()
}

type jIdentMapLiveKeySet[K comparable, V any] struct {
	Ext_java_util_AbstractSetBase
	m *IdentityOrderedMap[K, V]
}

func (s *jIdentMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
func (s *jIdentMapLiveKeySet[K, V]) Size() int32 {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.Size()
}
func (s *jIdentMapLiveKeySet[K, V]) IsEmpty() bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	return s.m.IsEmpty()
}
func (s *jIdentMapLiveKeySet[K, V]) Clear() {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	s.m.Clear()
}
func (s *jIdentMapLiveKeySet[K, V]) Contains(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	k, ok := tryUnboxOrderedMapKey[K](x)
	if !ok {
		return false
	}
	return s.m.raw.ContainsKey(normIdentKey(k))
}
func (s *jIdentMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	k, ok := tryUnboxOrderedMapKey[K](x)
	if !ok {
		return false
	}
	_, removed := s.m.raw.RemoveOk(normIdentKey(k))
	return removed
}
func (s *jIdentMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if k == NullString {
			return s.m.raw.ContainsKey(zero)
		}
		uk, ok := any(k).(K)
		return ok && s.m.raw.ContainsKey(uk)
	}
	if _, isStr := any(zero).(string); isStr {
		return s.m.raw.ContainsKey(any(k).(K))
	}
	if k == NullString {
		if uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
			return s.m.raw.ContainsKey(uk)
		}
	}
	return false
}
func (s *jIdentMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jIdentMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	var zero K
	if any(zero) == nil {
		if k == NullString {
			_, removed := s.m.raw.RemoveOk(zero)
			return removed
		}
		uk, ok := any(k).(K)
		if !ok {
			return false
		}
		_, removed := s.m.raw.RemoveOk(uk)
		return removed
	}
	if _, isStr := any(zero).(string); isStr {
		_, removed := s.m.raw.RemoveOk(any(k).(K))
		return removed
	}
	if k == NullString {
		if uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
			_, removed := s.m.raw.RemoveOk(uk)
			return removed
		}
	}
	return false
}
func (s *jIdentMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	if u == 0 {
		uk, ok := tryUnboxOrderedMapKey[K](nil)
		return ok && s.m.raw.ContainsKey(uk)
	}
	var zero K
	if any(zero) == nil {
		uk, ok := boxNodeById(u).(K)
		return ok && s.m.raw.ContainsKey(uk)
	}
	switch any(zero).(type) {
	case uint32:
		return s.m.raw.ContainsKey(any(u).(K))
	case nodestore.Node:
		return s.m.raw.ContainsKey(any(nodestore.Node(u)).(K))
	case jRawNodeProvider:
		if globalNodeBoxer != nil {
			if uk, ok := globalNodeBoxer(u).(K); ok {
				return s.m.raw.ContainsKey(uk)
			}
		}
		uk, ok := tryUnboxOrderedMapKey[K](u)
		return ok && s.m.raw.ContainsKey(uk)
	default:
		return false
	}
}
func (s *jIdentMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	panic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
}
func (s *jIdentMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return false
	}
	if u == 0 {
		uk, ok := tryUnboxOrderedMapKey[K](nil)
		if !ok {
			return false
		}
		_, removed := s.m.raw.RemoveOk(uk)
		return removed
	}
	var zero K
	if any(zero) == nil {
		uk, ok := boxNodeById(u).(K)
		if !ok {
			return false
		}
		_, removed := s.m.raw.RemoveOk(uk)
		return removed
	}
	switch any(zero).(type) {
	case uint32:
		_, removed := s.m.raw.RemoveOk(any(u).(K))
		return removed
	case nodestore.Node:
		_, removed := s.m.raw.RemoveOk(any(nodestore.Node(u)).(K))
		return removed
	case jRawNodeProvider:
		if globalNodeBoxer != nil {
			if uk, ok := globalNodeBoxer(u).(K); ok {
				_, removed := s.m.raw.RemoveOk(uk)
				return removed
			}
		}
		uk, ok := tryUnboxOrderedMapKey[K](u)
		if !ok {
			return false
		}
		_, removed := s.m.raw.RemoveOk(uk)
		return removed
	default:
		return false
	}
}
func (s *jIdentMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	p := ordMapNextActive(&s.m.raw, s.m.raw.firstIdx)
	return p, p >= 0
}
func (s *jIdentMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
	p := ordMapNextActive(&s.m.raw, pos)
	if p < 0 {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	return boxJavaElem(s.m.raw.entries[p].Key), ordMapNextActive(&s.m.raw, p+1)
}
func (s *jIdentMapLiveKeySet[K, V]) Iterator() any {
	if s == nil || s.m == nil {
		ThrowNPE()
	}
	if int32(len(s.m.raw.entries)) == s.m.raw.holes {
		return emptyMapIter_
	}
	return &jOrdMapKeyIter[K, V]{m: &s.m.raw, pos: s.m.raw.firstIdx}
}
func (m *IdentityOrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
	if m == nil {
		ThrowNPE()
	}
	s := &jIdentMapLiveKeySet[K, V]{m: m}
	s.BindExtSelf_(s)
	return s
}

func (m *IdentityOrderedMap[K, V]) FirstValue() V {
	if m == nil {
		ThrowNPE()
	}
	return m.raw.FirstValue()
}

func IdentityOrderedMapGetAny[K comparable, V any](m *IdentityOrderedMap[K, V], rawKey any) any {
	if m == nil {
		ThrowNPE()
	}
	k, ok := tryUnboxOrderedMapKey[K](rawKey)
	if !ok {
		return nil
	}
	v, found := m.raw.GetOk(normIdentKey(k))
	if !found {
		return nil
	}
	return boxJavaElem(v)
}

// IdentityOrderedSet is a slim 40-byte insertion-ordered set always using Go == reference identity
// with lazy slice allocation and orderedSetIndexThreshold = 8 before allocating the Go map index.
type IdentityOrderedSet[T comparable] struct {
	raw OrderedSet[T]
}

func NewIdentityOrderedSet[T comparable]() IdentityOrderedSet[T] {
	return IdentityOrderedSet[T]{}
}
func NewIdentityOrderedSetPtr[T comparable]() *IdentityOrderedSet[T] {
	return &IdentityOrderedSet[T]{}
}
func NewIdentityOrderedSetCap[T comparable](cap int32) IdentityOrderedSet[T] {
	return IdentityOrderedSet[T]{raw: NewOrderedSetCap[T](cap)}
}
func NewIdentityOrderedSetPtrCap[T comparable](cap int32) *IdentityOrderedSet[T] {
	s := NewIdentityOrderedSetCap[T](cap)
	return &s
}
func IdentityOrderedSetFromSet[T comparable](src IdentityOrderedSet[T]) IdentityOrderedSet[T] {
	return IdentityOrderedSet[T]{raw: OrderedSetFromSet(src.raw)}
}
func IdentityOrderedSetPtrFromSet[T comparable](src *IdentityOrderedSet[T]) *IdentityOrderedSet[T] {
	if src == nil {
		ThrowNPE()
	}
	dst := IdentityOrderedSetFromSet(*src)
	return &dst
}
func IdentityOrderedSetFromNodeCol[N ~uint32](col any) IdentityOrderedSet[N] {
	return IdentityOrderedSet[N]{raw: OrderedSetFromNodeCol[N](col)}
}
func IdentityOrderedSetPtrFromNodeCol[N ~uint32](col any) *IdentityOrderedSet[N] {
	s := IdentityOrderedSetFromNodeCol[N](col)
	return &s
}
func IdentityOrderedSetFromCol[T comparable](col any) IdentityOrderedSet[T] {
	var s IdentityOrderedSet[T]
	IdentityOrderedSetAddAllCol(&s, col)
	return s
}
func IdentityOrderedSetPtrFromCol[T comparable](col any) *IdentityOrderedSet[T] {
	s := IdentityOrderedSetFromCol[T](col)
	return &s
}
func (s *IdentityOrderedSet[T]) Init() {}
func (s *IdentityOrderedSet[T]) InitCap(cap int32) {
	if s == nil {
		ThrowNPE()
	}
	s.raw.InitCap(cap)
}
func (s *IdentityOrderedSet[T]) IsNil() bool      { return s == nil }
func (s *IdentityOrderedSet[T]) IsJavaNil_() bool { return s == nil }
func (s IdentityOrderedSet[T]) Size() int32       { return s.raw.Size() }
func (s IdentityOrderedSet[T]) IsEmpty() bool     { return s.raw.IsEmpty() }
func (s IdentityOrderedSet[T]) Contains(v T) bool { return s.raw.Contains(normIdentKey(v)) }
func (s IdentityOrderedSet[T]) GetFirst() T       { return s.raw.GetFirst() }
func (s *IdentityOrderedSet[T]) Add(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.raw.Add(normIdentKey(v))
}
func (s *IdentityOrderedSet[T]) Remove(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.raw.Remove(normIdentKey(v))
}
func (s *IdentityOrderedSet[T]) Clear() {
	if s == nil {
		ThrowNPE()
	}
	s.raw.Clear()
}
func (s *IdentityOrderedSet[T]) AddAllSet(other IdentityOrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	return s.raw.AddAllSet(other.raw)
}
func (s *IdentityOrderedSet[T]) AddAllSetPtr(other *IdentityOrderedSet[T]) bool {
	if s == nil || other == nil {
		ThrowNPE()
	}
	return s.raw.AddAllSet(other.raw)
}
func IdentityOrderedSetAddAllNode[N ~uint32](s *IdentityOrderedSet[N], col any) bool {
	if s == nil {
		ThrowNPE()
	}
	return OrderedSetAddAllNode(&s.raw, col)
}
func IdentityOrderedSetAddAllCol[T comparable](s *IdentityOrderedSet[T], col any) bool {
	if s == nil {
		ThrowNPE()
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
func (s *IdentityOrderedSet[T]) Iter() OrderedSetCursor[T] {
	if s == nil {
		ThrowNPE()
	}
	return s.raw.Iter()
}
func (s IdentityOrderedSet[T]) EqualsSet(other IdentityOrderedSet[T]) bool {
	return s.raw.EqualsSet(other.raw)
}
func (s *IdentityOrderedSet[T]) EqualsSetPtr(other *IdentityOrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	return s.raw.EqualsSet(other.raw)
}
func (s IdentityOrderedSet[T]) HashCode() int32 {
	if s.IsEmpty() {
		return 0
	}
	var h int32
	if s.raw.holes > 0 {
		for i, v := range s.raw.order {
			if s.raw.idx[v] == int32(i+1) {
				h += IdentityHash(boxJavaElem(v))
			}
		}
		return h
	}
	for _, v := range s.raw.order {
		h += IdentityHash(boxJavaElem(v))
	}
	return h
}
func (s *IdentityOrderedSet[T]) String() string {
	if s == nil {
		return "null"
	}
	return (&s.raw).String()
}
func (s IdentityOrderedSet[T]) ToJavaSet() Ext_java_util_Set {
	out := &Ext_java_util_HashSetBase{}
	out.BindExtSelf_(out)
	out.inlineSet.identity = true
	if s.IsEmpty() {
		return out
	}
	for cur := s.Iter(); cur.HasNext(); {
		out.Add_1(boxJavaElem(cur.Next()))
	}
	return out
}
func ArrayDequeFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) *Ext_java_util_ArrayDeque {
	return ArrayDequeFromOrderedSet(s.raw)
}
func ArrayDequeFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) *Ext_java_util_ArrayDeque {
	if s == nil {
		ThrowNPE()
	}
	return ArrayDequeFromOrderedSet(s.raw)
}
func DequeFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) Deque[T] {
	return DequeFromOrderedSet(s.raw)
}
func DequePtrFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) *Deque[T] {
	return DequePtrFromOrderedSet(s.raw)
}
func DequeFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) Deque[T] {
	if s == nil {
		ThrowNPE()
	}
	return DequeFromOrderedSet(s.raw)
}
func DequePtrFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) *Deque[T] {
	d := DequeFromIdentityOrderedSetPtr(s)
	return &d
}

// CustomOrderedSet is a 64-byte insertion-ordered set backed by CustomOrderedMap[T, struct{}]
// for custom .Equals/.HashCode types (JSType, ObjectType, *Var, AbstractVar, ScopedName, etc.)
// with lazy slice allocation and orderedMapIndexThreshold = 12 before bucket allocation.
type CustomOrderedSet[T FastCustomKeyComparable] struct {
	m CustomOrderedMap[T, struct{}]
}

type CustomOrderedSetCursor[T FastCustomKeyComparable] struct {
	raw CustomOrderedMapKeyCursor[T, struct{}]
}

func NewCustomOrderedSet[T FastCustomKeyComparable]() CustomOrderedSet[T] {
	return CustomOrderedSet[T]{}
}
func NewCustomOrderedSetPtr[T FastCustomKeyComparable]() *CustomOrderedSet[T] {
	return &CustomOrderedSet[T]{}
}
func NewCustomOrderedSetCap[T FastCustomKeyComparable](cap int32) CustomOrderedSet[T] {
	return CustomOrderedSet[T]{m: NewCustomOrderedMapCap[T, struct{}](cap)}
}
func NewCustomOrderedSetPtrCap[T FastCustomKeyComparable](cap int32) *CustomOrderedSet[T] {
	s := NewCustomOrderedSetCap[T](cap)
	return &s
}
func CustomOrderedSetFromSet[T FastCustomKeyComparable](src CustomOrderedSet[T]) CustomOrderedSet[T] {
	return CustomOrderedSet[T]{m: CustomOrderedMapFromPtr(&src.m)}
}
func CustomOrderedSetPtrFromSet[T FastCustomKeyComparable](src *CustomOrderedSet[T]) *CustomOrderedSet[T] {
	if src == nil {
		ThrowNPE()
	}
	dst := CustomOrderedSetFromSet(*src)
	return &dst
}
func CustomOrderedSetFromCol[T FastCustomKeyComparable](col any) CustomOrderedSet[T] {
	var s CustomOrderedSet[T]
	CustomOrderedSetAddAllCol(&s, col)
	return s
}
func CustomOrderedSetPtrFromCol[T FastCustomKeyComparable](col any) *CustomOrderedSet[T] {
	s := CustomOrderedSetFromCol[T](col)
	return &s
}
func (s *CustomOrderedSet[T]) Init() {}
func (s *CustomOrderedSet[T]) InitCap(cap int32) {
	if s == nil {
		ThrowNPE()
	}
	if cap < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
	}
	if cap > 8 && s.m.entries == nil {
		if cap >= 64 {
			s.m.entries = AcquirePooledSlice[CustomMapEntry[T, struct{}]](int(cap))
		} else {
			s.m.entries = make([]CustomMapEntry[T, struct{}], 0, cap)
		}
	}
}
func (s *CustomOrderedSet[T]) IsNil() bool      { return s == nil }
func (s *CustomOrderedSet[T]) IsJavaNil_() bool { return s == nil }
func (s CustomOrderedSet[T]) Size() int32       { return int32(len(s.m.entries)) - s.m.holes }
func (s CustomOrderedSet[T]) IsEmpty() bool     { return int32(len(s.m.entries)) == s.m.holes }
func (s CustomOrderedSet[T]) Contains(v T) bool {
	_, ok := s.m.GetOk(v)
	return ok
}
func (s CustomOrderedSet[T]) GetFirst() T {
	if int32(len(s.m.entries)) == s.m.holes {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 {
			return s.m.entries[i].Key
		}
	}
	panic(javalang.NewNoSuchElementException("empty"))
}
func (s *CustomOrderedSet[T]) Add(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	_, existed := s.m.PutOk(v, struct{}{})
	return !existed
}
func (s *CustomOrderedSet[T]) Remove(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	_, existed := s.m.RemoveOk(v)
	return existed
}
func (s *CustomOrderedSet[T]) Clear() {
	if s == nil {
		ThrowNPE()
	}
	s.m.Clear()
}
func (s *CustomOrderedSet[T]) AddAllSet(other CustomOrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	if int32(len(other.m.entries)) == other.m.holes {
		return false
	}
	changed := false
	for i := int(other.m.firstIdx); i < len(other.m.entries); i++ {
		if other.m.entries[i].HashNext >= 0 && s.Add(other.m.entries[i].Key) {
			changed = true
		}
	}
	return changed
}
func (s *CustomOrderedSet[T]) AddAllSetPtr(other *CustomOrderedSet[T]) bool {
	if s == nil || other == nil {
		ThrowNPE()
	}
	return s.AddAllSet(*other)
}
func CustomOrderedSetAddAllCol[T FastCustomKeyComparable](s *CustomOrderedSet[T], col any) bool {
	if s == nil {
		ThrowNPE()
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
func (s *CustomOrderedSet[T]) Iter() CustomOrderedSetCursor[T] {
	if s == nil {
		ThrowNPE()
	}
	return CustomOrderedSetCursor[T]{raw: s.m.IterKeys()}
}
func (c *CustomOrderedSetCursor[T]) HasNext() bool { return c.raw.HasNext() }
func (c *CustomOrderedSetCursor[T]) Next() T       { return c.raw.Next() }

func (s CustomOrderedSet[T]) EqualsSet(other CustomOrderedSet[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	if s.IsEmpty() {
		return true
	}
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 && !other.Contains(s.m.entries[i].Key) {
			return false
		}
	}
	return true
}
func (s *CustomOrderedSet[T]) EqualsSetPtr(other *CustomOrderedSet[T]) bool {
	if s == nil {
		ThrowNPE()
	}
	if other == nil {
		return false
	}
	return s.EqualsSet(*other)
}
func (s CustomOrderedSet[T]) HashCode() int32 {
	if s.IsEmpty() {
		return 0
	}
	var h int32
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		e := &s.m.entries[i]
		if e.HashNext >= 0 && int32(i+1) != s.m.nullSlot {
			h += e.Hash
		}
	}
	return h
}
func (s *CustomOrderedSet[T]) String() string {
	if s == nil {
		return "null"
	}
	if s.IsEmpty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(StrVal(boxJavaElem(s.m.entries[i].Key)))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}
func (s CustomOrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	return s.m.KeysToImmutableSet()
}
func (s CustomOrderedSet[T]) ToJavaSet() *Ext_java_util_LinkedHashSet {
	live := s.Size()
	obj := &Ext_java_util_LinkedHashSet{}
	obj.BindExtSelf_(obj)
	obj.inlineSet.InitExpectedSize(live)
	if live == 0 {
		return obj
	}
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 {
			obj.inlineSet.Add(boxJavaElem(s.m.entries[i].Key))
		}
	}
	return obj
}
func ArrayDequeFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) *Ext_java_util_ArrayDeque {
	live := s.Size()
	obj := &Ext_java_util_ArrayDeque{}
	obj.BindExtSelf_(obj)
	dc := &obj.inlineCore
	dc.InitCap(live)
	if live == 0 {
		return obj
	}
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 {
			dc.AddLast(boxJavaElem(s.m.entries[i].Key))
		}
	}
	return obj
}
func ArrayDequeFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) *Ext_java_util_ArrayDeque {
	if s == nil {
		ThrowNPE()
	}
	return ArrayDequeFromCustomOrderedSet(*s)
}
func DequeFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) Deque[T] {
	live := s.Size()
	if live == 0 {
		return Deque[T]{}
	}
	var dst Deque[T]
	dst.InitCap(live)
	for i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
		if s.m.entries[i].HashNext >= 0 {
			dst.AddLast(s.m.entries[i].Key)
		}
	}
	return dst
}
func DequePtrFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) *Deque[T] {
	d := DequeFromCustomOrderedSet(s)
	return &d
}
func DequeFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) Deque[T] {
	if s == nil {
		ThrowNPE()
	}
	return DequeFromCustomOrderedSet(*s)
}
func DequePtrFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) *Deque[T] {
	d := DequeFromCustomOrderedSetPtr(s)
	return &d
}

// Deque is a compact generic power-of-2 ring-buffer deque matching Java ArrayDeque
// with a 4-element inline buffer (0 heap allocations for 0..4 elements) and zero interface boxing.
type Deque[T any] struct {
	buf    []T
	head   int32
	count  int32
	inline [4]T
}

type DequeCursor[T any] struct {
	d   *Deque[T]
	off int32
}

func dequeNullZero[T any]() T {
	var zero T
	if _, isStr := any(zero).(string); isStr {
		return any(NullString).(T)
	}
	return zero
}

func (d *Deque[T]) capMask() int32 {
	if d.buf == nil {
		return 3
	}
	return int32(len(d.buf) - 1)
}

func (d *Deque[T]) slotGet(phys int32) T {
	if d.buf == nil {
		return d.inline[phys]
	}
	return d.buf[phys]
}

func (d *Deque[T]) slotSet(phys int32, v T) {
	if d.buf == nil {
		d.inline[phys] = v
	} else {
		d.buf[phys] = v
	}
}

func NewDeque[T any]() Deque[T] {
	return Deque[T]{}
}

func NewDequePtr[T any]() *Deque[T] {
	return &Deque[T]{}
}

func NewDequeCap[T any](initialCap int32) Deque[T] {
	var d Deque[T]
	d.InitCap(initialCap)
	return d
}

func NewDequePtrCap[T any](initialCap int32) *Deque[T] {
	d := NewDequeCap[T](initialCap)
	return &d
}

func DequeFromDeque[T any](src Deque[T]) Deque[T] {
	if src.count == 0 {
		return Deque[T]{}
	}
	if src.count <= 4 {
		var dst Deque[T]
		mask := src.capMask()
		for i := int32(0); i < src.count; i++ {
			dst.inline[i] = src.slotGet((src.head + i) & mask)
		}
		dst.count = src.count
		return dst
	}
	var dst Deque[T]
	dst.InitCap(src.count)
	firstPart := int32(len(src.buf)) - src.head
	if firstPart >= src.count {
		copy(dst.buf, src.buf[src.head:src.head+src.count])
	} else {
		copy(dst.buf, src.buf[src.head:])
		copy(dst.buf[firstPart:], src.buf[:src.count-firstPart])
	}
	dst.count = src.count
	return dst
}

func DequePtrFromDeque[T any](src *Deque[T]) *Deque[T] {
	if src == nil {
		ThrowNPE()
	}
	dst := DequeFromDeque(*src)
	return &dst
}

func DequeFromSlice[T any](src []T) Deque[T] {
	if src == nil {
		ThrowNPE()
	}
	if len(src) == 0 {
		return Deque[T]{}
	}
	var dst Deque[T]
	dst.InitCap(int32(len(src)))
	for _, v := range src {
		dst.AddLast(v)
	}
	return dst
}

func DequePtrFromSlice[T any](src []T) *Deque[T] {
	d := DequeFromSlice(src)
	return &d
}

func DequeFromOrderedSet[T comparable](s OrderedSet[T]) Deque[T] {
	live := s.Size()
	if live == 0 {
		return Deque[T]{}
	}
	var dst Deque[T]
	dst.InitCap(live)
	if s.holes > 0 {
		for i, v := range s.order {
			if s.idx[v] == int32(i+1) {
				dst.AddLast(v)
			}
		}
	} else {
		for _, v := range s.order {
			dst.AddLast(v)
		}
	}
	return dst
}

func DequePtrFromOrderedSet[T comparable](s OrderedSet[T]) *Deque[T] {
	d := DequeFromOrderedSet(s)
	return &d
}

func DequeFromOrderedSetPtr[T comparable](s *OrderedSet[T]) Deque[T] {
	if s == nil {
		ThrowNPE()
	}
	return DequeFromOrderedSet(*s)
}

func DequePtrFromOrderedSetPtr[T comparable](s *OrderedSet[T]) *Deque[T] {
	d := DequeFromOrderedSetPtr(s)
	return &d
}

func DequeFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) Deque[T] {
	if s == nil {
		ThrowNPE()
	}
	return DequeFromOrderedSet(s.toOrderedSetView())
}

func DequePtrFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) *Deque[T] {
	d := DequeFromLocalOrderedSet(s)
	return &d
}

func DequeFromCol[T any](col any) Deque[T] {
	CheckNotNull(col)
	var d Deque[T]
	d.AddAllCol(col)
	return d
}

func DequePtrFromCol[T any](col any) *Deque[T] {
	d := DequeFromCol[T](col)
	return &d
}

func (d *Deque[T]) Init() {}

func (d *Deque[T]) InitCap(initialCap int32) {
	if d == nil {
		ThrowNPE()
	}
	if initialCap <= 4 || d.buf != nil {
		return
	}
	c := int32(8)
	limit := initialCap
	if limit > (1 << 20) {
		limit = 1 << 20
	}
	for c < limit {
		c <<= 1
	}
	nb := make([]T, c)
	for i := int32(0); i < d.count; i++ {
		nb[i] = d.inline[(d.head+i)&3]
	}
	d.inline = [4]T{}
	d.buf = nb
	d.head = 0
}

func (d *Deque[T]) grow() {
	if len(d.buf) == 0 {
		nb := make([]T, 8)
		for i := int32(0); i < d.count; i++ {
			nb[i] = d.inline[(d.head+i)&3]
		}
		d.inline = [4]T{}
		d.buf = nb
		d.head = 0
		return
	}
	oldCap := int32(len(d.buf))
	if oldCap >= (1 << 30) {
		panic(javalang.NewIllegalStateException("Sorry, deque too big"))
	}
	newCap := oldCap << 1
	nb := make([]T, newCap)
	if d.count > 0 {
		firstPart := oldCap - d.head
		if firstPart >= d.count {
			copy(nb, d.buf[d.head:d.head+d.count])
		} else {
			copy(nb, d.buf[d.head:])
			copy(nb[firstPart:], d.buf[:d.count-firstPart])
		}
	}
	d.buf = nb
	d.head = 0
}

func (d *Deque[T]) IsNil() bool      { return d == nil }
func (d *Deque[T]) IsJavaNil_() bool { return d == nil }

func (d *Deque[T]) Size() int32 {
	if d == nil {
		ThrowNPE()
	}
	return d.count
}

func (d *Deque[T]) IsEmpty() bool {
	if d == nil {
		ThrowNPE()
	}
	return d.count == 0
}

func (d *Deque[T]) At(i int32) T {
	if d == nil {
		ThrowNPE()
	}
	CheckElementIndex(i, d.count)
	return d.slotGet((d.head + i) & d.capMask())
}

func (d *Deque[T]) AddFirst(v T) {
	if d == nil || isJavaNullMapVal(v) {
		ThrowNPE()
	}
	if d.buf == nil {
		if d.count < 4 {
			d.head = (d.head - 1) & 3
			d.inline[d.head] = v
			d.count++
			return
		}
		d.grow()
	} else if int(d.count) == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1) & int32(len(d.buf)-1)
	d.buf[d.head] = v
	d.count++
}

func (d *Deque[T]) Push(v T) {
	d.AddFirst(v)
}

func (d *Deque[T]) OfferFirst(v T) bool {
	d.AddFirst(v)
	return true
}

func (d *Deque[T]) AddLast(v T) {
	if d == nil || isJavaNullMapVal(v) {
		ThrowNPE()
	}
	if d.buf == nil {
		if d.count < 4 {
			d.inline[(d.head+d.count)&3] = v
			d.count++
			return
		}
		d.grow()
	} else if int(d.count) == len(d.buf) {
		d.grow()
	}
	d.buf[(d.head+d.count)&int32(len(d.buf)-1)] = v
	d.count++
}

func (d *Deque[T]) Add(v T) bool {
	d.AddLast(v)
	return true
}

func (d *Deque[T]) Offer(v T) bool {
	d.AddLast(v)
	return true
}

func (d *Deque[T]) OfferLast(v T) bool {
	d.AddLast(v)
	return true
}

func (d *Deque[T]) RemoveFirst() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	var zero T
	if d.buf == nil {
		v := d.inline[d.head]
		d.inline[d.head] = zero
		d.head = (d.head + 1) & 3
		d.count--
		return v
	}
	v := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) & int32(len(d.buf)-1)
	d.count--
	return v
}

func (d *Deque[T]) Pop() T {
	return d.RemoveFirst()
}

func (d *Deque[T]) RemoveLast() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	var zero T
	if d.buf == nil {
		idx := (d.head + d.count - 1) & 3
		v := d.inline[idx]
		d.inline[idx] = zero
		d.count--
		return v
	}
	idx := (d.head + d.count - 1) & int32(len(d.buf)-1)
	v := d.buf[idx]
	d.buf[idx] = zero
	d.count--
	return v
}

func (d *Deque[T]) GetFirst() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.slotGet(d.head)
}

func (d *Deque[T]) Element() T {
	return d.GetFirst()
}

func (d *Deque[T]) GetLast() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *Deque[T]) PollFirst() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return dequeNullZero[T]()
	}
	var zero T
	if d.buf == nil {
		v := d.inline[d.head]
		d.inline[d.head] = zero
		d.head = (d.head + 1) & 3
		d.count--
		return v
	}
	v := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) & int32(len(d.buf)-1)
	d.count--
	return v
}

func (d *Deque[T]) Poll() T {
	return d.PollFirst()
}

func (d *Deque[T]) PollLast() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return dequeNullZero[T]()
	}
	var zero T
	if d.buf == nil {
		idx := (d.head + d.count - 1) & 3
		v := d.inline[idx]
		d.inline[idx] = zero
		d.count--
		return v
	}
	idx := (d.head + d.count - 1) & int32(len(d.buf)-1)
	v := d.buf[idx]
	d.buf[idx] = zero
	d.count--
	return v
}

func (d *Deque[T]) PeekFirst() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return dequeNullZero[T]()
	}
	return d.slotGet(d.head)
}

func (d *Deque[T]) Peek() T {
	return d.PeekFirst()
}

func (d *Deque[T]) PeekLast() T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return dequeNullZero[T]()
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *Deque[T]) PollFirstUnbox() T {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.PollFirst()
}

func (d *Deque[T]) PollLastUnbox() T {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.PollLast()
}

func (d *Deque[T]) PeekFirstUnbox() T {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.dHead()
}

func (d *Deque[T]) dHead() T {
	return d.slotGet(d.head)
}

func (d *Deque[T]) PeekLastUnbox() T {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *Deque[T]) Clear() {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return
	}
	if d.buf == nil {
		d.inline = [4]T{}
		d.head = 0
		d.count = 0
		return
	}
	var zero T
	mask := int32(len(d.buf) - 1)
	for i := int32(0); i < d.count; i++ {
		d.buf[(d.head+i)&mask] = zero
	}
	d.head = 0
	d.count = 0
}

func (d *Deque[T]) deleteAtOffset(off int32) {
	mask := d.capMask()
	var zero T
	if off < (d.count >> 1) {
		for j := off; j > 0; j-- {
			d.slotSet((d.head+j)&mask, d.slotGet((d.head+j-1)&mask))
		}
		d.slotSet(d.head, zero)
		d.head = (d.head + 1) & mask
		d.count--
		return
	}
	last := d.count - 1
	for j := off; j < last; j++ {
		d.slotSet((d.head+j)&mask, d.slotGet((d.head+j+1)&mask))
	}
	d.slotSet((d.head+last)&mask, zero)
	d.count--
}

func (d *Deque[T]) Contains(v T) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || isJavaNullMapVal(v) {
		return false
	}
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
			return true
		}
	}
	return false
}

func (d *Deque[T]) ContainsAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil || isJavaNilAny(target) {
		return false
	}
	t := ToAny(target)
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
			return true
		}
	}
	return false
}

func (d *Deque[T]) RemoveFirstOccurrence(v T) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || isJavaNullMapVal(v) {
		return false
	}
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *Deque[T]) RemoveFirstOccurrenceAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil || isJavaNilAny(target) {
		return false
	}
	t := ToAny(target)
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *Deque[T]) RemoveLastOccurrence(v T) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || isJavaNullMapVal(v) {
		return false
	}
	mask := d.capMask()
	for i := d.count - 1; i >= 0; i-- {
		if javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *Deque[T]) RemoveLastOccurrenceAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil || isJavaNilAny(target) {
		return false
	}
	t := ToAny(target)
	mask := d.capMask()
	for i := d.count - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *Deque[T]) AddAllDeque(other *Deque[T]) bool {
	if d == nil || other == nil {
		ThrowNPE()
	}
	n := other.count
	if n == 0 {
		return false
	}
	if other == d {
		snap := d.IterSlice()
		for _, v := range snap {
			d.AddLast(v)
		}
		return true
	}
	mask := other.capMask()
	for i := int32(0); i < n; i++ {
		d.AddLast(other.slotGet((other.head + i) & mask))
	}
	return true
}

func (d *Deque[T]) AddAllSlice(s []T) bool {
	if d == nil || s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		d.AddLast(v)
	}
	return true
}

func (d *Deque[T]) AddAllCol(col any) bool {
	if d == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	if tv, ok := col.(*immTypedSliceView[T]); ok {
		return d.AddAllSlice(tv.s)
	}
	added := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		d.AddLast(unboxJavaElem[T](CheckNotNull(cur.Next())))
		added = true
	}
	return added
}

func (d *Deque[T]) Iter() DequeCursor[T] {
	if d == nil {
		ThrowNPE()
	}
	return DequeCursor[T]{d: d}
}

func (c *DequeCursor[T]) HasNext() bool {
	return c.off < c.d.count
}

func (c *DequeCursor[T]) Next() T {
	i := c.off
	c.off = i + 1
	return c.d.slotGet((c.d.head + i) & c.d.capMask())
}

func (d *Deque[T]) IterSlice() []T {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return nil
	}
	out := make([]T, d.count)
	if d.buf == nil {
		for i := int32(0); i < d.count; i++ {
			out[i] = d.inline[(d.head+i)&3]
		}
		return out
	}
	firstPart := int32(len(d.buf)) - d.head
	if firstPart >= d.count {
		copy(out, d.buf[d.head:d.head+d.count])
	} else {
		copy(out, d.buf[d.head:])
		copy(out[firstPart:], d.buf[:d.count-firstPart])
	}
	return out
}

func (d *Deque[T]) String() string {
	if d == nil {
		return "null"
	}
	if d.count == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(boxJavaElem(d.slotGet((d.head + i) & mask))))
	}
	sb.WriteByte(']')
	return sb.String()
}

func DequeJoinStrings(d *Deque[string], delim any) string {
	CheckNotNull(delim)
	if d == nil {
		ThrowNPE()
	}
	sep := StrVal(delim)
	if d.count == 0 {
		return ""
	}
	mask := d.capMask()
	if d.count == 1 {
		s := d.slotGet(d.head)
		if s == NullString {
			return "null"
		}
		return s
	}
	var sb strings.Builder
	for i := int32(0); i < d.count; i++ {
		if i > 0 {
			sb.WriteString(sep)
		}
		s := d.slotGet((d.head + i) & mask)
		if s == NullString {
			sb.WriteString("null")
		} else {
			sb.WriteString(s)
		}
	}
	return sb.String()
}

type typedDequeIter[T any] struct {
	d       *Deque[T]
	off     int32
	lastOff int32
}

func (*typedDequeIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *typedDequeIter[T]) IsJavaNil_() bool                       { return it == nil }
func (it *typedDequeIter[T]) HasNext() bool                          { return it.off < it.d.count }
func (it *typedDequeIter[T]) Next() any {
	if it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off++
	return boxJavaElem(v)
}
func (it *typedDequeIter[T]) NextTyped() T {
	if it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off++
	return v
}
func (it *typedDequeIter[T]) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.d.deleteAtOffset(it.lastOff)
	it.off = it.lastOff
	it.lastOff = -1
}

func DequeIterator[T any](d *Deque[T]) Ext_java_util_Iterator {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return emptyMapIter_
	}
	return &typedDequeIter[T]{d: d, off: 0, lastOff: -1}
}

type typedDequeDescIter[T any] struct {
	d       *Deque[T]
	off     int32
	lastOff int32
}

func (*typedDequeDescIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *typedDequeDescIter[T]) IsJavaNil_() bool                       { return it == nil }
func (it *typedDequeDescIter[T]) HasNext() bool                          { return it.off >= 0 && it.off < it.d.count }
func (it *typedDequeDescIter[T]) Next() any {
	if it.off < 0 || it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off--
	return boxJavaElem(v)
}
func (it *typedDequeDescIter[T]) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.d.deleteAtOffset(it.lastOff)
	it.lastOff = -1
}

func DequeDescendingIterator[T any](d *Deque[T]) Ext_java_util_Iterator {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return emptyMapIter_
	}
	return &typedDequeDescIter[T]{d: d, off: d.count - 1, lastOff: -1}
}

// NodeDeque is an unboxed ~uint32 Node specialization of Deque[N] with 1-instruction (v == 0)
// null checks and zero interface boxing across all ring-buffer operations.
type NodeDeque[N ~uint32] struct {
	buf    []N
	head   int32
	count  int32
	inline [4]N
}

type NodeDequeCursor[N ~uint32] struct {
	d   *NodeDeque[N]
	off int32
}

func (d *NodeDeque[N]) capMask() int32 {
	if d.buf == nil {
		return 3
	}
	return int32(len(d.buf) - 1)
}

func (d *NodeDeque[N]) slotGet(phys int32) N {
	if d.buf == nil {
		return d.inline[phys]
	}
	return d.buf[phys]
}

func (d *NodeDeque[N]) slotSet(phys int32, v N) {
	if d.buf == nil {
		d.inline[phys] = v
	} else {
		d.buf[phys] = v
	}
}

func NewNodeDeque[N ~uint32]() NodeDeque[N] {
	return NodeDeque[N]{}
}

func NewNodeDequePtr[N ~uint32]() *NodeDeque[N] {
	return &NodeDeque[N]{}
}

func NewNodeDequeCap[N ~uint32](initialCap int32) NodeDeque[N] {
	var d NodeDeque[N]
	d.InitCap(initialCap)
	return d
}

func NewNodeDequePtrCap[N ~uint32](initialCap int32) *NodeDeque[N] {
	d := NewNodeDequeCap[N](initialCap)
	return &d
}

func NodeDequeFromNodeDeque[N ~uint32](src NodeDeque[N]) NodeDeque[N] {
	if src.count == 0 {
		return NodeDeque[N]{}
	}
	if src.count <= 4 {
		var dst NodeDeque[N]
		mask := src.capMask()
		for i := int32(0); i < src.count; i++ {
			dst.inline[i] = src.slotGet((src.head + i) & mask)
		}
		dst.count = src.count
		return dst
	}
	var dst NodeDeque[N]
	dst.InitCap(src.count)
	firstPart := int32(len(src.buf)) - src.head
	if firstPart >= src.count {
		copy(dst.buf, src.buf[src.head:src.head+src.count])
	} else {
		copy(dst.buf, src.buf[src.head:])
		copy(dst.buf[firstPart:], src.buf[:src.count-firstPart])
	}
	dst.count = src.count
	return dst
}

func NodeDequeFromDeque[N ~uint32](src NodeDeque[N]) NodeDeque[N] {
	return NodeDequeFromNodeDeque(src)
}

func NodeDequePtrFromNodeDeque[N ~uint32](src *NodeDeque[N]) *NodeDeque[N] {
	if src == nil {
		ThrowNPE()
	}
	dst := NodeDequeFromNodeDeque(*src)
	return &dst
}

func NodeDequePtrFromDeque[N ~uint32](src *NodeDeque[N]) *NodeDeque[N] {
	return NodeDequePtrFromNodeDeque(src)
}

func NodeDequeFromSlice[N ~uint32](src []N) NodeDeque[N] {
	if src == nil {
		ThrowNPE()
	}
	if len(src) == 0 {
		return NodeDeque[N]{}
	}
	var dst NodeDeque[N]
	dst.InitCap(int32(len(src)))
	for _, v := range src {
		dst.AddLast(v)
	}
	return dst
}

func NodeDequePtrFromSlice[N ~uint32](src []N) *NodeDeque[N] {
	d := NodeDequeFromSlice(src)
	return &d
}

func NodeDequeFromOrderedSet[N ~uint32](s OrderedSet[N]) NodeDeque[N] {
	live := s.Size()
	if live == 0 {
		return NodeDeque[N]{}
	}
	var dst NodeDeque[N]
	dst.InitCap(live)
	if s.holes > 0 {
		for i, v := range s.order {
			if s.idx[v] == int32(i+1) {
				dst.AddLast(v)
			}
		}
	} else {
		for _, v := range s.order {
			dst.AddLast(v)
		}
	}
	return dst
}

func NodeDequePtrFromOrderedSet[N ~uint32](s OrderedSet[N]) *NodeDeque[N] {
	d := NodeDequeFromOrderedSet(s)
	return &d
}

func NodeDequeFromOrderedSetPtr[N ~uint32](s *OrderedSet[N]) NodeDeque[N] {
	if s == nil {
		ThrowNPE()
	}
	return NodeDequeFromOrderedSet(*s)
}

func NodeDequePtrFromOrderedSetPtr[N ~uint32](s *OrderedSet[N]) *NodeDeque[N] {
	d := NodeDequeFromOrderedSetPtr(s)
	return &d
}

func NodeDequeFromLocalOrderedSet[N ~uint32](s *LocalOrderedSet[N]) NodeDeque[N] {
	if s == nil {
		ThrowNPE()
	}
	return NodeDequeFromOrderedSet(s.toOrderedSetView())
}

func NodeDequePtrFromLocalOrderedSet[N ~uint32](s *LocalOrderedSet[N]) *NodeDeque[N] {
	d := NodeDequeFromLocalOrderedSet(s)
	return &d
}

func NodeDequeFromIdentityOrderedSet[N ~uint32](s IdentityOrderedSet[N]) NodeDeque[N] {
	return NodeDequeFromOrderedSet(s.raw)
}

func NodeDequePtrFromIdentityOrderedSet[N ~uint32](s IdentityOrderedSet[N]) *NodeDeque[N] {
	return NodeDequePtrFromOrderedSet(s.raw)
}

func NodeDequeFromIdentityOrderedSetPtr[N ~uint32](s *IdentityOrderedSet[N]) NodeDeque[N] {
	if s == nil {
		ThrowNPE()
	}
	return NodeDequeFromOrderedSet(s.raw)
}

func NodeDequePtrFromIdentityOrderedSetPtr[N ~uint32](s *IdentityOrderedSet[N]) *NodeDeque[N] {
	d := NodeDequeFromIdentityOrderedSetPtr(s)
	return &d
}

func NodeDequeFromCol[N ~uint32](col any) NodeDeque[N] {
	CheckNotNull(col)
	var d NodeDeque[N]
	d.AddAllCol(col)
	return d
}

func NodeDequePtrFromCol[N ~uint32](col any) *NodeDeque[N] {
	d := NodeDequeFromCol[N](col)
	return &d
}

func (d *NodeDeque[N]) Init() {}

func (d *NodeDeque[N]) InitCap(initialCap int32) {
	if d == nil {
		ThrowNPE()
	}
	if initialCap <= 4 || d.buf != nil {
		return
	}
	c := int32(8)
	limit := initialCap
	if limit > (1 << 20) {
		limit = 1 << 20
	}
	for c < limit {
		c <<= 1
	}
	nb := make([]N, c)
	for i := int32(0); i < d.count; i++ {
		nb[i] = d.inline[(d.head+i)&3]
	}
	d.inline = [4]N{}
	d.buf = nb
	d.head = 0
}

func (d *NodeDeque[N]) grow() {
	if len(d.buf) == 0 {
		nb := make([]N, 8)
		for i := int32(0); i < d.count; i++ {
			nb[i] = d.inline[(d.head+i)&3]
		}
		d.inline = [4]N{}
		d.buf = nb
		d.head = 0
		return
	}
	oldCap := int32(len(d.buf))
	if oldCap >= (1 << 30) {
		panic(javalang.NewIllegalStateException("Sorry, deque too big"))
	}
	newCap := oldCap << 1
	nb := make([]N, newCap)
	if d.count > 0 {
		firstPart := oldCap - d.head
		if firstPart >= d.count {
			copy(nb, d.buf[d.head:d.head+d.count])
		} else {
			copy(nb, d.buf[d.head:])
			copy(nb[firstPart:], d.buf[:d.count-firstPart])
		}
	}
	d.buf = nb
	d.head = 0
}

func (d *NodeDeque[N]) IsNil() bool      { return d == nil }
func (d *NodeDeque[N]) IsJavaNil_() bool { return d == nil }

func (d *NodeDeque[N]) Size() int32 {
	if d == nil {
		ThrowNPE()
	}
	return d.count
}

func (d *NodeDeque[N]) IsEmpty() bool {
	if d == nil {
		ThrowNPE()
	}
	return d.count == 0
}

func (d *NodeDeque[N]) At(i int32) N {
	if d == nil {
		ThrowNPE()
	}
	CheckElementIndex(i, d.count)
	return d.slotGet((d.head + i) & d.capMask())
}

func (d *NodeDeque[N]) AddFirst(v N) {
	if d == nil || v == 0 {
		ThrowNPE()
	}
	if d.buf == nil {
		if d.count < 4 {
			d.head = (d.head - 1) & 3
			d.inline[d.head] = v
			d.count++
			return
		}
		d.grow()
	} else if int(d.count) == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1) & int32(len(d.buf)-1)
	d.buf[d.head] = v
	d.count++
}

func (d *NodeDeque[N]) Push(v N) {
	d.AddFirst(v)
}

func (d *NodeDeque[N]) OfferFirst(v N) bool {
	d.AddFirst(v)
	return true
}

func (d *NodeDeque[N]) AddLast(v N) {
	if d == nil || v == 0 {
		ThrowNPE()
	}
	if d.buf == nil {
		if d.count < 4 {
			d.inline[(d.head+d.count)&3] = v
			d.count++
			return
		}
		d.grow()
	} else if int(d.count) == len(d.buf) {
		d.grow()
	}
	d.buf[(d.head+d.count)&int32(len(d.buf)-1)] = v
	d.count++
}

func (d *NodeDeque[N]) Add(v N) bool {
	d.AddLast(v)
	return true
}

func (d *NodeDeque[N]) Offer(v N) bool {
	d.AddLast(v)
	return true
}

func (d *NodeDeque[N]) OfferLast(v N) bool {
	d.AddLast(v)
	return true
}

func (d *NodeDeque[N]) RemoveFirst() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.buf == nil {
		v := d.inline[d.head]
		d.inline[d.head] = 0
		d.head = (d.head + 1) & 3
		d.count--
		return v
	}
	v := d.buf[d.head]
	d.buf[d.head] = 0
	d.head = (d.head + 1) & int32(len(d.buf)-1)
	d.count--
	return v
}

func (d *NodeDeque[N]) Pop() N {
	return d.RemoveFirst()
}

func (d *NodeDeque[N]) RemoveLast() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.buf == nil {
		idx := (d.head + d.count - 1) & 3
		v := d.inline[idx]
		d.inline[idx] = 0
		d.count--
		return v
	}
	idx := (d.head + d.count - 1) & int32(len(d.buf)-1)
	v := d.buf[idx]
	d.buf[idx] = 0
	d.count--
	return v
}

func (d *NodeDeque[N]) GetFirst() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.slotGet(d.head)
}

func (d *NodeDeque[N]) Element() N {
	return d.GetFirst()
}

func (d *NodeDeque[N]) GetLast() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *NodeDeque[N]) PollFirst() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return 0
	}
	if d.buf == nil {
		v := d.inline[d.head]
		d.inline[d.head] = 0
		d.head = (d.head + 1) & 3
		d.count--
		return v
	}
	v := d.buf[d.head]
	d.buf[d.head] = 0
	d.head = (d.head + 1) & int32(len(d.buf)-1)
	d.count--
	return v
}

func (d *NodeDeque[N]) Poll() N {
	return d.PollFirst()
}

func (d *NodeDeque[N]) PollLast() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return 0
	}
	if d.buf == nil {
		idx := (d.head + d.count - 1) & 3
		v := d.inline[idx]
		d.inline[idx] = 0
		d.count--
		return v
	}
	idx := (d.head + d.count - 1) & int32(len(d.buf)-1)
	v := d.buf[idx]
	d.buf[idx] = 0
	d.count--
	return v
}

func (d *NodeDeque[N]) PeekFirst() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return 0
	}
	return d.slotGet(d.head)
}

func (d *NodeDeque[N]) Peek() N {
	return d.PeekFirst()
}

func (d *NodeDeque[N]) PeekLast() N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return 0
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *NodeDeque[N]) PollFirstUnbox() N {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.PollFirst()
}

func (d *NodeDeque[N]) PollLastUnbox() N {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.PollLast()
}

func (d *NodeDeque[N]) PeekFirstUnbox() N {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.slotGet(d.head)
}

func (d *NodeDeque[N]) PeekLastUnbox() N {
	if d == nil || d.count == 0 {
		ThrowNPE()
	}
	return d.slotGet((d.head + d.count - 1) & d.capMask())
}

func (d *NodeDeque[N]) Clear() {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return
	}
	if d.buf == nil {
		d.inline = [4]N{}
		d.head = 0
		d.count = 0
		return
	}
	mask := int32(len(d.buf) - 1)
	for i := int32(0); i < d.count; i++ {
		d.buf[(d.head+i)&mask] = 0
	}
	d.head = 0
	d.count = 0
}

func (d *NodeDeque[N]) deleteAtOffset(off int32) {
	mask := d.capMask()
	if off < (d.count >> 1) {
		for j := off; j > 0; j-- {
			d.slotSet((d.head+j)&mask, d.slotGet((d.head+j-1)&mask))
		}
		d.slotSet(d.head, 0)
		d.head = (d.head + 1) & mask
		d.count--
		return
	}
	last := d.count - 1
	for j := off; j < last; j++ {
		d.slotSet((d.head+j)&mask, d.slotGet((d.head+j+1)&mask))
	}
	d.slotSet((d.head+last)&mask, 0)
	d.count--
}

func (d *NodeDeque[N]) Contains(v N) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || v == 0 {
		return false
	}
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if d.slotGet((d.head+i)&mask) == v {
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) ContainsAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil {
		return false
	}
	u, ok := extractDequeNodeId(target)
	if !ok || u == 0 {
		return false
	}
	needle := N(u)
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if d.slotGet((d.head+i)&mask) == needle {
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) Remove(v N) bool {
	return d.RemoveFirstOccurrence(v)
}

func (d *NodeDeque[N]) RemoveFirstOccurrence(v N) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || v == 0 {
		return false
	}
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if d.slotGet((d.head+i)&mask) == v {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) RemoveFirstOccurrenceAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil {
		return false
	}
	u, ok := extractDequeNodeId(target)
	if !ok || u == 0 {
		return false
	}
	needle := N(u)
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if d.slotGet((d.head+i)&mask) == needle {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) RemoveLastOccurrence(v N) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || v == 0 {
		return false
	}
	mask := d.capMask()
	for i := d.count - 1; i >= 0; i-- {
		if d.slotGet((d.head+i)&mask) == v {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) RemoveLastOccurrenceAny(target any) bool {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 || target == nil {
		return false
	}
	u, ok := extractDequeNodeId(target)
	if !ok || u == 0 {
		return false
	}
	needle := N(u)
	mask := d.capMask()
	for i := d.count - 1; i >= 0; i-- {
		if d.slotGet((d.head+i)&mask) == needle {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *NodeDeque[N]) AddAllNodeDeque(other *NodeDeque[N]) bool {
	if d == nil || other == nil {
		ThrowNPE()
	}
	n := other.count
	if n == 0 {
		return false
	}
	if other == d {
		snap := d.IterSlice()
		for _, v := range snap {
			d.AddLast(v)
		}
		return true
	}
	mask := other.capMask()
	for i := int32(0); i < n; i++ {
		d.AddLast(other.slotGet((other.head + i) & mask))
	}
	return true
}

func (d *NodeDeque[N]) AddAllDeque(other *NodeDeque[N]) bool {
	return d.AddAllNodeDeque(other)
}

func (d *NodeDeque[N]) AddAllSlice(s []N) bool {
	if d == nil || s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		d.AddLast(v)
	}
	return true
}

func (d *NodeDeque[N]) AddAllCol(col any) bool {
	if d == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	if tv, ok := col.(*immTypedSliceView[N]); ok {
		return d.AddAllSlice(tv.s)
	}
	added := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		d.AddLast(N(cur.NextNode()))
		added = true
	}
	return added
}

func (d *NodeDeque[N]) Iter() NodeDequeCursor[N] {
	if d == nil {
		ThrowNPE()
	}
	return NodeDequeCursor[N]{d: d}
}

func (c *NodeDequeCursor[N]) HasNext() bool {
	return c.off < c.d.count
}

func (c *NodeDequeCursor[N]) Next() N {
	i := c.off
	c.off = i + 1
	return c.d.slotGet((c.d.head + i) & c.d.capMask())
}

func (d *NodeDeque[N]) IterSlice() []N {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return nil
	}
	out := make([]N, d.count)
	if d.buf == nil {
		for i := int32(0); i < d.count; i++ {
			out[i] = d.inline[(d.head+i)&3]
		}
		return out
	}
	firstPart := int32(len(d.buf)) - d.head
	if firstPart >= d.count {
		copy(out, d.buf[d.head:d.head+d.count])
	} else {
		copy(out, d.buf[d.head:])
		copy(out[firstPart:], d.buf[:d.count-firstPart])
	}
	return out
}

func (d *NodeDeque[N]) ToArray_0() []any {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return []any{}
	}
	out := make([]any, d.count)
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		out[i] = boxNodeById(uint32(d.slotGet((d.head + i) & mask)))
	}
	return out
}

func (d *NodeDeque[N]) Iterator() Ext_java_util_Iterator {
	return NodeDequeIterator(d)
}

func (d *NodeDeque[N]) DescendingIterator() Ext_java_util_Iterator {
	return NodeDequeDescendingIterator(d)
}

func (d *NodeDeque[N]) String() string {
	if d == nil {
		return "null"
	}
	if d.count == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	mask := d.capMask()
	for i := int32(0); i < d.count; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(boxNodeById(uint32(d.slotGet((d.head + i) & mask)))))
	}
	sb.WriteByte(']')
	return sb.String()
}

type typedNodeDequeIter[N ~uint32] struct {
	d       *NodeDeque[N]
	off     int32
	lastOff int32
}

func (*typedNodeDequeIter[N]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *typedNodeDequeIter[N]) IsJavaNil_() bool                       { return it == nil }
func (it *typedNodeDequeIter[N]) HasNext() bool                          { return it.off < it.d.count }
func (it *typedNodeDequeIter[N]) Next() any {
	if it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off++
	return boxNodeById(uint32(v))
}
func (it *typedNodeDequeIter[N]) NextTyped() N {
	if it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off++
	return v
}
func (it *typedNodeDequeIter[N]) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.d.deleteAtOffset(it.lastOff)
	it.off = it.lastOff
	it.lastOff = -1
}

func NodeDequeIterator[N ~uint32](d *NodeDeque[N]) Ext_java_util_Iterator {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return emptyMapIter_
	}
	return &typedNodeDequeIter[N]{d: d, off: 0, lastOff: -1}
}

type typedNodeDequeDescIter[N ~uint32] struct {
	d       *NodeDeque[N]
	off     int32
	lastOff int32
}

func (*typedNodeDequeDescIter[N]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *typedNodeDequeDescIter[N]) IsJavaNil_() bool                       { return it == nil }
func (it *typedNodeDequeDescIter[N]) HasNext() bool                          { return it.off >= 0 && it.off < it.d.count }
func (it *typedNodeDequeDescIter[N]) Next() any {
	if it.off < 0 || it.off >= it.d.count {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
	it.lastOff = it.off
	it.off--
	return boxNodeById(uint32(v))
}
func (it *typedNodeDequeDescIter[N]) Remove() {
	if it.lastOff < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.d.deleteAtOffset(it.lastOff)
	it.lastOff = -1
}

func NodeDequeDescendingIterator[N ~uint32](d *NodeDeque[N]) Ext_java_util_Iterator {
	if d == nil {
		ThrowNPE()
	}
	if d.count == 0 {
		return emptyMapIter_
	}
	return &typedNodeDequeDescIter[N]{d: d, off: d.count - 1, lastOff: -1}
}

const (
	slabSmallChunkSize = 16
	slabChunkSize      = 256
	slabPoolMaxBytes   = 384 << 20
	slabSmallPoolMaxB  = 32 << 20
	slabGlobalMaxBytes = 1024 << 20
	slabTrimMinBytes   = 96 << 20
	slabMinEvictBytes  = 256 << 10
)

var (
	slabGlobalPooledBytes atomic.Int64
	slabTrimLastClock     atomic.Uint64
)

// Slab is a compact 32-byte pure bump-pointer arena allocator for package-owned AST/graph leaf structs
// whose lifetime matches their owner and are not pooled across runs.
type Slab[T any] struct {
	buf     []T
	nextCap int32
}

func (s *Slab[T]) Alloc() *T {
	if i := len(s.buf); i < cap(s.buf) {
		b := s.buf[:i+1]
		s.buf = b
		return &b[i]
	}
	return s.allocSlow()
}

func (s *Slab[T]) allocSlow() *T {
	n := s.nextCap
	if n == 0 {
		n = 16
		s.nextCap = 32
	} else if n < 256 {
		s.nextCap = n * 2
	} else {
		s.nextCap = 256
	}
	c := make([]T, 1, n)
	s.buf = c
	return &c[0]
}

func (s *Slab[T]) Release() {
	if s == nil {
		return
	}
	clear(s.buf)
	s.buf = nil
	s.nextCap = 0
}

// PooledSlab is a 48-byte two-tier pooled chunk arena allocator using a pooled 16-item
// slabSmallChunk[T] for small per-script owners before promoting to pooled 128-item slabChunk[T]
// chunks for large whole-AST owners, with zero unpooled bootstrap slices.
type PooledSlab[T any] struct {
	buf   []T
	small *slabSmallChunk[T]
	head  *slabChunk[T]
	pool  *slabChunkPool[T]
}

func (s *PooledSlab[T]) Alloc() *T {
	if i := len(s.buf); i < cap(s.buf) {
		b := s.buf[:i+1]
		s.buf = b
		return &b[i]
	}
	return s.allocSlow()
}

func (s *PooledSlab[T]) allocSlow() *T {
	p := s.pool
	if p == nil {
		p = getSlabChunkPool[T]()
		s.pool = p
	}
	if s.small == nil && s.head == nil {
		sc := p.allocSmall()
		s.small = sc
		s.buf = sc.items[:1:slabSmallChunkSize]
		return &sc.items[0]
	}
	c := p.alloc()
	c.next = s.head
	s.head = c
	s.buf = c.items[:1:slabChunkSize]
	return &c.items[0]
}

func (s *PooledSlab[T]) Reset() {
	if s == nil {
		return
	}
	sc := s.small
	head := s.head
	if head == nil {
		if sc == nil {
			s.buf = nil
			return
		}
		if n := len(s.buf); n > 0 {
			clear(sc.items[:n])
		}
		s.buf = sc.items[:0:slabSmallChunkSize]
		return
	}
	n := len(s.buf)
	s.head = nil
	if sc != nil {
		clear(sc.items[:])
		s.buf = sc.items[:0:slabSmallChunkSize]
	} else {
		s.buf = nil
	}
	if n > 0 {
		clear(head.items[:n])
	}
	for c := head.next; c != nil; c = c.next {
		clear(c.items[:])
	}
	p := s.pool
	if p == nil {
		p = getSlabChunkPool[T]()
		s.pool = p
	}
	p.releaseList(head)
}

func (s *PooledSlab[T]) Release() {
	if s == nil {
		return
	}
	sc := s.small
	head := s.head
	if sc == nil && head == nil {
		s.buf = nil
		return
	}
	n := len(s.buf)
	s.buf = nil
	s.small = nil
	s.head = nil
	p := s.pool
	if p == nil {
		p = getSlabChunkPool[T]()
		s.pool = p
	}
	if head == nil {
		if n > 0 {
			clear(sc.items[:n])
		}
		p.releaseSmall(sc)
		return
	}
	if sc != nil {
		clear(sc.items[:])
		p.releaseSmall(sc)
	}
	if n > 0 {
		clear(head.items[:n])
	}
	for c := head.next; c != nil; c = c.next {
		clear(c.items[:])
	}
	p.releaseList(head)
}

func (s *PooledSlab[T]) ReleaseWith(cleanup func(*T)) {
	if s == nil {
		return
	}
	if cleanup == nil {
		s.Release()
		return
	}
	sc := s.small
	head := s.head
	if sc == nil && head == nil {
		s.buf = nil
		return
	}
	n := len(s.buf)
	s.buf = nil
	s.small = nil
	s.head = nil
	p := s.pool
	if p == nil {
		p = getSlabChunkPool[T]()
		s.pool = p
	}
	if head == nil {
		for i := 0; i < n; i++ {
			cleanup(&sc.items[i])
		}
		if n > 0 {
			clear(sc.items[:n])
		}
		p.releaseSmall(sc)
		return
	}
	if sc != nil {
		for i := 0; i < slabSmallChunkSize; i++ {
			cleanup(&sc.items[i])
		}
		clear(sc.items[:])
		p.releaseSmall(sc)
	}
	for i := 0; i < n; i++ {
		cleanup(&head.items[i])
	}
	if n > 0 {
		clear(head.items[:n])
	}
	for c := head.next; c != nil; c = c.next {
		for i := 0; i < slabChunkSize; i++ {
			cleanup(&c.items[i])
		}
		clear(c.items[:])
	}
	p.releaseList(head)
}

const (
	pooledSliceMinShift        = 3
	pooledSliceMinCap          = 1 << pooledSliceMinShift // 8
	pooledSliceMaxShift        = 22
	pooledSliceMaxCap          = 1 << pooledSliceMaxShift                      // 4194304
	pooledSliceTiers           = pooledSliceMaxShift - pooledSliceMinShift + 1 // 20
	pooledSliceTierBudgetBytes = 4 << 20                                       // 4 MB max per tier per type
	pooledSliceTierMaxCount    = 65536
)

type pooledSliceTier[T any] struct {
	mu      sync.Mutex
	free    [][]T
	maxKeep int
}

type pooledSliceTypePool[T any] struct {
	tiers [pooledSliceTiers]pooledSliceTier[T]
}

var pooledSlicePools sync.Map

func newPooledSliceTypePool[T any]() *pooledSliceTypePool[T] {
	p := &pooledSliceTypePool[T]{}
	elemBytes := max(1, int(reflect.TypeFor[T]().Size()))
	for i := 0; i < pooledSliceTiers; i++ {
		tierCap := pooledSliceMinCap << i
		sliceBytes := tierCap * elemBytes
		var maxKeep int
		if sliceBytes <= 262144 {
			maxKeep = min(pooledSliceTierMaxCount, max(2, pooledSliceTierBudgetBytes/sliceBytes))
		} else {
			maxKeep = 2
		}
		p.tiers[i].maxKeep = maxKeep
	}
	return p
}

func (p *pooledSliceTypePool[T]) purge() {
	for i := 0; i < pooledSliceTiers; i++ {
		t := &p.tiers[i]
		t.mu.Lock()
		t.free = nil
		t.mu.Unlock()
	}
}

func getPooledSlicePool[T any]() *pooledSliceTypePool[T] {
	rt := reflect.TypeFor[T]()
	if v, ok := pooledSlicePools.Load(rt); ok {
		return v.(*pooledSliceTypePool[T])
	}
	p := newPooledSliceTypePool[T]()
	actual, loaded := pooledSlicePools.LoadOrStore(rt, p)
	if !loaded {
		purge := func() {
			p.purge()
		}
		registerCompPoolReset(purge)
		registerRuntimePoolPurger(purge)
	}
	return actual.(*pooledSliceTypePool[T])
}

func AcquirePooledSlice[T any](capHint int) []T {
	if capHint <= 0 {
		return nil
	}
	if capHint > pooledSliceMaxCap {
		return make([]T, 0, capHint)
	}
	tierIdx := 0
	allocCap := pooledSliceMinCap
	if capHint > pooledSliceMinCap {
		shift := bits.Len(uint(capHint - 1))
		tierIdx = shift - pooledSliceMinShift
		allocCap = 1 << shift
	}
	t := &getPooledSlicePool[T]().tiers[tierIdx]
	t.mu.Lock()
	if n := len(t.free); n > 0 {
		s := t.free[n-1]
		t.free[n-1] = nil
		t.free = t.free[:n-1]
		t.mu.Unlock()
		return s[:0]
	}
	t.mu.Unlock()
	return make([]T, 0, allocCap)
}

func AppendPooledSlice[T any](s *[]T, v T) {
	cur := *s
	if len(cur) < cap(cur) {
		*s = append(cur, v)
		return
	}
	appendPooledSliceGrow(s, v)
}

func appendPooledSliceGrow[T any](s *[]T, v T) {
	old := *s
	oldLen := len(old)
	newCap := cap(old) * 2
	if newCap < pooledSliceMinCap {
		newCap = pooledSliceMinCap
	}
	if newCap <= oldLen {
		newCap = oldLen + 1
	}
	next := AcquirePooledSlice[T](newCap)
	next = next[:oldLen+1]
	copy(next, old)
	next[oldLen] = v
	if cap(old) >= pooledSliceMinCap {
		ReleasePooledSlice(old)
	}
	*s = next
}

func ReleasePooledSlice[T any](s []T) {
	c := cap(s)
	if c < pooledSliceMinCap || c > pooledSliceMaxCap {
		return
	}
	clear(s[:c])
	tierIdx := bits.Len(uint(c)) - 1 - pooledSliceMinShift
	t := &getPooledSlicePool[T]().tiers[tierIdx]
	t.mu.Lock()
	if len(t.free) < t.maxKeep {
		t.free = append(t.free, s[:0])
	}
	t.mu.Unlock()
}

type jTypedViewCursorSource interface {
	typedViewLen_() (int, bool)
	typedViewAt_(i int) any
	typedViewNodeAt_(i int) (uint32, bool)
}

func typedElemRawNode[T any](v T) (uint32, bool) {
	var zero T
	switch any(zero).(type) {
	case uint32:
		return any(v).(uint32), true
	case nodestore.Node:
		return uint32(any(v).(nodestore.Node)), true
	}
	if rn, ok := any(v).(jRawNodeProvider); ok {
		return rn.RawNodeId_(), true
	}
	return 0, false
}

// jTypedArrayList wraps a native []T slice as a full mutable Ext_java_util_ArrayList without
// boxing elements into []any unless untyped mutation/subList fallback is invoked.
type jTypedArrayList[T any] struct {
	s  []T
	lc *JListCore
}

var _ Ext_java_util_ArrayList = (*jTypedArrayList[int32])(nil)

func WrapTypedArrayList[T any](s []T) Ext_java_util_ArrayList {
	if s == nil {
		return nil
	}
	return &jTypedArrayList[T]{s: s}
}

func (l *jTypedArrayList[T]) materialize() *JListCore {
	if l.lc == nil {
		items := make([]any, len(l.s))
		for i, v := range l.s {
			items[i] = boxJavaElem(v)
		}
		l.s = nil
		l.lc = &JListCore{Items: items}
	}
	return l.lc
}

func tryUnboxTypedArrayListElem[T any](raw any) (T, bool) {
	var zero T
	if raw == nil || isJavaNilAny(raw) {
		switch any(zero).(type) {
		case string:
			return any(NullString).(T), true
		case int32, int64, bool, float64, float32, uint16, int8, int16:
			return zero, false
		default:
			return zero, true
		}
	}
	switch any(zero).(type) {
	case string:
		if s, ok := raw.(string); ok {
			return any(s).(T), true
		}
		return zero, false
	case int32:
		if b, ok := raw.(*BoxInt); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(int32); ok {
			return any(v).(T), true
		}
		return zero, false
	case int64:
		if b, ok := raw.(*BoxLong); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(int64); ok {
			return any(v).(T), true
		}
		return zero, false
	case bool:
		if b, ok := raw.(*BoxBoolean); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(bool); ok {
			return any(v).(T), true
		}
		return zero, false
	case float64:
		if b, ok := raw.(*BoxDouble); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(float64); ok {
			return any(v).(T), true
		}
		return zero, false
	case float32:
		if b, ok := raw.(*BoxFloat); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(float32); ok {
			return any(v).(T), true
		}
		return zero, false
	case uint16:
		if b, ok := raw.(*BoxChar); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(uint16); ok {
			return any(v).(T), true
		}
		return zero, false
	case int8:
		if b, ok := raw.(*BoxByte); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(int8); ok {
			return any(v).(T), true
		}
		return zero, false
	case int16:
		if b, ok := raw.(*BoxShort); ok && b != nil {
			return any(b.Val).(T), true
		}
		if v, ok := raw.(int16); ok {
			return any(v).(T), true
		}
		return zero, false
	case uint32:
		if rn, ok := raw.(jRawNodeProvider); ok {
			return any(rn.RawNodeId_()).(T), true
		}
		if nd, ok := raw.(nodestore.Node); ok {
			return any(uint32(nd)).(T), true
		}
		if u, ok := raw.(uint32); ok {
			return any(u).(T), true
		}
		return zero, false
	case nodestore.Node:
		if rn, ok := raw.(jRawNodeProvider); ok {
			return any(nodestore.Node(rn.RawNodeId_())).(T), true
		}
		if nd, ok := raw.(nodestore.Node); ok {
			return any(nd).(T), true
		}
		if u, ok := raw.(uint32); ok {
			return any(nodestore.Node(u)).(T), true
		}
		return zero, false
	case jRawNodeProvider:
		if v, ok := raw.(T); ok {
			return v, true
		}
		var u uint32
		var haveU bool
		if rn, ok := raw.(jRawNodeProvider); ok {
			u = rn.RawNodeId_()
			haveU = true
		} else if nd, ok := raw.(nodestore.Node); ok {
			u = uint32(nd)
			haveU = true
		} else if rawU, ok := raw.(uint32); ok {
			u = rawU
			haveU = true
		}
		if !haveU {
			return zero, false
		}
		if u == 0 {
			return zero, true
		}
		if globalNodeBoxer != nil {
			if b, ok := globalNodeBoxer(u).(T); ok {
				return b, true
			}
		}
		rt := reflect.TypeFor[T]()
		if rt.Kind() == reflect.Uint32 {
			return reflect.ValueOf(u).Convert(rt).Interface().(T), true
		}
		return zero, false
	default:
		if v, ok := raw.(T); ok {
			return v, true
		}
		return zero, false
	}
}

func (l *jTypedArrayList[T]) IsJavaNil_() bool { return l == nil }
func (*jTypedArrayList[T]) JavaClassInfo_() (string, string) {
	return "java.util.ArrayList", "ArrayList"
}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_ArrayList()           {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_AbstractList()        {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_List()                {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_RandomAccess()        {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_AbstractCollection()  {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_SequencedCollection() {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_Collection()          {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_lang_Iterable()            {}
func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_io_Serializable()          {}
func (*jTypedArrayList[T]) GetExt_java_util_ArrayListBase() *Ext_java_util_ArrayListBase {
	return nil
}
func (*jTypedArrayList[T]) GetExt_java_util_AbstractListBase() *Ext_java_util_AbstractListBase {
	return nil
}
func (*jTypedArrayList[T]) GetExt_java_util_AbstractCollectionBase() *Ext_java_util_AbstractCollectionBase {
	return nil
}
func (l *jTypedArrayList[T]) typedViewLen_() (int, bool) {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return 0, false
	}
	return len(l.s), true
}
func (l *jTypedArrayList[T]) typedViewAt_(i int) any {
	return boxJavaElem(l.s[i])
}
func (l *jTypedArrayList[T]) typedViewNodeAt_(i int) (uint32, bool) {
	return typedElemRawNode(l.s[i])
}
func (l *jTypedArrayList[T]) IterSlice_() []any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.Items
	}
	if len(l.s) == 0 {
		return emptyAnySlice_
	}
	out := make([]any, len(l.s))
	for i, v := range l.s {
		out[i] = boxJavaElem(v)
	}
	return out
}
func (l *jTypedArrayList[T]) ListGetStr_(i int32) string {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return CastString(l.lc.Get(i))
	}
	CheckElementIndex(i, int32(len(l.s)))
	return CastString(boxJavaElem(l.s[i]))
}
func (l *jTypedArrayList[T]) Size() int32 {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.Size()
	}
	return int32(len(l.s))
}
func (l *jTypedArrayList[T]) IsEmpty() bool {
	return l.Size() == 0
}
func (l *jTypedArrayList[T]) Get(i int32) any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.Get(i)
	}
	CheckElementIndex(i, int32(len(l.s)))
	return boxJavaElem(l.s[i])
}
func (l *jTypedArrayList[T]) GetFirst() any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		if l.lc.IsEmpty() {
			panic(javalang.NewNoSuchElementException("empty"))
		}
		return l.lc.Get(0)
	}
	if len(l.s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return boxJavaElem(l.s[0])
}
func (l *jTypedArrayList[T]) GetLast() any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		if l.lc.IsEmpty() {
			panic(javalang.NewNoSuchElementException("empty"))
		}
		return l.lc.Get(l.lc.Size() - 1)
	}
	if len(l.s) == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return boxJavaElem(l.s[len(l.s)-1])
}
func (l *jTypedArrayList[T]) Set(i int32, x any) any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc == nil {
		if v, ok := tryUnboxTypedArrayListElem[T](x); ok {
			CheckElementIndex(i, int32(len(l.s)))
			old := boxJavaElem(l.s[i])
			l.s[i] = v
			return old
		}
	}
	return l.materialize().Set(i, x)
}
func (l *jTypedArrayList[T]) Add_1(x any) bool {
	if l == nil {
		ThrowNPE()
	}
	if l.lc == nil {
		if v, ok := tryUnboxTypedArrayListElem[T](x); ok {
			l.s = append(l.s, v)
			return true
		}
	}
	return l.materialize().Add(x)
}
func (l *jTypedArrayList[T]) Add_2(i int32, x any) {
	if l == nil {
		ThrowNPE()
	}
	if l.lc == nil {
		if v, ok := tryUnboxTypedArrayListElem[T](x); ok {
			l.s = SliceInsertAt(l.s, i, v)
			return
		}
	}
	l.materialize().AddAt(i, x)
}
func (l *jTypedArrayList[T]) AddFirst(x any) {
	l.Add_2(0, x)
}
func (l *jTypedArrayList[T]) AddLast(x any) {
	l.Add_1(x)
}
func (l *jTypedArrayList[T]) Remove_1_int32(i int32) any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.RemoveAt(i)
	}
	return boxJavaElem(SliceRemoveAtElem(&l.s, i))
}
func (l *jTypedArrayList[T]) RemoveFirst() any {
	if l.Size() == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Remove_1_int32(0)
}
func (l *jTypedArrayList[T]) RemoveLast() any {
	n := l.Size()
	if n == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return l.Remove_1_int32(n - 1)
}
func (l *jTypedArrayList[T]) Remove_1_any(x any) bool {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.RemoveVal(x)
	}
	t := ToAny(x)
	for i, v := range l.s {
		if javalang.ObjectsEquals(boxJavaElem(v), t) {
			SliceRemoveAtElem(&l.s, int32(i))
			return true
		}
	}
	return false
}
func (l *jTypedArrayList[T]) Clear() {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		l.lc.Clear()
		return
	}
	clear(l.s)
	l.s = l.s[:0]
}
func (l *jTypedArrayList[T]) Contains(x any) bool {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.Contains(x)
	}
	t := ToAny(x)
	for _, v := range l.s {
		if javalang.ObjectsEquals(boxJavaElem(v), t) {
			return true
		}
	}
	return false
}
func (l *jTypedArrayList[T]) ContainsAll(col Ext_java_util_Collection) bool {
	return colContainsAll(l.Contains, col)
}
func (l *jTypedArrayList[T]) IndexOf_1(x any) int32 {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.IndexOf(x)
	}
	t := ToAny(x)
	for i, v := range l.s {
		if javalang.ObjectsEquals(boxJavaElem(v), t) {
			return int32(i)
		}
	}
	return -1
}
func (l *jTypedArrayList[T]) LastIndexOf_1(x any) int32 {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return l.lc.LastIndexOf(x)
	}
	t := ToAny(x)
	for i := len(l.s) - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(boxJavaElem(l.s[i]), t) {
			return int32(i)
		}
	}
	return -1
}
func (l *jTypedArrayList[T]) AddAll_1(col Ext_java_util_Collection) bool {
	if l == nil {
		ThrowNPE()
	}
	CheckNotNull(col)
	if l.lc != nil {
		return listCoreAddAllIterable(l.lc, col)
	}
	added := false
	for cur := NewIterCursor(col); cur.HasNext(); {
		l.Add_1(cur.Next())
		added = true
	}
	return added
}
func (l *jTypedArrayList[T]) RemoveAll(col Ext_java_util_Collection) bool {
	if l == nil {
		ThrowNPE()
	}
	return listCoreRemoveAll(l.materialize(), col)
}
func (l *jTypedArrayList[T]) Sort(cmp Ext_java_util_Comparator) {
	if l == nil {
		ThrowNPE()
	}
	listCoreSort(l.materialize(), cmp)
}
func (l *jTypedArrayList[T]) SubList(fromIdx, toIdx int32) any {
	if l == nil {
		ThrowNPE()
	}
	return newSubListView(l.materialize(), fromIdx, toIdx)
}
func (l *jTypedArrayList[T]) ListIterator_1(idx int32) any {
	if l == nil {
		ThrowNPE()
	}
	return newFullListIter(l.materialize(), idx)
}
func (l *jTypedArrayList[T]) Iterator() any {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return newListIter(l.lc)
	}
	return &jTypedArrayListIter[T]{l: l, lastIdx: -1}
}
func (l *jTypedArrayList[T]) ForEach(action Ext_java_util_function_Consumer) {
	if l == nil {
		ThrowNPE()
	}
	CheckNotNull(action)
	if l.lc != nil {
		for _, x := range l.lc.Items {
			action.Accept(x)
		}
		return
	}
	for _, v := range l.s {
		action.Accept(boxJavaElem(v))
	}
}
func (l *jTypedArrayList[T]) ToArray_1_slice_any(dst []any) []any {
	return sliceToArray1(l.IterSlice_(), dst)
}
func (l *jTypedArrayList[T]) Stream() Ext_java_util_stream_Stream {
	return streamFromSlice(l.IterSlice_())
}
func (l *jTypedArrayList[T]) Equals(other any) bool {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return listEqualsSlice(l.lc.Items, other)
	}
	if any(l) == other {
		return true
	}
	return ImmSliceEquals(l.s, other)
}
func (l *jTypedArrayList[T]) HashCode() int32 {
	if l == nil {
		ThrowNPE()
	}
	if l.lc != nil {
		return listHashSlice(l.lc.Items)
	}
	return ImmSliceHashCode(l.s)
}
func (l *jTypedArrayList[T]) String() string {
	if l == nil {
		return "null"
	}
	if l.lc != nil {
		return sliceJavaString(l.lc.Items)
	}
	if len(l.s) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range l.s {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(StrVal(boxJavaElem(v)))
	}
	sb.WriteByte(']')
	return sb.String()
}

type jTypedArrayListIter[T any] struct {
	l       *jTypedArrayList[T]
	pos     int
	lastIdx int
}

func (*jTypedArrayListIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
func (it *jTypedArrayListIter[T]) IsJavaNil_() bool                       { return it == nil }
func (it *jTypedArrayListIter[T]) HasNext() bool                          { return int32(it.pos) < it.l.Size() }
func (it *jTypedArrayListIter[T]) Next() any {
	if int32(it.pos) >= it.l.Size() {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := it.l.Get(int32(it.pos))
	it.lastIdx = it.pos
	it.pos++
	return v
}
func (it *jTypedArrayListIter[T]) Remove() {
	if it.lastIdx < 0 {
		panic(javalang.NewIllegalStateException("remove before next"))
	}
	it.l.Remove_1_int32(int32(it.lastIdx))
	it.pos = it.lastIdx
	it.lastIdx = -1
}

// LocalStringBuilder is a stack-allocated append-only StringBuilder backed by strings.Builder
// with exact Java UTF-16 surrogate-pair stitching and WTF-8 lone-surrogate preservation.
type LocalStringBuilder struct {
	buf         strings.Builder
	pendingHigh uint16
}

func (b *LocalStringBuilder) InitCap(cap int32) *LocalStringBuilder {
	if cap < 0 {
		panic(javalang.NewRuntimeExceptionMsg("NegativeArraySizeException"))
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

func (b *LocalStringBuilder) AppendChars(chars []uint16) *LocalStringBuilder {
	CheckNotNullSlice(chars)
	for _, ch := range chars {
		b.AppendChar(ch)
	}
	return b
}

func (b *LocalStringBuilder) AppendSub(csq any, start, end int32) *LocalStringBuilder {
	return b.AppendStr(JStrSub(StrVal(csq), start, end))
}

func (b *LocalStringBuilder) AppendCodePoint(cp int32) *LocalStringBuilder {
	if cp >= 0 && cp <= 0xFFFF {
		return b.AppendChar(uint16(cp))
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		b.flushPendingHigh()
		b.buf.WriteRune(cp)
		return b
	}
	panic(javalang.NewIllegalArgumentException("Not a valid Unicode code point"))
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
	b.buf.WriteString(javalang.DoubleToString(v))
	return b
}

func (b *LocalStringBuilder) AppendFloat(v float32) *LocalStringBuilder {
	b.flushPendingHigh()
	b.buf.WriteString(javalang.FloatToString(v))
	return b
}

func (b *LocalStringBuilder) AppendAny(v any) *LocalStringBuilder {
	return b.AppendStr(StrVal(v))
}

func (b *LocalStringBuilder) IsEmpty() bool {
	return b.buf.Len() == 0 && b.pendingHigh == 0
}

func (b *LocalStringBuilder) Reset() {
	b.buf.Reset()
	b.pendingHigh = 0
}

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

type slabSmallChunk[T any] struct {
	next  *slabSmallChunk[T]
	items [slabSmallChunkSize]T
}

type slabChunk[T any] struct {
	next  *slabChunk[T]
	items [slabChunkSize]T
}

type slabEvictablePool interface {
	lastTouchVal() uint64
	pooledBytesApprox() int64
	evictBytes(maxBytes int64) int64
}

type slabChunkPool[T any] struct {
	mu           sync.Mutex
	free         []*slabChunk[T]
	smallFree    []*slabSmallChunk[T]
	pooledBytes  atomic.Int64
	lastTouch    atomic.Uint64
	freeLen      atomic.Int32
	smallLen     atomic.Int32
	maxFree      int
	maxSmallFree int
	chunkBytes   int64
	smallBytes   int64
}

var (
	slabClock      atomic.Uint64
	slabChunkPools sync.Map
)

func (p *slabChunkPool[T]) lastTouchVal() uint64 {
	return p.lastTouch.Load()
}

func (p *slabChunkPool[T]) pooledBytesApprox() int64 {
	return p.pooledBytes.Load()
}

func (p *slabChunkPool[T]) evictBytes(maxBytes int64) int64 {
	if maxBytes <= 0 || p.pooledBytes.Load() <= 0 {
		return 0
	}
	p.mu.Lock()
	var freed int64
	if nFree := len(p.free); nFree > 0 {
		rem := maxBytes
		var drop int
		if rem >= int64(nFree)*p.chunkBytes {
			drop = nFree
		} else {
			drop = int((rem + p.chunkBytes - 1) / p.chunkBytes)
		}
		if drop > 0 {
			newLen := nFree - drop
			clear(p.free[newLen:])
			if newLen == 0 {
				p.free = nil
			} else {
				p.free = p.free[:newLen]
			}
			p.freeLen.Store(int32(newLen))
			freed += int64(drop) * p.chunkBytes
		}
	}
	if freed < maxBytes {
		if nSmall := len(p.smallFree); nSmall > 0 {
			rem := maxBytes - freed
			var drop int
			if rem >= int64(nSmall)*p.smallBytes {
				drop = nSmall
			} else {
				drop = int((rem + p.smallBytes - 1) / p.smallBytes)
			}
			if drop > 0 {
				newLen := nSmall - drop
				clear(p.smallFree[newLen:])
				if newLen == 0 {
					p.smallFree = nil
				} else {
					p.smallFree = p.smallFree[:newLen]
				}
				p.smallLen.Store(int32(newLen))
				freed += int64(drop) * p.smallBytes
			}
		}
	}
	if freed > 0 {
		p.pooledBytes.Add(-freed)
		slabGlobalPooledBytes.Add(-freed)
	}
	p.mu.Unlock()
	return freed
}

type slabEvictCand struct {
	pool  slabEvictablePool
	touch uint64
	bytes int64
}

func slabEvictColdestPools(exclude slabEvictablePool, needBytes int64) int64 {
	if needBytes <= 0 {
		return 0
	}
	var cands []slabEvictCand
	slabChunkPools.Range(func(_, v any) bool {
		ep, ok := v.(slabEvictablePool)
		if !ok || ep == exclude {
			return true
		}
		if b := ep.pooledBytesApprox(); b > 0 {
			cands = append(cands, slabEvictCand{
				pool:  ep,
				touch: ep.lastTouchVal(),
				bytes: b,
			})
		}
		return true
	})
	if len(cands) == 0 {
		return 0
	}
	sort.Slice(cands, func(i, j int) bool {
		return cands[i].touch < cands[j].touch
	})
	var freed int64
	for i := range cands {
		stillNeeded := needBytes - freed
		if stillNeeded <= 0 {
			break
		}
		freed += cands[i].pool.evictBytes(stillNeeded)
	}
	return freed
}

func TrimStaleSlabPools(keepRecent int) {
	const maxInt64 = int64(^uint64(0) >> 1)
	if keepRecent <= 0 {
		if slabGlobalPooledBytes.Load() <= 0 {
			return
		}
		slabChunkPools.Range(func(_, v any) bool {
			if ep, ok := v.(slabEvictablePool); ok && ep.pooledBytesApprox() > 0 {
				ep.evictBytes(maxInt64)
			}
			return true
		})
		slabTrimLastClock.Store(slabClock.Load())
		return
	}
	curTotal := slabGlobalPooledBytes.Load()
	if curTotal < slabTrimMinBytes {
		return
	}
	curClock := slabClock.Load()
	if curClock == slabTrimLastClock.Load() && curTotal <= slabGlobalMaxBytes {
		return
	}
	slabTrimLastClock.Store(curClock)
	var cands []slabEvictCand
	slabChunkPools.Range(func(_, v any) bool {
		if ep, ok := v.(slabEvictablePool); ok {
			if b := ep.pooledBytesApprox(); b > 0 {
				cands = append(cands, slabEvictCand{
					pool:  ep,
					touch: ep.lastTouchVal(),
					bytes: b,
				})
			}
		}
		return true
	})
	if len(cands) == 0 {
		return
	}
	sort.Slice(cands, func(i, j int) bool {
		return cands[i].touch > cands[j].touch
	})
	var keptBytes int64
	warmLimit := min(keepRecent, len(cands))
	for i := 0; i < warmLimit; i++ {
		remBudget := int64(slabGlobalMaxBytes) - keptBytes
		if remBudget <= 0 {
			cands[i].pool.evictBytes(maxInt64)
			continue
		}
		if cands[i].bytes > remBudget {
			cands[i].pool.evictBytes(cands[i].bytes - remBudget)
			keptBytes = slabGlobalMaxBytes
		} else {
			keptBytes += cands[i].bytes
		}
	}
	for i := warmLimit; i < len(cands); i++ {
		cands[i].pool.evictBytes(maxInt64)
	}
}

func getSlabChunkPool[T any]() *slabChunkPool[T] {
	rt := reflect.TypeFor[T]()
	if v, ok := slabChunkPools.Load(rt); ok {
		return v.(*slabChunkPool[T])
	}
	chunkBytes := int64(max(1, int(reflect.TypeFor[slabChunk[T]]().Size())))
	smallBytes := int64(max(1, int(reflect.TypeFor[slabSmallChunk[T]]().Size())))
	maxFree := max(16, int(int64(slabPoolMaxBytes)/chunkBytes))
	maxSmallFree := max(64, int(int64(slabSmallPoolMaxB)/smallBytes))
	p := &slabChunkPool[T]{
		maxFree:      maxFree,
		maxSmallFree: maxSmallFree,
		chunkBytes:   chunkBytes,
		smallBytes:   smallBytes,
	}
	p.lastTouch.Store(slabClock.Add(1))
	actual, loaded := slabChunkPools.LoadOrStore(rt, p)
	if !loaded {
		purge := func() {
			p.evictBytes(int64(^uint64(0) >> 1))
			p.mu.Lock()
			if len(p.free) == 0 {
				p.free = nil
			}
			if len(p.smallFree) == 0 {
				p.smallFree = nil
			}
			p.mu.Unlock()
		}
		registerCompPoolReset(purge)
		registerRuntimePoolPurger(purge)
	}
	return actual.(*slabChunkPool[T])
}

func (p *slabChunkPool[T]) alloc() *slabChunk[T] {
	p.lastTouch.Store(slabClock.Add(1))
	p.mu.Lock()
	if n := len(p.free); n > 0 {
		c := p.free[n-1]
		p.free[n-1] = nil
		p.free = p.free[:n-1]
		p.freeLen.Store(int32(n - 1))
		p.pooledBytes.Add(-p.chunkBytes)
		slabGlobalPooledBytes.Add(-p.chunkBytes)
		p.mu.Unlock()
		return c
	}
	p.mu.Unlock()
	return new(slabChunk[T])
}

func (p *slabChunkPool[T]) allocSmall() *slabSmallChunk[T] {
	p.lastTouch.Store(slabClock.Add(1))
	p.mu.Lock()
	if n := len(p.smallFree); n > 0 {
		sc := p.smallFree[n-1]
		p.smallFree[n-1] = nil
		p.smallFree = p.smallFree[:n-1]
		p.smallLen.Store(int32(n - 1))
		p.pooledBytes.Add(-p.smallBytes)
		slabGlobalPooledBytes.Add(-p.smallBytes)
		p.mu.Unlock()
		return sc
	}
	p.mu.Unlock()
	return new(slabSmallChunk[T])
}

func (p *slabChunkPool[T]) releaseSmall(sc *slabSmallChunk[T]) {
	if sc == nil {
		return
	}
	sc.next = nil
	p.lastTouch.Store(slabClock.Add(1))
	sb := p.smallBytes
	if int(p.smallLen.Load()) < p.maxSmallFree {
		if curGlobal := slabGlobalPooledBytes.Load(); curGlobal+sb > slabGlobalMaxBytes {
			need := (curGlobal + sb) - slabGlobalMaxBytes
			slabEvictColdestPools(p, max(need, int64(slabMinEvictBytes)))
		}
	}
	p.mu.Lock()
	if len(p.smallFree) < p.maxSmallFree && slabGlobalPooledBytes.Load()+sb <= slabGlobalMaxBytes {
		if tot := slabGlobalPooledBytes.Add(sb); tot <= slabGlobalMaxBytes {
			p.smallFree = append(p.smallFree, sc)
			p.smallLen.Store(int32(len(p.smallFree)))
			p.pooledBytes.Add(sb)
			p.mu.Unlock()
			return
		} else {
			slabGlobalPooledBytes.Add(-sb)
		}
	}
	p.mu.Unlock()
}

func (p *slabChunkPool[T]) releaseList(first *slabChunk[T]) {
	if first == nil {
		return
	}
	p.lastTouch.Store(slabClock.Add(1))
	cb := p.chunkBytes
	if maxRoom := p.maxFree - int(p.freeLen.Load()); maxRoom > 0 {
		wantCnt := 0
		for c := first; c != nil && wantCnt < maxRoom; c = c.next {
			wantCnt++
		}
		wantBytes := int64(wantCnt) * cb
		if curGlobal := slabGlobalPooledBytes.Load(); curGlobal+wantBytes > slabGlobalMaxBytes {
			need := (curGlobal + wantBytes) - slabGlobalMaxBytes
			slabEvictColdestPools(p, max(need, int64(slabMinEvictBytes)))
		}
	}
	p.mu.Lock()
	free := p.free
	room := p.maxFree - len(free)
	if room > 0 {
		rem := int64(slabGlobalMaxBytes) - slabGlobalPooledBytes.Load()
		if rem <= 0 {
			room = 0
		} else if gr := int(rem / cb); gr < room {
			room = gr
		}
	}
	var added int64
	c := first
	for c != nil && room > 0 {
		cn := c.next
		c.next = nil
		free = append(free, c)
		c = cn
		room--
		added++
	}
	if added > 0 {
		addBytes := added * cb
		if tot := slabGlobalPooledBytes.Add(addBytes); tot > slabGlobalMaxBytes {
			drop := int((tot - slabGlobalMaxBytes + cb - 1) / cb)
			if drop > int(added) {
				drop = int(added)
			}
			if drop > 0 {
				newLen := len(free) - drop
				clear(free[newLen:])
				free = free[:newLen]
				dropBytes := int64(drop) * cb
				slabGlobalPooledBytes.Add(-dropBytes)
				addBytes -= dropBytes
			}
		}
		if addBytes > 0 {
			p.pooledBytes.Add(addBytes)
		}
	}
	p.free = free
	p.freeLen.Store(int32(len(free)))
	p.mu.Unlock()
	for c != nil {
		cn := c.next
		c.next = nil
		c = cn
	}
}

func allocSlabChunk[T any]() *slabChunk[T] {
	return getSlabChunkPool[T]().alloc()
}

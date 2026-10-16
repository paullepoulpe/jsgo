// Port of com/google/javascript/rhino/RhinoStringPool.java @ closure-compiler v20260415

package rhino

import (
	"iter"
	"unique"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RhinoStringPool is an interning pool for strings used by the Rhino
// package.
//
// As of 2021-03-16 and again as of 2024-12-09, this custom pool is
// measurably more performant than `String::intern`.
type RhinoStringPool struct{}

// The threadsafe datastructure that backs this pool.
//
// We use weak-refs, rather than strong-refs, to prevent a memory leak in
// server-like applications.
//
// Go-only: language-adapter — Guava's Interners.newWeakInterner() is Go's
// unique package (canonical, weakly held values); see RhinoStringPool_AddOrGet.

// RhinoStringPool_UncheckedEquals checks if two strings are the same
// according to interning.
//
// The caller is responsible for ensuring that strings passed to this method
// are actually interned. This method mainly exists to highlight where
// equality checks depend on interning for correctness.
//
// Ideally we could use something like a branded-type, for interned strings,
// to verify correct usage. However, Java doesn't support type-brands, and
// using like wrapper objects would undermine performance. This also needs to
// be ergonomic enough that callers don't resort to using `==` directly.
//
// Go-only: language-adapter — Java compares the references; Go compares the
// string values, which is the same for interned strings.
//
// Java: RhinoStringPool.uncheckedEquals(String, String)
func RhinoStringPool_UncheckedEquals(a string, b string) bool {
	return base.JSCompObjects_Identical(a, b)
}

// RhinoStringPool_UncheckedEqualsNullable is uncheckedEquals for Java
// Strings that may be null (*string in the port): null is identical only to
// null.
//
// Go-only: language-adapter — nullable String overload of
// RhinoStringPool.uncheckedEquals(String, String).
func RhinoStringPool_UncheckedEqualsNullable(a *string, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return RhinoStringPool_UncheckedEquals(*a, *b)
}

// Java: RhinoStringPool.addOrGet(String)
func RhinoStringPool_AddOrGet(s string) string {
	return unique.Make(s).Value()
}

// Java: RhinoStringPool()
func NewRhinoStringPool() *RhinoStringPool { return &RhinoStringPool{} }

// RhinoStringPool_LazyInternedStringList is a list of strings that is lazily
// interned into a RhinoStringPool as they are accessed.
//
// The only reason to use this class is for performance optimizations - the
// LazyInternedStringList#newStringNode method skips re-interning the string.
// This is a performance win if creating a lot of string nodes with the same
// string.
type RhinoStringPool_LazyInternedStringList struct {
	pool *jutil.ArrayList[string]

	isInterned *RhinoStringPool_WriteOnlyBitset
}

// Java: RhinoStringPool.LazyInternedStringList(List<String>)
func NewRhinoStringPool_LazyInternedStringList(pool []string) *RhinoStringPool_LazyInternedStringList {
	return &RhinoStringPool_LazyInternedStringList{
		pool:       jutil.NewArrayListFrom(pool),
		isInterned: NewRhinoStringPool_WriteOnlyBitset(len(pool)),
	}
}

// Get returns the interned string at offset.
//
// if Java supported type-brands we would want to use one here to indicate
// that the returned string is always interned but that's not possible
// without an expensive wrapper type - see also the comment on
// uncheckedEquals().
//
// Java: RhinoStringPool.LazyInternedStringList.get(int)
func (l *RhinoStringPool_LazyInternedStringList) Get(offset int) string {
	if !l.isInterned.Get(offset) {
		uninterned := l.pool.Get(offset)
		interned := RhinoStringPool_AddOrGet(uninterned)
		l.pool.Set(offset, interned)
		l.isInterned.Set(offset) // must happen after this.pool.set() to avoid a race condition
		return interned
	}
	return l.pool.Get(offset)
}

// Java: RhinoStringPool.LazyInternedStringList.stream()
func (l *RhinoStringPool_LazyInternedStringList) Stream() iter.Seq[string] {
	return l.pool.All()
}

// RhinoStringPool_WriteOnlyBitset is a bitset implementation that only
// supports setting bits from false => true, not vice versa, and is not
// resizable once created.
//
// This class was originally copied from the JDK BitSet implementation, but
// has been heavily modified to be thread-safe and remove unnecessary
// functionality.
//
// Go-only: language-adapter — the AtomicIntegerArray is a plain []int32
// (the port is single-threaded; see PORTING.md §5); out-of-range accesses
// throw ArrayIndexOutOfBoundsException like the JDK's.
type RhinoStringPool_WriteOnlyBitset struct {
	bits []int32
}

// BitSets are packed into arrays of "words." Currently a word is an int,
// which consists of 32 bits, requiring 5 address bits. The choice of word
// size comes from the convenience of using AtomicIntegerArray.
const RhinoStringPool_WriteOnlyBitset_ADDRESS_BITS_PER_WORD = 5

// Java: RhinoStringPool.WriteOnlyBitset(int)
func NewRhinoStringPool_WriteOnlyBitset(size int) *RhinoStringPool_WriteOnlyBitset {
	n := RhinoStringPool_WriteOnlyBitset_WordIndex(size-1) + 1
	if n < 0 {
		panic(jutil.NewRuntimeException("java.lang.NegativeArraySizeException: " + jutil.ToJavaString(n)))
	}
	return &RhinoStringPool_WriteOnlyBitset{bits: make([]int32, n)}
}

func (b *RhinoStringPool_WriteOnlyBitset) checkIndex(wordIndex int) {
	if wordIndex < 0 || wordIndex >= len(b.bits) {
		panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(wordIndex, len(b.bits)))
	}
}

// Java: RhinoStringPool.WriteOnlyBitset.set(int)
func (b *RhinoStringPool_WriteOnlyBitset) Set(offset int) {
	wordIndex := RhinoStringPool_WriteOnlyBitset_WordIndex(offset)
	b.checkIndex(wordIndex)
	b.bits[wordIndex] = b.bits[wordIndex] | jutil.Shl32(1, int64(offset))
}

// Java: RhinoStringPool.WriteOnlyBitset.get(int)
func (b *RhinoStringPool_WriteOnlyBitset) Get(offset int) bool {
	wordIndex := RhinoStringPool_WriteOnlyBitset_WordIndex(offset)
	b.checkIndex(wordIndex)
	return (b.bits[wordIndex] & jutil.Shl32(1, int64(offset))) != 0
}

// RhinoStringPool_WriteOnlyBitset_WordIndex returns, given a bit index, the
// word index containing it.
//
// Java: RhinoStringPool.WriteOnlyBitset.wordIndex(int)
func RhinoStringPool_WriteOnlyBitset_WordIndex(bitIndex int) int {
	return int(int32(bitIndex) >> RhinoStringPool_WriteOnlyBitset_ADDRESS_BITS_PER_WORD)
}

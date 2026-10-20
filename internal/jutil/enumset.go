package jutil

import (
	"iter"
	"math/bits"
	"strings"
)

// Go-only: language-adapter — java.util.EnumSet (and Guava's ImmutableEnumSet,
// which is what ImmutableSet.copyOf(EnumSet) returns): a bit set over an enum
// type's ordinals that iterates in ordinal order, like Java.

// EnumSet is java.util.EnumSet<E> for a ported enum type E (an int32 ordinal).
type EnumSet[E ~int32] struct {
	words []uint64
	buf   [2]uint64
}

// NewEnumSet is EnumSet.noneOf / EnumSet.of(e...).
func NewEnumSet[E ~int32](elems ...E) *EnumSet[E] {
	s := &EnumSet[E]{}
	s.words = s.buf[:0]
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// EnumSetCopyOf is EnumSet.copyOf(set).
func EnumSetCopyOf[E ~int32](other *EnumSet[E]) *EnumSet[E] {
	s := &EnumSet[E]{}
	if len(other.words) <= len(s.buf) {
		s.words = append(s.buf[:0], other.words...)
	} else {
		s.words = append([]uint64(nil), other.words...)
	}
	return s
}

// Add is Set.add; reports whether the set changed.
func (s *EnumSet[E]) Add(e E) bool {
	w, b := int(e)/64, uint(e)%64
	if s.words == nil && w < len(s.buf) {
		s.words = s.buf[:0]
	}
	for len(s.words) <= w {
		s.words = append(s.words, 0)
	}
	if s.words[w]&(1<<b) != 0 {
		return false
	}
	s.words[w] |= 1 << b
	return true
}

// Remove is Set.remove; reports whether the set changed.
func (s *EnumSet[E]) Remove(e E) bool {
	if !s.Contains(e) {
		return false
	}
	s.words[int(e)/64] &^= 1 << (uint(e) % 64)
	return true
}

// Contains is Set.contains.
func (s *EnumSet[E]) Contains(e E) bool {
	if s == nil || e < 0 {
		return false
	}
	w := int(e) / 64
	return w < len(s.words) && s.words[w]&(1<<(uint(e)%64)) != 0
}

// ContainsAll is Set.containsAll.
func (s *EnumSet[E]) ContainsAll(other *EnumSet[E]) bool {
	for i, w := range other.words {
		var mine uint64
		if i < len(s.words) {
			mine = s.words[i]
		}
		if w&^mine != 0 {
			return false
		}
	}
	return true
}

// AddAll is Set.addAll; reports whether the set changed.
func (s *EnumSet[E]) AddAll(other *EnumSet[E]) bool {
	changed := false
	if s.words == nil && len(other.words) <= len(s.buf) {
		s.words = s.buf[:0]
	}
	for len(s.words) < len(other.words) {
		s.words = append(s.words, 0)
	}
	for i, w := range other.words {
		if w&^s.words[i] != 0 {
			changed = true
		}
		s.words[i] |= w
	}
	return changed
}

// RemoveAll is Set.removeAll; reports whether the set changed.
func (s *EnumSet[E]) RemoveAll(other *EnumSet[E]) bool {
	changed := false
	for i := range s.words {
		if i < len(other.words) {
			if s.words[i]&other.words[i] != 0 {
				changed = true
			}
			s.words[i] &^= other.words[i]
		}
	}
	return changed
}

// Size is Set.size.
func (s *EnumSet[E]) Size() int {
	if s == nil {
		return 0
	}
	n := 0
	for _, w := range s.words {
		n += bits.OnesCount64(w)
	}
	return n
}

// IsEmpty is Set.isEmpty.
func (s *EnumSet[E]) IsEmpty() bool { return s.Size() == 0 }

// Equals is Set.equals for two enum sets.
func (s *EnumSet[E]) Equals(other *EnumSet[E]) bool {
	n := max(len(s.words), len(other.words))
	for i := 0; i < n; i++ {
		var a, b uint64
		if i < len(s.words) {
			a = s.words[i]
		}
		if i < len(other.words) {
			b = other.words[i]
		}
		if a != b {
			return false
		}
	}
	return true
}

// HashCode is a hash consistent with Equals.
//
// Go-only: language-adapter — Java's EnumSet.hashCode() sums the elements'
// identity hash codes, which are not reproducible; nothing may depend on it.
func (s *EnumSet[E]) HashCode() int32 {
	var h int32
	for e := range s.All() {
		h += int32(e)*31 + 7
	}
	return h
}

// All iterates the elements in ordinal order.
func (s *EnumSet[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		if s == nil {
			return
		}
		for i, w := range s.words {
			for w != 0 {
				b := bits.TrailingZeros64(w)
				w &^= 1 << uint(b)
				if !yield(E(i*64 + b)) {
					return
				}
			}
		}
	}
}

// ToSlice returns the elements in ordinal order.
func (s *EnumSet[E]) ToSlice() []E {
	var out []E
	for e := range s.All() {
		out = append(out, e)
	}
	return out
}

// String renders like AbstractCollection.toString(): "[a, b]" using each
// element's Java toString().
func (s *EnumSet[E]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e := range s.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

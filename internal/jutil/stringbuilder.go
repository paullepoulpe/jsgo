package jutil

import "strconv"

// Go-only: language-adapter — java.lang.StringBuilder. Indices are UTF-16
// code units, as in Java. Use strings.Builder instead when the ported code
// only appends.
type StringBuilder struct {
	units []uint16
}

// NewStringBuilder is new StringBuilder().
func NewStringBuilder() *StringBuilder { return &StringBuilder{} }

// NewStringBuilderWithCapacity is new StringBuilder(int capacity).
func NewStringBuilderWithCapacity(capacity int) *StringBuilder {
	if capacity < 0 {
		panic(NewNegativeArraySizeException(strconv.Itoa(capacity)))
	}
	return &StringBuilder{units: make([]uint16, 0, capacity)}
}

// grow makes room for n more units with Java's AbstractStringBuilder
// policy (the capacity at least doubles), so that appending a long text
// piece by piece copies it a bounded number of times, as in Java.
func (b *StringBuilder) grow(n int) {
	if need := len(b.units) + n; need > cap(b.units) {
		newCap := 2*cap(b.units) + 2
		if newCap < need {
			newCap = need
		}
		units := make([]uint16, len(b.units), newCap)
		copy(units, b.units)
		b.units = units
	}
}

// NewStringBuilderFrom is new StringBuilder(String).
func NewStringBuilderFrom(s string) *StringBuilder {
	return &StringBuilder{units: StringToUTF16(s)}
}

// Append is StringBuilder.append(String) (and append(Object) after
// String.valueOf).
func (b *StringBuilder) Append(s string) *StringBuilder {
	if isASCII(s) {
		b.grow(len(s))
		for i := 0; i < len(s); i++ {
			b.units = append(b.units, uint16(s[i]))
		}
		return b
	}
	u := StringToUTF16(s)
	b.grow(len(u))
	b.units = append(b.units, u...)
	return b
}

// AppendChar is StringBuilder.append(char).
func (b *StringBuilder) AppendChar(c uint16) *StringBuilder {
	b.grow(1)
	b.units = append(b.units, c)
	return b
}

// AppendInt is StringBuilder.append(int) / append(long).
func (b *StringBuilder) AppendInt(i int64) *StringBuilder {
	return b.Append(strconv.FormatInt(i, 10))
}

// Length is StringBuilder.length().
func (b *StringBuilder) Length() int { return len(b.units) }

// CharAt is StringBuilder.charAt(int).
func (b *StringBuilder) CharAt(i int) uint16 {
	if i < 0 || i >= len(b.units) {
		panic(NewStringIndexOutOfBoundsException(i, len(b.units)))
	}
	return b.units[i]
}

// SetCharAt is StringBuilder.setCharAt(int, char).
func (b *StringBuilder) SetCharAt(i int, c uint16) {
	if i < 0 || i >= len(b.units) {
		panic(NewStringIndexOutOfBoundsException(i, len(b.units)))
	}
	b.units[i] = c
}

// SetLength is StringBuilder.setLength(int): truncates, or pads with '\0'.
func (b *StringBuilder) SetLength(n int) {
	if n < 0 {
		panic(NewStringIndexOutOfBoundsException(n, len(b.units)))
	}
	for len(b.units) < n {
		b.units = append(b.units, 0)
	}
	b.units = b.units[:n]
}

// Insert is StringBuilder.insert(int, String).
func (b *StringBuilder) Insert(offset int, s string) *StringBuilder {
	if offset < 0 || offset > len(b.units) {
		panic(NewStringIndexOutOfBoundsException(offset, len(b.units)))
	}
	u := StringToUTF16(s)
	n := len(b.units)
	// Shift the tail in place (amortized growth, as Java's
	// AbstractStringBuilder does) instead of copying the whole buffer.
	b.grow(len(u))
	b.units = append(b.units, u...)
	copy(b.units[offset+len(u):], b.units[offset:n])
	copy(b.units[offset:], u)
	return b
}

// InsertChar is StringBuilder.insert(int, char).
func (b *StringBuilder) InsertChar(offset int, c uint16) *StringBuilder {
	if offset < 0 || offset > len(b.units) {
		panic(NewStringIndexOutOfBoundsException(offset, len(b.units)))
	}
	n := len(b.units)
	b.grow(1)
	b.units = append(b.units, 0)
	copy(b.units[offset+1:], b.units[offset:n])
	b.units[offset] = c
	return b
}

// String is StringBuilder.toString().
func (b *StringBuilder) String() string { return StringFromUTF16(b.units) }

// CharSequence is java.lang.CharSequence (the methods ported code reads).
type CharSequence interface {
	Length() int
	CharAt(index int) uint16
}

// AppendCharSequenceRange is StringBuilder.append(CharSequence s, int start,
// int end): appends s.charAt(start) .. s.charAt(end - 1). s may be b itself.
func (b *StringBuilder) AppendCharSequenceRange(s CharSequence, start, end int) *StringBuilder {
	if start < 0 || start > end || end > s.Length() {
		panic(NewIndexOutOfBoundsExceptionForRange(start, end, s.Length()))
	}
	for i := start; i < end; i++ {
		b.units = append(b.units, s.CharAt(i))
	}
	return b
}

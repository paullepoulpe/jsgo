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

// NewStringBuilderFrom is new StringBuilder(String).
func NewStringBuilderFrom(s string) *StringBuilder {
	return &StringBuilder{units: StringToUTF16(s)}
}

// Append is StringBuilder.append(String) (and append(Object) after
// String.valueOf).
func (b *StringBuilder) Append(s string) *StringBuilder {
	if isASCII(s) {
		for i := 0; i < len(s); i++ {
			b.units = append(b.units, uint16(s[i]))
		}
		return b
	}
	b.units = append(b.units, StringToUTF16(s)...)
	return b
}

// AppendChar is StringBuilder.append(char).
func (b *StringBuilder) AppendChar(c uint16) *StringBuilder {
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
	b.units = append(b.units, 0)
	copy(b.units[offset+1:], b.units[offset:n])
	b.units[offset] = c
	return b
}

// String is StringBuilder.toString().
func (b *StringBuilder) String() string { return StringFromUTF16(b.units) }

package jutil

// Go-only: language-adapter — java.lang.CharSequence and
// StringBuilder.append(CharSequence, int, int).

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

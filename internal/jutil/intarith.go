package jutil

// Go-only: language-adapter — Java int arithmetic on values the port keeps
// in Go int (64-bit): Java's + and - wrap at 32 bits. Ported code uses these
// where the operands can come from input and overflow is possible, so that
// sums and differences of Java-int-range values match Java exactly.

// IntAdd is Java's a + b on ints.
func IntAdd(a, b int) int { return int(int32(a) + int32(b)) }

// IntSub is Java's a - b on ints.
func IntSub(a, b int) int { return int(int32(a) - int32(b)) }

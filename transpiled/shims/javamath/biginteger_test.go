package javamath

import (
	"bytes"
	"math"
	"testing"
)

func TestBigIntegerOperations(t *testing.T) {
	a := NewBigInteger("-17")
	b := NewBigInteger("5")

	// Java truncated division & remainder: -17 / 5 = -3, -17 % 5 = -2
	if got := a.Divide(b).String(); got != "-3" {
		t.Fatalf("Divide(-17, 5) = %q, want -3", got)
	}
	if got := a.Remainder(b).String(); got != "-2" {
		t.Fatalf("Remainder(-17, 5) = %q, want -2", got)
	}
	if got := a.Mod(b).String(); got != "3" {
		t.Fatalf("Mod(-17, 5) = %q, want 3", got)
	}

	qr := a.DivideAndRemainder(b)
	if qr[0].String() != "-3" || qr[1].String() != "-2" {
		t.Fatalf("DivideAndRemainder = [%s, %s], want [-3, -2]", qr[0], qr[1])
	}

	// Byte array round-trip (two's complement)
	for _, s := range []string{"0", "1", "127", "128", "-1", "-128", "-129", "123456789012345678901234567890"} {
		x := NewBigInteger(s)
		raw := x.ToByteArray()
		rt := NewBigIntegerBytes(raw)
		if !x.Equals(rt) {
			t.Fatalf("round-trip failed for %s: bytes=%x -> %s", s, raw, rt)
		}
	}

	if !bytes.Equal(NewBigInteger("128").ToByteArray(), []byte{0x00, 0x80}) {
		t.Fatalf("128.ToByteArray() = %x, want 0080", NewBigInteger("128").ToByteArray())
	}
	if !bytes.Equal(NewBigInteger("-128").ToByteArray(), []byte{0x80}) {
		t.Fatalf("-128.ToByteArray() = %x, want 80", NewBigInteger("-128").ToByteArray())
	}

	// Pow & IntValueExact & DoubleValue
	p := BigInteger_ValueOf(5).Pow(3)
	if p.IntValueExact() != 125 {
		t.Fatalf("5^3 IntValueExact = %d, want 125", p.IntValueExact())
	}
	huge := BigInteger_TWO.Pow(1024)
	if !math.IsInf(huge.DoubleValue(), 1) {
		t.Fatalf("2^1024 DoubleValue = %v, want +Inf", huge.DoubleValue())
	}
}

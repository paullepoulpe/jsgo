package jutil

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type ptrHashableKey struct{ id int }

func (k *ptrHashableKey) Equals(o any) bool {
	if k == o {
		return true
	}
	x, ok := o.(*ptrHashableKey)
	return ok && x != nil && x.id == k.id
}
func (k *ptrHashableKey) HashCode() int32 { return int32(k.id*31 + 7) }

type testNodeID uint32

func (n testNodeID) IsNil() bool       { return n == 0 }
func (n testNodeID) HashCode() int32   { return int32(uint32(n) * 0x9E3779B1) }
func (n testNodeID) Equals(o any) bool { m, ok := o.(testNodeID); return ok && m == n }

func TestCheckNotNullZeroAllocs(t *testing.T) {
	k1 := &ptrHashableKey{42}
	node := testNodeID(7)

	if n := testing.AllocsPerRun(100, func() {
		_ = CheckNotNull("abc")
		_ = CheckNotNull(123)
		_ = CheckNotNull(k1)
		_ = CheckNotNull(node)
	}); n != 0 {
		t.Fatalf("expected 0 allocations, got %v", n)
	}
}

func TestTryParseHelpersMatchJava(t *testing.T) {
	for _, r := range readTSV(t, "parse_double.tsv") {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		got, ok := TryParseDouble(in)
		wantErr := strings.HasPrefix(r[1], "E:")
		if ok == wantErr {
			t.Fatalf("TryParseDouble(%q) ok=%v, Java=%s", in, ok, r[1])
		}
		if ok {
			if gotHex := strconv.FormatUint(math.Float64bits(got), 16); gotHex != r[1] {
				t.Fatalf("TryParseDouble(%q) = %s, Java %s", in, gotHex, r[1])
			}
		}
	}
	for _, r := range readTSV(t, "bigint_radix.tsv") {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		radix, _ := strconv.Atoi(r[1])
		got, ok := TryParseBigInt(in, radix)
		wantErr := strings.HasPrefix(r[2], "E:")
		if ok == wantErr {
			t.Fatalf("TryParseBigInt(%q, %d) ok=%v, Java=%s", in, radix, ok, r[2])
		}
		if ok && got.String() != r[2] {
			t.Fatalf("TryParseBigInt(%q, %d) = %s, Java %s", in, radix, got, r[2])
		}
	}
	if _, ok := TryParseLong("123", 1); ok {
		t.Fatal("radix 1 should fail")
	}
	if _, ok := TryParseLong("123", 37); ok {
		t.Fatal("radix 37 should fail")
	}
	if v, ok := TryParseLong("-9223372036854775808", 10); !ok || v != math.MinInt64 {
		t.Fatalf("MinInt64: got %d, %v", v, ok)
	}
	if _, ok := TryParseLong("9223372036854775808", 10); ok {
		t.Fatal("MaxInt64+1 should fail")
	}
}

// IdentityHashCodeOfPointer gives IdentityHashCode's value without
// allocating.
func TestIdentityHashCodeOfPointer(t *testing.T) {
	type node struct{ a, b int }
	ns := make([]*node, 64)
	for i := range ns {
		ns[i] = &node{a: i}
		if got, want := IdentityHashCodeOfPointer(ns[i]), IdentityHashCode(ns[i]); got != want {
			t.Fatalf("IdentityHashCodeOfPointer = %d, IdentityHashCode = %d", got, want)
		}
	}
	var none *node
	if got, want := IdentityHashCodeOfPointer(none), IdentityHashCode(none); got != want || got != 0 {
		t.Fatalf("nil: IdentityHashCodeOfPointer = %d, IdentityHashCode = %d", got, want)
	}
	if n := testing.AllocsPerRun(100, func() { _ = IdentityHashCodeOfPointer(ns[3]) }); n != 0 {
		t.Fatalf("allocated %v times", n)
	}
}

func TestPathElementOffsetsConcurrent(t *testing.T) {
	p := Path_Of("alpha/beta/gamma/delta.js")
	other := Path_Of("alpha/beta/other.js")
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := p.GetNameCount(); got != 4 {
				t.Errorf("GetNameCount = %d, want 4", got)
			}
			if got := p.GetName(2).String(); got != "gamma" {
				t.Errorf("GetName(2) = %q, want gamma", got)
			}
			if got := p.Relativize(other).String(); got != "../../other.js" {
				t.Errorf("Relativize = %q, want ../../other.js", got)
			}
		}()
	}
	wg.Wait()
}

// Port of test/com/google/javascript/jscomp/regex/CharRangesTest.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package regex

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: Long.parseLong(System.getProperty("junit.random.seed", "" + System.currentTimeMillis()))
// Go-only: language-adapter — the seed comes from $JUNIT_RANDOM_SEED, and
// math/rand stands in for java.util.Random (the test only needs random bits).
var charRangesTestSeed = func() int64 {
	if s := os.Getenv("JUNIT_RANDOM_SEED"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v
		}
	}
	return time.Now().UnixMilli()
}()

func TestCharRangesAgainstRegularImplementation(t *testing.T) {
	rnd := rand.New(rand.NewSource(charRangesTestSeed))

	for run := 10; run > 0; run-- {
		// Fill with bits in the range [0x1000, 0x3000).
		bs := make([]bool, 0x5000) // Java: BitSet
		for i := 0x1000; i > 0; i-- {
			bs[0x1000+rnd.Intn(0x3000)] = true
		}

		// Create an equivalent sparse bit set
		var members []int
		for i, set := range bs {
			if set {
				members = append(members, i)
			}
		}
		sbs := CharRanges_WithMembers(members...)

		// Check all bits including past the min/max bit
		for i := 0; i < 0x5000; i++ {
			if bs[i] != sbs.Contains(i) {
				t.Fatalf("seed=%d sbs=%s, difference at bit %d", charRangesTestSeed, sbs, i)
			}
		}
	}
}

func TestEmptyCharRanges(t *testing.T) {
	sbs := CharRanges_EMPTY
	for i := -1000; i < 1000; i++ {
		if sbs.Contains(i) {
			t.Fatalf("EMPTY contains %d", i)
		}
	}
	if got := sbs.String(); got != "[]" {
		t.Fatalf("EMPTY.toString() = %q", got)
	}
}

func TestCharRangesFactories(t *testing.T) {
	isbs := CharRanges_WithMembers(0, 1, 4, 9)
	isbs2 := CharRanges_WithMembers(0, 1, 4, 9)
	if got := isbs.String(); got != "[0x0-0x1 0x4 0x9]" {
		t.Fatalf("toString = %q", got)
	}

	esbs := CharRanges_WithMembers()

	if !isbs2.Equals(isbs) {
		t.Error("isbs2 != isbs")
	}
	if isbs.Equals(esbs) {
		t.Error("isbs == esbs")
	}
	if isbs == nil {
		t.Error("isbs is null")
	}
	if isbs.Equals(struct{}{}) {
		t.Error("isbs == new Object()")
	}

	if isbs2.HashCode() != isbs.HashCode() {
		t.Error("hash codes differ")
	}
	if isbs.HashCode() == esbs.HashCode() {
		t.Error("hash codes equal")
	}
}

func expectIllegalArgument(t *testing.T, what string, f func()) {
	t.Helper()
	caught := false
	jutil.Try(f, func(ex jutil.IIllegalArgumentException) { caught = true })
	if !caught {
		t.Errorf("%s: no IllegalArgumentException", what)
	}
}

func TestRangeConstructor(t *testing.T) {
	expectIllegalArgument(t, "Mismatched ranges", func() { CharRanges_WithRanges(1) })
	expectIllegalArgument(t, "Discontiguous ranges", func() { CharRanges_WithRanges(1, 4, 4, 5) })
	expectIllegalArgument(t, "Misordered ranges", func() { CharRanges_WithRanges(4, 5, 1, 3) })
	expectIllegalArgument(t, "Empty range", func() { CharRanges_WithRanges(0, 0) })
}

func TestDupeMembers(t *testing.T) {
	sbs1 := CharRanges_WithMembers(0, 1, 4, 9)
	if got := sbs1.String(); got != "[0x0-0x1 0x4 0x9]" {
		t.Fatalf("%s", got)
	}

	sbs2 := CharRanges_WithMembers(9, 1, 4, 1, 0)
	if got := sbs2.String(); got != "[0x0-0x1 0x4 0x9]" {
		t.Fatalf("%s", got)
	}

	// Java: new EqualsTester().addEqualityGroup(sbs1, sbs2).testEquals();
	for _, a := range []*CharRanges{sbs1, sbs2} {
		if a.Equals(nil) || a.Equals("not a CharRanges") {
			t.Error("equal to null or an unrelated object")
		}
		for _, b := range []*CharRanges{sbs1, sbs2} {
			if !a.Equals(b) || a.HashCode() != b.HashCode() {
				t.Errorf("%s and %s not equal", a, b)
			}
		}
	}

	for i := -10; i < 20; i++ {
		if sbs2.Contains(i) != sbs1.Contains(i) {
			t.Errorf("%d", i)
		}
	}
}

func assertString(t *testing.T, c *CharRanges, want string) {
	t.Helper()
	if got := c.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertBool(t *testing.T, what string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}

func TestDifference(t *testing.T) {
	//                     1               2               3
	//     0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0
	// b-a  DD         DD DDD        D      DDD
	// a      AAAAAAAAA      A A A A   A AAA   AAA A A
	// b    BBB  BBB  BBB BBB        B B    BBB
	// a-b     DD   DD       D D D D     DDD   DDD D D
	a := CharRanges_WithRanges(
		0x03, 0x0C, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1C, 0x1D, 0x1E, 0x21,
		0x24, 0x27, 0x28, 0x29, 0x2A, 0x2B)
	b := CharRanges_WithRanges(
		0x01, 0x04, 0x06, 0x09, 0x0B, 0x0E, 0x0F, 0x12, 0x1A, 0x1B, 0x1C, 0x1D, 0x21, 0x24)
	empty := CharRanges_WithMembers()

	assertBool(t, "empty.union(empty) == empty", empty.Union(empty).Equals(empty), true)
	assertBool(t, "a.union(empty) == a", a.Union(empty).Equals(a), true)
	assertBool(t, "empty.union(b) == b", empty.Union(b).Equals(b), true)

	aSb := a.Difference(b)
	assertString(t, aSb, "[0x4-0x5 0x9-0xa 0x12 0x14 0x16 0x18 0x1e-0x20 0x24-0x26 0x28 0x2a]")
	assertBool(t, "a.containsAll(aSb)", a.ContainsAll(aSb), true)
	assertBool(t, "aSb.containsAll(a)", aSb.ContainsAll(a), false)
	assertBool(t, "aSb.containsAll(b)", aSb.ContainsAll(b), false)

	bSa := b.Difference(a)
	assertString(t, bSa, "[0x1-0x2 0xc-0xd 0xf-0x11 0x1a 0x21-0x23]")
	assertBool(t, "b.containsAll(bSa)", b.ContainsAll(bSa), true)
	assertBool(t, "bSa.containsAll(a)", bSa.ContainsAll(a), false)
	assertBool(t, "bSa.containsAll(b)", bSa.ContainsAll(b), false)

	// Check that a and b not changed by operation
	assertString(t, a, "[0x3-0xb 0x12 0x14 0x16 0x18 0x1c 0x1e-0x20 0x24-0x26 0x28 0x2a]")
	assertString(t, b, "[0x1-0x3 0x6-0x8 0xb-0xd 0xf-0x11 0x1a 0x1c 0x21-0x23]")

	//    0 1 2 3 4 5 6 7 8 9 a b c d e f
	// m: * * * *     *     * *       * *
	// s:     *     * * *     * *   * *
	// d: * *   *           *           *
	m := CharRanges_WithMembers(0, 1, 2, 3, 6, 9, 0xa, 0xe, 0xf)
	s := CharRanges_WithMembers(2, 5, 6, 7, 0xa, 0xb, 0xd, 0xe)
	d := m.Difference(s)
	assertString(t, d, "[0x0-0x1 0x3 0x9 0xf]")
	assertBool(t, "m.containsAll(d)", m.ContainsAll(d), true)
	assertBool(t, "d.containsAll(m)", d.ContainsAll(m), false)
	assertBool(t, "d.containsAll(s)", d.ContainsAll(s), false)
	assertBool(t, "s.containsAll(d)", s.ContainsAll(d), false)
	assertBool(t, "d.containsAll(d)", d.ContainsAll(d), true)
}

func TestUnion(t *testing.T) {
	//                 1               2               3
	// 0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0
	//    AAAAAAAAA      A A A A   A AAA   AAA A A
	//  BBB  BBB  BBB BBB        B B    BBB
	//  UUUUUUUUUUUUU UUUU U U U U U UUUUUUUUU U U
	a := CharRanges_WithRanges(
		0x03, 0x0C, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1C, 0x1D, 0x1E, 0x21,
		0x24, 0x27, 0x28, 0x29, 0x2A, 0x2B)
	b := CharRanges_WithRanges(
		0x01, 0x04, 0x06, 0x09, 0x0B, 0x0E, 0x0F, 0x12, 0x1A, 0x1B, 0x1C, 0x1D, 0x21, 0x24)
	empty := CharRanges_WithMembers()

	assertBool(t, "empty.union(empty) == empty", empty.Union(empty).Equals(empty), true)
	assertBool(t, "a.union(empty) == a", a.Union(empty).Equals(a), true)
	assertBool(t, "empty.union(b) == b", empty.Union(b).Equals(b), true)

	aUb := a.Union(b)
	assertString(t, aUb, "[0x1-0xd 0xf-0x12 0x14 0x16 0x18 0x1a 0x1c 0x1e-0x26 0x28 0x2a]")
	assertBool(t, "b.union(a) == aUb", b.Union(a).Equals(aUb), true)
	assertBool(t, "aUb.containsAll(a)", aUb.ContainsAll(a), true)
	assertBool(t, "aUb.containsAll(b)", aUb.ContainsAll(b), true)
	assertBool(t, "a.containsAll(b)", a.ContainsAll(b), false)
	assertBool(t, "b.containsAll(a)", b.ContainsAll(a), false)
	assertBool(t, "a.containsAll(a)", a.ContainsAll(a), true)
	assertBool(t, "b.containsAll(b)", b.ContainsAll(b), true)
	assertBool(t, "aUb.containsAll(aUb)", aUb.ContainsAll(aUb), true)

	// Check that a and b not changed by operation
	assertString(t, a, "[0x3-0xb 0x12 0x14 0x16 0x18 0x1c 0x1e-0x20 0x24-0x26 0x28 0x2a]")
	assertString(t, b, "[0x1-0x3 0x6-0x8 0xb-0xd 0xf-0x11 0x1a 0x1c 0x21-0x23]")
}

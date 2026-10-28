// Port of com/google/javascript/jscomp/regex/CharRanges.java @ closure-compiler v20261021

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
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CharRanges is an immutable sparse bitset that deals well where the data is
// chunky: where P(bit[x+1] == bit[x]). E.g. [101,102,103,104,105,1001,1002,
// 1003,1004] is chunky.
//
// Java ints are Go ints here; the few sums that can wrap a Java int are
// computed in int32 (Go-only: language-adapter — Java int overflow).
type CharRanges struct {
	// A strictly increasing set of bit indices where even members are the
	// inclusive starts of ranges, and odd members are the exclusive ends.
	//
	// E.g., { 1, 5, 6, 10 } represents the set ( 1, 2, 3, 4, 6, 7, 8, 9 ).
	ranges []int
}

// CharRanges_EMPTY is the empty set.
var CharRanges_EMPTY = NewCharRanges([]int{})

// CharRanges_ALL_CODE_UNITS is the set of all UTF-16 code units.
var CharRanges_ALL_CODE_UNITS = NewCharRanges([]int{0, 0x10000})

// CharRanges_Inclusive returns the range [start, end].
//
// Java: CharRanges.inclusive(int, int)
func CharRanges_Inclusive(start, end int) *CharRanges {
	if start > end {
		panic(jutil.NewIndexOutOfBoundsException(strconv.Itoa(start) + " > " + strconv.Itoa(end)))
	}
	return NewCharRanges([]int{start, javaIntAdd(end, 1)})
}

// CharRanges_WithMembers returns an instance containing all and only the
// given members.
//
// Java: CharRanges.withMembers(int...)
func CharRanges_WithMembers(members ...int) *CharRanges {
	return NewCharRanges(CharRanges_IntArrayToRanges(members))
}

// CharRanges_WithRanges returns an instance containing the given ranges.
//
// ranges is an even-length ordered sequence of non-overlapping,
// non-contiguous, [inclusive start, exclusive end) ranges.
//
// Java: CharRanges.withRanges(int...)
func CharRanges_WithRanges(ranges ...int) *CharRanges {
	if (len(ranges) & 1) != 0 {
		panic(jutil.NewIllegalArgumentException())
	}
	for i := 1; i < len(ranges); i++ {
		if ranges[i] <= ranges[i-1] {
			panic(jutil.NewIllegalArgumentExceptionMsg(strconv.Itoa(ranges[i]) + " > " + strconv.Itoa(ranges[i-1])))
		}
	}
	return NewCharRanges(ranges)
}

// NewCharRanges wraps a ranges array.
//
// Java: CharRanges(int[])
func NewCharRanges(ranges []int) *CharRanges {
	c := &CharRanges{}
	c.ranges = ranges
	return c
}

// CharRanges_IntArrayToRanges converts a set of members to ranges. It sorts
// members in place, like Java's Arrays.sort.
//
// Java: CharRanges.intArrayToRanges(int[])
func CharRanges_IntArrayToRanges(members []int) []int {
	nMembers := len(members)
	if nMembers == 0 {
		return []int{}
	}

	sortInts(members)

	// Count the number of runs.
	nRuns := 1
	for i := 1; i < nMembers; i++ {
		current := members[i]
		last := members[i-1]
		if current == last {
			continue
		}
		if current != javaIntAdd(last, 1) {
			nRuns++
		}
	}

	ranges := make([]int, nRuns*2)
	ranges[0] = members[0]
	k := 0
	for i := 1; k+2 < len(ranges); i++ {
		current := members[i]
		last := members[i-1]
		if current == last {
			continue
		}
		if current != javaIntAdd(last, 1) {
			k++
			ranges[k] = javaIntAdd(last, 1) // add 1 to make end exclusive
			k++
			ranges[k] = current
		}
	}
	k++
	ranges[k] = javaIntAdd(members[nMembers-1], 1) // add 1 to make end exclusive
	return ranges
}

// Contains reports whether bit is in the set.
//
// Java: CharRanges.contains(int)
func (c *CharRanges) Contains(bit int) bool {
	return (jutil.ArraysBinarySearchInts(c.ranges, bit) & 1) == 0
	// By the contract of Arrays.binarySearch, its result is either the position
	// of bit in ranges or it is the bitwise inverse of the position of the
	// least element greater than bit.

	// Two cases
	// case (idx >= 0)
	//     We ended up exactly on a range boundary.
	//     Starts are inclusive and ends are both exclusive, so this contains
	//     bit iff idx is even.
	//
	// case (idx < 0)
	//     If the least element greater than bit is an odd element,
	//     then bit must be greater than a start and less than an end, so
	//     contained.
	//
	//     If bit is greater than all elements, then idx will be past the end of
	//     the array, and will be even since ranges.length is even.
	//
	//     Otherwise, bit must be in the space between two runs, so not
	//     contained.
	//
	//     In all cases, oddness is equivalent to containedness.

	// Those two cases lead to
	//     idx >= 0 ? ((idx & 1) == 0) : ((~idx & 1) == 1)

	// But ~n & bit == bit   <=>   n & bit == 0, so
	//     idx >= 0 ? ((idx & 1) == 0) : ((~idx & 1) == 1)
	// =>  idx >= 0 ? ((idx & 1) == 0) : ((idx & 1) == 0)
	// =>  (idx & 1) == 0
}

// IsEmpty reports whether the set is empty.
//
// Java: CharRanges.isEmpty()
func (c *CharRanges) IsEmpty() bool {
	return len(c.ranges) == 0
}

// GetNumRanges returns the number of ranges.
//
// Java: CharRanges.getNumRanges()
func (c *CharRanges) GetNumRanges() int { return len(c.ranges) >> 1 }

// Start returns the inclusive start of the i-th range.
//
// Java: CharRanges.start(int)
func (c *CharRanges) Start(i int) int { return c.ranges[i<<1] }

// End returns the exclusive end of the i-th range.
//
// Java: CharRanges.end(int)
func (c *CharRanges) End(i int) int { return c.ranges[(i<<1)|1] }

// Union returns the union of c and other.
//
// Java: CharRanges.union(CharRanges)
func (c *CharRanges) Union(other *CharRanges) *CharRanges {
	// Index of the input ranges
	q := c.ranges
	r := other.ranges
	// Lengths of the inputs
	m := len(q)
	n := len(r)

	if m == 0 {
		return other
	}
	if n == 0 {
		return c
	}

	// The output array.  The length is m+n in the worst case when all the
	// ranges in a are disjoint from the ranges in b.
	out := make([]int, m+n)

	// Indexes into the various arrays
	i := 0
	j := 0
	k := 0
	// Since there are three arrays, and indices into them the following
	// should never occur in this function:
	// (1) q[j] or q[k]                         -- q is indexed by i
	// (2) r[i] or r[k]                         -- r is indexed by j
	// (3) out[i] or out[j]                     -- out is indexed by k
	// (4) i < n or j < m                       -- index compared to wrong limit

	// This loop exits because we always increment at least one of i,j.
	for i < m && j < n {
		// Range starts and ends.
		a0 := q[i]
		a1 := q[i+1]
		b0 := r[j]
		b1 := r[j+1]
		if a1 < b0 { // [a0, a1) ends before [b0, b1) starts
			out[k] = a0
			k++
			out[k] = a1
			k++
			i += 2
		} else if b1 < a0 { // [b0, b1) ends before [a0, a1) starts
			out[k] = b0
			k++
			out[k] = b1
			k++
			j += 2
		} else { // ranges overlap
			// We need to compute a new range based on the set of ranges that
			// transitively overlap.
			//       AAAAAAAAA AAA
			//     BBB  BBB* BBB
			// In the range above, the start comes from one set, and the end from
			// another.  The range with the asterisk next to it is subsumed entirely
			// by a range from the other, and so not all ranges on the input
			// contribute a value to the output.
			// The last BBB run serves only as a bridge -- it overlaps two
			// disjoint ranges in the other one so establishes that they
			// transitively overlap.
			start := min(a0, b0)
			// Guess at the end, and lookahead to come up with a more complete
			// estimate.
			end := max(a1, b1)
			i += 2
			j += 2
			for i < m || j < n {
				if i < m && q[i] <= end {
					end = max(end, q[i+1])
					i += 2
				} else if j < n && r[j] <= end {
					end = max(end, r[j+1])
					j += 2
				} else {
					break
				}
			}
			out[k] = start
			k++
			out[k] = end
			k++
		}
	}
	// There may be unprocessed ranges at the end of one of the inputs.
	if i < m {
		copy(out[k:], q[i:m])
		k += m - i
	} else if j < n {
		copy(out[k:], r[j:n])
		k += n - j
	}
	// We guessed at the output length above.  Cut off the tail.
	return NewCharRanges(out[:k])
}

// Intersection returns the intersection of c and other.
//
// Java: CharRanges.intersection(CharRanges)
func (c *CharRanges) Intersection(other *CharRanges) *CharRanges {
	aRanges := c.ranges
	bRanges := other.ranges
	aLen := len(aRanges)
	bLen := len(bRanges)
	if aLen == 0 {
		return c
	}
	if bLen == 0 {
		return other
	}
	aIdx := 0
	bIdx := 0
	intersection := make([]int, min(aLen, bLen))
	intersectionIdx := 0
	pos := min(aRanges[0], bRanges[0])
	for aIdx < aLen && bIdx < bLen {
		if aRanges[aIdx+1] <= pos {
			aIdx += 2
		} else if bRanges[bIdx+1] <= pos {
			bIdx += 2
		} else {
			start := max(aRanges[aIdx], bRanges[bIdx])
			if pos < start { // Advance to start of common block.
				pos = start
			} else {
				// Now we know that pos is less than the ends of the two ranges and
				// greater or equal to the starts of the two ranges.
				end := min(aRanges[aIdx+1], bRanges[bIdx+1])
				if intersectionIdx != 0 && pos == intersection[intersectionIdx-1] {
					intersection[intersectionIdx-1] = end
				} else {
					if intersectionIdx == len(intersection) {
						newArr := make([]int, intersectionIdx*2)
						copy(newArr, intersection[:intersectionIdx])
						intersection = newArr
					}
					intersection[intersectionIdx] = pos
					intersectionIdx++
					intersection[intersectionIdx] = end
					intersectionIdx++
				}
				pos = end
			}
		}
	}
	return NewCharRanges(intersection[:intersectionIdx])
}

// Difference returns c minus subtrahendRanges.
//
// Java: CharRanges.difference(CharRanges)
func (c *CharRanges) Difference(subtrahendRanges *CharRanges) *CharRanges {
	// difference = minuend - subtrahend
	minuend := c.ranges
	subtrahend := subtrahendRanges.ranges

	mn := len(minuend)
	sn := len(subtrahend)
	if mn == 0 || sn == 0 {
		return c
	}

	difference := make([]int, len(minuend))

	// Indices into minuend.ranges, subtrahend.ranges, and difference.
	mIdx := 0
	sIdx := 0
	dIdx := 0

	pos := minuend[0]
	for mIdx < mn {
		if pos >= minuend[mIdx+1] {
			mIdx += 2
		} else if pos < minuend[mIdx] {
			// Skip gaps in the minuend.
			pos = minuend[mIdx]
		} else if sIdx < sn && pos >= subtrahend[sIdx] {
			// Skip over a removed part.
			pos = subtrahend[sIdx+1]
			sIdx += 2
		} else {
			// Now we know that pos is between [minuend[i], minuend[i + 1])
			// and outside [subtrahend[j], subtrahend[j + 1]).
			var end int
			if sIdx < sn {
				end = min(minuend[mIdx+1], subtrahend[sIdx])
			} else {
				end = minuend[mIdx+1]
			}
			if dIdx != 0 && difference[dIdx-1] == pos {
				difference[dIdx-1] = pos
			} else {
				if dIdx == len(difference) {
					newArr := make([]int, dIdx*2)
					copy(newArr, difference[:dIdx])
					difference = newArr
				}
				difference[dIdx] = pos
				dIdx++
				difference[dIdx] = end
				dIdx++
			}
			pos = end
		}
	}

	return NewCharRanges(difference[:dIdx])
}

// ContainsAll reports whether every member of sub is in c.
//
// Java: CharRanges.containsAll(CharRanges)
func (c *CharRanges) ContainsAll(sub *CharRanges) bool {
	superRanges := c.ranges
	subRanges := sub.ranges

	superIdx := 0
	subIdx := 0
	superLen := len(superRanges)
	subLen := len(subRanges)
	for subIdx < subLen {
		if superIdx == superLen {
			return false
		}
		if superRanges[superIdx+1] <= subRanges[subIdx] {
			// Super range ends before subRange starts.
			superIdx += 2
		} else if superRanges[superIdx] > subRanges[subIdx] {
			// Uncontained portion at start of sub range.
			return false
		} else if superRanges[superIdx+1] >= subRanges[subIdx+1] {
			// A sub range is completely contained in the super range.
			// We know this because of the above condition and we have already
			// ruled out that subRanges[subIdx] < superRanges[superIdx].
			subIdx += 2
		} else {
			// Uncontained portion at end of sub range.
			return false
		}
	}
	return subIdx == subLen
}

// Shift shifts the bits matched by the given delta. So if this has the bits
// (a, b, c, ..., z) set then the result has the bits ((a - delta),
// (b - delta), (c - delta), ...., (z - delta)) set.
//
// It throws IndexOutOfBoundsException if shifting by delta would cause an
// overflow or underflow in a 32 bit signed int range boundary. Since the end
// boundaries of ranges are exclusive, even if there is no range containing
// Integer.MAX_VALUE, shifting by a delta of 1 can cause an overflow.
//
// Java: CharRanges.shift(int)
func (c *CharRanges) Shift(delta int) *CharRanges {
	n := len(c.ranges)
	if delta == 0 || n == 0 {
		return c
	}
	// Test overflow/underflow
	if delta < 0 {
		lmin := int64(javaIntAdd(c.ranges[0], delta))
		if lmin < -1<<31 {
			panic(jutil.NewIndexOutOfBoundsExceptionNoMessage())
		}
	} else {
		lmax := int64(javaIntAdd(c.ranges[n-1], delta))
		if lmax > 1<<31-1 {
			panic(jutil.NewIndexOutOfBoundsExceptionNoMessage())
		}
	}
	// Create a shifted range.
	shiftedRanges := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		shiftedRanges[i] = javaIntAdd(c.ranges[i], delta)
	}
	return NewCharRanges(shiftedRanges)
}

// String renders the ranges, e.g. "[0x0-0x1 0x4 0x9]".
//
// Java: CharRanges.toString()
func (c *CharRanges) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i := 0; i < len(c.ranges); i++ {
		if (i&1) != 0 && c.ranges[i] == javaIntAdd(c.ranges[i-1], 1) {
			continue
		}
		if i != 0 {
			if (i & 1) == 0 {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('-')
			}
		}
		sb.WriteString("0x")
		sb.WriteString(strconv.FormatInt(int64(javaIntAdd(c.ranges[i], -(i&1))), 16))
	}
	sb.WriteByte(']')
	return sb.String()
}

// Equals is value equality on the ranges array.
//
// Java: CharRanges.equals(Object)
func (c *CharRanges) Equals(o any) bool {
	charRanges, ok := o.(*CharRanges)
	return ok && charRanges != nil && slices.Equal(c.ranges, charRanges.ranges)
}

// HashCode hashes the first 16 boundaries.
//
// Java: CharRanges.hashCode()
func (c *CharRanges) HashCode() int32 {
	var hc int32
	for i, n := 0, min(16, len(c.ranges)); i < n; i++ {
		hc = (hc << 2) + int32(c.ranges[i])
	}
	return hc
}

// javaIntAdd is Java int addition (wraps at 32 bits).
//
// Go-only: language-adapter — Java int overflow semantics on Go ints.
func javaIntAdd(a, b int) int {
	return int(int32(a) + int32(b))
}

// sortInts sorts in ascending order, in place (Java: Arrays.sort(int[])).
//
// Go-only: language-adapter — java.util.Arrays.sort(int[]).
func sortInts(a []int) {
	slices.Sort(a)
}

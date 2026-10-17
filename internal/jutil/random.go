package jutil

import (
	"sync/atomic"
	"time"
)

// Go-only: language-adapter — java.util.Random (the 48-bit linear
// congruential generator, bit-for-bit) and Collections.shuffle.

const (
	randomMultiplier int64 = 0x5DEECE66D
	randomAddend     int64 = 0xB
	randomMask       int64 = (1 << 48) - 1
)

// Random is java.util.Random.
type Random struct {
	seed int64
}

// randomInstances counts the generators NewRandom has created.
var randomInstances atomic.Int64

// NewRandom is new Random(). The Javadoc only asks for a seed that is very
// likely to differ from that of any other invocation; the JDK's seed is
// unspecified and time-dependent. Here it is the wall clock in nanoseconds
// plus the instance count times a large prime, so that generators created
// within one clock tick get different seeds too.
func NewRandom() *Random {
	return NewRandomWithSeed(time.Now().UnixNano() + randomInstances.Add(1)*1_000_003)
}

// NewRandomWithSeed is new Random(long seed).
func NewRandomWithSeed(seed int64) *Random {
	return &Random{seed: (seed ^ randomMultiplier) & randomMask}
}

// SetSeed is Random.setSeed(long).
func (r *Random) SetSeed(seed int64) {
	r.seed = (seed ^ randomMultiplier) & randomMask
}

// Next is Random.next(int bits).
func (r *Random) Next(bits int) int32 {
	r.seed = (r.seed*randomMultiplier + randomAddend) & randomMask
	return int32(uint64(r.seed) >> uint(48-bits))
}

// NextInt is Random.nextInt().
func (r *Random) NextInt() int32 { return r.Next(32) }

// NextIntBound is Random.nextInt(int bound): a value in [0, bound), using
// the algorithm that the Javadoc of Random.nextInt(int) gives. A power-of-two
// bound takes the high bits of next(31); any other bound takes next(31)
// modulo bound, drawing again whenever the 31-bit sample falls in the last,
// incomplete stretch of bound values (where the result would be biased).
func (r *Random) NextIntBound(bound int32) int32 {
	if bound <= 0 {
		panic(NewIllegalArgumentExceptionMsg("bound must be positive"))
	}
	if bound&-bound == bound {
		return int32((int64(bound) * int64(r.Next(31))) >> 31)
	}
	for {
		sample := r.Next(31)
		v := sample % bound
		// sample - v is the start of v's stretch; it is complete only if
		// all bound values from there fit below 2^31 (no int overflow).
		if int64(sample)-int64(v)+int64(bound-1) <= 0x7FFFFFFF {
			return v
		}
	}
}

// NextLong is Random.nextLong().
func (r *Random) NextLong() int64 {
	hi := int64(r.Next(32))
	lo := int64(r.Next(32))
	return (hi << 32) + lo
}

// NextBoolean is Random.nextBoolean().
func (r *Random) NextBoolean() bool { return r.Next(1) != 0 }

// NextDouble is Random.nextDouble().
func (r *Random) NextDouble() float64 {
	hi := int64(r.Next(26))
	lo := int64(r.Next(27))
	return float64((hi<<27)+lo) * (1.0 / float64(int64(1)<<53))
}

// NextFloat is Random.nextFloat().
func (r *Random) NextFloat() float32 {
	return float32(r.Next(24)) / float32(1<<24)
}

// Collections_Shuffle is Collections.shuffle(List, Random) on a
// random-access list.
func Collections_Shuffle[T any](list []T, rnd *Random) {
	for i := len(list); i > 1; i-- {
		j := rnd.NextIntBound(int32(i))
		list[i-1], list[j] = list[j], list[i-1]
	}
}

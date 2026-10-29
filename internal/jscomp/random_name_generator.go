// Port of com/google/javascript/jscomp/RandomNameGenerator.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Derived from com/google/common/hash/Murmur3_128HashFunction.java
// (Guava 33.4.6-jre): murmur3x64_128 and its helpers follow Guava's
// Murmur3_128Hasher (mixK1, mixK2, fmix64, and processRemaining's mixing of
// both halves of a partial block). MurmurHash3 itself was written by Austin
// Appleby, who placed it in the public domain.
//
// Copyright (C) 2011 The Guava Authors
// Licensed under the Apache License, Version 2.0; see LICENSES/Apache-2.0.txt.

package jscomp

import (
	"encoding/binary"
	"math/bits"
	"strings"
	"sync/atomic"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RandomNameGenerator is a class for generating unique, randomized
// JavaScript variable/property names.
//
// Unlike NameGenerator, names do not follow a predictable sequence such as
// a, b, ... z, aa, ab, ..., az, ba, ... but instead they are random, based
// on an external random seed. We do partially compromise for efficiency in
// that
//
//   - Generated names will have the same length as they would with
//     NameGenerator
//   - We don't use a completely different alphabet for each name prefix,
//     but instead choose among a few with a predictable formula.
//
// More precisely:
//
//   - We compute a random shuffle of the alphabet for "first characters",
//     and a small number of random shuffles of the alphabet for "non-first
//     characters". Then we do a typical number-to-text conversion of a
//     name's "index", where the alphabet for digits is not just 0 to 9.
//     The least significant digit comes first.
//   - We represent each digit using an appropriate alphabet. If it's not
//     the first character of the name (i.e. not the least significant one,
//     or there is a constant prefix) then we have several appropriate
//     alphabets to choose from; we choose one based a hash of the previous
//     digits of this name.
//
// This class is not thread safe.
type RandomNameGenerator struct {
	// The possible first characters, after reserved characters are removed
	firstChars []uint16

	// Possible non-first characters, after reserved characters are removed
	nonFirstChars []uint16

	firstAllowed    [128]bool
	nonFirstAllowed [128]bool

	// Source of randomness
	random *JavaRandom

	// Set of reserved names; these are not returned by generateNextName
	reservedNames map[string]struct{}

	// Prefix added to all generated names
	prefix         string
	prefixUTF16Len int

	// How many names have we issued so far (includes names that cannot be
	// used because they are reserved through 'reservedNames' or JavaScript
	// keywords)
	nameCount int

	// Randomly-shuffled version of firstChars
	shuffledFirst []uint16

	// Randomly-shuffled versions of nonFirstChars (there are NUM_SHUFFLES of
	// them)
	shuffledNonFirst [randomNameGenerator_NUM_SHUFFLES][]uint16

	firstBuf            [len(defaultNameGenerator_FIRST_CHAR_STR)]uint16
	nonFirstBuf         [len(defaultNameGenerator_NONFIRST_CHAR_STR)]uint16
	shuffledFirstBuf    [len(defaultNameGenerator_FIRST_CHAR_STR)]uint16
	shuffledNonFirstBuf [randomNameGenerator_NUM_SHUFFLES][len(defaultNameGenerator_NONFIRST_CHAR_STR)]uint16
	scratchBuf          [32]byte
}

// RandomNameGenerator_FIRST_CHAR: generate random names with this first
// character.
var RandomNameGenerator_FIRST_CHAR = ASCIICharSetFromString(defaultNameGenerator_FIRST_CHAR_STR)

// RandomNameGenerator_NONFIRST_CHAR: these appear after after the first
// character
var RandomNameGenerator_NONFIRST_CHAR = ASCIICharSetFromString(defaultNameGenerator_NONFIRST_CHAR_STR)

// How many shuffles of nonFirstChars to generate
const randomNameGenerator_NUM_SHUFFLES = 16

// Java: RandomNameGenerator(Random)
func NewRandomNameGenerator(random *JavaRandom) *RandomNameGenerator {
	g := &RandomNameGenerator{random: random}
	g.Reset(nil, "", ASCIICharSet{})
	return g
}

// Java: RandomNameGenerator(Set<String>, String, Set<Character>, Random)
func NewRandomNameGeneratorWithReserved(
	reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet, random *JavaRandom) *RandomNameGenerator {
	g := &RandomNameGenerator{random: random}
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
	return g
}

// NewRandomNameGeneratorWithFirstAndNonFirstChars creates a
// RandomNameGenerator.
//
// reservedNames is the set of names that are reserved; generated names
// will not include these names. This set is referenced rather than copied,
// so changes to the set will be reflected in how names are generated. All
// generated names begin with prefix (a name consisting of only this
// prefix, with no suffix, will not be generated). If specified,
// reservedFirstCharacters won't be used in generated names for the first
// character, and reservedNonFirstCharacters won't be used in generated
// names for characters after the first. random is the source of randomness
// when generating random names.
//
// Java: RandomNameGenerator(Set<String>, String, Set<Character>, Set<Character>, Random)
func NewRandomNameGeneratorWithFirstAndNonFirstChars(
	reservedNames jutil.Set[string],
	prefix string,
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	random *JavaRandom) *RandomNameGenerator {
	g := &RandomNameGenerator{random: random}
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
	return g
}

// Java: RandomNameGenerator.reset(Set<String>, String, Set<Character>)
func (g *RandomNameGenerator) Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) {
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java: RandomNameGenerator.reset(Set<String>, String, Set<Character>, Set<Character>)
func (g *RandomNameGenerator) ResetWithFirstAndNonFirstChars(
	reservedNames jutil.Set[string],
	prefix string,
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet) {
	g.reservedNames = randomNameGenerator_copyOf(reservedNames)
	g.prefix = prefix
	g.prefixUTF16Len = jutil.StringLength(prefix)
	g.nameCount = 0

	// Build the character sets to use
	g.firstChars, g.firstAllowed = randomNameGenerator_filterASCII(
		DefaultNameGenerator_FIRST_CHAR, reservedFirstCharacters, g.firstBuf[:0])
	g.nonFirstChars, g.nonFirstAllowed = randomNameGenerator_filterASCII(
		DefaultNameGenerator_NONFIRST_CHAR, reservedNonFirstCharacters, g.nonFirstBuf[:0])

	g.checkPrefix(prefix)
	g.shuffleAlphabets()
}

// randomNameGenerator_copyOf is ImmutableSet.copyOf(Collection).
func randomNameGenerator_copyOf(s jutil.Set[string]) map[string]struct{} {
	if s == nil || s.IsEmpty() {
		return nil
	}
	out := make(map[string]struct{}, s.Size())
	for e := range s.All() {
		out[e] = struct{}{}
	}
	return out
}

// randomNameGenerator_filterASCII filters an ASCII base alphabet (in order)
// against reserved, returning the remaining slice and an O(1) ASCII table.
func randomNameGenerator_filterASCII(
	base []uint16,
	reserved ASCIICharSet,
	dst []uint16,
) ([]uint16, [128]bool) {
	var allowed [128]bool
	if reserved.IsEmpty() {
		for _, c := range base {
			dst = append(dst, c)
			allowed[c] = true
		}
		return dst, allowed
	}
	for _, c := range base {
		if !reserved.Contains(c) {
			dst = append(dst, c)
			allowed[c] = true
		}
	}
	return dst, allowed
}

// Java: RandomNameGenerator.clone(Set<String>, String, Set<Character>)
func (g *RandomNameGenerator) Clone(
	reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) NameGenerator {
	return NewRandomNameGeneratorWithReserved(reservedNames, prefix, reservedCharacters, g.random)
}

// randomNameGenerator_joinChars is Joiner.on(", ").join(chars).
func randomNameGenerator_joinChars(chars []uint16) string {
	parts := make([]string, len(chars))
	for i, c := range chars {
		parts[i] = jutil.StringFromChar(c)
	}
	return strings.Join(parts, ", ")
}

// checkPrefix validates a name prefix.
//
// Java: RandomNameGenerator.checkPrefix(String)
func (g *RandomNameGenerator) checkPrefix(prefix string) {
	if prefix == "" {
		return
	}
	asciiValid := prefix[0] < 0x80 && g.firstAllowed[prefix[0]]
	for i := 1; asciiValid && i < len(prefix); i++ {
		asciiValid = prefix[i] < 0x80 && g.nonFirstAllowed[prefix[i]]
	}
	if asciiValid {
		return
	}
	units := jutil.StringToUTF16(prefix)
	if len(units) > 0 {
		// Make sure that prefix starts with a legal character.
		if units[0] >= 0x80 || !g.firstAllowed[units[0]] {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"prefix must start with one of: " + randomNameGenerator_joinChars(g.firstChars)))
		}
		for pos := 1; pos < len(units); pos++ {
			if units[pos] >= 0x80 || !g.nonFirstAllowed[units[pos]] {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"prefix has invalid characters, must be one of: " +
						randomNameGenerator_joinChars(g.nonFirstChars)))
			}
		}
	}
}

// shuffleAlphabets generates random shuffles of the alphabets.
//
// Java: RandomNameGenerator.shuffleAlphabets()
func (g *RandomNameGenerator) shuffleAlphabets() {
	first := g.shuffledFirstBuf[:len(g.firstChars)]
	copy(first, g.firstChars)
	shuffleSlice(first, g.random)
	g.shuffledFirst = first

	for i := 0; i < randomNameGenerator_NUM_SHUFFLES; i++ {
		nonFirst := g.shuffledNonFirstBuf[i][:len(g.nonFirstChars)]
		copy(nonFirst, g.nonFirstChars)
		shuffleSlice(nonFirst, g.random)
		g.shuffledNonFirst[i] = nonFirst
	}
}

// getNameLength computes the length (in digits) for a name suffix.
//
// Java: RandomNameGenerator.getNameLength(int, int)
func (g *RandomNameGenerator) getNameLength(position int, nameIdx int) int {
	length := 0
	nameIdx++
	for {
		nameIdx--
		var alphabetSize int
		if position == 0 {
			alphabetSize = len(g.firstChars)
		} else {
			alphabetSize = len(g.nonFirstChars)
		}
		if alphabetSize == 0 {
			panic(jutil.NewArithmeticException("/ by zero"))
		}
		nameIdx /= alphabetSize
		position++
		length++
		if !(nameIdx > 0) {
			break
		}
	}
	return length
}

// rngMurmur3HashLengthAndChars computes Hashing.murmur3_128().newHasher()
// .putInt(length).putUnencodedChars(name).hash().asInt() without heap allocation.
func rngMurmur3HashLengthAndChars(length int32, chars []uint16) int32 {
	n := 4 + 2*len(chars)
	var stackBuf [32]byte
	var data []byte
	if n <= len(stackBuf) {
		data = stackBuf[:n]
	} else {
		data = make([]byte, n)
	}
	binary.LittleEndian.PutUint32(data[0:4], uint32(length))
	for idx, c := range chars {
		binary.LittleEndian.PutUint16(data[4+2*idx:], c)
	}
	return murmur3_128BytesAsInt(data)
}

func (g *RandomNameGenerator) isReserved(cand string) bool {
	if len(g.reservedNames) == 0 {
		return false
	}
	_, ok := g.reservedNames[cand]
	return ok
}

// generateCandidate returns g.prefix + generateSuffix(prefixUTF16Len, nameIdx)
// using stack/inline buffers without intermediate string or slice allocations.
//
// Java: prefix + RandomNameGenerator.generateSuffix(int, int)
func (g *RandomNameGenerator) generateCandidate(nameIdx int) string {
	position := g.prefixUTF16Len
	length := g.getNameLength(position, nameIdx)
	var nameBuf [16]uint16
	name := nameBuf[:0]
	nameIdx++
	for {
		nameIdx--
		var alphabet []uint16
		if position == 0 {
			alphabet = g.shuffledFirst
		} else {
			hashInt := rngMurmur3HashLengthAndChars(int32(length), name)
			alphabetIdx := (hashInt & 0x7fffffff) % randomNameGenerator_NUM_SHUFFLES
			alphabet = g.shuffledNonFirst[alphabetIdx]
		}
		alphabetSize := len(alphabet)
		if alphabetSize == 0 {
			panic(jutil.NewArithmeticException("/ by zero"))
		}
		character := alphabet[nameIdx%alphabetSize]
		name = append(name, character)
		nameIdx /= alphabetSize
		position++
		if !(nameIdx > 0) {
			break
		}
	}
	if g.prefix == "" && len(name) == 1 {
		b := byte(name[0])
		return asciiByteTable[b : b+1]
	}
	if len(g.prefix)+len(name) <= len(g.scratchBuf) {
		n := copy(g.scratchBuf[:], g.prefix)
		for _, c := range name {
			g.scratchBuf[n] = byte(c)
			n++
		}
		return string(g.scratchBuf[:n])
	}
	return g.prefix + jutil.StringFromUTF16(name)
}

// GenerateNextName generates the next short name.
//
// This generates names of increasing length. To minimize output size,
// therefore, it's good to call it for the most used symbols first.
//
// Java: RandomNameGenerator.generateNextName()
func (g *RandomNameGenerator) GenerateNextName() string {
	for {
		nameIdx := g.nameCount
		g.nameCount++
		name := g.generateCandidate(nameIdx)

		// Make sure it's not a JS keyword or reserved name.
		if rhino.TokenStream_IsKeyword(name) ||
			g.isReserved(name) ||
			DefaultNameGenerator_IsBadName(name) {
			continue
		}

		return name
	}
}

const (
	murmur3C1 uint64 = 0x87c37b91114253d5
	murmur3C2 uint64 = 0x4cf5ad432745937f
)

func murmur3MixK1(k1 uint64) uint64 {
	k1 *= murmur3C1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= murmur3C2
	return k1
}

func murmur3MixK2(k2 uint64) uint64 {
	k2 *= murmur3C2
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= murmur3C1
	return k2
}

func murmur3Fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func murmur3x64_128(data []byte, seed uint32) (uint64, uint64) {
	h1, h2 := uint64(seed), uint64(seed)
	n := len(data)
	i := 0
	for ; i+16 <= n; i += 16 {
		k1 := binary.LittleEndian.Uint64(data[i:])
		k2 := binary.LittleEndian.Uint64(data[i+8:])
		h1 ^= murmur3MixK1(k1)
		h1 = bits.RotateLeft64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		h2 ^= murmur3MixK2(k2)
		h2 = bits.RotateLeft64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}
	if rem := n - i; rem > 0 {
		var k1, k2 uint64
		tail := data[i:]
		for j := rem - 1; j >= 8; j-- {
			k2 ^= uint64(tail[j]) << (8 * uint(j-8))
		}
		for j := min(rem, 8) - 1; j >= 0; j-- {
			k1 ^= uint64(tail[j]) << (8 * uint(j))
		}
		h1 ^= murmur3MixK1(k1)
		h2 ^= murmur3MixK2(k2)
	}
	h1 ^= uint64(n)
	h2 ^= uint64(n)
	h1 += h2
	h2 += h1
	h1 = murmur3Fmix64(h1)
	h2 = murmur3Fmix64(h2)
	h1 += h2
	h2 += h1
	return h1, h2
}

func murmur3_128BytesAsInt(data []byte) int32 {
	h1, _ := murmur3x64_128(data, 0)
	return int32(uint32(h1))
}

const (
	javaRandomMultiplier int64 = 0x5DEECE66D
	javaRandomAddend     int64 = 0xB
	javaRandomMask       int64 = (1 << 48) - 1
)

// JavaRandom emulates java.util.Random (the 48-bit linear congruential generator).
type JavaRandom struct {
	seed int64
}

var javaRandomInstances atomic.Int64

func NewJavaRandom() *JavaRandom {
	return NewJavaRandomWithSeed(time.Now().UnixNano() + javaRandomInstances.Add(1)*1_000_003)
}

func NewJavaRandomWithSeed(seed int64) *JavaRandom {
	return &JavaRandom{seed: (seed ^ javaRandomMultiplier) & javaRandomMask}
}

func (r *JavaRandom) SetSeed(seed int64) {
	r.seed = (seed ^ javaRandomMultiplier) & javaRandomMask
}

func (r *JavaRandom) Next(nbits int) int32 {
	r.seed = (r.seed*javaRandomMultiplier + javaRandomAddend) & javaRandomMask
	return int32(uint64(r.seed) >> uint(48-nbits))
}

func (r *JavaRandom) NextInt() int32 { return r.Next(32) }

func (r *JavaRandom) NextIntBound(bound int32) int32 {
	if bound <= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("bound must be positive"))
	}
	if bound&-bound == bound {
		return int32((int64(bound) * int64(r.Next(31))) >> 31)
	}
	for {
		sample := r.Next(31)
		v := sample % bound
		if int64(sample)-int64(v)+int64(bound-1) <= 0x7FFFFFFF {
			return v
		}
	}
}

func (r *JavaRandom) NextLong() int64 {
	hi := int64(r.Next(32))
	lo := int64(r.Next(32))
	return (hi << 32) + lo
}

func (r *JavaRandom) NextBoolean() bool { return r.Next(1) != 0 }

func (r *JavaRandom) NextDouble() float64 {
	hi := int64(r.Next(26))
	lo := int64(r.Next(27))
	return float64((hi<<27)+lo) * (1.0 / float64(int64(1)<<53))
}

func (r *JavaRandom) NextFloat() float32 {
	return float32(r.Next(24)) / float32(1<<24)
}

func shuffleSlice[T any](list []T, rnd *JavaRandom) {
	for i := len(list); i > 1; i-- {
		j := rnd.NextIntBound(int32(i))
		list[i-1], list[j] = list[j], list[i-1]
	}
}

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

package jscomp

import (
	"strings"

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
	firstChars *jutil.LinkedHashSet[uint16]

	// Possible non-first characters, after reserved characters are removed
	nonFirstChars *jutil.LinkedHashSet[uint16]

	// Source of randomness
	random *jutil.Random

	// List of reserved names; these are not returned by generateNextName
	reservedNames *jutil.LinkedHashSet[string]

	// Prefix added to all generated names
	prefix string

	// How many names have we issued so far (includes names that cannot be
	// used because they are reserved through 'reservedNames' or JavaScript
	// keywords)
	nameCount int

	// Randomly-shuffled version of firstChars
	shuffledFirst []uint16

	// Randomly-shuffled versions of nonFirstChars (there are NUM_SHUFFLES of
	// them)
	shuffledNonFirst [][]uint16
}

// RandomNameGenerator_FIRST_CHAR: generate random names with this first
// character.
var RandomNameGenerator_FIRST_CHAR = randomNameGenerator_asSet(DefaultNameGenerator_FIRST_CHAR)

// RandomNameGenerator_NONFIRST_CHAR: these appear after after the first
// character
var RandomNameGenerator_NONFIRST_CHAR = randomNameGenerator_asSet(DefaultNameGenerator_NONFIRST_CHAR)

// How many shuffles of nonFirstChars to generate
const randomNameGenerator_NUM_SHUFFLES = 16

// Java: RandomNameGenerator(Random)
func NewRandomNameGenerator(random *jutil.Random) *RandomNameGenerator {
	g := &RandomNameGenerator{random: random}
	g.Reset(jutil.NewLinkedHashSet[string](), "", jutil.NewLinkedHashSet[uint16]())
	return g
}

// Java: RandomNameGenerator(Set<String>, String, Set<Character>, Random)
func NewRandomNameGeneratorWithReserved(
	reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16], random *jutil.Random) *RandomNameGenerator {
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
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16],
	random *jutil.Random) *RandomNameGenerator {
	g := &RandomNameGenerator{random: random}
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
	return g
}

// Java: RandomNameGenerator.reset(Set<String>, String, Set<Character>)
func (g *RandomNameGenerator) Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) {
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java: RandomNameGenerator.reset(Set<String>, String, Set<Character>, Set<Character>)
func (g *RandomNameGenerator) ResetWithFirstAndNonFirstChars(
	reservedNames jutil.Set[string],
	prefix string,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16]) {
	g.reservedNames = randomNameGenerator_copyOf(reservedNames)
	g.prefix = prefix
	g.nameCount = 0

	// Build the character sets to use
	g.firstChars = randomNameGenerator_difference(RandomNameGenerator_FIRST_CHAR, reservedFirstCharacters)
	g.nonFirstChars = randomNameGenerator_difference(RandomNameGenerator_NONFIRST_CHAR, reservedNonFirstCharacters)

	g.checkPrefix(prefix)
	g.shuffleAlphabets()
}

// randomNameGenerator_copyOf is ImmutableSet.copyOf(Collection).
func randomNameGenerator_copyOf(s jutil.Set[string]) *jutil.LinkedHashSet[string] {
	out := jutil.NewLinkedHashSet[string]()
	for e := range s.All() {
		out.Add(e)
	}
	return out
}

// randomNameGenerator_difference is Sets.difference(set1,
// set2).immutableCopy(): the elements of set1 not in set2, in set1's
// order.
func randomNameGenerator_difference(set1 *jutil.LinkedHashSet[uint16], set2 jutil.Set[uint16]) *jutil.LinkedHashSet[uint16] {
	out := jutil.NewLinkedHashSet[uint16]()
	for e := range set1.All() {
		if !set2.Contains(e) {
			out.Add(e)
		}
	}
	return out
}

// Java: RandomNameGenerator.clone(Set<String>, String, Set<Character>)
func (g *RandomNameGenerator) Clone(
	reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) NameGenerator {
	return NewRandomNameGeneratorWithReserved(reservedNames, prefix, reservedCharacters, g.random)
}

// Java: RandomNameGenerator.asSet(char[])
func randomNameGenerator_asSet(chars []uint16) *jutil.LinkedHashSet[uint16] {
	if chars == nil {
		return jutil.NewLinkedHashSet[uint16]()
	}
	return jutil.NewLinkedHashSet(chars...)
}

// randomNameGenerator_joinChars is Joiner.on(", ").join(chars).
func randomNameGenerator_joinChars(chars *jutil.LinkedHashSet[uint16]) string {
	var parts []string
	for c := range chars.All() {
		parts = append(parts, jutil.StringFromChar(c))
	}
	return strings.Join(parts, ", ")
}

// checkPrefix validates a name prefix.
//
// Java: RandomNameGenerator.checkPrefix(String)
func (g *RandomNameGenerator) checkPrefix(prefix string) {
	units := jutil.StringToUTF16(prefix)
	if len(units) > 0 {
		// Make sure that prefix starts with a legal character.
		if !g.firstChars.Contains(units[0]) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"prefix must start with one of: " + randomNameGenerator_joinChars(g.firstChars)))
		}
		for pos := 1; pos < len(units); pos++ {
			if !g.nonFirstChars.Contains(units[pos]) {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"prefix has invalid characters, must be one of: " +
						randomNameGenerator_joinChars(g.nonFirstChars)))
			}
		}
	}
}

// Java: RandomNameGenerator.shuffleAndCopyAlphabet(Set<Character>, Random)
func randomNameGenerator_shuffleAndCopyAlphabet(input *jutil.LinkedHashSet[uint16], random *jutil.Random) []uint16 {
	shuffled := input.ToSlice()
	jutil.Collections_Shuffle(shuffled, random)
	return shuffled
}

// shuffleAlphabets generates random shuffles of the alphabets.
//
// Java: RandomNameGenerator.shuffleAlphabets()
func (g *RandomNameGenerator) shuffleAlphabets() {
	g.shuffledFirst = randomNameGenerator_shuffleAndCopyAlphabet(g.firstChars, g.random)
	builder := make([][]uint16, 0, randomNameGenerator_NUM_SHUFFLES)
	for i := 0; i < randomNameGenerator_NUM_SHUFFLES; i++ {
		builder = append(builder, randomNameGenerator_shuffleAndCopyAlphabet(g.nonFirstChars, g.random))
	}
	g.shuffledNonFirst = builder
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
			alphabetSize = g.firstChars.Size()
		} else {
			alphabetSize = g.nonFirstChars.Size()
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

// generateSuffix returns the nameIdx-th short name. This might be a
// reserved name. A user-requested prefix is not included, but the first
// returned character is supposed to go at position position in the final
// name.
//
// Java: RandomNameGenerator.generateSuffix(int, int)
func (g *RandomNameGenerator) generateSuffix(position int, nameIdx int) string {
	name := make([]uint16, 0, 8)
	length := g.getNameLength(position, nameIdx)
	nameIdx++
	for {
		nameIdx--
		var alphabet []uint16
		if position == 0 {
			alphabet = g.shuffledFirst
		} else {
			hasher := jutil.Hashing_Murmur3_128NewHasher()
			hasher.PutInt(int32(length))
			for _, c := range name {
				hasher.PutChar(c) // putUnencodedChars(name)
			}
			alphabetIdx := (hasher.Hash().AsInt() & 0x7fffffff) % randomNameGenerator_NUM_SHUFFLES
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
	return jutil.StringFromUTF16(name)
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
		name := g.prefix + g.generateSuffix(jutil.StringLength(g.prefix), nameIdx)

		// Make sure it's not a JS keyword or reserved name.
		if rhino.TokenStream_IsKeyword(name) ||
			g.reservedNames.Contains(name) ||
			DefaultNameGenerator_IsBadName(name) {
			continue
		}

		return name
	}
}

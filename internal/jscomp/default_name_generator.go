// Port of com/google/javascript/jscomp/DefaultNameGenerator.java @ closure-compiler v20261015

/*
 * Copyright 2005 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DefaultNameGenerator is a simple class for generating unique JavaScript
// variable/property names.
//
// This class is not thread safe.
type DefaultNameGenerator struct {
	// TODO(user): Maybe we don't need a LinkedHashMap to look up. I started
	// writing a skip-list like data-structure that would let us have O(1)
	// favors() and O(1) reset() but the code got very messy. Lets start
	// with a logical implementation first until performance becomes a
	// problem.
	priorityLookupMap *jutil.LinkedHashMap[uint16, *defaultNameGenerator_CharPriority]

	reservedNames jutil.Set[string]
	prefix        string
	nameCount     int

	firstChars    []*defaultNameGenerator_CharPriority
	nonFirstChars []*defaultNameGenerator_CharPriority
}

// defaultNameGenerator_CharPriority represents a char that can be used in
// renaming as well as how often that char appears in the generated code.
type defaultNameGenerator_CharPriority struct {
	name       uint16
	occurrence int

	// This is a tie-breaker when two chars occurrence count is the same.
	// When that happens, the 'natural' order prevails.
	order int
}

// Java: DefaultNameGenerator.CharPriority(char, int)
func newDefaultNameGenerator_CharPriority(name uint16, order int) *defaultNameGenerator_CharPriority {
	return &defaultNameGenerator_CharPriority{name: name, order: order, occurrence: 0}
}

// Java: DefaultNameGenerator.CharPriority.clone()
func (c *defaultNameGenerator_CharPriority) Clone() *defaultNameGenerator_CharPriority {
	result := newDefaultNameGenerator_CharPriority(c.name, c.order)
	result.occurrence = c.occurrence
	return result
}

// Java: DefaultNameGenerator.CharPriority.compareTo(CharPriority)
func (c *defaultNameGenerator_CharPriority) CompareTo(other *defaultNameGenerator_CharPriority) int {
	// Start out by putting the element with more occurrence first.
	result := other.occurrence - c.occurrence
	if result != 0 {
		return result
	}
	// If there is a tie, follow the order of FIRST_CHAR and NONFIRST_CHAR.
	result = c.order - other.order
	return result
}

// It is important that the ordering of FIRST_CHAR is as close to
// NONFIRST_CHAR as possible. Using the ASCII ordering is not a good idea.
// The reason is that we cannot use numbers as FIRST_CHAR yet the ACSII
// value of numbers is very small. If we picked numbers first in
// NONFIRST_CHAR, we would end up balancing the huffman tree and result is
// bad compression.

// DefaultNameGenerator_FIRST_CHAR: generate short name with this first
// character.
var DefaultNameGenerator_FIRST_CHAR = jutil.StringToUTF16(
	"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$")

// DefaultNameGenerator_NONFIRST_CHAR: these appear after the first
// character.
var DefaultNameGenerator_NONFIRST_CHAR = jutil.StringToUTF16(
	"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789$")

// defaultNameGenerator_BAD_NAMES initializes names that must not be
// generated. TokenStream.isKeyword() skips scanning new keywords such as
// "let", "await", etc. Hence, those names must be explicitly checked for on
// top of the successfully scanned ES3 reservedNames
var defaultNameGenerator_BAD_NAMES = jutil.NewLinkedHashSet("let", "yield", "await", "eval")

// Java: DefaultNameGenerator()
func NewDefaultNameGenerator() *DefaultNameGenerator {
	g := &DefaultNameGenerator{}
	g.buildPriorityLookupMap()
	reservedNames := jutil.NewLinkedHashSet[string]()
	g.Reset(reservedNames, "", jutil.NewLinkedHashSet[uint16]())
	return g
}

// Java: DefaultNameGenerator(Set<String>, String, Set<Character>)
func NewDefaultNameGeneratorWithReserved(
	reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) *DefaultNameGenerator {
	return NewDefaultNameGeneratorWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// NewDefaultNameGeneratorWithFirstAndNonFirstChars creates a
// DefaultNameGenerator.
//
// reservedNames is the set of names that are reserved; generated names
// will not include these names. This set is referenced rather than copied,
// so changes to the set will be reflected in how names are generated. All
// generated names begin with prefix. If specified, reservedFirstCharacters
// won't be used in generated names for the first character, and
// reservedNonFirstCharacters won't be used in generated names for
// characters after the first.
//
// Java: DefaultNameGenerator(Set<String>, String, Set<Character>, Set<Character>)
func NewDefaultNameGeneratorWithFirstAndNonFirstChars(
	reservedNames jutil.Set[string],
	prefix string,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16]) *DefaultNameGenerator {
	g := &DefaultNameGenerator{}
	g.buildPriorityLookupMap()
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
	return g
}

// Java: DefaultNameGenerator(Set<String>, String, Set<Character>, Map<Character, CharPriority>)
func newDefaultNameGeneratorCloned(
	reservedNames jutil.Set[string],
	prefix string,
	reservedCharacters jutil.Set[uint16],
	priorityLookupMap *jutil.LinkedHashMap[uint16, *defaultNameGenerator_CharPriority]) *DefaultNameGenerator {
	g := &DefaultNameGenerator{}
	// Clone the priorityLookupMap to preserve information about how often
	// characters are used.
	g.priorityLookupMap = jutil.NewLinkedHashMap[uint16, *defaultNameGenerator_CharPriority]()
	for key, value := range priorityLookupMap.All() {
		g.priorityLookupMap.Put(key, value.Clone())
	}

	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
	return g
}

// Java: DefaultNameGenerator.buildPriorityLookupMap()
func (g *DefaultNameGenerator) buildPriorityLookupMap() {
	g.priorityLookupMap = jutil.NewLinkedHashMap[uint16, *defaultNameGenerator_CharPriority]()
	order := 0
	for _, c := range DefaultNameGenerator_NONFIRST_CHAR {
		g.priorityLookupMap.Put(c, newDefaultNameGenerator_CharPriority(c, order))
		order++
	}
}

// DefaultNameGenerator_IsBadName returns if name is a badName (let, yield,
// eval, etc) that must not be generated.
//
// Java: DefaultNameGenerator.isBadName(String)
func DefaultNameGenerator_IsBadName(name string) bool {
	return defaultNameGenerator_BAD_NAMES.Contains(name)
}

// Java: DefaultNameGenerator.reset(Set<String>, String, Set<Character>)
func (g *DefaultNameGenerator) Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) {
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// ResetWithFirstAndNonFirstChars: note that the history of what characters
// are most used in the program (set through calls to 'favor') is not
// deleted. Upon 'reset', that history is taken into account for the names
// that will be generated later: it re-calculates how characters are
// prioritized based on how often the they appear in the final output.
//
// Java: DefaultNameGenerator.reset(Set<String>, String, Set<Character>, Set<Character>)
func (g *DefaultNameGenerator) ResetWithFirstAndNonFirstChars(
	reservedNames jutil.Set[string],
	prefix string,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16]) {

	g.reservedNames = reservedNames
	g.prefix = prefix
	g.nameCount = 0

	// build the character arrays to use
	g.firstChars = g.ReserveCharacters(DefaultNameGenerator_FIRST_CHAR, reservedFirstCharacters)
	g.nonFirstChars = g.ReserveCharacters(DefaultNameGenerator_NONFIRST_CHAR, reservedNonFirstCharacters)
	// Arrays.sort(Object[]) is a stable merge sort.
	slices.SortStableFunc(g.firstChars, (*defaultNameGenerator_CharPriority).CompareTo)
	slices.SortStableFunc(g.nonFirstChars, (*defaultNameGenerator_CharPriority).CompareTo)

	g.checkPrefix(prefix)
}

// Java: DefaultNameGenerator.clone(Set<String>, String, Set<Character>)
func (g *DefaultNameGenerator) Clone(
	reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) NameGenerator {
	return newDefaultNameGeneratorCloned(reservedNames, prefix, reservedCharacters,
		g.priorityLookupMap)
}

// Favors increases the prioritization of all the chars in a String. This
// information is not used until reset is called. A compiler would be able
// to generate names while changing the prioritization of the name
// generator for the next pass.
//
// Java: DefaultNameGenerator.favors(CharSequence)
func (g *DefaultNameGenerator) Favors(sequence string) {
	for _, ch := range jutil.StringToUTF16(sequence) {
		c := g.priorityLookupMap.Get(ch)
		if c != nil {
			c.occurrence++
		}
	}
}

// ReserveCharacters provides the array of available characters based on
// the specified arrays.
//
// chars is the list of characters that are legal; reservedCharacterSet the
// characters that should not be used. Returns an array of characters to
// use. Will return the chars array if reservedCharacters is null or empty,
// otherwise creates a new array.
//
// Java: DefaultNameGenerator.reserveCharacters(char[], Set<Character>)
func (g *DefaultNameGenerator) ReserveCharacters(chars []uint16, reservedCharacterSet jutil.Set[uint16]) []*defaultNameGenerator_CharPriority {
	if reservedCharacterSet.IsEmpty() {
		result := make([]*defaultNameGenerator_CharPriority, len(chars))
		for i := 0; i < len(chars); i++ {
			result[i] = g.priorityLookupMap.Get(chars[i])
		}
		return result
	}
	charSet := jutil.NewLinkedHashSet(chars...)
	for c := range reservedCharacterSet.All() {
		charSet.Remove(c)
	}

	result := make([]*defaultNameGenerator_CharPriority, charSet.Size())
	index := 0
	for c := range charSet.All() {
		result[index] = g.priorityLookupMap.Get(c)
		index++
	}
	return result
}

// checkPrefix validates a name prefix.
//
// Java: DefaultNameGenerator.checkPrefix(String)
func (g *DefaultNameGenerator) checkPrefix(prefix string) {
	units := jutil.StringToUTF16(prefix)
	if len(units) > 0 {
		// Make sure that prefix starts with a legal character.
		if !defaultNameGenerator_contains(g.firstChars, units[0]) {
			chars := make([]uint16, len(g.firstChars))
			for i := 0; i < len(chars); i++ {
				chars[i] = g.firstChars[i].name
			}
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"prefix must start with one of: " + jutil.CollectionToString(chars)))
		}
		for pos := 1; pos < len(units); pos++ {
			chars := make([]uint16, len(g.nonFirstChars))
			for i := 0; i < len(chars); i++ {
				chars[i] = g.nonFirstChars[i].name
			}
			if !defaultNameGenerator_contains(g.nonFirstChars, units[pos]) {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"prefix has invalid characters, must be one of: " +
						jutil.CollectionToString(chars)))
			}
		}
	}
}

// Java: DefaultNameGenerator.contains(CharPriority[], char)
func defaultNameGenerator_contains(arr []*defaultNameGenerator_CharPriority, c uint16) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i].name == c {
			return true
		}
	}
	return false
}

// GenerateNextName generates the next short name.
//
// Java: DefaultNameGenerator.generateNextName()
func (g *DefaultNameGenerator) GenerateNextName() string {
	var name string
	for {
		sb := jutil.NewStringBuilder()
		i := g.nameCount

		if g.prefix != "" {
			sb.Append(g.prefix)
		}

		if sb.Length() == 0 {
			if len(g.firstChars) == 0 {
				panic(jutil.NewArithmeticException("/ by zero"))
			}
			pos := i % len(g.firstChars)
			sb.AppendChar(g.firstChars[pos].name)
			i /= len(g.firstChars)
		}

		for i > 0 {
			i--
			if len(g.nonFirstChars) == 0 {
				panic(jutil.NewArithmeticException("/ by zero"))
			}
			pos := i % len(g.nonFirstChars)
			sb.AppendChar(g.nonFirstChars[pos].name)
			i /= len(g.nonFirstChars)
		}

		g.nameCount++
		name = sb.String()

		// Make sure it's not a JS keyword or reserved name.
		if !(rhino.TokenStream_IsKeyword(name) || g.reservedNames.Contains(name) || DefaultNameGenerator_IsBadName(name)) {
			break
		}
	}
	return name
}

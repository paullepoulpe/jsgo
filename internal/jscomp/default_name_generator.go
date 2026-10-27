// Port of com/google/javascript/jscomp/DefaultNameGenerator.java @ closure-compiler v20261021

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

const asciiByteTable = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
	"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
	" !\"#$%&'()*+,-./0123456789:;<=>?" +
	"@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
	"`abcdefghijklmnopqrstuvwxyz{|}~\x7f"

const (
	defaultNameGenerator_FIRST_CHAR_STR    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$"
	defaultNameGenerator_NONFIRST_CHAR_STR = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789$"
)

var defaultNameGenerator_orderByASCII = func() [128]int8 {
	var table [128]int8
	for i := range table {
		table[i] = -1
	}
	for order := 0; order < len(defaultNameGenerator_NONFIRST_CHAR_STR); order++ {
		table[defaultNameGenerator_NONFIRST_CHAR_STR[order]] = int8(order)
	}
	return table
}()

// DefaultNameGenerator is a simple class for generating unique JavaScript
// variable/property names.
//
// This class is not thread safe.
type DefaultNameGenerator struct {
	reservedNames jutil.Set[string]
	prefix        string
	nameCount     int

	firstChars    []*defaultNameGenerator_CharPriority
	nonFirstChars []*defaultNameGenerator_CharPriority
	priorities    [len(defaultNameGenerator_NONFIRST_CHAR_STR)]defaultNameGenerator_CharPriority
	firstBuf      [len(defaultNameGenerator_FIRST_CHAR_STR)]*defaultNameGenerator_CharPriority
	nonFirstBuf   [len(defaultNameGenerator_NONFIRST_CHAR_STR)]*defaultNameGenerator_CharPriority
	scratchBuf    [32]byte
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
var DefaultNameGenerator_FIRST_CHAR = jutil.StringToUTF16(defaultNameGenerator_FIRST_CHAR_STR)

// DefaultNameGenerator_NONFIRST_CHAR: these appear after the first
// character.
var DefaultNameGenerator_NONFIRST_CHAR = jutil.StringToUTF16(defaultNameGenerator_NONFIRST_CHAR_STR)

// Java: DefaultNameGenerator()
func NewDefaultNameGenerator() *DefaultNameGenerator {
	g := &DefaultNameGenerator{}
	g.initPriorities()
	g.Reset(nil, "", nil)
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
	g.initPriorities()
	g.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
	return g
}

func (g *DefaultNameGenerator) initPriorities() {
	for order := 0; order < len(defaultNameGenerator_NONFIRST_CHAR_STR); order++ {
		g.priorities[order] = defaultNameGenerator_CharPriority{
			name:  uint16(defaultNameGenerator_NONFIRST_CHAR_STR[order]),
			order: order,
		}
	}
}

// DefaultNameGenerator_IsBadName returns if name is a badName (let, yield,
// eval, etc) that must not be generated. TokenStream.isKeyword() skips
// scanning new keywords such as "let", "await", etc. Hence, those names must
// be explicitly checked on top of the successfully scanned ES3 reservedNames.
//
// Java: DefaultNameGenerator.isBadName(String)
func DefaultNameGenerator_IsBadName(name string) bool {
	switch name {
	case "let", "yield", "await", "eval":
		return true
	default:
		return false
	}
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
	g.firstChars = g.reserveCharactersInto(g.firstBuf[:0], defaultNameGenerator_FIRST_CHAR_STR, reservedFirstCharacters)
	g.nonFirstChars = g.reserveCharactersInto(g.nonFirstBuf[:0], defaultNameGenerator_NONFIRST_CHAR_STR, reservedNonFirstCharacters)
	// Arrays.sort(Object[]) is a stable merge sort.
	slices.SortStableFunc(g.firstChars, (*defaultNameGenerator_CharPriority).CompareTo)
	slices.SortStableFunc(g.nonFirstChars, (*defaultNameGenerator_CharPriority).CompareTo)

	g.checkPrefix(prefix)
}

// Java: DefaultNameGenerator.clone(Set<String>, String, Set<Character>)
func (g *DefaultNameGenerator) Clone(
	reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) NameGenerator {
	cloned := &DefaultNameGenerator{
		priorities: g.priorities,
	}
	cloned.ResetWithFirstAndNonFirstChars(reservedNames, prefix, reservedCharacters, reservedCharacters)
	return cloned
}

// Favors increases the prioritization of all the chars in a String. This
// information is not used until reset is called. A compiler would be able
// to generate names while changing the prioritization of the name
// generator for the next pass.
//
// Java: DefaultNameGenerator.favors(CharSequence)
func (g *DefaultNameGenerator) Favors(sequence string) {
	// All NONFIRST_CHAR entries are ASCII (< 0x80), and multi-byte UTF-8 / WTF-8
	// sequences never contain bytes < 0x80.
	for i := 0; i < len(sequence); i++ {
		b := sequence[i]
		if b < 0x80 {
			if idx := defaultNameGenerator_orderByASCII[b]; idx >= 0 {
				g.priorities[idx].occurrence++
			}
		}
	}
}

func (g *DefaultNameGenerator) reserveCharactersInto(
	dst []*defaultNameGenerator_CharPriority,
	chars string,
	reservedCharacterSet jutil.Set[uint16],
) []*defaultNameGenerator_CharPriority {
	if reservedCharacterSet == nil || reservedCharacterSet.IsEmpty() {
		for i := 0; i < len(chars); i++ {
			dst = append(dst, &g.priorities[defaultNameGenerator_orderByASCII[chars[i]]])
		}
		return dst
	}
	for i := 0; i < len(chars); i++ {
		ch := chars[i]
		if !reservedCharacterSet.Contains(uint16(ch)) {
			dst = append(dst, &g.priorities[defaultNameGenerator_orderByASCII[ch]])
		}
	}
	return dst
}

// checkPrefix validates a name prefix.
//
// Java: DefaultNameGenerator.checkPrefix(String)
func (g *DefaultNameGenerator) checkPrefix(prefix string) {
	if prefix == "" {
		return
	}
	units := jutil.StringToUTF16(prefix)
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
		if !defaultNameGenerator_contains(g.nonFirstChars, units[pos]) {
			chars := make([]uint16, len(g.nonFirstChars))
			for i := 0; i < len(chars); i++ {
				chars[i] = g.nonFirstChars[i].name
			}
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"prefix has invalid characters, must be one of: " +
					jutil.CollectionToString(chars)))
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
	buf := &g.scratchBuf
	for {
		i := g.nameCount
		n := 0
		asciiOK := true
		if g.prefix != "" {
			if len(g.prefix) > 16 {
				asciiOK = false
			} else {
				for k := 0; k < len(g.prefix); k++ {
					if g.prefix[k] >= 0x80 {
						asciiOK = false
						break
					}
				}
				if asciiOK {
					n = copy(buf[:], g.prefix)
				}
			}
		}

		if asciiOK {
			if n == 0 {
				if len(g.firstChars) == 0 {
					panic(jutil.NewArithmeticException("/ by zero"))
				}
				pos := i % len(g.firstChars)
				ch := g.firstChars[pos].name
				if ch >= 0x80 {
					asciiOK = false
				} else {
					buf[0] = byte(ch)
					n = 1
					i /= len(g.firstChars)
				}
			}
			for asciiOK && i > 0 {
				i--
				if len(g.nonFirstChars) == 0 {
					panic(jutil.NewArithmeticException("/ by zero"))
				}
				pos := i % len(g.nonFirstChars)
				ch := g.nonFirstChars[pos].name
				if ch >= 0x80 || n >= len(buf) {
					asciiOK = false
					break
				}
				buf[n] = byte(ch)
				n++
				i /= len(g.nonFirstChars)
			}
		}

		if asciiOK {
			g.nameCount++
			if n == 1 {
				b := buf[0]
				cand := asciiByteTable[b : b+1]
				if rhino.TokenStream_IsKeyword(cand) || (g.reservedNames != nil && g.reservedNames.Contains(cand)) || DefaultNameGenerator_IsBadName(cand) {
					continue
				}
				return cand
			}
			cand := string(buf[:n])
			if rhino.TokenStream_IsKeyword(cand) || (g.reservedNames != nil && g.reservedNames.Contains(cand)) || DefaultNameGenerator_IsBadName(cand) {
				continue
			}
			return cand
		}

		sb := jutil.NewStringBuilder()
		i = g.nameCount

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
		name := sb.String()

		// Make sure it's not a JS keyword or reserved name.
		if !(rhino.TokenStream_IsKeyword(name) || (g.reservedNames != nil && g.reservedNames.Contains(name)) || DefaultNameGenerator_IsBadName(name)) {
			return name
		}
	}
}

// Port of com/google/javascript/jscomp/DefaultNameGenerator.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/NameGenerator.java @ closure-compiler v20261021

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
	"iter"
	"math/bits"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// ASCIICharSet is a 128-bit value set of ASCII characters (0..127).
type ASCIICharSet [2]uint64

func NewASCIICharSet(chars ...byte) ASCIICharSet {
	var s ASCIICharSet
	for _, c := range chars {
		if c < 128 {
			s[c>>6] |= 1 << (c & 63)
		}
	}
	return s
}

func ASCIICharSetFromString(str string) ASCIICharSet {
	var s ASCIICharSet
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c < 128 {
			s[c>>6] |= 1 << (c & 63)
		}
	}
	return s
}

func (s *ASCIICharSet) Add(ch uint16) {
	if ch < 128 {
		s[ch>>6] |= 1 << (ch & 63)
	}
}

func (s ASCIICharSet) Contains(ch uint16) bool {
	if ch >= 128 {
		return false
	}
	return (s[ch>>6] & (1 << (ch & 63))) != 0
}

func (s ASCIICharSet) Union(other ASCIICharSet) ASCIICharSet {
	return ASCIICharSet{s[0] | other[0], s[1] | other[1]}
}

func (s ASCIICharSet) IsEmpty() bool {
	return s[0] == 0 && s[1] == 0
}

func (s ASCIICharSet) Size() int {
	return bits.OnesCount64(s[0]) + bits.OnesCount64(s[1])
}

func (s ASCIICharSet) Ptr() *ASCIICharSet {
	return &s
}

// stringMapSet adapts a plain Go map[string]struct{} to jutil.Set[string]
// without LinkedHashSet entry/order slice overhead.
type stringMapSet map[string]struct{}

func (s stringMapSet) Size() int              { return len(s) }
func (s stringMapSet) IsEmpty() bool          { return len(s) == 0 }
func (s stringMapSet) Contains(e string) bool { _, ok := s[e]; return ok }
func (s stringMapSet) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		// order-independent: unordered set view used only for membership snapshots
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

// setNodeStringID updates a StringNode's payload to pre-interned id without
// re-interning, while preserving SourceInformationAnnotator's lazy OriginalName
// spill when kindFlagOrigNameIsString is set.
func setNodeStringID(n rhino.Node, id strpool.ID) {
	if id == 0 {
		panic(jutil.NewNullPointerException())
	}
	curID := n.GetStringID()
	if curID == uint32(id) {
		return
	}
	if n.GetOriginalNameID() == curID {
		n.SetString(strpool.String(id))
		return
	}
	n.SetStringID(id)
}

// NameGenerator is a class that generates unique JavaScript
// variable/property names.
//
// Go-only: language-adapter — the Java Set parameters are jutil.Set views:
// reservedNames is referenced, not copied (it may be a live view such as a
// map's key set).
type NameGenerator interface {
	// Java: NameGenerator.reset(Set<String>, String, Set<Character>)
	Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet)

	// Java: NameGenerator.reset(Set<String>, String, Set<Character>, Set<Character>)
	ResetWithFirstAndNonFirstChars(
		reservedNames jutil.Set[string],
		prefix string,
		reservedFirstCharacters ASCIICharSet,
		reservedNonFirstCharacters ASCIICharSet)

	// Java: NameGenerator.clone(Set<String>, String, Set<Character>)
	Clone(reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) NameGenerator

	// Java: NameGenerator.generateNextName()
	GenerateNextName() string
}

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
	reservedMap   map[string]struct{}
	prefix        string
	nameCount     int
	hasFavored    bool

	firstChars      []*defaultNameGenerator_CharPriority
	nonFirstChars   []*defaultNameGenerator_CharPriority
	firstAllowed    [128]bool
	nonFirstAllowed [128]bool
	priorities      [len(defaultNameGenerator_NONFIRST_CHAR_STR)]defaultNameGenerator_CharPriority
	firstBuf        [len(defaultNameGenerator_FIRST_CHAR_STR)]*defaultNameGenerator_CharPriority
	nonFirstBuf     [len(defaultNameGenerator_NONFIRST_CHAR_STR)]*defaultNameGenerator_CharPriority
	scratchBuf      [32]byte
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
	g.Reset(nil, "", ASCIICharSet{})
	return g
}

// Java: DefaultNameGenerator(Set<String>, String, Set<Character>)
func NewDefaultNameGeneratorWithReserved(
	reservedNames jutil.Set[string], prefix string, reservedCharacters *ASCIICharSet) *DefaultNameGenerator {
	var chars ASCIICharSet
	if reservedCharacters != nil {
		chars = *reservedCharacters
	}
	return NewDefaultNameGeneratorWithFirstAndNonFirstChars(reservedNames, prefix, chars, chars)
}

func NewDefaultNameGeneratorWithReservedChars(
	reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) *DefaultNameGenerator {
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
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet) *DefaultNameGenerator {
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
func (g *DefaultNameGenerator) Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) {
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
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet) {

	g.reservedNames = reservedNames
	if m, ok := reservedNames.(stringMapSet); ok {
		g.reservedMap = m
	} else {
		g.reservedMap = nil
	}
	g.prefix = prefix
	g.nameCount = 0

	// build the character arrays to use
	g.firstChars, g.firstAllowed = g.reserveCharactersInto(g.firstBuf[:0], defaultNameGenerator_FIRST_CHAR_STR, reservedFirstCharacters)
	g.nonFirstChars, g.nonFirstAllowed = g.reserveCharactersInto(g.nonFirstBuf[:0], defaultNameGenerator_NONFIRST_CHAR_STR, reservedNonFirstCharacters)
	if g.hasFavored {
		// Arrays.sort(Object[]) is a stable merge sort.
		slices.SortStableFunc(g.firstChars, (*defaultNameGenerator_CharPriority).CompareTo)
		slices.SortStableFunc(g.nonFirstChars, (*defaultNameGenerator_CharPriority).CompareTo)
	}

	g.checkPrefix(prefix)
}

func (g *DefaultNameGenerator) isReserved(cand string) bool {
	if g.reservedMap != nil {
		_, ok := g.reservedMap[cand]
		return ok
	}
	return g.reservedNames != nil && g.reservedNames.Contains(cand)
}

// Java: DefaultNameGenerator.clone(Set<String>, String, Set<Character>)
func (g *DefaultNameGenerator) Clone(
	reservedNames jutil.Set[string], prefix string, reservedCharacters ASCIICharSet) NameGenerator {
	cloned := &DefaultNameGenerator{
		priorities: g.priorities,
		hasFavored: g.hasFavored,
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
				g.hasFavored = true
			}
		}
	}
}

func (g *DefaultNameGenerator) reserveCharactersInto(
	dst []*defaultNameGenerator_CharPriority,
	chars string,
	reservedCharacterSet ASCIICharSet,
) ([]*defaultNameGenerator_CharPriority, [128]bool) {
	var allowed [128]bool
	if reservedCharacterSet.IsEmpty() {
		for i := 0; i < len(chars); i++ {
			ch := chars[i]
			dst = append(dst, &g.priorities[defaultNameGenerator_orderByASCII[ch]])
			allowed[ch] = true
		}
		return dst, allowed
	}
	for i := 0; i < len(chars); i++ {
		ch := chars[i]
		if !reservedCharacterSet.Contains(uint16(ch)) {
			dst = append(dst, &g.priorities[defaultNameGenerator_orderByASCII[ch]])
			allowed[ch] = true
		}
	}
	return dst, allowed
}

func formatCharPriorities(arr []*defaultNameGenerator_CharPriority) string {
	if len(arr) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.Grow(2 + 3*len(arr))
	sb.WriteByte('[')
	for i, cp := range arr {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteByte(byte(cp.name))
	}
	sb.WriteByte(']')
	return sb.String()
}

// checkPrefix validates a name prefix.
//
// Java: DefaultNameGenerator.checkPrefix(String)
func (g *DefaultNameGenerator) checkPrefix(prefix string) {
	if prefix == "" {
		return
	}
	// Make sure that prefix starts with a legal character.
	if prefix[0] >= 0x80 || !g.firstAllowed[prefix[0]] {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"prefix must start with one of: " + formatCharPriorities(g.firstChars)))
	}
	for pos := 1; pos < len(prefix); pos++ {
		if prefix[pos] >= 0x80 || !g.nonFirstAllowed[prefix[pos]] {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"prefix has invalid characters, must be one of: " +
					formatCharPriorities(g.nonFirstChars)))
		}
	}
}

// GenerateNextName generates the next short name.
//
// Java: DefaultNameGenerator.generateNextName()
func (g *DefaultNameGenerator) GenerateNextName() string {
	for {
		i := g.nameCount
		b := g.scratchBuf[:0]
		if g.prefix != "" {
			b = append(b, g.prefix...)
		} else {
			if len(g.firstChars) == 0 {
				panic(jutil.NewArithmeticException("/ by zero"))
			}
			pos := i % len(g.firstChars)
			b = append(b, byte(g.firstChars[pos].name))
			i /= len(g.firstChars)
		}

		for i > 0 {
			i--
			if len(g.nonFirstChars) == 0 {
				panic(jutil.NewArithmeticException("/ by zero"))
			}
			pos := i % len(g.nonFirstChars)
			b = append(b, byte(g.nonFirstChars[pos].name))
			i /= len(g.nonFirstChars)
		}

		g.nameCount++
		var cand string
		if len(b) == 1 {
			c := b[0]
			cand = asciiByteTable[c : c+1]
		} else {
			cand = string(b)
		}
		if rhino.TokenStream_IsKeyword(cand) || g.isReserved(cand) || DefaultNameGenerator_IsBadName(cand) {
			continue
		}
		return cand
	}
}

// Port of com/google/javascript/jscomp/RandomNameGeneratorTest.java @ closure-compiler v20261015

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
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func generateRandomNames(ng *RandomNameGenerator, count int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = ng.GenerateNextName()
	}
	return result
}

func stringSet(names []string) map[string]bool {
	s := map[string]bool{}
	for _, n := range names {
		s[n] = true
	}
	return s
}

func TestRandomNameGenerator_ConstructorInvalidPrefixes(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)

	ex := assertThrows[jutil.IIllegalArgumentException](t, func() {
		NewRandomNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "123abc", noChars(), random)
	})
	// The error messages should contain meaningful information.
	if !strings.Contains(ex.GetMessage(), "W, X, Y, Z, $") {
		t.Errorf("message %q", ex.GetMessage())
	}

	ex = assertThrows[jutil.IIllegalArgumentException](t, func() {
		NewRandomNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "abc%", noChars(), random)
	})
	if !strings.Contains(ex.GetMessage(), "W, X, Y, Z, _, 0, 1") {
		t.Errorf("message %q", ex.GetMessage())
	}
}

func TestRandomNameGenerator_Generate(t *testing.T) {
	// Since there's a hash function involved, there's not much point in
	// mocking Random to get nicer values. Instead, let's just try to verify
	// the sanity of the results.
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet[string]()
	prefix := "prefix"
	prefixLen := len(prefix)
	// Add a prefix to avoid dropping JavaScript keywords.
	ng := NewRandomNameGeneratorWithReserved(reservedNames, prefix, noChars(), random)
	// Generate all 1- and 2-character names.
	// alphabet length, 1st digit
	len1 := RandomNameGenerator_NONFIRST_CHAR.Size()
	// alphabet length, 2nd digit
	len2 := RandomNameGenerator_NONFIRST_CHAR.Size()
	// len1 == len2 because we have a prefix
	count := len1 * (1 + len2)
	result := generateRandomNames(ng, count)
	resultSet := stringSet(result)

	// We got as many names as we asked for, and all are different.
	assertEqual(t, "distinct", len(resultSet), count)

	// First come names with length 1, then 2. No names are longer.
	for i := 0; i < len1; i++ {
		assertEqual(t, result[i], len(result[i]), prefixLen+1)
	}
	for i := len1; i < count; i++ {
		assertEqual(t, result[i], len(result[i]), prefixLen+2)
	}

	// We don't just have an alphabet for the first character and one for
	// the second, i.e. if we just had the valid characters A,B,C,D,E and
	// the shuffle was such that we started with C,D,A,B,E as names, the
	// following names would not just be
	// CC,DC,AC,BC,EC,CD,DD,AD,BD,ED,... Of course it could randomly happen
	// some time, which is why we check for all occurrences.

	// Verify that we don't get C,D,A,B,E,_C,_C,_C,_C,_C,_D,_D,_D,_D,_D_,...
	countPass := 0
	countTest := 0
	for i1 := 0; i1 < len1; i1++ {
		for i2 := 0; i2 < len2; i2, countTest = i2+1, countTest+1 {
			if result[i1][prefixLen] != result[len1+i1*len2+i2][prefixLen+1] {
				countPass++
			}
		}
	}
	if !(100.0*float64(countPass)/float64(countTest) > 80.0) { // arbitrary threshold
		t.Errorf("not shuffled enough: %d/%d", countPass, countTest)
	}

	// Names are not sorted (some might be, by chance)
	countPass = 0
	countTest = 0
	for i := 0; i < count-1; i++ {
		if jutil.StringCompareTo(result[i], result[i+1]) > 0 {
			countPass++
		}
	}
	if !(100.0*float64(countPass)/float64(countTest) > 25.0) { // arbitrary threshold
		t.Errorf("sorted: %d/%d", countPass, countTest)
	}
}

func TestRandomNameGenerator_FirstCharAlphabet(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet[string]()
	ng := NewRandomNameGeneratorWithReserved(reservedNames, "", noChars(), random)
	// Generate all 1- and 2-character names.
	len1 := RandomNameGenerator_FIRST_CHAR.Size()
	len2 := RandomNameGenerator_NONFIRST_CHAR.Size()
	count := len1 * (1 + len2)
	result := generateRandomNames(ng, count)
	resultSet := stringSet(result)

	// We got as many names as we asked for, and all are different.
	assertEqual(t, "distinct", len(resultSet), count)

	// First come len1 names with length 1, then 2.
	for i := 0; i < len1; i++ {
		assertEqual(t, result[i], len(result[i]), 1)
	}
	// Because there's no prefix, we use the "first characters" alphabet for
	// the first character, so there will be only len1 length-1 names, not
	// len2 as in testGenerate.

	for i := len1; i < count; i++ {
		// Because we don't have a prefix, some generated names will be
		// JavaScript keywords and will be skipped, so we'll get a few
		// length-3 names.
		if len(result[i]) < 2 {
			t.Errorf("%q is shorter than 2", result[i])
		}
	}
}

func TestRandomNameGenerator_Prefix(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet[string]()
	prefix := "prefix"
	ng := NewRandomNameGeneratorWithReserved(reservedNames, prefix, noChars(), random)
	// Generate all 1- and 2-character names.
	len1 := RandomNameGenerator_FIRST_CHAR.Size()
	len2 := RandomNameGenerator_NONFIRST_CHAR.Size()
	count := len1 * (1 + len2)
	result := generateRandomNames(ng, count)

	for i := 0; i < count; i++ {
		if !strings.HasPrefix(result[i], prefix) || len(result[i]) <= len(prefix) {
			t.Errorf("%q", result[i])
		}
	}
}

func TestRandomNameGenerator_Seeds(t *testing.T) {
	// Using different seeds should return different names.
	random0 := jutil.NewRandomWithSeed(0)
	random1 := jutil.NewRandomWithSeed(1)
	reservedNames := jutil.NewLinkedHashSet[string]()
	ng0 := NewRandomNameGeneratorWithReserved(reservedNames, "", noChars(), random0)
	ng1 := NewRandomNameGeneratorWithReserved(reservedNames, "", noChars(), random1)

	count := 1000
	results0 := generateRandomNames(ng0, count)
	results1 := generateRandomNames(ng1, count)

	countPass := 0
	countTest := 0
	for i := 0; i < count; i, countTest = i+1, countTest+1 {
		if results0[i] != results1[i] {
			countPass++
		}
	}
	if !(100.0*float64(countPass)/float64(countTest) > 90.0) { // arbitrary threshold
		t.Errorf("too similar: %d/%d", countPass, countTest)
	}
}

func TestRandomNameGenerator_ReservedNames(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet("x", "ba")
	ng := NewRandomNameGeneratorWithReserved(reservedNames, "", noChars(), random)
	// Generate all 1- and 2-character names (and a couple 3-character
	// names, because "x" and "ba", and keywords, shouldn't be used).
	count := RandomNameGenerator_FIRST_CHAR.Size() * (RandomNameGenerator_NONFIRST_CHAR.Size() + 1)
	result := stringSet(generateRandomNames(ng, count))

	assertEqual(t, "x", result["x"], false)
	assertEqual(t, "ba", result["ba"], false)

	// Even though we skipped "x" and "ba", we still got 'count' different
	// names. We know they are different because 'result' is a Set.
	assertEqual(t, "distinct", len(result), count)
}

func TestRandomNameGenerator_ReservedCharacters(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet[string]()
	ng := NewRandomNameGeneratorWithReserved(reservedNames, "", noChars('a', 'b'), random)
	// Generate all 1- and 2-character names (and also many 3-character
	// names, because "a" and "b" shouldn't be used).
	count := RandomNameGenerator_FIRST_CHAR.Size() * (RandomNameGenerator_NONFIRST_CHAR.Size() + 1)
	result := stringSet(generateRandomNames(ng, count))

	for _, name := range []string{"a", "b", "ax", "bx", "xa", "xb"} {
		assertEqual(t, name, result[name], false)
	}
	assertEqual(t, "x", result["x"], true)
	assertEqual(t, "xx", result["xx"], true)

	// Even though we skipped a few names with reserved characters, we
	// still got 'count' different names. We know they are different
	// because 'result' is a Set.
	assertEqual(t, "distinct", len(result), count)
}

func TestRandomNameGenerator_ES6KeywordsNotGenerated(t *testing.T) {
	random := jutil.NewRandomWithSeed(0)
	reservedNames := jutil.NewLinkedHashSet[string]()

	ng := NewRandomNameGeneratorWithReserved(reservedNames, "le", noChars('$'), random)
	result := generateRandomNames(ng, 106)
	assertEqual(t, "let", slices.Contains(result, "let"), false)

	ng = NewRandomNameGeneratorWithReserved(reservedNames, "awai", noChars('$'), random)
	result = generateRandomNames(ng, 106)
	assertEqual(t, "await", slices.Contains(result, "await"), false)
}

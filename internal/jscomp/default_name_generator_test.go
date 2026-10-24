// Port of com/google/javascript/jscomp/DefaultNameGeneratorTest.java @ closure-compiler v20261021

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var defaultNameGeneratorTest_RESERVED_NAMES = jutil.NewLinkedHashSet("ba", "xba")

func noChars(chars ...uint16) *jutil.LinkedHashSet[uint16] { return jutil.NewLinkedHashSet(chars...) }

func generateNames(t *testing.T, ng NameGenerator, prefix string, num int) []string {
	t.Helper()
	result := make([]string, num)
	for i := 0; i < num; i++ {
		result[i] = ng.GenerateNextName()
		if !strings.HasPrefix(result[i], prefix) {
			t.Fatalf("Error: %s", result[i])
		}
	}
	return result
}

func TestDefaultNameGenerator_NameGeneratorInvalidPrefixes(t *testing.T) {
	ex := assertThrows[jutil.IIllegalArgumentException](t, func() {
		NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "123abc", noChars())
	})
	// The error messages should contain meaningful information.
	if !strings.Contains(ex.GetMessage(), "W, X, Y, Z, $]") {
		t.Errorf("message %q", ex.GetMessage())
	}

	ex = assertThrows[jutil.IIllegalArgumentException](t, func() {
		NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "abc%", noChars())
	})
	if !strings.Contains(ex.GetMessage(), "W, X, Y, Z, _, 0, 1") {
		t.Errorf("message %q", ex.GetMessage())
	}
}

func TestDefaultNameGenerator_Generate(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result := generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "a")
	assertEqual(t, "25", result[25], "z")
	assertEqual(t, "26", result[26], "A")
	assertEqual(t, "51", result[51], "Z")
	assertEqual(t, "52", result[52], "$")
	assertEqual(t, "53", result[53], "aa")
	// ba is reserved
	assertEqual(t, "54", result[54], "ca")
	assertEqual(t, "104", result[104], "$a")

	ng = NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "x", noChars())
	result = generateNames(t, ng, "x", 132)

	// Expected: x, xa, ..., x$, xaa, ..., x$$
	assertEqual(t, "0", result[0], "x")
	assertEqual(t, "1", result[1], "xa")
	assertEqual(t, "64", result[64], "x$")
	assertEqual(t, "65", result[65], "xaa")
	// xba is reserved
	assertEqual(t, "66", result[66], "xca")
}

func TestDefaultNameGenerator_Reserve(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars('$'))
	result := generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "a")
	assertEqual(t, "25", result[25], "z")
	assertEqual(t, "26", result[26], "A")
	assertEqual(t, "51", result[51], "Z")
	assertEqual(t, "52", result[52], "aa")
	// ba is reserved
	assertEqual(t, "53", result[53], "ca")
	assertEqual(t, "103", result[103], "ab")
}

func TestDefaultNameGenerator_ES6KeywordsNotGenerated(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "le", noChars('$'))
	result := generateNames(t, ng, "le", 106)
	assertEqual(t, "19", result[19], "les")
	assertEqual(t, "20", result[20], "leu") // "let" keyword skipped
	assertEqual(t, "45", result[45], "leT") // "leT" not skipped

	ng = NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "awai", noChars('$'))
	result = generateNames(t, ng, "awai", 106)
	assertEqual(t, "19", result[19], "awais")
	assertEqual(t, "20", result[20], "awaiu") // "await" keyword skipped
	assertEqual(t, "45", result[45], "awaiT") // "awaiT" not skipped
}

func TestDefaultNameGenerator_GenerateWithPriority1(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result := generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "a")
	assertEqual(t, "25", result[25], "z")
	assertEqual(t, "26", result[26], "A")
	assertEqual(t, "51", result[51], "Z")
	assertEqual(t, "52", result[52], "$")
	assertEqual(t, "53", result[53], "aa")

	ng.Favors("b")
	ng.Reset(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result = generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "b")
	assertEqual(t, "1", result[1], "a")
	assertEqual(t, "2", result[2], "c")
	assertEqual(t, "3", result[3], "d")

	ng.Favors("cc")
	ng.Reset(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result = generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "c")
	assertEqual(t, "1", result[1], "b")
	assertEqual(t, "2", result[2], "a")
	assertEqual(t, "3", result[3], "d")
}

func TestDefaultNameGenerator_GenerateWithPriority2(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result := generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "a")
	assertEqual(t, "25", result[25], "z")
	assertEqual(t, "26", result[26], "A")
	assertEqual(t, "51", result[51], "Z")
	assertEqual(t, "52", result[52], "$")
	assertEqual(t, "53", result[53], "aa")

	ng.Favors("function")
	ng.Favors("function")
	ng.Favors("function")

	ng.Reset(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result = generateNames(t, ng, "", 106)

	// All the letters of function should come first. In alphabetical order.
	assertEqual(t, "0", result[0], "n")
	assertEqual(t, "1", result[1], "c")
	assertEqual(t, "2", result[2], "f")
	assertEqual(t, "3", result[3], "i")
	assertEqual(t, "4", result[4], "o")
	assertEqual(t, "5", result[5], "t")
	assertEqual(t, "6", result[6], "u")

	// Back to normal.
	assertEqual(t, "7", result[7], "a")
	assertEqual(t, "8", result[8], "b")

	// c has been prioritized.
	assertEqual(t, "9", result[9], "d")
	assertEqual(t, "10", result[10], "e")

	// This used to start with 'aa' but now n is prioritized over it.
	assertEqual(t, "53", result[53], "nn")
	assertEqual(t, "54", result[54], "cn")
}

func TestDefaultNameGenerator_GenerateWithPriority3(t *testing.T) {
	ng := NewDefaultNameGeneratorWithReserved(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	generateNames(t, ng, "", 106)
	ng.Favors("???")
	ng.Reset(defaultNameGeneratorTest_RESERVED_NAMES, "", noChars())
	result := generateNames(t, ng, "", 106)
	assertEqual(t, "0", result[0], "a")
}

// Java: DefaultNameGeneratorTest.testDeterministicNameGeneration()
func TestDefaultNameGenerator_DeterministicNameGeneration(t *testing.T) {
	reservedChars1 := jutil.NewLinkedHashSet[uint16]('a', 'b', 'c')
	reservedChars2 := jutil.NewLinkedHashSet[uint16]('c', 'b', 'a')

	ng1 := NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "", reservedChars1)
	ng2 := NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "", reservedChars2)

	result1 := generateNames(t, ng1, "", 10)
	result2 := generateNames(t, ng2, "", 10)

	if !slices.Equal(result1, result2) {
		t.Errorf("reserved characters: %v != %v", result1, result2)
	}

	reservedNames1 := jutil.NewLinkedHashSet("a", "b", "c")
	reservedNames2 := jutil.NewLinkedHashSet("c", "b", "a")

	ng1 = NewDefaultNameGeneratorWithReserved(reservedNames1, "", jutil.NewLinkedHashSet[uint16]())
	ng2 = NewDefaultNameGeneratorWithReserved(reservedNames2, "", jutil.NewLinkedHashSet[uint16]())

	result1 = generateNames(t, ng1, "", 10)
	result2 = generateNames(t, ng2, "", 10)

	if !slices.Equal(result1, result2) {
		t.Errorf("reserved names: %v != %v", result1, result2)
	}
}

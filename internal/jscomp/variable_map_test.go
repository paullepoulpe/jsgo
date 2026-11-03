// Port of com/google/javascript/jscomp/VariableMapTest.java @ closure-compiler v20261028

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"bytes"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func variableMapTestOf(kv ...string) *OrderedStringMap[string] {
	m := NewOrderedStringMapWithCapacity[string](len(kv) / 2)
	for i := 0; i < len(kv); i += 2 {
		if !m.PutIfAbsent(kv[i], kv[i+1]) {
			panic(jutil.NewIllegalArgumentExceptionMsg("duplicate key: " + kv[i]))
		}
	}
	return m
}

func TestVariableMap_Cycle1(t *testing.T) {
	variableMapCycleTest(t, variableMapTestOf("AAA", "a", "BBB", "b"))
	variableMapCycleTest(t, variableMapTestOf("AA:AA", "a", "BB:BB", "b"))
	variableMapCycleTest(t, variableMapTestOf("AAA", "a:a", "BBB", "b:b"))
}

func variableMapCycleTest(t *testing.T, m *OrderedStringMap[string]) {
	t.Helper()
	in := NewVariableMap(m)
	serialized := jutil.NewStringFromUTF8Bytes(in.ToBytes())
	out := VariableMap_FromBytes(jutil.StringGetBytesUTF8(serialized))
	variableMapAssertMapsEquals(t, in.ToMap(), out.ToMap())
}

func variableMapAssertMapsEquals(t *testing.T, expected, result *OrderedStringMap[string]) {
	t.Helper()
	assertEqual(t, "size", result.Size(), expected.Size())
	for _, e := range expected.Entries() {
		if got, ok := result.Lookup(e.Key); !ok || got != e.Value {
			t.Errorf("result %v does not contain %s=%s", result.Entries(), e.Key, e.Value)
		}
	}
}

func TestVariableMap_ToBytes(t *testing.T) {
	vm := NewVariableMap(variableMapTestOf("AAA", "a", "BBB", "b"))
	serialized := jutil.NewStringFromUTF8Bytes(vm.ToBytes())
	if !strings.HasSuffix(serialized, "\n") {
		t.Errorf("%q does not end with a newline", serialized)
	}

	lines := jutil.StringSplitLiteral(serialized, "\n")
	assertSameElements(t, "lines", lines, []string{"AAA:a", "BBB:b"})
}

func assertVariableMapNewName(t *testing.T, label string, vm *VariableMap, sourceName, want string) {
	t.Helper()
	got, ok := vm.LookupNewName(sourceName)
	assertEqual(t, label+" found", ok, true)
	assertEqual(t, label, got, want)
}

func assertVariableMapSourceName(t *testing.T, label string, vm *VariableMap, newName, want string) {
	t.Helper()
	got, ok := vm.LookupSourceName(newName)
	assertEqual(t, label+" found", ok, true)
	assertEqual(t, label, got, want)
}

func TestVariableMap_FromBytes(t *testing.T) {
	vm := VariableMap_FromBytes([]byte("AAA:a\nBBB:b\n"))
	assertEqual(t, "size", vm.GetOriginalNameToNewNameMap().Size(), 2)
	assertVariableMapNewName(t, "AAA", vm, "AAA", "a")
	assertVariableMapNewName(t, "BBB", vm, "BBB", "b")
	assertVariableMapSourceName(t, "a", vm, "a", "AAA")
	assertVariableMapSourceName(t, "b", vm, "b", "BBB")
}

func TestVariableMap_FromBytesWithEmptyValue(t *testing.T) {
	vm := VariableMap_FromBytes([]byte("AAA:"))
	assertVariableMapNewName(t, "AAA", vm, "AAA", "")
}

func TestVariableMap_FromStream(t *testing.T) {
	vm := VariableMap_FromStream(bytes.NewReader([]byte("AAA:a\nBBB:b\n")))
	assertEqual(t, "size", vm.GetOriginalNameToNewNameMap().Size(), 2)
	assertVariableMapNewName(t, "AAA", vm, "AAA", "a")
	assertVariableMapNewName(t, "BBB", vm, "BBB", "b")
}

func TestVariableMap_FileFormat1(t *testing.T) {
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("x\ny", "a")).ToBytes(), "x\\ny:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("x:y", "a")).ToBytes(), "x\\:y:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("x\ny", "a")).ToBytes(), "x\\ny:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("x\\y", "a")).ToBytes(), "x\\\\y:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("\n", "a")).ToBytes(), "\\n:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf(":", "a")).ToBytes(), "\\::a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("\n", "a")).ToBytes(), "\\n:a\n")
	variableMapAssertBytesEqual(t, NewVariableMap(variableMapTestOf("\\", "a")).ToBytes(), "\\\\:a\n")
}

func TestVariableMap_FromBytesComplex1(t *testing.T) {
	// Verify we get out what we put in.
	variableMapCycleTest(t, variableMapTestOf("AAA[':f']", "a"))

	// Verify the file format is as expected.
	in := NewVariableMap(variableMapTestOf("AAA[':f']", "a"))
	variableMapAssertBytesEqual(t, in.ToBytes(), "AAA['\\:f']:a\n")
}

func TestVariableMap_FromBytesComplex2(t *testing.T) {
	vm := VariableMap_FromBytes([]byte("AAA['\\:f']:a\n"))

	assertEqual(t, "size", vm.GetOriginalNameToNewNameMap().Size(), 1)
	assertVariableMapNewName(t, "lookup", vm, "AAA[':f']", "a")

	assertEqual(t, "inverse size", vm.GetNewNameToOriginalNameMap().Size(), 1)
	assertVariableMapSourceName(t, "reverse lookup", vm, "a", "AAA[':f']")
}

func variableMapAssertBytesEqual(t *testing.T, bytes1 []byte, s2 string) {
	t.Helper()
	if !bytes.Equal(bytes1, jutil.StringGetBytesUTF8(s2)) {
		t.Errorf("bytes %q, want %q", bytes1, s2)
	}
}

func TestVariableMap_ReverseThrowsErrorOnDuplicate(t *testing.T) {
	assertThrows[jutil.IIllegalArgumentException](t, func() {
		NewVariableMap(variableMapTestOf("AA", "b", "BB", "b"))
	})
}

func TestVariableMap_ReverseLookupOfNullFindsNoName(t *testing.T) {
	// Go-only: language-adapter — a Go string cannot be null; the lookup
	// of a name that is not a value finds no name either.
	vm := NewVariableMap(variableMapTestOf("AA", "a", "BB", "b"))
	if got, ok := vm.LookupSourceName("c"); ok {
		t.Errorf("got %q, want absent", got)
	}
}

// Java: VariableMapTest.testFromLines_leadingNewline()
func TestVariableMap_FromLines_leadingNewline(t *testing.T) {
	vm := VariableMap_FromBytes(jutil.StringGetBytesUTF8("\na:b\n"))
	assertVariableMapNewName(t, "lookup", vm, "a", "b")

	vmCrlf := VariableMap_FromBytes(jutil.StringGetBytesUTF8("\r\na:b\r\n"))
	assertVariableMapNewName(t, "lookup", vmCrlf, "a", "b")
}

// Java: VariableMapTest.testFromLines_colonAtStartThrowsParseException()
func TestVariableMap_FromLines_colonAtStartThrowsParseException(t *testing.T) {
	assertThrows[*VariableMapParseError](t, func() { VariableMap_FromBytes([]byte(":foo\n")) })
}

func TestVariableMap_ToBytesEscapingAndSurrogates(t *testing.T) {
	m := variableMapTestOf(
		"plain", "p",
		"a\\:b:c", "x:y\\z",
		"λ:π", "α\\β",
	)
	variableMapCycleTest(t, m)

	vm := NewVariableMap(variableMapTestOf("a\\:b\nc", "x\ny:z\\w"))
	variableMapAssertBytesEqual(t, vm.ToBytes(), "a\\\\\\:b\\nc:x\\ny\\:z\\\\w\n")

	// Unpaired WTF-8 surrogate must encode to '?' matching jutil.StringGetBytesUTF8.
	unpaired := jutil.StringFromUTF16([]uint16{0xD800, ':', 'k'})
	vmSurrogate := NewVariableMap(variableMapTestOf(unpaired, "v"))
	if !bytes.Equal(vmSurrogate.ToBytes(), []byte("?\\:k:v\n")) {
		t.Fatalf("got %q, want %q", vmSurrogate.ToBytes(), "?\\:k:v\n")
	}
}

// Port of com/google/javascript/jscomp/VariableMapTest.java @ closure-compiler v20261021

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

func variableMapTestOf(kv ...string) *jutil.LinkedHashMap[string, string] {
	b := jutil.NewImmutableMapBuilder[string, string]()
	for i := 0; i < len(kv); i += 2 {
		b.Put(kv[i], kv[i+1])
	}
	return b.BuildOrThrow()
}

func TestVariableMap_Cycle1(t *testing.T) {
	variableMapCycleTest(t, variableMapTestOf("AAA", "a", "BBB", "b"))
	variableMapCycleTest(t, variableMapTestOf("AA:AA", "a", "BB:BB", "b"))
	variableMapCycleTest(t, variableMapTestOf("AAA", "a:a", "BBB", "b:b"))
}

func variableMapCycleTest(t *testing.T, m *jutil.LinkedHashMap[string, string]) {
	t.Helper()
	in := NewVariableMap(m)
	serialized := jutil.NewStringFromUTF8Bytes(in.ToBytes())
	out := VariableMap_FromBytes(jutil.StringGetBytesUTF8(serialized))
	variableMapAssertMapsEquals(t, in.ToMap(), out.ToMap())
}

func variableMapAssertMapsEquals(t *testing.T, expected, result *jutil.LinkedHashMap[string, string]) {
	t.Helper()
	assertEqual(t, "size", result.Size(), expected.Size())
	for key, value := range expected.All() {
		if got, ok := result.Lookup(key); !ok || got != value {
			t.Errorf("result %s does not contain %s=%s", result, key, value)
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

func TestVariableMap_FromBytes(t *testing.T) {
	vm := VariableMap_FromBytes([]byte("AAA:a\nBBB:b\n"))
	assertEqual(t, "size", vm.GetOriginalNameToNewNameMap().Size(), 2)
	assertEqual(t, "AAA", strPtrValue(vm.LookupNewName("AAA")), "a")
	assertEqual(t, "BBB", strPtrValue(vm.LookupNewName("BBB")), "b")
	assertEqual(t, "a", strPtrValue(vm.LookupSourceName("a")), "AAA")
	assertEqual(t, "b", strPtrValue(vm.LookupSourceName("b")), "BBB")
}

func TestVariableMap_FromBytesWithEmptyValue(t *testing.T) {
	vm := VariableMap_FromBytes([]byte("AAA:"))
	assertEqual(t, "AAA", strPtrValue(vm.LookupNewName("AAA")), "")
}

func TestVariableMap_FromStream(t *testing.T) {
	vm := VariableMap_FromStream(bytes.NewReader([]byte("AAA:a\nBBB:b\n")))
	assertEqual(t, "size", vm.GetOriginalNameToNewNameMap().Size(), 2)
	assertEqual(t, "AAA", strPtrValue(vm.LookupNewName("AAA")), "a")
	assertEqual(t, "BBB", strPtrValue(vm.LookupNewName("BBB")), "b")
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
	assertEqual(t, "lookup", strPtrValue(vm.LookupNewName("AAA[':f']")), "a")

	assertEqual(t, "inverse size", vm.GetNewNameToOriginalNameMap().Size(), 1)
	assertEqual(t, "reverse lookup", strPtrValue(vm.LookupSourceName("a")), "AAA[':f']")
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
	if got := vm.LookupSourceName("c"); got != nil {
		t.Errorf("got %q, want null", *got)
	}
}

// Java: VariableMapTest.testFromLines_leadingNewline()
func TestVariableMap_FromLines_leadingNewline(t *testing.T) {
	vm := VariableMap_FromBytes(jutil.StringGetBytesUTF8("\na:b\n"))
	assertEqual(t, "lookup", strPtrValue(vm.LookupNewName("a")), "b")

	vmCrlf := VariableMap_FromBytes(jutil.StringGetBytesUTF8("\r\na:b\r\n"))
	assertEqual(t, "lookup", strPtrValue(vmCrlf.LookupNewName("a")), "b")
}

// Java: VariableMapTest.testFromLines_colonAtStartThrowsParseException()
func TestVariableMap_FromLines_colonAtStartThrowsParseException(t *testing.T) {
	assertThrows[*VariableMapParseError](t, func() { VariableMap_FromBytes([]byte(":foo\n")) })
}

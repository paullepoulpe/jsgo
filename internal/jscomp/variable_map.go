// Port of com/google/javascript/jscomp/VariableMap.java @ closure-compiler v20261021

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
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// VariableMapParseError represents a syntax error when parsing a serialized VariableMap.
type VariableMapParseError struct {
	jutil.JavaException
}

func newVariableMapParseError(msg string) *VariableMapParseError {
	e := &VariableMapParseError{}
	jutil.InitThrowable(e, "java.text.ParseException", &msg, nil)
	return e
}

type variableMapEntry struct {
	orig    string
	newName string
}

// VariableMap stores the bidirectional mapping between original variable names
// and new variable names.
//
// See RenameVars.
type VariableMap struct {
	origToNew map[string]string
	newToOrig map[string]string
	origKeys  []string
}

const variableMap_SEPARATOR = ':'

// Java: VariableMap(Map<String, String>)
func NewVariableMap(m *jutil.LinkedHashMap[string, string]) *VariableMap {
	n := m.Size()
	origToNew := make(map[string]string, n)
	newToOrig := make(map[string]string, n)
	origKeys := make([]string, 0, n)
	for k, v := range m.All() {
		if prevKey, dup := newToOrig[v]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same value: " + k + "=" + v + " and " + prevKey + "=" + v))
		}
		origToNew[k] = v
		newToOrig[v] = k
		origKeys = append(origKeys, k)
	}
	return &VariableMap{
		origToNew: origToNew,
		newToOrig: newToOrig,
		origKeys:  origKeys,
	}
}

func newVariableMapFromEntries(entries []variableMapEntry) *VariableMap {
	n := len(entries)
	origToNew := make(map[string]string, n)
	// Match Guava ImmutableMap.Builder.buildOrThrow() duplicate-key order (last-to-first).
	for i := n - 1; i >= 0; i-- {
		e := entries[i]
		if later, dup := origToNew[e.orig]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + e.orig + "=" + later + " and " + e.orig + "=" + e.newName))
		}
		origToNew[e.orig] = e.newName
	}
	newToOrig := make(map[string]string, n)
	origKeys := make([]string, n)
	for i, e := range entries {
		if prevOrig, dup := newToOrig[e.newName]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same value: " + e.orig + "=" + e.newName + " and " + prevOrig + "=" + e.newName))
		}
		newToOrig[e.newName] = e.orig
		origKeys[i] = e.orig
	}
	return &VariableMap{
		origToNew: origToNew,
		newToOrig: newToOrig,
		origKeys:  origKeys,
	}
}

// LookupNewName given an original variable name, looks up new name, returning
// ("", false) if it's not found.
//
// Java: VariableMap.lookupNewName(String)
func (vm *VariableMap) LookupNewName(sourceName string) (string, bool) {
	v, ok := vm.origToNew[sourceName]
	return v, ok
}

// LookupSourceName given a new variable name, looks up the source name,
// returning ("", false) if it's not found.
//
// Java: VariableMap.lookupSourceName(String)
func (vm *VariableMap) LookupSourceName(newName string) (string, bool) {
	v, ok := vm.newToOrig[newName]
	return v, ok
}

// GetOriginalNameToNewNameMap returns an immutable mapping from original
// names to new names.
//
// Java: VariableMap.getOriginalNameToNewNameMap()
func (vm *VariableMap) GetOriginalNameToNewNameMap() *jutil.LinkedHashMap[string, string] {
	// ImmutableSortedMap.copyOf(map): natural (String.compareTo) key order.
	keys := slices.Clone(vm.origKeys)
	slices.SortFunc(keys, jutil.StringCompareTo)
	out := jutil.NewLinkedHashMap[string, string]()
	for _, k := range keys {
		out.Put(k, vm.origToNew[k])
	}
	return out
}

// GetNewNameToOriginalNameMap returns an immutable mapping from new names to
// original names.
//
// Java: VariableMap.getNewNameToOriginalNameMap()
func (vm *VariableMap) GetNewNameToOriginalNameMap() *jutil.LinkedHashMap[string, string] {
	out := jutil.NewLinkedHashMap[string, string]()
	for _, orig := range vm.origKeys {
		out.Put(vm.origToNew[orig], orig)
	}
	return out
}

// Save saves the variable map to a file.
//
// Java: VariableMap.save(String)
func (vm *VariableMap) Save(filename string) {
	if err := os.WriteFile(filename, vm.ToBytes(), 0o666); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

// VariableMap_Load reads the variable map from a file written via Save.
//
// Java: VariableMap.load(String)
func VariableMap_Load(filename string) *VariableMap {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	vm, err := parseVariableMapLines(jutil.NewStringFromUTF8Bytes(b))
	if err != nil {
		// Wrap parse exception for backwards compatibility.
		panic(jutil.NewIOExceptionWithCause(newVariableMapParseError(err.Error())))
	}
	return vm
}

// ToBytes serializes the variable map to a byte array.
//
// Java: VariableMap.toBytes()
func (vm *VariableMap) ToBytes() []byte {
	// Go-only: language-adapter — the OutputStreamWriter's UTF-8 encoding
	// of the whole text; no surrogate pair can straddle two writes, so
	// encoding the concatenation is the same.
	var sb strings.Builder
	// The output order should be stable.
	keys := slices.Clone(vm.origKeys)
	slices.SortFunc(keys, jutil.StringCompareTo)
	for _, key := range keys {
		sb.WriteString(variableMap_escape(key))
		sb.WriteByte(variableMap_SEPARATOR)
		sb.WriteString(variableMap_escape(vm.origToNew[key]))
		sb.WriteByte('\n')
	}
	return jutil.StringGetBytesUTF8(sb.String())
}

// VariableMap_FromBytes deserializes the variable map from a byte array
// returned by ToBytes.
//
// Java: VariableMap.fromBytes(byte[])
func VariableMap_FromBytes(bytes []byte) *VariableMap {
	return VariableMap_FromLines(jutil.NewStringFromUTF8Bytes(bytes))
}

func parseVariableMapLine(line string) (orig, newName string, ok bool) {
	pos := variableMap_findIndexOfUnescapedChar(line, variableMap_SEPARATOR)
	if pos <= 0 {
		return "", "", false
	}
	if pos != len(line)-1 {
		newName = variableMap_unescape(line[pos+1:])
	}
	return variableMap_unescape(line[:pos]), newName, true
}

func parseVariableMapStream(stream io.Reader) (*VariableMap, error) {
	b, err := io.ReadAll(stream)
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	text := jutil.NewStringFromUTF8Bytes(b)
	var entries []variableMapEntry
	for pos := 0; ; {
		line, nextPos, ok := jutil.NextLineFromString(text, pos)
		if !ok {
			break
		}
		pos = nextPos
		if line == "" {
			continue
		}
		orig, newName, ok := parseVariableMapLine(line)
		if !ok {
			return nil, errors.New("Bad line: " + line)
		}
		entries = append(entries, variableMapEntry{orig: orig, newName: newName})
	}
	return newVariableMapFromEntries(entries), nil
}

func parseVariableMapLines(s string) (*VariableMap, error) {
	var entries []variableMapEntry
	startOfLine := 0
	for startOfLine < len(s) {
		newLine := strings.IndexByte(s[startOfLine:], '\n')
		if newLine == -1 {
			newLine = len(s)
		} else {
			newLine += startOfLine
		}

		endOfLine := newLine
		if newLine > 0 && s[newLine-1] == '\r' {
			newLine--
		}

		line := s[startOfLine:newLine]
		startOfLine = endOfLine + 1 // update index for next iteration

		if line == "" {
			continue
		}
		orig, newName, ok := parseVariableMapLine(line)
		if !ok {
			return nil, errors.New("Bad line: " + line)
		}
		entries = append(entries, variableMapEntry{orig: orig, newName: newName})
	}
	return newVariableMapFromEntries(entries), nil
}

// VariableMap_FromStream deserializes the variable map from an
// InputStream.
//
// Java: VariableMap.fromStream(InputStream)
func VariableMap_FromStream(stream io.Reader) *VariableMap {
	vm, err := parseVariableMapStream(stream)
	if err != nil {
		panic(newVariableMapParseError(err.Error()))
	}
	return vm
}

// VariableMap_FromLines deserializes the variable map from a string.
//
// Go-only: language-adapter — the Java indices are UTF-16 offsets; here
// they are byte offsets. Every index the method computes is at an ASCII
// character ('\n', '\r', ':', '\\'), so the lines and substrings are the
// same.
//
// Java: VariableMap.fromLines(String)
func VariableMap_FromLines(s string) *VariableMap {
	vm, err := parseVariableMapLines(s)
	if err != nil {
		panic(newVariableMapParseError(err.Error()))
	}
	return vm
}

// Java: VariableMap.escape(String)
func variableMap_escape(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "\\", "\\\\"), ":", "\\:"), "\n", "\\n")
}

// variableMap_findIndexOfUnescapedChar works on byte offsets (stopChar is
// ASCII; see VariableMap_FromLines).
//
// Java: VariableMap.findIndexOfUnescapedChar(String, char)
func variableMap_findIndexOfUnescapedChar(value string, stopChar byte) int {
	length := len(value)
	for i := 0; i < length; {
		stopCharIndex := strings.IndexByte(value[i:], stopChar)
		if stopCharIndex == -1 {
			return -1
		}
		stopCharIndex += i
		if stopCharIndex == 0 || value[stopCharIndex-1] != '\\' {
			// it isn't escaped, return
			return stopCharIndex
		}
		i = stopCharIndex + 1
	}
	return -1
}

// variableMap_unescape works on bytes: a backslash is dropped and the
// following character kept, whatever its UTF-8 length.
//
// Java: VariableMap.unescape(String)
func variableMap_unescape(value string) string {
	slashIndex := strings.IndexByte(value, '\\')
	if slashIndex == -1 {
		return value
	}

	var sb strings.Builder
	sb.Grow(len(value) - 1)
	sb.WriteString(value[:slashIndex])
	length := len(value)
	for i := slashIndex; i < length; i++ {
		c := value[i]
		if c == '\\' {
			i++
			if i < length {
				c = value[i]
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// VariableMap_FromMap initializes the variable map from an existing map.
//
// The map is the map to use from original names to generated names. It is
// copied and changes to the specified map will not affect the returned
// object.
//
// Java: VariableMap.fromMap(Map<String, String>)
func VariableMap_FromMap(m *jutil.LinkedHashMap[string, string]) *VariableMap {
	return NewVariableMap(m)
}

// Java: VariableMap.toMap()
func (vm *VariableMap) ToMap() *jutil.LinkedHashMap[string, string] {
	out := jutil.NewLinkedHashMap[string, string]()
	for _, orig := range vm.origKeys {
		out.Put(orig, vm.origToNew[orig])
	}
	return out
}

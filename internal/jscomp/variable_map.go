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
	"io"
	"os"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// VariableMap stores the mapping from original variable name to new
// variable names.
//
// See RenameVars.
type VariableMap struct {
	// Maps between original source name to new name
	map_ *jutil.BiMap[string, string]
}

const variableMap_SEPARATOR = ':'

// Java: VariableMap(Map<String, String>)
func NewVariableMap(m *jutil.LinkedHashMap[string, string]) *VariableMap {
	return &VariableMap{map_: jutil.ImmutableBiMap_CopyOf(m)}
}

// LookupNewName given an original variable name, looks up new name, may
// return null if it's not found.
//
// Java: VariableMap.lookupNewName(String)
func (vm *VariableMap) LookupNewName(sourceName string) *string {
	if v, ok := vm.map_.Lookup(sourceName); ok {
		return &v
	}
	return nil
}

// LookupSourceName given a new variable name, looks up the source name, may
// return null if it's not found.
//
// Java: VariableMap.lookupSourceName(String)
func (vm *VariableMap) LookupSourceName(newName string) *string {
	if v, ok := vm.map_.Inverse().Lookup(newName); ok {
		return &v
	}
	return nil
}

// GetOriginalNameToNewNameMap returns an immutable mapping from original
// names to new names.
//
// Java: VariableMap.getOriginalNameToNewNameMap()
func (vm *VariableMap) GetOriginalNameToNewNameMap() *jutil.LinkedHashMap[string, string] {
	// ImmutableSortedMap.copyOf(map): natural (String.compareTo) key order.
	keys := slices.Collect(vm.map_.Keys())
	slices.SortStableFunc(keys, jutil.StringCompareTo)
	out := jutil.NewLinkedHashMap[string, string]()
	for _, k := range keys {
		out.Put(k, vm.map_.Get(k))
	}
	return out
}

// GetNewNameToOriginalNameMap returns an immutable mapping from new names to
// original names.
//
// Java: VariableMap.getNewNameToOriginalNameMap()
func (vm *VariableMap) GetNewNameToOriginalNameMap() *jutil.LinkedHashMap[string, string] {
	return vm.map_.Inverse().ToLinkedHashMap()
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
	var result *VariableMap
	jutil.Try(func() {
		b, err := os.ReadFile(filename)
		if err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
		result = VariableMap_FromBytes(b)
	}, func(e jutil.IParseException) {
		// Wrap parse exception for backwards compatibility.
		panic(jutil.NewIOExceptionWithCause(e))
	})
	return result
}

// ToBytes serializes the variable map to a byte array.
//
// Java: VariableMap.toBytes()
func (vm *VariableMap) ToBytes() []byte {
	// Go-only: language-adapter — the OutputStreamWriter's UTF-8 encoding
	// of the whole text; no surrogate pair can straddle two writes, so
	// encoding the concatenation is the same.
	writer := jutil.NewStringBuilder()
	// The output order should be stable.
	keys := slices.Collect(vm.map_.Keys())
	slices.SortStableFunc(keys, jutil.StringCompareTo)
	for _, key := range keys {
		writer.Append(variableMap_escape(key))
		writer.AppendChar(variableMap_SEPARATOR)
		writer.Append(variableMap_escape(vm.map_.Get(key)))
		writer.AppendChar('\n')
	}
	return jutil.StringGetBytesUTF8(writer.String())
}

// VariableMap_FromBytes deserializes the variable map from a byte array
// returned by ToBytes.
//
// Java: VariableMap.fromBytes(byte[])
func VariableMap_FromBytes(bytes []byte) *VariableMap {
	return VariableMap_FromLines(jutil.NewStringFromUTF8Bytes(bytes))
}

// VariableMap_FromStream deserializes the variable map from an
// InputStream.
//
// Java: VariableMap.fromStream(InputStream)
func VariableMap_FromStream(stream io.Reader) *VariableMap {
	m := jutil.NewImmutableMapBuilder[string, string]()
	br := jutil.NewBufferedReaderUTF8(stream)
	for {
		line, ok := br.ReadLine()
		if !ok {
			break
		}
		if line == "" {
			continue
		}
		pos := variableMap_findIndexOfUnescapedChar(line, variableMap_SEPARATOR)
		if pos <= 0 {
			panic(jutil.NewParseException("Bad line: "+line, 0))
		}
		value := ""
		if pos != len(line)-1 {
			value = variableMap_unescape(line[pos+1:])
		}
		m.Put(variableMap_unescape(line[:pos]), value)
	}
	return NewVariableMap(m.BuildOrThrow())
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
	m := jutil.NewImmutableMapBuilder[string, string]()
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
		pos := variableMap_findIndexOfUnescapedChar(line, variableMap_SEPARATOR)
		if pos <= 0 {
			panic(jutil.NewParseException("Bad line: "+line, 0))
		}
		value := ""
		if pos != len(line)-1 {
			value = variableMap_unescape(line[pos+1:])
		}
		m.Put(variableMap_unescape(line[:pos]), value)
	}
	return NewVariableMap(m.BuildOrThrow())
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
	return vm.map_.ToLinkedHashMap()
}

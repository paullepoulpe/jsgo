// Port of com/google/javascript/jscomp/IdMappingUtil.java @ closure-compiler v20261021

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
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IdMappingUtil is a utility class for generating and parsing id mappings
// held by ReplaceIdGenerators.
type IdMappingUtil struct{}

// Java: IdMappingUtil.NEW_LINE
const IdMappingUtil_NEW_LINE uint16 = '\n'

// IdMappingSection holds bidirectional key/value entries for a single id
// generator section, preserving key insertion order.
type IdMappingSection struct {
	byKey   map[string]string
	byValue map[string]string
	keys    []string
}

func newIdMappingSection() *IdMappingSection {
	return &IdMappingSection{
		byKey:   make(map[string]string),
		byValue: make(map[string]string),
	}
}

func (s *IdMappingSection) Size() int           { return len(s.byKey) }
func (s *IdMappingSection) IsEmpty() bool       { return len(s.byKey) == 0 }
func (s *IdMappingSection) Get(k string) string { return s.byKey[k] }
func (s *IdMappingSection) Keys() []string      { return s.keys }
func (s *IdMappingSection) LookupByValue(v string) (string, bool) {
	k, ok := s.byValue[v]
	return k, ok
}

func (s *IdMappingSection) Put(k, v string) {
	if oldVal, hasKey := s.byKey[k]; hasKey && oldVal == v {
		return
	}
	if _, valExists := s.byValue[v]; valExists {
		panic(jutil.NewIllegalArgumentExceptionMsg("value already present: " + v))
	}
	if oldVal, hasKey := s.byKey[k]; hasKey {
		delete(s.byValue, oldVal)
		s.byKey[k] = v
		s.byValue[v] = k
		return
	}
	s.keys = append(s.keys, k)
	s.byKey[k] = v
	s.byValue[v] = k
}

// IdMappingUtil_GenerateSerializedIdMappings returns the serialized map of
// generators and their ids and their replacements.
//
// Java: IdMappingUtil.generateSerializedIdMappings(Map<String, Map<String, String>>)
func IdMappingUtil_GenerateSerializedIdMappings(idGeneratorMaps *jutil.LinkedHashMap[string, *jutil.LinkedHashMap[string, string]]) string {
	var sb strings.Builder
	for key, value := range idGeneratorMaps.All() {
		if !value.IsEmpty() {
			sb.WriteByte('[')
			sb.WriteString(key)
			sb.WriteString("]\n\n")

			for rKey, rValue := range value.All() {
				sb.WriteString(rKey)
				sb.WriteByte(':')
				sb.WriteString(rValue)
				sb.WriteByte('\n')
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// IdMappingUtil_MappingReader is a stateful pull parser for reading id
// mapping sections and entries line by line.
type IdMappingUtil_MappingReader struct {
	reader         *jutil.BufferedReader
	currentSection string
	currentKey     string
	currentValue   string
	isSection      bool
	lineIndex      int
}

// Java: IdMappingUtil.MappingReader(BufferedReader)
func NewIdMappingUtil_MappingReader(reader *jutil.BufferedReader) *IdMappingUtil_MappingReader {
	return &IdMappingUtil_MappingReader{reader: reader}
}

// Java: IdMappingUtil.MappingReader.next()
func (mr *IdMappingUtil_MappingReader) Next() bool {
	for {
		line, ok := mr.reader.ReadLine()
		if !ok {
			break
		}
		mr.lineIndex++
		if line == "" {
			continue
		}
		if line[0] == '[' {
			end := len(line) - 1
			if end >= 1 && line[end] >= 0x80 {
				_, w := utf8.DecodeLastRuneInString(line)
				end = len(line) - w
			}
			mr.currentSection = line[1:end]
			mr.currentKey = ""
			mr.currentValue = ""
			mr.isSection = true
			return true
		}
		split := strings.IndexByte(line, ':')
		if split != -1 {
			mr.currentKey = line[:split]
			mr.currentValue = line[split+1:]
			mr.isSection = false
			return true
		}
		panic(jutil.NewIllegalArgumentExceptionMsg(
			fmt.Sprintf("Cannot parse id map.\n Line: %s, lineIndex: %d", line, mr.lineIndex)))
	}
	return false
}

// Java: IdMappingUtil.MappingReader.isSection()
func (mr *IdMappingUtil_MappingReader) IsSection() bool {
	return mr.isSection
}

// Java: IdMappingUtil.MappingReader.getSection()
func (mr *IdMappingUtil_MappingReader) GetSection() string {
	return mr.currentSection
}

// Java: IdMappingUtil.MappingReader.getKey()
func (mr *IdMappingUtil_MappingReader) GetKey() string {
	return mr.currentKey
}

// Java: IdMappingUtil.MappingReader.getValue()
func (mr *IdMappingUtil_MappingReader) GetValue() string {
	return mr.currentValue
}

// Java: IdMappingUtil.MappingReader.getLineIndex()
func (mr *IdMappingUtil_MappingReader) GetLineIndex() int {
	return mr.lineIndex
}

// Java: IdMappingUtil.parseSectionAsStream(InputStream, String)
func IdMappingUtil_ParseSectionAsStream(stream io.Reader, sectionFilter string) *jutil.LinkedHashMap[string, string] {
	// ImmutableMap.Builder.buildKeepingLast(): a repeated key keeps its first
	// position and its last value, like LinkedHashMap.put.
	mapBuilder := jutil.NewLinkedHashMap[string, string]()
	br := jutil.NewBufferedReaderUTF8(stream)
	mr := NewIdMappingUtil_MappingReader(br)
	inSection := false
	for mr.Next() {
		if mr.IsSection() {
			if inSection {
				break // Next section found, stop
			}
			if mr.GetSection() == sectionFilter {
				inSection = true
			}
		} else if inSection {
			mapBuilder.Put(mr.GetKey(), mr.GetValue())
		}
	}
	return mapBuilder
}

// Java: IdMappingUtil.parseSerializedIdMappings(BufferedReader)
func IdMappingUtil_ParseSerializedIdMappingsFromReader(br *jutil.BufferedReader) map[string]*IdMappingSection {
	resultMap := make(map[string]*IdMappingSection)
	mr := NewIdMappingUtil_MappingReader(br)
	var currentSectionMap *IdMappingSection
	for mr.Next() {
		if mr.IsSection() {
			sectionName := mr.GetSection()
			if _, dup := resultMap[sectionName]; dup {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					fmt.Sprintf(
						"Cannot parse id map: Duplicate section %s\n lineIndex: %d",
						sectionName, mr.GetLineIndex())))
			}
			currentSectionMap = newIdMappingSection()
			resultMap[sectionName] = currentSectionMap
		} else {
			if currentSectionMap == nil {
				panic(jutil.NewIllegalArgumentExceptionMsg("Mapping entry outside of section"))
			}
			currentSectionMap.Put(mr.GetKey(), mr.GetValue()) // throws if duplicate values
		}
	}
	return resultMap
}

// Java: IdMappingUtil.parseSerializedIdMappings(String)
func IdMappingUtil_ParseSerializedIdMappings(idMappings *string) map[string]*IdMappingSection {
	if idMappings == nil || *idMappings == "" {
		return make(map[string]*IdMappingSection)
	}
	br := jutil.NewBufferedReaderFromString(*idMappings)
	defer br.Close()
	return IdMappingUtil_ParseSerializedIdMappingsFromReader(br)
}

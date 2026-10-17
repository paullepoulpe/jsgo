// Port of com/google/javascript/jscomp/IdMappingUtil.java @ closure-compiler v20261015

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
	"io"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IdMappingUtil is a utility class for generating and parsing id mappings
// held by ReplaceIdGenerators.
type IdMappingUtil struct{}

// Java: IdMappingUtil.NEW_LINE
const IdMappingUtil_NEW_LINE uint16 = '\n'

// IdMappingUtil_GenerateSerializedIdMappings returns the serialized map of
// generators and their ids and their replacements.
//
// Java: IdMappingUtil.generateSerializedIdMappings(Map<String, Map<String, String>>)
func IdMappingUtil_GenerateSerializedIdMappings(idGeneratorMaps *jutil.LinkedHashMap[string, *jutil.LinkedHashMap[string, string]]) string {
	sb := jutil.NewStringBuilder()
	for key, value := range idGeneratorMaps.All() {
		if !value.IsEmpty() {
			sb.AppendChar('[').Append(key).AppendChar(']').AppendChar(IdMappingUtil_NEW_LINE).AppendChar(IdMappingUtil_NEW_LINE)

			for rKey, rValue := range value.All() {
				sb.Append(rKey).
					AppendChar(':').
					Append(rValue).
					AppendChar(IdMappingUtil_NEW_LINE)
			}
			sb.AppendChar(IdMappingUtil_NEW_LINE)
		}
	}
	return sb.String()
}

// IdMappingUtil_MappingReader is a stateful pull parser for reading id
// mapping sections and entries line by line.
type IdMappingUtil_MappingReader struct {
	reader         *jutil.BufferedReader
	currentSection *string
	currentKey     *string
	currentValue   *string
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
		if jutil.StringCharAt(line, 0) == '[' {
			section := jutil.StringSubstring(line, 1, jutil.StringLength(line)-1)
			mr.currentSection = &section
			mr.currentKey = nil
			mr.currentValue = nil
			return true
		} else {
			split := jutil.StringIndexOfChar(line, ':')
			if split != -1 {
				key := jutil.StringSubstring(line, 0, split)
				value := jutil.StringSubstringFrom(line, split+1)
				mr.currentKey = &key
				mr.currentValue = &value
				return true
			} else {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					jutil.StringFormat("Cannot parse id map.\n Line: %s, lineIndex: %d", line, mr.lineIndex)))
			}
		}
	}
	return false
}

// Java: IdMappingUtil.MappingReader.isSection()
func (mr *IdMappingUtil_MappingReader) IsSection() bool {
	return mr.currentKey == nil
}

// Java: IdMappingUtil.MappingReader.getSection()
func (mr *IdMappingUtil_MappingReader) GetSection() *string {
	return mr.currentSection
}

// Java: IdMappingUtil.MappingReader.getKey()
func (mr *IdMappingUtil_MappingReader) GetKey() *string {
	return mr.currentKey
}

// Java: IdMappingUtil.MappingReader.getValue()
func (mr *IdMappingUtil_MappingReader) GetValue() *string {
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
			if *mr.GetSection() == sectionFilter {
				inSection = true
			}
		} else if inSection {
			mapBuilder.Put(*mr.GetKey(), *mr.GetValue())
		}
	}
	return mapBuilder
}

// Java: IdMappingUtil.parseSerializedIdMappings(BufferedReader)
func IdMappingUtil_ParseSerializedIdMappingsFromReader(br *jutil.BufferedReader) *jutil.LinkedHashMap[string, *jutil.BiMap[string, string]] {
	resultMap := jutil.NewLinkedHashMap[string, *jutil.BiMap[string, string]]()
	mr := NewIdMappingUtil_MappingReader(br)
	var currentSectionMap *jutil.BiMap[string, string]
	sectionNames := jutil.NewLinkedHashSet[string]()
	for mr.Next() {
		if mr.IsSection() {
			sectionName := *mr.GetSection()
			if !sectionNames.Add(sectionName) {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					jutil.StringFormat(
						"Cannot parse id map: Duplicate section %s\n lineIndex: %d",
						sectionName, mr.GetLineIndex())))
			}
			currentSectionMap = jutil.HashBiMap_Create[string, string]()
			resultMap.Put(sectionName, currentSectionMap)
		} else {
			if currentSectionMap == nil {
				panic(jutil.NewIllegalArgumentExceptionMsg("Mapping entry outside of section"))
			}
			currentSectionMap.Put(*mr.GetKey(), *mr.GetValue()) // throws if duplicate values
		}
	}
	return resultMap
}

// Java: IdMappingUtil.parseSerializedIdMappings(String)
func IdMappingUtil_ParseSerializedIdMappings(idMappings *string) *jutil.LinkedHashMap[string, *jutil.BiMap[string, string]] {
	if idMappings == nil || *idMappings == "" {
		return jutil.NewLinkedHashMap[string, *jutil.BiMap[string, string]]()
	}
	var result *jutil.LinkedHashMap[string, *jutil.BiMap[string, string]]
	jutil.Try(func() {
		br := jutil.NewBufferedReaderFromString(*idMappings)
		defer br.Close()
		result = IdMappingUtil_ParseSerializedIdMappingsFromReader(br)
	}, func(e jutil.IIOException) {
		ex := &jutil.IllegalArgumentException{}
		jutil.InitThrowable(ex, "java.lang.IllegalArgumentException", jutil.StringPtr("Failed to parse id mappings"), e)
		panic(ex)
	})
	return result
}

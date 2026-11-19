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
// Port of com/google/javascript/jscomp/IdMappingUtil.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.IdMappingUtil
type IdMappingUtil struct {
	idHash_ atomic.Int32
}

func (this_ *IdMappingUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *IdMappingUtil) IsJavaNil_() bool { return this_ == nil }
func (*IdMappingUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.IdMappingUtil", "IdMappingUtil"
}
func (this_ *IdMappingUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IdMappingUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.IdMappingUtil", "com.google.javascript.jscomp.IdMappingUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*IdMappingUtil](v_)
	})
}

func StaticGet_IdMappingUtil_NEW_LINE() uint16 { return uint16(10) }
func (this_ *IdMappingUtil) initFields_IdMappingUtil() {
	_ = this_
}

// Java constructor: IdMappingUtil()
func NewIdMappingUtil() *IdMappingUtil {
	this_ := &IdMappingUtil{}
	InitIdMappingUtil(this_)
	return this_
}

func InitIdMappingUtil(this_ *IdMappingUtil) {
	this_.initFields_IdMappingUtil()
}

// Java static method: IdMappingUtil.generateSerializedIdMappings(java.util.Map<java.lang.String,java.util.Map<java.lang.String,java.lang.String>> idGeneratorMaps)
func IdMappingUtil_GenerateSerializedIdMappings(idGeneratorMaps j2gort.Ext_java_util_Map) string {
	_ = idGeneratorMaps
	var sb j2gort.LocalStringBuilder

	{
		cur_1 := j2gort.NewMapEntryCursor(idGeneratorMaps)
		for cur_1.NextEntry() {
			if !(j2gort.AsIface[j2gort.Ext_java_util_Map](cur_1.GetValue()).IsEmpty()) {
				sb.AppendChar(uint16(91)).AppendStr(cur_1.KeyStr()).AppendChar(uint16(93)).AppendChar(uint16(10)).AppendChar(uint16(10))
				{
					cur_2 := j2gort.NewMapEntryCursor(j2gort.AsIface[j2gort.Ext_java_util_Map](cur_1.GetValue()))
					for cur_2.NextEntry() {
						sb.AppendStr(cur_2.KeyStr()).AppendChar(uint16(58)).AppendStr(j2gort.CastString(cur_2.GetValue())).AppendChar(uint16(10))
					}
				}
				sb.AppendChar(uint16(10))
			}
		}
	}
	return sb.String()
}

// Java static method: IdMappingUtil.parseSectionAsStream(java.io.InputStream stream, java.lang.String sectionFilter)
func IdMappingUtil_ParseSectionAsStream(stream j2gort.Ext_java_io_InputStream, sectionFilter string) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = stream
	_ = sectionFilter
	var mapBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = mapBuilder
	var br *j2gort.Ext_java_io_BufferedReader = j2gort.NewExt_java_io_BufferedReader_1(j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(stream, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8))
	_ = br
	var val_3 IdMappingUtil_MappingReader
	InitIdMappingUtil_MappingReader(&val_3, br)
	var mr *IdMappingUtil_MappingReader = &val_3
	_ = mr
	var inSection bool = false
	_ = inSection
	for mr.Next() {
		if mr.IsSection() {
			if inSection {
				break
			}
			if j2gort.CheckNotNullString(mr.GetSection()) == sectionFilter {
				inSection = true
			}
		} else if inSection {
			mapBuilder.Put_2(j2gort.BoxString(mr.GetKey()), j2gort.BoxString(mr.GetValue()))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](mapBuilder.BuildKeepingLast())
}

// Java static method: IdMappingUtil.parseSerializedIdMappings(java.io.BufferedReader br)
func IdMappingUtil_ParseSerializedIdMappings_1_Ext_java_io_BufferedReader(br *j2gort.Ext_java_io_BufferedReader) j2gort.Ext_java_util_Map {
	_ = br
	var resultMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = resultMap
	var val_4 IdMappingUtil_MappingReader
	InitIdMappingUtil_MappingReader(&val_4, br)
	var mr *IdMappingUtil_MappingReader = &val_4
	_ = mr
	var currentSectionMap j2gort.Ext_com_google_common_collect_BiMap = nil
	_ = currentSectionMap
	var sectionNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = sectionNames
	for mr.Next() {
		if mr.IsSection() {
			var sectionName string = mr.GetSection()
			_ = sectionName
			if !(sectionNames.Add(sectionName)) {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Cannot parse id map: Duplicate section %s\n lineIndex: %d", []any{j2gort.BoxString(sectionName), j2gort.ValueOfInt(mr.GetLineIndex())})))))
			}
			currentSectionMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.Ext_com_google_common_collect_HashBiMap_Create_0())
			j2gort.MapPutStr(resultMap, sectionName, currentSectionMap)
		} else {
			if currentSectionMap == nil {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Mapping entry outside of section")))
			}
			j2gort.MapPutStr(currentSectionMap, mr.GetKey(), j2gort.BoxString(mr.GetValue()))
		}
	}
	return resultMap
}

// Java static method: IdMappingUtil.parseSerializedIdMappings(java.lang.String idMappings)
func IdMappingUtil_ParseSerializedIdMappings_1_string(idMappings string) j2gort.Ext_java_util_Map {
	_ = idMappings
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(idMappings) {
		return j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
	}
	var didRet_5 bool
	var retVal_6 j2gort.Ext_java_util_Map
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javaio.IIOException](jex_8) {
					e := j2gort.AsIface[javaio.IIOException](jex_8)
					_ = e
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Failed to parse id mappings", e)))
				} else {
					panic(rec_7)
				}
			}
		}()
		var br *j2gort.Ext_java_io_BufferedReader = j2gort.NewExt_java_io_BufferedReader_1(j2gort.NewExt_java_io_StringReader(idMappings))
		_ = br
		defer j2gort.CloseResource(br)
		retVal_6 = IdMappingUtil_ParseSerializedIdMappings_1_Ext_java_io_BufferedReader(br)
		didRet_5 = true
		return
	}()
	if didRet_5 {
		return retVal_6
	}
	panic("unreachable")
}

// Java class: com.google.javascript.jscomp.IdMappingUtil.MappingReader
type IdMappingUtil_MappingReader struct {
	reader         *j2gort.Ext_java_io_BufferedReader `gson:"reader"`
	currentSection string                             `gson:"currentSection"`
	currentKey     string                             `gson:"currentKey"`
	currentValue   string                             `gson:"currentValue"`
	lineIndex      int32                              `gson:"lineIndex"`
	idHash_        atomic.Int32
}

func (this_ *IdMappingUtil_MappingReader) BindSelf_(s any) {
	_ = s
}
func (this_ *IdMappingUtil_MappingReader) IsJavaNil_() bool { return this_ == nil }
func (*IdMappingUtil_MappingReader) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.IdMappingUtil$MappingReader", "MappingReader"
}
func (this_ *IdMappingUtil_MappingReader) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IdMappingUtil_MappingReader) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.IdMappingUtil.MappingReader", "com.google.javascript.jscomp.IdMappingUtil$MappingReader", func(v_ any) bool {
		return j2gort.InstanceOf[*IdMappingUtil_MappingReader](v_)
	})
}

func (this_ *IdMappingUtil_MappingReader) initFields_IdMappingUtil_MappingReader() {
	_ = this_
	this_.currentSection = j2gort.NullString
	this_.currentKey = j2gort.NullString
	this_.currentValue = j2gort.NullString

}

// Java constructor: MappingReader(java.io.BufferedReader reader)
func NewIdMappingUtil_MappingReader(reader *j2gort.Ext_java_io_BufferedReader) *IdMappingUtil_MappingReader {
	this_ := &IdMappingUtil_MappingReader{}
	InitIdMappingUtil_MappingReader(this_, reader)
	return this_
}

func InitIdMappingUtil_MappingReader(this_ *IdMappingUtil_MappingReader, reader *j2gort.Ext_java_io_BufferedReader) {
	_ = reader
	this_.initFields_IdMappingUtil_MappingReader()
	this_.reader = reader
}

// Java method: MappingReader.next()
func (this_ *IdMappingUtil_MappingReader) Next() bool {
	var line string = j2gort.NullString
	_ = line
	for (func() string { var asgn_9 string = this_.reader.ReadLine(); line = asgn_9; return asgn_9 }()) != j2gort.NullString {
		this_.lineIndex = this_.lineIndex + 1
		if j2gort.CheckNotNullString(line) == "" {
			continue
		}
		if int32(j2gort.JStrCharAt(line, int32(0))) == int32(uint16(91)) {
			this_.currentSection = j2gort.JStrSub(line, int32(1), (j2gort.JStrLen(line) - int32(1)))
			this_.currentKey = j2gort.NullString
			this_.currentValue = j2gort.NullString
			return true
		} else {
			var split int32 = j2gort.JStrIndexOfCp(line, int32(uint16(58)))
			_ = split
			if split != int32(-1) {
				this_.currentKey = j2gort.JStrSub(line, int32(0), split)
				this_.currentValue = j2gort.JStrSubFrom(line, (split + int32(1)))
				return true
			} else {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Cannot parse id map.\n Line: %s, lineIndex: %d", []any{j2gort.BoxString(line), j2gort.ValueOfInt(this_.lineIndex)})))))
			}
		}
	}
	return false
}

// Java method: MappingReader.isSection()
func (this_ *IdMappingUtil_MappingReader) IsSection() bool {
	return (this_.currentKey == j2gort.NullString)
}

// Java method: MappingReader.getSection()
func (this_ *IdMappingUtil_MappingReader) GetSection() string {
	return this_.currentSection
}

// Java method: MappingReader.getKey()
func (this_ *IdMappingUtil_MappingReader) GetKey() string {
	return this_.currentKey
}

// Java method: MappingReader.getValue()
func (this_ *IdMappingUtil_MappingReader) GetValue() string {
	return this_.currentValue
}

// Java method: MappingReader.getLineIndex()
func (this_ *IdMappingUtil_MappingReader) GetLineIndex() int32 {
	return this_.lineIndex
}

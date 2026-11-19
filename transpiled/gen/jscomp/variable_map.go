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
// Port of com/google/javascript/jscomp/VariableMap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.VariableMap
type VariableMap struct {
	map_    j2gort.Ext_com_google_common_collect_ImmutableBiMap `gson:"map"`
	idHash_ atomic.Int32
}

func (this_ *VariableMap) BindSelf_(s any) {
	_ = s
}
func (this_ *VariableMap) IsJavaNil_() bool { return this_ == nil }
func (*VariableMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VariableMap", "VariableMap"
}
func (this_ *VariableMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VariableMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VariableMap", "com.google.javascript.jscomp.VariableMap", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableMap](v_)
	})
}

func StaticGet_VariableMap_SEPARATOR() uint16 { return uint16(58) }
func (this_ *VariableMap) initFields_VariableMap() {
	_ = this_
}

// Java constructor: VariableMap(java.util.Map<java.lang.String,java.lang.String> map)
func NewVariableMap(map_ j2gort.Ext_java_util_Map) *VariableMap {
	this_ := &VariableMap{}
	InitVariableMap(this_, map_)
	return this_
}

func InitVariableMap(this_ *VariableMap, map_ j2gort.Ext_java_util_Map) {
	_ = map_
	this_.initFields_VariableMap()
	this_.map_ = j2gort.Ext_com_google_common_collect_ImmutableBiMap_CopyOf_1_Ext_java_util_Map(map_)
}

// Java method: VariableMap.lookupNewName(java.lang.String sourceName)
func (this_ *VariableMap) LookupNewName(sourceName string) string {
	_ = sourceName
	return j2gort.CastString(j2gort.MapGetStr(this_.map_, sourceName))
}

// Java method: VariableMap.lookupSourceName(java.lang.String newName)
func (this_ *VariableMap) LookupSourceName(newName string) string {
	_ = newName
	return j2gort.CastString(j2gort.MapGetStr(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableBiMap](this_.map_.Inverse()), newName))
}

// Java method: VariableMap.getOriginalNameToNewNameMap()
func (this_ *VariableMap) GetOriginalNameToNewNameMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.Ext_com_google_common_collect_ImmutableSortedMap_CopyOf_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.map_)))
}

// Java method: VariableMap.getNewNameToOriginalNameMap()
func (this_ *VariableMap) GetNewNameToOriginalNameMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableBiMap](this_.map_.Inverse())
}

// Java method: VariableMap.save(java.lang.String filename)
func (this_ *VariableMap) Save(filename string) {
	_ = filename
	j2gort.Ext_com_google_common_io_Files_Write_2(this_.ToBytes(), j2gort.NewExt_java_io_File_1_string(filename))
}

// Java static method: VariableMap.load(java.lang.String filename)
func VariableMap_Load(filename string) *VariableMap {
	_ = filename
	var didRet_1 bool
	var retVal_2 *VariableMap
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[j2gort.Ext_java_text_ParseException](jex_4) {
					e := j2gort.AsIface[j2gort.Ext_java_text_ParseException](jex_4)
					_ = e
					panic(j2gort.AsIface[javaio.IIOException](j2gort.ConstructShimThrowable("java.io.IOException", e)))
				} else {
					panic(rec_3)
				}
			}
		}()
		retVal_2 = VariableMap_FromBytes(j2gort.Ext_com_google_common_io_Files_ToByteArray(j2gort.NewExt_java_io_File_1_string(filename)))
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java method: VariableMap.toBytes()
func (this_ *VariableMap) ToBytes() []int8 {
	var baos *j2gort.Ext_java_io_ByteArrayOutputStream = j2gort.NewExt_java_io_ByteArrayOutputStream()
	_ = baos
	var writer j2gort.Ext_java_io_Writer = j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_Charset(baos, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	_ = writer
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[javaio.IIOException](jex_6) {
					e := j2gort.AsIface[javaio.IIOException](jex_6)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_5)
				}
			}
		}()
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_CopyOf_2_Ext_java_util_Comparator_Ext_java_util_Collection(j2gort.Ext_java_util_Map_Entry_ComparingByKey_0(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.map_.EntrySet()))))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](rawElem_8)
				_ = entry
				writer.Write_1_string(VariableMap_Escape(j2gort.CastString(entry.GetKey())))
				writer.Write_1_int32(int32(uint16(58)))
				writer.Write_1_string(VariableMap_Escape(j2gort.CastString(entry.GetValue())))
				writer.Write_1_int32(int32(uint16(10)))
			}
		}
		writer.Close_0()
	}()
	return baos.ToByteArray()
}

// Java static method: VariableMap.fromBytes(byte[] bytes)
func VariableMap_FromBytes(bytes []int8) *VariableMap {
	_ = bytes
	return VariableMap_FromLines(func() string {
		b_9 := bytes
		return j2gort.NewStringFromBytes(b_9, 0, int32(len(b_9)), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	}())
}

// Java static method: VariableMap.fromStream(java.io.InputStream stream)
func VariableMap_FromStream(stream j2gort.Ext_java_io_InputStream) *VariableMap {
	_ = stream
	var map_ j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = map_
	var br *j2gort.Ext_java_io_BufferedReader = j2gort.NewExt_java_io_BufferedReader_1(j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(stream, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8))
	_ = br
	var line string = j2gort.NullString
	_ = line
	for (func() string { var asgn_10 string = br.ReadLine(); line = asgn_10; return asgn_10 }()) != j2gort.NullString {
		if j2gort.CheckNotNullString(line) == "" {
			continue
		}
		var pos int32 = VariableMap_FindIndexOfUnescapedChar(line, uint16(58))
		_ = pos
		if pos <= int32(0) {
			panic(j2gort.NewExt_java_text_ParseException(("Bad line: " + j2gort.StrValString(line)), int32(0)))
		}
		map_.Put_2(j2gort.BoxString(VariableMap_Unescape(j2gort.JStrSub(line, int32(0), pos))), j2gort.BoxString(func() string {
			if pos == (j2gort.JStrLen(line) - int32(1)) {
				return ""
			}
			return VariableMap_Unescape(j2gort.JStrSubFrom(line, (pos + int32(1))))
		}()))
	}
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](map_.BuildOrThrow()))
}

// Java static method: VariableMap.fromLines(java.lang.String string)
func VariableMap_FromLines(string_ string) *VariableMap {
	_ = string_
	var map_ j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = map_
	var startOfLine int32 = int32(0)
	_ = startOfLine
	for startOfLine < j2gort.JStrLen(string_) {
		var newLine int32 = j2gort.JStrIndexOfCpFrom(string_, int32(uint16(10)), startOfLine)
		_ = newLine
		if newLine == int32(-1) {
			newLine = j2gort.JStrLen(string_)
		}
		var endOfLine int32 = newLine
		_ = endOfLine
		if (newLine > int32(0)) && (int32(j2gort.JStrCharAt(string_, (newLine-int32(1)))) == int32(uint16(13))) {
			newLine = newLine - 1
		}
		var line string = j2gort.JStrSub(string_, startOfLine, newLine)
		_ = line
		startOfLine = (endOfLine + int32(1))
		if j2gort.CheckNotNullString(line) == "" {
			continue
		}
		var pos int32 = VariableMap_FindIndexOfUnescapedChar(line, uint16(58))
		_ = pos
		if pos <= int32(0) {
			panic(j2gort.NewExt_java_text_ParseException(("Bad line: " + j2gort.StrValString(line)), int32(0)))
		}
		map_.Put_2(j2gort.BoxString(VariableMap_Unescape(j2gort.JStrSub(line, int32(0), pos))), j2gort.BoxString(func() string {
			if pos == (j2gort.JStrLen(line) - int32(1)) {
				return ""
			}
			return VariableMap_Unescape(j2gort.JStrSubFrom(line, (pos + int32(1))))
		}()))
	}
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](map_.BuildOrThrow()))
}

// Java static method: VariableMap.escape(java.lang.String value)
func VariableMap_Escape(value string) string {
	_ = value
	return javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(value, "\\", "\\\\"), ":", "\\:"), "\n", "\\n")
}

// Java static method: VariableMap.findIndexOfUnescapedChar(java.lang.String value, char stopChar)
func VariableMap_FindIndexOfUnescapedChar(value string, stopChar uint16) int32 {
	_ = value
	_ = stopChar
	var len_ int32 = j2gort.JStrLen(value)
	_ = len_
	{
		var i int32 = int32(0)
		_ = i
		for i < len_ {
			var stopCharIndex int32 = j2gort.JStrIndexOfCpFrom(value, int32(stopChar), i)
			_ = stopCharIndex
			if stopCharIndex == int32(-1) {
				return int32(-1)
			}
			if (stopCharIndex == int32(0)) || (int32(j2gort.JStrCharAt(value, (stopCharIndex-int32(1)))) != int32(uint16(92))) {
				return stopCharIndex
			}
			i = (stopCharIndex + int32(1))
		}
	}
	return int32(-1)
}

// Java static method: VariableMap.unescape(java.lang.String value)
func VariableMap_Unescape(value string) string {
	_ = value
	var slashIndex int32 = j2gort.JStrIndexOfCp(value, int32(uint16(92)))
	_ = slashIndex
	if slashIndex == int32(-1) {
		return value
	}
	var sb j2gort.LocalStringBuilder
	sb.InitCap((j2gort.JStrLen(value) - int32(1)))

	sb.AppendSub(j2gort.BoxString(value), int32(0), slashIndex)
	var len_ int32 = j2gort.JStrLen(value)
	_ = len_
	{
		var i int32 = slashIndex
		_ = i
		for ; i < len_; i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(value, i)
			_ = c
			if (int32(c) == int32(uint16(92))) && (func() int32 { old_11 := i; nw_12 := old_11 + 1; i = nw_12; return nw_12 }() < len_) {
				c = j2gort.JStrCharAt(value, i)
			}
			sb.AppendChar(c)
		}
	}
	return sb.String()
}

// Java static method: VariableMap.fromMap(java.util.Map<java.lang.String,java.lang.String> map)
func VariableMap_FromMap(map_ j2gort.Ext_java_util_Map) *VariableMap {
	_ = map_
	return NewVariableMap(map_)
}

// Java method: VariableMap.toMap()
func (this_ *VariableMap) ToMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](this_.map_)
}

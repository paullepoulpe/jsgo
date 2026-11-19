/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/SourceMapGeneratorV3.java (generated mechanically by j2go).

package sourcemap

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3
type SourceMapGeneratorV3 struct {
	mappings             []*SourceMapGeneratorV3_Mapping     `gson:"mappings"`
	sourceFileMap        *j2gort.Ext_java_util_LinkedHashMap `gson:"sourceFileMap"`
	sourceFileContentMap j2gort.OrderedMap[string, string]   `gson:"sourceFileContentMap"`
	originalNameMap      *j2gort.Ext_java_util_LinkedHashMap `gson:"originalNameMap"`
	lastSourceFile       string                              `gson:"lastSourceFile"`
	lastMapping          *SourceMapGeneratorV3_Mapping       `gson:"lastMapping"`
	extensions           *j2gort.Ext_java_util_LinkedHashMap `gson:"extensions"`
	sourceRootPath       string                              `gson:"sourceRootPath"`
	lastSourceFileIndex  int32                               `gson:"lastSourceFileIndex"`
	offsetPositionLine   int32                               `gson:"offsetPositionLine"`
	offsetPositionColumn int32                               `gson:"offsetPositionColumn"`
	prefixPositionLine   int32                               `gson:"prefixPositionLine"`
	prefixPositionColumn int32                               `gson:"prefixPositionColumn"`
	idHash_              atomic.Int32
}

func (this_ *SourceMapGeneratorV3) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3", "SourceMapGeneratorV3"
}
func (this_ *SourceMapGeneratorV3) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3", "com.google.debugging.sourcemap.SourceMapGeneratorV3", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3](v_)
	})
}
func (*SourceMapGeneratorV3) Implements_sourcemap_SourceMapGenerator() {}

func StaticGet_SourceMapGeneratorV3_UNMAPPED() int32 { return int32(-1) }
func (this_ *SourceMapGeneratorV3) initFields_SourceMapGeneratorV3() {
	_ = this_
	this_.mappings = []*SourceMapGeneratorV3_Mapping{}
	this_.sourceFileMap = j2gort.NewExt_java_util_LinkedHashMap()
	this_.sourceFileContentMap = j2gort.NewOrderedMap[string, string]()
	this_.originalNameMap = j2gort.NewExt_java_util_LinkedHashMap()
	this_.lastSourceFile = j2gort.NullString
	this_.lastSourceFileIndex = int32(-1)

	this_.extensions = j2gort.NewExt_java_util_LinkedHashMap()
	this_.sourceRootPath = j2gort.NullString
}

// Java constructor: SourceMapGeneratorV3()
func NewSourceMapGeneratorV3() *SourceMapGeneratorV3 {
	this_ := &SourceMapGeneratorV3{}
	InitSourceMapGeneratorV3(this_)
	return this_
}

func InitSourceMapGeneratorV3(this_ *SourceMapGeneratorV3) {
	this_.initFields_SourceMapGeneratorV3()
}

// Java method: SourceMapGeneratorV3.reset()
func (this_ *SourceMapGeneratorV3) Reset() {
	clear(this_.mappings)
	this_.mappings = this_.mappings[:0]
	this_.lastMapping = nil
	this_.sourceFileMap.Clear()
	this_.originalNameMap.Clear()
	this_.lastSourceFile = j2gort.NullString
	this_.lastSourceFileIndex = int32(-1)
	this_.offsetPositionLine = int32(0)
	this_.offsetPositionColumn = int32(0)
	this_.prefixPositionLine = int32(0)
	this_.prefixPositionColumn = int32(0)
}

// Java method: SourceMapGeneratorV3.validate(boolean validate)
func (this_ *SourceMapGeneratorV3) Validate(validate bool) {
	_ = validate
}

// Java method: SourceMapGeneratorV3.setWrapperPrefix(java.lang.String prefix)
func (this_ *SourceMapGeneratorV3) SetWrapperPrefix(prefix string) {
	_ = prefix
	var prefixLine int32 = int32(0)
	_ = prefixLine
	var prefixIndex int32 = int32(0)
	_ = prefixIndex
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(prefix); i = i + 1 {
			if int32(j2gort.JStrCharAt(prefix, i)) == int32(uint16(10)) {
				prefixLine = prefixLine + 1
				prefixIndex = int32(0)
			} else {
				prefixIndex = prefixIndex + 1
			}
		}
	}
	this_.prefixPositionLine = prefixLine
	this_.prefixPositionColumn = prefixIndex
}

// Java method: SourceMapGeneratorV3.setStartingPosition(int offsetLine, int offsetIndex)
func (this_ *SourceMapGeneratorV3) SetStartingPosition(offsetLine int32, offsetIndex int32) {
	_ = offsetLine
	_ = offsetIndex
	if !(offsetLine >= int32(0)) {
		j2gort.CheckState(false)
	}
	if !(offsetIndex >= int32(0)) {
		j2gort.CheckState(false)
	}
	this_.offsetPositionLine = offsetLine
	this_.offsetPositionColumn = offsetIndex
}

// Java method: SourceMapGeneratorV3.addMapping(java.lang.String sourceName, java.lang.@org.jspecify.annotations.Nullable String symbolName, com.google.debugging.sourcemap.FilePosition sourceStartPosition, com.google.debugging.sourcemap.FilePosition startPosition, com.google.debugging.sourcemap.FilePosition endPosition)
func (this_ *SourceMapGeneratorV3) AddMapping_5(sourceName string, symbolName string, sourceStartPosition *FilePosition, startPosition *FilePosition, endPosition *FilePosition) {
	_ = sourceName
	_ = symbolName
	_ = sourceStartPosition
	_ = startPosition
	_ = endPosition
	this_.AddMapping_8(sourceName, symbolName, sourceStartPosition.GetLine(), sourceStartPosition.GetColumn(), startPosition.GetLine(), startPosition.GetColumn(), endPosition.GetLine(), endPosition.GetColumn())
}

// Java method: SourceMapGeneratorV3.addMapping(java.lang.String sourceName, java.lang.@org.jspecify.annotations.Nullable String symbolName, int sourceStartLine, int sourceStartColumn, int startLine, int startColumn, int endLine, int endColumn)
func (this_ *SourceMapGeneratorV3) AddMapping_8(sourceName string, symbolName string, sourceStartLine int32, sourceStartColumn int32, startLine int32, startColumn int32, endLine int32, endColumn int32) {
	_ = sourceName
	_ = symbolName
	_ = sourceStartLine
	_ = sourceStartColumn
	_ = startLine
	_ = startColumn
	_ = endLine
	_ = endColumn
	if (sourceName == j2gort.NullString) || (sourceStartLine < int32(0)) {
		return
	}
	var adjustedStartLine int32 = startLine
	_ = adjustedStartLine
	var adjustedStartColumn int32 = startColumn
	_ = adjustedStartColumn
	var adjustedEndLine int32 = endLine
	_ = adjustedEndLine
	var adjustedEndColumn int32 = endColumn
	_ = adjustedEndColumn
	if (this_.offsetPositionLine != int32(0)) || (this_.offsetPositionColumn != int32(0)) {
		var offsetLine int32 = this_.offsetPositionLine
		_ = offsetLine
		var startOffsetPosition int32
		_ = startOffsetPosition
		if startLine > int32(0) {
			startOffsetPosition = int32(0)
		} else {
			startOffsetPosition = this_.offsetPositionColumn
		}
		var endOffsetPosition int32
		_ = endOffsetPosition
		if endLine > int32(0) {
			endOffsetPosition = int32(0)
		} else {
			endOffsetPosition = this_.offsetPositionColumn
		}
		adjustedStartLine = (startLine + offsetLine)
		adjustedStartColumn = (startColumn + startOffsetPosition)
		adjustedEndLine = (endLine + offsetLine)
		adjustedEndColumn = (endColumn + endOffsetPosition)
	}
	var val_1 SourceMapGeneratorV3_Mapping
	InitSourceMapGeneratorV3_Mapping(&val_1)
	var mapping *SourceMapGeneratorV3_Mapping = &val_1
	_ = mapping
	{
		r_2 := mapping
		r_2.sourceFile = sourceName
	}
	{
		r_3 := mapping
		r_3.originalLine = sourceStartLine
	}
	{
		r_4 := mapping
		r_4.originalColumn = sourceStartColumn
	}
	{
		r_5 := mapping
		r_5.originalName = symbolName
	}
	{
		r_6 := mapping
		r_6.startLine = adjustedStartLine
	}
	{
		r_7 := mapping
		r_7.startColumn = adjustedStartColumn
	}
	{
		r_8 := mapping
		r_8.endLine = adjustedEndLine
	}
	{
		r_9 := mapping
		r_9.endColumn = adjustedEndColumn
	}
	if this_.lastMapping != nil {
		var lastLine int32 = this_.lastMapping.startLine
		_ = lastLine
		var lastColumn int32 = this_.lastMapping.startColumn
		_ = lastColumn
		var nextLine int32 = mapping.startLine
		_ = nextLine
		var nextColumn int32 = mapping.startColumn
		_ = nextColumn
		if !((nextLine > lastLine) || ((nextLine == lastLine) && (nextColumn >= lastColumn))) {
			j2gort.CheckStateFmt(false, "Incorrect source mappings order, previous : (%s,%s)\nnew : (%s,%s)", j2gort.ValueOfInt(lastLine), j2gort.ValueOfInt(lastColumn), j2gort.ValueOfInt(nextLine), j2gort.ValueOfInt(nextColumn))
		}
	}
	this_.lastMapping = mapping
	this_.mappings = append(this_.mappings, mapping)
}

// Java method: SourceMapGeneratorV3.addSourcesContent(java.lang.String source, java.lang.String content)
func (this_ *SourceMapGeneratorV3) AddSourcesContent(source string, content string) {
	_ = source
	_ = content
	this_.sourceFileContentMap.Put(source, content)
}

// Java method: SourceMapGeneratorV3.mergeMapSection(int line, int column, java.lang.String mapSectionContents)
func (this_ *SourceMapGeneratorV3) MergeMapSection_3(line int32, column int32, mapSectionContents string) {
	_ = line
	_ = column
	_ = mapSectionContents
	this_.SetStartingPosition(line, column)
	var val_10 SourceMapConsumerV3
	InitSourceMapConsumerV3(&val_10)
	var section *SourceMapConsumerV3 = &val_10
	_ = section
	section.Parse_1(mapSectionContents)
	section.VisitMappings(NewSourceMapGeneratorV3_ConsumerEntryVisitor(this_))
}

// Java method: SourceMapGeneratorV3.mergeMapSection(int line, int column, java.lang.String mapSectionContents, com.google.debugging.sourcemap.SourceMapGeneratorV3.ExtensionMergeAction mergeAction)
func (this_ *SourceMapGeneratorV3) MergeMapSection_4(line int32, column int32, mapSectionContents string, mergeAction SourceMapGeneratorV3_ExtensionMergeAction) {
	_ = line
	_ = column
	_ = mapSectionContents
	_ = mergeAction
	this_.SetStartingPosition(line, column)
	var val_11 SourceMapConsumerV3
	InitSourceMapConsumerV3(&val_11)
	var section *SourceMapConsumerV3 = &val_11
	_ = section
	section.Parse_1(mapSectionContents)
	section.VisitMappings(NewSourceMapGeneratorV3_ConsumerEntryVisitor(this_))
	{
		cur_12 := j2gort.NewMapEntryCursor(section.GetExtensions())
		for cur_12.NextEntry() {
			var extensionKey string = cur_12.KeyStr()
			_ = extensionKey
			if j2gort.MapContainsKeyStr(this_.extensions, extensionKey) {
				j2gort.MapPutStr(this_.extensions, extensionKey, mergeAction.Merge(extensionKey, j2gort.MapGetStr(this_.extensions, extensionKey), cur_12.GetValue()))
			} else {
				j2gort.MapPutStr(this_.extensions, extensionKey, cur_12.GetValue())
			}
		}
	}
}

// Java method: SourceMapGeneratorV3.appendTo(java.lang.Appendable out, java.lang.@org.jspecify.annotations.Nullable String name)
func (this_ *SourceMapGeneratorV3) AppendTo(out j2gort.Ext_java_lang_Appendable, name string) {
	_ = out
	_ = name
	var maxLine int32 = (this_.priv_SourceMapGeneratorV3_prepMappings_69() + int32(1))
	_ = maxLine
	out.Append_1_any("{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "version", "3")
	if name != j2gort.NullString {
		SourceMapGeneratorV3_AppendField(out, "file", j2gort.BoxString(SourceMapGeneratorV3_EscapeString(name)))
	}
	SourceMapGeneratorV3_AppendField(out, "lineCount", strconv.FormatInt(int64(maxLine), 10))
	if (this_.sourceRootPath != j2gort.NullString) && !(j2gort.CheckNotNullString(this_.sourceRootPath) == "") {
		SourceMapGeneratorV3_AppendField(out, "sourceRoot", j2gort.BoxString(SourceMapGeneratorV3_EscapeString(this_.sourceRootPath)))
	}
	SourceMapGeneratorV3_AppendFieldStart_2(out, "mappings")
	{
		var val_13 SourceMapGeneratorV3_LineMapper
		val_13.outer_ = this_
		InitSourceMapGeneratorV3_LineMapper(&val_13, out, maxLine)
		(&val_13).AppendLineMappings()
	}
	SourceMapGeneratorV3_AppendFieldEnd(out)
	SourceMapGeneratorV3_AppendFieldStart_2(out, "sources")
	out.Append_1_any("[")
	this_.priv_SourceMapGeneratorV3_addSourceNameMap_66(out)
	out.Append_1_any("]")
	SourceMapGeneratorV3_AppendFieldEnd(out)
	this_.priv_SourceMapGeneratorV3_addSourcesContentMap_67(out)
	SourceMapGeneratorV3_AppendFieldStart_2(out, "names")
	out.Append_1_any("[")
	this_.priv_SourceMapGeneratorV3_addSymbolNameMap_68(out)
	out.Append_1_any("]")
	SourceMapGeneratorV3_AppendFieldEnd(out)
	{
		cur_14 := j2gort.NewMapKeyStrCursor(this_.extensions)
		for cur_14.HasNext() {
			key := cur_14.NextStr()
			_ = key
			var objValue any = j2gort.MapGetStr(this_.extensions, key)
			_ = objValue
			var value string = j2gort.NullString
			_ = value
			var string_ string
			_ = string_
			var ok_15 bool
			string_, ok_15 = j2gort.TryCastString(objValue)
			if ok_15 {
				value = SourceMapGeneratorV3_EscapeString(string_)
			} else {
				value = javalang.ToJavaString(j2gort.CheckNotNull(objValue))
			}
			SourceMapGeneratorV3_AppendField(out, key, j2gort.BoxString(value))
		}
	}
	out.Append_1_any("\n}\n")
}

// Java method: SourceMapGeneratorV3.setSourceRoot(java.lang.String path)
func (this_ *SourceMapGeneratorV3) SetSourceRoot(path string) {
	_ = path
	this_.sourceRootPath = path
}

// Java method: SourceMapGeneratorV3.addExtension(java.lang.String name, java.lang.Object object)
func (this_ *SourceMapGeneratorV3) AddExtension(name string, object any) {
	_ = name
	_ = object
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(name), "x_")) {
		panic(NewSourceMapParseException(("Extension '" + j2gort.StrValString(name) + "' must start with 'x_'")))
	}
	j2gort.MapPutStr(this_.extensions, name, object)
}

// Java method: SourceMapGeneratorV3.removeExtension(java.lang.String name)
func (this_ *SourceMapGeneratorV3) RemoveExtension(name string) {
	_ = name
	if j2gort.MapContainsKeyStr(this_.extensions, name) {
		j2gort.MapRemoveStr(this_.extensions, name)
	}
}

// Java method: SourceMapGeneratorV3.hasExtension(java.lang.String name)
func (this_ *SourceMapGeneratorV3) HasExtension(name string) bool {
	_ = name
	return j2gort.MapContainsKeyStr(this_.extensions, name)
}

// Java method: SourceMapGeneratorV3.getExtension(java.lang.String name)
func (this_ *SourceMapGeneratorV3) GetExtension(name string) any {
	_ = name
	return j2gort.MapGetStr(this_.extensions, name)
}

// Java method: SourceMapGeneratorV3.addSourceNameMap(java.lang.Appendable out)
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_addSourceNameMap_66(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	SourceMapGeneratorV3_AddNameMap(out, j2gort.AsIface[j2gort.Ext_java_util_Map](this_.sourceFileMap))
}

// Java method: SourceMapGeneratorV3.addSourcesContentMap(java.lang.Appendable out)
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_addSourcesContentMap_67(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	var found bool = false
	_ = found
	var size int32 = this_.sourceFileMap.Size()
	_ = size
	var contents []string = j2gort.MakeNonNilSlice[string](size)
	_ = contents
	contents = j2gort.AppendStrCol(contents, j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Collections_NCopies(size, "")))
	{
		cur_16 := j2gort.NewMapEntryCursor(this_.sourceFileMap)
		for cur_16.NextEntry() {
			var index *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](cur_16.GetValue())
			_ = index
			if !(j2gort.UnboxInt(index) < size) {
				j2gort.CheckState(false)
			}
			var content string = j2gort.OrderedMapGetStr(this_.sourceFileContentMap, cur_16.KeyStr())
			_ = content
			if content != j2gort.NullString {
				contents[j2gort.UnboxInt(index)] = content
				found = true
			}
		}
	}
	if !(found) {
		return
	}
	SourceMapGeneratorV3_AppendFieldStart_2(out, "sourcesContent")
	out.Append_1_any("[")
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			if i != int32(0) {
				out.Append_1_any(",")
			}
			var sourceContent string = contents[i]
			_ = sourceContent
			out.Append_1_any(j2gort.BoxString(SourceMapGeneratorV3_EscapeString(j2gort.Ext_com_google_common_base_Strings_NullToEmpty(sourceContent))))
		}
	}
	out.Append_1_any("]")
	SourceMapGeneratorV3_AppendFieldEnd(out)
}

// Java method: SourceMapGeneratorV3.addSymbolNameMap(java.lang.Appendable out)
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_addSymbolNameMap_68(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	SourceMapGeneratorV3_AddNameMap(out, j2gort.AsIface[j2gort.Ext_java_util_Map](this_.originalNameMap))
}

// Java static method: SourceMapGeneratorV3.addNameMap(java.lang.Appendable out, java.util.Map<java.lang.String,java.lang.Integer> map)
func SourceMapGeneratorV3_AddNameMap(out j2gort.Ext_java_lang_Appendable, map_ j2gort.Ext_java_util_Map) {
	_ = out
	_ = map_
	var i int32 = int32(0)
	_ = i
	{
		cur_17 := j2gort.NewMapEntryCursor(map_)
		for cur_17.NextEntry() {
			var key string = cur_17.KeyStr()
			_ = key
			if i != int32(0) {
				out.Append_1_any(",")
			}
			out.Append_1_any(j2gort.BoxString(SourceMapGeneratorV3_EscapeString(key)))
			i = i + 1
		}
	}
}

// Java static method: SourceMapGeneratorV3.escapeString(java.lang.String value)
func SourceMapGeneratorV3_EscapeString(value string) string {
	_ = value
	return Util_EscapeString(value)
}

// Java static method: SourceMapGeneratorV3.appendFirstField(java.lang.Appendable out, java.lang.String name, java.lang.CharSequence value)
func SourceMapGeneratorV3_AppendFirstField(out j2gort.Ext_java_lang_Appendable, name string, value any) {
	_ = out
	_ = name
	_ = value
	SourceMapGeneratorV3_AppendFieldStart_3(out, name, true)
	out.Append_1_any(value)
}

// Java static method: SourceMapGeneratorV3.appendField(java.lang.Appendable out, java.lang.String name, java.lang.CharSequence value)
func SourceMapGeneratorV3_AppendField(out j2gort.Ext_java_lang_Appendable, name string, value any) {
	_ = out
	_ = name
	_ = value
	SourceMapGeneratorV3_AppendFieldStart_3(out, name, false)
	out.Append_1_any(value)
}

// Java static method: SourceMapGeneratorV3.appendFieldStart(java.lang.Appendable out, java.lang.String name)
func SourceMapGeneratorV3_AppendFieldStart_2(out j2gort.Ext_java_lang_Appendable, name string) {
	_ = out
	_ = name
	SourceMapGeneratorV3_AppendFieldStart_3(out, name, false)
}

// Java static method: SourceMapGeneratorV3.appendFieldStart(java.lang.Appendable out, java.lang.String name, boolean first)
func SourceMapGeneratorV3_AppendFieldStart_3(out j2gort.Ext_java_lang_Appendable, name string, first bool) {
	_ = out
	_ = name
	_ = first
	if !(first) {
		out.Append_1_any(",\n")
	}
	out.Append_1_any("\"")
	out.Append_1_any(j2gort.BoxString(name))
	out.Append_1_any("\"")
	out.Append_1_any(":")
}

// Java static method: SourceMapGeneratorV3.appendFieldEnd(java.lang.Appendable out)
func SourceMapGeneratorV3_AppendFieldEnd(out j2gort.Ext_java_lang_Appendable) {
	_ = out
}

// Java method: SourceMapGeneratorV3.prepMappings()
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_prepMappings_69() int32 {
	{
		var val_18 SourceMapGeneratorV3_MappingTraversal
		val_18.outer_ = this_
		InitSourceMapGeneratorV3_MappingTraversal(&val_18)
		(&val_18).Traverse(NewSourceMapGeneratorV3_UsedMappingCheck())
	}
	var id int32 = int32(0)
	_ = id
	var maxLine int32 = int32(0)
	_ = maxLine
	for _, m := range this_.mappings {
		_ = m
		if m.used {
			{
				r_23 := m
				r_23.id = func() int32 { old_21 := id; nw_22 := old_21 + 1; id = nw_22; return old_21 }()
			}
			var endPositionLine int32 = m.endLine
			_ = endPositionLine
			maxLine = max(maxLine, endPositionLine)
		}
	}

	return (maxLine + this_.prefixPositionLine)
}

// Java method: SourceMapGeneratorV3.appendIndexMapTo(java.lang.Appendable out, java.lang.String name, java.util.List<com.google.debugging.sourcemap.SourceMapSection> sections)
func (this_ *SourceMapGeneratorV3) AppendIndexMapTo(out j2gort.Ext_java_lang_Appendable, name string, sections j2gort.Ext_java_util_List) {
	_ = out
	_ = name
	_ = sections
	out.Append_1_any("{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "version", "3")
	SourceMapGeneratorV3_AppendField(out, "file", j2gort.BoxString(SourceMapGeneratorV3_EscapeString(name)))
	SourceMapGeneratorV3_AppendFieldStart_2(out, "sections")
	out.Append_1_any("[\n")
	var first bool = true
	_ = first
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(sections))
		for cur_24.HasNext() {
			rawElem_25 := cur_24.Next()
			section := j2gort.CastPtr[SourceMapSection](rawElem_25)
			_ = section
			if first {
				first = false
			} else {
				out.Append_1_any(",\n")
			}
			out.Append_1_any("{\n")
			SourceMapGeneratorV3_AppendFieldStart_3(out, "offset", true)
			SourceMapGeneratorV3_AppendOffsetValue(out, section.GetLine(), section.GetColumn())
			if section.GetSectionType() == StaticGet_SourceMapSection_SectionType_URL() {
				SourceMapGeneratorV3_AppendField(out, "url", j2gort.BoxString(SourceMapGeneratorV3_EscapeString(section.GetSectionValue())))
			} else if section.GetSectionType() == StaticGet_SourceMapSection_SectionType_MAP() {
				SourceMapGeneratorV3_AppendField(out, "map", j2gort.BoxString(section.GetSectionValue()))
			} else {
				panic(j2gort.AsIface[javaio.IIOException](j2gort.ConstructShimThrowable("java.io.IOException", "Unexpected section type")))
			}
			out.Append_1_any("\n}")
		}
	}
	out.Append_1_any("\n]")
	SourceMapGeneratorV3_AppendFieldEnd(out)
	out.Append_1_any("\n}\n")
}

// Java static method: SourceMapGeneratorV3.appendOffsetValue(java.lang.Appendable out, int line, int column)
func SourceMapGeneratorV3_AppendOffsetValue(out j2gort.Ext_java_lang_Appendable, line int32, column int32) {
	_ = out
	_ = line
	_ = column
	out.Append_1_any("{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "line", strconv.FormatInt(int64(line), 10))
	SourceMapGeneratorV3_AppendField(out, "column", strconv.FormatInt(int64(column), 10))
	out.Append_1_any("\n}")
}

// Java method: SourceMapGeneratorV3.getSourceId(java.lang.String sourceName)
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_getSourceId_70(sourceName string) int32 {
	_ = sourceName
	if !(sourceName == this_.lastSourceFile) {
		this_.lastSourceFile = sourceName
		var index *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](j2gort.MapGetStr(this_.sourceFileMap, sourceName))
		_ = index
		if index != nil {
			this_.lastSourceFileIndex = j2gort.UnboxInt(index)
		} else {
			this_.lastSourceFileIndex = this_.sourceFileMap.Size()
			j2gort.MapPutStr(this_.sourceFileMap, sourceName, j2gort.ValueOfInt(this_.lastSourceFileIndex))
		}
	}
	return this_.lastSourceFileIndex
}

// Java method: SourceMapGeneratorV3.getNameId(java.lang.String symbolName)
func (this_ *SourceMapGeneratorV3) priv_SourceMapGeneratorV3_getNameId_71(symbolName string) int32 {
	_ = symbolName
	var originalNameIndex int32 = int32(0)
	_ = originalNameIndex
	var index *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](j2gort.MapGetStr(this_.originalNameMap, symbolName))
	_ = index
	if index != nil {
		originalNameIndex = j2gort.UnboxInt(index)
	} else {
		originalNameIndex = this_.originalNameMap.Size()
		j2gort.MapPutStr(this_.originalNameMap, symbolName, j2gort.ValueOfInt(originalNameIndex))
	}
	return originalNameIndex
}

// Java interface: com.google.debugging.sourcemap.SourceMapGeneratorV3.ExtensionMergeAction
type SourceMapGeneratorV3_ExtensionMergeAction interface {
	Implements_sourcemap_SourceMapGeneratorV3_ExtensionMergeAction()
	IsJavaNil_() bool
	Merge(string, any, any) any
}

func ToIface_SourceMapGeneratorV3_ExtensionMergeAction[T interface {
	SourceMapGeneratorV3_ExtensionMergeAction
	comparable
}](v T) SourceMapGeneratorV3_ExtensionMergeAction {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceMapGeneratorV3_ExtensionMergeAction[T SourceMapGeneratorV3_ExtensionMergeAction](v T) SourceMapGeneratorV3_ExtensionMergeAction {
	var vi SourceMapGeneratorV3_ExtensionMergeAction = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.ExtensionMergeAction", "com.google.debugging.sourcemap.SourceMapGeneratorV3$ExtensionMergeAction", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapGeneratorV3_ExtensionMergeAction](v_)
	})
}

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3.ConsumerEntryVisitor
type SourceMapGeneratorV3_ConsumerEntryVisitor struct {
	outer_  *SourceMapGeneratorV3
	idHash_ atomic.Int32
}

func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3_ConsumerEntryVisitor) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3$ConsumerEntryVisitor", "ConsumerEntryVisitor"
}
func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3_ConsumerEntryVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.ConsumerEntryVisitor", "com.google.debugging.sourcemap.SourceMapGeneratorV3$ConsumerEntryVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3_ConsumerEntryVisitor](v_)
	})
}
func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) SetOuter_(o *SourceMapGeneratorV3) {
	this_.outer_ = o
}
func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) GetOuter_() *SourceMapGeneratorV3 {
	return this_.outer_
}
func (*SourceMapGeneratorV3_ConsumerEntryVisitor) Implements_sourcemap_SourceMapConsumerV3_EntryVisitor() {
}

func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) initFields_SourceMapGeneratorV3_ConsumerEntryVisitor() {
	_ = this_
}

// Java constructor: ConsumerEntryVisitor()
func NewSourceMapGeneratorV3_ConsumerEntryVisitor(outer_ *SourceMapGeneratorV3) *SourceMapGeneratorV3_ConsumerEntryVisitor {
	this_ := &SourceMapGeneratorV3_ConsumerEntryVisitor{}
	this_.outer_ = outer_
	InitSourceMapGeneratorV3_ConsumerEntryVisitor(this_)
	return this_
}

func InitSourceMapGeneratorV3_ConsumerEntryVisitor(this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) {
	this_.initFields_SourceMapGeneratorV3_ConsumerEntryVisitor()
}

// Java method: ConsumerEntryVisitor.visit(java.lang.String sourceName, java.lang.String symbolName, com.google.debugging.sourcemap.FilePosition sourceStartPosition, com.google.debugging.sourcemap.FilePosition startPosition, com.google.debugging.sourcemap.FilePosition endPosition)
func (this_ *SourceMapGeneratorV3_ConsumerEntryVisitor) Visit(sourceName string, symbolName string, sourceStartPosition *FilePosition, startPosition *FilePosition, endPosition *FilePosition) {
	_ = sourceName
	_ = symbolName
	_ = sourceStartPosition
	_ = startPosition
	_ = endPosition
	this_.outer_.AddMapping_5(sourceName, symbolName, sourceStartPosition, startPosition, endPosition)
}

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping
type SourceMapGeneratorV3_Mapping struct {
	sourceFile     string `gson:"sourceFile"`
	originalName   string `gson:"originalName"`
	id             int32  `gson:"id"`
	originalLine   int32  `gson:"originalLine"`
	originalColumn int32  `gson:"originalColumn"`
	startLine      int32  `gson:"startLine"`
	startColumn    int32  `gson:"startColumn"`
	endLine        int32  `gson:"endLine"`
	endColumn      int32  `gson:"endColumn"`
	idHash_        atomic.Int32
	used           bool `gson:"used"`
}

func (this_ *SourceMapGeneratorV3_Mapping) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3_Mapping) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3_Mapping) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3$Mapping", "Mapping"
}
func (this_ *SourceMapGeneratorV3_Mapping) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3_Mapping) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping", "com.google.debugging.sourcemap.SourceMapGeneratorV3$Mapping", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3_Mapping](v_)
	})
}

func (this_ *SourceMapGeneratorV3_Mapping) initFields_SourceMapGeneratorV3_Mapping() {
	_ = this_
	this_.id = int32(-1)
	this_.sourceFile = j2gort.NullString
	this_.originalName = j2gort.NullString

}

// Java constructor: Mapping()
func NewSourceMapGeneratorV3_Mapping() *SourceMapGeneratorV3_Mapping {
	this_ := &SourceMapGeneratorV3_Mapping{}
	InitSourceMapGeneratorV3_Mapping(this_)
	return this_
}

func InitSourceMapGeneratorV3_Mapping(this_ *SourceMapGeneratorV3_Mapping) {
	this_.initFields_SourceMapGeneratorV3_Mapping()
}

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3.UsedMappingCheck
type SourceMapGeneratorV3_UsedMappingCheck struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapGeneratorV3_UsedMappingCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3_UsedMappingCheck) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3_UsedMappingCheck) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3$UsedMappingCheck", "UsedMappingCheck"
}
func (this_ *SourceMapGeneratorV3_UsedMappingCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3_UsedMappingCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.UsedMappingCheck", "com.google.debugging.sourcemap.SourceMapGeneratorV3$UsedMappingCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3_UsedMappingCheck](v_)
	})
}
func (*SourceMapGeneratorV3_UsedMappingCheck) Implements_sourcemap_SourceMapGeneratorV3_MappingVisitor() {
}

func (this_ *SourceMapGeneratorV3_UsedMappingCheck) initFields_SourceMapGeneratorV3_UsedMappingCheck() {
	_ = this_
}

// Java constructor: UsedMappingCheck()
func NewSourceMapGeneratorV3_UsedMappingCheck() *SourceMapGeneratorV3_UsedMappingCheck {
	this_ := &SourceMapGeneratorV3_UsedMappingCheck{}
	InitSourceMapGeneratorV3_UsedMappingCheck(this_)
	return this_
}

func InitSourceMapGeneratorV3_UsedMappingCheck(this_ *SourceMapGeneratorV3_UsedMappingCheck) {
	this_.initFields_SourceMapGeneratorV3_UsedMappingCheck()
}

// Java method: UsedMappingCheck.visit(com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m, int line, int col, int nextLine, int nextCol)
func (this_ *SourceMapGeneratorV3_UsedMappingCheck) Visit(m *SourceMapGeneratorV3_Mapping, line int32, col int32, nextLine int32, nextCol int32) {
	_ = m
	_ = line
	_ = col
	_ = nextLine
	_ = nextCol
	if m != nil {
		{
			r_26 := m
			r_26.used = true
		}
	}
}

// Java interface: com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor
type SourceMapGeneratorV3_MappingVisitor interface {
	Implements_sourcemap_SourceMapGeneratorV3_MappingVisitor()
	IsJavaNil_() bool
	Visit(*SourceMapGeneratorV3_Mapping, int32, int32, int32, int32)
}

func ToIface_SourceMapGeneratorV3_MappingVisitor[T interface {
	SourceMapGeneratorV3_MappingVisitor
	comparable
}](v T) SourceMapGeneratorV3_MappingVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceMapGeneratorV3_MappingVisitor[T SourceMapGeneratorV3_MappingVisitor](v T) SourceMapGeneratorV3_MappingVisitor {
	var vi SourceMapGeneratorV3_MappingVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor", "com.google.debugging.sourcemap.SourceMapGeneratorV3$MappingVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapGeneratorV3_MappingVisitor](v_)
	})
}

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingTraversal
type SourceMapGeneratorV3_MappingTraversal struct {
	outer_  *SourceMapGeneratorV3
	line    int32 `gson:"line"`
	col     int32 `gson:"col"`
	idHash_ atomic.Int32
}

func (this_ *SourceMapGeneratorV3_MappingTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3_MappingTraversal) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3_MappingTraversal) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3$MappingTraversal", "MappingTraversal"
}
func (this_ *SourceMapGeneratorV3_MappingTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3_MappingTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingTraversal", "com.google.debugging.sourcemap.SourceMapGeneratorV3$MappingTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3_MappingTraversal](v_)
	})
}
func (this_ *SourceMapGeneratorV3_MappingTraversal) SetOuter_(o *SourceMapGeneratorV3) {
	this_.outer_ = o
}
func (this_ *SourceMapGeneratorV3_MappingTraversal) GetOuter_() *SourceMapGeneratorV3 {
	return this_.outer_
}

func (this_ *SourceMapGeneratorV3_MappingTraversal) initFields_SourceMapGeneratorV3_MappingTraversal() {
	_ = this_
}

// Java constructor: MappingTraversal()
func NewSourceMapGeneratorV3_MappingTraversal(outer_ *SourceMapGeneratorV3) *SourceMapGeneratorV3_MappingTraversal {
	this_ := &SourceMapGeneratorV3_MappingTraversal{}
	this_.outer_ = outer_
	InitSourceMapGeneratorV3_MappingTraversal(this_)
	return this_
}

func InitSourceMapGeneratorV3_MappingTraversal(this_ *SourceMapGeneratorV3_MappingTraversal) {
	this_.initFields_SourceMapGeneratorV3_MappingTraversal()
}

// Java method: MappingTraversal.traverse(com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor v)
func (this_ *SourceMapGeneratorV3_MappingTraversal) Traverse(v SourceMapGeneratorV3_MappingVisitor) {
	_ = v
	var stack j2gort.Deque[*SourceMapGeneratorV3_Mapping]
	_ = stack
	for _, m := range this_.outer_.mappings {
		_ = m
		for !(stack.IsEmpty()) && !(this_.priv_MappingTraversal_isOverlapped_77(stack.Peek(), m)) {
			var previous *SourceMapGeneratorV3_Mapping = stack.Pop()
			_ = previous
			this_.priv_MappingTraversal_maybeVisit_78(v, previous)
		}
		var parent *SourceMapGeneratorV3_Mapping = stack.Peek()
		_ = parent
		this_.priv_MappingTraversal_maybeVisitParent_79(v, parent, m)
		stack.Push(m)
	}

	for !(stack.IsEmpty()) {
		var m_2 *SourceMapGeneratorV3_Mapping = stack.Pop()
		_ = m_2
		this_.priv_MappingTraversal_maybeVisit_78(v, m_2)
	}
}

// Java method: MappingTraversal.getAdjustedLine(int rawLine)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_getAdjustedLine_75(rawLine int32) int32 {
	_ = rawLine
	return (rawLine + this_.outer_.prefixPositionLine)
}

// Java method: MappingTraversal.getAdjustedCol(int rawLine, int rawCol)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_getAdjustedCol_76(rawLine int32, rawCol int32) int32 {
	_ = rawLine
	_ = rawCol
	if rawLine != int32(0) {
		return rawCol
	} else {
		return (rawCol + this_.outer_.prefixPositionColumn)
	}
}

// Java method: MappingTraversal.isOverlapped(com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m1, com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m2)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_isOverlapped_77(m1 *SourceMapGeneratorV3_Mapping, m2 *SourceMapGeneratorV3_Mapping) bool {
	_ = m1
	_ = m2
	var l1 int32 = m1.endLine
	_ = l1
	var l2 int32 = m2.startLine
	_ = l2
	var c1 int32 = m1.endColumn
	_ = c1
	var c2 int32 = m2.startColumn
	_ = c2
	return (((l1 == l2) && (c1 >= c2)) || (l1 > l2))
}

// Java method: MappingTraversal.maybeVisit(com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor v, com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_maybeVisit_78(v SourceMapGeneratorV3_MappingVisitor, m *SourceMapGeneratorV3_Mapping) {
	_ = v
	_ = m
	var nextLine int32 = this_.priv_MappingTraversal_getAdjustedLine_75(m.endLine)
	_ = nextLine
	var nextCol int32 = this_.priv_MappingTraversal_getAdjustedCol_76(m.endLine, m.endColumn)
	_ = nextCol
	if (this_.line < nextLine) || ((this_.line == nextLine) && (this_.col < nextCol)) {
		this_.priv_MappingTraversal_visit_80(v, m, nextLine, nextCol)
	}
}

// Java method: MappingTraversal.maybeVisitParent(com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor v, com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping parent, com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_maybeVisitParent_79(v SourceMapGeneratorV3_MappingVisitor, parent *SourceMapGeneratorV3_Mapping, m *SourceMapGeneratorV3_Mapping) {
	_ = v
	_ = parent
	_ = m
	var nextLine int32 = this_.priv_MappingTraversal_getAdjustedLine_75(m.startLine)
	_ = nextLine
	var nextCol int32 = this_.priv_MappingTraversal_getAdjustedCol_76(m.startLine, m.startColumn)
	_ = nextCol
	if !((this_.line < nextLine) || (this_.col <= nextCol)) {
		j2gort.CheckState(false)
	}
	if (this_.line < nextLine) || ((this_.line == nextLine) && (this_.col < nextCol)) {
		this_.priv_MappingTraversal_visit_80(v, parent, nextLine, nextCol)
	}
}

// Java method: MappingTraversal.visit(com.google.debugging.sourcemap.SourceMapGeneratorV3.MappingVisitor v, com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m, int nextLine, int nextCol)
func (this_ *SourceMapGeneratorV3_MappingTraversal) priv_MappingTraversal_visit_80(v SourceMapGeneratorV3_MappingVisitor, m *SourceMapGeneratorV3_Mapping, nextLine int32, nextCol int32) {
	_ = v
	_ = m
	_ = nextLine
	_ = nextCol
	if !(this_.line <= nextLine) {
		j2gort.CheckState(false)
	}
	if !((this_.line < nextLine) || (this_.col < nextCol)) {
		j2gort.CheckState(false)
	}
	if (this_.line == nextLine) && (this_.col == nextCol) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	v.Visit(m, this_.line, this_.col, nextLine, nextCol)
	this_.line = nextLine
	this_.col = nextCol
}

// Java class: com.google.debugging.sourcemap.SourceMapGeneratorV3.LineMapper
type SourceMapGeneratorV3_LineMapper struct {
	outer_               *SourceMapGeneratorV3
	out                  j2gort.Ext_java_lang_Appendable `gson:"out"`
	maxLine              int32                           `gson:"maxLine"`
	previousLine         int32                           `gson:"previousLine"`
	previousColumn       int32                           `gson:"previousColumn"`
	previousSourceFileId int32                           `gson:"previousSourceFileId"`
	previousSourceLine   int32                           `gson:"previousSourceLine"`
	previousSourceColumn int32                           `gson:"previousSourceColumn"`
	previousNameId       int32                           `gson:"previousNameId"`
	idHash_              atomic.Int32
}

func (this_ *SourceMapGeneratorV3_LineMapper) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapGeneratorV3_LineMapper) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapGeneratorV3_LineMapper) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapGeneratorV3$LineMapper", "LineMapper"
}
func (this_ *SourceMapGeneratorV3_LineMapper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapGeneratorV3_LineMapper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapGeneratorV3.LineMapper", "com.google.debugging.sourcemap.SourceMapGeneratorV3$LineMapper", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapGeneratorV3_LineMapper](v_)
	})
}
func (this_ *SourceMapGeneratorV3_LineMapper) SetOuter_(o *SourceMapGeneratorV3)                   { this_.outer_ = o }
func (this_ *SourceMapGeneratorV3_LineMapper) GetOuter_() *SourceMapGeneratorV3                    { return this_.outer_ }
func (*SourceMapGeneratorV3_LineMapper) Implements_sourcemap_SourceMapGeneratorV3_MappingVisitor() {}

func (this_ *SourceMapGeneratorV3_LineMapper) initFields_SourceMapGeneratorV3_LineMapper() {
	_ = this_
	this_.previousLine = int32(-1)

}

// Java constructor: LineMapper(java.lang.Appendable out, int maxLine)
func NewSourceMapGeneratorV3_LineMapper(outer_ *SourceMapGeneratorV3, out j2gort.Ext_java_lang_Appendable, maxLine int32) *SourceMapGeneratorV3_LineMapper {
	this_ := &SourceMapGeneratorV3_LineMapper{}
	this_.outer_ = outer_
	InitSourceMapGeneratorV3_LineMapper(this_, out, maxLine)
	return this_
}

func InitSourceMapGeneratorV3_LineMapper(this_ *SourceMapGeneratorV3_LineMapper, out j2gort.Ext_java_lang_Appendable, maxLine int32) {
	_ = out
	_ = maxLine
	this_.initFields_SourceMapGeneratorV3_LineMapper()
	this_.out = out
	this_.maxLine = maxLine
}

// Java method: LineMapper.visit(com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m, int line, int col, int nextLine, int nextCol)
func (this_ *SourceMapGeneratorV3_LineMapper) Visit(m *SourceMapGeneratorV3_Mapping, line int32, col int32, nextLine int32, nextCol int32) {
	_ = m
	_ = line
	_ = col
	_ = nextLine
	_ = nextCol
	if this_.previousLine != line {
		this_.previousColumn = int32(0)
	}
	if (line != nextLine) || (col != nextCol) {
		if line < this_.maxLine {
			if this_.previousLine == line {
				this_.out.Append_1_uint16(uint16(44))
			}
			this_.WriteEntry(m, col)
			this_.previousLine = line
			this_.previousColumn = col
		} else {
			if !(m == nil) {
				j2gort.CheckState(false)
			}
		}
	}
	{
		var i int32 = line
		_ = i
		for ; (i <= nextLine) && (i < this_.maxLine); i = i + 1 {
			if i == nextLine {
				break
			}
			this_.priv_LineMapper_closeLine_84(false)
			this_.priv_LineMapper_openLine_83(false)
		}
	}
}

// Java method: LineMapper.writeEntry(com.google.debugging.sourcemap.SourceMapGeneratorV3.Mapping m, int column)
func (this_ *SourceMapGeneratorV3_LineMapper) WriteEntry(m *SourceMapGeneratorV3_Mapping, column int32) {
	_ = m
	_ = column
	Base64VLQ_Encode(this_.out, (column - this_.previousColumn))
	this_.previousColumn = column
	if m != nil {
		var sourceId int32 = this_.outer_.priv_SourceMapGeneratorV3_getSourceId_70(m.sourceFile)
		_ = sourceId
		Base64VLQ_Encode(this_.out, (sourceId - this_.previousSourceFileId))
		this_.previousSourceFileId = sourceId
		var srcline int32 = m.originalLine
		_ = srcline
		var srcColumn int32 = m.originalColumn
		_ = srcColumn
		Base64VLQ_Encode(this_.out, (srcline - this_.previousSourceLine))
		this_.previousSourceLine = srcline
		Base64VLQ_Encode(this_.out, (srcColumn - this_.previousSourceColumn))
		this_.previousSourceColumn = srcColumn
		if m.originalName != j2gort.NullString {
			var nameId int32 = this_.outer_.priv_SourceMapGeneratorV3_getNameId_71(m.originalName)
			_ = nameId
			Base64VLQ_Encode(this_.out, (nameId - this_.previousNameId))
			this_.previousNameId = nameId
		}
	}
}

// Java method: LineMapper.appendLineMappings()
func (this_ *SourceMapGeneratorV3_LineMapper) AppendLineMappings() {
	this_.priv_LineMapper_openLine_83(true)
	{
		var val_29 SourceMapGeneratorV3_MappingTraversal
		val_29.outer_ = this_.outer_
		InitSourceMapGeneratorV3_MappingTraversal(&val_29)
		(&val_29).Traverse(this_)
	}
	this_.priv_LineMapper_closeLine_84(true)
}

// Java method: LineMapper.openLine(boolean firstEntry)
func (this_ *SourceMapGeneratorV3_LineMapper) priv_LineMapper_openLine_83(firstEntry bool) {
	_ = firstEntry
	if firstEntry {
		this_.out.Append_1_uint16(uint16(34))
	}
}

// Java method: LineMapper.closeLine(boolean finalEntry)
func (this_ *SourceMapGeneratorV3_LineMapper) priv_LineMapper_closeLine_84(finalEntry bool) {
	_ = finalEntry
	this_.out.Append_1_uint16(uint16(59))
	if finalEntry {
		this_.out.Append_1_uint16(uint16(34))
	}
}

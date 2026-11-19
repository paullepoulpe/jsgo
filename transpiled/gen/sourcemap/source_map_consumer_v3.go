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
// Port of com/google/debugging/sourcemap/SourceMapConsumerV3.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_source_map_consumer_v3_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_source_map_consumer_v3_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_20 any) any {
		k := j2gort.CastString(lp_20)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashMap()
	})
}

// Java class: com.google.debugging.sourcemap.SourceMapConsumerV3
type SourceMapConsumerV3 struct {
	sources              []string                                             `gson:"sources"`
	sourcesContent       []string                                             `gson:"sourcesContent"`
	names                []string                                             `gson:"names"`
	file                 string                                               `gson:"file"`
	mappings             *SourceMapConsumerV3_Mappings                        `gson:"mappings"`
	reverseSourceMapping *j2gort.OrderedMap[string, j2gort.Ext_java_util_Map] `gson:"reverseSourceMapping"`
	sourceRoot           string                                               `gson:"sourceRoot"`
	extensions           j2gort.Ext_java_util_Map                             `gson:"extensions"`
	lineCount            int32                                                `gson:"lineCount"`
	idHash_              atomic.Int32
}

func (this_ *SourceMapConsumerV3) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerV3) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerV3) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerV3", "SourceMapConsumerV3"
}
func (this_ *SourceMapConsumerV3) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerV3) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3", "com.google.debugging.sourcemap.SourceMapConsumerV3", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerV3](v_)
	})
}
func (*SourceMapConsumerV3) Implements_sourcemap_SourceMapConsumer()       {}
func (*SourceMapConsumerV3) Implements_sourcemap_SourceMappingReversable() {}
func (*SourceMapConsumerV3) Implements_sourcemap_SourceMapping()           {}

func StaticGet_SourceMapConsumerV3_UNMAPPED() int32 { return int32(-1) }
func (this_ *SourceMapConsumerV3) initFields_SourceMapConsumerV3() {
	_ = this_
	this_.file = j2gort.NullString

	this_.sourceRoot = j2gort.NullString
	this_.extensions = j2gort.NewExt_java_util_LinkedHashMap()
}

// Java constructor: SourceMapConsumerV3()
func NewSourceMapConsumerV3() *SourceMapConsumerV3 {
	this_ := &SourceMapConsumerV3{}
	InitSourceMapConsumerV3(this_)
	return this_
}

func InitSourceMapConsumerV3(this_ *SourceMapConsumerV3) {
	this_.initFields_SourceMapConsumerV3()
}

// Java method: SourceMapConsumerV3.parse(java.lang.String contents)
func (this_ *SourceMapConsumerV3) Parse_1(contents string) {
	_ = contents
	var sourceMapObject *SourceMapObject = SourceMapObjectParser_Parse(contents)
	_ = sourceMapObject
	this_.Parse_2(sourceMapObject, nil)
}

// Java method: SourceMapConsumerV3.parse(com.google.debugging.sourcemap.SourceMapObject sourceMapObject, com.google.debugging.sourcemap.@org.jspecify.annotations.Nullable SourceMapSupplier sectionSupplier)
func (this_ *SourceMapConsumerV3) Parse_2(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	_ = sourceMapObject
	_ = sectionSupplier
	if sourceMapObject.GetVersion() != int32(3) {
		panic(NewSourceMapParseException(("Unknown version: " + j2gort.StrValInt32(sourceMapObject.GetVersion()))))
	}
	this_.file = sourceMapObject.GetFile()
	if (this_.file != j2gort.NullString) && (j2gort.CheckNotNullString(this_.file) == "") {
		panic(NewSourceMapParseException("File entry is empty"))
	}
	if sourceMapObject.GetSections() != nil {
		this_.priv_SourceMapConsumerV3_parseMetaMap_8(sourceMapObject, sectionSupplier)
		return
	}
	this_.lineCount = sourceMapObject.GetLineCount()
	this_.sourceRoot = sourceMapObject.GetSourceRoot()
	this_.sources = sourceMapObject.GetSources()
	this_.sourcesContent = sourceMapObject.GetSourcesContent()
	this_.names = sourceMapObject.GetNames()
	var useCompactMappings bool = ((int32(len(this_.sources)) < int32(65535)) && ((this_.names == nil) || (int32(len(this_.names)) < int32(65535))))
	_ = useCompactMappings
	var estimatedEntries int32
	_ = estimatedEntries
	if this_.lineCount >= int32(0) {
		estimatedEntries = (this_.lineCount * int32(2))
	} else {
		estimatedEntries = int32(1000)
	}
	if sourceMapObject.GetMappings() != j2gort.NullString {
		estimatedEntries = max(estimatedEntries, (j2gort.JStrLen(sourceMapObject.GetMappings()) / int32(6)))
	}
	this_.extensions.PutAll(sourceMapObject.GetExtensions())
	this_.mappings = NewSourceMapConsumerV3_MappingBuilder(this_, sourceMapObject.GetMappings(), this_.lineCount, useCompactMappings, estimatedEntries).Build()
}

// Java method: SourceMapConsumerV3.parseMetaMap(com.google.debugging.sourcemap.SourceMapObject sourceMapObject, com.google.debugging.sourcemap.SourceMapSupplier sectionSupplier)
func (this_ *SourceMapConsumerV3) priv_SourceMapConsumerV3_parseMetaMap_8(sourceMapObject *SourceMapObject, sectionSupplier SourceMapSupplier) {
	_ = sourceMapObject
	_ = sectionSupplier
	if sectionSupplier == nil || sectionSupplier.IsJavaNil_() {
		sectionSupplier = NewSourceMapConsumerV3_DefaultSourceMapSupplier()
	}
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					ex := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = ex
					panic(NewSourceMapParseException(("IO exception: " + j2gort.StrVal(ex))))
				} else {
					panic(rec_1)
				}
			}
		}()
		if (((sourceMapObject.GetLineCount() >= int32(0)) || (sourceMapObject.GetMappings() != j2gort.NullString)) || ((sourceMapObject.GetSources() != nil) && (int32(len(sourceMapObject.GetSources())) > int32(0)))) || (sourceMapObject.GetNames() != nil) {
			panic(NewSourceMapParseException("Invalid map format"))
		}
		var val_3 SourceMapGeneratorV3
		InitSourceMapGeneratorV3(&val_3)
		var generator *SourceMapGeneratorV3 = &val_3
		_ = generator
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(sourceMapObject.GetSections()))
			for cur_4.HasNext() {
				rawElem_5 := cur_4.Next()
				section := j2gort.CastPtr[SourceMapSection](rawElem_5)
				_ = section
				var mapSectionContents string = section.GetSectionValue()
				_ = mapSectionContents
				if section.GetSectionType() == StaticGet_SourceMapSection_SectionType_URL() {
					mapSectionContents = sectionSupplier.GetSourceMap(section.GetSectionValue())
				}
				if mapSectionContents == j2gort.NullString {
					panic(NewSourceMapParseException(("Unable to retrieve: " + j2gort.StrValString(section.GetSectionValue()))))
				}
				generator.MergeMapSection_3(section.GetLine(), section.GetColumn(), mapSectionContents)
			}
		}
		var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = sb
		generator.AppendTo(sb, sourceMapObject.GetFile())
		this_.Parse_1(sb.String())
	}()
}

// Java method: SourceMapConsumerV3.getMappingForLine(int lineNumber, int column)
func (this_ *SourceMapConsumerV3) GetMappingForLine(lineNumber int32, column int32) *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping {
	_ = lineNumber
	_ = column
	lineNumber = lineNumber - 1
	column = column - 1
	if (lineNumber < int32(0)) || (lineNumber >= this_.mappings.GetParsedLineCount()) {
		return nil
	}
	var start int32 = this_.mappings.GetLineStart(lineNumber)
	_ = start
	var end int32 = this_.mappings.GetLineStart((lineNumber + int32(1)))
	_ = end
	if start == end {
		return this_.priv_SourceMapConsumerV3_getPreviousMapping_19(lineNumber)
	}
	if this_.mappings.GetGeneratedColumn(start) > column {
		return this_.priv_SourceMapConsumerV3_getPreviousMapping_19(lineNumber)
	}
	var index int32 = this_.priv_SourceMapConsumerV3_search_18(column, start, (end - this_.mappings.GetEntrySize()))
	_ = index
	if !(index >= int32(0)) {
		j2gort.CheckStateFmt(false, "unexpected:%s", j2gort.ValueOfInt(index))
	}
	return this_.priv_SourceMapConsumerV3_getOriginalMappingForEntry_20(index, j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision_F_EXACT)
}

// Java method: SourceMapConsumerV3.getOriginalSources()
func (this_ *SourceMapConsumerV3) GetOriginalSources() j2gort.Ext_java_util_Collection {
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_6 := this_.sources
		if arr_6 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_7 := make([]any, len(arr_6))
		for i_, el_8 := range arr_6 {
			out_7[i_] = j2gort.BoxString(el_8)
		}
		iifeArg_1 = out_7
		break iifeSwitch_1
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(iifeArg_1))
}

// Java method: SourceMapConsumerV3.getOriginalSourcesContent()
func (this_ *SourceMapConsumerV3) GetOriginalSourcesContent() j2gort.Ext_java_util_Collection {
	if this_.sourcesContent == nil {
		return nil
	} else {
		var iifeArg_2 []any
	iifeSwitch_2:
		switch {
		default:

			arr_9 := this_.sourcesContent
			if arr_9 == nil {
				iifeArg_2 = nil
				break iifeSwitch_2
			}
			out_10 := make([]any, len(arr_9))
			for i_, el_11 := range arr_9 {
				out_10[i_] = j2gort.BoxString(el_11)
			}
			iifeArg_2 = out_10
			break iifeSwitch_2
		}
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(iifeArg_2))
	}
}

// Java method: SourceMapConsumerV3.getOriginalNames()
func (this_ *SourceMapConsumerV3) GetOriginalNames() j2gort.Ext_java_util_List {
	var iifeArg_3 []any
iifeSwitch_3:
	switch {
	default:

		arr_12 := this_.names
		if arr_12 == nil {
			iifeArg_3 = nil
			break iifeSwitch_3
		}
		out_13 := make([]any, len(arr_12))
		for i_, el_14 := range arr_12 {
			out_13[i_] = j2gort.BoxString(el_14)
		}
		iifeArg_3 = out_13
		break iifeSwitch_3
	}
	return j2gort.Ext_java_util_Arrays_AsList(iifeArg_3)
}

// Java method: SourceMapConsumerV3.getFile()
func (this_ *SourceMapConsumerV3) GetFile() string {
	return this_.file
}

// Java method: SourceMapConsumerV3.getLineCount()
func (this_ *SourceMapConsumerV3) GetLineCount() int32 {
	return this_.lineCount
}

// Java method: SourceMapConsumerV3.getReverseMapping(java.lang.String originalFile, int line, int column)
func (this_ *SourceMapConsumerV3) GetReverseMapping(originalFile string, line int32, column int32) j2gort.Ext_java_util_Collection {
	_ = originalFile
	_ = line
	_ = column
	if this_.reverseSourceMapping == nil {
		this_.priv_SourceMapConsumerV3_createReverseMapping_21()
	}
	var sourceLineToCollectionMap j2gort.Ext_java_util_Map = this_.reverseSourceMapping.Get(originalFile)
	_ = sourceLineToCollectionMap
	if sourceLineToCollectionMap == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Collections_EmptyList())
	} else {
		var mappings j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](sourceLineToCollectionMap.Get(j2gort.ValueOfInt(line)))
		_ = mappings
		if mappings == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Collections_EmptyList())
		} else {
			return mappings
		}
	}
}

// Java method: SourceMapConsumerV3.getSourceRoot()
func (this_ *SourceMapConsumerV3) GetSourceRoot() string {
	return this_.sourceRoot
}

// Java method: SourceMapConsumerV3.getExtensions()
func (this_ *SourceMapConsumerV3) GetExtensions() j2gort.Ext_java_util_Map {
	return this_.extensions
}

// Java method: SourceMapConsumerV3.search(int target, int start, int end)
func (this_ *SourceMapConsumerV3) priv_SourceMapConsumerV3_search_18(target int32, start int32, end int32) int32 {
	_ = target
	_ = start
	_ = end
	var entrySize int32 = this_.mappings.GetEntrySize()
	_ = entrySize
	for {
		var mid int32 = ((((end - start) / (entrySize * int32(2))) * entrySize) + start)
		_ = mid
		var compare int32 = (this_.mappings.GetGeneratedColumn(mid) - target)
		_ = compare
		if compare == int32(0) {
			return mid
		} else if compare < int32(0) {
			start = (mid + entrySize)
			if start > end {
				return end
			}
		} else {
			end = (mid - entrySize)
			if end < start {
				return end
			}
		}
	}
	panic("unreachable")
}

// Java method: SourceMapConsumerV3.getPreviousMapping(int lineNumber)
func (this_ *SourceMapConsumerV3) priv_SourceMapConsumerV3_getPreviousMapping_19(lineNumber int32) *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping {
	_ = lineNumber
	for doFirst_15 := true; doFirst_15 || (this_.mappings.GetLineStart(lineNumber) == this_.mappings.GetLineStart((lineNumber + int32(1)))); doFirst_15 = false {
		if lineNumber == int32(0) {
			return nil
		}
		lineNumber = lineNumber - 1
	}
	var index int32 = (this_.mappings.GetLineStart((lineNumber + int32(1))) - this_.mappings.GetEntrySize())
	_ = index
	return this_.priv_SourceMapConsumerV3_getOriginalMappingForEntry_20(index, j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision_F_APPROXIMATE_LINE)
}

// Java method: SourceMapConsumerV3.getOriginalMappingForEntry(int index, com.google.debugging.sourcemap.proto.Mapping.OriginalMapping.Precision precision)
func (this_ *SourceMapConsumerV3) priv_SourceMapConsumerV3_getOriginalMappingForEntry_20(index int32, precision *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision) *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping {
	_ = index
	_ = precision
	var sourceFileId int32 = this_.mappings.GetSourceFileId(index)
	_ = sourceFileId
	if sourceFileId == int32(-1) {
		return nil
	} else {
		var x *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder = j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_NewBuilder_0().SetOriginalFile(this_.sources[sourceFileId]).SetLineNumber((this_.mappings.GetSourceLine(index) + int32(1))).SetColumnPosition((this_.mappings.GetSourceColumn(index) + int32(1))).SetPrecision(precision)
		_ = x
		var nameId int32 = this_.mappings.GetNameId(index)
		_ = nameId
		if nameId != int32(-1) {
			x.SetIdentifier(this_.names[nameId])
		}
		return j2gort.CastPtr[j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping](x.Build())
	}
}

// Java method: SourceMapConsumerV3.createReverseMapping()
func (this_ *SourceMapConsumerV3) priv_SourceMapConsumerV3_createReverseMapping_21() {
	this_.reverseSourceMapping = &j2gort.OrderedMap[string, j2gort.Ext_java_util_Map]{}
	{
		var targetLine int32 = int32(0)
		_ = targetLine
		for ; targetLine < this_.mappings.GetParsedLineCount(); targetLine = targetLine + 1 {
			var start int32 = this_.mappings.GetLineStart(targetLine)
			_ = start
			var end int32 = this_.mappings.GetLineStart((targetLine + int32(1)))
			_ = end
			{
				var i int32 = start
				_ = i
				for ; i < end; i = (i + this_.mappings.GetEntrySize()) {
					var sourceFileId int32 = this_.mappings.GetSourceFileId(i)
					_ = sourceFileId
					var sourceLine int32 = this_.mappings.GetSourceLine(i)
					_ = sourceLine
					if (sourceFileId != int32(-1)) && (sourceLine != int32(-1)) {
						var originalFile string = this_.sources[sourceFileId]
						_ = originalFile
						var lineToCollectionMap j2gort.Ext_java_util_Map = j2gort.OrderedMapComputeIfAbsentFn(this_.reverseSourceMapping, originalFile, nil, lambdaSingleton_source_map_consumer_v3_1)
						_ = lineToCollectionMap
						if !(lineToCollectionMap.ContainsKey(j2gort.ValueOfInt(sourceLine))) {
							lineToCollectionMap.Put_2(j2gort.ValueOfInt(sourceLine), j2gort.NewExt_java_util_ArrayList_1_int32(int32(1)))
						}
						var mappingsLine j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](lineToCollectionMap.Get(j2gort.ValueOfInt(sourceLine)))
						_ = mappingsLine
						var builder *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder = j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_NewBuilder_0().SetLineNumber(targetLine).SetColumnPosition(this_.mappings.GetGeneratedColumn(i))
						_ = builder
						mappingsLine.Add_1(j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping](builder.Build())))
					}
				}
			}
		}
	}
}

// Java method: SourceMapConsumerV3.visitMappings(com.google.debugging.sourcemap.SourceMapConsumerV3.EntryVisitor visitor)
func (this_ *SourceMapConsumerV3) VisitMappings(visitor SourceMapConsumerV3_EntryVisitor) {
	_ = visitor
	var pending bool = false
	_ = pending
	var sourceName string = j2gort.NullString
	_ = sourceName
	var symbolName string = j2gort.NullString
	_ = symbolName
	var sourceStartPosition *FilePosition = nil
	_ = sourceStartPosition
	var startPosition *FilePosition = nil
	_ = startPosition
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.mappings.GetParsedLineCount(); i = i + 1 {
			var start int32 = this_.mappings.GetLineStart(i)
			_ = start
			var end int32 = this_.mappings.GetLineStart((i + int32(1)))
			_ = end
			if start != end {
				{
					var j int32 = start
					_ = j
					for ; j < end; j = (j + this_.mappings.GetEntrySize()) {
						if pending {
							var endPosition *FilePosition = NewFilePosition(i, this_.mappings.GetGeneratedColumn(j))
							_ = endPosition
							visitor.Visit(sourceName, symbolName, sourceStartPosition, startPosition, endPosition)
							pending = false
						}
						var sourceFileId int32 = this_.mappings.GetSourceFileId(j)
						_ = sourceFileId
						if sourceFileId != int32(-1) {
							pending = true
							sourceName = this_.sources[sourceFileId]
							var nameId int32 = this_.mappings.GetNameId(j)
							_ = nameId
							if nameId != int32(-1) {
								symbolName = this_.names[nameId]
							} else {
								symbolName = j2gort.NullString
							}
							sourceStartPosition = NewFilePosition(this_.mappings.GetSourceLine(j), this_.mappings.GetSourceColumn(j))
							startPosition = NewFilePosition(i, this_.mappings.GetGeneratedColumn(j))
						}
					}
				}
			}
		}
	}
	if pending {
		var val_21 FilePosition
		InitFilePosition(&val_21, startPosition.GetLine(), (startPosition.GetColumn() + int32(1)))
		var endPosition_2 *FilePosition = &val_21
		_ = endPosition_2
		visitor.Visit(sourceName, symbolName, sourceStartPosition, startPosition, endPosition_2)
	}
}

// Java class: com.google.debugging.sourcemap.SourceMapConsumerV3.DefaultSourceMapSupplier
type SourceMapConsumerV3_DefaultSourceMapSupplier struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerV3_DefaultSourceMapSupplier) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerV3$DefaultSourceMapSupplier", "DefaultSourceMapSupplier"
}
func (this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerV3_DefaultSourceMapSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3.DefaultSourceMapSupplier", "com.google.debugging.sourcemap.SourceMapConsumerV3$DefaultSourceMapSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerV3_DefaultSourceMapSupplier](v_)
	})
}
func (*SourceMapConsumerV3_DefaultSourceMapSupplier) Implements_sourcemap_SourceMapSupplier() {}

func (this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) initFields_SourceMapConsumerV3_DefaultSourceMapSupplier() {
	_ = this_
}

// Java constructor: DefaultSourceMapSupplier()
func NewSourceMapConsumerV3_DefaultSourceMapSupplier() *SourceMapConsumerV3_DefaultSourceMapSupplier {
	this_ := &SourceMapConsumerV3_DefaultSourceMapSupplier{}
	InitSourceMapConsumerV3_DefaultSourceMapSupplier(this_)
	return this_
}

func InitSourceMapConsumerV3_DefaultSourceMapSupplier(this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) {
	this_.initFields_SourceMapConsumerV3_DefaultSourceMapSupplier()
}

// Java method: DefaultSourceMapSupplier.getSourceMap(java.lang.String url)
func (this_ *SourceMapConsumerV3_DefaultSourceMapSupplier) GetSourceMap(url string) string {
	_ = url
	return j2gort.NullString
}

// Java class: com.google.debugging.sourcemap.SourceMapConsumerV3.MappingBuilder
type SourceMapConsumerV3_MappingBuilder struct {
	outer_            *SourceMapConsumerV3
	content           *SourceMapConsumerV3_StringCharIterator `gson:"content"`
	flatEntries       []int32                                 `gson:"flatEntries"`
	lineStart         []int32                                 `gson:"lineStart"`
	lineCount         int32                                   `gson:"lineCount"`
	entrySize         int32                                   `gson:"entrySize"`
	line              int32                                   `gson:"line"`
	previousCol       int32                                   `gson:"previousCol"`
	previousSrcId     int32                                   `gson:"previousSrcId"`
	previousSrcLine   int32                                   `gson:"previousSrcLine"`
	previousSrcColumn int32                                   `gson:"previousSrcColumn"`
	previousNameId    int32                                   `gson:"previousNameId"`
	idHash_           atomic.Int32
}

func (this_ *SourceMapConsumerV3_MappingBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerV3_MappingBuilder) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerV3_MappingBuilder) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerV3$MappingBuilder", "MappingBuilder"
}
func (this_ *SourceMapConsumerV3_MappingBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerV3_MappingBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3.MappingBuilder", "com.google.debugging.sourcemap.SourceMapConsumerV3$MappingBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerV3_MappingBuilder](v_)
	})
}
func (this_ *SourceMapConsumerV3_MappingBuilder) SetOuter_(o *SourceMapConsumerV3) { this_.outer_ = o }
func (this_ *SourceMapConsumerV3_MappingBuilder) GetOuter_() *SourceMapConsumerV3 {
	return this_.outer_
}

func StaticGet_SourceMapConsumerV3_MappingBuilder_MAX_ENTRY_VALUES() int32 { return int32(5) }
func (this_ *SourceMapConsumerV3_MappingBuilder) initFields_SourceMapConsumerV3_MappingBuilder() {
	_ = this_

}

// Java constructor: MappingBuilder(java.lang.String lineMap, int lineCount, boolean useCompactMappings, int estimatedEntries)
func NewSourceMapConsumerV3_MappingBuilder(outer_ *SourceMapConsumerV3, lineMap string, lineCount int32, useCompactMappings bool, estimatedEntries int32) *SourceMapConsumerV3_MappingBuilder {
	this_ := &SourceMapConsumerV3_MappingBuilder{}
	this_.outer_ = outer_
	InitSourceMapConsumerV3_MappingBuilder(this_, lineMap, lineCount, useCompactMappings, estimatedEntries)
	return this_
}

func InitSourceMapConsumerV3_MappingBuilder(this_ *SourceMapConsumerV3_MappingBuilder, lineMap string, lineCount int32, useCompactMappings bool, estimatedEntries int32) {
	_ = lineMap
	_ = lineCount
	_ = useCompactMappings
	_ = estimatedEntries
	this_.initFields_SourceMapConsumerV3_MappingBuilder()
	this_.content = NewSourceMapConsumerV3_StringCharIterator(lineMap)
	this_.lineCount = lineCount
	if useCompactMappings {
		this_.entrySize = int32(4)
	} else {
		this_.entrySize = int32(5)
	}
	this_.flatEntries = make([]int32, (estimatedEntries * this_.entrySize))
	var estimatedLines int32
	_ = estimatedLines
	if lineCount >= int32(0) {
		estimatedLines = lineCount
	} else {
		estimatedLines = int32(1000)
	}
	this_.lineStart = make([]int32, (estimatedLines + int32(1)))
}

// Java method: MappingBuilder.build()
func (this_ *SourceMapConsumerV3_MappingBuilder) Build() *SourceMapConsumerV3_Mappings {
	var temp []int32 = make([]int32, int32(5))
	_ = temp
	var entriesCount int32 = int32(0)
	_ = entriesCount
	var linesCount int32 = int32(0)
	_ = linesCount
	this_.priv_MappingBuilder_setLineStart_26(int32(0), int32(0))
	for this_.content.HasNext() {
		if this_.priv_MappingBuilder_tryConsumeToken_37(uint16(59)) {
			linesCount = linesCount + 1
			this_.priv_MappingBuilder_setLineStart_26(linesCount, entriesCount)
			this_.line = this_.line + 1
			this_.previousCol = int32(0)
		} else {
			var entryValues int32 = int32(0)
			_ = entryValues
			for !(this_.priv_MappingBuilder_entryComplete_38()) {
				{
					a_22 := temp
					i_23 := entryValues
					a_22[i_23] = this_.priv_MappingBuilder_nextValue_39()
				}
				entryValues = entryValues + 1
			}
			this_.priv_MappingBuilder_setEntry_27(entriesCount, temp, entryValues, this_.previousCol, this_.previousSrcId, this_.previousSrcLine, this_.previousSrcColumn, this_.previousNameId)
			this_.priv_MappingBuilder_validateEntry_35(entriesCount, this_.line)
			this_.previousCol = this_.priv_MappingBuilder_getGeneratedColumn_32(entriesCount)
			if entryValues >= int32(4) {
				this_.previousSrcId = this_.priv_MappingBuilder_getSourceFileId_30(entriesCount)
				this_.previousSrcLine = this_.priv_MappingBuilder_getSourceLine_33(entriesCount)
				this_.previousSrcColumn = this_.priv_MappingBuilder_getSourceColumn_34(entriesCount)
			}
			if entryValues == int32(5) {
				this_.previousNameId = this_.priv_MappingBuilder_getNameId_31(entriesCount)
			}
			entriesCount = (entriesCount + this_.entrySize)
			this_.priv_MappingBuilder_tryConsumeToken_37(uint16(44))
		}
	}
	if entriesCount > this_.priv_MappingBuilder_getLineStart_25(linesCount) {
		linesCount = linesCount + 1
		this_.priv_MappingBuilder_setLineStart_26(linesCount, entriesCount)
	}
	this_.priv_MappingBuilder_trim_36(entriesCount, linesCount)
	return NewSourceMapConsumerV3_Mappings(this_.flatEntries, this_.lineStart, this_.lineCount, this_.entrySize)
}

// Java method: MappingBuilder.getLineStart(int line)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getLineStart_25(line int32) int32 {
	_ = line
	return this_.lineStart[line]
}

// Java method: MappingBuilder.setLineStart(int line, int entriesCount)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_setLineStart_26(line int32, entriesCount int32) {
	_ = line
	_ = entriesCount
	if line >= int32(len(this_.lineStart)) {
		this_.lineStart = func() []int32 {
			arrSrc_24 := this_.lineStart
			arrLen_25 := (int32(len(this_.lineStart)) * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_24)
			if arrLen_25 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_26 := make([]int32, arrLen_25)
			copy(arrOut_26, arrSrc_24)
			return arrOut_26
		}()
	}
	{
		a_27 := this_.lineStart
		i_28 := line
		a_27[i_28] = entriesCount
	}
}

// Java method: MappingBuilder.setEntry(int entriesCount, int[] vals, int entryValues, int previousCol, int previousSrcId, int previousSrcLine, int previousSrcColumn, int previousNameId)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_setEntry_27(entriesCount int32, vals []int32, entryValues int32, previousCol int32, previousSrcId int32, previousSrcLine int32, previousSrcColumn int32, previousNameId int32) {
	_ = entriesCount
	_ = vals
	_ = entryValues
	_ = previousCol
	_ = previousSrcId
	_ = previousSrcLine
	_ = previousSrcColumn
	_ = previousNameId
	if (entriesCount + this_.entrySize) > int32(len(this_.flatEntries)) {
		this_.flatEntries = func() []int32 {
			arrSrc_29 := this_.flatEntries
			arrLen_30 := (int32(len(this_.flatEntries)) * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_29)
			if arrLen_30 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_31 := make([]int32, arrLen_30)
			copy(arrOut_31, arrSrc_29)
			return arrOut_31
		}()
	}
	if this_.entrySize == int32(4) {
		this_.priv_MappingBuilder_setEntryCompact_28(entriesCount, vals, entryValues, previousCol, previousSrcId, previousSrcLine, previousSrcColumn, previousNameId)
	} else {
		this_.priv_MappingBuilder_setEntryNormal_29(entriesCount, vals, entryValues, previousCol, previousSrcId, previousSrcLine, previousSrcColumn, previousNameId)
	}
}

// Java method: MappingBuilder.setEntryCompact(int entriesCount, int[] vals, int entryValues, int previousCol, int previousSrcId, int previousSrcLine, int previousSrcColumn, int previousNameId)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_setEntryCompact_28(entriesCount int32, vals []int32, entryValues int32, previousCol int32, previousSrcId int32, previousSrcLine int32, previousSrcColumn int32, previousNameId int32) {
	_ = entriesCount
	_ = vals
	_ = entryValues
	_ = previousCol
	_ = previousSrcId
	_ = previousSrcLine
	_ = previousSrcColumn
	_ = previousNameId
	switch entryValues {
	case int32(1):
		{
			a_32 := this_.flatEntries
			i_33 := entriesCount
			a_32[i_33] = (vals[int32(0)] + previousCol)
		}
		{
			a_34 := this_.flatEntries
			i_35 := (entriesCount + int32(1))
			a_34[i_35] = int32(-1)
		}
		{
			a_36 := this_.flatEntries
			i_37 := (entriesCount + int32(2))
			a_36[i_37] = int32(-1)
		}
		{
			a_38 := this_.flatEntries
			i_39 := (entriesCount + int32(3))
			a_38[i_39] = int32(-1)
		}
	case int32(4):
		{
			a_40 := this_.flatEntries
			i_41 := entriesCount
			a_40[i_41] = (vals[int32(0)] + previousCol)
		}
		{
			a_42 := this_.flatEntries
			i_43 := (entriesCount + int32(1))
			a_42[i_43] = (vals[int32(2)] + previousSrcLine)
		}
		{
			a_44 := this_.flatEntries
			i_45 := (entriesCount + int32(2))
			a_44[i_45] = (vals[int32(3)] + previousSrcColumn)
		}
		var srcId int32 = (vals[int32(1)] + previousSrcId)
		_ = srcId
		{
			a_46 := this_.flatEntries
			i_47 := (entriesCount + int32(3))
			a_46[i_47] = ((javalang.Shl32(srcId, int(int32(16)))) | int32(65535))
		}
	case int32(5):
		{
			a_48 := this_.flatEntries
			i_49 := entriesCount
			a_48[i_49] = (vals[int32(0)] + previousCol)
		}
		{
			a_50 := this_.flatEntries
			i_51 := (entriesCount + int32(1))
			a_50[i_51] = (vals[int32(2)] + previousSrcLine)
		}
		{
			a_52 := this_.flatEntries
			i_53 := (entriesCount + int32(2))
			a_52[i_53] = (vals[int32(3)] + previousSrcColumn)
		}
		var srcId_2 int32 = (vals[int32(1)] + previousSrcId)
		_ = srcId_2
		var nameId int32 = (vals[int32(4)] + previousNameId)
		_ = nameId
		{
			a_54 := this_.flatEntries
			i_55 := (entriesCount + int32(3))
			a_54[i_55] = ((javalang.Shl32(srcId_2, int(int32(16)))) | (nameId & int32(65535)))
		}
	default:
		panic(NewSourceMapParseException(("Unexpected number of values for entry:" + j2gort.StrValInt32(entryValues))))
	}
}

// Java method: MappingBuilder.setEntryNormal(int entriesCount, int[] vals, int entryValues, int previousCol, int previousSrcId, int previousSrcLine, int previousSrcColumn, int previousNameId)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_setEntryNormal_29(entriesCount int32, vals []int32, entryValues int32, previousCol int32, previousSrcId int32, previousSrcLine int32, previousSrcColumn int32, previousNameId int32) {
	_ = entriesCount
	_ = vals
	_ = entryValues
	_ = previousCol
	_ = previousSrcId
	_ = previousSrcLine
	_ = previousSrcColumn
	_ = previousNameId
	switch entryValues {
	case int32(1):
		{
			a_56 := this_.flatEntries
			i_57 := entriesCount
			a_56[i_57] = (vals[int32(0)] + previousCol)
		}
		{
			a_58 := this_.flatEntries
			i_59 := (entriesCount + int32(1))
			a_58[i_59] = int32(-1)
		}
		{
			a_60 := this_.flatEntries
			i_61 := (entriesCount + int32(2))
			a_60[i_61] = int32(-1)
		}
		{
			a_62 := this_.flatEntries
			i_63 := (entriesCount + int32(3))
			a_62[i_63] = int32(-1)
		}
		{
			a_64 := this_.flatEntries
			i_65 := (entriesCount + int32(4))
			a_64[i_65] = int32(-1)
		}
	case int32(4):
		{
			a_66 := this_.flatEntries
			i_67 := entriesCount
			a_66[i_67] = (vals[int32(0)] + previousCol)
		}
		{
			a_68 := this_.flatEntries
			i_69 := (entriesCount + int32(1))
			a_68[i_69] = (vals[int32(2)] + previousSrcLine)
		}
		{
			a_70 := this_.flatEntries
			i_71 := (entriesCount + int32(2))
			a_70[i_71] = (vals[int32(3)] + previousSrcColumn)
		}
		{
			a_72 := this_.flatEntries
			i_73 := (entriesCount + int32(3))
			a_72[i_73] = (vals[int32(1)] + previousSrcId)
		}
		{
			a_74 := this_.flatEntries
			i_75 := (entriesCount + int32(4))
			a_74[i_75] = int32(-1)
		}
	case int32(5):
		{
			a_76 := this_.flatEntries
			i_77 := entriesCount
			a_76[i_77] = (vals[int32(0)] + previousCol)
		}
		{
			a_78 := this_.flatEntries
			i_79 := (entriesCount + int32(1))
			a_78[i_79] = (vals[int32(2)] + previousSrcLine)
		}
		{
			a_80 := this_.flatEntries
			i_81 := (entriesCount + int32(2))
			a_80[i_81] = (vals[int32(3)] + previousSrcColumn)
		}
		{
			a_82 := this_.flatEntries
			i_83 := (entriesCount + int32(3))
			a_82[i_83] = (vals[int32(1)] + previousSrcId)
		}
		{
			a_84 := this_.flatEntries
			i_85 := (entriesCount + int32(4))
			a_84[i_85] = (vals[int32(4)] + previousNameId)
		}
	default:
		panic(NewSourceMapParseException(("Unexpected number of values for entry:" + j2gort.StrValInt32(entryValues))))
	}
}

// Java method: MappingBuilder.getSourceFileId(int index)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getSourceFileId_30(index int32) int32 {
	_ = index
	if this_.entrySize == int32(4) {
		var val int32 = javalang.Shr32(this_.flatEntries[(index+int32(3))], int(int32(16)))
		_ = val
		if (val == int32(-1)) || ((val & int32(65535)) == int32(65535)) {
			return int32(-1)
		} else {
			return (val & int32(65535))
		}
	}
	return this_.flatEntries[(index + int32(3))]
}

// Java method: MappingBuilder.getNameId(int index)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getNameId_31(index int32) int32 {
	_ = index
	if this_.entrySize == int32(4) {
		var val int32 = (this_.flatEntries[(index+int32(3))] & int32(65535))
		_ = val
		if val == int32(65535) {
			return int32(-1)
		} else {
			return val
		}
	}
	return this_.flatEntries[(index + int32(4))]
}

// Java method: MappingBuilder.getGeneratedColumn(int index)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getGeneratedColumn_32(index int32) int32 {
	_ = index
	return this_.flatEntries[index]
}

// Java method: MappingBuilder.getSourceLine(int index)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getSourceLine_33(index int32) int32 {
	_ = index
	return this_.flatEntries[(index + int32(1))]
}

// Java method: MappingBuilder.getSourceColumn(int index)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_getSourceColumn_34(index int32) int32 {
	_ = index
	return this_.flatEntries[(index + int32(2))]
}

// Java method: MappingBuilder.validateEntry(int entriesCount, int line)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_validateEntry_35(entriesCount int32, line int32) {
	_ = entriesCount
	_ = line
	if !((this_.lineCount < int32(0)) || (line < this_.lineCount)) {
		j2gort.CheckStateFmt(false, "line=%s, lineCount=%s", j2gort.ValueOfInt(line), j2gort.ValueOfInt(this_.lineCount))
	}
	var sourceFileId int32 = this_.priv_MappingBuilder_getSourceFileId_30(entriesCount)
	_ = sourceFileId
	var nameId int32 = this_.priv_MappingBuilder_getNameId_31(entriesCount)
	_ = nameId
	if !((sourceFileId == int32(-1)) || (sourceFileId < int32(len(this_.outer_.sources)))) {
		j2gort.CheckState(false)
	}
	if !((nameId == int32(-1)) || (nameId < int32(len(this_.outer_.names)))) {
		j2gort.CheckState(false)
	}
}

// Java method: MappingBuilder.trim(int entriesCount, int linesCount)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_trim_36(entriesCount int32, linesCount int32) {
	_ = entriesCount
	_ = linesCount
	if int32(len(this_.flatEntries)) > entriesCount {
		this_.flatEntries = func() []int32 {
			arrSrc_86 := this_.flatEntries
			arrLen_87 := entriesCount
			j2gort.CheckNotNullSlice(arrSrc_86)
			if arrLen_87 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_88 := make([]int32, arrLen_87)
			copy(arrOut_88, arrSrc_86)
			return arrOut_88
		}()
	}
	if int32(len(this_.lineStart)) > (linesCount + int32(1)) {
		this_.lineStart = func() []int32 {
			arrSrc_89 := this_.lineStart
			arrLen_90 := (linesCount + int32(1))
			j2gort.CheckNotNullSlice(arrSrc_89)
			if arrLen_90 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_91 := make([]int32, arrLen_90)
			copy(arrOut_91, arrSrc_89)
			return arrOut_91
		}()
	}
	this_.lineCount = linesCount
}

// Java method: MappingBuilder.tryConsumeToken(char token)
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_tryConsumeToken_37(token uint16) bool {
	_ = token
	if this_.content.HasNext() && (int32(this_.content.Peek()) == int32(token)) {
		this_.content.Next()
		return true
	}
	return false
}

// Java method: MappingBuilder.entryComplete()
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_entryComplete_38() bool {
	if !(this_.content.HasNext()) {
		return true
	}
	var c uint16 = this_.content.Peek()
	_ = c
	return ((int32(c) == int32(uint16(59))) || (int32(c) == int32(uint16(44))))
}

// Java method: MappingBuilder.nextValue()
func (this_ *SourceMapConsumerV3_MappingBuilder) priv_MappingBuilder_nextValue_39() int32 {
	return Base64VLQ_Decode(ToIface_Base64VLQ_CharIterator(this_.content))
}

// Java class: com.google.debugging.sourcemap.SourceMapConsumerV3.StringCharIterator
type SourceMapConsumerV3_StringCharIterator struct {
	content string `gson:"content"`
	length  int32  `gson:"length"`
	current int32  `gson:"current"`
	idHash_ atomic.Int32
}

func (this_ *SourceMapConsumerV3_StringCharIterator) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerV3_StringCharIterator) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerV3_StringCharIterator) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerV3$StringCharIterator", "StringCharIterator"
}
func (this_ *SourceMapConsumerV3_StringCharIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerV3_StringCharIterator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3.StringCharIterator", "com.google.debugging.sourcemap.SourceMapConsumerV3$StringCharIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerV3_StringCharIterator](v_)
	})
}
func (*SourceMapConsumerV3_StringCharIterator) Implements_sourcemap_Base64VLQ_CharIterator() {}

func (this_ *SourceMapConsumerV3_StringCharIterator) initFields_SourceMapConsumerV3_StringCharIterator() {
	_ = this_
	this_.content = j2gort.NullString

}

// Java constructor: StringCharIterator(java.lang.String content)
func NewSourceMapConsumerV3_StringCharIterator(content string) *SourceMapConsumerV3_StringCharIterator {
	this_ := &SourceMapConsumerV3_StringCharIterator{}
	InitSourceMapConsumerV3_StringCharIterator(this_, content)
	return this_
}

func InitSourceMapConsumerV3_StringCharIterator(this_ *SourceMapConsumerV3_StringCharIterator, content string) {
	_ = content
	this_.initFields_SourceMapConsumerV3_StringCharIterator()
	this_.content = content
	this_.length = j2gort.JStrLen(content)
}

// Java method: StringCharIterator.next()
func (this_ *SourceMapConsumerV3_StringCharIterator) Next() uint16 {
	return j2gort.JStrCharAt(this_.content, func() int32 { old_92 := this_.current; nw_93 := old_92 + 1; this_.current = nw_93; return old_92 }())
}

// Java method: StringCharIterator.peek()
func (this_ *SourceMapConsumerV3_StringCharIterator) Peek() uint16 {
	return j2gort.JStrCharAt(this_.content, this_.current)
}

// Java method: StringCharIterator.hasNext()
func (this_ *SourceMapConsumerV3_StringCharIterator) HasNext() bool {
	return (this_.current < this_.length)
}

// Java class: com.google.debugging.sourcemap.SourceMapConsumerV3.Mappings
type SourceMapConsumerV3_Mappings struct {
	flatEntries []int32 `gson:"flatEntries"`
	lineStart   []int32 `gson:"lineStart"`
	lineCount   int32   `gson:"lineCount"`
	entrySize   int32   `gson:"entrySize"`
	idHash_     atomic.Int32
}

func (this_ *SourceMapConsumerV3_Mappings) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapConsumerV3_Mappings) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapConsumerV3_Mappings) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapConsumerV3$Mappings", "Mappings"
}
func (this_ *SourceMapConsumerV3_Mappings) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapConsumerV3_Mappings) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3.Mappings", "com.google.debugging.sourcemap.SourceMapConsumerV3$Mappings", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapConsumerV3_Mappings](v_)
	})
}

func (this_ *SourceMapConsumerV3_Mappings) initFields_SourceMapConsumerV3_Mappings() {
	_ = this_
}

// Java constructor: Mappings(int[] flatEntries, int[] lineStart, int lineCount, int entrySize)
func NewSourceMapConsumerV3_Mappings(flatEntries []int32, lineStart []int32, lineCount int32, entrySize int32) *SourceMapConsumerV3_Mappings {
	this_ := &SourceMapConsumerV3_Mappings{}
	InitSourceMapConsumerV3_Mappings(this_, flatEntries, lineStart, lineCount, entrySize)
	return this_
}

func InitSourceMapConsumerV3_Mappings(this_ *SourceMapConsumerV3_Mappings, flatEntries []int32, lineStart []int32, lineCount int32, entrySize int32) {
	_ = flatEntries
	_ = lineStart
	_ = lineCount
	_ = entrySize
	this_.initFields_SourceMapConsumerV3_Mappings()
	this_.flatEntries = flatEntries
	this_.lineStart = lineStart
	this_.lineCount = lineCount
	this_.entrySize = entrySize
}

// Java method: Mappings.getEntrySize()
func (this_ *SourceMapConsumerV3_Mappings) GetEntrySize() int32 {
	return this_.entrySize
}

// Java method: Mappings.getGeneratedColumn(int index)
func (this_ *SourceMapConsumerV3_Mappings) GetGeneratedColumn(index int32) int32 {
	_ = index
	return this_.flatEntries[index]
}

// Java method: Mappings.getSourceFileId(int index)
func (this_ *SourceMapConsumerV3_Mappings) GetSourceFileId(index int32) int32 {
	_ = index
	if this_.entrySize == int32(4) {
		var val int32 = javalang.Shr32(this_.flatEntries[(index+int32(3))], int(int32(16)))
		_ = val
		if (val == int32(-1)) || ((val & int32(65535)) == int32(65535)) {
			return int32(-1)
		} else {
			return (val & int32(65535))
		}
	}
	return this_.flatEntries[(index + int32(3))]
}

// Java method: Mappings.getSourceLine(int index)
func (this_ *SourceMapConsumerV3_Mappings) GetSourceLine(index int32) int32 {
	_ = index
	return this_.flatEntries[(index + int32(1))]
}

// Java method: Mappings.getSourceColumn(int index)
func (this_ *SourceMapConsumerV3_Mappings) GetSourceColumn(index int32) int32 {
	_ = index
	return this_.flatEntries[(index + int32(2))]
}

// Java method: Mappings.getNameId(int index)
func (this_ *SourceMapConsumerV3_Mappings) GetNameId(index int32) int32 {
	_ = index
	if this_.entrySize == int32(4) {
		var val int32 = (this_.flatEntries[(index+int32(3))] & int32(65535))
		_ = val
		if val == int32(65535) {
			return int32(-1)
		} else {
			return val
		}
	}
	return this_.flatEntries[(index + int32(4))]
}

// Java method: Mappings.getParsedLineCount()
func (this_ *SourceMapConsumerV3_Mappings) GetParsedLineCount() int32 {
	return this_.lineCount
}

// Java method: Mappings.getLineStart(int line)
func (this_ *SourceMapConsumerV3_Mappings) GetLineStart(line int32) int32 {
	_ = line
	return this_.lineStart[line]
}

// Java interface: com.google.debugging.sourcemap.SourceMapConsumerV3.EntryVisitor
type SourceMapConsumerV3_EntryVisitor interface {
	Implements_sourcemap_SourceMapConsumerV3_EntryVisitor()
	IsJavaNil_() bool
	Visit(string, string, *FilePosition, *FilePosition, *FilePosition)
}

func ToIface_SourceMapConsumerV3_EntryVisitor[T interface {
	SourceMapConsumerV3_EntryVisitor
	comparable
}](v T) SourceMapConsumerV3_EntryVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceMapConsumerV3_EntryVisitor[T SourceMapConsumerV3_EntryVisitor](v T) SourceMapConsumerV3_EntryVisitor {
	var vi SourceMapConsumerV3_EntryVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapConsumerV3.EntryVisitor", "com.google.debugging.sourcemap.SourceMapConsumerV3$EntryVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapConsumerV3_EntryVisitor](v_)
	})
}

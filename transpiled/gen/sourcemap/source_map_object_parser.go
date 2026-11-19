/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/SourceMapObjectParser.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapObjectParser
type SourceMapObjectParser struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapObjectParser) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapObjectParser) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapObjectParser) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapObjectParser", "SourceMapObjectParser"
}
func (this_ *SourceMapObjectParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapObjectParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapObjectParser", "com.google.debugging.sourcemap.SourceMapObjectParser", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapObjectParser](v_)
	})
}

var static_SourceMapObjectParser_gson *j2gort.Ext_com_google_gson_Gson

func StaticGet_SourceMapObjectParser_Gson() *j2gort.Ext_com_google_gson_Gson {
	clinit_SourceMapObjectParser()
	return static_SourceMapObjectParser_gson
}
func StaticSet_SourceMapObjectParser_Gson(v *j2gort.Ext_com_google_gson_Gson) {
	clinit_SourceMapObjectParser()
	static_SourceMapObjectParser_gson = v
}

var clinitDone_SourceMapObjectParser uint32
var clinitState_SourceMapObjectParser j2gort.ClinitState

func clinit_SourceMapObjectParser() {
	if atomic.LoadUint32(&clinitDone_SourceMapObjectParser) != 0 {
		return
	}
	clinitSlow_SourceMapObjectParser()
}

func clinitSlow_SourceMapObjectParser() {
	if !j2gort.EnterClinit(&clinitState_SourceMapObjectParser, &clinitDone_SourceMapObjectParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceMapObjectParser, &clinitDone_SourceMapObjectParser)
	static_SourceMapObjectParser_gson = j2gort.NewExt_com_google_gson_Gson()
}

func (this_ *SourceMapObjectParser) initFields_SourceMapObjectParser() {
	_ = this_
}

// Java constructor: SourceMapObjectParser()
func NewSourceMapObjectParser() *SourceMapObjectParser {
	clinit_SourceMapObjectParser()
	this_ := &SourceMapObjectParser{}
	InitSourceMapObjectParser(this_)
	return this_
}

func InitSourceMapObjectParser(this_ *SourceMapObjectParser) {
	this_.initFields_SourceMapObjectParser()
}

// Java static method: SourceMapObjectParser.parse(java.lang.String contents)
func SourceMapObjectParser_Parse(contents string) *SourceMapObject {
	clinit_SourceMapObjectParser()
	_ = contents
	var builder *SourceMapObject_Builder = SourceMapObject_S_Builder()
	_ = builder
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[j2gort.Ext_com_google_gson_JsonParseException](jex_2) {
					ex := j2gort.AsIface[j2gort.Ext_com_google_gson_JsonParseException](jex_2)
					_ = ex
					panic(NewSourceMapParseException(("JSON parse exception: " + j2gort.StrVal(ex))))
				} else {
					panic(rec_1)
				}
			}
		}()
		var sourceMapRoot *j2gort.Ext_com_google_gson_JsonObject = j2gort.CastPtr[j2gort.Ext_com_google_gson_JsonObject](StaticGet_SourceMapObjectParser_Gson().FromJson_2_string_Class(contents, j2gort.ClassLiteral("com.google.gson.JsonObject", "JsonObject")))
		_ = sourceMapRoot
		builder.SetVersion(sourceMapRoot.Get("version").GetAsInt())
		builder.SetFile(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "file"))
		var iifeArg_1 int32
	iifeSwitch_1:
		switch {
		default:

			if sourceMapRoot.Has("lineCount") {
				iifeArg_1 = sourceMapRoot.Get("lineCount").GetAsInt()
				break iifeSwitch_1
			}
			iifeArg_1 = int32(-1)
			break iifeSwitch_1
		}
		builder.SetLineCount(iifeArg_1)
		builder.SetMappings(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "mappings"))
		builder.SetSourceRoot(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "sourceRoot"))
		if sourceMapRoot.Has("sections") {
			var listBuilder []*SourceMapSection = []*SourceMapSection{}
			_ = listBuilder
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(sourceMapRoot.Get("sections").GetAsJsonArray_0()))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					each := j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](rawElem_4)
					_ = each
					listBuilder = append(listBuilder, j2gort.CheckNotNullElem(SourceMapObjectParser_BuildSection(each.GetAsJsonObject_0())))
				}
			}
			builder.SetSections(j2gort.ImmListFromTypedSlice(listBuilder))
		}
		builder.SetSources(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("sources")))
		builder.SetSourcesContent(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("sourcesContent")))
		builder.SetNames(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("names")))
		var extensions j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
		_ = extensions
		{
			cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(sourceMapRoot.EntrySet()))
			for cur_5.HasNext() {
				rawElem_6 := cur_5.Next()
				entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](rawElem_6)
				_ = entry
				if javalang.StringStartsWith(j2gort.CheckNotNullString(j2gort.CastString(entry.GetKey())), "x_") {
					j2gort.MapPutStr(extensions, j2gort.CastString(entry.GetKey()), j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](entry.GetValue()))
				}
			}
		}
		builder.SetExtensions(j2gort.Ext_java_util_Collections_UnmodifiableMap(extensions))
	}()
	return builder.Build()
}

// Java static method: SourceMapObjectParser.parseFast(java.lang.String contents)
func SourceMapObjectParser_ParseFast(contents string) *SourceMapObject {
	clinit_SourceMapObjectParser()
	_ = contents
	var builder *SourceMapObject_Builder = SourceMapObject_S_Builder()
	_ = builder
	builder.SetLineCount(int32(-1))
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javalang.IException](jex_8) {
					ex := j2gort.AsIface[javalang.IException](jex_8)
					_ = ex
					panic(NewSourceMapParseException(("Parse exception: " + j2gort.StrVal(ex))))
				} else {
					panic(rec_7)
				}
			}
		}()
		var val_9 SourceMapJsonLexer
		InitSourceMapJsonLexer(&val_9, contents)
		var lexer *SourceMapJsonLexer = &val_9
		_ = lexer
		lexer.BeginObject()
		var extensions j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
		_ = extensions
		for lexer.HasNext() {
			var name string = lexer.NextName()
			_ = name
			switch j2gort.CheckNotNullString(name) {
			case "version":
				builder.SetVersion(lexer.NextInt())
			case "file":
				builder.SetFile(lexer.NextStringOrNull())
			case "lineCount":
				builder.SetLineCount(lexer.NextInt())
			case "mappings":
				builder.SetMappings(lexer.NextStringOrNull())
			case "sourceRoot":
				builder.SetSourceRoot(lexer.NextStringOrNull())
			case "sections":
				var rawSections string = lexer.NextRawValue()
				_ = rawSections
				var sectionsArray *j2gort.Ext_com_google_gson_JsonArray = j2gort.CastPtr[j2gort.Ext_com_google_gson_JsonArray](StaticGet_SourceMapObjectParser_Gson().FromJson_2_string_Class(rawSections, j2gort.ClassLiteral("com.google.gson.JsonArray", "JsonArray")))
				_ = sectionsArray
				var listBuilder []*SourceMapSection = []*SourceMapSection{}
				_ = listBuilder
				{
					cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(sectionsArray))
					for cur_10.HasNext() {
						rawElem_11 := cur_10.Next()
						each := j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](rawElem_11)
						_ = each
						listBuilder = append(listBuilder, j2gort.CheckNotNullElem(SourceMapObjectParser_BuildSection(each.GetAsJsonObject_0())))
					}
				}
				builder.SetSections(j2gort.ImmListFromTypedSlice(listBuilder))
			case "sources":
				builder.SetSources(SourceMapObjectParser_ReadStringArray(lexer))
			case "sourcesContent":
				builder.SetSourcesContent(SourceMapObjectParser_ReadStringArray(lexer))
			case "names":
				builder.SetNames(SourceMapObjectParser_ReadStringArray(lexer))
			default:
				if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "x_") {
					j2gort.MapPutStr(extensions, name, j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](StaticGet_SourceMapObjectParser_Gson().FromJson_2_string_Class(lexer.NextRawValue(), j2gort.ClassLiteral("com.google.gson.JsonElement", "JsonElement"))))
				} else {
					lexer.SkipValue()
				}
			}
			lexer.CheckComma()
		}
		lexer.EndObject()
		lexer.SkipWhitespace()
		if lexer.pos < lexer.length {
			panic(NewSourceMapParseException("Unexpected trailing characters"))
		}
		builder.SetExtensions(j2gort.Ext_java_util_Collections_UnmodifiableMap(extensions))
	}()
	return builder.Build()
}

// Java static method: SourceMapObjectParser.buildSection(com.google.gson.JsonObject section)
func SourceMapObjectParser_BuildSection(section *j2gort.Ext_com_google_gson_JsonObject) *SourceMapSection {
	clinit_SourceMapObjectParser()
	_ = section
	var offset *j2gort.Ext_com_google_gson_JsonObject = section.Get("offset").GetAsJsonObject_0()
	_ = offset
	var line int32 = offset.Get("line").GetAsInt()
	_ = line
	var column int32 = offset.Get("column").GetAsInt()
	_ = column
	if section.Has("map") && section.Has("url") {
		panic(NewSourceMapParseException("Invalid map format: section may not have both 'map' and 'url'"))
	} else if section.Has("url") {
		return SourceMapSection_ForURL(section.Get("url").GetAsString(), line, column)
	} else if section.Has("map") {
		var map_ j2gort.Ext_com_google_gson_JsonElement = section.Get("map")
		_ = map_
		var mapStr string
		_ = mapStr
		if map_.IsJsonPrimitive() && map_.GetAsJsonPrimitive_0().IsString() {
			mapStr = map_.GetAsString()
		} else {
			mapStr = map_.String()
		}
		return SourceMapSection_ForMap(mapStr, line, column)
	}
	panic(NewSourceMapParseException("Invalid map format: section must have either 'map' or 'url'"))
}

// Java static method: SourceMapObjectParser.getStringOrNull(com.google.gson.JsonObject object, java.lang.String key)
func SourceMapObjectParser_GetStringOrNull(object *j2gort.Ext_com_google_gson_JsonObject, key string) string {
	clinit_SourceMapObjectParser()
	_ = object
	_ = key
	if object.Has(key) && !(object.Get(key).IsJsonNull()) {
		return object.Get(key).GetAsString()
	} else {
		return j2gort.NullString
	}
}

// Java static method: SourceMapObjectParser.getJavaStringArray(com.google.gson.JsonElement element)
func SourceMapObjectParser_GetJavaStringArray(element j2gort.Ext_com_google_gson_JsonElement) []string {
	clinit_SourceMapObjectParser()
	_ = element
	if (element == nil) || j2gort.InstanceOf[*j2gort.Ext_com_google_gson_JsonNull](element) {
		return nil
	}
	var array *j2gort.Ext_com_google_gson_JsonArray = element.GetAsJsonArray_0()
	_ = array
	var len_ int32 = array.Size()
	_ = len_
	var result []string = j2gort.NewStringArray(len_)
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var item j2gort.Ext_com_google_gson_JsonElement = array.Get(i)
			_ = item
			{
				a_12 := result
				i_13 := i
				a_12[i_13] = func() string {
					if item.IsJsonNull() {
						return j2gort.NullString
					}
					return item.GetAsString()
				}()
			}
		}
	}
	return result
}

// Java static method: SourceMapObjectParser.readStringArray(com.google.debugging.sourcemap.SourceMapJsonLexer lexer)
func SourceMapObjectParser_ReadStringArray(lexer *SourceMapJsonLexer) []string {
	clinit_SourceMapObjectParser()
	_ = lexer
	lexer.SkipWhitespace()
	if (lexer.pos < lexer.length) && j2gort.JStrStartsWithOffset(lexer.json, "null", lexer.pos) {
		r_14 := lexer
		r_14.pos = (r_14.pos + int32(4))
		return nil
	}
	var list j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = list
	lexer.BeginArray()
	for lexer.HasNext() {
		list.Add_1(j2gort.BoxString(lexer.NextStringOrNull()))
		lexer.CheckComma()
	}
	lexer.EndArray()
	return j2gort.ColToNewStringArray(list)
}

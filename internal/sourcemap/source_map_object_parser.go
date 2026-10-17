// Port of com/google/debugging/sourcemap/SourceMapObjectParser.java @ closure-compiler v20261015

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

package sourcemap

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapObjectParser is the Java implementation of the source map
// parser.
type SourceMapObjectParser struct{}

// Gson objects are safe to share across threads and are somewhat expensive
// to construct. Use a single instance. This is Safe per gson docs:
// https://google.github.io/gson/apidocs/com/google/gson/Gson.html
var sourceMapObjectParser_gson = jutil.NewGson()

// SourceMapObjectParser_Parse parses a source map with Gson. It panics with
// a *SourceMapParseException on malformed JSON.
//
// Java: SourceMapObjectParser.parse(String)
func SourceMapObjectParser_Parse(contents string) *SourceMapObject {
	builder := SourceMapObject_Builder_()

	jutil.Try(func() {
		sourceMapRoot := sourceMapObjectParser_gson.FromJsonObject(contents)
		builder.SetVersion(int(sourceMapRoot.Get("version").GetAsInt()))
		builder.SetFile(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "file"))
		if sourceMapRoot.Has("lineCount") {
			builder.SetLineCount(int(sourceMapRoot.Get("lineCount").GetAsInt()))
		} else {
			builder.SetLineCount(-1)
		}
		builder.SetMappings(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "mappings"))
		builder.SetSourceRoot(SourceMapObjectParser_GetStringOrNull(sourceMapRoot, "sourceRoot"))

		if sourceMapRoot.Has("sections") {
			listBuilder := []*SourceMapSection{}
			for each := range sourceMapRoot.Get("sections").GetAsJsonArray().All() {
				listBuilder = append(listBuilder, SourceMapObjectParser_BuildSection(each.GetAsJsonObject()))
			}
			builder.SetSections(listBuilder)
		}

		builder.SetSources(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("sources")))
		builder.SetSourcesContent(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("sourcesContent")))
		builder.SetNames(SourceMapObjectParser_GetJavaStringArray(sourceMapRoot.Get("names")))

		extensions := jutil.NewLinkedHashMap[string, any]()
		for key, value := range sourceMapRoot.EntrySet().All() {
			if strings.HasPrefix(key, "x_") {
				extensions.Put(key, value)
			}
		}
		// Java wraps the map in Collections.unmodifiableMap; no caller
		// modifies it.
		builder.SetExtensions(extensions)
	}, func(ex jutil.IJsonParseException) {
		panic(NewSourceMapParseException("JSON parse exception: " + ex.Error()))
	})

	return builder.Build()
}

// SourceMapObjectParser_ParseFast parses a source map with
// SourceMapJsonLexer. It panics with a *SourceMapParseException.
//
// TODO: fenghaolw - Migrate everyone to use this method and rename it to
// parse.
//
// Java: SourceMapObjectParser.parseFast(String)
func SourceMapObjectParser_ParseFast(contents string) *SourceMapObject {
	builder := SourceMapObject_Builder_()
	builder.SetLineCount(-1)

	jutil.Try(func() {
		lexer := NewSourceMapJsonLexer(contents)
		lexer.BeginObject()
		extensions := jutil.NewLinkedHashMap[string, any]()

		for lexer.HasNext() {
			name := lexer.NextName()
			switch name {
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
				rawSections := lexer.NextRawValue()
				sectionsArray := sourceMapObjectParser_gson.FromJsonArray(rawSections)
				listBuilder := []*SourceMapSection{}
				for each := range sectionsArray.All() {
					listBuilder = append(listBuilder, SourceMapObjectParser_BuildSection(each.GetAsJsonObject()))
				}
				builder.SetSections(listBuilder)
			case "sources":
				builder.SetSources(SourceMapObjectParser_ReadStringArray(lexer))
			case "sourcesContent":
				builder.SetSourcesContent(SourceMapObjectParser_ReadStringArray(lexer))
			case "names":
				builder.SetNames(SourceMapObjectParser_ReadStringArray(lexer))
			default:
				if strings.HasPrefix(name, "x_") {
					extensions.Put(name, sourceMapObjectParser_gson.FromJsonElement(lexer.NextRawValue()))
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
		builder.SetExtensions(extensions)

	}, func(ex jutil.Throwable) {
		// catch (Exception ex): every Throwable but java.lang.Error.
		if _, isError := ex.(jutil.IError); isError {
			panic(ex)
		}
		panic(NewSourceMapParseException("Parse exception: " + ex.Error()))
	})

	return builder.Build()
}

// Java: SourceMapObjectParser.buildSection(JsonObject)
func SourceMapObjectParser_BuildSection(section *jutil.JsonObject) *SourceMapSection {
	offset := sourceMapObjectParser_nonNull(section.Get("offset"), "getAsJsonObject()").GetAsJsonObject()
	line := int(sourceMapObjectParser_nonNull(offset.Get("line"), "getAsInt()").GetAsInt())
	column := int(sourceMapObjectParser_nonNull(offset.Get("column"), "getAsInt()").GetAsInt())

	if section.Has("map") && section.Has("url") {
		panic(NewSourceMapParseException(
			"Invalid map format: section may not have both 'map' and 'url'"))
	} else if section.Has("url") {
		return SourceMapSection_ForURL(section.Get("url").GetAsString(), line, column)
	} else if section.Has("map") {
		m := section.Get("map")
		var mapStr string
		if m.IsJsonPrimitive() && m.GetAsJsonPrimitive().IsString() {
			mapStr = m.GetAsString()
		} else {
			mapStr = m.String()
		}
		return SourceMapSection_ForMap(mapStr, line, column)
	}

	panic(NewSourceMapParseException(
		"Invalid map format: section must have either 'map' or 'url'"))
}

// sourceMapObjectParser_nonNull returns e, the result of JsonObject.get,
// for the call of method on it.
//
// Go-only: language-adapter — calling method on a null JsonObject.get
// result throws a NullPointerException whose (JDK 21 "helpful") message
// names the call; parseFast embeds that message in its
// SourceMapParseException, so it is reproduced here.
func sourceMapObjectParser_nonNull(e jutil.JsonElement, method string) jutil.JsonElement {
	if e == nil {
		panic(jutil.NewHelpfulNullPointerException("Cannot invoke \"" + jutil.GsonPackage + ".JsonElement." + method +
			"\" because the return value of \"" + jutil.GsonPackage + ".JsonObject.get(String)\" is null"))
	}
	return e
}

// Java: SourceMapObjectParser.getStringOrNull(JsonObject, String)
func SourceMapObjectParser_GetStringOrNull(object *jutil.JsonObject, key string) *string {
	if object.Has(key) && !object.Get(key).IsJsonNull() {
		s := object.Get(key).GetAsString()
		return &s
	}
	return nil
}

// SourceMapObjectParser_GetJavaStringArray returns the array's strings (nil
// elements for JSON nulls), or nil for an absent or null element.
//
// Java: SourceMapObjectParser.getJavaStringArray(JsonElement)
func SourceMapObjectParser_GetJavaStringArray(element jutil.JsonElement) []*string {
	if element == nil {
		return nil
	}
	if _, isNull := element.(*jutil.JsonNull); isNull {
		return nil
	}
	array := element.GetAsJsonArray()
	length := array.Size()
	result := make([]*string, length)
	for i := 0; i < length; i++ {
		item := array.Get(i)
		if item.IsJsonNull() {
			result[i] = nil
		} else {
			s := item.GetAsString()
			result[i] = &s
		}
	}
	return result
}

// Java: SourceMapObjectParser.readStringArray(SourceMapJsonLexer)
func SourceMapObjectParser_ReadStringArray(lexer *SourceMapJsonLexer) []*string {
	lexer.SkipWhitespace()
	if lexer.pos < lexer.length && lexer.startsWith("null", lexer.pos) {
		lexer.pos += 4
		return nil
	}
	list := []*string{}
	lexer.BeginArray()
	for lexer.HasNext() {
		list = append(list, lexer.NextStringOrNull())
		lexer.CheckComma()
	}
	lexer.EndArray()
	return list
}

// Port of com/google/debugging/sourcemap/SourceMapObjectParser.java @ closure-compiler v20261021

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

	sourceMapRoot, parseEx := sourceMapObjectParser_gson.TryFromJsonObject(contents)
	if parseEx != nil {
		panic(NewSourceMapParseException("JSON parse exception: " + parseEx.Error()))
	}

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
	builder.SetExtensions(extensions)

	return builder.Build()
}

// SourceMapObjectParser_ParseFast parses a source map with
// SourceMapJsonLexer. It panics with a *SourceMapParseException.
//
// Java: SourceMapObjectParser.parseFast(String)
func SourceMapObjectParser_ParseFast(contents string) *SourceMapObject {
	obj, err := parseFast(contents)
	if err != nil {
		panic(err)
	}
	return obj
}

func wrapFastParseErr(err error) *SourceMapParseException {
	return NewSourceMapParseException("Parse exception: " + err.Error())
}

func parseFast(contents string) (*SourceMapObject, *SourceMapParseException) {
	builder := SourceMapObject_Builder_()
	builder.SetLineCount(-1)

	lexer := NewSourceMapJsonLexer(contents)
	if err := lexer.beginObject(); err != nil {
		return nil, wrapFastParseErr(err)
	}
	extensions := jutil.NewLinkedHashMap[string, any]()

	for lexer.HasNext() {
		name, err := lexer.nextName()
		if err != nil {
			return nil, wrapFastParseErr(err)
		}
		switch name {
		case "version":
			v, err := lexer.nextInt()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetVersion(v)
		case "file":
			s, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetFile(s)
		case "lineCount":
			lc, err := lexer.nextInt()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetLineCount(lc)
		case "mappings":
			m, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetMappings(m)
		case "sourceRoot":
			sr, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetSourceRoot(sr)
		case "sections":
			rawSections, err := lexer.nextRawValue()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			sectionsArray, parseEx := sourceMapObjectParser_gson.TryFromJsonArray(rawSections)
			if parseEx != nil {
				return nil, wrapFastParseErr(parseEx)
			}
			listBuilder := make([]*SourceMapSection, 0, sectionsArray.Size())
			for each := range sectionsArray.All() {
				if !each.IsJsonObject() {
					return nil, wrapFastParseErr(jutil.NewIllegalStateExceptionMsg("Not a JSON Object: " + each.String()))
				}
				sec, err := buildSection(each.GetAsJsonObject())
				if err != nil {
					return nil, wrapFastParseErr(err)
				}
				listBuilder = append(listBuilder, sec)
			}
			builder.SetSections(listBuilder)
		case "sources":
			arr, err := readStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetSources(arr)
		case "sourcesContent":
			arr, err := readStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetSourcesContent(arr)
		case "names":
			arr, err := readStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			builder.SetNames(arr)
		default:
			if strings.HasPrefix(name, "x_") {
				rawVal, err := lexer.nextRawValue()
				if err != nil {
					return nil, wrapFastParseErr(err)
				}
				elem, parseEx := sourceMapObjectParser_gson.TryFromJsonElement(rawVal)
				if parseEx != nil {
					return nil, wrapFastParseErr(parseEx)
				}
				extensions.Put(name, elem)
			} else {
				if err := lexer.skipValue(); err != nil {
					return nil, wrapFastParseErr(err)
				}
			}
		}
		if err := lexer.checkComma(); err != nil {
			return nil, wrapFastParseErr(err)
		}
	}
	if err := lexer.endObject(); err != nil {
		return nil, wrapFastParseErr(err)
	}
	lexer.SkipWhitespace()
	if lexer.pos < lexer.length {
		return nil, wrapFastParseErr(NewSourceMapParseException("Unexpected trailing characters"))
	}
	builder.SetExtensions(extensions)

	return builder.Build(), nil
}

// Java: SourceMapObjectParser.buildSection(JsonObject)
func SourceMapObjectParser_BuildSection(section *jutil.JsonObject) *SourceMapSection {
	sec, err := buildSection(section)
	if err != nil {
		panic(err)
	}
	return sec
}

func buildSection(section *jutil.JsonObject) (*SourceMapSection, jutil.Throwable) {
	offsetElem, err := sourceMapObjectParser_nonNullErr(section.Get("offset"), "getAsJsonObject()")
	if err != nil {
		return nil, err
	}
	if !offsetElem.IsJsonObject() {
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Object: " + offsetElem.String())
	}
	offset := offsetElem.GetAsJsonObject()

	lineElem, err := sourceMapObjectParser_nonNullErr(offset.Get("line"), "getAsInt()")
	if err != nil {
		return nil, err
	}
	line, err := lineElem.TryGetAsInt()
	if err != nil {
		return nil, err
	}

	colElem, err := sourceMapObjectParser_nonNullErr(offset.Get("column"), "getAsInt()")
	if err != nil {
		return nil, err
	}
	column, err := colElem.TryGetAsInt()
	if err != nil {
		return nil, err
	}

	if section.Has("map") && section.Has("url") {
		return nil, NewSourceMapParseException(
			"Invalid map format: section may not have both 'map' and 'url'")
	} else if section.Has("url") {
		urlStr, err := section.Get("url").TryGetAsString()
		if err != nil {
			return nil, err
		}
		return SourceMapSection_ForURL(urlStr, int(line), int(column)), nil
	} else if section.Has("map") {
		m := section.Get("map")
		var mapStr string
		if m.IsJsonPrimitive() && m.GetAsJsonPrimitive().IsString() {
			mapStr = m.GetAsString()
		} else {
			mapStr = m.String()
		}
		return SourceMapSection_ForMap(mapStr, int(line), int(column)), nil
	}

	return nil, NewSourceMapParseException(
		"Invalid map format: section must have either 'map' or 'url'")
}

// sourceMapObjectParser_nonNullErr returns e, or the JDK 21 helpful NPE naming
// the method invocation when JsonObject.get(String) returned null.
func sourceMapObjectParser_nonNullErr(e jutil.JsonElement, method string) (jutil.JsonElement, jutil.Throwable) {
	if e == nil {
		return nil, jutil.NewHelpfulNullPointerException("Cannot invoke \"" + jutil.GsonPackage + ".JsonElement." + method +
			"\" because the return value of \"" + jutil.GsonPackage + ".JsonObject.get(String)\" is null")
	}
	return e, nil
}

// Java: SourceMapObjectParser.getStringOrNull(JsonObject, String)
func SourceMapObjectParser_GetStringOrNull(object *jutil.JsonObject, key string) *string {
	if elem := object.Get(key); elem != nil && !elem.IsJsonNull() {
		s := elem.GetAsString()
		return &s
	}
	return nil
}

// SourceMapObjectParser_GetJavaStringArray returns the array's strings (nil
// elements for JSON nulls), or nil for an absent or null element.
//
// Java: SourceMapObjectParser.getJavaStringArray(JsonElement)
func SourceMapObjectParser_GetJavaStringArray(element jutil.JsonElement) []*string {
	if element == nil || element.IsJsonNull() {
		return nil
	}
	array := element.GetAsJsonArray()
	length := array.Size()
	result := make([]*string, length)
	for i := 0; i < length; i++ {
		item := array.Get(i)
		if !item.IsJsonNull() {
			s := item.GetAsString()
			result[i] = &s
		}
	}
	return result
}

func readStringArray(lexer *SourceMapJsonLexer) ([]*string, error) {
	lexer.SkipWhitespace()
	if lexer.pos < lexer.length && lexer.startsWith("null", lexer.pos) {
		lexer.pos += 4
		return nil, nil
	}
	list := []*string{}
	if err := lexer.beginArray(); err != nil {
		return nil, err
	}
	for lexer.HasNext() {
		s, err := lexer.nextStringOrNull()
		if err != nil {
			return nil, err
		}
		list = append(list, s)
		if err := lexer.checkComma(); err != nil {
			return nil, err
		}
	}
	if err := lexer.endArray(); err != nil {
		return nil, err
	}
	return list, nil
}

// Java: SourceMapObjectParser.readStringArray(SourceMapJsonLexer)
func SourceMapObjectParser_ReadStringArray(lexer *SourceMapJsonLexer) []*string {
	list, err := readStringArray(lexer)
	if err != nil {
		panic(err)
	}
	return list
}

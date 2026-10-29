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

// Gson objects are safe to share across threads and are somewhat expensive
// to construct. Use a single instance. This is Safe per gson docs:
// https://google.github.io/gson/apidocs/com/google/gson/Gson.html
var sourceMapObjectParser_gson = jutil.NewGson()

// SourceMapObjectParser_Parse parses a source map with Gson. It panics with
// a *SourceMapParseException on malformed JSON.
//
// Java: SourceMapObjectParser.parse(String)
func SourceMapObjectParser_Parse(contents string) *SourceMapObject {
	obj, err := parseSourceMapObject(contents)
	if err != nil {
		panic(err)
	}
	return obj
}

func parseSourceMapObject(contents string) (*SourceMapObject, *SourceMapParseException) {
	if obj, ok := tryParseStrictFast(contents); ok {
		return obj, nil
	}
	return parseSourceMapObjectGson(contents)
}

// tryParseStrictFast parses standard strict-JSON V3 source maps directly via
// SourceMapJsonLexer without allocating a Gson DOM. It returns (_, false) on
// index maps ("sections"), missing "version", or non-strict / malformed JSON
// so parseSourceMapObjectGson preserves exact Gson diagnostics and coercions.
func tryParseStrictFast(contents string) (*SourceMapObject, bool) {
	obj := &SourceMapObject{lineCount: -1}
	lexer := NewSourceMapJsonLexer(contents)
	if err := lexer.beginObject(); err != nil {
		return nil, false
	}
	var extensions []SourceMapExtension
	hasVersion := false

	for lexer.HasNext() {
		name, ok := strictNextName(lexer)
		if !ok {
			return nil, false
		}
		switch name {
		case "version":
			v, ok := strictNextInt(lexer)
			if !ok {
				return nil, false
			}
			obj.version = v
			hasVersion = true
		case "file":
			s, ok := strictNextStringOrNull(lexer)
			if !ok {
				return nil, false
			}
			obj.file = s
		case "lineCount":
			lc, ok := strictNextInt(lexer)
			if !ok {
				return nil, false
			}
			obj.lineCount = lc
		case "mappings":
			m, ok := strictNextStringOrNull(lexer)
			if !ok {
				return nil, false
			}
			obj.mappings = m
		case "sourceRoot":
			sr, ok := strictNextStringOrNull(lexer)
			if !ok {
				return nil, false
			}
			obj.sourceRoot = sr
		case "sections":
			return nil, false
		case "sources":
			arr, ok := strictReadOptionalStringArray(lexer)
			if !ok {
				return nil, false
			}
			obj.sources = arr
		case "sourcesContent":
			arr, ok := strictReadOptionalStringArray(lexer)
			if !ok {
				return nil, false
			}
			obj.sourcesContent = arr
		case "names":
			arr, ok := strictReadOptionalStringArray(lexer)
			if !ok {
				return nil, false
			}
			obj.names = arr
		default:
			if strings.HasPrefix(name, "x_") {
				lexer.SkipWhitespace()
				start := lexer.pos
				if !strictSkipValue(lexer) {
					return nil, false
				}
				elem, parseEx := sourceMapObjectParser_gson.TryFromJsonElement(lexer.json[start:lexer.pos])
				if parseEx != nil {
					return nil, false
				}
				extensions = putExtension(extensions, name, elem)
			} else if !strictSkipValue(lexer) {
				return nil, false
			}
		}
		if err := lexer.checkComma(); err != nil {
			return nil, false
		}
	}
	if !hasVersion || lexer.endObject() != nil {
		return nil, false
	}
	lexer.SkipWhitespace()
	if lexer.pos != lexer.length {
		return nil, false
	}
	obj.extensions = extensions
	return obj, true
}

func strictNextString(l *SourceMapJsonLexer) (string, bool) {
	l.SkipWhitespace()
	if l.pos >= l.length || l.json[l.pos] != '"' {
		return "", false
	}
	start := l.pos + 1
	pos := start
	hasEscape := false
	for pos < l.length {
		c := l.json[pos]
		if c < 0x20 {
			return "", false
		}
		if c == '"' {
			break
		}
		if c == '\\' {
			hasEscape = true
			if pos+1 >= l.length {
				return "", false
			}
			switch l.json[pos+1] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				pos += 2
			case 'u':
				if pos+6 > l.length {
					return "", false
				}
				if _, ok := parseHex4ASCII(l.json[pos+2 : pos+6]); !ok {
					return "", false
				}
				pos += 6
			default:
				return "", false
			}
			continue
		}
		pos++
	}
	if pos >= l.length {
		return "", false
	}
	var res string
	if !hasEscape {
		res = l.json[start:pos]
	} else {
		var err error
		res, err = sourceMapJsonLexer_unescape(l.json, start, pos)
		if err != nil {
			return "", false
		}
	}
	l.pos = pos + 1
	return res, true
}

func strictNextName(l *SourceMapJsonLexer) (string, bool) {
	name, ok := strictNextString(l)
	if !ok {
		return "", false
	}
	l.SkipWhitespace()
	if l.pos >= l.length || l.json[l.pos] != ':' {
		return "", false
	}
	l.pos++
	return name, true
}

func strictNextStringOrNull(l *SourceMapJsonLexer) (optionalString, bool) {
	l.SkipWhitespace()
	if l.startsWith("null", l.pos) {
		l.pos += 4
		return optionalString{}, true
	}
	s, ok := strictNextString(l)
	if !ok {
		return optionalString{}, false
	}
	return optionalString{text: s, present: true}, true
}

func strictNextInt(l *SourceMapJsonLexer) (int, bool) {
	l.SkipWhitespace()
	start := l.pos
	if l.pos < l.length && l.json[l.pos] == '-' {
		l.pos++
	}
	if l.pos >= l.length || l.json[l.pos] < '0' || l.json[l.pos] > '9' {
		return 0, false
	}
	if l.json[l.pos] == '0' {
		l.pos++
	} else {
		for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
			l.pos++
		}
	}
	if l.pos < l.length {
		c := l.json[l.pos]
		if c == '.' || c == 'e' || c == 'E' || (c >= '0' && c <= '9') {
			return 0, false
		}
	}
	v, ok := sourceMapJsonLexer_parseInt(l.json, start, l.pos, 10)
	if !ok {
		return 0, false
	}
	return int(v), true
}

func strictReadOptionalStringArray(l *SourceMapJsonLexer) ([]optionalString, bool) {
	l.SkipWhitespace()
	if l.startsWith("null", l.pos) {
		l.pos += 4
		return nil, true
	}
	if l.beginArray() != nil {
		return nil, false
	}
	list := []optionalString{}
	for l.HasNext() {
		s, ok := strictNextStringOrNull(l)
		if !ok {
			return nil, false
		}
		list = append(list, s)
		if l.checkComma() != nil {
			return nil, false
		}
	}
	if l.endArray() != nil {
		return nil, false
	}
	return list, true
}

func strictSkipValue(l *SourceMapJsonLexer) bool {
	l.SkipWhitespace()
	if l.pos >= l.length {
		return false
	}
	switch c := l.json[l.pos]; {
	case c == '{':
		l.pos++
		for l.HasNext() {
			if _, ok := strictNextName(l); !ok {
				return false
			}
			if !strictSkipValue(l) || l.checkComma() != nil {
				return false
			}
		}
		return l.endObject() == nil
	case c == '[':
		l.pos++
		for l.HasNext() {
			if !strictSkipValue(l) || l.checkComma() != nil {
				return false
			}
		}
		return l.endArray() == nil
	case c == '"':
		_, ok := strictNextString(l)
		return ok
	case c == 't' && l.startsWith("true", l.pos):
		l.pos += 4
		return true
	case c == 'f' && l.startsWith("false", l.pos):
		l.pos += 5
		return true
	case c == 'n' && l.startsWith("null", l.pos):
		l.pos += 4
		return true
	case c == '-' || (c >= '0' && c <= '9'):
		start := l.pos
		if c == '-' {
			l.pos++
		}
		if l.pos >= l.length || l.json[l.pos] < '0' || l.json[l.pos] > '9' {
			return false
		}
		if l.json[l.pos] == '0' {
			l.pos++
			if l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
				return false
			}
		} else {
			for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
				l.pos++
			}
		}
		if l.pos < l.length && l.json[l.pos] == '.' {
			l.pos++
			if l.pos >= l.length || l.json[l.pos] < '0' || l.json[l.pos] > '9' {
				return false
			}
			for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
				l.pos++
			}
		}
		if l.pos < l.length && (l.json[l.pos] == 'e' || l.json[l.pos] == 'E') {
			l.pos++
			if l.pos < l.length && (l.json[l.pos] == '+' || l.json[l.pos] == '-') {
				l.pos++
			}
			if l.pos >= l.length || l.json[l.pos] < '0' || l.json[l.pos] > '9' {
				return false
			}
			expStart := l.pos
			for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
				l.pos++
			}
			if l.pos-expStart > 3 {
				return false
			}
		}
		return l.pos-start <= 32
	default:
		return false
	}
}

func parseSourceMapObjectGson(contents string) (*SourceMapObject, *SourceMapParseException) {
	sourceMapRoot, parseEx := sourceMapObjectParser_gson.TryFromJsonObject(contents)
	if parseEx != nil {
		return nil, NewSourceMapParseException("JSON parse exception: " + parseEx.Error())
	}

	obj := &SourceMapObject{
		version:   int(sourceMapRoot.Get("version").GetAsInt()),
		file:      getOptionalString(sourceMapRoot, "file"),
		lineCount: -1,
	}
	if sourceMapRoot.Has("lineCount") {
		obj.lineCount = int(sourceMapRoot.Get("lineCount").GetAsInt())
	}
	obj.mappings = getOptionalString(sourceMapRoot, "mappings")
	obj.sourceRoot = getOptionalString(sourceMapRoot, "sourceRoot")

	if sourceMapRoot.Has("sections") {
		sectionsArray := sourceMapRoot.Get("sections").GetAsJsonArray()
		listBuilder := make([]*SourceMapSection, 0, sectionsArray.Size())
		for each := range sectionsArray.All() {
			sec, err := buildSection(each.GetAsJsonObject())
			if err != nil {
				if smErr, ok := err.(*SourceMapParseException); ok {
					return nil, smErr
				}
				panic(err)
			}
			listBuilder = append(listBuilder, sec)
		}
		obj.sections = listBuilder
	}

	obj.sources = getOptionalStringArray(sourceMapRoot.Get("sources"))
	obj.sourcesContent = getOptionalStringArray(sourceMapRoot.Get("sourcesContent"))
	obj.names = getOptionalStringArray(sourceMapRoot.Get("names"))

	var extensions []SourceMapExtension
	for key, value := range sourceMapRoot.EntrySet().All() {
		if strings.HasPrefix(key, "x_") {
			extensions = putExtension(extensions, key, value)
		}
	}
	obj.extensions = extensions

	return obj, nil
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
	obj := &SourceMapObject{lineCount: -1}

	lexer := NewSourceMapJsonLexer(contents)
	if err := lexer.beginObject(); err != nil {
		return nil, wrapFastParseErr(err)
	}
	var extensions []SourceMapExtension

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
			obj.version = v
		case "file":
			s, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.file = s
		case "lineCount":
			lc, err := lexer.nextInt()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.lineCount = lc
		case "mappings":
			m, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.mappings = m
		case "sourceRoot":
			sr, err := lexer.nextStringOrNull()
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.sourceRoot = sr
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
			obj.sections = listBuilder
		case "sources":
			arr, err := readOptionalStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.sources = arr
		case "sourcesContent":
			arr, err := readOptionalStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.sourcesContent = arr
		case "names":
			arr, err := readOptionalStringArray(lexer)
			if err != nil {
				return nil, wrapFastParseErr(err)
			}
			obj.names = arr
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
				extensions = putExtension(extensions, name, elem)
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
	obj.extensions = extensions

	return obj, nil
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

func getOptionalString(object *jutil.JsonObject, key string) optionalString {
	if elem := object.Get(key); elem != nil && !elem.IsJsonNull() {
		return optionalString{text: elem.GetAsString(), present: true}
	}
	return optionalString{}
}

func getOptionalStringArray(element jutil.JsonElement) []optionalString {
	if element == nil || element.IsJsonNull() {
		return nil
	}
	array := element.GetAsJsonArray()
	length := array.Size()
	result := make([]optionalString, length)
	for i := 0; i < length; i++ {
		item := array.Get(i)
		if !item.IsJsonNull() {
			result[i] = optionalString{text: item.GetAsString(), present: true}
		}
	}
	return result
}

func readOptionalStringArray(lexer *SourceMapJsonLexer) ([]optionalString, error) {
	lexer.SkipWhitespace()
	if lexer.pos < lexer.length && lexer.startsWith("null", lexer.pos) {
		lexer.pos += 4
		return nil, nil
	}
	list := []optionalString{}
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

// Port of com/google/debugging/sourcemap/SourceMapObjectParser.java @ closure-compiler v20261028
// Port of com/google/debugging/sourcemap/SourceMapJsonLexer.java @ closure-compiler v20261028

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

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func wrapScannerErr(err jutil.Throwable) jutil.Throwable {
	switch err.(type) {
	case jutil.IIOException, jutil.IIllegalStateException:
		return NewSourceMapParseException(
			"JSON parse exception: " + jutil.NewJsonSyntaxExceptionWithCause(err).Error())
	default:
		return err
	}
}

func wrapGsonSyntaxDirect(msg string) *SourceMapParseException {
	return NewSourceMapParseException(
		"JSON parse exception: " + jutil.GsonPackage + ".JsonSyntaxException: " + msg)
}

// ParseSourceMapObject parses a source map JSON document.
func ParseSourceMapObject(contents string) (*SourceMapObject, error) {
	obj, err := parseSourceMapObject(contents)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func parseSourceMapObject(contents string) (*SourceMapObject, jutil.Throwable) {
	var s jutil.JsonScanner
	s.Init(contents, true)

	rootTok, emptyEOF, err := s.DoPeek(true)
	if err != nil {
		return nil, wrapScannerErr(err)
	}
	if emptyEOF {
		return nil, jutil.NewHelpfulNullPointerException(
			"Cannot invoke \"" + jutil.GsonPackage + ".JsonObject.get(String)\" because \"sourceMapRoot\" is null")
	}

	if rootTok != jutil.JsonPeekBeginObject {
		elem, err := s.ReadJsonElement()
		if err != nil {
			return nil, wrapScannerErr(err)
		}
		var actualSimple string
		switch elem.(type) {
		case *jutil.JsonArray:
			actualSimple = "JsonArray"
		case *jutil.JsonPrimitive:
			actualSimple = "JsonPrimitive"
		default:
			actualSimple = "JsonNull"
		}
		return nil, wrapGsonSyntaxDirect(
			"Expected a " + jutil.GsonPackage + ".JsonObject but was " + jutil.GsonPackage + "." + actualSimple + "; at path $")
	}

	if err := s.BeginObject(); err != nil {
		return nil, wrapScannerErr(err)
	}

	const (
		errVersion = iota
		errFile
		errLineCount
		errMappings
		errSourceRoot
		errSections
		errSources
		errSourcesContent
		errNames
	)
	obj := &SourceMapObject{lineCount: -1}
	hasVersion := false
	var fieldErrs [9]jutil.Throwable
	var extensions []SourceMapExtension

	for {
		more, err := s.HasNext()
		if err != nil {
			return nil, wrapScannerErr(err)
		}
		if !more {
			break
		}
		name, err := s.NextName()
		if err != nil {
			return nil, wrapScannerErr(err)
		}
		vtok, err := s.Peek()
		if err != nil {
			return nil, wrapScannerErr(err)
		}

		switch name {
		case "version":
			v, cerr, serr := readStreamInt(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.version = v
			fieldErrs[errVersion] = cerr
			hasVersion = true

		case "file":
			str, cerr, serr := readStreamOptionalString(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.file = str
			fieldErrs[errFile] = cerr

		case "lineCount":
			v, cerr, serr := readStreamInt(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.lineCount = v
			fieldErrs[errLineCount] = cerr

		case "mappings":
			str, cerr, serr := readStreamOptionalString(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.mappings = str
			fieldErrs[errMappings] = cerr

		case "sourceRoot":
			str, cerr, serr := readStreamOptionalString(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.sourceRoot = str
			fieldErrs[errSourceRoot] = cerr

		case "sections":
			secs, cerr, serr := readStreamSections(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.sections = secs
			fieldErrs[errSections] = cerr

		case "sources":
			arr, cerr, serr := readStreamOptionalStringArray(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.sources = arr
			fieldErrs[errSources] = cerr

		case "sourcesContent":
			arr, cerr, serr := readStreamOptionalStringArray(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.sourcesContent = arr
			fieldErrs[errSourcesContent] = cerr

		case "names":
			arr, cerr, serr := readStreamOptionalStringArray(&s, vtok)
			if serr != nil {
				return nil, wrapScannerErr(serr)
			}
			obj.names = arr
			fieldErrs[errNames] = cerr

		default:
			if strings.HasPrefix(name, "x_") {
				elem, err := s.ReadJsonElement()
				if err != nil {
					return nil, wrapScannerErr(err)
				}
				extensions = putExtension(extensions, name, elem)
			} else if err := s.SkipValueAsElement(); err != nil {
				return nil, wrapScannerErr(err)
			}
		}
	}

	if err := s.EndObject(); err != nil {
		return nil, wrapScannerErr(err)
	}

	// Assert full consumption in strict mode (Gson.assertFullConsumption).
	s.SetLenient(false)
	endTok, err := s.Peek()
	if err != nil {
		return nil, wrapScannerErr(err)
	}
	if endTok != jutil.JsonPeekEOF {
		return nil, wrapGsonSyntaxDirect("JSON document was not fully consumed.")
	}

	if !hasVersion {
		return nil, sourceMapParserMissingFieldNPE("getAsInt()")
	}
	for _, ferr := range fieldErrs {
		if ferr != nil {
			return nil, ferr
		}
	}

	obj.extensions = extensions
	return obj, nil
}

func sourceMapParserMissingFieldNPE(method string) jutil.Throwable {
	return jutil.NewHelpfulNullPointerException(
		"Cannot invoke \"" + jutil.GsonPackage + ".JsonElement." + method +
			"\" because the return value of \"" + jutil.GsonPackage + ".JsonObject.get(String)\" is null")
}

func readStreamInt(s *jutil.JsonScanner, vtok uint8) (val int, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok == jutil.JsonPeekLong && s.PeekedLong() >= math.MinInt32 && s.PeekedLong() <= math.MaxInt32 {
		v := int(s.PeekedLong())
		s.ConsumeScalar()
		return v, nil, nil
	}
	elem, err := s.ReadJsonElement()
	if err != nil {
		return 0, nil, err
	}
	v, cerr := elem.TryGetAsInt()
	return int(v), cerr, nil
}

func readStreamOptionalString(s *jutil.JsonScanner, vtok uint8) (val OptionalString, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	switch vtok {
	case jutil.JsonPeekNull:
		s.ConsumeScalar()
		return OptionalString{}, nil, nil
	case jutil.JsonPeekDoubleQuoted, jutil.JsonPeekSingleQuoted, jutil.JsonPeekUnquoted:
		str, err := s.NextString()
		if err != nil {
			return OptionalString{}, nil, err
		}
		return SomeString(str), nil, nil
	default:
		elem, err := s.ReadJsonElement()
		if err != nil {
			return OptionalString{}, nil, err
		}
		str, cerr := elem.TryGetAsString()
		if cerr != nil {
			return OptionalString{}, cerr, nil
		}
		return SomeString(str), nil, nil
	}
}

func readStreamOptionalStringArray(s *jutil.JsonScanner, vtok uint8) (val []OptionalString, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok == jutil.JsonPeekNull {
		s.ConsumeScalar()
		return nil, nil, nil
	}
	if vtok != jutil.JsonPeekBeginArray {
		elem, err := s.ReadJsonElement()
		if err != nil {
			return nil, nil, err
		}
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Array: " + elem.String()), nil
	}
	if err := s.BeginArray(); err != nil {
		return nil, nil, err
	}
	list := []OptionalString{}
	var firstErr jutil.Throwable
	for {
		more, err := s.HasNext()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		etok, err := s.Peek()
		if err != nil {
			return nil, nil, err
		}
		switch etok {
		case jutil.JsonPeekNull:
			s.ConsumeScalar()
			list = append(list, OptionalString{})
		case jutil.JsonPeekDoubleQuoted, jutil.JsonPeekSingleQuoted, jutil.JsonPeekUnquoted:
			str, err := s.NextString()
			if err != nil {
				return nil, nil, err
			}
			list = append(list, SomeString(str))
		default:
			elem, err := s.ReadJsonElement()
			if err != nil {
				return nil, nil, err
			}
			if firstErr == nil {
				str, cerr := elem.TryGetAsString()
				if cerr != nil {
					firstErr = cerr
				} else {
					list = append(list, SomeString(str))
				}
			}
		}
	}
	if err := s.EndArray(); err != nil {
		return nil, nil, err
	}
	if firstErr != nil {
		return nil, firstErr, nil
	}
	return list, nil, nil
}

func readStreamSections(s *jutil.JsonScanner, vtok uint8) (secs []*SourceMapSection, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok != jutil.JsonPeekBeginArray {
		elem, err := s.ReadJsonElement()
		if err != nil {
			return nil, nil, err
		}
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Array: " + elem.String()), nil
	}
	if err := s.BeginArray(); err != nil {
		return nil, nil, err
	}
	sections := []*SourceMapSection{}
	var firstErr jutil.Throwable
	for {
		more, err := s.HasNext()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		etok, err := s.Peek()
		if err != nil {
			return nil, nil, err
		}
		if etok != jutil.JsonPeekBeginObject {
			elem, err := s.ReadJsonElement()
			if err != nil {
				return nil, nil, err
			}
			if firstErr == nil {
				firstErr = jutil.NewIllegalStateExceptionMsg("Not a JSON Object: " + elem.String())
			}
			continue
		}
		sec, cerr, serr := readStreamSectionObject(s)
		if serr != nil {
			return nil, nil, serr
		}
		if cerr != nil {
			if firstErr == nil {
				firstErr = cerr
			}
		} else {
			sections = append(sections, sec)
		}
	}
	if err := s.EndArray(); err != nil {
		return nil, nil, err
	}
	if firstErr != nil {
		return nil, firstErr, nil
	}
	return sections, nil, nil
}

func readStreamSectionObject(s *jutil.JsonScanner) (sec *SourceMapSection, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if err := s.BeginObject(); err != nil {
		return nil, nil, err
	}

	hasOffset := false
	offsetIsObj := false
	var offsetStr string
	hasLine, hasCol := false, false
	var lineVal, colVal int
	var lineErr, colErr jutil.Throwable

	hasURL, hasMap := false, false
	var urlStr, mapStr string
	var urlErr jutil.Throwable

	for {
		more, err := s.HasNext()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		key, err := s.NextName()
		if err != nil {
			return nil, nil, err
		}
		ktok, err := s.Peek()
		if err != nil {
			return nil, nil, err
		}

		switch key {
		case "offset":
			hasOffset = true
			if ktok != jutil.JsonPeekBeginObject {
				elem, err := s.ReadJsonElement()
				if err != nil {
					return nil, nil, err
				}
				offsetIsObj = false
				offsetStr = elem.String()
				continue
			}
			offsetIsObj = true
			hasLine, hasCol = false, false
			lineErr, colErr = nil, nil
			if err := s.BeginObject(); err != nil {
				return nil, nil, err
			}
			for {
				omore, err := s.HasNext()
				if err != nil {
					return nil, nil, err
				}
				if !omore {
					break
				}
				okey, err := s.NextName()
				if err != nil {
					return nil, nil, err
				}
				ovtok, err := s.Peek()
				if err != nil {
					return nil, nil, err
				}
				switch okey {
				case "line":
					v, cerr, serr := readStreamInt(s, ovtok)
					if serr != nil {
						return nil, nil, serr
					}
					lineVal = v
					lineErr = cerr
					hasLine = true
				case "column":
					v, cerr, serr := readStreamInt(s, ovtok)
					if serr != nil {
						return nil, nil, serr
					}
					colVal = v
					colErr = cerr
					hasCol = true
				default:
					if err := s.SkipValueAsElement(); err != nil {
						return nil, nil, err
					}
				}
			}
			if err := s.EndObject(); err != nil {
				return nil, nil, err
			}

		case "url":
			hasURL = true
			if ktok == jutil.JsonPeekDoubleQuoted || ktok == jutil.JsonPeekSingleQuoted || ktok == jutil.JsonPeekUnquoted {
				str, err := s.NextString()
				if err != nil {
					return nil, nil, err
				}
				urlStr = str
				urlErr = nil
			} else {
				elem, err := s.ReadJsonElement()
				if err != nil {
					return nil, nil, err
				}
				str, cerr := elem.TryGetAsString()
				urlStr = str
				urlErr = cerr
			}

		case "map":
			hasMap = true
			if ktok == jutil.JsonPeekDoubleQuoted || ktok == jutil.JsonPeekSingleQuoted || ktok == jutil.JsonPeekUnquoted {
				str, err := s.NextString()
				if err != nil {
					return nil, nil, err
				}
				mapStr = str
			} else {
				str, err := readStreamCompactJson(s)
				if err != nil {
					return nil, nil, err
				}
				mapStr = str
			}

		default:
			if err := s.SkipValueAsElement(); err != nil {
				return nil, nil, err
			}
		}
	}

	if err := s.EndObject(); err != nil {
		return nil, nil, err
	}

	if !hasOffset {
		return nil, sourceMapParserMissingFieldNPE("getAsJsonObject()"), nil
	}
	if !offsetIsObj {
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Object: " + offsetStr), nil
	}
	if !hasLine {
		return nil, sourceMapParserMissingFieldNPE("getAsInt()"), nil
	}
	if lineErr != nil {
		return nil, lineErr, nil
	}
	if !hasCol {
		return nil, sourceMapParserMissingFieldNPE("getAsInt()"), nil
	}
	if colErr != nil {
		return nil, colErr, nil
	}
	if hasMap && hasURL {
		return nil, NewSourceMapParseException("Invalid map format: section may not have both 'map' and 'url'"), nil
	}
	if hasURL {
		if urlErr != nil {
			return nil, urlErr, nil
		}
		return SourceMapSection_ForURL(urlStr, lineVal, colVal), nil, nil
	}
	if hasMap {
		return SourceMapSection_ForMap(mapStr, lineVal, colVal), nil, nil
	}
	return nil, NewSourceMapParseException("Invalid map format: section must have either 'map' or 'url'"), nil
}

// readStreamCompactJson formats the next JSON value into Gson's compact
// JsonElement.toString() form without allocating a DOM when object keys are
// unique, falling back to ReadJsonElement on duplicate keys.
func readStreamCompactJson(s *jutil.JsonScanner) (string, jutil.Throwable) {
	snap := s.Clone()
	var sb strings.Builder
	dup, err := appendStreamCompactValue(s, &sb)
	if err != nil {
		return "", err
	}
	if !dup {
		return sb.String(), nil
	}
	*s = snap
	elem, err := s.ReadJsonElement()
	if err != nil {
		return "", err
	}
	return elem.String(), nil
}

func appendStreamCompactValue(s *jutil.JsonScanner, sb *strings.Builder) (dupKey bool, err jutil.Throwable) {
	tok, err := s.Peek()
	if err != nil {
		return false, err
	}
	switch tok {
	case jutil.JsonPeekBeginObject:
		if err := s.BeginObject(); err != nil {
			return false, err
		}
		sb.WriteByte('{')
		var seenBuf [12]string
		seen := seenBuf[:0]
		first := true
		for {
			more, err := s.HasNext()
			if err != nil {
				return false, err
			}
			if !more {
				break
			}
			name, err := s.NextName()
			if err != nil {
				return false, err
			}
			for _, prev := range seen {
				if prev == name {
					return true, nil
				}
			}
			seen = append(seen, name)
			if !first {
				sb.WriteByte(',')
			}
			first = false
			_ = jutil.WriteGsonString(sb, name, false)
			sb.WriteByte(':')
			dup, err := appendStreamCompactValue(s, sb)
			if err != nil || dup {
				return dup, err
			}
		}
		if err := s.EndObject(); err != nil {
			return false, err
		}
		sb.WriteByte('}')
		return false, nil

	case jutil.JsonPeekBeginArray:
		if err := s.BeginArray(); err != nil {
			return false, err
		}
		sb.WriteByte('[')
		first := true
		for {
			more, err := s.HasNext()
			if err != nil {
				return false, err
			}
			if !more {
				break
			}
			if !first {
				sb.WriteByte(',')
			}
			first = false
			dup, err := appendStreamCompactValue(s, sb)
			if err != nil || dup {
				return dup, err
			}
		}
		if err := s.EndArray(); err != nil {
			return false, err
		}
		sb.WriteByte(']')
		return false, nil

	case jutil.JsonPeekDoubleQuoted, jutil.JsonPeekSingleQuoted, jutil.JsonPeekUnquoted:
		str, err := s.NextString()
		if err != nil {
			return false, err
		}
		_ = jutil.WriteGsonString(sb, str, false)
		return false, nil

	case jutil.JsonPeekLong:
		sb.WriteString(strconv.FormatInt(s.PeekedLong(), 10))
		s.ConsumeScalar()
		return false, nil

	case jutil.JsonPeekNumber:
		str, err := s.NextString()
		if err != nil {
			return false, err
		}
		sb.WriteString(str)
		return false, nil

	case jutil.JsonPeekTrue:
		sb.WriteString("true")
		s.ConsumeScalar()
		return false, nil

	case jutil.JsonPeekFalse:
		sb.WriteString("false")
		s.ConsumeScalar()
		return false, nil

	case jutil.JsonPeekNull:
		sb.WriteString("null")
		s.ConsumeScalar()
		return false, nil

	default:
		return false, jutil.NewIllegalStateExceptionMsg("Unexpected token")
	}
}

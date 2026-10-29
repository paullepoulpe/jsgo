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
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapObjectParser_Parse parses a source map JSON document. It panics with
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

func parseSourceMapObject(contents string) (*SourceMapObject, jutil.Throwable) {
	var l SourceMapJsonLexer
	l.initStream(contents)

	rootTok, emptyEOF, err := l.doPeek(true)
	if err != nil {
		return nil, err
	}
	if emptyEOF {
		return nil, jutil.NewHelpfulNullPointerException(
			"Cannot invoke \"" + jutil.GsonPackage + ".JsonObject.get(String)\" because \"sourceMapRoot\" is null")
	}

	if rootTok != smPeekBeginObject {
		elem, err := l.readJsonElement()
		if err != nil {
			return nil, err
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

	if err := l.beginObjectStream(); err != nil {
		return nil, err
	}

	obj := &SourceMapObject{lineCount: -1}
	hasVersion := false
	var verErr, fileErr, lcErr, mappingsErr, sourceRootErr jutil.Throwable
	var sectionsErr, sourcesErr, sourcesContentErr, namesErr jutil.Throwable
	var extensions []SourceMapExtension

	for {
		more, err := l.hasNextStream()
		if err != nil {
			return nil, err
		}
		if !more {
			break
		}
		name, err := l.nextNameStream()
		if err != nil {
			return nil, err
		}
		vtok, err := l.peekStream()
		if err != nil {
			return nil, err
		}

		switch name {
		case "version":
			v, cerr, serr := readStreamInt(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.version = v
			verErr = cerr
			hasVersion = true

		case "file":
			s, cerr, serr := readStreamOptionalString(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.file = s
			fileErr = cerr

		case "lineCount":
			v, cerr, serr := readStreamInt(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.lineCount = v
			lcErr = cerr

		case "mappings":
			s, cerr, serr := readStreamOptionalString(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.mappings = s
			mappingsErr = cerr

		case "sourceRoot":
			s, cerr, serr := readStreamOptionalString(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.sourceRoot = s
			sourceRootErr = cerr

		case "sections":
			secs, cerr, serr := readStreamSections(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.sections = secs
			sectionsErr = cerr

		case "sources":
			arr, cerr, serr := readStreamOptionalStringArray(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.sources = arr
			sourcesErr = cerr

		case "sourcesContent":
			arr, cerr, serr := readStreamOptionalStringArray(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.sourcesContent = arr
			sourcesContentErr = cerr

		case "names":
			arr, cerr, serr := readStreamOptionalStringArray(&l, vtok)
			if serr != nil {
				return nil, serr
			}
			obj.names = arr
			namesErr = cerr

		default:
			if strings.HasPrefix(name, "x_") {
				elem, err := l.readJsonElement()
				if err != nil {
					return nil, err
				}
				extensions = putExtension(extensions, name, elem)
			} else if err := l.skipValueStream(); err != nil {
				return nil, err
			}
		}
	}

	if err := l.endObjectStream(); err != nil {
		return nil, err
	}

	// Assert full consumption in strict mode (Gson.assertFullConsumption).
	l.lenient = false
	endTok, err := l.peekStream()
	if err != nil {
		return nil, err
	}
	if endTok != smPeekEOF {
		return nil, wrapGsonSyntaxDirect("JSON document was not fully consumed.")
	}

	if !hasVersion {
		return nil, sourceMapParserMissingFieldNPE("getAsInt()")
	}
	if verErr != nil {
		return nil, verErr
	}
	if fileErr != nil {
		return nil, fileErr
	}
	if lcErr != nil {
		return nil, lcErr
	}
	if mappingsErr != nil {
		return nil, mappingsErr
	}
	if sourceRootErr != nil {
		return nil, sourceRootErr
	}
	if sectionsErr != nil {
		return nil, sectionsErr
	}
	if sourcesErr != nil {
		return nil, sourcesErr
	}
	if sourcesContentErr != nil {
		return nil, sourcesContentErr
	}
	if namesErr != nil {
		return nil, namesErr
	}

	obj.extensions = extensions
	return obj, nil
}

func sourceMapParserMissingFieldNPE(method string) jutil.Throwable {
	return jutil.NewHelpfulNullPointerException(
		"Cannot invoke \"" + jutil.GsonPackage + ".JsonElement." + method +
			"\" because the return value of \"" + jutil.GsonPackage + ".JsonObject.get(String)\" is null")
}

func readStreamInt(l *SourceMapJsonLexer, vtok uint8) (val int, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok == smPeekLong && l.peekedLong >= math.MinInt32 && l.peekedLong <= math.MaxInt32 {
		v := int(l.peekedLong)
		l.consumePeekedScalar()
		return v, nil, nil
	}
	elem, err := l.readJsonElement()
	if err != nil {
		return 0, nil, err
	}
	v, cerr := elem.TryGetAsInt()
	return int(v), cerr, nil
}

func readStreamOptionalString(l *SourceMapJsonLexer, vtok uint8) (val optionalString, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	switch vtok {
	case smPeekNull:
		l.consumePeekedScalar()
		return optionalString{}, nil, nil
	case smPeekDoubleQuoted, smPeekSingleQuoted, smPeekUnquoted:
		s, err := l.nextStringStream()
		if err != nil {
			return optionalString{}, nil, err
		}
		return optionalString{text: s, present: true}, nil, nil
	default:
		elem, err := l.readJsonElement()
		if err != nil {
			return optionalString{}, nil, err
		}
		s, cerr := elem.TryGetAsString()
		if cerr != nil {
			return optionalString{}, cerr, nil
		}
		return optionalString{text: s, present: true}, nil, nil
	}
}

func readStreamOptionalStringArray(l *SourceMapJsonLexer, vtok uint8) (val []optionalString, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok == smPeekNull {
		l.consumePeekedScalar()
		return nil, nil, nil
	}
	if vtok != smPeekBeginArray {
		elem, err := l.readJsonElement()
		if err != nil {
			return nil, nil, err
		}
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Array: " + elem.String()), nil
	}
	if err := l.beginArrayStream(); err != nil {
		return nil, nil, err
	}
	list := []optionalString{}
	var firstErr jutil.Throwable
	for {
		more, err := l.hasNextStream()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		etok, err := l.peekStream()
		if err != nil {
			return nil, nil, err
		}
		switch etok {
		case smPeekNull:
			l.consumePeekedScalar()
			list = append(list, optionalString{})
		case smPeekDoubleQuoted, smPeekSingleQuoted, smPeekUnquoted:
			s, err := l.nextStringStream()
			if err != nil {
				return nil, nil, err
			}
			list = append(list, optionalString{text: s, present: true})
		default:
			elem, err := l.readJsonElement()
			if err != nil {
				return nil, nil, err
			}
			if firstErr == nil {
				s, cerr := elem.TryGetAsString()
				if cerr != nil {
					firstErr = cerr
				} else {
					list = append(list, optionalString{text: s, present: true})
				}
			}
		}
	}
	if err := l.endArrayStream(); err != nil {
		return nil, nil, err
	}
	if firstErr != nil {
		return nil, firstErr, nil
	}
	return list, nil, nil
}

func readStreamSections(l *SourceMapJsonLexer, vtok uint8) (secs []*SourceMapSection, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if vtok != smPeekBeginArray {
		elem, err := l.readJsonElement()
		if err != nil {
			return nil, nil, err
		}
		return nil, jutil.NewIllegalStateExceptionMsg("Not a JSON Array: " + elem.String()), nil
	}
	if err := l.beginArrayStream(); err != nil {
		return nil, nil, err
	}
	sections := []*SourceMapSection{}
	var firstErr jutil.Throwable
	for {
		more, err := l.hasNextStream()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		etok, err := l.peekStream()
		if err != nil {
			return nil, nil, err
		}
		if etok != smPeekBeginObject {
			elem, err := l.readJsonElement()
			if err != nil {
				return nil, nil, err
			}
			if firstErr == nil {
				firstErr = jutil.NewIllegalStateExceptionMsg("Not a JSON Object: " + elem.String())
			}
			continue
		}
		sec, cerr, serr := readStreamSectionObject(l)
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
	if err := l.endArrayStream(); err != nil {
		return nil, nil, err
	}
	if firstErr != nil {
		return nil, firstErr, nil
	}
	return sections, nil, nil
}

func readStreamSectionObject(l *SourceMapJsonLexer) (sec *SourceMapSection, coerceErr jutil.Throwable, syntaxErr jutil.Throwable) {
	if err := l.beginObjectStream(); err != nil {
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
		more, err := l.hasNextStream()
		if err != nil {
			return nil, nil, err
		}
		if !more {
			break
		}
		key, err := l.nextNameStream()
		if err != nil {
			return nil, nil, err
		}
		ktok, err := l.peekStream()
		if err != nil {
			return nil, nil, err
		}

		switch key {
		case "offset":
			hasOffset = true
			if ktok != smPeekBeginObject {
				elem, err := l.readJsonElement()
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
			if err := l.beginObjectStream(); err != nil {
				return nil, nil, err
			}
			for {
				omore, err := l.hasNextStream()
				if err != nil {
					return nil, nil, err
				}
				if !omore {
					break
				}
				okey, err := l.nextNameStream()
				if err != nil {
					return nil, nil, err
				}
				ovtok, err := l.peekStream()
				if err != nil {
					return nil, nil, err
				}
				switch okey {
				case "line":
					v, cerr, serr := readStreamInt(l, ovtok)
					if serr != nil {
						return nil, nil, serr
					}
					lineVal = v
					lineErr = cerr
					hasLine = true
				case "column":
					v, cerr, serr := readStreamInt(l, ovtok)
					if serr != nil {
						return nil, nil, serr
					}
					colVal = v
					colErr = cerr
					hasCol = true
				default:
					if err := l.skipValueStream(); err != nil {
						return nil, nil, err
					}
				}
			}
			if err := l.endObjectStream(); err != nil {
				return nil, nil, err
			}

		case "url":
			hasURL = true
			if ktok == smPeekDoubleQuoted || ktok == smPeekSingleQuoted || ktok == smPeekUnquoted {
				s, err := l.nextStringStream()
				if err != nil {
					return nil, nil, err
				}
				urlStr = s
				urlErr = nil
			} else {
				elem, err := l.readJsonElement()
				if err != nil {
					return nil, nil, err
				}
				s, cerr := elem.TryGetAsString()
				urlStr = s
				urlErr = cerr
			}

		case "map":
			hasMap = true
			if ktok == smPeekDoubleQuoted || ktok == smPeekSingleQuoted || ktok == smPeekUnquoted {
				s, err := l.nextStringStream()
				if err != nil {
					return nil, nil, err
				}
				mapStr = s
			} else {
				s, err := readStreamCompactJson(l)
				if err != nil {
					return nil, nil, err
				}
				mapStr = s
			}

		default:
			if err := l.skipValueStream(); err != nil {
				return nil, nil, err
			}
		}
	}

	if err := l.endObjectStream(); err != nil {
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
// unique, falling back to readJsonElement on duplicate keys.
func readStreamCompactJson(l *SourceMapJsonLexer) (string, jutil.Throwable) {
	snap := l.cloneState()
	var sb strings.Builder
	dup, err := appendStreamCompactValue(l, &sb)
	if err != nil {
		return "", err
	}
	if !dup {
		return sb.String(), nil
	}
	*l = snap
	elem, err := l.readJsonElement()
	if err != nil {
		return "", err
	}
	return elem.String(), nil
}

func appendStreamCompactValue(l *SourceMapJsonLexer, sb *strings.Builder) (dupKey bool, err jutil.Throwable) {
	tok, err := l.peekStream()
	if err != nil {
		return false, err
	}
	switch tok {
	case smPeekBeginObject:
		if err := l.beginObjectStream(); err != nil {
			return false, err
		}
		sb.WriteByte('{')
		var seenBuf [12]string
		seen := seenBuf[:0]
		first := true
		for {
			more, err := l.hasNextStream()
			if err != nil {
				return false, err
			}
			if !more {
				break
			}
			name, err := l.nextNameStream()
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
			writeCompactJsonString(sb, name)
			sb.WriteByte(':')
			dup, err := appendStreamCompactValue(l, sb)
			if err != nil || dup {
				return dup, err
			}
		}
		if err := l.endObjectStream(); err != nil {
			return false, err
		}
		sb.WriteByte('}')
		return false, nil

	case smPeekBeginArray:
		if err := l.beginArrayStream(); err != nil {
			return false, err
		}
		sb.WriteByte('[')
		first := true
		for {
			more, err := l.hasNextStream()
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
			dup, err := appendStreamCompactValue(l, sb)
			if err != nil || dup {
				return dup, err
			}
		}
		if err := l.endArrayStream(); err != nil {
			return false, err
		}
		sb.WriteByte(']')
		return false, nil

	case smPeekDoubleQuoted, smPeekSingleQuoted, smPeekUnquoted:
		s, err := l.nextStringStream()
		if err != nil {
			return false, err
		}
		writeCompactJsonString(sb, s)
		return false, nil

	case smPeekLong:
		sb.WriteString(strconv.FormatInt(l.peekedLong, 10))
		l.consumePeekedScalar()
		return false, nil

	case smPeekNumber:
		s, err := l.nextStringStream()
		if err != nil {
			return false, err
		}
		sb.WriteString(s)
		return false, nil

	case smPeekTrue:
		sb.WriteString("true")
		l.consumePeekedScalar()
		return false, nil

	case smPeekFalse:
		sb.WriteString("false")
		l.consumePeekedScalar()
		return false, nil

	case smPeekNull:
		sb.WriteString("null")
		l.consumePeekedScalar()
		return false, nil

	default:
		return false, jutil.NewIllegalStateExceptionMsg("Unexpected token")
	}
}

func writeCompactJsonString(sb *strings.Builder, s string) {
	sb.WriteByte('"')
	last := 0
	for i := 0; i < len(s); {
		b := s[i]
		if b < 0x80 {
			var esc string
			switch b {
			case '"':
				esc = `\"`
			case '\\':
				esc = `\\`
			case '\b':
				esc = `\b`
			case '\f':
				esc = `\f`
			case '\n':
				esc = `\n`
			case '\r':
				esc = `\r`
			case '\t':
				esc = `\t`
			default:
				if b < 0x20 {
					const hex = "0123456789abcdef"
					if last < i {
						sb.WriteString(s[last:i])
					}
					var buf [6]byte
					buf[0] = '\\'
					buf[1] = 'u'
					buf[2] = '0'
					buf[3] = '0'
					buf[4] = hex[b>>4]
					buf[5] = hex[b&0x0f]
					sb.Write(buf[:])
					i++
					last = i
					continue
				}
				i++
				continue
			}
			if last < i {
				sb.WriteString(s[last:i])
			}
			sb.WriteString(esc)
			i++
			last = i
			continue
		}
		if b == 0xe2 && i+2 < len(s) && s[i+1] == 0x80 && (s[i+2] == 0xa8 || s[i+2] == 0xa9) {
			if last < i {
				sb.WriteString(s[last:i])
			}
			if s[i+2] == 0xa8 {
				sb.WriteString(`\u2028`)
			} else {
				sb.WriteString(`\u2029`)
			}
			i += 3
			last = i
			continue
		}
		i++
	}
	if last < len(s) {
		sb.WriteString(s[last:])
	}
	sb.WriteByte('"')
}

// skipValueStream consumes and syntax-validates one JSON value while updating
// path state identically to Gson's TypeAdapters.JSON_ELEMENT.read without
// allocating DOM nodes.
func (l *SourceMapJsonLexer) skipValueStream() jutil.Throwable {
	p, err := l.peekStream()
	if err != nil {
		return err
	}
	switch p {
	case smPeekBeginObject:
		_ = l.beginObjectStream()
	case smPeekBeginArray:
		_ = l.beginArrayStream()
	case smPeekDoubleQuoted, smPeekSingleQuoted:
		if _, err := l.nextStringStream(); err != nil {
			return err
		}
		return nil
	case smPeekUnquoted:
		_ = l.readUnquotedValue()
		l.peeked = smPeekNone
		l.incTopIndex()
		return nil
	case smPeekLong, smPeekNumber, smPeekTrue, smPeekFalse, smPeekNull:
		l.consumePeekedScalar()
		return nil
	default:
		return jutil.NewIllegalStateExceptionMsg("Unexpected token")
	}

	depth := 1
	for depth > 0 {
		for {
			more, err := l.hasNextStream()
			if err != nil {
				return err
			}
			if !more {
				break
			}
			if l.topScope() != smScopeEmptyArray && l.topScope() != smScopeNonEmptyArray {
				if _, err := l.nextNameStream(); err != nil {
					return err
				}
			}
			vp, err := l.peekStream()
			if err != nil {
				return err
			}
			switch vp {
			case smPeekBeginObject:
				_ = l.beginObjectStream()
				depth++
			case smPeekBeginArray:
				_ = l.beginArrayStream()
				depth++
			case smPeekDoubleQuoted, smPeekSingleQuoted:
				if _, err := l.nextStringStream(); err != nil {
					return err
				}
			case smPeekUnquoted:
				_ = l.readUnquotedValue()
				l.peeked = smPeekNone
				l.incTopIndex()
			case smPeekLong, smPeekNumber, smPeekTrue, smPeekFalse, smPeekNull:
				l.consumePeekedScalar()
			default:
				return jutil.NewIllegalStateExceptionMsg("Unexpected token")
			}
		}
		if l.topScope() == smScopeEmptyArray || l.topScope() == smScopeNonEmptyArray {
			if err := l.endArrayStream(); err != nil {
				return err
			}
		} else {
			if err := l.endObjectStream(); err != nil {
				return err
			}
		}
		depth--
	}
	return nil
}

// readJsonElement reads one Gson JsonElement from the stream (used for "x_..."
// extensions and coerced/invalid field values).
func (l *SourceMapJsonLexer) readJsonElement() (jutil.JsonElement, jutil.Throwable) {
	p, err := l.peekStream()
	if err != nil {
		return nil, err
	}
	readScalar := func(tok uint8) (jutil.JsonElement, jutil.Throwable) {
		switch tok {
		case smPeekDoubleQuoted, smPeekSingleQuoted, smPeekUnquoted:
			s, err := l.nextStringStream()
			if err != nil {
				return nil, err
			}
			return jutil.NewJsonPrimitiveString(s), nil
		case smPeekLong, smPeekNumber:
			s, err := l.nextStringStream()
			if err != nil {
				return nil, err
			}
			return jutil.NewJsonPrimitiveNumber(jutil.NewLazilyParsedNumber(s)), nil
		case smPeekTrue:
			l.consumePeekedScalar()
			return jutil.NewJsonPrimitiveBoolean(true), nil
		case smPeekFalse:
			l.consumePeekedScalar()
			return jutil.NewJsonPrimitiveBoolean(false), nil
		case smPeekNull:
			l.consumePeekedScalar()
			return jutil.JsonNull_INSTANCE, nil
		default:
			return nil, jutil.NewIllegalStateExceptionMsg("Unexpected token")
		}
	}

	var current jutil.JsonElement
	switch p {
	case smPeekBeginArray:
		_ = l.beginArrayStream()
		current = jutil.NewJsonArray()
	case smPeekBeginObject:
		_ = l.beginObjectStream()
		current = jutil.NewJsonObject()
	default:
		return readScalar(p)
	}

	var stack []jutil.JsonElement
	for {
		for {
			more, err := l.hasNextStream()
			if err != nil {
				return nil, err
			}
			if !more {
				break
			}
			var name string
			_, isObj := current.(*jutil.JsonObject)
			if isObj {
				name, err = l.nextNameStream()
				if err != nil {
					return nil, err
				}
			}
			vp, err := l.peekStream()
			if err != nil {
				return nil, err
			}
			var value jutil.JsonElement
			isNesting := false
			switch vp {
			case smPeekBeginArray:
				_ = l.beginArrayStream()
				value = jutil.NewJsonArray()
				isNesting = true
			case smPeekBeginObject:
				_ = l.beginObjectStream()
				value = jutil.NewJsonObject()
				isNesting = true
			default:
				value, err = readScalar(vp)
				if err != nil {
					return nil, err
				}
			}
			if isObj {
				current.(*jutil.JsonObject).Add(name, value)
			} else {
				current.(*jutil.JsonArray).Add(value)
			}
			if isNesting {
				stack = append(stack, current)
				current = value
			}
		}
		if _, isArr := current.(*jutil.JsonArray); isArr {
			if err := l.endArrayStream(); err != nil {
				return nil, err
			}
		} else {
			if err := l.endObjectStream(); err != nil {
				return nil, err
			}
		}
		if len(stack) == 0 {
			return current, nil
		}
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

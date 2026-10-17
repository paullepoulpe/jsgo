// Port of com/google/javascript/jscomp/resources/PropertiesParser.java @ closure-compiler v20261015

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

package resources

import (
	"regexp"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PropertiesParser is a simple parser for .properties files: one key and
// value per line, separated by ':', '=' or (failing both) the first space;
// lines starting with '#' or '!' are comments; a trailing backslash
// continues the value on the next line.
//
// Java: final class PropertiesParser (only static members).
type PropertiesParser struct{}

// propertiesParser_lineSeparator is the regex of source.split("\r?\n").
//
// Go-only: language-adapter — the java.util.regex pattern as an RE2 regexp
// (it has no backreferences or lookaround, so both engines agree).
var propertiesParser_lineSeparator = regexp.MustCompile("\r?\n")

// PropertiesParser_Parse parses source into an ordered key/value map. It
// throws IllegalArgumentException if a key occurs twice.
//
// Java: PropertiesParser.parse(String)
func PropertiesParser_Parse(source string) *jutil.LinkedHashMap[string, string] {
	builder := jutil.NewImmutableMapBuilder[string, string]()

	lines := jutil.StringSplitRegexp(source, propertiesParser_lineSeparator)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue // skip if empty or starts with # or !
		}

		delimeterIndex := PropertiesParser_FindDelimiter(line)
		if delimeterIndex == -1 {
			continue
		}

		data := ""
		// Remove whitespace on both sides of key.
		key := jutil.StringTrim(jutil.StringSubstring(line, 0, delimeterIndex))
		// Remove whitespace only on left side of data value. Trailing white space is data.
		line = PropertiesParser_TrimLeft(jutil.StringSubstringFrom(line, delimeterIndex+1))
		for {
			if strings.HasSuffix(line, "\\") {
				data += jutil.StringSubstring(line, 0, jutil.StringLength(line)-1)
				if i+1 == len(lines) {
					break
				}
				i++
				line = PropertiesParser_TrimLeft(lines[i])
			} else {
				data += line
				break
			}
		}

		builder.Put(key, data)
	}

	return builder.BuildOrThrow()
}

// PropertiesParser_TrimLeft removes leading spaces (only ' ').
//
// Java: PropertiesParser.trimLeft(String)
func PropertiesParser_TrimLeft(str string) string {
	chars := jutil.StringToUTF16(str) // Java: str.charAt(i)
	for i := 0; i < len(chars); i++ {
		if chars[i] != ' ' {
			return jutil.StringSubstringFrom(str, i)
		}
	}
	return str
}

// PropertiesParser_FindDelimiter returns the index of the first ':' or '=',
// or else of the first space, or -1.
//
// Java: PropertiesParser.findDelimiter(String)
func PropertiesParser_FindDelimiter(line string) int {
	chars := jutil.StringToUTF16(line) // Java: line.charAt(i)
	for i := 0; i < len(chars); i++ {
		switch chars[i] {
		case ':', '=':
			return i
		default:
		}
	}

	// If no : or =, delimiter is first whitespace.
	return jutil.StringIndexOf(line, " ")
}

// NewPropertiesParser is the private constructor.
//
// Java: PropertiesParser()
func NewPropertiesParser() *PropertiesParser {
	return &PropertiesParser{}
}

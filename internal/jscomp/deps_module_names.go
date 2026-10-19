// Port of com/google/javascript/jscomp/deps/ModuleNames.java @ closure-compiler v20261015

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleNames holds static methods related to module names.
type ModuleNames struct{}

// ModuleNames_MODULE_SLASH is the delimiter: according to the spec, the
// forward slash should be the delimiter on all platforms.
const ModuleNames_MODULE_SLASH = "/"

// ModuleNames_FileToModuleName returns a module name for an absolute path,
// with no resolution or checking.
//
// Java: ModuleNames.fileToModuleName(String)
func ModuleNames_FileToModuleName(path string) string {
	return ModuleNames_ToModuleName(ModuleNames_EscapePath(path))
}

// ModuleNames_FileToJsIdentifier returns a module name for an absolute path,
// with no resolution or checking.
//
// Java: ModuleNames.fileToJsIdentifier(String)
func ModuleNames_FileToJsIdentifier(path string) string {
	return ModuleNames_ToJSIdentifier(ModuleNames_EscapePath(path))
}

// ModuleNames_EscapePath escapes the given input path.
//
// Java: ModuleNames.escapePath(String)
func ModuleNames_EscapePath(input string) string {
	// Handle special characters
	encodedInput := strings.ReplaceAll(input, ":", "-")
	encodedInput = strings.ReplaceAll(encodedInput, "\\", "/")
	encodedInput = strings.ReplaceAll(encodedInput, " ", "%20")
	encodedInput = strings.ReplaceAll(encodedInput, "[", "%5B")
	encodedInput = strings.ReplaceAll(encodedInput, "]", "%5D")
	encodedInput = strings.ReplaceAll(encodedInput, "<", "%3C")
	encodedInput = strings.ReplaceAll(encodedInput, ">", "%3E")

	return ModuleNames_CanonicalizePath(encodedInput)
}

// Java: ModuleNames.toJSIdentifier(String)
func ModuleNames_ToJSIdentifier(path string) string {
	s := moduleNames_stripJsExtension(path)
	// replaceFirst("^\\./", "")
	s = strings.TrimPrefix(s, "./")
	var b strings.Builder
	last := 0
	for i := 0; i < len(s); i++ {
		var replacement byte
		width := 1
		switch s[i] {
		case '/', '\\', '@', '+', '!':
			replacement = '$'
		case '-', ':', '.':
			replacement = '_'
		case '%':
			if i+2 >= len(s) || s[i+1] != '2' || s[i+2] != '0' {
				continue
			}
			replacement = '_'
			width = 3
		default:
			continue
		}
		if last == 0 {
			b.Grow(len(s))
		}
		b.WriteString(s[last:i])
		b.WriteByte(replacement)
		i += width - 1
		last = i + 1
	}
	if last == 0 {
		return s
	}
	b.WriteString(s[last:])
	return b.String()
}

// Java: ModuleNames.toModuleName(String)
func ModuleNames_ToModuleName(path string) string {
	if strings.HasPrefix(path, ModuleNames_MODULE_SLASH) {
		path = jutil.StringSubstringFrom(path, 1)
	}
	return "module$" + ModuleNames_ToJSIdentifier(path)
}

func moduleNames_stripJsExtension(fileName string) string {
	if strings.HasSuffix(fileName, ".js") {
		fileName = fileName[:len(fileName)-len(".js")]
		// internal intermediary file format ".js.i.js"
		if strings.HasSuffix(fileName, ".js.i") {
			fileName = fileName[:len(fileName)-len(".js.i")]
		}
	}
	return fileName
}

// ModuleNames_CanonicalizePath canonicalizes a given path, removing
// segments containing "." and consuming segments for "..".
//
// If no segment could be consumed for "..", retains the segment.
//
// Java: ModuleNames.canonicalizePath(String)
func ModuleNames_CanonicalizePath(path string) string {
	if moduleNames_pathIsCanonical(path) {
		return path
	}
	parts := jutil.StringSplitLiteral(path, ModuleNames_MODULE_SLASH)
	buffer := make([]string, len(parts))
	position := 0
	available := 0

	absolutePath := len(parts) > 1 && parts[0] == ""
	if absolutePath {
		// If the path starts with "/" (so the left side, index zero, is
		// empty), then the path will always remain absolute. Make the first
		// segment unavailable to touch.
		available--
	}

	for _, part := range parts {
		if part == "." {
			continue
		}

		if part == ".." {
			if available > 0 {
				// Consume the previous segment.
				position--
				available--
				buffer[position] = ""
			} else if !absolutePath {
				// If this is a relative path, retain "..", as it can't be
				// consumed on the left.
				buffer[position] = part
				position++
			}
			continue
		}

		buffer[position] = part
		position++
		available++
	}

	if absolutePath && position == 1 {
		return ModuleNames_MODULE_SLASH // special-case single absolute segment as joining [""] doesn't work
	}
	return strings.Join(buffer[:position], ModuleNames_MODULE_SLASH)
}

func moduleNames_pathIsCanonical(path string) bool {
	if strings.HasSuffix(path, ModuleNames_MODULE_SLASH) {
		return false
	}
	for i := 0; i < len(path); i++ {
		if path[i] != '.' || (i != 0 && path[i-1] != '/') {
			continue
		}
		end := i + 1
		if end < len(path) && path[end] == '.' {
			end++
		}
		if end == len(path) || path[end] == '/' {
			return false
		}
	}
	return true
}

// Java: ModuleNames()
func newModuleNames() *ModuleNames { return &ModuleNames{} } // Static methods only; do not instantiate.

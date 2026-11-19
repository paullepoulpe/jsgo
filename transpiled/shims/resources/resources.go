/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Package resources embeds Closure Compiler v20261104 classpath resources
// (externs.zip, runtime_libs.typedast, runtime JS polyfills, ParserConfig.properties,
// and lib/*.js) and provides Class.getResourceAsStream / ResourceLoader / ResourceBundle shims.
package resources

import (
	"embed"
	"fmt"
	"io/fs"
	slashpath "path"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// JsDir is the embedded directory root for Closure JS runtime libraries and tables.
const JsDir = "com/google/javascript/jscomp/js/"

//go:embed com lib externs.zip build-data.properties
var classpath embed.FS

//go:embed runtime_libs.typedast
var runtimeLibsTypedAst []byte

//go:embed externs.zip
var externsZip []byte

const (
	runtimeLibsTypedAstName = "runtime_libs.typedast"
	externsZipName          = "externs.zip"
)

var (
	overridesMu    sync.RWMutex
	overridesMap   map[string][]byte
	overridesEpoch atomic.Uint64

	bundlesMu sync.RWMutex
	bundles   = make(map[string]*ResourceBundle)
)

// ExternsZipBytes returns the raw bytes of the embedded /externs.zip archive.
func ExternsZipBytes() []byte {
	return externsZip
}

// RuntimeLibsTypedAstBytes returns the raw bytes of the embedded /runtime_libs.typedast archive.
func RuntimeLibsTypedAstBytes() []byte {
	return runtimeLibsTypedAst
}

// RegisterOverride registers an in-memory resource override at path (with any
// leading '/' stripped), taking precedence over embedded resources.
func RegisterOverride(path string, data []byte) {
	name := strings.TrimPrefix(path, "/")
	cloned := slices.Clone(data)
	if cloned == nil {
		cloned = []byte{}
	}
	overridesMu.Lock()
	if overridesMap == nil {
		overridesMap = make(map[string][]byte)
	}
	overridesMap[name] = cloned
	overridesEpoch.Add(1)
	overridesMu.Unlock()
}

// UnregisterOverride removes any in-memory resource override registered at path.
func UnregisterOverride(path string) {
	name := strings.TrimPrefix(path, "/")
	overridesMu.Lock()
	if len(overridesMap) > 0 {
		delete(overridesMap, name)
		overridesEpoch.Add(1)
	}
	overridesMu.Unlock()
}

// ClearOverrides removes all registered in-memory resource overrides.
func ClearOverrides() {
	overridesMu.Lock()
	if len(overridesMap) > 0 {
		clear(overridesMap)
		overridesEpoch.Add(1)
	}
	overridesMu.Unlock()
}

func resolveClasspathName(basePackage, resourcePath string) string {
	if strings.HasPrefix(resourcePath, "/") {
		cleaned := slashpath.Clean(strings.TrimPrefix(resourcePath, "/"))
		if cleaned == "." {
			return ""
		}
		return cleaned
	}
	pkgDir := strings.ReplaceAll(strings.Trim(basePackage, "."), ".", "/")
	joined := resourcePath
	if pkgDir != "" {
		joined = pkgDir + "/" + resourcePath
	}
	cleaned := slashpath.Clean(joined)
	if cleaned == "." {
		return ""
	}
	return cleaned
}

// GetResourceBytes resolves resourcePath against basePackage (following Java Class.getResource rules)
// and returns the resource bytes if found.
func GetResourceBytes(basePackage, resourcePath string) ([]byte, bool) {
	return getResource(resolveClasspathName(basePackage, resourcePath))
}

// GetResourceAsStream implements Java Class.getResourceAsStream(path) for a class in basePackage.
// Returns nil when the resource does not exist.
func GetResourceAsStream(basePackage, resourcePath string) javaio.InputStream {
	data, ok := GetResourceBytes(basePackage, resourcePath)
	if !ok {
		return nil
	}
	return javaio.NewByteArrayInputStream(data)
}

// ResourceExists implements Java Class.getResource(path) != null / ResourceLoader.resourceExists.
func ResourceExists(basePackage, resourcePath string) bool {
	_, ok := GetResourceBytes(basePackage, resourcePath)
	return ok
}

// LoadTextResource implements Closure's ResourceLoader.loadTextResource(Class<?>, String).
func LoadTextResource(basePackage, resourcePath string) string {
	data, ok := GetResourceBytes(basePackage, resourcePath)
	if !ok {
		panic(javalang.NewRuntimeExceptionMsg("No such resource: " + resourcePath))
	}
	return string(data)
}

// LoadGlobalResource loads a UTF-8 resource at resourcePath (relative to basePackage or absolute with '/'),
// panicking with RuntimeException if not found.
func LoadGlobalResource(basePackage, resourcePath string) string {
	return LoadTextResource(basePackage, resourcePath)
}

// LoadBytes returns the bytes of the embedded resource at path (root-relative), or an error wrapping fs.ErrNotExist.
func LoadBytes(path string) ([]byte, error) {
	data, ok := getResource(strings.TrimPrefix(path, "/"))
	if !ok {
		return nil, fmt.Errorf("no such resource %q: %w", path, fs.ErrNotExist)
	}
	return data, nil
}

func getResource(name string) ([]byte, bool) {
	if overridesEpoch.Load() != 0 {
		overridesMu.RLock()
		data, ok := overridesMap[name]
		overridesMu.RUnlock()
		if ok {
			return data, true
		}
	}
	if name == runtimeLibsTypedAstName {
		return runtimeLibsTypedAst, true
	}
	if name == externsZipName {
		return externsZip, true
	}
	if data, err := fs.ReadFile(classpath, name); err == nil {
		return data, true
	}
	dir := strings.TrimSuffix(name, "/")
	if dir != "" {
		if info, err := fs.Stat(classpath, dir); err == nil && info.IsDir() {
			return []byte{}, true
		}
	}
	return nil, false
}

// ---------------------------------------------------------------------------
// java.util.ResourceBundle shim (used by ParserConfiguration.java)
// ---------------------------------------------------------------------------

// ResourceBundle provides read-only key/value access to an embedded .properties bundle.
type ResourceBundle struct {
	baseName string
	entries  map[string]string
}

// ResourceBundle_GetBundle loads the embedded "<baseName>.properties" resource
// (converting '.' in baseName to '/'). Panics with RuntimeException if missing.
func ResourceBundle_GetBundle(baseName string) *ResourceBundle {
	bundlesMu.RLock()
	b, ok := bundles[baseName]
	bundlesMu.RUnlock()
	if ok {
		return b
	}

	propPath := strings.ReplaceAll(baseName, ".", "/") + ".properties"
	data, found := getResource(propPath)
	if !found {
		panic(javalang.NewRuntimeExceptionMsg("Can't find bundle for base name " + baseName))
	}
	parsed := ParseJavaProperties(string(data))
	b = &ResourceBundle{
		baseName: baseName,
		entries:  parsed,
	}

	bundlesMu.Lock()
	bundles[baseName] = b
	bundlesMu.Unlock()
	return b
}

// GetString returns the string for the given key from this resource bundle,
// panicking if the key is missing (Java MissingResourceException extends RuntimeException).
func (b *ResourceBundle) GetString(key string) string {
	val, ok := b.entries[key]
	if !ok {
		panic(javalang.NewRuntimeExceptionMsg("Can't find resource for bundle " + b.baseName + ", key " + key))
	}
	return val
}

// ContainsKey reports whether key is contained in this ResourceBundle.
func (b *ResourceBundle) ContainsKey(key string) bool {
	_, ok := b.entries[key]
	return ok
}

// ParseJavaProperties parses a standard Java .properties text (including line continuations and \uXXXX escapes).
func ParseJavaProperties(text string) map[string]string {
	result := make(map[string]string)
	pos := 0
lines:
	for pos < len(text) {
		start, end, next := propertiesNaturalLine(text, pos)
		pos = next
		i := propertiesSkipSpace(text, start, end)
		if i == end || text[i] == '#' || text[i] == '!' {
			continue
		}
		var logical string
		if !propertiesContinues(text, i, end) {
			logical = text[i:end]
		} else {
			var sb strings.Builder
			for {
				if !propertiesContinues(text, i, end) {
					sb.WriteString(text[i:end])
					break
				}
				sb.WriteString(text[i : end-1])
				if pos >= len(text) {
					if sb.Len() == 0 && next-end == 2 {
						continue lines
					}
					break
				}
				if sb.Len() == 0 {
					continue lines
				}
				start, end, next = propertiesNaturalLine(text, pos)
				pos = next
				i = propertiesSkipSpace(text, start, end)
				if i == end {
					break
				}
			}
			logical = sb.String()
		}
		key, value := propertiesSplitLine(logical)
		result[propertiesUnescape(key)] = propertiesUnescape(value)
	}
	return result
}

func propertiesNaturalLine(in string, pos int) (start, end, next int) {
	end = pos
	for end < len(in) && in[end] != '\n' && in[end] != '\r' {
		end++
	}
	next = end
	if next < len(in) {
		if in[next] == '\r' && next+1 < len(in) && in[next+1] == '\n' {
			next++
		}
		next++
	}
	return pos, end, next
}

func propertiesIsSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\f' }

func propertiesSkipSpace(in string, i, end int) int {
	for i < end && propertiesIsSpace(in[i]) {
		i++
	}
	return i
}

func propertiesContinues(in string, i, end int) bool {
	n := 0
	for j := end - 1; j >= i && in[j] == '\\'; j-- {
		n++
	}
	return n%2 == 1
}

func propertiesSplitLine(line string) (key, value string) {
	k := 0
	for k < len(line) {
		c := line[k]
		if c == '\\' {
			k++
			if k < len(line) {
				_, w := utf8.DecodeRuneInString(line[k:])
				k += w
			}
			continue
		}
		if c == '=' || c == ':' || propertiesIsSpace(c) {
			break
		}
		k++
	}
	if k > len(line) {
		k = len(line)
	}
	v := propertiesSkipSpace(line, k, len(line))
	if v < len(line) && (line[v] == '=' || line[v] == ':') {
		v = propertiesSkipSpace(line, v+1, len(line))
	}
	return line[:k], line[v:]
}

func propertiesUnescape(s string) string {
	if strings.IndexByte(s, '\\') < 0 {
		return s
	}
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			out.WriteByte(c)
			continue
		}
		i++
		if i == len(s) {
			break
		}
		switch c = s[i]; c {
		case 't':
			out.WriteByte('\t')
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 'f':
			out.WriteByte('\f')
		case 'u':
			v := propertiesParseHex4(s, i+1)
			i += 4
			if v >= 0xD800 && v <= 0xDBFF && i+6 < len(s) && s[i+1] == '\\' && s[i+2] == 'u' {
				v2 := propertiesParseHex4(s, i+3)
				if v2 >= 0xDC00 && v2 <= 0xDFFF {
					out.WriteRune(utf16.DecodeRune(rune(v), rune(v2)))
					i += 6
					continue
				}
			}
			if v >= 0xD800 && v <= 0xDFFF {
				out.WriteRune(utf8.RuneError)
			} else {
				out.WriteRune(rune(v))
			}
		default:
			if c < 0x80 {
				out.WriteByte(c)
			} else {
				r, w := utf8.DecodeRuneInString(s[i:])
				out.WriteRune(r)
				i += w - 1
			}
		}
	}
	return out.String()
}

func propertiesParseHex4(s string, pos int) uint16 {
	var v uint16
	for k := 0; k < 4; k++ {
		d := -1
		if pos+k < len(s) {
			d = propertiesHexDigit(s[pos+k])
		}
		if d < 0 {
			panic(javalang.NewIllegalArgumentException(`Malformed \uxxxx encoding.`))
		}
		v = v<<4 | uint16(d)
	}
	return v
}

func propertiesHexDigit(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}

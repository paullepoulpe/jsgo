// Port of com/google/javascript/jscomp/transpile/TranspileResult.java @ closure-compiler v20261021

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

package transpile

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TranspileResult is the result of transpiling a single file. Includes
// transpiled and original source, and an optional source map.
type TranspileResult struct {
	path       *jutil.URI
	original   string
	transpiled string
	sourceMap  string
}

// Java: TranspileResult(URI, String, String, String)
func NewTranspileResult(path *jutil.URI, original string, transpiled string, sourceMap string) *TranspileResult {
	return &TranspileResult{
		path:       jutil.CheckNotNull(path),
		original:   original,
		transpiled: transpiled,
		sourceMap:  sourceMap,
	}
}

// Java: TranspileResult.path()
func (r *TranspileResult) Path() *jutil.URI { return r.path }

// TODO(sdh): how to specify naming scheme for sourcemap and original?
//  - if we're not embedding then we need to know how to find them.
//  - probabaly we'll keep 'path' pointing to the original, since that's
//    what's more important and harder to find - we can add some suffix
//    on for the transpiled version as necessary.

// Java: TranspileResult.original()
func (r *TranspileResult) Original() string { return r.original }

// Java: TranspileResult.sourceMap()
func (r *TranspileResult) SourceMap() string { return r.sourceMap }

// TODO(sdh): might be a bit nicer for this to be a suffix rather than a
// whole path?

// Java: TranspileResult.embedSourcemapUrl(String)
func (r *TranspileResult) EmbedSourcemapUrl(url string) *TranspileResult {
	// If there's no sourcemap, don't reference it.
	if r.sourceMap == "" {
		return r
	}
	embedded := r.transpiled + "\n//# sourceMappingURL=" + url + "\n"
	return NewTranspileResult(r.path, r.original, embedded, r.sourceMap)
}

const sourceMapPercentHex = "0123456789ABCDEF"

func isSafeSourceMapPercentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '*'
}

// percentEscapeSourceMap implements PercentEscaper("-_.*", false): ASCII
// alphanumerics and "-_.*" are left unescaped; all other code points (including
// ':' and ' ') are UTF-8 percent-encoded with uppercase hex digits.
func percentEscapeSourceMap(s string) string {
	firstUnsafe := -1
	for i := 0; i < len(s); i++ {
		if !isSafeSourceMapPercentByte(s[i]) {
			firstUnsafe = i
			break
		}
	}
	if firstUnsafe < 0 {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + (len(s)-firstUnsafe)*2)
	sb.WriteString(s[:firstUnsafe])

	for i := firstUnsafe; i < len(s); {
		b := s[i]
		if isSafeSourceMapPercentByte(b) {
			start := i
			i++
			for i < len(s) && isSafeSourceMapPercentByte(s[i]) {
				i++
			}
			sb.WriteString(s[start:i])
			continue
		}
		if b < 0x80 {
			sb.WriteByte('%')
			sb.WriteByte(sourceMapPercentHex[b>>4])
			sb.WriteByte(sourceMapPercentHex[b&0x0F])
			i++
			continue
		}
		// Reject unpaired WTF-8 surrogates (0xED 0xA0..0xBF 0x80..0xBF) to match
		// Guava UnicodeEscaper.codePointAt.
		if b == 0xED && i+2 < len(s) && s[i+1] >= 0xA0 && s[i+1] <= 0xBF && s[i+2] >= 0x80 && s[i+2] <= 0xBF {
			panic(jutil.NewIllegalArgumentExceptionMsg("Unpaired surrogate in percent-escaped string"))
		}
		_, width := utf8.DecodeRuneInString(s[i:])
		for k := 0; k < width; k++ {
			oct := s[i+k]
			sb.WriteByte('%')
			sb.WriteByte(sourceMapPercentHex[oct>>4])
			sb.WriteByte(sourceMapPercentHex[oct&0x0F])
		}
		i += width
	}
	return sb.String()
}

// Java: TranspileResult.embedSourcemap()
func (r *TranspileResult) EmbedSourcemap() *TranspileResult {
	if r.sourceMap == "" {
		return r
	}
	embedded := r.transpiled + "\n//# sourceMappingURL=data:," + percentEscapeSourceMap(r.sourceMap) + "\n"
	return NewTranspileResult(r.path, r.original, embedded, "")
}

// Java: TranspileResult.embedSourcemapBase64()
func (r *TranspileResult) EmbedSourcemapBase64() *TranspileResult {
	if r.sourceMap == "" {
		return r
	}
	embedded := r.transpiled +
		"\n//# sourceMappingURL=data:application/json;base64," +
		base64.StdEncoding.EncodeToString(jutil.StringGetBytesUTF8(r.sourceMap)) +
		"\n"
	return NewTranspileResult(r.path, r.original, embedded, "")
}

// Java: TranspileResult.transpiled()
func (r *TranspileResult) Transpiled() string { return r.transpiled }

// Java: TranspileResult.wasTranspiled()
func (r *TranspileResult) WasTranspiled() bool { return r.transpiled != r.original }

// Java: TranspileResult.equals(Object)
func (r *TranspileResult) Equals(other any) bool {
	transpileResult, ok := other.(*TranspileResult)
	return ok && transpileResult != nil &&
		transpileResult.path.Equals(r.path) &&
		transpileResult.original == r.original &&
		transpileResult.transpiled == r.transpiled &&
		transpileResult.sourceMap == r.sourceMap
}

// Java: TranspileResult.hashCode() (Objects.hash)
func (r *TranspileResult) HashCode() int32 {
	h := int32(1)
	h = 31*h + r.path.HashCode()
	h = 31*h + jutil.StringHashCode(r.original)
	h = 31*h + jutil.StringHashCode(r.transpiled)
	h = 31*h + jutil.StringHashCode(r.sourceMap)
	return h
}

// Java: TranspileResult.toString()
func (r *TranspileResult) String() string {
	return jutil.StringFormat(
		"TranspileResut{path=%s, original=%s, transpiled=%s, sourceMapURL=%s}",
		r.path, r.original, r.transpiled, r.sourceMap)
}

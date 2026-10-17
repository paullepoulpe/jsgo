// Port of com/google/javascript/jscomp/transpile/TranspileResult.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

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

// sourcemaps must escape :, and cannot escape space as +, so we need a
// custom escaper.
var transpileResult_ESCAPER jutil.Escaper = jutil.NewPercentEscaper("-_.*" /* plusForSpace= */, false)

// Java: TranspileResult.embedSourcemap()
func (r *TranspileResult) EmbedSourcemap() *TranspileResult {
	if r.sourceMap == "" {
		return r
	}
	embedded := r.transpiled + "\n//# sourceMappingURL=data:," + transpileResult_ESCAPER.Escape(r.sourceMap) + "\n"
	return NewTranspileResult(r.path, r.original, embedded, "")
}

// Java: TranspileResult.embedSourcemapBase64()
func (r *TranspileResult) EmbedSourcemapBase64() *TranspileResult {
	if r.sourceMap == "" {
		return r
	}
	embedded := r.transpiled +
		"\n//# sourceMappingURL=data:application/json;base64," +
		jutil.BaseEncoding_Base64().Encode(jutil.StringGetBytesUTF8(r.sourceMap)) +
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

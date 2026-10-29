// Port of com/google/javascript/jscomp/SourceMapResolver.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SourceMapResolver is a utility class for resolving source maps and files
// referenced in source maps.
type SourceMapResolver struct{}

const sourceMapResolver_BASE64_URL_PREFIX = "data:"
const sourceMapResolver_BASE64_START = "base64,"

// For now only accept UTF-8. We could use the actual charset information in
// the future. Default is UTF-8 if unspecified.
func sourceMapResolver_isAcceptedMediaType(mediaType string) bool {
	return mediaType == "application/json;charset=utf-8;" || mediaType == "application/json;"
}

// SourceMapResolver_ExtractSourceMap locates the appropriate sourcemap on
// disk for a given //# sourceMappingUrl. This is use for sourcemap merging
// (--apply_input_source_maps) and for error resolution. parseInlineSourceMaps
// says whether to parse Base64 encoded source maps. It may return nil.
//
// Java: SourceMapResolver.extractSourceMap(SourceFile, String, boolean)
func SourceMapResolver_ExtractSourceMap(jsFile *SourceFile, sourceMapURL string, parseInlineSourceMaps bool) *SourceFile {
	if parseInlineSourceMaps && strings.HasPrefix(sourceMapURL, sourceMapResolver_BASE64_URL_PREFIX) {
		if rawBase64, canonicalUTF8, ok := sourceMapResolver_extractRawBase64(sourceMapURL); ok {
			return newInlineBase64SourceFile(
				jsFile.GetName()+".inline.map", rawBase64, canonicalUTF8)
		}
		return nil
	}
	// TODO(tdeegan): Handle absolute urls here.  The compiler needs to come up with a scheme for
	// properly resolving absolute urls from http:// or the root /some/abs/path/... See b/62544959.
	if sourceMapResolver_isAbsolute(sourceMapURL) {
		return nil
	}
	// If not absolute, its relative.
	// TODO(tdeegan): Handle urls relative to //# sourceURL. See the sourcemap spec.
	return SourceMapResolver_GetRelativePath(jsFile.GetName(), sourceMapURL)
}

func sourceMapResolver_extractRawBase64(url string) (rawBase64 string, canonicalUTF8 bool, ok bool) {
	if strings.HasPrefix(url, sourceMapResolver_BASE64_URL_PREFIX) {
		base64StartIndex := strings.Index(url, sourceMapResolver_BASE64_START)
		if base64StartIndex >= 0 {
			mediaType := url[len(sourceMapResolver_BASE64_URL_PREFIX):base64StartIndex]
			if sourceMapResolver_isAcceptedMediaType(mediaType) {
				raw := url[base64StartIndex+len(sourceMapResolver_BASE64_START):]
				canonical := jutil.BaseEncoding_Base64().ValidateAndCheckCanonicalUTF8(raw)
				return strings.Clone(raw), canonical, true
			}
		}
	}
	return "", false, false
}

// sourceMapResolver_extractBase64String is based on
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URIs.
// It returns nil when url is not an accepted data URL.
//
// Java: SourceMapResolver.extractBase64String(String)
func sourceMapResolver_extractBase64String(url string) *string {
	if strings.HasPrefix(url, sourceMapResolver_BASE64_URL_PREFIX) && strings.Contains(url, sourceMapResolver_BASE64_START) {
		base64StartIndex := strings.Index(url, sourceMapResolver_BASE64_START)
		mediaType := url[len(sourceMapResolver_BASE64_URL_PREFIX):base64StartIndex]
		if sourceMapResolver_isAcceptedMediaType(mediaType) {
			data := jutil.BaseEncoding_Base64().Decode(
				url[base64StartIndex+len(sourceMapResolver_BASE64_START):])
			s := jutil.NewStringFromUTF8Bytes(data)
			return &s
		}
	}
	return nil
}

// Java: SourceMapResolver.isAbsolute(String)
func sourceMapResolver_isAbsolute(url string) bool {
	u, err := jutil.ParseURI(url)
	if err != nil {
		cause, _ := err.(jutil.Throwable)
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("Sourcemap url was invalid: "+url, cause))
	}
	return u.IsAbsolute() || strings.HasPrefix(url, "/")
}

// SourceMapResolver_GetRelativePath returns the relative path, resolved
// relative to the base path, where the base path is interpreted as a
// filename rather than a directory. E.g.: getRelativeTo("../foo/bar.js",
// "baz/bam/qux.js") --> "baz/foo/bar.js"
//
// Java: SourceMapResolver.getRelativePath(String, String)
func SourceMapResolver_GetRelativePath(baseFilePath string, relativePath string) *SourceFile {
	return SourceFile_Builder_().
		WithPathFromPath(
			jutil.Path_Of(baseFilePath).ResolveSiblingString(relativePath).Normalize()).
		WithKind(rhino.StaticSourceFile_SourceKind_NON_CODE).
		Build()
}

// SourceMapResolver_AddBase64PrefixToEncodedSourceMap returns the encoded
// source map with the base64 prefix attached. (i.e. prefix with
// "data:application/json;base64,")
//
// E.g. turns "eyJ2ZXJzaW9uI..." into "data:application/json;base64,eyJ2ZXJzaW9uI..."
//
// Java: SourceMapResolver.addBase64PrefixToEncodedSourceMap(String)
func SourceMapResolver_AddBase64PrefixToEncodedSourceMap(encoded string) string {
	return sourceMapResolver_BASE64_URL_PREFIX + "application/json;" + sourceMapResolver_BASE64_START + encoded
}

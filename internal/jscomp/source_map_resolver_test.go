// Port of test/com/google/javascript/jscomp/SourceMapResolverTest.java @ closure-compiler v20261021

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
	"encoding/base64"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestSourceMapResolver_ResolveBase64Inline(t *testing.T) {
	sourceMap := "{map: 'asdfasdf'}"
	encoded := base64.StdEncoding.EncodeToString([]byte(sourceMap))
	url := "data:application/json;base64," + encoded
	code := "console.log('asdf')\n//# sourceMappingURL=" + url
	s := SourceMapResolver_ExtractSourceMap(SourceFile_FromCode("somePath/hello.js", code), url, true)
	jsaEq(t, s.GetCode(), sourceMap)
	jsaEq(t, s.GetName(), "somePath/hello.js.inline.map")

	// --parse_inline_source_maps=false
	noInline := SourceMapResolver_ExtractSourceMap(SourceFile_FromCode("somePath/hello.js", code), url, false)
	jsaTrue(t, noInline == nil)
}

func TestSourceMapResolver_ResolveBase64WithCharsetInline(t *testing.T) {
	sourceMap := "{map: 'asdfasdf'}"
	encoded := base64.StdEncoding.EncodeToString([]byte(sourceMap))
	url := "data:application/json;charset=utf-8;base64," + encoded
	code := "console.log('asdf')\n//# sourceMappingURL=" + url
	s := SourceMapResolver_ExtractSourceMap(SourceFile_FromCode("somePath/hello.js", code), url, true)
	jsaEq(t, s.GetCode(), sourceMap)
	jsaEq(t, s.GetName(), "somePath/hello.js.inline.map")

	// Try non supported charset.
	dataURLWithBadCharset := "data:application/json;charset=asdf;base64," + encoded
	charsetCode := "console.log('asdf')\n//# sourceMappingURL=" + dataURLWithBadCharset
	result := SourceMapResolver_ExtractSourceMap(
		SourceFile_FromCode("somePath/hello.js", charsetCode), dataURLWithBadCharset, true)
	jsaTrue(t, result == nil)
}

func TestSourceMapResolver_Absolute(t *testing.T) {
	jsFile := SourceFile_FromCode("somePath/hello.js", "console.log(1)")
	// We cannot reslove absolute urls.
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "/asdf/asdf.js", true) == nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "/asdf/.././asdf.js", true) == nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "http://google.com/asdf/asdf.js", true) == nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "https://google.com/asdf/asdf.js", true) == nil)

	// We can resolve relative urls
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "asdf.js", true) != nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "asdf/asdf.js", true) != nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "asdf/.././asdf.js", true) != nil)
	jsaTrue(t, SourceMapResolver_ExtractSourceMap(jsFile, "not/.././a/js/file.txt", true) != nil)
}

// Go-only: extra coverage — an invalid URL is wrapped in a RuntimeException.
func TestSourceMapResolver_InvalidUrl(t *testing.T) {
	jsFile := SourceFile_FromCode("somePath/hello.js", "console.log(1)")
	e := jsaThrows[jutil.IRuntimeException](t, func() { SourceMapResolver_ExtractSourceMap(jsFile, "a b", true) })
	if e != nil {
		jsaEq(t, e.Error(), "java.lang.RuntimeException: Sourcemap url was invalid: a b")
	}
}

func TestSourceMapResolver_RelativePaths(t *testing.T) {
	jsaEq(t, SourceMapResolver_GetRelativePath("basefile.js", "basefile.js.map").GetName(), "basefile.js.map")
	jsaEq(t, SourceMapResolver_GetRelativePath("path/basefile.js", "relative/path/basefile.js.map").GetName(),
		"path/relative/path/basefile.js.map")
	jsaEq(t, SourceMapResolver_GetRelativePath("some/longer/path/basefile.js", "../sourcemap.js.map").String(),
		"some/longer/sourcemap.js.map")
	jsaEq(t, SourceMapResolver_GetRelativePath("some/longer/path/basefile.js", ".././../sourcemap.js.map").GetName(),
		"some/sourcemap.js.map")
	jsaEq(t, SourceMapResolver_GetRelativePath("basefile.js", "../basefile.js.map").GetName(), "../basefile.js.map")
	jsaEq(t, SourceMapResolver_GetRelativePath("baz/bam/qux.js", "../foo/bar.js").GetName(), "baz/foo/bar.js")
}

func TestSourceMapResolver_Integration(t *testing.T) {
	url := "relative/path/to/sourcemap/hello.js.map"
	s := SourceMapResolver_ExtractSourceMap(SourceFile_FromCode("somePath/hello.js", ""), url, false)
	jsaEq(t, s.GetName(), "somePath/relative/path/to/sourcemap/hello.js.map")
}

func TestSourceMapResolver_AddBase64PrefixToEncodedSourceMap(t *testing.T) {
	base64Prefix := "data:application/json;base64,"
	encodedSourceMap := "eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZm9vLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiZm9vLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBO0lBR0UsV0FBWSxLQUFhO1FBQ3ZCLElBQUksQ0FBQyxDQUFDLEdBQUcsS0FBSyxDQUFDO0lBQ2pCLENBQUM7SUFDSCxRQUFDO0FBQUQsQ0FBQyxBQU5ELElBTUM7QUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMifQ=="
	base64EncodedSourceMap := SourceMapResolver_AddBase64PrefixToEncodedSourceMap(encodedSourceMap)
	jsaEq(t, base64EncodedSourceMap, base64Prefix+encodedSourceMap)
}

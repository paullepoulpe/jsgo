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
// Port of com/google/javascript/jscomp/SourceMapResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SourceMapResolver
type SourceMapResolver struct {
	idHash_ atomic.Int32
}

func (this_ *SourceMapResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapResolver) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMapResolver", "SourceMapResolver"
}
func (this_ *SourceMapResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMapResolver", "com.google.javascript.jscomp.SourceMapResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapResolver](v_)
	})
}

func StaticGet_SourceMapResolver_BASE64_URL_PREFIX() string { return "data:" }
func StaticGet_SourceMapResolver_BASE64_START() string      { return "base64," }

var static_SourceMapResolver_ACCEPTED_MEDIA_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_SourceMapResolver_ACCEPTED_MEDIA_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_SourceMapResolver()
	return static_SourceMapResolver_ACCEPTED_MEDIA_TYPES
}
func StaticSet_SourceMapResolver_ACCEPTED_MEDIA_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_SourceMapResolver()
	static_SourceMapResolver_ACCEPTED_MEDIA_TYPES = v
}

var clinitDone_SourceMapResolver uint32
var clinitState_SourceMapResolver j2gort.ClinitState

func clinit_SourceMapResolver() {
	if atomic.LoadUint32(&clinitDone_SourceMapResolver) != 0 {
		return
	}
	clinitSlow_SourceMapResolver()
}

func clinitSlow_SourceMapResolver() {
	if !j2gort.EnterClinit(&clinitState_SourceMapResolver, &clinitDone_SourceMapResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceMapResolver, &clinitDone_SourceMapResolver)
	static_SourceMapResolver_ACCEPTED_MEDIA_TYPES = j2gort.ImmSetOfStr2("application/json;charset=utf-8;", "application/json;")
}

func (this_ *SourceMapResolver) initFields_SourceMapResolver() {
	_ = this_
}

// Java constructor: SourceMapResolver()
func NewSourceMapResolver() *SourceMapResolver {
	this_ := &SourceMapResolver{}
	InitSourceMapResolver(this_)
	return this_
}

func InitSourceMapResolver(this_ *SourceMapResolver) {
	this_.initFields_SourceMapResolver()
}

// Java static method: SourceMapResolver.extractSourceMap(com.google.javascript.jscomp.SourceFile jsFile, java.lang.String sourceMapURL, boolean parseInlineSourceMaps)
func SourceMapResolver_ExtractSourceMap(jsFile *SourceFile, sourceMapURL string, parseInlineSourceMaps bool) *SourceFile {
	_ = jsFile
	_ = sourceMapURL
	_ = parseInlineSourceMaps
	if parseInlineSourceMaps && javalang.StringStartsWith(j2gort.CheckNotNullString(sourceMapURL), "data:") {
		var extractedString string = SourceMapResolver_ExtractBase64String(sourceMapURL)
		_ = extractedString
		if extractedString != j2gort.NullString {
			return SourceFile_FromCode_3((j2gort.StrValString(jsFile.GetName()) + ".inline.map"), extractedString, pkg_rhino.StaticGet_StaticSourceFile_SourceKind_NON_CODE())
		}
		return nil
	}
	if SourceMapResolver_IsAbsolute(sourceMapURL) {
		return nil
	}
	return SourceMapResolver_GetRelativePath(jsFile.GetName(), sourceMapURL)
}

// Java static method: SourceMapResolver.extractBase64String(java.lang.String url)
func SourceMapResolver_ExtractBase64String(url string) string {
	_ = url
	if javalang.StringStartsWith(j2gort.CheckNotNullString(url), "data:") && (j2gort.JStrIndexOf(url, "base64,") >= 0) {
		var base64StartIndex int32 = j2gort.JStrIndexOf(url, "base64,")
		_ = base64StartIndex
		var mediaType string = j2gort.JStrSub(url, j2gort.JStrLen("data:"), base64StartIndex)
		_ = mediaType
		if j2gort.SetContainsStr(StaticGet_SourceMapResolver_ACCEPTED_MEDIA_TYPES(), mediaType) {
			var data []int8 = j2gort.Ext_com_google_common_io_BaseEncoding_Base64().Decode(j2gort.BoxString(j2gort.JStrSubFrom(url, (base64StartIndex + j2gort.JStrLen("base64,")))))
			_ = data
			{
				b_1 := data
				return j2gort.NewStringFromBytes(b_1, 0, int32(len(b_1)), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
			}
		}
	}
	return j2gort.NullString
}

// Java static method: SourceMapResolver.isAbsolute(java.lang.String url)
func SourceMapResolver_IsAbsolute(url string) bool {
	_ = url
	var didRet_2 bool
	var retVal_3 bool
	func() {
		defer func() {
			if rec_4 := recover(); rec_4 != nil {
				jex_5 := javalang.RecoverException(rec_4)
				_ = jex_5
				if j2gort.InstanceOf[j2gort.Ext_java_net_URISyntaxException](jex_5) {
					e := j2gort.AsIface[j2gort.Ext_java_net_URISyntaxException](jex_5)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Sourcemap url was invalid: "+j2gort.StrValString(url), e)))
				} else {
					panic(rec_4)
				}
			}
		}()
		retVal_3 = (j2gort.NewExt_java_net_URI_1(url).IsAbsolute() || javalang.StringStartsWith(j2gort.CheckNotNullString(url), "/"))
		didRet_2 = true
		return
	}()
	if didRet_2 {
		return retVal_3
	}
	panic("unreachable")
}

// Java static method: SourceMapResolver.getRelativePath(java.lang.String baseFilePath, java.lang.String relativePath)
func SourceMapResolver_GetRelativePath(baseFilePath string, relativePath string) *SourceFile {
	_ = baseFilePath
	_ = relativePath
	return SourceFile_S_Builder().WithPath_1_Ext_java_nio_file_Path(j2gort.Ext_java_nio_file_FileSystems_GetDefault().GetPath(baseFilePath, []string{}).ResolveSibling_1_string(relativePath).Normalize()).WithKind(pkg_rhino.StaticGet_StaticSourceFile_SourceKind_NON_CODE()).Build()
}

// Java static method: SourceMapResolver.addBase64PrefixToEncodedSourceMap(java.lang.String encoded)
func SourceMapResolver_AddBase64PrefixToEncodedSourceMap(encoded string) string {
	_ = encoded
	return ("data:" + "application/json;" + "base64," + j2gort.StrValString(encoded))
}

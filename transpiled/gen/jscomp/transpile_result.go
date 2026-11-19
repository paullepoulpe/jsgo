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
// Port of com/google/javascript/jscomp/transpile/TranspileResult.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.transpile.TranspileResult
type TranspileResult struct {
	path       *j2gort.Ext_java_net_URI `gson:"path"`
	original   string                   `gson:"original"`
	transpiled string                   `gson:"transpiled"`
	sourceMap  string                   `gson:"sourceMap"`
}

func (this_ *TranspileResult) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspileResult) IsJavaNil_() bool { return this_ == nil }
func (*TranspileResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.TranspileResult", "TranspileResult"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.TranspileResult", "com.google.javascript.jscomp.transpile.TranspileResult", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspileResult](v_)
	})
}

var static_TranspileResult_ESCAPER j2gort.Ext_com_google_common_escape_Escaper

func StaticGet_TranspileResult_ESCAPER() j2gort.Ext_com_google_common_escape_Escaper {
	clinit_TranspileResult()
	return static_TranspileResult_ESCAPER
}
func StaticSet_TranspileResult_ESCAPER(v j2gort.Ext_com_google_common_escape_Escaper) {
	clinit_TranspileResult()
	static_TranspileResult_ESCAPER = v
}

var clinitDone_TranspileResult uint32
var clinitState_TranspileResult j2gort.ClinitState

func clinit_TranspileResult() {
	if atomic.LoadUint32(&clinitDone_TranspileResult) != 0 {
		return
	}
	clinitSlow_TranspileResult()
}

func clinitSlow_TranspileResult() {
	if !j2gort.EnterClinit(&clinitState_TranspileResult, &clinitDone_TranspileResult) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TranspileResult, &clinitDone_TranspileResult)
	static_TranspileResult_ESCAPER = j2gort.NewExt_com_google_common_net_PercentEscaper("-_.*", false)
}

func (this_ *TranspileResult) initFields_TranspileResult() {
	_ = this_
	this_.original = j2gort.NullString
	this_.transpiled = j2gort.NullString
	this_.sourceMap = j2gort.NullString
}

// Java constructor: TranspileResult(java.net.URI path, java.lang.String original, java.lang.String transpiled, java.lang.String sourceMap)
func NewTranspileResult(path *j2gort.Ext_java_net_URI, original string, transpiled string, sourceMap string) *TranspileResult {
	clinit_TranspileResult()
	this_ := &TranspileResult{}
	InitTranspileResult(this_, path, original, transpiled, sourceMap)
	return this_
}

func InitTranspileResult(this_ *TranspileResult, path *j2gort.Ext_java_net_URI, original string, transpiled string, sourceMap string) {
	_ = path
	_ = original
	_ = transpiled
	_ = sourceMap
	this_.initFields_TranspileResult()
	this_.path = j2gort.CheckNotNullPtr(path)
	this_.original = j2gort.CheckNotNullString(original)
	this_.transpiled = j2gort.CheckNotNullString(transpiled)
	this_.sourceMap = j2gort.CheckNotNullString(sourceMap)
}

// Java method: TranspileResult.path()
func (this_ *TranspileResult) Path() *j2gort.Ext_java_net_URI {
	return this_.path
}

// Java method: TranspileResult.original()
func (this_ *TranspileResult) Original() string {
	return this_.original
}

// Java method: TranspileResult.sourceMap()
func (this_ *TranspileResult) SourceMap() string {
	return this_.sourceMap
}

// Java method: TranspileResult.embedSourcemapUrl(java.lang.String url)
func (this_ *TranspileResult) EmbedSourcemapUrl(url string) *TranspileResult {
	_ = url
	if j2gort.CheckNotNullString(this_.sourceMap) == "" {
		return this_
	}
	var embedded string = (j2gort.StrValString(this_.transpiled) + "\n//# sourceMappingURL=" + j2gort.StrValString(url) + "\n")
	_ = embedded
	return NewTranspileResult(this_.path, this_.original, embedded, this_.sourceMap)
}

// Java method: TranspileResult.embedSourcemap()
func (this_ *TranspileResult) EmbedSourcemap() *TranspileResult {
	if j2gort.CheckNotNullString(this_.sourceMap) == "" {
		return this_
	}
	var embedded string = (j2gort.StrValString(this_.transpiled) + "\n//# sourceMappingURL=data:," + j2gort.StrValString(StaticGet_TranspileResult_ESCAPER().Escape_1_string(this_.sourceMap)) + "\n")
	_ = embedded
	return NewTranspileResult(this_.path, this_.original, embedded, "")
}

// Java method: TranspileResult.embedSourcemapBase64()
func (this_ *TranspileResult) EmbedSourcemapBase64() *TranspileResult {
	if j2gort.CheckNotNullString(this_.sourceMap) == "" {
		return this_
	}
	var embedded string = (j2gort.StrValString(this_.transpiled) + "\n//# sourceMappingURL=data:application/json;base64," + j2gort.StrValString(j2gort.Ext_com_google_common_io_BaseEncoding_Base64().Encode_1(func() []int8 {
		_ = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
		return j2gort.JStrGetBytes(this_.sourceMap)
	}())) + "\n")
	_ = embedded
	return NewTranspileResult(this_.path, this_.original, embedded, "")
}

// Java method: TranspileResult.transpiled()
func (this_ *TranspileResult) Transpiled() string {
	return this_.transpiled
}

// Java method: TranspileResult.wasTranspiled()
func (this_ *TranspileResult) WasTranspiled() bool {
	return !(func() bool {
		sr_1 := this_.transpiled
		se_2 := this_.original
		return j2gort.CheckNotNullString(sr_1) == se_2
	}())
}

// Java method: TranspileResult.equals(java.lang.Object other)
func (this_ *TranspileResult) Equals(arg_3 any) bool {
	other := arg_3
	_ = other
	var transpileResult *TranspileResult
	_ = transpileResult
	var ok_4 bool
	transpileResult, ok_4 = j2gort.TryCast[*TranspileResult](other)
	return ok_4 && (transpileResult.path.Equals(j2gort.PtrToAny(this_.path))) && (func() bool {
		sr_5 := transpileResult.original
		se_6 := this_.original
		return j2gort.CheckNotNullString(sr_5) == se_6
	}()) && (func() bool {
		sr_7 := transpileResult.transpiled
		se_8 := this_.transpiled
		return j2gort.CheckNotNullString(sr_7) == se_8
	}()) && (func() bool {
		sr_9 := transpileResult.sourceMap
		se_10 := this_.sourceMap
		return j2gort.CheckNotNullString(sr_9) == se_10
	}())
}

// Java method: TranspileResult.hashCode()
func (this_ *TranspileResult) HashCode() int32 {
	return func() int32 {
		var oh_11 int32 = 1
		oh_11 = 31*oh_11 + j2gort.ObjectHash(j2gort.PtrToAny(this_.path))
		{
			hs_12 := this_.original
			oh_11 = 31 * oh_11
			if hs_12 != j2gort.NullString {
				oh_11 += javalang.StringHashCode(hs_12)
			}
		}
		{
			hs_13 := this_.transpiled
			oh_11 = 31 * oh_11
			if hs_13 != j2gort.NullString {
				oh_11 += javalang.StringHashCode(hs_13)
			}
		}
		{
			hs_14 := this_.sourceMap
			oh_11 = 31 * oh_11
			if hs_14 != j2gort.NullString {
				oh_11 += javalang.StringHashCode(hs_14)
			}
		}
		return oh_11
	}()
}

// Java method: TranspileResult.toString()
func (this_ *TranspileResult) String() string {
	return j2gort.Ext_java_lang_String_Format_2("TranspileResut{path=%s, original=%s, transpiled=%s, sourceMapURL=%s}", []any{j2gort.PtrToAny(this_.path), j2gort.BoxString(this_.original), j2gort.BoxString(this_.transpiled), j2gort.BoxString(this_.sourceMap)})
}

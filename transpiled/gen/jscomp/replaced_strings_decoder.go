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
// Port of com/google/javascript/jscomp/ReplacedStringsDecoder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ReplacedStringsDecoder
type ReplacedStringsDecoder struct {
	originalToNewNameMap j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"originalToNewNameMap"`
	idHash_              atomic.Int32
}

func (this_ *ReplacedStringsDecoder) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplacedStringsDecoder) IsJavaNil_() bool { return this_ == nil }
func (*ReplacedStringsDecoder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplacedStringsDecoder", "ReplacedStringsDecoder"
}
func (this_ *ReplacedStringsDecoder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplacedStringsDecoder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplacedStringsDecoder", "com.google.javascript.jscomp.ReplacedStringsDecoder", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplacedStringsDecoder](v_)
	})
}

func StaticGet_ReplacedStringsDecoder_ARGUMENT_PLACE_HOLDER() string { return "`" }

var static_ReplacedStringsDecoder_NULL_DECODER *ReplacedStringsDecoder

func StaticGet_ReplacedStringsDecoder_NULL_DECODER() *ReplacedStringsDecoder {
	clinit_ReplacedStringsDecoder()
	return static_ReplacedStringsDecoder_NULL_DECODER
}
func StaticSet_ReplacedStringsDecoder_NULL_DECODER(v *ReplacedStringsDecoder) {
	clinit_ReplacedStringsDecoder()
	static_ReplacedStringsDecoder_NULL_DECODER = v
}

var clinitDone_ReplacedStringsDecoder uint32
var clinitState_ReplacedStringsDecoder j2gort.ClinitState

func clinit_ReplacedStringsDecoder() {
	if atomic.LoadUint32(&clinitDone_ReplacedStringsDecoder) != 0 {
		return
	}
	clinitSlow_ReplacedStringsDecoder()
}

func clinitSlow_ReplacedStringsDecoder() {
	if !j2gort.EnterClinit(&clinitState_ReplacedStringsDecoder, &clinitDone_ReplacedStringsDecoder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplacedStringsDecoder, &clinitDone_ReplacedStringsDecoder)
	static_ReplacedStringsDecoder_NULL_DECODER = NewReplacedStringsDecoder(VariableMap_FromMap(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())))
}

func (this_ *ReplacedStringsDecoder) initFields_ReplacedStringsDecoder() {
	_ = this_
}

// Java constructor: ReplacedStringsDecoder(com.google.javascript.jscomp.VariableMap variableMap)
func NewReplacedStringsDecoder(variableMap *VariableMap) *ReplacedStringsDecoder {
	clinit_ReplacedStringsDecoder()
	this_ := &ReplacedStringsDecoder{}
	InitReplacedStringsDecoder(this_, variableMap)
	return this_
}

func InitReplacedStringsDecoder(this_ *ReplacedStringsDecoder, variableMap *VariableMap) {
	_ = variableMap
	this_.initFields_ReplacedStringsDecoder()
	this_.originalToNewNameMap = variableMap.GetOriginalNameToNewNameMap()
}

// Java method: ReplacedStringsDecoder.decode(java.lang.String encodedStr)
func (this_ *ReplacedStringsDecoder) Decode(encodedStr string) string {
	_ = encodedStr
	var suppliedBits []string = j2gort.JStrSplitLimit(encodedStr, "`", int32(-1))
	_ = suppliedBits
	var originalStr string = j2gort.CastString(j2gort.MapGetStr(this_.originalToNewNameMap, suppliedBits[int32(0)]))
	_ = originalStr
	if originalStr == j2gort.NullString {
		return encodedStr
	}
	var originalBits []string = j2gort.JStrSplitLimit(originalStr, "`", int32(-1))
	_ = originalBits
	var sb j2gort.LocalStringBuilder
	sb.InitStr(originalBits[int32(0)])

	{
		var i int32 = int32(1)
		_ = i
		for ; i < max(int32(len(originalBits)), int32(len(suppliedBits))); i = i + 1 {
			sb.AppendStr(func() string {
				if i < int32(len(suppliedBits)) {
					return suppliedBits[i]
				}
				return "-"
			}())
			sb.AppendStr(func() string {
				if i < int32(len(originalBits)) {
					return originalBits[i]
				}
				return "-"
			}())
		}
	}
	return sb.String()
}

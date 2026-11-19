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
// Port of com/google/javascript/jscomp/OutputCharsetEncoder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.OutputCharsetEncoder
type OutputCharsetEncoder struct {
	encoder j2gort.Ext_java_nio_charset_CharsetEncoder `gson:"encoder"`
	idHash_ atomic.Int32
}

func (this_ *OutputCharsetEncoder) BindSelf_(s any) {
	_ = s
}
func (this_ *OutputCharsetEncoder) IsJavaNil_() bool { return this_ == nil }
func (*OutputCharsetEncoder) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.OutputCharsetEncoder", "OutputCharsetEncoder" }
func (this_ *OutputCharsetEncoder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OutputCharsetEncoder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OutputCharsetEncoder", "com.google.javascript.jscomp.OutputCharsetEncoder", func(v_ any) bool {
		return j2gort.InstanceOf[*OutputCharsetEncoder](v_)
	})
}

func (this_ *OutputCharsetEncoder) initFields_OutputCharsetEncoder() {
	_ = this_
}

// Java constructor: OutputCharsetEncoder(java.nio.charset.Charset outputCharset)
func NewOutputCharsetEncoder(outputCharset j2gort.Ext_java_nio_charset_Charset) *OutputCharsetEncoder {
	this_ := &OutputCharsetEncoder{}
	InitOutputCharsetEncoder(this_, outputCharset)
	return this_
}

func InitOutputCharsetEncoder(this_ *OutputCharsetEncoder, outputCharset j2gort.Ext_java_nio_charset_Charset) {
	_ = outputCharset
	this_.initFields_OutputCharsetEncoder()
	if (((outputCharset == nil) || j2gort.RefEquals(outputCharset, j2gort.Ext_java_nio_charset_StandardCharsets_F_US_ASCII))) {
		this_.encoder = nil
	} else {
		this_.encoder = outputCharset.NewEncoder()
	}
}

// Java method: OutputCharsetEncoder.canEncode(char c)
func (this_ *OutputCharsetEncoder) CanEncode(c uint16) bool {
	_ = c
	return ((this_.encoder != nil) && this_.encoder.CanEncode_1_uint16(c))
}


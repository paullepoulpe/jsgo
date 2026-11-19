/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/GoogleJsMessageIdGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.GoogleJsMessageIdGenerator
type GoogleJsMessageIdGenerator struct {
	projectId string `gson:"projectId"`
	idHash_   atomic.Int32
}

func (this_ *GoogleJsMessageIdGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *GoogleJsMessageIdGenerator) IsJavaNil_() bool { return this_ == nil }
func (*GoogleJsMessageIdGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GoogleJsMessageIdGenerator", "GoogleJsMessageIdGenerator"
}
func (this_ *GoogleJsMessageIdGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GoogleJsMessageIdGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GoogleJsMessageIdGenerator", "com.google.javascript.jscomp.GoogleJsMessageIdGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogleJsMessageIdGenerator](v_)
	})
}
func (*GoogleJsMessageIdGenerator) Implements_jscomp_JsMessage_IdGenerator() {}

func (this_ *GoogleJsMessageIdGenerator) initFields_GoogleJsMessageIdGenerator() {
	_ = this_
	this_.projectId = j2gort.NullString
}

// Java constructor: GoogleJsMessageIdGenerator(java.lang.String projectId)
func NewGoogleJsMessageIdGenerator(projectId string) *GoogleJsMessageIdGenerator {
	this_ := &GoogleJsMessageIdGenerator{}
	InitGoogleJsMessageIdGenerator(this_, projectId)
	return this_
}

func InitGoogleJsMessageIdGenerator(this_ *GoogleJsMessageIdGenerator, projectId string) {
	_ = projectId
	this_.initFields_GoogleJsMessageIdGenerator()
	this_.projectId = projectId
}

// Java method: GoogleJsMessageIdGenerator.generateId(java.lang.String meaning, java.util.List<com.google.javascript.jscomp.JsMessage.Part> messageParts)
func (this_ *GoogleJsMessageIdGenerator) GenerateId(meaning string, messageParts j2gort.Ext_java_util_List) string {
	_ = meaning
	_ = messageParts
	if !(meaning != j2gort.NullString) {
		j2gort.CheckState(false)
	}
	var sb j2gort.LocalStringBuilder

	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(messageParts))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			part := j2gort.AsIface[JsMessage_Part](rawElem_2)
			_ = part
			if part.IsPlaceholder() {
				sb.AppendStr(part.GetCanonicalPlaceholderName())
			} else {
				sb.AppendStr(part.GetString())
			}
		}
	}
	var tcValue string = sb.String()
	_ = tcValue
	var projectScopedMeaning string = (j2gort.StrValString(func() string {
		if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.projectId) {
			return ""
		}
		return (j2gort.StrValString(this_.projectId) + ": ")
	}()) + j2gort.StrValString(meaning))
	_ = projectScopedMeaning
	return strconv.FormatInt(GoogleJsMessageIdGenerator_MessageId_GenerateId(tcValue, projectScopedMeaning), 10)
}

// Java class: com.google.javascript.jscomp.GoogleJsMessageIdGenerator.FP
type GoogleJsMessageIdGenerator_FP struct {
	idHash_ atomic.Int32
}

func (this_ *GoogleJsMessageIdGenerator_FP) BindSelf_(s any) {
	_ = s
}
func (this_ *GoogleJsMessageIdGenerator_FP) IsJavaNil_() bool { return this_ == nil }
func (*GoogleJsMessageIdGenerator_FP) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GoogleJsMessageIdGenerator$FP", "FP"
}
func (this_ *GoogleJsMessageIdGenerator_FP) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GoogleJsMessageIdGenerator_FP) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GoogleJsMessageIdGenerator.FP", "com.google.javascript.jscomp.GoogleJsMessageIdGenerator$FP", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogleJsMessageIdGenerator_FP](v_)
	})
}

func (this_ *GoogleJsMessageIdGenerator_FP) initFields_GoogleJsMessageIdGenerator_FP() {
	_ = this_
}

// Java constructor: FP()
func NewGoogleJsMessageIdGenerator_FP() *GoogleJsMessageIdGenerator_FP {
	this_ := &GoogleJsMessageIdGenerator_FP{}
	InitGoogleJsMessageIdGenerator_FP(this_)
	return this_
}

func InitGoogleJsMessageIdGenerator_FP(this_ *GoogleJsMessageIdGenerator_FP) {
	this_.initFields_GoogleJsMessageIdGenerator_FP()
}

// Java static method: FP.fingerprint(byte[] str, int start, int limit)
func GoogleJsMessageIdGenerator_FP_Fingerprint_3(str []int8, start int32, limit int32) int64 {
	_ = str
	_ = start
	_ = limit
	var hi int32 = GoogleJsMessageIdGenerator_FP_Hash32(str, start, limit, int32(0))
	_ = hi
	var lo int32 = GoogleJsMessageIdGenerator_FP_Hash32(str, start, limit, int32(102072))
	_ = lo
	if (hi == int32(0)) && ((lo == int32(0)) || (lo == int32(1))) {
		hi = (hi ^ int32(319790063))
		lo = (lo ^ int32(-1801410264))
	}
	return ((javalang.Shl64((int64(hi)), int(int32(32)))) | (int64(lo) & int64(4294967295)))
}

// Java static method: FP.fingerprint(java.lang.String str)
func GoogleJsMessageIdGenerator_FP_Fingerprint_1(str string) int64 {
	_ = str
	var tmp []int8
	{
		_ = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
		tmp = j2gort.JStrGetBytes(str)
	}

	_ = tmp
	return GoogleJsMessageIdGenerator_FP_Fingerprint_3(tmp, int32(0), int32(len(tmp)))
}

// Java static method: FP.hash32(byte[] str, int start, int limit, int c)
func GoogleJsMessageIdGenerator_FP_Hash32(str []int8, start int32, limit int32, c int32) int32 {
	_ = str
	_ = start
	_ = limit
	_ = c
	var a int32 = int32(-1640531527)
	_ = a
	var b int32 = int32(-1640531527)
	_ = b
	var i int32 = int32(0)
	_ = i
	{
		i = start
		for ; (i + int32(12)) <= limit; i = (i + int32(12)) {
			a = (a + ((((javalang.Shl32((int32(str[(i+int32(0))]) & int32(255)), int(int32(0)))) | (javalang.Shl32((int32(str[(i+int32(1))]) & int32(255)), int(int32(8))))) | (javalang.Shl32((int32(str[(i+int32(2))]) & int32(255)), int(int32(16))))) | (javalang.Shl32((int32(str[(i+int32(3))]) & int32(255)), int(int32(24))))))
			b = (b + ((((javalang.Shl32((int32(str[(i+int32(4))]) & int32(255)), int(int32(0)))) | (javalang.Shl32((int32(str[(i+int32(5))]) & int32(255)), int(int32(8))))) | (javalang.Shl32((int32(str[(i+int32(6))]) & int32(255)), int(int32(16))))) | (javalang.Shl32((int32(str[(i+int32(7))]) & int32(255)), int(int32(24))))))
			c = (c + ((((javalang.Shl32((int32(str[(i+int32(8))]) & int32(255)), int(int32(0)))) | (javalang.Shl32((int32(str[(i+int32(9))]) & int32(255)), int(int32(8))))) | (javalang.Shl32((int32(str[(i+int32(10))]) & int32(255)), int(int32(16))))) | (javalang.Shl32((int32(str[(i+int32(11))]) & int32(255)), int(int32(24))))))
			a = (a - b)
			a = (a - c)
			a = (a ^ (javalang.Ushr32(c, int(int32(13)))))
			b = (b - c)
			b = (b - a)
			b = (b ^ (javalang.Shl32(a, int(int32(8)))))
			c = (c - a)
			c = (c - b)
			c = (c ^ (javalang.Ushr32(b, int(int32(13)))))
			a = (a - b)
			a = (a - c)
			a = (a ^ (javalang.Ushr32(c, int(int32(12)))))
			b = (b - c)
			b = (b - a)
			b = (b ^ (javalang.Shl32(a, int(int32(16)))))
			c = (c - a)
			c = (c - b)
			c = (c ^ (javalang.Ushr32(b, int(int32(5)))))
			a = (a - b)
			a = (a - c)
			a = (a ^ (javalang.Ushr32(c, int(int32(3)))))
			b = (b - c)
			b = (b - a)
			b = (b ^ (javalang.Shl32(a, int(int32(10)))))
			c = (c - a)
			c = (c - b)
			c = (c ^ (javalang.Ushr32(b, int(int32(15)))))
		}
	}
	c = (c + (limit - start))
	var tmp int32 = (limit - i)
	_ = tmp
	if tmp == int32(11) {
		c = (c + javalang.Shl32((int32(str[(i+int32(10))])&int32(255)), int(int32(24))))
	}
	if tmp >= int32(10) {
		c = (c + javalang.Shl32((int32(str[(i+int32(9))])&int32(255)), int(int32(16))))
	}
	if tmp >= int32(9) {
		c = (c + javalang.Shl32((int32(str[(i+int32(8))])&int32(255)), int(int32(8))))
	}
	if tmp >= int32(8) {
		b = (b + javalang.Shl32((int32(str[(i+int32(7))])&int32(255)), int(int32(24))))
	}
	if tmp >= int32(7) {
		b = (b + javalang.Shl32((int32(str[(i+int32(6))])&int32(255)), int(int32(16))))
	}
	if tmp >= int32(6) {
		b = (b + javalang.Shl32((int32(str[(i+int32(5))])&int32(255)), int(int32(8))))
	}
	if tmp >= int32(5) {
		b = (b + (int32(str[(i+int32(4))]) & int32(255)))
	}
	if tmp >= int32(4) {
		a = (a + javalang.Shl32((int32(str[(i+int32(3))])&int32(255)), int(int32(24))))
	}
	if tmp >= int32(3) {
		a = (a + javalang.Shl32((int32(str[(i+int32(2))])&int32(255)), int(int32(16))))
	}
	if tmp >= int32(2) {
		a = (a + javalang.Shl32((int32(str[(i+int32(1))])&int32(255)), int(int32(8))))
	}
	if tmp >= int32(1) {
		a = (a + (int32(str[(i+int32(0))]) & int32(255)))
	}
	a = (a - b)
	a = (a - c)
	a = (a ^ (javalang.Ushr32(c, int(int32(13)))))
	b = (b - c)
	b = (b - a)
	b = (b ^ (javalang.Shl32(a, int(int32(8)))))
	c = (c - a)
	c = (c - b)
	c = (c ^ (javalang.Ushr32(b, int(int32(13)))))
	a = (a - b)
	a = (a - c)
	a = (a ^ (javalang.Ushr32(c, int(int32(12)))))
	b = (b - c)
	b = (b - a)
	b = (b ^ (javalang.Shl32(a, int(int32(16)))))
	c = (c - a)
	c = (c - b)
	c = (c ^ (javalang.Ushr32(b, int(int32(5)))))
	a = (a - b)
	a = (a - c)
	a = (a ^ (javalang.Ushr32(c, int(int32(3)))))
	b = (b - c)
	b = (b - a)
	b = (b ^ (javalang.Shl32(a, int(int32(10)))))
	c = (c - a)
	c = (c - b)
	c = (c ^ (javalang.Ushr32(b, int(int32(15)))))
	return c
}

// Java class: com.google.javascript.jscomp.GoogleJsMessageIdGenerator.MessageId
type GoogleJsMessageIdGenerator_MessageId struct {
	idHash_ atomic.Int32
}

func (this_ *GoogleJsMessageIdGenerator_MessageId) BindSelf_(s any) {
	_ = s
}
func (this_ *GoogleJsMessageIdGenerator_MessageId) IsJavaNil_() bool { return this_ == nil }
func (*GoogleJsMessageIdGenerator_MessageId) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GoogleJsMessageIdGenerator$MessageId", "MessageId"
}
func (this_ *GoogleJsMessageIdGenerator_MessageId) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GoogleJsMessageIdGenerator_MessageId) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GoogleJsMessageIdGenerator.MessageId", "com.google.javascript.jscomp.GoogleJsMessageIdGenerator$MessageId", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogleJsMessageIdGenerator_MessageId](v_)
	})
}

func (this_ *GoogleJsMessageIdGenerator_MessageId) initFields_GoogleJsMessageIdGenerator_MessageId() {
	_ = this_
}

// Java constructor: MessageId()
func NewGoogleJsMessageIdGenerator_MessageId() *GoogleJsMessageIdGenerator_MessageId {
	this_ := &GoogleJsMessageIdGenerator_MessageId{}
	InitGoogleJsMessageIdGenerator_MessageId(this_)
	return this_
}

func InitGoogleJsMessageIdGenerator_MessageId(this_ *GoogleJsMessageIdGenerator_MessageId) {
	this_.initFields_GoogleJsMessageIdGenerator_MessageId()
}

// Java static method: MessageId.generateId(java.lang.String message, java.lang.String meaning)
func GoogleJsMessageIdGenerator_MessageId_GenerateId(message string, meaning string) int64 {
	_ = message
	_ = meaning
	var fp int64 = GoogleJsMessageIdGenerator_FP_Fingerprint_1(message)
	_ = fp
	if (meaning != j2gort.NullString) && (j2gort.JStrLen(meaning) > int32(0)) {
		var fp2 int64 = GoogleJsMessageIdGenerator_FP_Fingerprint_1(meaning)
		_ = fp2
		fp = ((fp2 + (javalang.Shl64(fp, int(int32(1))))) + int64((func() int32 {
			if fp < int64(int32(0)) {
				return int32(1)
			}
			return int32(0)
		}())))
	}
	return (fp & int64(9223372036854775807))
}

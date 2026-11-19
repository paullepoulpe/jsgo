/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/SourceMapJsonLexer.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.SourceMapJsonLexer
type SourceMapJsonLexer struct {
	json    string `gson:"json"`
	pos     int32  `gson:"pos"`
	length  int32  `gson:"length"`
	idHash_ atomic.Int32
}

func (this_ *SourceMapJsonLexer) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapJsonLexer) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapJsonLexer) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.SourceMapJsonLexer", "SourceMapJsonLexer"
}
func (this_ *SourceMapJsonLexer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapJsonLexer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapJsonLexer", "com.google.debugging.sourcemap.SourceMapJsonLexer", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapJsonLexer](v_)
	})
}

func (this_ *SourceMapJsonLexer) initFields_SourceMapJsonLexer() {
	_ = this_
	this_.json = j2gort.NullString

}

// Java constructor: SourceMapJsonLexer(java.lang.String json)
func NewSourceMapJsonLexer(json string) *SourceMapJsonLexer {
	this_ := &SourceMapJsonLexer{}
	InitSourceMapJsonLexer(this_, json)
	return this_
}

func InitSourceMapJsonLexer(this_ *SourceMapJsonLexer, json string) {
	_ = json
	this_.initFields_SourceMapJsonLexer()
	this_.json = json
	this_.length = j2gort.JStrLen(json)
}

// Java method: SourceMapJsonLexer.skipWhitespace()
func (this_ *SourceMapJsonLexer) SkipWhitespace() {
	for this_.pos < this_.length {
		var c uint16 = j2gort.JStrCharAt(this_.json, this_.pos)
		_ = c
		if (((int32(c) == int32(uint16(32))) || (int32(c) == int32(uint16(9)))) || (int32(c) == int32(uint16(10)))) || (int32(c) == int32(uint16(13))) {
			this_.pos = this_.pos + 1
		} else {
			break
		}
	}
}

// Java method: SourceMapJsonLexer.hasNext()
func (this_ *SourceMapJsonLexer) HasNext() bool {
	this_.SkipWhitespace()
	if this_.pos >= this_.length {
		return false
	}
	var c uint16 = j2gort.JStrCharAt(this_.json, this_.pos)
	_ = c
	return ((int32(c) != int32(uint16(93))) && (int32(c) != int32(uint16(125))))
}

// Java method: SourceMapJsonLexer.beginObject()
func (this_ *SourceMapJsonLexer) BeginObject() {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(123))) {
		this_.pos = this_.pos + 1
	} else {
		panic(NewSourceMapParseException("Expected '{'"))
	}
}

// Java method: SourceMapJsonLexer.endObject()
func (this_ *SourceMapJsonLexer) EndObject() {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(125))) {
		this_.pos = this_.pos + 1
	} else {
		panic(NewSourceMapParseException("Expected '}'"))
	}
}

// Java method: SourceMapJsonLexer.beginArray()
func (this_ *SourceMapJsonLexer) BeginArray() {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(91))) {
		this_.pos = this_.pos + 1
	} else {
		panic(NewSourceMapParseException("Expected '['"))
	}
}

// Java method: SourceMapJsonLexer.endArray()
func (this_ *SourceMapJsonLexer) EndArray() {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(93))) {
		this_.pos = this_.pos + 1
	} else {
		panic(NewSourceMapParseException("Expected ']'"))
	}
}

// Java method: SourceMapJsonLexer.nextName()
func (this_ *SourceMapJsonLexer) NextName() string {
	var name string = this_.NextString()
	_ = name
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(58))) {
		this_.pos = this_.pos + 1
	} else {
		panic(NewSourceMapParseException("Expected ':'"))
	}
	return name
}

// Java method: SourceMapJsonLexer.nextStringOrNull()
func (this_ *SourceMapJsonLexer) NextStringOrNull() string {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && j2gort.JStrStartsWithOffset(this_.json, "null", this_.pos) {
		this_.pos = (this_.pos + int32(4))
		return j2gort.NullString
	}
	return this_.NextString()
}

// Java method: SourceMapJsonLexer.nextString()
func (this_ *SourceMapJsonLexer) NextString() string {
	this_.SkipWhitespace()
	if (this_.pos >= this_.length) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) != int32(uint16(34))) {
		panic(NewSourceMapParseException("Expected string"))
	}
	this_.pos = this_.pos + 1
	var start int32 = this_.pos
	_ = start
	var hasEscape bool = false
	_ = hasEscape
	for this_.pos < this_.length {
		var c uint16 = j2gort.JStrCharAt(this_.json, this_.pos)
		_ = c
		if int32(c) == int32(uint16(34)) {
			break
		}
		if int32(c) == int32(uint16(92)) {
			hasEscape = true
			this_.pos = this_.pos + 1
		}
		this_.pos = this_.pos + 1
	}
	if this_.pos >= this_.length {
		panic(NewSourceMapParseException("Unterminated string"))
	}
	var result string = j2gort.NullString
	_ = result
	if !(hasEscape) {
		result = j2gort.JStrSub(this_.json, start, this_.pos)
	} else {
		result = SourceMapJsonLexer_Unescape(this_.json, start, this_.pos)
	}
	this_.pos = this_.pos + 1
	return result
}

// Java method: SourceMapJsonLexer.nextInt()
func (this_ *SourceMapJsonLexer) NextInt() int32 {
	this_.SkipWhitespace()
	var start int32 = this_.pos
	_ = start
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(45))) {
		this_.pos = this_.pos + 1
	}
	for ((this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) >= int32(uint16(48)))) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) <= int32(uint16(57))) {
		this_.pos = this_.pos + 1
	}
	if start == this_.pos {
		panic(NewSourceMapParseException("Expected integer"))
	}
	var didRet_1 bool
	var retVal_2 int32
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_4) {
					e := j2gort.AsIface[javalang.INumberFormatException](jex_4)
					_ = e
					panic(NewSourceMapParseException(("Invalid integer: " + j2gort.StrValString(j2gort.JStrSub(this_.json, start, this_.pos)))))
				} else {
					panic(rec_3)
				}
			}
		}()
		retVal_2 = j2gort.BoxInt_ParseInt_4(j2gort.BoxString(this_.json), start, this_.pos, int32(10))
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java method: SourceMapJsonLexer.nextRawValue()
func (this_ *SourceMapJsonLexer) NextRawValue() string {
	this_.SkipWhitespace()
	var start int32 = this_.pos
	_ = start
	this_.SkipValue()
	return j2gort.JStrSub(this_.json, start, this_.pos)
}

// Java method: SourceMapJsonLexer.skipValue()
func (this_ *SourceMapJsonLexer) SkipValue() {
	this_.SkipWhitespace()
	if this_.pos >= this_.length {
		panic(NewSourceMapParseException("Unexpected end of input"))
	}
	var c uint16 = j2gort.JStrCharAt(this_.json, this_.pos)
	_ = c
	if int32(c) == int32(uint16(123)) {
		this_.pos = this_.pos + 1
		for this_.HasNext() {
			this_.NextName()
			this_.SkipValue()
			this_.CheckComma()
		}
		this_.EndObject()
	} else if int32(c) == int32(uint16(91)) {
		this_.pos = this_.pos + 1
		for this_.HasNext() {
			this_.SkipValue()
			this_.CheckComma()
		}
		this_.EndArray()
	} else if int32(c) == int32(uint16(34)) {
		this_.NextString()
	} else if (int32(c) == int32(uint16(116))) && j2gort.JStrStartsWithOffset(this_.json, "true", this_.pos) {
		this_.pos = (this_.pos + int32(4))
	} else if (int32(c) == int32(uint16(102))) && j2gort.JStrStartsWithOffset(this_.json, "false", this_.pos) {
		this_.pos = (this_.pos + int32(5))
	} else if (int32(c) == int32(uint16(110))) && j2gort.JStrStartsWithOffset(this_.json, "null", this_.pos) {
		this_.pos = (this_.pos + int32(4))
	} else if (int32(c) == int32(uint16(45))) || ((int32(c) >= int32(uint16(48))) && (int32(c) <= int32(uint16(57)))) {
		if int32(c) == int32(uint16(45)) {
			this_.pos = this_.pos + 1
		}
		for (this_.pos < this_.length) && (((((((int32(j2gort.JStrCharAt(this_.json, this_.pos)) >= int32(uint16(48))) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) <= int32(uint16(57)))) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(46)))) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(101)))) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(69)))) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(43)))) || (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(45)))) {
			this_.pos = this_.pos + 1
		}
	} else {
		panic(NewSourceMapParseException(("Unexpected character: " + j2gort.StrValChar(c))))
	}
}

// Java method: SourceMapJsonLexer.checkComma()
func (this_ *SourceMapJsonLexer) CheckComma() {
	var hasComma bool = this_.TryConsumeComma()
	_ = hasComma
	if hasComma && !(this_.HasNext()) {
		panic(NewSourceMapParseException("Unexpected trailing comma"))
	}
	if !(hasComma) && this_.HasNext() {
		panic(NewSourceMapParseException("Expected ','"))
	}
}

// Java method: SourceMapJsonLexer.tryConsumeComma()
func (this_ *SourceMapJsonLexer) TryConsumeComma() bool {
	this_.SkipWhitespace()
	if (this_.pos < this_.length) && (int32(j2gort.JStrCharAt(this_.json, this_.pos)) == int32(uint16(44))) {
		this_.pos = this_.pos + 1
		return true
	}
	return false
}

// Java static method: SourceMapJsonLexer.unescape(java.lang.String s, int start, int end)
func SourceMapJsonLexer_Unescape(s string, start int32, end int32) string {
	_ = s
	_ = start
	_ = end
	var sb j2gort.LocalStringBuilder
	sb.InitCap((end - start))

	{
		var i int32 = start
		_ = i
		for ; i < end; i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			if (int32(c) == int32(uint16(92))) && ((i + int32(1)) < end) {
				i = i + 1
				var next uint16 = j2gort.JStrCharAt(s, i)
				_ = next
				switch next {
				case uint16(34):
					sb.AppendChar(uint16(34))
				case uint16(92):
					sb.AppendChar(uint16(92))
				case uint16(47):
					sb.AppendChar(uint16(47))
				case uint16(98):
					sb.AppendChar(uint16(8))
				case uint16(102):
					sb.AppendChar(uint16(12))
				case uint16(110):
					sb.AppendChar(uint16(10))
				case uint16(114):
					sb.AppendChar(uint16(13))
				case uint16(116):
					sb.AppendChar(uint16(9))
				case uint16(117):
					if (i + int32(4)) >= end {
						panic(NewSourceMapParseException("Invalid unicode escape sequence"))
					}
					func() {
						defer func() {
							if rec_5 := recover(); rec_5 != nil {
								jex_6 := javalang.RecoverException(rec_5)
								_ = jex_6
								if j2gort.InstanceOf[javalang.INumberFormatException](jex_6) {
									e := j2gort.AsIface[javalang.INumberFormatException](jex_6)
									_ = e
									panic(NewSourceMapParseException(("Invalid unicode escape: " + j2gort.StrValString(j2gort.JStrSub(s, (i+int32(1)), (i+int32(5)))))))
								} else {
									panic(rec_5)
								}
							}
						}()
						sb.AppendChar(uint16(j2gort.BoxInt_ParseInt_4(j2gort.BoxString(s), (i + int32(1)), (i + int32(5)), int32(16))))
						i = (i + int32(4))
					}()
				default:
					sb.AppendChar(next)
				}
			} else {
				sb.AppendChar(c)
			}
		}
	}
	return sb.String()
}

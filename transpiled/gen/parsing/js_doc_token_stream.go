/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/JsDocTokenStream.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.JsDocTokenStream
type JsDocTokenStream struct {
	string_            string   `gson:"string"`
	stringBuffer       []uint16 `gson:"stringBuffer"`
	sourceString       string   `gson:"sourceString"`
	stringBufferTop    int32    `gson:"stringBufferTop"`
	rawTokenStart      int32    `gson:"rawTokenStart"`
	unget0             int32    `gson:"unget0"`
	unget1             int32    `gson:"unget1"`
	unget2             int32    `gson:"unget2"`
	ungetCursor        int32    `gson:"ungetCursor"`
	lineStart          int32    `gson:"lineStart"`
	lineEndChar        int32    `gson:"lineEndChar"`
	lineno             int32    `gson:"lineno"`
	charno             int32    `gson:"charno"`
	initCharno         int32    `gson:"initCharno"`
	initLineno         int32    `gson:"initLineno"`
	sourceEnd          int32    `gson:"sourceEnd"`
	sourceCursor       int32    `gson:"sourceCursor"`
	cursor             int32    `gson:"cursor"`
	idHash_            atomic.Int32
	rawTokenContiguous bool `gson:"rawTokenContiguous"`
}

func (this_ *JsDocTokenStream) BindSelf_(s any) {
	_ = s
}
func (this_ *JsDocTokenStream) IsJavaNil_() bool { return this_ == nil }
func (*JsDocTokenStream) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocTokenStream", "JsDocTokenStream"
}
func (this_ *JsDocTokenStream) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsDocTokenStream) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocTokenStream", "com.google.javascript.jscomp.parsing.JsDocTokenStream", func(v_ any) bool {
		return j2gort.InstanceOf[*JsDocTokenStream](v_)
	})
}

func StaticGet_JsDocTokenStream_EOF_CHAR() int32 { return int32(-1) }

var static_JsDocTokenStream_EMPTY_CHAR_ARRAY []uint16

func StaticGet_JsDocTokenStream_EMPTY_CHAR_ARRAY() []uint16 {
	clinit_JsDocTokenStream()
	return static_JsDocTokenStream_EMPTY_CHAR_ARRAY
}
func StaticSet_JsDocTokenStream_EMPTY_CHAR_ARRAY(v []uint16) {
	clinit_JsDocTokenStream()
	static_JsDocTokenStream_EMPTY_CHAR_ARRAY = v
}

var clinitDone_JsDocTokenStream uint32
var clinitState_JsDocTokenStream j2gort.ClinitState

func clinit_JsDocTokenStream() {
	if atomic.LoadUint32(&clinitDone_JsDocTokenStream) != 0 {
		return
	}
	clinitSlow_JsDocTokenStream()
}

func clinitSlow_JsDocTokenStream() {
	if !j2gort.EnterClinit(&clinitState_JsDocTokenStream, &clinitDone_JsDocTokenStream) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsDocTokenStream, &clinitDone_JsDocTokenStream)
	static_JsDocTokenStream_EMPTY_CHAR_ARRAY = make([]uint16, int32(0))
}

func (this_ *JsDocTokenStream) initFields_JsDocTokenStream() {
	_ = this_

	this_.stringBuffer = StaticGet_JsDocTokenStream_EMPTY_CHAR_ARRAY()
	this_.rawTokenStart = int32(-1)

	this_.lineEndChar = int32(-1)
	this_.charno = int32(-1)
	this_.sourceString = j2gort.NullString
}

// Java constructor: JsDocTokenStream(java.lang.String sourceString)
func NewJsDocTokenStream_1(sourceString string) *JsDocTokenStream {
	clinit_JsDocTokenStream()
	this_ := &JsDocTokenStream{}
	InitJsDocTokenStream_1(this_, sourceString)
	return this_
}

func InitJsDocTokenStream_1(this_ *JsDocTokenStream, sourceString string) {
	_ = sourceString
	InitJsDocTokenStream_2(this_, sourceString, int32(0))
}

// Java constructor: JsDocTokenStream(java.lang.String sourceString, int lineno)
func NewJsDocTokenStream_2(sourceString string, lineno int32) *JsDocTokenStream {
	clinit_JsDocTokenStream()
	this_ := &JsDocTokenStream{}
	InitJsDocTokenStream_2(this_, sourceString, lineno)
	return this_
}

func InitJsDocTokenStream_2(this_ *JsDocTokenStream, sourceString string, lineno int32) {
	_ = sourceString
	_ = lineno
	InitJsDocTokenStream_3(this_, sourceString, lineno, int32(0))
}

// Java constructor: JsDocTokenStream(java.lang.String sourceString, int lineno, int initCharno)
func NewJsDocTokenStream_3(sourceString string, lineno int32, initCharno int32) *JsDocTokenStream {
	clinit_JsDocTokenStream()
	this_ := &JsDocTokenStream{}
	InitJsDocTokenStream_3(this_, sourceString, lineno, initCharno)
	return this_
}

func InitJsDocTokenStream_3(this_ *JsDocTokenStream, sourceString string, lineno int32, initCharno int32) {
	_ = sourceString
	_ = lineno
	_ = initCharno
	this_.initFields_JsDocTokenStream()
	this_.Reset(sourceString, lineno, initCharno)
}

// Java method: JsDocTokenStream.reset(java.lang.String sourceString, int lineno, int initCharno)
func (this_ *JsDocTokenStream) Reset(sourceString string, lineno int32, initCharno int32) {
	_ = sourceString
	_ = lineno
	_ = initCharno
	j2gort.CheckNotNullString(sourceString)
	this_.string_ = ""
	this_.stringBufferTop = int32(0)
	this_.rawTokenStart = int32(-1)
	this_.rawTokenContiguous = false
	this_.ungetCursor = int32(0)
	this_.lineStart = int32(0)
	this_.lineEndChar = int32(-1)
	this_.lineno = lineno
	this_.charno = int32(-1)
	this_.sourceString = sourceString
	this_.sourceEnd = j2gort.JStrLen(sourceString)
	this_.sourceCursor = func() int32 { var asgn_1 int32 = int32(0); this_.cursor = asgn_1; return asgn_1 }()
	this_.initLineno = lineno
	this_.initCharno = initCharno
}

// Java method: JsDocTokenStream.clearSource()
func (this_ *JsDocTokenStream) ClearSource() {
	this_.sourceString = ""
	this_.sourceEnd = int32(0)
	this_.string_ = ""
	this_.stringBufferTop = int32(0)
	this_.rawTokenStart = int32(-1)
	this_.rawTokenContiguous = false
}

// Java method: JsDocTokenStream.resetStringBuffer()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_resetStringBuffer_7105() {
	this_.stringBufferTop = int32(0)
	this_.rawTokenStart = int32(-1)
	this_.rawTokenContiguous = false
}

// Java method: JsDocTokenStream.getJsDocToken()
func (this_ *JsDocTokenStream) GetJsDocToken() JsDocToken {
	var c int32 = int32(0)
	_ = c
	this_.priv_JsDocTokenStream_resetStringBuffer_7105()
	for {
		for {
			this_.charno = int32(-1)
			c = this_.GetChar()
			if c == int32(-1) {
				return StaticGet_JsDocToken_EOF()
			} else if c == int32(uint16(10)) {
				return StaticGet_JsDocToken_EOL()
			} else if !(pkg_rhino.TokenUtil_IsJSSpace(c)) {
				break
			}
		}
		switch c {
		case int32(uint16(64)):
			for doFirst_2 := true; doFirst_2 || (true); doFirst_2 = false {
				c = this_.GetChar()
				if JsDocTokenStream_IsAlpha(c) {
					this_.priv_JsDocTokenStream_addToString_7112(c)
				} else {
					this_.UngetChar(c)
					this_.string_ = this_.priv_JsDocTokenStream_getStringFromBuffer_7111()
					this_.priv_JsDocTokenStream_resetStringBuffer_7105()
					return StaticGet_JsDocToken_ANNOTATION()
				}
			}
			fallthrough
		case int32(uint16(42)):
			if this_.priv_JsDocTokenStream_matchChar_7117(int32(uint16(47))) {
				return StaticGet_JsDocToken_EOC()
			} else {
				return StaticGet_JsDocToken_STAR()
			}
		case int32(uint16(44)):
			return StaticGet_JsDocToken_COMMA()
		case int32(uint16(62)):
			return StaticGet_JsDocToken_RIGHT_ANGLE()
		case int32(uint16(40)):
			return StaticGet_JsDocToken_LEFT_PAREN()
		case int32(uint16(41)):
			return StaticGet_JsDocToken_RIGHT_PAREN()
		case int32(uint16(123)):
			return StaticGet_JsDocToken_LEFT_CURLY()
		case int32(uint16(125)):
			return StaticGet_JsDocToken_RIGHT_CURLY()
		case int32(uint16(91)):
			return StaticGet_JsDocToken_LEFT_SQUARE()
		case int32(uint16(93)):
			return StaticGet_JsDocToken_RIGHT_SQUARE()
		case int32(uint16(63)):
			return StaticGet_JsDocToken_QMARK()
		case int32(uint16(33)):
			return StaticGet_JsDocToken_BANG()
		case int32(uint16(58)):
			return StaticGet_JsDocToken_COLON()
		case int32(uint16(61)):
			return StaticGet_JsDocToken_EQUALS()
		case int32(uint16(124)):
			return StaticGet_JsDocToken_PIPE()
		case int32(uint16(60)):
			return StaticGet_JsDocToken_LEFT_ANGLE()
		case int32(uint16(46)):
			c = this_.GetChar()
			if c == int32(uint16(60)) {
				return StaticGet_JsDocToken_LEFT_ANGLE()
			} else {
				if c == int32(uint16(46)) {
					c = this_.GetChar()
					if c == int32(uint16(46)) {
						return StaticGet_JsDocToken_ITER_REST()
					} else {
						this_.priv_JsDocTokenStream_addToString_7112(int32(uint16(46)))
					}
				}
				this_.priv_JsDocTokenStream_pushUnget_7113(c)
				c = int32(uint16(46))
			}
			fallthrough
		default:
			{
				var c1 int32 = c
				_ = c1
				this_.priv_JsDocTokenStream_addToString_7112(c)
				var c2 int32 = this_.GetChar()
				_ = c2
				if !(JsDocTokenStream_IsJSDocString(c2)) {
					this_.UngetChar(c2)
					this_.string_ = this_.priv_JsDocTokenStream_getStringFromBuffer_7111()
					this_.priv_JsDocTokenStream_resetStringBuffer_7105()
					return StaticGet_JsDocToken_STRING()
				} else {
					for doFirst_3 := true; doFirst_3 || (true); doFirst_3 = false {
						c1 = c2
						c2 = this_.GetChar()
						if (c1 == int32(uint16(46))) && (c2 == int32(uint16(60))) {
							this_.UngetChar(c2)
							this_.UngetChar(c1)
							this_.string_ = this_.priv_JsDocTokenStream_getStringFromBuffer_7111()
							this_.priv_JsDocTokenStream_resetStringBuffer_7105()
							return StaticGet_JsDocToken_STRING()
						} else {
							if JsDocTokenStream_IsJSDocString(c2) {
								this_.priv_JsDocTokenStream_addToString_7112(c1)
							} else {
								this_.UngetChar(c2)
								this_.priv_JsDocTokenStream_addToString_7112(c1)
								this_.string_ = this_.priv_JsDocTokenStream_getStringFromBuffer_7111()
								this_.priv_JsDocTokenStream_resetStringBuffer_7105()
								return StaticGet_JsDocToken_STRING()
							}
						}
					}
				}
			}
		}
	}
	panic("unreachable")
}

// Java method: JsDocTokenStream.getRemainingJSDocLine()
func (this_ *JsDocTokenStream) GetRemainingJSDocLine() string {
	var c int32 = int32(0)
	_ = c
	for {
		c = this_.GetChar()
		switch c {
		case int32(uint16(42)):
			if this_.priv_JsDocTokenStream_peekChar_7119() != int32(uint16(47)) {
				this_.priv_JsDocTokenStream_addToString_7112(c)
				break
			}
			fallthrough
		case int32(-1), int32(uint16(10)):
			this_.UngetChar(c)
			this_.string_ = this_.priv_JsDocTokenStream_getStringFromBuffer_7111()
			this_.priv_JsDocTokenStream_resetStringBuffer_7105()
			return this_.string_
		default:
			this_.priv_JsDocTokenStream_addToString_7112(c)
		}
	}
	panic("unreachable")
}

// Java method: JsDocTokenStream.getLineno()
func (this_ *JsDocTokenStream) GetLineno() int32 {
	return this_.lineno
}

// Java method: JsDocTokenStream.getCharno()
func (this_ *JsDocTokenStream) GetCharno() int32 {
	if this_.lineno == this_.initLineno {
		return (this_.initCharno + this_.charno)
	} else {
		return this_.charno
	}
}

// Java method: JsDocTokenStream.getString()
func (this_ *JsDocTokenStream) GetString() string {
	return this_.string_
}

// Java method: JsDocTokenStream.getStringFromBuffer()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_getStringFromBuffer_7111() string {
	var top int32 = this_.stringBufferTop
	_ = top
	if top == int32(0) {
		return ""
	}
	if this_.rawTokenContiguous {
		return j2gort.JStrSub(this_.sourceString, this_.rawTokenStart, (this_.rawTokenStart + top))
	}
	return func() string {
		ca_4 := this_.stringBuffer
		co_5 := int32(0)
		cc_6 := top
		return javalang.StringFromUTF16(ca_4[co_5 : co_5+cc_6])
	}()
}

// Java method: JsDocTokenStream.addToString(int c)
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_addToString_7112(c int32) {
	_ = c
	var n int32 = this_.stringBufferTop
	_ = n
	if n == int32(0) {
		var start int32 = ((this_.sourceCursor - this_.ungetCursor) - int32(1))
		_ = start
		if ((start >= int32(0)) && (start < this_.sourceEnd)) && (int32(j2gort.JStrCharAt(this_.sourceString, start)) == int32(uint16(c))) {
			this_.rawTokenStart = start
			this_.rawTokenContiguous = true
		} else {
			this_.rawTokenStart = int32(-1)
			this_.rawTokenContiguous = false
		}
	} else if this_.rawTokenContiguous {
		var pos int32 = (this_.rawTokenStart + n)
		_ = pos
		if (pos >= this_.sourceEnd) || (int32(j2gort.JStrCharAt(this_.sourceString, pos)) != int32(uint16(c))) {
			this_.rawTokenContiguous = false
		}
	}
	if n == int32(len(this_.stringBuffer)) {
		var newCap int32
		_ = newCap
		if n == int32(0) {
			newCap = int32(64)
		} else {
			newCap = (n * int32(2))
		}
		var tmp []uint16 = make([]uint16, newCap)
		_ = tmp
		if n > int32(0) {
			func() {
				acSrc_7 := this_.stringBuffer
				acSP_8 := int32(0)
				acDst_9 := tmp
				acDP_10 := int32(0)
				acLen_11 := n
				j2gort.CheckNotNullSlice(acSrc_7)
				j2gort.CheckNotNullSlice(acDst_9)
				if acSP_8 < 0 || acDP_10 < 0 || acLen_11 < 0 || int64(acSP_8)+int64(acLen_11) > int64(len(acSrc_7)) || int64(acDP_10)+int64(acLen_11) > int64(len(acDst_9)) {
					panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
				}
				copy(acDst_9[int(acDP_10):int(acDP_10)+int(acLen_11)], acSrc_7[int(acSP_8):int(acSP_8)+int(acLen_11)])
			}()
		}
		this_.stringBuffer = tmp
	}
	{
		a_12 := this_.stringBuffer
		i_13 := n
		a_12[i_13] = uint16(c)
	}
	this_.stringBufferTop = (n + int32(1))
}

// Java method: JsDocTokenStream.pushUnget(int c)
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_pushUnget_7113(c int32) {
	_ = c
	switch func() int32 {
		old_14 := this_.ungetCursor
		nw_15 := old_14 + 1
		this_.ungetCursor = nw_15
		return old_14
	}() {
	case int32(0):
		this_.unget0 = c
	case int32(1):
		this_.unget1 = c
	case int32(2):
		this_.unget2 = c
	default:
		panic(j2gort.AsIface[javalang.IArrayIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException", j2gort.ValueOfInt((this_.ungetCursor - int32(1))))))
	}
}

// Java method: JsDocTokenStream.popUnget()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_popUnget_7114() int32 {
	switch func() int32 {
		old_16 := this_.ungetCursor
		nw_17 := old_16 - 1
		this_.ungetCursor = nw_17
		return nw_17
	}() {
	case int32(0):
		return this_.unget0
	case int32(1):
		return this_.unget1
	case int32(2):
		return this_.unget2
	default:
		panic(j2gort.AsIface[javalang.IArrayIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException", j2gort.ValueOfInt(this_.ungetCursor))))
	}
	panic("unreachable switch expression")
}

// Java method: JsDocTokenStream.lastUnget()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_lastUnget_7115() int32 {
	switch this_.ungetCursor - int32(1) {
	case int32(0):
		return this_.unget0
	case int32(1):
		return this_.unget1
	default:
		return this_.unget2
	}
	panic("unreachable switch expression")
}

// Java method: JsDocTokenStream.ungetChar(int c)
func (this_ *JsDocTokenStream) UngetChar(c int32) {
	_ = c
	if !(!((this_.ungetCursor != int32(0)) && (this_.priv_JsDocTokenStream_lastUnget_7115() == int32(uint16(10))))) {
		panic(javalang.NewAssertionError(nil))
	}
	this_.priv_JsDocTokenStream_pushUnget_7113(c)
	this_.cursor = this_.cursor - 1
}

// Java method: JsDocTokenStream.matchChar(int test)
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_matchChar_7117(test int32) bool {
	_ = test
	var c int32 = this_.priv_JsDocTokenStream_getCharIgnoreLineEnd_7121()
	_ = c
	if c == test {
		return true
	} else {
		this_.priv_JsDocTokenStream_ungetCharIgnoreLineEnd_7122(c)
		return false
	}
}

// Java static method: JsDocTokenStream.isAlpha(int c)
func JsDocTokenStream_IsAlpha(c int32) bool {
	clinit_JsDocTokenStream()
	_ = c
	if c <= int32(uint16(90)) {
		return (int32(uint16(65)) <= c)
	} else {
		return ((int32(uint16(97)) <= c) && (c <= int32(uint16(122))))
	}
}

// Java static method: JsDocTokenStream.isJSDocString(int c)
func JsDocTokenStream_IsJSDocString(c int32) bool {
	clinit_JsDocTokenStream()
	_ = c
	switch c {
	case int32(uint16(64)), int32(uint16(42)), int32(uint16(44)), int32(uint16(60)), int32(uint16(62)), int32(uint16(58)), int32(uint16(40)), int32(uint16(41)), int32(uint16(123)), int32(uint16(125)), int32(uint16(91)), int32(uint16(93)), int32(uint16(63)), int32(uint16(33)), int32(uint16(124)), int32(uint16(61)), int32(-1), int32(uint16(10)):
		return false
	default:
		return !(pkg_rhino.TokenUtil_IsJSSpace(c))
	}
	panic("unreachable switch expression")
}

// Java method: JsDocTokenStream.update()
func (this_ *JsDocTokenStream) Update() {
	this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
}

// Java method: JsDocTokenStream.peekChar()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_peekChar_7119() int32 {
	var c int32 = this_.GetChar()
	_ = c
	this_.UngetChar(c)
	return c
}

// Java method: JsDocTokenStream.getChar()
func (this_ *JsDocTokenStream) GetChar() int32 {
	if this_.ungetCursor != int32(0) {
		this_.cursor = this_.cursor + 1
		var c int32 = this_.priv_JsDocTokenStream_popUnget_7114()
		_ = c
		if this_.charno == int32(-1) {
			this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
		}
		return c
	}
	for {
		var c_2 int32 = int32(0)
		_ = c_2
		if this_.sourceCursor == this_.sourceEnd {
			if this_.charno == int32(-1) {
				this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
			}
			return int32(-1)
		}
		this_.cursor = this_.cursor + 1
		c_2 = int32(j2gort.JStrCharAt(this_.sourceString, func() int32 {
			old_18 := this_.sourceCursor
			nw_19 := old_18 + 1
			this_.sourceCursor = nw_19
			return old_18
		}()))
		if this_.lineEndChar >= int32(0) {
			if (this_.lineEndChar == int32(uint16(13))) && (c_2 == int32(uint16(10))) {
				this_.lineEndChar = int32(uint16(10))
				continue
			}
			this_.lineEndChar = int32(-1)
			this_.lineStart = (this_.sourceCursor - int32(1))
			this_.lineno = this_.lineno + 1
		}
		if c_2 <= int32(127) {
			if (c_2 == int32(uint16(10))) || (c_2 == int32(uint16(13))) {
				this_.lineEndChar = c_2
				c_2 = int32(uint16(10))
			}
		} else {
			if pkg_rhino.TokenUtil_IsJSFormatChar(c_2) {
				continue
			}
			if JsDocTokenStream_IsJSLineTerminator(c_2) {
				this_.lineEndChar = c_2
				c_2 = int32(uint16(10))
			}
		}
		if this_.charno == int32(-1) {
			this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
		}
		return c_2
	}
	panic("unreachable")
}

// Java method: JsDocTokenStream.getCharIgnoreLineEnd()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_getCharIgnoreLineEnd_7121() int32 {
	if this_.ungetCursor != int32(0) {
		this_.cursor = this_.cursor + 1
		var c int32 = this_.priv_JsDocTokenStream_popUnget_7114()
		_ = c
		if this_.charno == int32(-1) {
			this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
		}
		return c
	}
	for {
		var c_2 int32 = int32(0)
		_ = c_2
		if this_.sourceCursor == this_.sourceEnd {
			if this_.charno == int32(-1) {
				this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
			}
			return int32(-1)
		}
		this_.cursor = this_.cursor + 1
		c_2 = int32(j2gort.JStrCharAt(this_.sourceString, func() int32 {
			old_20 := this_.sourceCursor
			nw_21 := old_20 + 1
			this_.sourceCursor = nw_21
			return old_20
		}()))
		if c_2 <= int32(127) {
			if (c_2 == int32(uint16(10))) || (c_2 == int32(uint16(13))) {
				this_.lineEndChar = c_2
				c_2 = int32(uint16(10))
			}
		} else {
			if pkg_rhino.TokenUtil_IsJSFormatChar(c_2) {
				continue
			}
			if JsDocTokenStream_IsJSLineTerminator(c_2) {
				this_.lineEndChar = c_2
				c_2 = int32(uint16(10))
			}
		}
		if this_.charno == int32(-1) {
			this_.charno = this_.priv_JsDocTokenStream_getLineOffset_7123()
		}
		return c_2
	}
	panic("unreachable")
}

// Java static method: JsDocTokenStream.isJSLineTerminator(int c)
func JsDocTokenStream_IsJSLineTerminator(c int32) bool {
	clinit_JsDocTokenStream()
	_ = c
	if (c & int32(57296)) != int32(0) {
		return false
	}
	return ((((c == int32(uint16(10))) || (c == int32(uint16(13)))) || (c == int32(8232))) || (c == int32(8233)))
}

// Java method: JsDocTokenStream.ungetCharIgnoreLineEnd(int c)
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_ungetCharIgnoreLineEnd_7122(c int32) {
	_ = c
	this_.priv_JsDocTokenStream_pushUnget_7113(c)
	this_.cursor = this_.cursor - 1
}

// Java method: JsDocTokenStream.getLineOffset()
func (this_ *JsDocTokenStream) priv_JsDocTokenStream_getLineOffset_7123() int32 {
	return (((this_.sourceCursor - this_.lineStart) - this_.ungetCursor) - int32(1))
}

// Java method: JsDocTokenStream.getCursor()
func (this_ *JsDocTokenStream) GetCursor() int32 {
	return this_.cursor
}

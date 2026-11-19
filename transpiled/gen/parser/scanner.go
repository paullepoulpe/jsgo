/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.Scanner
type Scanner struct {
	errorReporter      ErrorReporter           `gson:"errorReporter"`
	source             *SourceFile             `gson:"source"`
	lineNumberScanner  *LineNumberScanner      `gson:"lineNumberScanner"`
	contents           string                  `gson:"contents"`
	currentTokens      []Token                 `gson:"currentTokens"`
	commentRecorder    Scanner_CommentRecorder `gson:"commentRecorder"`
	contentsLength     int32                   `gson:"contentsLength"`
	index              int32                   `gson:"index"`
	typeParameterLevel int32                   `gson:"typeParameterLevel"`
	idHash_            atomic.Int32
}

func (this_ *Scanner) BindSelf_(s any) {
	_ = s
}
func (this_ *Scanner) IsJavaNil_() bool { return this_ == nil }
func (*Scanner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Scanner", "Scanner"
}
func (this_ *Scanner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Scanner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Scanner", "com.google.javascript.jscomp.parsing.parser.Scanner", func(v_ any) bool {
		return j2gort.InstanceOf[*Scanner](v_)
	})
}

func (this_ *Scanner) initFields_Scanner() {
	_ = this_
	this_.contents = j2gort.NullString
	this_.currentTokens = []Token{}
}

// Java constructor: Scanner(com.google.javascript.jscomp.parsing.parser.util.ErrorReporter errorReporter, com.google.javascript.jscomp.parsing.parser.Scanner.CommentRecorder commentRecorder, com.google.javascript.jscomp.parsing.parser.SourceFile file, int offset)
func NewScanner(errorReporter ErrorReporter, commentRecorder Scanner_CommentRecorder, file *SourceFile, offset int32) *Scanner {
	this_ := &Scanner{}
	InitScanner(this_, errorReporter, commentRecorder, file, offset)
	return this_
}

func InitScanner(this_ *Scanner, errorReporter ErrorReporter, commentRecorder Scanner_CommentRecorder, file *SourceFile, offset int32) {
	_ = errorReporter
	_ = commentRecorder
	_ = file
	_ = offset
	this_.initFields_Scanner()
	this_.errorReporter = errorReporter
	this_.commentRecorder = commentRecorder
	this_.source = file
	this_.lineNumberScanner = NewLineNumberScanner(this_.source)
	this_.contents = file.contents
	this_.contentsLength = j2gort.JStrLen(file.contents)
	this_.index = offset
	this_.typeParameterLevel = int32(0)
}

// Java method: Scanner.getFile()
func (this_ *Scanner) GetFile() *SourceFile {
	return this_.source
}

// Java method: Scanner.getOffset()
func (this_ *Scanner) GetOffset() int32 {
	if len(this_.currentTokens) == 0 {
		return this_.index
	} else {
		return this_.PeekToken_0().location.start.offset
	}
}

// Java method: Scanner.setPosition(com.google.javascript.jscomp.parsing.parser.util.SourcePosition position)
func (this_ *Scanner) SetPosition(position SourcePosition) {
	_ = position
	this_.lineNumberScanner.RewindTo(position)
	clear(this_.currentTokens)
	this_.currentTokens = this_.currentTokens[:0]
	this_.index = position.offset
}

// Java method: Scanner.getPosition()
func (this_ *Scanner) GetPosition_0() SourcePosition {
	if len(this_.currentTokens) == 0 {
		return this_.priv_Scanner_getPosition_7469(this_.index)
	} else {
		return this_.PeekToken_0().location.start
	}
}

// Java method: Scanner.getPosition(int offset)
func (this_ *Scanner) priv_Scanner_getPosition_7469(offset int32) SourcePosition {
	_ = offset
	return this_.lineNumberScanner.GetSourcePosition(offset)
}

// Java method: Scanner.getTokenRange(int startOffset)
func (this_ *Scanner) priv_Scanner_getTokenRange_7470(startOffset int32) SourceRange {
	_ = startOffset
	return this_.lineNumberScanner.GetSourceRange(startOffset, this_.index)
}

// Java method: Scanner.getTokenRange(com.google.javascript.jscomp.parsing.parser.util.SourcePosition position)
func (this_ *Scanner) priv_Scanner_getTokenRange_7471(position SourcePosition) SourceRange {
	_ = position
	this_.lineNumberScanner.RewindTo(position)
	return this_.lineNumberScanner.GetSourceRange(position.offset, this_.index)
}

// Java method: Scanner.nextToken()
func (this_ *Scanner) NextToken() Token {
	this_.PeekToken_0()
	return j2gort.SliceRemoveAtElem(&this_.currentTokens, int32(0))
}

// Java method: Scanner.clearTokenLookahead()
func (this_ *Scanner) priv_Scanner_clearTokenLookahead_7473() {
	if len(this_.currentTokens) != 0 {
		this_.SetPosition(this_.PeekToken_0().location.start)
	}
}

// Java method: Scanner.nextRegularExpressionLiteralToken()
func (this_ *Scanner) NextRegularExpressionLiteralToken() Token {
	this_.priv_Scanner_clearTokenLookahead_7473()
	var beginToken int32 = this_.index
	_ = beginToken
	this_.priv_Scanner_nextChar_7520()
	if !(this_.priv_Scanner_skipRegularExpressionBody_7476()) {
		return Token_Literal(StaticGet_TokenType_REGULAR_EXPRESSION(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	}
	if int32(this_.priv_Scanner_peekChar_7522()) != int32(uint16(47)) {
		this_.priv_Scanner_reportError_7524("Expected '/' in regular expression literal", []any{})
		return Token_Literal(StaticGet_TokenType_REGULAR_EXPRESSION(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	}
	this_.priv_Scanner_nextChar_7520()
	for Identifiers_IsIdentifierPart(this_.priv_Scanner_peekChar_7522()) {
		this_.priv_Scanner_nextChar_7520()
	}
	return Token_Literal(StaticGet_TokenType_REGULAR_EXPRESSION(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
}

// Java method: Scanner.nextTemplateLiteralToken()
func (this_ *Scanner) NextTemplateLiteralToken() Token {
	var token Token = this_.NextToken()
	_ = token
	if this_.priv_Scanner_isAtEnd_7484() || (token.type_ != StaticGet_TokenType_CLOSE_CURLY()) {
		this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(this_.index), "Expected '}' after expression in template literal", []any{})
	}
	return this_.priv_Scanner_nextTemplateLiteralTokenShared_7503(StaticGet_TokenType_TEMPLATE_TAIL(), StaticGet_TokenType_TEMPLATE_MIDDLE())
}

// Java method: Scanner.skipRegularExpressionBody()
func (this_ *Scanner) priv_Scanner_skipRegularExpressionBody_7476() bool {
	if !(Scanner_IsRegularExpressionFirstChar(this_.priv_Scanner_peekChar_7522())) {
		this_.priv_Scanner_reportError_7524("Expected regular expression first char", []any{})
		return false
	}
	if !(this_.priv_Scanner_skipRegularExpressionChar_7477()) {
		return false
	}
	for !(this_.priv_Scanner_isAtEnd_7484()) && Scanner_IsRegularExpressionChar(this_.priv_Scanner_peekChar_7522()) {
		if !(this_.priv_Scanner_skipRegularExpressionChar_7477()) {
			return false
		}
	}
	return true
}

// Java method: Scanner.skipRegularExpressionChar()
func (this_ *Scanner) priv_Scanner_skipRegularExpressionChar_7477() bool {
	switch this_.priv_Scanner_peekChar_7522() {
	case uint16(92):
		return this_.priv_Scanner_skipRegularExpressionBackslashSequence_7478()
	case uint16(91):
		return this_.priv_Scanner_skipRegularExpressionClass_7479()
	default:
		this_.priv_Scanner_nextChar_7520()
		return true
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.skipRegularExpressionBackslashSequence()
func (this_ *Scanner) priv_Scanner_skipRegularExpressionBackslashSequence_7478() bool {
	this_.priv_Scanner_nextChar_7520()
	if Scanner_IsLineTerminator(this_.priv_Scanner_peekChar_7522()) {
		this_.priv_Scanner_reportError_7524("New line not allowed in regular expression literal", []any{})
		return false
	}
	this_.priv_Scanner_nextChar_7520()
	return true
}

// Java method: Scanner.skipRegularExpressionClass()
func (this_ *Scanner) priv_Scanner_skipRegularExpressionClass_7479() bool {
	this_.priv_Scanner_nextChar_7520()
	for !(this_.priv_Scanner_isAtEnd_7484()) && this_.priv_Scanner_peekRegularExpressionClassChar_7480() {
		if !(this_.priv_Scanner_skipRegularExpressionClassChar_7481()) {
			return false
		}
	}
	if int32(this_.priv_Scanner_peekChar_7522()) != int32(uint16(93)) {
		this_.priv_Scanner_reportError_7524("']' expected", []any{})
		return false
	}
	this_.priv_Scanner_nextChar_7520()
	return true
}

// Java method: Scanner.peekRegularExpressionClassChar()
func (this_ *Scanner) priv_Scanner_peekRegularExpressionClassChar_7480() bool {
	return ((int32(this_.priv_Scanner_peekChar_7522()) != int32(uint16(93))) && !(Scanner_IsLineTerminator(this_.priv_Scanner_peekChar_7522())))
}

// Java method: Scanner.skipRegularExpressionClassChar()
func (this_ *Scanner) priv_Scanner_skipRegularExpressionClassChar_7481() bool {
	if this_.priv_Scanner_peek_7521(uint16(92)) {
		return this_.priv_Scanner_skipRegularExpressionBackslashSequence_7478()
	}
	this_.priv_Scanner_nextChar_7520()
	return true
}

// Java static method: Scanner.isRegularExpressionFirstChar(char ch)
func Scanner_IsRegularExpressionFirstChar(ch uint16) bool {
	_ = ch
	return (Scanner_IsRegularExpressionChar(ch) && (int32(ch) != int32(uint16(42))))
}

// Java static method: Scanner.isRegularExpressionChar(char ch)
func Scanner_IsRegularExpressionChar(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(47):
		return false
	case uint16(92), uint16(91):
		return true
	default:
		return !(Scanner_IsLineTerminator(ch))
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.peekToken()
func (this_ *Scanner) PeekToken_0() Token {
	return this_.PeekToken_1(int32(0))
}

// Java method: Scanner.peekToken(int index)
func (this_ *Scanner) PeekToken_1(index int32) Token {
	_ = index
	for int32(len(this_.currentTokens)) <= index {
		this_.currentTokens = append(this_.currentTokens, this_.priv_Scanner_scanToken_7495())
	}
	return this_.currentTokens[index]
}

// Java method: Scanner.isAtEnd()
func (this_ *Scanner) priv_Scanner_isAtEnd_7484() bool {
	return !(this_.priv_Scanner_isValidIndex_7485(this_.index))
}

// Java method: Scanner.isValidIndex(int index)
func (this_ *Scanner) priv_Scanner_isValidIndex_7485(index int32) bool {
	_ = index
	return ((index >= int32(0)) && (index < this_.contentsLength))
}

// Java method: Scanner.skipWhitespace()
func (this_ *Scanner) priv_Scanner_skipWhitespace_7486() bool {
	var foundLineTerminator bool = false
	_ = foundLineTerminator
	for !(this_.priv_Scanner_isAtEnd_7484()) && this_.priv_Scanner_peekWhitespace_7487() {
		if Scanner_IsLineTerminator(this_.priv_Scanner_nextChar_7520()) {
			foundLineTerminator = true
		}
	}
	return foundLineTerminator
}

// Java method: Scanner.peekWhitespace()
func (this_ *Scanner) priv_Scanner_peekWhitespace_7487() bool {
	return Scanner_IsWhitespace(this_.priv_Scanner_peekChar_7522())
}

// Java static method: Scanner.isWhitespace(char ch)
func Scanner_IsWhitespace(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(9), uint16(11), uint16(12), uint16(32), uint16(160), uint16(65279), uint16(10), uint16(13), uint16(8232), uint16(8233), uint16(12288):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: Scanner.isLineTerminator(char ch)
func Scanner_IsLineTerminator(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(10), uint16(13), uint16(8232), uint16(8233):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: Scanner.isStringLineTerminator(char ch)
func Scanner_IsStringLineTerminator(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(8232), uint16(8233):
		return false
	default:
		return Scanner_IsLineTerminator(ch)
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.skipComments()
func (this_ *Scanner) priv_Scanner_skipComments_7488() {
	for this_.priv_Scanner_skipComment_7489() {
	}
}

// Java method: Scanner.skipComment()
func (this_ *Scanner) priv_Scanner_skipComment_7489() bool {
	var isStartOfLine bool = this_.priv_Scanner_skipWhitespace_7486()
	_ = isStartOfLine
	if !(this_.priv_Scanner_isAtEnd_7484()) {
		switch this_.priv_Scanner_peekChar_7523(int32(0)) {
		case uint16(47):
			switch this_.priv_Scanner_peekChar_7523(int32(1)) {
			case uint16(47):
				this_.priv_Scanner_skipSingleLineComment_7491()
				return true
			case uint16(42):
				this_.priv_Scanner_skipMultiLineComment_7494()
				return true
			default:
			}
		case uint16(60):
			if ((int32(this_.priv_Scanner_peekChar_7523(int32(1))) == int32(uint16(33))) && (int32(this_.priv_Scanner_peekChar_7523(int32(2))) == int32(uint16(45)))) && (int32(this_.priv_Scanner_peekChar_7523(int32(3))) == int32(uint16(45))) {
				this_.priv_Scanner_reportHtmlCommentWarning_7490()
				this_.priv_Scanner_skipSingleLineComment_7491()
				return true
			}
		case uint16(45):
			if (isStartOfLine && (int32(this_.priv_Scanner_peekChar_7523(int32(1))) == int32(uint16(45)))) && (int32(this_.priv_Scanner_peekChar_7523(int32(2))) == int32(uint16(62))) {
				this_.priv_Scanner_reportHtmlCommentWarning_7490()
				this_.priv_Scanner_skipSingleLineComment_7491()
				return true
			}
		case uint16(35):
			if (this_.index == int32(0)) && (int32(this_.priv_Scanner_peekChar_7523(int32(1))) == int32(uint16(33))) {
				this_.priv_Scanner_skipSingleLineComment_7492(StaticGet_Comment_Type_SHEBANG())
				return true
			}
		default:
		}
	}
	return false
}

// Java method: Scanner.reportHtmlCommentWarning()
func (this_ *Scanner) priv_Scanner_reportHtmlCommentWarning_7490() {
	this_.priv_Scanner_reportWarning_7526("In some cases, '<!--' and '-->' are treated as a '//' for legacy reasons. Removing this from your code is safe for all browsers currently in use.", []any{})
}

// Java method: Scanner.skipSingleLineComment()
func (this_ *Scanner) priv_Scanner_skipSingleLineComment_7491() {
	this_.priv_Scanner_skipSingleLineComment_7492(StaticGet_Comment_Type_LINE())
}

// Java method: Scanner.skipSingleLineComment(com.google.javascript.jscomp.parsing.parser.trees.Comment.Type type)
func (this_ *Scanner) priv_Scanner_skipSingleLineComment_7492(type_ Comment_Type) {
	_ = type_
	var startOffset int32 = this_.index
	_ = startOffset
	for !(this_.priv_Scanner_isAtEnd_7484()) && !(Scanner_IsLineTerminator(this_.priv_Scanner_peekChar_7522())) {
		this_.priv_Scanner_nextChar_7520()
	}
	var range_ SourceRange = this_.lineNumberScanner.GetSourceRange(startOffset, this_.index)
	_ = range_
	var value string = j2gort.JStrSub(this_.contents, startOffset, this_.index)
	_ = value
	this_.priv_Scanner_recordComment_7493(type_, range_, value)
}

// Java method: Scanner.recordComment(com.google.javascript.jscomp.parsing.parser.trees.Comment.Type type, com.google.javascript.jscomp.parsing.parser.util.SourceRange range, java.lang.String value)
func (this_ *Scanner) priv_Scanner_recordComment_7493(type_ Comment_Type, range_ SourceRange, value string) {
	_ = type_
	_ = range_
	_ = value
	this_.commentRecorder.RecordComment(type_, range_, value)
}

// Java method: Scanner.skipMultiLineComment()
func (this_ *Scanner) priv_Scanner_skipMultiLineComment_7494() {
	var startOffset int32 = this_.index
	_ = startOffset
	this_.priv_Scanner_nextChar_7520()
	this_.priv_Scanner_nextChar_7520()
	for !(this_.priv_Scanner_isAtEnd_7484()) && ((int32(this_.priv_Scanner_peekChar_7522()) != int32(uint16(42))) || (int32(this_.priv_Scanner_peekChar_7523(int32(1))) != int32(uint16(47)))) {
		this_.priv_Scanner_nextChar_7520()
	}
	if !(this_.priv_Scanner_isAtEnd_7484()) {
		this_.priv_Scanner_nextChar_7520()
		this_.priv_Scanner_nextChar_7520()
		var type_ Comment_Type = StaticGet_Comment_Type_BLOCK()
		_ = type_
		if (this_.index - startOffset) > int32(4) {
			if int32(j2gort.JStrCharAt(this_.contents, (startOffset+int32(2)))) == int32(uint16(42)) {
				type_ = StaticGet_Comment_Type_JSDOC()
			} else if int32(j2gort.JStrCharAt(this_.contents, (startOffset+int32(2)))) == int32(uint16(33)) {
				type_ = StaticGet_Comment_Type_IMPORTANT()
			}
		}
		var range_ SourceRange = this_.lineNumberScanner.GetSourceRange(startOffset, this_.index)
		_ = range_
		var value string = j2gort.JStrSub(this_.contents, startOffset, this_.index)
		_ = value
		this_.priv_Scanner_recordComment_7493(type_, range_, value)
	} else {
		this_.priv_Scanner_reportError_7524("unterminated comment", []any{})
	}
}

// Java method: Scanner.scanToken()
func (this_ *Scanner) priv_Scanner_scanToken_7495() Token {
	this_.priv_Scanner_skipComments_7488()
	var beginToken int32 = this_.index
	_ = beginToken
	if this_.priv_Scanner_isAtEnd_7484() {
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_END_OF_FILE(), beginToken)
	}
	var ch uint16 = this_.priv_Scanner_nextChar_7520()
	_ = ch
	switch ch {
	case uint16(123):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OPEN_CURLY(), beginToken)
	case uint16(125):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CLOSE_CURLY(), beginToken)
	case uint16(40):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OPEN_PAREN(), beginToken)
	case uint16(41):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CLOSE_PAREN(), beginToken)
	case uint16(91):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OPEN_SQUARE(), beginToken)
	case uint16(93):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CLOSE_SQUARE(), beginToken)
	case uint16(46):
		if Scanner_IsDecimalDigit(this_.priv_Scanner_peekChar_7522()) {
			return this_.priv_Scanner_scanNumberPostPeriod_7496(beginToken)
		}
		if this_.priv_Scanner_peek_7521(uint16(46)) && (int32(this_.priv_Scanner_peekChar_7523(int32(1))) == int32(uint16(46))) {
			this_.priv_Scanner_nextChar_7520()
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ELLIPSIS(), beginToken)
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PERIOD(), beginToken)
	case uint16(59):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_SEMI_COLON(), beginToken)
	case uint16(44):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_COMMA(), beginToken)
	case uint16(126):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_TILDE(), beginToken)
	case uint16(63):
		if this_.priv_Scanner_peek_7521(uint16(63)) {
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_QUESTION_QUESTION_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_QUESTION_QUESTION(), beginToken)
		}
		if this_.priv_Scanner_peek_7521(uint16(46)) {
			if !(Scanner_IsDecimalDigit(this_.priv_Scanner_peekChar_7523(int32(1)))) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_QUESTION_DOT(), beginToken)
			}
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_QUESTION(), beginToken)
	case uint16(58):
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_COLON(), beginToken)
	case uint16(60):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(60):
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_LEFT_SHIFT_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_LEFT_SHIFT(), beginToken)
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_LESS_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OPEN_ANGLE(), beginToken)
		}
	case uint16(62):
		if this_.typeParameterLevel > int32(0) {
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CLOSE_ANGLE(), beginToken)
		}
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(62):
			this_.priv_Scanner_nextChar_7520()
			switch this_.priv_Scanner_peekChar_7522() {
			case uint16(61):
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_RIGHT_SHIFT_EQUAL(), beginToken)
			case uint16(62):
				this_.priv_Scanner_nextChar_7520()
				if this_.priv_Scanner_peek_7521(uint16(61)) {
					this_.priv_Scanner_nextChar_7520()
					return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL(), beginToken)
				}
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_UNSIGNED_RIGHT_SHIFT(), beginToken)
			default:
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_RIGHT_SHIFT(), beginToken)
			}
			fallthrough
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_GREATER_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CLOSE_ANGLE(), beginToken)
		}
	case uint16(61):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_EQUAL_EQUAL_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_EQUAL_EQUAL(), beginToken)
		case uint16(62):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ARROW(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_EQUAL(), beginToken)
		}
	case uint16(33):
		if this_.priv_Scanner_peek_7521(uint16(61)) {
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_NOT_EQUAL_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_NOT_EQUAL(), beginToken)
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_BANG(), beginToken)
	case uint16(42):
		if this_.priv_Scanner_peek_7521(uint16(61)) {
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_STAR_EQUAL(), beginToken)
		} else if this_.priv_Scanner_peek_7521(uint16(42)) {
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_STAR_STAR_EQUAL(), beginToken)
			} else {
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_STAR_STAR(), beginToken)
			}
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_STAR(), beginToken)
	case uint16(37):
		if this_.priv_Scanner_peek_7521(uint16(61)) {
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PERCENT_EQUAL(), beginToken)
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PERCENT(), beginToken)
	case uint16(94):
		if this_.priv_Scanner_peek_7521(uint16(61)) {
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CARET_EQUAL(), beginToken)
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_CARET(), beginToken)
	case uint16(47):
		if this_.priv_Scanner_peek_7521(uint16(61)) {
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_SLASH_EQUAL(), beginToken)
		}
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_SLASH(), beginToken)
	case uint16(43):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(43):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PLUS_PLUS(), beginToken)
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PLUS_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_PLUS(), beginToken)
		}
	case uint16(45):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(45):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_MINUS_MINUS(), beginToken)
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_MINUS_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_MINUS(), beginToken)
		}
	case uint16(38):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(38):
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_AND_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_AND(), beginToken)
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_AMPERSAND_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_AMPERSAND(), beginToken)
		}
	case uint16(124):
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(124):
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(61)) {
				this_.priv_Scanner_nextChar_7520()
				return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OR_EQUAL(), beginToken)
			}
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_OR(), beginToken)
		case uint16(61):
			this_.priv_Scanner_nextChar_7520()
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_BAR_EQUAL(), beginToken)
		default:
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_BAR(), beginToken)
		}
	case uint16(35):
		if this_.priv_Scanner_peek_7521(uint16(33)) {
			this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(this_.index), "Shebang comment must be at the start of the file", []any{})
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ERROR(), beginToken)
		}
		return this_.priv_Scanner_scanIdentifierOrKeyword_7500(beginToken, ch)
	case uint16(48):
		return this_.priv_Scanner_scanPostZero_7498(beginToken)
	case uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
		return this_.priv_Scanner_scanPostDigit_7497(beginToken)
	case uint16(34), uint16(39):
		return this_.priv_Scanner_scanStringLiteral_7501(beginToken, ch)
	case uint16(96):
		return this_.priv_Scanner_scanTemplateLiteral_7502(beginToken)
	default:
		return this_.priv_Scanner_scanIdentifierOrKeyword_7500(beginToken, ch)
	}
	panic("unreachable")
}

// Java method: Scanner.scanNumberPostPeriod(int beginToken)
func (this_ *Scanner) priv_Scanner_scanNumberPostPeriod_7496(beginToken int32) Token {
	_ = beginToken
	this_.priv_Scanner_skipDecimalDigits_7514()
	return this_.priv_Scanner_scanExponentOfNumericLiteral_7513(beginToken)
}

// Java method: Scanner.scanPostDigit(int beginToken)
func (this_ *Scanner) priv_Scanner_scanPostDigit_7497(beginToken int32) Token {
	_ = beginToken
	this_.priv_Scanner_skipDecimalDigits_7514()
	if this_.priv_Scanner_peek_7521(uint16(110)) {
		this_.priv_Scanner_nextChar_7520()
		return Token_Literal(StaticGet_TokenType_BIGINT(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	}
	return this_.priv_Scanner_scanFractionalNumericLiteral_7512(beginToken)
}

// Java method: Scanner.scanPostZero(int beginToken)
func (this_ *Scanner) priv_Scanner_scanPostZero_7498(beginToken int32) Token {
	_ = beginToken
	var isBigInt bool = false
	_ = isBigInt
	switch this_.priv_Scanner_peekChar_7522() {
	case uint16(98), uint16(66):
		this_.priv_Scanner_nextChar_7520()
		if !(Scanner_IsBinaryDigit(this_.priv_Scanner_peekChar_7522())) {
			this_.priv_Scanner_reportError_7524("Binary Integer Literal must contain at least one digit", []any{})
		}
		this_.priv_Scanner_skipBinaryDigits_7519()
		isBigInt = this_.priv_Scanner_peek_7521(uint16(110))
		if isBigInt {
			this_.priv_Scanner_nextChar_7520()
		}
		var iifeArg_1 TokenType
	iifeSwitch_1:
		switch {
		default:

			if isBigInt {
				iifeArg_1 = StaticGet_TokenType_BIGINT()
				break iifeSwitch_1
			}
			iifeArg_1 = StaticGet_TokenType_NUMBER()
			break iifeSwitch_1
		}
		return Token_Literal(iifeArg_1,

			this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	case uint16(111), uint16(79):
		this_.priv_Scanner_nextChar_7520()
		if !(Scanner_IsOctalDigit(this_.priv_Scanner_peekChar_7522())) {
			this_.priv_Scanner_reportError_7524("Octal Integer Literal must contain at least one digit", []any{})
		}
		this_.priv_Scanner_skipOctalDigits_7518()
		if this_.priv_Scanner_peek_7521(uint16(56)) || this_.priv_Scanner_peek_7521(uint16(57)) {
			this_.priv_Scanner_reportError_7524("Invalid octal digit in octal literal.", []any{})
		}
		isBigInt = this_.priv_Scanner_peek_7521(uint16(110))
		if isBigInt {
			this_.priv_Scanner_nextChar_7520()
		}
		var iifeArg_2 TokenType
	iifeSwitch_2:
		switch {
		default:

			if isBigInt {
				iifeArg_2 = StaticGet_TokenType_BIGINT()
				break iifeSwitch_2
			}
			iifeArg_2 = StaticGet_TokenType_NUMBER()
			break iifeSwitch_2
		}
		return Token_Literal(iifeArg_2,

			this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	case uint16(120), uint16(88):
		this_.priv_Scanner_nextChar_7520()
		if !(this_.priv_Scanner_peekHexDigit_7516()) {
			this_.priv_Scanner_reportError_7524("Hex Integer Literal must contain at least one digit", []any{})
		}
		this_.priv_Scanner_skipHexDigits_7517()
		isBigInt = this_.priv_Scanner_peek_7521(uint16(110))
		if isBigInt {
			this_.priv_Scanner_nextChar_7520()
		}
		var iifeArg_3 TokenType
	iifeSwitch_3:
		switch {
		default:

			if isBigInt {
				iifeArg_3 = StaticGet_TokenType_BIGINT()
				break iifeSwitch_3
			}
			iifeArg_3 = StaticGet_TokenType_NUMBER()
			break iifeSwitch_3
		}
		return Token_Literal(iifeArg_3,

			this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	case uint16(101), uint16(69):
		return this_.priv_Scanner_scanExponentOfNumericLiteral_7513(beginToken)
	case uint16(46):
		return this_.priv_Scanner_scanFractionalNumericLiteral_7512(beginToken)
	case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
		this_.priv_Scanner_skipDecimalDigits_7514()
		if this_.priv_Scanner_peek_7521(uint16(46)) {
			this_.priv_Scanner_nextChar_7520()
			this_.priv_Scanner_skipDecimalDigits_7515(false)
		}
		if this_.priv_Scanner_peek_7521(uint16(110)) {
			this_.priv_Scanner_reportError_7524("SyntaxError: nonzero BigInt can't have leading zero", []any{})
		}
		return Token_Literal(StaticGet_TokenType_NUMBER(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	case uint16(110):
		this_.priv_Scanner_nextChar_7520()
		return Token_Literal(StaticGet_TokenType_BIGINT(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	default:
		return Token_Literal(StaticGet_TokenType_NUMBER(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.createToken(com.google.javascript.jscomp.parsing.parser.TokenType type, int beginToken)
func (this_ *Scanner) priv_Scanner_createToken_7499(type_ TokenType, beginToken int32) Token {
	_ = type_
	_ = beginToken
	return NewToken_2(type_, this_.priv_Scanner_getTokenRange_7470(beginToken))
}

// Java method: Scanner.scanIdentifierOrKeyword(int beginToken, char ch)
func (this_ *Scanner) priv_Scanner_scanIdentifierOrKeyword_7500(beginToken int32, ch uint16) Token {
	_ = beginToken
	_ = ch
	var valueStartIndex int32 = (this_.index - int32(1))
	_ = valueStartIndex
	var containsUnicodeEscape bool = (int32(ch) == int32(uint16(92)))
	_ = containsUnicodeEscape
	var bracedUnicodeEscape bool = false
	_ = bracedUnicodeEscape
	var isPrivateIdentifier bool = (int32(ch) == int32(uint16(35)))
	_ = isPrivateIdentifier
	var unicodeEscapeLen int32
	_ = unicodeEscapeLen
	if containsUnicodeEscape {
		unicodeEscapeLen = int32(1)
	} else {
		unicodeEscapeLen = int32(0)
	}
	ch = this_.priv_Scanner_peekChar_7522()
	for ((Identifiers_IsIdentifierPart(ch) || (int32(ch) == int32(uint16(92)))) || ((int32(ch) == int32(uint16(123))) && (unicodeEscapeLen == int32(2)))) || ((int32(ch) == int32(uint16(125))) && bracedUnicodeEscape) {
		if int32(ch) == int32(uint16(92)) {
			containsUnicodeEscape = true
		}
		if (int32(ch) == int32(uint16(92))) || (unicodeEscapeLen > int32(0)) {
			unicodeEscapeLen = unicodeEscapeLen + 1
		}
		if int32(ch) == int32(uint16(123)) {
			bracedUnicodeEscape = true
		}
		if (int32(ch) == int32(uint16(125))) || ((unicodeEscapeLen >= int32(6)) && !(bracedUnicodeEscape)) {
			bracedUnicodeEscape = false
			unicodeEscapeLen = int32(0)
		}
		this_.priv_Scanner_nextChar_7520()
		ch = this_.priv_Scanner_peekChar_7522()
	}
	var value string = j2gort.JStrSub(this_.contents, valueStartIndex, this_.index)
	_ = value
	if isPrivateIdentifier && (j2gort.CheckNotNullString(value) == "#") {
		this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(beginToken), "Invalid usage of #", []any{})
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ERROR(), beginToken)
	}
	if containsUnicodeEscape {
		value = Scanner_ProcessUnicodeEscapes(value)
		if value == j2gort.NullString {
			this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(this_.index), "Invalid escape sequence", []any{})
			return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ERROR(), beginToken)
		}
	}
	var start uint16 = j2gort.JStrCharAt(value, int32(0))
	_ = start
	if isPrivateIdentifier {
		start = j2gort.JStrCharAt(value, int32(1))
	}
	if !(Identifiers_IsIdentifierStart(start)) {
		this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(beginToken), "Character '%c' (U+%04X) is not a valid identifier start char", []any{j2gort.ValueOfChar(start), j2gort.ValueOfInt(int32(start))})
		return this_.priv_Scanner_createToken_7499(StaticGet_TokenType_ERROR(), beginToken)
	}
	var k *Keywords = Keywords_Get_1_string(value)
	_ = k
	if (k != nil) && !(containsUnicodeEscape) {
		return NewToken_2(k.type_, this_.priv_Scanner_getTokenRange_7470(beginToken))
	}
	return Token_Identifier(this_.priv_Scanner_getTokenRange_7470(beginToken), value)
}

// Java static method: Scanner.processUnicodeEscapes(java.lang.String value)
func Scanner_ProcessUnicodeEscapes(value string) string {
	_ = value
	for j2gort.JStrIndexOf(value, "\\") >= 0 {
		var escapeStart int32 = j2gort.JStrIndexOfCp(value, int32(uint16(92)))
		_ = escapeStart
		var didRet_1 bool
		var retVal_2 string
		func() {
			defer func() {
				if rec_3 := recover(); rec_3 != nil {
					jex_4 := javalang.RecoverException(rec_3)
					_ = jex_4
					if j2gort.InstanceOf[javalang.INumberFormatException](jex_4) || j2gort.InstanceOf[javalang.IStringIndexOutOfBoundsException](jex_4) {
						e := jex_4
						_ = e
						retVal_2 = j2gort.NullString
						didRet_1 = true
						return
					} else {
						panic(rec_3)
					}
				}
			}()
			if int32(j2gort.JStrCharAt(value, (escapeStart+int32(1)))) != int32(uint16(117)) {
				retVal_2 = j2gort.NullString
				didRet_1 = true
				return
			}
			var hexDigits string = j2gort.NullString
			_ = hexDigits
			var escapeEnd int32 = int32(0)
			_ = escapeEnd
			if int32(j2gort.JStrCharAt(value, (escapeStart+int32(2)))) != int32(uint16(123)) {
				escapeEnd = (escapeStart + int32(6))
				if escapeEnd > j2gort.JStrLen(value) {
					retVal_2 = j2gort.NullString
					didRet_1 = true
					return
				}
				hexDigits = j2gort.JStrSub(value, (escapeStart + int32(2)), escapeEnd)
			} else {
				escapeEnd = (escapeStart + int32(3))
				for Scanner_IsHexDigit(j2gort.JStrCharAt(value, escapeEnd)) {
					escapeEnd = escapeEnd + 1
				}
				if int32(j2gort.JStrCharAt(value, escapeEnd)) != int32(uint16(125)) {
					retVal_2 = j2gort.NullString
					didRet_1 = true
					return
				}
				hexDigits = j2gort.JStrSub(value, (escapeStart + int32(3)), escapeEnd)
				escapeEnd = escapeEnd + 1
			}
			var codePoint int32 = j2gort.ParseInt(hexDigits, int32(16))
			_ = codePoint
			if codePoint > int32(65535) {
				retVal_2 = j2gort.NullString
				didRet_1 = true
				return
			}
			var ch uint16 = uint16(codePoint)
			_ = ch
			if !(Identifiers_IsIdentifierPart(ch)) {
				retVal_2 = j2gort.NullString
				didRet_1 = true
				return
			}
			value = (j2gort.StrValString(j2gort.JStrSub(value, int32(0), escapeStart)) + j2gort.StrValChar(ch) + j2gort.StrValString(j2gort.JStrSubFrom(value, escapeEnd)))
		}()
		if didRet_1 {
			return retVal_2
		}
	}
	return value
}

// Java method: Scanner.scanStringLiteral(int beginIndex, char terminator)
func (this_ *Scanner) priv_Scanner_scanStringLiteral_7501(beginIndex int32, terminator uint16) Token {
	_ = beginIndex
	_ = terminator
	var startingPosition SourcePosition = this_.priv_Scanner_getPosition_7469(beginIndex)
	_ = startingPosition
	var hasUnescapedUnicodeLineOrParagraphSeparator bool = false
	_ = hasUnescapedUnicodeLineOrParagraphSeparator
	for this_.priv_Scanner_peekStringLiteralChar_7505(terminator) {
		var c uint16 = this_.priv_Scanner_peekChar_7522()
		_ = c
		hasUnescapedUnicodeLineOrParagraphSeparator = ((hasUnescapedUnicodeLineOrParagraphSeparator || (int32(c) == int32(uint16(8232)))) || (int32(c) == int32(uint16(8233))))
		if !(this_.priv_Scanner_skipStringLiteralChar_7506()) {
			return Token_StringLiteral(this_.priv_Scanner_getTokenString_7504(beginIndex), this_.priv_Scanner_getTokenRange_7471(startingPosition), hasUnescapedUnicodeLineOrParagraphSeparator)
		}
	}
	if int32(this_.priv_Scanner_peekChar_7522()) != int32(terminator) {
		this_.priv_Scanner_reportError_7525(startingPosition, "Unterminated string literal", []any{})
	} else {
		this_.priv_Scanner_nextChar_7520()
	}
	return Token_StringLiteral(this_.priv_Scanner_getTokenString_7504(beginIndex), this_.priv_Scanner_getTokenRange_7471(startingPosition), hasUnescapedUnicodeLineOrParagraphSeparator)
}

// Java method: Scanner.scanTemplateLiteral(int beginIndex)
func (this_ *Scanner) priv_Scanner_scanTemplateLiteral_7502(beginIndex int32) Token {
	_ = beginIndex
	if this_.priv_Scanner_isAtEnd_7484() {
		this_.priv_Scanner_reportError_7525(this_.priv_Scanner_getPosition_7469(beginIndex), "Unterminated template literal", []any{})
	}
	return this_.priv_Scanner_nextTemplateLiteralTokenShared_7503(StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE(), StaticGet_TokenType_TEMPLATE_HEAD())
}

// Java method: Scanner.nextTemplateLiteralTokenShared(com.google.javascript.jscomp.parsing.parser.TokenType endType, com.google.javascript.jscomp.parsing.parser.TokenType middleType)
func (this_ *Scanner) priv_Scanner_nextTemplateLiteralTokenShared_7503(endType TokenType, middleType TokenType) Token {
	_ = endType
	_ = middleType
	var beginIndex int32 = this_.index
	_ = beginIndex
	var startingPosition SourcePosition = this_.priv_Scanner_getPosition_7469(beginIndex)
	_ = startingPosition
	var skipTemplateCharactersResult *Scanner_SkipTemplateCharactersResult = this_.priv_Scanner_skipTemplateCharacters_7507()
	_ = skipTemplateCharactersResult
	if this_.priv_Scanner_isAtEnd_7484() {
		this_.priv_Scanner_reportError_7525(startingPosition, "Unterminated template literal", []any{})
	}
	var value string = this_.priv_Scanner_getTokenString_7504(beginIndex)
	_ = value
	switch this_.priv_Scanner_peekChar_7522() {
	case uint16(96):
		this_.priv_Scanner_nextChar_7520()
		return Token_TemplateLiteral(endType, value, skipTemplateCharactersResult.GetErrorMessage(), skipTemplateCharactersResult.GetErrorLevel(), skipTemplateCharactersResult.GetPosition(), this_.priv_Scanner_getTokenRange_7471(startingPosition))
	case uint16(36):
		this_.priv_Scanner_nextChar_7520()
		this_.priv_Scanner_nextChar_7520()
		return Token_TemplateLiteral(middleType, value, skipTemplateCharactersResult.GetErrorMessage(), skipTemplateCharactersResult.GetErrorLevel(), skipTemplateCharactersResult.GetPosition(), this_.priv_Scanner_getTokenRange_7471(startingPosition))
	default:
		return Token_TemplateLiteral(endType, value, skipTemplateCharactersResult.GetErrorMessage(), skipTemplateCharactersResult.GetErrorLevel(), skipTemplateCharactersResult.GetPosition(), this_.priv_Scanner_getTokenRange_7471(startingPosition))
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.getTokenString(int beginIndex)
func (this_ *Scanner) priv_Scanner_getTokenString_7504(beginIndex int32) string {
	_ = beginIndex
	return j2gort.JStrSub(this_.contents, beginIndex, this_.index)
}

// Java method: Scanner.peekStringLiteralChar(char terminator)
func (this_ *Scanner) priv_Scanner_peekStringLiteralChar_7505(terminator uint16) bool {
	_ = terminator
	return ((!(this_.priv_Scanner_isAtEnd_7484()) && (int32(this_.priv_Scanner_peekChar_7522()) != int32(terminator))) && !(Scanner_IsStringLineTerminator(this_.priv_Scanner_peekChar_7522())))
}

// Java method: Scanner.skipStringLiteralChar()
func (this_ *Scanner) priv_Scanner_skipStringLiteralChar_7506() bool {
	if this_.priv_Scanner_peek_7521(uint16(92)) {
		return this_.priv_Scanner_skipStringLiteralEscapeSequence_7509()
	}
	this_.priv_Scanner_nextChar_7520()
	return true
}

// Java method: Scanner.skipTemplateCharacters()
func (this_ *Scanner) priv_Scanner_skipTemplateCharacters_7507() *Scanner_SkipTemplateCharactersResult {
	var result *Scanner_SkipTemplateCharactersResult = this_.priv_Scanner_createSkipTemplateCharactersResult_7529(j2gort.NullString, TemplateLiteralToken_ErrorLevel(0))
	_ = result
	for !(this_.priv_Scanner_isAtEnd_7484()) {
		{
			var newError *Scanner_SkipTemplateCharactersResult = nil
			_ = newError
			switch this_.priv_Scanner_peekChar_7522() {
			case uint16(96):
				return result
			case uint16(92):
				newError = this_.priv_Scanner_skipTemplateLiteralEscapeSequence_7508()
				if newError != nil {
					if !(result.HasMessage()) || (!(result.HasError()) && newError.HasError()) {
						result = newError
					}
				}
			case uint16(36):
				if int32(this_.priv_Scanner_peekChar_7523(int32(1))) == int32(uint16(123)) {
					return result
				}
				fallthrough
			default:
				this_.priv_Scanner_nextChar_7520()
			}
		}
	}
	return result
}

// Java method: Scanner.skipTemplateLiteralEscapeSequence()
func (this_ *Scanner) priv_Scanner_skipTemplateLiteralEscapeSequence_7508() *Scanner_SkipTemplateCharactersResult {
	this_.priv_Scanner_nextChar_7520()
	if this_.priv_Scanner_isAtEnd_7484() {
		this_.priv_Scanner_reportError_7524("Unterminated template literal escape sequence", []any{})
		return nil
	}
	if Scanner_IsLineTerminator(this_.priv_Scanner_peekChar_7522()) {
		this_.priv_Scanner_skipLineTerminator_7511()
		return nil
	}
	var next uint16 = this_.priv_Scanner_nextChar_7520()
	_ = next
	switch next {
	case uint16(48):
		if Scanner_IsDecimalDigit(this_.priv_Scanner_peekChar_7522()) {
			return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Invalid escape sequence", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
		}
		return nil
	case uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
		return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Invalid escape sequence", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
	case uint16(120):
		var doubleHexDigit bool = (this_.priv_Scanner_skipHexDigit_7510() && this_.priv_Scanner_skipHexDigit_7510())
		_ = doubleHexDigit
		if !(doubleHexDigit) {
			return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Hex digit expected", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
		}
		return nil
	case uint16(117):
		if this_.priv_Scanner_peek_7521(uint16(123)) {
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(125)) {
				return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Empty unicode escape", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
			}
			var allHexDigits bool = true
			_ = allHexDigits
			for !(this_.priv_Scanner_peek_7521(uint16(125))) && allHexDigits {
				allHexDigits = (allHexDigits && this_.priv_Scanner_skipHexDigit_7510())
			}
			if !(allHexDigits) {
				return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Hex digit expected", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
			}
			this_.priv_Scanner_nextChar_7520()
			return nil
		} else {
			var quadHexDigit bool = (((this_.priv_Scanner_skipHexDigit_7510() && this_.priv_Scanner_skipHexDigit_7510()) && this_.priv_Scanner_skipHexDigit_7510()) && this_.priv_Scanner_skipHexDigit_7510())
			_ = quadHexDigit
			if !(quadHexDigit) {
				return this_.priv_Scanner_createSkipTemplateCharactersResult_7529("Hex digit expected", StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
			}
			return nil
		}
	case uint16(92), uint16(98), uint16(102), uint16(110), uint16(114), uint16(116), uint16(118), uint16(36), uint16(96):
		return nil
	case uint16(39):
		return this_.priv_Scanner_createSkipTemplateCharactersResult_7529(j2gort.Ext_com_google_common_base_Strings_LenientFormat("Unnecessary escape: \"\\%s\" is equivalent to just \"%s\"", []any{j2gort.ValueOfChar(next), j2gort.ValueOfChar(next)}), StaticGet_TemplateLiteralToken_ErrorLevel_WARNING())
	default:
		return this_.priv_Scanner_createSkipTemplateCharactersResult_7529(j2gort.Ext_com_google_common_base_Strings_LenientFormat("Unnecessary escape: '\\%s' is equivalent to just '%s'", []any{j2gort.ValueOfChar(next), j2gort.ValueOfChar(next)}), StaticGet_TemplateLiteralToken_ErrorLevel_WARNING())
	}
	panic("unreachable")
}

// Java method: Scanner.skipStringLiteralEscapeSequence()
func (this_ *Scanner) priv_Scanner_skipStringLiteralEscapeSequence_7509() bool {
	this_.priv_Scanner_nextChar_7520()
	if this_.priv_Scanner_isAtEnd_7484() {
		this_.priv_Scanner_reportError_7524("Unterminated string literal escape sequence", []any{})
		return false
	}
	if Scanner_IsLineTerminator(this_.priv_Scanner_peekChar_7522()) {
		this_.priv_Scanner_skipLineTerminator_7511()
		return true
	}
	var next uint16 = this_.priv_Scanner_nextChar_7520()
	_ = next
	switch next {
	case uint16(39), uint16(34), uint16(96), uint16(92), uint16(98), uint16(102), uint16(110), uint16(114), uint16(116), uint16(118), uint16(48):
		return true
	case uint16(120):
		var doubleHexDigit bool = (this_.priv_Scanner_skipHexDigit_7510() && this_.priv_Scanner_skipHexDigit_7510())
		_ = doubleHexDigit
		if !(doubleHexDigit) {
			this_.priv_Scanner_reportError_7524("Hex digit expected", []any{})
		}
		return doubleHexDigit
	case uint16(117):
		if this_.priv_Scanner_peek_7521(uint16(123)) {
			this_.priv_Scanner_nextChar_7520()
			if this_.priv_Scanner_peek_7521(uint16(125)) {
				this_.priv_Scanner_reportError_7524("Empty unicode escape", []any{})
				return false
			}
			var allHexDigits bool = true
			_ = allHexDigits
			for !(this_.priv_Scanner_peek_7521(uint16(125))) && allHexDigits {
				allHexDigits = (allHexDigits && this_.priv_Scanner_skipHexDigit_7510())
			}
			if !(allHexDigits) {
				this_.priv_Scanner_reportError_7524("Hex digit expected", []any{})
			}
			this_.priv_Scanner_nextChar_7520()
			return allHexDigits
		} else {
			var quadHexDigit bool = (((this_.priv_Scanner_skipHexDigit_7510() && this_.priv_Scanner_skipHexDigit_7510()) && this_.priv_Scanner_skipHexDigit_7510()) && this_.priv_Scanner_skipHexDigit_7510())
			_ = quadHexDigit
			if !(quadHexDigit) {
				this_.priv_Scanner_reportError_7524("Hex digit expected", []any{})
			}
			return quadHexDigit
		}
	default:
	}
	if int32(next) == int32(uint16(47)) {
	} else {
		this_.priv_Scanner_reportWarning_7526("Unnecessary escape: '\\%s' is equivalent to just '%s'", []any{j2gort.ValueOfChar(next), j2gort.ValueOfChar(next)})
	}
	return true
}

// Java method: Scanner.skipHexDigit()
func (this_ *Scanner) priv_Scanner_skipHexDigit_7510() bool {
	if !(this_.priv_Scanner_peekHexDigit_7516()) {
		return false
	}
	this_.priv_Scanner_nextChar_7520()
	return true
}

// Java method: Scanner.skipLineTerminator()
func (this_ *Scanner) priv_Scanner_skipLineTerminator_7511() {
	var first uint16 = this_.priv_Scanner_nextChar_7520()
	_ = first
	if (int32(first) == int32(uint16(13))) && this_.priv_Scanner_peek_7521(uint16(10)) {
		this_.priv_Scanner_nextChar_7520()
	}
}

// Java method: Scanner.scanFractionalNumericLiteral(int beginToken)
func (this_ *Scanner) priv_Scanner_scanFractionalNumericLiteral_7512(beginToken int32) Token {
	_ = beginToken
	if this_.priv_Scanner_peek_7521(uint16(46)) {
		this_.priv_Scanner_nextChar_7520()
		this_.priv_Scanner_skipDecimalDigits_7515(false)
	}
	return this_.priv_Scanner_scanExponentOfNumericLiteral_7513(beginToken)
}

// Java method: Scanner.scanExponentOfNumericLiteral(int beginToken)
func (this_ *Scanner) priv_Scanner_scanExponentOfNumericLiteral_7513(beginToken int32) Token {
	_ = beginToken
	switch this_.priv_Scanner_peekChar_7522() {
	case uint16(101), uint16(69):
		this_.priv_Scanner_nextChar_7520()
		switch this_.priv_Scanner_peekChar_7522() {
		case uint16(43), uint16(45):
			this_.priv_Scanner_nextChar_7520()
		default:
		}
		if !(Scanner_IsDecimalDigit(this_.priv_Scanner_peekChar_7522())) {
			this_.priv_Scanner_reportError_7524("Exponent part must contain at least one digit", []any{})
		}
		this_.priv_Scanner_skipDecimalDigits_7514()
	default:
	}
	return Token_Literal(StaticGet_TokenType_NUMBER(), this_.priv_Scanner_getTokenString_7504(beginToken), this_.priv_Scanner_getTokenRange_7470(beginToken))
}

// Java method: Scanner.skipDecimalDigits()
func (this_ *Scanner) priv_Scanner_skipDecimalDigits_7514() {
	this_.priv_Scanner_skipDecimalDigits_7515(true)
}

// Java method: Scanner.skipDecimalDigits(boolean hasLeadingDigit)
func (this_ *Scanner) priv_Scanner_skipDecimalDigits_7515(hasLeadingDigit bool) {
	_ = hasLeadingDigit
	var hasDigit bool = hasLeadingDigit
	_ = hasDigit
	var ch uint16 = this_.priv_Scanner_peekChar_7522()
	_ = ch
	for Scanner_IsDecimalDigit(ch) || (int32(ch) == int32(uint16(95))) {
		this_.priv_Scanner_nextChar_7520()
		if int32(ch) == int32(uint16(95)) {
			if hasDigit && Scanner_IsDecimalDigit(this_.priv_Scanner_peekChar_7522()) {
				this_.priv_Scanner_nextChar_7520()
			} else {
				this_.priv_Scanner_reportError_7524("Trailing numeric separator", []any{})
			}
		} else {
			hasDigit = true
		}
		ch = this_.priv_Scanner_peekChar_7522()
	}
}

// Java static method: Scanner.isDecimalDigit(char ch)
func Scanner_IsDecimalDigit(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.peekHexDigit()
func (this_ *Scanner) priv_Scanner_peekHexDigit_7516() bool {
	return Scanner_IsHexDigit(this_.priv_Scanner_peekChar_7522())
}

// Java static method: Scanner.isHexDigit(char ch)
func Scanner_IsHexDigit(ch uint16) bool {
	_ = ch
	return (j2gort.BoxChar_Digit_2_uint16_int32(ch, int32(16)) >= int32(0))
}

// Java method: Scanner.skipHexDigits()
func (this_ *Scanner) priv_Scanner_skipHexDigits_7517() {
	var ch uint16 = this_.priv_Scanner_peekChar_7522()
	_ = ch
	for Scanner_IsHexDigit(ch) || (int32(ch) == int32(uint16(95))) {
		this_.priv_Scanner_nextChar_7520()
		if int32(ch) == int32(uint16(95)) {
			if this_.priv_Scanner_peekHexDigit_7516() {
				this_.priv_Scanner_nextChar_7520()
			} else {
				this_.priv_Scanner_reportError_7524("Trailing numeric separator", []any{})
			}
		}
		ch = this_.priv_Scanner_peekChar_7522()
	}
}

// Java method: Scanner.skipOctalDigits()
func (this_ *Scanner) priv_Scanner_skipOctalDigits_7518() {
	var ch uint16 = this_.priv_Scanner_peekChar_7522()
	_ = ch
	for Scanner_IsOctalDigit(ch) || (int32(ch) == int32(uint16(95))) {
		this_.priv_Scanner_nextChar_7520()
		if int32(ch) == int32(uint16(95)) {
			if Scanner_IsOctalDigit(this_.priv_Scanner_peekChar_7522()) {
				this_.priv_Scanner_nextChar_7520()
			} else {
				this_.priv_Scanner_reportError_7524("Trailing numeric separator", []any{})
			}
		}
		ch = this_.priv_Scanner_peekChar_7522()
	}
}

// Java static method: Scanner.isOctalDigit(char ch)
func Scanner_IsOctalDigit(ch uint16) bool {
	_ = ch
	return (Scanner_ValueOfOctalDigit(ch) >= int32(0))
}

// Java static method: Scanner.valueOfOctalDigit(char ch)
func Scanner_ValueOfOctalDigit(ch uint16) int32 {
	_ = ch
	switch ch {
	case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55):
		return (int32(ch) - int32(uint16(48)))
	default:
		return int32(-1)
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.skipBinaryDigits()
func (this_ *Scanner) priv_Scanner_skipBinaryDigits_7519() {
	var ch uint16 = this_.priv_Scanner_peekChar_7522()
	_ = ch
	for Scanner_IsBinaryDigit(ch) || (int32(ch) == int32(uint16(95))) {
		this_.priv_Scanner_nextChar_7520()
		if int32(ch) == int32(uint16(95)) {
			if Scanner_IsBinaryDigit(this_.priv_Scanner_peekChar_7522()) {
				this_.priv_Scanner_nextChar_7520()
			} else {
				this_.priv_Scanner_reportError_7524("Trailing numeric separator", []any{})
			}
		}
		ch = this_.priv_Scanner_peekChar_7522()
	}
}

// Java static method: Scanner.isBinaryDigit(char ch)
func Scanner_IsBinaryDigit(ch uint16) bool {
	_ = ch
	return (Scanner_ValueOfBinaryDigit(ch) >= int32(0))
}

// Java static method: Scanner.valueOfBinaryDigit(char ch)
func Scanner_ValueOfBinaryDigit(ch uint16) int32 {
	_ = ch
	switch ch {
	case uint16(48):
		return int32(0)
	case uint16(49):
		return int32(1)
	default:
		return int32(-1)
	}
	panic("unreachable switch expression")
}

// Java method: Scanner.nextChar()
func (this_ *Scanner) priv_Scanner_nextChar_7520() uint16 {
	if this_.priv_Scanner_isAtEnd_7484() {
		return uint16(0)
	}
	return j2gort.JStrCharAt(this_.contents, func() int32 { old_5 := this_.index; nw_6 := old_5 + 1; this_.index = nw_6; return old_5 }())
}

// Java method: Scanner.peek(char ch)
func (this_ *Scanner) priv_Scanner_peek_7521(ch uint16) bool {
	_ = ch
	return (int32(this_.priv_Scanner_peekChar_7522()) == int32(ch))
}

// Java method: Scanner.peekChar()
func (this_ *Scanner) priv_Scanner_peekChar_7522() uint16 {
	return this_.priv_Scanner_peekChar_7523(int32(0))
}

// Java method: Scanner.peekChar(int offset)
func (this_ *Scanner) priv_Scanner_peekChar_7523(offset int32) uint16 {
	_ = offset
	if !(this_.priv_Scanner_isValidIndex_7485((this_.index + offset))) {
		return uint16(0)
	} else {
		return j2gort.JStrCharAt(this_.contents, (this_.index + offset))
	}
}

// Java method: Scanner.reportError(java.lang.String format, java.lang.Object[] arguments)
func (this_ *Scanner) priv_Scanner_reportError_7524(format string, arguments []any) {
	_ = format
	_ = arguments
	this_.priv_Scanner_reportError_7525(this_.GetPosition_0(), format, arguments)
}

// Java method: Scanner.reportError(com.google.javascript.jscomp.parsing.parser.util.SourcePosition position, java.lang.String format, java.lang.Object[] arguments)
func (this_ *Scanner) priv_Scanner_reportError_7525(position SourcePosition, format string, arguments []any) {
	_ = position
	_ = format
	_ = arguments
	this_.errorReporter.GetErrorReporterBase().ReportError_3(position, format, arguments)
}

// Java method: Scanner.reportWarning(java.lang.String format, java.lang.Object[] arguments)
func (this_ *Scanner) priv_Scanner_reportWarning_7526(format string, arguments []any) {
	_ = format
	_ = arguments
	this_.errorReporter.GetErrorReporterBase().ReportWarning_3(this_.GetPosition_0(), format, arguments)
}

// Java method: Scanner.incTypeParameterLevel()
func (this_ *Scanner) IncTypeParameterLevel() {
	this_.typeParameterLevel = this_.typeParameterLevel + 1
}

// Java method: Scanner.decTypeParameterLevel()
func (this_ *Scanner) DecTypeParameterLevel() {
	this_.typeParameterLevel = this_.typeParameterLevel - 1
}

// Java method: Scanner.createSkipTemplateCharactersResult(java.lang.@org.jspecify.annotations.Nullable String message, com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.@org.jspecify.annotations.Nullable ErrorLevel errorLevel)
func (this_ *Scanner) priv_Scanner_createSkipTemplateCharactersResult_7529(message string, errorLevel TemplateLiteralToken_ErrorLevel) *Scanner_SkipTemplateCharactersResult {
	_ = message
	_ = errorLevel
	return NewScanner_SkipTemplateCharactersResult(message, errorLevel, this_.GetPosition_0())
}

// Java interface: com.google.javascript.jscomp.parsing.parser.Scanner.CommentRecorder
type Scanner_CommentRecorder interface {
	Implements_parser_Scanner_CommentRecorder()
	IsJavaNil_() bool
	RecordComment(Comment_Type, SourceRange, string)
}

func ToIface_Scanner_CommentRecorder[T interface {
	Scanner_CommentRecorder
	comparable
}](v T) Scanner_CommentRecorder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Scanner_CommentRecorder[T Scanner_CommentRecorder](v T) Scanner_CommentRecorder {
	var vi Scanner_CommentRecorder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Scanner.CommentRecorder", "com.google.javascript.jscomp.parsing.parser.Scanner$CommentRecorder", func(v_ any) bool {
		return j2gort.InstanceOf[Scanner_CommentRecorder](v_)
	})
}

// Java class: com.google.javascript.jscomp.parsing.parser.Scanner.SkipTemplateCharactersResult
type Scanner_SkipTemplateCharactersResult struct {
	errorMessage string         `gson:"errorMessage"`
	position     SourcePosition `gson:"position"`
	idHash_      atomic.Int32
	errorLevel   TemplateLiteralToken_ErrorLevel `gson:"errorLevel"`
}

func (this_ *Scanner_SkipTemplateCharactersResult) BindSelf_(s any) {
	_ = s
}
func (this_ *Scanner_SkipTemplateCharactersResult) IsJavaNil_() bool { return this_ == nil }
func (*Scanner_SkipTemplateCharactersResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Scanner$SkipTemplateCharactersResult", "SkipTemplateCharactersResult"
}
func (this_ *Scanner_SkipTemplateCharactersResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Scanner_SkipTemplateCharactersResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Scanner.SkipTemplateCharactersResult", "com.google.javascript.jscomp.parsing.parser.Scanner$SkipTemplateCharactersResult", func(v_ any) bool {
		return j2gort.InstanceOf[*Scanner_SkipTemplateCharactersResult](v_)
	})
}

func (this_ *Scanner_SkipTemplateCharactersResult) initFields_Scanner_SkipTemplateCharactersResult() {
	_ = this_
	this_.errorMessage = j2gort.NullString
}

// Java constructor: SkipTemplateCharactersResult(java.lang.@org.jspecify.annotations.Nullable String message, com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.@org.jspecify.annotations.Nullable ErrorLevel errorLevel, com.google.javascript.jscomp.parsing.parser.util.SourcePosition position)
func NewScanner_SkipTemplateCharactersResult(message string, errorLevel TemplateLiteralToken_ErrorLevel, position SourcePosition) *Scanner_SkipTemplateCharactersResult {
	this_ := &Scanner_SkipTemplateCharactersResult{}
	InitScanner_SkipTemplateCharactersResult(this_, message, errorLevel, position)
	return this_
}

func InitScanner_SkipTemplateCharactersResult(this_ *Scanner_SkipTemplateCharactersResult, message string, errorLevel TemplateLiteralToken_ErrorLevel, position SourcePosition) {
	_ = message
	_ = errorLevel
	_ = position
	this_.initFields_Scanner_SkipTemplateCharactersResult()
	this_.errorMessage = message
	this_.errorLevel = errorLevel
	this_.position = position
}

// Java method: SkipTemplateCharactersResult.getErrorMessage()
func (this_ *Scanner_SkipTemplateCharactersResult) GetErrorMessage() string {
	return this_.errorMessage
}

// Java method: SkipTemplateCharactersResult.getErrorLevel()
func (this_ *Scanner_SkipTemplateCharactersResult) GetErrorLevel() TemplateLiteralToken_ErrorLevel {
	return this_.errorLevel
}

// Java method: SkipTemplateCharactersResult.getPosition()
func (this_ *Scanner_SkipTemplateCharactersResult) GetPosition() SourcePosition {
	return this_.position
}

// Java method: SkipTemplateCharactersResult.hasMessage()
func (this_ *Scanner_SkipTemplateCharactersResult) HasMessage() bool {
	return (this_.errorMessage != j2gort.NullString)
}

// Java method: SkipTemplateCharactersResult.hasError()
func (this_ *Scanner_SkipTemplateCharactersResult) HasError() bool {
	return (this_.errorLevel == StaticGet_TemplateLiteralToken_ErrorLevel_ERROR())
}

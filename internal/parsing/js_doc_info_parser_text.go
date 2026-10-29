// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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

package parsing

import (
	"bytes"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsDocInfoParser_ExtractionInfo is a tuple for returning both the string
// extracted and the new token following a call to any of the extract*Block
// methods.
//
// Go-only: language-adapter — a small immutable tuple, passed by value
// instead of allocated.
type JsDocInfoParser_ExtractionInfo struct {
	string_ string
	token   JsDocToken
}

// Java: JsDocInfoParser.ExtractionInfo(String, JsDocToken)
func NewJsDocInfoParser_ExtractionInfo(string_ string, token JsDocToken) JsDocInfoParser_ExtractionInfo {
	return JsDocInfoParser_ExtractionInfo{string_: string_, token: token}
}

// ExtractSingleLineBlock extracts the text found on the current line starting
// at token. Note that token = token.info; should be called after this method
// is used to update the token properly in the parser.
//
// Returns the extraction information.
//
// Java: JsDocInfoParser.extractSingleLineBlock()
func (p *JsDocInfoParser) ExtractSingleLineBlock() JsDocInfoParser_ExtractionInfo {

	// Get the current starting point.
	p.stream.Update()
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno() + 1

	line := jutil.StringTrim(p.GetRemainingJSDocLine())

	// Record the textual description.
	if len(line) > 0 {
		p.jsdocBuilder.MarkText(line, lineno, charno, lineno, charno+jutil.StringLength(line))
	}

	return NewJsDocInfoParser_ExtractionInfo(line, p.Next())
}

// Java: JsDocInfoParser.extractMultilineTextualBlock(JsDocToken)
func (p *JsDocInfoParser) ExtractMultilineTextualBlock(token JsDocToken) JsDocInfoParser_ExtractionInfo {
	return p.ExtractMultilineTextualBlockWithOption(
		token, p.GetWhitespaceOption(JsDocInfoParser_WhitespaceOption_SINGLE_LINE), false)
}

// ExtractMultilineTextualBlockWithOption extracts the text found on the
// current line and all subsequent until either an annotation, end of comment
// or end of file is reached. Note that if this method detects an end of line
// as the first token, it will quit immediately (indicating that there is no
// text where it was expected). Note that token = info.token; should be called
// after this method is used to update the token properly in the parser.
//
// token: The start token. option: How to handle whitespace.
// includeAnnotations: Whether the extracted text may include annotations. If
// set to false, text extraction will stop on the first encountered
// annotation token. Returns the extraction information.
//
// Java: JsDocInfoParser.extractMultilineTextualBlock(JsDocToken, WhitespaceOption, boolean)
func (p *JsDocInfoParser) ExtractMultilineTextualBlockWithOption(
	token JsDocToken, option JsDocInfoParser_WhitespaceOption, includeAnnotations bool) JsDocInfoParser_ExtractionInfo {
	if token == JsDocToken_EOC || token == JsDocToken_EOL || token == JsDocToken_EOF {
		return NewJsDocInfoParser_ExtractionInfo("", token)
	}
	return p.ExtractMultilineComment(token, option, true, includeAnnotations)
}

// Java: JsDocInfoParser.getWhitespaceOption(WhitespaceOption)
func (p *JsDocInfoParser) GetWhitespaceOption(defaultValue JsDocInfoParser_WhitespaceOption) JsDocInfoParser_WhitespaceOption {
	if p.preserveWhitespace {
		return JsDocInfoParser_WhitespaceOption_PRESERVE
	}
	return defaultValue
}

// JsDocInfoParser_WhitespaceOption is JsDocInfoParser.WhitespaceOption.
type JsDocInfoParser_WhitespaceOption int32

const (
	// JsDocInfoParser_WhitespaceOption_PRESERVE preserves all whitespace and
	// formatting. Needed for licenses and purposely formatted text.
	JsDocInfoParser_WhitespaceOption_PRESERVE JsDocInfoParser_WhitespaceOption = iota

	// JsDocInfoParser_WhitespaceOption_TRIM preserves newlines but trims the
	// output.
	JsDocInfoParser_WhitespaceOption_TRIM

	// JsDocInfoParser_WhitespaceOption_SINGLE_LINE removes newlines and turns
	// the output into a single line string.
	JsDocInfoParser_WhitespaceOption_SINGLE_LINE
)

var jsDocInfoParser_WhitespaceOption_names = [...]string{"PRESERVE", "TRIM", "SINGLE_LINE"}

// Name is Enum.name().
//
// Java: JsDocInfoParser.WhitespaceOption.name()
func (o JsDocInfoParser_WhitespaceOption) Name() string {
	return jsDocInfoParser_WhitespaceOption_names[o]
}

// String is Enum.toString().
//
// Java: JsDocInfoParser.WhitespaceOption.toString()
func (o JsDocInfoParser_WhitespaceOption) String() string { return o.Name() }

// ExtractBlockComment extracts the top-level block comment from the JsDoc
// comment, if any. This method differs from the extractMultilineTextualBlock
// in that it terminates under different conditions (it doesn't have the same
// prechecks), it does not first read in the remaining of the current line and
// its conditions for ignoring the "*" (STAR) are different.
//
// token: The starting token. Returns the extraction information.
//
// Java: JsDocInfoParser.extractBlockComment(JsDocToken)
func (p *JsDocInfoParser) ExtractBlockComment(token JsDocToken) JsDocInfoParser_ExtractionInfo {
	return p.ExtractMultilineComment(token, p.GetWhitespaceOption(JsDocInfoParser_WhitespaceOption_TRIM), false, false)
}

// ExtractMultilineComment extracts text from the stream until the end of the
// comment, end of the file, or an annotation token is encountered. If the
// text is being extracted for a JSDoc marker, the first line in the stream
// will always be included in the extract text.
//
// token: The starting token. option: How to handle whitespace. isMarker:
// Whether the extracted text is for a JSDoc marker or a block comment.
// includeAnnotations: Whether the extracted text may include annotations. If
// set to false, text extraction will stop on the first encountered
// annotation token. Returns the extraction information.
//
// Go-only: language-adapter — the Java StringBuilder is a bytes.Buffer of
// WTF-8: it is only appended to, and only tested for emptiness and for a
// trailing '\n', which byte length and the last byte answer exactly.
//
// Java: JsDocInfoParser.extractMultilineComment(JsDocToken, WhitespaceOption, boolean, boolean)
func (p *JsDocInfoParser) ExtractMultilineComment(
	token JsDocToken, option JsDocInfoParser_WhitespaceOption, isMarker bool, includeAnnotations bool) JsDocInfoParser_ExtractionInfo {

	var builder bytes.Buffer
	startLineno := -1
	startCharno := -1

	if isMarker {
		p.stream.Update()
		startLineno = p.stream.GetLineno()
		startCharno = p.stream.GetCharno() + 1

		line := p.GetRemainingJSDocLine()
		if option != JsDocInfoParser_WhitespaceOption_PRESERVE {
			line = jutil.StringTrim(line)
		}
		builder.WriteString(line)

		p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
		token = p.Next()
	}

	ignoreStar := false

	// Track the start of the line to count whitespace that
	// the tokenizer skipped. Because this case is rare, it's easier
	// to do this here than in the tokenizer.
	lineStartChar := -1

	for {
		switch token {
		case JsDocToken_STAR:
			if ignoreStar {
				// Mark the position after the star as the new start of the line.
				lineStartChar = p.stream.GetCharno() + 1
				ignoreStar = false
			} else {
				// The star is part of the comment.
				p.PadLine(&builder, lineStartChar, option)
				lineStartChar = -1
				builder.WriteByte('*')
			}

			token = p.Next()
			for token == JsDocToken_STAR {
				if lineStartChar != -1 {
					p.PadLine(&builder, lineStartChar, option)
					lineStartChar = -1
				}
				builder.WriteByte('*')
				token = p.Next()
			}
			continue
		case JsDocToken_EOL:
			if option != JsDocInfoParser_WhitespaceOption_SINGLE_LINE {
				builder.WriteByte('\n')
			}

			ignoreStar = true
			lineStartChar = 0
			token = p.Next()
			continue
		default:
			ignoreStar = false
			p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION

			isEOC := token == JsDocToken_EOC
			if !isEOC {
				p.PadLine(&builder, lineStartChar, option)
				lineStartChar = -1
			}

			if token == JsDocToken_EOC ||
				token == JsDocToken_EOF ||
				// When we're capturing a license block, annotations
				// in the block are OK.
				(token == JsDocToken_ANNOTATION && !includeAnnotations) {
				multilineText := builder.String()

				if option != JsDocInfoParser_WhitespaceOption_PRESERVE {
					multilineText = jutil.StringTrim(multilineText)
				}

				if isMarker && multilineText != "" {
					endLineno := p.stream.GetLineno()
					endCharno := p.stream.GetCharno()
					p.jsdocBuilder.MarkText(multilineText, startLineno, startCharno, endLineno, endCharno)
				}

				return NewJsDocInfoParser_ExtractionInfo(multilineText, token)
			}

			builder.WriteString(p.ToStringToken(token))

			line := p.GetRemainingJSDocLine()

			if option != JsDocInfoParser_WhitespaceOption_PRESERVE {
				line = JsDocInfoParser_TrimEnd(line)
			}

			builder.WriteString(line)
			token = p.Next()
		}
	}
}

// Java: JsDocInfoParser.padLine(StringBuilder, int, WhitespaceOption)
func (p *JsDocInfoParser) PadLine(builder *bytes.Buffer, lineStartChar int, option JsDocInfoParser_WhitespaceOption) {
	if lineStartChar != -1 && option == JsDocInfoParser_WhitespaceOption_PRESERVE {
		numSpaces := p.stream.GetCharno() - lineStartChar
		for i := 0; i < numSpaces; i++ {
			builder.WriteByte(' ')
		}
	} else if builder.Len() > 0 {
		if builder.Bytes()[builder.Len()-1] != '\n' || option == JsDocInfoParser_WhitespaceOption_PRESERVE {
			builder.WriteByte(' ')
		}
	}
}

// JsDocInfoParser_TrimEnd trims characters from only the end of a string.
// This method will remove all whitespace characters (defined by
// TokenUtil.isWhitespace(char), in addition to the characters provided, from
// the end of the provided string.
//
// s: String to be trimmed. Returns the string with whitespace and characters
// in extraChars removed from the end.
//
// Go-only: language-adapter — walks the WTF-8 string backwards by code
// point instead of by UTF-16 char: every char that isWhitespace accepts is a
// BMP non-surrogate, so a supplementary code point (a surrogate pair in
// Java) or a lone surrogate stops the scan exactly where Java's would.
//
// Java: JsDocInfoParser.trimEnd(String)
func JsDocInfoParser_TrimEnd(s string) string {
	end := len(s)
	for end > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:end])
		if (r == utf8.RuneError && size <= 1) || r > 0xFFFF {
			break
		}
		if rhino.TokenUtil_IsWhitespace(int(r)) {
			end -= size
		} else {
			break
		}
	}

	if end == len(s) {
		return s
	}
	return s[0:end]
}

// Java: JsDocInfoParser.eatUntilEOLIfNotAnnotation()
func (p *JsDocInfoParser) EatUntilEOLIfNotAnnotation() JsDocToken {
	return p.EatUntilEOLIfNotAnnotationFrom(p.Next())
}

// Java: JsDocInfoParser.eatUntilEOLIfNotAnnotation(JsDocToken)
func (p *JsDocInfoParser) EatUntilEOLIfNotAnnotationFrom(token JsDocToken) JsDocToken {
	if token == JsDocToken_ANNOTATION {
		p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
		return token
	}
	return p.EatTokensUntilEOLFrom(token)
}

// EatTokensUntilEOL eats tokens until JsDocToken_EOL included, and switches
// back the state to JsDocInfoParser_State_SEARCHING_ANNOTATION.
//
// Java: JsDocInfoParser.eatTokensUntilEOL()
func (p *JsDocInfoParser) EatTokensUntilEOL() JsDocToken {
	return p.EatTokensUntilEOLFrom(p.Next())
}

// EatTokensUntilEOLFrom eats tokens until JsDocToken_EOL included, and
// switches back the state to JsDocInfoParser_State_SEARCHING_ANNOTATION.
//
// Java: JsDocInfoParser.eatTokensUntilEOL(JsDocToken)
func (p *JsDocInfoParser) EatTokensUntilEOLFrom(token JsDocToken) JsDocToken {
	for {
		if token == JsDocToken_EOL || token == JsDocToken_EOC || token == JsDocToken_EOF {
			p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
			return token
		}
		token = p.Next()
	}
}

// JsDocInfoParser_NO_UNREAD_TOKEN is the specific value indicating that the
// unreadToken contains no token.
//
// Go-only: language-adapter — Java uses null; JsDocToken is a Go value
// type, so an out-of-range ordinal stands for it.
//
// Java: JsDocInfoParser.NO_UNREAD_TOKEN
const JsDocInfoParser_NO_UNREAD_TOKEN JsDocToken = -1

// RestoreLookAhead restores the lookahead token to the token stream
//
// Java: JsDocInfoParser.restoreLookAhead(JsDocToken)
func (p *JsDocInfoParser) RestoreLookAhead(token JsDocToken) {
	p.unreadToken = token
}

// Match tests whether the next symbol of the token stream matches the
// specific token.
//
// Java: JsDocInfoParser.match(JsDocToken)
func (p *JsDocInfoParser) Match(token JsDocToken) bool {
	p.unreadToken = p.Next()
	return p.unreadToken == token
}

// MatchEither tests that the next symbol of the token stream matches one of
// the specified tokens.
//
// Java: JsDocInfoParser.match(JsDocToken, JsDocToken)
func (p *JsDocInfoParser) MatchEither(token1 JsDocToken, token2 JsDocToken) bool {
	p.unreadToken = p.Next()
	return p.unreadToken == token1 || p.unreadToken == token2
}

// EatIfMatch consumes the next token and returns true if it matches
// expected.
//
// Java: JsDocInfoParser.eatIfMatch(JsDocToken)
func (p *JsDocInfoParser) EatIfMatch(token JsDocToken) bool {
	if p.Match(token) {
		p.Next()
		return true
	}
	return false
}

// Next gets the next token of the token stream or the buffered token if a
// matching was previously made.
//
// Java: JsDocInfoParser.next()
func (p *JsDocInfoParser) Next() JsDocToken {
	if p.unreadToken == JsDocInfoParser_NO_UNREAD_TOKEN {
		return p.stream.GetJsDocToken()
	} else {
		return p.Current()
	}
}

// Current gets the current token, invalidating it in the process.
//
// Java: JsDocInfoParser.current()
func (p *JsDocInfoParser) Current() JsDocToken {
	t := p.unreadToken
	p.unreadToken = JsDocInfoParser_NO_UNREAD_TOKEN
	return t
}

// SkipEOLs skips all EOLs and all empty lines in the JSDoc. Call this method
// if you want the JSDoc entry to span multiple lines.
//
// Java: JsDocInfoParser.skipEOLs()
func (p *JsDocInfoParser) SkipEOLs() {
	for p.Match(JsDocToken_EOL) {
		p.Next()
		if p.Match(JsDocToken_STAR) {
			p.Next()
		}
	}
}

// GetRemainingJSDocLine returns the remainder of the line.
//
// Java: JsDocInfoParser.getRemainingJSDocLine()
func (p *JsDocInfoParser) GetRemainingJSDocLine() string {
	result := p.stream.GetRemainingJSDocLine()
	p.unreadToken = JsDocInfoParser_NO_UNREAD_TOKEN
	return result
}

// HasParsedFileOverviewDocInfo determines whether the parser has been
// populated with docinfo with a fileoverview tag.
//
// Java: JsDocInfoParser.hasParsedFileOverviewDocInfo()
func (p *JsDocInfoParser) HasParsedFileOverviewDocInfo() bool {
	return p.jsdocBuilder.IsPopulatedWithFileOverview()
}

// Java: JsDocInfoParser.retrieveAndResetParsedJSDocInfo()
func (p *JsDocInfoParser) RetrieveAndResetParsedJSDocInfo() *rhino.JSDocInfo {
	return p.jsdocBuilder.Build()
}

// GetFileOverviewJSDocInfo gets the fileoverview JSDocInfo, if any.
//
// Java: JsDocInfoParser.getFileOverviewJSDocInfo()
func (p *JsDocInfoParser) GetFileOverviewJSDocInfo() *rhino.JSDocInfo {
	return p.fileOverviewJSDocInfo
}

// LookAheadForType looks ahead for a type annotation by advancing the
// character stream. Does not modify the token stream. This is kind of a hack,
// and is only necessary because we use the token stream to parse types, but
// need the underlying character stream to get JsDoc descriptions.
//
// Returns whether we found a type annotation.
//
// Java: JsDocInfoParser.lookAheadForType()
func (p *JsDocInfoParser) LookAheadForType() bool {
	return p.LookAheadFor('{')
}

// Java: JsDocInfoParser.lookAheadForAnnotation()
func (p *JsDocInfoParser) LookAheadForAnnotation() bool {
	return p.LookAheadFor('@')
}

// LookAheadFor looks ahead by advancing the character stream. Does not modify
// the token stream.
//
// Returns whether we found the char.
//
// Java: JsDocInfoParser.lookAheadFor(char)
func (p *JsDocInfoParser) LookAheadFor(expect uint16) bool {
	matched := false
	var c int
	for {
		c = p.stream.GetChar()
		if c == ' ' {
			continue
		} else if c == int(expect) {
			matched = true
			break
		} else {
			break
		}
	}
	p.stream.UngetChar(c)
	return matched
}

// Java: JsDocInfoParser.parseClosureUnawareCode()
func (p *JsDocInfoParser) ParseClosureUnawareCode() JsDocToken {
	// This warning that is always reported is translated into an off-by-default error,
	// and that is conditionally elevated into a build-blocking error using the relevant
	// DiagnosticGroup flags.
	p.AddParserWarning(rhino.Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID)
	if !p.Match(JsDocToken_LEFT_CURLY) {
		if !p.jsdocBuilder.RecordClosureUnawareCode() {
			p.AddParserWarning(rhino.Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA)
		}
		return p.EatUntilEOLIfNotAnnotationFrom(p.Current())
	}
	token := p.Next()
	if p.Match(JsDocToken_STRING) {
		value := p.stream.GetString()
		var mode rhino.JSDocInfo_PerFileClosureUnawareMode
		switch value {
		case "WHITESPACE":
			mode = rhino.JSDocInfo_PerFileClosureUnawareMode_WHITESPACE
		case "SIMPLE":
			mode = rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE
		default:
			p.AddParserWarningWithArg(rhino.Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE, value)
			mode = rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE
		}
		if !p.jsdocBuilder.RecordClosureUnawareCodeWithMode(mode) {
			p.AddParserWarning(rhino.Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA)
		}
	}
	return p.EatUntilEOLIfNotAnnotationFrom(token)
}

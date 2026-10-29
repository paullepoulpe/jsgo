// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20261021

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

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Destructuring (aka pattern matching); see
// http://wiki.ecmascript.org/doku.php?id=harmony:destructuring

// Parser_PatternKind is Parser.PatternKind: kinds of destructuring patterns.
type Parser_PatternKind int32

const (
	// A var, let, const; catch head; or formal parameter list--only
	// identifiers are allowed as lvalues
	Parser_PatternKind_INITIALIZER Parser_PatternKind = iota
	// An assignment or for-in initializer--any lvalue is allowed
	Parser_PatternKind_ANY
)

var parser_PatternKind_names = [...]string{"INITIALIZER", "ANY"}

// Parser_PatternKind_values is PatternKind.values().
//
// Java: Parser.PatternKind.values()
func Parser_PatternKind_values() []Parser_PatternKind {
	return []Parser_PatternKind{Parser_PatternKind_INITIALIZER, Parser_PatternKind_ANY}
}

// Name is Enum.name().
//
// Java: Parser.PatternKind.name()
func (k Parser_PatternKind) Name() string { return parser_PatternKind_names[k] }

// String is Enum.toString().
//
// Java: Parser.PatternKind.toString()
func (k Parser_PatternKind) String() string { return k.Name() }

// Java: Parser.peekPatternStart()
func (p *Parser) PeekPatternStart() bool {
	return p.Peek(TokenType_OPEN_SQUARE) || p.Peek(TokenType_OPEN_CURLY)
}

// Java: Parser.parsePattern(PatternKind)
func (p *Parser) ParsePattern(kind Parser_PatternKind) ParseTree {
	switch p.PeekType() {
	case TokenType_OPEN_SQUARE:
		return p.ParseArrayPattern(kind)
	case TokenType_OPEN_CURLY:
		fallthrough
	default:
		return p.ParseObjectPattern(kind)
	}
}

// Java: Parser.peekArrayPatternElement()
func (p *Parser) PeekArrayPatternElement() bool {
	return p.PeekExpression()
}

// Java: Parser.parseIterRest(PatternKind)
func (p *Parser) ParseIterRest(patternKind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_ELLIPSIS)
	patternAssignmentTarget := p.ParseRestAssignmentTarget(patternKind)
	return NewIterRestTree(p.GetTreeLocation(start), patternAssignmentTarget)
}

// Java: Parser.parseRestAssignmentTarget(PatternKind)
func (p *Parser) ParseRestAssignmentTarget(patternKind Parser_PatternKind) ParseTree {
	patternAssignmentTarget := p.ParsePatternAssignmentTargetNoDefault(patternKind)
	if p.Peek(TokenType_EQUAL) {
		p.ReportError("A default value cannot be specified after '...'")
	}
	return patternAssignmentTarget
}

// ParseArrayPattern parses Pattern ::= ... | "[" Element? ("," Element?)* "]"
//
// Java: Parser.parseArrayPattern(PatternKind)
func (p *Parser) ParseArrayPattern(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	elements := []ParseTree{}
	p.Eat(TokenType_OPEN_SQUARE)
	for p.Peek(TokenType_COMMA) || p.PeekArrayPatternElement() {
		if p.Peek(TokenType_COMMA) {
			nullStart := p.GetTreeStartLocation()
			p.Eat(TokenType_COMMA)
			elements = append(elements, NewNullTree(p.GetTreeLocation(nullStart)))
		} else {
			elements = append(elements, p.ParsePatternAssignmentTarget(kind))

			if p.Peek(TokenType_COMMA) {
				// Consume the comma separator
				p.Eat(TokenType_COMMA)
			} else {
				// Otherwise we must be done
				break
			}
		}
	}
	if p.Peek(TokenType_ELLIPSIS) {
		p.RecordFeatureUsed(FeatureSet_Feature_ARRAY_PATTERN_REST)
		elements = append(elements, p.ParseIterRest(kind))
	}
	if p.Eat(TokenType_CLOSE_SQUARE).IsNull() {
		// If we get no closing bracket then return invalid tree to avoid compiler tripping
		// downstream. It's needed only for IDE mode where compiler continues processing even if
		// source has syntactic errors.
		return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
	}
	return NewArrayPatternTree(p.GetTreeLocation(start), elements)
}

// ParseObjectPattern parses
// Pattern ::= "{" (Field ("," Field)* ","?)? "}" | ...
//
// Java: Parser.parseObjectPattern(PatternKind)
func (p *Parser) ParseObjectPattern(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	fields := []ParseTree{}
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekObjectPatternField() {
		fields = append(fields, p.ParseObjectPatternField(kind))

		if p.Peek(TokenType_COMMA) {
			// Consume the comma separator
			p.Eat(TokenType_COMMA)
		} else {
			// Otherwise we must be done
			break
		}
	}
	if p.Peek(TokenType_ELLIPSIS) {
		p.RecordFeatureUsed(FeatureSet_Feature_OBJECT_PATTERN_REST)
		restStart := p.GetTreeStartLocation()
		p.Eat(TokenType_ELLIPSIS)
		patternAssignmentTarget := p.ParseRestAssignmentTarget(kind)
		fields = append(fields, NewObjectRestTree(p.GetTreeLocation(restStart), patternAssignmentTarget))
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return NewObjectPatternTree(p.GetTreeLocation(start), fields)
}

// Java: Parser.peekObjectPatternField()
func (p *Parser) PeekObjectPatternField() bool {
	return p.PeekPropertyNameOrComputedProp(0)
}

// Java: Parser.parseObjectPatternField(PatternKind)
func (p *Parser) ParseObjectPatternField(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekType() == TokenType_OPEN_SQUARE {
		key := p.ParseComputedPropertyName()
		p.Eat(TokenType_COLON)
		value := p.ParsePatternAssignmentTarget(kind)
		return NewComputedPropertyDefinitionTree(p.GetTreeLocation(start), key, value)
	}

	var name Token
	if p.PeekIdOrKeyword() {
		idToken := p.EatIdOrKeywordAsId()
		if !p.Peek(TokenType_COLON) {
			if idToken.IsKeyword() {
				p.ReportError("cannot use keyword '%s' here.", idToken)
			}
			if p.Peek(TokenType_EQUAL) {
				idTree := p.newIdentifierExpressionTree(p.GetTreeLocation(start), idToken)
				p.Eat(TokenType_EQUAL)
				defaultValue := p.ParseAssignmentExpression()
				return NewDefaultParameterTree(p.GetTreeLocation(start), idTree, defaultValue)
			}
			return p.newPropertyNameAssignmentTree(p.GetTreeLocation(start), idToken, nil)
		}
		name = idToken
	} else {
		name = p.ParseLiteralExpression().LiteralToken
	}

	p.Eat(TokenType_COLON)
	value := p.ParsePatternAssignmentTarget(kind)
	return p.newPropertyNameAssignmentTree(p.GetTreeLocation(start), name, value)
}

// ParsePatternAssignmentTarget parses a PatternAssignmentTarget: the
// location where the assigned value gets stored, including an optional
// default value.
//
//	Spec AssignmentElement === PatternAssignmentTarget(PatternKind.ANY)
//	  Valid in an assignment that is not a formal parameter list or variable
//	  declaration. Sub-patterns and arbitrary left hand side expressions are
//	  allowed.
//	Spec BindingElement === PatternAssignmentElement(PatternKind.INITIALIZER)
//	  Valid in a formal parameter list or variable declaration statement. Only
//	  sub-patterns and identifiers are allowed.
//
// Examples:
//
//	[a, {foo: b = 'default'}] = someArray;          // valid
//	[x.a, {foo: x.b = 'default'}] = someArray;      // valid
//
//	let [a, {foo: b = 'default'}] = someArray;      // valid
//	let [x.a, {foo: x.b = 'default'}] = someArray;  // invalid
//
//	function f([a, {foo: b = 'default'}]) {...}     // valid
//	function f([x.a, {foo: x.b = 'default'}]) {...} // invalid
//
// Java: Parser.parsePatternAssignmentTarget(PatternKind)
func (p *Parser) ParsePatternAssignmentTarget(patternKind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	var assignmentTarget ParseTree
	assignmentTarget = p.ParsePatternAssignmentTargetNoDefault(patternKind)

	if p.Peek(TokenType_EQUAL) {
		p.Eat(TokenType_EQUAL)
		defaultValue := p.ParseAssignmentExpression()
		assignmentTarget = NewDefaultParameterTree(p.GetTreeLocation(start), assignmentTarget, defaultValue)
	}
	return assignmentTarget
}

// Java: Parser.parsePatternAssignmentTargetNoDefault(PatternKind)
func (p *Parser) ParsePatternAssignmentTargetNoDefault(kind Parser_PatternKind) ParseTree {
	var assignmentTarget ParseTree
	if p.PeekPatternStart() {
		assignmentTarget = p.ParsePattern(kind)
	} else {
		assignmentTarget = p.ParseLeftHandSideExpression()
		if !assignmentTarget.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
		}
		if kind == Parser_PatternKind_INITIALIZER &&
			assignmentTarget.GetParseTreeBase().Type != ParseTreeType_IDENTIFIER_EXPRESSION {
			// We're in the context of a formal parameter list or a variable declaration statement
			p.ReportError("Only an identifier or destructuring pattern is allowed here.")
		}
	}
	return assignmentTarget
}

// EatPossiblyImplicitSemiColon consumes a (possibly implicit) semi-colon.
// Reports an error if a semi-colon is not present.
//
// Java: Parser.eatPossiblyImplicitSemiColon()
func (p *Parser) EatPossiblyImplicitSemiColon() {
	if p.Peek(TokenType_SEMI_COLON) {
		p.Eat(TokenType_SEMI_COLON)
		return
	}
	if p.PeekImplicitSemiColon() {
		return
	}

	p.ReportError("Semi-colon expected")
}

// PeekImplicitSemiColon returns true if an implicit or explicit semi colon
// is at the current location.
//
// Java: Parser.peekImplicitSemiColon()
func (p *Parser) PeekImplicitSemiColon() bool {
	return p.PeekImplicitSemiColonAt(0)
}

// Java: Parser.peekImplicitSemiColon(int)
func (p *Parser) PeekImplicitSemiColonAt(index int) bool {
	var lineAdvanced bool
	if index == 0 {
		lineAdvanced = p.GetNextLine() > p.GetLastLine()
	} else {
		lineAdvanced = p.line(p.PeekTokenAt(index).Location.Start) > p.line(p.PeekTokenAt(index-1).Location.End)
	}
	return lineAdvanced ||
		p.PeekAt(index, TokenType_SEMI_COLON) ||
		p.PeekAt(index, TokenType_CLOSE_CURLY) ||
		p.PeekAt(index, TokenType_END_OF_FILE)
}

// GetLastLine returns the line number of the most recently consumed token.
//
// Java: Parser.getLastLine()
func (p *Parser) GetLastLine() int {
	return p.line(p.lastSourcePosition)
}

// line returns the 0-based line of pos.
//
// Go-only: positions are offsets (see SourcePosition).
func (p *Parser) line(pos SourcePosition) int {
	return p.scanner.source.Line(pos)
}

// linePos returns pos with its line and column, for error messages.
func (p *Parser) linePos(pos SourcePosition) LinePosition {
	return p.scanner.source.LinePosition(pos)
}

// GetNextLine returns the line number of the next token.
//
// Java: Parser.getNextLine()
func (p *Parser) GetNextLine() int {
	return p.line(p.PeekToken().Location.Start)
}

// EatOpt consumes the next token if it is of the expected type. Otherwise
// returns Token{}. Never reports errors.
//
// Returns the consumed token, or Token{} if the next token is not of the
// expected type.
//
// Java: Parser.eatOpt(TokenType)
func (p *Parser) EatOpt(expectedTokenType TokenType) Token {
	if p.Peek(expectedTokenType) {
		return p.Eat(expectedTokenType)
	}
	return Token{}
}

// Java: Parser.inStrictContext()
func (p *Parser) InStrictContext() bool {
	// TODO(johnlenz): track entering strict scripts/modules/functions.
	return p.config.isStrictMode
}

// Java: Parser.peekId()
func (p *Parser) PeekId() bool {
	return p.PeekIdAt(0)
}

// PeekIdAt returns whether the next token is an identifier.
//
// Java: Parser.peekId(int)
func (p *Parser) PeekIdAt(index int) bool {
	token := p.PeekTokenAt(index)
	typeToCheck := token.Type()

	if typeToCheck == TokenType_IDENTIFIER {
		if !token.IsKeyword() {
			return true
		}
		// Swap the type to the matched keyword's type
		k := Keywords_Get(token.GetMaybePrivateValue())
		typeToCheck = k.Type()
	}

	// There is one special case to handle here: outside of strict-mode code, strict-mode keywords
	// can be used as identifiers
	return !p.InStrictContext() && Keywords_IsStrictKeyword(typeToCheck)
}

// Java: Parser.peekIdOrKeyword()
func (p *Parser) PeekIdOrKeyword() bool {
	return p.PeekIdOrKeywordAt(0)
}

// Java: Parser.peekIdOrKeyword(int)
func (p *Parser) PeekIdOrKeywordAt(index int) bool {
	type_ := p.PeekTypeAt(index)
	return TokenType_IDENTIFIER == type_ || Keywords_IsKeywordTokenType(type_)
}

// EatIdOpt is shorthand for eatOpt(TokenType.IDENTIFIER).
//
// Java: Parser.eatIdOpt()
func (p *Parser) EatIdOpt() Token {
	if p.PeekId() {
		return p.EatIdOrKeywordAsId()
	}
	return Token{}
}

// EatId consumes an identifier token that is not a reserved word.
//
// See http://www.ecma-international.org/ecma-262/5.1/#sec-7.6
//
// Java: Parser.eatId()
func (p *Parser) EatId() Token {
	if p.PeekId() {
		return p.EatIdOrKeywordAsId()
	} else {
		p.ReportExpectedError(p.PeekToken(), TokenType_IDENTIFIER)
		if p.PeekIdOrKeyword() {
			return p.EatIdOrKeywordAsId()
		} else {
			return Token{}
		}
	}
}

// Java: Parser.eatObjectLiteralPropertyName()
func (p *Parser) EatObjectLiteralPropertyName() Token {
	token := p.PeekToken()
	switch token.Type() {
	case TokenType_STRING, TokenType_NUMBER, TokenType_BIGINT:
		return p.NextToken()
	case TokenType_IDENTIFIER:
		fallthrough
	default:
		return p.EatIdOrKeywordAsId()
	}
}

// EatIdOrKeywordAsId consumes an identifier token that may be a reserved
// word, i.e. an IdentifierName, not necessarily an Identifier.
//
// See http://www.ecma-international.org/ecma-262/5.1/#sec-7.6
//
// Java: Parser.eatIdOrKeywordAsId()
func (p *Parser) EatIdOrKeywordAsId() Token {
	token := p.NextToken()
	if token.Type() == TokenType_IDENTIFIER {
		return token
	} else if Keywords_IsKeywordTokenType(token.Type()) {
		kw := Keywords_GetByTokenType(token.Type())
		return newIdentifierTokenFromID(token.Location, p.scanner.internIdent(kw.Value()), tokenFlagIdentIsKeyword)
	} else {
		p.ReportExpectedError(token, TokenType_IDENTIFIER)
	}
	return Token{}
}

// Eat consumes the next token. If the consumed token is not of the expected
// type then report an error and return Token{}. Otherwise return the consumed
// token.
//
// Returns the consumed token, or Token{} if the next token is not of the
// expected type.
//
// Java: Parser.eat(TokenType)
func (p *Parser) Eat(expectedTokenType TokenType) Token {
	token := p.NextToken()
	if token.Type() != expectedTokenType {
		p.ReportExpectedError(token, expectedTokenType)
		return Token{}
	}
	return token
}

// ReportExpectedError reports a 'X' expected error message.
//
// token: the location to report the message at. expected: the thing that
// was expected.
//
// Java: Parser.reportExpectedError(Token, Object)
func (p *Parser) ReportExpectedError(token Token, expected any) {
	p.ReportErrorAtToken(token, "'%s' expected", expected)
}

// GetTreeStartLocation returns a SourcePosition for the start of a parse
// tree that starts at the current location.
//
// Java: Parser.getTreeStartLocation()
func (p *Parser) GetTreeStartLocation() SourcePosition {
	return p.PeekToken().Location.Start
}

// GetTreeEndLocation returns a SourcePosition for the end of a parse tree
// that ends at the current location.
//
// Java: Parser.getTreeEndLocation()
func (p *Parser) GetTreeEndLocation() SourcePosition {
	return p.lastSourcePosition
}

// GetTreeLocation returns a SourceRange for a parse tree that starts at
// {start} and ends at the current location.
//
// Java: Parser.getTreeLocation(SourcePosition)
func (p *Parser) GetTreeLocation(start SourcePosition) SourceRange {
	return NewSourceRange(start, p.GetTreeEndLocation())
}

// NextToken consumes the next token and returns it. Will return a never
// ending stream of TokenType.END_OF_FILE at the end of the file so callers
// don't have to check for EOF explicitly.
//
// Tokenizing is contextual. NextToken() will never return a regular
// expression literal.
//
// Java: Parser.nextToken()
func (p *Parser) NextToken() Token {
	token := p.scanner.NextToken()
	p.lastSourcePosition = token.Location.End
	return token
}

// NextRegularExpressionLiteralToken consumes a regular expression literal
// token and returns it.
//
// Java: Parser.nextRegularExpressionLiteralToken()
func (p *Parser) NextRegularExpressionLiteralToken() Token {
	token := p.scanner.NextRegularExpressionLiteralToken()
	p.lastSourcePosition = token.Location.End
	return token
}

// NextTemplateLiteralToken consumes a template literal token and returns
// it.
//
// Java: Parser.nextTemplateLiteralToken()
func (p *Parser) NextTemplateLiteralToken() Token {
	token := p.scanner.NextTemplateLiteralToken()
	p.lastSourcePosition = token.Location.End
	return token
}

// Peek returns true if the next token is of the expected type. Does not
// consume the token.
//
// Java: Parser.peek(TokenType)
func (p *Parser) Peek(expectedType TokenType) bool {
	return p.PeekAt(0, expectedType)
}

// PeekAt returns true if the index-th next token is of the expected type.
// Does not consume any tokens.
//
// Java: Parser.peek(int, TokenType)
func (p *Parser) PeekAt(index int, expectedType TokenType) bool {
	return p.PeekTypeAt(index) == expectedType
}

// PeekType returns the TokenType of the next token. Does not consume any
// tokens.
//
// Java: Parser.peekType()
func (p *Parser) PeekType() TokenType {
	return p.PeekTypeAt(0)
}

// PeekTypeAt returns the TokenType of the index-th next token. Does not
// consume any tokens.
//
// Java: Parser.peekType(int)
func (p *Parser) PeekTypeAt(index int) TokenType {
	return p.PeekTokenAt(index).Type()
}

// PeekToken returns the next token. Does not consume any tokens.
//
// Java: Parser.peekToken()
func (p *Parser) PeekToken() Token {
	return p.PeekTokenAt(0)
}

// PeekTokenAt returns the index-th next token. Does not consume any tokens.
//
// Java: Parser.peekToken(int)
func (p *Parser) PeekTokenAt(index int) Token {
	return p.scanner.PeekTokenAt(index)
}

// ReportErrorAtToken reports an error message at a given token.
//
// token: the location to report the message at. message: the message to
// report in String.format style. arguments: the arguments to fill in the
// message format.
//
// Java: Parser.reportError(Token, String, Object...)
func (p *Parser) ReportErrorAtToken(token Token, message string, arguments ...any) {
	if token.IsNull() {
		p.ReportError(message, arguments...)
	} else {
		reportErrorFormat(p.errorReporter, p.linePos(token.GetStart()), message, arguments...)
	}
}

// ReportErrorAtTree reports an error message at a given parse tree's
// location.
//
// parseTree: the location to report the message at. message: the message to
// report in String.format style. arguments: the arguments to fill in the
// message format.
//
// Java: Parser.reportError(ParseTree, String, Object...)
func (p *Parser) ReportErrorAtTree(parseTree ParseTree, message string, arguments ...any) {
	if parseTree == nil {
		p.ReportError(message, arguments...)
	} else {
		reportErrorFormat(p.errorReporter, p.linePos(parseTree.GetParseTreeBase().Location.Start), message, arguments...)
	}
}

// ReportError reports an error at the current location.
//
// message: the message to report in String.format style. arguments: the
// arguments to fill in the message format.
//
// Java: Parser.reportError(String, Object...)
func (p *Parser) ReportError(message string, arguments ...any) {
	reportErrorFormat(p.errorReporter, p.linePos(p.scanner.GetPosition()), message, arguments...)
}

// ReportErrorAtPosition reports an error at the specified location.
//
// position: the position of the error. message: the message to report in
// String.format style. arguments: the arguments to fill in the message
// format.
//
// Java: Parser.reportError(SourcePosition, String, Object...)
func (p *Parser) ReportErrorAtPosition(position SourcePosition, message string, arguments ...any) {
	reportErrorFormat(p.errorReporter, p.linePos(position), message, arguments...)
}

// Java: Parser.reportTemplateErrorIfPresent(TemplateLiteralToken)
func (p *Parser) ReportTemplateErrorIfPresent(templateToken Token) {
	if templateToken.aux == 0 {
		return
	}
	errInfo := p.scanner.TemplateError(templateToken)
	if !errInfo.HasMsg {
		return
	}
	// Java: switch (templateToken.errorLevel) — a null errorLevel would throw
	// a NullPointerException here.
	switch *errInfo.ErrorLevelPtr() {
	case TemplateLiteralToken_ErrorLevel_WARNING:
		p.errorReporter.ReportWarning(p.linePos(errInfo.Position), errInfo.Message)
		return
	case TemplateLiteralToken_ErrorLevel_ERROR:
		p.ReportErrorAtPosition(errInfo.Position, "%s", errInfo.Message)
		return
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: Parser.recordFeatureUsed(Feature)
func (p *Parser) RecordFeatureUsed(feature FeatureSet_Feature) *Parser {
	p.features = p.features.With(feature)
	return p
}

// variableDeclarationListAsParseTree widens a possibly-nil
// *VariableDeclarationListTree to ParseTree.
func variableDeclarationListAsParseTree(t *VariableDeclarationListTree) ParseTree {
	if t == nil {
		return nil
	}
	return t
}

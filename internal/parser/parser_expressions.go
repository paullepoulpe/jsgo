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

// 11.1 Primary Expressions

// Java: Parser.parsePrimaryExpression()
func (p *Parser) ParsePrimaryExpression() ParseTree {
	switch p.PeekType() {
	case TokenType_CLASS:
		return p.ParseClassExpression()
	case TokenType_SUPER:
		return p.ParseSuperExpression()
	case TokenType_THIS:
		return p.ParseThisExpression()
	case TokenType_IMPORT:
		return p.ParseDynamicImportExpression()
	case TokenType_IDENTIFIER, TokenType_TYPE, TokenType_DECLARE, TokenType_MODULE, TokenType_NAMESPACE:
		return p.ParseIdentifierExpression()
	case TokenType_NUMBER, TokenType_STRING, TokenType_BIGINT, TokenType_TRUE, TokenType_FALSE, TokenType_NULL:
		return p.ParseLiteralExpression()
	case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
		return p.ParseTemplateLiteral(nil)
	case TokenType_OPEN_SQUARE:
		return p.ParseArrayLiteral()
	case TokenType_OPEN_CURLY:
		return p.ParseObjectLiteral()
	case TokenType_OPEN_PAREN:
		return p.ParseCoverParenthesizedExpressionAndArrowParameterList()
	case TokenType_SLASH, TokenType_SLASH_EQUAL:
		return p.ParseRegularExpressionLiteral()
	default:
		return p.ParseMissingPrimaryExpression()
	}
}

// Java: Parser.parseSuperExpression()
func (p *Parser) ParseSuperExpression() *SuperExpressionTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SUPER)
	if p.Peek(TokenType_QUESTION_DOT) { // super?.() not allowed
		p.ReportError("Optional chaining is forbidden in super?.")
	}
	return NewSuperExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseThisExpression()
func (p *Parser) ParseThisExpression() *ThisExpressionTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_THIS)
	return p.newThisExpressionTree(p.GetTreeLocation(start))
}

// ParseDynamicImportExpression parses import(...).
// https://tc39.github.io/proposal-dynamic-import
//
// Java: Parser.parseDynamicImportExpression()
func (p *Parser) ParseDynamicImportExpression() *DynamicImportTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IMPORT)
	if p.Peek(TokenType_QUESTION_DOT) { // import?.() not allowed
		p.ReportError("Optional chaining is forbidden in import?.")
	}
	p.Eat(TokenType_OPEN_PAREN)
	argument := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	p.RecordFeatureUsed(FeatureSet_Feature_DYNAMIC_IMPORT)
	return NewDynamicImportTree(p.GetTreeLocation(start), argument)
}

// Java: Parser.parseIdentifierExpression()
func (p *Parser) ParseIdentifierExpression() *IdentifierExpressionTree {
	start := p.GetTreeStartLocation()
	identifier := p.EatId()
	return p.newIdentifierExpressionTree(p.GetTreeLocation(start), identifier)
}

// Java: Parser.parseLiteralExpression()
func (p *Parser) ParseLiteralExpression() *LiteralExpressionTree {
	start := p.GetTreeStartLocation()
	literal := p.NextLiteralToken()

	switch literal.Type() {
	case TokenType_STRING:
		if literal.HasUnescapedUnicodeLineOrParagraphSeparator() {
			p.RecordFeatureUsed(FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)
		}
	case TokenType_NUMBER:
		if literal.HasNumericSeparator() {
			p.RecordFeatureUsed(FeatureSet_Feature_NUMERIC_SEPARATOR)
		}
	case TokenType_BIGINT:
		p.RecordFeatureUsed(FeatureSet_Feature_BIGINT)
	}

	return p.newLiteralExpressionTree(p.GetTreeLocation(start), literal)
}

// ParseTemplateLiteral constructs a template literal expression tree.
// "operand" is used to handle the case like "foo`bar`", which is a
// CallExpression or MemberExpression that calls the function foo() with the
// template literal as the argument (with extra handling). In this case,
// operand would be "foo", which is the callsite.
//
// We store this operand in the TemplateLiteralExpressionTree and generate a
// TAGGED_TEMPLATELIT node if it's not null later when transpiling.
//
// operand: a non-null value would represent the callsite. Returns the
// template literal expression.
//
// Java: Parser.parseTemplateLiteral(ParseTree)
func (p *Parser) ParseTemplateLiteral(operand ParseTree) *TemplateLiteralExpressionTree {
	var start SourcePosition
	if operand == nil {
		start = p.GetTreeStartLocation()
	} else {
		start = operand.GetParseTreeBase().Location.Start
	}
	token := p.NextToken()
	if !token.IsTemplateLiteral() {
		p.ReportErrorAtToken(token, "Unexpected template literal token %s.", token.Type().String())
	}
	isTaggedTemplate := operand != nil
	templateToken := token
	if !isTaggedTemplate {
		p.ReportTemplateErrorIfPresent(templateToken)
	}
	elements := []ParseTree{}
	elements = append(elements, NewTemplateLiteralPortionTree(templateToken.Location, templateToken))
	if templateToken.Type() == TokenType_NO_SUBSTITUTION_TEMPLATE {
		return NewTemplateLiteralExpressionTree(p.GetTreeLocation(start), operand, elements)
	}

	// `abc${
	expression := p.ParseExpression()
	elements = append(elements, NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression))
	for !p.errorReporter.HadError() {
		templateToken = p.NextTemplateLiteralToken()
		if templateToken.Type() == TokenType_ERROR || templateToken.Type() == TokenType_END_OF_FILE {
			break
		}
		if !isTaggedTemplate {
			p.ReportTemplateErrorIfPresent(templateToken)
		}
		elements = append(elements, NewTemplateLiteralPortionTree(templateToken.Location, templateToken))
		if templateToken.Type() == TokenType_TEMPLATE_TAIL {
			break
		}

		expression = p.ParseExpression()
		elements = append(elements, NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression))
	}

	return NewTemplateLiteralExpressionTree(p.GetTreeLocation(start), operand, elements)
}

// Java: Parser.nextLiteralToken()
func (p *Parser) NextLiteralToken() Token {
	return p.NextToken()
}

// Java: Parser.parseRegularExpressionLiteral()
func (p *Parser) ParseRegularExpressionLiteral() ParseTree {
	start := p.GetTreeStartLocation()
	literal := p.NextRegularExpressionLiteralToken()
	p.RecordFeatureUsed(FeatureSet_Feature_REGEXP_SYNTAX)
	return p.newLiteralExpressionTree(p.GetTreeLocation(start), literal)
}

// 11.1.4 Array Literal Expression

// Java: Parser.parseArrayLiteral()
func (p *Parser) ParseArrayLiteral() ParseTree {
	// ArrayLiteral :
	//   [ Elisionopt ]
	//   [ ElementList ]
	//   [ ElementList , Elisionopt ]
	//
	// ElementList :
	//   Elisionopt AssignmentOrSpreadExpression
	//   ElementList , Elisionopt AssignmentOrSpreadExpression
	//
	// Elision :
	//   ,
	//   Elision ,

	start := p.GetTreeStartLocation()
	elements := []ParseTree{}

	p.Eat(TokenType_OPEN_SQUARE)
	var trailingCommaToken Token
	for p.Peek(TokenType_COMMA) || p.Peek(TokenType_ELLIPSIS) || p.PeekAssignmentExpression() {
		trailingCommaToken = Token{}
		if p.Peek(TokenType_COMMA) {
			commaStart := p.GetTreeStartLocation()
			trailingCommaToken = p.Eat(TokenType_COMMA)
			// Consider the empty element to start & end immediately before the comma token.
			elements = append(elements, NewNullTree(NewSourceRange(commaStart, commaStart)))

		} else {
			if p.Peek(TokenType_ELLIPSIS) {
				p.RecordFeatureUsed(FeatureSet_Feature_SPREAD_EXPRESSIONS)
				elements = append(elements, p.ParseIterSpread())
			} else {
				elements = append(elements, p.ParseAssignmentExpression())
			}
			if !p.Peek(TokenType_CLOSE_SQUARE) {
				trailingCommaToken = p.Eat(TokenType_COMMA)
			}
		}
	}
	p.Eat(TokenType_CLOSE_SQUARE)

	p.MaybeReportTrailingComma(trailingCommaToken)

	return NewArrayLiteralExpressionTree(
		p.GetTreeLocation(start), elements, trailingCommaToken.IsPresent())
}

// 11.1.4 Object Literal Expression

// Java: Parser.parseObjectLiteral()
func (p *Parser) ParseObjectLiteral() ParseTree {
	start := p.GetTreeStartLocation()
	base := len(p.treeStack)

	p.Eat(TokenType_OPEN_CURLY)
	var commaToken Token
	for p.Peek(TokenType_ELLIPSIS) || p.PeekPropertyNameOrComputedProp(0) || p.Peek(TokenType_STAR) {
		prop := p.ParseObjectLiteralPropertyAssignment()
		p.treeStack = append(p.treeStack, prop)
		commaToken = p.EatOpt(TokenType_COMMA)
		if commaToken.IsNull() {
			break
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)

	p.MaybeReportTrailingComma(commaToken)

	result := p.popTrees(base)
	return NewObjectLiteralExpressionTree(
		p.GetTreeLocation(start), result, commaToken.IsPresent())
}

// Java: Parser.maybeReportTrailingComma(Token)
func (p *Parser) MaybeReportTrailingComma(commaToken Token) {
	if commaToken.IsPresent() {
		p.RecordFeatureUsed(FeatureSet_Feature_TRAILING_COMMA)
		if p.config.warnTrailingCommas {
			// In ES3 mode warn about trailing commas which aren't accepted by
			// older browsers (such as IE8).
			p.errorReporter.ReportWarningFormat(
				p.linePos(commaToken.Location.Start),
				"Trailing comma is not legal in an ECMA-262 object initializer")
		}
	}
}

// Java: Parser.peekPropertyNameOrComputedProp(int)
func (p *Parser) PeekPropertyNameOrComputedProp(tokenIndex int) bool {
	return p.PeekPropertyName(tokenIndex) || p.PeekTypeAt(tokenIndex) == TokenType_OPEN_SQUARE
}

// Java: Parser.peekPropertyName(int)
func (p *Parser) PeekPropertyName(tokenIndex int) bool {
	type_ := p.PeekTypeAt(tokenIndex)
	switch type_ {
	case TokenType_IDENTIFIER, TokenType_STRING, TokenType_NUMBER, TokenType_BIGINT:
		return true
	default:
		return Keywords_IsKeywordTokenType(type_)
	}
}

// Java: Parser.parseObjectLiteralPropertyAssignment()
func (p *Parser) ParseObjectLiteralPropertyAssignment() ParseTree {
	type_ := p.PeekType()
	if type_ == TokenType_STAR {
		return p.ParseObjectLiteralPropertyAssignmentGenerator()
	} else if type_ == TokenType_ELLIPSIS {
		p.RecordFeatureUsed(FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_ELLIPSIS)
		operand := p.ParseAssignmentExpression()
		return NewObjectSpreadTree(p.GetTreeLocation(start), operand)
	} else if type_ == TokenType_STRING ||
		type_ == TokenType_NUMBER ||
		type_ == TokenType_BIGINT ||
		type_ == TokenType_IDENTIFIER ||
		Keywords_IsKeywordTokenType(type_) {
		if p.PeekGetAccessor() {
			return p.ParseObjectLiteralGetAccessor()
		} else if p.PeekSetAccessor() {
			return p.ParseObjectLiteralSetAccessor()
		} else if p.PeekAsyncMethod() {
			return p.ParseObjectLiteralAsyncMethod()
		} else if p.PeekTypeAt(1) == TokenType_OPEN_PAREN {
			return p.ParseObjectLiteralMethodDeclaration()
		} else {
			return p.ParseObjectLiteralPropertyNameAssignment()
		}
	} else if type_ == TokenType_OPEN_SQUARE {
		start := p.GetTreeStartLocation()
		name := p.ParseComputedPropertyName()

		if p.Peek(TokenType_COLON) {
			p.Eat(TokenType_COLON)
			value := p.ParseAssignmentExpression()
			return NewComputedPropertyDefinitionTree(p.GetTreeLocation(start), name, value)
		} else {
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION)
			p.ParseFunctionTail(builder, Parser_FunctionFlavor_NORMAL)
			var value ParseTree = builder.Build(p.GetTreeLocation(start))
			return NewComputedPropertyMethodTree(p.GetTreeLocation(start), name, value)
		}
	} else {
		panic(jutil.NewRuntimeException("unreachable"))
	}
}

// Java: Parser.parseObjectLiteralPropertyAssignmentGenerator()
func (p *Parser) ParseObjectLiteralPropertyAssignmentGenerator() ParseTree {
	type_ := p.PeekTypeAt(1)
	if type_ == TokenType_STRING ||
		type_ == TokenType_NUMBER ||
		type_ == TokenType_IDENTIFIER ||
		Keywords_IsKeywordTokenType(type_) {
		// parseMethodDeclaration will consume the '*'.
		return p.ParseObjectLiteralMethodDeclaration()
	} else {
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_STAR)

		name := p.ParseComputedPropertyName()
		builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION)

		p.ParseFunctionTail(builder, Parser_FunctionFlavor_GENERATOR)
		var value ParseTree = builder.Build(p.GetTreeLocation(start))
		return NewComputedPropertyMethodTree(p.GetTreeLocation(start), name, value)
	}
}

// Java: Parser.parseComputedPropertyName()
func (p *Parser) ParseComputedPropertyName() ParseTree {

	p.Eat(TokenType_OPEN_SQUARE)
	assign := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_SQUARE)
	return assign
}

// Java: Parser.peekGetAccessor()
func (p *Parser) PeekGetAccessor() bool {
	return p.PeekPredefinedString(PredefinedName_GET) && p.PeekPropertyNameOrComputedProp(1)
}

// Java: Parser.peekPredefinedString(String)
func (p *Parser) PeekPredefinedString(string_ string) bool {
	return p.PeekPredefinedStringAt(0, string_)
}

// EatPredefinedString consumes an identifier with the given value, or
// reports an error and returns Token{}.
//
// Java: Parser.eatPredefinedString(String)
func (p *Parser) EatPredefinedString(string_ string) Token {
	token := p.EatId()
	if token.IsNull() || !token.ValueEquals(string_) {
		p.ReportExpectedError(token, string_)
		return Token{}
	}
	return token
}

// Java: Parser.peekPredefinedString(int, String)
func (p *Parser) PeekPredefinedStringAt(index int, string_ string) bool {
	return p.PeekAt(index, TokenType_IDENTIFIER) &&
		p.PeekTokenAt(index).ValueEquals(string_)
}

// Java: Parser.parseObjectLiteralGetAccessor()
func (p *Parser) ParseObjectLiteralGetAccessor() ParseTree {
	return p.ParseGetAccessor(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseGetAccessor(ClassOrObjectElementInfo)
func (p *Parser) ParseGetAccessor(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(PredefinedName_GET)

	if p.PeekPropertyName(0) {
		propertyName := p.EatObjectLiteralPropertyName()
		p.Eat(TokenType_OPEN_PAREN)
		p.Eat(TokenType_CLOSE_PAREN)
		body := p.ParseFunctionBody()
		p.RecordFeatureUsed(FeatureSet_Feature_GETTER)
		return NewGetAccessorTree(
			p.GetTreeLocation(elementInfo.start),
			propertyName,
			elementInfo.isClassMember,
			elementInfo.isStatic,
			body)
	} else {
		property := p.ParseComputedPropertyName()
		p.Eat(TokenType_OPEN_PAREN)
		p.Eat(TokenType_CLOSE_PAREN)
		body := p.ParseFunctionBody()
		p.RecordFeatureUsed(FeatureSet_Feature_GETTER)
		return NewComputedPropertyGetterTree(
			p.GetTreeLocation(elementInfo.start), property, elementInfo.isStatic, body)
	}
}

// Java: Parser.peekSetAccessor()
func (p *Parser) PeekSetAccessor() bool {
	return p.PeekPredefinedString(PredefinedName_SET) && p.PeekPropertyNameOrComputedProp(1)
}

// Java: Parser.parseObjectLiteralSetAccessor()
func (p *Parser) ParseObjectLiteralSetAccessor() ParseTree {
	return p.ParseSetAccessor(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseSetAccessor(ClassOrObjectElementInfo)
func (p *Parser) ParseSetAccessor(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(PredefinedName_SET)
	if p.PeekPropertyName(0) {
		propertyName := p.EatObjectLiteralPropertyName()
		parameter := p.ParseSetterParameterList()

		body := p.ParseFunctionBody()

		p.RecordFeatureUsed(FeatureSet_Feature_SETTER)
		return NewSetAccessorTree(
			p.GetTreeLocation(elementInfo.start),
			propertyName,
			elementInfo.isClassMember,
			elementInfo.isStatic,
			parameter,
			body)
	} else {
		property := p.ParseComputedPropertyName()
		parameter := p.ParseSetterParameterList()
		body := p.ParseFunctionBody()

		p.RecordFeatureUsed(FeatureSet_Feature_SETTER)
		return NewComputedPropertySetterTree(
			p.GetTreeLocation(elementInfo.start), property, elementInfo.isStatic, parameter, body)
	}
}

// Java: Parser.parseObjectLiteralPropertyNameAssignment()
func (p *Parser) ParseObjectLiteralPropertyNameAssignment() ParseTree {
	start := p.GetTreeStartLocation()
	name := p.EatObjectLiteralPropertyName()
	colon := p.EatOpt(TokenType_COLON)
	if colon.IsNull() {
		if name.Type() != TokenType_IDENTIFIER {
			p.ReportExpectedError(p.PeekToken(), TokenType_COLON)
		} else if name.IsKeyword() {
			p.ReportErrorAtToken(name, "Cannot use keyword in short object literal")
		} else if p.Peek(TokenType_EQUAL) {
			idTree := p.newIdentifierExpressionTree(p.GetTreeLocation(start), name)
			p.Eat(TokenType_EQUAL)
			defaultValue := p.ParseAssignmentExpression()
			return NewDefaultParameterTree(p.GetTreeLocation(start), idTree, defaultValue)
		}
	}
	var value ParseTree
	if colon.IsPresent() {
		value = p.ParseAssignmentExpression()
	}
	return p.newPropertyNameAssignmentTree(p.GetTreeLocation(start), name, value)
}

// 12.2 Primary Expression
//   CoverParenthesizedExpressionAndArrowParameterList ::=
//     ( Expression )
//     ( Expression, )
//     ( )
//     ( ... BindingIdentifier )
//     ( Expression , ... BindingIdentifier )

// Java: Parser.parseCoverParenthesizedExpressionAndArrowParameterList()
func (p *Parser) ParseCoverParenthesizedExpressionAndArrowParameterList() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_PAREN)
	// Case ( )
	if p.Peek(TokenType_CLOSE_PAREN) {
		p.Eat(TokenType_CLOSE_PAREN)
		if p.Peek(TokenType_ARROW) {
			return NewFormalParameterListTree(
				p.GetTreeLocation(start),
				[]ParseTree{},
				/* hasTrailingComma= */ false,
				[]SourcePosition{})
		} else {
			p.ReportError("invalid parenthesized expression")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
		}
	}
	// Case ( ... BindingIdentifier )
	if p.Peek(TokenType_ELLIPSIS) {
		params := []ParseTree{p.ParseParameter()}
		p.Eat(TokenType_CLOSE_PAREN)
		if p.Peek(TokenType_ARROW) {
			return NewFormalParameterListTree(
				p.GetTreeLocation(start),
				params,
				/* hasTrailingComma= */ false,
				[]SourcePosition{})
		} else {
			p.ReportError("invalid parenthesized expression")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
		}
	}
	// For either of the three remaining cases:
	//     ( Expression )
	//     ( Expression, )
	//     ( Expression, ...BindingIdentifier )
	// we can parse as an expression.
	result := p.ParseExpression()
	// If it follows with a comma, we must be in either of two cases
	//     ( Expression, )
	//     ( Expression, ...BindingIdentifier )
	// case.
	if p.Peek(TokenType_COMMA) {
		if p.PeekAt(1, TokenType_CLOSE_PAREN) {
			// Create the formal parameter list here so we can record
			// the trailing comma
			p.ResetScanner(start)
			// If we fail to parse as an ArrowFunction parameter list then
			// parseFormalParameterList will take care of reporting errors.
			return p.ParseFormalParameterList()
		} else {
			p.Eat(TokenType_COMMA)
			// Since we already parsed as an expression, we will guaranteed reparse this expression
			// as an arrow function parameter list, but just leave it as a comma expression for now.
			result = NewCommaExpressionTree(
				p.GetTreeLocation(start), []ParseTree{result, p.ParseParameter()})
		}
	}
	p.Eat(TokenType_CLOSE_PAREN)
	return NewParenExpressionTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseMissingPrimaryExpression()
func (p *Parser) ParseMissingPrimaryExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.NextToken()
	p.ReportError("primary expression expected")
	return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
}

// Parser_Expression is Parser.Expression: differentiates between parsing for
// 'In' vs. 'NoIn' Variants of expression grammars.
type Parser_Expression int32

const (
	Parser_Expression_NO_IN Parser_Expression = iota
	Parser_Expression_NORMAL
)

var parser_Expression_names = [...]string{"NO_IN", "NORMAL"}

// Parser_Expression_values is Expression.values().
//
// Java: Parser.Expression.values()
func Parser_Expression_values() []Parser_Expression {
	return []Parser_Expression{Parser_Expression_NO_IN, Parser_Expression_NORMAL}
}

// Name is Enum.name().
//
// Java: Parser.Expression.name()
func (e Parser_Expression) Name() string { return parser_Expression_names[e] }

// String is Enum.toString().
//
// Java: Parser.Expression.toString()
func (e Parser_Expression) String() string { return e.Name() }

// 11.14 Expressions

// Java: Parser.parseExpressionNoIn()
func (p *Parser) ParseExpressionNoIn() ParseTree {
	return p.Parse(Parser_Expression_NO_IN)
}

// Java: Parser.parseExpression()
func (p *Parser) ParseExpression() ParseTree {
	return p.Parse(Parser_Expression_NORMAL)
}

// Java: Parser.peekExpression()
func (p *Parser) PeekExpression() bool {
	switch p.PeekType() {
	case TokenType_BANG,
		TokenType_CLASS,
		TokenType_DELETE,
		TokenType_FALSE,
		TokenType_FUNCTION,
		TokenType_IDENTIFIER,
		TokenType_TYPE,
		TokenType_DECLARE,
		TokenType_MODULE,
		TokenType_NAMESPACE,
		TokenType_MINUS,
		TokenType_MINUS_MINUS,
		TokenType_NEW,
		TokenType_NULL,
		TokenType_NUMBER,
		TokenType_BIGINT,
		TokenType_OPEN_CURLY,
		TokenType_OPEN_PAREN,
		TokenType_OPEN_SQUARE,
		TokenType_PLUS,
		TokenType_PLUS_PLUS,
		TokenType_SLASH,       // regular expression literal
		TokenType_SLASH_EQUAL, // regular expression literal
		TokenType_STRING,
		TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_SUPER,
		TokenType_THIS,
		TokenType_TILDE,
		TokenType_TRUE,
		TokenType_TYPEOF,
		TokenType_VOID,
		TokenType_YIELD:
		return true
	case TokenType_IMPORT:
		return p.PeekImportCall() || p.PeekImportDot()
	default:
		return false
	}
}

// Java: Parser.parse(Expression)
func (p *Parser) Parse(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	result := p.ParseAssignment(expressionIn)
	if p.Peek(TokenType_COMMA) && !p.PeekAt(1, TokenType_ELLIPSIS) && !p.PeekAt(1, TokenType_CLOSE_PAREN) {
		exprs := []ParseTree{}
		exprs = append(exprs, result)
		for p.Peek(TokenType_COMMA) &&
			!p.PeekAt(1, TokenType_ELLIPSIS) &&
			!p.PeekAt(1, TokenType_CLOSE_PAREN) {
			p.Eat(TokenType_COMMA)
			exprs = append(exprs, p.ParseAssignment(expressionIn))
		}
		return NewCommaExpressionTree(p.GetTreeLocation(start), exprs)
	}
	return result
}

// 12.14 Assignment operators

// Java: Parser.parseAssignmentExpression()
func (p *Parser) ParseAssignmentExpression() ParseTree {
	return p.ParseAssignment(Parser_Expression_NORMAL)
}

// Java: Parser.peekAssignmentExpression()
func (p *Parser) PeekAssignmentExpression() bool {
	return p.PeekExpression()
}

// Java: Parser.parseAssignment(Expression)
func (p *Parser) ParseAssignment(expressionIn Parser_Expression) ParseTree {
	if p.Peek(TokenType_YIELD) && p.InGeneratorContext() {
		return p.ParseYield(expressionIn)
	}

	start := p.GetTreeStartLocation()
	left := p.ParseConditional(expressionIn)

	if p.IsStartOfAsyncArrowFunction(left) {
		// re-evaluate as an async arrow function.
		p.ResetScannerToTree(left)
		return p.ParseAsyncArrowFunction(expressionIn)
	}
	if p.Peek(TokenType_ARROW) {
		return p.CompleteAssignmentExpressionParseAtArrowWithExpression(left, expressionIn)
	}

	if p.PeekAssignmentOperator() {
		if !p.Peek(TokenType_EQUAL) {
			// not the vanilla assignment operator `=`, but a special equals operator (`+=`, `-=`,
			// `**=`, etc)
			if !left.IsValidNonVanillaAssignmentTarget() {
				p.ReportError("invalid assignment target")
				return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
			}
		}
		left = p.TransformLeftHandSideExpression(left)
		if !left.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
		}
		operator := p.NextToken()
		right := p.ParseAssignment(expressionIn)
		return p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.isStartOfAsyncArrowFunction(ParseTree)
func (p *Parser) IsStartOfAsyncArrowFunction(partialExpression ParseTree) bool {
	if partialExpression.GetParseTreeBase().Type == ParseTreeType_IDENTIFIER_EXPRESSION {
		identifierToken := partialExpression.(*IdentifierExpressionTree).IdentifierToken
		// partialExpression is `async`
		// followed by `[no newline] bindingIdentifier [no newline] =>`
		return identifierToken.ValueEquals(Parser_ASYNC) &&
			!p.PeekImplicitSemiColonAt(0) &&
			p.PeekId() &&
			!p.PeekImplicitSemiColonAt(1) &&
			p.PeekAt(1, TokenType_ARROW)
	} else if partialExpression.GetParseTreeBase().Type == ParseTreeType_CALL_EXPRESSION {
		callExpression := partialExpression.(*CallExpressionTree)
		callee := callExpression.Operand
		var arguments ParseTree = callExpression.Arguments
		// partialExpression is `async [no newline] (parameters)`
		// followed by `[no newline] =>`
		return callee.GetParseTreeBase().Type == ParseTreeType_IDENTIFIER_EXPRESSION &&
			callee.(*IdentifierExpressionTree).IdentifierToken.ValueEquals(Parser_ASYNC) &&
			p.line(callee.GetParseTreeBase().Location.End) == p.line(arguments.GetParseTreeBase().Location.Start) &&
			!p.PeekImplicitSemiColon() &&
			p.Peek(TokenType_ARROW)
	} else {
		return false
	}
}

// Java: Parser.completeAssignmentExpressionParseAtArrow(ParseTree, Expression)
func (p *Parser) CompleteAssignmentExpressionParseAtArrowWithExpression(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	if leftOfArrow.GetParseTreeBase().Type == ParseTreeType_CALL_EXPRESSION {
		//   ... someAssignmentExpression // implicit semicolon
		//   (args) =>
		return p.CompleteAssignmentExpressionParseAtArrow(leftOfArrow.(*CallExpressionTree))
	} else {
		return p.CompleteArrowFunctionParseAtArrow(leftOfArrow, expressionIn)
	}
}

// Java: Parser.completeArrowFunctionParseAtArrow(ParseTree, Expression)
func (p *Parser) CompleteArrowFunctionParseAtArrow(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	arrowFormalParameters := p.TransformToArrowFormalParameters(leftOfArrow)
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed before '=>'")
	}
	p.Eat(TokenType_ARROW)
	arrowFunctionBody := p.ParseArrowFunctionBody(expressionIn, Parser_FunctionFlavor_NORMAL)

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_ARROW).
		SetFormalParameterList(arrowFormalParameters).
		SetFunctionBody(arrowFunctionBody)
	return builder.Build(p.GetTreeLocation(arrowFormalParameters.Location.Start))
}

// Java: Parser.transformToArrowFormalParameters(ParseTree)
func (p *Parser) TransformToArrowFormalParameters(leftOfArrow ParseTree) *FormalParameterListTree {
	var arrowParameterList *FormalParameterListTree
	switch leftOfArrow.GetParseTreeBase().Type {
	case ParseTreeType_FORMAL_PARAMETER_LIST:
		arrowParameterList = leftOfArrow.(*FormalParameterListTree)
	case ParseTreeType_IDENTIFIER_EXPRESSION:
		// e.g. x => x + 1
		arrowParameterList = NewFormalParameterListTree(
			leftOfArrow.GetParseTreeBase().Location,
			[]ParseTree{leftOfArrow},
			/* hasTrailingComma= */ false,
			[]SourcePosition{})
	case ParseTreeType_ARGUMENT_LIST, ParseTreeType_PAREN_EXPRESSION:
		// e.g. (x) => x + 1
		p.ResetScannerToTree(leftOfArrow)
		// If we fail to parse as an ArrowFunction parameter list then
		// parseFormalParameterList will take care of reporting errors.
		arrowParameterList = p.ParseFormalParameterList()
	default:
		p.ReportErrorAtTree(leftOfArrow, "invalid arrow function parameters")
		arrowParameterList = Parser_NewEmptyFormalParameterList(leftOfArrow.GetParseTreeBase().Location)
	}
	return arrowParameterList
}

// Java: Parser.completeAssignmentExpressionParseAtArrow(CallExpressionTree)
func (p *Parser) CompleteAssignmentExpressionParseAtArrow(callExpression *CallExpressionTree) ParseTree {
	operand := callExpression.Operand
	var arguments ParseTree = callExpression.Arguments
	var result ParseTree
	if p.line(operand.GetParseTreeBase().Location.End) < p.line(arguments.GetParseTreeBase().Location.Start) {
		// break at the implicit semicolon
		// Example:
		// foo.bar // operand and implicit semicolon
		// () => { doSomething; };
		p.ResetScannerAfter(operand)
		result = operand
	} else {
		p.ReportError("'=>' unexpected")
		result = callExpression
	}
	return result
}

// Java: Parser.parseAsyncArrowFunction(Expression)
func (p *Parser) ParseAsyncArrowFunction(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	p.EatPredefinedString(Parser_ASYNC)
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed between `async` and arrow function parameter list")
	}
	var arrowParameterList *FormalParameterListTree
	if p.Peek(TokenType_OPEN_PAREN) {
		// async (...) =>
		arrowParameterList = p.ParseFormalParameterList()
	} else {
		// async arg =>
		singleParameter := p.ParseIdentifierExpression()
		arrowParameterList = NewFormalParameterListTree(
			singleParameter.Location,
			[]ParseTree{singleParameter},
			/* hasTrailingComma= */ false,
			[]SourcePosition{})
	}
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed before '=>'")
	}
	p.Eat(TokenType_ARROW)
	arrowFunctionBody := p.ParseArrowFunctionBody(expressionIn, Parser_FunctionFlavor_ASYNCHRONOUS)

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_ARROW).
		SetAsync(true).
		SetFormalParameterList(arrowParameterList).
		SetFunctionBody(arrowFunctionBody)
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseArrowFunctionBody(Expression, FunctionFlavor)
func (p *Parser) ParseArrowFunctionBody(expressionIn Parser_Expression, functionFlavor Parser_FunctionFlavor) ParseTree {
	p.functionContextStack = append(p.functionContextStack, functionFlavor)
	var arrowFunctionBody ParseTree
	if p.Peek(TokenType_OPEN_CURLY) {
		arrowFunctionBody = p.ParseFunctionBody()
	} else {
		arrowFunctionBody = p.ParseAssignment(expressionIn)
	}
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
	return arrowFunctionBody
}

// Java: Parser.newEmptyFormalParameterList(SourceRange)
func Parser_NewEmptyFormalParameterList(location SourceRange) *FormalParameterListTree {
	return NewFormalParameterListTree(
		location,
		[]ParseTree{},
		/* hasTrailingComma= */ false,
		[]SourcePosition{})
}

// TransformLeftHandSideExpression transforms a LeftHandSideExpression into
// a LeftHandSidePattern if possible. This returns the transformed tree if it
// parses as a LeftHandSidePattern, otherwise it returns the original tree.
//
// Java: Parser.transformLeftHandSideExpression(ParseTree)
func (p *Parser) TransformLeftHandSideExpression(tree ParseTree) ParseTree {
	switch tree.GetParseTreeBase().Type {
	case ParseTreeType_ARRAY_LITERAL_EXPRESSION, ParseTreeType_OBJECT_LITERAL_EXPRESSION:
		p.ResetScannerToTree(tree)
		// If we fail to parse as an LeftHandSidePattern then
		// parseLeftHandSidePattern will take care reporting errors.
		return p.ParseLeftHandSidePattern()
	default:
		return tree
	}
}

// Java: Parser.parseLeftHandSidePattern()
func (p *Parser) ParseLeftHandSidePattern() ParseTree {
	return p.ParsePattern(Parser_PatternKind_ANY)
}

// Java: Parser.resetScanner(SourcePosition)
func (p *Parser) ResetScanner(start SourcePosition) {
	// TODO(bradfordcsmith): lastSourcePosition should really point to the end of the last token
	//     before the tree to correctly detect implicit semicolons, but it doesn't matter for the
	//     current use case.
	p.lastSourcePosition = start
	p.scanner.SetPosition(p.lastSourcePosition)
}

// Java: Parser.resetScanner(ParseTree)
func (p *Parser) ResetScannerToTree(tree ParseTree) {
	p.scanner.SetPosition(tree.GetParseTreeBase().Location.Start)
}

// Java: Parser.resetScannerAfter(ParseTree)
func (p *Parser) ResetScannerAfter(parseTree ParseTree) {
	p.lastSourcePosition = parseTree.GetParseTreeBase().Location.End
	// NOTE: The "end" position for a parseTree actually points to the first character after the
	//     last token in the tree, so this is not an off-by-one error.
	p.scanner.SetPosition(p.lastSourcePosition)
}

// Java: Parser.peekAssignmentOperator()
func (p *Parser) PeekAssignmentOperator() bool {
	switch p.PeekType() {
	case TokenType_EQUAL,
		TokenType_STAR_EQUAL,
		TokenType_STAR_STAR_EQUAL,
		TokenType_SLASH_EQUAL,
		TokenType_PERCENT_EQUAL,
		TokenType_PLUS_EQUAL,
		TokenType_MINUS_EQUAL,
		TokenType_LEFT_SHIFT_EQUAL,
		TokenType_RIGHT_SHIFT_EQUAL,
		TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL,
		TokenType_AMPERSAND_EQUAL,
		TokenType_CARET_EQUAL,
		TokenType_BAR_EQUAL,
		TokenType_OR_EQUAL,
		TokenType_AND_EQUAL,
		TokenType_QUESTION_QUESTION_EQUAL:
		return true
	default:
		return false
	}
}

// Java: Parser.inGeneratorContext()
func (p *Parser) InGeneratorContext() bool {
	// disallow yield outside of generators
	return p.functionContextStack[len(p.functionContextStack)-1].IsGenerator()
}

// ParseYield parses yield [no line terminator] (*)? AssignExpression
// https://people.mozilla.org/~jorendorff/es6-draft.html#sec-generator-function-definitions-runtime-semantics-evaluation
//
// Java: Parser.parseYield(Expression)
func (p *Parser) ParseYield(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_YIELD)
	isYieldAll := false
	var expression ParseTree
	if !p.PeekImplicitSemiColon() {
		isYieldAll = p.EatOpt(TokenType_STAR).IsPresent()
		if p.PeekAssignmentExpression() {
			expression = p.ParseAssignment(expressionIn)
		} else if isYieldAll {
			p.ReportError("yield* requires an expression")
		}
	}
	return NewYieldExpressionTree(p.GetTreeLocation(start), isYieldAll, expression)
}

// 11.12 Conditional Expression

// Java: Parser.parseConditional(Expression)
func (p *Parser) ParseConditional(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	condition := p.ParseShortCircuit(expressionIn)
	if p.Peek(TokenType_QUESTION) {
		p.Eat(TokenType_QUESTION)
		left := p.ParseAssignment(expressionIn)
		p.Eat(TokenType_COLON)
		right := p.ParseAssignment(expressionIn)
		return NewConditionalExpressionTree(p.GetTreeLocation(start), condition, left, right)
	}
	return condition
}

// Java: Parser.parseShortCircuit(Expression)
func (p *Parser) ParseShortCircuit(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseLogicalOR(expressionIn)
	if p.Peek(TokenType_QUESTION_QUESTION) {
		if left.GetParseTreeBase().Type == ParseTreeType_BINARY_OPERATOR {
			binaryTree := left.(*BinaryOperatorTree)
			if binaryTree.Operator.Type() == TokenType_AND || binaryTree.Operator.Type() == TokenType_OR {
				p.ReportError("Logical OR and logical AND require parentheses when used with '??'")
			}
		}
		return p.ParseNullishCoalesce(expressionIn, left, start)
	} else {
		return left
	}
}

// Java: Parser.parseNullishCoalesce(Expression, ParseTree, SourcePosition)
func (p *Parser) ParseNullishCoalesce(expressionIn Parser_Expression, left ParseTree, start SourcePosition) ParseTree {
	for p.Peek(TokenType_QUESTION_QUESTION) {
		operator := p.Eat(TokenType_QUESTION_QUESTION)
		right := p.ParseBitwiseOR(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	if p.Peek(TokenType_AND) || p.Peek(TokenType_OR) {
		p.ReportError("Logical OR and logical AND require parentheses when used with '??'")
	}
	return left
}

// 11.11 Logical OR

// Java: Parser.parseLogicalOR(Expression)
func (p *Parser) ParseLogicalOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseLogicalAND(expressionIn)
	for p.Peek(TokenType_OR) {
		operator := p.Eat(TokenType_OR)
		right := p.ParseLogicalAND(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.11 Logical AND

// Java: Parser.parseLogicalAND(Expression)
func (p *Parser) ParseLogicalAND(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseOR(expressionIn)
	for p.Peek(TokenType_AND) {
		operator := p.Eat(TokenType_AND)
		right := p.ParseBitwiseOR(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise OR

// Java: Parser.parseBitwiseOR(Expression)
func (p *Parser) ParseBitwiseOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseXOR(expressionIn)
	for p.Peek(TokenType_BAR) {
		operator := p.Eat(TokenType_BAR)
		right := p.ParseBitwiseXOR(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise XOR

// Java: Parser.parseBitwiseXOR(Expression)
func (p *Parser) ParseBitwiseXOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseAND(expressionIn)
	for p.Peek(TokenType_CARET) {
		operator := p.Eat(TokenType_CARET)
		right := p.ParseBitwiseAND(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise AND

// Java: Parser.parseBitwiseAND(Expression)
func (p *Parser) ParseBitwiseAND(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseEquality(expressionIn)
	for p.Peek(TokenType_AMPERSAND) {
		operator := p.Eat(TokenType_AMPERSAND)
		right := p.ParseEquality(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.9 Equality Expression

// Java: Parser.parseEquality(Expression)
func (p *Parser) ParseEquality(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseRelational(expressionIn)
	for p.PeekEqualityOperator() {
		operator := p.NextToken()
		right := p.ParseRelational(expressionIn)
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekEqualityOperator()
func (p *Parser) PeekEqualityOperator() bool {
	switch p.PeekType() {
	case TokenType_EQUAL_EQUAL, TokenType_NOT_EQUAL, TokenType_EQUAL_EQUAL_EQUAL, TokenType_NOT_EQUAL_EQUAL:
		return true
	default:
		return false
	}
}

// 11.8 Relational

// Java: Parser.parseRelational(Expression)
func (p *Parser) ParseRelational(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseShiftExpression()
	for p.PeekRelationalOperator(expressionIn) {
		operator := p.NextToken()
		right := p.ParseShiftExpression()
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekRelationalOperator(Expression)
func (p *Parser) PeekRelationalOperator(expressionIn Parser_Expression) bool {
	switch p.PeekType() {
	case TokenType_OPEN_ANGLE, TokenType_CLOSE_ANGLE, TokenType_GREATER_EQUAL, TokenType_LESS_EQUAL, TokenType_INSTANCEOF:
		return true
	case TokenType_IN:
		return expressionIn == Parser_Expression_NORMAL
	default:
		return false
	}
}

// 11.7 Shift Expression

// Java: Parser.parseShiftExpression()
func (p *Parser) ParseShiftExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseAdditiveExpression()
	for p.PeekShiftOperator() {
		operator := p.NextToken()
		right := p.ParseAdditiveExpression()
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekShiftOperator()
func (p *Parser) PeekShiftOperator() bool {
	switch p.PeekType() {
	case TokenType_LEFT_SHIFT, TokenType_RIGHT_SHIFT, TokenType_UNSIGNED_RIGHT_SHIFT:
		return true
	default:
		return false
	}
}

// 11.6 Additive Expression

// Java: Parser.parseAdditiveExpression()
func (p *Parser) ParseAdditiveExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseMultiplicativeExpression()
	for p.PeekAdditiveOperator() {
		operator := p.NextToken()
		right := p.ParseMultiplicativeExpression()
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekAdditiveOperator()
func (p *Parser) PeekAdditiveOperator() bool {
	switch p.PeekType() {
	case TokenType_PLUS, TokenType_MINUS:
		return true
	default:
		return false
	}
}

// 11.5 Multiplicative Expression

// Java: Parser.parseMultiplicativeExpression()
func (p *Parser) ParseMultiplicativeExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseExponentiationExpression()
	for p.PeekMultiplicativeOperator() {
		operator := p.NextToken()
		right := p.ParseExponentiationExpression()
		left = p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekMultiplicativeOperator()
func (p *Parser) PeekMultiplicativeOperator() bool {
	switch p.PeekType() {
	case TokenType_STAR, TokenType_SLASH, TokenType_PERCENT:
		return true
	default:
		return false
	}
}

// Java: Parser.parseExponentiationExpression()
func (p *Parser) ParseExponentiationExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseUnaryExpression()
	if p.Peek(TokenType_STAR_STAR) {
		// ExponentiationExpression does not allow a UnaryExpression before '**'.
		// Parentheses are required to disambiguate:
		//   (-x)**y is valid
		//   -(x**y) is valid
		//   -x**y is a syntax error
		if left.GetParseTreeBase().Type == ParseTreeType_UNARY_EXPRESSION {
			p.ReportError(
				"Unary operator '%s' requires parentheses before '**'",
				left.(*UnaryExpressionTree).Operator)
		} else if left.GetParseTreeBase().Type == ParseTreeType_AWAIT_EXPRESSION {
			p.ReportError("Unary operator 'await' requires parentheses before '**'")
		}
		operator := p.NextToken()
		right := p.ParseExponentiationExpression()
		return p.newBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	} else {
		return left
	}
}

// 11.4 Unary Operator

// Java: Parser.parseUnaryExpression()
func (p *Parser) ParseUnaryExpression() ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekUnaryOperator() {
		operator := p.NextToken()
		operand := p.ParseUnaryExpression()
		return NewUnaryExpressionTree(p.GetTreeLocation(start), operator, operand)
	} else if p.PeekAwaitExpression() {
		return p.ParseAwaitExpression()
	} else {
		return p.ParseUpdateExpression()
	}
}

// Java: Parser.peekUnaryOperator()
func (p *Parser) PeekUnaryOperator() bool {
	switch p.PeekType() {
	case TokenType_DELETE, TokenType_VOID, TokenType_TYPEOF, TokenType_PLUS, TokenType_MINUS, TokenType_TILDE, TokenType_BANG:
		return true
	default:
		return false
	}
}

// Parser_AWAIT is Parser.AWAIT.
const Parser_AWAIT = "await"

// Java: Parser.peekAwaitExpression()
func (p *Parser) PeekAwaitExpression() bool {
	return p.PeekPredefinedString(Parser_AWAIT)
}

// Java: Parser.parseAwaitExpression()
func (p *Parser) ParseAwaitExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatPredefinedString(Parser_AWAIT)
	expression := p.ParseUnaryExpression()
	return NewAwaitExpressionTree(p.GetTreeLocation(start), expression)
}

// Java: Parser.parseUpdateExpression()
func (p *Parser) ParseUpdateExpression() ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekUpdateOperator() {
		operator := p.NextToken()
		operand := p.ParseUnaryExpression()
		return UpdateExpressionTree_Prefix(p.GetTreeLocation(start), operator, operand)
	} else {
		lhs := p.ParseLeftHandSideExpression()
		if p.PeekUpdateOperator() && !p.PeekImplicitSemiColon() {
			// newline not allowed before an update operator.
			operator := p.NextToken()
			return UpdateExpressionTree_Postfix(p.GetTreeLocation(start), operator, lhs)
		} else {
			return lhs
		}
	}
}

// Java: Parser.peekUpdateOperator()
func (p *Parser) PeekUpdateOperator() bool {
	switch p.PeekType() {
	case TokenType_PLUS_PLUS, TokenType_MINUS_MINUS:
		return true
	default:
		return false
	}
}

// Java: Parser.peekImportCall()
func (p *Parser) PeekImportCall() bool {
	return p.Peek(TokenType_IMPORT) && p.PeekAt(1, TokenType_OPEN_PAREN)
}

// Java: Parser.peekImportDot()
func (p *Parser) PeekImportDot() bool {
	return p.Peek(TokenType_IMPORT) && p.PeekAt(1, TokenType_PERIOD)
}

// ParseLeftHandSideExpression parses a LeftHandSideExpression.
//
// Java: Parser.parseLeftHandSideExpression()
func (p *Parser) ParseLeftHandSideExpression() ParseTree {
	start := p.GetTreeStartLocation()
	// We have these possible productions.
	// LeftHandSideExpression -> NewExpression
	//                        -> CallExpression
	//                        -> MemberExpression
	//                        -> OptionalExpression
	//
	// NewExpression -> new NewExpression
	//               -> MemberExpression
	//
	// CallExpression -> MemberExpression Arguments
	//                -> CallExpression ... see below
	//
	// OptionalExpression -> MemberExpression OptionalChain
	//                    -> CallExpression OptionalChain
	//                    -> OptionalExpression OptionalChain
	//
	// We try parsing a NewExpression, here, because that will include parsing MemberExpression.
	// If what we really have is a CallExpression or OptionalExpression, then the MemberExpression
	// we get back from parseNewExpression will be the first part of it, and we'll build the
	// rest later.
	operand := p.ParseNewExpression()

	// this test is equivalent to is member expression
	if newExpressionTree, ok := operand.(*NewExpressionTree); !ok || newExpressionTree.Arguments != nil {
		// We have a MemberExpression, but it may actually be just the first part of a CallExpression
		// Attempt to gather the rest of the CallExpression, if so.
		for p.PeekCallSuffix() {
			switch p.PeekType() {
			case TokenType_OPEN_PAREN:
				arguments := p.ParseArguments()
				operand = p.newCallExpressionTree(p.GetTreeLocation(start), operand, arguments)
			case TokenType_OPEN_SQUARE:
				p.Eat(TokenType_OPEN_SQUARE)
				member := p.ParseExpression()
				p.Eat(TokenType_CLOSE_SQUARE)
				operand = p.newMemberLookupExpressionTree(p.GetTreeLocation(start), operand, member)
			case TokenType_PERIOD:
				p.Eat(TokenType_PERIOD)
				id := p.EatIdOrKeywordAsId()
				operand = p.newMemberExpressionTree(p.GetTreeLocation(start), operand, id)
			case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
				operand = p.ParseTemplateLiteral(operand)
			default:
				panic(jutil.NewAssertionError("unexpected case: " + p.PeekType().String()))
			}
		}
		operand = p.MaybeParseOptionalExpression(operand)
	}
	return operand
}

// Java: Parser.peekCallSuffix()
func (p *Parser) PeekCallSuffix() bool {
	return p.Peek(TokenType_OPEN_PAREN) ||
		p.Peek(TokenType_OPEN_SQUARE) ||
		p.Peek(TokenType_PERIOD) ||
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) ||
		p.Peek(TokenType_TEMPLATE_HEAD)
}

// MaybeParseOptionalExpression tries to parse the expression as an optional
// expression.
//
// `operand?.identifier` or `operand?.[expression]` or `operand?.(arg1, arg2)`
//
// returns parse tree after trying to parse it as an optional expression
//
// Java: Parser.maybeParseOptionalExpression(ParseTree)
func (p *Parser) MaybeParseOptionalExpression(operand ParseTree) ParseTree {
	// The optional chain's source info should cover the lhs operand also
	start := operand.GetParseTreeBase().Location.Start

	for p.Peek(TokenType_QUESTION_DOT) {
		p.Eat(TokenType_QUESTION_DOT)
		switch p.PeekType() {
		case TokenType_OPEN_PAREN:
			arguments := p.ParseArguments()
			operand = NewOptChainCallExpressionTree(
				p.GetTreeLocation(start),
				operand,
				arguments,
				/* isStartOfOptionalChain= */ true,
				arguments.HasTrailingComma)
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			operand = NewOptionalMemberLookupExpressionTree(
				p.GetTreeLocation(start), operand, member /* isStartOfOptionalChain= */, true)
		case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
			p.ReportError("template literal cannot be used within optional chaining")
		default:
			if p.PeekIdOrKeyword() {
				id := p.EatIdOrKeywordAsId()
				operand = NewOptionalMemberExpressionTree(
					p.GetTreeLocation(start), operand, id /* isStartOfOptionalChain= */, true)
			} else {
				p.ReportError("syntax error: %s not allowed in optional chain", p.PeekType())
			}
		}
		operand = p.ParseRemainingOptionalChainSegment(operand)
	}
	return operand
}

// ParseRemainingOptionalChainSegment parses the remaining components of an
// optional chain till the current chain's end, or a new chain's start.
//
// `optionalExpression.identifier`, `optionalExpression[expression]`,
// `optionalExpression(arg1, arg2)`, or
// `optionalExpression?.optionalExpression`
//
// returns parse tree after trying to parse it as an optional chain
//
// Java: Parser.parseRemainingOptionalChainSegment(ParseTree)
func (p *Parser) ParseRemainingOptionalChainSegment(optionalExpression ParseTree) ParseTree {
	// The optional chain's source info should cover the lhs operand also
	start := optionalExpression.GetParseTreeBase().Location.Start
	for p.PeekOptionalChainSuffix() {
		if p.PeekType() == TokenType_NO_SUBSTITUTION_TEMPLATE ||
			p.PeekType() == TokenType_TEMPLATE_HEAD {
			p.ReportError("template literal cannot be used within optional chaining")
			break
		}
		switch p.PeekType() {
		case TokenType_PERIOD:
			p.Eat(TokenType_PERIOD)
			id := p.EatIdOrKeywordAsId()
			optionalExpression = NewOptionalMemberExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				id,
				/* isStartOfOptionalChain= */ false)
		case TokenType_OPEN_PAREN:
			arguments := p.ParseArguments()
			optionalExpression = NewOptChainCallExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				arguments,
				/* isStartOfOptionalChain= */ false,
				arguments.HasTrailingComma)
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			optionalExpression = NewOptionalMemberLookupExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				member,
				/* isStartOfOptionalChain= */ false)
		default:
			panic(jutil.NewAssertionError("unexpected case: " + p.PeekType().String()))
		}
	}
	return optionalExpression
}

// PeekOptionalChainSuffix: tokens that indicate continuation of an optional
// chain.
//
// Java: Parser.peekOptionalChainSuffix()
func (p *Parser) PeekOptionalChainSuffix() bool {
	return p.Peek(TokenType_OPEN_PAREN) || // a?.b( ...
		p.Peek(TokenType_OPEN_SQUARE) || // a?.b[ ...
		p.Peek(TokenType_PERIOD) || // a?.b. ...
		// TEMPLATE_HEAD and NO_SUBSTITUTION_TEMPLATE are actually not allowed within optional
		// chaining and leads to an early error as dictated by the spec.
		// https://tc39.es/proposal-optional-chaining/#sec-left-hand-side-expressions-static-semantics-early-errors
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) || // a?.b`text`
		p.Peek(TokenType_TEMPLATE_HEAD) // a?.b`text ${substitution} text`
}

// Parser_ASYNC is Parser.ASYNC.
const Parser_ASYNC = "async"

// 11.2 Member Expression without the new production

// Java: Parser.parseMemberExpressionNoNew()
func (p *Parser) ParseMemberExpressionNoNew() ParseTree {
	start := p.GetTreeStartLocation()
	var operand ParseTree
	if p.PeekImportDot() {
		operand = p.ParseImportDotMeta()
	} else if p.PeekAsyncFunctionStart() {
		operand = p.ParseAsyncFunctionExpression()
	} else if p.PeekFunction() {
		operand = p.ParseFunctionExpression()
	} else {
		operand = p.ParsePrimaryExpression()
	}
	for p.PeekMemberExpressionSuffix() {
		switch p.PeekType() {
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			operand = p.newMemberLookupExpressionTree(p.GetTreeLocation(start), operand, member)
		case TokenType_PERIOD:
			p.Eat(TokenType_PERIOD)
			id := p.EatIdOrKeywordAsId()
			operand = p.newMemberExpressionTree(p.GetTreeLocation(start), operand, id)
		case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
			operand = p.ParseTemplateLiteral(operand)
		default:
			panic(jutil.NewRuntimeException("unreachable"))
		}
	}
	return operand
}

// Java: Parser.peekMemberExpressionSuffix()
func (p *Parser) PeekMemberExpressionSuffix() bool {
	return p.Peek(TokenType_OPEN_SQUARE) ||
		p.Peek(TokenType_PERIOD) ||
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) ||
		p.Peek(TokenType_TEMPLATE_HEAD)
}

// Java: Parser.parseNewExpression()
func (p *Parser) ParseNewExpression() ParseTree {
	if !p.Peek(TokenType_NEW) {
		return p.ParseMemberExpressionNoNew()
	} else if p.PeekAt(1, TokenType_PERIOD) {
		return p.ParseNewDotSomething()
	} else {
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_NEW)
		if p.Peek(TokenType_QUESTION_DOT) { // new?.target not allowed
			p.ReportError("Optional chaining is forbidden in `new?.target` contexts.")
		}
		operand := p.ParseNewExpression()
		if p.Peek(TokenType_QUESTION_DOT) { // new a?.() not allowed
			p.ReportError("Optional chaining is forbidden in construction contexts.")
		}
		var arguments *ArgumentListTree
		if p.Peek(TokenType_OPEN_PAREN) {
			arguments = p.ParseArguments()
		}
		return NewNewExpressionTree(
			p.GetTreeLocation(start),
			operand,
			arguments,
			arguments != nil && arguments.HasTrailingComma)
	}
}

// Java: Parser.parseNewDotSomething()
func (p *Parser) ParseNewDotSomething() ParseTree {
	// currently only "target" is valid after "new."
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_NEW)
	p.Eat(TokenType_PERIOD)
	p.EatPredefinedString("target")
	return NewNewTargetExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseImportDotMeta()
func (p *Parser) ParseImportDotMeta() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IMPORT)
	p.Eat(TokenType_PERIOD)
	p.EatPredefinedString("meta")
	return NewImportMetaExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseArguments()
func (p *Parser) ParseArguments() *ArgumentListTree {
	// ArgumentList :
	//   AssignmentOrSpreadExpression
	//   ArgumentList , AssignmentOrSpreadExpression
	//
	// AssignmentOrSpreadExpression :
	//   ... AssignmentExpression
	//   AssignmentExpression

	start := p.GetTreeStartLocation()
	treeBase := len(p.treeStack)
	trailingComma := false
	posBase := len(p.posStack)

	p.Eat(TokenType_OPEN_PAREN)
	for p.PeekAssignmentOrSpread() {
		arg := p.ParseAssignmentOrSpread()
		p.treeStack = append(p.treeStack, arg)

		if !p.Peek(TokenType_CLOSE_PAREN) {
			comma := p.Eat(TokenType_COMMA)
			if comma.IsPresent() {
				p.posStack = append(p.posStack, comma.GetStart())
			}
			if p.Peek(TokenType_CLOSE_PAREN) {
				if !p.config.atLeast8 {
					p.ReportErrorAtToken(comma, "Invalid trailing comma in arguments list")
				}
				trailingComma = true
			}
		}
	}
	p.Eat(TokenType_CLOSE_PAREN)
	arguments := p.popTrees(treeBase)
	commaPositions := popStackChunked(&p.posStack, posBase, &p.posLists, parser_posListChunk)
	return p.newArgumentListTree(
		p.GetTreeLocation(start), arguments, trailingComma, commaPositions)
}

// PeekAssignmentOrSpread reports whether we have a spread expression or an
// assignment next.
//
// This does not peek the operand for the spread expression. This means that
// ParseAssignmentOrSpread might still fail when this returns true.
//
// Java: Parser.peekAssignmentOrSpread()
func (p *Parser) PeekAssignmentOrSpread() bool {
	return p.Peek(TokenType_ELLIPSIS) || p.PeekAssignmentExpression()
}

// Java: Parser.parseAssignmentOrSpread()
func (p *Parser) ParseAssignmentOrSpread() ParseTree {
	if p.Peek(TokenType_ELLIPSIS) {
		return p.ParseIterSpread()
	}
	return p.ParseAssignmentExpression()
}

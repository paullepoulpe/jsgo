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
	portion := NewTemplateLiteralPortionTree(templateToken.Location, templateToken)
	if templateToken.Type() == TokenType_NO_SUBSTITUTION_TEMPLATE {
		return NewTemplateLiteralExpressionTree(p.GetTreeLocation(start), operand, p.singleTreeList(portion))
	}

	base := len(p.treeStack)
	p.treeStack = append(p.treeStack, portion)

	// `abc${
	expression := p.ParseExpression()
	sub := NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression)
	p.treeStack = append(p.treeStack, sub)
	for !p.errorReporter.HadError() {
		templateToken = p.NextTemplateLiteralToken()
		if templateToken.Type() == TokenType_ERROR || templateToken.Type() == TokenType_END_OF_FILE {
			break
		}
		if !isTaggedTemplate {
			p.ReportTemplateErrorIfPresent(templateToken)
		}
		portion = NewTemplateLiteralPortionTree(templateToken.Location, templateToken)
		p.treeStack = append(p.treeStack, portion)
		if templateToken.Type() == TokenType_TEMPLATE_TAIL {
			break
		}

		expression = p.ParseExpression()
		sub = NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression)
		p.treeStack = append(p.treeStack, sub)
	}

	elements := p.popTrees(base)
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
	base := len(p.treeStack)

	p.Eat(TokenType_OPEN_SQUARE)
	var trailingCommaToken Token
	for p.Peek(TokenType_COMMA) || p.Peek(TokenType_ELLIPSIS) || p.PeekAssignmentExpression() {
		trailingCommaToken = Token{}
		if p.Peek(TokenType_COMMA) {
			commaStart := p.GetTreeStartLocation()
			trailingCommaToken = p.Eat(TokenType_COMMA)
			// Consider the empty element to start & end immediately before the comma token.
			nullTree := NewNullTree(NewSourceRange(commaStart, commaStart))
			p.treeStack = append(p.treeStack, nullTree)
		} else {
			if p.Peek(TokenType_ELLIPSIS) {
				p.RecordFeatureUsed(FeatureSet_Feature_SPREAD_EXPRESSIONS)
				spread := p.ParseIterSpread()
				p.treeStack = append(p.treeStack, spread)
			} else {
				elem := p.ParseAssignmentExpression()
				p.treeStack = append(p.treeStack, elem)
			}
			if !p.Peek(TokenType_CLOSE_SQUARE) {
				trailingCommaToken = p.Eat(TokenType_COMMA)
			}
		}
	}
	p.Eat(TokenType_CLOSE_SQUARE)

	p.MaybeReportTrailingComma(trailingCommaToken)

	elements := p.popTrees(base)
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
			p.errorReporter.ReportWarning(
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
	switch {
	case type_ == TokenType_STAR:
		return p.ParseObjectLiteralPropertyAssignmentGenerator()
	case type_ == TokenType_ELLIPSIS:
		p.RecordFeatureUsed(FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_ELLIPSIS)
		operand := p.ParseAssignmentExpression()
		return NewObjectSpreadTree(p.GetTreeLocation(start), operand)
	case type_ == TokenType_STRING ||
		type_ == TokenType_NUMBER ||
		type_ == TokenType_BIGINT ||
		type_ == TokenType_IDENTIFIER ||
		Keywords_IsKeywordTokenType(type_):
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
	case type_ == TokenType_OPEN_SQUARE:
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
	default:
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
		param := p.ParseParameter()
		params := p.singleTreeList(param)
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
			base := len(p.treeStack)
			p.treeStack = append(p.treeStack, result)
			param := p.ParseParameter()
			p.treeStack = append(p.treeStack, param)
			exprs := p.popTrees(base)
			result = NewCommaExpressionTree(p.GetTreeLocation(start), exprs)
		}
	}
	p.Eat(TokenType_CLOSE_PAREN)
	return NewParenExpressionTree(p.GetTreeLocation(start), result)
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
			p.singleTreeList(leftOfArrow),
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
			p.singleTreeList(singleParameter),
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

// Java: Parser.parseFunctionExpression()
func (p *Parser) ParseFunctionExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(Keywords_FUNCTION.Type())
	isGenerator := p.EatOpt(TokenType_STAR).IsPresent()

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
		SetName(p.EatIdOpt())
	p.ParseFunctionTail(builder, functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))

	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseAsyncFunctionExpression()
func (p *Parser) ParseAsyncFunctionExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatAsyncFunctionStart()

	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
		SetName(p.EatIdOpt()).
		SetAsync(true)

	p.ParseFunctionTail(
		builder, functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.peekClassDeclaration()
func (p *Parser) PeekClassDeclaration() bool {
	return p.Peek(TokenType_CLASS)
}

// Java: Parser.parseClassDeclaration()
func (p *Parser) ParseClassDeclaration() ParseTree {
	return p.ParseClass( /* isExpression= */ false)
}

// Java: Parser.parseClassExpression()
func (p *Parser) ParseClassExpression() ParseTree {
	return p.ParseClass( /* isExpression= */ true)
}

// Java: Parser.parseClass(boolean)
func (p *Parser) ParseClass(isExpression bool) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CLASS)
	var name Token
	if !isExpression || p.PeekId() {
		name = p.EatId()
	}

	var superClass ParseTree
	if p.Peek(TokenType_EXTENDS) {
		p.Eat(TokenType_EXTENDS)
		superClass = p.ParseLeftHandSideExpression()
	}

	p.Eat(TokenType_OPEN_CURLY)
	elements := p.ParseClassElements()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewClassDeclarationTree(p.GetTreeLocation(start), name, superClass, elements)
}

// Java: Parser.parseClassElements()
func (p *Parser) ParseClassElements() []ParseTree {
	base := len(p.treeStack)
	for {
		token := p.PeekToken()
		if token.Type() == TokenType_SEMI_COLON {
			p.Eat(TokenType_SEMI_COLON)
			continue
		}
		if !p.IsClassElementStart(token) {
			return p.popTrees(base)
		}
		elem := p.ParseClassElement()
		p.treeStack = append(p.treeStack, elem)
	}
}

// Java: Parser.isClassElementStart(Token)
func (p *Parser) IsClassElementStart(token Token) bool {
	switch token.Type() {
	case TokenType_IDENTIFIER, TokenType_NUMBER, TokenType_BIGINT, TokenType_STAR, TokenType_STATIC, TokenType_STRING, TokenType_OPEN_SQUARE:
		return true
	default:
		return Keywords_IsKeywordTokenType(token.Type())
	}
}

// Java: Parser.createObjectLiteralElementInfo()
func (p *Parser) CreateObjectLiteralElementInfo() *ClassOrObjectElementInfo {
	return ClassOrObjectElementInfo_CreateObjectLiteralElementInfo(p.GetTreeStartLocation())
}

// Java: Parser.parseClassElement()
func (p *Parser) ParseClassElement() ParseTree {
	if p.Peek(TokenType_SEMI_COLON) {
		return p.ParseEmptyStatement()
	}

	if p.PeekClassStaticInitializerBlock() {
		return p.ParseClassStaticInitializerBlock()
	}

	elementInfo := ClassOrObjectElementInfo_CreateClassMemberInfo(
		p.GetTreeStartLocation(), p.EatStaticIfNotElementName())
	if p.PeekGetAccessor() {
		return p.ParseGetAccessor(elementInfo)
	} else if p.PeekSetAccessor() {
		return p.ParseSetAccessor(elementInfo)
	} else if p.PeekAsyncMethod() {
		return p.ParseAsyncMethod(elementInfo)
	} else {
		return p.ParseClassMemberDeclaration(elementInfo)
	}
}

// Java: Parser.eatStaticIfNotElementName()
func (p *Parser) EatStaticIfNotElementName() bool {
	// only eat `static` if it being used as a keyword and not
	// a member name.
	if p.Peek(TokenType_STATIC) && p.IsClassElementStart(p.PeekTokenAt(1)) {
		p.Eat(TokenType_STATIC)
		return true
	}
	return false
}

// Java: Parser.peekAsyncMethod()
func (p *Parser) PeekAsyncMethod() bool {
	return p.PeekPredefinedString(Parser_ASYNC) &&
		!p.PeekImplicitSemiColonAt(1) &&
		(p.PeekPropertyNameOrComputedProp(1) ||
			(p.PeekAt(1, TokenType_STAR) && p.PeekPropertyNameOrComputedProp(2)))
}

// Java: Parser.peekClassStaticInitializerBlock()
func (p *Parser) PeekClassStaticInitializerBlock() bool {
	return p.Peek(TokenType_STATIC) && p.PeekAt(1, TokenType_OPEN_CURLY)
}

// Java: Parser.parseClassElementName(ClassOrObjectElementInfo)
func (p *Parser) ParseClassElementName(elementInfo *ClassOrObjectElementInfo) {
	if p.PeekPropertyName(0) {
		if p.PeekIdOrKeyword() {
			elementInfo.SetName(p.EatIdOrKeywordAsId())
			if elementInfo.GetName().IsKeyword() {
				p.RecordFeatureUsed(FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
			}
		} else {
			// { 'str'() {} }
			// { 123() {} }
			// Treat these as if they were computed properties.
			// TODO(b/123769080): Stop making this assumption!
			elementInfo.SetNameExpr(p.ParseLiteralExpression())
		}
	} else {
		elementInfo.SetNameExpr(p.ParseComputedPropertyName())
	}
}

// Java: Parser.parseFieldDefinition(ClassOrObjectElementInfo)
func (p *Parser) ParseFieldDefinition(elementInfo *ClassOrObjectElementInfo) ParseTree {
	var initializer ParseTree
	if p.Peek(TokenType_EQUAL) {
		initializer = p.ParseInitializer(Parser_Expression_NORMAL)
	}

	p.EatPossiblyImplicitSemiColon()
	if elementInfo.HasName() {
		return NewFieldDeclarationTree(
			p.GetTreeLocation(elementInfo.start),
			elementInfo.GetName(),
			elementInfo.isStatic,
			initializer)
	} else {
		jutil.CheckState(elementInfo.HasNameExpr())
		return NewComputedPropertyFieldTree(
			p.GetTreeLocation(elementInfo.start),
			elementInfo.GetNameExpr(),
			elementInfo.isStatic,
			initializer)
	}
}

// Java: Parser.parseMethodDefinition(ClassOrObjectElementInfo, boolean)
func (p *Parser) ParseMethodDefinition(elementInfo *ClassOrObjectElementInfo, isGenerator bool) ParseTree {
	var kind FunctionDeclarationTree_Kind
	if elementInfo.HasNameExpr() {
		kind = FunctionDeclarationTree_Kind_EXPRESSION
	} else {
		kind = FunctionDeclarationTree_Kind_MEMBER
	}

	builder := FunctionDeclarationTree_NewBuilder(kind).
		SetIsClassMember(elementInfo.isClassMember).
		SetStatic(elementInfo.isStatic)
	if elementInfo.HasName() {
		builder.SetName(elementInfo.GetName())
	}
	p.ParseFunctionTail(builder, functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))

	var function ParseTree = builder.Build(p.GetTreeLocation(elementInfo.start))
	if kind == FunctionDeclarationTree_Kind_MEMBER {
		return function
	} else {
		jutil.CheckState(elementInfo.HasNameExpr())
		return NewComputedPropertyMethodTree(
			p.GetTreeLocation(elementInfo.start), elementInfo.GetNameExpr(), function)
	}
}

// Java: Parser.parseObjectLiteralMethodDeclaration()
func (p *Parser) ParseObjectLiteralMethodDeclaration() ParseTree {
	return p.ParseMethodDeclaration(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseMethodDeclaration(ClassOrObjectElementInfo)
func (p *Parser) ParseMethodDeclaration(elementInfo *ClassOrObjectElementInfo) ParseTree {
	isGenerator := p.EatOpt(TokenType_STAR).IsPresent()
	p.ParseClassElementName(elementInfo)
	return p.ParseMethodDefinition(elementInfo, isGenerator)
}

// Java: Parser.parseClassMemberDeclaration(ClassOrObjectElementInfo)
func (p *Parser) ParseClassMemberDeclaration(elementInfo *ClassOrObjectElementInfo) ParseTree {
	isGenerator := p.EatOpt(TokenType_STAR).IsPresent()
	p.ParseClassElementName(elementInfo)
	if p.PeekTypeAt(0) == TokenType_OPEN_PAREN {
		return p.ParseMethodDefinition(elementInfo, isGenerator)
	} else {
		return p.ParseFieldDefinition(elementInfo)
	}
}

// Java: Parser.parseObjectLiteralAsyncMethod()
func (p *Parser) ParseObjectLiteralAsyncMethod() ParseTree {
	return p.ParseAsyncMethod(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseAsyncMethod(ClassOrObjectElementInfo)
func (p *Parser) ParseAsyncMethod(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(Parser_ASYNC)
	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}
	if p.PeekPropertyName(0) {
		if p.PeekIdOrKeyword() {
			name := p.EatIdOrKeywordAsId()
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_MEMBER).
				SetAsync(true).
				SetGenerator(generator).
				SetIsClassMember(elementInfo.isClassMember).
				SetStatic(elementInfo.isStatic).
				SetName(name)
			p.ParseFunctionTail(
				builder,
				functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

			return builder.Build(p.GetTreeLocation(name.GetStart()))
		} else {
			// { 'str'() {} }
			// { 123() {} }
			// Treat these as if they were computed properties.
			var nameExpr ParseTree = p.ParseLiteralExpression()
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
				SetAsync(true).
				SetGenerator(generator).
				SetIsClassMember(elementInfo.isClassMember).
				SetStatic(elementInfo.isStatic)
			p.ParseFunctionTail(
				builder,
				functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

			var function ParseTree = builder.Build(p.GetTreeLocation(nameExpr.GetStart()))
			return NewComputedPropertyMethodTree(
				p.GetTreeLocation(nameExpr.GetStart()), nameExpr, function)
		}
	} else {
		// expect '[' to start computed property name
		nameExpr := p.ParseComputedPropertyName()
		builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
			SetAsync(true).
			SetGenerator(generator).
			SetIsClassMember(elementInfo.isClassMember).
			SetStatic(elementInfo.isStatic)
		p.ParseFunctionTail(
			builder, functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

		var function ParseTree = builder.Build(p.GetTreeLocation(nameExpr.GetStart()))
		return NewComputedPropertyMethodTree(
			p.GetTreeLocation(nameExpr.GetStart()), nameExpr, function)
	}
}

// Java: Parser.parseClassStaticInitializerBlock()
func (p *Parser) ParseClassStaticInitializerBlock() ParseTree {
	p.Eat(TokenType_STATIC)
	return p.ParseBlock()
}

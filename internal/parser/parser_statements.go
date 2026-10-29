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

// 12 Statements

// ParseStatement: in V8, all source elements may appear where statements
// occur in the grammar.
//
// Java: Parser.parseStatement()
func (p *Parser) ParseStatement() ParseTree {
	return p.ParseSourceElement()
}

// ParseStatementStandard reflects the ECMA standard. Most places use
// peekStatement instead.
//
// Java: Parser.parseStatementStandard()
func (p *Parser) ParseStatementStandard() ParseTree {
	switch p.PeekType() {
	case TokenType_OPEN_CURLY:
		return p.ParseBlock()
	case TokenType_CONST, TokenType_VAR:
		return p.ParseVariableStatement()
	case TokenType_SEMI_COLON:
		return p.ParseEmptyStatement()
	case TokenType_IF:
		return p.ParseIfStatement()
	case TokenType_DO:
		return p.ParseDoWhileStatement()
	case TokenType_WHILE:
		return p.ParseWhileStatement()
	case TokenType_FOR:
		return p.ParseForStatement()
	case TokenType_CONTINUE:
		return p.ParseContinueStatement()
	case TokenType_BREAK:
		return p.ParseBreakStatement()
	case TokenType_RETURN:
		return p.ParseReturnStatement()
	case TokenType_WITH:
		return p.ParseWithStatement()
	case TokenType_SWITCH:
		return p.ParseSwitchStatement()
	case TokenType_THROW:
		return p.ParseThrowStatement()
	case TokenType_TRY:
		return p.ParseTryStatement()
	case TokenType_DEBUGGER:
		return p.ParseDebuggerStatement()
	default:
		if p.PeekLabelledStatement() {
			return p.ParseLabelledStatement()
		}
		return p.ParseExpressionStatement()
	}
}

// PeekStatement: in V8 all source elements may appear where statements
// appear in the grammar.
//
// Java: Parser.peekStatement()
func (p *Parser) PeekStatement() bool {
	return p.PeekSourceElement()
}

// PeekStatementStandard reflects the ECMA standard. Most places use
// peekStatement instead.
//
// Java: Parser.peekStatementStandard()
func (p *Parser) PeekStatementStandard() bool {
	switch p.PeekType() {
	case TokenType_OPEN_CURLY,
		TokenType_VAR,
		TokenType_CONST,
		TokenType_SEMI_COLON,
		TokenType_IF,
		TokenType_DO,
		TokenType_WHILE,
		TokenType_FOR,
		TokenType_CONTINUE,
		TokenType_BREAK,
		TokenType_RETURN,
		TokenType_WITH,
		TokenType_SWITCH,
		TokenType_THROW,
		TokenType_TRY,
		TokenType_DEBUGGER,
		TokenType_YIELD,
		TokenType_IDENTIFIER,
		TokenType_TYPE,
		TokenType_DECLARE,
		TokenType_MODULE,
		TokenType_NAMESPACE,
		TokenType_THIS,
		TokenType_CLASS,
		TokenType_SUPER,
		TokenType_NUMBER,
		TokenType_BIGINT,
		TokenType_STRING,
		TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_NULL,
		TokenType_TRUE,
		TokenType_SLASH,       // regular expression literal
		TokenType_SLASH_EQUAL, // regular expression literal
		TokenType_FALSE,
		TokenType_OPEN_SQUARE,
		TokenType_OPEN_PAREN,
		TokenType_NEW,
		TokenType_DELETE,
		TokenType_VOID,
		TokenType_TYPEOF,
		TokenType_PLUS_PLUS,
		TokenType_MINUS_MINUS,
		TokenType_PLUS,
		TokenType_MINUS,
		TokenType_TILDE,
		TokenType_BANG,
		TokenType_IMPORT:
		return true
	default:
		return false
	}
}

// 12.1 Block

// Java: Parser.parseBlock()
func (p *Parser) ParseBlock() *BlockTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_CURLY)
	// Spec says Statement list. However functions are also embedded in the wild.
	result := p.ParseSourceElementList()
	p.Eat(TokenType_CLOSE_CURLY)
	return p.newBlockTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseStatementList()
func (p *Parser) ParseStatementList() []ParseTree {
	base := len(p.treeStack)
	for p.PeekStatement() {
		stmt := p.ParseStatement()
		p.treeStack = append(p.treeStack, stmt)
	}
	return p.popTrees(base)
}

// 12.2 Variable Statement

// Java: Parser.parseVariableStatement()
func (p *Parser) ParseVariableStatement() *VariableStatementTree {
	start := p.GetTreeStartLocation()
	declarations := p.ParseVariableDeclarationList()
	p.EatPossiblyImplicitSemiColon()
	return p.newVariableStatementTree(p.GetTreeLocation(start), declarations)
}

// Java: Parser.parseVariableDeclarationList()
func (p *Parser) ParseVariableDeclarationList() *VariableDeclarationListTree {
	return p.ParseVariableDeclarationListWithExpression(Parser_Expression_NORMAL)
}

// Java: Parser.parseVariableDeclarationListNoIn()
func (p *Parser) ParseVariableDeclarationListNoIn() *VariableDeclarationListTree {
	return p.ParseVariableDeclarationListWithExpression(Parser_Expression_NO_IN)
}

// ParseVariableDeclarationListWithExpression returns nil (Java null) after
// reporting an error if no declaration keyword is next.
//
// Java: Parser.parseVariableDeclarationList(Expression)
func (p *Parser) ParseVariableDeclarationListWithExpression(expressionIn Parser_Expression) *VariableDeclarationListTree {
	start := p.GetTreeStartLocation()
	token := p.PeekType()

	switch token {
	case TokenType_CONST, TokenType_LET, TokenType_VAR:
		p.Eat(token)
	default:
		p.ReportErrorAtToken(p.PeekToken(), "expected declaration")
		return nil
	}

	base := len(p.varDeclStack)

	decl := p.ParseVariableDeclaration(token, expressionIn)
	p.varDeclStack = append(p.varDeclStack, decl)
	for p.Peek(TokenType_COMMA) {
		p.Eat(TokenType_COMMA)
		decl = p.ParseVariableDeclaration(token, expressionIn)
		p.varDeclStack = append(p.varDeclStack, decl)
	}
	declarations := popStackChunked(&p.varDeclStack, base, &p.varDeclLists, parser_treeListChunk)
	return p.newVariableDeclarationListTree(p.GetTreeLocation(start), token, declarations)
}

// Java: Parser.parseVariableDeclaration(TokenType, Expression)
func (p *Parser) ParseVariableDeclaration(binding TokenType, expressionIn Parser_Expression) *VariableDeclarationTree {

	start := p.GetTreeStartLocation()
	var lvalue ParseTree
	if p.PeekPatternStart() {
		lvalue = p.ParsePattern(Parser_PatternKind_INITIALIZER)
	} else {
		lvalue = p.ParseIdentifierExpression()
	}

	var initializer ParseTree
	if p.Peek(TokenType_EQUAL) {
		initializer = p.ParseInitializer(expressionIn)
	} else if expressionIn != Parser_Expression_NO_IN {
		// NOTE(blickly): this is a bit of a hack, declarations outside of for statements allow "in",
		// and by chance, also must have initializers for const/destructuring. Vanilla for loops
		// also require intializers, but are handled separately in checkVanillaForInitializers
		p.MaybeReportNoInitializer(binding, lvalue)
	}
	return p.newVariableDeclarationTree(p.GetTreeLocation(start), lvalue, initializer)
}

// Java: Parser.parseInitializer(Expression)
func (p *Parser) ParseInitializer(expressionIn Parser_Expression) ParseTree {
	p.Eat(TokenType_EQUAL)
	return p.ParseAssignment(expressionIn)
}

// 12.3 Empty Statement

// Java: Parser.parseEmptyStatement()
func (p *Parser) ParseEmptyStatement() *EmptyStatementTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SEMI_COLON)
	return NewEmptyStatementTree(p.GetTreeLocation(start))
}

// 12.4 Expression Statement

// Java: Parser.parseExpressionStatement()
func (p *Parser) ParseExpressionStatement() *ExpressionStatementTree {
	start := p.GetTreeStartLocation()
	expression := p.ParseExpression()
	p.EatPossiblyImplicitSemiColon()
	return p.newExpressionStatementTree(p.GetTreeLocation(start), expression)
}

// 12.5 If Statement

// Java: Parser.parseIfStatement()
func (p *Parser) ParseIfStatement() *IfStatementTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IF)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	ifClause := p.ParseStatement()
	var elseClause ParseTree
	if p.Peek(TokenType_ELSE) {
		p.Eat(TokenType_ELSE)
		elseClause = p.ParseStatement()
	}
	return NewIfStatementTree(p.GetTreeLocation(start), condition, ifClause, elseClause)
}

// 12.6 Iteration Statements

// 12.6.1 The do-while Statement

// Java: Parser.parseDoWhileStatement()
func (p *Parser) ParseDoWhileStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_DO)
	body := p.ParseStatement()
	p.Eat(TokenType_WHILE)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	// The semicolon after the "do-while" is optional.
	if p.Peek(TokenType_SEMI_COLON) {
		p.Eat(TokenType_SEMI_COLON)
	}
	return NewDoWhileStatementTree(p.GetTreeLocation(start), body, condition)
}

// 12.6.2 The while Statement

// Java: Parser.parseWhileStatement()
func (p *Parser) ParseWhileStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_WHILE)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewWhileStatementTree(p.GetTreeLocation(start), condition, body)
}

// 12.6.3 The for Statement
// 12.6.4 The for-in Statement
// The for-of Statement
// The for-await-of Statement

// Java: Parser.parseForStatement()
func (p *Parser) ParseForStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_FOR)
	awaited := p.PeekPredefinedString(Parser_AWAIT)
	if awaited {
		p.EatPredefinedString(Parser_AWAIT)
	}
	p.Eat(TokenType_OPEN_PAREN)
	if p.PeekVariableDeclarationList() {
		variables := p.ParseVariableDeclarationListNoIn()
		if p.Peek(TokenType_IN) {
			if awaited {
				p.ReportError("for-await-of is the only allowed asynchronous iteration")
			}
			// for-in: only one declaration allowed
			if len(variables.Declarations) > 1 {
				p.ReportError("for-in statement may not have more than one variable declaration")
			}
			declaration := variables.Declarations[0]
			if declaration.Initializer != nil {
				// An initializer is allowed here in ES5 and below, but not in ES6.
				// Warn about it, to encourage people to eliminate it from their code.
				// http://esdiscuss.org/topic/initializer-expression-on-for-in-syntax-subject
				if p.config.atLeast6 {
					p.ReportError("for-in statement may not have initializer")
				} else {
					p.errorReporter.ReportWarning(
						p.linePos(declaration.Location.Start), "for-in statement should not have initializer")
				}
			}

			return p.ParseForInStatement(start, variableDeclarationListAsParseTree(variables))
		} else if p.PeekPredefinedString(PredefinedName_OF) {
			// for-of: only one declaration allowed
			if len(variables.Declarations) > 1 {
				if awaited {
					p.ReportError("for-await-of statement may not have more than one variable declaration")
				} else {
					p.ReportError("for-of statement may not have more than one variable declaration")
				}
			}
			// for-of: initializer is illegal
			declaration := variables.Declarations[0]
			if declaration.Initializer != nil {
				if awaited {
					p.ReportError("for-await-of statement may not have initializer")
				} else {
					p.ReportError("for-of statement may not have initializer")
				}
			}

			if awaited {
				return p.ParseForAwaitOfStatement(start, variableDeclarationListAsParseTree(variables))
			} else {
				return p.ParseForOfStatement(start, variableDeclarationListAsParseTree(variables))
			}
		} else {
			// "Vanilla" for statement: const/destructuring must have initializer
			p.CheckVanillaForInitializers(variables)
			return p.ParseForStatementWithInitializer(start, variableDeclarationListAsParseTree(variables))
		}
	}

	if p.Peek(TokenType_SEMI_COLON) {
		return p.ParseForStatementWithInitializer(start, nil)
	}

	initializer := p.ParseExpressionNoIn()
	if p.Peek(TokenType_IN) || p.Peek(TokenType_EQUAL) || p.PeekPredefinedString(PredefinedName_OF) {
		initializer = p.TransformLeftHandSideExpression(initializer)
		if !initializer.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
		}
	}

	if p.Peek(TokenType_IN) || p.PeekPredefinedString(PredefinedName_OF) {
		if initializer.GetParseTreeBase().Type != ParseTreeType_BINARY_OPERATOR &&
			initializer.GetParseTreeBase().Type != ParseTreeType_COMMA_EXPRESSION {
			if p.Peek(TokenType_IN) {
				return p.ParseForInStatement(start, initializer)
			} else {
				// for {await}? ( _ of _ )
				if awaited {
					return p.ParseForAwaitOfStatement(start, initializer)
				} else {
					return p.ParseForOfStatement(start, initializer)
				}
			}
		}
	}

	return p.ParseForStatementWithInitializer(start, initializer)
}

// The for-of Statement
// for  (  { let | var }?  identifier  of  expression  )  statement

// Java: Parser.parseForOfStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForOfStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.EatPredefinedString(PredefinedName_OF)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForOfStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// Java: Parser.parseForAwaitOfStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForAwaitOfStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.EatPredefinedString(PredefinedName_OF)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForAwaitOfStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// CheckVanillaForInitializers checks variable declarations in for
// statements.
//
// Java: Parser.checkVanillaForInitializers(VariableDeclarationListTree)
func (p *Parser) CheckVanillaForInitializers(variables *VariableDeclarationListTree) {
	for _, declaration := range variables.Declarations {
		if declaration.Initializer == nil {
			p.MaybeReportNoInitializer(variables.DeclarationType, declaration.Lvalue)
		}
	}
}

// MaybeReportNoInitializer reports if declaration requires an initializer,
// assuming initializer is absent.
//
// Java: Parser.maybeReportNoInitializer(TokenType, ParseTree)
func (p *Parser) MaybeReportNoInitializer(token TokenType, lvalue ParseTree) {
	if token == TokenType_CONST {
		p.ReportError("const variables must have an initializer")
	} else if lvalue.IsPattern() {
		p.ReportError("destructuring must have an initializer")
	}
}

// Java: Parser.peekVariableDeclarationList()
func (p *Parser) PeekVariableDeclarationList() bool {
	switch p.PeekType() {
	case TokenType_VAR, TokenType_CONST, TokenType_LET:
		return true
	default:
		return false
	}
}

// 12.6.3 The for Statement

// Java: Parser.parseForStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForStatementWithInitializer(start SourcePosition, initializer ParseTree) ParseTree {
	if initializer == nil {
		initializer = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_SEMI_COLON)

	var condition ParseTree
	if !p.Peek(TokenType_SEMI_COLON) {
		condition = p.ParseExpression()
	} else {
		condition = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_SEMI_COLON)

	var increment ParseTree
	if !p.Peek(TokenType_CLOSE_PAREN) {
		increment = p.ParseExpression()
	} else {
		increment = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForStatementTree(p.GetTreeLocation(start), initializer, condition, increment, body)
}

// 12.6.4 The for-in Statement

// Java: Parser.parseForInStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForInStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.Eat(TokenType_IN)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForInStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// 12.7 The continue Statement

// Java: Parser.parseContinueStatement()
func (p *Parser) ParseContinueStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CONTINUE)
	var name Token
	if !p.PeekImplicitSemiColon() {
		name = p.EatIdOpt()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewContinueStatementTree(p.GetTreeLocation(start), name)
}

// 12.8 The break Statement

// Java: Parser.parseBreakStatement()
func (p *Parser) ParseBreakStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_BREAK)
	var name Token
	if !p.PeekImplicitSemiColon() {
		name = p.EatIdOpt()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewBreakStatementTree(p.GetTreeLocation(start), name)
}

// 12.9 The return Statement

// Java: Parser.parseReturnStatement()
func (p *Parser) ParseReturnStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_RETURN)
	var expression ParseTree
	if !p.PeekImplicitSemiColon() {
		expression = p.ParseExpression()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewReturnStatementTree(p.GetTreeLocation(start), expression)
}

// 12.10 The with Statement

// Java: Parser.parseWithStatement()
func (p *Parser) ParseWithStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_WITH)
	p.Eat(TokenType_OPEN_PAREN)
	expression := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewWithStatementTree(p.GetTreeLocation(start), expression, body)
}

// 12.11 The switch Statement

// Java: Parser.parseSwitchStatement()
func (p *Parser) ParseSwitchStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SWITCH)
	p.Eat(TokenType_OPEN_PAREN)
	expression := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	p.Eat(TokenType_OPEN_CURLY)
	caseClauses := p.ParseCaseClauses()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewSwitchStatementTree(p.GetTreeLocation(start), expression, caseClauses)
}

// Java: Parser.parseCaseClauses()
func (p *Parser) ParseCaseClauses() []ParseTree {
	foundDefaultClause := false
	base := len(p.treeStack)

	for {
		start := p.GetTreeStartLocation()
		switch p.PeekType() {
		case TokenType_CASE:
			p.Eat(TokenType_CASE)
			expression := p.ParseExpression()
			p.Eat(TokenType_COLON)
			statements := p.ParseCaseStatementsOpt()
			clause := NewCaseClauseTree(p.GetTreeLocation(start), expression, statements)
			p.treeStack = append(p.treeStack, clause)
		case TokenType_DEFAULT:
			if foundDefaultClause {
				p.ReportError("Switch statements may have at most one default clause")
			} else {
				foundDefaultClause = true
			}
			p.Eat(TokenType_DEFAULT)
			p.Eat(TokenType_COLON)
			clause := NewDefaultClauseTree(p.GetTreeLocation(start), p.ParseCaseStatementsOpt())
			p.treeStack = append(p.treeStack, clause)
		default:
			return p.popTrees(base)
		}
	}
}

// Java: Parser.parseCaseStatementsOpt()
func (p *Parser) ParseCaseStatementsOpt() []ParseTree {
	return p.ParseStatementList()
}

// 12.12 Labelled Statement

// Java: Parser.parseLabelledStatement()
func (p *Parser) ParseLabelledStatement() ParseTree {
	start := p.GetTreeStartLocation()
	name := p.EatId()
	p.Eat(TokenType_COLON)
	return NewLabelledStatementTree(p.GetTreeLocation(start), name, p.ParseStatement())
}

// Java: Parser.peekLabelledStatement()
func (p *Parser) PeekLabelledStatement() bool {
	return p.PeekId() && p.PeekAt(1, TokenType_COLON)
}

// 12.13 Throw Statement

// Java: Parser.parseThrowStatement()
func (p *Parser) ParseThrowStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_THROW)
	var value ParseTree
	if p.PeekImplicitSemiColon() {
		p.ReportError("semicolon/newline not allowed after 'throw'")
	} else {
		value = p.ParseExpression()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewThrowStatementTree(p.GetTreeLocation(start), value)
}

// 12.14 Try Statement

// Java: Parser.parseTryStatement()
func (p *Parser) ParseTryStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_TRY)
	var body ParseTree = p.ParseBlock()
	var catchBlock ParseTree
	if p.Peek(TokenType_CATCH) {
		catchBlock = p.ParseCatch()
	}
	var finallyBlock ParseTree
	if p.Peek(TokenType_FINALLY) {
		finallyBlock = p.ParseFinallyBlock()
	}
	if catchBlock == nil && finallyBlock == nil {
		p.ReportError("'catch' or 'finally' expected.")
	}
	return NewTryStatementTree(p.GetTreeLocation(start), body, catchBlock, finallyBlock)
}

// Java: Parser.parseCatch()
func (p *Parser) ParseCatch() *CatchTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CATCH)

	var exception ParseTree = NewEmptyStatementTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))

	if p.PeekToken().Type() == TokenType_OPEN_PAREN {
		p.Eat(TokenType_OPEN_PAREN)
		if p.PeekPatternStart() {
			exception = p.ParsePattern(Parser_PatternKind_INITIALIZER)
		} else {
			exception = p.ParseIdentifierExpression()
		}
		p.Eat(TokenType_CLOSE_PAREN)
	} else {
		p.RecordFeatureUsed(FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
	}

	catchBody := p.ParseBlock()
	return NewCatchTree(p.GetTreeLocation(start), exception, catchBody)
}

// Java: Parser.parseFinallyBlock()
func (p *Parser) ParseFinallyBlock() *FinallyTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_FINALLY)
	finallyBlock := p.ParseBlock()
	return NewFinallyTree(p.GetTreeLocation(start), finallyBlock)
}

// 12.15 The Debugger Statement

// Java: Parser.parseDebuggerStatement()
func (p *Parser) ParseDebuggerStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_DEBUGGER)
	p.EatPossiblyImplicitSemiColon()

	return NewDebuggerStatementTree(p.GetTreeLocation(start))
}

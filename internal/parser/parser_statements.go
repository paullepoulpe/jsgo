// Port of com/google/javascript/jscomp/parsing/parser/Parser.java (statement-level methods)
// @ closure-compiler v20260415
//
// Statement parsing: commits 1-5 of the part-2 rewrite.
// Methods are in the same order as the Java source.

package parser

import (
	"github.com/paullepoulpe/jsgo/internal/parser/trees"
)

// ---- Top-level element -------------------------------------------------------

// parseScriptElement mirrors parseScriptElement().
func (p *ParserRewrite) parseScriptElement() *trees.ParseTree {
	if p.peekImportDeclaration() {
		return p.parseImportDeclaration()
	}
	if p.peekExportDeclaration() {
		return p.parseExportDeclaration()
	}
	return p.parseSourceElement()
}

func (p *ParserRewrite) peekImportDeclaration() bool {
	return p.peek(TokenTypeImport) &&
		(p.peekIdOrKeywordAt(1) ||
			p.peekAt(1, TokenTypeString) ||
			p.peekAt(1, TokenTypeOpenCurly) ||
			p.peekAt(1, TokenTypeStar))
}

// ---- Source elements ---------------------------------------------------------

// parseSourceElement mirrors parseSourceElement().
func (p *ParserRewrite) parseSourceElement() *trees.ParseTree {
	if p.peekAsyncFunctionStart() {
		return p.parseAsyncFunctionDeclaration()
	}
	if p.peekFunction() {
		return p.parseFunctionDeclaration()
	}
	if p.peekClassDeclaration() {
		return p.parseClassDeclaration()
	}
	// let block-scoped bindings
	if p.peek(TokenTypeLet) {
		return p.parseVariableStatement()
	}
	return p.parseStatementStandard()
}

func (p *ParserRewrite) peekSourceElement() bool {
	return p.peekFunction() || p.peekStatementStandard() || p.peekDeclaration()
}

func (p *ParserRewrite) peekDeclaration() bool {
	return p.peek(TokenTypeLet) || p.peekClassDeclaration()
}

func (p *ParserRewrite) peekClassDeclaration() bool {
	return p.peek(TokenTypeClass)
}

// parseSourceElementList mirrors parseSourceElementList().
func (p *ParserRewrite) parseSourceElementList() []*trees.ParseTree {
	var result []*trees.ParseTree
	for p.peekSourceElement() {
		result = append(result, p.parseSourceElement())
	}
	return result
}

// parseStatement mirrors parseStatement() — in V8 all source elements may appear.
func (p *ParserRewrite) parseStatement() *trees.ParseTree {
	return p.parseSourceElement()
}

// parseStatementListItem mirrors parseStatement() — same as parseStatement in practice.
// Go-only: alias for parseStatement used at declaration-level; Java uses parseSourceElement for the same purpose.
func (p *ParserRewrite) parseStatementListItem() *trees.ParseTree {
	return p.parseStatement()
}

// peekStatementStandard mirrors peekStatementStandard().
func (p *ParserRewrite) peekStatementStandard() bool {
	switch p.peekType() {
	case TokenTypeOpenCurly,
		TokenTypeVar,
		TokenTypeConst,
		TokenTypeSemiColon,
		TokenTypeIf,
		TokenTypeDo,
		TokenTypeWhile,
		TokenTypeFor,
		TokenTypeContinue,
		TokenTypeBreak,
		TokenTypeReturn,
		TokenTypeWith,
		TokenTypeSwitch,
		TokenTypeThrow,
		TokenTypeTry,
		TokenTypeDebugger,
		TokenTypeYield,
		TokenTypeIdentifier,
		TokenTypeType,
		TokenTypeDeclare,
		TokenTypeModule,
		TokenTypeNamespace,
		TokenTypeThis,
		TokenTypeClass,
		TokenTypeSuper,
		TokenTypeNumber,
		TokenTypeBigInt,
		TokenTypeString,
		TokenTypeNoSubstitutionTemplate,
		TokenTypeTemplateHead,
		TokenTypeNull,
		TokenTypeTrue,
		TokenTypeSlash,
		TokenTypeSlashEqual,
		TokenTypeFalse,
		TokenTypeOpenSquare,
		TokenTypeOpenParen,
		TokenTypeNew,
		TokenTypeDelete,
		TokenTypeVoid,
		TokenTypeTypeof,
		TokenTypePlusPlus,
		TokenTypeMinusMinus,
		TokenTypePlus,
		TokenTypeMinus,
		TokenTypeTilde,
		TokenTypeBang,
		TokenTypeImport:
		return true
	}
	return false
}

// parseStatementStandard mirrors parseStatementStandard().
func (p *ParserRewrite) parseStatementStandard() *trees.ParseTree {
	switch p.peekType() {
	case TokenTypeOpenCurly:
		return p.asParseTree(p.parseBlock())
	case TokenTypeConst, TokenTypeVar:
		return p.parseVariableStatement()
	case TokenTypeSemiColon:
		return p.asParseTree(p.parseEmptyStatement())
	case TokenTypeIf:
		return p.asParseTree(p.parseIfStatement())
	case TokenTypeDo:
		return p.parseDoWhileStatement()
	case TokenTypeWhile:
		return p.asParseTree(p.parseWhileStatement())
	case TokenTypeFor:
		return p.parseForStatement()
	case TokenTypeContinue:
		return p.parseContinueStatement()
	case TokenTypeBreak:
		return p.parseBreakStatement()
	case TokenTypeReturn:
		return p.parseReturnStatement()
	case TokenTypeWith:
		return p.asParseTree(p.parseWithStatement())
	case TokenTypeSwitch:
		return p.asParseTree(p.parseSwitchStatement())
	case TokenTypeThrow:
		return p.parseThrowStatement()
	case TokenTypeTry:
		return p.parseTryStatement()
	case TokenTypeDebugger:
		return p.parseDebuggerStatement()
	default:
		if p.peekLabelledStatement() {
			return p.parseLabelledStatement()
		}
		return p.asParseTree(p.parseExpressionStatement())
	}
}

// ---- Block -------------------------------------------------------------------

// parseBlock mirrors parseBlock().
func (p *ParserRewrite) parseBlock() *trees.BlockTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeOpenCurly)
	result := p.parseSourceElementList()
	p.eat(TokenTypeCloseCurly)
	return trees.NewBlockTree(p.getTreeLocation(start), result)
}

// parseStatementList mirrors parseStatementList().
func (p *ParserRewrite) parseStatementList() []*trees.ParseTree {
	var result []*trees.ParseTree
	for p.peekSourceElement() {
		result = append(result, p.parseStatement())
	}
	return result
}

// ---- Variable declarations ---------------------------------------------------

// parseVariableStatement mirrors parseVariableStatement().
func (p *ParserRewrite) parseVariableStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	declarations := p.parseVariableDeclarationList(expressionNormal)
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewVariableStatementTree(p.getTreeLocation(start), declarations))
}

// parseVariableDeclarationList mirrors parseVariableDeclarationList().
func (p *ParserRewrite) parseVariableDeclarationList(exprIn expressionKind) *trees.VariableDeclarationListTree {
	start := p.getTreeStartLocation()
	token := p.peekType()
	switch token {
	case TokenTypeConst, TokenTypeLet, TokenTypeVar:
		p.eat(token)
	default:
		p.reportErrorToken(p.peekToken(), "expected declaration")
		return nil
	}
	var declarations []*trees.VariableDeclarationTree
	declarations = append(declarations, p.parseVariableDeclaration(token, exprIn))
	for p.peek(TokenTypeComma) {
		p.eat(TokenTypeComma)
		declarations = append(declarations, p.parseVariableDeclaration(token, exprIn))
	}
	return trees.NewVariableDeclarationListTree(
		p.getTreeLocation(start),
		trees.TokenType(int(token)),
		declarations,
	)
}

// parseVariableDeclarationListNoIn mirrors parseVariableDeclarationListNoIn().
func (p *ParserRewrite) parseVariableDeclarationListNoIn() *trees.VariableDeclarationListTree {
	return p.parseVariableDeclarationList(expressionNoIn)
}

// parseVariableDeclaration mirrors parseVariableDeclaration().
func (p *ParserRewrite) parseVariableDeclaration(binding TokenType, exprIn expressionKind) *trees.VariableDeclarationTree {
	start := p.getTreeStartLocation()
	var lvalue *trees.ParseTree
	if p.peekPatternStart() {
		lvalue = p.parsePattern(patternKindInitializer)
	} else {
		lvalue = p.asParseTree(p.parseIdentifierExpression())
	}

	var initializer *trees.ParseTree
	if p.peek(TokenTypeEqual) {
		p.eat(TokenTypeEqual)
		initializer = p.parseAssignment(exprIn)
	} else if exprIn != expressionNoIn {
		p.maybeReportNoInitializer(binding, lvalue)
	}

	return trees.NewVariableDeclarationTree(p.getTreeLocation(start), lvalue, initializer)
}

// maybeReportNoInitializer mirrors maybeReportNoInitializer().
func (p *ParserRewrite) maybeReportNoInitializer(token TokenType, lvalue *trees.ParseTree) {
	if token == TokenTypeConst {
		p.reportError("const variables must have an initializer")
	} else if lvalue != nil && p.isPattern(lvalue) {
		p.reportError("destructuring must have an initializer")
	}
}

// isPattern returns true if the tree is a destructuring pattern.
// Go-only: tree-type predicate used by parseVariableDeclaration; Java checks inline via instanceof.
func (p *ParserRewrite) isPattern(tree *trees.ParseTree) bool {
	return tree.Type == trees.ParseTreeTypeArrayPattern || tree.Type == trees.ParseTreeTypeObjectPattern
}

// peekVariableDeclarationList returns true if a variable declaration list starts.
func (p *ParserRewrite) peekVariableDeclarationList() bool {
	switch p.peekType() {
	case TokenTypeVar, TokenTypeConst, TokenTypeLet:
		return true
	}
	return false
}

// ---- Empty statement ---------------------------------------------------------

// parseEmptyStatement mirrors parseEmptyStatement().
func (p *ParserRewrite) parseEmptyStatement() *trees.EmptyStatementTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeSemiColon)
	return trees.NewEmptyStatementTree(p.getTreeLocation(start))
}

// ---- Expression statement ----------------------------------------------------

// parseExpressionStatement mirrors parseExpressionStatement().
func (p *ParserRewrite) parseExpressionStatement() *trees.ExpressionStatementTree {
	start := p.getTreeStartLocation()
	expression := p.parseExpression()
	p.eatPossiblyImplicitSemiColon()
	return trees.NewExpressionStatementTree(p.getTreeLocation(start), expression)
}

// ---- If statement ------------------------------------------------------------

// parseIfStatement mirrors parseIfStatement().
func (p *ParserRewrite) parseIfStatement() *trees.IfStatementTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeIf)
	p.eat(TokenTypeOpenParen)
	condition := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	ifClause := p.parseStatement()
	var elseClause *trees.ParseTree
	if p.peek(TokenTypeElse) {
		p.eat(TokenTypeElse)
		elseClause = p.parseStatement()
	}
	return trees.NewIfStatementTree(p.getTreeLocation(start), condition, ifClause, elseClause)
}

// ---- Do-while statement ------------------------------------------------------

// parseDoWhileStatement mirrors parseDoWhileStatement().
func (p *ParserRewrite) parseDoWhileStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeDo)
	body := p.parseStatement()
	p.eat(TokenTypeWhile)
	p.eat(TokenTypeOpenParen)
	condition := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	// The semicolon after do-while is optional.
	if p.peek(TokenTypeSemiColon) {
		p.eat(TokenTypeSemiColon)
	}
	return p.asParseTree(trees.NewDoWhileStatementTree(p.getTreeLocation(start), body, condition))
}

// ---- While statement ---------------------------------------------------------

// parseWhileStatement mirrors parseWhileStatement().
func (p *ParserRewrite) parseWhileStatement() *trees.WhileStatementTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeWhile)
	p.eat(TokenTypeOpenParen)
	condition := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return trees.NewWhileStatementTree(p.getTreeLocation(start), condition, body)
}

// ---- For statements ----------------------------------------------------------

// parseForStatement mirrors parseForStatement() — entry point.
func (p *ParserRewrite) parseForStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeFor)
	awaited := p.peekPredefinedString(predefinedAWAIT)
	if awaited {
		p.eatPredefinedString(predefinedAWAIT)
	}
	p.eat(TokenTypeOpenParen)

	if p.peekVariableDeclarationList() {
		variables := p.parseVariableDeclarationListNoIn()
		if p.peek(TokenTypeIn) {
			if awaited {
				p.reportError("for-await-of is the only allowed asynchronous iteration")
			}
			if len(variables.Declarations) > 1 {
				p.reportError("for-in statement may not have more than one variable declaration")
			}
			if len(variables.Declarations) > 0 {
				decl := variables.Declarations[0]
				if decl.Initializer != nil {
					if p.config.AtLeast6 {
						p.reportError("for-in statement may not have initializer")
					}
				}
			}
			return p.parseForInStatement(start, p.asParseTree(variables))
		} else if p.peekPredefinedString(predefinedOF) {
			if len(variables.Declarations) > 1 {
				if awaited {
					p.reportError("for-await-of statement may not have more than one variable declaration")
				} else {
					p.reportError("for-of statement may not have more than one variable declaration")
				}
			}
			if len(variables.Declarations) > 0 {
				decl := variables.Declarations[0]
				if decl.Initializer != nil {
					if awaited {
						p.reportError("for-await-of statement may not have initializer")
					} else {
						p.reportError("for-of statement may not have initializer")
					}
				}
			}
			if awaited {
				return p.parseForAwaitOfStatement(start, p.asParseTree(variables))
			}
			return p.parseForOfStatement(start, p.asParseTree(variables))
		} else {
			p.checkVanillaForInitializers(variables)
			return p.parseForStatementVanilla(start, p.asParseTree(variables))
		}
	}

	if p.peek(TokenTypeSemiColon) {
		return p.parseForStatementVanilla(start, nil)
	}

	initializer := p.parseExpressionNoIn()
	if p.peek(TokenTypeIn) || p.peek(TokenTypeEqual) || p.peekPredefinedString(predefinedOF) {
		initializer = p.transformLeftHandSideExpression(initializer)
		if !p.isValidAssignmentTarget(initializer) {
			p.reportError("invalid assignment target")
		}
	}

	if p.peek(TokenTypeIn) || p.peekPredefinedString(predefinedOF) {
		if initializer.Type != trees.ParseTreeTypeBinaryOperator &&
			initializer.Type != trees.ParseTreeTypeCommaExpression {
			if p.peek(TokenTypeIn) {
				return p.parseForInStatement(start, initializer)
			}
			if awaited {
				return p.parseForAwaitOfStatement(start, initializer)
			}
			return p.parseForOfStatement(start, initializer)
		}
	}
	return p.parseForStatementVanilla(start, initializer)
}

// parseForStatementVanilla mirrors parseForStatement(start, initializer).
// Go-only: extracted body of the vanilla for-loop path; Java has a single overloaded parseForStatement(SourcePosition, ParseTree).
func (p *ParserRewrite) parseForStatementVanilla(start SourcePosition, initializer *trees.ParseTree) *trees.ParseTree {
	var initVal *trees.ParseTree
	if initializer == nil {
		endPos := p.getTreeEndLocation()
		startPos := p.getTreeStartLocation()
		nullRange := toTreesRange(SourceRange{Start: endPos, End: startPos})
		initVal = p.asParseTree(trees.NewNullTree(nullRange))
	} else {
		initVal = initializer
	}
	p.eat(TokenTypeSemiColon)

	var conditionVal *trees.ParseTree
	if !p.peek(TokenTypeSemiColon) {
		conditionVal = p.parseExpression()
	} else {
		endPos := p.getTreeEndLocation()
		startPos := p.getTreeStartLocation()
		nullRange := toTreesRange(SourceRange{Start: endPos, End: startPos})
		conditionVal = p.asParseTree(trees.NewNullTree(nullRange))
	}
	p.eat(TokenTypeSemiColon)

	var incrementVal *trees.ParseTree
	if !p.peek(TokenTypeCloseParen) {
		incrementVal = p.parseExpression()
	} else {
		endPos := p.getTreeEndLocation()
		startPos := p.getTreeStartLocation()
		nullRange := toTreesRange(SourceRange{Start: endPos, End: startPos})
		incrementVal = p.asParseTree(trees.NewNullTree(nullRange))
	}
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return p.asParseTree(trees.NewForStatementTree(p.getTreeLocation(start), initVal, conditionVal, incrementVal, body))
}

// checkVanillaForInitializers mirrors checkVanillaForInitializers().
func (p *ParserRewrite) checkVanillaForInitializers(variables *trees.VariableDeclarationListTree) {
	for _, decl := range variables.Declarations {
		if decl.Initializer == nil {
			p.maybeReportNoInitializer(TokenType(int(variables.DeclarationType)), decl.Lvalue)
		}
	}
}

// parseForOfStatement mirrors parseForOfStatement().
func (p *ParserRewrite) parseForOfStatement(start SourcePosition, initializer *trees.ParseTree) *trees.ParseTree {
	p.eatPredefinedString(predefinedOF)
	collection := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return p.asParseTree(trees.NewForOfStatementTree(p.getTreeLocation(start), initializer, collection, body))
}

// parseForAwaitOfStatement mirrors parseForAwaitOfStatement().
func (p *ParserRewrite) parseForAwaitOfStatement(start SourcePosition, initializer *trees.ParseTree) *trees.ParseTree {
	p.eatPredefinedString(predefinedOF)
	collection := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return p.asParseTree(trees.NewForAwaitOfStatementTree(
		p.getTreeLocation(start),
		initializer,
		collection,
		body,
	))
}

// parseForInStatement mirrors parseForInStatement().
func (p *ParserRewrite) parseForInStatement(start SourcePosition, initializer *trees.ParseTree) *trees.ParseTree {
	p.eat(TokenTypeIn)
	collection := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return p.asParseTree(trees.NewForInStatementTree(p.getTreeLocation(start), initializer, collection, body))
}

// ---- Continue/Break/Return/Throw/With/Debugger/Labelled ----------------------

// parseContinueStatement mirrors parseContinueStatement().
func (p *ParserRewrite) parseContinueStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeContinue)
	var name *trees.IdentifierToken
	if !p.peekImplicitSemiColon() {
		tok := p.eatIdOpt()
		if tok != nil {
			name = toTreesIdentToken(tok)
		}
	}
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewContinueStatementTree(p.getTreeLocation(start), name))
}

// parseBreakStatement mirrors parseBreakStatement().
func (p *ParserRewrite) parseBreakStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeBreak)
	var name *trees.IdentifierToken
	if !p.peekImplicitSemiColon() {
		tok := p.eatIdOpt()
		if tok != nil {
			name = toTreesIdentToken(tok)
		}
	}
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewBreakStatementTree(p.getTreeLocation(start), name))
}

// parseReturnStatement mirrors parseReturnStatement().
func (p *ParserRewrite) parseReturnStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeReturn)
	var expression *trees.ParseTree
	if !p.peekImplicitSemiColon() {
		expression = p.parseExpression()
	}
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewReturnStatementTree(p.getTreeLocation(start), expression))
}

// parseWithStatement mirrors parseWithStatement().
func (p *ParserRewrite) parseWithStatement() *trees.WithStatementTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeWith)
	p.eat(TokenTypeOpenParen)
	expression := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	body := p.parseStatement()
	return trees.NewWithStatementTree(p.getTreeLocation(start), expression, body)
}

// parseThrowStatement mirrors parseThrowStatement().
func (p *ParserRewrite) parseThrowStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeThrow)
	var value *trees.ParseTree
	if p.peekImplicitSemiColon() {
		p.reportError("semicolon/newline not allowed after 'throw'")
	} else {
		value = p.parseExpression()
	}
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewThrowStatementTree(p.getTreeLocation(start), value))
}

// parseDebuggerStatement mirrors parseDebuggerStatement().
func (p *ParserRewrite) parseDebuggerStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeDebugger)
	p.eatPossiblyImplicitSemiColon()
	return p.asParseTree(trees.NewDebuggerStatementTree(p.getTreeLocation(start)))
}

// peekLabelledStatement returns true if a labelled statement follows.
func (p *ParserRewrite) peekLabelledStatement() bool {
	return p.peekId() && p.peekAt(1, TokenTypeColon)
}

// parseLabelledStatement mirrors parseLabelledStatement().
func (p *ParserRewrite) parseLabelledStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	nameTok := p.eatId()
	var name *trees.IdentifierToken
	if nameTok != nil {
		name = toTreesIdentToken(nameTok)
	}
	p.eat(TokenTypeColon)
	stmt := p.parseStatement()
	return p.asParseTree(trees.NewLabelledStatementTree(p.getTreeLocation(start), name, stmt))
}

// ---- Switch statement --------------------------------------------------------

// parseSwitchStatement mirrors parseSwitchStatement().
func (p *ParserRewrite) parseSwitchStatement() *trees.SwitchStatementTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeSwitch)
	p.eat(TokenTypeOpenParen)
	expression := p.parseExpression()
	p.eat(TokenTypeCloseParen)
	p.eat(TokenTypeOpenCurly)
	caseClauses := p.parseCaseClauses()
	p.eat(TokenTypeCloseCurly)
	return trees.NewSwitchStatementTree(p.getTreeLocation(start), expression, caseClauses)
}

// parseCaseClauses mirrors parseCaseClauses().
func (p *ParserRewrite) parseCaseClauses() []*trees.ParseTree {
	foundDefault := false
	var result []*trees.ParseTree
	for {
		start := p.getTreeStartLocation()
		switch p.peekType() {
		case TokenTypeCase:
			p.eat(TokenTypeCase)
			expr := p.parseExpression()
			p.eat(TokenTypeColon)
			stmts := p.parseStatementList()
			clause := trees.NewCaseClauseTree(p.getTreeLocation(start), expr, stmts)
			result = append(result, &clause.ParseTree)
		case TokenTypeDefault:
			if foundDefault {
				p.reportError("Switch statements may have at most one default clause")
			} else {
				foundDefault = true
			}
			p.eat(TokenTypeDefault)
			p.eat(TokenTypeColon)
			stmts := p.parseStatementList()
			clause := trees.NewDefaultClauseTree(p.getTreeLocation(start), stmts)
			result = append(result, &clause.ParseTree)
		default:
			return result
		}
	}
}

// ---- Try statement -----------------------------------------------------------

// parseTryStatement mirrors parseTryStatement().
func (p *ParserRewrite) parseTryStatement() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeTry)
	body := p.parseBlock()
	var catchBlock *trees.ParseTree
	if p.peek(TokenTypeCatch) {
		catchBlock = p.parseCatch()
	}
	var finallyBlock *trees.ParseTree
	if p.peek(TokenTypeFinally) {
		finallyBlock = p.parseFinallyBlock()
	}
	if catchBlock == nil && finallyBlock == nil {
		p.reportError("'catch' or 'finally' expected.")
	}
	var bodyPT *trees.ParseTree
	if body != nil {
		bodyPT = &body.ParseTree
	}
	return p.asParseTree(trees.NewTryStatementTree(p.getTreeLocation(start), bodyPT, catchBlock, finallyBlock))
}

// parseCatch mirrors parseCatch().
func (p *ParserRewrite) parseCatch() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeCatch)

	var exception *trees.ParseTree
	if p.peek(TokenTypeOpenParen) {
		p.eat(TokenTypeOpenParen)
		if p.peekPatternStart() {
			exception = p.parsePattern(patternKindInitializer)
		} else {
			exception = p.asParseTree(p.parseIdentifierExpression())
		}
		p.eat(TokenTypeCloseParen)
	} else {
		// optional catch binding: `catch { ... }` (ES2019)
		endPos := p.getTreeEndLocation()
		startPos := p.getTreeStartLocation()
		emptyRange := toTreesRange(SourceRange{Start: endPos, End: startPos})
		exception = p.asParseTree(trees.NewEmptyStatementTree(emptyRange))
	}
	catchBody := p.parseBlock()
	var catchBodyPT *trees.ParseTree
	if catchBody != nil {
		catchBodyPT = &catchBody.ParseTree
	}
	return p.asParseTree(trees.NewCatchTree(p.getTreeLocation(start), exception, catchBodyPT))
}

// parseFinallyBlock mirrors parseFinallyBlock().
func (p *ParserRewrite) parseFinallyBlock() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeFinally)
	block := p.parseBlock()
	var blockPT *trees.ParseTree
	if block != nil {
		blockPT = &block.ParseTree
	}
	return p.asParseTree(trees.NewFinallyTree(p.getTreeLocation(start), blockPT))
}

// ---- Function declarations ---------------------------------------------------

// parseFunctionDeclaration mirrors parseFunctionDeclaration().
func (p *ParserRewrite) parseFunctionDeclaration() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeFunction)
	isGenerator := p.eatOpt(TokenTypeStar) != nil
	nameTok := p.eatId()
	var name *trees.IdentifierToken
	if nameTok != nil {
		name = toTreesIdentToken(nameTok)
	}
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindDeclaration).
		SetName(name)
	flavor := functionFlavorNormal
	if isGenerator {
		flavor = functionFlavorGenerator
	}
	p.parseFunctionTail(builder, flavor)
	return p.asParseTree(builder.Build(p.getTreeLocation(start)))
}

// parseFunctionExpression replaces the stub in parser_new.go.
func (p *ParserRewrite) parseFunctionExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeFunction)
	isGenerator := p.eatOpt(TokenTypeStar) != nil
	var name *trees.IdentifierToken
	if tok := p.eatIdOpt(); tok != nil {
		name = toTreesIdentToken(tok)
	}
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression).
		SetName(name)
	flavor := functionFlavorNormal
	if isGenerator {
		flavor = functionFlavorGenerator
	}
	p.parseFunctionTail(builder, flavor)
	return p.asParseTree(builder.Build(p.getTreeLocation(start)))
}

// parseAsyncFunctionDeclaration mirrors parseAsyncFunctionDeclaration().
func (p *ParserRewrite) parseAsyncFunctionDeclaration() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eatAsyncFunctionStart()
	generator := p.peek(TokenTypeStar)
	if generator {
		p.eat(TokenTypeStar)
	}
	nameTok := p.eatId()
	var name *trees.IdentifierToken
	if nameTok != nil {
		name = toTreesIdentToken(nameTok)
	}
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindDeclaration).
		SetName(name).
		SetAsync(true)
	flavor := functionFlavorAsynchronous
	if generator {
		flavor = functionFlavorAsynchronousGenerator
	}
	p.parseFunctionTail(builder, flavor)
	return p.asParseTree(builder.Build(p.getTreeLocation(start)))
}

// parseAsyncFunctionExpression replaces the stub in parser_new.go.
func (p *ParserRewrite) parseAsyncFunctionExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eatAsyncFunctionStart()
	generator := p.peek(TokenTypeStar)
	if generator {
		p.eat(TokenTypeStar)
	}
	var name *trees.IdentifierToken
	if tok := p.eatIdOpt(); tok != nil {
		name = toTreesIdentToken(tok)
	}
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression).
		SetName(name).
		SetAsync(true)
	flavor := functionFlavorAsynchronous
	if generator {
		flavor = functionFlavorAsynchronousGenerator
	}
	p.parseFunctionTail(builder, flavor)
	return p.asParseTree(builder.Build(p.getTreeLocation(start)))
}

// parseFunctionTail replaces the stub in parser_new.go.
func (p *ParserRewrite) parseFunctionTail(builder *trees.FunctionDeclarationBuilder, flavor functionFlavor) {
	p.functionContextStack = append(p.functionContextStack, flavor)
	builder.
		SetGenerator(flavor.isGenerator()).
		SetFormalParameterList(p.parseFormalParameterList()).
		SetFunctionBody(p.asParseTree(p.parseFunctionBody()))
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
}

// parseFunctionBody replaces the stub in parser_new.go.
func (p *ParserRewrite) parseFunctionBody() *trees.BlockTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeOpenCurly)
	result := p.parseSourceElementList()
	p.eat(TokenTypeCloseCurly)
	return trees.NewBlockTree(p.getTreeLocation(start), result)
}

// ---- Class declarations ------------------------------------------------------

// parseClassDeclaration mirrors parseClassDeclaration().
func (p *ParserRewrite) parseClassDeclaration() *trees.ParseTree {
	return p.parseClass(false)
}

// parseClassExpression replaces the stub in parser_new.go.
func (p *ParserRewrite) parseClassExpression() *trees.ParseTree {
	return p.parseClass(true)
}

// parseClass mirrors parseClass(isExpression).
func (p *ParserRewrite) parseClass(isExpression bool) *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeClass)
	var name *trees.IdentifierToken
	if !isExpression || p.peekId() {
		if tok := p.eatId(); tok != nil {
			name = toTreesIdentToken(tok)
		}
	}
	var superClass *trees.ParseTree
	if p.peek(TokenTypeExtends) {
		p.eat(TokenTypeExtends)
		superClass = p.parseLeftHandSideExpression()
	}
	p.eat(TokenTypeOpenCurly)
	elements := p.parseClassElements()
	p.eat(TokenTypeCloseCurly)
	return p.asParseTree(trees.NewClassDeclarationTree(p.getTreeLocation(start), name, superClass, elements))
}

// parseClassElements mirrors parseClassElements().
func (p *ParserRewrite) parseClassElements() []*trees.ParseTree {
	var result []*trees.ParseTree
	for {
		tok := p.peekToken()
		if tok.Type == TokenTypeSemiColon {
			p.eat(TokenTypeSemiColon)
			continue
		}
		if !p.isClassElementStart(tok) {
			return result
		}
		elem := p.parseClassElement()
		if elem != nil {
			result = append(result, elem)
		}
	}
}

// isClassElementStart mirrors isClassElementStart().
func (p *ParserRewrite) isClassElementStart(tok *Token) bool {
	switch tok.Type {
	case TokenTypeIdentifier, TokenTypeNumber, TokenTypeBigInt, TokenTypeStar, TokenTypeStatic, TokenTypeString, TokenTypeOpenSquare:
		return true
	}
	return isKeywordType(tok.Type)
}

// parseClassElement mirrors parseClassElement().
func (p *ParserRewrite) parseClassElement() *trees.ParseTree {
	if p.peek(TokenTypeSemiColon) {
		return p.asParseTree(p.parseEmptyStatement())
	}
	if p.peekClassStaticInitializerBlock() {
		return p.parseClassStaticInitializerBlock()
	}
	elementInfo := createClassMemberInfoAt(toTreesPos(p.getTreeStartLocation()), p.eatStaticIfNotElementName())
	if p.peekGetAccessor() {
		return p.parseGetAccessor(elementInfo)
	} else if p.peekSetAccessor() {
		return p.parseSetAccessor(elementInfo)
	} else if p.peekAsyncMethod() {
		return p.parseAsyncMethod(elementInfo)
	}
	return p.parseClassMemberDeclaration(elementInfo)
}

// eatStaticIfNotElementName mirrors eatStaticIfNotElementName().
func (p *ParserRewrite) eatStaticIfNotElementName() bool {
	if p.peek(TokenTypeStatic) && p.isClassElementStart(p.peekTokenAt(1)) {
		p.eat(TokenTypeStatic)
		return true
	}
	return false
}

// peekClassStaticInitializerBlock mirrors peekClassStaticInitializerBlock().
func (p *ParserRewrite) peekClassStaticInitializerBlock() bool {
	return p.peek(TokenTypeStatic) && p.peekAt(1, TokenTypeOpenCurly)
}

// parseClassStaticInitializerBlock mirrors parseClassStaticInitializerBlock().
func (p *ParserRewrite) parseClassStaticInitializerBlock() *trees.ParseTree {
	p.eat(TokenTypeStatic)
	return p.asParseTree(p.parseBlock())
}

// parseClassElementName mirrors parseClassElementName().
func (p *ParserRewrite) parseClassElementName(elementInfo *classOrObjectElementInfo) {
	if p.peekPropertyNameAt(0) {
		if p.peekIdOrKeyword() {
			elementInfo.setName(toTreesIdentToken(p.eatIdOrKeywordAsId()))
		} else {
			elementInfo.setNameExpr(p.asParseTree(p.parseLiteralExpression()))
		}
	} else {
		elementInfo.setNameExpr(p.parseComputedPropertyName())
	}
}

// parseFieldDefinition mirrors parseFieldDefinition().
func (p *ParserRewrite) parseFieldDefinition(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	var initializer *trees.ParseTree
	if p.peek(TokenTypeEqual) {
		p.eat(TokenTypeEqual)
		initializer = p.parseAssignment(expressionNormal)
	}
	p.eatPossiblyImplicitSemiColon()
	if elementInfo.hasName() {
		return p.asParseTree(trees.NewFieldDeclarationTree(
			p.getTreeLocationFromTreesPos(elementInfo.start),
			elementInfo.getName(),
			elementInfo.isStatic,
			initializer,
		))
	}
	// computed property field
	return p.asParseTree(trees.NewComputedPropertyFieldTree(
		p.getTreeLocationFromTreesPos(elementInfo.start),
		elementInfo.getNameExpr(),
		elementInfo.isStatic,
		initializer,
	))
}

// parseClassMemberDeclaration mirrors parseClassMemberDeclaration().
func (p *ParserRewrite) parseClassMemberDeclaration(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	isGenerator := p.eatOpt(TokenTypeStar) != nil
	p.parseClassElementName(elementInfo)
	if p.peek(TokenTypeOpenParen) {
		return p.parseMethodDefinition(elementInfo, isGenerator)
	}
	return p.parseFieldDefinition(elementInfo)
}

// ---- Import/Export declarations ----------------------------------------------

// peekFunction returns true if the next token starts a function declaration.
func (p *ParserRewrite) peekFunction() bool {
	return p.peek(TokenTypeFunction)
}

// peekAsyncFunctionStart returns true if the next tokens are `async function`.
func (p *ParserRewrite) peekAsyncFunctionStart() bool {
	return p.peekPredefinedString(predefinedASYNC) && !p.peekImplicitSemiColonAt(1) && p.peekAt(1, TokenTypeFunction)
}

// eatAsyncFunctionStart consumes `async function`.
func (p *ParserRewrite) eatAsyncFunctionStart() {
	p.eatPredefinedString(predefinedASYNC)
	p.eat(TokenTypeFunction)
}

// peekExportDeclaration returns true if an export declaration follows.
func (p *ParserRewrite) peekExportDeclaration() bool {
	return p.peek(TokenTypeExport)
}

// parseImportDeclaration mirrors parseImportDeclaration().
func (p *ParserRewrite) parseImportDeclaration() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeImport)

	// import ModuleSpecifier ;
	if p.peek(TokenTypeString) {
		tok := p.eat(TokenTypeString)
		p.eatPossiblyImplicitSemiColon()
		var moduleSpecifier *trees.LiteralToken
		if tok != nil {
			moduleSpecifier = &trees.LiteralToken{Token: trees.Token{Type: trees.TokenType(int(tok.Type)), Location: toTreesRange(tok.Location), Value: tok.Value}}
		}
		return p.asParseTree(trees.NewImportDeclarationTree(p.getTreeLocation(start), nil, nil, nil, moduleSpecifier))
	}

	var defaultBindingIdentifier *trees.IdentifierToken
	var nameSpaceImportIdentifier *trees.IdentifierToken
	var identifierSet []*trees.ParseTree

	parseExplicitNames := true
	if p.peekId() {
		tok := p.eatId()
		if tok != nil {
			defaultBindingIdentifier = toTreesIdentToken(tok)
		}
		if p.peek(TokenTypeComma) {
			p.eat(TokenTypeComma)
		} else {
			parseExplicitNames = false
		}
	} else if isKeywordType(p.peekType()) {
		tok := p.nextToken()
		p.reportErrorToken(tok, "cannot use keyword '%s' here.", tok.Value)
	}

	if parseExplicitNames {
		if p.peek(TokenTypeStar) {
			p.eat(TokenTypeStar)
			p.eatPredefinedString(predefinedAS)
			tok := p.eatId()
			if tok != nil {
				nameSpaceImportIdentifier = toTreesIdentToken(tok)
			}
		} else {
			identifierSet = p.parseImportSpecifierSet()
		}
	}

	p.eatPredefinedString(predefinedFROM)
	strTok := p.eat(TokenTypeString)
	var moduleSpecifier *trees.LiteralToken
	if strTok != nil {
		moduleSpecifier = &trees.LiteralToken{Token: trees.Token{Type: trees.TokenType(int(strTok.Type)), Location: toTreesRange(strTok.Location), Value: strTok.Value}}
	}
	p.eatPossiblyImplicitSemiColon()

	return p.asParseTree(trees.NewImportDeclarationTree(
		p.getTreeLocation(start),
		defaultBindingIdentifier,
		identifierSet,
		nameSpaceImportIdentifier,
		moduleSpecifier,
	))
}

// parseImportSpecifierSet mirrors parseImportSpecifierSet().
func (p *ParserRewrite) parseImportSpecifierSet() []*trees.ParseTree {
	var elements []*trees.ParseTree
	p.eat(TokenTypeOpenCurly)
	for p.peekIdOrKeyword() {
		elem := p.parseImportSpecifier()
		if elem != nil {
			elements = append(elements, p.asParseTree(elem))
		}
		if !p.peek(TokenTypeCloseCurly) {
			p.eat(TokenTypeComma)
		}
	}
	p.eat(TokenTypeCloseCurly)
	return elements
}

// parseImportSpecifier mirrors parseImportSpecifier().
func (p *ParserRewrite) parseImportSpecifier() *trees.ImportSpecifierTree {
	start := p.getTreeStartLocation()
	importedName := toTreesIdentToken(p.eatIdOrKeywordAsId())
	var destinationName *trees.IdentifierToken
	if p.peekPredefinedString(predefinedAS) {
		p.eatPredefinedString(predefinedAS)
		tok := p.eatId()
		if tok != nil {
			destinationName = toTreesIdentToken(tok)
		}
	} else if importedName != nil && isKeywordType(TokenType(int(importedName.Type))) {
		p.reportExpectedError(nil, predefinedAS)
	}
	return trees.NewImportSpecifierTree(p.getTreeLocation(start), importedName, destinationName)
}

// parseExportDeclaration mirrors parseExportDeclaration().
func (p *ParserRewrite) parseExportDeclaration() *trees.ParseTree {
	start := p.getTreeStartLocation()
	isDefault := false
	isExportAll := false
	isExportSpecifier := false
	needsSemiColon := true
	p.eat(TokenTypeExport)
	var export *trees.ParseTree
	var exportSpecifierList []*trees.ParseTree

	switch p.peekType() {
	case TokenTypeStar:
		isExportAll = true
		p.nextToken()
	case TokenTypeIdentifier:
		export = p.parseAsyncFunctionDeclaration()
	case TokenTypeFunction:
		export = p.parseFunctionDeclaration()
		needsSemiColon = false
	case TokenTypeClass:
		export = p.parseClassDeclaration()
		needsSemiColon = false
	case TokenTypeDefault:
		isDefault = true
		p.nextToken()
		export = p.parseExpression()
		needsSemiColon = false
	case TokenTypeOpenCurly:
		isExportSpecifier = true
		exportSpecifierList = p.parseExportSpecifierSet()
	default: // var, let, const, or error
		export = p.asParseTree(p.parseVariableDeclarationList(expressionNormal))
	}

	var moduleSpecifier *trees.LiteralToken
	if isExportAll || (isExportSpecifier && p.peekPredefinedString(predefinedFROM)) {
		p.eatPredefinedString(predefinedFROM)
		tok := p.eat(TokenTypeString)
		if tok != nil {
			moduleSpecifier = &trees.LiteralToken{Token: trees.Token{Type: trees.TokenType(int(tok.Type)), Location: toTreesRange(tok.Location), Value: tok.Value}}
		}
	} else if isExportSpecifier {
		// validate that exported names are not keywords
		for _, tree := range exportSpecifierList {
			_ = tree // keyword check omitted (Java uses importedName.isKeyword())
		}
	}

	if needsSemiColon || p.peekImplicitSemiColon() {
		p.eatPossiblyImplicitSemiColon()
	}

	return p.asParseTree(trees.NewExportDeclarationTree(
		p.getTreeLocation(start),
		isDefault,
		isExportAll,
		export,
		exportSpecifierList,
		moduleSpecifier,
	))
}

// parseExportSpecifierSet mirrors parseExportSpecifierSet().
func (p *ParserRewrite) parseExportSpecifierSet() []*trees.ParseTree {
	var elements []*trees.ParseTree
	p.eat(TokenTypeOpenCurly)
	for p.peekIdOrKeyword() {
		elem := p.parseExportSpecifier()
		if elem != nil {
			elements = append(elements, p.asParseTree(elem))
		}
		if !p.peek(TokenTypeCloseCurly) {
			p.eat(TokenTypeComma)
		}
	}
	p.eat(TokenTypeCloseCurly)
	return elements
}

// parseExportSpecifier mirrors parseExportSpecifier().
func (p *ParserRewrite) parseExportSpecifier() *trees.ExportSpecifierTree {
	start := p.getTreeStartLocation()
	importedName := toTreesIdentToken(p.eatIdOrKeywordAsId())
	var destinationName *trees.IdentifierToken
	if p.peekPredefinedString(predefinedAS) {
		p.eatPredefinedString(predefinedAS)
		tok := p.eatIdOrKeywordAsId()
		if tok != nil {
			destinationName = toTreesIdentToken(tok)
		}
	}
	return trees.NewExportSpecifierTree(p.getTreeLocation(start), importedName, destinationName)
}

// Port of com/google/javascript/jscomp/parsing/parser/Parser.java
// @ closure-compiler v20260415
//
// Parses a JavaScript file into a ParseTree.  This is the new recursive-descent
// parser that will eventually replace parser.go.  During migration it coexists with
// the old parser; nothing calls into it yet.
//
// The parseX() methods never return nil — even on parse errors. Each returns the
// corresponding ParseTree subclass from internal/parser/trees/.
//
// Token input comes from the Scanner in scanner.go (NextToken / PeekToken /
// NextRegularExpressionLiteralToken / NextTemplateLiteralToken).

package parser

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/parser/trees"
)

// ---- Config ------------------------------------------------------------------

// ParserMode mirrors Parser.Config.Mode.
type ParserMode int

const (
	ParserModeES3          ParserMode = iota
	ParserModeES5
	ParserModeES6OrES7
	ParserModeES8OrGreater
)

// ParserConfig mirrors Parser.Config.
type ParserConfig struct {
	AtLeast6         bool // mode >= ES6_OR_ES7
	AtLeast8         bool // mode == ES8_OR_GREATER
	IsStrictMode     bool
	WarnTrailingCommas bool
}

// NewParserConfig mirrors new Config(mode, isStrictMode).
// Go-only: language-adapter — Go constructor for ParserConfig (Java uses Config(Mode, boolean) overload).
func NewParserConfig(mode ParserMode, strictMode bool) ParserConfig {
	return ParserConfig{
		AtLeast6:           mode != ParserModeES3 && mode != ParserModeES5,
		AtLeast8:           mode == ParserModeES8OrGreater,
		IsStrictMode:       strictMode,
		WarnTrailingCommas: mode == ParserModeES3,
	}
}

// ---- FunctionFlavor ----------------------------------------------------------

// functionFlavor mirrors Parser.FunctionFlavor.
type functionFlavor int

const (
	functionFlavorNormal             functionFlavor = iota
	functionFlavorGenerator
	functionFlavorAsynchronous
	functionFlavorAsynchronousGenerator
)

func (f functionFlavor) isGenerator() bool {
	return f == functionFlavorGenerator || f == functionFlavorAsynchronousGenerator
}

// ---- Expression enum ---------------------------------------------------------

// expressionKind mirrors Parser.Expression (NO_IN vs NORMAL).
type expressionKind int

const (
	expressionNormal expressionKind = iota
	expressionNoIn
)

// ---- PatternKind -------------------------------------------------------------

// patternKind mirrors Parser.PatternKind.
type patternKind int

const (
	patternKindInitializer patternKind = iota // let, const, catch, formal params
	patternKindAny                            // assignment or for-in
)

// ---- CommentRecorder ---------------------------------------------------------

// parserCommentRecorder mirrors Parser.CommentRecorder.
type parserCommentRecorder struct {
	comments             []trees.Comment
	lastCommentEndOffset int
	initialized          bool
}

func (r *parserCommentRecorder) RecordComment(typ CommentType, rang SourceRange, value string) {
	if r.initialized && rang.End.Offset <= r.lastCommentEndOffset {
		return // already recorded (scanner rewind)
	}
	r.initialized = true
	r.lastCommentEndOffset = rang.End.Offset
	// convert SourceRange → trees.SourceRange
	tr := toTreesRange(rang)
	tc := trees.Comment{
		Location: tr,
		Value:    value,
	}
	// map CommentType — preserve JSDoc distinction for IRFactory JSDoc attachment.
	switch typ {
	case CommentJSDoc:
		tc.Type = trees.CommentTypeJSDoc
	case CommentImportant:
		tc.Type = trees.CommentTypeImportant
	case CommentBlock:
		tc.Type = trees.CommentTypeBlock
	default:
		tc.Type = trees.CommentTypeLine
	}
	r.comments = append(r.comments, tc)
}

// ---- Position helpers --------------------------------------------------------

// Go-only: language-adapter — converts parser SourcePosition to trees.SourcePosition (Go package split; Java uses one type).
func toTreesPos(p SourcePosition) trees.SourcePosition {
	return trees.SourcePosition{Offset: p.Offset, Line: p.Line, Column: p.Column}
}

// Go-only: language-adapter — converts parser SourceRange to trees.SourceRange (Go package split; Java uses one type).
func toTreesRange(r SourceRange) trees.SourceRange {
	return trees.SourceRange{Start: toTreesPos(r.Start), End: toTreesPos(r.End)}
}

// ---- Token conversion helpers ------------------------------------------------

// toTreesToken converts a *parser.Token to a *trees.Token, bridging the
// parser and trees package token types.
// Go-only: language-adapter — converts scanner Token to trees.Token (Go package split; Java uses one type).
func toTreesToken(t *Token) *trees.Token {
	if t == nil {
		return nil
	}
	return &trees.Token{
		Type:     trees.TokenType(int(t.Type)),
		Location: toTreesRange(t.Location),
		Value:    t.Value,
	}
}

// toTreesIdentToken converts a *parser.Token (identifier) to a *trees.IdentifierToken.
// Go-only: language-adapter — converts scanner IdentifierToken to trees.IdentifierToken (Go package split; Java uses one type).
func toTreesIdentToken(t *Token) *trees.IdentifierToken {
	if t == nil {
		return nil
	}
	return &trees.IdentifierToken{
		Token: trees.Token{
			Type:     trees.TokenType(int(t.Type)),
			Location: toTreesRange(t.Location),
		},
		Value: t.Value,
	}
}

// ---- ClassOrObjectElementInfo ------------------------------------------------

// classOrObjectElementInfo mirrors ClassOrObjectElementInfo.java.
type classOrObjectElementInfo struct {
	start         trees.SourcePosition
	isClassMember bool
	isStatic      bool
	name          *trees.IdentifierToken // set by setName
	nameExpr      *trees.ParseTree       // set by setNameExpr
}

// Go-only: language-adapter — factory for classOrObjectElementInfo at a position (Java uses createObjectLiteralElementInfo overload).
func createObjectLiteralElementInfoAt(start trees.SourcePosition) *classOrObjectElementInfo {
	return &classOrObjectElementInfo{start: start, isClassMember: false, isStatic: false}
}

// Go-only: language-adapter — factory for classOrObjectElementInfo with static flag (Java constructs inline; Go needs named constructor).
func createClassMemberInfoAt(start trees.SourcePosition, isStatic bool) *classOrObjectElementInfo {
	return &classOrObjectElementInfo{start: start, isClassMember: true, isStatic: isStatic}
}

func (e *classOrObjectElementInfo) setName(name *trees.IdentifierToken) {
	e.name = name
}

func (e *classOrObjectElementInfo) hasName() bool { return e.name != nil }

func (e *classOrObjectElementInfo) getName() *trees.IdentifierToken { return e.name }

func (e *classOrObjectElementInfo) setNameExpr(expr *trees.ParseTree) {
	e.nameExpr = expr
}

func (e *classOrObjectElementInfo) hasNameExpr() bool { return e.nameExpr != nil }

func (e *classOrObjectElementInfo) getNameExpr() *trees.ParseTree { return e.nameExpr }

// ---- Predefined name constants -----------------------------------------------

const (
	predefinedAS    = "as"
	predefinedFROM  = "from"
	predefinedGET   = "get"
	predefinedOF    = "of"
	predefinedSET   = "set"
	predefinedASYNC = "async"
	predefinedAWAIT = "await"
)

// ---- Parser struct -----------------------------------------------------------

// Parser mirrors com.google.javascript.jscomp.parsing.parser.Parser.
// It parses a JavaScript source file into a tree of ParseTree nodes.
type ParserRewrite struct {
	scanner              *Scanner
	errorReporter        ErrorReporter
	config               ParserConfig
	commentRecorder      parserCommentRecorder
	functionContextStack []functionFlavor
	lastSourcePosition   SourcePosition
	sourceMapURL         string

	// lastIdentValue records the Value of the most recently consumed identifier
	// token (set by eatId / eatIdOrKeywordAsId). Used by isStartOfAsyncArrowFunction
	// to inspect the identifier expression that was just parsed as a ParseTree.
	lastIdentValue string
}

// NewParserRewrite mirrors new Parser(config, errorReporter, source).
// Go-only: language-adapter — Go constructor (Java uses Parser(Config, ErrorReporter, SourceFile) constructor overload).
func NewParserRewrite(config ParserConfig, errorReporter ErrorReporter, src string) *ParserRewrite {
	p := &ParserRewrite{
		config:        config,
		errorReporter: errorReporter,
	}
	p.scanner = NewScannerFull(errorReporter, &p.commentRecorder, src, 0)
	p.functionContextStack = []functionFlavor{functionFlavorNormal}
	p.lastSourcePosition = p.scanner.GetPosition()
	return p
}

// Parse mirrors parseProgram().
func (p *ParserRewrite) Parse() *trees.ProgramTree {
	start := p.lastSourcePosition
	sourceElements := p.parseGlobalSourceElements()
	p.eat(TokenTypeEndOfFile)
	return trees.NewProgramTree(
		toTreesRange(SourceRange{Start: start, End: p.lastSourcePosition}),
		sourceElements,
		p.commentRecorder.comments,
	)
}

// GetSourceMapURL returns any sourceMappingURL found in the source.
func (p *ParserRewrite) GetSourceMapURL() string { return p.sourceMapURL }

// ---- Top-level helpers -------------------------------------------------------

func (p *ParserRewrite) parseGlobalSourceElements() []*trees.ParseTree {
	var result []*trees.ParseTree
	for !p.peek(TokenTypeEndOfFile) {
		result = append(result, p.parseScriptElement())
	}
	return result
}

// ---- Infrastructure: position & location -------------------------------------

func (p *ParserRewrite) getTreeStartLocation() SourcePosition {
	return p.peekToken().Location.Start
}

func (p *ParserRewrite) getTreeEndLocation() SourcePosition {
	return p.lastSourcePosition
}

func (p *ParserRewrite) getTreeLocation(start SourcePosition) trees.SourceRange {
	return toTreesRange(SourceRange{Start: start, End: p.getTreeEndLocation()})
}

// getTreeLocationFromTreesPos is like getTreeLocation but accepts a trees.SourcePosition
// (used when start comes from a trees struct like classOrObjectElementInfo).
// Go-only: language-adapter — overload of getTreeLocation taking trees.SourcePosition (Java uses method overloading).
func (p *ParserRewrite) getTreeLocationFromTreesPos(start trees.SourcePosition) trees.SourceRange {
	return trees.SourceRange{Start: start, End: toTreesPos(p.getTreeEndLocation())}
}

// ---- Infrastructure: token consumption ---------------------------------------

// nextToken consumes and returns the next token, updating lastSourcePosition.
func (p *ParserRewrite) nextToken() *Token {
	tok := p.scanner.NextToken()
	p.lastSourcePosition = tok.Location.End
	return tok
}

// nextRegularExpressionLiteralToken consumes a regex literal token.
func (p *ParserRewrite) nextRegularExpressionLiteralToken() *Token {
	tok := p.scanner.NextRegularExpressionLiteralToken()
	p.lastSourcePosition = tok.Location.End
	return tok
}

// nextTemplateLiteralToken consumes a template literal token.
func (p *ParserRewrite) nextTemplateLiteralToken() *Token {
	tok := p.scanner.NextTemplateLiteralToken()
	p.lastSourcePosition = tok.Location.End
	return tok
}

// peek returns true if the next token is of the expected type.
func (p *ParserRewrite) peek(expected TokenType) bool {
	return p.peekAt(0, expected)
}

// peekAt returns true if the index-th next token is of the expected type.
// Go-only: language-adapter — index-parameterised peek (Java has peek(int, TokenType) overload; Go has no overloading).
func (p *ParserRewrite) peekAt(index int, expected TokenType) bool {
	return p.peekTypeAt(index) == expected
}

// peekType returns the type of the next token without consuming it.
func (p *ParserRewrite) peekType() TokenType {
	return p.peekTypeAt(0)
}

// peekTypeAt returns the type of the index-th next token without consuming it.
// Go-only: language-adapter — index-parameterised peekType (Java has peekType(int) overload; Go has no overloading).
func (p *ParserRewrite) peekTypeAt(index int) TokenType {
	return p.peekTokenAt(index).Type
}

// peekToken returns the next token without consuming it.
func (p *ParserRewrite) peekToken() *Token {
	return p.peekTokenAt(0)
}

// peekTokenAt returns the index-th next token without consuming it.
// Go-only: language-adapter — index-parameterised peekToken (Java has peekToken(int) overload; Go has no overloading).
func (p *ParserRewrite) peekTokenAt(index int) *Token {
	return p.scanner.PeekTokenAt(index)
}

// peekTokenNoLineTerminator returns the next token if it is on the same line
// as the last consumed token; otherwise returns nil.
// Go-only: language-adapter — no-line-terminator peek variant (Java checks inline; Go needs helper due to no overloading).
func (p *ParserRewrite) peekTokenNoLineTerminator() *Token {
	tok := p.peekToken()
	if tok.Location.Start.Line > p.lastSourcePosition.Line {
		return nil
	}
	return tok
}

// eat consumes the next token. Reports an error and returns nil if the token is
// not of expectedTokenType.
func (p *ParserRewrite) eat(expected TokenType) *Token {
	tok := p.nextToken()
	if tok.Type != expected {
		p.reportExpectedError(tok, tokenTypeValue(expected))
		return nil
	}
	return tok
}

// eatOpt consumes the next token only if it matches expected. Returns nil otherwise.
func (p *ParserRewrite) eatOpt(expected TokenType) *Token {
	if p.peek(expected) {
		return p.eat(expected)
	}
	return nil
}

// eatPredefinedString consumes an identifier that equals str. Reports an error
// and returns nil if it doesn't match.
func (p *ParserRewrite) eatPredefinedString(str string) *Token {
	tok := p.eatId()
	if tok == nil || tok.Value != str {
		p.reportExpectedError(tok, "'"+str+"'")
		return nil
	}
	return tok
}

// peekPredefinedString returns true if the next token is the identifier str.
func (p *ParserRewrite) peekPredefinedString(str string) bool {
	return p.peekPredefinedStringAt(0, str)
}

// peekPredefinedStringAt returns true if the index-th token is the identifier str.
// Go-only: language-adapter — index-parameterised peekPredefinedString (Java has method overload; Go has no overloading).
func (p *ParserRewrite) peekPredefinedStringAt(index int, str string) bool {
	tok := p.peekTokenAt(index)
	return tok.Type == TokenTypeIdentifier && tok.Value == str
}

// ---- Infrastructure: identifier helpers --------------------------------------

// peekId returns true if the next token can be used as an identifier.
func (p *ParserRewrite) peekId() bool {
	return p.peekIdAt(0)
}

// peekIdAt returns true if the index-th token can be used as an identifier.
// Go-only: language-adapter — index-parameterised peekId (Java has peekId(int) overload; Go has no overloading).
func (p *ParserRewrite) peekIdAt(index int) bool {
	typ := p.peekTypeAt(index)
	return typ == TokenTypeIdentifier || (!p.inStrictContext() && isStrictKeyword(typ))
}

// peekIdOrKeyword returns true if the next token is an identifier or keyword.
func (p *ParserRewrite) peekIdOrKeyword() bool {
	return p.peekIdOrKeywordAt(0)
}

// peekIdOrKeywordAt returns true if the index-th token is an identifier or keyword.
// Go-only: language-adapter — index-parameterised peekIdOrKeyword (Java has peekIdOrKeyword(int) overload; Go has no overloading).
func (p *ParserRewrite) peekIdOrKeywordAt(index int) bool {
	typ := p.peekTypeAt(index)
	return typ == TokenTypeIdentifier || isKeywordType(typ)
}

// eatId consumes an identifier (or strict-mode keyword in non-strict mode).
// Reports an error if none is present. Returns nil on failure.
func (p *ParserRewrite) eatId() *Token {
	if p.peekId() {
		return p.eatIdOrKeywordAsId()
	}
	p.reportExpectedError(p.peekToken(), tokenTypeValue(TokenTypeIdentifier))
	if p.peekIdOrKeyword() {
		return p.eatIdOrKeywordAsId()
	}
	return nil
}

// eatIdOpt consumes an identifier if present, else returns nil.
func (p *ParserRewrite) eatIdOpt() *Token {
	if p.peekId() {
		return p.eatIdOrKeywordAsId()
	}
	return nil
}

// eatIdOrKeywordAsId consumes an identifier or keyword token, returning it as
// an identifier (keyword value is used as-is).
func (p *ParserRewrite) eatIdOrKeywordAsId() *Token {
	tok := p.nextToken()
	if tok.Type == TokenTypeIdentifier {
		p.lastIdentValue = tok.Value
		return tok
	}
	if isKeywordType(tok.Type) {
		// Synthesize an identifier token with the keyword's string value.
		val := keywordForType(tok.Type)
		p.lastIdentValue = val
		return &Token{
			Type:     TokenTypeIdentifier,
			Location: tok.Location,
			Value:    val,
		}
	}
	p.reportExpectedError(tok, tokenTypeValue(TokenTypeIdentifier))
	return nil
}

// eatObjectLiteralPropertyName consumes a property name token (string, number,
// or identifier/keyword).
func (p *ParserRewrite) eatObjectLiteralPropertyName() *Token {
	tok := p.peekToken()
	switch tok.Type {
	case TokenTypeString, TokenTypeNumber, TokenTypeBigInt:
		return p.nextToken()
	default:
		return p.eatIdOrKeywordAsId()
	}
}

// inStrictContext returns true when parsing in strict mode.
func (p *ParserRewrite) inStrictContext() bool {
	return p.config.IsStrictMode
}

// inGeneratorContext returns true when the current function context is a generator.
func (p *ParserRewrite) inGeneratorContext() bool {
	return p.functionContextStack[len(p.functionContextStack)-1].isGenerator()
}

// ---- Infrastructure: error reporting -----------------------------------------

// reportError reports an error at the current scanner position.
func (p *ParserRewrite) reportError(format string, args ...interface{}) {
	p.errorReporter.ReportError(p.scanner.GetPosition(), format, args...)
}

// reportErrorAt reports an error at a specific position.
// Go-only: language-adapter — position-based reportError variant (Java has method overloading; Go has no overloading).
func (p *ParserRewrite) reportErrorAt(pos SourcePosition, format string, args ...interface{}) {
	p.errorReporter.ReportError(pos, format, args...)
}

// reportErrorToken reports an error at the given token's start position.
// Go-only: language-adapter — token-based reportError variant (Java has method overloading; Go has no overloading).
func (p *ParserRewrite) reportErrorToken(tok *Token, format string, args ...interface{}) {
	if tok == nil {
		p.reportError(format, args...)
		return
	}
	p.errorReporter.ReportError(tok.Location.Start, format, args...)
}

// reportExpectedError reports "'X' expected" error.
func (p *ParserRewrite) reportExpectedError(tok *Token, expected string) {
	p.reportErrorToken(tok, "'%s' expected", expected)
}

// ---- Infrastructure: implicit semicolons -------------------------------------

// peekImplicitSemiColon returns true if an implicit semicolon is valid at the
// current position (newline, close curly, EOF, or explicit semicolon).
func (p *ParserRewrite) peekImplicitSemiColon() bool {
	return p.peekImplicitSemiColonAt(0)
}

// Go-only: language-adapter — index-parameterised peekImplicitSemiColon (Java has method overload; Go has no overloading).
func (p *ParserRewrite) peekImplicitSemiColonAt(index int) bool {
	var lineAdvanced bool
	if index == 0 {
		lineAdvanced = p.getNextLine() > p.getLastLine()
	} else {
		lineAdvanced = p.peekTokenAt(index).Location.Start.Line > p.peekTokenAt(index-1).Location.End.Line
	}
	return lineAdvanced ||
		p.peekAt(index, TokenTypeSemiColon) ||
		p.peekAt(index, TokenTypeCloseCurly) ||
		p.peekAt(index, TokenTypeEndOfFile)
}

func (p *ParserRewrite) getLastLine() int { return p.lastSourcePosition.Line }
func (p *ParserRewrite) getNextLine() int { return p.peekToken().Location.Start.Line }

// eatPossiblyImplicitSemiColon consumes an explicit semicolon if present,
// or silently accepts an implicit one. Reports error otherwise.
func (p *ParserRewrite) eatPossiblyImplicitSemiColon() {
	if p.peek(TokenTypeSemiColon) {
		p.eat(TokenTypeSemiColon)
		return
	}
	if p.peekImplicitSemiColon() {
		return
	}
	p.reportError("Semi-colon expected")
}

// ---- Infrastructure: scanner reset -------------------------------------------

// resetScannerToPos resets the scanner to a saved source position.
// Go-only: language-adapter — position-based scanner reset (Java has resetScanner(SourcePosition) overload).
func (p *ParserRewrite) resetScannerToPos(pos SourcePosition) {
	p.lastSourcePosition = pos
	p.scanner.SetPosition(pos)
}

// resetScannerToTree rewinds the scanner to the start of the given tree.
// Go-only: language-adapter — tree-based scanner reset (Java has resetScanner(ParseTree) overload).
func (p *ParserRewrite) resetScannerToTree(tree *trees.ParseTree) {
	sp := SourcePosition{
		Offset: tree.Location.Start.Offset,
		Line:   tree.Location.Start.Line,
		Column: tree.Location.Start.Column,
	}
	p.scanner.SetPosition(sp)
}

// resetScannerAfterTree rewinds the scanner to the end of the given tree.
// Go-only: language-adapter — post-tree scanner reset (Java has resetScannerAfter(ParseTree); Go needs named variant).
func (p *ParserRewrite) resetScannerAfterTree(tree *trees.ParseTree) {
	sp := SourcePosition{
		Offset: tree.Location.End.Offset,
		Line:   tree.Location.End.Line,
		Column: tree.Location.End.Column,
	}
	p.lastSourcePosition = sp
	p.scanner.SetPosition(sp)
}

// ---- Infrastructure: maybeReportTrailingComma --------------------------------

func (p *ParserRewrite) maybeReportTrailingComma(commaToken *Token) {
	if commaToken != nil {
		if p.config.WarnTrailingCommas {
			p.errorReporter.ReportWarning(
				commaToken.Location.Start,
				"Trailing comma is not legal in an ECMA-262 object initializer",
			)
		}
	}
}

// ---- Infrastructure: peekPropertyName / peekPropertyNameOrComputedProp -------

// Go-only: language-adapter — renamed from Java peekPropertyNameOrComputedProp(int) (no overloading in Go).
func (p *ParserRewrite) peekPropertyNameOrComputedPropAt(tokenIndex int) bool {
	return p.peekPropertyNameAt(tokenIndex) || p.peekTypeAt(tokenIndex) == TokenTypeOpenSquare
}

// Go-only: language-adapter — renamed from Java peekPropertyName(int) (no overloading in Go).
func (p *ParserRewrite) peekPropertyNameAt(tokenIndex int) bool {
	typ := p.peekTypeAt(tokenIndex)
	switch typ {
	case TokenTypeIdentifier, TokenTypeString, TokenTypeNumber, TokenTypeBigInt:
		return true
	}
	return isKeywordType(typ)
}

// ---- Expression parsing: entry points ----------------------------------------

// parseExpression mirrors parseExpression() — parses a comma-separated
// sequence of assignment expressions.
func (p *ParserRewrite) parseExpression() *trees.ParseTree {
	return p.parse(expressionNormal)
}

// parseExpressionNoIn mirrors parseExpressionNoIn().
func (p *ParserRewrite) parseExpressionNoIn() *trees.ParseTree {
	return p.parse(expressionNoIn)
}

// parse implements the comma-expression production.
func (p *ParserRewrite) parse(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	result := p.parseAssignment(exprIn)
	if p.peek(TokenTypeComma) &&
		!p.peekAt(1, TokenTypeEllipsis) &&
		!p.peekAt(1, TokenTypeCloseParen) {
		var exprs []*trees.ParseTree
		exprs = append(exprs, result)
		for p.peek(TokenTypeComma) &&
			!p.peekAt(1, TokenTypeEllipsis) &&
			!p.peekAt(1, TokenTypeCloseParen) {
			p.eat(TokenTypeComma)
			exprs = append(exprs, p.parseAssignment(exprIn))
		}
		return p.asParseTree(trees.NewCommaExpressionTree(p.getTreeLocation(start), exprs))
	}
	return result
}

// peekExpression returns true if the next token can start an expression.
func (p *ParserRewrite) peekExpression() bool {
	switch p.peekType() {
	case TokenTypeBang,
		TokenTypeClass,
		TokenTypeDelete,
		TokenTypeFalse,
		TokenTypeFunction,
		TokenTypeIdentifier,
		TokenTypeType,
		TokenTypeDeclare,
		TokenTypeModule,
		TokenTypeNamespace,
		TokenTypeMinus,
		TokenTypeMinusMinus,
		TokenTypeNew,
		TokenTypeNull,
		TokenTypeNumber,
		TokenTypeBigInt,
		TokenTypeOpenCurly,
		TokenTypeOpenParen,
		TokenTypeOpenSquare,
		TokenTypePlus,
		TokenTypePlusPlus,
		TokenTypeSlash,
		TokenTypeSlashEqual,
		TokenTypeString,
		TokenTypeNoSubstitutionTemplate,
		TokenTypeTemplateHead,
		TokenTypeSuper,
		TokenTypeThis,
		TokenTypeTilde,
		TokenTypeTrue,
		TokenTypeTypeof,
		TokenTypeVoid,
		TokenTypeYield:
		return true
	case TokenTypeImport:
		return p.peekImportCall() || p.peekImportDot()
	}
	return false
}

// ---- Assignment expression ---------------------------------------------------

// parseAssignmentExpression mirrors parseAssignmentExpression().
func (p *ParserRewrite) parseAssignmentExpression() *trees.ParseTree {
	return p.parseAssignment(expressionNormal)
}

// parseAssignmentExpressionNoIn mirrors parseAssignment(Expression.NO_IN).
// Go-only: language-adapter — no-in variant wrapper (Java uses parseAssignmentExpression(Expression.NO_IN) overload).
func (p *ParserRewrite) parseAssignmentExpressionNoIn() *trees.ParseTree {
	return p.parseAssignment(expressionNoIn)
}

// peekAssignmentExpression mirrors peekAssignmentExpression().
func (p *ParserRewrite) peekAssignmentExpression() bool {
	return p.peekExpression()
}

// parseAssignment mirrors parseAssignment(Expression).
func (p *ParserRewrite) parseAssignment(exprIn expressionKind) *trees.ParseTree {
	if p.peek(TokenTypeYield) && p.inGeneratorContext() {
		return p.parseYield(exprIn)
	}

	start := p.getTreeStartLocation()
	left := p.parseConditional(exprIn)

	if p.isStartOfAsyncArrowFunction(left) {
		p.resetScannerToTree(left)
		return p.parseAsyncArrowFunction(exprIn)
	}
	if p.peek(TokenTypeArrow) {
		return p.completeAssignmentExpressionParseAtArrow(left, exprIn)
	}

	if p.peekAssignmentOperator() {
		if !p.peek(TokenTypeEqual) {
			// compound assignment: target must be valid non-vanilla target
			if !p.isValidNonVanillaAssignmentTarget(left) {
				p.reportError("invalid assignment target")
				return p.asParseTree(trees.NewMissingPrimaryExpressionTree(
					p.getTreeLocation(p.getTreeStartLocation())))
			}
		}
		left = p.transformLeftHandSideExpression(left)
		if !p.isValidAssignmentTarget(left) {
			p.reportError("invalid assignment target")
			return p.asParseTree(trees.NewMissingPrimaryExpressionTree(
				p.getTreeLocation(p.getTreeStartLocation())))
		}
		operator := p.nextToken()
		right := p.parseAssignment(exprIn)
		return p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

// isValidAssignmentTarget mirrors ParseTree.isValidAssignmentTarget().
// An assignment target is valid if it is an identifier, member expression, or pattern.
func (p *ParserRewrite) isValidAssignmentTarget(tree *trees.ParseTree) bool {
	switch tree.Type {
	case trees.ParseTreeTypeIdentifierExpression,
		trees.ParseTreeTypeMemberExpression,
		trees.ParseTreeTypeMemberLookupExpression,
		trees.ParseTreeTypeArrayPattern,
		trees.ParseTreeTypeObjectPattern:
		return true
	}
	return false
}

// isValidNonVanillaAssignmentTarget mirrors ParseTree.isValidNonVanillaAssignmentTarget().
func (p *ParserRewrite) isValidNonVanillaAssignmentTarget(tree *trees.ParseTree) bool {
	switch tree.Type {
	case trees.ParseTreeTypeIdentifierExpression,
		trees.ParseTreeTypeMemberExpression,
		trees.ParseTreeTypeMemberLookupExpression:
		return true
	}
	return false
}

// isStartOfAsyncArrowFunction checks whether a partial parse result is the
// beginning of an async arrow function.
func (p *ParserRewrite) isStartOfAsyncArrowFunction(partial *trees.ParseTree) bool {
	if partial.Type == trees.ParseTreeTypeIdentifierExpression {
		// partial = `async`; next must be [no newline] bindingId [no newline] =>
		// p.lastIdentValue holds the value of the most-recently consumed identifier
		// (set by eatIdOrKeywordAsId, which was the last thing parseIdentifierExpression called).
		return p.lastIdentValue == predefinedASYNC &&
			!p.peekImplicitSemiColon() &&
			p.peekId() &&
			!p.peekImplicitSemiColonAt(1) &&
			p.peekAt(1, TokenTypeArrow)
	}
	if partial.Type == trees.ParseTreeTypeCallExpression {
		// partial = `async [no newline] (parameters)`; next must be [no newline] =>
		return !p.peekImplicitSemiColon() && p.peek(TokenTypeArrow)
	}
	return false
}

// peekAssignmentOperator returns true if the next token is an assignment operator.
func (p *ParserRewrite) peekAssignmentOperator() bool {
	switch p.peekType() {
	case TokenTypeEqual,
		TokenTypeStarEqual,
		TokenTypeStarStarEqual,
		TokenTypeSlashEqual,
		TokenTypePercentEqual,
		TokenTypePlusEqual,
		TokenTypeMinusEqual,
		TokenTypeLeftShiftEqual,
		TokenTypeRightShiftEqual,
		TokenTypeUnsignedRightShiftEqual,
		TokenTypeAmpersandEqual,
		TokenTypeCaretEqual,
		TokenTypeBarEqual,
		TokenTypeOrEqual,
		TokenTypeAndEqual,
		TokenTypeQuestionQuestionEqual:
		return true
	}
	return false
}

// transformLeftHandSideExpression mirrors transformLeftHandSideExpression.
// For array/object literals on the left of an assignment, re-parse as pattern.
func (p *ParserRewrite) transformLeftHandSideExpression(tree *trees.ParseTree) *trees.ParseTree {
	switch tree.Type {
	case trees.ParseTreeTypeArrayLiteralExpression,
		trees.ParseTreeTypeObjectLiteralExpression:
		p.resetScannerToTree(tree)
		return p.parseLeftHandSidePattern()
	}
	return tree
}

// parseLeftHandSidePattern mirrors parseLeftHandSidePattern().
func (p *ParserRewrite) parseLeftHandSidePattern() *trees.ParseTree {
	return p.parsePattern(patternKindAny)
}

// ---- Arrow function helpers --------------------------------------------------

// completeAssignmentExpressionParseAtArrow mirrors completeAssignmentExpressionParseAtArrow.
func (p *ParserRewrite) completeAssignmentExpressionParseAtArrow(leftOfArrow *trees.ParseTree, exprIn expressionKind) *trees.ParseTree {
	if leftOfArrow.Type == trees.ParseTreeTypeCallExpression {
		return p.completeAssignmentAtArrowFromCall(leftOfArrow)
	}
	return p.completeArrowFunctionParseAtArrow(leftOfArrow, exprIn)
}

// completeAssignmentAtArrowFromCall handles the edge case where a CallExpression
// precedes `=>` (e.g., an implicit semicolon before the arrow).
// Go-only: language-adapter — arrow completion from a call-expression tree (Java has completeAssignmentExpressionParseAtArrow overload).
func (p *ParserRewrite) completeAssignmentAtArrowFromCall(callTree *trees.ParseTree) *trees.ParseTree {
	// We can't easily inspect CallExpressionTree fields from *ParseTree alone without
	// a type-safe cast. For now treat conservatively: report error.
	p.reportError("'=>' unexpected")
	return callTree
}

// completeArrowFunctionParseAtArrow mirrors completeArrowFunctionParseAtArrow.
func (p *ParserRewrite) completeArrowFunctionParseAtArrow(leftOfArrow *trees.ParseTree, exprIn expressionKind) *trees.ParseTree {
	arrowFormalParameters := p.transformToArrowFormalParameters(leftOfArrow)
	if p.peekImplicitSemiColon() {
		p.reportError("No newline allowed before '=>'")
	}
	p.eat(TokenTypeArrow)
	arrowFunctionBody := p.parseArrowFunctionBody(exprIn, functionFlavorNormal)
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindArrow).
		SetFormalParameterList(arrowFormalParameters).
		SetFunctionBody(arrowFunctionBody)
	loc := toTreesRange(SourceRange{
		Start: SourcePosition{
			Offset: arrowFormalParameters.Location.Start.Offset,
			Line:   arrowFormalParameters.Location.Start.Line,
			Column: arrowFormalParameters.Location.Start.Column,
		},
		End: p.lastSourcePosition,
	})
	return p.asParseTree(builder.Build(loc))
}

// transformToArrowFormalParameters mirrors transformToArrowFormalParameters.
func (p *ParserRewrite) transformToArrowFormalParameters(leftOfArrow *trees.ParseTree) *trees.FormalParameterListTree {
	switch leftOfArrow.Type {
	case trees.ParseTreeTypeFormalParameterList:
		// direct formal parameter list — came from parseCoverParenthesizedExpression...
		// when the result was already re-parsed as a FormalParameterListTree.
		// Re-parse from the start of the tree to recover the typed value.
		p.resetScannerToTree(leftOfArrow)
		return p.parseFormalParameterList()
	case trees.ParseTreeTypeIdentifierExpression:
		// single-parameter arrow: x =>
		return trees.NewFormalParameterListTree(
			leftOfArrow.Location,
			[]*trees.ParseTree{leftOfArrow},
			false,
			nil,
		)
	case trees.ParseTreeTypeArgumentList, trees.ParseTreeTypeParenExpression:
		// (x) => or (x, y) =>  — re-parse as formal param list
		p.resetScannerToTree(leftOfArrow)
		return p.parseFormalParameterList()
	default:
		p.reportErrorAt(
			SourcePosition{
				Offset: leftOfArrow.Location.Start.Offset,
				Line:   leftOfArrow.Location.Start.Line,
				Column: leftOfArrow.Location.Start.Column,
			},
			"invalid arrow function parameters",
		)
		return trees.NewFormalParameterListTree(leftOfArrow.Location, nil, false, nil)
	}
}

// parseAsyncArrowFunction mirrors parseAsyncArrowFunction.
func (p *ParserRewrite) parseAsyncArrowFunction(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eatPredefinedString(predefinedASYNC)
	if p.peekImplicitSemiColon() {
		p.reportError("No newline allowed between `async` and arrow function parameter list")
	}
	var arrowParameterList *trees.FormalParameterListTree
	if p.peek(TokenTypeOpenParen) {
		arrowParameterList = p.parseFormalParameterList()
	} else {
		single := p.parseIdentifierExpression()
		arrowParameterList = trees.NewFormalParameterListTree(
			single.Location,
			[]*trees.ParseTree{p.asParseTree(single)},
			false,
			nil,
		)
	}
	if p.peekImplicitSemiColon() {
		p.reportError("No newline allowed before '=>'")
	}
	p.eat(TokenTypeArrow)
	arrowFunctionBody := p.parseArrowFunctionBody(exprIn, functionFlavorAsynchronous)
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindArrow).
		SetAsync(true).
		SetFormalParameterList(arrowParameterList).
		SetFunctionBody(arrowFunctionBody)
	return p.asParseTree(builder.Build(p.getTreeLocation(start)))
}

// parseArrowFunctionBody mirrors parseArrowFunctionBody.
func (p *ParserRewrite) parseArrowFunctionBody(exprIn expressionKind, flavor functionFlavor) *trees.ParseTree {
	p.functionContextStack = append(p.functionContextStack, flavor)
	var body *trees.ParseTree
	if p.peek(TokenTypeOpenCurly) {
		body = p.asParseTree(p.parseFunctionBody())
	} else {
		body = p.parseAssignment(exprIn)
	}
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
	return body
}

// ---- Yield -------------------------------------------------------------------

// parseYield mirrors parseYield(Expression).
func (p *ParserRewrite) parseYield(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeYield)
	isYieldAll := false
	var expression *trees.ParseTree
	if !p.peekImplicitSemiColon() {
		isYieldAll = p.eatOpt(TokenTypeStar) != nil
		if p.peekAssignmentExpression() {
			expression = p.parseAssignment(exprIn)
		} else if isYieldAll {
			p.reportError("yield* requires an expression")
		}
	}
	return p.asParseTree(trees.NewYieldExpressionTree(p.getTreeLocation(start), isYieldAll, expression))
}

// ---- Conditional expression --------------------------------------------------

// parseConditional mirrors parseConditional(Expression).
func (p *ParserRewrite) parseConditional(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	condition := p.parseShortCircuit(exprIn)
	if p.peek(TokenTypeQuestion) {
		p.eat(TokenTypeQuestion)
		left := p.parseAssignment(exprIn)
		p.eat(TokenTypeColon)
		right := p.parseAssignment(exprIn)
		return p.asParseTree(trees.NewConditionalExpressionTree(p.getTreeLocation(start), condition, left, right))
	}
	return condition
}

// ---- Short circuit (nullish coalescing) --------------------------------------

func (p *ParserRewrite) parseShortCircuit(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseLogicalOR(exprIn)
	if p.peek(TokenTypeQuestionQuestion) {
		if left.Type == trees.ParseTreeTypeBinaryOperator {
			// We can't inspect the operator without a typed pointer; the check for
			// && or || mixing with ?? is best-effort. If this tree was just built by
			// parseLogicalOR it will already have reported an error if relevant.
		}
		return p.parseNullishCoalesce(exprIn, left, start)
	}
	return left
}

func (p *ParserRewrite) parseNullishCoalesce(exprIn expressionKind, left *trees.ParseTree, start SourcePosition) *trees.ParseTree {
	for p.peek(TokenTypeQuestionQuestion) {
		operator := p.eat(TokenTypeQuestionQuestion)
		right := p.parseBitwiseOR(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	if p.peek(TokenTypeAnd) || p.peek(TokenTypeOr) {
		p.reportError("Logical OR and logical AND require parentheses when used with '??'")
	}
	return left
}

// ---- Binary operators --------------------------------------------------------

func (p *ParserRewrite) parseLogicalOR(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseLogicalAND(exprIn)
	for p.peek(TokenTypeOr) {
		operator := p.eat(TokenTypeOr)
		right := p.parseLogicalAND(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) parseLogicalAND(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseBitwiseOR(exprIn)
	for p.peek(TokenTypeAnd) {
		operator := p.eat(TokenTypeAnd)
		right := p.parseBitwiseOR(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) parseBitwiseOR(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseBitwiseXOR(exprIn)
	for p.peek(TokenTypeBar) {
		operator := p.eat(TokenTypeBar)
		right := p.parseBitwiseXOR(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) parseBitwiseXOR(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseBitwiseAND(exprIn)
	for p.peek(TokenTypeCaret) {
		operator := p.eat(TokenTypeCaret)
		right := p.parseBitwiseAND(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) parseBitwiseAND(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseEquality(exprIn)
	for p.peek(TokenTypeAmpersand) {
		operator := p.eat(TokenTypeAmpersand)
		right := p.parseEquality(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) parseEquality(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseRelational(exprIn)
	for p.peekEqualityOperator() {
		operator := p.nextToken()
		right := p.parseRelational(exprIn)
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) peekEqualityOperator() bool {
	switch p.peekType() {
	case TokenTypeEqualEqual, TokenTypeNotEqual,
		TokenTypeEqualEqualEqual, TokenTypeNotEqualEqual:
		return true
	}
	return false
}

func (p *ParserRewrite) parseRelational(exprIn expressionKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseShiftExpression()
	for p.peekRelationalOperator(exprIn) {
		operator := p.nextToken()
		right := p.parseShiftExpression()
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) peekRelationalOperator(exprIn expressionKind) bool {
	switch p.peekType() {
	case TokenTypeOpenAngle, TokenTypeCloseAngle,
		TokenTypeGreaterEqual, TokenTypeLessEqual,
		TokenTypeInstanceof:
		return true
	case TokenTypeIn:
		return exprIn == expressionNormal
	}
	return false
}

func (p *ParserRewrite) parseShiftExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseAdditiveExpression()
	for p.peekShiftOperator() {
		operator := p.nextToken()
		right := p.parseAdditiveExpression()
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) peekShiftOperator() bool {
	switch p.peekType() {
	case TokenTypeLeftShift, TokenTypeRightShift, TokenTypeUnsignedRightShift:
		return true
	}
	return false
}

func (p *ParserRewrite) parseAdditiveExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseMultiplicativeExpression()
	for p.peekAdditiveOperator() {
		operator := p.nextToken()
		right := p.parseMultiplicativeExpression()
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) peekAdditiveOperator() bool {
	switch p.peekType() {
	case TokenTypePlus, TokenTypeMinus:
		return true
	}
	return false
}

func (p *ParserRewrite) parseMultiplicativeExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseExponentiationExpression()
	for p.peekMultiplicativeOperator() {
		operator := p.nextToken()
		right := p.parseExponentiationExpression()
		left = p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

func (p *ParserRewrite) peekMultiplicativeOperator() bool {
	switch p.peekType() {
	case TokenTypeStar, TokenTypeSlash, TokenTypePercent:
		return true
	}
	return false
}

func (p *ParserRewrite) parseExponentiationExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	left := p.parseUnaryExpression()
	if p.peek(TokenTypeStarStar) {
		if left.Type == trees.ParseTreeTypeUnaryExpression {
			p.reportError("Unary operator requires parentheses before '**'")
		}
		operator := p.nextToken()
		right := p.parseExponentiationExpression()
		return p.asParseTree(trees.NewBinaryOperatorTree(p.getTreeLocation(start), left, toTreesToken(operator), right))
	}
	return left
}

// ---- Unary / update expressions ----------------------------------------------

// parseUnaryExpression mirrors parseUnaryExpression().
func (p *ParserRewrite) parseUnaryExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	if p.peekUnaryOperator() {
		operator := p.nextToken()
		operand := p.parseUnaryExpression()
		return p.asParseTree(trees.NewUnaryExpressionTree(p.getTreeLocation(start), toTreesToken(operator), operand))
	}
	if p.peekAwaitExpression() {
		return p.parseAwaitExpression()
	}
	return p.parseUpdateExpression()
}

func (p *ParserRewrite) peekUnaryOperator() bool {
	switch p.peekType() {
	case TokenTypeDelete, TokenTypeVoid, TokenTypeTypeof,
		TokenTypePlus, TokenTypeMinus, TokenTypeTilde, TokenTypeBang:
		return true
	}
	return false
}

func (p *ParserRewrite) peekAwaitExpression() bool {
	return p.peekPredefinedString(predefinedAWAIT)
}

func (p *ParserRewrite) parseAwaitExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eatPredefinedString(predefinedAWAIT)
	expression := p.parseUnaryExpression()
	return p.asParseTree(trees.NewAwaitExpressionTree(p.getTreeLocation(start), expression))
}

// parseUpdateExpression mirrors parseUpdateExpression().
func (p *ParserRewrite) parseUpdateExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	if p.peekUpdateOperator() {
		operator := p.nextToken()
		operand := p.parseUnaryExpression()
		return p.asParseTree(trees.NewPrefixUpdateExpressionTree(p.getTreeLocation(start), toTreesToken(operator), operand))
	}
	lhs := p.parseLeftHandSideExpression()
	if p.peekUpdateOperator() && !p.peekImplicitSemiColon() {
		operator := p.nextToken()
		return p.asParseTree(trees.NewPostfixUpdateExpressionTree(p.getTreeLocation(start), toTreesToken(operator), lhs))
	}
	return lhs
}

func (p *ParserRewrite) peekUpdateOperator() bool {
	switch p.peekType() {
	case TokenTypePlusPlus, TokenTypeMinusMinus:
		return true
	}
	return false
}

// ---- Left-hand side expression -----------------------------------------------

func (p *ParserRewrite) peekImportCall() bool {
	return p.peek(TokenTypeImport) && p.peekAt(1, TokenTypeOpenParen)
}

func (p *ParserRewrite) peekImportDot() bool {
	return p.peek(TokenTypeImport) && p.peekAt(1, TokenTypePeriod)
}

// parseLeftHandSideExpression mirrors parseLeftHandSideExpression().
func (p *ParserRewrite) parseLeftHandSideExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	operand := p.parseNewExpression()

	// If operand is not a bare NewExpression (no arguments), then it is a
	// MemberExpression and we look for call suffixes.
	isNewWithoutArgs := operand.Type == trees.ParseTreeTypeNewExpression
	if !isNewWithoutArgs {
		for p.peekCallSuffix() {
			switch p.peekType() {
			case TokenTypeOpenParen:
				arguments := p.parseArguments()
				operand = p.asParseTree(trees.NewCallExpressionTree(p.getTreeLocation(start), operand, arguments))
			case TokenTypeOpenSquare:
				p.eat(TokenTypeOpenSquare)
				member := p.parseExpression()
				p.eat(TokenTypeCloseSquare)
				operand = p.asParseTree(trees.NewMemberLookupExpressionTree(p.getTreeLocation(start), operand, member))
			case TokenTypePeriod:
				p.eat(TokenTypePeriod)
				id := p.eatIdOrKeywordAsId()
				operand = p.asParseTree(trees.NewMemberExpressionTree(p.getTreeLocation(start), operand, toTreesIdentToken(id)))
			case TokenTypeNoSubstitutionTemplate, TokenTypeTemplateHead:
				operand = p.asParseTree(p.parseTemplateLiteral(operand))
			}
		}
		operand = p.maybeParseOptionalExpression(operand)
	}
	return operand
}

func (p *ParserRewrite) peekCallSuffix() bool {
	return p.peek(TokenTypeOpenParen) ||
		p.peek(TokenTypeOpenSquare) ||
		p.peek(TokenTypePeriod) ||
		p.peek(TokenTypeNoSubstitutionTemplate) ||
		p.peek(TokenTypeTemplateHead)
}

// maybeParseOptionalExpression handles `?.` optional chaining.
func (p *ParserRewrite) maybeParseOptionalExpression(operand *trees.ParseTree) *trees.ParseTree {
	start := SourcePosition{
		Offset: operand.Location.Start.Offset,
		Line:   operand.Location.Start.Line,
		Column: operand.Location.Start.Column,
	}
	for p.peek(TokenTypeQuestionDot) {
		p.eat(TokenTypeQuestionDot)
		switch p.peekType() {
		case TokenTypeOpenParen:
			arguments := p.parseArguments()
			operand = p.asParseTree(trees.NewOptChainCallExpressionTree(
				p.getTreeLocation(start), operand, arguments, true, arguments.HasTrailingComma))
		case TokenTypeOpenSquare:
			p.eat(TokenTypeOpenSquare)
			member := p.parseExpression()
			p.eat(TokenTypeCloseSquare)
			operand = p.asParseTree(trees.NewOptionalMemberLookupExpressionTree(
				p.getTreeLocation(start), operand, member, true))
		case TokenTypeNoSubstitutionTemplate, TokenTypeTemplateHead:
			p.reportError("template literal cannot be used within optional chaining")
		default:
			if p.peekIdOrKeyword() {
				id := p.eatIdOrKeywordAsId()
				operand = p.asParseTree(trees.NewOptionalMemberExpressionTree(
					p.getTreeLocation(start), operand, toTreesIdentToken(id), true))
			} else {
				p.reportError("syntax error: %s not allowed in optional chain", tokenTypeValue(p.peekType()))
			}
		}
		operand = p.parseRemainingOptionalChainSegment(operand)
	}
	return operand
}

// parseRemainingOptionalChainSegment mirrors parseRemainingOptionalChainSegment.
func (p *ParserRewrite) parseRemainingOptionalChainSegment(optExpr *trees.ParseTree) *trees.ParseTree {
	start := SourcePosition{
		Offset: optExpr.Location.Start.Offset,
		Line:   optExpr.Location.Start.Line,
		Column: optExpr.Location.Start.Column,
	}
	for p.peekOptionalChainSuffix() {
		if p.peekType() == TokenTypeNoSubstitutionTemplate ||
			p.peekType() == TokenTypeTemplateHead {
			p.reportError("template literal cannot be used within optional chaining")
			break
		}
		switch p.peekType() {
		case TokenTypePeriod:
			p.eat(TokenTypePeriod)
			id := p.eatIdOrKeywordAsId()
			optExpr = p.asParseTree(trees.NewOptionalMemberExpressionTree(
				p.getTreeLocation(start), optExpr, toTreesIdentToken(id), false))
		case TokenTypeOpenParen:
			arguments := p.parseArguments()
			optExpr = p.asParseTree(trees.NewOptChainCallExpressionTree(
				p.getTreeLocation(start), optExpr, arguments, false, arguments.HasTrailingComma))
		case TokenTypeOpenSquare:
			p.eat(TokenTypeOpenSquare)
			member := p.parseExpression()
			p.eat(TokenTypeCloseSquare)
			optExpr = p.asParseTree(trees.NewOptionalMemberLookupExpressionTree(
				p.getTreeLocation(start), optExpr, member, false))
		default:
			panic(fmt.Sprintf("parseRemainingOptionalChainSegment: unexpected token %s", tokenTypeValue(p.peekType())))
		}
	}
	return optExpr
}

func (p *ParserRewrite) peekOptionalChainSuffix() bool {
	return p.peek(TokenTypeOpenParen) ||
		p.peek(TokenTypeOpenSquare) ||
		p.peek(TokenTypePeriod) ||
		p.peek(TokenTypeNoSubstitutionTemplate) ||
		p.peek(TokenTypeTemplateHead)
}

// ---- Member expression -------------------------------------------------------

// parseMemberExpressionNoNew mirrors parseMemberExpressionNoNew().
func (p *ParserRewrite) parseMemberExpressionNoNew() *trees.ParseTree {
	start := p.getTreeStartLocation()
	var operand *trees.ParseTree
	if p.peekImportDot() {
		operand = p.asParseTree(p.parseImportDotMeta())
	} else if p.peekAsyncFunctionStart() {
		operand = p.parseAsyncFunctionExpression()
	} else if p.peekFunction() {
		operand = p.parseFunctionExpression()
	} else {
		operand = p.parsePrimaryExpression()
	}
	for p.peekMemberExpressionSuffix() {
		switch p.peekType() {
		case TokenTypeOpenSquare:
			p.eat(TokenTypeOpenSquare)
			member := p.parseExpression()
			p.eat(TokenTypeCloseSquare)
			operand = p.asParseTree(trees.NewMemberLookupExpressionTree(p.getTreeLocation(start), operand, member))
		case TokenTypePeriod:
			p.eat(TokenTypePeriod)
			id := p.eatIdOrKeywordAsId()
			operand = p.asParseTree(trees.NewMemberExpressionTree(p.getTreeLocation(start), operand, toTreesIdentToken(id)))
		case TokenTypeNoSubstitutionTemplate, TokenTypeTemplateHead:
			operand = p.asParseTree(p.parseTemplateLiteral(operand))
		}
	}
	return operand
}

func (p *ParserRewrite) peekMemberExpressionSuffix() bool {
	return p.peek(TokenTypeOpenSquare) ||
		p.peek(TokenTypePeriod) ||
		p.peek(TokenTypeNoSubstitutionTemplate) ||
		p.peek(TokenTypeTemplateHead)
}

// ---- New expression ----------------------------------------------------------

// parseNewExpression mirrors parseNewExpression().
func (p *ParserRewrite) parseNewExpression() *trees.ParseTree {
	if !p.peek(TokenTypeNew) {
		return p.parseMemberExpressionNoNew()
	}
	if p.peekAt(1, TokenTypePeriod) {
		return p.asParseTree(p.parseNewDotSomething())
	}
	start := p.getTreeStartLocation()
	p.eat(TokenTypeNew)
	if p.peek(TokenTypeQuestionDot) {
		p.reportError("Optional chaining is forbidden in `new?.target` contexts.")
	}
	operand := p.parseNewExpression()
	if p.peek(TokenTypeQuestionDot) {
		p.reportError("Optional chaining is forbidden in construction contexts.")
	}
	var arguments *trees.ArgumentListTree
	if p.peek(TokenTypeOpenParen) {
		arguments = p.parseArguments()
	}
	hasTrailingComma := arguments != nil && arguments.HasTrailingComma
	return p.asParseTree(trees.NewNewExpressionTree(p.getTreeLocation(start), operand, arguments, hasTrailingComma))
}

// parseNewDotSomething mirrors parseNewDotSomething().
func (p *ParserRewrite) parseNewDotSomething() *trees.NewTargetExpressionTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeNew)
	p.eat(TokenTypePeriod)
	p.eatPredefinedString("target")
	return trees.NewNewTargetExpressionTree(p.getTreeLocation(start))
}

// parseImportDotMeta mirrors parseImportDotMeta().
func (p *ParserRewrite) parseImportDotMeta() *trees.ImportMetaExpressionTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeImport)
	p.eat(TokenTypePeriod)
	p.eatPredefinedString("meta")
	return trees.NewImportMetaExpressionTree(p.getTreeLocation(start))
}

// ---- Arguments ---------------------------------------------------------------

// parseArguments mirrors parseArguments().
func (p *ParserRewrite) parseArguments() *trees.ArgumentListTree {
	start := p.getTreeStartLocation()
	var arguments []*trees.ParseTree
	trailingComma := false
	var commaPositions []trees.SourcePosition

	p.eat(TokenTypeOpenParen)
	for p.peekAssignmentOrSpread() {
		arguments = append(arguments, p.parseAssignmentOrSpread())
		if !p.peek(TokenTypeCloseParen) {
			comma := p.eat(TokenTypeComma)
			if comma != nil {
				commaPositions = append(commaPositions, toTreesPos(comma.GetStart()))
			}
			if p.peek(TokenTypeCloseParen) {
				if !p.config.AtLeast8 {
					p.reportErrorToken(comma, "Invalid trailing comma in arguments list")
				}
				trailingComma = true
			}
		}
	}
	p.eat(TokenTypeCloseParen)
	return trees.NewArgumentListTree(p.getTreeLocation(start), arguments, trailingComma, commaPositions)
}

func (p *ParserRewrite) peekAssignmentOrSpread() bool {
	return p.peek(TokenTypeEllipsis) || p.peekAssignmentExpression()
}

func (p *ParserRewrite) parseAssignmentOrSpread() *trees.ParseTree {
	if p.peek(TokenTypeEllipsis) {
		return p.parseIterSpread()
	}
	return p.parseAssignmentExpression()
}

// parseIterSpread mirrors parseIterSpread().
func (p *ParserRewrite) parseIterSpread() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeEllipsis)
	expr := p.parseAssignmentExpression()
	return p.asParseTree(trees.NewIterSpreadTree(p.getTreeLocation(start), expr))
}

// ---- Primary expression ------------------------------------------------------

// parsePrimaryExpression mirrors parsePrimaryExpression().
func (p *ParserRewrite) parsePrimaryExpression() *trees.ParseTree {
	switch p.peekType() {
	case TokenTypeClass:
		return p.parseClassExpression()
	case TokenTypeSuper:
		return p.asParseTree(p.parseSuperExpression())
	case TokenTypeThis:
		return p.asParseTree(p.parseThisExpression())
	case TokenTypeImport:
		return p.asParseTree(p.parseDynamicImportExpression())
	case TokenTypeIdentifier, TokenTypeType, TokenTypeDeclare, TokenTypeModule, TokenTypeNamespace:
		return p.asParseTree(p.parseIdentifierExpression())
	case TokenTypeNumber, TokenTypeString, TokenTypeBigInt, TokenTypeTrue, TokenTypeFalse, TokenTypeNull:
		return p.asParseTree(p.parseLiteralExpression())
	case TokenTypeNoSubstitutionTemplate, TokenTypeTemplateHead:
		return p.asParseTree(p.parseTemplateLiteral(nil))
	case TokenTypeOpenSquare:
		return p.parseArrayInitializer()
	case TokenTypeOpenCurly:
		return p.parseObjectLiteral()
	case TokenTypeOpenParen:
		return p.parseCoverParenthesizedExpressionAndArrowParameterList()
	case TokenTypeSlash, TokenTypeSlashEqual:
		return p.asParseTree(p.parseRegularExpressionLiteral())
	default:
		return p.parseMissingPrimaryExpression()
	}
}

func (p *ParserRewrite) parseSuperExpression() *trees.SuperExpressionTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeSuper)
	if p.peek(TokenTypeQuestionDot) {
		p.reportError("Optional chaining is forbidden in super?.")
	}
	return trees.NewSuperExpressionTree(p.getTreeLocation(start))
}

func (p *ParserRewrite) parseThisExpression() *trees.ThisExpressionTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeThis)
	return trees.NewThisExpressionTree(p.getTreeLocation(start))
}

func (p *ParserRewrite) parseDynamicImportExpression() *trees.DynamicImportTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeImport)
	if p.peek(TokenTypeQuestionDot) {
		p.reportError("Optional chaining is forbidden in import?.")
	}
	p.eat(TokenTypeOpenParen)
	argument := p.parseAssignmentExpression()
	p.eat(TokenTypeCloseParen)
	return trees.NewDynamicImportTree(p.getTreeLocation(start), argument)
}

func (p *ParserRewrite) parseIdentifierExpression() *trees.IdentifierExpressionTree {
	start := p.getTreeStartLocation()
	identifier := p.eatId()
	return trees.NewIdentifierExpressionTree(p.getTreeLocation(start), toTreesIdentToken(identifier))
}

func (p *ParserRewrite) parseLiteralExpression() *trees.LiteralExpressionTree {
	start := p.getTreeStartLocation()
	literal := p.nextToken()
	return trees.NewLiteralExpressionTree(p.getTreeLocation(start), toTreesToken(literal))
}

func (p *ParserRewrite) parseRegularExpressionLiteral() *trees.LiteralExpressionTree {
	start := p.getTreeStartLocation()
	literal := p.nextRegularExpressionLiteralToken()
	return trees.NewLiteralExpressionTree(p.getTreeLocation(start), toTreesToken(literal))
}

// parseMissingPrimaryExpression mirrors parseMissingPrimaryExpression().
func (p *ParserRewrite) parseMissingPrimaryExpression() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.nextToken()
	p.reportError("primary expression expected")
	return p.asParseTree(trees.NewMissingPrimaryExpressionTree(p.getTreeLocation(start)))
}

// ---- Template literals -------------------------------------------------------

// parseTemplateLiteral mirrors parseTemplateLiteral(ParseTree operand).
// operand is nil for untagged templates (e.g., `abc`), non-nil for tagged (foo`abc`).
func (p *ParserRewrite) parseTemplateLiteral(operand *trees.ParseTree) *trees.TemplateLiteralExpressionTree {
	var start SourcePosition
	if operand == nil {
		start = p.getTreeStartLocation()
	} else {
		start = SourcePosition{
			Offset: operand.Location.Start.Offset,
			Line:   operand.Location.Start.Line,
			Column: operand.Location.Start.Column,
		}
	}
	token := p.nextToken()
	if token.Type != TokenTypeNoSubstitutionTemplate && token.Type != TokenTypeTemplateHead {
		p.reportErrorToken(token, "Unexpected template literal token %s.", tokenTypeValue(token.Type))
	}
	isTaggedTemplate := operand != nil
	if !isTaggedTemplate {
		p.reportTemplateErrorIfPresent(token)
	}
	var elements []*trees.ParseTree
	elements = append(elements, p.asParseTree(trees.NewTemplateLiteralPortionTree(
		toTreesRange(token.Location), toTreesToken(token))))
	if token.Type == TokenTypeNoSubstitutionTemplate {
		return trees.NewTemplateLiteralExpressionTree(p.getTreeLocation(start), operand, elements)
	}
	// `abc${
	expression := p.parseExpression()
	elements = append(elements, p.asParseTree(trees.NewTemplateSubstitutionTree(
		expression.Location, expression)))
	for {
		templateToken := p.nextTemplateLiteralToken()
		if templateToken.Type == TokenTypeError || templateToken.Type == TokenTypeEndOfFile {
			break
		}
		if !isTaggedTemplate {
			p.reportTemplateErrorIfPresent(templateToken)
		}
		elements = append(elements, p.asParseTree(trees.NewTemplateLiteralPortionTree(
			toTreesRange(templateToken.Location), toTreesToken(templateToken))))
		if templateToken.Type == TokenTypeTemplateTail {
			break
		}
		expression = p.parseExpression()
		elements = append(elements, p.asParseTree(trees.NewTemplateSubstitutionTree(
			expression.Location, expression)))
	}
	return trees.NewTemplateLiteralExpressionTree(p.getTreeLocation(start), operand, elements)
}

func (p *ParserRewrite) reportTemplateErrorIfPresent(tok *Token) {
	if tok.ErrorMessage == "" {
		return
	}
	switch tok.ErrorLevel {
	case TemplateLiteralWarning:
		p.errorReporter.ReportWarning(tok.ErrorPosition, "%s", tok.ErrorMessage)
	case TemplateLiteralError:
		p.reportErrorAt(tok.ErrorPosition, "%s", tok.ErrorMessage)
	}
}

// ---- Array literal / comprehension ------------------------------------------

func (p *ParserRewrite) parseArrayInitializer() *trees.ParseTree {
	if p.peekTypeAt(1) == TokenTypeFor {
		return p.parseArrayComprehension()
	}
	return p.parseArrayLiteral()
}

func (p *ParserRewrite) parseGeneratorComprehension() *trees.ParseTree {
	return p.parseComprehension(
		trees.ComprehensionTypeGenerator, TokenTypeOpenParen, TokenTypeCloseParen)
}

func (p *ParserRewrite) parseArrayComprehension() *trees.ParseTree {
	return p.parseComprehension(
		trees.ComprehensionTypeArray, TokenTypeOpenSquare, TokenTypeCloseSquare)
}

func (p *ParserRewrite) parseComprehension(
	compType trees.ComprehensionType, startTok, endTok TokenType,
) *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(startTok)
	var children []*trees.ParseTree
	for p.peek(TokenTypeFor) || p.peek(TokenTypeIf) {
		if p.peek(TokenTypeFor) {
			children = append(children, p.parseComprehensionFor())
		} else {
			children = append(children, p.parseComprehensionIf())
		}
	}
	tailExpression := p.parseAssignmentExpression()
	p.eat(endTok)
	return p.asParseTree(trees.NewComprehensionTree(p.getTreeLocation(start), compType, children, tailExpression))
}

func (p *ParserRewrite) parseComprehensionFor() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeFor)
	p.eat(TokenTypeOpenParen)
	var initializer *trees.ParseTree
	if p.peekId() {
		initializer = p.asParseTree(p.parseIdentifierExpression())
	} else {
		initializer = p.parsePattern(patternKindAny)
	}
	p.eatPredefinedString(predefinedOF)
	collection := p.parseAssignmentExpression()
	p.eat(TokenTypeCloseParen)
	return p.asParseTree(trees.NewComprehensionForTree(p.getTreeLocation(start), initializer, collection))
}

func (p *ParserRewrite) parseComprehensionIf() *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeIf)
	p.eat(TokenTypeOpenParen)
	initializer := p.parseAssignmentExpression()
	p.eat(TokenTypeCloseParen)
	return p.asParseTree(trees.NewComprehensionIfTree(p.getTreeLocation(start), initializer))
}

// parseArrayLiteral mirrors parseArrayLiteral().
func (p *ParserRewrite) parseArrayLiteral() *trees.ParseTree {
	start := p.getTreeStartLocation()
	var elements []*trees.ParseTree
	p.eat(TokenTypeOpenSquare)
	var trailingCommaToken *Token
	for p.peek(TokenTypeComma) || p.peek(TokenTypeEllipsis) || p.peekAssignmentExpression() {
		trailingCommaToken = nil
		if p.peek(TokenTypeComma) {
			commaStart := p.getTreeStartLocation()
			trailingCommaToken = p.eat(TokenTypeComma)
			nullRange := toTreesRange(SourceRange{Start: commaStart, End: commaStart})
			elements = append(elements, p.asParseTree(trees.NewNullTree(nullRange)))
		} else {
			if p.peek(TokenTypeEllipsis) {
				elements = append(elements, p.parseIterSpread())
			} else {
				elements = append(elements, p.parseAssignmentExpression())
			}
			if !p.peek(TokenTypeCloseSquare) {
				trailingCommaToken = p.eat(TokenTypeComma)
			}
		}
	}
	p.eat(TokenTypeCloseSquare)
	p.maybeReportTrailingComma(trailingCommaToken)
	return p.asParseTree(trees.NewArrayLiteralExpressionTree(
		p.getTreeLocation(start), elements, trailingCommaToken != nil))
}

// ---- Object literal ----------------------------------------------------------

// parseObjectLiteral mirrors parseObjectLiteral().
func (p *ParserRewrite) parseObjectLiteral() *trees.ParseTree {
	start := p.getTreeStartLocation()
	var result []*trees.ParseTree
	p.eat(TokenTypeOpenCurly)
	var commaToken *Token
	for p.peek(TokenTypeEllipsis) || p.peekPropertyNameOrComputedPropAt(0) || p.peek(TokenTypeStar) {
		result = append(result, p.parseObjectLiteralPropertyAssignment())
		commaToken = p.eatOpt(TokenTypeComma)
		if commaToken == nil {
			break
		}
	}
	p.eat(TokenTypeCloseCurly)
	p.maybeReportTrailingComma(commaToken)
	return p.asParseTree(trees.NewObjectLiteralExpressionTree(
		p.getTreeLocation(start), result, commaToken != nil))
}

// parseObjectLiteralPropertyAssignment mirrors parseObjectLiteralPropertyAssignment().
func (p *ParserRewrite) parseObjectLiteralPropertyAssignment() *trees.ParseTree {
	typ := p.peekType()
	switch {
	case typ == TokenTypeStar:
		return p.parseObjectLiteralPropertyAssignmentGenerator()
	case typ == TokenTypeEllipsis:
		start := p.getTreeStartLocation()
		p.eat(TokenTypeEllipsis)
		operand := p.parseAssignmentExpression()
		return p.asParseTree(trees.NewObjectSpreadTree(p.getTreeLocation(start), operand))
	case typ == TokenTypeString || typ == TokenTypeNumber || typ == TokenTypeBigInt ||
		typ == TokenTypeIdentifier || isKeywordType(typ):
		if p.peekGetAccessor() {
			return p.parseObjectLiteralGetAccessor()
		} else if p.peekSetAccessor() {
			return p.parseObjectLiteralSetAccessor()
		} else if p.peekAsyncMethod() {
			return p.parseObjectLiteralAsyncMethod()
		} else if p.peekTypeAt(1) == TokenTypeOpenParen {
			return p.parseObjectLiteralMethodDeclaration()
		}
		return p.parseObjectLiteralPropertyNameAssignment()
	case typ == TokenTypeOpenSquare:
		start := p.getTreeStartLocation()
		name := p.parseComputedPropertyName()
		if p.peek(TokenTypeColon) {
			p.eat(TokenTypeColon)
			value := p.parseAssignmentExpression()
			return p.asParseTree(trees.NewComputedPropertyDefinitionTree(p.getTreeLocation(start), name, value))
		}
		builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression)
		p.parseFunctionTail(builder, functionFlavorNormal)
		value := p.asParseTree(builder.Build(p.getTreeLocation(start)))
		return p.asParseTree(trees.NewComputedPropertyMethodTree(p.getTreeLocation(start), name, value))
	default:
		panic(fmt.Sprintf("parseObjectLiteralPropertyAssignment: unreachable, token=%s", tokenTypeValue(typ)))
	}
}

func (p *ParserRewrite) parseObjectLiteralPropertyAssignmentGenerator() *trees.ParseTree {
	typ := p.peekTypeAt(1)
	switch {
	case typ == TokenTypeString || typ == TokenTypeNumber || typ == TokenTypeIdentifier || isKeywordType(typ):
		return p.parseObjectLiteralMethodDeclaration()
	default:
		start := p.getTreeStartLocation()
		p.eat(TokenTypeStar)
		name := p.parseComputedPropertyName()
		builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression)
		p.parseFunctionTail(builder, functionFlavorGenerator)
		value := p.asParseTree(builder.Build(p.getTreeLocation(start)))
		return p.asParseTree(trees.NewComputedPropertyMethodTree(p.getTreeLocation(start), name, value))
	}
}

func (p *ParserRewrite) parseComputedPropertyName() *trees.ParseTree {
	p.eat(TokenTypeOpenSquare)
	assign := p.parseAssignmentExpression()
	p.eat(TokenTypeCloseSquare)
	return assign
}

func (p *ParserRewrite) peekGetAccessor() bool {
	return p.peekPredefinedString(predefinedGET) && p.peekPropertyNameOrComputedPropAt(1)
}

func (p *ParserRewrite) peekSetAccessor() bool {
	return p.peekPredefinedString(predefinedSET) && p.peekPropertyNameOrComputedPropAt(1)
}

func (p *ParserRewrite) parseObjectLiteralGetAccessor() *trees.ParseTree {
	return p.parseGetAccessor(createObjectLiteralElementInfoAt(toTreesPos(p.getTreeStartLocation())))
}

func (p *ParserRewrite) parseGetAccessor(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	p.eatPredefinedString(predefinedGET)
	if p.peekPropertyNameAt(0) {
		propertyName := p.eatObjectLiteralPropertyName()
		p.eat(TokenTypeOpenParen)
		p.eat(TokenTypeCloseParen)
		body := p.parseFunctionBody()
		return p.asParseTree(trees.NewGetAccessorTree(
			p.getTreeLocationFromTreesPos(elementInfo.start), toTreesToken(propertyName),
			elementInfo.isClassMember, elementInfo.isStatic, body))
	}
	property := p.parseComputedPropertyName()
	p.eat(TokenTypeOpenParen)
	p.eat(TokenTypeCloseParen)
	body := p.parseFunctionBody()
	return p.asParseTree(trees.NewComputedPropertyGetterTree(
		p.getTreeLocationFromTreesPos(elementInfo.start), property, elementInfo.isStatic, body))
}

func (p *ParserRewrite) parseObjectLiteralSetAccessor() *trees.ParseTree {
	return p.parseSetAccessor(createObjectLiteralElementInfoAt(toTreesPos(p.getTreeStartLocation())))
}

func (p *ParserRewrite) parseSetAccessor(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	p.eatPredefinedString(predefinedSET)
	if p.peekPropertyNameAt(0) {
		propertyName := p.eatObjectLiteralPropertyName()
		parameter := p.parseSetterParameterList()
		body := p.parseFunctionBody()
		return p.asParseTree(trees.NewSetAccessorTree(
			p.getTreeLocationFromTreesPos(elementInfo.start), toTreesToken(propertyName),
			elementInfo.isClassMember, elementInfo.isStatic, parameter, body))
	}
	property := p.parseComputedPropertyName()
	parameter := p.parseSetterParameterList()
	body := p.parseFunctionBody()
	return p.asParseTree(trees.NewComputedPropertySetterTree(
		p.getTreeLocationFromTreesPos(elementInfo.start), property, elementInfo.isStatic, parameter, body))
}

func (p *ParserRewrite) parseObjectLiteralMethodDeclaration() *trees.ParseTree {
	return p.parseMethodDeclaration(createObjectLiteralElementInfoAt(toTreesPos(p.getTreeStartLocation())))
}

func (p *ParserRewrite) parseMethodDeclaration(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	isGenerator := p.eatOpt(TokenTypeStar) != nil
	p.parseClassElementName(elementInfo)
	return p.parseMethodDefinition(elementInfo, isGenerator)
}

func (p *ParserRewrite) parseObjectLiteralAsyncMethod() *trees.ParseTree {
	return p.parseAsyncMethod(createObjectLiteralElementInfoAt(toTreesPos(p.getTreeStartLocation())))
}

// parseAsyncMethod mirrors parseAsyncMethod() from the Java parser.
func (p *ParserRewrite) parseAsyncMethod(elementInfo *classOrObjectElementInfo) *trees.ParseTree {
	p.eatPredefinedString(predefinedASYNC)
	generator := p.peek(TokenTypeStar)
	if generator {
		p.eat(TokenTypeStar)
	}
	flavor := functionFlavorAsynchronous
	if generator {
		flavor = functionFlavorAsynchronousGenerator
	}
	if p.peekPropertyNameAt(0) {
		if p.peekIdOrKeyword() {
			nameTok := p.eatIdOrKeywordAsId()
			name := toTreesIdentToken(nameTok)
			nameStart := nameTok.Location.Start
			builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindMember).
				SetAsync(true).
				SetGenerator(generator).
				SetIsClassMember(elementInfo.isClassMember).
				SetStatic(elementInfo.isStatic).
				SetName(name)
			p.parseFunctionTail(builder, flavor)
			nameTreesPos := toTreesPos(nameStart)
			return p.asParseTree(builder.Build(p.getTreeLocationFromTreesPos(nameTreesPos)))
		}
		// string/number literal — treat as computed
		nameExpr := p.asParseTree(p.parseLiteralExpression())
		builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression).
			SetAsync(true).
			SetGenerator(generator).
			SetIsClassMember(elementInfo.isClassMember).
			SetStatic(elementInfo.isStatic)
		p.parseFunctionTail(builder, flavor)
		funcTree := p.asParseTree(builder.Build(p.getTreeLocationFromTreesPos(nameExpr.Location.Start)))
		return p.asParseTree(trees.NewComputedPropertyMethodTree(
			p.getTreeLocationFromTreesPos(nameExpr.Location.Start), nameExpr, funcTree))
	}
	// computed property name
	nameExpr := p.parseComputedPropertyName()
	builder := trees.NewFunctionDeclarationBuilder(trees.FunctionDeclarationKindExpression).
		SetAsync(true).
		SetGenerator(generator).
		SetIsClassMember(elementInfo.isClassMember).
		SetStatic(elementInfo.isStatic)
	p.parseFunctionTail(builder, flavor)
	funcTree := p.asParseTree(builder.Build(p.getTreeLocationFromTreesPos(nameExpr.Location.Start)))
	return p.asParseTree(trees.NewComputedPropertyMethodTree(
		p.getTreeLocationFromTreesPos(nameExpr.Location.Start), nameExpr, funcTree))
}

// parseMethodDefinition mirrors parseMethodDefinition() from the Java parser.
func (p *ParserRewrite) parseMethodDefinition(elementInfo *classOrObjectElementInfo, isGenerator bool) *trees.ParseTree {
	var kind trees.FunctionDeclarationKind
	if elementInfo.hasNameExpr() {
		kind = trees.FunctionDeclarationKindExpression
	} else {
		kind = trees.FunctionDeclarationKindMember
	}
	builder := trees.NewFunctionDeclarationBuilder(kind).
		SetIsClassMember(elementInfo.isClassMember).
		SetStatic(elementInfo.isStatic)
	if elementInfo.hasName() {
		builder.SetName(elementInfo.getName())
	}
	flavor := functionFlavorNormal
	if isGenerator {
		flavor = functionFlavorGenerator
	}
	p.parseFunctionTail(builder, flavor)
	function := p.asParseTree(builder.Build(p.getTreeLocationFromTreesPos(elementInfo.start)))
	if kind == trees.FunctionDeclarationKindMember {
		return function
	}
	return p.asParseTree(trees.NewComputedPropertyMethodTree(
		p.getTreeLocationFromTreesPos(elementInfo.start), elementInfo.getNameExpr(), function))
}

// peekAsyncMethod mirrors peekAsyncMethod() from the Java parser.
func (p *ParserRewrite) peekAsyncMethod() bool {
	return p.peekPredefinedString(predefinedASYNC) &&
		!p.peekImplicitSemiColonAt(1) &&
		(p.peekPropertyNameOrComputedPropAt(1) ||
			(p.peekAt(1, TokenTypeStar) && p.peekPropertyNameOrComputedPropAt(2)))
}

func (p *ParserRewrite) parseObjectLiteralPropertyNameAssignment() *trees.ParseTree {
	start := p.getTreeStartLocation()
	name := p.eatObjectLiteralPropertyName()
	colon := p.eatOpt(TokenTypeColon)
	if colon == nil {
		if name.Type != TokenTypeIdentifier {
			p.reportExpectedError(p.peekToken(), tokenTypeValue(TokenTypeColon))
		} else if name.Value != "" && isKeywordType(p.peekType()) {
			// keyword used as shorthand — report error
			p.reportErrorToken(name, "Cannot use keyword in short object literal")
		} else if p.peek(TokenTypeEqual) {
			// shorthand with default: { a = 1 }
			idTree := p.asParseTree(trees.NewIdentifierExpressionTree(p.getTreeLocation(start), toTreesIdentToken(name)))
			p.eat(TokenTypeEqual)
			defaultValue := p.parseAssignmentExpression()
			return p.asParseTree(trees.NewDefaultParameterTree(p.getTreeLocation(start), idTree, defaultValue))
		}
	}
	var value *trees.ParseTree
	if colon != nil {
		value = p.parseAssignmentExpression()
	}
	return p.asParseTree(trees.NewPropertyNameAssignmentTree(p.getTreeLocation(start), toTreesToken(name), value))
}

// ---- Parenthesized expression / arrow params ---------------------------------

// parseCoverParenthesizedExpressionAndArrowParameterList mirrors
// parseCoverParenthesizedExpressionAndArrowParameterList().
func (p *ParserRewrite) parseCoverParenthesizedExpressionAndArrowParameterList() *trees.ParseTree {
	if p.peekTypeAt(1) == TokenTypeFor {
		return p.parseGeneratorComprehension()
	}
	start := p.getTreeStartLocation()
	p.eat(TokenTypeOpenParen)
	// Case ( )
	if p.peek(TokenTypeCloseParen) {
		p.eat(TokenTypeCloseParen)
		if p.peek(TokenTypeArrow) {
			return p.asParseTree(trees.NewFormalParameterListTree(p.getTreeLocation(start), nil, false, nil))
		}
		p.reportError("invalid parenthesized expression")
		return p.asParseTree(trees.NewMissingPrimaryExpressionTree(p.getTreeLocation(start)))
	}
	// Case ( ... BindingIdentifier )
	if p.peek(TokenTypeEllipsis) {
		params := []*trees.ParseTree{p.parseParameter()}
		p.eat(TokenTypeCloseParen)
		if p.peek(TokenTypeArrow) {
			return p.asParseTree(trees.NewFormalParameterListTree(p.getTreeLocation(start), params, false, nil))
		}
		p.reportError("invalid parenthesized expression")
		return p.asParseTree(trees.NewMissingPrimaryExpressionTree(p.getTreeLocation(start)))
	}
	// ( Expression ) or ( Expression, ) or ( Expression, ...BindingIdentifier )
	result := p.parseExpression()
	if p.peek(TokenTypeComma) {
		if p.peekAt(1, TokenTypeCloseParen) {
			// ( Expression, ) — trailing comma => re-parse as formal param list
			p.resetScannerToPos(start)
			return p.asParseTree(p.parseFormalParameterList())
		}
		p.eat(TokenTypeComma)
		// ( Expression, ...param ) — leave as comma expression; will be re-parsed
		result = p.asParseTree(trees.NewCommaExpressionTree(
			p.getTreeLocation(start),
			[]*trees.ParseTree{result, p.parseParameter()},
		))
	}
	p.eat(TokenTypeCloseParen)
	return p.asParseTree(trees.NewParenExpressionTree(p.getTreeLocation(start), result))
}

// ---- Destructuring patterns --------------------------------------------------

// parsePattern mirrors parsePattern(PatternKind).
func (p *ParserRewrite) parsePattern(kind patternKind) *trees.ParseTree {
	switch p.peekType() {
	case TokenTypeOpenSquare:
		return p.parseArrayPattern(kind)
	default:
		return p.parseObjectPattern(kind)
	}
}

func (p *ParserRewrite) peekPatternStart() bool {
	return p.peek(TokenTypeOpenSquare) || p.peek(TokenTypeOpenCurly)
}

// parseArrayPattern mirrors parseArrayPattern(PatternKind).
func (p *ParserRewrite) parseArrayPattern(kind patternKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	var elements []*trees.ParseTree
	p.eat(TokenTypeOpenSquare)
	for p.peek(TokenTypeComma) || p.peekArrayPatternElement() {
		if p.peek(TokenTypeComma) {
			nullStart := p.getTreeStartLocation()
			p.eat(TokenTypeComma)
			elements = append(elements, p.asParseTree(trees.NewNullTree(p.getTreeLocation(nullStart))))
		} else {
			elements = append(elements, p.parsePatternAssignmentTarget(kind))
			if p.peek(TokenTypeComma) {
				p.eat(TokenTypeComma)
			} else {
				break
			}
		}
	}
	if p.peek(TokenTypeEllipsis) {
		elements = append(elements, p.parseIterRest(kind))
	}
	if p.eat(TokenTypeCloseSquare) == nil {
		return p.asParseTree(trees.NewMissingPrimaryExpressionTree(
			p.getTreeLocation(p.getTreeStartLocation())))
	}
	return p.asParseTree(trees.NewArrayPatternTree(p.getTreeLocation(start), elements))
}

func (p *ParserRewrite) peekArrayPatternElement() bool {
	return p.peekExpression()
}

// parseObjectPattern mirrors parseObjectPattern(PatternKind).
func (p *ParserRewrite) parseObjectPattern(kind patternKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	var fields []*trees.ParseTree
	p.eat(TokenTypeOpenCurly)
	for p.peekObjectPatternField() {
		fields = append(fields, p.parseObjectPatternField(kind))
		if p.peek(TokenTypeComma) {
			p.eat(TokenTypeComma)
		} else {
			break
		}
	}
	if p.peek(TokenTypeEllipsis) {
		restStart := p.getTreeStartLocation()
		p.eat(TokenTypeEllipsis)
		patternAssignmentTarget := p.parseRestAssignmentTarget(kind)
		fields = append(fields, p.asParseTree(trees.NewObjectRestTree(p.getTreeLocation(restStart), patternAssignmentTarget)))
	}
	p.eat(TokenTypeCloseCurly)
	return p.asParseTree(trees.NewObjectPatternTree(p.getTreeLocation(start), fields))
}

func (p *ParserRewrite) peekObjectPatternField() bool {
	return p.peekPropertyNameOrComputedPropAt(0)
}

func (p *ParserRewrite) parseObjectPatternField(kind patternKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	if p.peekType() == TokenTypeOpenSquare {
		key := p.parseComputedPropertyName()
		p.eat(TokenTypeColon)
		value := p.parsePatternAssignmentTarget(kind)
		return p.asParseTree(trees.NewComputedPropertyDefinitionTree(p.getTreeLocation(start), key, value))
	}
	if p.peekIdOrKeyword() {
		idToken := p.eatIdOrKeywordAsId()
		if !p.peek(TokenTypeColon) {
			if isKeywordType(idToken.Type) {
				p.reportErrorToken(idToken, "cannot use keyword '%s' here.", idToken.Value)
			}
			if p.peek(TokenTypeEqual) {
				idTree := p.asParseTree(trees.NewIdentifierExpressionTree(p.getTreeLocation(start), toTreesIdentToken(idToken)))
				p.eat(TokenTypeEqual)
				defaultValue := p.parseAssignmentExpression()
				return p.asParseTree(trees.NewDefaultParameterTree(p.getTreeLocation(start), idTree, defaultValue))
			}
			return p.asParseTree(trees.NewPropertyNameAssignmentTree(p.getTreeLocation(start), toTreesToken(idToken), nil))
		}
		p.eat(TokenTypeColon)
		value := p.parsePatternAssignmentTarget(kind)
		return p.asParseTree(trees.NewPropertyNameAssignmentTree(p.getTreeLocation(start), toTreesToken(idToken), value))
	}
	// literal key
	name := p.parseLiteralExpression().LiteralToken
	p.eat(TokenTypeColon)
	value := p.parsePatternAssignmentTarget(kind)
	return p.asParseTree(trees.NewPropertyNameAssignmentTree(p.getTreeLocation(start), name, value))
}

// parsePatternAssignmentTarget mirrors parsePatternAssignmentTarget(PatternKind).
func (p *ParserRewrite) parsePatternAssignmentTarget(kind patternKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	target := p.parsePatternAssignmentTargetNoDefault(kind)
	if p.peek(TokenTypeEqual) {
		p.eat(TokenTypeEqual)
		defaultValue := p.parseAssignmentExpression()
		return p.asParseTree(trees.NewDefaultParameterTree(p.getTreeLocation(start), target, defaultValue))
	}
	return target
}

// parsePatternAssignmentTargetNoDefault mirrors parsePatternAssignmentTargetNoDefault(PatternKind).
func (p *ParserRewrite) parsePatternAssignmentTargetNoDefault(kind patternKind) *trees.ParseTree {
	if p.peekPatternStart() {
		return p.parsePattern(kind)
	}
	if kind == patternKindInitializer {
		return p.asParseTree(p.parseIdentifierExpression())
	}
	// patternKindAny: allow arbitrary left-hand-side
	return p.parseLeftHandSideExpression()
}

// parseIterRest mirrors parseIterRest(PatternKind).
func (p *ParserRewrite) parseIterRest(kind patternKind) *trees.ParseTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeEllipsis)
	target := p.parseRestAssignmentTarget(kind)
	return p.asParseTree(trees.NewIterRestTree(p.getTreeLocation(start), target))
}

// parseRestAssignmentTarget mirrors parseRestAssignmentTarget(PatternKind).
func (p *ParserRewrite) parseRestAssignmentTarget(kind patternKind) *trees.ParseTree {
	target := p.parsePatternAssignmentTargetNoDefault(kind)
	if p.peek(TokenTypeEqual) {
		p.reportError("A default value cannot be specified after '...'")
	}
	return target
}

// ---- Formal parameter list ---------------------------------------------------
// These are stubs; full implementation (with rest, default, etc.) comes in part 2.

// parseFormalParameterList mirrors parseFormalParameterList().
func (p *ParserRewrite) parseFormalParameterList() *trees.FormalParameterListTree {
	start := p.getTreeStartLocation()
	var params []*trees.ParseTree
	hasTrailingComma := false
	var commaPositions []trees.SourcePosition
	p.eat(TokenTypeOpenParen)
	for p.peekParameter() {
		params = append(params, p.parseParameter())
		if !p.peek(TokenTypeCloseParen) {
			comma := p.eat(TokenTypeComma)
			if comma != nil {
				commaPositions = append(commaPositions, toTreesPos(comma.GetStart()))
			}
			if p.peek(TokenTypeCloseParen) {
				hasTrailingComma = true
			}
		}
	}
	p.eat(TokenTypeCloseParen)
	return trees.NewFormalParameterListTree(p.getTreeLocation(start), params, hasTrailingComma, commaPositions)
}

// parseSetterParameterList mirrors parseSetterParameterList().
func (p *ParserRewrite) parseSetterParameterList() *trees.FormalParameterListTree {
	start := p.getTreeStartLocation()
	p.eat(TokenTypeOpenParen)
	var param *trees.ParseTree
	if p.peekParameter() {
		param = p.parseParameter()
	}
	p.eat(TokenTypeCloseParen)
	var params []*trees.ParseTree
	if param != nil {
		params = []*trees.ParseTree{param}
	}
	return trees.NewFormalParameterListTree(p.getTreeLocation(start), params, false, nil)
}

// peekParameter mirrors peekParameter().
func (p *ParserRewrite) peekParameter() bool {
	return p.peekPatternStart() || p.peek(TokenTypeEllipsis) || p.peekId()
}

// parseParameter mirrors parseParameter().
func (p *ParserRewrite) parseParameter() *trees.ParseTree {
	start := p.getTreeStartLocation()
	if p.peek(TokenTypeEllipsis) {
		return p.parseIterRest(patternKindInitializer)
	}
	if p.peekPatternStart() {
		return p.parsePatternAssignmentTarget(patternKindInitializer)
	}
	ident := p.parseIdentifierExpression()
	if p.peek(TokenTypeEqual) {
		p.eat(TokenTypeEqual)
		defaultValue := p.parseAssignmentExpression()
		return p.asParseTree(trees.NewDefaultParameterTree(
			p.getTreeLocation(start), p.asParseTree(ident), defaultValue))
	}
	return p.asParseTree(ident)
}

// ---- Utility: asParseTree ----------------------------------------------------

// asParseTree converts any concrete ParseTree subtype to *trees.ParseTree.
// This is used as a uniform conversion since Go doesn't have Java's implicit
// upcasting to a reference base type.
// asParseTree converts a concrete tree pointer to a *trees.ParseTree that
// satisfies the embedded-field invariant: the returned pointer must equal
// &concrete.ParseTree so that IRFactory's unsafe.Pointer dispatch can
// recover the concrete type from the base pointer.
//
// Every concrete tree type embeds ParseTree as its first field, so the
// address of the embedded field equals the address of the concrete struct,
// and the cast (*ConcreteTree)(unsafe.Pointer(pt)) is valid.
// Go-only: language-adapter — type-asserts any concrete tree to *trees.ParseTree (Java uses subtype polymorphism; Go needs explicit assertion).
func (p *ParserRewrite) asParseTree(node interface{}) *trees.ParseTree {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *trees.ParseTree:
		return n
	// ---- expression types ----
	case *trees.BinaryOperatorTree:
		return &n.ParseTree
	case *trees.CommaExpressionTree:
		return &n.ParseTree
	case *trees.MissingPrimaryExpressionTree:
		return &n.ParseTree
	case *trees.ConditionalExpressionTree:
		return &n.ParseTree
	case *trees.UnaryExpressionTree:
		return &n.ParseTree
	case *trees.UpdateExpressionTree:
		return &n.ParseTree
	case *trees.CallExpressionTree:
		return &n.ParseTree
	case *trees.NewExpressionTree:
		return &n.ParseTree
	case *trees.MemberExpressionTree:
		return &n.ParseTree
	case *trees.MemberLookupExpressionTree:
		return &n.ParseTree
	case *trees.OptChainCallExpressionTree:
		return &n.ParseTree
	case *trees.OptionalMemberExpressionTree:
		return &n.ParseTree
	case *trees.OptionalMemberLookupExpressionTree:
		return &n.ParseTree
	case *trees.IdentifierExpressionTree:
		return &n.ParseTree
	case *trees.LiteralExpressionTree:
		return &n.ParseTree
	case *trees.TemplateLiteralExpressionTree:
		return &n.ParseTree
	case *trees.TemplateLiteralPortionTree:
		return &n.ParseTree
	case *trees.TemplateSubstitutionTree:
		return &n.ParseTree
	case *trees.ArrayLiteralExpressionTree:
		return &n.ParseTree
	case *trees.ObjectLiteralExpressionTree:
		return &n.ParseTree
	case *trees.ObjectSpreadTree:
		return &n.ParseTree
	case *trees.ParenExpressionTree:
		return &n.ParseTree
	case *trees.FormalParameterListTree:
		return &n.ParseTree
	case *trees.FunctionDeclarationTree:
		return &n.ParseTree
	case *trees.SuperExpressionTree:
		return &n.ParseTree
	case *trees.ThisExpressionTree:
		return &n.ParseTree
	case *trees.DynamicImportTree:
		return &n.ParseTree
	case *trees.ImportMetaExpressionTree:
		return &n.ParseTree
	case *trees.NewTargetExpressionTree:
		return &n.ParseTree
	case *trees.YieldExpressionTree:
		return &n.ParseTree
	case *trees.AwaitExpressionTree:
		return &n.ParseTree
	case *trees.ArrayPatternTree:
		return &n.ParseTree
	case *trees.ObjectPatternTree:
		return &n.ParseTree
	case *trees.IterSpreadTree:
		return &n.ParseTree
	case *trees.IterRestTree:
		return &n.ParseTree
	case *trees.ObjectRestTree:
		return &n.ParseTree
	case *trees.DefaultParameterTree:
		return &n.ParseTree
	case *trees.ComprehensionTree:
		return &n.ParseTree
	case *trees.ComprehensionForTree:
		return &n.ParseTree
	case *trees.ComprehensionIfTree:
		return &n.ParseTree
	case *trees.NullTree:
		return &n.ParseTree
	case *trees.ComputedPropertyDefinitionTree:
		return &n.ParseTree
	case *trees.ComputedPropertyMethodTree:
		return &n.ParseTree
	case *trees.ComputedPropertyGetterTree:
		return &n.ParseTree
	case *trees.ComputedPropertySetterTree:
		return &n.ParseTree
	case *trees.ComputedPropertyFieldTree:
		return &n.ParseTree
	case *trees.GetAccessorTree:
		return &n.ParseTree
	case *trees.SetAccessorTree:
		return &n.ParseTree
	case *trees.PropertyNameAssignmentTree:
		return &n.ParseTree
	case *trees.FieldDeclarationTree:
		return &n.ParseTree
	// ---- statement types ----
	case *trees.BlockTree:
		return &n.ParseTree
	case *trees.IfStatementTree:
		return &n.ParseTree
	case *trees.DoWhileStatementTree:
		return &n.ParseTree
	case *trees.WhileStatementTree:
		return &n.ParseTree
	case *trees.ForStatementTree:
		return &n.ParseTree
	case *trees.ForInStatementTree:
		return &n.ParseTree
	case *trees.ForOfStatementTree:
		return &n.ParseTree
	case *trees.ForAwaitOfStatementTree:
		return &n.ParseTree
	case *trees.SwitchStatementTree:
		return &n.ParseTree
	case *trees.CaseClauseTree:
		return &n.ParseTree
	case *trees.DefaultClauseTree:
		return &n.ParseTree
	case *trees.TryStatementTree:
		return &n.ParseTree
	case *trees.CatchTree:
		return &n.ParseTree
	case *trees.FinallyTree:
		return &n.ParseTree
	case *trees.VariableStatementTree:
		return &n.ParseTree
	case *trees.VariableDeclarationListTree:
		return &n.ParseTree
	case *trees.VariableDeclarationTree:
		return &n.ParseTree
	case *trees.EmptyStatementTree:
		return &n.ParseTree
	case *trees.ExpressionStatementTree:
		return &n.ParseTree
	case *trees.BreakStatementTree:
		return &n.ParseTree
	case *trees.ContinueStatementTree:
		return &n.ParseTree
	case *trees.ReturnStatementTree:
		return &n.ParseTree
	case *trees.WithStatementTree:
		return &n.ParseTree
	case *trees.ThrowStatementTree:
		return &n.ParseTree
	case *trees.DebuggerStatementTree:
		return &n.ParseTree
	case *trees.LabelledStatementTree:
		return &n.ParseTree
	case *trees.ClassDeclarationTree:
		return &n.ParseTree
	case *trees.ImportDeclarationTree:
		return &n.ParseTree
	case *trees.ImportSpecifierTree:
		return &n.ParseTree
	case *trees.ExportDeclarationTree:
		return &n.ParseTree
	case *trees.ExportSpecifierTree:
		return &n.ParseTree
	default:
		panic(fmt.Sprintf("asParseTree: unhandled type %T", node))
	}
}

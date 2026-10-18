package parsing

import (
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// This opt-in experiment builds final Rhino nodes directly from Scanner tokens.
// The original parser remains the default and the fallback for every construct
// outside the explicit subset below. Read once before concurrent compilations.
var nativeParserExperimentEnabled = os.Getenv("JSGO_EXPERIMENT_NATIVE_PARSER") == "1"

// nativeParserDiagnostics never publishes speculative diagnostics. A lexical or
// IR warning/error rejects the entire native attempt, so the original pipeline
// produces the complete diagnostic sequence, with its original ordering.
type nativeParserDiagnostics struct {
	parser.ErrorReporterBase
	failed bool
}

func (r *nativeParserDiagnostics) ReportError(*parser.SourcePosition, string)   { r.failed = true }
func (r *nativeParserDiagnostics) ReportWarning(*parser.SourcePosition, string) { r.failed = true }
func (r *nativeParserDiagnostics) Warning(string, string, int, int)             { r.failed = true }
func (r *nativeParserDiagnostics) Error(string, string, int, int)               { r.failed = true }
func (r *nativeParserDiagnostics) RecordComment(parser.Comment_Type, *parser.SourceRange, string) {
	// Whole-file comments, JSDoc/fileoverview and closure-unaware annotations
	// require the original comment discovery/attachment pipeline, even if the
	// comment occurs after code that the native sink already constructed.
	r.failed = true
}

type nativeExpression struct {
	node       *rhino.Node
	start, end *parser.SourcePosition
}

// nativeParser has no intermediate per-expression objects. Its expression
// frames are returned by value, and only final Rhino nodes survive acceptance.
// Grammar: semicolon-terminated var/expression/empty statements; ordinary ES5
// unary/binary/assignment/conditional expressions, calls and member accesses;
// decimal numbers, unescaped strings, names/this/null/booleans; arrays without
// holes/trailing commas; object literals with explicit non-keyword keys.
// Functions, blocks, modules, ASI, patterns and extensions use the old parser.
type nativeParser struct {
	scanner     *parser.Scanner
	current     parser.Token
	last        parser.Token
	factory     *IRFactory
	diagnostics *nativeParserDiagnostics
	depth       int
}

func parserRunnerParseWithNativeFallback(sourceFile rhino.StaticSourceFile, source string,
	config *Config, reporter rhino.ErrorReporter) *ParserRunner_ParseResult {
	if result, ok := tryNativeParser(sourceFile, source, config); ok {
		return result
	}
	return parserRunnerParseUncached(sourceFile, source, config, reporter)
}

func tryNativeParser(sourceFile rhino.StaticSourceFile, source string, config *Config) (result *ParserRunner_ParseResult, accepted bool) {
	jutil.Try(func() {
		result, accepted = nativeParserAttempt(sourceFile, source, config)
	}, func(jutil.Throwable) {
		// The original entry point owns exception wrapping as well as normal
		// diagnostics. No speculative exception may bypass that boundary.
		result, accepted = nil, false
	})
	return result, accepted
}

func nativeParserAttempt(sourceFile rhino.StaticSourceFile, source string, config *Config) (*ParserRunner_ParseResult, bool) {
	if sourceFile.IsTypeScriptSource() || config.JsDocParsingMode().ShouldParseDescriptions() {
		return nil, false
	}
	file := parser.NewSourceFile(sourceFile.GetName(), source)
	diagnostics := &nativeParserDiagnostics{}
	diagnostics.InitErrorReporterBase(diagnostics)
	p := &nativeParser{diagnostics: diagnostics}
	p.scanner = parser.NewScanner(diagnostics, diagnostics, file, 0)
	p.factory = NewIRFactory(sourceFile, config, diagnostics, nil, file, JsDocInfoParser_JsDocSourceKind_NORMAL)
	start := p.scanner.GetPosition()
	p.advance()
	root := p.factory.NewNode(rhino.Token_SCRIPT)
	for !p.failed() && p.kind() != parser.TokenType_END_OF_FILE {
		statement := p.statement()
		if p.failed() {
			return nil, false
		}
		root.AddChildToBack(statement.node)
	}
	if p.failed() {
		return nil, false
	}
	p.factory.SetSourceInfoFromPositions(root, start, p.current.GetTokenBase().Location.End)
	p.factory.transformDispatcher.ParseDirectives(root)
	// goog.module() has whole-file module wrapping semantics; leave it on the
	// original path rather than publishing an incomplete module AST.
	if p.factory.transformDispatcher.IsGoogModuleFile(root) {
		return nil, false
	}
	p.factory.SetFileOverviewJsDoc(root)
	p.factory.ValidateAll(root)
	if p.failed() {
		return nil, false
	}
	features := p.factory.GetFeatures()
	root.PutProp(rhino.Node_FEATURE_SET, features)
	return NewParserRunner_ParseResultWithTopLevelStatementRanges(root, []*parser.Comment{}, []*parser.SourceRange{}, features, nil), true
}

func (p *nativeParser) failed() bool           { return p.diagnostics.failed }
func (p *nativeParser) reject()                { p.diagnostics.failed = true }
func (p *nativeParser) kind() parser.TokenType { return p.current.GetTokenBase().Type }
func (p *nativeParser) advance() {
	p.last = p.current
	p.current = p.scanner.NextToken()
	switch p.kind() {
	case parser.TokenType_IDENTIFIER:
		token := p.current.AsIdentifier()
		name := token.GetMaybePrivateValue()
		if token.IsPrivateIdentifier() || token.IsKeyword() || rhino.TokenStream_IsKeyword(name) ||
			name == "await" || name == "eval" || name == "arguments" || IRFactory_ES5_STRICT_RESERVED_KEYWORDS.Contains(name) {
			p.reject()
		}
	case parser.TokenType_STRING:
		text := p.current.AsLiteral().GetLiteralTokenBase().Value
		if strings.ContainsAny(text, "\\\u2028\u2029") {
			p.reject()
		}
	case parser.TokenType_NUMBER:
		text := p.current.AsLiteral().GetLiteralTokenBase().Value
		if strings.ContainsAny(text, "_xXbBoO") || len(text) > 1 && text[0] == '0' && text[1] >= '0' && text[1] <= '9' {
			p.reject()
		}
	}
}
func (p *nativeParser) take(kind parser.TokenType) bool {
	if p.failed() || p.kind() != kind {
		p.reject()
		return false
	}
	p.advance()
	return !p.failed()
}
func (p *nativeParser) located(node *rhino.Node, start, end *parser.SourcePosition) nativeExpression {
	p.factory.SetSourceInfoFromPositions(node, start, end)
	return nativeExpression{node, start, end}
}
func (p *nativeParser) statement() nativeExpression {
	start := p.current.GetStart()
	if p.kind() == parser.TokenType_SEMI_COLON {
		p.advance()
		return p.located(p.factory.NewNode(rhino.Token_EMPTY), start, p.last.GetTokenBase().Location.End)
	}
	if p.kind() == parser.TokenType_VAR {
		p.advance()
		declaration := p.factory.NewNode(rhino.Token_VAR)
		for !p.failed() {
			if p.kind() != parser.TokenType_IDENTIFIER {
				p.reject()
				break
			}
			token := p.current.AsIdentifier()
			name := p.factory.transformDispatcher.ProcessNameWithInlineCommentsOfToken(token, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
			p.advance()
			if p.kind() == parser.TokenType_EQUAL {
				p.advance()
				initializer := p.expression(1)
				if p.failed() {
					break
				}
				name.AddChildToBack(initializer.node)
				p.factory.SetLength(name, token.GetStart(), initializer.end)
			}
			declaration.AddChildToBack(name)
			if p.kind() != parser.TokenType_COMMA {
				break
			}
			p.advance()
		}
		if !p.take(parser.TokenType_SEMI_COLON) {
			return nativeExpression{}
		}
		return p.located(declaration, start, p.last.GetTokenBase().Location.End)
	}
	// At statement start, an object-looking brace denotes a block/label.
	if p.kind() == parser.TokenType_OPEN_CURLY {
		p.reject()
		return nativeExpression{}
	}
	expr := p.expression(1)
	if !p.take(parser.TokenType_SEMI_COLON) {
		return nativeExpression{}
	}
	return p.located(p.factory.NewNodeWithChild(rhino.Token_EXPR_RESULT, expr.node), start, p.last.GetTokenBase().Location.End)
}

func nativeBinaryPrecedence(kind parser.TokenType) int {
	switch kind {
	case parser.TokenType_EQUAL:
		return 1
	case parser.TokenType_OR:
		return 3
	case parser.TokenType_AND:
		return 4
	case parser.TokenType_BAR:
		return 5
	case parser.TokenType_CARET:
		return 6
	case parser.TokenType_AMPERSAND:
		return 7
	case parser.TokenType_EQUAL_EQUAL, parser.TokenType_NOT_EQUAL, parser.TokenType_EQUAL_EQUAL_EQUAL, parser.TokenType_NOT_EQUAL_EQUAL:
		return 8
	case parser.TokenType_OPEN_ANGLE, parser.TokenType_CLOSE_ANGLE, parser.TokenType_LESS_EQUAL, parser.TokenType_GREATER_EQUAL, parser.TokenType_IN, parser.TokenType_INSTANCEOF:
		return 9
	case parser.TokenType_LEFT_SHIFT, parser.TokenType_RIGHT_SHIFT, parser.TokenType_UNSIGNED_RIGHT_SHIFT:
		return 10
	case parser.TokenType_PLUS, parser.TokenType_MINUS:
		return 11
	case parser.TokenType_STAR, parser.TokenType_SLASH, parser.TokenType_PERCENT:
		return 12
	default:
		return 0
	}
}
func (p *nativeParser) expression(minimum int) nativeExpression {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > 128 || p.failed() {
		p.reject()
		return nativeExpression{}
	}
	left := p.prefix()
	for !p.failed() {
		kind := p.kind()
		if kind == parser.TokenType_QUESTION && minimum <= 2 {
			p.advance()
			middle := p.expression(1)
			if !p.take(parser.TokenType_COLON) {
				return nativeExpression{}
			}
			right := p.expression(1)
			if p.failed() {
				return nativeExpression{}
			}
			left = p.located(p.factory.NewNodeWith3Children(rhino.Token_HOOK, left.node, middle.node, right.node), left.start, right.end)
			continue
		}
		precedence := nativeBinaryPrecedence(kind)
		if precedence == 0 || precedence < minimum {
			break
		}
		if kind == parser.TokenType_EQUAL && !left.node.IsValidAssignmentTarget() {
			p.reject()
			return nativeExpression{}
		}
		p.advance()
		next := precedence + 1
		if kind == parser.TokenType_EQUAL {
			next = precedence
		}
		right := p.expression(next)
		if p.failed() {
			return nativeExpression{}
		}
		left = p.located(p.factory.NewNodeWith2Children(IRFactory_TransformBinaryTokenType(kind), left.node, right.node), left.start, right.end)
	}
	return left
}
func (p *nativeParser) prefix() nativeExpression {
	if p.failed() {
		return nativeExpression{}
	}
	token := p.current
	start := token.GetStart()
	kind := p.kind()
	switch kind {
	case parser.TokenType_PLUS, parser.TokenType_MINUS, parser.TokenType_BANG, parser.TokenType_TILDE, parser.TokenType_VOID, parser.TokenType_TYPEOF:
		p.advance()
		operand := p.expression(13)
		if p.failed() {
			return nativeExpression{}
		}
		return p.located(p.factory.NewNodeWithChild(IRFactory_TransformUnaryTokenType(kind), operand.node), start, operand.end)
	}
	var result nativeExpression
	switch kind {
	case parser.TokenType_IDENTIFIER:
		node := p.factory.transformDispatcher.ProcessNameOfToken(token.AsIdentifier(), rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
		p.advance()
		result = nativeExpression{node, start, token.GetTokenBase().Location.End}
	case parser.TokenType_NUMBER:
		node := p.factory.NewNumberNode(p.factory.NormalizeNumber(token.AsLiteral()))
		p.advance()
		result = p.located(node, start, token.GetTokenBase().Location.End)
	case parser.TokenType_STRING:
		node := p.factory.transformDispatcher.ProcessString(token.AsLiteral())
		p.advance()
		result = nativeExpression{node, start, token.GetTokenBase().Location.End}
	case parser.TokenType_TRUE, parser.TokenType_FALSE, parser.TokenType_NULL, parser.TokenType_THIS:
		var nodeKind rhino.Token
		switch kind {
		case parser.TokenType_TRUE:
			nodeKind = rhino.Token_TRUE
		case parser.TokenType_FALSE:
			nodeKind = rhino.Token_FALSE
		case parser.TokenType_NULL:
			nodeKind = rhino.Token_NULL
		case parser.TokenType_THIS:
			nodeKind = rhino.Token_THIS
		}
		p.advance()
		result = p.located(p.factory.NewNode(nodeKind), start, token.GetTokenBase().Location.End)
	case parser.TokenType_OPEN_PAREN:
		p.advance()
		result = p.expression(1)
		if !p.take(parser.TokenType_CLOSE_PAREN) {
			return nativeExpression{}
		}
		result.node.SetIsParenthesized(true)
		result.start, result.end = start, p.last.GetTokenBase().Location.End
	case parser.TokenType_OPEN_SQUARE, parser.TokenType_OPEN_CURLY:
		result = p.literalContainer()
	default:
		p.reject()
		return nativeExpression{}
	}
	for !p.failed() {
		switch p.kind() {
		case parser.TokenType_PERIOD:
			p.advance()
			if p.failed() || p.kind() != parser.TokenType_IDENTIFIER {
				p.reject()
				return nativeExpression{}
			}
			property := p.current
			node := p.factory.transformDispatcher.NewStringNodeFromIdentifier(rhino.Token_GETPROP, IRFactory_TransformDispatcher_IdentifierType_STANDARD, property.AsIdentifier())
			node.AddChildToBack(result.node)
			p.factory.SetSourceInfoFromToken(node, property)
			p.advance()
			result = nativeExpression{node, result.start, property.GetTokenBase().Location.End}
		case parser.TokenType_OPEN_SQUARE:
			p.advance()
			index := p.expression(1)
			if !p.take(parser.TokenType_CLOSE_SQUARE) {
				return nativeExpression{}
			}
			result = p.located(p.factory.NewNodeWith2Children(rhino.Token_GETELEM, result.node, index.node), result.start, p.last.GetTokenBase().Location.End)
		case parser.TokenType_OPEN_PAREN:
			p.advance()
			call := p.factory.NewNodeWithChild(rhino.Token_CALL, result.node)
			for !p.failed() && p.kind() != parser.TokenType_CLOSE_PAREN {
				arg := p.expression(1)
				if p.failed() {
					return nativeExpression{}
				}
				call.AddChildToBack(arg.node)
				if p.kind() != parser.TokenType_COMMA {
					break
				}
				p.advance()
				if p.kind() == parser.TokenType_CLOSE_PAREN {
					p.reject()
					return nativeExpression{}
				}
			}
			if !p.take(parser.TokenType_CLOSE_PAREN) {
				return nativeExpression{}
			}
			p.factory.transformDispatcher.AnnotateCalls(call)
			result = p.located(call, result.start, p.last.GetTokenBase().Location.End)
		default:
			return result
		}
	}
	return nativeExpression{}
}
func (p *nativeParser) literalContainer() nativeExpression {
	start := p.current.GetStart()
	object := p.kind() == parser.TokenType_OPEN_CURLY
	endKind, nodeKind := parser.TokenType_CLOSE_SQUARE, rhino.Token_ARRAYLIT
	if object {
		endKind, nodeKind = parser.TokenType_CLOSE_CURLY, rhino.Token_OBJECTLIT
	}
	p.advance()
	container := p.factory.NewNode(nodeKind)
	for !p.failed() && p.kind() != endKind {
		var key *rhino.Node
		if object {
			switch p.kind() {
			case parser.TokenType_IDENTIFIER, parser.TokenType_STRING, parser.TokenType_NUMBER:
				key = p.factory.transformDispatcher.ProcessObjectLitKey(p.current, rhino.Token_STRING_KEY, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
			default:
				p.reject()
				return nativeExpression{}
			}
			p.advance()
			if !p.take(parser.TokenType_COLON) {
				return nativeExpression{}
			}
		}
		value := p.expression(1)
		if p.failed() {
			return nativeExpression{}
		}
		if object {
			key.AddChildToBack(value.node)
			container.AddChildToBack(key)
		} else {
			container.AddChildToBack(value.node)
		}
		if p.kind() != parser.TokenType_COMMA {
			break
		}
		p.advance()
		if p.kind() == endKind {
			p.reject()
			return nativeExpression{}
		}
	}
	if !p.take(endKind) {
		return nativeExpression{}
	}
	return p.located(container, start, p.last.GetTokenBase().Location.End)
}

// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20260415
//
// Statement-loop parser. The grammar is still narrow — enough for
// Tier-0 (primitive @type decls + primitive literal assigns) and
// Tier-1 (function decls with @param/@return + simple call sites).
// Unrecognised statements are skipped to the next boundary so a single
// unsupported construct doesn't poison neighbouring statements we can
// handle. Parse never panics (guardrail #6).

package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type Parser struct {
	src     []rune
	pos     int
	line    int // 1-indexed
	lineBOL int // absolute pos of start of current line, for column calc
}

func Parse(src string) (n *rhino.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse panic recovered: %v", r)
		}
	}()
	p := &Parser{src: []rune(src), line: 1}
	script := rhino.New(rhino.SCRIPT)
	script.SetPos(1, 0)
	p.parseStatements(script, 0)
	return script, nil
}

// parseStatements fills parent with top-level / block statements until
// EOF (if endTok == 0) or until a '}' is hit at the current depth.
func (p *Parser) parseStatements(parent *rhino.Node, endTok rune) {
	for {
		p.skipTrivia()
		if p.atEnd() {
			return
		}
		if endTok != 0 && p.src[p.pos] == endTok {
			return
		}
		start := p.pos
		stmt, ok := p.tryStatement()
		if ok && stmt != nil {
			parent.AddChildToBack(stmt)
			continue
		}
		if p.pos == start {
			p.skipToStatementBoundary()
		}
	}
}

// tryStatement attempts one statement. Returns (stmt, true) on success
// or (nil, false) otherwise; in either case p.pos is left past the
// consumed input (or rewound on speculative failures).
func (p *Parser) tryStatement() (*rhino.Node, bool) {
	var jsdoc *rhino.JSDocInfo
	if body, ok := p.consumeBlockComment(); ok {
		jsdoc = rhino.ParseJSDoc(body)
		p.skipTrivia()
	}

	if p.consumeKeyword("function") {
		return p.finishFunction(jsdoc)
	}
	if declToken, ok := p.consumeDeclKeyword(); ok {
		return p.finishDeclStatement(declToken, jsdoc)
	}
	if p.consumeKeyword("return") {
		return p.finishReturn(), true
	}

	// JSDoc found but not attached to a supported declaration; skip.
	if jsdoc != nil {
		p.skipToStatementBoundary()
		return nil, false
	}

	if stmt, ok := p.tryCallOrAssignStatement(); ok {
		return stmt, true
	}
	return nil, false
}

func (p *Parser) finishDeclStatement(declToken rhino.Token, jsdoc *rhino.JSDocInfo) (*rhino.Node, bool) {
	declNode := rhino.New(declToken)
	if jsdoc != nil {
		declNode.SetJSDocInfo(jsdoc)
	}

	p.skipTrivia()
	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		p.skipToStatementBoundary()
		return nil, false
	}
	nameNode := rhino.NewString(rhino.NAME, ident)
	nameNode.SetPos(nameLine, nameCol)
	declNode.AddChildToBack(nameNode)

	p.skipTrivia()
	if p.consumeByte('=') {
		if !p.atEnd() && p.src[p.pos] == '=' {
			p.skipToStatementBoundary()
			return declNode, true
		}
		p.skipTrivia()
		initLine, initCol := p.line, p.col()
		init, ok := p.tryExpression()
		if !ok {
			p.skipToStatementBoundary()
			return declNode, true
		}
		p.skipTrivia()
		if p.atEnd() || p.src[p.pos] != ';' {
			p.skipToStatementBoundary()
			return declNode, true
		}
		init.SetPos(initLine, initCol)
		nameNode.AddChildToBack(init)
	}
	p.skipTrivia()
	if !p.consumeByte(';') {
		p.skipToStatementBoundary()
	}
	return declNode, true
}

// finishFunction parses `function NAME(params) { body }` after the
// `function` keyword is consumed. Attaches the caller-provided JSDoc
// (carrying @param / @return).
func (p *Parser) finishFunction(jsdoc *rhino.JSDocInfo) (*rhino.Node, bool) {
	fn := rhino.New(rhino.FUNCTION)
	if jsdoc != nil {
		fn.SetJSDocInfo(jsdoc)
	}

	p.skipTrivia()
	nameLine, nameCol := p.line, p.col()
	name, ok := p.consumeIdent()
	nameNode := rhino.NewString(rhino.NAME, "")
	if ok {
		nameNode = rhino.NewString(rhino.NAME, name)
	}
	nameNode.SetPos(nameLine, nameCol)
	fn.AddChildToBack(nameNode)

	p.skipTrivia()
	if !p.consumeByte('(') {
		p.skipToStatementBoundary()
		return fn, true
	}
	params := rhino.New(rhino.PARAM_LIST)
	for {
		p.skipTrivia()
		if p.atEnd() {
			return fn, true
		}
		if p.src[p.pos] == ')' {
			p.advance()
			break
		}
		pLine, pCol := p.line, p.col()
		pname, ok := p.consumeIdent()
		if !ok {
			p.skipToStatementBoundary()
			return fn, true
		}
		pNode := rhino.NewString(rhino.NAME, pname)
		pNode.SetPos(pLine, pCol)
		params.AddChildToBack(pNode)
		p.skipTrivia()
		if p.consumeByte(',') {
			continue
		}
	}
	fn.AddChildToBack(params)

	p.skipTrivia()
	if !p.consumeByte('{') {
		p.skipToStatementBoundary()
		return fn, true
	}
	block := rhino.New(rhino.BLOCK)
	p.parseStatements(block, '}')
	p.consumeByte('}')
	fn.AddChildToBack(block)
	return fn, true
}

func (p *Parser) finishReturn() *rhino.Node {
	ret := rhino.New(rhino.RETURN)
	p.skipTrivia()
	if !p.atEnd() && p.src[p.pos] != ';' {
		eline, ecol := p.line, p.col()
		if expr, ok := p.tryExpression(); ok {
			expr.SetPos(eline, ecol)
			ret.AddChildToBack(expr)
		} else {
			p.skipToStatementBoundary()
			return ret
		}
	}
	p.skipTrivia()
	p.consumeByte(';')
	return ret
}

// tryCallOrAssignStatement parses either `IDENT(args);` or
// `IDENT = LITERAL;` — both speculative, rewinding on mismatch.
func (p *Parser) tryCallOrAssignStatement() (*rhino.Node, bool) {
	save := p.pos
	saveLine := p.line
	saveBOL := p.lineBOL
	rewind := func() {
		p.pos = save
		p.line = saveLine
		p.lineBOL = saveBOL
	}

	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		rewind()
		return nil, false
	}
	p.skipTrivia()
	if p.atEnd() {
		rewind()
		return nil, false
	}

	switch p.src[p.pos] {
	case '(':
		p.advance()
		call := rhino.New(rhino.CALL)
		call.SetPos(nameLine, nameCol)
		callee := rhino.NewString(rhino.NAME, ident)
		callee.SetPos(nameLine, nameCol)
		call.AddChildToBack(callee)
		for {
			p.skipTrivia()
			if p.atEnd() {
				rewind()
				return nil, false
			}
			if p.src[p.pos] == ')' {
				p.advance()
				break
			}
			argLine, argCol := p.line, p.col()
			arg, ok := p.tryExpression()
			if !ok {
				rewind()
				return nil, false
			}
			arg.SetPos(argLine, argCol)
			call.AddChildToBack(arg)
			p.skipTrivia()
			if p.consumeByte(',') {
				continue
			}
		}
		p.skipTrivia()
		if !p.consumeByte(';') {
			// Call without trailing `;` (e.g. last statement in block)
			// — still keep the node; its position is already set.
		}
		stmt := rhino.New(rhino.EXPR_RESULT)
		stmt.SetPos(nameLine, nameCol)
		stmt.AddChildToBack(call)
		return stmt, true

	case '=':
		p.advance()
		if !p.atEnd() && (p.src[p.pos] == '=' || p.src[p.pos] == '>') {
			rewind()
			return nil, false
		}
		p.skipTrivia()
		rhsLine, rhsCol := p.line, p.col()
		rhs, ok := p.tryExpression()
		if !ok {
			rewind()
			return nil, false
		}
		p.skipTrivia()
		if p.atEnd() || p.src[p.pos] != ';' {
			rewind()
			return nil, false
		}
		p.advance()
		nameNode := rhino.NewString(rhino.NAME, ident)
		nameNode.SetPos(nameLine, nameCol)
		rhs.SetPos(rhsLine, rhsCol)
		assign := rhino.New(rhino.ASSIGN)
		assign.SetPos(nameLine, nameCol)
		assign.AddChildToBack(nameNode)
		assign.AddChildToBack(rhs)
		stmt := rhino.New(rhino.EXPR_RESULT)
		stmt.SetPos(nameLine, nameCol)
		stmt.AddChildToBack(assign)
		return stmt, true
	}
	rewind()
	return nil, false
}

// tryExpression parses one expression: a literal, an identifier
// reference, or a call. Anything else returns (nil, false).
func (p *Parser) tryExpression() (*rhino.Node, bool) {
	if p.atEnd() {
		return nil, false
	}
	c := p.src[p.pos]
	switch {
	case c == '\'' || c == '"':
		return p.parseStringLiteral(c)
	case unicode.IsDigit(c):
		return p.parseNumberLiteral()
	case c == 't' || c == 'f':
		if n, ok := p.parseBooleanLiteral(); ok {
			return n, true
		}
	case c == 'n':
		if p.consumeKeyword("null") {
			return rhino.New(rhino.NULL), true
		}
	case c == 'u':
		if p.consumeKeyword("undefined") {
			return rhino.New(rhino.UNDEFINED), true
		}
	}
	if unicode.IsLetter(c) || c == '_' || c == '$' {
		return p.parseNameOrCall()
	}
	return nil, false
}

// parseNameOrCall parses `IDENT` or `IDENT(args)`.
func (p *Parser) parseNameOrCall() (*rhino.Node, bool) {
	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		return nil, false
	}
	name := rhino.NewString(rhino.NAME, ident)
	name.SetPos(nameLine, nameCol)
	p.skipTrivia()
	if p.atEnd() || p.src[p.pos] != '(' {
		return name, true
	}
	p.advance()
	call := rhino.New(rhino.CALL)
	call.SetPos(nameLine, nameCol)
	call.AddChildToBack(name)
	for {
		p.skipTrivia()
		if p.atEnd() {
			return nil, false
		}
		if p.src[p.pos] == ')' {
			p.advance()
			return call, true
		}
		argLine, argCol := p.line, p.col()
		arg, ok := p.tryExpression()
		if !ok {
			return nil, false
		}
		arg.SetPos(argLine, argCol)
		call.AddChildToBack(arg)
		p.skipTrivia()
		if p.consumeByte(',') {
			continue
		}
	}
}

func (p *Parser) parseBooleanLiteral() (*rhino.Node, bool) {
	if p.consumeKeyword("true") {
		return rhino.New(rhino.TRUE), true
	}
	if p.consumeKeyword("false") {
		return rhino.New(rhino.FALSE), true
	}
	return nil, false
}

func (p *Parser) parseStringLiteral(quote rune) (*rhino.Node, bool) {
	p.advance() // opening quote
	start := p.pos
	for !p.atEnd() && p.src[p.pos] != quote {
		if p.src[p.pos] == '\n' {
			return nil, false
		}
		if p.src[p.pos] == '\\' && p.pos+1 < len(p.src) {
			p.advance()
			p.advance()
			continue
		}
		p.advance()
	}
	if p.atEnd() {
		return nil, false
	}
	val := string(p.src[start:p.pos])
	p.advance() // closing quote
	return rhino.NewString(rhino.STRINGLIT, val), true
}

func (p *Parser) parseNumberLiteral() (*rhino.Node, bool) {
	start := p.pos
	for !p.atEnd() && (unicode.IsDigit(p.src[p.pos]) || p.src[p.pos] == '.') {
		p.advance()
	}
	if p.pos == start {
		return nil, false
	}
	v, err := strconv.ParseFloat(string(p.src[start:p.pos]), 64)
	if err != nil {
		return nil, false
	}
	return rhino.NewNumber(v), true
}

func (p *Parser) consumeBlockComment() (string, bool) {
	if p.pos+1 >= len(p.src) || p.src[p.pos] != '/' || p.src[p.pos+1] != '*' {
		return "", false
	}
	p.advance()
	p.advance()
	start := p.pos
	for p.pos+1 < len(p.src) && !(p.src[p.pos] == '*' && p.src[p.pos+1] == '/') {
		p.advance()
	}
	if p.pos+1 >= len(p.src) {
		return "", false
	}
	body := string(p.src[start:p.pos])
	p.advance()
	p.advance()
	return body, true
}

func (p *Parser) consumeLineComment() bool {
	if p.pos+1 < len(p.src) && p.src[p.pos] == '/' && p.src[p.pos+1] == '/' {
		for !p.atEnd() && p.src[p.pos] != '\n' {
			p.advance()
		}
		return true
	}
	return false
}

func (p *Parser) consumeDeclKeyword() (rhino.Token, bool) {
	if p.consumeKeyword("var") {
		return rhino.VAR, true
	}
	if p.consumeKeyword("let") {
		return rhino.LET, true
	}
	if p.consumeKeyword("const") {
		return rhino.CONST, true
	}
	return rhino.ERROR, false
}

func (p *Parser) consumeKeyword(kw string) bool {
	rs := []rune(kw)
	if p.pos+len(rs) > len(p.src) {
		return false
	}
	for i, r := range rs {
		if p.src[p.pos+i] != r {
			return false
		}
	}
	if p.pos+len(rs) < len(p.src) {
		nxt := p.src[p.pos+len(rs)]
		if unicode.IsLetter(nxt) || unicode.IsDigit(nxt) || nxt == '_' {
			return false
		}
	}
	for i := 0; i < len(rs); i++ {
		p.advance()
	}
	return true
}

func (p *Parser) consumeIdent() (string, bool) {
	if p.atEnd() {
		return "", false
	}
	c := p.src[p.pos]
	if !(unicode.IsLetter(c) || c == '_' || c == '$') {
		return "", false
	}
	start := p.pos
	for !p.atEnd() {
		c := p.src[p.pos]
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$') {
			break
		}
		p.advance()
	}
	return string(p.src[start:p.pos]), true
}

func (p *Parser) consumeByte(b byte) bool {
	if p.atEnd() || byte(p.src[p.pos]) != b {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) skipTrivia() {
	for !p.atEnd() {
		c := p.src[p.pos]
		if unicode.IsSpace(c) {
			p.advance()
			continue
		}
		if p.consumeLineComment() {
			continue
		}
		if p.pos+2 < len(p.src) && c == '/' && p.src[p.pos+1] == '*' && p.src[p.pos+2] != '*' {
			p.advance()
			p.advance()
			for p.pos+1 < len(p.src) && !(p.src[p.pos] == '*' && p.src[p.pos+1] == '/') {
				p.advance()
			}
			if p.pos+1 < len(p.src) {
				p.advance()
				p.advance()
			}
			continue
		}
		break
	}
}

func (p *Parser) skipToStatementBoundary() {
	depth := 0
	for !p.atEnd() {
		c := p.src[p.pos]
		switch c {
		case ';':
			p.advance()
			if depth == 0 {
				return
			}
		case '{':
			depth++
			p.advance()
		case '}':
			if depth == 0 {
				return
			}
			depth--
			p.advance()
			if depth == 0 {
				return
			}
		case '/':
			if p.pos+1 < len(p.src) && p.src[p.pos+1] == '/' {
				for !p.atEnd() && p.src[p.pos] != '\n' {
					p.advance()
				}
			} else if p.pos+1 < len(p.src) && p.src[p.pos+1] == '*' {
				p.advance()
				p.advance()
				for p.pos+1 < len(p.src) && !(p.src[p.pos] == '*' && p.src[p.pos+1] == '/') {
					p.advance()
				}
				if p.pos+1 < len(p.src) {
					p.advance()
					p.advance()
				}
			} else {
				p.advance()
			}
		case '"', '\'':
			p.advance()
			for !p.atEnd() && p.src[p.pos] != c {
				if p.src[p.pos] == '\\' && p.pos+1 < len(p.src) {
					p.advance()
					p.advance()
					continue
				}
				if p.src[p.pos] == '\n' {
					break
				}
				p.advance()
			}
			if !p.atEnd() {
				p.advance()
			}
		default:
			p.advance()
		}
	}
}

func (p *Parser) advance() {
	if p.atEnd() {
		return
	}
	if p.src[p.pos] == '\n' {
		p.line++
		p.lineBOL = p.pos + 1
	}
	p.pos++
}

func (p *Parser) col() int    { return p.pos - p.lineBOL }
func (p *Parser) atEnd() bool { return p.pos >= len(p.src) }

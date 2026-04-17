// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20260415
//
// Statement-loop parser. Each top-level statement is tried; on
// recognition failure we skip to the next statement boundary so that
// unrecognised syntax doesn't poison statements we *can* handle. The
// grammar currently recognises:
//
//     /** @type {<name>} */
//     (var|let|const) IDENT [= LITERAL] ;
//
// Literals: string, number, boolean. Anything else is skipped to the
// next ';' or '}'. Parse never panics (guardrail #6).

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

	// Parse statements until EOF. Unrecognised statements are silently
	// skipped so that subsequent recognised statements still contribute.
	for {
		p.skipTrivia()
		if p.atEnd() {
			break
		}
		start := p.pos
		stmt, ok := p.tryStatement()
		if ok && stmt != nil {
			script.AddChildToBack(stmt)
			continue
		}
		// Didn't recognise; advance to next boundary to avoid looping.
		if p.pos == start {
			p.skipToStatementBoundary()
		}
	}
	return script, nil
}

// tryStatement attempts to parse one statement. Returns (stmt, true) on
// success, or (nil, false) with p.pos advanced past the current
// statement when the input is unrecognised.
func (p *Parser) tryStatement() (*rhino.Node, bool) {
	var jsdoc *rhino.JSDocInfo
	if body, ok := p.consumeBlockComment(); ok {
		jsdoc = rhino.ParseJSDoc(body)
		p.skipTrivia()
	}

	if declToken, ok := p.consumeDeclKeyword(); ok {
		return p.finishDeclStatement(declToken, jsdoc)
	}

	// JSDoc found but not attached to a decl we handle; skip statement.
	if jsdoc != nil {
		p.skipToStatementBoundary()
		return nil, false
	}

	// Try bare assignment: IDENT = LITERAL ;
	if stmt, ok := p.tryAssignStatement(); ok {
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
		// Refuse `==`/`===`/... — we only handle simple assignment.
		if !p.atEnd() && p.src[p.pos] == '=' {
			p.skipToStatementBoundary()
			return declNode, true
		}
		p.skipTrivia()
		initLine, initCol := p.line, p.col()
		init, ok := p.tryLiteral()
		if !ok {
			p.skipToStatementBoundary()
			return declNode, true
		}
		// Only attach if the statement ends right after. Stops us from
		// mis-typing `var x = 1 < 2;` as `var x = 1;`.
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

// tryAssignStatement recognises `IDENT = LITERAL ;` with nothing after
// the literal before `;`. Returns EXPR_RESULT(ASSIGN(NAME, LITERAL)) to
// mirror Closure's AST shape. Rewinds on mismatch so the caller can try
// other statement forms.
func (p *Parser) tryAssignStatement() (*rhino.Node, bool) {
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
	if !p.consumeByte('=') {
		rewind()
		return nil, false
	}
	// Reject `==`, `===`, `=>`, `+=`, etc. The caller gets us back to
	// start so the statement can be skipped.
	if !p.atEnd() && (p.src[p.pos] == '=' || p.src[p.pos] == '>') {
		rewind()
		return nil, false
	}
	p.skipTrivia()
	rhsLine, rhsCol := p.line, p.col()
	rhs, ok := p.tryLiteral()
	if !ok {
		rewind()
		return nil, false
	}
	p.skipTrivia()
	if p.atEnd() || p.src[p.pos] != ';' {
		rewind()
		return nil, false
	}
	p.advance() // consume ';'

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

func (p *Parser) tryLiteral() (*rhino.Node, bool) {
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
	return nil, false
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

// skipTrivia skips whitespace plus line/block comments that aren't JSDoc.
// Block comments starting with /** are left in place for tryStatement to
// pick up as JSDoc.
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
			// Plain block comment (not JSDoc). Skip.
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

// skipToStatementBoundary walks forward until just past the next `;` or
// `}`, or to end of input. Uses a small balanced-brace scan so we don't
// stop inside braces. Strings and comments are handled coarsely.
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
				p.advance()
				return
			}
			depth--
			p.advance()
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

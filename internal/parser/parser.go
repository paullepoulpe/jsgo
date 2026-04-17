// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20260415
//
// Tier-0 stub. Recognises only:
//
//     /** @type {<name>} */
//     (var|let|const) IDENT [= LITERAL] ;
//
// where LITERAL is a string, number, or boolean literal. Anything else
// returns a parse error without panicking (guardrail #6). The grammar
// grows incrementally as corpus tiers require it.

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

func Parse(src string) (*rhino.Node, error) {
	p := &Parser{src: []rune(src), line: 1}
	script := rhino.New(rhino.SCRIPT)
	script.SetPos(1, 0)

	p.skipWhitespaceOnly()
	var jsdoc *rhino.JSDocInfo
	if body, ok := p.consumeBlockComment(); ok {
		jsdoc = rhino.ParseJSDoc(body)
	}

	p.skipWhitespaceOnly()
	declToken, ok := p.consumeDeclKeyword()
	if !ok {
		return nil, p.errorf("expected 'var', 'let' or 'const'")
	}

	declNode := rhino.New(declToken)
	if jsdoc != nil {
		declNode.SetJSDocInfo(jsdoc)
	}
	script.AddChildToBack(declNode)

	p.skipWhitespaceOnly()
	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		return nil, p.errorf("expected identifier")
	}
	nameNode := rhino.NewString(rhino.NAME, ident)
	nameNode.SetPos(nameLine, nameCol)
	declNode.AddChildToBack(nameNode)

	p.skipWhitespaceOnly()
	if p.consumeByte('=') {
		p.skipWhitespaceOnly()
		initLine, initCol := p.line, p.col()
		init, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		init.SetPos(initLine, initCol)
		nameNode.AddChildToBack(init)
		p.skipWhitespaceOnly()
	}

	if !p.consumeByte(';') {
		return nil, p.errorf("expected ';'")
	}

	return script, nil
}

func (p *Parser) parseLiteral() (*rhino.Node, error) {
	if p.atEnd() {
		return nil, p.errorf("expected literal")
	}
	c := p.src[p.pos]
	switch {
	case c == '\'' || c == '"':
		return p.parseStringLiteral(c)
	case unicode.IsDigit(c):
		return p.parseNumberLiteral()
	case c == 't' || c == 'f':
		if n, ok := p.parseBooleanLiteral(); ok {
			return n, nil
		}
	}
	return nil, p.errorf("unsupported literal starting with %q", c)
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

func (p *Parser) parseStringLiteral(quote rune) (*rhino.Node, error) {
	p.advance() // opening quote
	start := p.pos
	for !p.atEnd() && p.src[p.pos] != quote {
		if p.src[p.pos] == '\n' {
			return nil, p.errorf("unterminated string literal")
		}
		p.advance()
	}
	if p.atEnd() {
		return nil, p.errorf("unterminated string literal")
	}
	val := string(p.src[start:p.pos])
	p.advance() // closing quote
	return rhino.NewString(rhino.STRINGLIT, val), nil
}

func (p *Parser) parseNumberLiteral() (*rhino.Node, error) {
	start := p.pos
	for !p.atEnd() && (unicode.IsDigit(p.src[p.pos]) || p.src[p.pos] == '.') {
		p.advance()
	}
	v, err := strconv.ParseFloat(string(p.src[start:p.pos]), 64)
	if err != nil {
		return nil, p.errorf("invalid number literal: %v", err)
	}
	return rhino.NewNumber(v), nil
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
	// require non-ident char after keyword
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

func (p *Parser) skipWhitespaceOnly() {
	for !p.atEnd() && unicode.IsSpace(p.src[p.pos]) {
		p.advance()
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

func (p *Parser) errorf(f string, a ...any) error {
	return fmt.Errorf("parse error at line %d col %d: %s", p.line, p.col(), fmt.Sprintf(f, a...))
}

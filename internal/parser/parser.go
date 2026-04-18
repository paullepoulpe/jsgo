// Port of com/google/javascript/jscomp/parsing/parser/Parser.java
// @ closure-compiler v20260415 (narrow; not a faithful grammar port).
//
// The Java class is ~10k lines of a full ES parser built on a
// Scanner + token stream (parseSourceElement / parseStatement /
// parseExpression / parseAssignmentExpression /
// parsePrimaryExpression / …). This file is NOT a port of that
// structure — it's a hand-rolled statement-loop scanner that handles
// the grammar jsgo's checker currently uses (var/let/const +
// class/function decl, call/new/assign/return statements,
// arrow/function expressions, basic literals). Parse never panics
// (guardrail #6); unrecognised constructs skip to the next statement
// boundary. A real Parser.java port is a separate multi-session
// undertaking and will land at internal/jscomp/parsing/parser/
// alongside the JsDocInfoParser port.

package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/paullepoulpe/jsgo/internal/jscomp/parsing"
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
		jsdoc = parsing.ParseJSDocInfo(body)
		p.skipTrivia()
	}

	if p.consumeKeyword("function") {
		return p.finishFunction(jsdoc)
	}
	if p.consumeKeyword("class") {
		return p.finishClass(jsdoc)
	}
	if declToken, ok := p.consumeDeclKeyword(); ok {
		return p.finishDeclStatement(declToken, jsdoc)
	}
	if p.consumeKeyword("return") {
		return p.finishReturn(), true
	}

	// JSDoc may also annotate an expression statement (most commonly
	// `/** @type {T} */ this.x = …;`). Try the same statement forms,
	// attaching the comment to the inner ASSIGN/CALL/NEW so the scope
	// creator can pick it up.
	if stmt, ok := p.tryCallOrAssignStatement(); ok {
		if jsdoc != nil {
			attachJSDocToInner(stmt, jsdoc)
		}
		return stmt, true
	}
	if stmt, ok := p.tryExpressionStatement(); ok {
		if jsdoc != nil {
			attachJSDocToInner(stmt, jsdoc)
		}
		return stmt, true
	}
	if jsdoc != nil {
		p.skipToStatementBoundary()
		return nil, false
	}
	return nil, false
}

// attachJSDocToInner places the JSDocInfo on the immediately wrapped
// ASSIGN/CALL/NEW, matching how Closure attaches JSDoc to the operative
// node rather than the EXPR_RESULT wrapper.
func attachJSDocToInner(stmt *rhino.Node, info *rhino.JSDocInfo) {
	if stmt.Token() == rhino.EXPR_RESULT {
		if inner := stmt.FirstChild(); inner != nil {
			inner.SetJSDocInfo(info)
			return
		}
	}
	stmt.SetJSDocInfo(info)
}

// tryExpressionStatement wraps any expression followed by `;` into an
// EXPR_RESULT. Covers forms the IDENT-first speculative
// tryCallOrAssignStatement can't see — notably `new Foo(args);`.
func (p *Parser) tryExpressionStatement() (*rhino.Node, bool) {
	save := p.pos
	saveLine := p.line
	saveBOL := p.lineBOL
	rewind := func() {
		p.pos = save
		p.line = saveLine
		p.lineBOL = saveBOL
	}
	exprLine, exprCol := p.line, p.col()
	expr, ok := p.tryExpression()
	if !ok {
		rewind()
		return nil, false
	}
	p.skipTrivia()
	if !p.consumeByte(';') {
		rewind()
		return nil, false
	}
	stmt := rhino.New(rhino.EXPR_RESULT)
	stmt.SetPos(exprLine, exprCol)
	stmt.AddChildToBack(expr)
	return stmt, true
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

// finishClass parses `class NAME [extends BASE] { body }` after the
// `class` keyword has been consumed. Produces a CLASS node whose
// children are:
//
//   [0] NAME — class name (empty string when anonymous)
//   [1] NAME — super name (empty string when there's no `extends`)
//   [2] BLOCK — class body; each child is either a FUNCTION
//       (constructor or method) or a CLASS_FIELD_DEF.
//
// Mirrors Parser.java's parseClassDeclaration / parseClassBody output
// shape closely enough for TypedScopeCreator to process classes with
// the same two-pass pre-registration it already uses for @constructor
// functions.
func (p *Parser) finishClass(jsdoc *rhino.JSDocInfo) (*rhino.Node, bool) {
	class := rhino.New(rhino.CLASS)
	if jsdoc != nil {
		class.SetJSDocInfo(jsdoc)
	}

	p.skipTrivia()
	nameLine, nameCol := p.line, p.col()
	name, _ := p.consumeIdent() // anonymous class allowed
	nameNode := rhino.NewString(rhino.NAME, name)
	nameNode.SetPos(nameLine, nameCol)
	class.AddChildToBack(nameNode)

	p.skipTrivia()
	superNode := rhino.NewString(rhino.NAME, "")
	if p.consumeKeyword("extends") {
		p.skipTrivia()
		sLine, sCol := p.line, p.col()
		sname, ok := p.consumeIdent()
		if ok {
			superNode = rhino.NewString(rhino.NAME, sname)
			superNode.SetPos(sLine, sCol)
		}
		p.skipTrivia()
	}
	class.AddChildToBack(superNode)

	body := rhino.New(rhino.BLOCK)
	if !p.consumeByte('{') {
		p.skipToStatementBoundary()
		return class, true
	}
	for {
		p.skipTrivia()
		if p.atEnd() {
			break
		}
		if p.src[p.pos] == '}' {
			p.advance()
			break
		}
		p.parseClassMember(body)
	}
	class.AddChildToBack(body)
	return class, true
}

// parseClassMember parses a single entry inside a class body: an
// optional JSDoc comment, then either a method (IDENT(params) {body})
// or a field (IDENT [= expr];). Unknown shapes advance to the next
// boundary so the whole body doesn't get lost.
func (p *Parser) parseClassMember(body *rhino.Node) {
	var jsdoc *rhino.JSDocInfo
	if b, ok := p.consumeBlockComment(); ok {
		jsdoc = parsing.ParseJSDocInfo(b)
		p.skipTrivia()
	}
	if p.atEnd() || p.src[p.pos] == '}' {
		return
	}
	// `static` toggles the JSDoc's isStatic flag so the scope creator
	// lands the member on the class's static slot rather than the
	// instance prototype. Synthesises a JSDocInfo when none was
	// present on the source.
	isStatic := false
	if p.consumeKeyword("static") {
		isStatic = true
		p.skipTrivia()
	}
	if isStatic {
		if jsdoc == nil {
			jsdoc = &rhino.JSDocInfo{}
		}
		jsdoc.SetStatic(true)
	}
	memberLine, memberCol := p.line, p.col()
	name, ok := p.consumeIdent()
	if !ok {
		p.skipToStatementBoundary()
		return
	}
	p.skipTrivia()
	switch {
	case !p.atEnd() && p.src[p.pos] == '(':
		// Method: `name(params) { body }`.
		p.advance()
		fn := rhino.New(rhino.FUNCTION)
		if jsdoc != nil {
			fn.SetJSDocInfo(jsdoc)
		}
		fn.SetPos(memberLine, memberCol)
		nameNode := rhino.NewString(rhino.NAME, name)
		nameNode.SetPos(memberLine, memberCol)
		fn.AddChildToBack(nameNode)
		params := rhino.New(rhino.PARAM_LIST)
		for {
			p.skipTrivia()
			if p.atEnd() {
				return
			}
			if p.src[p.pos] == ')' {
				p.advance()
				break
			}
			pLine, pCol := p.line, p.col()
			pname, ok := p.consumeIdent()
			if !ok {
				p.skipToStatementBoundary()
				return
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
		methodBody := rhino.New(rhino.BLOCK)
		if p.consumeByte('{') {
			p.parseStatements(methodBody, '}')
			p.consumeByte('}')
		}
		fn.AddChildToBack(methodBody)
		body.AddChildToBack(fn)
	default:
		// Field: `name [= expr];` — captures both typed-only declarations
		// and typed-with-initializer class properties.
		field := rhino.New(rhino.CLASS_FIELD_DEF)
		if jsdoc != nil {
			field.SetJSDocInfo(jsdoc)
		}
		field.SetPos(memberLine, memberCol)
		nameNode := rhino.NewString(rhino.NAME, name)
		nameNode.SetPos(memberLine, memberCol)
		field.AddChildToBack(nameNode)
		if p.consumeByte('=') {
			p.skipTrivia()
			initLine, initCol := p.line, p.col()
			init, ok := p.tryExpression()
			if ok {
				init.SetPos(initLine, initCol)
				field.AddChildToBack(init)
			}
		}
		p.skipTrivia()
		p.consumeByte(';')
		body.AddChildToBack(field)
	}
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

// tryCallOrAssignStatement parses one of:
//
//   <lhs> = <expr>;    // ASSIGN, where lhs is NAME or GETPROP chain
//   <lhs>(args);       // CALL on a NAME/GETPROP chain
//   this.x = <expr>;   // bare-this assign — same shape as the GETPROP case
//
// Speculative: rewinds on mismatch so the caller can fall back to other
// statement forms.
func (p *Parser) tryCallOrAssignStatement() (*rhino.Node, bool) {
	save := p.pos
	saveLine := p.line
	saveBOL := p.lineBOL
	rewind := func() {
		p.pos = save
		p.line = saveLine
		p.lineBOL = saveBOL
	}

	startLine, startCol := p.line, p.col()
	lhs, ok := p.tryLHSExpression()
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
		// Call-as-statement on the whole LHS (handles foo.bar.baz(...)).
		p.advance()
		call := rhino.New(rhino.CALL)
		call.SetPos(startLine, startCol)
		call.AddChildToBack(lhs)
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
		p.consumeByte(';')
		stmt := rhino.New(rhino.EXPR_RESULT)
		stmt.SetPos(startLine, startCol)
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
		rhs.SetPos(rhsLine, rhsCol)
		assign := rhino.New(rhino.ASSIGN)
		assign.SetPos(startLine, startCol)
		assign.AddChildToBack(lhs)
		assign.AddChildToBack(rhs)
		stmt := rhino.New(rhino.EXPR_RESULT)
		stmt.SetPos(startLine, startCol)
		stmt.AddChildToBack(assign)
		return stmt, true
	}
	rewind()
	return nil, false
}

// tryLHSExpression parses a NAME or GETPROP chain starting at the
// current position. Used by tryCallOrAssignStatement to support
// `obj.x = …` and `this.x = …` shapes uniformly. Does NOT consume a
// trailing `(args)` — the caller decides whether the LHS is being
// called or assigned to.
func (p *Parser) tryLHSExpression() (*rhino.Node, bool) {
	if p.atEnd() {
		return nil, false
	}
	startLine, startCol := p.line, p.col()
	c := p.src[p.pos]
	if c == 't' && p.consumeKeyword("this") {
		base := rhino.New(rhino.THIS)
		base.SetPos(startLine, startCol)
		return p.extendWithProperty(base), true
	}
	if !(unicode.IsLetter(c) || c == '_' || c == '$') {
		return nil, false
	}
	ident, ok := p.consumeIdent()
	if !ok {
		return nil, false
	}
	base := rhino.NewString(rhino.NAME, ident)
	base.SetPos(startLine, startCol)
	return p.extendWithProperty(base), true
}

// tryExpression parses one expression: an arrow function, function
// expression, literal, identifier/property-access chain, call, or new.
// Anything else returns (nil, false).
func (p *Parser) tryExpression() (*rhino.Node, bool) {
	if p.atEnd() {
		return nil, false
	}
	if arrow, ok := p.tryArrow(); ok {
		return arrow, true
	}
	if p.consumeKeyword("function") {
		// Function expression — `function [name] (params) { body }`.
		// Shares the body-parse path with statement-level function
		// declarations; the caller never sees the difference.
		if fn, ok := p.finishFunction(nil); ok {
			return fn, true
		}
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
		if p.consumeKeyword("new") {
			return p.finishNew()
		}
		if p.consumeKeyword("null") {
			return rhino.New(rhino.NULL), true
		}
	case c == 'u':
		if p.consumeKeyword("undefined") {
			return rhino.New(rhino.UNDEFINED), true
		}
	case c == 't':
		if p.consumeKeyword("this") {
			return p.extendWithProperty(rhino.New(rhino.THIS)), true
		}
		if n, ok := p.parseBooleanLiteral(); ok {
			return n, true
		}
	}
	if unicode.IsLetter(c) || c == '_' || c == '$' {
		return p.parseNameOrCall()
	}
	return nil, false
}

// extendWithProperty consumes any trailing `.IDENT` chain after the
// given base, building GETPROP(base, STRINGLIT) nodes as in Closure.
// Used by the `this` and NAME paths so member access is uniform.
func (p *Parser) extendWithProperty(base *rhino.Node) *rhino.Node {
	for {
		p.skipTrivia()
		if p.atEnd() || p.src[p.pos] != '.' {
			return base
		}
		p.advance() // '.'
		p.skipTrivia()
		propLine, propCol := p.line, p.col()
		ident, ok := p.consumeIdent()
		if !ok {
			return base
		}
		gp := rhino.New(rhino.GETPROP)
		gp.SetPos(base.Lineno(), base.Charno())
		gp.AddChildToBack(base)
		propNode := rhino.NewString(rhino.STRINGLIT, ident)
		propNode.SetPos(propLine, propCol)
		gp.AddChildToBack(propNode)
		base = gp
	}
}

// tryArrow speculatively parses `(params) => body` or `ident => body`.
// Rewinds on failure so the caller can fall back to other expression
// forms. Arrow functions surface as FUNCTION nodes so the same
// scope/typecheck machinery applies to them.
func (p *Parser) tryArrow() (*rhino.Node, bool) {
	save := p.pos
	saveLine := p.line
	saveBOL := p.lineBOL
	rewind := func() {
		p.pos = save
		p.line = saveLine
		p.lineBOL = saveBOL
	}

	startLine, startCol := p.line, p.col()
	var params []string

	if p.consumeByte('(') {
		p.skipTrivia()
		if !p.consumeByte(')') {
			for {
				ident, ok := p.consumeIdent()
				if !ok {
					rewind()
					return nil, false
				}
				params = append(params, ident)
				p.skipTrivia()
				if p.consumeByte(',') {
					p.skipTrivia()
					continue
				}
				if p.consumeByte(')') {
					break
				}
				rewind()
				return nil, false
			}
		}
	} else {
		c := p.src[p.pos]
		if !(unicode.IsLetter(c) || c == '_' || c == '$') {
			return nil, false
		}
		ident, ok := p.consumeIdent()
		if !ok {
			rewind()
			return nil, false
		}
		params = []string{ident}
	}

	p.skipTrivia()
	if !p.consumeByte('=') || !p.consumeByte('>') {
		rewind()
		return nil, false
	}
	p.skipTrivia()

	fn := rhino.New(rhino.FUNCTION)
	fn.SetPos(startLine, startCol)
	nameNode := rhino.NewString(rhino.NAME, "")
	nameNode.SetPos(startLine, startCol)
	fn.AddChildToBack(nameNode)

	paramList := rhino.New(rhino.PARAM_LIST)
	for _, pn := range params {
		pNode := rhino.NewString(rhino.NAME, pn)
		paramList.AddChildToBack(pNode)
	}
	fn.AddChildToBack(paramList)

	body := rhino.New(rhino.BLOCK)
	if p.consumeByte('{') {
		p.parseStatements(body, '}')
		p.consumeByte('}')
	} else {
		// Single-expression body: wrap as implicit `return EXPR`.
		eline, ecol := p.line, p.col()
		expr, ok := p.tryExpression()
		if !ok {
			rewind()
			return nil, false
		}
		expr.SetPos(eline, ecol)
		ret := rhino.New(rhino.RETURN)
		ret.AddChildToBack(expr)
		body.AddChildToBack(ret)
	}
	fn.AddChildToBack(body)
	return fn, true
}

// finishNew parses `IDENT ( args )` after the `new` keyword is consumed.
// Mirrors com.google.javascript.jscomp.parsing.parser.Parser's
// parseMemberExpressionNoNew path that produces a NEW node.
func (p *Parser) finishNew() (*rhino.Node, bool) {
	p.skipTrivia()
	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		return nil, false
	}
	p.skipTrivia()
	if !p.consumeByte('(') {
		// `new Foo` without args isn't supported yet; treat as a plain
		// name reference so callers aren't blocked.
		name := rhino.NewString(rhino.NAME, ident)
		name.SetPos(nameLine, nameCol)
		return name, true
	}
	n := rhino.New(rhino.NEW)
	n.SetPos(nameLine, nameCol)
	callee := rhino.NewString(rhino.NAME, ident)
	callee.SetPos(nameLine, nameCol)
	n.AddChildToBack(callee)
	for {
		p.skipTrivia()
		if p.atEnd() {
			return nil, false
		}
		if p.src[p.pos] == ')' {
			p.advance()
			return n, true
		}
		argLine, argCol := p.line, p.col()
		arg, ok := p.tryExpression()
		if !ok {
			return nil, false
		}
		arg.SetPos(argLine, argCol)
		n.AddChildToBack(arg)
		p.skipTrivia()
		if p.consumeByte(',') {
			continue
		}
	}
}

// parseNameOrCall parses `IDENT [ .IDENT ]* [ ( args ) ]`. The chain of
// `.IDENT` becomes GETPROP nodes; a trailing `(...)` becomes CALL on
// whatever expression came before. Mirrors Parser.java's
// parseLeftHandSideExpressionAllowCall.
func (p *Parser) parseNameOrCall() (*rhino.Node, bool) {
	nameLine, nameCol := p.line, p.col()
	ident, ok := p.consumeIdent()
	if !ok {
		return nil, false
	}
	expr := rhino.NewString(rhino.NAME, ident)
	expr.SetPos(nameLine, nameCol)
	expr = p.extendWithProperty(expr)
	p.skipTrivia()
	if p.atEnd() || p.src[p.pos] != '(' {
		return expr, true
	}
	p.advance()
	call := rhino.New(rhino.CALL)
	call.SetPos(nameLine, nameCol)
	call.AddChildToBack(expr)
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

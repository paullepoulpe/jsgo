// Port of com/google/javascript/jscomp/parsing/parser/Parser.java
// @ closure-compiler v20260415 (narrow grammar; structure tracks Java).
//
// The Java class is ~10k lines of a full ES parser. This file carries
// the recursive-descent shape — a Parser that owns a Scanner, with
// speculative parse paths implemented via Scanner.Save/Restore — but
// handles only the grammar jsgo's checker currently consumes
// (var/let/const + class/function decl, call/new/assign/return
// statements, arrow/function expressions, basic literals). The
// character-level primitives live in scanner.go, matching Java's
// Scanner.java → Parser.java split. Token-stream-driven parsing
// (proper TokenType enum) lands when we port the grammar further —
// for now the Parser still peeks at runes via p.scanner.Peek().

package parser

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/paullepoulpe/jsgo/internal/jscomp/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Parser wraps a Scanner and produces a SCRIPT-rooted rhino.Node.
// Structure mirrors Java's Parser class: a Parser field of type
// Scanner (com.google.javascript.jscomp.parsing.parser.Scanner),
// recursive-descent methods on the parse path, and speculative
// branches implemented via Scanner.Save / Restore.
type Parser struct {
	scanner *Scanner
	// pendingStatements is a FIFO of nodes that should be consumed by
	// parseStatements before parsing fresh input. Used by
	// swallowForHead to surface the synthetic VAR decl for a
	// `for (let x of …)` binding as a real scope-building child.
	pendingStatements []*rhino.Node
}

func Parse(src string) (n *rhino.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse panic recovered: %v", r)
		}
	}()
	p := &Parser{scanner: NewScanner(src)}
	script := rhino.New(rhino.SCRIPT)
	script.SetPos(1, 0)
	p.parseStatements(script, 0)
	return script, nil
}

// parseStatements fills parent with top-level / block statements until
// EOF (if endTok == 0) or until a '}' is hit at the current depth.
// Always makes forward progress: if neither tryStatement nor
// SkipToStatementBoundary advances, the loop consumes one rune as a
// safety valve to guarantee termination on malformed input (matches
// the recovery intent of Java's Parser's error-resync path).
func (p *Parser) parseStatements(parent *rhino.Node, endTok rune) {
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return
		}
		if endTok != 0 && p.scanner.Peek() == endTok {
			return
		}
		start := p.scanner.Pos()
		stmt, ok := p.tryStatement()
		if ok && stmt != nil {
			parent.AddChildToBack(stmt)
			continue
		}
		if p.scanner.Pos() == start {
			p.scanner.SkipToStatementBoundary()
		}
		if p.scanner.Pos() == start {
			// Hard-stuck (e.g. a bare `}` at depth zero). Advance one
			// rune so the loop always makes progress.
			p.scanner.Advance()
		}
	}
}

// mergeSuppressions copies every @suppress name from `src` into
// `dst`. Used to propagate stacked-JSDoc suppressions onto the
// statement that ultimately owns them, even when another block
// in the stack "wins" the primary attachment.
func mergeSuppressions(dst, src *rhino.JSDocInfo) {
	for _, name := range []string{"missingProperties", "checkTypes", "duplicate", "undefinedVars", "missingReturn"} {
		if src.HasSuppression(name) && !dst.HasSuppression(name) {
			dst.AddSuppression(name)
		}
	}
}

// jsdocCarriesMoreInfo reports whether `b` should displace `a` as the
// JSDoc attached to the next statement. Priority: declaration-shaped
// tags (@interface, @constructor, @record, @extends) beat plain
// type / param tags. Matches Closure's "the last declaration comment
// wins" heuristic when multiple JSDoc blocks precede a single decl.
func jsdocCarriesMoreInfo(b, a *rhino.JSDocInfo) bool {
	decl := func(i *rhino.JSDocInfo) bool {
		return i.IsConstructor() || i.IsInterface() || i.IsRecord() || i.BaseType() != nil
	}
	if decl(b) && !decl(a) {
		return true
	}
	return false
}

// tryStatement attempts one statement. Returns (stmt, true) on success
// or (nil, false) otherwise; in either case p.scanner.Pos() is left past the
// consumed input (or rewound on speculative failures).
func (p *Parser) tryStatement() (*rhino.Node, bool) {
	var jsdoc *rhino.JSDocInfo
	// Consume any number of back-to-back block comments. When a
	// JSDoc carries declaration-shaped tags (@interface, @extends,
	// @constructor, @param, @return), that's the one that attaches
	// to the following statement — ignore earlier stray JSDocs
	// which Java silently drops in the same scenario.
	for {
		body, ok := p.scanner.ConsumeBlockComment()
		if !ok {
			break
		}
		parsed := parsing.ParseJSDocInfo(body)
		if parsed != nil {
			if jsdoc == nil {
				jsdoc = parsed
			} else if jsdocCarriesMoreInfo(parsed, jsdoc) {
				mergeSuppressions(parsed, jsdoc)
				jsdoc = parsed
			} else {
				// Losing block still contributes its @suppress
				// entries — @fileoverview comments typically carry
				// these and must not be dropped silently.
				mergeSuppressions(jsdoc, parsed)
			}
		}
		p.scanner.SkipTrivia()
	}

	// Capture the `function` / `class` keyword's position so the
	// resulting FUNCTION / CLASS node has a stable lineno/charno
	// matching Java (passes like SuperclassCheck report at the
	// declaration's start column).
	kwLine, kwCol := p.scanner.Line(), p.scanner.Col()
	if p.scanner.ConsumeKeyword("function") {
		fn, ok := p.finishFunction(jsdoc)
		if fn != nil {
			fn.SetPos(kwLine, kwCol)
		}
		return fn, ok
	}
	if p.scanner.ConsumeKeyword("class") {
		cls, ok := p.finishClass(jsdoc)
		if cls != nil {
			cls.SetPos(kwLine, kwCol)
		}
		return cls, ok
	}
	if declToken, ok := p.consumeDeclKeyword(); ok {
		return p.finishDeclStatement(declToken, jsdoc)
	}
	if p.scanner.ConsumeKeyword("return") {
		return p.finishReturn(), true
	}
	if p.scanner.ConsumeKeyword("if") {
		return p.finishIf(), true
	}
	if p.scanner.ConsumeKeyword("while") {
		return p.finishWhile(), true
	}
	if p.scanner.ConsumeKeyword("do") {
		return p.finishDoWhile(), true
	}
	if p.scanner.ConsumeKeyword("for") {
		return p.finishFor(), true
	}
	if p.scanner.ConsumeKeyword("switch") {
		return p.finishSwitch(), true
	}
	if p.scanner.ConsumeKeyword("try") {
		return p.finishTry(), true
	}
	if p.scanner.ConsumeKeyword("throw") {
		return p.finishThrow(), true
	}
	if p.scanner.ConsumeKeyword("break") {
		return p.finishBreakContinue(rhino.BREAK), true
	}
	if p.scanner.ConsumeKeyword("continue") {
		return p.finishBreakContinue(rhino.CONTINUE), true
	}
	if p.scanner.ConsumeKeyword("debugger") {
		return p.finishDebugger(kwLine, kwCol), true
	}
	// Labeled statement: `IDENT : STMT`. Speculative — rewinds if the
	// cursor doesn't end up at a label colon.
	if stmt, ok := p.tryLabeledStatement(jsdoc); ok {
		return stmt, true
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
		p.scanner.SkipToStatementBoundary()
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
	save := p.scanner.Save()
	rewind := func() {
		p.scanner.Restore(save)
	}
	exprLine, exprCol := p.scanner.Line(), p.scanner.Col()
	expr, ok := p.tryExpression()
	if !ok {
		rewind()
		return nil, false
	}
	p.scanner.SkipTrivia()
	// Normal form: `expr;`. Java's Parser also handles ASI — a newline
	// or EOF / `}` closes the statement without a literal semicolon.
	// We accept the same lax terminators so single-line fixtures like
	// `'foo' in Object` surface as statements.
	if !p.scanner.ConsumeByte(';') {
		if !p.scanner.AtEnd() && p.scanner.Peek() != '}' {
			rewind()
			return nil, false
		}
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

	p.scanner.SkipTrivia()
	// Inline JSDoc between the decl keyword and the name —
	// `var /** null */ n = …;` — treat the brace-only type
	// expression as an @type annotation on the resulting VAR.
	// Matches how Java's parser synthesises a @type JSDoc from
	// an inline bracketed comment at this position.
	if body, ok := p.scanner.ConsumeBlockComment(); ok {
		if inline := parseInlineTypeJSDoc(body); inline != nil && declNode.JSDocInfo() == nil {
			declNode.SetJSDocInfo(inline)
		}
		p.scanner.SkipTrivia()
	}
	// One iteration per comma-separated binding. Matches Java's
	// parseVariableDeclarationList loop.
	for {
		nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
		ident, ok := p.scanner.ConsumeIdent()
		if !ok {
			p.scanner.SkipToStatementBoundary()
			return nil, false
		}
		nameNode := rhino.NewString(rhino.NAME, ident)
		nameNode.SetPos(nameLine, nameCol)
		declNode.AddChildToBack(nameNode)

		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte('=') {
			if !p.scanner.AtEnd() && p.scanner.Peek() == '=' {
				p.scanner.SkipToStatementBoundary()
				return declNode, true
			}
			p.scanner.SkipTrivia()
			initLine, initCol := p.scanner.Line(), p.scanner.Col()
			init, ok := p.tryExpression()
			if !ok {
				p.scanner.SkipToStatementBoundary()
				return declNode, true
			}
			init.SetPos(initLine, initCol)
			nameNode.AddChildToBack(init)
		}
		p.scanner.SkipTrivia()
		if !p.scanner.ConsumeByte(',') {
			break
		}
		p.scanner.SkipTrivia()
		// Allow an inline JSDoc per-binding (`var a, /** @type T */ b`).
		if body, ok := p.scanner.ConsumeBlockComment(); ok {
			if inline := parseInlineTypeJSDoc(body); inline != nil {
				_ = inline // discarded — Java drops per-binding type hints in this pos
			}
			p.scanner.SkipTrivia()
		}
	}
	sawNewline := p.scanner.Line() != p.scanner.LastTokenLine()
	if !p.scanner.ConsumeByte(';') {
		// ASI (ES5 §7.9): if we're at EOF / `}` / saw a newline
		// before the next token, treat the VAR as terminated.
		// Otherwise we're mid-statement and recovery is needed.
		if sawNewline || p.scanner.AtEnd() || p.scanner.Peek() == '}' {
			return declNode, true
		}
		p.scanner.SkipToStatementBoundary()
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

	p.scanner.SkipTrivia()
	nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
	name, ok := p.scanner.ConsumeIdent()
	nameNode := rhino.NewString(rhino.NAME, "")
	if ok {
		nameNode = rhino.NewString(rhino.NAME, name)
	}
	nameNode.SetPos(nameLine, nameCol)
	fn.AddChildToBack(nameNode)

	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		p.scanner.SkipToStatementBoundary()
		return fn, true
	}
	params := rhino.New(rhino.PARAM_LIST)
	for {
		p.scanner.SkipTrivia()
		// Inline JSDoc on parameters: `function f(/** T */ x)` — parse
		// the brace body as an @type annotation and attach to the
		// resulting param NAME node. Matches Java's Parser which stores
		// the inline type on the param's JSDocInfo so the scope builder
		// can pick it up even without an outer @param.
		var inlineInfo *rhino.JSDocInfo
		if body, ok := p.scanner.ConsumeBlockComment(); ok {
			if parsed := parseInlineTypeJSDoc(body); parsed != nil {
				inlineInfo = parsed
			}
			p.scanner.SkipTrivia()
		}
		if p.scanner.AtEnd() {
			return fn, true
		}
		if p.scanner.Peek() == ')' {
			p.scanner.Advance()
			break
		}
		pLine, pCol := p.scanner.Line(), p.scanner.Col()
		pname, ok := p.scanner.ConsumeIdent()
		if !ok {
			p.scanner.SkipToStatementBoundary()
			return fn, true
		}
		pNode := rhino.NewString(rhino.NAME, pname)
		pNode.SetPos(pLine, pCol)
		if inlineInfo != nil {
			pNode.SetJSDocInfo(inlineInfo)
		}
		params.AddChildToBack(pNode)
		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte(',') {
			continue
		}
	}
	fn.AddChildToBack(params)

	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('{') {
		p.scanner.SkipToStatementBoundary()
		return fn, true
	}
	block := rhino.New(rhino.BLOCK)
	p.parseStatements(block, '}')
	p.scanner.ConsumeByte('}')
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

	p.scanner.SkipTrivia()
	nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
	name, _ := p.scanner.ConsumeIdent() // anonymous class allowed
	nameNode := rhino.NewString(rhino.NAME, name)
	nameNode.SetPos(nameLine, nameCol)
	class.AddChildToBack(nameNode)

	p.scanner.SkipTrivia()
	superNode := rhino.NewString(rhino.NAME, "")
	if p.scanner.ConsumeKeyword("extends") {
		p.scanner.SkipTrivia()
		sLine, sCol := p.scanner.Line(), p.scanner.Col()
		sname, ok := p.scanner.ConsumeIdent()
		if ok {
			superNode = rhino.NewString(rhino.NAME, sname)
			superNode.SetPos(sLine, sCol)
		}
		p.scanner.SkipTrivia()
	}
	class.AddChildToBack(superNode)

	body := rhino.New(rhino.BLOCK)
	if !p.scanner.ConsumeByte('{') {
		p.scanner.SkipToStatementBoundary()
		return class, true
	}
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			break
		}
		if p.scanner.Peek() == '}' {
			p.scanner.Advance()
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
	if b, ok := p.scanner.ConsumeBlockComment(); ok {
		jsdoc = parsing.ParseJSDocInfo(b)
		p.scanner.SkipTrivia()
	}
	if p.scanner.AtEnd() || p.scanner.Peek() == '}' {
		return
	}
	// `static` toggles the JSDoc's isStatic flag so the scope creator
	// lands the member on the class's static slot rather than the
	// instance prototype. Synthesises a JSDocInfo when none was
	// present on the source. The member's start position is captured
	// BEFORE consuming `static` so downstream diagnostics point at the
	// whole member declaration (Java's Node.getLineno() for a
	// CLASS_FIELD_DEF is the keyword / name column, inclusive of
	// `static`).
	memberLine, memberCol := p.scanner.Line(), p.scanner.Col()
	isStatic := false
	if p.scanner.ConsumeKeyword("static") {
		isStatic = true
		p.scanner.SkipTrivia()
	}
	if isStatic {
		if jsdoc == nil {
			jsdoc = &rhino.JSDocInfo{}
		}
		jsdoc.SetStatic(true)
	}
	nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
	name, ok := p.scanner.ConsumeIdent()
	if !ok {
		p.scanner.SkipToStatementBoundary()
		return
	}
	p.scanner.SkipTrivia()
	switch {
	case !p.scanner.AtEnd() && p.scanner.Peek() == '(':
		// Method: `name(params) { body }`.
		p.scanner.Advance()
		fn := rhino.New(rhino.FUNCTION)
		if jsdoc != nil {
			fn.SetJSDocInfo(jsdoc)
		}
		fn.SetPos(memberLine, memberCol)
		nameNode := rhino.NewString(rhino.NAME, name)
		nameNode.SetPos(nameLine, nameCol)
		fn.AddChildToBack(nameNode)
		params := rhino.New(rhino.PARAM_LIST)
		for {
			p.scanner.SkipTrivia()
			if p.scanner.AtEnd() {
				return
			}
			if p.scanner.Peek() == ')' {
				p.scanner.Advance()
				break
			}
			pLine, pCol := p.scanner.Line(), p.scanner.Col()
			pname, ok := p.scanner.ConsumeIdent()
			if !ok {
				p.scanner.SkipToStatementBoundary()
				return
			}
			pNode := rhino.NewString(rhino.NAME, pname)
			pNode.SetPos(pLine, pCol)
			params.AddChildToBack(pNode)
			p.scanner.SkipTrivia()
			if p.scanner.ConsumeByte(',') {
				continue
			}
		}
		fn.AddChildToBack(params)
		p.scanner.SkipTrivia()
		methodBody := rhino.New(rhino.BLOCK)
		if p.scanner.ConsumeByte('{') {
			p.parseStatements(methodBody, '}')
			p.scanner.ConsumeByte('}')
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
		nameNode.SetPos(nameLine, nameCol)
		field.AddChildToBack(nameNode)
		if p.scanner.ConsumeByte('=') {
			p.scanner.SkipTrivia()
			initLine, initCol := p.scanner.Line(), p.scanner.Col()
			init, ok := p.tryExpression()
			if ok {
				init.SetPos(initLine, initCol)
				field.AddChildToBack(init)
			}
		}
		p.scanner.SkipTrivia()
		p.scanner.ConsumeByte(';')
		body.AddChildToBack(field)
	}
}

// finishIf parses `if ( cond ) thenStmt [else elseStmt]`. Builds an
// IF node whose children are (cond, thenBlock, elseBlock?). thenStmt
// and elseStmt may be either a brace-delimited block or a single
// statement. Mirrors Parser.java's parseIfStatement.
func (p *Parser) finishIf() *rhino.Node {
	ifNode := rhino.New(rhino.IF)
	ifNode.SetPos(p.scanner.Line(), p.scanner.Col())
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		p.scanner.SkipToStatementBoundary()
		return ifNode
	}
	p.scanner.SkipTrivia()
	cond, ok := p.tryExpression()
	if !ok {
		p.scanner.SkipToStatementBoundary()
		return ifNode
	}
	ifNode.AddChildToBack(cond)
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte(')') {
		p.scanner.SkipToStatementBoundary()
		return ifNode
	}
	ifNode.AddChildToBack(p.parseControlBody())
	p.scanner.SkipTrivia()
	if p.scanner.ConsumeKeyword("else") {
		ifNode.AddChildToBack(p.parseControlBody())
	}
	return ifNode
}

// parseControlBody reads a block `{ … }` or a single statement into a
// BLOCK node. Helpful both for `if` branches and for the loop/try
// coarse handlers.
//
// When there are no explicit braces, the returned BLOCK has
// PropAddedBlock set (mirrors Java Node.IS_ADDED_BLOCK) so that
// CheckSuspiciousCode can detect misplaced semicolons like `if (x);`.
func (p *Parser) parseControlBody() *rhino.Node {
	p.scanner.SkipTrivia()
	body := rhino.New(rhino.BLOCK)
	if p.scanner.ConsumeByte('{') {
		p.parseStatements(body, '}')
		p.scanner.ConsumeByte('}')
		return body
	}
	// Synthesized block — mark it so CheckSuspiciousCode can detect
	// `if (x);` (PropAddedBlock + no children = misplaced semicolon).
	body.PutBooleanProp(rhino.PropAddedBlock, true)
	if stmt, ok := p.tryStatement(); ok && stmt != nil {
		body.AddChildToBack(stmt)
	}
	return body
}

// finishWhile parses `while (cond) stmt` into WHILE(cond, body).
func (p *Parser) finishWhile() *rhino.Node {
	n := rhino.New(rhino.WHILE)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	p.scanner.SkipTrivia()
	cond := p.consumeParenExpression()
	if cond == nil {
		cond = rhino.New(rhino.TRUE)
	}
	n.AddChildToBack(cond)
	n.AddChildToBack(p.parseControlBody())
	return n
}

// finishDoWhile parses `do stmt while (cond);` into
// DO_WHILE(body, cond).
func (p *Parser) finishDoWhile() *rhino.Node {
	n := rhino.New(rhino.DO_WHILE)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	body := p.parseControlBody()
	p.scanner.SkipTrivia()
	p.scanner.ConsumeKeyword("while")
	cond := p.consumeParenExpression()
	if cond == nil {
		cond = rhino.New(rhino.TRUE)
	}
	n.AddChildToBack(body)
	n.AddChildToBack(cond)
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
	return n
}

// finishFor parses `for (init; cond; upd) stmt` / `for (X of Y) stmt`
// / `for (X in Y) stmt`. Returns FOR / FOR_OF / FOR_IN with the
// appropriate child layout.
//
// FOR children: (initOrEmpty, condOrEmpty, updOrEmpty, body)
// FOR_IN / FOR_OF: (declOrLhs, iterableExpr, body)
func (p *Parser) finishFor() *rhino.Node {
	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		n := rhino.New(rhino.FOR)
		n.SetPos(startLine, startCol)
		return n
	}
	p.scanner.SkipTrivia()
	// Parse the init piece — either a var/let/const binding or an
	// expression or empty.
	var init *rhino.Node
	if declTok, ok := p.consumeDeclKeyword(); ok {
		p.scanner.SkipTrivia()
		nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
		ident, ok := p.scanner.ConsumeIdent()
		if ok {
			nameNode := rhino.NewString(rhino.NAME, ident)
			nameNode.SetPos(nameLine, nameCol)
			decl := rhino.New(declTok)
			decl.SetPos(nameLine, nameCol)
			decl.AddChildToBack(nameNode)
			p.scanner.SkipTrivia()
			// Optional `= expr` initializer.
			if p.scanner.ConsumeByte('=') {
				p.scanner.SkipTrivia()
				if expr, ok := p.tryExpression(); ok {
					nameNode.AddChildToBack(expr)
				}
			}
			init = decl
		}
	} else if p.scanner.Peek() != ';' {
		init, _ = p.tryExpression()
	}
	p.scanner.SkipTrivia()
	// Disambiguate for-in / for-of / classic for.
	if init != nil && p.scanner.AtKeyword("in") {
		p.scanner.ConsumeKeywordRaw("in")
		return p.finishForInOf(rhino.FOR_IN, init, startLine, startCol)
	}
	if init != nil && p.scanner.AtKeyword("of") {
		p.scanner.ConsumeKeywordRaw("of")
		return p.finishForInOf(rhino.FOR_OF, init, startLine, startCol)
	}
	// Classic `for (init; cond; upd) body`.
	p.scanner.ConsumeByte(';')
	p.scanner.SkipTrivia()
	var cond *rhino.Node
	if p.scanner.Peek() != ';' {
		cond, _ = p.tryExpression()
	}
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
	p.scanner.SkipTrivia()
	var upd *rhino.Node
	if p.scanner.Peek() != ')' {
		upd, _ = p.tryExpression()
	}
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(')')
	n := rhino.New(rhino.FOR)
	n.SetPos(startLine, startCol)
	if init == nil {
		init = rhino.New(rhino.BLOCK)
	}
	if cond == nil {
		cond = rhino.New(rhino.TRUE)
	}
	if upd == nil {
		upd = rhino.New(rhino.BLOCK)
	}
	n.AddChildToBack(init)
	n.AddChildToBack(cond)
	n.AddChildToBack(upd)
	n.AddChildToBack(p.parseControlBody())
	return n
}

// finishForInOf parses the `of EXPR) body` or `in EXPR) body` tail
// of a for-in/for-of loop after the keyword has been consumed.
func (p *Parser) finishForInOf(tok rhino.Token, lhs *rhino.Node, line, col int) *rhino.Node {
	p.scanner.SkipTrivia()
	rhs, _ := p.tryExpression()
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(')')
	n := rhino.New(tok)
	n.SetPos(line, col)
	n.AddChildToBack(lhs)
	if rhs == nil {
		rhs = rhino.New(rhino.UNDEFINED)
	}
	n.AddChildToBack(rhs)
	n.AddChildToBack(p.parseControlBody())
	return n
}

// finishSwitch parses `switch (e) { case v1: stmts... case v2: stmts...
// default: stmts... }` into
//   SWITCH(discriminant, CASE(val, BLOCK), CASE(val, BLOCK), DEFAULT_CASE(BLOCK))
func (p *Parser) finishSwitch() *rhino.Node {
	n := rhino.New(rhino.SWITCH)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	disc := p.consumeParenExpression()
	if disc == nil {
		disc = rhino.New(rhino.UNDEFINED)
	}
	n.AddChildToBack(disc)
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('{') {
		return n
	}
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() || p.scanner.Peek() == '}' {
			break
		}
		caseLine, caseCol := p.scanner.Line(), p.scanner.Col()
		if p.scanner.ConsumeKeyword("case") {
			p.scanner.SkipTrivia()
			val, _ := p.tryExpression()
			p.scanner.SkipTrivia()
			p.scanner.ConsumeByte(':')
			body := p.parseCaseBody()
			c := rhino.New(rhino.CASE)
			c.SetPos(caseLine, caseCol)
			if val == nil {
				val = rhino.New(rhino.UNDEFINED)
			}
			c.AddChildToBack(val)
			c.AddChildToBack(body)
			n.AddChildToBack(c)
			continue
		}
		if p.scanner.ConsumeKeyword("default") {
			p.scanner.SkipTrivia()
			p.scanner.ConsumeByte(':')
			body := p.parseCaseBody()
			d := rhino.New(rhino.DEFAULT_CASE)
			d.SetPos(caseLine, caseCol)
			d.AddChildToBack(body)
			n.AddChildToBack(d)
			continue
		}
		// Unrecognised tokens inside the switch — advance one rune
		// so we don't infinite-loop.
		p.scanner.Advance()
	}
	p.scanner.ConsumeByte('}')
	return n
}

// parseCaseBody reads statements up to the next `case` / `default` /
// `}`. Returned as a BLOCK so the CFG builder can chain them.
func (p *Parser) parseCaseBody() *rhino.Node {
	body := rhino.New(rhino.BLOCK)
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() || p.scanner.Peek() == '}' {
			return body
		}
		if p.scanner.AtKeyword("case") || p.scanner.AtKeyword("default") {
			return body
		}
		start := p.scanner.Pos()
		stmt, ok := p.tryStatement()
		if ok && stmt != nil {
			body.AddChildToBack(stmt)
			continue
		}
		if p.scanner.Pos() == start {
			p.scanner.Advance()
		}
	}
}

// finishTry parses `try { … } catch (e) { … } [finally { … }]` into
// TRY(body, CATCH(name, handler), FINALLY(block)). Catch clause is
// optional (modern JS allows `try { } finally { }` without catch).
func (p *Parser) finishTry() *rhino.Node {
	n := rhino.New(rhino.TRY)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	body := p.parseControlBody()
	n.AddChildToBack(body)
	p.scanner.SkipTrivia()
	if p.scanner.ConsumeKeyword("catch") {
		p.scanner.SkipTrivia()
		var param *rhino.Node
		if p.scanner.ConsumeByte('(') {
			p.scanner.SkipTrivia()
			_, _ = p.scanner.ConsumeBlockComment()
			p.scanner.SkipTrivia()
			nl, nc := p.scanner.Line(), p.scanner.Col()
			if ident, ok := p.scanner.ConsumeIdent(); ok {
				param = rhino.NewString(rhino.NAME, ident)
				param.SetPos(nl, nc)
			}
			p.scanner.SkipTrivia()
			p.scanner.ConsumeByte(')')
		}
		handler := p.parseControlBody()
		c := rhino.New(rhino.CATCH)
		c.SetPos(p.scanner.Line(), p.scanner.Col())
		if param == nil {
			param = rhino.NewString(rhino.NAME, "")
		}
		c.AddChildToBack(param)
		c.AddChildToBack(handler)
		n.AddChildToBack(c)
	}
	p.scanner.SkipTrivia()
	if p.scanner.ConsumeKeyword("finally") {
		finallyBlock := p.parseControlBody()
		f := rhino.New(rhino.FINALLY)
		f.AddChildToBack(finallyBlock)
		n.AddChildToBack(f)
	}
	return n
}

// finishThrow parses `throw expr;` → THROW(expr).
func (p *Parser) finishThrow() *rhino.Node {
	n := rhino.New(rhino.THROW)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	p.scanner.SkipTrivia()
	expr, _ := p.tryExpression()
	if expr == nil {
		expr = rhino.New(rhino.UNDEFINED)
	}
	n.AddChildToBack(expr)
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
	return n
}

// finishBreakContinue parses `break [label];` / `continue [label];`.
// Emits BREAK / CONTINUE with an optional NAME child carrying the
// label.
func (p *Parser) finishBreakContinue(tok rhino.Token) *rhino.Node {
	n := rhino.New(tok)
	n.SetPos(p.scanner.Line(), p.scanner.Col())
	p.scanner.SkipTrivia()
	if !p.scanner.AtEnd() && p.scanner.Peek() != ';' && p.scanner.Peek() != '}' && p.scanner.Peek() != '\n' {
		if ident, ok := p.scanner.ConsumeIdent(); ok {
			label := rhino.NewString(rhino.NAME, ident)
			n.AddChildToBack(label)
		}
	}
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
	return n
}

// finishDebugger parses `debugger;` → DEBUGGER node.
// Position is taken from the keyword (kwLine, kwCol) captured by the
// caller before ConsumeKeyword advanced the scanner.
func (p *Parser) finishDebugger(kwLine, kwCol int) *rhino.Node {
	n := rhino.New(rhino.DEBUGGER)
	n.SetPos(kwLine, kwCol)
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
	return n
}

// tryLabeledStatement speculatively parses `IDENT : STMT` into
// LABEL(NAME, inner). Rewinds when the cursor doesn't form a label
// (so the same position can fall through to the call/assign/expr
// statement parsers).
func (p *Parser) tryLabeledStatement(jsdoc *rhino.JSDocInfo) (*rhino.Node, bool) {
	save := p.scanner.Save()
	line, col := p.scanner.Line(), p.scanner.Col()
	ident, ok := p.scanner.ConsumeIdent()
	if !ok {
		p.scanner.Restore(save)
		return nil, false
	}
	p.scanner.SkipTrivia()
	if p.scanner.Peek() != ':' || p.scanner.PeekAt(1) == ':' || p.scanner.PeekAt(1) == '=' {
		p.scanner.Restore(save)
		return nil, false
	}
	p.scanner.Advance() // ':'
	inner, ok := p.tryStatement()
	if !ok {
		p.scanner.Restore(save)
		return nil, false
	}
	label := rhino.New(rhino.LABEL)
	label.SetPos(line, col)
	if jsdoc != nil {
		label.SetJSDocInfo(jsdoc)
	}
	nameNode := rhino.NewString(rhino.NAME, ident)
	nameNode.SetPos(line, col)
	label.AddChildToBack(nameNode)
	if inner != nil {
		label.AddChildToBack(inner)
	}
	return label, true
}

// consumeParenExpression reads `( expr )` at the cursor, returning
// the inner expression (or nil if malformed). Used by while / do /
// switch header parsers.
func (p *Parser) consumeParenExpression() *rhino.Node {
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		return nil
	}
	p.scanner.SkipTrivia()
	expr, ok := p.tryExpression()
	if !ok {
		// Recover: skip to the matching close.
		depth := 1
		for depth > 0 && !p.scanner.AtEnd() {
			c := p.scanner.Peek()
			switch c {
			case '(':
				depth++
			case ')':
				depth--
			}
			p.scanner.Advance()
		}
		return nil
	}
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(')')
	return expr
}

func (p *Parser) finishReturn() *rhino.Node {
	ret := rhino.New(rhino.RETURN)
	p.scanner.SkipTrivia()
	if !p.scanner.AtEnd() && p.scanner.Peek() != ';' {
		eline, ecol := p.scanner.Line(), p.scanner.Col()
		if expr, ok := p.tryExpression(); ok {
			expr.SetPos(eline, ecol)
			ret.AddChildToBack(expr)
		} else {
			p.scanner.SkipToStatementBoundary()
			return ret
		}
	}
	p.scanner.SkipTrivia()
	p.scanner.ConsumeByte(';')
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
	save := p.scanner.Save()
	rewind := func() {
		p.scanner.Restore(save)
	}

	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	lhs, ok := p.tryLHSExpression()
	if !ok {
		rewind()
		return nil, false
	}
	p.scanner.SkipTrivia()
	if p.scanner.AtEnd() {
		rewind()
		return nil, false
	}

	switch p.scanner.Peek() {
	case '(':
		// Call-as-statement on the whole LHS (handles foo.bar.baz(...)).
		p.scanner.Advance()
		call := rhino.New(rhino.CALL)
		call.SetPos(startLine, startCol)
		call.AddChildToBack(lhs)
		for {
			p.scanner.SkipTrivia()
			if p.scanner.AtEnd() {
				rewind()
				return nil, false
			}
			if p.scanner.Peek() == ')' {
				p.scanner.Advance()
				break
			}
			argLine, argCol := p.scanner.Line(), p.scanner.Col()
			arg, ok := p.tryExpression()
			if !ok {
				rewind()
				return nil, false
			}
			arg.SetPos(argLine, argCol)
			call.AddChildToBack(arg)
			p.scanner.SkipTrivia()
			if p.scanner.ConsumeByte(',') {
				continue
			}
		}
		p.scanner.SkipTrivia()
		// Call-as-statement must end with `;` (or `}` / EOF). If we
		// see a postfix continuation — `.prop` / `(args)` / `[i]` —
		// or a binary-op keyword (`in`, `instanceof`) the CALL is a
		// sub-expression; rewind so tryExpressionStatement picks up
		// the full form.
		if !p.scanner.AtEnd() {
			switch p.scanner.Peek() {
			case '.', '(', '[':
				rewind()
				return nil, false
			}
			if p.scanner.AtKeyword("in") || p.scanner.AtKeyword("instanceof") {
				rewind()
				return nil, false
			}
			// Binary-op character that isn't a statement terminator
			// (e.g. `+`, `-`, `*`, `<`, `==`) — same story.
			if tok, _, _ := p.peekBinOp(); tok != rhino.ERROR {
				rewind()
				return nil, false
			}
		}
		p.scanner.ConsumeByte(';')
		stmt := rhino.New(rhino.EXPR_RESULT)
		stmt.SetPos(startLine, startCol)
		stmt.AddChildToBack(call)
		return stmt, true

	case '=':
		p.scanner.Advance()
		if !p.scanner.AtEnd() && (p.scanner.Peek() == '=' || p.scanner.Peek() == '>') {
			rewind()
			return nil, false
		}
		p.scanner.SkipTrivia()
		rhsLine, rhsCol := p.scanner.Line(), p.scanner.Col()
		rhs, ok := p.tryExpression()
		if !ok {
			rewind()
			return nil, false
		}
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() || p.scanner.Peek() != ';' {
			rewind()
			return nil, false
		}
		p.scanner.Advance()
		if rhs.Lineno() == 0 && rhs.Charno() == 0 {
			rhs.SetPos(rhsLine, rhsCol)
		}
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
	if p.scanner.AtEnd() {
		return nil, false
	}
	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	c := p.scanner.Peek()
	if c == 't' && p.scanner.ConsumeKeyword("this") {
		base := rhino.New(rhino.THIS)
		base.SetPos(startLine, startCol)
		return p.extendWithProperty(base), true
	}
	if c == 's' && p.scanner.ConsumeKeyword("super") {
		base := rhino.New(rhino.SUPER)
		base.SetPos(startLine, startCol)
		return p.extendWithProperty(base), true
	}
	if !(unicode.IsLetter(c) || c == '_' || c == '$') {
		return nil, false
	}
	ident, ok := p.scanner.ConsumeIdent()
	if !ok {
		return nil, false
	}
	base := rhino.NewString(rhino.NAME, ident)
	base.SetPos(startLine, startCol)
	return p.extendWithProperty(base), true
}

// tryExpression is the top-level expression entry — parses the widest
// production (AssignmentExpression in Java). Composed of a precedence-
// climbing chain: tryAssignment → tryConditional → tryBinary →
// tryUnary → tryPostfix → tryPrimary. Mirrors Parser.java's shape
// (parseAssignmentExpression → parseConditional → parseBinary →
// parseUnary → parseLeftHandSideExpression → parsePrimaryExpression).
func (p *Parser) tryExpression() (*rhino.Node, bool) {
	return p.tryAssignment()
}

// tryAssignment parses `lhs [= += -= …] rhs`. Right-associative, so the
// RHS recursively descends back through tryAssignment.
func (p *Parser) tryAssignment() (*rhino.Node, bool) {
	if p.scanner.AtEnd() {
		return nil, false
	}
	lhs, ok := p.tryConditional()
	if !ok {
		return nil, false
	}
	p.scanner.SkipTrivia()
	tok, consumed := p.peekAssignOp()
	if consumed == 0 {
		return lhs, true
	}
	p.consumeRunes(consumed)
	p.scanner.SkipTrivia()
	rhsLine, rhsCol := p.scanner.Line(), p.scanner.Col()
	rhs, ok := p.tryAssignment()
	if !ok {
		return nil, false
	}
	if rhs.Lineno() == 0 && rhs.Charno() == 0 {
		rhs.SetPos(rhsLine, rhsCol)
	}
	n := rhino.New(tok)
	n.SetPos(lhs.Lineno(), lhs.Charno())
	n.AddChildToBack(lhs)
	n.AddChildToBack(rhs)
	return n, true
}

// tryConditional parses `cond ? then : else`. Ternary is right-assoc,
// so both branches recurse into tryAssignment (not tryConditional).
func (p *Parser) tryConditional() (*rhino.Node, bool) {
	cond, ok := p.tryBinary(0)
	if !ok {
		return nil, false
	}
	p.scanner.SkipTrivia()
	if p.scanner.AtEnd() || p.scanner.Peek() != '?' {
		return cond, true
	}
	// Disambiguate from optional chaining `?.` — we don't consume it
	// and leave the `?` for a future optional-chaining port.
	if p.scanner.PeekAt(1) == '.' {
		return cond, true
	}
	// And from nullish coalescing `??`, handled at tryBinary level.
	if p.scanner.PeekAt(1) == '?' {
		return cond, true
	}
	p.scanner.Advance() // '?'
	p.scanner.SkipTrivia()
	thenExpr, ok := p.tryAssignment()
	if !ok {
		return nil, false
	}
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte(':') {
		return nil, false
	}
	p.scanner.SkipTrivia()
	elseExpr, ok := p.tryAssignment()
	if !ok {
		return nil, false
	}
	n := rhino.New(rhino.HOOK)
	n.SetPos(cond.Lineno(), cond.Charno())
	n.AddChildToBack(cond)
	n.AddChildToBack(thenExpr)
	n.AddChildToBack(elseExpr)
	return n, true
}

// tryBinary runs precedence climbing for binary operators. minPrec is
// the minimum precedence allowed at the current level — higher values
// correspond to tighter binding. Matches the canonical Pratt/Crockford
// loop structure, with operators looked up via peekBinOp.
func (p *Parser) tryBinary(minPrec int) (*rhino.Node, bool) {
	left, ok := p.tryUnary()
	if !ok {
		return nil, false
	}
	for {
		p.scanner.SkipTrivia()
		// Keyword binary ops (`instanceof`, `in`) are handled by name
		// lookup since the symbolic peekBinOp can only see characters.
		// Both live at the relational precedence level (8).
		if minPrec <= 8 && p.scanner.AtKeyword("instanceof") {
			p.scanner.ConsumeKeywordRaw("instanceof")
			p.scanner.SkipTrivia()
			right, ok := p.tryBinary(9)
			if !ok {
				return nil, false
			}
			n := rhino.New(rhino.INSTANCEOF)
			n.SetPos(left.Lineno(), left.Charno())
			n.AddChildToBack(left)
			n.AddChildToBack(right)
			left = n
			continue
		}
		if minPrec <= 8 && p.scanner.AtKeyword("in") {
			p.scanner.ConsumeKeywordRaw("in")
			p.scanner.SkipTrivia()
			right, ok := p.tryBinary(9)
			if !ok {
				return nil, false
			}
			n := rhino.New(rhino.IN)
			n.SetPos(left.Lineno(), left.Charno())
			n.AddChildToBack(left)
			n.AddChildToBack(right)
			left = n
			continue
		}
		tok, runes, prec := p.peekBinOp()
		if runes == 0 || prec < minPrec {
			return left, true
		}
		p.consumeRunes(runes)
		p.scanner.SkipTrivia()
		rhsLine, rhsCol := p.scanner.Line(), p.scanner.Col()
		right, ok := p.tryBinary(prec + 1) // left-assoc
		if !ok {
			return nil, false
		}
		// Only backfill when tryBinary produced an unpositioned
		// node (e.g. a bare NAME literal). A paren-wrapped
		// expression like \`(new Bar)\` already has a precise
		// position from its inner tryPostfix, and shouldn't be
		// overwritten with the \`(\` column.
		if right.Lineno() == 0 && right.Charno() == 0 {
			right.SetPos(rhsLine, rhsCol)
		}
		n := rhino.New(tok)
		n.SetPos(left.Lineno(), left.Charno())
		n.AddChildToBack(left)
		n.AddChildToBack(right)
		left = n
	}
}

// tryUnary handles prefix unary operators. Keywords `typeof`, `void`,
// `delete`, and `new` (MemberExpression new) also live here; the
// NewExpression case delegates to finishNew, whose NEW node is then
// allowed to pick up further `.prop` / `(args)` chains via the
// postfix loop.
func (p *Parser) tryUnary() (*rhino.Node, bool) {
	p.scanner.SkipTrivia()
	if p.scanner.AtEnd() {
		return nil, false
	}
	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	c := p.scanner.Peek()
	switch c {
	case '!':
		p.scanner.Advance()
		if p.scanner.Peek() == '=' { // not a unary
			// rewind the '!'; callers won't see this path (peekBinOp
			// covers !=), but we keep the branch total by returning a
			// plain postfix.
		}
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		n := rhino.New(rhino.NOT)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	case '~':
		p.scanner.Advance()
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		n := rhino.New(rhino.BITNOT)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	case '-', '+':
		// Only treat as unary when the next non-space char isn't a
		// digit-shaped continuation of a binary op. tryBinary ensures
		// binary `-` is consumed before we get here.
		p.scanner.Advance()
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		tok := rhino.NEG
		if c == '+' {
			tok = rhino.POS
		}
		n := rhino.New(tok)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	}
	// Keyword-prefix unaries.
	if p.scanner.ConsumeKeyword("typeof") {
		p.scanner.SkipTrivia()
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		n := rhino.New(rhino.TYPEOF)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	}
	if p.scanner.ConsumeKeyword("void") {
		p.scanner.SkipTrivia()
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		n := rhino.New(rhino.VOID)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	}
	if p.scanner.ConsumeKeyword("delete") {
		p.scanner.SkipTrivia()
		operand, ok := p.tryUnary()
		if !ok {
			return nil, false
		}
		n := rhino.New(rhino.DELPROP)
		n.SetPos(startLine, startCol)
		n.AddChildToBack(operand)
		return n, true
	}
	return p.tryPostfix()
}

// tryPostfix handles the left-hand-side postfix chain — property access
// (`.ident`), call (`(args)`), and `new` expressions. Returns a
// LHSExpression suitable for the binary / assignment layers above.
func (p *Parser) tryPostfix() (*rhino.Node, bool) {
	p.scanner.SkipTrivia()
	if p.scanner.AtEnd() {
		return nil, false
	}
	var expr *rhino.Node
	newLine, newCol := p.scanner.Line(), p.scanner.Col()
	if p.scanner.ConsumeKeyword("new") {
		n, ok := p.finishNew()
		if !ok {
			return nil, false
		}
		// Position the NEW node at the `new` keyword, not at the
		// callee identifier — Java's TypeCheck reports at the
		// expression's start column.
		n.SetPos(newLine, newCol)
		expr = n
	} else {
		prim, ok := p.tryPrimary()
		if !ok {
			return nil, false
		}
		expr = prim
	}
	// Postfix chain: property access (`.ident`), computed index access
	// (`[expr]`) and call (`(args)`). Mirrors how Parser.java's
	// parseMemberExpression flattens postfix operators into a left-leaning
	// chain.
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return expr, true
		}
		c := p.scanner.Peek()
		if c == '.' {
			p.scanner.Advance()
			p.scanner.SkipTrivia()
			propLine, propCol := p.scanner.Line(), p.scanner.Col()
			ident, ok := p.scanner.ConsumeIdent()
			if !ok {
				return expr, true
			}
			gp := rhino.New(rhino.GETPROP)
			gp.SetPos(expr.Lineno(), expr.Charno())
			gp.AddChildToBack(expr)
			propNode := rhino.NewString(rhino.STRINGLIT, ident)
			propNode.SetPos(propLine, propCol)
			gp.AddChildToBack(propNode)
			expr = gp
			continue
		}
		if c == '[' {
			// `x[y]` → GETELEM(x, y). Matches Parser.java's
			// parseMemberExpression case for `[`.
			p.scanner.Advance()
			p.scanner.SkipTrivia()
			idxLine, idxCol := p.scanner.Line(), p.scanner.Col()
			idx, ok := p.tryExpression()
			if !ok {
				return nil, false
			}
			idx.SetPos(idxLine, idxCol)
			p.scanner.SkipTrivia()
			if !p.scanner.ConsumeByte(']') {
				return nil, false
			}
			ge := rhino.New(rhino.GETELEM)
			ge.SetPos(expr.Lineno(), expr.Charno())
			ge.AddChildToBack(expr)
			ge.AddChildToBack(idx)
			expr = ge
			continue
		}
		if c == '(' {
			p.scanner.Advance()
			call := rhino.New(rhino.CALL)
			call.SetPos(expr.Lineno(), expr.Charno())
			call.AddChildToBack(expr)
			for {
				p.scanner.SkipTrivia()
				if p.scanner.AtEnd() {
					return nil, false
				}
				if p.scanner.Peek() == ')' {
					p.scanner.Advance()
					break
				}
				argLine, argCol := p.scanner.Line(), p.scanner.Col()
				arg, ok := p.tryExpression()
				if !ok {
					return nil, false
				}
				arg.SetPos(argLine, argCol)
				call.AddChildToBack(arg)
				p.scanner.SkipTrivia()
				if p.scanner.ConsumeByte(',') {
					continue
				}
			}
			expr = call
			continue
		}
		return expr, true
	}
}

// tryPrimary parses literals, identifiers, `this`, array/object
// literals, parenthesized expressions, and arrow / function
// expressions. Mirrors Parser.java's parsePrimaryExpression.
func (p *Parser) tryPrimary() (*rhino.Node, bool) {
	p.scanner.SkipTrivia()
	if p.scanner.AtEnd() {
		return nil, false
	}
	startLine, startCol := p.scanner.Line(), p.scanner.Col()

	// Arrow and parenthesized expressions both start with '('. Try the
	// arrow speculatively; on failure the Save/Restore in tryArrow
	// leaves the scanner ready for a paren parse.
	if p.scanner.Peek() == '(' {
		if arrow, ok := p.tryArrow(); ok {
			return arrow, true
		}
		return p.tryParen()
	}
	fnLine, fnCol := p.scanner.Line(), p.scanner.Col()
	if p.scanner.ConsumeKeyword("function") {
		fn, ok := p.finishFunction(nil)
		if fn != nil {
			fn.SetPos(fnLine, fnCol)
		}
		return fn, ok
	}
	c := p.scanner.Peek()
	switch {
	case c == '\'' || c == '"':
		n, ok := p.parseStringLiteral(c)
		if !ok {
			return nil, false
		}
		n.SetPos(startLine, startCol)
		return n, true
	case unicode.IsDigit(c) || (c == '.' && unicode.IsDigit(p.scanner.PeekAt(1))):
		n, ok := p.parseNumberLiteral()
		if !ok {
			return nil, false
		}
		n.SetPos(startLine, startCol)
		return n, true
	case c == '[':
		return p.tryArrayLit()
	case c == '{':
		return p.tryObjectLit()
	case c == 't' || c == 'f':
		if n, ok := p.parseBooleanLiteral(); ok {
			n.SetPos(startLine, startCol)
			return n, true
		}
	case c == 'n':
		if p.scanner.ConsumeKeyword("null") {
			n := rhino.New(rhino.NULL)
			n.SetPos(startLine, startCol)
			return n, true
		}
	case c == 'u':
		if p.scanner.ConsumeKeyword("undefined") {
			n := rhino.New(rhino.UNDEFINED)
			n.SetPos(startLine, startCol)
			return n, true
		}
	}
	if c == 't' && p.scanner.ConsumeKeyword("this") {
		n := rhino.New(rhino.THIS)
		n.SetPos(startLine, startCol)
		return n, true
	}
	if c == 's' && p.scanner.ConsumeKeyword("super") {
		n := rhino.New(rhino.SUPER)
		n.SetPos(startLine, startCol)
		return n, true
	}
	if unicode.IsLetter(c) || c == '_' || c == '$' {
		ident, ok := p.scanner.ConsumeIdent()
		if !ok {
			return nil, false
		}
		n := rhino.NewString(rhino.NAME, ident)
		n.SetPos(startLine, startCol)
		return n, true
	}
	if c == '/' {
		n, ok := p.parseRegularExpressionLiteral()
		if ok {
			n.SetPos(startLine, startCol)
		}
		return n, ok
	}
	return nil, false
}

// parseRegularExpressionLiteral mirrors Parser.java's
// parseRegularExpressionLiteral. The cursor must be on the leading `/`
// and the call site must have determined that a regex is expected (i.e.
// we are in a primary-expression position, not after an expression
// where `/` would be division — the Pratt binary loop calls peekBinOp,
// which handles division, so this function is only ever reached from
// tryPrimary and is therefore always in the right context).
func (p *Parser) parseRegularExpressionLiteral() (*rhino.Node, bool) {
	text, ok := p.scanner.ScanRegularExpressionLiteral()
	if !ok {
		return nil, false
	}
	return rhino.NewString(rhino.REGEXP, text), true
}

// tryParen consumes `( expr )`. Called from tryPrimary after arrow
// detection fails. Not exposed elsewhere because parenthesized groups
// don't exist at any lower layer.
func (p *Parser) tryParen() (*rhino.Node, bool) {
	save := p.scanner.Save()
	openLine, openCol := p.scanner.Line(), p.scanner.Col()
	if !p.scanner.ConsumeByte('(') {
		p.scanner.Restore(save)
		return nil, false
	}
	p.scanner.SkipTrivia()
	inner, ok := p.tryExpression()
	if !ok {
		p.scanner.Restore(save)
		return nil, false
	}
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte(')') {
		p.scanner.Restore(save)
		return nil, false
	}
	// If the inner expression didn't get a position (e.g. an
	// anonymous function-expression wrapped in parens), use the
	// opening `(` to give it one. Java uses the `(`-column for the
	// whole paren group at diagnostic time.
	if inner.Lineno() == 0 && inner.Charno() == 0 {
		inner.SetPos(openLine, openCol)
	}
	return inner, true
}

// tryArrayLit parses `[ e1, e2, … ]`. Holes (consecutive commas) are
// handled by emitting UNDEFINED sentinels, matching how Closure
// represents ArrayLiteral elisions.
func (p *Parser) tryArrayLit() (*rhino.Node, bool) {
	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	if !p.scanner.ConsumeByte('[') {
		return nil, false
	}
	arr := rhino.New(rhino.ARRAYLIT)
	arr.SetPos(startLine, startCol)
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return nil, false
		}
		if p.scanner.Peek() == ']' {
			p.scanner.Advance()
			return arr, true
		}
		if p.scanner.Peek() == ',' {
			hole := rhino.New(rhino.UNDEFINED)
			hole.SetPos(p.scanner.Line(), p.scanner.Col())
			arr.AddChildToBack(hole)
			p.scanner.Advance()
			continue
		}
		elLine, elCol := p.scanner.Line(), p.scanner.Col()
		el, ok := p.tryExpression()
		if !ok {
			return nil, false
		}
		el.SetPos(elLine, elCol)
		arr.AddChildToBack(el)
		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte(',') {
			continue
		}
	}
}

// tryObjectLit parses `{ key: value, … }`. Keys may be identifiers,
// strings, or numbers. Each pair becomes a STRING_KEY(name) node with
// the value as its sole child. Shorthand (`{x}` → `{x: x}`) and
// computed keys (`{[k]: v}`) are deferred.
func (p *Parser) tryObjectLit() (*rhino.Node, bool) {
	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	if !p.scanner.ConsumeByte('{') {
		return nil, false
	}
	obj := rhino.New(rhino.OBJECTLIT)
	obj.SetPos(startLine, startCol)
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return nil, false
		}
		if p.scanner.Peek() == '}' {
			p.scanner.Advance()
			return obj, true
		}
		// Optional JSDoc on the property value — ignored for now, since
		// record-literal element typing doesn't consult it at Tier-3.
		_, _ = p.scanner.ConsumeBlockComment()
		p.scanner.SkipTrivia()
		// Accessor / method shorthand — `get foo(){…}`, `set foo(x){…}`,
		// `foo(args){…}`. Fully consume the member (params + body) and
		// attach a synthetic value to the object so subsequent members
		// keep parsing. The member shape is lossy (we don't carry
		// getter/setter flags yet), but it's enough to not stall the
		// parser when the corpus sprinkles these into record-style
		// object literals.
		if p.looksLikeObjectMethod() {
			if !p.consumeObjectMethod(obj) {
				return nil, false
			}
			p.scanner.SkipTrivia()
			if p.scanner.ConsumeByte(',') {
				continue
			}
			continue
		}
		keyLine, keyCol := p.scanner.Line(), p.scanner.Col()
		var keyName string
		c := p.scanner.Peek()
		switch {
		case c == '\'' || c == '"':
			p.scanner.Advance()
			start := p.scanner.Pos()
			for !p.scanner.AtEnd() && p.scanner.Peek() != c {
				if p.scanner.Peek() == '\\' && p.scanner.PeekAt(1) != 0 {
					p.scanner.Advance()
					p.scanner.Advance()
					continue
				}
				p.scanner.Advance()
			}
			keyName = p.scanner.Text(start, p.scanner.Pos())
			if !p.scanner.AtEnd() {
				p.scanner.Advance()
			}
		case unicode.IsDigit(c):
			start := p.scanner.Pos()
			for !p.scanner.AtEnd() && (unicode.IsDigit(p.scanner.Peek()) || p.scanner.Peek() == '.') {
				p.scanner.Advance()
			}
			keyName = p.scanner.Text(start, p.scanner.Pos())
		default:
			ident, ok := p.scanner.ConsumeIdent()
			if !ok {
				return nil, false
			}
			keyName = ident
		}
		p.scanner.SkipTrivia()
		// Shorthand (`{x}`) — treat as `{x: x}`.
		if p.scanner.Peek() == ',' || p.scanner.Peek() == '}' {
			key := rhino.NewString(rhino.STRING_KEY, keyName)
			key.SetPos(keyLine, keyCol)
			ref := rhino.NewString(rhino.NAME, keyName)
			ref.SetPos(keyLine, keyCol)
			key.AddChildToBack(ref)
			obj.AddChildToBack(key)
			if p.scanner.ConsumeByte(',') {
				continue
			}
			continue
		}
		if !p.scanner.ConsumeByte(':') {
			return nil, false
		}
		p.scanner.SkipTrivia()
		valLine, valCol := p.scanner.Line(), p.scanner.Col()
		val, ok := p.tryExpression()
		if !ok {
			return nil, false
		}
		val.SetPos(valLine, valCol)
		key := rhino.NewString(rhino.STRING_KEY, keyName)
		key.SetPos(keyLine, keyCol)
		key.AddChildToBack(val)
		obj.AddChildToBack(key)
		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte(',') {
			continue
		}
	}
}

// looksLikeObjectMethod peeks ahead to decide whether the upcoming
// object-literal entry is a method / accessor (`get foo(){}`,
// `set foo(x){}`, `foo(args){}`) rather than a `key: value` pair.
// Non-consuming; restores the scanner on return.
func (p *Parser) looksLikeObjectMethod() bool {
	save := p.scanner.Save()
	defer p.scanner.Restore(save)
	// `get NAME(` or `set NAME(`.
	if p.scanner.ConsumeKeyword("get") || p.scanner.ConsumeKeyword("set") {
		p.scanner.SkipTrivia()
		_, ok := p.scanner.ConsumeIdent()
		if !ok {
			return false
		}
		p.scanner.SkipTrivia()
		return !p.scanner.AtEnd() && p.scanner.Peek() == '('
	}
	// Shorthand method: IDENT(...).
	_, ok := p.scanner.ConsumeIdent()
	if !ok {
		return false
	}
	p.scanner.SkipTrivia()
	return !p.scanner.AtEnd() && p.scanner.Peek() == '('
}

// consumeObjectMethod absorbs a getter / setter / shorthand-method
// entry and attaches a STRING_KEY(<name>) with a FUNCTION child to
// `obj`. Returns false on malformed input so the caller can abort the
// entire object literal.
func (p *Parser) consumeObjectMethod(obj *rhino.Node) bool {
	keyLine, keyCol := p.scanner.Line(), p.scanner.Col()
	// Swallow the `get` / `set` prefix if present; we don't mark the
	// kind today, matching our coarse method-shorthand handling.
	_ = p.scanner.ConsumeKeyword("get") || p.scanner.ConsumeKeyword("set")
	p.scanner.SkipTrivia()
	name, ok := p.scanner.ConsumeIdent()
	if !ok {
		return false
	}
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		return false
	}
	params := rhino.New(rhino.PARAM_LIST)
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return false
		}
		if p.scanner.Peek() == ')' {
			p.scanner.Advance()
			break
		}
		pline, pcol := p.scanner.Line(), p.scanner.Col()
		pn, ok := p.scanner.ConsumeIdent()
		if !ok {
			return false
		}
		pnode := rhino.NewString(rhino.NAME, pn)
		pnode.SetPos(pline, pcol)
		params.AddChildToBack(pnode)
		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte(',') {
			continue
		}
	}
	p.scanner.SkipTrivia()
	body := rhino.New(rhino.BLOCK)
	if p.scanner.ConsumeByte('{') {
		p.parseStatements(body, '}')
		p.scanner.ConsumeByte('}')
	}
	fn := rhino.New(rhino.FUNCTION)
	fn.SetPos(keyLine, keyCol)
	nameInFn := rhino.NewString(rhino.NAME, "")
	nameInFn.SetPos(keyLine, keyCol)
	fn.AddChildToBack(nameInFn)
	fn.AddChildToBack(params)
	fn.AddChildToBack(body)
	key := rhino.NewString(rhino.STRING_KEY, name)
	key.SetPos(keyLine, keyCol)
	key.AddChildToBack(fn)
	obj.AddChildToBack(key)
	return true
}

// peekAssignOp recognises the assignment-op family at the current
// scanner position. Returns the matching Token (ASSIGN / ASSIGN_*) and
// the number of runes that should be consumed on success. The scanner
// is NOT advanced here — the caller is responsible for calling
// consumeRunes.
func (p *Parser) peekAssignOp() (rhino.Token, int) {
	if p.scanner.AtEnd() {
		return rhino.ERROR, 0
	}
	c0 := p.scanner.Peek()
	c1 := p.scanner.PeekAt(1)
	c2 := p.scanner.PeekAt(2)
	c3 := p.scanner.PeekAt(3)
	// Three-rune forms first to avoid misreading `>>>=` as `>>` etc.
	if c0 == '>' && c1 == '>' && c2 == '>' && c3 == '=' {
		return rhino.ASSIGN_URSH, 4
	}
	if c0 == '>' && c1 == '>' && c2 == '=' {
		return rhino.ASSIGN_RSH, 3
	}
	if c0 == '<' && c1 == '<' && c2 == '=' {
		return rhino.ASSIGN_LSH, 3
	}
	// Two-rune forms.
	if c1 == '=' {
		switch c0 {
		case '+':
			return rhino.ASSIGN_ADD, 2
		case '-':
			return rhino.ASSIGN_SUB, 2
		case '*':
			return rhino.ASSIGN_MUL, 2
		case '/':
			return rhino.ASSIGN_DIV, 2
		case '%':
			return rhino.ASSIGN_MOD, 2
		case '&':
			return rhino.ASSIGN_BITAND, 2
		case '|':
			return rhino.ASSIGN_BITOR, 2
		case '^':
			return rhino.ASSIGN_BITXOR, 2
		}
	}
	// Bare `=` — but NOT `==` / `===` / `=>`.
	if c0 == '=' && c1 != '=' && c1 != '>' {
		return rhino.ASSIGN, 1
	}
	return rhino.ERROR, 0
}

// peekBinOp inspects the current position for a binary operator and
// returns (Token, runes, precedence). When no operator is at the
// cursor, returns (ERROR, 0, -1). Assignment and ternary operators
// are deliberately excluded — they're handled at higher layers.
func (p *Parser) peekBinOp() (rhino.Token, int, int) {
	if p.scanner.AtEnd() {
		return rhino.ERROR, 0, -1
	}
	c0 := p.scanner.Peek()
	c1 := p.scanner.PeekAt(1)
	c2 := p.scanner.PeekAt(2)
	// Three-char: `>>>`, `===`, `!==`.
	if c0 == '>' && c1 == '>' && c2 == '>' && p.scanner.PeekAt(3) != '=' {
		return rhino.URSH, 3, 9
	}
	if c0 == '=' && c1 == '=' && c2 == '=' {
		return rhino.SHEQ, 3, 7
	}
	if c0 == '!' && c1 == '=' && c2 == '=' {
		return rhino.SHNE, 3, 7
	}
	// Two-char: `==`, `!=`, `<=`, `>=`, `<<`, `>>`, `&&`, `||`, `??`.
	// Skip these when followed by `=` (compound assign handled above).
	if c2 == '=' && ((c0 == '<' && c1 == '<') || (c0 == '>' && c1 == '>')) {
		return rhino.ERROR, 0, -1
	}
	if c0 == '=' && c1 == '=' {
		return rhino.EQ, 2, 7
	}
	if c0 == '!' && c1 == '=' {
		return rhino.NE, 2, 7
	}
	if c0 == '<' && c1 == '=' {
		return rhino.LE, 2, 8
	}
	if c0 == '>' && c1 == '=' {
		return rhino.GE, 2, 8
	}
	if c0 == '<' && c1 == '<' {
		return rhino.LSH, 2, 9
	}
	if c0 == '>' && c1 == '>' {
		return rhino.RSH, 2, 9
	}
	if c0 == '&' && c1 == '&' {
		return rhino.AND, 2, 3
	}
	if c0 == '|' && c1 == '|' {
		return rhino.OR, 2, 2
	}
	if c0 == '?' && c1 == '?' {
		return rhino.COALESCE, 2, 2
	}
	// One-char forms. Skip when followed by `=` (compound assign) or
	// when they're tokens we don't treat as binary (`&` → bitand only
	// when not `&&=`, which we don't support yet).
	if c1 == '=' {
		// `X=` is an assignment.
		return rhino.ERROR, 0, -1
	}
	switch c0 {
	case '+':
		return rhino.ADD, 1, 10
	case '-':
		return rhino.SUB, 1, 10
	case '*':
		return rhino.MUL, 1, 11
	case '/':
		// `/` could start a regex literal; without regex support we
		// still treat it as divide, which is what the corpus expects.
		// A `/*` / `//` pair is a comment, not division — guard so a
		// trailing JSDoc (`/**`) after an expression doesn't get
		// parsed as the start of a division op.
		if c1 == '*' || c1 == '/' {
			return rhino.ERROR, 0, -1
		}
		return rhino.DIV, 1, 11
	case '%':
		return rhino.MOD, 1, 11
	case '<':
		return rhino.LT, 1, 8
	case '>':
		return rhino.GT, 1, 8
	case '&':
		return rhino.BITAND, 1, 6
	case '|':
		return rhino.BITOR, 1, 4
	case '^':
		return rhino.BITXOR, 1, 5
	}
	return rhino.ERROR, 0, -1
}

// consumeRunes advances the scanner by n runes. Small helper so the
// op-lookup functions above can separate recognition from consumption.
func (p *Parser) consumeRunes(n int) {
	for i := 0; i < n; i++ {
		p.scanner.Advance()
	}
}

// parseInlineTypeJSDoc turns a brace-only inline JSDoc body
// (`{T}` / `?`-prefixed / etc.) into a `@type {T}`-style
// JSDocInfo. Returns nil when the body isn't a bare type
// expression — the full JSDoc parser already handles comments
// with tags.
func parseInlineTypeJSDoc(body string) *rhino.JSDocInfo {
	// Cheap shape filter — require a `{…}` or a bare type
	// name. If the comment contains an `@` tag, defer to the
	// full parser and return nil so the caller keeps its own
	// path.
	for i := 0; i < len(body); i++ {
		if body[i] == '@' {
			return nil
		}
	}
	// Trim whitespace + leading `*` markers common in `/** … */`.
	trimmed := ""
	inLineStar := true
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c == '\n' || c == '\r' {
			trimmed += " "
			inLineStar = true
			continue
		}
		if inLineStar && (c == ' ' || c == '\t' || c == '*') {
			continue
		}
		inLineStar = false
		trimmed += string(c)
	}
	trimmed = trimSpace(trimmed)
	if trimmed == "" {
		return nil
	}
	info := &rhino.JSDocInfo{}
	info.SetType(rhino.NewJSTypeExpression(trimmed))
	return info
}

// trimSpace is a small helper matching strings.TrimSpace but
// avoiding a package import in the parser's hot path.
func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}

// extendWithProperty consumes any trailing `.IDENT` or `[expr]` chain
// after the given base, building GETPROP / GETELEM nodes as in Closure.
// Used by the `this` and NAME paths so member access is uniform with
// the postfix chain in tryPostfix.
func (p *Parser) extendWithProperty(base *rhino.Node) *rhino.Node {
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return base
		}
		c := p.scanner.Peek()
		if c == '.' {
			p.scanner.Advance()
			p.scanner.SkipTrivia()
			propLine, propCol := p.scanner.Line(), p.scanner.Col()
			ident, ok := p.scanner.ConsumeIdent()
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
			continue
		}
		if c == '[' {
			p.scanner.Advance()
			p.scanner.SkipTrivia()
			idxLine, idxCol := p.scanner.Line(), p.scanner.Col()
			idx, ok := p.tryExpression()
			if !ok {
				return base
			}
			idx.SetPos(idxLine, idxCol)
			p.scanner.SkipTrivia()
			if !p.scanner.ConsumeByte(']') {
				return base
			}
			ge := rhino.New(rhino.GETELEM)
			ge.SetPos(base.Lineno(), base.Charno())
			ge.AddChildToBack(base)
			ge.AddChildToBack(idx)
			base = ge
			continue
		}
		return base
	}
}

// tryArrow speculatively parses `(params) => body` or `ident => body`.
// Rewinds on failure so the caller can fall back to other expression
// forms. Arrow functions surface as FUNCTION nodes so the same
// scope/typecheck machinery applies to them.
func (p *Parser) tryArrow() (*rhino.Node, bool) {
	save := p.scanner.Save()
	rewind := func() {
		p.scanner.Restore(save)
	}

	startLine, startCol := p.scanner.Line(), p.scanner.Col()
	var params []string

	if p.scanner.ConsumeByte('(') {
		p.scanner.SkipTrivia()
		if !p.scanner.ConsumeByte(')') {
			for {
				ident, ok := p.scanner.ConsumeIdent()
				if !ok {
					rewind()
					return nil, false
				}
				params = append(params, ident)
				p.scanner.SkipTrivia()
				if p.scanner.ConsumeByte(',') {
					p.scanner.SkipTrivia()
					continue
				}
				if p.scanner.ConsumeByte(')') {
					break
				}
				rewind()
				return nil, false
			}
		}
	} else {
		c := p.scanner.Peek()
		if !(unicode.IsLetter(c) || c == '_' || c == '$') {
			return nil, false
		}
		ident, ok := p.scanner.ConsumeIdent()
		if !ok {
			rewind()
			return nil, false
		}
		params = []string{ident}
	}

	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('=') || !p.scanner.ConsumeByte('>') {
		rewind()
		return nil, false
	}
	p.scanner.SkipTrivia()

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
	if p.scanner.ConsumeByte('{') {
		p.parseStatements(body, '}')
		p.scanner.ConsumeByte('}')
	} else {
		// Single-expression body: wrap as implicit `return EXPR`.
		eline, ecol := p.scanner.Line(), p.scanner.Col()
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
	p.scanner.SkipTrivia()
	nameLine, nameCol := p.scanner.Line(), p.scanner.Col()
	ident, ok := p.scanner.ConsumeIdent()
	if !ok {
		return nil, false
	}
	// `new Foo.Bar(args)` — pick up any `.prop` chain on the callee
	// before the optional `(args)`. Matches Parser.java's
	// parseNewExpression where the MemberExpression can itself be a
	// property chain.
	callee := rhino.NewString(rhino.NAME, ident)
	callee.SetPos(nameLine, nameCol)
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() || p.scanner.Peek() != '.' {
			break
		}
		p.scanner.Advance()
		p.scanner.SkipTrivia()
		pl, pc := p.scanner.Line(), p.scanner.Col()
		pi, ok := p.scanner.ConsumeIdent()
		if !ok {
			break
		}
		gp := rhino.New(rhino.GETPROP)
		gp.SetPos(nameLine, nameCol)
		gp.AddChildToBack(callee)
		pn := rhino.NewString(rhino.STRINGLIT, pi)
		pn.SetPos(pl, pc)
		gp.AddChildToBack(pn)
		callee = gp
	}
	n := rhino.New(rhino.NEW)
	n.SetPos(nameLine, nameCol)
	n.AddChildToBack(callee)
	p.scanner.SkipTrivia()
	if !p.scanner.ConsumeByte('(') {
		// `new Foo` without arg-list is a zero-arg construction.
		return n, true
	}
	for {
		p.scanner.SkipTrivia()
		if p.scanner.AtEnd() {
			return nil, false
		}
		if p.scanner.Peek() == ')' {
			p.scanner.Advance()
			return n, true
		}
		argLine, argCol := p.scanner.Line(), p.scanner.Col()
		arg, ok := p.tryExpression()
		if !ok {
			return nil, false
		}
		arg.SetPos(argLine, argCol)
		n.AddChildToBack(arg)
		p.scanner.SkipTrivia()
		if p.scanner.ConsumeByte(',') {
			continue
		}
	}
}

func (p *Parser) parseBooleanLiteral() (*rhino.Node, bool) {
	if p.scanner.ConsumeKeyword("true") {
		return rhino.New(rhino.TRUE), true
	}
	if p.scanner.ConsumeKeyword("false") {
		return rhino.New(rhino.FALSE), true
	}
	return nil, false
}

func (p *Parser) parseStringLiteral(quote rune) (*rhino.Node, bool) {
	p.scanner.Advance() // opening quote
	start := p.scanner.Pos()
	for !p.scanner.AtEnd() && p.scanner.Peek() != quote {
		if p.scanner.Peek() == '\n' {
			return nil, false
		}
		if p.scanner.Peek() == '\\' && p.scanner.PeekAt(1) != 0 {
			p.scanner.Advance()
			p.scanner.Advance()
			continue
		}
		p.scanner.Advance()
	}
	if p.scanner.AtEnd() {
		return nil, false
	}
	val := p.scanner.Text(start, p.scanner.Pos())
	p.scanner.Advance() // closing quote
	p.scanner.MarkTokenConsumed()
	return rhino.NewString(rhino.STRINGLIT, val), true
}

func (p *Parser) parseNumberLiteral() (*rhino.Node, bool) {
	start := p.scanner.Pos()
	for !p.scanner.AtEnd() && (unicode.IsDigit(p.scanner.Peek()) || p.scanner.Peek() == '.') {
		p.scanner.Advance()
	}
	if p.scanner.Pos() == start {
		return nil, false
	}
	v, err := strconv.ParseFloat(p.scanner.Text(start, p.scanner.Pos()), 64)
	if err != nil {
		return nil, false
	}
	p.scanner.MarkTokenConsumed()
	return rhino.NewNumber(v), true
}

// consumeDeclKeyword reads `var` / `let` / `const` and returns the
// matching AST Token. Mirrors the three checkKeywordAndConsume
// dispatches Closure's parseVariableStatement does in Java.
func (p *Parser) consumeDeclKeyword() (rhino.Token, bool) {
	if p.scanner.ConsumeKeyword("var") {
		return rhino.VAR, true
	}
	if p.scanner.ConsumeKeyword("let") {
		return rhino.LET, true
	}
	if p.scanner.ConsumeKeyword("const") {
		return rhino.CONST, true
	}
	return rhino.ERROR, false
}


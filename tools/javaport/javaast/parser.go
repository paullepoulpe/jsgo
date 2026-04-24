package javaast

import (
	"strings"
	"unicode"
)

// ParseBody parses the raw text between the method's outer braces and
// returns a BlockStmt. Any portion that cannot be parsed is preserved
// as a RawStmt so callers can always emit something.
func ParseBody(src string) *BlockStmt {
	p := &parser{src: src}
	stmts := p.parseStmtList()
	return &BlockStmt{Stmts: stmts}
}

// ---- parser internals --------------------------------------------------

type parser struct {
	src string
	pos int
}

func (p *parser) rest() string { return p.src[p.pos:] }
func (p *parser) eof() bool   { return p.pos >= len(p.src) }

// peek returns the byte at the current position, or 0 at EOF.
func (p *parser) peek() byte {
	if p.pos >= len(p.src) {
		return 0
	}
	return p.src[p.pos]
}

// advance moves the cursor forward by n bytes.
func (p *parser) advance(n int) { p.pos += n }

// skipWS skips whitespace and line/block comments.
func (p *parser) skipWS() {
	for !p.eof() {
		c := p.src[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			p.pos++
			continue
		}
		if strings.HasPrefix(p.rest(), "//") {
			nl := strings.Index(p.rest(), "\n")
			if nl < 0 {
				p.pos = len(p.src)
				return
			}
			p.pos += nl + 1
			continue
		}
		if strings.HasPrefix(p.rest(), "/*") {
			end := strings.Index(p.rest(), "*/")
			if end < 0 {
				p.pos = len(p.src)
				return
			}
			p.pos += end + 2
			continue
		}
		break
	}
}

// consume advances past the literal string s (after skipping WS) and
// returns true, or returns false and leaves pos unchanged on mismatch.
func (p *parser) consume(s string) bool {
	p.skipWS()
	if strings.HasPrefix(p.rest(), s) {
		p.pos += len(s)
		return true
	}
	return false
}

// consumeKeyword is like consume but requires that the match is not
// immediately followed by an identifier character.
func (p *parser) consumeKeyword(kw string) bool {
	p.skipWS()
	if strings.HasPrefix(p.rest(), kw) {
		after := p.pos + len(kw)
		if after < len(p.src) && isIdentChar(p.src[after]) {
			return false
		}
		p.pos = after
		return true
	}
	return false
}

// peekKeyword reports whether the next token is kw (without consuming).
func (p *parser) peekKeyword(kw string) bool {
	saved := p.pos
	ok := p.consumeKeyword(kw)
	p.pos = saved
	return ok
}

// readIdent reads one Java identifier (including `$`).
func (p *parser) readIdent() (string, bool) {
	p.skipWS()
	if p.eof() || !isIdentStart(p.peek()) {
		return "", false
	}
	start := p.pos
	for !p.eof() && isIdentChar(p.src[p.pos]) {
		p.pos++
	}
	return p.src[start:p.pos], true
}

// peekIdent returns the next identifier token without consuming it.
func (p *parser) peekIdent() (string, bool) {
	saved := p.pos
	id, ok := p.readIdent()
	p.pos = saved
	return id, ok
}

// readUntilStop reads from the current position until it finds `stop`
// at nesting depth 0 (where '(' increments depth, ')' decrements it,
// and similarly for '['/']' and '{'/'}').  The inner text (before
// `stop`) is returned and pos is left just past the stop byte.
// Returns ("", false) if stop is never found.
//
// The caller is responsible for having already consumed the matching
// open delimiter (e.g. '(' before calling for stop=')').
func (p *parser) readUntilStop(stop byte) (string, bool) {
	saved := p.pos
	depth := 0
	for !p.eof() {
		c := p.src[p.pos]
		// Skip string / char literals.
		if c == '"' || c == '\'' {
			p.pos = skipStringLitP(p.src, p.pos) + 1
			continue
		}
		// Skip line comments.
		if c == '/' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '/' {
			nl := strings.Index(p.src[p.pos:], "\n")
			if nl < 0 {
				p.pos = len(p.src)
				break
			}
			p.pos += nl + 1
			continue
		}
		// Skip block comments.
		if c == '/' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '*' {
			end := strings.Index(p.src[p.pos:], "*/")
			if end < 0 {
				p.pos = len(p.src)
				break
			}
			p.pos += end + 2
			continue
		}
		if c == stop && depth == 0 {
			text := p.src[saved:p.pos]
			p.pos++ // consume the stop delimiter
			return text, true
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
			// at depth==0 and c!=stop we just keep scanning; this can
			// happen for unbalanced source we're doing best-effort on.
		}
		p.pos++
	}
	p.pos = saved
	return "", false
}

// readParenContent reads `(content)` and returns `content`.
func (p *parser) readParenContent() (string, bool) {
	if !p.consume("(") {
		return "", false
	}
	return p.readUntilStop(')')
}

// readBraceBlock reads `{ ... }` and returns the inner text.
func (p *parser) readBraceBlock() (string, bool) {
	if !p.consume("{") {
		return "", false
	}
	return p.readUntilStop('}')
}

// ---- Statement list ----------------------------------------------------

func (p *parser) parseStmtList() []Stmt {
	var stmts []Stmt
	for {
		p.skipWS()
		if p.eof() {
			break
		}
		// End of enclosing block
		if p.peek() == '}' {
			break
		}
		s := p.parseStmt()
		if s == nil {
			break
		}
		stmts = append(stmts, s)
	}
	return stmts
}

// parseStmt dispatches on the leading keyword or token.
func (p *parser) parseStmt() Stmt {
	p.skipWS()
	if p.eof() || p.peek() == '}' {
		return nil
	}

	// Annotations on their own line — skip.
	if p.peek() == '@' {
		nl := strings.Index(p.rest(), "\n")
		if nl < 0 {
			p.pos = len(p.src)
			return nil
		}
		p.pos += nl + 1
		return p.parseStmt()
	}

	// Block
	if p.peek() == '{' {
		inner, ok := p.readBraceBlock()
		if !ok {
			return p.rawLine()
		}
		sub := &parser{src: inner}
		return &BlockStmt{Stmts: sub.parseStmtList()}
	}

	// Semicolon → empty statement; skip.
	if p.peek() == ';' {
		p.pos++
		return p.parseStmt()
	}

	// if
	if p.peekKeyword("if") {
		return p.parseIf()
	}
	// for
	if p.peekKeyword("for") {
		return p.parseFor()
	}
	// while
	if p.peekKeyword("while") {
		return p.parseWhile()
	}
	// do
	if p.peekKeyword("do") {
		return p.parseDo()
	}
	// switch
	if p.peekKeyword("switch") {
		return p.parseSwitch()
	}
	// try
	if p.peekKeyword("try") {
		return p.parseTry()
	}
	// return
	if p.peekKeyword("return") {
		return p.parseReturn()
	}
	// throw
	if p.peekKeyword("throw") {
		return p.parseThrow()
	}
	// break
	if p.peekKeyword("break") {
		return p.parseBreak()
	}
	// continue
	if p.peekKeyword("continue") {
		return p.parseContinue()
	}
	// final / var / type-initiated local var
	if s, ok := p.tryLocalVar(); ok {
		return s
	}
	// label: stmt
	if s, ok := p.tryLabel(); ok {
		return s
	}
	// expression statement (assignment, call, pre/post inc)
	return p.parseExprStmt()
}

// ---- if ----------------------------------------------------------------

func (p *parser) parseIf() Stmt {
	if !p.consumeKeyword("if") {
		return p.rawLine()
	}
	condText, ok := p.readParenContent()
	if !ok {
		return p.rawLine()
	}
	cond := parseExpr(condText)
	p.skipWS()
	then := p.parseStmt()
	if then == nil {
		return &IfStmt{Cond: cond, Then: &BlockStmt{}}
	}
	p.skipWS()
	var els Stmt
	if p.peekKeyword("else") {
		p.consumeKeyword("else")
		p.skipWS()
		els = p.parseStmt()
	}
	return &IfStmt{Cond: cond, Then: then, Else: els}
}

// ---- for ---------------------------------------------------------------

func (p *parser) parseFor() Stmt {
	if !p.consumeKeyword("for") {
		return p.rawLine()
	}
	headerText, ok := p.readParenContent()
	if !ok {
		return p.rawLine()
	}
	p.skipWS()
	body := p.parseStmt()
	if body == nil {
		body = &BlockStmt{}
	}

	// Detect for-each: `Type name : expr`
	if colon := findColonSep(headerText); colon >= 0 {
		left := strings.TrimSpace(headerText[:colon])
		right := strings.TrimSpace(headerText[colon+1:])
		typeName, varName := splitTypeVar(left)
		return &ForEachStmt{
			Type: typeName,
			Name: varName,
			Coll: parseExpr(right),
			Body: body,
		}
	}

	// Classic for: split on `;`
	parts := splitAtSemicolon(headerText)
	if len(parts) != 3 {
		return &RawStmt{Text: "for (" + headerText + ") { ... }"}
	}
	initText := strings.TrimSpace(parts[0])
	condText := strings.TrimSpace(parts[1])
	postText := strings.TrimSpace(parts[2])

	var initStmt, postStmt Stmt
	if initText != "" {
		ip := &parser{src: initText + ";"}
		initStmt = ip.parseStmt()
	}
	var condExpr Expr
	if condText != "" {
		condExpr = parseExpr(condText)
	}
	if postText != "" {
		ep := &parser{src: postText + ";"}
		postStmt = ep.parseStmt()
	}
	return &ForStmt{Init: initStmt, Cond: condExpr, Post: postStmt, Body: body}
}

// findColonSep returns the index of a `:` in s that is NOT inside
// parens/angles (i.e. a for-each colon, not a ternary colon). Returns -1 if
// none found.
func findColonSep(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case ':':
			if depth == 0 {
				// Make sure it's not `::` (method ref) or part of ternary deeply nested.
				if i+1 < len(s) && s[i+1] == ':' {
					i++
					continue
				}
				return i
			}
		}
	}
	return -1
}

// splitAtSemicolon splits s on `;` outside paren/bracket depth.
// Intentionally does NOT track `<`/`>` since those appear as comparison
// operators in for-loop headers (e.g. `i < n`).
func splitAtSemicolon(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// splitTypeVar splits "Type name" into (type, name), handling generics.
func splitTypeVar(s string) (string, string) {
	s = strings.TrimSpace(s)
	// strip `final`
	if strings.HasPrefix(s, "final ") {
		s = strings.TrimSpace(s[6:])
	}
	sp := lastTopLevelSpaceP(s)
	if sp < 0 {
		return s, ""
	}
	return strings.TrimSpace(s[:sp]), strings.TrimSpace(s[sp+1:])
}

func lastTopLevelSpaceP(s string) int {
	depth := 0
	last := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case ' ':
			if depth == 0 {
				last = i
			}
		}
	}
	return last
}

// ---- while / do-while --------------------------------------------------

func (p *parser) parseWhile() Stmt {
	if !p.consumeKeyword("while") {
		return p.rawLine()
	}
	condText, ok := p.readParenContent()
	if !ok {
		return p.rawLine()
	}
	p.skipWS()
	body := p.parseStmt()
	if body == nil {
		body = &BlockStmt{}
	}
	return &WhileStmt{Cond: parseExpr(condText), Body: body}
}

func (p *parser) parseDo() Stmt {
	if !p.consumeKeyword("do") {
		return p.rawLine()
	}
	p.skipWS()
	body := p.parseStmt()
	if body == nil {
		body = &BlockStmt{}
	}
	if !p.consumeKeyword("while") {
		return body
	}
	condText, ok := p.readParenContent()
	if !ok {
		return body
	}
	p.consume(";")
	return &DoWhileStmt{Body: body, Cond: parseExpr(condText)}
}

// ---- switch ------------------------------------------------------------

func (p *parser) parseSwitch() Stmt {
	if !p.consumeKeyword("switch") {
		return p.rawLine()
	}
	exprText, ok := p.readParenContent()
	if !ok {
		return p.rawLine()
	}
	body, ok := p.readBraceBlock()
	if !ok {
		return p.rawLine()
	}
	sw := &SwitchStmt{X: parseExpr(exprText)}
	sw.Cases = parseSwitchBody(body)
	return sw
}

func parseSwitchBody(body string) []SwitchCase {
	// Split by `case X:` and `default:` at depth 0.
	var cases []SwitchCase
	sp := &parser{src: body}
	for {
		sp.skipWS()
		if sp.eof() {
			break
		}
		if sp.peekKeyword("case") {
			sp.consumeKeyword("case")
			// Read until `:` (at depth 0, mind nested ternaries)
			val, ok := sp.readUntilColon()
			if !ok {
				break
			}
			stmts, hasFall := sp.parseCaseBody()
			cases = append(cases, SwitchCase{
				Value:    parseExpr(strings.TrimSpace(val)),
				Stmts:    stmts,
				FallThru: hasFall,
			})
		} else if sp.peekKeyword("default") {
			sp.consumeKeyword("default")
			sp.consume(":")
			stmts, hasFall := sp.parseCaseBody()
			cases = append(cases, SwitchCase{
				Value:    nil,
				Stmts:    stmts,
				FallThru: hasFall,
			})
		} else {
			// Unexpected content — consume a line.
			nl := strings.Index(sp.rest(), "\n")
			if nl < 0 {
				break
			}
			sp.pos += nl + 1
		}
	}
	return cases
}

// readUntilColon reads text up to and including the next `:` at depth 0.
func (p *parser) readUntilColon() (string, bool) {
	start := p.pos
	depth := 0
	for !p.eof() {
		c := p.src[p.pos]
		switch c {
		case '(', '[', '<':
			depth++
		case ')', ']', '>':
			depth--
		case ':':
			if depth == 0 {
				text := p.src[start:p.pos]
				p.pos++ // consume ':'
				return text, true
			}
		}
		p.pos++
	}
	return "", false
}

// parseCaseBody reads statements until the next `case`/`default`/`}`.
// Returns the stmts and whether the last stmt has fall-through.
func (p *parser) parseCaseBody() ([]Stmt, bool) {
	var stmts []Stmt
	for {
		p.skipWS()
		if p.eof() || p.peek() == '}' {
			break
		}
		if p.peekKeyword("case") || p.peekKeyword("default") {
			break
		}
		s := p.parseStmt()
		if s == nil {
			break
		}
		stmts = append(stmts, s)
	}
	// Fall-through if last stmt is not break/return.
	hasFall := false
	if len(stmts) > 0 {
		last := stmts[len(stmts)-1]
		switch last.(type) {
		case *BreakStmt, *ReturnStmt:
			hasFall = false
		default:
			hasFall = true
		}
	}
	return stmts, hasFall
}

// ---- try ---------------------------------------------------------------

func (p *parser) parseTry() Stmt {
	if !p.consumeKeyword("try") {
		return p.rawLine()
	}
	// skip optional resources: `try (X x = ...)` — treat as opaque
	p.skipWS()
	if p.peek() == '(' {
		p.readParenContent()
	}
	bodyText, ok := p.readBraceBlock()
	if !ok {
		return p.rawLine()
	}
	bodyBlock := &parser{src: bodyText}
	body := &BlockStmt{Stmts: bodyBlock.parseStmtList()}
	ts := &TryStmt{Body: body}
	for p.peekKeyword("catch") {
		p.consumeKeyword("catch")
		catchHead, ok := p.readParenContent()
		if !ok {
			break
		}
		catchBodyText, ok := p.readBraceBlock()
		if !ok {
			break
		}
		tp, name := splitTypeVar(catchHead)
		cb := &parser{src: catchBodyText}
		ts.Catches = append(ts.Catches, CatchClause{
			Type: tp,
			Name: name,
			Body: &BlockStmt{Stmts: cb.parseStmtList()},
		})
	}
	if p.peekKeyword("finally") {
		p.consumeKeyword("finally")
		finText, ok := p.readBraceBlock()
		if !ok {
			return ts
		}
		fb := &parser{src: finText}
		ts.Finally = &BlockStmt{Stmts: fb.parseStmtList()}
	}
	return ts
}

// ---- return / throw / break / continue ---------------------------------

func (p *parser) parseReturn() Stmt {
	p.consumeKeyword("return")
	p.skipWS()
	if p.peek() == ';' {
		p.pos++
		return &ReturnStmt{}
	}
	exprText, ok := p.readUntilStop(';')
	if !ok {
		return &ReturnStmt{}
	}
	return &ReturnStmt{X: parseExpr(strings.TrimSpace(exprText))}
}

func (p *parser) parseThrow() Stmt {
	p.consumeKeyword("throw")
	p.skipWS()
	exprText, ok := p.readUntilStop(';')
	if !ok {
		return &ThrowStmt{X: &RawExpr{Text: "???"}}
	}
	return &ThrowStmt{X: parseExpr(strings.TrimSpace(exprText))}
}

func (p *parser) parseBreak() Stmt {
	p.consumeKeyword("break")
	p.skipWS()
	if p.peek() == ';' {
		p.pos++
		return &BreakStmt{}
	}
	label, _ := p.readIdent()
	p.consume(";")
	return &BreakStmt{Label: label}
}

func (p *parser) parseContinue() Stmt {
	p.consumeKeyword("continue")
	p.skipWS()
	if p.peek() == ';' {
		p.pos++
		return &ContinueStmt{}
	}
	label, _ := p.readIdent()
	p.consume(";")
	return &ContinueStmt{Label: label}
}

// ---- local var ---------------------------------------------------------

// tryLocalVar tries to parse `[final] Type name [= init];`.
// Returns (stmt, true) on success or (nil, false) if this doesn't look
// like a local variable declaration.
func (p *parser) tryLocalVar() (Stmt, bool) {
	saved := p.pos
	// Optional `final`.
	hasFinal := p.consumeKeyword("final")
	// Try to read a type token.
	id, ok := p.readIdent()
	if !ok {
		p.pos = saved
		return nil, false
	}
	// Consume generic suffix: `<...>`.
	typeStr := id
	p.skipWS()
	if p.peek() == '<' {
		p.pos++ // consume <
		inner, ok2 := p.readUntilStop('>')
		if ok2 {
			typeStr = id + "<" + inner + ">"
		}
	}
	// Array suffix
	for p.consume("[]") {
		typeStr += "[]"
	}
	// Read variable name.
	varName, ok2 := p.readIdent()
	if !ok2 {
		p.pos = saved
		return nil, false
	}
	// The next non-WS token must be `=`, `;`, or `,` for it to be a var decl.
	// If it's `(` then this is a method call, not a var decl.
	p.skipWS()
	if p.peek() == '(' {
		p.pos = saved
		return nil, false
	}
	// Also: if typeStr is a keyword-only word that could be a method name,
	// and varName looks like a call arg, bail.
	if !hasFinal && !isJavaType(typeStr) && !isUpperStart(typeStr) {
		// Might be a method call like `foo.bar(...)` — let expr handle it.
		p.pos = saved
		return nil, false
	}
	var init Expr
	if p.consume("=") {
		exprText, ok3 := p.readUntilStop(';')
		if !ok3 {
			p.pos = saved
			return nil, false
		}
		init = parseExpr(strings.TrimSpace(exprText))
	} else if p.consume(";") {
		// bare decl
	} else {
		p.pos = saved
		return nil, false
	}
	return &LocalVarStmt{Type: typeStr, Name: varName, Init: init}, true
}

// isJavaType returns true for primitive types and common collection names
// that we know must be types (not method names) when starting a statement.
func isJavaType(s string) bool {
	switch s {
	case "int", "long", "boolean", "double", "float", "short", "byte", "char",
		"String", "Object", "var",
		"List", "Map", "Set", "Collection", "Iterable",
		"ImmutableList", "ImmutableMap", "ImmutableSet",
		"Optional", "Iterator", "ArrayList", "HashMap", "HashSet",
		"LinkedList", "TreeMap", "TreeSet":
		return true
	}
	return false
}

func isUpperStart(s string) bool {
	return len(s) > 0 && unicode.IsUpper(rune(s[0]))
}

// ---- label -------------------------------------------------------------

func (p *parser) tryLabel() (Stmt, bool) {
	saved := p.pos
	id, ok := p.readIdent()
	if !ok {
		return nil, false
	}
	p.skipWS()
	if p.peek() != ':' {
		p.pos = saved
		return nil, false
	}
	p.pos++ // consume ':'
	p.skipWS()
	inner := p.parseStmt()
	return &LabeledStmt{Label: id, Stmt: inner}, true
}

// ---- expression statement ----------------------------------------------

func (p *parser) parseExprStmt() Stmt {
	exprText, ok := p.readUntilStop(';')
	if !ok {
		return p.rawLine()
	}
	text := strings.TrimSpace(exprText)
	if text == "" {
		return p.parseStmt()
	}
	// Detect compound assignment.
	op, lhs, rhs, isAssign := splitAssignment(text)
	if isAssign {
		return &AssignStmt{LHS: parseExpr(lhs), Op: op, RHS: parseExpr(rhs)}
	}
	return &ExprStmt{X: parseExpr(text)}
}

// rawLine consumes to the next newline or semicolon and returns a RawStmt.
func (p *parser) rawLine() Stmt {
	start := p.pos
	for !p.eof() {
		c := p.src[p.pos]
		if c == '\n' || c == ';' {
			p.pos++
			break
		}
		p.pos++
	}
	return &RawStmt{Text: strings.TrimSpace(p.src[start:p.pos])}
}

// ---- assignment detection ----------------------------------------------

var compoundOps = []string{">>>=", "<<=", ">>=", "+=", "-=", "*=", "/=", "%=", "|=", "&=", "^=", "||=", "&&=", "="}

// splitAssignment tries to find an assignment operator at the top level of s.
// Returns (op, lhs, rhs, true) or ("", "", "", false).
func splitAssignment(s string) (string, string, string, bool) {
	depth := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' || c == '\'' {
			i = skipStringLitP(s, i)
			continue
		}
		switch c {
		case '(', '[', '{', '<':
			depth++
		case ')', ']', '}', '>':
			if depth > 0 {
				depth--
			}
		}
		if depth > 0 {
			continue
		}
		for _, op := range compoundOps {
			if strings.HasPrefix(s[i:], op) {
				// exclude `==`, `!=`, `<=`, `>=`
				before := ""
				if i > 0 {
					before = string(s[i-1])
				}
				if op == "=" && (before == "!" || before == "<" || before == ">" || before == "=") {
					continue
				}
				after := i + len(op)
				if op == "=" && after < len(s) && s[after] == '=' {
					continue
				}
				lhs := strings.TrimSpace(s[:i])
				rhs := strings.TrimSpace(s[after:])
				return op, lhs, rhs, true
			}
		}
	}
	return "", "", "", false
}

// ---- expression parser -------------------------------------------------

// parseExpr parses a full Java expression string into an Expr tree.
// It handles: ternary, binary operators, instanceof, unary, postfix
// (calls, field access, array index), cast, new, literals, identifiers.
func parseExpr(src string) Expr {
	src = strings.TrimSpace(src)
	if src == "" {
		return &LitExpr{Val: ""}
	}
	ep := &exprParser{src: src}
	e := ep.parseTernary()
	return e
}

type exprParser struct {
	src string
	pos int
}

func (ep *exprParser) rest() string { return ep.src[ep.pos:] }
func (ep *exprParser) eof() bool   { return ep.pos >= len(ep.src) }
func (ep *exprParser) peek() byte {
	if ep.pos >= len(ep.src) {
		return 0
	}
	return ep.src[ep.pos]
}

func (ep *exprParser) skipWS() {
	for !ep.eof() {
		c := ep.src[ep.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			ep.pos++
		} else {
			break
		}
	}
}

func (ep *exprParser) consume(s string) bool {
	ep.skipWS()
	if strings.HasPrefix(ep.rest(), s) {
		ep.pos += len(s)
		return true
	}
	return false
}

func (ep *exprParser) parseTernary() Expr {
	cond := ep.parseBinary(0)
	ep.skipWS()
	if ep.peek() != '?' {
		return cond
	}
	ep.pos++ // consume '?'
	then := ep.parseTernary()
	if !ep.consume(":") {
		return &TernaryExpr{Cond: cond, Then: then, Else: &RawExpr{Text: "???"}}
	}
	els := ep.parseTernary()
	return &TernaryExpr{Cond: cond, Then: then, Else: els}
}

// Operator precedence table (higher value = tighter binding).
var binaryOps = [][]string{
	{"||"},
	{"&&"},
	{"|"},
	{"^"},
	{"&"},
	{"==", "!="},
	{"<", ">", "<=", ">="},
	{"<<", ">>>", ">>"},
	{"+", "-"},
	{"*", "/", "%"},
}

func (ep *exprParser) parseBinary(prec int) Expr {
	if prec >= len(binaryOps) {
		return ep.parseUnary()
	}
	left := ep.parseBinary(prec + 1)
	for {
		ep.skipWS()
		// instanceof
		if prec == 5 { // same level as comparison
			if strings.HasPrefix(ep.rest(), "instanceof ") || strings.HasPrefix(ep.rest(), "instanceof\t") {
				ep.pos += len("instanceof")
				ep.skipWS()
				start := ep.pos
				for !ep.eof() && (isIdentChar(ep.src[ep.pos]) || ep.src[ep.pos] == '.') {
					ep.pos++
				}
				typ := strings.TrimSpace(ep.src[start:ep.pos])
				left = &TypeCheckExpr{X: left, Type: typ}
				continue
			}
		}
		op := ep.matchOp(binaryOps[prec])
		if op == "" {
			break
		}
		right := ep.parseBinary(prec + 1)
		left = &BinaryExpr{X: left, Op: op, Y: right}
	}
	return left
}

func (ep *exprParser) matchOp(ops []string) string {
	ep.skipWS()
	for _, op := range ops {
		if strings.HasPrefix(ep.rest(), op) {
			// Make sure we're not matching a longer operator.
			after := op
			rest := ep.rest()
			if len(rest) > len(after) {
				next := rest[len(after)]
				// Avoid matching `=` when it's `==`, `<` when `<<` etc.
				if (op == "=" && next == '=') ||
					(op == "!" && next == '=') ||
					(op == "<" && (next == '<' || next == '=')) ||
					(op == ">" && (next == '>' || next == '=')) ||
					(op == "&" && next == '&') ||
					(op == "|" && next == '|') ||
					(op == "+" && next == '+') ||
					(op == "-" && next == '-') {
					continue
				}
			}
			ep.pos += len(op)
			return op
		}
	}
	return ""
}

func (ep *exprParser) parseUnary() Expr {
	ep.skipWS()
	// Prefix operators
	for _, op := range []string{"!", "-", "~", "++", "--"} {
		if strings.HasPrefix(ep.rest(), op) {
			// Don't consume `--` inside `--foo` vs `a - -b`
			ep.pos += len(op)
			x := ep.parsePostfix()
			return &UnaryExpr{Op: op, X: x}
		}
	}
	// Cast: `(Type) expr` — heuristic: `(Identifier)` followed by non-operator
	if ep.peek() == '(' {
		saved := ep.pos
		ep.pos++ // consume '('
		ep.skipWS()
		castStart := ep.pos
		for !ep.eof() && (isIdentChar(ep.src[ep.pos]) || ep.src[ep.pos] == '.' || ep.src[ep.pos] == '[' || ep.src[ep.pos] == ']' || ep.src[ep.pos] == '<' || ep.src[ep.pos] == '>' || ep.src[ep.pos] == ',') {
			ep.pos++
		}
		ep.skipWS()
		if ep.peek() == ')' && ep.pos > castStart {
			typeStr := strings.TrimSpace(ep.src[castStart:ep.pos])
			ep.pos++ // consume ')'
			ep.skipWS()
			// Check it looks like a type (starts uppercase, or primitive)
			if isJavaCastType(typeStr) && !ep.eof() && ep.peek() != ')' && ep.peek() != '+' && ep.peek() != '-' && ep.peek() != '*' && ep.peek() != '/' {
				x := ep.parsePostfix()
				return &CastExpr{Type: typeStr, X: x}
			}
		}
		ep.pos = saved
	}
	return ep.parsePostfix()
}

func isJavaCastType(s string) bool {
	if s == "" {
		return false
	}
	switch s {
	case "int", "long", "boolean", "double", "float", "short", "byte", "char", "String", "Object":
		return true
	}
	return isUpperStart(s) || strings.Contains(s, "[]")
}

func (ep *exprParser) parsePostfix() Expr {
	base := ep.parsePrimary()
	for {
		ep.skipWS()
		if ep.eof() {
			break
		}
		c := ep.peek()
		if c == '.' {
			ep.pos++
			ep.skipWS()
			// method ref `::` skipped
			if strings.HasPrefix(ep.rest(), ":") {
				break
			}
			id := ep.readIdentStr()
			if id == "" {
				break
			}
			ep.skipWS()
			if ep.peek() == '(' {
				args := ep.readArgs()
				base = &CallExpr{
					Fun:  &SelectorExpr{X: base, Sel: id},
					Args: args,
				}
			} else {
				base = &SelectorExpr{X: base, Sel: id}
			}
		} else if c == '[' {
			ep.pos++ // consume '['
			inner, ok := ep.readUntil(']')
			if !ok {
				break
			}
			base = &IndexExpr{X: base, Idx: parseExpr(inner)}
		} else if c == '+' && ep.pos+1 < len(ep.src) && ep.src[ep.pos+1] == '+' {
			ep.pos += 2
			base = &UnaryExpr{Op: "++", X: base}
		} else if c == '-' && ep.pos+1 < len(ep.src) && ep.src[ep.pos+1] == '-' {
			ep.pos += 2
			base = &UnaryExpr{Op: "--", X: base}
		} else {
			break
		}
	}
	return base
}

func (ep *exprParser) parsePrimary() Expr {
	ep.skipWS()
	if ep.eof() {
		return &RawExpr{Text: ""}
	}

	// new
	if strings.HasPrefix(ep.rest(), "new ") || strings.HasPrefix(ep.rest(), "new\t") {
		ep.pos += 3
		ep.skipWS()
		typeName := ep.readTypeStr()
		ep.skipWS()
		if ep.peek() == '[' {
			ep.pos++ // consume '['
			if ep.peek() == ']' {
				ep.pos++ // consume ']'
				// Array initialiser: new Type[]{...} or just new Type[]
				ep.skipWS()
				if ep.peek() == '{' {
					ep.pos++ // consume '{'
					inner, _ := ep.readUntil('}')
					elems := splitCommaExprs(inner)
					return &NewArrayExpr{ElemType: typeName, Elems: elems}
				}
				return &NewArrayExpr{ElemType: typeName}
			}
			sizeText, _ := ep.readUntil(']')
			return &NewArrayExpr{ElemType: typeName, Size: parseExpr(sizeText)}
		}
		if ep.peek() == '(' {
			args := ep.readArgs()
			// Anonymous class body — skip it.
			ep.skipWS()
			if ep.peek() == '{' {
				ep.pos++ // consume '{'
				ep.readUntil('}')
			}
			return &NewExpr{Type: typeName, Args: args}
		}
		return &NewExpr{Type: typeName}
	}

	// Array literal `{...}`
	if ep.peek() == '{' {
		ep.pos++
		inner, _ := ep.readUntil('}')
		return &ArrayLitExpr{Elems: splitCommaExprs(inner)}
	}

	// Parenthesised expression (or cast — handled by parseUnary already)
	if ep.peek() == '(' {
		ep.pos++ // consume '('
		inner, ok := ep.readUntil(')')
		if !ok {
			return &RawExpr{Text: ep.rest()}
		}
		return parseExpr(inner)
	}

	// String literal
	if ep.peek() == '"' {
		start := ep.pos
		ep.pos = skipStringLitP(ep.src, ep.pos) + 1
		return &LitExpr{Val: ep.src[start:ep.pos]}
	}

	// Char literal
	if ep.peek() == '\'' {
		start := ep.pos
		ep.pos = skipStringLitP(ep.src, ep.pos) + 1
		return &LitExpr{Val: ep.src[start:ep.pos]}
	}

	// Number literal
	if ep.peek() >= '0' && ep.peek() <= '9' {
		return ep.readNumber()
	}

	// Identifier / keyword
	id := ep.readIdentStr()
	if id == "" {
		// Unknown — consume the rest as raw.
		raw := ep.rest()
		ep.pos = len(ep.src)
		return &RawExpr{Text: raw}
	}
	switch id {
	case "true", "false":
		return &LitExpr{Val: id}
	case "null":
		return &LitExpr{Val: "null"}
	case "this":
		return &IdentExpr{Name: "this"}
	case "super":
		return &IdentExpr{Name: "super"}
	}
	ep.skipWS()
	if ep.peek() == '(' {
		args := ep.readArgs()
		return &CallExpr{Fun: &IdentExpr{Name: id}, Args: args}
	}
	return &IdentExpr{Name: id}
}

func (ep *exprParser) readNumber() Expr {
	start := ep.pos
	for !ep.eof() {
		c := ep.src[ep.pos]
		if (c >= '0' && c <= '9') || c == '.' || c == 'x' || c == 'X' ||
			c == 'L' || c == 'l' || c == 'f' || c == 'F' || c == 'e' || c == 'E' ||
			(c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || c == '_' {
			ep.pos++
		} else {
			break
		}
	}
	return &LitExpr{Val: ep.src[start:ep.pos]}
}

func (ep *exprParser) readIdentStr() string {
	ep.skipWS()
	if ep.eof() || !isIdentStart(ep.peek()) {
		return ""
	}
	start := ep.pos
	for !ep.eof() && isIdentChar(ep.src[ep.pos]) {
		ep.pos++
	}
	return ep.src[start:ep.pos]
}

// readTypeStr reads a Java type name including generics: `Foo<Bar, Baz>`.
func (ep *exprParser) readTypeStr() string {
	id := ep.readIdentStr()
	ep.skipWS()
	if ep.peek() == '<' {
		ep.pos++ // consume '<'
		inner, _ := ep.readUntil('>')
		id += "<" + inner + ">"
	}
	// Array suffix
	for strings.HasPrefix(ep.rest(), "[]") {
		ep.pos += 2
		id += "[]"
	}
	return id
}

// readArgs reads `(arg1, arg2, ...)` and returns the parsed args.
func (ep *exprParser) readArgs() []Expr {
	if ep.peek() != '(' {
		return nil
	}
	ep.pos++ // consume '('
	inner, _ := ep.readUntil(')')
	return splitCommaExprs(inner)
}

// readUntil reads until char stop at depth 0, returns inner text.
func (ep *exprParser) readUntil(stop byte) (string, bool) {
	depth := 0
	open := ep.matchingOpen(stop)
	start := ep.pos
	for !ep.eof() {
		c := ep.src[ep.pos]
		if c == '"' || c == '\'' {
			ep.pos = skipStringLitP(ep.src, ep.pos) + 1
			continue
		}
		if open != 0 && c == open {
			depth++
			ep.pos++
			continue
		}
		if c == stop {
			if depth == 0 {
				text := ep.src[start:ep.pos]
				ep.pos++ // consume stop
				return text, true
			}
			depth--
		}
		ep.pos++
	}
	return ep.src[start:ep.pos], false
}

func (ep *exprParser) matchingOpen(close byte) byte {
	switch close {
	case ')':
		return '('
	case ']':
		return '['
	case '}':
		return '{'
	case '>':
		return '<'
	}
	return 0
}

// splitCommaExprs splits a comma-separated expr list at depth 0.
func splitCommaExprs(s string) []Expr {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' || c == '\'' {
			i = skipStringLitP(s, i)
			continue
		}
		switch c {
		case '(', '[', '{', '<':
			depth++
		case ')', ']', '}', '>':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(s[start:]))
	var out []Expr
	for _, p := range parts {
		if p != "" {
			out = append(out, parseExpr(p))
		}
	}
	return out
}

// ---- helpers -----------------------------------------------------------

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '$'
}

func isIdentChar(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

// skipStringLitP returns the index of the closing quote.
func skipStringLitP(s string, i int) int {
	quote := s[i]
	i++
	for i < len(s) {
		switch s[i] {
		case '\\':
			i += 2
			continue
		case quote:
			return i
		}
		i++
	}
	return len(s) - 1
}

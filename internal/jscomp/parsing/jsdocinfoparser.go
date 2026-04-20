// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Scope: Tier-0 through Tier-3 tags only (@type, @param, @return,
// @extends, @implements, @typedef, @template, @constructor, @interface,
// @record). Tag dispatch + type-expression recursive descent mirror the
// Java class exactly: parseAnnotation → parseTypeExpressionAnnotation
// → parseTopLevelTypeExpression → parseTypeExpression →
// parseBasicTypeExpression → parseTypeName / parseRecordType /
// parseFunctionType. Type expressions are serialised back to text so
// the existing TypedScope.ResolveExpr continues to work; the richer
// TypeAST → JSType resolution lands in a later slice.

package parsing

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsDocInfoParser walks a JsDocTokenStream and populates a JSDocInfo.
// Matches the Java class's public surface (Parse + accessors).
type JsDocInfoParser struct {
	stream      *JsDocTokenStream
	info        *rhino.JSDocInfo
	unread      JsDocToken
	hasUnread   bool
	anyRecorded bool

	// collectedNames accumulates every identifier token seen while
	// parsing the current type expression, tagged with its (line, col)
	// inside the comment body. Reset at the top of each
	// parseTypeExpressionAnnotation call. Mirrors how Java's
	// JSDocInfoParser carries the positions along on the Node tree it
	// produces for each type expression.
	collectedNames []rhino.NamedType

	// typeExprStartLine / typeExprStartCol snapshot the position where
	// parseTypeExpressionAnnotation began (just inside the opening
	// `{`). Used when recording a JSDoc parse error so the diagnostic
	// points at the top of the expression.
	typeExprStartLine int
	typeExprStartCol  int
}

// NewJsDocInfoParser mirrors Java's constructor. sourceString is the
// comment body *without* the /** ... */ delimiters.
func NewJsDocInfoParser(sourceString string) *JsDocInfoParser {
	return &JsDocInfoParser{
		stream: NewJsDocTokenStream(sourceString, 0, 0),
		info:   &rhino.JSDocInfo{},
	}
}

// ParseJSDocInfo is the package-level convenience wrapper that mirrors
// how callers in jscomp/parsing/IRFactory use the parser in Java.
// Returns nil when the comment carries no recognised tag, matching the
// previous rhino.ParseJSDoc contract.
func ParseJSDocInfo(body string) *rhino.JSDocInfo {
	p := NewJsDocInfoParser(body)
	p.Parse()
	if !p.anyRecorded {
		return nil
	}
	return p.info
}

// Parse mirrors Java's public boolean parse(). Walks the token stream
// and dispatches annotations to parseAnnotation. Returns true when any
// recognised annotation was recorded (the Java version additionally
// returns false on hard parse errors; we don't surface errors at this
// layer yet).
func (p *JsDocInfoParser) Parse() bool {
	p.skipEOLs()
	token := p.next()
	return p.parseHelperLoop(token)
}

// parseHelperLoop mirrors the Java method of the same name. The loop
// alternates between ANNOTATION dispatch (via parseAnnotation) and
// discarding non-annotation content up to the next EOL / EOC / EOF.
func (p *JsDocInfoParser) parseHelperLoop(token JsDocToken) bool {
	for {
		switch token {
		case JsDocTokenAnnotation:
			token = p.parseAnnotation(token)
		case JsDocTokenEOC, JsDocTokenEOF:
			return p.anyRecorded
		case JsDocTokenEOL:
			p.skipEOLs()
			token = p.next()
		default:
			// Description text / block content — skip to next line.
			token = p.eatUntilEOLIfNotAnnotation()
		}
	}
}

// parseAnnotation mirrors the Java dispatch. The Java version has
// dozens of branches; we handle the tags jsgo's type checker needs.
// Unknown annotations are skipped to end-of-line so they don't derail
// the rest of the comment.
func (p *JsDocInfoParser) parseAnnotation(_ JsDocToken) JsDocToken {
	annotation := strings.ToLower(p.stream.String_())
	switch annotation {
	case "type":
		return p.parseTypeTag()
	case "param":
		return p.parseParamTag()
	case "return", "returns":
		return p.parseReturnTag()
	case "extends":
		return p.parseExtendsTag()
	case "implements":
		return p.parseImplementsTag()
	case "typedef":
		return p.parseTypedefTag()
	case "template":
		return p.parseTemplateTag()
	case "constructor":
		p.info.SetConstructor(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "interface":
		p.info.SetInterface(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "record":
		p.info.SetRecord(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "suppress":
		return p.parseSuppressTag()
	}
	return p.eatUntilEOLIfNotAnnotation()
}

// parseSuppressTag captures the parenthesised list of suppression
// names from `@suppress {name1, name2, …}` (or the brace-only form
// `@suppress {name}`). Mirrors Java's parseSuppressTag — names are
// added to the JSDocInfo's suppression set so downstream passes
// (VariableReferenceCheck etc.) can skip diagnostics gated on them.
func (p *JsDocInfoParser) parseSuppressTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	if t != JsDocTokenLeftCurly {
		return p.eatUntilEOLIfNotAnnotation()
	}
	for {
		t = p.next()
		switch t {
		case JsDocTokenString:
			p.info.AddSuppression(p.stream.String_())
			p.anyRecorded = true
		case JsDocTokenComma:
			continue
		case JsDocTokenRightCurly, JsDocTokenEOF, JsDocTokenEOC:
			return p.eatUntilEOLIfNotAnnotation()
		}
	}
}

// pushBack un-consumes the given token so the next call to next() /
// match() returns it. Used by tag parsers that peek past the expected
// `{` to let the fallback path re-scan the same characters.
func (p *JsDocInfoParser) pushBack(t JsDocToken) {
	p.unread = t
	p.hasUnread = true
}

// bareIdentExpr constructs a JSTypeExpression from a single STRING token
// already produced by the stream. Used by `@extends Foo` / `@implements
// Foo` without braces, which Closure parses silently as a nominal
// reference.
func (p *JsDocInfoParser) bareIdentExpr() *rhino.JSTypeExpression {
	name := p.stream.String_()
	line, col := p.stream.TokenStartLine(), p.stream.TokenStartCol()
	return &rhino.JSTypeExpression{
		TypeName:   name,
		SourceLine: line,
		SourceCol:  col,
		Names: []rhino.NamedType{{
			UnresolvedName: name,
			SourceLine:     line,
			SourceCol:      col,
		}},
	}
}

// parseTypeTag ports the @type branch of the Java parseAnnotation
// switch: parseTypeExpressionAnnotation, record on the builder.
func (p *JsDocInfoParser) parseTypeTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	// `@type` with no following `{` isn't a parse error in Closure's
	// lenient mode — it just omits the type binding. Swallow the rest of
	// the line without recording a diagnostic.
	if t != JsDocTokenLeftCurly {
		p.pushBack(t)
		return p.eatUntilEOLIfNotAnnotation()
	}
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		// Record a type-parse error anchored at the `{` so the resolver
		// can emit JSC_TYPE_PARSE_ERROR at the right JSDoc location.
		p.info.SetType(p.buildParseErrorExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetType(p.buildTypeExpr(expr))
	p.anyRecorded = true
	return next
}

// parseParamTag is the direct port of Java's parseParamTag: a type
// expression in braces, then a parameter name on the same logical
// line. The Java side has additional optional / rest / default-value
// handling; we ignore those (the Tier-0 checker treats them as plain
// names).
func (p *JsDocInfoParser) parseParamTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseParamTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	typeExpr := p.buildTypeExpr(expr)
	// parseParamTypeExpressionAnnotation already advanced to the
	// token immediately following `}`. Java's parser expects that to
	// be the parameter name (STRING), possibly wrapped in optional-
	// param `[ … ]` syntax. Take the name from stream.String_ when
	// the follower is a STRING, otherwise skip to end-of-line.
	if next == JsDocTokenString {
		name := p.stream.String_()
		// Optional-param wrappers arrive as raw identifiers; strip a
		// leading `[` if the tokenizer merged it into the name.
		name = strings.TrimLeft(name, "[")
		if i := strings.IndexAny(name, "=]"); i >= 0 {
			name = name[:i]
		}
		if isJsIdent(name) {
			p.info.AddParam(name, typeExpr)
			p.anyRecorded = true
		}
	}
	return p.eatUntilEOLIfNotAnnotation()
}

func (p *JsDocInfoParser) parseReturnTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	// `@return` without a `{TYPE}` is legal in Closure and simply omits
	// the return-type entry. Don't fire a parse error.
	if t != JsDocTokenLeftCurly {
		p.pushBack(t)
		return p.eatUntilEOLIfNotAnnotation()
	}
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		p.info.SetReturnType(p.buildParseErrorExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetReturnType(p.buildTypeExpr(expr))
	p.anyRecorded = true
	return next
}

func (p *JsDocInfoParser) parseExtendsTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	if t == JsDocTokenString {
		// `@extends Foo` without braces — Closure parses the bare name
		// silently. Capture the identifier as the base type without a
		// parse error.
		p.info.SetBaseType(p.bareIdentExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	if t != JsDocTokenLeftCurly {
		p.pushBack(t)
		return p.eatUntilEOLIfNotAnnotation()
	}
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		p.info.SetBaseType(p.buildParseErrorExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetBaseType(p.buildTypeExpr(expr))
	p.anyRecorded = true
	return next
}

func (p *JsDocInfoParser) parseImplementsTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	if t == JsDocTokenString {
		p.info.AddImplementation(p.bareIdentExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	if t != JsDocTokenLeftCurly {
		p.pushBack(t)
		return p.eatUntilEOLIfNotAnnotation()
	}
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		p.info.AddImplementation(p.buildParseErrorExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.AddImplementation(p.buildTypeExpr(expr))
	p.anyRecorded = true
	return next
}

func (p *JsDocInfoParser) parseTypedefTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		p.info.SetTypedefType(p.buildParseErrorExpr())
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetTypedefType(p.buildTypeExpr(expr))
	p.anyRecorded = true
	return next
}

// buildTypeExpr constructs a JSTypeExpression carrying the captured
// name-position list and the start of the expression inside the JSDoc
// comment body. Mirrors the Java JSDocInfoBuilder that stashes a
// Node.TypeAST root with its position metadata before resolution.
func (p *JsDocInfoParser) buildTypeExpr(expr string) *rhino.JSTypeExpression {
	names := p.collectedNames
	p.collectedNames = nil
	return rhino.NewJSTypeExpressionAt(expr, p.typeExprStartLine, p.typeExprStartCol, names)
}

// buildParseErrorExpr produces a JSTypeExpression flagged as a parse
// failure. The resolver emits JSC_TYPE_PARSE_ERROR at ParseError's
// recorded position. Mirrors the `report(…, TYPE_PARSE_ERROR, …)` path
// in Java's JSDocInfoParser.parseTypeExpressionAnnotation.
func (p *JsDocInfoParser) buildParseErrorExpr() *rhino.JSTypeExpression {
	pe := &rhino.NamedType{SourceLine: p.typeExprStartLine, SourceCol: p.typeExprStartCol}
	p.collectedNames = nil
	return &rhino.JSTypeExpression{
		TypeName:   "",
		SourceLine: p.typeExprStartLine,
		SourceCol:  p.typeExprStartCol,
		ParseError: pe,
	}
}

// parseTemplateTag mirrors the Java method: read the remainder of the
// line, split on commas, register each as a template name.
func (p *JsDocInfoParser) parseTemplateTag() JsDocToken {
	rest := p.stream.GetRemainingJSDocLine()
	for _, part := range strings.Split(rest, ",") {
		part = strings.TrimSpace(part)
		// Templates can carry bounded-generic syntax like `{number} T`
		// that we don't parse yet; strip to the bare ident.
		if idx := strings.LastIndexByte(part, ' '); idx >= 0 {
			part = part[idx+1:]
		}
		part = strings.TrimSpace(part)
		if part != "" && isJsIdent(part) {
			p.info.AddTemplateName(part)
			p.anyRecorded = true
		}
	}
	return p.next()
}

// ---------------------------------------------------------------- //
// Type-expression recursive descent. Method names mirror Java.     //
// ---------------------------------------------------------------- //

// parseTypeExpressionAnnotation: `{ TopLevelTypeExpression }`. Returns
// the expression source + the next unconsumed token.
func (p *JsDocInfoParser) parseTypeExpressionAnnotation(token JsDocToken) (string, JsDocToken, bool) {
	if token != JsDocTokenLeftCurly {
		return "", token, false
	}
	// Record the position of the first character AFTER `{` as the
	// expression origin. Java's JSDocInfoParser uses the same anchor
	// when reporting type-expression parse errors.
	openLine, openCol := p.stream.TokenStartLine(), p.stream.TokenStartCol()
	p.typeExprStartLine, p.typeExprStartCol = openLine, openCol+1
	p.collectedNames = nil
	p.skipEOLs()
	expr, ok := p.parseTopLevelTypeExpression(p.next())
	if !ok {
		return "", p.next(), false
	}
	p.skipEOLs()
	if !p.match(JsDocTokenRightCurly) {
		return "", p.next(), false
	}
	p.next() // consume `}`
	return expr, p.next(), true
}

// parseParamTypeExpressionAnnotation: same as
// parseTypeExpressionAnnotation with optional `=` / `...` modifiers;
// we drop those as a Tier-0 simplification.
func (p *JsDocInfoParser) parseParamTypeExpressionAnnotation(token JsDocToken) (string, JsDocToken, bool) {
	if token != JsDocTokenLeftCurly {
		return "", token, false
	}
	openLine, openCol := p.stream.TokenStartLine(), p.stream.TokenStartCol()
	p.typeExprStartLine, p.typeExprStartCol = openLine, openCol+1
	p.collectedNames = nil
	p.skipEOLs()
	// `{...T}` rest + `{=T}` default. Consume and ignore the modifier.
	if p.match(JsDocTokenIterRest) {
		p.next()
		p.skipEOLs()
	}
	expr, ok := p.parseTopLevelTypeExpression(p.next())
	if !ok {
		return "", p.next(), false
	}
	p.skipEOLs()
	// Optional trailing `=` for default-value params.
	if p.match(JsDocTokenEquals) {
		p.next()
		p.skipEOLs()
	}
	if !p.match(JsDocTokenRightCurly) {
		return "", p.next(), false
	}
	p.next()
	return expr, p.next(), true
}

// parseTopLevelTypeExpression: TypeExpression | TypeExpression '|' ...
// Mirrors Java exactly: top-level unions need no parentheses.
func (p *JsDocInfoParser) parseTopLevelTypeExpression(token JsDocToken) (string, bool) {
	first, ok := p.parseTypeExpression(token)
	if !ok {
		return "", false
	}
	if !p.match(JsDocTokenPipe) {
		return first, true
	}
	alts := []string{first}
	for p.match(JsDocTokenPipe) {
		p.next()
		p.skipEOLs()
		alt, ok := p.parseTypeExpression(p.next())
		if !ok {
			return "", false
		}
		alts = append(alts, alt)
	}
	return strings.Join(alts, "|"), true
}

// parseTypeExpression: BasicTypeExpression with optional leading `?`
// (nullable) or `!` (non-null). Mirrors Java.
func (p *JsDocInfoParser) parseTypeExpression(token JsDocToken) (string, bool) {
	switch token {
	case JsDocTokenQMark:
		// Look ahead: lone `?` with no following type means "unknown".
		if p.matchAny(JsDocTokenRightCurly, JsDocTokenComma, JsDocTokenRightParen,
			JsDocTokenPipe, JsDocTokenRightAngle, JsDocTokenEquals, JsDocTokenEOC, JsDocTokenEOF) {
			return "?", true
		}
		p.skipEOLs()
		inner, ok := p.parseBasicTypeExpression(p.next())
		if !ok {
			return "", false
		}
		return "?" + inner, true
	case JsDocTokenBang:
		p.skipEOLs()
		inner, ok := p.parseBasicTypeExpression(p.next())
		if !ok {
			return "", false
		}
		return "!" + inner, true
	}
	base, ok := p.parseBasicTypeExpression(token)
	if !ok {
		return "", false
	}
	// Postfix `?` / `!`.
	if p.match(JsDocTokenQMark) {
		p.next()
		return "?" + base, true
	}
	if p.match(JsDocTokenBang) {
		p.next()
		return "!" + base, true
	}
	return base, true
}

// parseBasicTypeExpression: NAME / `*` / `(` TYPE `)` / `{` RECORD `}`
// / `function` FunctionType. Java also handles TYPEOF here; we skip.
func (p *JsDocInfoParser) parseBasicTypeExpression(token JsDocToken) (string, bool) {
	switch token {
	case JsDocTokenStar:
		return "*", true
	case JsDocTokenLeftCurly:
		return p.parseRecordType()
	case JsDocTokenLeftParen:
		p.skipEOLs()
		inner, ok := p.parseTopLevelTypeExpression(p.next())
		if !ok {
			return "", false
		}
		p.skipEOLs()
		if !p.match(JsDocTokenRightParen) {
			return "", false
		}
		p.next()
		return "(" + inner + ")", true
	case JsDocTokenString:
		name := p.stream.String_()
		// Snapshot position BEFORE consuming further tokens (e.g.
		// template-arg brackets); parseTypeName advances.
		nameLine := p.stream.TokenStartLine()
		nameCol := p.stream.TokenStartCol()
		if name == "function" {
			return p.parseFunctionType()
		}
		if name == "typeof" {
			// typeof NAME — Closure's TypeOf form. The value-name isn't a
			// type name, so we swallow the following STRING without
			// recording it as a NamedType (avoids false-positive
			// UNRECOGNIZED_TYPE_ERROR). Mirrors Java's JsDocInfoParser
			// behaviour of returning a TypeOf node without registering a
			// NamedType.
			p.skipEOLs()
			tok := p.next()
			if tok != JsDocTokenString {
				return "", false
			}
			return "typeof " + p.stream.String_(), true
		}
		out, ok := p.parseTypeName(name)
		if !ok {
			return "", false
		}
		// Record the leading identifier for post-parse resolution.
		// Template-arg pieces inside `Foo<Bar>` are harvested recursively
		// when parseTypeName descends into parseTopLevelTypeExpression.
		p.collectedNames = append(p.collectedNames, rhino.NamedType{
			UnresolvedName: name,
			SourceLine:     nameLine,
			SourceCol:      nameCol,
		})
		return out, true
	}
	return "", false
}

// parseTypeName: IDENT possibly followed by `.<` generics. Emits the
// reconstructed source text.
func (p *JsDocInfoParser) parseTypeName(first string) (string, bool) {
	name := first
	// Java supports `Foo.<Bar>` and `Foo<Bar>` interchangeably; the
	// token stream already emitted STRING for `Foo` and treats the
	// lone `.` before `<` specially (see collectString). Either way,
	// we now expect an optional `<`.
	if p.match(JsDocTokenLeftAngle) {
		p.next()
		p.skipEOLs()
		var params []string
		for {
			param, ok := p.parseTopLevelTypeExpression(p.next())
			if !ok {
				return "", false
			}
			params = append(params, param)
			p.skipEOLs()
			if !p.match(JsDocTokenComma) {
				break
			}
			p.next()
			p.skipEOLs()
		}
		if !p.match(JsDocTokenRightAngle) {
			return "", false
		}
		p.next()
		name += "<" + strings.Join(params, ",") + ">"
	}
	return name, true
}

// parseRecordType: `{` fieldName (`:` TypeExpression)? (`,` …)* `}`.
// LEFT_CURLY already consumed by the caller.
func (p *JsDocInfoParser) parseRecordType() (string, bool) {
	var parts []string
	for {
		p.skipEOLs()
		if p.match(JsDocTokenRightCurly) {
			break
		}
		tok := p.next()
		if tok != JsDocTokenString {
			return "", false
		}
		field := p.stream.String_()
		p.skipEOLs()
		if p.match(JsDocTokenColon) {
			p.next()
			p.skipEOLs()
			typ, ok := p.parseTypeExpression(p.next())
			if !ok {
				return "", false
			}
			parts = append(parts, field+":"+typ)
		} else {
			parts = append(parts, field)
		}
		p.skipEOLs()
		if !p.match(JsDocTokenComma) {
			break
		}
		p.next()
	}
	p.skipEOLs()
	if !p.match(JsDocTokenRightCurly) {
		return "", false
	}
	p.next()
	return "{" + strings.Join(parts, ",") + "}", true
}

// parseFunctionType: `function(params): returnType`. Shrunk from
// Java's more complete version: we skip `this:` / `new:` prefixes and
// treat every param as a positional type.
func (p *JsDocInfoParser) parseFunctionType() (string, bool) {
	if !p.match(JsDocTokenLeftParen) {
		return "Function", true // bare `function`
	}
	p.next()
	var params []string
	for {
		p.skipEOLs()
		if p.match(JsDocTokenRightParen) {
			break
		}
		tok := p.next()
		// Allow leading `this:` / `new:` qualifier; swallow it.
		if tok == JsDocTokenString {
			name := p.stream.String_()
			if (name == "this" || name == "new") && p.match(JsDocTokenColon) {
				p.next()
				p.skipEOLs()
				tok = p.next()
			}
		}
		// Allow a leading `...` (varargs). Mirrors Java's handling in
		// parseFunctionType. The popped collectedNames entry (if any) for
		// the type that follows is unaffected.
		if tok == JsDocTokenIterRest {
			p.skipEOLs()
			tok = p.next()
			// `function(...?)` — bare `?` after `...` means "any".
			if tok == JsDocTokenRightParen {
				params = append(params, "?")
				break
			}
		}
		expr, ok := p.parseTypeExpression(tok)
		if !ok {
			return "", false
		}
		// Optional-param marker `=` after the type (e.g. `number=`).
		if p.match(JsDocTokenEquals) {
			p.next()
			expr += "="
		}
		params = append(params, expr)
		p.skipEOLs()
		if p.match(JsDocTokenComma) {
			p.next()
		}
	}
	p.next() // consume `)`
	p.skipEOLs()
	result := ""
	if p.match(JsDocTokenColon) {
		p.next()
		p.skipEOLs()
		res, ok := p.parseTypeExpression(p.next())
		if !ok {
			return "", false
		}
		result = res
	}
	out := "function(" + strings.Join(params, ",") + ")"
	if result != "" {
		out += ":" + result
	}
	return out, true
}

// ---------------------------------------------------------------- //
// Small helpers: match / next / skipEOLs / line-eaters.            //
// ---------------------------------------------------------------- //

func (p *JsDocInfoParser) match(tok JsDocToken) bool {
	p.unread = p.next()
	p.hasUnread = true
	return p.unread == tok
}

func (p *JsDocInfoParser) matchAny(toks ...JsDocToken) bool {
	p.unread = p.next()
	p.hasUnread = true
	for _, t := range toks {
		if p.unread == t {
			return true
		}
	}
	return false
}

func (p *JsDocInfoParser) next() JsDocToken {
	if p.hasUnread {
		p.hasUnread = false
		return p.unread
	}
	return p.stream.GetJsDocToken()
}

func (p *JsDocInfoParser) skipEOLs() {
	for p.match(JsDocTokenEOL) {
		p.next()
		if p.match(JsDocTokenStar) {
			p.next()
		}
	}
}

// eatUntilEOLIfNotAnnotation consumes tokens until EOL / EOC / EOF or
// until a fresh ANNOTATION starts. Returns the token that terminated
// the scan so the caller can re-dispatch.
func (p *JsDocInfoParser) eatUntilEOLIfNotAnnotation() JsDocToken {
	for {
		t := p.next()
		switch t {
		case JsDocTokenEOL, JsDocTokenEOC, JsDocTokenEOF, JsDocTokenAnnotation:
			return t
		}
	}
}

// readParamName consumes the identifier that follows a `@param {T}`
// type annotation, tolerating `[name]` and `[name=default]`.
func (p *JsDocInfoParser) readParamName() string {
	// Consume whitespace (still in the stream's char layer).
	rest := p.stream.GetRemainingJSDocLine()
	rest = strings.TrimLeft(rest, " \t*[")
	end := 0
	for end < len(rest) {
		c := rest[end]
		if isJsIdentChar(c) {
			end++
			continue
		}
		break
	}
	return rest[:end]
}

// isJsIdent / isJsIdentChar: duplicate the small helpers the Java side
// delegates to TokenUtil.
func isJsIdent(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isJsIdentChar(s[i]) {
			return false
		}
	}
	return len(s) > 0
}

func isJsIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '$'
}

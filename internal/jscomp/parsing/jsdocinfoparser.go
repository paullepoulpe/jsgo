// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415 (subset)
//
// Scope: Tier-0 through Tier-3 tags (@type, @param, @return, @extends,
// @implements, @typedef, @template, @constructor, @interface, @record,
// @abstract) plus the side-effect markers CheckSideEffects consults
// (@nosideeffects / @noside). Tag dispatch + type-expression recursive
// descent mirror the Java class exactly: parseAnnotation →
// parseTypeExpressionAnnotation → parseTopLevelTypeExpression →
// parseTypeExpression → parseBasicTypeExpression → parseTypeName /
// parseRecordType / parseFunctionType. Type expressions are serialised
// back to text so the existing TypedScope.ResolveExpr continues to
// work; the richer TypeAST → JSType resolution lands in a later slice.

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

	// pendingDupField* is the parallel channel for deferred parse
	// errors that survive a successful parse. parseRecordType
	// populates these when a record body contains a duplicate field —
	// the record itself still parses so downstream
	// UNRECOGNIZED_TYPE_ERROR emission for unknown field types
	// continues to fire; buildTypeExpr reads the pending fields after
	// receiving the ok==true return and stamps them onto the built
	// JSTypeExpression as a ParseError marker. Mirrors Java's
	// parseFieldTypeList which emits a report() alongside the
	// resolved TypeAST instead of aborting the whole parse.
	pendingDupFieldLine int
	pendingDupFieldCol  int
	pendingDupFieldMsg  string
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
	case "struct":
		p.info.SetStruct(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "dict":
		p.info.SetDict(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "enum":
		return p.parseEnumTag()
	case "abstract":
		p.info.SetAbstract(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "unrestricted":
		// `@unrestricted` opts out of both dict and struct discipline.
		// No explicit flag today; just record that we saw a tag so the
		// parser doesn't discard the info.
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "private":
		p.info.SetVisibility(rhino.VisibilityPrivate)
		p.anyRecorded = true
		return p.parseVisibilityRest()
	case "protected":
		p.info.SetVisibility(rhino.VisibilityProtected)
		p.anyRecorded = true
		return p.parseVisibilityRest()
	case "public":
		p.info.SetVisibility(rhino.VisibilityPublic)
		p.anyRecorded = true
		return p.parseVisibilityRest()
	case "package":
		p.info.SetVisibility(rhino.VisibilityPackage)
		p.anyRecorded = true
		return p.parseVisibilityRest()
	case "suppress":
		return p.parseSuppressTag()
	case "nosideeffects", "noside":
		// `@nosideeffects` / `@noside` mark the function as
		// observably pure. Mirrors the equivalent branch in Java's
		// JsDocInfoParser.parseAnnotation.
		p.info.SetNoSideEffects(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	return p.eatUntilEOLIfNotAnnotation()
}

// parseVisibilityRest mirrors Java's parseAnnotation branch for
// `@private`/`@protected`/`@public`/`@package`: an optional inline
// type expression `{T}` may follow (e.g. `@private {number}`), which
// we capture into `@type` so the downstream checker sees the variable
// declaration's type.
func (p *JsDocInfoParser) parseVisibilityRest() JsDocToken {
	p.skipEOLs()
	if !p.match(JsDocTokenLeftCurly) {
		return p.eatUntilEOLIfNotAnnotation()
	}
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	if !p.info.HasType() {
		p.info.SetType(rhino.NewJSTypeExpression(expr))
	}
	return next
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
// line. Modifier handling (`{T=}` optional / `{...T}` rest) matches
// JSDocInfoParser.java#parseParamTag — the flags are attached to the
// stored JSTypeExpression so downstream passes (CheckArity) can see
// them per-parameter.
func (p *JsDocInfoParser) parseParamTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, isRest, isOptional, next, ok := p.parseParamTypeAnnotationFull(t)
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
		// Optional-param `[name]` / `[name=default]` wrappers mark the
		// param as optional independently of the `{T=}` form — mirrors
		// JSDocInfoParser's handling of the bracketed name syntax.
		if strings.HasPrefix(name, "[") {
			isOptional = true
		}
		name = strings.TrimLeft(name, "[")
		if i := strings.IndexAny(name, "=]"); i >= 0 {
			name = name[:i]
		}
		if isJsIdent(name) {
			p.info.AddParam(name, typeExpr)
			if isRest {
				p.info.MarkRestParam(name)
			}
			if isOptional {
				p.info.MarkOptionalParam(name)
			}
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

// parseEnumTag handles `@enum [{T}]`. The braced type argument is
// optional — bare `@enum` implies an `@enum {number}` in Java's
// EnumType defaulting path. Either way we flip the `enumMarked`
// flag on the JSDocInfo so TypedScopeCreator's enum branch fires.
func (p *JsDocInfoParser) parseEnumTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	// If the next token isn't `{`, no type argument — still mark
	// the enum and keep parsing the rest of the annotation block.
	if t != JsDocTokenLeftCurly {
		p.info.MarkEnum()
		p.anyRecorded = true
		// Put the token back and let the outer loop re-dispatch.
		p.unread = t
		p.hasUnread = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		p.info.MarkEnum()
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetEnumParameterType(rhino.NewJSTypeExpression(expr))
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
	te := rhino.NewJSTypeExpressionAt(expr, p.typeExprStartLine, p.typeExprStartCol, names)
	// Surface any deferred duplicate-record-field error captured during
	// parseRecordType. The expression keeps its Names (so downstream
	// UNRECOGNIZED_TYPE_ERROR emission continues for preceding fields)
	// and simultaneously carries a ParseError with the Java-matching
	// message + position.
	if p.pendingDupFieldMsg != "" {
		te.ParseError = &rhino.NamedType{
			SourceLine: p.pendingDupFieldLine,
			SourceCol:  p.pendingDupFieldCol,
			Message:    p.pendingDupFieldMsg,
		}
		p.pendingDupFieldMsg = ""
	}
	return te
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
// we drop those as a Tier-0 simplification. Returns the expr + the
// post-annotation token.
func (p *JsDocInfoParser) parseParamTypeExpressionAnnotation(token JsDocToken) (string, JsDocToken, bool) {
	expr, _, next, ok := p.parseParamTypeExpressionAnnotationWithRest(token)
	return expr, next, ok
}

// parseParamTypeExpressionAnnotationWithRest is the modifier-aware
// variant — callers who need to know whether the brace carried
// `{...T}` (rest) or `{T=}` (optional) use this so the param's
// flags thread through to the JSDocInfo.
func (p *JsDocInfoParser) parseParamTypeExpressionAnnotationWithRest(token JsDocToken) (string, bool, JsDocToken, bool) {
	expr, isRest, _, next, ok := p.parseParamTypeAnnotationFull(token)
	return expr, isRest, next, ok
}

// parseParamTypeAnnotationFull returns (expr, isRest, isOptional,
// nextToken, ok). Splits the two brace-position modifiers Java's
// JsDocInfoParser tracks.
func (p *JsDocInfoParser) parseParamTypeAnnotationFull(token JsDocToken) (string, bool, bool, JsDocToken, bool) {
	if token != JsDocTokenLeftCurly {
		return "", false, false, token, false
	}
	openLine, openCol := p.stream.TokenStartLine(), p.stream.TokenStartCol()
	p.typeExprStartLine, p.typeExprStartCol = openLine, openCol+1
	p.collectedNames = nil
	p.skipEOLs()
	isRest := false
	isOptional := false
	if p.match(JsDocTokenIterRest) {
		p.next()
		p.skipEOLs()
		isRest = true
	}
	expr, ok := p.parseTopLevelTypeExpression(p.next())
	if !ok {
		return "", isRest, isOptional, p.next(), false
	}
	p.skipEOLs()
	if p.match(JsDocTokenEquals) {
		p.next()
		p.skipEOLs()
		isOptional = true
	}
	if !p.match(JsDocTokenRightCurly) {
		return "", isRest, isOptional, p.next(), false
	}
	p.next()
	return expr, isRest, isOptional, p.next(), true
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
//
// Duplicate field names surface as a deferred JSC_TYPE_PARSE_ERROR
// anchored at the closing `}` position — mirrors Java's
// parseFieldTypeList report() for `Duplicate record field <name>`.
// The duplicate's body (and anything after it) is swallowed without
// contributing to collectedNames so the resolver only emits unknown-
// type errors for names Java would also report.
func (p *JsDocInfoParser) parseRecordType() (string, bool) {
	var parts []string
	seen := map[string]bool{}
	dupField := ""
	dupLine, dupCol := 0, 0
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
		isDup := seen[field] && dupField == ""
		seen[field] = true
		if isDup {
			dupField = field
			// Discard the duplicate's body (and any trailing fields)
			// without contributing Names to collectedNames — Java's
			// parseFieldTypeList aborts on duplicate detection, so the
			// resolver should only flag unknown types from fields that
			// appeared *before* the duplicate.
			saved := append([]rhino.NamedType(nil), p.collectedNames...)
			p.skipEOLs()
			if p.match(JsDocTokenColon) {
				p.next()
				p.skipEOLs()
				_, _ = p.parseTypeExpression(p.next())
			}
			p.skipEOLs()
			for p.match(JsDocTokenComma) {
				p.next()
				p.skipEOLs()
				if p.match(JsDocTokenRightCurly) {
					break
				}
				if p.next() != JsDocTokenString {
					break
				}
				p.skipEOLs()
				if p.match(JsDocTokenColon) {
					p.next()
					p.skipEOLs()
					_, _ = p.parseTypeExpression(p.next())
				}
				p.skipEOLs()
			}
			p.collectedNames = saved
			// Peek at the closing `}` so its (line, col) becomes the
			// most recent TokenStart — that's where Java's report()
			// anchors the Duplicate record field diagnostic.
			if p.match(JsDocTokenRightCurly) {
				dupLine = p.stream.TokenStartLine()
				dupCol = p.stream.TokenStartCol()
			}
			break
		}
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
	if dupField != "" {
		p.pendingDupFieldLine = dupLine
		p.pendingDupFieldCol = dupCol
		p.pendingDupFieldMsg = "Bad type annotation. Duplicate record field " + dupField + ". See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information."
	}
	return "{" + strings.Join(parts, ",") + "}", true
}

// parseFunctionType: `function(params): returnType`. Mirrors Java's
// `parseParametersType` branch: supports `this:T` / `new:T` qualifier
// on the first param, `...T` leading-rest on the last param, and
// `T=` trailing-optional per-param. The textual shape is preserved
// — rest params gain a leading `...` and optional params a trailing
// `=` in the serialised form — so `scope.ResolveExpr` can round-trip
// the string back into a *jstype.FunctionType with accurate arity.
func (p *JsDocInfoParser) parseFunctionType() (string, bool) {
	if !p.match(JsDocTokenLeftParen) {
		return "Function", true // bare `function`
	}
	p.next()
	var params []string
	first := true
	sawRest := false
	isCtorShape := false
	for {
		p.skipEOLs()
		if p.match(JsDocTokenRightParen) {
			break
		}
		// Leading `...T` rest marker inside the paren list. Mirrors
		// Java's `parseParametersType` ITER_REST branch.
		isRest := false
		if p.match(JsDocTokenIterRest) {
			p.next()
			p.skipEOLs()
			isRest = true
			sawRest = true
		}
		tok := p.next()
		// Allow leading `this:` / `new:` qualifier on the first param
		// slot. Java tracks each separately on the ArrowType; for our
		// textual form we swallow `this:` and propagate `new:` so the
		// resolved FunctionType can flip its IsConstructor bit.
		skipFirstParam := false
		if first && !isRest && tok == JsDocTokenString {
			name := p.stream.String_()
			if (name == "this" || name == "new") && p.match(JsDocTokenColon) {
				p.next()
				p.skipEOLs()
				if name == "new" {
					isCtorShape = true
					// `new:T` is not a real parameter — swallow the
					// type that follows and move on to the next slot.
					_, ok := p.parseTypeExpression(p.next())
					if !ok {
						return "", false
					}
					p.skipEOLs()
					skipFirstParam = true
				} else {
					tok = p.next()
				}
			}
		}
		first = false
		if skipFirstParam {
			if p.match(JsDocTokenComma) {
				p.next()
			}
			continue
		}
		expr, ok := p.parseTypeExpression(tok)
		if !ok {
			return "", false
		}
		p.skipEOLs()
		// Per-param trailing `=` marks the param optional. Java's
		// parser wraps the type in Token.EQUALS; we keep the `=`
		// suffix in the text form.
		isOptional := false
		if p.match(JsDocTokenEquals) {
			p.next()
			p.skipEOLs()
			isOptional = true
		}
		switch {
		case isRest:
			params = append(params, "..."+expr)
		case isOptional:
			params = append(params, expr+"=")
		default:
			params = append(params, expr)
		}
		p.skipEOLs()
		if sawRest {
			// After a rest param, no more params are allowed.
			break
		}
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
	inner := strings.Join(params, ",")
	if isCtorShape {
		// Preserve the `new:` marker in the textual form so scope
		// resolution can flip IsConstructor on the resulting
		// FunctionType. Mirrors how Java's ArrowType tracks the
		// bound `new` this-type separately from the param list.
		if inner != "" {
			inner = "new:," + inner
		} else {
			inner = "new:"
		}
	}
	out := "function(" + inner + ")"
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

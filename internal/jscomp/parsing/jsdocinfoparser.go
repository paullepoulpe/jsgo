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

// parseTypeTag ports the @type branch of the Java parseAnnotation
// switch: parseTypeExpressionAnnotation, record on the builder.
func (p *JsDocInfoParser) parseTypeTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetType(rhino.NewJSTypeExpression(expr))
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
			p.info.AddParam(name, rhino.NewJSTypeExpression(expr))
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
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetReturnType(rhino.NewJSTypeExpression(expr))
	p.anyRecorded = true
	return next
}

func (p *JsDocInfoParser) parseExtendsTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetBaseType(rhino.NewJSTypeExpression(expr))
	p.anyRecorded = true
	return next
}

func (p *JsDocInfoParser) parseImplementsTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.AddImplementation(rhino.NewJSTypeExpression(expr))
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
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.info.SetTypedefType(rhino.NewJSTypeExpression(expr))
	p.anyRecorded = true
	return next
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
		if name == "function" {
			return p.parseFunctionType()
		}
		return p.parseTypeName(name)
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

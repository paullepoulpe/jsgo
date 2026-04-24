// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Scope: Tier-0 through Tier-3 tags only (@type, @param, @return,
// @extends, @implements, @typedef, @template, @constructor, @interface,
// @record, @abstract, @struct, @dict, @private, @protected, @public,
// @nosideeffects). Tag dispatch + type-expression recursive descent
// mirror the Java class exactly: parseAnnotation →
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
// startLine is the 1-indexed file line of the `/**` comment opener; it is
// used to compute absolute positions for JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES
// (mirrors JsDocInfoParser.java's lineno/charno capture at parseAnnotation entry).
// Returns nil when the comment carries no recognised tag, matching the
// previous rhino.ParseJSDoc contract.
func ParseJSDocInfo(body string, startLine int) *rhino.JSDocInfo {
	p := NewJsDocInfoParser(body)
	p.Parse()
	if !p.anyRecorded {
		return nil
	}
	// Scan the raw body for bounded-generic @template {B} T annotations and
	// record their absolute positions. This mirrors the
	// addParserWarning(Msg.JSDOC_TEMPLATE_BOUNDEDGENERICS_USED, lineno, charno)
	// call in JsDocInfoParser.java's TEMPLATE case — Java captures lineno/charno
	// from stream.getLineno()/getCharno() at parseAnnotation entry, which gives
	// the position of the '@' character.
	// TODO: wire proper lineno/charno through JsDocTokenStream instead of scanning.
	recordBoundedTemplatePositions(p.info, body, startLine)
	recordInvalidParamPositions(p.info, body, startLine)
	return p.info
}

// recordBoundedTemplatePositions scans the raw JSDoc body for every
// `@template` annotation that carries a bound (`{`). For each one it records
// the absolute (line, col) of the `@` character into info.
// startLine is the 1-indexed file line of the `/**` opening delimiter.
func recordBoundedTemplatePositions(info *rhino.JSDocInfo, body string, startLine int) {
	line := startLine // current absolute line (1-indexed)
	i := 0
	for i < len(body) {
		// Track newlines to maintain absolute line number.
		if body[i] == '\n' {
			line++
			i++
			continue
		}
		// Look for `@template` followed (after optional spaces) by `{`.
		if body[i] != '@' {
			i++
			continue
		}
		// Compute column: scan back to start of current line.
		lineStart := strings.LastIndex(body[:i], "\n") + 1 // byte offset of BOL
		col := i - lineStart                               // 0-indexed column of '@'
		// Match "template" (case-insensitive per JsDocInfoParser).
		rest := body[i+1:]
		if !strings.HasPrefix(strings.ToLower(rest), "template") {
			i++
			continue
		}
		afterKw := rest[len("template"):]
		// Skip optional whitespace between "@template" and the bound "{".
		j := 0
		for j < len(afterKw) && (afterKw[j] == ' ' || afterKw[j] == '\t') {
			j++
		}
		if j < len(afterKw) && afterKw[j] == '{' {
			// This is a bounded @template {B} T — record position of '@'.
			info.AddBoundedTemplatePosition(rhino.BoundedTemplatePosition{
				Line: line,
				Col:  col,
			})
		}
		i++ // advance past '@'
	}
}

// recordInvalidParamPositions fills in absolute (line, col) for each
// InvalidParamPosition stored during parseParamTag, by scanning the raw body.
func recordInvalidParamPositions(info *rhino.JSDocInfo, body string, startLine int) {
	invalids := info.InvalidParamPositions()
	if len(invalids) == 0 {
		return
	}
	need := make(map[string]int, len(invalids))
	for i, iv := range invalids {
		need[iv.Name] = i
	}
	line, i := startLine, 0
	for i < len(body) {
		if body[i] == '\n' { line++; i++; continue }
		if body[i] != '@' || !strings.HasPrefix(body[i+1:], "param") { i++; continue }
		after := body[i+1+len("param"):]
		j := 0
		for j < len(after) && (after[j] == ' ' || after[j] == '\t') { j++ }
		if j >= len(after) || after[j] != '{' { i++; continue }
		j++
		for depth := 1; j < len(after) && depth > 0; j++ {
			if after[j] == '{' { depth++ } else if after[j] == '}' { depth-- }
		}
		for j < len(after) && (after[j] == ' ' || after[j] == '\t') { j++ }
		k := j
		for k < len(after) && after[k] != ' ' && after[k] != '\t' && after[k] != '\n' && after[k] != '*' { k++ }
		name := strings.TrimLeft(after[j:k], "[")
		if idx := strings.IndexAny(name, "=]"); idx >= 0 { name = name[:idx] }
		if idx, ok := need[name]; ok {
			lineStart := strings.LastIndex(body[:i], "\n") + 1
			info.UpdateInvalidParamPosition(idx, line, i+1+len("param")+j-lineStart)
			delete(need, name)
		}
		i++
	}
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
	case "abstract":
		// Mirrors Java's ABSTRACT case: record the flag; no type argument.
		p.info.SetAbstract(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
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
	case "private":
		// Mirrors Java's PRIVATE case: record visibility; optional type
		// argument on the same tag is ignored at this layer (parse-and-store
		// only — no inline-type side-channel needed for Tier-0..3).
		p.info.RecordVisibility(rhino.VisibilityPrivate)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "protected":
		p.info.RecordVisibility(rhino.VisibilityProtected)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "public":
		p.info.RecordVisibility(rhino.VisibilityPublic)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "nosideeffects":
		// Mirrors Java's NO_SIDE_EFFECTS case.
		p.info.SetNoSideEffects(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "deprecated":
		// Mirrors Java's DEPRECATED case: recordDeprecated() then
		// extractMultilineTextualBlock to capture the reason.
		// We read the rest of the current line as the reason text and
		// trim leading/trailing whitespace (Java's extractMultilineTextualBlock
		// with WhitespaceOption.TRIM trims each line before joining).
		reason := strings.TrimSpace(p.stream.GetRemainingJSDocLine())
		p.info.SetDeprecated(reason)
		p.anyRecorded = true
		return p.next()
	case "final":
		// Mirrors Java's FINAL case: record the @final flag; no argument.
		p.info.SetFinal(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "package":
		// @package visibility — analogous to @private/@protected/@public.
		p.info.RecordVisibility(rhino.VisibilityPackage)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "enum":
		return p.parseEnumTag()
	case "suppress":
		return p.parseSuppressTag()
	case "this":
		return p.parseThisTag()
	case "override":
		p.info.SetOverride(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "lends":
		return p.parseLendsTag()
	case "const":
		p.info.SetConst(true)
		p.anyRecorded = true
		// Java's parser calls recordType(type) when a `{T}` follows @const
		// (CONSTANT is in isAlternateTypeAnnotation; lookAheadForType()→true
		// → parseAndRecordTypeNode → recordType). Mirror that: attempt to
		// parse an optional `{T}` and store it as the declared type so that
		// @const {string} behaves like @type {string}.
		p.skipEOLs()
		t := p.next()
		if expr, next, ok := p.parseTypeExpressionAnnotation(t); ok {
			p.info.SetType(rhino.NewJSTypeExpression(expr))
			return next
		}
		// No `{T}` present — t is the first non-EOL token (already consumed).
		// Return it directly so the outer loop re-dispatches it.
		return t
	case "implicitcast":
		// Mirrors Java's IMPLICIT_CAST case: recordImplicitCast().
		p.info.SetBit(rhino.BitImplicitCast, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "define":
		// Mirrors Java's DEFINE case: recordDefineType() with optional {T}.
		p.info.SetBit(rhino.BitDefine, true)
		p.anyRecorded = true
		p.skipEOLs()
		t := p.next()
		if expr, next, ok := p.parseTypeExpressionAnnotation(t); ok {
			p.info.SetType(rhino.NewJSTypeExpression(expr))
			return next
		}
		return t
	case "nocollapse":
		// Mirrors Java's NO_COLLAPSE case: recordNoCollapse().
		p.info.SetBit(rhino.BitNoCollapse, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "closureprimitive":
		// Mirrors Java's CLOSURE_PRIMITIVE case: skipEOLs then
		// delegate to parseClosurePrimitiveTag.
		p.skipEOLs()
		return p.parseClosurePrimitiveTag(p.next())
	case "usedviadotconstructor":
		// Mirrors Java's USED_VIA_DOT_CONSTRUCTOR case.
		p.info.SetBit(rhino.BitUsedViaDotCtor, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "typesummary":
		// Mirrors Java's TYPE_SUMMARY case: recordTypeSummary().
		p.info.SetBit(rhino.BitTypeSummary, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "nginject":
		// Mirrors Java's NG_INJECT case: recordNgInject(true).
		p.info.SetBit(rhino.BitNgInject, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "wizaction":
		// Mirrors Java's WIZACTION case: recordWizaction().
		p.info.SetBit(rhino.BitWizAction, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "polymerbehavior":
		// Mirrors Java's POLYMER_BEHAVIOR case: recordPolymerBehavior().
		p.info.SetBit(rhino.BitPolymerBehavior, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "customelement":
		// Mirrors Java's CUSTOM_ELEMENT case: recordCustomElement().
		p.info.SetBit(rhino.BitCustomElement, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "throws":
		// Mirrors Java's THROWS case: optionally parse a {type} expression,
		// then capture the remaining text as the throws annotation.
		return p.parseThrowsTag()
	case "modifies":
		// Mirrors Java's MODIFIES case: delegate to parseModifiesTag.
		return p.parseModifiesTag(p.next())
	case "polymer":
		// Mirrors Java's POLYMER case: recordPolymer().
		p.info.SetBit(rhino.BitPolymer, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "mixinclass":
		// Mirrors Java's MIXIN_CLASS case: recordMixinClass().
		p.info.SetBit(rhino.BitMixinClass, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "mixinfunction":
		// Mirrors Java's MIXIN_FUNCTION case: recordMixinFunction().
		p.info.SetBit(rhino.BitMixinFunction, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "externs":
		// Mirrors Java's EXTERNS case: recordExterns().
		p.info.SetBit(rhino.BitExterns, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "nocompile":
		// Mirrors Java's NO_COMPILE case: recordNoCompile().
		p.info.SetBit(rhino.BitNoCompile, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "nocoverage":
		// Mirrors Java's NO_COVERAGE case: recordNoCoverage().
		p.info.SetBit(rhino.BitNoCoverage, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "nodts":
		// Mirrors Java's NO_DTS case: recordNoDts().
		p.info.SetBit(rhino.BitNoDts, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "fileoverview":
		// Mirrors Java's FILE_OVERVIEW case: recordFileOverview().
		p.info.SetBit(rhino.BitFileoverview, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "pureorbreakmycode":
		// Mirrors Java's PURE_OR_BREAK_MY_CODE case: recordPureOrBreakMyCode().
		p.info.SetBit(rhino.BitPureOrBreakMyCode, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "collapsibleorbreakmycode":
		// Mirrors Java's COLLAPSIBLE_OR_BREAK_MY_CODE case.
		p.info.SetBit(rhino.BitCollapsibleOrBreak, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "export":
		// Mirrors Java's EXPORT case: recordExport().
		p.info.SetBit(rhino.BitExport, true)
		p.anyRecorded = true
		// @export may optionally carry a type (like @const); skip it.
		p.skipEOLs()
		t := p.next()
		if _, next, ok := p.parseTypeExpressionAnnotation(t); ok {
			return next
		}
		return t
	case "unrestricted":
		// Mirrors Java's UNRESTRICTED case: recordUnrestricted().
		p.info.SetBit(rhino.BitUnrestricted, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "wizcallback":
		// Mirrors Java's WIZCALLBACK case.
		p.info.SetBit(rhino.BitWizCallback, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "requireinlining":
		// Mirrors Java's REQUIRE_INLINING case.
		p.info.SetBit(rhino.BitRequireInlining, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "encourageinlining":
		// Mirrors Java's ENCOURAGE_INLINING case.
		p.info.SetBit(rhino.BitEncourageInlining, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "providegoog":
		// Mirrors Java's PROVIDE_GOOG case.
		p.info.SetBit(rhino.BitProvideGoog, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "providealreadyprovided":
		// Mirrors Java's PROVIDE_ALREADY_PROVIDED case.
		p.info.SetBit(rhino.BitProvideAlreadyProvided, true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "idgenerator":
		// Mirrors Java's IDGENERATOR case: delegate to parseIdGeneratorTag.
		p.skipEOLs()
		return p.parseIdGeneratorTag(p.next())
	case "see":
		// Mirrors Java's SEE case: skip to end of line.
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "author":
		// Mirrors Java's AUTHOR case: skip to end of line.
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "desc":
		// Mirrors Java's DESC case: skip to end of line (documentation not collected here).
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	case "inheritdoc":
		// Mirrors Java's INHERIT_DOC case — treat like @override.
		p.info.SetOverride(true)
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
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

// parseClosurePrimitiveTag ports Java's parseClosurePrimitiveTag: expects
// `{id}` and records the closure primitive id. Mirrors Java exactly:
// token must be LEFT_CURLY; next must be STRING; then RIGHT_CURLY.
func (p *JsDocInfoParser) parseClosurePrimitiveTag(token JsDocToken) JsDocToken {
	if token != JsDocTokenLeftCurly {
		return p.eatUntilEOLIfNotAnnotation()
	}
	// Peek at the next token; if it's a STRING, record the id.
	inner := p.next()
	if inner != JsDocTokenString {
		return p.eatUntilEOLIfNotAnnotation()
	}
	id := p.stream.String_()
	p.info.SetClosurePrimitiveId(id)
	p.anyRecorded = true
	// Expect closing '}'.
	if !p.match(JsDocTokenRightCurly) {
		return p.eatUntilEOLIfNotAnnotation()
	}
	p.next() // consume '}'
	return p.eatUntilEOLIfNotAnnotation()
}

// parseThrowsTag ports Java's THROWS case. Java's implementation reads
// free text via extractMultilineTextualBlock (no type-expression parse).
// We extend it: if a `{T}` follows @throws, we parse the type expression
// and record it; otherwise we record an empty sentinel so HasThrows() is true.
// Mirrors Java's lookAheadForAnnotation check: if the next token is
// ANNOTATION/EOC/EOF we skip recording (no actual throw clause follows).
func (p *JsDocInfoParser) parseThrowsTag() JsDocToken {
	// Peek at the first non-EOL token.
	p.skipEOLs()
	// Use match to peek without consuming.
	if p.matchAny(JsDocTokenAnnotation, JsDocTokenEOC, JsDocTokenEOF) {
		// lookAheadForAnnotation path: no throw clause.
		return p.next() // returns the peeked annotation token
	}
	// Check for a type expression `{T}`.
	if p.match(JsDocTokenLeftCurly) {
		t := p.next() // consume `{`
		if expr, next, ok := p.parseTypeExpressionAnnotation(t); ok {
			p.info.AddThrowsAnnotation(expr)
			p.anyRecorded = true
			return next
		}
		// Type parse failed — record sentinel and eat to EOL.
		p.info.AddThrowsAnnotation("")
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	// Free-text path: read remaining line text. The token peek left a
	// non-`{` token in unread; GetRemainingJSDocLine works at the char
	// level so we read what's left on the line after the peeked token.
	// The peeked token is already consumed from the stream — clear unread
	// and read the rest of the line via GetRemainingJSDocLine.
	_ = p.next() // consume the peeked token
	text := strings.TrimSpace(p.stream.GetRemainingJSDocLine())
	if text != "" {
		p.info.AddThrowsAnnotation(text)
	} else {
		p.info.AddThrowsAnnotation("")
	}
	p.anyRecorded = true
	return p.next()
}

// parseModifiesTag ports Java's parseModifiesTag: `@modifies{this|arguments}`.
// Mirrors Java: expects LEFT_CURLY, then STRING items separated by PIPE,
// then RIGHT_CURLY. Calls JSDocInfo.SetModifies with the collected set.
func (p *JsDocInfoParser) parseModifiesTag(token JsDocToken) JsDocToken {
	if token != JsDocTokenLeftCurly {
		return p.eatUntilEOLIfNotAnnotation()
	}
	modifies := map[string]bool{}
	for {
		if p.match(JsDocTokenString) {
			name := p.stream.String_()
			modifies[name] = true
			token = p.next()
		} else {
			return p.eatUntilEOLIfNotAnnotation()
		}
		if p.match(JsDocTokenPipe) {
			token = p.next()
		} else {
			break
		}
	}
	if !p.match(JsDocTokenRightCurly) {
		return p.eatUntilEOLIfNotAnnotation()
	}
	token = p.next() // consume '}'
	p.info.SetModifies(modifies)
	p.anyRecorded = true
	return token
}

// parseThisTag ports the @this branch of the Java parseAnnotation:
// consumes the optional `{Type}` argument and marks the JSDocInfo as
// having a @this annotation. Mirrors Java's JsDocInfoParser handling
// of @this — which simply records it as hasThisType (the actual type
// expression is not resolved by CheckGlobalThis, only its presence matters).
func (p *JsDocInfoParser) parseThisTag() JsDocToken {
	p.info.SetHasThisType(true)
	p.anyRecorded = true
	return p.eatUntilEOLIfNotAnnotation()
}

// parseLendsTag ports the @lends branch of the Java parseAnnotation:
// consumes the `{Name}` argument and records the lend target on JSDocInfo.
// Mirrors Java's JsDocInfoParser.parseLendsTag — CheckGlobalThis uses
// getLendsName().getRoot().getString() to detect prototype lending.
func (p *JsDocInfoParser) parseLendsTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	if t != JsDocTokenLeftCurly {
		return p.eatUntilEOLIfNotAnnotation()
	}
	typStr, ok := p.parseTopLevelTypeExpression(p.next())
	if ok && typStr != "" {
		p.info.SetLendsName(typStr)
		p.anyRecorded = true
	}
	return p.eatUntilEOLIfNotAnnotation()
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
// line. The Java side has additional optional / rest / default-value
// handling; we ignore those (the Tier-0 checker treats them as plain
// names).
func (p *JsDocInfoParser) parseParamTag() JsDocToken {
	p.skipEOLs()
	t := p.next()
	expr, isRest, isOptional, next, ok := p.parseParamTypeAnnotationFull(t)
	if !ok {
		return p.eatUntilEOLIfNotAnnotation()
	}
	// The token after `}` is either:
	//   STRING          — the bare parameter name
	//   LEFT_SQUARE     — the optional-param wrapper `[name]`
	// Java handles LEFT_SQUARE via isBracketedParam. We mirror that here.
	isBracketedParam := next == JsDocTokenLeftSquare
	if isBracketedParam {
		next = p.next() // consume `[` and advance to the name
		isOptional = true
	}
	if next == JsDocTokenString {
		name := p.stream.String_()
		// Optional-param wrappers arrive as raw identifiers; strip a
		// leading `[` if the tokenizer merged it into the name.
		name = strings.TrimLeft(name, "[")
		if i := strings.IndexAny(name, "=]"); i >= 0 {
			name = name[:i]
		}
		if isBracketedParam {
			// Consume optional default-value syntax: `[name=default]`.
			// Eat up to the matching `]`.
			tok := p.next()
			if tok == JsDocTokenEquals {
				// consume any default value tokens up to `]`
				for tok != JsDocTokenRightSquare && tok != JsDocTokenEOC && tok != JsDocTokenEOF {
					tok = p.next()
				}
			} else if tok == JsDocTokenRightSquare {
				// `[name]` — just the closing bracket
			} else {
				// unexpected; continue
				_ = tok
			}
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
		} else {
			// Port of JsDocInfoParser.addParserWarning(JSDOC_INVALID_PARAM_NAME).
			// Store the invalid name; positions are resolved by
			// recordInvalidParamPositions after parsing completes.
			p.info.AddInvalidParamPosition(rhino.InvalidParamPosition{Name: name})
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
	// Braced form: `@extends {Foo}` — standard.
	if expr, next, ok := p.parseTypeExpressionAnnotation(t); ok {
		p.info.SetBaseType(rhino.NewJSTypeExpression(expr))
		p.anyRecorded = true
		return next
	}
	// Bare form: `@extends Foo` (no curly braces). Java's JsDocInfoParser
	// accepts this for @extends and @implements. Mirrors the
	// parseNameExpression fallback path in the Java source.
	if t == JsDocTokenString {
		expr := p.stream.String_()
		p.info.SetBaseType(rhino.NewJSTypeExpression(expr))
		p.anyRecorded = true
		return p.eatUntilEOLIfNotAnnotation()
	}
	return p.eatUntilEOLIfNotAnnotation()
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
// line, split on commas, register each as a template name. Detects TTL
// type-transformation syntax `T := expr =:` and records T as a TTL
// template name rather than a plain @template name (mirrors
// JsDocInfoParser.java's recordTypeTransformation path).
func (p *JsDocInfoParser) parseTemplateTag() JsDocToken {
	rest := p.stream.GetRemainingJSDocLine()
	// Detect TTL: `T := … =:` on this line (single-name constraint).
	// Java calls recordTypeTransformation(name, ttlAst); we just track the name.
	if idx := strings.Index(rest, ":="); idx >= 0 {
		name := strings.TrimSpace(rest[:idx])
		// Strip bounded-generic `{B}` prefix if present.
		if lb := strings.LastIndexByte(name, ' '); lb >= 0 {
			name = strings.TrimSpace(name[lb+1:])
		}
		if name != "" && isJsIdent(name) {
			p.info.AddTTLTemplateName(name)
			p.anyRecorded = true
		}
		return p.next()
	}
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

// parseFunctionType: `function(params): returnType`. Shrunk from
// Java's more complete version.
//
// Param modifiers recognised:
//   - Leading `...`  (JsDocTokenIterRest) marks the param as variadic;
//     emitted as `...T` in the serialised string.
//   - Trailing `=`   (JsDocTokenEquals) marks the param as optional;
//     emitted as `T=` in the serialised string.
//
// The `this:T` qualifier is captured and emitted as a `this:T` prefix
// on the first param slot so resolveFunctionTypeExpr can set typeOfThis.
// `new:T` is consumed without recording (it means the function must be
// called with `new` — not the same as `this:T`).
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
		// Handle leading `this:` / `new:` qualifier. Each is emitted as
		// a `this:T` / `new:T` pseudo-param so resolveFunctionTypeExpr
		// can set typeOfThis (this:) or IsConstructor (new:).
		if tok == JsDocTokenString {
			name := p.stream.String_()
			if (name == "this" || name == "new") && p.match(JsDocTokenColon) {
				p.next() // consume `:`
				p.skipEOLs()
				tok = p.next()
				typeExpr, ok := p.parseTypeExpression(tok)
				if !ok {
					return "", false
				}
				// Emit `this:T` and `new:T` as pseudo-params so the
				// resolver can set typeOfThis (this:) or IsConstructor
				// (new:). Mirrors FunctionType.builder.withTypeOfThis /
				// .forConstructor when these qualifiers are parsed.
				params = append(params, name+":"+typeExpr)
				p.skipEOLs()
				if p.match(JsDocTokenComma) {
					p.next() // consume `,`
				}
				continue
			}
		}
		// Leading `...` marks a variadic (rest) param.
		isRest := false
		if tok == JsDocTokenIterRest {
			isRest = true
			p.skipEOLs()
			tok = p.next()
		}
		expr, ok := p.parseTypeExpression(tok)
		if !ok {
			return "", false
		}
		p.skipEOLs()
		// Trailing `=` marks an optional param.
		isOptional := false
		if p.match(JsDocTokenEquals) {
			p.next()
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

// ParseInlineTypeDoc mirrors Java's parseInlineTypeDoc() — parses an inline
// type annotation (e.g. `/** @type {T} */`) and returns a JSDocInfo with the
// type recorded and the inline-type bit set. Returns nil if no type was parsed.
func (p *JsDocInfoParser) ParseInlineTypeDoc() *rhino.JSDocInfo {
	p.skipEOLs()
	t := p.next()
	expr, next, ok := p.parseTypeExpressionAnnotation(t)
	if !ok {
		return nil
	}
	p.info.SetType(rhino.NewJSTypeExpression(expr))
	p.info.RecordInlineType()
	p.anyRecorded = true
	_ = next
	return p.RetrieveAndResetParsedJSDocInfo()
}

// ParseJsdoc mirrors Java's static parseJsdoc(String) — parses a complete
// JsDoc comment and returns the resulting JSDocInfo, or nil if empty.
func ParseJsdoc(toParse string) *rhino.JSDocInfo {
	p := NewJsDocInfoParser(toParse)
	p.Parse()
	return p.RetrieveAndResetParsedJSDocInfo()
}

// ParseTypeString mirrors Java's static parseTypeString(String) — parses
// a JsDoc type expression and returns it as a string. Returns "" on failure.
func ParseTypeString(typeString string) string {
	p := NewJsDocInfoParser(typeString)
	p.skipEOLs()
	t := p.next()
	// Try to parse as a top-level type expression.
	expr, _ := p.parseTopLevelTypeExpression(t)
	return expr
}

// RetrieveAndResetParsedJSDocInfo mirrors Java's
// retrieveAndResetParsedJSDocInfo() — returns the parsed JSDocInfo if any
// annotation was recorded, otherwise nil.
func (p *JsDocInfoParser) RetrieveAndResetParsedJSDocInfo() *rhino.JSDocInfo {
	if !p.anyRecorded {
		return nil
	}
	return p.info
}

// GetFileOverviewJSDocInfo mirrors Java's getFileOverviewJSDocInfo() —
// returns the fileoverview JSDocInfo if one was parsed. Not implemented
// in the Go port (fileoverview detection requires a separate pass over
// the parsed info). Returns nil by default.
func (p *JsDocInfoParser) GetFileOverviewJSDocInfo() *rhino.JSDocInfo {
	if p.info.HasFileOverview() {
		return p.info
	}
	return nil
}

// lookAheadForType mirrors Java's lookAheadForType() — peeks ahead in the
// character stream (without consuming) to see if a `{` follows.
func (p *JsDocInfoParser) lookAheadForType() bool {
	return p.lookAheadFor('{')
}

// lookAheadForAnnotation mirrors Java's lookAheadForAnnotation() — peeks
// ahead to see if an `@` follows (ignoring leading spaces).
func (p *JsDocInfoParser) lookAheadForAnnotation() bool {
	return p.lookAheadFor('@')
}

// lookAheadFor peeks ahead in the character stream for a specific rune,
// skipping spaces only. Mirrors Java's lookAheadFor(char).
func (p *JsDocInfoParser) lookAheadFor(expect byte) bool {
	rest := p.stream.GetRemainingJSDocLine()
	for _, c := range rest {
		if c == ' ' || c == '\t' {
			continue
		}
		return byte(c) == expect
	}
	return false
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

// eatUntilEOLIfNotAnnotationFrom mirrors Java's
// eatUntilEOLIfNotAnnotation(JsDocToken) — starts from the given token.
// If token is ANNOTATION, returns it immediately; otherwise eats until EOL/EOC/EOF.
func (p *JsDocInfoParser) eatUntilEOLIfNotAnnotationFrom(token JsDocToken) JsDocToken {
	if token == JsDocTokenAnnotation {
		return token
	}
	return p.eatTokensUntilEOL(token)
}

// eatTokensUntilEOL mirrors Java's eatTokensUntilEOL() — consumes tokens
// until EOL, EOC, or EOF (but not ANNOTATION), then returns that terminating token.
func (p *JsDocInfoParser) eatTokensUntilEOL(token JsDocToken) JsDocToken {
	for {
		switch token {
		case JsDocTokenEOL, JsDocTokenEOC, JsDocTokenEOF:
			return token
		}
		token = p.next()
	}
}

// eatIfMatch mirrors Java's eatIfMatch(JsDocToken) — if the next token
// matches, consume it and return true.
func (p *JsDocInfoParser) eatIfMatch(tok JsDocToken) bool {
	if p.match(tok) {
		p.next()
		return true
	}
	return false
}

// parseIdGeneratorTag ports Java's parseIdGeneratorTag: expects an optional
// `{kind}` and records the id generator kind. Mirrors Java exactly.
func (p *JsDocInfoParser) parseIdGeneratorTag(token JsDocToken) JsDocToken {
	idgenKind := "unique"
	if token == JsDocTokenLeftCurly {
		if p.match(JsDocTokenString) {
			idgenKind = p.stream.String_()
			token = p.next()
		} else {
			return p.eatUntilEOLIfNotAnnotation()
		}
		if !p.match(JsDocTokenRightCurly) {
			return p.eatUntilEOLIfNotAnnotation()
		}
		token = p.next()
	}
	switch idgenKind {
	case "unique":
		p.info.SetIdGenerator(rhino.IdGeneratorUnique)
	case "consistent":
		p.info.SetIdGenerator(rhino.IdGeneratorConsistent)
	case "stable":
		p.info.SetIdGenerator(rhino.IdGeneratorStable)
	case "xid":
		p.info.SetIdGenerator(rhino.IdGeneratorXid)
	case "mapped":
		p.info.SetIdGenerator(rhino.IdGeneratorMapped)
	}
	p.anyRecorded = true
	return token
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

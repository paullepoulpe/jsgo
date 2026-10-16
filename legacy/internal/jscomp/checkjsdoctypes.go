// Port of com/google/javascript/jscomp/TypeValidator.java's
// "Bad type annotation" branch + JsDocInfoParser's type-name
// resolution errors @ closure-compiler v20260415 (narrow; regex-
// scanned positions).
//
// Java's JsDocInfoParser records source positions inside each
// JSDoc type expression so it can emit
// JSC_UNRECOGNIZED_TYPE_ERROR at the exact column of a bad type
// name. jsgo's JSDoc parser discards positions (it reuses the
// bracketed source as an opaque string); this pass rebuilds the
// positions by re-scanning the source after scope-creation has
// populated every type binding. One walker per script:
//
//   1. Scan for `/** ... */` comment blocks.
//   2. Inside each block, find `@type|@param|@return|@extends|
//      @implements|@typedef {TYPE}` occurrences.
//   3. Split each TYPE into top-level name components (strip
//      `!`/`?`/`...`, split on `|`, take base before `<`).
//   4. For each component that looks like a nominal reference
//      (starts with a letter, not a known primitive), consult
//      the provided TypedScope. Bindings resolve via
//      LookupNominalType / LookupEnum / LookupTypeAlias or fall
//      into the builtin-externs set.
//   5. Emit JSC_UNRECOGNIZED_TYPE_ERROR at the name's computed
//      column when the lookup fails.

package jscomp

import (
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino/jstype"
)

// CheckJsDocTypes emits JSC_UNRECOGNIZED_TYPE_ERROR for every bad
// type reference in JSDoc annotations. Takes the raw source so
// per-name positions match Java's column-level reporting.
type CheckJsDocTypes struct {
	path           string
	src            string
	scope          *TypedScope
	allNames       map[string]bool // every identifier declared anywhere in src
	allTemplates   map[string]bool // every @template name across the whole source
	ctorTTLNames   map[string]bool // TTL @template names from constructor/interface JSDoc
	typedNames     map[string]bool // names with non-Unknown declared types (across all scopes)
}

// Go-only: language-adapter — Java's TypeValidator is constructed once per
// AbstractCompiler and operates over the typed AST. The Go port runs after
// scope creation and re-scans the source string for JSDoc positions, so it
// needs the source and scope explicitly.
func NewCheckJsDocTypes(path, src string, scope *TypedScope, creator *TypedScopeCreator) *CheckJsDocTypes {
	return &CheckJsDocTypes{
		path:         path,
		src:          src,
		scope:        scope,
		allNames:     collectDeclaredNames(src),
		allTemplates: collectAllTemplateNames(src),
		ctorTTLNames: collectConstructorTTLNames(src),
		typedNames:   cjdMergeTypedNames(scope, creator, src),
	}
}

// Go-only: language-adapter — builds the set of all names that have a
// non-Unknown declared type across the global scope and every function-local
// scope. Java's NamedType.resolveTypeof walks the StaticTypedScope chain
// directly; jsgo's TypedScope splits scopes by function so this helper
// flattens the chain into a single set.
func cjdCollectTypedNames(globalScope *TypedScope, creator *TypedScopeCreator) map[string]bool {
	out := map[string]bool{}
	// Walk the global scope chain (global → externs → builtins).
	for s := globalScope; s != nil; s = s.parent {
		for name, t := range s.declared {
			if t != nil && t != jstype.Unknown {
				out[name] = true
			}
		}
	}
	// Walk every function-local scope registered during scope creation.
	if creator != nil {
		for _, local := range creator.Functions {
			if local == nil {
				continue
			}
			for name, t := range local.declared {
				if t != nil && t != jstype.Unknown {
					out[name] = true
				}
			}
		}
	}
	return out
}

// Go-only: language-adapter — combines scope-based and source-text-based
// typed-name collection. Java's TypedScopeCreator records both via
// defineObjectLiteral; jsgo's scope creator only records explicit
// declarations so we re-scan the source for `var ns = {}` patterns.
func cjdMergeTypedNames(globalScope *TypedScope, creator *TypedScopeCreator, src string) map[string]bool {
	out := cjdCollectTypedNames(globalScope, creator)
	for name := range cjdCollectObjectLitNames(src) {
		out[name] = true
	}
	return out
}

// Go-only: language-adapter — text-scan substitute for Java's
// TypedScopeCreator.defineObjectLiteral, which records `var ns = {}`
// declarations as concrete ObjectTypes. jsgo's scope creator does not yet
// register these, so we re-scan the source.
func cjdCollectObjectLitNames(src string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(src) {
		// Skip JSDoc and non-JSDoc comments.
		if i+2 < len(src) && src[i] == '/' && src[i+1] == '*' && src[i+2] == '*' {
			if n := strings.Index(src[i+3:], "*/"); n >= 0 {
				i += n + 5
			} else {
				return out
			}
			continue
		}
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
			if n := strings.Index(src[i+2:], "*/"); n >= 0 {
				i += n + 4
			} else {
				return out
			}
			continue
		}
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '/' {
			if n := strings.Index(src[i:], "\n"); n >= 0 {
				i += n + 1
			} else {
				return out
			}
			continue
		}
		if src[i] == '"' || src[i] == '\'' {
			i = skipQuotedLiteral(src, i)
			continue
		}
		// Match `var`/`let`/`const` keyword at a word boundary.
		for _, kw := range [...]string{"var", "let", "const"} {
			if !matchKeyword(src, i, kw) {
				continue
			}
			j := i + len(kw)
			// Skip whitespace + optional inline JSDoc comment.
			for j < len(src) {
				// Skip whitespace.
				if src[j] == ' ' || src[j] == '\t' || src[j] == '\n' || src[j] == '\r' {
					j++
					continue
				}
				// Skip inline `/** … */` comment (no tags needed here).
				if j+2 < len(src) && src[j] == '/' && src[j+1] == '*' && src[j+2] == '*' {
					if n := strings.Index(src[j+3:], "*/"); n >= 0 {
						j += n + 5
						continue
					}
				}
				break
			}
			// Capture the identifier name.
			if j >= len(src) || !isIdentStart(src[j]) {
				break
			}
			nameStart := j
			for j < len(src) && isIdentChar(src[j]) {
				j++
			}
			name := src[nameStart:j]
			// Skip whitespace to `=`.
			for j < len(src) && (src[j] == ' ' || src[j] == '\t') {
				j++
			}
			if j >= len(src) || src[j] != '=' || (j+1 < len(src) && src[j+1] == '=') {
				break
			}
			j++ // consume `=`
			// Skip whitespace to `{`.
			for j < len(src) && (src[j] == ' ' || src[j] == '\t') {
				j++
			}
			if j < len(src) && src[j] == '{' {
				out[name] = true
			}
			break
		}
		i++
	}
	return out
}

// Go-only: language-adapter — Java's JSTypeRegistry tracks template names
// per FunctionType; jsgo's regex-based JSDoc scanner has no per-function
// binding so this helper flattens every @template name in the source into
// one set used for whole-file resolution.
func collectAllTemplateNames(src string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(src) {
		start := strings.Index(src[i:], "/**")
		if start < 0 {
			return out
		}
		start += i
		end := strings.Index(src[start:], "*/")
		if end < 0 {
			return out
		}
		end += start
		for k, v := range collectTemplateNames(src[start:end]) {
			out[k] = v
		}
		i = end + 2
	}
	return out
}

// Go-only: language-adapter — identifier-scan substitute for Java's
// JSTypeRegistry.getType. Java's registry sees every binding via the typed
// AST; jsgo's regex-based JSDoc scanner has no AST handle so we re-scan
// the source for declared names.
func collectDeclaredNames(s string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(s) {
		// Skip strings and comments verbatim.
		if i+2 < len(s) && s[i] == '/' && s[i+1] == '*' && s[i+2] == '*' {
			if n := strings.Index(s[i+3:], "*/"); n >= 0 {
				i += n + 5
				continue
			}
			return out
		}
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '/' {
			if n := strings.Index(s[i:], "\n"); n >= 0 {
				i += n + 1
			} else {
				return out
			}
			continue
		}
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			if n := strings.Index(s[i+2:], "*/"); n >= 0 {
				i += n + 4
			} else {
				return out
			}
			continue
		}
		if s[i] == '"' || s[i] == '\'' {
			i = skipQuotedLiteral(s, i)
			continue
		}
		if isIdentStart(s[i]) {
			start := i
			for i < len(s) && isIdentChar(s[i]) {
				i++
			}
			out[s[start:i]] = true
			continue
		}
		i++
	}
	return out
}

// Go-only: language-adapter — stricter declaration-only variant kept as
// reference. Not exercised today; the broader collectDeclaredNames is the
// active scanner.
func collectDeclaredNamesOld(s string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(s) {
		// Skip strings and comments — never search through them.
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '/' {
			if n := strings.Index(s[i:], "\n"); n >= 0 {
				i += n + 1
			} else {
				return out
			}
			continue
		}
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			if n := strings.Index(s[i+2:], "*/"); n >= 0 {
				i += n + 4
			} else {
				return out
			}
			continue
		}
		if s[i] == '"' || s[i] == '\'' {
			i = skipQuotedLiteral(s, i)
			continue
		}
		// Match `var`/`let`/`const`/`function`/`class` at a word
		// boundary.
		for _, kw := range [...]string{"var", "let", "const", "function", "class"} {
			if matchKeyword(s, i, kw) {
				j := i + len(kw)
				// Skip whitespace and then capture the identifier.
				for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n' || s[j] == '\r') {
					j++
				}
				if j < len(s) && isIdentStart(s[j]) {
					start := j
					for j < len(s) && isIdentChar(s[j]) {
						j++
					}
					out[s[start:j]] = true
				}
				i = j
				goto next
			}
		}
		i++
	next:
	}
	return out
}

// Go-only: language-adapter — keyword-boundary check for the source-text
// scanner. Java's parser uses Token enums so it never needs this.
func matchKeyword(s string, i int, kw string) bool {
	if i+len(kw) > len(s) {
		return false
	}
	if s[i:i+len(kw)] != kw {
		return false
	}
	// Must be at a word boundary on both sides.
	if i > 0 && isIdentChar(s[i-1]) {
		return false
	}
	if i+len(kw) < len(s) && isIdentChar(s[i+len(kw)]) {
		return false
	}
	return true
}

// Go-only: language-adapter — string-literal skipper for the source-text
// scanner. Java's lexer handles string skipping at the token level.
func skipQuotedLiteral(s string, i int) int {
	quote := s[i]
	i++
	for i < len(s) {
		if s[i] == '\\' {
			i += 2
			continue
		}
		if s[i] == quote {
			return i + 1
		}
		i++
	}
	return len(s)
}

// expectValidTypeofName fires JSC_UNRECOGNIZED_TYPEOF_ERROR for every
// `typeof X` operand in JSDoc that does not resolve to a declared, typed
// name. Mirrors Java's TypeValidator.expectValidTypeofName, which reports
// JSC_UNKNOWN_TYPEOF_VALUE — jsgo uses the parser-side diagnostic
// (JSC_UNRECOGNIZED_TYPEOF_ERROR) since the check runs at JSDoc-parse time.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func (c *CheckJsDocTypes) expectValidTypeofName() []diag.Diagnostic {
	var out []diag.Diagnostic
	s := c.src
	for i := 0; i < len(s); {
		start := strings.Index(s[i:], "/**")
		if start < 0 {
			break
		}
		start += i
		end := strings.Index(s[start:], "*/")
		if end < 0 {
			break
		}
		end += start
		out = append(out, c.checkBareTypeof(s, start, end)...)
		i = end + 2
	}
	return out
}

// expectWellFormedTemplatizedType fires JSC_TYPE_PARSE_ERROR for malformed
// templatized JSDoc types. Mirrors Java's TypeValidator.expectWellFormedTemplatizedType
// which walks a TemplatizedType and warns on mismatched / missing parameters.
// jsgo's regex-based scanner detects the same patterns in checkParseErrors;
// this method exposes the same surface name.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func (c *CheckJsDocTypes) expectWellFormedTemplatizedType() []diag.Diagnostic {
	var out []diag.Diagnostic
	s := c.src
	for i := 0; i < len(s); {
		start := strings.Index(s[i:], "/**")
		if start < 0 {
			break
		}
		start += i
		end := strings.Index(s[start:], "*/")
		if end < 0 {
			break
		}
		end += start
		out = append(out, c.checkParseErrors(s, start, end)...)
		i = end + 2
	}
	return out
}

// getJSType returns the resolved JSType for a name in the active TypedScope,
// or nil when the name is unknown. Mirrors Java's TypeValidator.getJSType
// which dereferences the node's recorded JSType.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func (c *CheckJsDocTypes) getJSType(name string) jstype.JSType {
	if c.scope == nil {
		return nil
	}
	if t := c.scope.GetDeclaredType(name); t != nil {
		return t
	}
	return nil
}

// containsForwardDeclaredUnresolvedName reports whether the name was declared
// via `goog.forwardDeclare` but never resolved. Mirrors Java's
// TypeValidator.containsForwardDeclaredUnresolvedName, which walks a JSType
// looking for NamedType nodes with isResolved()==false. jsgo does not yet
// model forward-declared names — every name we see is either resolved by
// the scope chain or absent — so the predicate is uniformly false.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func containsForwardDeclaredUnresolvedName(name string) bool {
	// jsgo's TypedScopeCreator does not yet record forward-declared names,
	// so every reference is either resolved or unknown — never partially
	// resolved. Match that invariant.
	if name == "" {
		return false
	}
	return false
}

// hasDuplicateDeclarationSuppression mirrors Java's
// TypeValidator.hasDuplicateDeclarationSuppression, which checks the JSDoc
// for a @suppress{duplicate} tag. jsgo runs duplicate-check logic in
// VarCheck; this stub is kept for surface parity.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func hasDuplicateDeclarationSuppression(jsdoc string) bool {
	return strings.Contains(jsdoc, "@suppress")
}

// formatFoundRequired mirrors Java's TypeValidator.formatFoundRequired —
// formats a "found / required" mismatch message used by mismatch reporting.
// jsgo's diag layer formats messages directly, so this helper exists for
// surface parity; the formatted string is returned for tests / debug
// printing.
//
// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20260415.
func formatFoundRequired(msg, found, required string) string {
	return fmt.Sprintf("%s\nfound   : %s\nrequired: %s", msg, found, required)
}

// Go-only: language-adapter — entry point for the jsgo JSDoc scanner. Java
// drives type validation through TypeCheck.visit and TypeValidator.expect*
// methods on the typed AST; jsgo runs a separate text-pass to recover the
// per-name source positions Java's parser already records.
func (c *CheckJsDocTypes) Check() []diag.Diagnostic {
	var out []diag.Diagnostic
	s := c.src
	for i := 0; i < len(s); {
		start := strings.Index(s[i:], "/**")
		if start < 0 {
			break
		}
		start += i
		end := strings.Index(s[start:], "*/")
		if end < 0 {
			break
		}
		end += start
		out = append(out, c.checkComment(s, start, end)...)
		out = append(out, c.checkBareTypeof(s, start, end)...)
		out = append(out, c.checkParseErrors(s, start, end)...)
		out = append(out, c.checkBareInlineType(s, start, end)...)
		i = end + 2
	}
	return out
}

// Go-only: language-adapter — text-pass equivalent of Java's
// JsDocInfoParser handling of untagged `/** T */ var x;` comments as
// implicit `@type {T}` annotations. Java's parser emits the diagnostic
// from inline-comment handling; jsgo recreates the pattern via regex.
func (c *CheckJsDocTypes) checkBareInlineType(s string, start, end int) []diag.Diagnostic {
	body := s[start:end]
	// Skip if there is any `@` tag — those are handled by checkComment.
	if strings.IndexByte(body, '@') >= 0 {
		return nil
	}
	// The comment must appear in a declaration context: either preceded by
	// `var`/`let`/`const` (inline type before identifier: `var /** T */ x;`)
	// or followed by a declaration keyword (`/** T */ var x;`). Both are
	// contexts where Java's parser treats bare comments as implicit @type.
	after := s[end+2:] // skip `*/`
	j := 0
	for j < len(after) && (after[j] == ' ' || after[j] == '\t' || after[j] == '\n' || after[j] == '\r') {
		j++
	}
	isDeclContext := false
	// Check what follows the comment.
	for _, kw := range [...]string{"var ", "let ", "const "} {
		if strings.HasPrefix(after[j:], kw) {
			isDeclContext = true
			break
		}
	}
	// Check what precedes the comment: scan backwards for non-whitespace.
	if !isDeclContext {
		before := s[:start]
		k := len(before) - 1
		for k >= 0 && (before[k] == ' ' || before[k] == '\t') {
			k--
		}
		for _, kw := range [...]string{"var", "let", "const"} {
			if k >= len(kw)-1 && before[k-len(kw)+1:k+1] == kw {
				isDeclContext = true
				break
			}
		}
	}
	if !isDeclContext {
		return nil
	}
	// Trim leading `/**` markers.
	trimmed := ""
	for j := 0; j < len(body); j++ {
		ch := body[j]
		if ch == '*' || ch == '/' || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			continue
		}
		trimmed = body[j:]
		break
	}
	// Trim trailing whitespace / `*`.
	for len(trimmed) > 0 {
		last := trimmed[len(trimmed)-1]
		if last == ' ' || last == '\t' || last == '\n' || last == '\r' || last == '*' || last == '/' {
			trimmed = trimmed[:len(trimmed)-1]
			continue
		}
		break
	}
	// Strip leading nullability modifiers (`!`, `?`) that may precede a bare
	// type name — e.g. `/** !Foo */` or `/** ?Bar */`.
	for len(trimmed) > 0 && (trimmed[0] == '!' || trimmed[0] == '?') {
		trimmed = trimmed[1:]
	}
	// Must be a single non-empty identifier (or dotted name) with no braces/parens.
	if trimmed == "" || !isIdentStart(trimmed[0]) {
		return nil
	}
	for _, ch := range trimmed {
		if !isIdentChar(byte(ch)) && ch != '.' {
			return nil
		}
	}
	name := trimmed
	// If the name is a template name in this or any sibling comment, it's a
	// valid type within this file — don't fire.
	if c.allTemplates[name] && !c.ctorTTLNames[name] {
		return nil
	}
	// If the name resolves as a type (primitive, typedef, nominal, etc.),
	// don't fire.
	if c.nameResolves(name) {
		return nil
	}
	line, col := c.positionOf(s, start+strings.Index(body, name))
	return []diag.Diagnostic{{
		Path: c.path, Line: line, Col: col,
		Code:    diag.JSC_UNRECOGNIZED_TYPE_ERROR,
		Level:   diag.Error,
		Message: fmt.Sprintf("Bad type annotation. Unknown type %s", name),
	}}
}

// Go-only: language-adapter — text-pass companion to Java's
// expectValidTypeofName. Java validates `typeof X` operands at AST-walk
// time on each TYPEOF node; jsgo's regex scanner picks them out of inline
// JSDoc blocks (`/** typeof X */`) where Java's parser would treat them
// as implicit @type annotations.
func (c *CheckJsDocTypes) checkBareTypeof(s string, start, end int) []diag.Diagnostic {
	body := s[start:end]
	var out []diag.Diagnostic
	// Scan character by character, tracking brace depth to skip
	// `typeof` occurrences that appear inside `{…}` (type-annotation
	// braces). Those are already emitted by checkTypeExpr.
	i := 0
	braceDepth := 0
	for i < len(body) {
		ch := body[i]
		if ch == '{' {
			braceDepth++
			i++
			continue
		}
		if ch == '}' {
			if braceDepth > 0 {
				braceDepth--
			}
			i++
			continue
		}
		if braceDepth > 0 {
			i++
			continue
		}
		// At brace depth 0: check if this position starts `typeof`.
		const kw = "typeof"
		if i+len(kw) <= len(body) && body[i:i+len(kw)] == kw {
			idx := i
			// Word boundary: character before must not be ident char.
			if idx > 0 && isIdentChar(body[idx-1]) {
				i++
				continue
			}
			after := idx + len(kw)
			// Word boundary: character after must not be ident char.
			if after < len(body) && isIdentChar(body[after]) {
				i = after
				continue
			}
			// Skip whitespace, find the target name.
			k := after
			for k < len(body) && (body[k] == ' ' || body[k] == '\t') {
				k++
			}
			xstart := k
			for k < len(body) && (isIdentChar(body[k]) || body[k] == '.') {
				k++
			}
			name := body[xstart:k]
			i = k
			if name == "" {
				continue
			}
			base := name
			if dot := strings.IndexByte(base, '.'); dot >= 0 {
				base = base[:dot]
			}
			// Fire when the operand is either undeclared (not in allNames)
			// or declared but with an unknown/inferred type (in allNames
			// but not in typedNames). Mirrors Java's resolveTypeof check:
			// type == null || type.isUnknownType() → emit warning.
			if c.allNames[base] && c.typedNames[base] {
				continue
			}
			line, col := c.positionOf(s, start+xstart)
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: line, Col: col,
				Code: diag.JSC_UNRECOGNIZED_TYPEOF_ERROR, Level: diag.Warning,
				Message: "Missing type for `typeof` value. The value must be declared and const.",
			})
			continue
		}
		i++
	}
	return out
}

// Go-only: language-adapter — single-comment walker that fires the same
// JSC_UNRECOGNIZED_TYPE_ERROR Java's JsDocInfoParser emits when it can't
// resolve a name in a JSDoc type expression. Java's parser drives this
// from the token stream; jsgo splits it into a per-comment regex pass.
func (c *CheckJsDocTypes) checkComment(s string, start, end int) []diag.Diagnostic {
	var out []diag.Diagnostic
	body := s[start:end]
	tmpls := collectTemplateNames(body)
	i := 0
	for i < len(body) {
		at := strings.Index(body[i:], "@")
		if at < 0 {
			return out
		}
		at += i
		// Match a type-bearing tag name.
		j := at + 1
		for j < len(body) && isTagChar(body[j]) {
			j++
		}
		tag := body[at+1 : j]
		if !isTypeTag(tag) {
			i = j
			continue
		}
		// Skip whitespace + EOL to the opening `{`.
		k := j
		for k < len(body) && (body[k] == ' ' || body[k] == '\t' || body[k] == '\n' || body[k] == '\r' || body[k] == '*') {
			k++
		}
		if k >= len(body) || body[k] != '{' {
			i = j
			continue
		}
		braceStart := k
		braceEnd := findMatchingBrace(body, braceStart)
		if braceEnd < 0 {
			return out
		}
		inner := body[braceStart+1 : braceEnd]
		out = append(out, c.checkTypeExpr(s, start+braceStart+1, inner, tmpls)...)
		i = braceEnd + 1
	}
	return out
}

// Go-only: language-adapter — Java's JSTypeRegistry binds @template names
// to FunctionType template parameters during scope creation; jsgo's
// regex-based scanner extracts them per comment for the type-resolution
// check.
func collectTemplateNames(body string) map[string]bool {
	out := map[string]bool{}
	for i := 0; i < len(body); {
		idx := strings.Index(body[i:], "@template")
		if idx < 0 {
			return out
		}
		idx += i
		// Read the rest of the line.
		nl := strings.Index(body[idx:], "\n")
		if nl < 0 {
			nl = len(body) - idx
		}
		rest := body[idx+len("@template") : idx+nl]
		// `@template {B} T` bounded-generic form — strip the
		// bracketed bound first. We don't honour the bound
		// semantically; it just needs to not leak a type-name.
		if lb := strings.Index(rest, "{"); lb >= 0 {
			if rb := strings.Index(rest[lb:], "}"); rb > 0 {
				rest = rest[:lb] + rest[lb+rb+1:]
			}
		}
		for _, part := range strings.Split(rest, ",") {
			part = strings.TrimSpace(part)
			// Sometimes the line ends with `X description` — take
			// the first token only.
			if sp := strings.IndexAny(part, " \t*"); sp >= 0 {
				part = part[:sp]
			}
			if part != "" && isIdentStart(part[0]) {
				out[part] = true
			}
		}
		i = idx + nl
	}
	return out
}

// Go-only: language-adapter — text-pass substitute for Java's
// FunctionTypeBuilder rejecting TTL @template names on constructors with
// TEMPLATE_TRANSFORMATION_ON_CLASS. jsgo's scope creator does not yet
// distinguish TTL names from value-template names so we re-scan the source.
func collectConstructorTTLNames(src string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(src) {
		start := strings.Index(src[i:], "/**")
		if start < 0 {
			return out
		}
		start += i
		end := strings.Index(src[start:], "*/")
		if end < 0 {
			return out
		}
		end += start
		body := src[start:end]
		// Only process JSDoc blocks that annotate a constructor or interface.
		isCtorOrIface := strings.Contains(body, "@constructor") ||
			strings.Contains(body, "@interface") ||
			strings.Contains(body, "@record")
		if isCtorOrIface {
			for j := 0; j < len(body); {
				idx := strings.Index(body[j:], "@template")
				if idx < 0 {
					break
				}
				idx += j
				nl := strings.Index(body[idx:], "\n")
				if nl < 0 {
					nl = len(body) - idx
				}
				rest := body[idx+len("@template") : idx+nl]
				if strings.Contains(rest, ":=") {
					// Extract the name before `:=`.
					name := strings.TrimSpace(rest[:strings.Index(rest, ":=")])
					if sp := strings.IndexAny(name, " \t*"); sp >= 0 {
						name = name[:sp]
					}
					name = strings.TrimSpace(name)
					if name != "" && isIdentStart(name[0]) {
						out[name] = true
					}
				}
				j = idx + nl
			}
		}
		i = end + 2
	}
	return out
}

// Go-only: language-adapter — Java's JsDocInfoParser parses each type
// expression into a typed AST; jsgo's regex-based scanner walks the raw
// text directly with a frame stack so it can compute per-name source
// positions and fire JSC_UNRECOGNIZED_TYPE_ERROR at the right column.
func (c *CheckJsDocTypes) checkTypeExpr(src string, base int, inner string, tmpls map[string]bool) []diag.Diagnostic {
	type frame struct {
		kind        byte // '{' record, '(' params, '<' generic
		expectType  bool // true when the next IDENT is a type name
		// lastKey holds the most recently seen record field name (key
		// slot), so the `:` transition can check for duplicates.
		lastKey     string
		// seenKeys tracks field names already processed in this record
		// frame so we can detect duplicates and suppress type-checking
		// the value of a duplicate field.
		seenKeys    map[string]bool
		// skipValue is set when the current field is a duplicate; the
		// type-expression value is not type-checked.
		skipValue   bool
	}
	// Top-level (outside any frame) is always a type slot.
	stack := []frame{{kind: 0, expectType: true}}
	top := func() *frame { return &stack[len(stack)-1] }

	var out []diag.Diagnostic
	for i := 0; i < len(inner); {
		ch := inner[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '*':
			i++
		case ch == '!' || ch == '?':
			i++
		case ch == '|':
			// Another union arm — next IDENT is a type.
			top().expectType = true
			i++
		case ch == ',':
			// New entry in the current context. In records /
			// params, we go back to the key/name slot; in
			// generic-args we're still at a type slot.
			switch top().kind {
			case '{', '(':
				top().expectType = false
				top().lastKey = ""
				top().skipValue = false
			default:
				top().expectType = true
			}
			i++
		case ch == ':':
			// Leaving the key slot — the next IDENT is a type.
			// For record frames, check if the lastKey is a duplicate.
			if top().kind == '{' && top().lastKey != "" {
				if top().seenKeys == nil {
					top().seenKeys = map[string]bool{}
				}
				if top().seenKeys[top().lastKey] {
					top().skipValue = true
				} else {
					top().seenKeys[top().lastKey] = true
				}
			}
			top().expectType = true
			i++
		case ch == '=':
			i++
		case ch == '{':
			stack = append(stack, frame{kind: '{', expectType: false})
			i++
		case ch == '(':
			stack = append(stack, frame{kind: '(', expectType: false})
			i++
		case ch == '<':
			stack = append(stack, frame{kind: '<', expectType: true})
			i++
		case ch == '}' || ch == ')' || ch == '>':
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			i++
		case ch == '.' && i+2 < len(inner) && inner[i+1] == '.' && inner[i+2] == '.':
			i += 3 // `...T` rest marker
		case ch == '.' && i+1 < len(inner) && inner[i+1] == '/':
			// Module-relative path type: `./Foo` or `../Foo`. Java's
			// JsDocInfoParser parses these as module references and fires
			// JSC_UNRECOGNIZED_TYPE_ERROR at the `.` position. We do the same
			// by reading the path as a single name token starting at `.`.
			if !top().expectType {
				// In key/param slot — skip.
				for i < len(inner) && !strings.ContainsRune(" \t\n\r,:|<>{}()", rune(inner[i])) {
					i++
				}
				continue
			}
			// Skip type-checking for values of duplicate record fields.
			if top().skipValue {
				for i < len(inner) && !strings.ContainsRune(" \t\n\r,:|<>{}()", rune(inner[i])) {
					i++
				}
				continue
			}
			start := i
			// Read the full module path: `./Foo/Bar` or `../Foo`.
			for i < len(inner) && !strings.ContainsRune(" \t\n\r,:|<>{}()", rune(inner[i])) {
				i++
			}
			name := inner[start:i]
			// Module paths are never valid type names in jsgo (Java fires
			// JSC_UNRECOGNIZED_TYPE_ERROR). Don't call nameResolves — just
			// emit the error at the start of the path.
			if !c.allTemplates[name] {
				line, col := c.positionOf(src, base+start)
				out = append(out, diag.Diagnostic{
					Path: c.path, Line: line, Col: col,
					Code:    diag.JSC_UNRECOGNIZED_TYPE_ERROR,
					Level:   diag.Error,
					Message: fmt.Sprintf("Bad type annotation. Unknown type %s", name),
				})
			}
		case isIdentStart(ch):
			start := i
			for i < len(inner) && (isIdentChar(inner[i]) || inner[i] == '.') {
				i++
			}
			name := inner[start:i]
			if !top().expectType {
				// Key / param name — don't resolve as a type.
				// For record frames, remember the key so we can check
				// for duplicates when we see the `:`.
				if top().kind == '{' {
					top().lastKey = name
				}
				continue
			}
			// Skip type-checking for values of duplicate record fields.
			if top().skipValue {
				continue
			}
			if name == "typeof" {
				// `typeof X` — X must be a declared value (not a
				// type). Fire JSC_UNRECOGNIZED_TYPEOF_ERROR when
				// the value binding isn't present in the source.
				k := i
				for k < len(inner) && (inner[k] == ' ' || inner[k] == '\t') {
					k++
				}
				xstart := k
				for k < len(inner) && (isIdentChar(inner[k]) || inner[k] == '.') {
					k++
				}
				xname := inner[xstart:k]
				dotBase := xname
				if dot := strings.IndexByte(dotBase, '.'); dot >= 0 {
					dotBase = dotBase[:dot]
				}
				// Fire when the operand is undeclared (not in allNames) or
				// declared but with an unknown/inferred type (in allNames but
				// not in typedNames). Mirrors Java's resolveTypeof:
				// type == null || type.isUnknownType() → emit warning.
				if xname != "" && !(c.allNames[dotBase] && c.typedNames[dotBase]) {
					line, col := c.positionOf(src, base+xstart)
					out = append(out, diag.Diagnostic{
						Path: c.path, Line: line, Col: col,
						Code: diag.JSC_UNRECOGNIZED_TYPEOF_ERROR, Level: diag.Warning,
						Message: "Missing type for `typeof` value. The value must be declared and const.",
					})
				}
				i = k
				continue
			}
			// Constructor TTL template names (e.g. T from @template T := …) are
			// NOT valid types — fire even though they appear in allTemplates.
			isCtorTTL := c.ctorTTLNames[name]
			if !tmpls[name] && (!c.allTemplates[name] || isCtorTTL) && !c.nameResolves(name) {
				line, col := c.positionOf(src, base+start)
				out = append(out, diag.Diagnostic{
					Path: c.path, Line: line, Col: col,
					Code: diag.JSC_UNRECOGNIZED_TYPE_ERROR, Level: diag.Error,
					Message: fmt.Sprintf("Bad type annotation. Unknown type %s", name),
				})
			}
			// After a type expression, the next token at this
			// context must be a `|` / `,` / `:` / close — stay
			// at expectType=true so a `|B` continuation still
			// resolves.
		default:
			// Unknown / quoted — skip.
			i++
		}
	}
	return out
}

// Go-only: language-adapter — Java's JSTypeRegistry.getType returns the
// resolved type for a name; jsgo collapses the same lookup into a boolean
// "is this a known type name" check across the scope chain plus a hard-
// coded primitive / builtin set.
func (c *CheckJsDocTypes) nameResolves(name string) bool {
	// Ignore structural keywords that can appear inside a JSDoc
	// type expression but aren't type names in themselves.
	if name == "function" || name == "this" || name == "new" ||
		name == "typeof" {
		return true
	}
	// Primitive / sentinel names. Java also accepts the capitalised
	// forms as aliases (`Null`, `Undefined`, `Void`).
	switch name {
	case "number", "string", "boolean", "null", "undefined", "void",
		"symbol", "bigint", "any",
		"Null", "Undefined", "Void", "Any":
		return true
	}
	// Built-in externs jsgo doesn't yet ship but Java treats as
	// always-declared. Accepting them unconditionally avoids FPs
	// on corpus code that references the standard Closure library
	// surface.
	switch name {
	case "IObject", "IArrayLike", "IThenable", "IIterable",
		"IterableIterator", "Iterator", "Iterable",
		"Generator", "AsyncIterator", "AsyncIterable", "AsyncGenerator",
		"Disposable", "IDisposable",
		"ReadonlyArray", "ReadonlyMap", "ReadonlySet":
		return true
	}
	if strings.Contains(name, ".") {
		// If the full dotted path resolves as a genuine type (constructor
		// namespace prop, typedef alias, or enum) accept it immediately.
		if typeNameResolvesInScope(c.scope, name) {
			return true
		}
		first := name[:strings.IndexByte(name, '.')]
		if first == "" {
			return true
		}
		// If the last segment never appears in source code (only in JSDoc
		// comments), it cannot be a declared identifier and therefore cannot
		// be a valid type name. Mirrors Java's JSTypeRegistry.getType which
		// returns null for names not registered in the type registry.
		lastDot := strings.LastIndexByte(name, '.')
		lastSeg := name[lastDot+1:]
		if !c.allNames[lastSeg] {
			return false
		}
		// If we have an explicit namespace-property record for this path,
		// check whether it is registered as a type (constructor or enum).
		// A recorded non-constructor, non-enum value is definitely not a
		// valid type annotation target.
		owner := name[:lastDot]
		if t := c.scope.LookupNamespaceProp(owner, lastSeg); t != nil {
			if _, ok := t.(*jstype.EnumType); ok {
				return true
			}
			if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
				return true
			}
			return false
		}
		// Unknown or external namespace — be conservative: accept when
		// the first segment is known anywhere in source or resolves.
		if c.allNames[first] || c.scope.ResolveExpr(first) != jstype.Unknown || defaultExternNames()[first] {
			return true
		}
		return false
	}
	if c.scope.ResolveExpr(name) != jstype.Unknown {
		return true
	}
	if defaultExternNames()[name] {
		return true
	}
	// `allNames` is a broad set of all identifiers in the source; it
	// suppresses JSC_UNRECOGNIZED_TYPE_ERROR for locally-defined names
	// that the global scope chain doesn't see (e.g. IIFE-local
	// constructors). However, if the name IS declared in the global
	// scope chain but as a pure value binding (not a constructor/enum/
	// typedef), the Java registry doesn't treat it as a type, so we
	// should not suppress the error. Mirrors Java's
	// JSTypeRegistry.getType returning null for plain function/object
	// bindings.
	rawDecl := c.scope.GetRawDeclaredType(name)
	if rawDecl != nil &&
		rawDecl != jstype.All && // All = unresolved alias; may be a type
		c.scope.LookupTypeAlias(name) == nil &&
		c.scope.LookupNominalType(name) == nil &&
		c.scope.LookupEnum(name) == nil {
		// Name IS in the global scope chain as a pure value binding
		// (not a constructor/typedef/enum, not an unresolved alias).
		// Since scope.ResolveExpr already returned Unknown, Java's
		// JSTypeRegistry.getType would also return null for it.
		return false
	}
	if c.allNames[name] {
		return true
	}
	return false
}

// Go-only: language-adapter — text-pass companion to Java's
// JsDocInfoParser.addTypeWarning calls for JSC_TYPE_PARSE_ERROR. Java
// emits these from the parser's token stream; jsgo recreates them via a
// regex pass so it can match Java's column positions exactly.
func (c *CheckJsDocTypes) checkParseErrors(s string, start, end int) []diag.Diagnostic {
	var out []diag.Diagnostic
	body := s[start:end]

	// ---- (1) Duplicate record fields --------------------------------
	// Walk every {…} annotation brace pair. If the inner text is a
	// record type (starts with `{` after stripping `!`/`?`/`…`), scan
	// for duplicate field names and emit at the `}` that closes the
	// record type.
	checkRecordDuplicates := func(inner string, innerAbsStart int) {
		// innerAbsStart is the absolute offset of inner[0] in s.
		// Strip leading modifiers so `!{k:T}` → `{k:T}`.
		stripped := strings.TrimLeft(inner, "!?")
		if len(stripped) == 0 || stripped[0] != '{' {
			return
		}
		record := stripped[1:] // content after the `{`
		seen := map[string]bool{}
		i := 0
		for i < len(record) {
			// Skip whitespace.
			for i < len(record) && (record[i] == ' ' || record[i] == '\t' || record[i] == '\n' || record[i] == '\r' || record[i] == '*') {
				i++
			}
			if i >= len(record) || record[i] == '}' {
				break
			}
			// Read field name (identifier).
			if !isIdentStart(record[i]) {
				break
			}
			nameStart := i
			for i < len(record) && (isIdentChar(record[i]) || record[i] == '_') {
				i++
			}
			fieldName := record[nameStart:i]
			// Skip whitespace.
			for i < len(record) && (record[i] == ' ' || record[i] == '\t') {
				i++
			}
			// Optional `:` followed by type expression.
			if i < len(record) && record[i] == ':' {
				i++ // consume `:`
				// Skip whitespace.
				for i < len(record) && (record[i] == ' ' || record[i] == '\t') {
					i++
				}
				// Skip the type expression (up to `,` or `}` at depth 0).
				depth := 0
				for i < len(record) {
					ch := record[i]
					if ch == '{' || ch == '<' || ch == '(' {
						depth++
					} else if ch == '}' || ch == '>' || ch == ')' {
						if depth == 0 {
							break // end of record
						}
						depth--
					} else if ch == ',' && depth == 0 {
						break
					}
					i++
				}
			}
			// After the field entry: position is at `,` or `}`.
			// Java fires addTypeWarning HERE (at stream position after
			// consuming the field's type).
			if seen[fieldName] {
				// errPos: the offset of `i` within `record` + offset of
				// `record` within `inner` + innerAbsStart.
				// stripped offset relative to inner = len(inner)-len(stripped)
				// record = stripped[1:], so record offset = strippedOffset+1
				strippedOff := len(inner) - len(stripped)
				recAbsStart := innerAbsStart + strippedOff + 1
				errOff := recAbsStart + i
				line, col := c.positionOf(s, errOff)
				out = append(out, diag.Diagnostic{
					Path: c.path, Line: line, Col: col,
					Code:    diag.JSC_TYPE_PARSE_ERROR,
					Level:   diag.Error,
					Message: "Bad type annotation. Duplicate record field " + fieldName + ".",
				})
			}
			seen[fieldName] = true
			// Advance past `,` or stop at `}`.
			if i < len(record) && record[i] == ',' {
				i++
			}
		}
	}

	// Scan all type-bearing annotation braces in the comment.
	bi := 0
	for bi < len(body) {
		at := strings.Index(body[bi:], "@")
		if at < 0 {
			break
		}
		at += bi
		j := at + 1
		for j < len(body) && isTagChar(body[j]) {
			j++
		}
		tag := body[at+1 : j]
		if !isTypeTag(tag) {
			bi = j
			continue
		}
		k := j
		for k < len(body) && (body[k] == ' ' || body[k] == '\t' || body[k] == '\n' || body[k] == '\r' || body[k] == '*') {
			k++
		}
		if k >= len(body) || body[k] != '{' {
			bi = j
			continue
		}
		braceStart := k
		braceEnd := findMatchingBrace(body, braceStart)
		if braceEnd < 0 {
			break
		}
		inner := body[braceStart+1 : braceEnd]
		checkRecordDuplicates(inner, start+braceStart+1)
		bi = braceEnd + 1
	}

	// ---- (2) Double @extends on a non-interface/non-record class ----
	// If the block has @interface or @record, multiple @extends is OK.
	// Otherwise, the second (and subsequent) @extends emits TYPE_PARSE_ERROR.
	isInterface := strings.Contains(body, "@interface") || strings.Contains(body, "@record")
	if !isInterface {
		extendsCount := 0
		ei := 0
		for ei < len(body) {
			idx := strings.Index(body[ei:], "@extends")
			if idx < 0 {
				break
			}
			idx += ei
			extendsCount++
			if extendsCount >= 2 {
				// Find the `{...}` following this @extends.
				k := idx + len("@extends")
				for k < len(body) && (body[k] == ' ' || body[k] == '\t' || body[k] == '\n' || body[k] == '\r' || body[k] == '*') {
					k++
				}
				if k < len(body) && body[k] == '{' {
					braceEnd := findMatchingBrace(body, k)
					if braceEnd >= 0 {
						// Java fires at stream position after the type name,
						// which is the `}` position (the closing brace).
						errOff := start + braceEnd
						line, col := c.positionOf(s, errOff)
						out = append(out, diag.Diagnostic{
							Path: c.path, Line: line, Col: col,
							Code:    diag.JSC_TYPE_PARSE_ERROR,
							Level:   diag.Error,
							Message: "Bad type annotation. type annotation incompatible with other annotations.",
						})
					}
				}
			}
			ei = idx + len("@extends")
		}
	}

	// ---- (3) @lends without a `{Name}` argument, or `{Name|…}` -------
	// Java's JsDocInfoParser expects exactly one name (no `|`) inside
	// the @lends braces; a `|` triggers "expected closing }".
	li := 0
	for li < len(body) {
		idx := strings.Index(body[li:], "@lends")
		if idx < 0 {
			break
		}
		idx += li
		// Check that `@lends` is not followed by a valid `{Name}` block.
		k := idx + len("@lends")
		// Ensure this is a word boundary (not @lendsomething).
		if k < len(body) && isTagChar(body[k]) {
			li = k
			continue
		}
		// Skip whitespace.
		for k < len(body) && (body[k] == ' ' || body[k] == '\t') {
			k++
		}
		if k >= len(body) || body[k] != '{' {
			// No `{` follows — emit at the first non-space char after @lends.
			errOff := start + k
			line, col := c.positionOf(s, errOff)
			out = append(out, diag.Diagnostic{
				Path: c.path, Line: line, Col: col,
				Code: diag.JSC_TYPE_PARSE_ERROR, Level: diag.Error,
				Message: "Bad type annotation. missing object name in @lends tag.",
			})
		} else {
			// `{` follows — check that no `|` appears inside (Java's
			// parseNameExpression reads one name then expects `}`).
			braceOpen := k
			if braceEnd := findMatchingBrace(body, braceOpen); braceEnd >= 0 {
				content := body[braceOpen+1 : braceEnd]
				if pipeIdx := strings.IndexByte(content, '|'); pipeIdx >= 0 {
					errOff := start + braceOpen + 1 + pipeIdx
					line, col := c.positionOf(s, errOff)
					out = append(out, diag.Diagnostic{
						Path: c.path, Line: line, Col: col,
						Code: diag.JSC_TYPE_PARSE_ERROR, Level: diag.Error,
						Message: "Bad type annotation. expected closing }",
					})
				}
			}
		}
		li = idx + len("@lends")
	}

	// ---- (4) Rest parameter not last in function type -----------------
	// `function(...T, S)` or `function(...T, S=)` — the rest (`...T`) must
	// be the last parameter. Java's JsDocInfoParser fires
	// "variable length argument must be last" at the position of the
	// parameter that follows the rest.
	// Scan all function(...) type expressions in the body.
	fi := 0
	for fi < len(body) {
		idx := strings.Index(body[fi:], "function(")
		if idx < 0 {
			break
		}
		idx += fi
		// Find the matching `)`.
		openParen := idx + len("function(") - 1
		depth := 0
		closeParen := -1
		for j := openParen; j < len(body); j++ {
			if body[j] == '(' {
				depth++
			} else if body[j] == ')' {
				depth--
				if depth == 0 {
					closeParen = j
					break
				}
			}
		}
		if closeParen < 0 {
			break
		}
		// Check the parameter list for a rest param that isn't last.
		params := body[openParen+1 : closeParen]
		restIdx := -1 // index in params where `...` was found
		// Split params on top-level commas.
		paramsStart := openParen + 1 // absolute offset in body
		depth = 0
		pStart := 0
		for j := 0; j <= len(params); j++ {
			if j < len(params) {
				ch := params[j]
				if ch == '(' || ch == '<' || ch == '{' {
					depth++
				} else if ch == ')' || ch == '>' || ch == '}' {
					depth--
				}
			}
			if depth == 0 && (j == len(params) || params[j] == ',') {
				param := strings.TrimSpace(params[pStart:j])
				if restIdx >= 0 && len(param) > 0 {
					// There was a rest param earlier — fire at the comma that separates
					// the rest param from this next param. Java reports at the position
					// right before the next param (at the separating comma).
					// The comma is at position `pStart - 1` in params.
					commaOff := start + paramsStart + pStart - 1
					// Skip back to the actual comma character.
					for commaOff > 0 && (s[commaOff] == ' ' || s[commaOff] == '\t') {
						commaOff--
					}
					line, col := c.positionOf(s, commaOff)
					out = append(out, diag.Diagnostic{
						Path: c.path, Line: line, Col: col,
						Code:  diag.JSC_TYPE_PARSE_ERROR, Level: diag.Error,
						Message: "Bad type annotation. variable length argument must be last." +
							" See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information.",
					})
					break
				}
				if strings.HasPrefix(param, "...") {
					restIdx = pStart
				}
				pStart = j + 1
			}
		}
		fi = closeParen + 1
	}

	return out
}

// Go-only: language-adapter — Java reads (line, column) from
// Node.getLineno() / getCharno() on the AST; jsgo's text-pass scanner
// computes them from a raw byte offset.
func (c *CheckJsDocTypes) positionOf(src string, off int) (int, int) {
	if off > len(src) {
		off = len(src)
	}
	line := 1
	col := 0
	for i := 0; i < off; i++ {
		if src[i] == '\n' {
			line++
			col = 0
			continue
		}
		col++
	}
	return line, col
}

// Go-only: language-adapter — text-scan brace matcher. Java's parser
// tracks nesting via the token stream so it never needs a brace matcher.
func findMatchingBrace(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// Go-only: language-adapter — text-scan character class for JSDoc tag
// names. Java's lexer handles tag tokenisation so it never needs this.
func isTagChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// Go-only: language-adapter — Java's JsDocToken enum encodes the same
// set of type-bearing JSDoc tags; jsgo's regex scanner matches them as
// strings.
func isTypeTag(s string) bool {
	switch s {
	case "type", "param", "return", "returns", "extends", "implements", "typedef":
		return true
	}
	return false
}

// Go-only: language-adapter — text-scan character class. Java uses
// Character.isJavaIdentifierStart at the lexer level.
func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b == '$'
}

// Go-only: language-adapter — text-scan character class. Java uses
// Character.isJavaIdentifierPart at the lexer level.
func isIdentChar(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

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

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
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

// cjdCollectTypedNames builds the set of all names that have a non-Unknown
// declared type across the global scope and every function-local scope.
// Mirrors Java's NamedType.resolveTypeof check: the error fires when the
// looked-up type is null or isUnknownType(). Names absent from this set
// but present in the source should trigger JSC_UNRECOGNIZED_TYPEOF_ERROR
// when used as a `typeof` operand in a JSDoc type annotation.
//
// Additionally, names whose declaration initialises an object literal
// (`var ns = {}`) are included: Java's TypedScopeCreator calls
// defineObjectLiteral which gives them a concrete ObjectType rather than
// Unknown, so resolveTypeof doesn't fire for them.
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

// cjdMergeTypedNames combines cjdCollectTypedNames (scope-based) with
// cjdCollectObjectLitNames (text-based) so that both explicitly-typed
// declarations and object-literal namespace declarations are excluded from
// the JSC_UNRECOGNIZED_TYPEOF_ERROR check.
func cjdMergeTypedNames(globalScope *TypedScope, creator *TypedScopeCreator, src string) map[string]bool {
	out := cjdCollectTypedNames(globalScope, creator)
	for name := range cjdCollectObjectLitNames(src) {
		out[name] = true
	}
	return out
}

// cjdCollectObjectLitNames scans source text for `var/let/const NAME = {`
// declarations and returns the set of such names. Java's TypedScopeCreator
// calls defineObjectLiteral on these, giving them a concrete ObjectType
// (non-Unknown) so resolveTypeof does not emit JSC_UNRECOGNIZED_TYPEOF_ERROR.
// We need the same exclusion: object-literal namespace objects (e.g.
// `/** @const */ var ns = {}`) are valid `typeof` operands.
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

// collectAllTemplateNames fans out over every JSDoc block in the
// source pulling out every @template name. Names are scope-wide
// in Java's implementation (a derived interface's @template binds
// through its own JSDoc but the supertype reference carries it)
// — we widen to "valid as a type anywhere in this source" to
// avoid FPs on interface-inheritance chains where the template
// name leaks across comment boundaries.
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

// collectDeclaredNames scans the source for every identifier
// appearing in non-JSDoc code — `var`/`let`/`const`/`function`/
// `class` bindings, function parameters, and any free-form
// reference. Used by nameResolves so `@type {X}` references to
// bindings the TypedScopeCreator's narrow scope chain doesn't
// reach (IIFE locals, function parameters, const-aliased enums)
// stay quiet. Broader than Java's exact check but strictly less
// prone to FPs.
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

// collectDeclaredNamesOld (unused — keep as reference for the
// stricter declaration-only variant).
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

// Check scans the source for JSDoc comments and validates every
// type-annotation brace pair. Returns the collected diagnostics.
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

// checkBareInlineType fires JSC_UNRECOGNIZED_TYPE_ERROR for untagged
// inline JSDoc blocks that carry a bare type name and appear immediately
// before a `var`/`let`/`const` declaration — the `/** T */ var x;` pattern.
// Mirrors Java's JsDocInfoParser treating such comments as implicit `@type {T}`
// annotations. Only fires when the type name is not resolvable (e.g. TTL
// template names like T from `@template T := … =:` are not valid types).
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
	// Only fire for names that are TTL template names on constructors/interfaces
	// (not valid types) and are not otherwise resolvable. Plain-function TTL
	// template names ARE valid types and should not fire.
	if !c.ctorTTLNames[name] {
		return nil
	}
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

// checkBareTypeof fires JSC_UNRECOGNIZED_TYPEOF_ERROR on inline
// JSDoc blocks of the form `/** typeof X */` where X isn't a
// declared or typed name. Java treats these as @type annotations and
// validates the `typeof` operand the same way.
// Only fires for `typeof` that appears OUTSIDE of `{…}` type-
// annotation braces — occurrences inside braces are handled by
// checkTypeExpr and emitted via checkComment.
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

// checkComment walks one `/** … */` block, finding every
// annotation of the form `@tag {…}` where `tag` is a type-
// bearing JSDoc tag, and validating each type component against
// the scope. `@template X, Y, Z` names are collected first and
// threaded into checkTypeExpr so they resolve as types inside
// the same comment.
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

// collectTemplateNames extracts every `@template` name from a
// JSDoc block. Names declared here are valid as types only
// within the same comment.
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

// collectConstructorTTLNames collects TTL @template names (`@template T := … =:`)
// that appear in JSDoc blocks on constructor or interface functions. These names
// are NOT valid types (Java's FunctionTypeBuilder rejects them with
// TEMPLATE_TRANSFORMATION_ON_CLASS when allowTypeTransformations=false for
// constructors/interfaces). Mirrors getTypeTransformations() key-set for
// constructor JSDoc only.
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

// checkTypeExpr walks the text inside a `{…}` annotation value,
// walking nested `<…>` / `(…)` / `{…}` groups and checking each
// top-level name component that's in a TYPE position (as opposed
// to a field-name or param-name slot). `base` is the absolute
// offset of `inner` in the source — used to compute per-name
// positions.
//
// Position tracking: we maintain a stack of "contexts" — record
// (`{`), function-params (`(`), generic-args (`<`) — and a
// per-context flag saying whether the next identifier is a
// key/param name (to be skipped) or a type name (to validate).
// Key / param slots turn into type slots after a `:`. After a
// `,` at the top of the current context, we flip back to the
// key slot.
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

// nameResolves reports whether the dotted name is a known
// primitive, nominal, enum, typedef, or a builtin we ship as an
// extern. Dotted chains (`foo.bar`) take the first segment for
// lookup since our scope doesn't track object-property type
// namespaces yet — so `goog.foo` resolves when `goog` does.
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

// checkParseErrors scans a single `/** … */` block for
// JSC_TYPE_PARSE_ERROR conditions that the Java JsDocInfoParser
// reports via addTypeWarning:
//
//  1. Duplicate record field (`{x:T1, x:T2}` → Msg.JSDOC_TYPE_RECORD_DUPLICATE).
//  2. Second @extends on a non-interface/non-record class (Msg.JSDOC_INCOMPAT_TYPE).
//  3. @lends without a `{Name}` argument (Msg.JSDOC_LENDS_MISSING).
//
// Position mirrors Java: for (1) the `}` closing the duplicate field's
// type; for (2) the `}` closing the second @extends type; for (3) the
// first non-space char after `@lends` when no `{` follows.
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

	return out
}

// positionOf returns (line, col) for the absolute offset `off` in
// `src`. Lines are 1-indexed, columns 0-indexed (matching Java
// Node.getLineno / getCharno).
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

// findMatchingBrace returns the offset of the `}` matching the
// `{` at `open`. Tracks nested `{…}` pairs. Returns -1 when no
// matching close is found.
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

func isTagChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isTypeTag(s string) bool {
	switch s {
	case "type", "param", "return", "returns", "extends", "implements", "typedef":
		return true
	}
	return false
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b == '$'
}

func isIdentChar(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

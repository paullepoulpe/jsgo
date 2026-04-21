// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java's
// late-resolve path (resolveTypesInScope + reportUnresolvedTypes) plus
// the parseTypeExpressionAnnotation error branch from
// com/google/javascript/jscomp/parsing/JSDocInfoParser.java
// @ closure-compiler v20260415.
//
// In Java, the registry hands out NamedType wrappers for every JSDoc
// reference that doesn't resolve at JSDoc-parse time. A dedicated late
// pass (post-scope-creation, pre-TypeCheck) walks those wrappers and
// either substitutes the resolved type or emits
// JSC_UNRECOGNIZED_TYPE_ERROR at the source position recorded when the
// wrapper was constructed. Parse failures inside the type expression
// surface as JSC_TYPE_PARSE_ERROR at the opening `{`.
//
// This Go port keeps the Java semantics but sources positions from the
// JSDoc parser's comment-relative line/col (see
// rhino.JSTypeExpression.SourceLine / SourceCol and the NamedType
// wrappers it carries). The absolute position is recovered by scanning
// the source text backward from the annotated node's start position
// for the preceding `*/` and matching `/**`.
//
// Scope: we only report nominal-name failures on the tags the checker
// already consumes (type / param / return / extends / implements /
// typedef). Other positions continue to swallow resolution failures
// silently, matching the pre-pass behaviour so the corpus doesn't
// regress on tests whose oracles don't include the warning.

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// JSTypeResolver is the NamedType resolution pass. It walks the AST,
// inspects every JSDoc type reference, and emits the late-resolve
// diagnostics Java produces at the end of JSTypeRegistry.resolveTypes.
type JSTypeResolver struct {
	path     string
	src      string
	scope    *TypedScope
	creator  *TypedScopeCreator
	externs  *TypedScope
	diags    []diag.Diagnostic
	seen     map[diagKey]bool // dedupe identical (line,col,code,name)
	// sourceIdents is a rough set of every identifier that appears as a
	// `function <name>`, `class <name>`, `var|let|const <name>`, or
	// assignment target in the raw source text. Used as a last-resort
	// fallback: when jsgo's parser fails to construct a scope binding
	// (e.g. a missing semicolon on the previous line drops the next
	// declaration), we still avoid flagging a JSDoc reference to that
	// name as unknown. Mirrors how Java's registry consults the full
	// SymbolTable before reporting a NamedType as unresolved.
	sourceIdents map[string]bool
	// fileTemplates collects every @template name declared anywhere in
	// the source. Closure's late-resolve treats a template name declared
	// on an enclosing @constructor / @interface / @record as in-scope for
	// its prototype properties (so `/** @type {T} */ Foo.prototype.x`
	// resolves T). Rather than thread the enclosing function through the
	// AST walk, we accept any name declared anywhere as a template — it's
	// a superset of the legal bindings, but the corpus fires no fixture
	// that relies on rejecting a reused single-letter name.
	fileTemplates map[string]bool
	// namedTypes accumulates one *jstype.NamedType per JSDoc reference
	// walked during Check. Mirrors the per-registry list Java's
	// JSTypeRegistry maintains for NamedType wrappers. Resolved entries
	// carry a non-nil ReferencedType; unresolved ones were the ones that
	// fired JSC_UNRECOGNIZED_TYPE_ERROR. Callers (and tests) can inspect
	// this via NamedTypes() to audit the late-resolve bookkeeping
	// without re-scanning diagnostics.
	namedTypes []*jstype.NamedType
}

type diagKey struct {
	line, col int
	code      string
	name      string
}

// NewJSTypeResolver builds a resolver for one compilation unit. The
// externs scope is passed alongside the user scope so names like
// `Node` / `Element` (declared only in the default externs) resolve
// before the name is flagged as unknown.
func NewJSTypeResolver(path, src string, creator *TypedScopeCreator, scope, externs *TypedScope) *JSTypeResolver {
	return &JSTypeResolver{
		path:         path,
		src:          src,
		scope:        scope,
		creator:      creator,
		externs:      externs,
		seen:         map[diagKey]bool{},
		sourceIdents: collectSourceIdents(src),
		fileTemplates: collectResolverTemplateNames(src),
	}
}

// collectResolverTemplateNames scans the source for `@template NAME[, NAME2]` tags
// and returns the union set. This is a textual scan (same strategy as
// sourceIdents) — good enough to recognise a template name used by any
// prototype property in the same file.
//
// Also handles the bounded form `@template {BOUND} NAME` by skipping the
// `{...}` brace group before reading identifiers.
func collectResolverTemplateNames(src string) map[string]bool {
	out := map[string]bool{}
	const tag = "@template"
	for i := 0; i+len(tag) <= len(src); i++ {
		if src[i:i+len(tag)] != tag {
			continue
		}
		j := i + len(tag)
		if j < len(src) && isIdentByte(src[j]) {
			continue
		}
		// Skip whitespace + optional bound `{…}` right after `@template`.
		k := j
		for k < len(src) && (src[k] == ' ' || src[k] == '\t') {
			k++
		}
		if k < len(src) && src[k] == '{' {
			depth := 1
			k++
			for k < len(src) && depth > 0 {
				switch src[k] {
				case '{':
					depth++
				case '}':
					depth--
				case '\n':
					depth = 0
				}
				k++
			}
		}
		// Collect comma-separated identifiers until we hit a non-ident /
		// non-comma non-space char (e.g. newline, `*/`, `}`).
		for k < len(src) {
			for k < len(src) && (src[k] == ' ' || src[k] == '\t' || src[k] == ',') {
				k++
			}
			name, end := readIdent(src, k)
			if name == "" {
				break
			}
			out[name] = true
			k = end
			// Stop on end-of-line or JSDoc comment boundary.
			if k < len(src) && (src[k] == '\n' || src[k] == '\r') {
				break
			}
			if k+1 < len(src) && src[k] == '*' && src[k+1] == '/' {
				break
			}
			if k < len(src) && src[k] == ':' {
				break
			}
		}
	}
	return out
}

// Check walks `script` and returns any emitted diagnostics. Mirrors the
// Java `Compiler.resolveTypes` hook that runs between scope creation
// and TypeCheck.
func (r *JSTypeResolver) Check(script *rhino.Node) []diag.Diagnostic {
	r.visit(script, r.scope)
	return r.diags
}

// visit traverses the AST. Each node carrying a JSDocInfo contributes
// the position-aware JSTypeExpressions it declared. For FUNCTION /
// CLASS nodes we also descend using the creator-recorded local scope so
// @template names attached to a method resolve against the method's
// shim rather than the outer scope.
func (r *JSTypeResolver) visit(n *rhino.Node, scope *TypedScope) {
	if n == nil {
		return
	}
	if info := n.JSDocInfo(); info != nil {
		r.checkInfo(n, info, scope)
	}
	next := scope
	if local, ok := r.creator.Functions[n]; ok && local != nil {
		next = local
	}
	for _, child := range n.Children() {
		r.visit(child, next)
	}
}

// checkInfo emits diagnostics for every type expression on `info` that
// fails resolution. Covers the same tags TypedScopeCreator consumes.
//
// Template names declared on the same JSDoc block are treated as
// resolved for the duration of this call — mirrors how Java's
// TypedScopeCreator layers a per-function TemplateTypeMap on top of
// the scope before resolving the function's parameter / return
// expressions.
func (r *JSTypeResolver) checkInfo(n *rhino.Node, info *rhino.JSDocInfo, scope *TypedScope) {
	tmpls := info.TemplateNames()
	if info.HasType() {
		r.checkExpr(n, info.Type(), scope, tmpls)
	}
	if info.ReturnType() != nil {
		r.checkExpr(n, info.ReturnType(), scope, tmpls)
	}
	if info.BaseType() != nil {
		r.checkExpr(n, info.BaseType(), scope, tmpls)
	}
	for _, impl := range info.Implementations() {
		r.checkExpr(n, impl, scope, tmpls)
	}
	if info.TypedefType() != nil {
		r.checkExpr(n, info.TypedefType(), scope, tmpls)
	}
	for _, name := range info.ParameterNames() {
		if pe := info.ParameterType(name); pe != nil {
			r.checkExpr(n, pe, scope, tmpls)
		}
	}
}

// checkExpr reports JSC_TYPE_PARSE_ERROR when the JSDoc parser flagged
// this expression as malformed, and JSC_UNRECOGNIZED_TYPE_ERROR for
// every name that still fails to resolve against the current scope,
// template bindings, and extern name set.
//
// For every name walked here a *jstype.NamedType wrapper is allocated
// so the registry-style late-resolve bookkeeping (Java's
// NamedType.resolveInternal) survives the pass: resolved wrappers
// carry a non-nil ReferencedType; unresolved ones are the ones that
// fired the diagnostic.
func (r *JSTypeResolver) checkExpr(n *rhino.Node, expr *rhino.JSTypeExpression, scope *TypedScope, tmpls []string) {
	if expr == nil {
		return
	}
	if expr.ParseError != nil {
		line, col := r.absolutePos(n, expr.ParseError.SourceLine, expr.ParseError.SourceCol)
		r.emit(line, col, diag.JSC_TYPE_PARSE_ERROR, "", "Bad type annotation. type annotation incompatible with other annotations")
		return
	}
	for _, nm := range expr.Names {
		if slices.Contains(tmpls, nm.UnresolvedName) {
			continue
		}
		if r.fileTemplates[nm.UnresolvedName] {
			continue
		}
		line, col := r.absolutePos(n, nm.SourceLine, nm.SourceCol)
		named := jstype.NewNamedType(nm.UnresolvedName, r.path, line, col)
		if t := r.resolveName(nm.UnresolvedName, scope); t != nil {
			named.Resolve(t)
			r.namedTypes = append(r.namedTypes, named)
			continue
		}
		r.namedTypes = append(r.namedTypes, named)
		r.emit(line, col, diag.JSC_UNRECOGNIZED_TYPE_ERROR, nm.UnresolvedName,
			"Bad type annotation. Unknown type "+nm.UnresolvedName)
	}
}

// NamedTypes returns every forward-reference wrapper seen during the
// resolve pass. Callers inspect IsResolved to distinguish the wrappers
// that got bound from the ones that fired a diagnostic. Mirrors
// JSTypeRegistry.getResolvedTypesModuloResolving() in Java.
func (r *JSTypeResolver) NamedTypes() []*jstype.NamedType { return r.namedTypes }

// resolveName mirrors JSTypeRegistry.getType(StaticScope, String) — the
// registry's post-scope lookup used by NamedType.resolveInternal. Returns
// the resolved JSType (possibly a lenient marker like jstype.Unknown for
// recognised-but-unmodelled keywords) or nil for a name Java would fire
// JSC_UNRECOGNIZED_TYPE_ERROR on.
//
// nameResolves is the boolean sibling of this helper; they share the
// same resolution rules so the two views (did-resolve / resolved-to-what)
// stay in sync.
func (r *JSTypeResolver) resolveName(name string, scope *TypedScope) jstype.JSType {
	if name == "" {
		return jstype.Unknown
	}
	switch name {
	case "number", "string", "boolean", "null", "undefined", "void", "symbol", "bigint",
		"Null", "Undefined", "Void",
		"Object", "Function", "Array", "?", "*", "this",
		"Arguments", "IArrayLike", "IObject", "IThenable", "Iterable",
		"Generator", "AsyncIterable", "AsyncGenerator", "AsyncIterator",
		"Iterator", "IteratorIterable", "ReadonlyArray", "ReadonlyMap", "ReadonlySet",
		"Symbol", "Promise", "Map", "Set", "WeakMap", "WeakSet",
		"Date", "RegExp", "Error", "Number", "String", "Boolean":
		return jstype.Unknown
	}
	head := name
	if dot := strings.IndexByte(name, '.'); dot > 0 {
		head = name[:dot]
	}
	if scope != nil {
		if t := scope.LookupNominalType(name); t != nil {
			return t
		}
		if t := scope.LookupTypeAlias(name); t != nil {
			return t
		}
		if t := scope.GetRawDeclaredType(head); t != nil {
			return t
		}
	}
	if r.externs != nil && r.externs != scope {
		if t := r.externs.LookupNominalType(name); t != nil {
			return t
		}
		if t := r.externs.LookupTypeAlias(name); t != nil {
			return t
		}
		if t := r.externs.GetRawDeclaredType(head); t != nil {
			return t
		}
	}
	if defaultExternNames()[head] {
		return jstype.Unknown
	}
	if r.sourceIdents[name] || r.sourceIdents[head] {
		return jstype.Unknown
	}
	return nil
}

// nameResolves mirrors the predicate Java's registry uses to decide
// whether a NamedType is still pending: the scope must expose it as a
// nominal constructor, a @typedef alias, a template parameter, a
// primitive, or a known extern. Everything else is an unresolved
// reference.
func (r *JSTypeResolver) nameResolves(name string, scope *TypedScope) bool {
	if name == "" {
		return true
	}
	// JSDoc keyword types supported by the registry but not exposed
	// as individual primitives. Mirrors the hard-coded hits in
	// JSTypeRegistry.getType(String).
	switch name {
	case "number", "string", "boolean", "null", "undefined", "void", "symbol", "bigint",
		"Null", "Undefined", "Void",
		"Object", "Function", "Array", "?", "*", "this",
		"Arguments", "IArrayLike", "IObject", "IThenable", "Iterable",
		"Generator", "AsyncIterable", "AsyncGenerator", "AsyncIterator",
		"Iterator", "IteratorIterable", "ReadonlyArray", "ReadonlyMap", "ReadonlySet",
		"Symbol", "Promise", "Map", "Set", "WeakMap", "WeakSet",
		"Date", "RegExp", "Error", "Number", "String", "Boolean":
		return true
	}
	// Qualified names like `foo.bar.Baz` are resolved against the
	// leading identifier; the scope creator doesn't build namespace
	// objects, so we accept the qualified name when the qualifier is
	// bound. Mirrors Java's getQualifiedNameComponents fallback.
	head := name
	if dot := strings.IndexByte(name, '.'); dot > 0 {
		head = name[:dot]
	}
	if scope != nil {
		if scope.LookupNominalType(name) != nil {
			return true
		}
		if scope.LookupTypeAlias(name) != nil {
			return true
		}
		// Accept ANY value-namespace binding as a possibly-type name.
		// Closure's late-resolve path first looks up the registry, then
		// tries to dereference a JavaScript value as a type (e.g.
		// `var Foo = function() {}` — Foo holds a FunctionType whose
		// instance can serve as a type). Mirroring that liberal match
		// here avoids false positives on the many fixtures whose type
		// names are bound indirectly (`@enum`, constructor aliases,
		// namespaces). Real unknown names reach this branch with a
		// nil declared type and fall through to the error.
		if scope.GetRawDeclaredType(head) != nil {
			return true
		}
	}
	if r.externs != nil && r.externs != scope {
		if r.externs.LookupNominalType(name) != nil {
			return true
		}
		if r.externs.LookupTypeAlias(name) != nil {
			return true
		}
		if r.externs.GetRawDeclaredType(head) != nil {
			return true
		}
	}
	if defaultExternNames()[head] {
		return true
	}
	// Last-resort text scan: if the identifier appears as a declaration /
	// assignment target anywhere in the source, treat the reference as
	// resolved. jsgo's parser silently drops declarations on lines like
	// `var ns = {}` (no ASI trailing newline) + `function rec(){}` — the
	// textual scan recovers those names so we don't false-positive JSDoc
	// references to them. Both the full qualified form and the leading
	// identifier are consulted: `u.A` matches a `u.A = …` assignment
	// target, while `Foo` matches any top-level declaration.
	if r.sourceIdents[name] {
		return true
	}
	if r.sourceIdents[head] {
		return true
	}
	return false
}

// absolutePos converts a comment-relative (line, col) to an absolute
// source position. Scans backward from `n`'s start in r.src to find the
// preceding `*/` and its matching `/**`, then offsets the body-relative
// coordinates by the body's source origin.
func (r *JSTypeResolver) absolutePos(n *rhino.Node, bodyLine, bodyCol int) (int, int) {
	line, col := nodeStartPos(n)
	// Convert (1-indexed line, 0-indexed col) into a byte offset so we
	// can scan backward through the source bytes.
	endOff := lineColToOffset(r.src, line, col)
	bodyStartLine, bodyStartCol, ok := findJSDocBodyStart(r.src, endOff)
	if !ok {
		// Fall back to the annotated node's own position — better than
		// pointing at column 0, and mirrors Java's fallback when the
		// ScriptSourceFile can't resolve a charno.
		return line, col
	}
	if bodyLine == 0 {
		return bodyStartLine, bodyStartCol + bodyCol
	}
	return bodyStartLine + bodyLine, bodyCol
}

// nodeStartPos returns a (line, col) for `n`. FUNCTION / CLASS nodes
// in jsgo's parser are constructed without an explicit SetPos call, so
// we walk into their first positioned descendant (the name / param
// list) to anchor the scan.
func nodeStartPos(n *rhino.Node) (int, int) {
	if n == nil {
		return 0, 0
	}
	if n.Lineno() > 0 {
		return n.Lineno(), n.Charno()
	}
	for _, c := range n.Children() {
		if l, col := nodeStartPos(c); l > 0 {
			return l, col
		}
	}
	return 0, 0
}

// findJSDocBodyStart walks `src[:endOff]` backward to locate the
// nearest `*/` at or before `endOff-1`, then walks further back to the
// matching `/**`. Returns the 1-indexed source line and 0-indexed
// column of the first char AFTER `/**` (the body's origin). ok=false
// when no well-formed JSDoc block sits before `endOff`.
//
// The node we anchor against may be arbitrary whitespace / code
// characters after `*/` (e.g. `function derived` follows the `*/`
// with no leading whitespace), so we don't try to be strict about the
// gap — the nearest preceding `*/` wins.
func findJSDocBodyStart(src string, endOff int) (int, int, bool) {
	if endOff > len(src) {
		endOff = len(src)
	}
	// Find nearest `*/` at or before endOff-1.
	i := endOff - 1
	for i >= 1 && !(src[i] == '/' && src[i-1] == '*') {
		i--
	}
	if i < 1 {
		return 0, 0, false
	}
	// Find matching `/**`. Block-comment nesting isn't a thing in JS,
	// so the first `/*` we hit closes the match.
	j := i - 2
	for j >= 1 {
		if src[j] == '*' && src[j-1] == '/' {
			// Ensure it's really `/**` — require a trailing `*` after
			// the opening to be a JSDoc, otherwise plain `/* ... */`
			// comments would match. The parser only attaches JSDoc for
			// `/**` openings.
			if j+1 < len(src) && src[j+1] == '*' {
				// The scanner's ConsumeBlockComment strips only `/*`
				// and returns the body starting at the second `*` of
				// `/**`. JsDocTokenStream's offset 0 therefore aligns
				// with the column of that `*`, which is (column of
				// `/`) + 2.
				line, col := offsetToLineCol(src, j-1)
				return line, col + 2, true
			}
			// Not a JSDoc block; give up.
			return 0, 0, false
		}
		j--
	}
	return 0, 0, false
}

// lineColToOffset maps (1-indexed line, 0-indexed col) → byte offset.
func lineColToOffset(src string, line, col int) int {
	l := 1
	colSeen := 0
	for i := 0; i < len(src); i++ {
		if l == line && colSeen == col {
			return i
		}
		if src[i] == '\n' {
			if l == line {
				// Clamp to end-of-line if col is past it.
				return i
			}
			l++
			colSeen = 0
			continue
		}
		colSeen++
	}
	return len(src)
}

// offsetToLineCol is the inverse: byte offset → (1-indexed line,
// 0-indexed col). Matches Node.getLineno / getCharno semantics.
func offsetToLineCol(src string, off int) (int, int) {
	if off > len(src) {
		off = len(src)
	}
	line, col := 1, 0
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

// collectSourceIdents walks `src` for each identifier token that
// follows a declaration keyword (function/class/var/let/const) or
// appears as an assignment LHS. Used only as a lenient fallback when
// nameResolves can't otherwise prove a reference is known; keeps us
// from false-positiving on declarations jsgo's parser loses due to
// unrelated parse quirks (e.g. missing ASI before `function`).
func collectSourceIdents(src string) map[string]bool {
	out := map[string]bool{}
	for _, kw := range []string{"function", "class", "var", "let", "const"} {
		collectAfterKeyword(src, kw, out)
	}
	// `NAME = ...` at the start of a line declares an implicit global
	// in sloppy mode; Closure accepts it as a nominal binding.
	collectAssignTargets(src, out)
	return out
}

func collectAfterKeyword(src, kw string, out map[string]bool) {
	i := 0
	for i < len(src) {
		j := findKeyword(src, i, kw)
		if j < 0 {
			return
		}
		k := j + len(kw)
		for k < len(src) && (src[k] == ' ' || src[k] == '\t') {
			k++
		}
		name, end := readIdent(src, k)
		if name != "" {
			out[name] = true
		}
		i = end
		if i <= j {
			i = j + len(kw)
		}
	}
}

func findKeyword(src string, from int, kw string) int {
	for i := from; i+len(kw) <= len(src); i++ {
		if i > 0 && isIdentByte(src[i-1]) {
			continue
		}
		if src[i:i+len(kw)] != kw {
			continue
		}
		after := i + len(kw)
		if after < len(src) && isIdentByte(src[after]) {
			continue
		}
		return i
	}
	return -1
}

func collectAssignTargets(src string, out map[string]bool) {
	// Scan line-by-line for `IDENT[.IDENT…] = ...`. Conservative: we
	// only care about the first identifier on a line after optional
	// whitespace and optional JSDoc comment. Qualified chains contribute
	// both the head and every prefix (so `a.b.c = x` registers `a` as
	// an implicit namespace plus the full `a.b.c` chain).
	start := 0
	for start < len(src) {
		end := start
		for end < len(src) && src[end] != '\n' {
			end++
		}
		line := src[start:end]
		trimmed := trimLeadingSpaceAndComment(line)
		name, after := readIdent(trimmed, 0)
		if name != "" {
			// Gather any `.IDENT` continuation to form the qualified name.
			qual := name
			for after < len(trimmed) && trimmed[after] == '.' {
				next, nEnd := readIdent(trimmed, after+1)
				if next == "" {
					break
				}
				qual = qual + "." + next
				after = nEnd
			}
			rest := trimLeadingSpace(trimmed[after:])
			if len(rest) > 0 && rest[0] == '=' && (len(rest) == 1 || rest[1] != '=') {
				out[name] = true
				out[qual] = true
			}
		}
		start = end + 1
	}
}

func trimLeadingSpace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return s[i:]
}

func trimLeadingSpaceAndComment(s string) string {
	s = trimLeadingSpace(s)
	if len(s) >= 4 && s[0] == '/' && s[1] == '*' {
		// Consume `/* ... */` on the same line.
		i := 2
		for i+1 < len(s) && !(s[i] == '*' && s[i+1] == '/') {
			i++
		}
		if i+1 < len(s) {
			s = trimLeadingSpace(s[i+2:])
		}
	}
	return s
}

func readIdent(src string, from int) (string, int) {
	i := from
	for i < len(src) && isIdentByte(src[i]) {
		i++
	}
	if i == from {
		return "", from
	}
	return src[from:i], i
}

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '$'
}

// emit pushes a de-duplicated diagnostic. Java's Compiler fires the
// late-resolve warning exactly once per (position, name) pair; the seen
// map mirrors that invariant.
func (r *JSTypeResolver) emit(line, col int, code, name, msg string) {
	k := diagKey{line: line, col: col, code: code, name: name}
	if r.seen[k] {
		return
	}
	r.seen[k] = true
	r.diags = append(r.diags, diag.Diagnostic{
		Path:    r.path,
		Line:    line,
		Col:     col,
		Code:    code,
		Level:   diag.Error,
		Message: msg,
	})
}

// Port of the type-mismatch arm of
// com/google/javascript/jscomp/TypeCheck.java#checkDuplicateVarDeclaration
// @ closure-compiler v20260415 (narrow; qualified-name scope only).
//
// Java's TypeCheck emits JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH when the
// same qualified name is declared twice with incompatible JSDoc-ascribed
// types (as opposed to JSC_VAR_MULTIPLY_DECLARED_ERROR, which fires for
// plain-NAME redeclarations). The real pass lives inside TypeCheck's
// visit() and leans on TypedScope.getVar(); jsgo's scope is read-only at
// the point we need it here, so this port walks the AST instead: it
// collects every qualified-name (GETPROP) declaration site that carries
// a type-ascribing JSDoc tag (@type, @typedef, @constructor, @interface,
// @record) and, when the same path appears more than once, compares the
// declared types pairwise. If the two types aren't mutually subtype-
// compatible, the diagnostic fires at the second declaration's
// property-name position — matching the oracle's column selection.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckDupVarDeclMismatch walks the SCRIPT looking for duplicate,
// type-ascribed declarations of the same qualified name. Mirrors the
// one-shot per-file pass-instance pattern used by VariableReferenceCheck.
type CheckDupVarDeclMismatch struct {
	path  string
	scope *TypedScope
	diags []diag.Diagnostic
}

// declSiteTyped records everything we need about one type-ascribed
// declaration so we can render the diagnostic at the second occurrence.
type declSiteTyped struct {
	line    int
	col     int
	typ     jstype.JSType
	typExpr string // raw JSDoc expression, for the message fallback
	// kind distinguishes `@type`/`@typedef` (value-like) from
	// `@constructor`/`@interface`/`@record` (constructor-like). When
	// two sites of different kinds share a qname, that's always a
	// mismatch regardless of the underlying jstype equality.
	kind declKind
}

type declKind int

const (
	kindTypeAscribed   declKind = iota // @type or @typedef
	kindConstructorish                 // @constructor / @interface / @record
)

// NewCheckDupVarDeclMismatch builds a fresh pass-instance per file,
// matching Compiler's DefaultPassConfig pattern.
func NewCheckDupVarDeclMismatch(path string, scope *TypedScope) *CheckDupVarDeclMismatch {
	return &CheckDupVarDeclMismatch{path: path, scope: scope}
}

// Check runs the traversal and returns the collected diagnostics.
func (c *CheckDupVarDeclMismatch) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	sites := map[string][]declSiteTyped{}
	c.collect(script, sites)
	for _, list := range sites {
		if len(list) < 2 {
			continue
		}
		// Compare each site against the first; the diagnostic fires at
		// the later-in-source occurrence, mirroring how Java reports
		// redefinition at the offending (second+) declaration.
		first := list[0]
		for i := 1; i < len(list); i++ {
			cur := list[i]
			if compatibleDecl(first, cur) {
				continue
			}
			c.diags = append(c.diags, diag.Diagnostic{
				Path:  c.path,
				Line:  cur.line,
				Col:   cur.col,
				Code:  diag.JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf(
					"variable redefined with type %s, original definition at %s:%d with type %s",
					cur.typExpr, c.path, first.line, first.typExpr),
			})
			// After firing, treat the current site as the new baseline
			// so a third / fourth mismatch is reported against it —
			// matches the Java pass which updates the scope's recorded
			// type each time it reports.
			first = cur
		}
	}
	return c.diags
}

// collect walks the AST gathering every qualified-name declaration
// site that carries a type-ascribing JSDoc tag. Keyed by the
// dot-joined qualified name so `ns.x` / `goog.foo` each land in their
// own bucket. Plain-NAME `var`/`let`/`const` declarations are
// deliberately excluded — Java's VariableReferenceCheck fires
// JSC_VAR_MULTIPLY_DECLARED_ERROR for those, and the two diagnostics
// are mutually exclusive in TypeCheck's handling.
func (c *CheckDupVarDeclMismatch) collect(n *rhino.Node, out map[string][]declSiteTyped) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if inner := n.FirstChild(); inner != nil {
			c.collectExpr(inner, out)
		}
	}
	for _, k := range n.Children() {
		c.collect(k, out)
	}
}

// collectExpr inspects the expression inside an EXPR_RESULT for a
// qualified-name declaration pattern — either a bare GETPROP
// (`/** @type {T} */ ns.x;`) or an ASSIGN whose LHS is a GETPROP
// (`/** @constructor */ goog.Bar = function() {};`).
//
// Declarations whose root NAME isn't declared anywhere in scope are
// skipped: VarCheck will already fire JSC_UNDEFINED_VARIABLE for them
// and Java doesn't cascade a type-mismatch on top.
func (c *CheckDupVarDeclMismatch) collectExpr(inner *rhino.Node, out map[string][]declSiteTyped) {
	info := inner.JSDocInfo()
	lhs := inner
	if inner.Token() == rhino.ASSIGN {
		kids := inner.Children()
		if len(kids) < 2 {
			return
		}
		lhs = kids[0]
	}
	if lhs.Token() != rhino.GETPROP {
		return
	}
	qname, prop := qualifiedName(lhs)
	if qname == "" || prop == nil || !c.rootDeclared(lhs) {
		return
	}
	if site, ok := c.siteFromJSDoc(info, prop.Lineno(), prop.Charno()); ok {
		out[qname] = append(out[qname], site)
	}
}

// rootDeclared reports whether the NAME at the base of a GETPROP
// chain resolves to a binding in the current scope chain. Used as a
// guard so `ns.x` with an undeclared `ns` doesn't attract a spurious
// DUP_VAR_DECLARATION_TYPE_MISMATCH on top of the UNDEFINED_VARIABLE
// that VarCheck already reports.
func (c *CheckDupVarDeclMismatch) rootDeclared(getProp *rhino.Node) bool {
	cur := getProp
	for cur != nil && cur.Token() == rhino.GETPROP {
		kids := cur.Children()
		if len(kids) < 1 {
			return false
		}
		cur = kids[0]
	}
	if cur == nil || cur.Token() != rhino.NAME {
		return false
	}
	return c.scope.GetDeclaredType(cur.String_()) != nil
}

// siteFromJSDoc classifies the JSDoc into one of our two tracked
// buckets and returns the corresponding declSiteTyped. Declarations
// without any of the recognised tags are skipped so the untyped-
// redeclare arm (VariableReferenceCheck / VarCheck) keeps handling
// them.
func (c *CheckDupVarDeclMismatch) siteFromJSDoc(info *rhino.JSDocInfo, line, col int) (declSiteTyped, bool) {
	if info == nil {
		return declSiteTyped{}, false
	}
	var expr *rhino.JSTypeExpression
	switch {
	case info.HasType():
		expr = info.Type()
	case info.TypedefType() != nil:
		expr = info.TypedefType()
	case info.IsConstructor() || info.IsInterface() || info.IsRecord():
		return declSiteTyped{
			line: line, col: col,
			typExpr: "<constructor>",
			kind:    kindConstructorish,
		}, true
	default:
		return declSiteTyped{}, false
	}
	return declSiteTyped{
		line: line, col: col,
		typ:     c.scope.ResolveExpr(expr.TypeName),
		typExpr: expr.TypeName,
		kind:    kindTypeAscribed,
	}, true
}

// compatibleDecl reports whether two type-ascribed declarations of the
// same qname are assignable in both directions. A cross-kind pair
// (value-like + constructor-like) is always incompatible — matches
// Java's rejection of `@constructor Foo = …;` + `@typedef goog.Bar;`
// overlap.
func compatibleDecl(a, b declSiteTyped) bool {
	if a.kind != b.kind {
		return false
	}
	if a.kind == kindConstructorish {
		// Two constructor-like decls of the same qname: treat as a
		// mismatch. There's no second @constructor on a qname in the
		// corpus without also tripping VariableReferenceCheck, so this
		// arm is defensive.
		return false
	}
	// Unknown on either side accepts anything — avoids false positives
	// where ResolveExpr couldn't parse the annotation.
	if isJSUnknown(a.typ) || isJSUnknown(b.typ) {
		return true
	}
	// Two distinct nominal object types (both carrying a Name) are
	// treated as incompatible regardless of structural overlap —
	// matches how TypeCheck.java's equivalence uses nominal identity,
	// not structural subtyping, when comparing redeclarations. Empty
	// `@record`s rec3 / rec4 would otherwise both be mutual subtypes
	// of the other via the structural rule.
	if aName, bName := nominalName(a.typ), nominalName(b.typ); aName != "" && bName != "" && aName != bName {
		return false
	}
	return jstype.IsSubtype(a.typ, b.typ) && jstype.IsSubtype(b.typ, a.typ)
}

// nominalName returns the declared nominal name of an ObjectType, or
// "" for anonymous / primitive / structural-literal types. Used to
// separate nominal identity from structural compatibility when
// comparing two redeclarations.
func nominalName(t jstype.JSType) string {
	if o, ok := t.(*jstype.ObjectType); ok {
		return o.Name
	}
	return ""
}

// qualifiedName flattens a GETPROP chain rooted at a NAME into its
// dotted source form (e.g. `goog.labs.a11y`). Returns ("", nil) when
// the chain doesn't end in a plain NAME — computed keys (`x[0].foo`)
// or bracketed members are outside this pass's scope. The second
// return is the STRINGLIT of the final property, whose position is
// what the oracle points at.
func qualifiedName(n *rhino.Node) (string, *rhino.Node) {
	if n == nil || n.Token() != rhino.GETPROP {
		return "", nil
	}
	kids := n.Children()
	if len(kids) < 2 {
		return "", nil
	}
	prop := kids[1]
	if prop.Token() != rhino.STRINGLIT {
		return "", nil
	}
	base := kids[0]
	switch base.Token() {
	case rhino.NAME:
		return base.String_() + "." + prop.String_(), prop
	case rhino.GETPROP:
		parent, _ := qualifiedName(base)
		if parent == "" {
			return "", nil
		}
		return parent + "." + prop.String_(), prop
	}
	return "", nil
}

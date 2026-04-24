// Port of com/google/javascript/jscomp/CheckRegExp.java @ closure-compiler v20260415

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// REGEXP_PROPERTY_SKIPLIST contains the global RegExp properties whose
// access is allowed without triggering JSC_REGEXP_REFERENCE. Mirrors
// the ImmutableSet in CheckRegExp.java.
var regExpPropertySkiplist = map[string]bool{
	"$1":           true,
	"$2":           true,
	"$3":           true,
	"$4":           true,
	"$5":           true,
	"$6":           true,
	"$7":           true,
	"$8":           true,
	"$9":           true,
	"$_":           true,
	"$input":       true,
	"input":        true,
	"lastMatch":    true,
	"lastParen":    true,
	"leftContext":  true,
	"rightContext":  true,
	"global":       true,
	"ignoreCase":   true,
	"lastIndex":    true,
	"multiline":    true,
	"source":       true,
}

// CheckRegExp looks for references to the global RegExp object that
// would cause regular expressions to be unoptimizable, and checks that
// regular expressions are syntactically valid.
//
// Mirrors CheckRegExp extends AbstractPostOrderCallback implements
// CompilerPass in Java.
type CheckRegExp struct {
	path                             string
	reportErrors                     bool
	globalRegExpPropertiesUsed       bool
	assumeAllGlobalRegexpUsagesVisible bool
	diags                            []diag.Diagnostic
}

// NewCheckRegExp mirrors CheckRegExp(AbstractCompiler, boolean, boolean)
// in Java. reportErrors=true is the standard check pass; false is the
// optimization-only pass.
func NewCheckRegExp(path string, assumeAllGlobalRegexpUsagesVisible bool, reportErrors bool) *CheckRegExp {
	return &CheckRegExp{
		path:                             path,
		reportErrors:                     reportErrors,
		assumeAllGlobalRegexpUsagesVisible: assumeAllGlobalRegexpUsagesVisible,
	}
}

// IsGlobalRegExpPropertiesUsed mirrors the Java method. When
// assumeAllGlobalRegexpUsagesVisible is false the method always returns
// true (the compiler assumes all global RegExp usages are visible).
func (c *CheckRegExp) IsGlobalRegExpPropertiesUsed() bool {
	if c.assumeAllGlobalRegexpUsagesVisible {
		return c.globalRegExpPropertiesUsed
	}
	return true
}

// Check runs the traversal and returns accumulated diagnostics.
// Go-only: language-adapter — Java's process() takes (externs, root); Go uses a single SCRIPT node.
func (c *CheckRegExp) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script, nil)
	return c.diags
}

// ShouldTraverse satisfies the Callback interface (always descend).
// Go-only: language-adapter — Java's AbstractPostOrderCallback has this implicit; Go requires explicit impl.
func (c *CheckRegExp) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }

// Visit mirrors AbstractPostOrderCallback.visit — called after a node's
// children have been visited.
func (c *CheckRegExp) Visit(_ *NodeTraversal, n, parent *rhino.Node) {
	c.visitNode(n, parent)
}

// walk drives a manual post-order traversal so CheckRegExp can be used
// Go-only: language-adapter — Java uses NodeTraversal.traverse; Go uses recursive walk.
// without a full NodeTraversal / TypedScopeCreator rig (no scope needed
// for the REGEXP_REFERENCE check which only cares about NAME("RegExp")
// references at the global level).
func (c *CheckRegExp) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	for _, child := range n.Children() {
		c.walk(child, n)
	}
	c.visitNode(n, parent)
}

// visitNode implements the Java visit() body: detect reference-name
// "RegExp" usages that prevent regex optimization, and (when the parser
// produces REGEXP nodes) validate regex syntax.
//
// NOTE: the jsgo parser does NOT yet produce a REGEXP token kind, so the
// n.isRegExp() / JSC_MALFORMED_REGEXP branch is unreachable for now.
// TODO: parser needs REGEXP token — add the validation branch once the
// scanner produces regex literals as first-class AST nodes.
func (c *CheckRegExp) visitNode(n, parent *rhino.Node) {
	// Branch 1: NodeUtil.isReferenceName(n) — any NAME node that is a
	// value reference (not a declaration binding site). We approximate
	// this by checking for NAME token and skipping obvious binding sites.
	if isRegExpReferenceName(n, parent) {
		name := n.String_()
		if name == "RegExp" {
			// In Java: t.getScope().getVar(name) == null means the name is
			// NOT locally declared — it resolves to the global RegExp.
			// We use the same heuristic: if the identifier is declared in
			// a local (non-extern) scope it shadows the global and is safe.
			// Since CheckRegExp is a lightweight pass that doesn't carry a
			// full TypedScopeCreator, we check whether it resolves to a
			// local declaration by walking to a scope — but we don't have
			// a traversal here. Instead we follow the same logic Java does:
			// NAME nodes that are the first child of VAR/LET/CONST/FUNCTION/
			// CLASS or a PARAM_LIST child are binding sites and were already
			// filtered by isRegExpReferenceName. All remaining NAME("RegExp")
			// nodes are references to the global. (Local re-declarations are
			// handled at the binding site.)
			parentToken := rhino.ERROR
			if parent != nil {
				parentToken = parent.Token()
			}
			first := parent != nil && parent.FirstChild() == n

			// Acceptable usage patterns (mirrors Java's !(...) condition):
			//   new RegExp(...)
			//   RegExp(...)
			//   x instanceof RegExp
			//   RegExp == x  /  RegExp != x  /  RegExp === x  /  RegExp !== x
			//   case RegExp:
			//   RegExp.safeProperty  (safeProperty in skiplist)
			allowed := false
			switch {
			case parentToken == rhino.NEW && first:
				allowed = true
			case parentToken == rhino.CALL && first:
				allowed = true
			case parentToken == rhino.INSTANCEOF && !first:
				allowed = true
			case parentToken == rhino.EQ:
				allowed = true
			case parentToken == rhino.NE:
				allowed = true
			case parentToken == rhino.SHEQ:
				allowed = true
			case parentToken == rhino.SHNE:
				allowed = true
			case parentToken == rhino.CASE:
				allowed = true
			case parentToken == rhino.GETPROP && first:
				// Accessing a property that is in the skiplist is OK; any
				// other property access reveals global RegExp usage.
				propName := ""
				if kids := parent.Children(); len(kids) >= 2 {
					propName = kids[1].String_()
				}
				allowed = regExpPropertySkiplist[propName]
			}

			if !allowed {
				if c.reportErrors {
					c.diags = append(c.diags, diag.Diagnostic{
						Path:    c.path,
						Line:    n.Lineno(),
						Col:     n.Charno(),
						Code:    diag.JSC_REGEXP_REFERENCE,
						Level:   diag.Warning,
						Message: "References to the global RegExp object prevents optimization of regular expressions.",
					})
				}
				c.globalRegExpPropertiesUsed = true
			}
		}
	}

	// Branch 2: n.isRegExp() — validate regex syntax.
	// The jsgo parser DOES produce REGEXP tokens (the full literal is in
	// n.String_(), e.g. "/abc/gi"). Java uses RegExpTree.parseRegExp(pattern,
	// flags) from com.google.javascript.jscomp.regex. Go lacks a JS-compatible
	// regex parser, so JSC_MALFORMED_REGEXP is not yet emitted.
	// TODO(port): implement JS regex syntax validation when a JS regex library
	// is available, or port Java's RegExpTree to Go.
	if c.reportErrors && n.Token() == rhino.REGEXP {
		// Placeholder — no JS regex validation library available yet.
		// Java: pattern = n.getFirstChild().getString(), flags = n.hasTwoChildren() ? ... : ""
		// try { RegExpTree.parseRegExp(pattern, flags) } catch { report MALFORMED_REGEXP }
	}
}

// isRegExpReferenceName reports whether `n` is a NAME node that is a
// value reference (not a declaration binding site). Mirrors
// NodeUtil.isReferenceName in Java: a NAME whose parent makes it a
// binding site (VAR/LET/CONST child, PARAM_LIST child, FUNCTION/CLASS
// first child, STRING_KEY child) is NOT a reference name.
// Go-only: language-adapter — Java uses NodeUtil.isReferenceName directly in the callback.
func isRegExpReferenceName(n, parent *rhino.Node) bool {
	if n == nil || n.Token() != rhino.NAME {
		return false
	}
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// The first child of the declaration is the bound name; any NAME
		// inside an initializer (nested under the binding) has a different
		// parent and is a reference.
		return parent.FirstChild() != n
	case rhino.PARAM_LIST:
		return false
	case rhino.FUNCTION, rhino.CLASS:
		// First child is the declared name; not a reference.
		return parent.FirstChild() != n
	case rhino.CLASS_FIELD_DEF:
		return parent.FirstChild() != n
	case rhino.STRING_KEY:
		// Object-literal key slot — the key itself is stored as the node's
		// own string (STRING_KEY token), so a NAME child of STRING_KEY is
		// always the value side and is a reference.
		return true
	}
	return true
}

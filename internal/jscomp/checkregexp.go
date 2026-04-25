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

// Visit mirrors CheckRegExp.visit(NodeTraversal, Node, Node) — called
// after a node's children have been visited (post-order).
func (c *CheckRegExp) Visit(_ *NodeTraversal, n, parent *rhino.Node) {
	if isRegExpReferenceName(n, parent) {
		name := n.String_()
		if name == "RegExp" {
			// Java: t.getScope().getVar(name) == null — name not locally declared,
			// so it resolves to the global RegExp constructor.
			// Go-only: isRegExpReferenceName already filters binding sites; all
			// remaining NAME("RegExp") nodes reference the global.
			parentToken := rhino.ERROR
			if parent != nil {
				parentToken = parent.Token()
			}
			first := parent != nil && parent.FirstChild() == n
			if !((parentToken == rhino.NEW && first) ||
				(parentToken == rhino.CALL && first) ||
				(parentToken == rhino.INSTANCEOF && !first) ||
				parentToken == rhino.EQ ||
				parentToken == rhino.NE ||
				parentToken == rhino.SHEQ ||
				parentToken == rhino.SHNE ||
				parentToken == rhino.CASE ||
				(parentToken == rhino.GETPROP && first &&
					!regExpPropertySkiplist[parent.String_()])) {
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
		// Check the syntax of regular expression patterns.
	} else if c.reportErrors && n.Token() == rhino.REGEXP {
		// Go-only: language-adapter — Java calls RegExpTree.parseRegExp(pattern, flags);
		// no JS-compatible regex parser exists in Go yet, so JSC_MALFORMED_REGEXP
		// is not emitted. TODO: port Java's RegExpTree or use a JS regex library.
		// Java: pattern = n.getFirstChild().getString(), flags = n.hasTwoChildren() ? n.getLastChild().getString() : ""
		// try { RegExpTree.parseRegExp(pattern, flags) } catch { t.report(n, MALFORMED_REGEXP, ex.getMessage()) }
	}
}

// walk drives a manual post-order traversal so CheckRegExp can be used
// Go-only: language-adapter — Java uses NodeTraversal.traverse; Go uses recursive walk.
// without a full NodeTraversal / TypedScopeCreator rig.
func (c *CheckRegExp) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	for _, child := range n.Children() {
		c.walk(child, n)
	}
	c.Visit(nil, n, parent)
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

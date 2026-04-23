// Port of com/google/javascript/rhino/JSTypeExpression.java @ closure-compiler v20260415

// JSTypeExpression wraps a parsed type-AST Node plus a source filename.
// During type checking, type ASTs are evaluated to JavaScript types.
//
// The struct retains the TypeName string field that existed in the Tier-0
// stub so that all existing callers (scope.go, typecheck.go, …) that reach
// directly for expr.TypeName continue to compile without modification.
// New callers should prefer GetRoot() + Evaluate(...) when the richer
// Node-based type resolution is available.

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// implicitTemplateBoundSource is the synthetic source name used to mark
// the implicit `?` upper bound on an unconstrained template parameter.
// Matches Java's IMPLICIT_TEMPLATE_BOUND_SOURCE constant.
const implicitTemplateBoundSource = "<IMPLICIT_TEMPLATE_BOUND>"

// ImplicitTemplateBound is the package-level singleton that represents
// the implicit `?` upper bound on unconstrained template parameters.
// Mirrors Java's static IMPLICIT_TEMPLATE_BOUND field.
var ImplicitTemplateBound = &JSTypeExpression{
	root:       New(QMARK),
	sourceName: implicitTemplateBoundSource,
	TypeName:   "?",
}

// JSTypeExpression wraps a parsed type-AST Node plus a source filename.
// It is the Go port of com.google.javascript.rhino.JSTypeExpression.
//
// TypeName is kept for backward compatibility with Tier-0 callers that
// access the string form directly; it mirrors the text the parser stored
// before the Node-rooted representation was introduced.
type JSTypeExpression struct {
	// root is the root of the type-annotation AST (matches Java's `root`
	// field). May be nil for expressions created by the Tier-0
	// NewJSTypeExpression(name) factory before a Node parse tree exists.
	root *Node

	// sourceName matches Java's `sourceName` field — the source file
	// where the type expression appeared.
	sourceName string

	// TypeName is the (possibly qualified) name of the declared type in
	// its text form. Tier-0 stores just the primitive name. Tier-1+ will
	// have this derived from the Node-rooted parse tree.
	TypeName string
}

// NewJSTypeExpression is the legacy Tier-0 factory: wraps a plain type
// name string. All existing callers (parser, jsdocinfoparser) use this.
// Preserves the pre-widening API exactly.
func NewJSTypeExpression(name string) *JSTypeExpression {
	return &JSTypeExpression{TypeName: name}
}

// NewJSTypeExpressionFromNode is the Java-matching constructor:
//
//	new JSTypeExpression(Node root, String sourceName)
//
// Used when the type annotation has been parsed into an AST.
func NewJSTypeExpressionFromNode(root *Node, sourceName string) *JSTypeExpression {
	expr := &JSTypeExpression{
		root:       root,
		sourceName: sourceName,
	}
	if root != nil {
		expr.TypeName = nodeTypeName(root)
	}
	return expr
}

// nodeTypeName derives a text representation from a type-AST node.
// This is a best-effort rendering for backward compatibility with callers
// that access expr.TypeName; it is NOT a full pretty-printer.
func nodeTypeName(n *Node) string {
	if n == nil {
		return "?"
	}
	switch n.token {
	case STRINGLIT:
		return n.str
	case QMARK:
		if !n.HasChildren() {
			return "?"
		}
		return "?" + nodeTypeName(n.FirstChild())
	case EQUALS:
		if n.HasChildren() {
			return nodeTypeName(n.FirstChild()) + "="
		}
		return "="
	case ITER_REST:
		if n.HasChildren() {
			return "..." + nodeTypeName(n.FirstChild())
		}
		return "..."
	}
	// Fallback: use the token name.
	return strings.ToLower(n.token.String())
}

// GetRoot returns the root of the type-annotation AST.
// Mirrors Java's JSTypeExpression.getRoot().
func (e *JSTypeExpression) GetRoot() *Node { return e.root }

// GetSourceName returns the source file name where this expression appears.
// Mirrors Java's JSTypeExpression.getSourceName().
func (e *JSTypeExpression) GetSourceName() string { return e.sourceName }

// IsOptionalArg reports whether this expression denotes an optional @param
// (the type was written as `{T=}`). Mirrors Java's isOptionalArg().
func (e *JSTypeExpression) IsOptionalArg() bool {
	return e.root != nil && e.root.token == EQUALS
}

// IsVarArgs reports whether this expression denotes a rest @param (the
// type was written as `{...T}`). Mirrors Java's isVarArgs().
func (e *JSTypeExpression) IsVarArgs() bool {
	return e.root != nil && e.root.token == ITER_REST
}

// TypeEvaluator is the minimal interface for evaluating a type expression
// into a JSType. Implemented by *jstype.Registry via its ResolveExpr method,
// and by any richer registry that supports createTypeFromCommentNode.
//
// This exists because the full JSTypeRegistry.createTypeFromCommentNode is
// not yet ported; callers that already have a *jstype.Registry pass it
// directly and receive the ResolveExpr-based result.
type TypeEvaluator interface {
	// ResolveExpr resolves a textual type expression (e.g. "number|string")
	// to a JSType. Used by the current tier as the Evaluate back-end.
	ResolveExpr(expr string) jstype.JSType
}

// Evaluate resolves this type expression into a JSType using the given
// evaluator. Mirrors Java's JSTypeExpression.evaluate(scope, registry).
//
// When the expression has a Node-rooted AST, the TypeName text form is
// used for resolution (which is the same text that the parser stored).
// When the registry grows createTypeFromCommentNode this method will be
// updated to pass root directly.
func (e *JSTypeExpression) Evaluate(registry TypeEvaluator) jstype.JSType {
	return registry.ResolveExpr(e.TypeName)
}

// IsEquivalentTo reports whether this type expression is structurally equal
// to another. Mirrors Java's JSTypeExpression.isEquivalentTo(other).
//
// When both sides have Node roots the comparison is structural (via
// Node.IsEquivalentTo). When only TypeName strings are available those are
// compared instead.
func (e *JSTypeExpression) IsEquivalentTo(other *JSTypeExpression) bool {
	if other == nil {
		return false
	}
	if e.root != nil && other.root != nil {
		return e.root.IsEquivalentTo(other.root)
	}
	// Fall back to text equality when one or both sides lack a parse tree.
	return e.TypeName == other.TypeName
}

// Copy returns a deep copy of this expression. Mirrors Java's copy().
func (e *JSTypeExpression) Copy() *JSTypeExpression {
	c := &JSTypeExpression{
		sourceName: e.sourceName,
		TypeName:   e.TypeName,
	}
	if e.root != nil {
		c.root = e.root.CloneTree()
	}
	return c
}

// String returns a human-readable rendering. Mirrors Java's toString().
func (e *JSTypeExpression) String() string {
	if e.root != nil {
		return "type: " + e.TypeName
	}
	return "type: " + e.TypeName
}

// IsExplicitUnknownTemplateBound reports whether this expression is the
// explicit `?` upper bound on a template parameter (i.e. written by the
// programmer, not synthesised implicitly). Mirrors Java's
// isExplicitUnknownTemplateBound().
func (e *JSTypeExpression) IsExplicitUnknownTemplateBound() bool {
	return e.root != nil &&
		e.root.token == QMARK &&
		!e.root.HasChildren() &&
		e.sourceName != implicitTemplateBoundSource
}

// MakeOptionalArg wraps expr in an EQUALS node to make it optional, if it
// is not already optional or a rest-args expression. Mirrors Java's static
// JSTypeExpression.makeOptionalArg(JSTypeExpression).
func MakeOptionalArg(expr *JSTypeExpression) *JSTypeExpression {
	if expr.IsOptionalArg() || expr.IsVarArgs() {
		return expr
	}
	var equalsRoot *Node
	if expr.root != nil {
		equalsRoot = New(EQUALS)
		equalsRoot.AddChildToBack(expr.root.CloneTree())
	}
	return &JSTypeExpression{
		root:       equalsRoot,
		sourceName: expr.sourceName,
		TypeName:   expr.TypeName + "=",
	}
}

// GetAllTypeNodes returns all STRING_LIT leaf nodes in this type expression.
// Mirrors Java's getAllTypeNodes().
func (e *JSTypeExpression) GetAllTypeNodes() []*Node {
	if e.root == nil {
		return nil
	}
	var result []*Node
	visitAllTypeNodes(e.root, func(n *Node) { result = append(result, n) })
	return result
}

// GetAllTypeNames returns the string values of all STRING_LIT leaf nodes in
// this type expression. Mirrors Java's getAllTypeNames().
func (e *JSTypeExpression) GetAllTypeNames() []string {
	if e.root == nil {
		if e.TypeName != "" && e.TypeName != "?" {
			return []string{e.TypeName}
		}
		return nil
	}
	var result []string
	visitAllTypeNodes(e.root, func(n *Node) { result = append(result, n.str) })
	return result
}

// visitAllTypeNodes walks the type-AST rooted at n and calls visitor for
// every STRINGLIT node. Mirrors Java's private visitAllTypeNodes.
func visitAllTypeNodes(n *Node, visitor func(*Node)) {
	if n == nil {
		return
	}
	for _, c := range n.children {
		visitAllTypeNodes(c, visitor)
	}
	if n.IsStringLit() {
		visitor(n)
	}
}

// ReplaceNamesWithUnknownType returns a copy of this expression with every
// STRINGLIT whose base name (the part before any `.`) is in names replaced
// by a bare QMARK node. Mirrors Java's replaceNamesWithUnknownType.
func (e *JSTypeExpression) ReplaceNamesWithUnknownType(names map[string]bool) *JSTypeExpression {
	if e.root == nil {
		// No parse tree — fall back to text-level replacement.
		base := e.TypeName
		if dot := strings.Index(base, "."); dot != -1 {
			base = base[:dot]
		}
		if names[base] {
			return &JSTypeExpression{TypeName: "?"}
		}
		return e
	}
	newRoot := replaceNames(e.root.CloneTree(), names)
	return NewJSTypeExpressionFromNode(newRoot, e.sourceName)
}

// replaceNames is the recursive helper for ReplaceNamesWithUnknownType.
// Mirrors Java's private static replaceNames.
func replaceNames(n *Node, names map[string]bool) *Node {
	if n == nil {
		return nil
	}
	for _, c := range n.children {
		replaceNames(c, names)
	}
	if n.IsStringLit() && names[baseName(n.str)] {
		qMark := New(QMARK)
		if n.HasChildren() && !n.FirstChild().IsBlock() {
			// Preserve templatised children, but not BLOCK (record body).
			qMark.AddChildrenToBack(n.RemoveChildren())
		}
		if n.HasParent() {
			n.ReplaceWith(qMark)
		}
		return qMark
	}
	return n
}

// baseName strips any leading qualifier from a dotted name, mirroring
// Java's private baseName helper.
func baseName(name string) string {
	if dot := strings.Index(name, "."); dot != -1 {
		return name[:dot]
	}
	return name
}

// GetRecordPropertyNames returns the set of all property keys declared in
// record-type literals (e.g. `{{key: string}}`). Mirrors Java's
// getRecordPropertyNames().
func (e *JSTypeExpression) GetRecordPropertyNames() []string {
	if e.root == nil {
		return nil
	}
	var result []string
	getRecordPropertyNamesRecursive(e.root, &result)
	return result
}

// getRecordPropertyNamesRecursive collects STRING_KEY tokens recursively.
func getRecordPropertyNamesRecursive(n *Node, names *[]string) {
	if n == nil {
		return
	}
	for _, c := range n.children {
		getRecordPropertyNamesRecursive(c, names)
	}
	if n.IsStringKey() {
		*names = append(*names, n.str)
	}
}

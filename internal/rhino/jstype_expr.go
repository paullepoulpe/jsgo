// Port of com/google/javascript/rhino/JSTypeExpression.java +
// com/google/javascript/rhino/jstype/NamedType.java
// @ closure-compiler v20260415
//
// In Java, JSTypeExpression wraps a parsed type-AST Node plus a source
// filename, and is later resolved against a JSTypeRegistry. NamedType is
// the forward-reference wrapper the registry hands out for a name that
// hasn't yet been declared at resolution time. Tier-0 stored only a flat
// TypeName string; to emit UNRECOGNIZED_TYPE_ERROR at the JSDoc location
// (as Java's late-resolve diagnostic does) we add comment-relative line
// and column offsets captured by JsDocInfoParser, plus a list of
// NamedType wrappers gathered per type token so the late resolver can
// revisit them.
//
// The offsets are relative to the containing `/** ... */` comment's
// start; jstype_resolver converts them to absolute source positions by
// scanning backward from the annotated AST node for the enclosing
// JSDoc block.

package rhino

type JSTypeExpression struct {
	// TypeName is the (possibly qualified) name of the declared type.
	// Tier-0 stores just the primitive name. Tier-1+ will replace this
	// with a Node-rooted parse tree, matching the Java structure.
	TypeName string

	// SourceLine / SourceCol are the position of TypeName *inside* the
	// containing JSDoc comment body (0-indexed line, 0-indexed column
	// from the start of the comment body, matching the semantics of
	// JsDocTokenStream.getLineno / getCharno). Callers add the comment
	// start position to recover absolute coordinates.
	SourceLine int
	SourceCol  int

	// Names lists every name token seen while parsing this expression,
	// each tagged with the (line, col) at which it appeared inside the
	// comment body. The resolver consults this after scope creation to
	// emit JSC_UNRECOGNIZED_TYPE_ERROR at the exact JSDoc position, not
	// at the owning AST node's position.
	Names []NamedType

	// ParseError is set by JsDocInfoParser when a type expression is
	// malformed (missing `}`, unmatched paren, etc). The resolver turns
	// this into JSC_TYPE_PARSE_ERROR at the captured position.
	ParseError *NamedType
}

// NamedType mirrors com/google/javascript/rhino/jstype/NamedType.java:
// a wrapper around an unresolved reference name plus the position at
// which it appeared. The resolver walks these after all scope bindings
// exist; still-unresolved names turn into diagnostics.
type NamedType struct {
	UnresolvedName string
	SourceLine     int
	SourceCol      int
}

// Name returns the wrapped name. Mirrors NamedType.getReferenceName().
func (n NamedType) Name() string { return n.UnresolvedName }

func NewJSTypeExpression(name string) *JSTypeExpression {
	return &JSTypeExpression{TypeName: name}
}

// NewJSTypeExpressionAt mirrors JSTypeExpression's richer constructor:
// carries the comment-relative origin of the expression alongside its
// serialized form.
func NewJSTypeExpressionAt(name string, line, col int, names []NamedType) *JSTypeExpression {
	return &JSTypeExpression{
		TypeName:   name,
		SourceLine: line,
		SourceCol:  col,
		Names:      names,
	}
}

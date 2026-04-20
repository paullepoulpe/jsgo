// Port of com/google/javascript/rhino/JSTypeExpression.java @ closure-compiler v20260415
//
// In Java, JSTypeExpression wraps a parsed type-AST Node plus a source
// filename, and is later resolved against a JSTypeRegistry. For Tier-0 we
// only need to carry a simple type name ("number", "string"); the parse
// tree of richer expressions (unions, templates) is added later.

package rhino

type JSTypeExpression struct {
	// TypeName is the (possibly qualified) name of the declared type.
	// Tier-0 stores just the primitive name. Tier-1+ will replace this
	// with a Node-rooted parse tree, matching the Java structure.
	TypeName string

	// Optional mirrors the trailing `=` on a parameter type
	// (`@param {T=} name`) — marks the param as legally omittable
	// at call sites.
	Optional bool

	// Rest mirrors the leading `...` on a parameter type
	// (`@param {...T} name`) — marks the param as a
	// variadic rest-slot at call sites.
	Rest bool
}

func NewJSTypeExpression(name string) *JSTypeExpression {
	return &JSTypeExpression{TypeName: name}
}

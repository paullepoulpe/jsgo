// Forward declarations of com.google.javascript.rhino.jstype types that the
// core rhino files reference. rhino and rhino.jstype are one Go package
// (PORTING.md §2); until jstype is ported, each type it provides is an
// interface declaring only the methods the ported code calls. This file is
// deleted when jstype is ported.

package rhino

// JSType stands in for com.google.javascript.rhino.jstype.JSType.
//
// Node stores a JSType or a colors.Color in one field and tells them apart
// with `instanceof JSType`; the calls it makes are toString() (when
// printing) and equals() (Objects.equals in isEquivalentTo).
type JSType interface {
	// Java: JSType.toString()
	String() string

	// Java: JSType.equals(Object)
	Equals(o any) bool

	// Go-only: language-adapter — an unexported marker so that values of
	// other types with String/Equals methods (e.g. *colors.Color) do not
	// satisfy `instanceof JSType` by Go's structural typing. The real JSType
	// interface excludes them through its JSTypeBase accessor instead.
	isForwardDeclaredJSType()
}

// JSTypeRegistry stands in for
// com.google.javascript.rhino.jstype.JSTypeRegistry.
type JSTypeRegistry interface {
	// Java: JSTypeRegistry.createTypeFromCommentNode(Node, String, StaticTypedScope)
	CreateTypeFromCommentNode(n *Node, sourceName string, scope StaticTypedScope) JSType
}

// StaticTypedScope stands in for
// com.google.javascript.rhino.jstype.StaticTypedScope. The ported code only
// passes it through (JSTypeExpression.evaluate), so no methods are declared.
type StaticTypedScope interface{}

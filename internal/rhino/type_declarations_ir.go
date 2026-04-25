// Port of com/google/javascript/rhino/TypeDeclarationsIR.java @ closure-compiler v20260415
//
// Note: Java's TypeDeclarationsIR resides in the rhino package and calls
// methods on IR (also in rhino). In Go the IR factory functions are the
// IR-prefixed package-level functions in ir.go.

package rhino

// SetDeclaredTypeExpression sets the declared type AST on node n.
// Mirrors Node.setDeclaredTypeExpression(Node).
func (n *Node) SetDeclaredTypeExpression(typeExpression *Node) {
	n.props[PropDeclaredTypeExpr] = typeExpression
}

// GetDeclaredTypeExpression returns the declared type AST on node n, or nil.
// Mirrors Node.getDeclaredTypeExpression().
func (n *Node) GetDeclaredTypeExpression() *Node {
	if n.props == nil {
		return nil
	}
	v, _ := n.props[PropDeclaredTypeExpr].(*Node)
	return v
}

// IRStringType returns a new node representing the string built-in type.
// Mirrors TypeDeclarationsIR.stringType().
func IRStringType() *Node { return New(STRING_TYPE) }

// IRNumberType returns a new node representing the number built-in type.
// Mirrors TypeDeclarationsIR.numberType().
func IRNumberType() *Node { return New(NUMBER_TYPE) }

// IRBooleanType returns a new node representing the boolean built-in type.
// Mirrors TypeDeclarationsIR.booleanType().
func IRBooleanType() *Node { return New(BOOLEAN_TYPE) }

// IRAnyType returns a new node representing any type (the UNKNOWN type in
// Closure, expressed with {?}). Mirrors TypeDeclarationsIR.anyType().
func IRAnyType() *Node { return New(ANY_TYPE) }

// IRVoidType returns a new node representing the Void type.
// Mirrors TypeDeclarationsIR.voidType().
func IRVoidType() *Node { return New(VOID_TYPE) }

// IRUndefinedType returns a new node representing the Undefined type.
// Mirrors TypeDeclarationsIR.undefinedType().
func IRUndefinedType() *Node { return New(UNDEFINED_TYPE) }

// IRNamedType splits a '.' separated qualified name into a tree of type
// segments. Mirrors TypeDeclarationsIR.namedType(String).
func IRNamedType(typeName string) *Node {
	segments := splitDot(typeName)
	return IRNamedTypeFromSegments(segments)
}

// IRNamedTypeFromSegments produces a NAMED_TYPE node tree from the given
// segments. Mirrors TypeDeclarationsIR.namedType(Iterable<String>).
func IRNamedTypeFromSegments(segments []string) *Node {
	if len(segments) == 0 {
		return New(NAMED_TYPE)
	}
	node := IRName(segments[0])
	for _, seg := range segments[1:] {
		node = IRGetprop(node, seg)
	}
	return newNodes(NAMED_TYPE, node)
}

// IRRecordType constructs a RECORD_TYPE node from a map of property names to
// property type nodes. Mirrors TypeDeclarationsIR.recordType(LinkedHashMap).
func IRRecordType(properties map[string]*Node) *Node {
	node := New(RECORD_TYPE)
	for name, typ := range properties {
		stringKey := IRStringKey(name)
		node.AddChildToBack(stringKey)
		if typ != nil {
			stringKey.AddChildToFront(typ)
		}
	}
	return node
}

// irMaybeAddType is a helper that sets the declared type expression on node if
// typ is non-nil. Mirrors TypeDeclarationsIR.maybeAddType(Node, Node).
func irMaybeAddType(node, typ *Node) *Node {
	if typ != nil {
		node.SetDeclaredTypeExpression(typ)
	}
	return node
}

// IRFunctionType constructs a FUNCTION_TYPE node with the given return type
// and parameters. Mirrors TypeDeclarationsIR.functionType(...).
//
// requiredParams and optionalParams are ordered pairs of (name, type);
// pass nil for the type to leave it untyped.
func IRFunctionType(
	returnType *Node,
	requiredParams [][2]*Node,
	optionalParams [][2]*Node,
	restName string,
	restType *Node,
) *Node {
	node := newNodes(FUNCTION_TYPE, returnType)
	for _, p := range requiredParams {
		name := IRName(p[0].GetString())
		node.AddChildToBack(irMaybeAddType(name, p[1]))
	}
	for _, p := range optionalParams {
		name := IRName(p[0].GetString())
		name.PutBooleanProp(PropOptEs6Typed, true)
		node.AddChildToBack(irMaybeAddType(name, p[1]))
	}
	if restName != "" {
		rest := newNodes(ITER_REST, IRName(restName))
		node.AddChildToBack(irMaybeAddType(rest, restType))
	}
	return node
}

// IRParameterizedType constructs a PARAMETERIZED_TYPE node.
// Mirrors TypeDeclarationsIR.parameterizedType(Node, Iterable<Node>).
func IRParameterizedType(baseType *Node, typeParameters []*Node) *Node {
	if len(typeParameters) == 0 {
		return baseType
	}
	node := newNodes(PARAMETERIZED_TYPE, baseType)
	for _, tp := range typeParameters {
		node.AddChildToBack(tp)
	}
	return node
}

// IRArrayType constructs an ARRAY_TYPE node.
// Mirrors TypeDeclarationsIR.arrayType(Node).
func IRArrayType(elementType *Node) *Node {
	return newNodes(ARRAY_TYPE, elementType)
}

// IRUnionType constructs a UNION_TYPE node from the given options.
// Mirrors TypeDeclarationsIR.unionType(Iterable<Node>).
func IRUnionType(options ...*Node) *Node {
	node := New(UNION_TYPE)
	for _, o := range options {
		node.AddChildToBack(o)
	}
	return node
}

// IROptionalParameter wraps parameterType in an OPTIONAL_PARAMETER node.
// Mirrors TypeDeclarationsIR.optionalParameter(Node).
func IROptionalParameter(parameterType *Node) *Node {
	return newNodes(OPTIONAL_PARAMETER, parameterType)
}

// splitDot splits s on '.'.
func splitDot(s string) []string {
	if s == "" {
		return nil
	}
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

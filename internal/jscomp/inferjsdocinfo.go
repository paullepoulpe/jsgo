// Port of com/google/javascript/jscomp/InferJsDocInfo.java
// @ closure-compiler v20260415 (subset).
//
// Sets the JSDocInfo on all JSType objects — including their properties —
// using the JSDoc on the node defining that type or property.
//
// This pass propagates JSDocs across the type graph, but not across the
// symbol graph. For example:
//
//	/**
//	 * I'm a user type!
//	 * @constructor
//	 */
//	var Foo = function() { };
//	var Bar = Foo;
//
// assigns the "I'm a user type!" JSDoc from the `Foo` node to type "Foo"
// and type "ctor{Foo}". This pass says nothing about JSDocs being
// propagated to the `Bar` node.
//
// JSDoc is initially attached to AST Nodes at parse time. There are 3
// cases where JSDocs get attached to the type system, always from the
// associated declaration node:
//
//  1. Nominal types (constructors, interfaces, enums).
//  2. Object type properties (e.g. `Foo.prototype.bar`).
//  3. Anonymous structural types, including functions.
//
// Zero diagnostics are emitted; the pass is a pure annotation step.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// InferJsDocInfo mirrors Java's InferJSDocInfo class — it implements
// AbstractPostOrderCallback and CompilerPass. It walks the AST post-order
// and annotates each JSType with the JSDoc from its declaring node.
type InferJsDocInfo struct {
	creator *TypedScopeCreator
	global  *TypedScope
}

// NewInferJsDocInfo mirrors InferJSDocInfo(AbstractCompiler). The scope
// creator + global scope replace the compiler reference because jsgo
// resolves types from the scope rather than from node annotations.
func NewInferJsDocInfo(creator *TypedScopeCreator, global *TypedScope) *InferJsDocInfo {
	return &InferJsDocInfo{creator: creator, global: global}
}

// Process mirrors InferJSDocInfo.process(Node, Node). Traverses externs
// then root in post-order, annotating type objects with their JSDoc.
func (p *InferJsDocInfo) Process(externs, root *rhino.Node) {
	trav := NewNodeTraversal(p, p.creator, p.global)
	if externs != nil {
		trav.Traverse(externs)
	}
	if root != nil {
		trav.Traverse(root)
	}
}

// ShouldTraverse always returns true — this is a post-order-only pass.
func (p *InferJsDocInfo) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }

// Visit mirrors InferJSDocInfo.visit(NodeTraversal, Node, Node). Dispatches
// to the appropriate infer* helper based on the node's token.
//
// The Go AST is a subset of the Java AST — Java tokens GETTER_DEF,
// SETTER_DEF, MEMBER_FUNCTION_DEF, COMPUTED_PROP, and CLASS_MEMBERS are
// not present in our token set. Instead, class methods appear as FUNCTION
// children of CLASS, and class fields as CLASS_FIELD_DEF children.
// Object literal properties use STRING_KEY.
func (p *InferJsDocInfo) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.NAME:
		// Infer JSDocInfo on types of all type declarations on variables.
		p.ijiInferJSDocForName(n, parent, t.CurrentScope())
	case rhino.STRING_KEY:
		// Object literal key: `{ /** @type {T} */ prop: value }`.
		p.ijiInferJSDocForObjectKeyOrMember(n, parent, t.CurrentScope())
	case rhino.GETPROP:
		p.ijiInferJSDocForGetProp(n, parent, t.CurrentScope())
	}
}

// ijiInferJSDocForName mirrors InferJSDocInfo.inferJSDocForName. Handles
// NAME nodes that appear in var/let/const declarations, function
// declarations, class declarations, and simple assigns.
func (p *InferJsDocInfo) ijiInferJSDocForName(n, parent *rhino.Node, scope *TypedScope) {
	if parent == nil {
		return
	}

	var typeDoc *rhino.JSDocInfo
	var inferredType jstype.JSType
	var aliasedType jstype.JSType // non-nil when RHS is a qualified name

	switch {
	case ijiIsNameDeclaration(parent):
		// Case: `/** ... */ (var|let|const) x = function() { ... }`.
		// Case: `(var|let|const) /** ... */ x = function() { ... }`.
		nameInfo := n.JSDocInfo()
		if nameInfo != nil {
			typeDoc = nameInfo
		} else {
			typeDoc = parent.JSDocInfo()
		}
		inferredType = ijiInferNameType(n, scope)
		// For `const y = x.y;`, the type of `x.y`.
		if value := n.FirstChild(); value != nil && value.IsQualifiedName() {
			aliasedType = ijiInferExprType(value, scope)
		}

	case ijiIsFunctionDeclaration(parent) || ijiIsClassDeclaration(parent):
		// Case: `/** ... */ function f() { ... }`.
		// Case: `/** ... */ class Foo() { ... }`.
		typeDoc = parent.JSDocInfo()
		inferredType = ijiInferNameType(n, scope)
		aliasedType = nil

	case parent.Token() == rhino.ASSIGN && n.IsFirstChildOf(parent):
		typeDoc = parent.JSDocInfo()
		inferredType = ijiInferNameType(n, scope)
		// For `y = x.y;`, the type of the RHS.
		if rhs := n.Next(); rhs != nil && rhs.IsQualifiedName() {
			aliasedType = ijiInferExprType(rhs, scope)
		}

	default:
		return
	}

	if typeDoc == nil {
		return
	}

	objType := ijiDereferenced(inferredType)
	if ijiShouldAttachJSDocToNominalTypeOrShape(objType, aliasedType) {
		ijiAttachJSDocInfoToNominalTypeOrShape(objType, typeDoc, n.String_())
	}
}

// ijiInferJSDocForObjectKeyOrMember mirrors
// InferJSDocInfo.inferJSDocForObjectKeyOrClassField. Handles STRING_KEY
// nodes inside object literals. In the Go AST, class methods appear as
// FUNCTION children of CLASS (not as MEMBER_FUNCTION_DEF), so the class-
// member branch from the Java pass doesn't apply here.
func (p *InferJsDocInfo) ijiInferJSDocForObjectKeyOrMember(n, parent *rhino.Node, scope *TypedScope) {
	typeDoc := n.JSDocInfo()

	// Object literal: STRING_KEY inside an OBJECTLIT. The owning type is
	// the type of the OBJECTLIT parent.
	if parent == nil {
		return
	}

	owningType := ijiDereferenced(ijiInferExprType(parent, scope))
	// When the OBJECTLIT itself isn't typed (e.g. in Foo.prototype = {…}),
	// try resolving the prototype type from the enclosing assignment.
	// Mirrors Java's inferJSDocForObjectKeyOrClassField walking up to
	// find the owningType when the OBJECTLIT is a bulk-prototype override:
	//   Foo.prototype = { bar: function() { … } }
	if owningType == nil {
		owningType = ijiPrototypeLitOwner(parent, scope)
	}

	if owningType == nil {
		return
	}

	propName := n.String_()
	if propName == "" {
		return
	}

	// Determine the JSDoc to attach: prefer the STRING_KEY's own doc;
	// fall back to the enclosing ASSIGN's doc for the prototype-literal
	// case (branch 2: inherit ASSIGN JSDoc to each STRING_KEY value).
	if typeDoc == nil {
		if assign := ijiEnclosingPrototypeAssign(parent); assign != nil {
			typeDoc = assign.JSDocInfo()
		}
	}
	if typeDoc == nil {
		return
	}

	if owningType.HasOwnProperty(propName) && owningType.GetPropertyJSDocInfo(propName) == nil {
		owningType.SetPropertyJSDocInfo(propName, typeDoc)
	}
}

// ijiPrototypeLitOwner returns the ObjectType for the prototype when
// parent is an OBJECTLIT that is the RHS of `Foo.prototype = {…}`.
// Mirrors how Java's inferJSDocForObjectKeyOrClassField resolves the
// owning type for object-literal prototype assignments.
func ijiPrototypeLitOwner(objectlit *rhino.Node, scope *TypedScope) *jstype.ObjectType {
	if objectlit == nil || objectlit.Token() != rhino.OBJECTLIT {
		return nil
	}
	assign := objectlit.Parent()
	if assign == nil || assign.Token() != rhino.ASSIGN {
		return nil
	}
	kids := assign.Children()
	if len(kids) < 2 || kids[1] != objectlit {
		return nil
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return nil
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 || lhsKids[1].String_() != "prototype" {
		return nil
	}
	// Resolve the type of the receiver to get the constructor's prototype.
	recv := ijiInferExprType(lhsKids[0], scope)
	if fn, ok := recv.(*jstype.FunctionType); ok && fn.Instance != nil {
		return fn.Instance
	}
	return nil
}

// ijiEnclosingPrototypeAssign returns the ASSIGN node when objectlit is
// the RHS of a `Foo.prototype = {…}` assignment, nil otherwise.
func ijiEnclosingPrototypeAssign(objectlit *rhino.Node) *rhino.Node {
	if objectlit == nil {
		return nil
	}
	assign := objectlit.Parent()
	if assign == nil || assign.Token() != rhino.ASSIGN {
		return nil
	}
	kids := assign.Children()
	if len(kids) < 2 || kids[1] != objectlit {
		return nil
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return nil
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 || lhsKids[1].String_() != "prototype" {
		return nil
	}
	return assign
}

// ijiInferJSDocForGetProp mirrors InferJSDocInfo.inferJSDocForProperty.
// Handles GETPROP nodes (property reads / declarations).
func (p *InferJsDocInfo) ijiInferJSDocForGetProp(n, parent *rhino.Node, scope *TypedScope) {
	var typeDoc *rhino.JSDocInfo
	var aliasedType jstype.JSType

	switch {
	case parent != nil && parent.Token() == rhino.ASSIGN && n.IsFirstChildOf(parent):
		// Case: `/** @deprecated */ obj.prop = ...;`
		typeDoc = parent.JSDocInfo()
		rhs := n.Next()
		if rhs != nil && rhs.IsQualifiedName() {
			aliasedType = ijiInferExprType(rhs, scope)
		}
	case parent != nil && parent.Token() == rhino.EXPR_RESULT:
		// Case: `/** @deprecated */ obj.prop;`
		typeDoc = n.JSDocInfo()
		aliasedType = nil
	default:
		return
	}

	if typeDoc == nil || !typeDoc.ContainsDeclaration() {
		return
	}

	// Resolve the receiver's type.
	// GETPROP has two children: the receiver (children[0]) and the
	// property-name STRINGLIT (children[1]). Use children[1].String_()
	// rather than n.String_() because the Go AST stores the property
	// name in the child node, not in the GETPROP's own str field
	// (unlike Java's Node.getString() which mirrors the str field).
	nKids := n.Children()
	if len(nKids) < 2 {
		return
	}
	recv := nKids[0]
	propName := nKids[1].String_()
	if propName == "" {
		return
	}
	lhsType := ijiDereferenced(ijiInferExprType(recv, scope))
	if lhsType == nil {
		return
	}

	// Put the JSDoc in the property slot, if there is one.
	if lhsType.HasOwnProperty(propName) && lhsType.GetPropertyJSDocInfo(propName) == nil {
		lhsType.SetPropertyJSDocInfo(propName, typeDoc)
	}

	propTypeRaw := lhsType.GetPropertyType(propName)
	propType := ijiDereferenced(propTypeRaw)
	if ijiShouldAttachJSDocToNominalTypeOrShape(propType, aliasedType) {
		ijiAttachJSDocInfoToNominalTypeOrShape(propType, typeDoc, n.GetQualifiedName())
	}
	// Branch 1: `Foo.prototype.bar = function() { … }`.
	// In Java, FunctionType is a subtype of ObjectType, so
	// attachJSDocInfoToNominalTypeOrShape handles it. In Go they are
	// separate; attach the JSDoc directly to the FunctionType here when
	// the property type is a function and has no doc yet.
	if fn, ok := propTypeRaw.(*jstype.FunctionType); ok && fn.GetJSDocInfo() == nil {
		fn.SetJSDocInfo(typeDoc)
	}
}

// ijiDereferenced mirrors InferJSDocInfo's private dereferenced() — null-
// safe wrapper for ObjectType dereferencing.
func ijiDereferenced(t jstype.JSType) *jstype.ObjectType {
	return jstype.Dereference(t)
}

// ijiShouldAttachJSDocToNominalTypeOrShape mirrors the private
// shouldAttachJSDocToNominalTypeOrShape. Returns true when the type is
// non-nil, has no existing JSDoc, and is not merely an alias for another
// type (aliasedType == nil or != type).
func ijiShouldAttachJSDocToNominalTypeOrShape(t *jstype.ObjectType, aliasedType jstype.JSType) bool {
	if t == nil {
		return false
	}
	if t.GetJSDocInfo() != nil {
		return false
	}
	// Alias guard: don't attach the `/** @const */` on `var cbAlias = cb`
	// to the type of `cb`.
	if aliasedType != nil {
		if aliased, ok := aliasedType.(*jstype.ObjectType); ok && aliased == t {
			return false
		}
	}
	return true
}

// ijiAttachJSDocInfoToNominalTypeOrShape mirrors the private
// attachJSDocInfoToNominalTypeOrShape — handles cases #1 and #3 in the
// class-level doc.
//
// In Java the method distinguishes constructors, interfaces, and enum
// types via FunctionType.isConstructor() / isInterface() and
// EnumType.getElementsType(). In our Go model ObjectType doesn't encode
// those flags directly; we approximate by checking whether the type has a
// reference name (nominal) or not (anonymous structural), skipping native
// externs types that should not be clobbered.
func ijiAttachJSDocInfoToNominalTypeOrShape(objType *jstype.ObjectType, docInfo *rhino.JSDocInfo, qName string) {
	if objType.IsNativeObjectType() {
		return
	}

	if objType.HasReferenceName() {
		// Case #1: nominal types (constructors, interfaces, enums).
		// Only attach when the reference name matches the declaring name,
		// guarding against attaching e.g. the alias's JSDoc to the aliased
		// type. Mirrors isReferenceNameOf in Java.
		if !ijiIsReferenceNameOf(objType, qName) {
			return
		}
		objType.SetJSDocInfo(docInfo)
	} else {
		// Case #3: anonymous structural types (object literals, anonymous
		// function expressions). Every anonymous instance is unique so we
		// can safely attach without a name check.
		objType.SetJSDocInfo(docInfo)
	}
}

// ijiIsReferenceNameOf mirrors InferJSDocInfo.isReferenceNameOf — true
// when `type.getReferenceName() == name`.
func ijiIsReferenceNameOf(t *jstype.ObjectType, name string) bool {
	return t.HasReferenceName() && t.GetReferenceName() == name
}

// ijiIsNameDeclaration mirrors NodeUtil.isNameDeclaration — true for VAR,
// LET, CONST parents.
func ijiIsNameDeclaration(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		return true
	}
	return false
}

// ijiIsFunctionDeclaration mirrors NodeUtil.isFunctionDeclaration — true
// for a FUNCTION node that is a statement (not an expression).
func ijiIsFunctionDeclaration(n *rhino.Node) bool {
	return n.Token() == rhino.FUNCTION
}

// ijiIsClassDeclaration mirrors NodeUtil.isClassDeclaration — true for a
// CLASS node that is a statement.
func ijiIsClassDeclaration(n *rhino.Node) bool {
	return n.Token() == rhino.CLASS
}

// ijiInferNameType resolves the declared type of a NAME node from the
// scope. Mirrors Node.getJSType() on a NAME reference — returns the type
// attached by TypeInference in the Java pipeline.
func ijiInferNameType(n *rhino.Node, scope *TypedScope) jstype.JSType {
	if n == nil || scope == nil {
		return nil
	}
	return scope.GetDeclaredType(n.String_())
}

// ijiInferExprType resolves the type of an expression node. Used to
// determine the "aliasedType" when the RHS of an assignment is a
// qualified name, and to look up the receiver type for GETPROP /
// STRING_KEY / MEMBER_FUNCTION_DEF nodes. Mirrors Node.getJSType() for
// expression nodes in the Java pipeline.
func ijiInferExprType(n *rhino.Node, scope *TypedScope) jstype.JSType {
	if n == nil || scope == nil {
		return nil
	}
	switch n.Token() {
	case rhino.NAME:
		return scope.GetDeclaredType(n.String_())
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) < 2 {
			return nil
		}
		recv := ijiInferExprType(kids[0], scope)
		propName := kids[1].String_()
		// ObjectType: standard property lookup.
		if obj, ok := recv.(*jstype.ObjectType); ok {
			return obj.GetPropertyType(propName)
		}
		// FunctionType: in Java, FunctionType extends ObjectType so
		// getJSType() on `Foo` returns the FunctionType and `.prototype`
		// is special-cased. Mirror that here: "prototype" yields the
		// constructor's Instance (prototype ObjectType); other names
		// look into the Instance's property map.
		if fn, ok := recv.(*jstype.FunctionType); ok && fn.Instance != nil {
			if propName == "prototype" {
				return fn.Instance
			}
			return fn.Instance.GetPropertyType(propName)
		}
	}
	return nil
}

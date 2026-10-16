// Port of com/google/javascript/jscomp/InferJSDocInfo.java @ closure-compiler v20260415

/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InferJSDocInfo sets the JSDocInfo on all JSTypes, including their
// properties, using the JSDoc on the node defining that type or property.
//
// This pass propagates JSDocs across the type graph, but not across the
// symbol graph. For example:
//
//	/**
//	 * I'm a user type!
//	 * @constructor
//	 *\/
//	var Foo = function() { };
//
//	var Bar = Foo;
//
// will assign the "I'm a user type!" JSDoc from the `Foo` node to type
// "Foo" and type "ctor{Foo}". However, this pass says nothing about JSDocs
// being propagated to the `Bar` node.
//
// JSDoc is initially attached to AST Nodes at parse time. There are 3 cases
// where JSDocs get attached to the type system, always from the associated
// declaration node:
//
//   - Nominal types (e.g. constructors, interfaces, enums).
//   - Object type properties (e.g. `Foo.prototype.bar`).
//   - Anonymous structural types, including functions. (Some function types
//     with the same signature are unique and therefore can have distinct
//     JSDocs.) (b/111070482: it's unclear why we support this.)
//
// #1 is fairly straight-forward with the additional detail that JSDocs are
// propagated to both the instance type *and* the declaration type (i.e.
// the ctor or enum type). #2 should also be mostly self-explanatory; it
// covers scenarios like the following:
//
//	/**
//	 * I'm a method!
//	 * @param {number} x
//	 * @return {number}
//	 *\/
//	Foo.prototype.bar = function(x) { ... };
//
// in which JSDocInfo will appear on the "bar" slot of `Foo.prototype` and
// `Foo`. The function type used as the RHS of the assignments (i.e.
// `function(this:Foo, number): number`) is not considered. Note that this
// example would work equally well if `bar` were declared abstract. #3 is a
// bit trickier; it covers types such as the following declarations:
//
//	/** I'm an anonymous structural object type! *\/
//	var myObject = {a: 5, b: 'Hello'};
//
//	/**
//	 * I'm an anonymous structural function type!
//	 * @param {number} x
//	 * @return {string}
//	 *\/
//	var myFunction = function(x) { ... };
//
// which define unique types with their own JSDoc attributes. Object literal
// or function types with the same structure will get different JSDocs
// despite possibly comparing equal. Additionally, when assigning instances
// of these types as properties of nominal types (e.g. using `myFunction` as
// the RHS of #2) the structural type JSDoc plays no part.
//
// Java: class InferJSDocInfo extends AbstractPostOrderCallback implements CompilerPass
type InferJSDocInfo struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: InferJSDocInfo(AbstractCompiler)
func NewInferJSDocInfo(compiler AbstractCompiler) *InferJSDocInfo {
	return &InferJSDocInfo{compiler: compiler}
}

// Java: InferJSDocInfo.process(Node, Node)
func (p *InferJSDocInfo) Process(externs *rhino.Node, root *rhino.Node) {
	if externs != nil {
		NodeTraversal_Traverse(p.compiler, externs, p)
	}
	if root != nil {
		NodeTraversal_Traverse(p.compiler, root, p)
	}
}

// Java: InferJSDocInfo.visit(NodeTraversal, Node, Node)
func (p *InferJSDocInfo) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		// Infer JSDocInfo on types of all type declarations on variables.
		p.inferJSDocForName(n, parent)
		return
	case rhino.Token_STRING_KEY,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_COMPUTED_FIELD_DEF:
		p.inferJSDocForObjectKeyOrClassField(n, parent)
		return
	case rhino.Token_GETPROP:
		p.inferJSDocForProperty(n, parent)
		return
	default:
		return
	}
}

// Java: InferJSDocInfo.inferJSDocForName(Node, Node)
func (p *InferJSDocInfo) inferJSDocForName(n *rhino.Node, parent *rhino.Node) {
	if parent == nil {
		return
	}

	// Only allow JSDoc on variable declarations, named functions, named classes, and assigns.
	var typeDoc *rhino.JSDocInfo
	var aliasedType rhino.JSType // if the right-hand side is a qualified name, its type
	var inferredType rhino.JSType
	if NodeUtil_IsNameDeclaration(parent) {
		// Case: `/** ... */ (var|let|const) x = function() { ... }`.
		// Case: `(var|let|const) /** ... */ x = function() { ... }`.
		nameInfo := n.GetJSDocInfo()
		if nameInfo != nil {
			typeDoc = nameInfo
		} else {
			typeDoc = parent.GetJSDocInfo()
		}

		inferredType = n.GetJSType()
		// Example: for `const y = x.y;`, the type of `x.y`.
		value := n.GetFirstChild()
		if value != nil && value.IsQualifiedName() {
			aliasedType = value.GetJSType()
		}
	} else if NodeUtil_IsFunctionDeclaration(parent) || NodeUtil_IsClassDeclaration(parent) {
		// Case: `/** ... */ function f() { ... }`.
		// Case: `/** ... */ class Foo() { ... }`.
		typeDoc = parent.GetJSDocInfo()
		inferredType = parent.GetJSType()
		aliasedType = nil // the value is definitely a class or function, not a qualified name.
	} else if parent.IsAssign() && n.IsFirstChildOf(parent) {
		typeDoc = parent.GetJSDocInfo()
		inferredType = n.GetJSType()
		// Example: y = x.y;
		value := n.GetNext()
		if value.IsQualifiedName() {
			aliasedType = value.GetJSType()
		}
	} else {
		return
	}

	if typeDoc == nil {
		return
	}

	objType := inferJSDocInfo_dereferenced(inferredType)
	if inferJSDocInfo_shouldAttachJSDocToNominalTypeOrShape(objType, aliasedType) {
		name := n.GetString()
		inferJSDocInfo_attachJSDocInfoToNominalTypeOrShape(objType, typeDoc, &name)
	}
}

// Java: InferJSDocInfo.inferJSDocForObjectKeyOrClassField(Node, Node)
func (p *InferJSDocInfo) inferJSDocForObjectKeyOrClassField(n *rhino.Node, parent *rhino.Node) {
	typeDoc := n.GetJSDocInfo()
	if typeDoc == nil {
		return
	}

	var owningType rhino.ObjectType
	if parent.IsClassMembers() {
		ctorType := rhino.JSType_ToMaybeFunctionType(parent.GetParent().GetJSType())
		if ctorType == nil {
			return
		}

		if n.IsStaticMember() {
			owningType = ctorType
		} else if n.IsMemberFieldDef() {
			owningType = ctorType.GetInstanceType()
		} else {
			owningType = ctorType.GetPrototype()
		}
	} else {
		owningType = inferJSDocInfo_dereferenced(parent.GetJSType())
	}
	if owningType == nil {
		return
	}

	var propName rhino.Property_Key
	if n.IsComputedFieldDef() || n.IsComputedProp() {
		keyType := n.GetFirstChild().GetJSType()
		if keyType == nil || !keyType.IsKnownSymbolValueType() {
			return
		}
		propName = rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType())
	} else {
		propName = rhino.NewProperty_StringKey(n.GetString())
	}
	if owningType.HasOwnPropertyByKey(propName) && owningType.GetPropertyJSDocInfoByKey(propName) == nil {
		owningType.SetPropertyJSDocInfoByKey(propName, typeDoc)
	}

	// NOTE(lharker): it seems surprising that this doesn't also call
	// attachJSDocInfoToNominalTypeOrShape. I don't know for sure this is a bug, but leaving a
	// comment to document that this is probably not intentional.
}

// Java: InferJSDocInfo.inferJSDocForProperty(Node, Node)
func (p *InferJSDocInfo) inferJSDocForProperty(n *rhino.Node, parent *rhino.Node) {
	// Infer JSDocInfo on properties.
	// There are two ways to write doc comments on a property.
	var typeDoc *rhino.JSDocInfo
	var aliasedType rhino.JSType // if the rhs is a qualified name, its type
	if parent.IsAssign() && n.IsFirstChildOf(parent) {
		// Case: `/** @deprecated */ obj.prop = ...;`
		typeDoc = parent.GetJSDocInfo()
		rhs := n.GetNext()
		// Example: for `/** @deprecated */ obj.prop = obj.newProp;`, the type of `obj.newProp`.
		if rhs.IsQualifiedName() {
			aliasedType = rhs.GetJSType()
		}
	} else if parent.IsExprResult() {
		// Case: `/** @deprecated */ obj.prop;`
		typeDoc = n.GetJSDocInfo()
		aliasedType = nil // there's no value being assigned here
	} else {
		return
	}

	if typeDoc == nil || !typeDoc.ContainsDeclaration() {
		return
	}

	lhsType := inferJSDocInfo_dereferenced(n.GetFirstChild().GetJSType())
	if lhsType == nil {
		return
	}

	// Put the JSDoc in the property slot, if there is one.
	propName := n.GetString()
	if lhsType.HasOwnProperty(propName) && lhsType.GetPropertyJSDocInfo(propName) == nil {
		lhsType.SetPropertyJSDocInfo(propName, typeDoc)
	}

	propType := inferJSDocInfo_dereferenced(lhsType.GetPropertyType(propName))
	if inferJSDocInfo_shouldAttachJSDocToNominalTypeOrShape(propType, aliasedType) {
		inferJSDocInfo_attachJSDocInfoToNominalTypeOrShape(propType, typeDoc, n.GetQualifiedName())
	}
}

// inferJSDocInfo_dereferenced is a nullsafe wrapper for JSType#dereference().
//
// Java: InferJSDocInfo.dereferenced(JSType)
func inferJSDocInfo_dereferenced(type_ rhino.JSType) rhino.ObjectType {
	if type_ == nil {
		return nil
	}
	return type_.Dereference()
}

// inferJSDocInfo_shouldAttachJSDocToNominalTypeOrShape: type_ is the type to
// which we're considering attaching JSDoc, or null; aliasedType, if not
// null, is the JSType of the qualified name on the right-hand side of this
// assignment. e.g. for `const x = y;`, the type of `y`. null if there's no
// value (e.g. `let x;` with no right-hand side) or if the right-hand side is
// not a qualified name (e.g. `const C = class {};` or `const result =
// fn();`).
//
// Java: InferJSDocInfo.shouldAttachJSDocToNominalTypeOrShape(ObjectType, JSType)
func inferJSDocInfo_shouldAttachJSDocToNominalTypeOrShape(
	type_ rhino.ObjectType, aliasedType rhino.JSType) bool {
	// If we have no type, or the type already has a JSDocInfo, then no need to attach `typeDoc`.
	// Also check if we're just aliasing something else's type rather than defining a new type.
	// For example, consider this declaration:
	//   /** @const */
	//   var cbAlias = cb;
	// This is a declaration of the name `cbAlias` but doesn't declare a new type. So `type` should
	// equal `aliasedType`, and we shouldn't attach the `/** @const */`
	// JSDoc to the type of `cb`.
	return type_ != nil && type_.GetJSDocInfo() == nil && !type_.Equals(aliasedType)
}

// inferJSDocInfo_attachJSDocInfoToNominalTypeOrShape handles cases #1 and
// #3 in the class doc.
//
// Java: InferJSDocInfo.attachJSDocInfoToNominalTypeOrShape(ObjectType, JSDocInfo, String)
func inferJSDocInfo_attachJSDocInfoToNominalTypeOrShape(
	objType rhino.ObjectType, docInfo *rhino.JSDocInfo, qName *string) {

	if objType.IsConstructor() || objType.IsInterface() {
		if !inferJSDocInfo_isReferenceNameOf(objType, qName) {
			return
		}

		objType.SetJSDocInfo(docInfo)
		rhino.JSType_ToMaybeFunctionType(objType).GetInstanceType().SetJSDocInfo(docInfo)
	} else if objType.IsEnumType() {
		// Given: `/** @enum {number} */ MyEnum = { FOO: 0 };`
		// Then: typeOf(MyEnum).referenceName() == "enum{MyEnum}"
		// Then: typeOf(MyEnum.FOO).referenceName() == "MyEnum"
		elementType := objType.ToMaybeEnumType().GetElementsType()
		if !inferJSDocInfo_isReferenceNameOf(elementType, qName) {
			return
		}

		objType.SetJSDocInfo(docInfo)
		elementType.SetJSDocInfo(docInfo)
	} else if !objType.IsNativeObjectType() && objType.IsFunctionType() {
		// Anonymous function types identified by their parameter and return types. Remember there can
		// be many unique but equal instances.
		objType.SetJSDocInfo(docInfo)
	}
}

// Java: InferJSDocInfo.isReferenceNameOf(ObjectType, String)
func inferJSDocInfo_isReferenceNameOf(type_ rhino.ObjectType, name *string) bool {
	if !type_.HasReferenceName() {
		return false
	}
	// type.getReferenceName().equals(name): a null name is never equal.
	return name != nil && *type_.GetReferenceName() == *name
}

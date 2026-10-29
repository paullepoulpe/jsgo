// Port of com/google/javascript/jscomp/RemoveUnusedCode.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TraverseNormalOrOptChainGetProp traverses `expr.prop` or `expr?.prop`.
//
// Note that this method is called only for RHS nodes. Property references
// that are being assigned to are handled by the logic traversing their
// parent (e.g. ASSIGN) node.
//
// The primary purpose of this method is to make sure the property reference
// is correctly recorded.
//
// Java: RemoveUnusedCode.traverseNormalOrOptChainGetProp(Node, Scope)
func (r *RemoveUnusedCode) TraverseNormalOrOptChainGetProp(getProp rhino.Node, scope *Scope) {
	access, ok := ast.AsNormalOrOptChainGetProp(getProp)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", getProp)
	}
	objectNode := access.Object().Node()
	propID := getProp.GetStringID()

	if r.polyBits.Get(int(propID)) {
		r.considerPolyfillRef(propID, getProp)
	}

	if NodeUtil_IsExpressionResultUsed(getProp) ||
		r.ConsiderForAccessorSideEffects(getProp, AccessorSummary_PropertyAccessKind_GETTER_ONLY) {
		// must record as reference to the property and continue traversal.
		r.markPropertyIDPinned(propID)
		r.TraverseNode(objectNode, scope)
	} else if objectNode.IsThis() {
		// This is probably the declaration of a class field in a constructor.
		// /** @private {number} */
		// this.propName;
		// We don't want to consider this a real usage that should prevent removal.
		r.ConsiderForIndependentRemoval(r.newUnusedReadRefRemovable(getProp, remFlagThisDotPropRef))
	} else if owner, ok := removeUnusedCode_prototypeOwner(objectNode); ok {
		// (objExpression).prototype.propName;
		objExpression := owner.Node()
		if objExpression.IsName() {
			// name.prototype.propName;
			varInfo := r.TraverseNameNode(objExpression, scope)
			varInfo.AddRemovable(r.newUnusedReadRefRemovable(getProp, remFlagProtoDotPropRef))
		} else {
			// (objExpression).prototype.propName;
			rem := r.newUnusedReadRefRemovable(getProp, remFlagProtoDotPropRef)
			if r.MayHaveSideEffects(objExpression) {
				r.TraverseNode(objExpression, scope)
			} else {
				rem.addContinuation(r, objExpression, scope)
			}
			r.ConsiderForIndependentRemoval(rem)
		}
	} else {
		// TODO(bradfordcsmith): add removal of `varName.propName;`
		r.markPropertyIDPinned(propID)
		r.TraverseNode(objectNode, scope)
	}
}

// TraverseObjectDefinePropertiesCall traverses
// `Object.defineProperties(someObject, propertyDefinitions);`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesCall(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesCall(callNode rhino.Node, scope *Scope) {
	// First child is Object.defineProperties or some equivalent of it.
	call := ast.MustNormalOrOptChainCall(callNode)
	callee := call.Callee().Node()
	target, _ := call.FirstArg() // isObjectDefinePropertiesDefinition: two arguments
	definitions, _ := target.Next()
	targetObject, propertyDefinitions := target.Node(), definitions.Node()

	// The name that the target object is or is the prototype of.
	nameNode, ok := target.AsName()
	if !ok {
		nameNode, ok = removeUnusedCode_nameDotPrototype(targetObject)
	}
	if ok && !NodeUtil_IsExpressionResultUsed(callNode) {
		// NOTE: Object.defineProperties() returns its first argument, so if its return value is used
		// that counts as a use of the targetObject.
		varInfo := r.TraverseNameNode(nameNode.Node(), scope)
		rem := r.newClassSetupCallRemovable(callNode, false)
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		rem.addContinuation(r, callee, scope)
		if r.MayHaveSideEffects(propertyDefinitions) {
			r.TraverseNode(propertyDefinitions, scope)
		} else {
			rem.addContinuation(r, propertyDefinitions, scope)
		}
		varInfo.AddRemovable(rem)
	} else {
		// TODO(bradfordcsmith): Is it really necessary to traverse the callee
		// (aka. Object.defineProperties)?
		r.TraverseNode(callee, scope)
		r.TraverseNode(targetObject, scope)
		r.TraverseNode(propertyDefinitions, scope)
	}
}

// TraverseObjectDefinePropertiesLiteral traverses the object literal passed
// as the second argument to `Object.defineProperties()`.
//
// Java: RemoveUnusedCode.traverseObjectDefinePropertiesLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectDefinePropertiesLiteral(propertyDefinitions rhino.Node, scope *Scope) {
	for member := range ast.MustObjectLit(propertyDefinitions).Members() {
		property := member.Node()
		if key, ok := member.AsStringKey(); ok && key.IsQuoted() {
			// Quoted property name counts as a reference to the property and protects it from removal.
			r.markPropertyIDPinned(property.GetStringID())
			r.TraverseNode(key.Value().Node(), scope)
		} else if ok {
			definition := key.Value().Node()
			if r.MayHaveSideEffects(definition) {
				r.TraverseNode(definition, scope)
			} else {
				rem := r.newObjectDefinePropsDefRemovable(property)
				rem.addContinuation(r, definition, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
		} else {
			// TODO(bradfordcsmith): Maybe report error for anything other than a computed property,
			// since getters, setters, and methods don't make much sense in this context.
			r.TraverseNode(property, scope)
		}
	}
}

// Java: RemoveUnusedCode.traverseObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	if !objectLiteral.IsObjectLit() {
		jutil.CheckArgumentMsg1(false, "%s", objectLiteral)
	}
	// Is this an object literal that is assigned directly to a 'prototype' property?
	if r.IsAssignmentToPrototype(objectLiteral.GetParent()) {
		r.TraversePrototypeLiteral(objectLiteral, scope)
	} else if r.IsObjectDefinePropertiesSecondArgument(objectLiteral) {
		// TODO(bradfordcsmith): Consider restricting special handling of the properties literal to
		// cases where the target object is a known class, prototype, or this.
		r.TraverseObjectDefinePropertiesLiteral(objectLiteral, scope)
	} else {
		r.TraverseNonPrototypeObjectLiteral(objectLiteral, scope)
	}
}

// Java: RemoveUnusedCode.isObjectDefinePropertiesSecondArgument(Node)
func (r *RemoveUnusedCode) IsObjectDefinePropertiesSecondArgument(n rhino.Node) bool {
	call, ok := ast.AsCall(n.GetParent())
	if !ok {
		return false
	}
	last, ok := call.LastArg()
	return ok && last.Is(n) && NodeUtil_IsObjectDefinePropertiesDefinition(call.Node())
}

// Java: RemoveUnusedCode.traverseNonPrototypeObjectLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraverseNonPrototypeObjectLiteral(objectLiteral rhino.Node, scope *Scope) {
	for member := range ast.MustObjectLit(objectLiteral).Members() {
		if key, ok := member.AsStringKey(); ok {
			// A property name in an object literal counts as a reference,
			// because of some reflection patterns.
			// Note that we are intentionally treating both quoted and unquoted keys as
			// references.
			r.markPropertyIDPinned(member.Node().GetStringID())
			r.TraverseNode(key.Value().Node(), scope)
		} else {
			r.TraverseNode(member.Node(), scope)
		}
	}
}

// Java: RemoveUnusedCode.traversePrototypeLiteral(Node, Scope)
func (r *RemoveUnusedCode) TraversePrototypeLiteral(objectLiteral rhino.Node, scope *Scope) {
	for propertyNode := objectLiteral.GetFirstChild(); !propertyNode.IsNil(); propertyNode = propertyNode.GetNext() {
		if propertyNode.IsComputedProp() || propertyNode.IsQuotedStringKey() {
			r.TraverseChildren(propertyNode, scope)
		} else {
			valueNode := propertyNode.GetOnlyChild()
			if r.MayHaveSideEffects(valueNode) {
				// TODO(bradfordcsmith): Ideally we should preserve the side-effect without keeping the
				// property itself alive.
				r.TraverseNode(valueNode, scope)
			} else {
				// If we've come this far, we already know we're keeping the prototype literal itself,
				// but we may be able to remove unreferenced properties in it.
				rem := r.newClassOrProtoPropRemovable(propertyNode)
				rem.addContinuation(r, valueNode, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
		}
	}
}

// Java: RemoveUnusedCode.isAssignmentToPrototype(Node)
func (r *RemoveUnusedCode) IsAssignmentToPrototype(n rhino.Node) bool {
	assign, ok := ast.AsAssign(n)
	return ok && RemoveUnusedCode_IsDotPrototype(assign.Target().Node())
}

// RemoveUnusedCode_IsDotPrototype: true for `someExpression.prototype`.
//
// Java: RemoveUnusedCode.isDotPrototype(Node)
func RemoveUnusedCode_IsDotPrototype(n rhino.Node) bool {
	_, ok := removeUnusedCode_prototypeOwner(n)
	return ok
}

// removeUnusedCode_prototypeOwner returns someExpression of n =
// someExpression.prototype (or ?.prototype); ok is false for any other n.
func removeUnusedCode_prototypeOwner(n rhino.Node) (ast.Expr, bool) {
	access, ok := ast.AsNormalOrOptChainGetProp(n)
	if !ok || n.GetStringID() != rucPropPrototypeID {
		return ast.Expr{}, false
	}
	return access.Object(), true
}

// removeUnusedCode_nameDotPrototype returns name of n = name.prototype; ok
// is false for any other n.
func removeUnusedCode_nameDotPrototype(n rhino.Node) (ast.Name, bool) {
	getProp, ok := ast.AsGetProp(n)
	if !ok || n.GetStringID() != rucPropPrototypeID {
		return ast.Name{}, false
	}
	return getProp.Object().AsName()
}

// Java: RemoveUnusedCode.traverseAssign(Node, Scope)
func (r *RemoveUnusedCode) TraverseAssign(assignNode rhino.Node, scope *Scope) {
	assign, ok := ast.AsAssignment(assignNode)
	if !ok {
		jutil.CheckState(false)
	}

	lhs := assign.Target().Node()
	valueNode := assign.Value().Node()
	if _, ok := assign.Target().AsName(); ok {
		// varName = something
		varInfo := r.TraverseNameNode(lhs, scope)
		rem := r.newVarAssignRemovable(assignNode, varInfo)
		r.traverseRemovableAssignValue(rem, valueNode, scope)
		varInfo.AddRemovable(rem)
	} else if getElem, ok := assign.Target().AsGetElem(); ok {
		getElemObj := getElem.Object().Node()
		getElemKey := getElem.Index().Node()
		varName, ok := getElem.Object().AsName()
		if !ok {
			varName, ok = removeUnusedCode_nameDotPrototype(getElemObj)
		}

		if ok {
			// varName[someExpression] = someValue
			// OR
			// varName.prototype[someExpression] = someValue
			varInfo := r.TraverseNameNode(varName.Node(), scope)
			rem := r.newComputedPropAssignRemovable(assignNode, getElemKey, varInfo)
			if r.MayHaveSideEffects(getElemKey) {
				r.TraverseNode(getElemKey, scope)
			} else {
				rem.addContinuation(r, getElemKey, scope)
			}
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			varInfo.AddRemovable(rem)
		} else {
			r.TraverseNode(getElemObj, scope)
			r.TraverseNode(getElemKey, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else if getProp, ok := assign.Target().AsGetProp(); ok {
		getPropLhs := getProp.Object().Node()
		// Assignments `Foo.prototype.bar = function() {`
		owner, isDotPrototypeLhs := removeUnusedCode_prototypeOwner(getPropLhs)
		isPrototypeMethodDef := isDotPrototypeLhs && valueNode.IsFunction()

		if !isPrototypeMethodDef &&
			r.ConsiderForAccessorSideEffects(lhs, AccessorSummary_PropertyAccessKind_SETTER_ONLY) {
			// And the possible side-effects mean we can't do any removal. We don't use the
			// `AstAnalyzer` because we only want to consider side-effect from the assignment, not the
			// entire l-value subtree.
			// Assume prototype method assignments never trigger setters, matching ES class semantics
			r.TraverseNode(getPropLhs, scope) // Don't re-traverse the GETPROP as a read.
			r.TraverseNode(valueNode, scope)
		} else if getPropLhs.IsName() {
			// varName.propertyName = someValue
			varInfo := r.TraverseNameNode(getPropLhs, scope)
			rem := r.newNamedPropAssignRemovable(assignNode, lhs, varInfo, 0)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			varInfo.AddRemovable(rem)
		} else if isDotPrototypeLhs {
			// objExpression.prototype.propertyName = someValue
			objExpression := owner.Node()
			if objExpression.IsName() {
				// varName.prototype.propertyName = someValue
				rem := r.newNamedPropAssignRemovable(assignNode, lhs, nil, remFlagProtoDotPropRef)
				r.traverseRemovableAssignValue(rem, valueNode, scope)
				varInfo := r.TraverseNameNode(objExpression, scope)
				rem.varInfo = varInfo
				varInfo.AddRemovable(rem)
			} else {
				// (someExpression).prototype.propertyName = someValue
				rem := r.newAnonProtoPropAssignRemovable(assignNode, lhs)
				r.traverseRemovableAssignValue(rem, valueNode, scope)
				if r.MayHaveSideEffects(objExpression) {
					r.TraverseNode(objExpression, scope)
				} else {
					rem.addContinuation(r, objExpression, scope)
				}
				r.ConsiderForIndependentRemoval(rem)
			}
		} else if getPropLhs.IsThis() {
			// this.propertyName = someValue
			rem := r.newNamedPropAssignRemovable(assignNode, lhs, nil, remFlagThisDotPropRef)
			r.traverseRemovableAssignValue(rem, valueNode, scope)
			r.ConsiderForIndependentRemoval(rem)
		} else {
			r.TraverseNode(lhs, scope)
			r.TraverseNode(valueNode, scope)
		}
	} else {
		// no other cases are removable
		r.TraverseNode(lhs, scope)
		r.TraverseNode(valueNode, scope)
	}
}

// Java: RemoveUnusedCode.isNameDotPrototype(Node)
func (r *RemoveUnusedCode) IsNameDotPrototype(n rhino.Node) bool {
	_, ok := removeUnusedCode_nameDotPrototype(n)
	return ok
}

// Java: RemoveUnusedCode.traverseObjectPattern(Node, Scope)
func (r *RemoveUnusedCode) TraverseObjectPattern(pattern rhino.Node, scope *Scope) {
	objectPattern, ok := ast.AsObjectPattern(pattern)
	if !ok {
		jutil.CheckStateMsg1(false, "%s", pattern)
	}

	for prop := range objectPattern.Props() {
		elem := prop.Node()
		switch prop.Token() {
		case rhino.Token_COMPUTED_PROP:
			value, _ := prop.MustComputedProp().Value() // always there in a pattern
			r.TraverseIndirectAssignment(elem, value.Node(), scope)
		case rhino.Token_STRING_KEY:
			key := prop.MustStringKey()
			if !key.IsQuoted() {
				r.markPropertyIDPinned(elem.GetStringID())
			}
			r.TraverseIndirectAssignment(elem, key.Value().Node(), scope)
		case rhino.Token_OBJECT_REST:
			// Recall that the rest target can be any l-value expression
			r.TraverseIndirectAssignment(elem, prop.MustObjectRest().Target().Node(), scope)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of " + pattern.GetToken().String() + ": " + elem.ToStringTree()))
		}
	}
}

// TraverseClass handles a class that is not the RHS child of an assignment
// or a variable declaration initializer.
//
// Java: RemoveUnusedCode.traverseClass(Node, Scope)
func (r *RemoveUnusedCode) TraverseClass(classNode rhino.Node, scope *Scope) {
	jutil.CheckArgument(classNode.IsClass())
	if NodeUtil_IsClassDeclaration(classNode) {
		r.TraverseClassDeclaration(classNode, scope)
	} else {
		r.TraverseClassExpression(classNode, scope)
	}
}

// Java: RemoveUnusedCode.traverseClassDeclaration(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassDeclaration(classNode rhino.Node, scope *Scope) {
	class, ok := ast.AsClass(classNode)
	jutil.CheckArgument(ok)
	className, _ := class.Name() // a declaration has one
	// A missing extends (an EMPTY node) has nothing to traverse.
	baseClass, hasBaseClass := class.Extends()
	classBodyNode := class.Body().Node()
	// Class declarations never declare the class name in the inner class scope (only class
	// expressions do), so creating an inner scope here would always produce an empty scope.
	classScope := scope

	varInfo := r.TraverseNameNode(className.Node(), scope)
	if classNode.GetParent().IsExport() {
		// Cannot remove an exported class.
		varInfo.pin(rucPinExportedClass)
		if hasBaseClass {
			r.TraverseNode(baseClass.Node(), scope)
		}
		// Use traverseChildren() here, because we should not consider any properties on the exported
		// class to be removable.
		r.TraverseChildren(classBodyNode, classScope)
	} else if hasBaseClass && r.MayHaveSideEffects(baseClass.Node()) {
		// TODO(bradfordcsmith): implement removal without losing side-effects for this case
		varInfo.pin(rucPinBaseClassSideEffects)
		r.TraverseNode(baseClass.Node(), scope)
		r.TraverseClassMembers(classBodyNode, classScope)
	} else if r.MayHaveSideEffects(classBodyNode) {
		varInfo.pin(rucPinClassBodySideEffects)
		if hasBaseClass {
			r.TraverseNode(baseClass.Node(), scope)
		}
		r.TraverseClassMembers(classBodyNode, classScope)
	} else {
		rem := r.newClassDeclRemovable(classNode)
		if hasBaseClass {
			rem.addContinuation(r, baseClass.Node(), classScope)
		}
		rem.addContinuation(r, classBodyNode, classScope)
		varInfo.AddRemovable(rem)
	}
}

// Java: RemoveUnusedCode.traverseClassExpression(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassExpression(classNode rhino.Node, scope *Scope) {
	class, ok := ast.AsClass(classNode)
	jutil.CheckArgument(ok)
	classScope := scope

	if className, ok := class.Name(); ok {
		classScope = r.scopeCreator.CreateScopeOfScope(classNode, scope)
		// We may be able to remove the name node if nothing ends up referring to it.
		varInfo := r.TraverseNameNode(className.Node(), classScope)
		// The class is non-local, because it is accessible by unknown code outside
		// of the scope where InnerName is defined.
		// e.g. `use(class InnerName {})`
		varInfo.SetHasNonLocalOrNonLiteralValue()
		varInfo.AddRemovable(r.newNamedClassExprRemovable(classNode))
	}
	// If we're traversing the class expression, we've already decided we cannot remove it.
	if baseClass, ok := class.Extends(); ok {
		r.TraverseNode(baseClass.Node(), scope)
	}
	r.TraverseClassMembers(class.Body().Node(), classScope)
}

// Java: RemoveUnusedCode.traverseClassMembers(Node, Scope)
func (r *RemoveUnusedCode) TraverseClassMembers(node rhino.Node, scope *Scope) {
	jutil.CheckArgumentMsg1(node.IsClassMembers(), "%s", node)
	if !r.removeUnusedPrototypeProperties {
		r.TraverseChildren(node, scope)
		return
	}

	for member := node.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		switch member.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			// If we get as far as traversing the members of a class, we've already decided that
			// we cannot remove the class itself, so just consider individual members for removal.
			rem := r.newClassOrProtoPropRemovable(member)
			rem.addContinuation(r, member, scope)
			r.ConsiderForIndependentRemoval(rem)
		case rhino.Token_MEMBER_FIELD_DEF:
			// TODO(bradfordcsmith): currently if the RHS of a field has side effects, we do not
			// remove any part of the field. The proper behavior of class C { x = alert(); }
			// would be to remove x, leaving class C { constructor() { alert(); } }
			// but currently we aren't removing anything.
			value, hasValue := ast.MustMemberFieldDef(member).Value()
			if !hasValue || !r.MayHaveSideEffects(value.Node()) {
				rem := r.newClassOrProtoPropRemovable(member)
				rem.addContinuation(r, member, scope)
				r.ConsiderForIndependentRemoval(rem)
			}
			// TODO: b/354704593 - remove the entire class when it is not referenced
			if hasValue {
				r.TraverseChildren(member, scope)
			}
		case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
			r.TraverseChildren(member, scope)
		case rhino.Token_BLOCK:
			r.TraverseChildren(member, r.scopeCreator.CreateScopeOfScope(member, scope))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of CLASS_MEMBERS: " + member.ToStringTree()))
		}
	}
}

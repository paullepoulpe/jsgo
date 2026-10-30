// Port of com/google/javascript/jscomp/TypeInference.java @ closure-compiler v20261021

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

// TypeInference_TypeDeclaringCallback abstracts logic for declaring an
// lvalue in a particular scope.
//
// It updates the given scope upon seeing an assignment or declaration:
// scope is the scope we are in, lvalue the value being updated (a NAME,
// GETPROP, GETELEM, or CAST), type_ the type we've inferred for the lvalue;
// it returns the updated flow scope.
//
// Java: TypeInference.TypeDeclaringCallback.declareTypeInScope(FlowScope, Node, JSType)
//
// Go-only: language-adapter — a single-method interface implemented only by
// lambdas, so a Go func type.
type TypeInference_TypeDeclaringCallback func(scope *LinkedFlowScope, lvalue rhino.Node, type_ rhino.JSType) *LinkedFlowScope

// InferParameters infers all of a function's parameters if their types
// aren't declared.
//
// Java: TypeInference.inferParameters(FlowScope)
func (ti *TypeInference) InferParameters(entryFlowScope *LinkedFlowScope) *LinkedFlowScope {
	functionNode := ti.containerScope.GetRootNode()
	fn, ok := ast.AsFunction(functionNode)
	if !ok {
		return entryFlowScope // we're in the global scope
	} else if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		// Pretend the function literal in `goog.loadModule(function(exports) {` does not exist.
		return entryFlowScope
	}
	astParameters := fn.Params()
	var iifeArgumentNode rhino.Node

	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}

	functionType := rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	parameterTypes := functionType.GetParameters()
	nextParameter := 0
	var parameter rhino.FunctionType_Parameter
	if nextParameter < len(parameterTypes) {
		parameter = parameterTypes[nextParameter]
		nextParameter++
	}

	// This really iterates over three different things at once:
	//   - the actual AST parameter nodes (which may be REST, DEFAULT_VALUE, etc.)
	//   - the argument nodes in an IIFE
	//   - the parameter type nodes from the FunctionType on the FUNCTION node
	// Always visit every AST parameter once, regardless of how many IIFE arguments or
	// FunctionType param nodes there are.
	for astParamItr := range astParameters.Params() {
		astParam := astParamItr.Node()
		if !iifeArgumentNode.IsNil() && iifeArgumentNode.IsSpread() {
			// block inference on all parameters that might possibly be set by a spread, e.g. `z` in
			// (function f(x, y, z = 1))(...[1, 2], 'foo')
			iifeArgumentNode = rhino.NilNode
		}

		// Running variable for the type of the param within the body of the function. We use the
		// existing type on the param node as the default, and then transform it according to the
		// declaration syntax.
		inferredType := ti.GetJSType(astParam)

		if !iifeArgumentNode.IsNil() {
			if iifeArgumentNode.GetJSType() != nil {
				inferredType = iifeArgumentNode.GetJSType()
			}
		} else if parameter.GetJSType() != nil {
			inferredType = parameter.GetJSType()
		}

		var defaultValue rhino.Node
		if defVal, ok := astParamItr.AsDefaultValue(); ok {
			defaultValue = defVal.Value().Node()
			// must call `traverse` to correctly type the default value
			entryFlowScope = ti.Traverse(defaultValue, entryFlowScope)
			astParam = defVal.Target().Node()
		} else if rest, ok := astParamItr.AsIterRest(); ok {
			// e.g. `function f(p1, ...restParamName) {}`
			// set astParam = restParamName
			astParam = rest.Target().Node()
			// convert 'number' into 'Array<number>' for rest parameters
			inferredType =
				ti.registry.CreateTemplatizedType(ti.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), inferredType)
		}

		if !defaultValue.IsNil() {
			// The param could possibly be the default type, and `undefined` args won't propagate in.
			inferredType =
				ti.registry.CreateUnionType(
					inferredType.RestrictByNotUndefined(), ti.GetJSType(defaultValue))
		}

		if astParam.IsDestructuringPattern() {
			// even if the inferredType is null, we still need to type all the nodes inside the
			// destructuring pattern. (e.g. in computed properties or default value expressions)
			entryFlowScope = ti.UpdateDestructuringParameter(astParam, inferredType, entryFlowScope)
		} else {
			// for simple named parameters, we only need to update the scope/AST if we have a new
			// inferred type.
			entryFlowScope =
				ti.UpdateNamedParameter(astParam, !defaultValue.IsNil(), inferredType, entryFlowScope)
		}

		parameter = rhino.FunctionType_Parameter{}
		if nextParameter < len(parameterTypes) {
			parameter = parameterTypes[nextParameter]
			nextParameter++
		}
		if !iifeArgumentNode.IsNil() {
			iifeArgumentNode = iifeArgumentNode.GetNext()
		}
	}

	return entryFlowScope
}

// UpdateDestructuringParameter sets the types of a un-named/destructuring
// function parameter to an inferred type.
//
// This method is responsible for typing:
//   - The scope slot
//   - The pattern nodes
//
// Java: TypeInference.updateDestructuringParameter(Node, JSType, FlowScope)
func (ti *TypeInference) UpdateDestructuringParameter(
	pattern rhino.Node, inferredType rhino.JSType, entryFlowScope *LinkedFlowScope) *LinkedFlowScope {
	// look through all expressions and lvalues in the pattern.
	// given an lvalue, change its type if either a) it's inferred (not declared in
	// TypedScopeCreator) or b) it has a default value
	entryFlowScope =
		ti.TraverseDestructuringPatternHelper(
			pattern,
			entryFlowScope,
			inferredType,
			func(scope *LinkedFlowScope, lvalue rhino.Node, type_ rhino.JSType) *LinkedFlowScope {
				var_ := jutil.CheckNotNullPtr(ti.containerScope.GetVar(lvalue.GetString()))
				// This condition will trigger on cases like
				//   (function f({x}) {})({x: 3})
				// where `x` is of unknown type during the typed scope creation phase, but
				// here we can infer that it is of type `number`
				if var_.IsTypeInferred() {
					var_.SetType(type_)
					lvalue.SetJSType(type_)
				}
				if lvalue.GetParent().IsDefaultValue() {
					// Given
					//   /** @param {{age: (number|undefined)}} data */
					//   function f({age = 99}) {}
					// infer that `age` is now a `number` and not `number|undefined`
					// treat this similarly to if there was an assignment inside the function body
					// TODO(b/117162687): allow people to narrow the declared type to
					// exclude 'undefined' inside the function body.
					scope = ti.UpdateScopeForAssignment(scope, lvalue, type_, TypeInference_AssignmentType_ASSIGN)
				}
				return scope
			})

	return entryFlowScope
}

// UpdateNamedParameter sets the types of a named/non-destructuring function
// parameter to an inferred type.
//
// This method is responsible for typing:
//   - The scope slot
//   - The param node
//
// Java: TypeInference.updateNamedParameter(Node, boolean, JSType, FlowScope)
func (ti *TypeInference) UpdateNamedParameter(
	paramName rhino.Node, hasDefaultValue bool, inferredType rhino.JSType, entryFlowScope *LinkedFlowScope) *LinkedFlowScope {
	var_ := ti.containerScope.GetVar(paramName.GetString())
	if var_ == nil {
		panic(jutil.NewNullPointerExceptionMsg(jutil.StringFormat("Missing var for parameter %s", paramName)))
	}

	paramName.SetJSType(inferredType)

	if var_.IsTypeInferred() {
		var_.SetType(inferredType)
	} else if hasDefaultValue {
		// If this is a declared type with a default value, update the LinkedFlowScope slots but not
		// the actual TypedVar. This is similar to what would happen if the default value was moved
		// into an assignment in the fn body
		entryFlowScope = ti.RedeclareSimpleVar(entryFlowScope, paramName, inferredType)
	}
	return entryFlowScope
}

// Java: TypeInference.traverseAssignUnsignedRightShift(Node, FlowScope)
func (ti *TypeInference) TraverseAssignUnsignedRightShift(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	left := ast.MustAssignOp(n).Target().Node()
	scope = ti.TraverseUnsignedRightShift(n, scope)
	return ti.UpdateScopeForAssignmentWithUpdateNode(
		scope, left, ti.GetJSType(n) /* updateNode= */, rhino.NilNode, TypeInference_AssignmentType_ASSIGN)
}

// TraverseCatch: any value can be thrown, so it's really impossible to
// determine the type of a CATCH param. Treat it as the UNKNOWN type.
//
// Java: TypeInference.traverseCatch(Node, FlowScope)
func (ti *TypeInference) TraverseCatch(catchNode rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	catch := ast.MustCatch(catchNode)
	param, ok := catch.Param()
	if !ok {
		// ES2019 allows `try {} catch {}` with no catch expression
		return scope
	}
	if name, isName := param.AsName(); isName {
		nameNode := name.Node()
		type_ := nameNode.GetJSType()
		// If the catch expression name was declared in the catch in TypedScopeCreator use that type.
		// Otherwise use "unknown".
		if type_ == nil {
			type_ = ti.unknownType
			nameNode.SetJSType(ti.unknownType)
		}
		return ti.RedeclareSimpleVar(scope, nameNode, type_)
	}
	return ti.TraverseDestructuringPattern(param.Node(), scope, ti.unknownType, TypeInference_AssignmentType_DECLARATION)
}

// Java: TypeInference.traverseAssign(Node, FlowScope)
func (ti *TypeInference) TraverseAssign(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	assign := ast.MustAssign(n)
	target := assign.Target().Node()
	value := assign.Value().Node()
	if target.IsDestructuringPattern() {
		scope = ti.Traverse(value, scope)
		valueType := ti.GetJSType(value)
		n.SetJSType(valueType)
		return ti.TraverseDestructuringPattern(target, scope, valueType, TypeInference_AssignmentType_ASSIGN)
	}
	scope = ti.TraverseChildren(n, scope)

	valueType := ti.GetJSType(value)
	n.SetJSType(valueType)

	return ti.UpdateScopeForAssignment(scope, target, valueType, TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseAssignOp(Node, FlowScope)
func (ti *TypeInference) TraverseAssignOp(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	left := ast.MustAssignOp(n).Target().Node()
	scope = ti.TraverseBigIntCompatibleBinaryOperator(n, scope)

	// The lhs is both an input and an output, so don't update the input type here.
	return ti.UpdateScopeForAssignmentWithUpdateNode(
		scope, left, ti.GetJSType(n) /* updateNode= */, rhino.NilNode, TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.isInExternFile(Node)
func TypeInference_IsInExternFile(n rhino.Node) bool {
	return NodeUtil_GetSourceFile(n).IsExtern()
}

// Java: TypeInference.isPossibleMixinApplication(Node, Node)
func TypeInference_IsPossibleMixinApplication(lvalue rhino.Node, rvalue rhino.Node) bool {
	if TypeInference_IsInExternFile(lvalue) {
		return true
	}

	jsdoc := NodeUtil_GetBestJSDocInfo(lvalue)
	return jsdoc != nil &&
		jsdoc.IsConstructor() &&
		jsdoc.GetImplementedInterfaceCount() > 0 &&
		lvalue.IsQualifiedName() &&
		(!rvalue.IsNil() && rvalue.IsCall())
}

// TypeInference_AddMissingInterfaceProperties: constructor is a constructor
// function defined by a call, which may be a mixin application. The
// constructor implements at least one interface. If the constructor is
// missing some properties of the inherited interfaces, this method declares
// these properties.
//
// Java: TypeInference.addMissingInterfaceProperties(JSType)
func TypeInference_AddMissingInterfaceProperties(constructor rhino.JSType) {
	if constructor == nil || !constructor.IsConstructor() {
		return
	}
	f := constructor.ToMaybeFunctionType()
	proto := f.GetPrototype()
	for _, interf := range f.GetImplementedInterfaces() {
		for _, pname := range interf.GetPropertyNames() {
			if !proto.HasProperty(pname) {
				proto.DefineDeclaredProperty(pname, interf.GetPropertyType(pname), rhino.NilNode)
			}
		}
	}
}

// TypeInference_AssignmentType is either a combined
// declaration/initialization or a regular assignment.
type TypeInference_AssignmentType int32

const (
	TypeInference_AssignmentType_DECLARATION TypeInference_AssignmentType = iota // var x = 3;
	TypeInference_AssignmentType_ASSIGN                                          // `a.b.c = d;` or `x = 4;`
)

var typeInference_AssignmentType_names = [...]string{"DECLARATION", "ASSIGN"}

// Name is Java's Enum.name().
func (t TypeInference_AssignmentType) Name() string { return typeInference_AssignmentType_names[t] }

// String is Java's Enum.toString().
func (t TypeInference_AssignmentType) String() string { return t.Name() }

// Ordinal is Java's Enum.ordinal().
func (t TypeInference_AssignmentType) Ordinal() int { return int(t) }

// Java: TypeInference.updateScopeForAssignment(FlowScope, Node, JSType, AssignmentType)
func (ti *TypeInference) UpdateScopeForAssignment(
	scope *LinkedFlowScope, target rhino.Node, resultType rhino.JSType, type_ TypeInference_AssignmentType) *LinkedFlowScope {
	return ti.UpdateScopeForAssignmentWithUpdateNode(scope, target, resultType, target, type_)
}

// UpdateScopeForAssignmentWithUpdateNode updates the scope according to the
// result of an assignment.
//
// Java: TypeInference.updateScopeForAssignment(FlowScope, Node, JSType, Node, AssignmentType)
func (ti *TypeInference) UpdateScopeForAssignmentWithUpdateNode(
	scope *LinkedFlowScope,
	target rhino.Node,
	resultType rhino.JSType,
	updateNode rhino.Node,
	type_ TypeInference_AssignmentType) *LinkedFlowScope {
	if resultType == nil {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState(updateNode.IsNil() || updateNode == target)

	targetType := target.GetJSType() // may be null

	right := NodeUtil_GetRValueOfLValue(target)
	if TypeInference_IsPossibleMixinApplication(target, right) {
		TypeInference_AddMissingInterfaceProperties(targetType)
	}

	switch target.GetToken() {
	case rhino.Token_NAME:
		name := ast.MustName(target)
		varName := name.Ident()
		var_ := TypeInference_GetDeclaredVar(scope, varName)
		var varType rhino.JSType
		if var_ != nil {
			varType = var_.GetType()
		}
		isVarDeclaration :=
			type_ == TypeInference_AssignmentType_DECLARATION &&
				var_ != nil &&
				!var_.IsTypeInferred() &&
				!var_.GetNameNode().IsNil() // implicit vars (like arguments) have no nameNode

		// Whether this variable is declared not because it has JSDoc with a declaration, but
		// because it is const and the right-hand-side is easily inferrable.
		// e.g. these are 'typeless const declarations':
		//   const x = 0;
		//   /** @const */
		//   a.b.c = SomeOtherConstructor;
		// but these are not:
		//    let x = 0;
		//    /** @const @constructor */
		//    a.b.c = someMixin();
		// This is messy, since the definition of 'typeless const' is duplicated in
		// TypedScopeCreator and this code.
		isTypelessConstDecl :=
			isVarDeclaration &&
				var_.GetNameNode().IsName() && // ignore redeclarations of implicit globals
				NodeUtil_IsConstantDeclaration(var_.GetJSDocInfo(), var_.GetNameNode()) &&
				!(var_.GetJSDocInfo() != nil &&
					var_.GetJSDocInfo().ContainsDeclarationExcludingTypelessConst()) &&
				!(var_.GetType() != nil && var_.GetType().IsNoResolvedType())

		// When looking at VAR initializers for declared VARs, we tend
		// to use the declared type over the type it's being
		// initialized to in the global scope.
		//
		// For example,
		// /** @param {number} */ var f = goog.abstractMethod;
		// it's obvious that the programmer wants you to use
		// the declared function signature, not the inferred signature.
		//
		// Or,
		// /** @type {Object.<string>} */ var x = {};
		// the one-time anonymous object on the right side
		// is as narrow as it can possibly be, but we need to make
		// sure we back-infer the <string> element constraint on
		// the left hand side, so we use the left hand side.

		isVarTypeBetter :=
			isVarDeclaration &&
				// Makes it easier to check for NPEs.
				!resultType.IsNullType() &&
				!resultType.IsVoidType() &&
				// Do not use the var type if the declaration looked like
				// /** @const */ var x = 3;
				// because this type was computed from the RHS
				!isTypelessConstDecl

		// TODO(nicksantos): This might be a better check once we have
		// back-inference of object/array constraints.  It will probably
		// introduce more type warnings.  It uses the result type iff it's
		// strictly narrower than the declared var type.
		//
		// boolean isVarTypeBetter = isVarDeclaration &&
		//    (varType.restrictByNotNullOrUndefined().isSubtype(resultType)
		//     || !resultType.isSubtype(varType));

		if isVarTypeBetter {
			scope = ti.redeclareSimpleVar(scope, target, var_, varType)
		} else {
			scope = ti.redeclareSimpleVar(scope, target, var_, resultType)
		}

		if !updateNode.IsNil() {
			updateNode.SetJSType(resultType)
		}

		_, hasInit := name.Init()
		if var_ != nil &&
			var_.IsTypeInferred() &&
			// Don't change the typed scope to include "undefined" upon seeing "let foo;", because
			// this is incompatible with how we currently handle VARs and breaks existing code.
			// TODO(sdh): remove this condition after cleaning up code depending on it.
			!(target.GetParent().IsLet() && !hasInit) {
			oldType := var_.GetType()
			if oldType == nil {
				var_.SetType(resultType)
			} else {
				var_.SetType(oldType.GetLeastSupertype(resultType))
			}
		} else if isTypelessConstDecl {
			// /** @const */ var x = y;
			// should be redeclared, so that the type of y
			// gets propagated to inner scopes.
			var_.SetType(resultType)
		}
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(target)
		qualifiedName, isQName := target.GetQualifiedNameString()
		if isQName {
			declaredSlotType := false
			rawObjType := gp.Object().Node().GetJSType()
			if rawObjType != nil {
				objType := rhino.ObjectType_Cast(rawObjType.RestrictByNotNullOrUndefined())
				if objType != nil {
					propName := gp.Property()
					declaredSlotType = objType.IsPropertyTypeDeclared(propName)
				}
			}
			var safeLeftType rhino.JSType
			if targetType == nil {
				safeLeftType = ti.unknownType
			} else {
				safeLeftType = targetType
			}
			scope =
				scope.InferQualifiedSlot(
					target, qualifiedName, safeLeftType, resultType, declaredSlotType)
		}

		if !updateNode.IsNil() {
			updateNode.SetJSType(resultType)
		}
		ti.ensurePropertyDefined(target, qualifiedName, isQName, resultType, scope)
	default:
	}
	return scope
}

// EnsurePropertyDefined defines a property if the property has not been
// defined yet.
//
// Java: TypeInference.ensurePropertyDefined(Node, JSType, FlowScope)
func (ti *TypeInference) EnsurePropertyDefined(getprop rhino.Node, rightType rhino.JSType, scope *LinkedFlowScope) {
	qName, isQName := getprop.GetQualifiedNameString()
	ti.ensurePropertyDefined(getprop, qName, isQName, rightType, scope)
}

func (ti *TypeInference) ensurePropertyDefined(getprop rhino.Node, qName string, isQName bool, rightType rhino.JSType, scope *LinkedFlowScope) {
	gp := ast.MustGetProp(getprop)
	propName := gp.Property()
	obj := gp.Object().Node()
	nodeType := ti.GetJSType(obj)
	objectType := rhino.ObjectType_Cast(nodeType.RestrictByNotNullOrUndefined())
	propCreationInConstructor :=
		obj.IsThis() && ti.GetJSType(ti.containerScope.GetRootNode()).IsConstructor()

	if objectType == nil {
		ti.registry.RegisterPropertyOnType(propName, nodeType)
	} else {
		if nodeType.IsStruct() && !objectType.HasProperty(propName) {
			// In general, we don't want to define a property on a struct object,
			// b/c TypeCheck will later check for improper property creation on
			// structs. There are two exceptions.
			// 1) If it's a property created inside the constructor, on the newly
			//    created instance, allow it.
			// 2) If it's a prototype property, allow it. For example:
			//    Foo.prototype.bar = baz;
			//    where Foo.prototype is a struct and the assignment happens at the
			//    top level and the constructor Foo is defined in the same file.
			staticPropCreation := false
			maybeAssignStm := getprop.GetGrandparent()
			if ti.containerScope.IsGlobal() && NodeUtil_IsPrototypePropertyDeclaration(maybeAssignStm) {
				propCreationFilename := maybeAssignStm.GetSourceFileName()
				ctor := objectType.GetOwnerFunction().GetSource()
				if !ctor.IsNil() {
					// ctor.getSourceFileName().equals(propCreationFilename)
					ctorFileName := *ctor.GetSourceFileName()
					if propCreationFilename != nil && ctorFileName == *propCreationFilename {
						staticPropCreation = true
					}
				}
			}
			if !propCreationInConstructor && !staticPropCreation {
				return // Early return to avoid creating the property below.
			}
		}

		if ti.ensurePropertyDeclaredHelper(getprop, qName, isQName, objectType, scope) {
			return
		}

		if !objectType.IsPropertyTypeDeclared(propName) {
			// We do not want a "stray" assign to define an inferred property
			// for every object of this type in the program. So we use a heuristic
			// approach to determine whether to infer the property.
			//
			// 1) If the property is already defined, join it with the previously
			//    inferred type.
			// 2) If this isn't an instance object, define it.
			// 3) If the property of an object is being assigned in the constructor,
			//    define it.
			// 4) If this is a stub, define it.
			// 5) Otherwise, do not define the type, but declare it in the registry
			//    so that we can use it for missing property checks.
			if objectType.HasProperty(propName) || !objectType.IsInstanceType() {
				if "prototype" == propName {
					ti.DefineDeclaredProperty(objectType, propName, rightType, getprop)
				} else {
					objectType.DefineInferredProperty(propName, rightType, getprop)
				}
			} else if propCreationInConstructor {
				objectType.DefineInferredProperty(propName, rightType, getprop)
			} else {
				ti.registry.RegisterPropertyOnType(propName, objectType)
			}
		}
	}
}

// Java: TypeInference.defineDeclaredProperty(ObjectType, String, JSType, Node)
func (ti *TypeInference) DefineDeclaredProperty(
	objectType rhino.ObjectType, propName string, rightType rhino.JSType, getprop rhino.Node) bool {
	if propName == "prototype" && !getprop.GetParent().IsExprResult() {
		functionType := objectType.ToMaybeFunctionType()
		if functionType != nil &&
			!functionType.GetSource().IsNil() &&
			functionType.GetSource().IsClass() {
			ti.compiler.Report(JSError_MakeWithNode(getprop.GetParent(), TypeInference_REASSIGN_CLASS_PROTOTYPE))
			return true
		}
	}
	return objectType.DefineDeclaredProperty(propName, rightType, getprop)
}

// EnsurePropertyDeclaredHelper declares a property on its owner, if
// necessary.
//
// Returns true if a property was declared.
//
// Java: TypeInference.ensurePropertyDeclaredHelper(Node, ObjectType, FlowScope)
func (ti *TypeInference) EnsurePropertyDeclaredHelper(
	getprop rhino.Node, objectType rhino.ObjectType, scope *LinkedFlowScope) bool {
	qName, isQName := getprop.GetQualifiedNameString()
	return ti.ensurePropertyDeclaredHelper(getprop, qName, isQName, objectType, scope)
}

func (ti *TypeInference) ensurePropertyDeclaredHelper(
	getprop rhino.Node, qName string, isQName bool, objectType rhino.ObjectType, scope *LinkedFlowScope) bool {
	if isQName {
		propName := getprop.GetString()
		var_ := TypeInference_GetDeclaredVar(scope, qName)
		if var_ != nil && !var_.IsTypeInferred() {
			// Handle normal declarations that could not be addressed earlier.
			if propName == "prototype" ||
				// Handle prototype declarations that could not be addressed earlier.
				(!objectType.HasOwnProperty(propName) &&
					(!objectType.IsInstanceType() ||
						(var_.IsExtern() && !objectType.IsNativeObjectType()))) {
				return ti.DefineDeclaredProperty(objectType, propName, var_.GetType(), getprop)
			}
		}
	}
	return false
}

// Java: TypeInference.traverseDeclaration(Node, FlowScope)
func (ti *TypeInference) TraverseDeclaration(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	for declarationChild := n.GetFirstChild(); !declarationChild.IsNil(); declarationChild = declarationChild.GetNext() {
		scope = ti.TraverseDeclarationChild(declarationChild, scope)
	}

	return scope
}

// Java: TypeInference.traverseDeclarationChild(Node, FlowScope)
func (ti *TypeInference) TraverseDeclarationChild(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	if n.IsName() {
		return ti.TraverseName(n, scope)
	}

	lhs := ast.MustDestructuringLhs(n)
	init, _ := lhs.Init()
	scope = ti.Traverse(init.Node(), scope)
	return ti.TraverseDestructuringPattern(
		lhs.Pattern().Node(), scope, ti.GetJSType(init.Node()), TypeInference_AssignmentType_DECLARATION)
}

// TraverseDestructuringPattern traverses a destructuring pattern in an
// assignment or declaration.
//
// Java: TypeInference.traverseDestructuringPattern(Node, FlowScope, JSType, AssignmentType)
func (ti *TypeInference) TraverseDestructuringPattern(
	pattern rhino.Node, scope *LinkedFlowScope, patternType rhino.JSType, assignmentType TypeInference_AssignmentType) *LinkedFlowScope {
	return ti.TraverseDestructuringPatternHelper(
		pattern,
		scope,
		patternType,
		func(flowScope *LinkedFlowScope, targetNode rhino.Node, targetType rhino.JSType) *LinkedFlowScope {
			if targetType == nil {
				targetType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			}
			return ti.UpdateScopeForAssignment(flowScope, targetNode, targetType, assignmentType)
		})
}

// TraverseDestructuringPatternHelper traverses a destructuring pattern, and
// calls declarer on each lvalue.
//
// The purpose of the callback is to abstract different logic for declaring
// lvalues in function parameters vs. in regular assignments/declarations.
//
// declarer contains a callback called on every lvalue.
//
// Java: TypeInference.traverseDestructuringPatternHelper(Node, FlowScope, JSType, TypeDeclaringCallback)
func (ti *TypeInference) TraverseDestructuringPatternHelper(
	pattern rhino.Node, scope *LinkedFlowScope, patternType rhino.JSType, declarer TypeInference_TypeDeclaringCallback) *LinkedFlowScope {
	jutil.CheckArgumentMsg1(pattern.IsDestructuringPattern(), "%s", pattern)
	if patternType == nil {
		panic(jutil.NewNullPointerException())
	}
	ForEachDestructuredTarget(ti.registry, patternType, pattern, func(target DestructuredTarget) {
		// The computed property is always evaluated first.
		if target.HasComputedProperty() {
			scope = ti.Traverse(target.GetComputedProperty().GetFirstChild(), scope)
		}
		targetNode := target.GetNode()

		if targetNode.IsDestructuringPattern() {
			if target.HasDefaultValue() {
				ti.Traverse(target.GetDefaultValue(), scope)
			}

			// traverse into nested patterns
			targetType := target.InferType()
			if targetType == nil {
				targetType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			}
			scope = ti.TraverseDestructuringPatternHelper(targetNode, scope, targetType, declarer)
		} else {
			scope = ti.Traverse(targetNode, scope)

			if target.HasDefaultValue() {
				// TODO(lharker): what do we do with the inferred slots in the scope?
				// throw them away or join them with the previous scope?
				ti.Traverse(target.GetDefaultValue(), scope)
			}

			// declare in the scope
			scope = declarer(scope, targetNode, target.InferType())
		}
	})
	// put the `inferred type` of a pattern on it, to make it easier to do typechecking
	pattern.SetJSType(patternType)
	return scope
}

// Java: TypeInference.traverseName(Node, FlowScope)
func (ti *TypeInference) TraverseName(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	name := ast.MustName(n)
	varName := name.Ident()
	type_ := n.GetJSType()
	if value, ok := name.Init(); ok {
		// The only case where `value` isn't null is when we are in a name declaration/initialization
		//     var x = 3;
		scope = ti.Traverse(value.Node(), scope)
		return ti.UpdateScopeForAssignment(scope, n, ti.GetJSType(value.Node()), TypeInference_AssignmentType_DECLARATION)
	}

	if NodeUtil_IsNameDeclaration(n.GetParent()) &&
		TypeInference_IsPossibleMixinApplication(n /* rvalue= */, rhino.NilNode) {
		TypeInference_AddMissingInterfaceProperties(type_)
	}

	if n.GetParent().IsLet() {
		// Whenever we see a LET, we're guaranteed it's not yet in the scope, and we don't need to
		// worry about it being from an outer scope.  In this case, it has no child, so the actual
		// type should be undefined, but we make a special allowance for type-annotated variables.
		// In that case, we use the annotated type instead.
		// TODO(sdh): I would have thought that #updateScopeForTypeChange would handle using the
		// declared type correctly, but for some reason it doesn't so we handle it here.
		var resultType rhino.JSType
		if type_ != nil {
			resultType = type_
		} else {
			resultType = ti.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		}
		scope = ti.UpdateScopeForAssignment(scope, n, resultType, TypeInference_AssignmentType_DECLARATION)
		type_ = resultType
	} else {
		declaredVar, slotVar, flowType, isInferred, ok := scope.LookupSlotType(varName)
		if ok {
			// There are two situations where we don't want to use type information
			// from the scope, even if we have it.

			// 1) The var is escaped and assigned in an inner scope, e.g.,
			// function f() { var x = 3; function g() { x = null } (x); }
			unflowable := isInferred && ti.IsUnflowable(declaredVar)

			// 2) We're reading type information from another scope for an
			// inferred variable. That variable is assigned more than once,
			// and we can't know which type we're getting.
			//
			// var t = null; function f() { (t); } doStuff(); t = {};
			//
			// Notice that this heuristic isn't perfect. For example, you might
			// have:
			//
			// function f() { (t); } f(); var t = 3;
			//
			// In this case, we would infer the first reference to t as
			// type {number}, even though it's undefined.
			var maybeOuterVar *TypedVar
			if isInferred && ti.containerScope.IsLocal() {
				maybeOuterVar = ti.containerScope.GetParent().GetVar(varName)
			}
			// Go-only: language-adapter — Java's var.equals(maybeOuterVar) is
			// TypedVar's identity equality; a null maybeOuterVar is never equal.
			nonLocalInferredSlot :=
				maybeOuterVar != nil && slotVar == maybeOuterVar &&
					!maybeOuterVar.IsMarkedAssignedExactlyOnce()

			if !unflowable && !nonLocalInferredSlot {
				type_ = flowType
				if type_ == nil {
					type_ = ti.unknownType
				}
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// Java: TypeInference.traverseImportSpec(FlowScope, Node)
func (ti *TypeInference) TraverseImportSpec(scope *LinkedFlowScope, spec rhino.Node) *LinkedFlowScope {
	impSpec := ast.MustImportSpec(spec)
	exportedName := impSpec.Imported().Node()
	localName := impSpec.Local().Node()

	scope = ti.Traverse(localName, scope)

	exportedName.SetJSType(localName.GetJSType())
	return scope
}

// Java: TypeInference.redeclareSimpleVar(FlowScope, Node, JSType)
func (ti *TypeInference) RedeclareSimpleVar(scope *LinkedFlowScope, nameNode rhino.Node, varType rhino.JSType) *LinkedFlowScope {
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
	varName := nameNode.GetString()
	return ti.redeclareSimpleVar(scope, nameNode, TypeInference_GetDeclaredVar(scope, varName), varType)
}

func (ti *TypeInference) redeclareSimpleVar(scope *LinkedFlowScope, nameNode rhino.Node, declaredVar *TypedVar, varType rhino.JSType) *LinkedFlowScope {
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
	varName := nameNode.GetString()
	if varType == nil {
		varType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	if ti.IsUnflowable(declaredVar) {
		return scope
	}
	return scope.InferSlotType(varName, varType)
}

// Java: TypeInference.getDeclaredVar(FlowScope, String)
func TypeInference_GetDeclaredVar(scope *LinkedFlowScope, name string) *TypedVar {
	if scope.syntacticScope == nil {
		return nil
	}
	return scope.syntacticScope.GetVar(name)
}

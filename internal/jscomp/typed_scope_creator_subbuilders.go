// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// visitPreorder dispatches pre-order scope-builder actions by scopeBuilderKind.
//
// Java: TypedScopeCreator.{NormalScopeBuilder,FunctionScopeBuilder,ClassScopeBuilder}.visitPreorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) visitPreorder(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch b.kind {
	case scopeBuilderNormal:
		// Create any child block scopes "pre-order" as we see them.
		//
		// This is required because hoisted or qualified names defined in earlier blocks might be
		// referred to later outside the block. This isn't a big deal in most cases since a NamedType
		// will be created and resolved later, but if a NamedType is used for a superclass, we lose a
		// lot of valuable checking. Recursing into child blocks immediately prevents this from being
		// a problem.
		//
		// We don't traverse into CLASSes because we haven't yet have created the class-type on which
		// to assign members. We'll do this on the way back up (post-order) instead, after the
		// class-type has been attached to the AST.
		if !parent.IsNil() && NodeUtil_CreatesBlockScope(n) && !n.IsClass() {
			b.outer.CreateScope(n, b.currentScope)
		}

		// All other functions (and classes, etc) are handled when we see the actual function node.
		if n.IsFunction() {
			b.DefineFunctionLiteral(n)
		}
	case scopeBuilderFunction:
		if parent.IsNil() {
			b.HandleFunctionInputs()
		} else if n.IsFunction() {
			b.DefineFunctionLiteral(n)
		}
	case scopeBuilderClass:
		// These are not descended into, so must be done preorder
		if fn, ok := ast.AsFunction(n); ok {
			if fn.IsEs6Constructor() {
				// Constructor has already been analyzed, so pull that here.
				b.outer.SetDeferredType(n, b.currentScope.GetRootNode().GetJSType())
			} else {
				b.DefineFunctionLiteral(n)
			}
		}
	}
}

// visitPostorder dispatches post-order scope-builder actions by scopeBuilderKind.
//
// Java: TypedScopeCreator.{NormalScopeBuilder,FunctionScopeBuilder,ClassScopeBuilder}.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) visitPostorder(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch b.kind {
	case scopeBuilderNormal:
		b.visitNormalPostorder(t, n, parent)
	case scopeBuilderClass:
		b.visitClassPostorder(t, n, parent)
	}
}

// Java: TypedScopeCreator.NormalScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) visitNormalPostorder(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		b.CheckForClassDefiningCalls(n)
	case rhino.Token_ASSIGN:
		// Handle initialization of properties.
		// We only allow qualified name declarations of the form
		//   /** @type {number} */ a.b.c = rhs;
		// or, for symbol-typed properties:
		//   /** @type {number} */ a.b[Symbol.something] = rhs;
		// TODO(b/77597706): Ensure that CheckJSDoc warns for JSDoc on assignments not to
		// qualified names, e.g.
		//   /** @type {number} */ [a.b.c] = someArr;
		assign := ast.MustAssign(n)
		lhs := assign.Target().Node()
		rhs := assign.Value().Node()
		if lhs.IsGetProp() && lhs.IsQualifiedName() {
			b.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), lhs, n, rhs)
		} else if ge, ok := ast.AsGetElem(lhs); ok && ge.Object().Node().IsQualifiedName() {
			b.MaybeDeclareGetElem(n.GetJSDocInfo(), lhs, rhs)
		} else if _, ok := b.outer.undeclaredNamesForClosure[lhs]; ok {
			b.DefineAssignAsIfVarDeclaration(n)
		}
	case rhino.Token_CATCH:
		b.DefineCatch(n)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		b.DefineVars(n)
	case rhino.Token_GETPROP:
		// Handle stubbed properties.
		if parent.IsExprResult() && n.IsQualifiedName() {
			b.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), n, parent, rhino.NilNode)
		}
	case rhino.Token_GETELEM:
		// Handle stubbed properties.
		if parent.IsExprResult() && ast.MustGetElem(n).Object().Node().IsQualifiedName() {
			b.MaybeDeclareGetElem(n.GetJSDocInfo(), n, rhino.NilNode)
		}
	case rhino.Token_CLASS:
		// Analyse CLASS child-scopes now because later code in this scope may assign
		// properties to these class-types. We want to ensure declarations within the CLASS have
		// priority.
		b.outer.CreateScope(n, b.currentScope)
	case rhino.Token_EXPR_RESULT:
		for _, name := range b.outer.providedNamesFromCall[n] {
			b.DeclareProvidedNs(n, name)
		}
	case rhino.Token_EXPORT:
		exp := ast.MustExport(n)
		if exp.IsDefault() {
			// Define a dummy var for "export default <someExpr>" so that other utilities have
			// access to the type.
			// Problem: according to the debugger, if we have export default SomeName
			// SomeName does not have a type attached so we don't really export it and can't use it
			// as a type later. So we need to call getDeclaredType() instead to get the type.
			declaredType := b.GetDeclaredType(nil, n, exp.Body().Node(), nil)
			b.defineSlot(typedScopeSlotSpec{
				declNode:                n,
				varName:                 Export_DEFAULT_EXPORT_NAME,
				scope:                   b.currentScope,
				type_:                   declaredType,
				allowLaterTypeInference: declaredType == nil,
			})
		}
	default:
	}
}

// HandleFunctionInputs handles bleeding functions and function parameters.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.handleFunctionInputs()
func (b *TypedScopeCreator_scopeBuilder) HandleFunctionInputs() {
	// Handle bleeding functions. These are defined as function expressions which have a non-empty
	// name, which we declare in the FUNCTION scope. Function declarations are hoisted and are
	// already declared in the containing scope; ignore those.
	fnNode := b.currentScope.GetRootNode()
	fn := ast.MustFunction(fnNode)
	if fnNameNode, ok := fn.Name(); ok && NodeUtil_IsFunctionExpression(fnNode) {
		b.defineSlot(typedScopeSlotSpec{
			declNode:                fnNameNode.Node(),
			varName:                 fnNameNode.Ident(),
			scope:                   b.currentScope,
			type_:                   fnNode.GetJSType(),
			allowLaterTypeInference: false,
		})
	}

	b.DeclareParameters(fn)
}

// DeclareParameters declares all of a function's parameters inside the
// function's scope.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareParameters(Node)
func (b *TypedScopeCreator_scopeBuilder) DeclareParameters(fn ast.Function) {
	functionNode := fn.Node()
	if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		// Skip declaring 'exports' for a goog.loadModule(function(exports) {.
		// We pretend that any assignments to 'exports' in the body are actually declarations.
		return
	}

	var iifeArgumentNode rhino.Node

	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}

	functionType := rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	if functionType != nil {
		// Go-only: language-adapter — the Iterator<Parameter> is an index into
		// the parameter list.
		jsdocParameters := functionType.GetParameters()
		jsdocParametersIndex := 0
		var jsDocParameter *rhino.FunctionType_Parameter
		if jsdocParametersIndex < len(jsdocParameters) {
			jsDocParameter = jsdocParameters[jsdocParametersIndex]
			jsdocParametersIndex++
		}

		for astParameter := range fn.Params().Params() {
			if !iifeArgumentNode.IsNil() && iifeArgumentNode.IsSpread() {
				// don't try inferring types from spreads in iifes because we don't know how
				// many items are in the iterable.
				iifeArgumentNode = rhino.NilNode
			}
			var declaredType rhino.JSType
			if jsDocParameter == nil {
				declaredType = b.outer.unknownType
			} else {
				declaredType = jsDocParameter.GetJSType()
			}
			b.DeclareNamesInPositionalParameter(astParameter, declaredType, iifeArgumentNode)
			if jsDocParameter != nil {
				jsDocParameter = nil
				if jsdocParametersIndex < len(jsdocParameters) {
					jsDocParameter = jsdocParameters[jsdocParametersIndex]
					jsdocParametersIndex++
				}
			}
			if !iifeArgumentNode.IsNil() {
				iifeArgumentNode = iifeArgumentNode.GetNext()
			}
		}

		b.DeclareTemplatesInFunctionBody(functionNode, functionType)
	}
}

// Java: TypedScopeCreator.FunctionScopeBuilder.declareTemplatesInFunctionBody(Node, FunctionType)
func (b *TypedScopeCreator_scopeBuilder) DeclareTemplatesInFunctionBody(functionNode rhino.Node, functionType rhino.FunctionType) {
	// Add template params to the scope.
	// This lets JSTypeRegistry resolve references to template types within the function body
	// scope. (they were already registered as types by FunctionTypeBuilder, but they need to be
	// in a TypedScope to let JSTypeRegistry correctly handle scoping).
	info := NodeUtil_GetBestJSDocInfo(functionNode)
	var templateNames *jutil.LinkedHashSet[string]

	isPossiblePrototypeMethod :=
		!functionType.HasInstanceType() && (info == nil || !info.HasThisType())
	if isPossiblePrototypeMethod {
		// If this is a prototype method, add any @template types from the prototype owner.
		// e.g. given
		//   /** @template T */
		//   class Foo {
		//     bar() {}
		//   }
		// We want to declare 'T' in the scope of 'bar'
		ownerNode := NodeUtil_GetBestLValue(functionNode)
		if !ownerNode.IsNil() {
			ownerKeys := b.FindOwnerTypeKeys(ownerNode)
			if len(ownerKeys) > 0 {
				templateNames = jutil.NewLinkedHashSet[string]()
				for _, templateKey := range ownerKeys {
					templateNames.Add(*templateKey.GetReferenceName())
				}
			}
		}
	}
	if info != nil {
		// Add @template parameters from the JSDoc on this function directly. For example, given
		//   /** @template U */
		//   function foo() {}
		// This declares 'U'.
		infoTemplateNames := info.GetTemplateTypeEntries()
		if len(infoTemplateNames) > 0 {
			if templateNames == nil {
				templateNames = jutil.NewLinkedHashSet[string]()
			}
			for _, entry := range infoTemplateNames {
				templateNames.Add(entry.Name)
			}
		}
		typeTransformations := jsDocTypeTransformations(info)
		if typeTransformations != nil && !typeTransformations.IsEmpty() {
			if templateNames == nil {
				templateNames = jutil.NewLinkedHashSet[string]()
			}
			for name := range typeTransformations.Keys() {
				templateNames.Add(name)
			}
		}
	}
	if templateNames == nil || templateNames.IsEmpty() {
		return
	}
	input := b.GetCompilerInput()
	voidType := b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	// Declare any template names in the function scope. This means that if someone shadows
	// an outer variable FOO with a @template FOO and refers to FOO inside the method, we
	// will treat it as undefined, rather than the correct type, which could lead to weird
	// errors. Ideally we'd have a "don't use me" type that gives an error at use.
	for name := range templateNames.All() {
		if !b.currentScope.CanDeclare(name) {
			b.outer.validator.ExpectUndeclaredVariable(
				typedScopeCreator_sourceNameArg(NodeUtil_GetSourceName(functionNode)),
				input,
				functionNode,
				functionNode.GetParent(),
				b.currentScope.GetVar(name),
				name,
				voidType)
		}
		b.currentScope.Declare(name, functionNode, voidType, input /* inferred= */, false)
	}
}

// DeclareNamesInPositionalParameter declares the name(s) in a positional AST
// parameter in the scope.
//
// astParameter is the positional parameter node; declaredParameterType the
// declared parameter type, or the unknown type if there is none;
// iifeArgumentNode the corresponding argument from the iife, if in an iife.
// e.g. for `(function (x) {}(3);` this would be `3`.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareNamesInPositionalParameter(Node, JSType, Node)
func (b *TypedScopeCreator_scopeBuilder) DeclareNamesInPositionalParameter(
	astParameter ast.Param, declaredParameterType rhino.JSType, iifeArgumentNode rhino.Node) {
	paramType := declaredParameterType
	isInferred := paramType.Equals(b.outer.unknownType)

	if !iifeArgumentNode.IsNil() && isInferred {
		var argumentVar *TypedVar
		if b.currentScope.GetParent() != nil {
			if argumentName, ok := iifeArgumentNode.GetQualifiedNameString(); ok {
				argumentVar = b.currentScope.GetParent().GetVar(argumentName)
			}
		}
		if argumentVar != nil && !argumentVar.IsTypeInferred() {
			paramType = argumentVar.GetType()
		}
	}

	if paramType == nil {
		paramType = b.outer.unknownType
	}

	switch astParameter.Token() {
	case rhino.Token_NAME:
		// function f(x) {}
		b.DeclareSingleParameterName(isInferred, astParameter.Node(), paramType)
	case rhino.Token_ITER_REST:
		// function f(...x) {}
		// rest parameter is actually an array of the type specified in the JSDoc
		target := astParameter.MustIterRest().Target()
		arrayType := b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
		restParamType := b.outer.typeRegistry.CreateTemplatizedType(arrayType, paramType)
		if target.Token() == rhino.Token_NAME {
			b.DeclareSingleParameterName(isInferred, target.Node(), restParamType)
		} else {
			// function f(...{length}) {}
			b.DeclareDestructuringParameter(isInferred, target.Node(), restParamType)
		}
	case rhino.Token_DEFAULT_VALUE:
		// function f(x = 3) {} or function f([x] = []) {}
		actualParam := astParameter.MustDefaultValue().Target()
		if actualParam.Token() == rhino.Token_NAME {
			b.DeclareSingleParameterName(isInferred, actualParam.Node(), paramType)
		} else {
			b.DeclareDestructuringParameter(isInferred, actualParam.Node(), paramType)
		}
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		// function f([x]) {}
		// function f({x}) {}
		b.DeclareDestructuringParameter(isInferred, astParameter.Node(), paramType)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected function parameter node " + nodeStringOrNull(astParameter.Node())))
	}
}

// DeclareDestructuringParameter declares all names inside a destructuring
// pattern in a parameter list in the scope if we can find a non-unknown type
// for them.
//
// Unknown typed parameters are always treated as inferred, not declared.
// TypeInference may later give them a better inferred type than unknown,
// but they will never become declared.
//
// NOTE: currently, there are some less-than-ideal aspects to how we do
// this. If the pattern type is an unresolved NamedType, then we can't
// lookup properties on it (to find the individual parameter types) until
// after name resolution. The current state is to defer to TypeInference to
// type those parameters, with the drawback that they are 'inferred', not
// 'declared', and so any type can be assigned to them. In the future we
// will just enforce typing each parameter individually: relevant issue
// https://github.com/google/closure-compiler/issues/1781
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareDestructuringParameter(boolean, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) DeclareDestructuringParameter(
	isInferred bool, pattern rhino.Node, patternType rhino.JSType) {
	ForEachDestructuredTarget(b.outer.typeRegistry, patternType, pattern, func(target DestructuredTarget) {
		parameterType := target.InferTypeWithoutUsingDefaultValue()

		if target.GetNode().IsDestructuringPattern() {
			b.DeclareDestructuringParameter(isInferred, target.GetNode(), parameterType)
		} else {
			paramName := target.GetNode()
			jutil.CheckStateMsg1(paramName.IsName(), "Expected all parameters to be names, got %s", paramName)

			if parameterType == nil || parameterType.IsUnknownType() {
				paramJSDoc := paramName.GetJSDocInfo()
				if paramJSDoc != nil && paramJSDoc.HasType() {
					// see if the parameter has its own inline JSDoc, and use that unless we already have
					// a type from @param JSDoc.
					// TODO(b/112651122): this should happen inside FunctionTypeBuilder, so that we can
					// check that calls to the function match the inline JSDoc.
					// TODO(b/111523967): we should also report a
					// warning if the inline and non-inline JSDoc conflict.
					parameterType =
						b.outer.typeRegistry.EvaluateTypeExpression(paramJSDoc.GetType(), b.currentScope)
					isInferred = false
				} else {
					// note - these parameters may get better types during TypeInference
					isInferred = true
				}
			}
			b.DeclareSingleParameterName(isInferred, paramName, parameterType)
		}
	})
}

// Java: TypedScopeCreator.FunctionScopeBuilder.declareSingleParameterName(boolean, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) DeclareSingleParameterName(isInferred bool, name rhino.Node, type_ rhino.JSType) {
	b.defineSlot(typedScopeSlotSpec{
		declNode:                name,
		varName:                 name.GetString(),
		scope:                   b.currentScope,
		type_:                   type_,
		allowLaterTypeInference: isInferred,
	})
} // end FunctionScopeBuilder

// Java: TypedScopeCreator.ClassScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) visitClassPostorder(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if name, ok := ast.AsName(n); ok &&
		parent == b.currentScope.GetRootNode() &&
		NodeUtil_IsClassExpression(parent) {
		// Declare bleeding class name in scope.  Pull the type off the AST.
		jutil.CheckState(name.Ident() != "") // anonymous classes have EMPTY nodes, not NAME
		b.defineSlot(typedScopeSlotSpec{
			declNode:                n,
			varName:                 slotVarNameFromDeclNode(n),
			scope:                   b.currentScope,
			type_:                   parent.GetJSType(),
			allowLaterTypeInference: false,
		})
	} else if mf, ok := ast.AsMemberFunctionDef(n); ok && parent.IsClassMembers() {
		// Ignore "constructor" since it has special handling in `createClassTypeFromNodes()`.
		if !mf.IsEs6Constructor() {
			b.DefineMemberFunction(mf)
		}
	} else if f, ok := ast.AsMemberFieldDef(n); ok {
		// public fields are roots of their own scope so the parent doesn't get passed into
		// visitPostorder
		b.DefineMemberField(f)
	} else if gs, ok := ast.AsGetterOrSetterDef(n); ok && parent.IsClassMembers() {
		b.DefineGetterSetter(gs)
	} else if n.IsComputedFieldDef() {
		b.DefineComputedMemberField(ast.MustComputedMember(n))
	} else if cp, ok := ast.AsComputedProp(n); ok && cp.IsMethod() && parent.IsClassMembers() {
		b.DefineComputedMemberField(cp.AsComputedMember())
	}
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineMemberFunction(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineMemberFunction(mf ast.MemberFunctionDef) {
	n := mf.Node()
	ownerType := b.DetermineOwnerTypeForClassMember(n)
	ownerType.DefineDeclaredProperty(mf.Key(), mf.Function().Node().GetJSType(), n)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineMemberField(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineMemberField(f ast.MemberFieldDef) {
	n := f.Node()
	var rhs rhino.Node
	if val, ok := f.Value(); ok {
		rhs = val.Node()
	}
	declaredType := b.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	// When there's no JSDoc type declaration for a field:
	// If we can infer a type other than UNKNOWN for an initial value, use that as the field type
	// Otherwise, set the field's type to the ALL type ('*')
	if declaredType == nil {
		var rhsType rhino.JSType
		if !rhs.IsNil() {
			rhsType = b.GetDeclaredRValueType(rhino.NilNode, rhs)
		}
		if rhsType != nil && !rhsType.IsUnknownType() {
			declaredType = rhsType
		} else {
			declaredType = b.outer.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
		}
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	ownerType.DefineDeclaredProperty(f.Key(), declaredType, n)
	n.SetJSType(declaredType)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineComputedMemberField(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineComputedMemberField(cm ast.ComputedMember) {
	n := cm.Node()
	key := cm.Key().Node()
	if !key.IsQualifiedName() {
		return
	}
	keyType := b.currentScope.LookupQualifiedName(key.GetQualifiedNameObject())
	if keyType == nil || !keyType.IsKnownSymbolValueType() {
		return
	}
	var rhs rhino.Node
	if val, ok := cm.Value(); ok {
		rhs = val.Node()
	}
	declaredType := b.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	// When there's no JSDoc type declaration for a field:
	// If we can infer a type other than UNKNOWN for an initial value, use that as the field type
	// Otherwise, set the field's type to the ALL type ('*')
	if declaredType == nil {
		var rhsType rhino.JSType
		if !rhs.IsNil() {
			rhsType = b.GetDeclaredRValueType(rhino.NilNode, rhs)
		}
		if rhsType != nil && !rhsType.IsUnknownType() {
			declaredType = rhsType
		} else {
			declaredType = b.outer.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
		}
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	ownerType.DefineDeclaredPropertyByKey(
		rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), declaredType, n)
	n.SetJSType(declaredType)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineGetterSetter(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineGetterSetter(gs ast.GetterOrSetterDef) {
	n := gs.Node()
	name := gs.Key()
	methodType := gs.Function().Node().GetJSType().ToMaybeFunctionType()

	var propertyType rhino.JSType
	switch gs.Token() {
	case rhino.Token_GETTER_DEF:
		// TODO(sdh): consider only falling back on unknown if the function body is empty?
		// But
		// we need to not report a conflicting type error if there's different unknowns.
		if methodType.IsReturnTypeInferred() {
			propertyType = b.outer.unknownType
		} else {
			propertyType = methodType.GetReturnType()
		}
	case rhino.Token_SETTER_DEF:
		if len(methodType.GetParameters()) == 0 {
			propertyType = b.outer.unknownType
		} else {
			propertyType = methodType.GetParameters()[0].GetJSType()
		}
	default:
		panic(jutil.NewAssertionError(n.ToStringTree()))
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	// TODO(b/116797078): correctly model getters/setters and stop treating this as a normal
	// property.
	ownerType.DefineDeclaredProperty(name, propertyType, n)
}

// DetermineOwnerTypeForClassMember returns the owner type for a class member
// function, getter, or setter.
//
// For a member on class C, this is either `C` for a static member or
// `C.prototype` for a nonstatic member.
//
// Java: TypedScopeCreator.ClassScopeBuilder.determineOwnerTypeForClassMember(Node)
func (b *TypedScopeCreator_scopeBuilder) DetermineOwnerTypeForClassMember(member rhino.Node) rhino.ObjectType {
	// MEMBER_FUNCTION_DEF -> CLASS_MEMBERS -> CLASS or
	// MEMBER_FIELD_DEF -> CLASS_MEMBERS -> CLASS or
	// GETTER_DEF -> CLASS_MEMBERS -> CLASS or
	// SETTER_DEF -> CLASS_MEMBERS -> CLASS
	ownerNode := member.GetGrandparent()
	jutil.CheckState(ownerNode.IsClass())
	ownerType := ownerNode.GetJSType().ToMaybeFunctionType()
	if member.IsStaticMember() {
		return ownerType
	} else if member.IsMemberFieldDef() {
		return ownerType.GetInstanceType()
	} else {
		jutil.CheckStateMsg1(
			member.IsMemberFunctionDef() ||
				member.IsGetterDef() ||
				member.IsSetterDef() ||
				member.IsComputedProp() ||
				member.IsComputedFieldDef(),
			"%s",
			member)
		return ownerType.GetPrototype()
	}
} // end ClassScopeBuilder

// TypedScopeCreator_FirstOrderFunctionAnalyzer does a first-order function
// analysis that just looks at simple things like what variables are
// escaped, and whether 'this' is used.
//
// The syntactic scopes created in this traversal are also stored for later
// use.
type TypedScopeCreator_FirstOrderFunctionAnalyzer struct {
	NodeTraversal_AbstractScopedCallbackBase

	// Go-only: language-adapter — the enclosing TypedScopeCreator instance of
	// this inner class.
	outer                    *TypedScopeCreator
	identifyEnumsAndTypedefs TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable
	localVarState            map[*Var]uint8
	localContainerDepth      int
	globalDeclEstimate       int
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer()
func NewTypedScopeCreator_FirstOrderFunctionAnalyzer(outer *TypedScopeCreator) *TypedScopeCreator_FirstOrderFunctionAnalyzer {
	return &TypedScopeCreator_FirstOrderFunctionAnalyzer{outer: outer}
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer.enterScope(NodeTraversal)
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) EnterScope(t *NodeTraversal) {
	a.outer.scopeCount++
	if t.InGlobalScope() {
		if t.scopes[0].untypedScope == nil {
			gs := Scope_CreateGlobalScope(t.GetScopeRoot())
			gs.scanDone = true
			t.scopes[0].scope = gs
			t.scopes[0].untypedScope = gs
			t.scopes[0].lookupScope = gs
		}
		return
	}
	root := t.GetScopeRoot()
	if root.IsFunction() || root.IsModuleBody() || NodeUtil_IsClassStaticBlock(root) {
		a.localContainerDepth++
	}
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer.exitScope(NodeTraversal)
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) ExitScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		return
	}
	root := t.GetScopeRoot()
	if root.IsFunction() || root.IsModuleBody() || NodeUtil_IsClassStaticBlock(root) {
		a.localContainerDepth--
	}
	scope := t.scopes[len(t.scopes)-1].untypedScope
	if scope == nil {
		if fn, ok := ast.AsFunction(root); !root.IsModuleBody() && (!ok || fn.NameNode().Ident() == "") {
			return
		}
		scope = t.GetScope()
	}
	if scope.implicitVars != nil {
		if ev := scope.implicitVars[AbstractScope_ImplicitVar_EXPORTS]; ev != nil {
			h := a.localVarState[ev]
			delete(a.localVarState, ev)
			var flags uint8
			if (h & localVarEscapedBit) != 0 {
				flags |= typedScopeCreator_VarFlagEscaped
			}
			if (h & localVarAssignCountMask) != localVarAssignedOnce {
				flags |= typedScopeCreator_VarFlagNotSingleAssigned
			}
			if flags != 0 {
				a.outer.varFlags[typedScopeCreator_ScopedNameKey{name: "exports", scopeRoot: root}] |= flags
			}
			if (h & localVarAssignCountMask) == localVarAssignedOnce {
				list := a.outer.reservedNamesForScope[root]
				if list == nil {
					list = make([]string, 0, scope.GetVarCount()+1)
				}
				list = append(list, "exports")
				a.outer.reservedNamesForScope[root] = list
			}
		}
	}
	varCount := scope.GetVarCount()
	if varCount == 0 {
		return
	}
	list := a.outer.reservedNamesForScope[root]
	if list == nil {
		list = make([]string, 0, varCount)
	}
	for _, symbol := range scope.varList {
		if symbol == nil {
			continue
		}
		name := symbol.GetName()
		h := a.localVarState[symbol]
		delete(a.localVarState, symbol)
		if h != localVarAssignedOnce {
			var flags uint8
			if (h & localVarEscapedBit) != 0 {
				flags |= typedScopeCreator_VarFlagEscaped
			}
			if (h & localVarAssignCountMask) != localVarAssignedOnce {
				flags |= typedScopeCreator_VarFlagNotSingleAssigned
			}
			a.outer.varFlags[typedScopeCreator_ScopedNameKey{name: name, scopeRoot: root}] |= flags
		}
		if !slices.Contains(list, name) {
			list = append(list, name)
		}
	}
	a.outer.reservedNamesForScope[root] = list
}

const (
	localVarAssignCountMask  uint8 = 3
	localVarAssignedOnce     uint8 = 1
	localVarAssignedMultiple uint8 = 2
	localVarEscapedBit       uint8 = 4
)

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer.visit(NodeTraversal, Node, Node)
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if a.identifyEnumsAndTypedefs.registry != nil {
		a.identifyEnumsAndTypedefs.Visit(t, n, parent)
	}
	if t.InGlobalScope() {
		switch n.GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_ASSIGN:
			a.globalDeclEstimate++
		}
		// The first-order function analyzer looks at two types of variables:
		//
		// 1) Local variables that are assigned in inner scopes ("escaped vars")
		//
		// 2) Local variables that are assigned more than once.
		//
		// We treat all global variables as escaped by default, so there's
		// no reason to do this extra computation for them.
		return
	}

	// Record function with returns or arrow functions without bodies
	if (n.IsReturn() && n.HasChildren()) || NodeUtil_IsBlocklessArrowFunctionResult(n) {
		for i := len(t.scopes) - 1; i >= 0; i-- {
			rootNode := t.scopes[i].node
			if !NodeUtil_CreatesBlockScope(rootNode) {
				a.outer.functionsWithNonEmptyReturns[rootNode] = struct{}{}
				break
			}
		}
	}

	// Be careful of bleeding functions, which create variables
	// in the inner scope, not the scope where the name appears.
	if n.IsName() && NodeUtil_IsLValue(n) && !NodeUtil_IsBleedingFunctionName(n) {
		name := n.GetString()
		scope := t.resolveLookupScope(len(t.scopes) - 1)
		var_ := scope.GetVar(name)
		// TODO(sdh): consider checking hasSameHoistScope instead of container scope here and
		// below. This will detect function(a) { a.foo = bar } as an escaped qualified name,
		// which seems like the right thing to do (but could possibly break things?)
		// Doing so will allow removing the warning on TypeCheckTest#testIssue1024b.
		if var_ != nil {
			ownerScope := var_.scope
			if ownerScope.IsLocal() {
				if a.localVarState == nil {
					a.localVarState = make(map[*Var]uint8, 64)
				}
				h := a.localVarState[var_]
				if (h & localVarAssignCountMask) < localVarAssignedMultiple {
					h++
				}
				if a.localContainerDepth > 1 {
					containerScope := scope.GetClosestContainerScope()
					if !containerScope.HasSameContainerScope(ownerScope) {
						h |= localVarEscapedBit
					}
				}
				a.localVarState[var_] = h
			}
		}
	} else if n.IsGetProp() && n.IsUnscopedQualifiedName() && NodeUtil_IsLValue(n) {
		scope := t.resolveLookupScope(len(t.scopes) - 1)
		if a.localContainerDepth > 1 {
			name := NodeUtil_GetRootOfQualifiedName(n).GetString()
			var_ := scope.GetVar(name)
			if var_ != nil {
				ownerScope := var_.scope
				if ownerScope.IsLocal() {
					containerScope := scope.GetClosestContainerScope()
					if !containerScope.HasSameContainerScope(ownerScope) {
						qName, _ := n.GetQualifiedNameString()
						a.outer.varFlags[typedScopeCreator_ScopedNameKey{name: qName, scopeRoot: ownerScope.GetRootNode()}] |= typedScopeCreator_VarFlagEscaped
					}
				}
			}
		}
	}
}

// Java: TypedScopeCreator.containingGoogModuleIdOf(TypedScope)
func TypedScopeCreator_ContainingGoogModuleIdOf(scope *TypedScope) *string {
	for scope != nil {
		module := scope.GetModule()
		if module == nil {
			scope = scope.GetParent()
			continue
		}

		// Stop recursing once we've hit a module scope.
		metadata := module.Metadata()
		jutil.CheckStateMsg1(metadata.IsModule(), "%s", metadata)

		/*
		 * This module may not have a goog.module/goog.declareModuleId. Also don't crash if it's
		 * malformed with multiple module IDs.
		 */
		if ns := module.ClosureNamespace(); ns != nil {
			return ns
		}
		for namespace := range metadata.GoogNamespaces().All() {
			return strpool.StringPtr(strpool.Intern(namespace))
		}
		return nil
	}
	return nil
}

// Java: TypedScopeCreator.extractKnownSymbolKey(TypedScope, Node)
func TypedScopeCreator_ExtractKnownSymbolKey(currentScope *TypedScope, key rhino.Node) *rhino.KnownSymbolType {
	if !key.IsQualifiedName() {
		return nil
	}
	type_ := currentScope.LookupQualifiedName(key.GetQualifiedNameObject())
	if type_ != nil {
		return type_.ToMaybeKnownSymbolType()
	}
	return nil
}

var (
	_ ScopeCreator                                  = (*TypedScopeCreator)(nil)
	_ rhino.StaticSymbolTable[*TypedVar, *TypedVar] = (*TypedScopeCreator)(nil)
	_ NodeTraversal_Callback                        = (*TypedScopeCreator_scopeBuilder)(nil)
	_ NodeTraversal_Callback                        = (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable)(nil)
	_ NodeTraversal_ScopedCallback                  = (*TypedScopeCreator_FirstOrderFunctionAnalyzer)(nil)
)

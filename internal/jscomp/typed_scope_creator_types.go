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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ShouldUseFunctionLiteralType: if a variable is assigned a function
// literal in the global scope, make that a declared type (even if there's
// no doc info). There's only one exception to this rule: if the return type
// is inferred, and we're in a local scope, we should assume the whole
// function is inferred.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.shouldUseFunctionLiteralType(FunctionType, JSDocInfo, Node)
func (b *TypedScopeCreator_scopeBuilder) ShouldUseFunctionLiteralType(type_ rhino.FunctionType, info *rhino.JSDocInfo, lValue rhino.Node) bool {
	if info != nil {
		return true
	}
	if !lValue.IsNil() && NodeUtil_MayBeObjectLitKey(lValue) {
		return false
	}
	// TODO(johnlenz): consider unifying global and local behavior
	return b.IsLValueRootedInGlobalScope(lValue) || !type_.IsReturnTypeInferred()
}

// CreateClassTypeFromNodes creates a new class type from the given class
// literal. This function does not need to worry about stubs and aliases
// because they are handled by createFunctionTypeFromNodes instead.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createClassTypeFromNodes(Node, String, String, JSDocInfo, Node)
func (b *TypedScopeCreator_scopeBuilder) CreateClassTypeFromNodes(
	clazz rhino.Node,
	name *string,
	syntacticName *string,
	info *rhino.JSDocInfo,
	lvalueNode rhino.Node) rhino.FunctionType {
	cls := ast.MustClass(clazz)

	declarationScope := b.currentScope
	if !lvalueNode.IsNil() {
		declarationScope = b.GetLValueRootScope(lvalueNode)
	}
	builder := NewFunctionTypeBuilder(name, b.outer.compiler, clazz, b.currentScope).
		UsingClassSyntax().
		SetSyntacticFunctionName(syntacticName).
		SetAstContents(clazz, false).
		SetDeclarationScope(declarationScope).
		InferKind(info).
		InferTemplateTypeName(info, nil)

	extendsClause := cls.ExtendsNode()

	// Look at the extends clause and/or JSDoc info to find a super class.  Use generics from the
	// JSDoc to supplement the extends type when available.
	baseType := b.FindSuperClassFromNodes(extendsClause, info)
	builder.InferInheritance(info, baseType)

	// Look for an explicit constructor.
	var constructor rhino.Node
	if ctorDef, ok := cls.Constructor(); ok {
		ctorFn := ctorDef.Function()
		constructor = ctorFn.Node()
		ctorInfo := NodeUtil_GetBestJSDocInfo(constructor)
		builder.InferConstructorParameters(ctorFn.Params().Node(), ctorInfo)
	} else if extendsClause.IsEmpty() {
		// No explicit constructor and no superclass: constructor is no-args.
		builder.InferImplicitConstructorParameters(emptyFunctionParametersSlice)
	} else {
		// No explicit constructor, but we have a superclass.  If we know its type, then copy its
		// constructor arguments (and templates).  If not, make the constructor arguments unknown.
		// TODO(sdh): consider allowing attaching constructor @param annotations somewhere else?
		var extendsCtor rhino.FunctionType
		if baseType != nil {
			extendsCtor = baseType.GetConstructor()
		}
		if extendsCtor != nil {
			// Known superclass: copy the parameters node.
			builder.InferImplicitConstructorParameters(extendsCtor.GetParameters())
		} else {
			// Unresolveable extends clause: suppress typechecking.
			builder.InferImplicitConstructorParameters(
				b.outer.typeRegistry.CreateParametersWithVarArgs(
					b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)))
		}
	}

	// TODO(sdh): Handle template parameters.  The constructor should store all parameters,
	// while the instance type should only have the class-level parameters?

	// Add the type for the "constructor" property.
	classType := builder.BuildAndRegister()
	if classType.IsConstructor() {
		// NOTE: This logic is similar to the goog.inherits handling in
		// ClosureCodingConvention#applySubclassRelationship. If this logic is modified
		// it is likely that code needs to be modified as well.

		// Notice that constructor functions do not need to be covariant on the superclass.
		// So if G extends F, new G() and new F() can accept completely different argument
		// types, but G.prototype.constructor needs to be covariant on F.prototype.constructor.
		// To get around this, we just turn off type-checking on arguments and return types
		// of G.prototype.constructor.

		// NOTE: For final classes we could do better here and retain the parameter types.

		qmarkCtor := classType.ForgetParameterAndReturnTypes()
		classPrototypeType := classType.GetPrototypeProperty()
		classPrototypeType.DefineDeclaredProperty("constructor", qmarkCtor, constructor)
	}
	if classType.HasInstanceType() {
		classPrototype := classType.GetSlot("prototype")
		// SymbolTable users expect the class prototype and actual class to have the same
		// declaration node.
		if !lvalueNode.IsNil() {
			classPrototype.SetNode(lvalueNode)
		} else {
			classPrototype.SetNode(classPrototype.GetNode())
		}
	}

	return classType
}

// FindSuperClassFromNodes looks at the extends clause to find the instance
// type being extended. Returns null if there is no such clause, and unknown
// if the type cannot be determined.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.findSuperClassFromNodes(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) FindSuperClassFromNodes(extendsNode rhino.Node, info *rhino.JSDocInfo) rhino.ObjectType {
	if extendsNode.IsEmpty() {
		// No extends clause: return null.
		return nil
	}
	ctorType := extendsNode.GetJSType()
	if ctorType == nil {
		if extendsNode.IsQualifiedName() {
			superclass, _ := extendsNode.GetQualifiedNameString()
			// Look up qualified names in the scope (types won't be set on the AST until inference).
			ctorType = rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, superclass)
			if ctorType == nil {
				var_ := b.currentScope.GetVar(superclass)
				if var_ != nil {
					ctorType = var_.GetType()
				} else {
					ctorType = nil
				}
			}
			// If that doesn't work, then fall back on the registry
			if ctorType == nil {
				return rhino.ObjectType_Cast(
					b.outer.typeRegistry.GetTypeWithLocation(
						b.currentScope,
						superclass,
						extendsNode.GetSourceFileName(),
						extendsNode.GetLineno(),
						extendsNode.GetCharno()))
			}
		} else if b.IsGoogModuleGetProperty(extendsNode) {
			ctorType = b.GetCtorForGoogModuleGet(extendsNode)
		} else {
			// Anything TypedScopeCreator can infer has already been read off the AST.  This is likely
			// a CALL or GETELEM, which are unknown until TypeInference.  Instead, ignore it for now,
			// require an @extends tag in the JSDoc, and verify correctness in TypeCheck.
			if info == nil || !info.HasBaseType() {
				b.outer.Report(JSError_MakeWithNode(extendsNode, TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC))
			}
		}
	}

	if ctorType != nil {
		if ctorType.IsConstructor() || ctorType.IsInterface() {
			return ctorType.ToMaybeFunctionType().GetInstanceType()
		} else if ctorType.IsUnknownType() || ctorType.IsNoResolvedType() {
			// The constructor could have an unknown type for cases where it is dynamically
			// created or passed in from elsewhere.
			// e.g. with a mixin pattern
			// function mixinSomething(ctor) {
			//   return class extends ctor { ... };
			// }
			// In that case consider the super class instance type to be unknown.
			// On the other hand, the constructor can be a NoResolvedType if running in Clutz's
			// partial compilation mode.
			return ctorType.ToMaybeObjectType()
		}
	}

	// We couldn't determine the type, so for TypedScope creation purposes we will treat it as if
	// there were no extends clause.  TypeCheck will give a more precise error later.
	return nil
}

// Java: TypedScopeCreator.AbstractScopeBuilder.isGoogModuleGetProperty(Node)
func (b *TypedScopeCreator_scopeBuilder) IsGoogModuleGetProperty(extendsClause rhino.Node) bool {
	if extendsClause.IsCall() {
		return ModuleImportResolver_IsGoogModuleDependencyCall(extendsClause)
	} else if gp, ok := ast.AsGetProp(extendsClause); ok {
		return b.IsGoogModuleGetProperty(gp.Object().Node())
	} else {
		return false
	}
}

// GetCtorForGoogModuleGet looks up the type of a goog.module.get property
// access.
//
// extendsNode is `goog.module.get('x');` or a getprop
// `goog.module.get('x').y.z`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getCtorForGoogModuleGet(Node)
func (b *TypedScopeCreator_scopeBuilder) GetCtorForGoogModuleGet(extendsNode rhino.Node) rhino.JSType {
	call := extendsNode
	properties := []string{}
	for !call.IsCall() {
		gp := ast.MustGetProp(call)
		properties = slices.Insert(properties, 0, gp.Property())
		call = gp.Object().Node()
	}
	extendsCall := b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(call)
	if extendsCall == nil || !b.outer.memoized.ContainsKey(extendsCall.GetScopeRoot()) {
		// Handle invalid goog.module.get calls
		return nil
	}
	moduleScope := b.outer.memoized.Get(extendsCall.GetScopeRoot())
	properties = slices.Insert(properties, 0, extendsCall.GetName())
	return rhino.StaticTypedScope_LookupQualifiedNameString(moduleScope, strings.Join(properties, "."))
}

// CreateFunctionTypeFromNodes creates a new function type, based on the
// given nodes.
//
// This handles two cases that are semantically very different, but are not
// mutually exclusive: - A function literal that needs a type attached to it
// (called from defineClassLiteral with a non-null FUNCTION node for rValue).
// - An assignment expression with function-type info in the JsDoc (called
// from getDeclaredType on a stub (rValue == null) or alias (rValue is a
// qualified name).
//
// All parameters are optional, and we will do the best we can to create a
// function type.
//
// This function will always create a function type, so only call it if
// you're sure that's what you want.
//
// rValue is the function node; name the function's name; info the JSDocInfo
// attached to the function definition; lvalueNode the node where this
// function is being assigned. For example, `A.prototype.foo = ...` would be
// used to determine that this function is a method of A.prototype. May be
// null to indicate that this is not being assigned to a qualified name.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createFunctionTypeFromNodes(Node, String, JSDocInfo, Node)
func (b *TypedScopeCreator_scopeBuilder) CreateFunctionTypeFromNodes(
	rValue rhino.Node,
	name *string,
	info *rhino.JSDocInfo,
	lvalueNode rhino.Node) rhino.FunctionType {
	// Check for an alias.
	if !rValue.IsNil() && !lvalueNode.IsNil() {
		if rvalQname, ok := rValue.GetQualifiedNameString(); ok {
			var_ := b.currentScope.GetVar(rvalQname)
			if var_ != nil && var_.GetType() != nil && var_.GetType().IsFunctionType() {
				aliasedType := var_.GetType().ToMaybeFunctionType()
				if aliasedType.IsConstructor() || aliasedType.IsInterface() {
					// TODO(nick): Remove this. This should already be handled by normal type resolution.
					if name != nil {
						b.outer.typeRegistry.DeclareType(b.currentScope, *name, aliasedType.GetInstanceType())
					}
					b.CheckFunctionAliasAnnotations(lvalueNode, aliasedType, info)
					return aliasedType
				}
			}
		}
	}

	// No alias: look for an explicit @type in JSDocInfo.
	if info != nil && info.HasType() {
		type_ := info.GetType().Evaluate(b.currentScope, b.outer.typeRegistry)

		// Known to be not null since we have the FUNCTION token there.
		type_ = type_.RestrictByNotNullOrUndefined()
		if type_.IsFunctionType() {
			functionType := type_.ToMaybeFunctionType()
			functionType.SetJSDocInfo(info)
			return functionType
		}
	}

	// No alias or explicit @type, so look for a function literal, or @param/@return.
	errorRoot := rValue
	if rValue.IsNil() {
		errorRoot = lvalueNode
	}
	var fnRoot rhino.Node
	var parametersNode rhino.Node
	if fn, ok := ast.AsFunction(rValue); ok {
		fnRoot = rValue
		parametersNode = fn.Params().Node()
	}

	// If this function is being assigned as a property on a type, try finding the owner type
	// and the property name.
	// This is easy to do for class members because the owner type is on the CLASS node and
	// the property name is the MEMBER_FUNCTION_DEF/GETTER_DEF/SETTER_DEF string.
	// For other functions, we rely on NodeUtil.getBestLValueOwner.
	var classRoot rhino.Node
	if !lvalueNode.IsNil() && lvalueNode.GetParent().IsClassMembers() {
		classRoot = lvalueNode.GetGrandparent()
	}
	ownerNode := NodeUtil_GetBestLValueOwner(lvalueNode)

	var ownerType rhino.ObjectType
	var propName string
	var hasPropName bool
	if !classRoot.IsNil() {
		// Static members are owned by the constructor, non-statics are owned by the prototype.
		if ownerFn := rhino.JSType_ToMaybeFunctionType(classRoot.GetJSType()); ownerFn != nil {
			ownerType = ownerFn
		}
		if !lvalueNode.IsStaticMember() && ownerType != nil {
			ownerType = ownerType.(rhino.FunctionType).GetPrototype()
		}
		if !lvalueNode.IsComputedProp() {
			propName = lvalueNode.GetString()
			hasPropName = true
		}
	} else {
		var ownerName string
		var hasOwnerName bool
		if name != nil {
			if gp, ok := ast.AsGetProp(lvalueNode); ok && ownerNode == gp.Object().Node() {
				propName = gp.Property()
				if len(*name) > len(propName)+1 && (*name)[len(*name)-len(propName)-1] == '.' {
					ownerName = (*name)[:len(*name)-len(propName)-1]
					hasOwnerName = true
					hasPropName = true
				}
			}
		}
		if !hasOwnerName {
			ownerName, hasOwnerName = b.getOwnerLValueNameString(ownerNode)
		}
		if hasOwnerName {
			ownerType =
				rhino.ObjectType_Cast(rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, ownerName))
		}

		if hasOwnerName && name != nil && !hasPropName {
			// TODO(b/111621092): Use the AST rather than manipulating strings here.
			jutil.CheckStateMsg2(
				jutil.StringStartsWith(*name, ownerName), "Expected \"%s\" to start with \"%s\"", *name, ownerName)
			propName = jutil.StringSubstringFrom(*name, jutil.StringLength(ownerName)+1)
			hasPropName = true
		}
	}

	prototypeOwner := b.GetPrototypeOwnerType(ownerType)
	var prototypeOwnerTypeMap *rhino.TemplateTypeMap
	if prototypeOwner != nil && prototypeOwner.GetTypeOfThis() != nil {
		prototypeOwnerTypeMap = prototypeOwner.GetTypeOfThis().GetTemplateTypeMap()
	}

	// Find the type of any overridden function.
	var overriddenType rhino.FunctionType
	if ownerType != nil && hasPropName {
		// the type of the property this overrides, not necessarily a function.
		overriddenPropType :=
			b.FindOverriddenProperty(ownerType, propName, prototypeOwnerTypeMap)
		if overriddenPropType != nil {
			// Overridden getters and setters need special handling because we declare
			// getters/setters as simple properties with their respective return/parameter type. This
			// causes a split during inference where left and right sides of a getter/setter
			// declaration will be inferred to have different types; if the left side has type `T`,
			// the right side will be some function type involving `T`.
			if lvalueNode.IsGetterDef() {
				// Convert `number` to `function(): number`
				overriddenType = b.outer.typeRegistry.CreateFunctionType(overriddenPropType)
			} else if lvalueNode.IsSetterDef() {
				// Convert `number` to `function(number): undefined`
				overriddenType =
					b.outer.typeRegistry.CreateFunctionType(b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE), overriddenPropType)
			} else if overriddenPropType.IsFunctionType() {
				// for cases where we override a non-method (e.g. a number) with a method, don't put the
				// non-method type (e.g. number) on the function.
				// Instead do some basic inference to create a function type.
				// we will warn during typechecking for an invalid override, but we don't want to put a
				// non-function type on this function because that will interfere with type inference
				// inside the function.
				overriddenType = overriddenPropType.ToMaybeFunctionType()
			}
		}
	}

	var hasNonEmptyReturn bool
	if !fnRoot.IsNil() {
		_, hasNonEmptyReturn = b.outer.functionsWithNonEmptyReturns[fnRoot]
	}

	bestTypeName := b.GetBestTypeName(lvalueNode, name)
	declarationScope := b.currentScope
	if !lvalueNode.IsNil() {
		declarationScope = b.GetLValueRootScope(lvalueNode)
	}
	// Go-only: language-adapter — a nil ObjectType must be passed as a nil
	// JSType interface.
	var prototypeOwnerOrNil rhino.JSType
	if prototypeOwner != nil {
		prototypeOwnerOrNil = prototypeOwner
	}
	builder := NewFunctionTypeBuilder(bestTypeName, b.outer.compiler, errorRoot, b.currentScope).
		SetSyntacticFunctionName(name).
		SetAstContents(fnRoot, hasNonEmptyReturn).
		SetDeclarationScope(declarationScope).
		InferFromOverriddenFunction(overriddenType, parametersNode).
		InferKind(info).
		InferClosurePrimitive(info).
		InferTemplateTypeName(info, prototypeOwnerOrNil).
		InferInheritance(info, nil)

	if info == nil || !info.HasReturnType() {
		// when there is no @return annotation, look for inline return type declaration
		if fn, ok := ast.AsFunction(rValue); ok && rValue.HasChildren() {
			nameDocInfo := fn.NameNode().Node().GetJSDocInfo()
			builder.InferReturnType(nameDocInfo, true)
		}
	} else {
		builder.InferReturnType(info, false)
	}

	// Infer the context type.
	var fallbackReceiverType rhino.JSType
	if ownerType != nil &&
		ownerType.IsFunctionPrototypeType() &&
		ownerType.GetOwnerFunction().HasInstanceType() {
		fallbackReceiverType = ownerType.GetOwnerFunction().GetInstanceType()
	} else if ownerType != nil &&
		ownerType.IsFunctionType() &&
		ownerType.ToMaybeFunctionType().HasInstanceType() &&
		!lvalueNode.IsNil() &&
		lvalueNode.IsStaticMember() {
		// Limit this case to members of ctors and interfaces decalared using `static`. Most
		// namespaces, like object literals, are assumed to declare free functions, so we exclude
		// them. Additionally, methods *assigned* to a ctor, especially an ES5 ctor, were never
		// designed with static polymorphism in mind, so excluding them preserves their assumptions.
		fallbackReceiverType = ownerType
	} else if !ownerNode.IsNil() && ownerNode.IsThis() {
		fallbackReceiverType = b.currentScope.GetTypeOfThis()
	}

	fnType := builder.
		InferThisTypeWithFallback(info, fallbackReceiverType).
		InferParameterTypesWithParamsParent(parametersNode, info).
		BuildAndRegister()

	// Do some additional validation for constructors and interfaces.
	if fnType.HasInstanceType() && !lvalueNode.IsNil() {
		prototypeSlot := fnType.GetSlot("prototype")

		// We want to make sure that the function and its prototype are declared at the same node.
		// This consistency is helpful to users of SymbolTable, because everything gets declared at
		// the same place.
		prototypeSlot.SetNode(lvalueNode)
	}
	return fnType
}

// CheckFunctionAliasAnnotations checks that the annotations in info are
// compatible with the aliased type. Any errors will be reported at n, which
// should be the qualified name node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.checkFunctionAliasAnnotations(Node, FunctionType, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) CheckFunctionAliasAnnotations(n rhino.Node, type_ rhino.FunctionType, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	var annotation *string
	if info.UsesImplicitMatch() {
		if !type_.IsStructuralInterface() {
			annotation = jutil.StringPtr("@record")
		}
	} else if info.IsInterface() {
		if !type_.IsInterface() {
			annotation = jutil.StringPtr("@interface")
		}
	} else if info.IsConstructor() && !type_.IsConstructor() {
		annotation = jutil.StringPtr("@constructor")
	}
	// TODO(sdh): consider checking @template, @param, @return, and/or @this.
	if annotation != nil &&
		// TODO(sdh): Remove this extra check once TypeScript stops passing us duplicate
		// conflicting externs.  In particular, TS considers everything an interface, but Closure
		// externs mark most things as @constructor.  The load order is not always the same, so
		// the error can show up in either the generated TS externs file or in our own extern.
		(!n.IsFromExterns() || *annotation == "@record") {
		b.outer.Report(JSError_MakeWithNode(n, TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION, *annotation, jutil.ToJavaString(n.GetQualifiedName())))
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getPrototypeOwnerType(ObjectType)
func (b *TypedScopeCreator_scopeBuilder) GetPrototypeOwnerType(ownerType rhino.ObjectType) rhino.ObjectType {
	if ownerType != nil && ownerType.IsFunctionPrototypeType() {
		if ownerFunction := ownerType.GetOwnerFunction(); ownerFunction != nil {
			return ownerFunction
		}
		return nil
	}
	return nil
}

// FindOverriddenProperty finds the property that's being overridden on this
// type, if any.
//
// Said property could be a method, field, getter, or setter. We don't
// distinguish between these when looking up a property type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.findOverriddenProperty(ObjectType, String, TemplateTypeMap)
func (b *TypedScopeCreator_scopeBuilder) FindOverriddenProperty(
	ownerType rhino.ObjectType, propName string, typeMap *rhino.TemplateTypeMap) rhino.JSType {
	var result rhino.JSType

	// First, check to see if the property is implemented
	// on a superclass.
	propType := ownerType.GetPropertyType(propName)
	if propType != nil && !propType.IsUnknownType() {
		result = propType
	} else {
		// If it's not, then check to see if it's implemented
		// on an implemented interface.
		for _, iface := range ownerType.GetCtorImplementedInterfaces() {
			propType = iface.GetPropertyType(propName)
			if propType != nil && !propType.IsUnknownType() {
				result = propType
				break
			}
		}
	}

	if result != nil && typeMap != nil && !typeMap.IsEmpty() {
		replacer := rhino.TemplateTypeReplacer_MakeForPartialReplacement(b.outer.typeRegistry, typeMap)
		result = replacer.VisitJSType(result)
	}

	return result
}

// CreateEnumTypeFromNodes creates a new enum type, based on the given nodes.
//
// This handles two cases that are semantically very different, but are not
// mutually exclusive: (1) An object literal that needs an enum type attached
// to it. (2) An assignment expression with an enum tag in the JsDoc.
//
// This function will always create an enum type, so only call it if you're
// sure that's what you want.
//
// rValue is the right-hand side of the enum, or null if none; lValue the
// left-hand side of the enum; name the qualified name of the enum; info the
// JSDocInfo attached to the enum definition.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createEnumTypeFromNodes(Node, String, Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) CreateEnumTypeFromNodes(
	rValue rhino.Node, name *string, lValue rhino.Node, info *rhino.JSDocInfo) *rhino.EnumType {
	jutil.CheckNotNull(info)
	jutil.CheckState(info.HasEnumParameterType())
	jutil.CheckStateMsg(
		!lValue.IsNil() || !rValue.IsNil(),
		"An enum initializer should come from either an lvalue or rvalue")

	var enumType *rhino.EnumType
	if !rValue.IsNil() {
		if rvalQname, ok := rValue.GetQualifiedNameString(); ok {
			// Handle an aliased enum. Note that putting @enum on an enum alias is optional. If the
			// rValue is not an enum, then this assignment errors during TypeCheck.
			var_ := b.currentScope.GetVar(rvalQname)
			if var_ != nil && var_.GetType() != nil && var_.GetType().IsEnumType() {
				enumType = var_.GetType().ToMaybeEnumType()
			}
		}
	}

	if enumType == nil {
		enumType = rhino.EnumType_Builder_(b.outer.typeRegistry).
			SetName(b.GetBestTypeName(lValue, name)).
			SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(b.currentScope)).
			SetSource(rValue).
			SetElementType(info.GetEnumParameterType().Evaluate(b.currentScope, b.outer.typeRegistry)).
			Build()

		if lit, ok := ast.AsObjectLit(rValue); ok {
			// collect enum elements
			for member := range lit.Members() {
				key := member.Node()
				if _, isComputed := member.AsComputedProp(); isComputed {
					b.outer.Report(JSError_MakeWithNode(key, TypedScopeCreator_INVALID_ENUM_KEY))
					continue
				}
				keyName := key.GetString()
				jutil.CheckNotNullMsg1(keyName, "Invalid enum key: %s", key)
				enumType.DefineElement(keyName, key)
			}
		}
	}

	if name != nil {
		b.outer.typeRegistry.DeclareType(b.currentScope, *name, enumType.GetElementsType())
	}

	if rValue.IsNil() || !(rValue.IsObjectLit() || rValue.IsQualifiedName()) {
		errorNode := rValue
		if !lValue.IsNil() {
			errorNode = lValue
		}
		b.outer.Report(JSError_MakeWithNode(errorNode, TypedScopeCreator_ENUM_INITIALIZER))
	}
	return enumType
}

// GetDeclaredType looks for a type declaration on a property assignment (in
// an ASSIGN or an object literal key).
//
// info is the doc info for this property; lValue the l-value node; rValue
// the node that n is being initialized to, or null if this is a stub
// declaration; declaredRValueTypeSupplier a supplier for the declared type
// of the rvalue, used for destructuring declarations where we have to do
// additional work on the rvalue.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredType(JSDocInfo, Node, Node, Supplier)
func (b *TypedScopeCreator_scopeBuilder) GetDeclaredType(
	info *rhino.JSDocInfo,
	lValue rhino.Node,
	rValue rhino.Node,
	declaredRValueTypeSupplier func() typedScopeRValueInfo) rhino.JSType {
	if info != nil && info.HasType() {
		type_ := b.GetDeclaredTypeInAnnotation(lValue, info)

		if gp, ok := ast.AsGetProp(lValue); ok &&
			type_.IsSymbol() &&
			lValue.IsFromExterns() &&
			gp.Object().Node().MatchesName("Symbol") &&
			b.currentScope.IsGlobal() {
			// Create a unique subtype type for this name.
			// Currently we only look specifically for extern properties on the global Symbol object
			// itself - so
			//   /** @const {symbol} */
			//   Symbol.iterator;
			// but not:
			//   /** @const {symbol} */
			//   const x = Symbol();
			// We could consider relaxing that in the future & also defining known symbols from
			// source code.
			lvalQname, ok := lValue.GetQualifiedNameString()
			if !ok {
				panic(jutil.NewNullPointerException())
			}
			return rhino.NewKnownSymbolType(b.outer.typeRegistry, lvalQname)
		}
		return type_
	} else if !rValue.IsNil() &&
		rValue.IsFunction() &&
		b.ShouldUseFunctionLiteralType(
			rhino.JSType_ToMaybeFunctionType(rValue.GetJSType()), info, lValue) {
		return rValue.GetJSType()
	} else if !rValue.IsNil() && rValue.IsClass() {
		return rValue.GetJSType()
	} else if info != nil {
		if info.HasEnumParameterType() {
			if !rValue.IsNil() && rValue.IsObjectLit() {
				return rValue.GetJSType()
			} else {
				return b.CreateEnumTypeFromNodes(rValue, lValue.GetQualifiedName(), lValue, info)
			}
		} else if info.IsConstructorOrInterface() {
			fnType :=
				b.CreateFunctionTypeFromNodes(rValue, lValue.GetQualifiedName(), info, lValue)
			if rValue.IsNil() && !lValue.IsFromExterns() {
				diagnostic := TypedScopeCreator_IFACE_INITIALIZER
				if fnType.IsConstructor() {
					diagnostic = TypedScopeCreator_CTOR_INITIALIZER
				}
				b.outer.Report(
					JSError_MakeWithNode(
						lValue,
						diagnostic,
						jutil.ToJavaString(lValue.GetQualifiedName())))
			}
			return fnType
		}
	}

	// Check if this is constant and if it has a known type.
	isDefaultExport := false
	if exp, ok := ast.AsExport(lValue); ok && exp.IsDefault() {
		isDefaultExport = true
	}
	if isDefaultExport ||
		NodeUtil_IsConstantDeclaration(info, lValue) ||
		b.IsGoogModuleExports(lValue) {
		if !rValue.IsNil() {
			rValueType := b.GetDeclaredRValueType(lValue, rValue)
			b.DeclareAliasTypeIfRvalueIsAliasable(
				lValue, rValue.GetQualifiedNameObject(), rValueType, b.currentScope)
			if rValueType != nil {
				return rValueType
			}
		} else if declaredRValueTypeSupplier != nil {
			rvalueInfo := declaredRValueTypeSupplier()
			b.DeclareAliasTypeIfRvalueIsAliasable(
				lValue, rvalueInfo.qualifiedName, rvalueInfo.type_, b.currentScope)
			if rvalueInfo.type_ != nil {
				return rvalueInfo.type_
			}
		}
	}

	if assign, ok := ast.AsAssign(rValue); ok {
		// Handle nested assignments. For example, TypeScript generates code like this:
		//   var Foo_1;
		//   let Foo = Foo_1 = class Foo {}
		//   Foo = Foo_1 = tslib_1.decorate(..., Foo);
		return b.GetDeclaredType(info, lValue, assign.Value().Node(), nil)
	}

	if info != nil && FunctionTypeBuilder_IsFunctionTypeDeclaration(info) {
		fnName := lValue.GetQualifiedName()
		return b.CreateFunctionTypeFromNodes(rhino.NilNode, fnName, info, lValue)
	}

	if b.IsValidTypedefDeclaration(lValue, info) {
		return b.outer.GetNativeType(rhino.JSTypeNative_NO_TYPE)
	}

	return nil
}

// DeclareAliasTypeIfRvalueIsAliasable: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// This method recognizes three kinds of type aliases: @typedefs,
// @constructor/@interface types, and @enums.
//
// Given any of those three types, this method redeclares the aliasing name
// in the typeRegistry. For @typedefs and global @enums, this method also
// marks the qualified name referring to the type as non-nullable by
// default.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(Node, QualifiedName, JSType, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DeclareAliasTypeIfRvalueIsAliasable(
	lValue rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope) {
	if rValue == nil {
		return
	}
	b.DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
		lValue.GetQualifiedName(), lValue, rValue, rValueType, rValueLookupScope, b.currentScope)
}

// DeclareAliasTypeIfRvalueIsAliasableWithLValueName: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// NOTE: in most cases, call the version with fewer arguments. This version
// only exists to handle goog.declareLegacyNamespace, which is strange
// compared to normal type aliasing because 1) there's no GETPROP node
// representing the lvalue and 2) the type is declared in the global scope,
// not the current module-local scope.
//
// lValueName is the fully qualified lValue name, if any. If null, all this
// method will do is propagate the @typedef Node annotation to
// actualLvalueNode. aliasDeclarationScope is the scope in which to declare
// the alias name. In most cases, this should just be the currentScope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(String, Node, QualifiedName, JSType, TypedScope, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
	lValueName *string,
	actualLvalueNode rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope,
	aliasDeclarationScope *TypedScope) {
	// NOTE: this allows some strange patterns such allowing instance properties
	// to be aliases of constructors, and then creating a local alias of that to be
	// used as a type name.  Consider restricting this.

	if rValue == nil {
		return
	}

	// Look for a @typedef annotation on the definition node
	definitionNode := b.GetDefinitionNode(rValue, rValueLookupScope)
	if !definitionNode.IsNil() {
		typedefType := definitionNode.GetTypedefTypeProp()
		if typedefType != nil {
			// Propagate typedef type to typedef aliases.
			actualLvalueNode.SetTypedefTypeProp(typedefType)
			if lValueName != nil {
				b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
				b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, typedefType)
			}
			return
		}
	}

	if lValueName == nil {
		return
	}

	// Check if the provided rValueType indicates that we should declare this type
	// Note that we only look for enums and constructors/interfaces here: this step cannot work
	// for @typedefs. The 'type' of the TypedVar representing a @typedef'd name is the None type,
	// not the @typedef'd type.
	if rValueType != nil &&
		rValueType.IsFunctionType() &&
		rValueType.ToMaybeFunctionType().HasInstanceType() {
		// Look for @constructor/@interface by checking if the RHS has an instance type
		functionType := rValueType.ToMaybeFunctionType()
		b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, functionType.GetInstanceType())
		return
	}

	if rValueType != nil && rValueType.IsEnumType() {
		// Look for cases where the rValue is an Enum namespace
		b.outer.typeRegistry.DeclareType(
			aliasDeclarationScope, *lValueName, rValueType.ToMaybeEnumType().GetElementsType())
		b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
	}
}

// IsGoogModuleExports reports whether this lvalue is either `exports`,
// `exports.x`, or a string key in `exports = {x}`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isGoogModuleExports(Node)
func (b *TypedScopeCreator_scopeBuilder) IsGoogModuleExports(lValue rhino.Node) bool {
	if b.GetModule() == nil || lValue.IsNil() {
		return false
	}
	if _, ok := b.outer.undeclaredNamesForClosure[lValue]; ok {
		// includes "exports" and "exports.x"
		return true
	}
	// ASSIGN
	//   NAME "exports"
	//   OBJECT_LIT
	//     STRING_KEY "x"
	//     [value]
	if !lValue.IsStringKey() || !lValue.GetParent().IsObjectLit() {
		return false
	}
	assign, ok := ast.AsAssign(lValue.GetGrandparent())
	if !ok {
		return false
	}
	target := assign.Target().Node()
	if !target.MatchesName("exports") {
		return false
	}
	_, ok = b.outer.undeclaredNamesForClosure[target]
	return ok
}

// GetDefinitionNode returns the AST node associated with the definition, if
// any.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDefinitionNode(QualifiedName, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) GetDefinitionNode(qname rhino.QualifiedName, scope *TypedScope) rhino.Node {
	if qname.IsSimple() {
		var_ := scope.GetVar(qname.GetComponent())
		if var_ != nil {
			return var_.GetNameNode()
		}
		return rhino.NilNode
	}
	parent := rhino.ObjectType_Cast(scope.LookupQualifiedName(qname.GetOwner()))
	if parent != nil {
		return parent.GetPropertyDefSite(qname.GetComponent())
	}
	return rhino.NilNode
}

// GetDeclaredRValueType checks for common idioms of a typed R-value
// assigned to a const L-value.
//
// Normally, we would only want this sort of propagation to happen under
// type inference. But we want a declared const to be nameable in a type
// annotation, so we need to figure out the type before we try to resolve
// the annotation.
//
// lValue is the lvalue node if this is a simple assignment, null for
// destructuring.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredRValueType(Node, Node)
func (b *TypedScopeCreator_scopeBuilder) GetDeclaredRValueType(lValue rhino.Node, rValue rhino.Node) rhino.JSType {
	// If rValue has a type-cast, we use the type in the type-cast.
	rValueInfo := rValue.GetJSDocInfo()
	if rValue.IsCast() && rValueInfo != nil && rValueInfo.HasType() {
		return rValueInfo.GetType().Evaluate(b.currentScope, b.outer.typeRegistry)
	}

	// Check if the type has already been computed during scope-creation.
	// This is mostly useful for literals like BOOLEAN, NUMBER, STRING, and
	// OBJECT_LITERAL
	type_ := rValue.GetJSType()
	if type_ != nil && !type_.IsUnknownType() {
		return type_
	}

	// If rValue is a name, try looking it up in the current scope.
	if qname, ok := rValue.GetQualifiedNameString(); ok {
		return rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, qname)
	}

	// Check for simple invariant operations, such as "!x" or "+x" or "''+x"
	if NodeUtil_IsBooleanResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	}

	if NodeUtil_IsNumericResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	}

	if NodeUtil_IsBigIntResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	}

	if NodeUtil_IsStringResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	}

	if nw, ok := ast.AsNew(rValue); ok {
		if qname, ok := nw.Callee().Node().GetQualifiedNameString(); ok {
			targetType := rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, qname)
			if targetType != nil {
				fnType := targetType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
				if fnType != nil && fnType.HasInstanceType() {
					return fnType.GetInstanceType()
				}
			}
		}
	}

	if call, ok := ast.AsCall(rValue); ok && call.Callee().Node().MatchesName("Symbol") {
		// Type calls like `Symbol('foo')`
		symbolVar := b.currentScope.GetVar("Symbol")
		if symbolVar != nil && symbolVar.GetScope().(*TypedScope).IsGlobal() {
			return b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
		}
	}

	// Check for a very specific JS idiom:
	// var x = x || TYPE;
	// This is used by Closure's base namespace for esoteric
	// reasons, so we only really care about that case.
	if rValue.IsOr() {
		orExpr := ast.MustBinary(rValue)
		if leftName, ok := orExpr.Left().AsName(); ok {
			if lvalName, ok := ast.AsName(lValue); ok && leftName.Ident() == lvalName.Ident() {
				type_ = orExpr.Right().Node().GetJSType()
				if type_ != nil && !type_.IsUnknownType() {
					return type_
				}
			}
		}
	}

	return nil
}

// CheckForClassDefiningCalls looks for class-defining calls. Because JS has
// no 'native' syntax for defining classes, this is often very
// coding-convention dependent and business-logic heavy.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.checkForClassDefiningCalls(Node)
func (b *TypedScopeCreator_scopeBuilder) CheckForClassDefiningCalls(n rhino.Node) {
	relationship := b.outer.codingConvention.GetClassesDefinedByCall(n)
	if relationship != nil {
		superClass :=
			TypeValidator_GetInstanceOfCtor(
				rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, relationship.superclassName))
		subClass :=
			TypeValidator_GetInstanceOfCtor(
				rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, relationship.subclassName))
		if superClass != nil && subClass != nil {
			// superCtor and subCtor might be structural constructors
			// (like {function(new:Object)}) so we need to resolve them back
			// to the original ctor objects.
			superCtor := superClass.GetConstructor()
			subCtor := subClass.GetConstructor()
			if superCtor != nil && subCtor != nil {
				b.outer.codingConvention.ApplySubclassRelationship(
					rhino.NewNominalTypeBuilder(superCtor, superClass),
					rhino.NewNominalTypeBuilder(subCtor, subClass),
					relationship.type_)
			}
		}
	}

	singletonGetterClassName := b.outer.codingConvention.GetSingletonGetterClassName(n)
	if singletonGetterClassName != nil {
		objectType :=
			rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *singletonGetterClassName))
		if objectType != nil {
			functionType := objectType.GetConstructor()

			if functionType != nil {
				getterType := b.outer.typeRegistry.CreateFunctionType(objectType)
				b.outer.codingConvention.ApplySingletonGetter(
					rhino.NewNominalTypeBuilder(functionType, objectType), getterType)
			}
		}
	}

	objectLiteralCast := b.outer.codingConvention.GetObjectLiteralCast(n)
	if objectLiteralCast != nil {
		if objectLiteralCast.diagnosticType == nil {
			type_ :=
				rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *objectLiteralCast.typeName))
			if type_ != nil && type_.GetConstructor() != nil {
				b.outer.SetDeferredType(objectLiteralCast.objectNode, type_)
				objectLiteralCast.objectNode.PutBooleanProp(rhino.Node_REFLECTED_OBJECT, true)
			} else {
				b.outer.Report(JSError_MakeWithNode(n, TypedScopeCreator_CONSTRUCTOR_EXPECTED))
			}
		} else {
			b.outer.Report(JSError_MakeWithNode(n, objectLiteralCast.diagnosticType))
		}
	}
}

// MaybeDeclareGetElem declares a computed property on its owner type.
//
// info is the doc info for this property; n a top-level GETELEM node (it
// should not be contained inside another GETPROP); rhsValue the node that n
// is being initialized to, or null if this is a stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareGetElem(JSDocInfo, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) MaybeDeclareGetElem(info *rhino.JSDocInfo, n rhino.Node, rhsValue rhino.Node) {
	ge := ast.MustGetElem(n)
	ownerNode := ge.Object().Node()
	key := ge.Index().Node()
	keyType := TypedScopeCreator_ExtractKnownSymbolKey(b.currentScope, key)
	if keyType == nil {
		return
	}
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil {
		return
	}
	ownerName, ok := ownerNode.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	ownerType := b.GetObjectSlot(ownerName)
	if ownerType == nil {
		return
	}
	b.DeclarePropertyIfNamespaceType(
		ownerType,
		ownerNode,
		rhino.NewProperty_SymbolKey(keyType),
		valueType,
		n,
		/* jsdocInfo= */ nil)
}

// MaybeDeclareQualifiedName declares the symbol for a qualified name in the
// current scope.
//
// info is the doc info for this property; n a top-level GETPROP node (it
// should not be contained inside another GETPROP); parent the parent of n;
// rhsValue the node that n is being initialized to, or null if this is a
// stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareQualifiedName(NodeTraversal, JSDocInfo, Node, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) MaybeDeclareQualifiedName(
	t *NodeTraversal, info *rhino.JSDocInfo, n rhino.Node, parent rhino.Node, rhsValue rhino.Node) {
	isTypedef := b.IsValidTypedefDeclaration(n, info)
	if isTypedef {
		b.DeclareTypedefType(n, info)
	}

	gp := ast.MustGetProp(n)
	ownerNode := gp.Object().Node()
	var qName string
	var hasQName bool
	if !rhsValue.IsNil() && rhsValue.IsFunction() {
		if origID := rhsValue.GetOriginalNameID(); origID != 0 {
			origStr := rhino.RhinoStringPool_StringFromID(origID)
			if n.MatchesQualifiedName(origStr) {
				qName = origStr
				hasQName = true
			}
		}
	}
	if !hasQName {
		qName, hasQName = n.GetQualifiedNameString()
	}
	propName := gp.Property()
	jutil.CheckArgument(hasQName)
	ownerNameStr := qName[:len(qName)-len(propName)-1]

	// Precedence of type information on GETPROPs:
	// 1) @type annotation / @enum annotation
	// 2) ASSIGN to FUNCTION literal
	// 3) @param/@return annotation (with no function literal)
	// 4) ASSIGN to something marked @const
	// 5) ASSIGN to anything else
	//
	// 1, 3, and 4 are declarations, 5 is inferred, and 2 is a declaration iff
	// the function has JsDoc or has not been declared before.
	//
	// FUNCTION literals are special because TypedScopeCreator is very smart
	// about getting as much type information as possible for them.

	// Determining type for #1 + #2 + #3 + #4
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil && !rhsValue.IsNil() {
		// Determining type for #5
		valueType = rhsValue.GetJSType()
	}

	// Function prototypes are special.
	// It's a common JS idiom to do:
	// F.prototype = { ... };
	// So if F does not have an explicitly declared super type,
	// allow F.prototype to be redefined arbitrarily.
	if "prototype" == propName {
		qVar := b.currentScope.GetVar(qName)
		if qVar != nil {
			// If the programmer has declared that F inherits from Super,
			// and they assign F.prototype to an object literal,
			// then they are responsible for making sure that the object literal's
			// implicit prototype is set up appropriately. We just obey
			// the @extends tag.
			qVarType := rhino.ObjectType_Cast(qVar.GetType())
			if qVarType != nil && !rhsValue.IsNil() && rhsValue.IsObjectLit() {
				b.outer.typeRegistry.ResetImplicitPrototype(
					rhsValue.GetJSType(), qVarType.GetImplicitPrototype())
			} else if !qVar.IsTypeInferred() {
				// If the programmer has declared that F inherits from Super,
				// and they assign F.prototype to some arbitrary expression,
				// there's not much we can do. We just ignore the expression,
				// and hope they've annotated their code in a way to tell us
				// what props are going to be on that prototype.
				return
			}

			qVar.GetScope().(*TypedScope).Undeclare(qVar)
		}
	}

	if valueType == nil {
		if parent.IsExprResult() {
			// t is mutable so make sure to capture the current state before the lambda.
			isExtern := t.GetInput() != nil && t.GetInput().IsExtern()
			b.deferAction(b.currentScope.GetRootNode(), func() { b.ResolveStubDeclaration(n, isExtern, ownerNameStr) })
		}

		return
	}

	inferred := b.IsQualifiedNameInferred(qName, hasQName, n, info, rhsValue, valueType)
	if !inferred {
		ownerType := b.GetObjectSlot(ownerNameStr)
		if ownerType != nil {
			b.DeclarePropertyIfNamespaceType(
				ownerType, ownerNode, rhino.NewProperty_StringKey(propName), valueType, n /* jsdocInfo= */, nil)
		}

		// this is a memory optimization: we don't need to declare .prototype props in the scope.
		// NOTE: in theory we could possibly extend this to more kinds of declarations and only
		// declare simple (non-qualified) names in the scope, but that seems to cause a lot more
		// issues with older code.
		declareInScope := ownerType == nil || !ownerType.IsFunctionPrototypeType()

		// If the property is already declared, the error will be caught when we try to declare it
		// in the current scope.
		if declareInScope {
			b.defineSlot(typedScopeSlotSpec{
				declNode:                n,
				varName:                 qName,
				scope:                   b.GetLValueRootScope(n),
				type_:                   valueType,
				allowLaterTypeInference: inferred,
			})
		}
	}
}

// IsQualifiedNameInferred determines whether a qualified name is inferred.
// NOTE(nicksantos): Determining whether a property is declared or not is
// really really obnoxious.
//
// The problem is that there are two (equally valid) coding styles:
//
// (function() { /* The authoritative definition of goog.bar. / goog.bar =
// function() {}; })();
//
// function f() { goog.bar(); /* Reset goog.bar to a no-op. / goog.bar =
// function() {}; }
//
// In a dynamic language with first-class functions, it's very difficult to
// know which one the user intended without looking at lots of contextual
// information (the second example demonstrates a small case of this, but
// there are some really pathological cases as well).
//
// The current algorithm checks if either the declaration has JsDoc type
// information, or @const with a known type, or a function literal with a
// name we haven't seen before.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isQualifiedNameInferred(String, Node, JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) IsQualifiedNameInferred(
	qName string,
	hasQName bool,
	n rhino.Node,
	info *rhino.JSDocInfo,
	rhsValue rhino.Node,
	valueType rhino.JSType) bool {
	// Prototypes of constructors and interfaces are always declared.
	if hasQName && strings.HasSuffix(qName, ".prototype") {
		className := qName[:len(qName)-len(".prototype")]
		slot := b.currentScope.GetVar(className)
		var classType rhino.JSType
		if slot != nil {
			classType = slot.GetType()
		}
		if classType != nil && (classType.IsConstructor() || classType.IsInterface()) {
			return false
		}
	}
	// treat "foo = bar = <VAL>" the same as "foo = <VAL>"
	rhsValue = b.UnwrapIfAssign(rhsValue)

	// If the jsdoc or RHS specifies a concrete type, it's not inferred.
	if (info != nil &&
		(info.HasType() ||
			info.HasEnumParameterType() ||
			b.IsValidTypedefDeclaration(n, info) ||
			FunctionTypeBuilder_IsFunctionTypeDeclaration(info) ||
			(!rhsValue.IsNil() && rhsValue.IsFunction()))) ||
		b.IsTypedConstantDeclaration(info, n, valueType) {
		return false
	}

	// At this point, we're pretty sure it's inferred, since there's neither
	// useful jsdoc info, nor a useful const or doc'd function RHS.  But
	// there's still one case where it may still not be: if the RHS is a
	// class or function that is not
	//   (1) a scoped qualified name (i.e. this.b.c or super.b.c),
	//   (2) already declared in a scope,
	//   (3) assigned in a conditional block, or
	//   (4) escaped to a closure,
	// then we treat it as if it is declared, rather than inferred.
	// Stubs and other values are always considered inferred at this point.
	if rhsValue.IsNil() || (!rhsValue.IsFunction() && !rhsValue.IsClass()) {
		return true
	}

	// "Scoped" qualified names (e.g. this.b.c or super.d) are inferred.
	if !n.IsUnscopedQualifiedName() {
		return true
	}

	// If this qname is already declared then treat this definition as inferred.
	ownerScope := b.GetLValueRootScope(n)
	if ownerScope != nil && ownerScope.HasOwnSlot(qName) {
		return true
	}

	// Check if this is in a conditional block.
	// Functions assigned in conditional blocks are inferred.
	if b.HasControlStructureAncestor(n.GetParent()) {
		return true
	}

	// Check if this is assigned in an inner scope.
	// Functions assigned in inner scopes are inferred.
	if ownerScope != nil && len(b.outer.varFlags) > 0 {
		if (b.outer.varFlags[typedScopeCreator_ScopedNameKey{name: qName, scopeRoot: ownerScope.GetRootNode()}] & typedScopeCreator_VarFlagEscaped) != 0 {
			return true
		}
	}

	return false
}

// Java: TypedScopeCreator.AbstractScopeBuilder.unwrapIfAssign(Node)
func (b *TypedScopeCreator_scopeBuilder) UnwrapIfAssign(maybeAssign rhino.Node) rhino.Node {
	for {
		assign, ok := ast.AsAssign(maybeAssign)
		if !ok {
			return maybeAssign
		}
		maybeAssign = assign.Value().Node()
	}
}

// DeclareProvidedNs: given a `goog.provide()` or legacy `goog.module()` call
// and implicit ProvidedName, declares the name in the global scope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareProvidedNs(Node, ProvidedName)
func (b *TypedScopeCreator_scopeBuilder) DeclareProvidedNs(provideCall rhino.Node, providedName *ProcessClosureProvidesAndRequires_ProvidedName) {
	// Redefine this name if we haven't already added a provide definition.
	// Note: in some cases, this will cause a redefinition error.
	anonymousObjectType := b.outer.typeRegistry.CreateAnonymousObjectType(nil)
	b.defineSlot(typedScopeSlotSpec{
		declNode:                provideCall,
		varName:                 providedName.GetNamespace(),
		scope:                   b.currentScope.GetGlobalScope(),
		type_:                   anonymousObjectType,
		allowLaterTypeInference: false,
		forGoogProvidedName:     true,
	})

	namespace := rhino.QualifiedName_Of(providedName.GetNamespace())
	if !namespace.IsSimple() {
		ownerType := b.currentScope.LookupQualifiedName(namespace.GetOwner())
		if ownerType != nil && ownerType.IsObjectType() {
			ownerType.
				ToMaybeObjectType().
				DefineDeclaredProperty(namespace.GetComponent(), anonymousObjectType, provideCall)
		}
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.isTypedConstantDeclaration(JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) IsTypedConstantDeclaration(info *rhino.JSDocInfo, n rhino.Node, valueType rhino.JSType) bool {
	return (NodeUtil_IsConstantDeclaration(info, n) || b.IsGoogModuleExports(n)) &&
		valueType != nil
}

// Java: TypedScopeCreator.AbstractScopeBuilder.hasControlStructureAncestor(Node)
func (b *TypedScopeCreator_scopeBuilder) HasControlStructureAncestor(n rhino.Node) bool {
	for !(n.IsScript() || n.IsFunction()) {
		if NodeUtil_IsControlStructure(n) {
			return true
		}
		n = n.GetParent()
	}
	return false
}

// GetObjectSlot finds the ObjectType associated with the given slot.
//
// slotName is the name of the slot to find the type in. Returns an object
// type, or null if this slot does not contain an object.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getObjectSlot(String)
func (b *TypedScopeCreator_scopeBuilder) GetObjectSlot(slotName string) rhino.ObjectType {
	ownerVar := b.currentScope.GetVar(slotName)
	if ownerVar != nil {
		ownerVarType := ownerVar.GetType()
		if ownerVarType == nil {
			return rhino.ObjectType_Cast(nil)
		}
		return rhino.ObjectType_Cast(ownerVarType.RestrictByNotNullOrUndefined())
	}
	return nil
}

// GetInheritedInterfacePropertyType: when a class has a stub for a
// property, and the property exists on a super interface, use that type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getInheritedInterfacePropertyType(ObjectType, String)
func (b *TypedScopeCreator_scopeBuilder) GetInheritedInterfacePropertyType(obj rhino.ObjectType, propName string) rhino.JSType {
	if obj != nil && obj.IsFunctionPrototypeType() {
		f := obj.GetOwnerFunction()
		for _, i := range f.GetImplementedInterfaces() {
			if i.HasProperty(propName) {
				return i.GetPropertyType(propName)
			}
		}
	}
	return nil
}

// ResolveStubDeclaration resolves any type-less stub declarations to
// unknown types if we could not find types for them during traversal. This
// method is only called as a deferred action after the root node is visted.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.resolveStubDeclaration(Node, boolean, String)
func (b *TypedScopeCreator_scopeBuilder) ResolveStubDeclaration(n rhino.Node, isExtern bool, ownerName string) {
	qName, ok := n.GetQualifiedNameString()
	propName := n.GetString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}

	// TODO(b/111216910): should this be getLValueRoot(n).hasOwnSlot(qName)?
	if b.currentScope.HasOwnSlot(qName) {
		return
	}

	// If we see a stub property, make sure to register this property
	// in the type registry.
	ownerType := b.GetObjectSlot(ownerName)
	inheritedType := b.GetInheritedInterfacePropertyType(ownerType, propName)
	var stubType rhino.JSType
	if inheritedType == nil {
		stubType = b.outer.unknownType
	} else {
		stubType = inheritedType
	}
	b.defineSlot(typedScopeSlotSpec{
		declNode:                n,
		varName:                 slotVarNameFromDeclNode(n),
		scope:                   b.GetLValueRootScope(n),
		type_:                   stubType,
		allowLaterTypeInference: true,
	})

	if ownerType != nil && (isExtern || ownerType.IsFunctionPrototypeType()) {
		// If this is a stub for a prototype, just declare it
		// as an unknown type. These are seen often in externs.
		ownerType.DefineInferredProperty(propName, stubType, n)
	} else {
		var registerOn rhino.JSType
		if ownerType == nil {
			registerOn = stubType
		} else {
			registerOn = ownerType
		}
		b.outer.typeRegistry.RegisterPropertyOnType(propName, registerOn)
	}
}

// IsValidTypedefDeclaration returns whether this is a valid declaration of a
// @typedef.
//
// candidate is a qualified name node; info JSDoc comments.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isValidTypedefDeclaration(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) IsValidTypedefDeclaration(candidate rhino.Node, info *rhino.JSDocInfo) bool {
	if info == nil || !info.HasTypedefType() {
		return false
	}
	// `isUnscopedQualifiedName` excludes `this` and `super` properties.
	return candidate.IsUnscopedQualifiedName() && !NodeUtil_IsPrototypeProperty(candidate)
}

// DeclareTypedefType declares a typedef'd name in the JSTypeRegistry.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareTypedefType(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) DeclareTypedefType(candidate rhino.Node, info *rhino.JSDocInfo) {
	typedef, ok := candidate.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}

	// TODO(nicksantos|user): This is a terrible, terrible hack
	// to bail out on recursive typedefs. We'll eventually need
	// to handle these properly.
	b.outer.typeRegistry.DeclareType(b.currentScope, typedef, b.outer.unknownType)

	realType := info.GetTypedefType().Evaluate(b.currentScope, b.outer.typeRegistry)
	if realType == nil {
		b.outer.Report(JSError_MakeWithNode(candidate, TypedScopeCreator_MALFORMED_TYPEDEF, typedef))
	} else {
		candidate.SetTypedefTypeProp(realType)
	}

	b.outer.typeRegistry.OverwriteDeclaredType(b.currentScope, typedef, realType)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.declarePropertyIfNamespaceType(ObjectType, Node, Property.Key, JSType, Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) DeclarePropertyIfNamespaceType(
	ownerType rhino.ObjectType,
	ownerNode rhino.Node,
	propName rhino.Property_Key,
	valueType rhino.JSType,
	declarationNode rhino.Node,
	jsdocInfo *rhino.JSDocInfo) {
	// Only declare this as an official property if it has not been
	// declared yet.
	if ownerType.HasOwnPropertyByKey(propName) && !ownerType.IsPropertyTypeInferredByKey(propName) {
		return
	}
	// Define the property if any of the following are true:
	//   (1) it's a non-native extern type. Native types are excluded here because we don't
	//       want externs of the form "/** @type {!Object} */ var api = {}; api.foo;" to
	//       cause a property "foo" to be declared on Object.
	//   (2) it's a non-instance type. This primarily covers static properties on
	//       constructors (which are FunctionTypes, not InstanceTypes).
	//   (3) it's an assignment to 'this', which covers instance properties assigned in
	//       constructors or other methods.
	isNonNativeExtern :=
		b.GetCompilerInput() != nil &&
			b.GetCompilerInput().IsExtern() &&
			!ownerType.IsNativeObjectType()
	if isNonNativeExtern || !ownerType.IsInstanceType() || ownerNode.IsThis() {
		// If the property is undeclared or inferred, declare it now.
		ownerType.DefineDeclaredPropertyByKey(propName, valueType, declarationNode)
		ownerType.SetPropertyJSDocInfoByKey(propName, jsdocInfo)
	}
} // end AbstractScopeBuilder

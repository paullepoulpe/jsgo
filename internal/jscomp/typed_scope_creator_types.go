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
	var builder FunctionTypeBuilder
	builder.InitFunctionTypeBuilder(name, b.outer.compiler, clazz, b.currentScope)
	builder.UsingClassSyntax()
	builder.SetSyntacticFunctionName(syntacticName)
	builder.SetAstContents(clazz, false)
	builder.SetDeclarationScope(declarationScope)
	builder.InferKind(info)
	builder.InferTemplateTypeName(info, nil)

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
	if extendsCall == nil {
		// Handle invalid goog.module.get calls
		return nil
	}
	moduleScope, ok := b.outer.memoized[extendsCall.GetScopeRoot()]
	if !ok {
		// Handle invalid goog.module.get calls
		return nil
	}
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
	var builder FunctionTypeBuilder
	builder.InitFunctionTypeBuilder(bestTypeName, b.outer.compiler, errorRoot, b.currentScope)
	builder.SetSyntacticFunctionName(name)
	builder.SetAstContents(fnRoot, hasNonEmptyReturn)
	builder.SetDeclarationScope(declarationScope)
	builder.InferFromOverriddenFunction(overriddenType, parametersNode)
	builder.InferKind(info)
	builder.InferClosurePrimitive(info)
	builder.InferTemplateTypeName(info, prototypeOwnerOrNil)
	builder.InferInheritance(info, nil)

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

	builder.InferThisTypeWithFallback(info, fallbackReceiverType)
	builder.InferParameterTypesWithParamsParent(parametersNode, info)
	fnType := builder.BuildAndRegister()

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

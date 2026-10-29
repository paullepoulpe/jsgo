// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CreateOptionalType creates a type representing optional values of the
// given type. Returns the union of the type and the void type
//
// Java: JSTypeRegistry.createOptionalType(JSType)
func (r *JSTypeRegistry) CreateOptionalType(type_ JSType) JSType {
	if _, ok := type_.(*UnknownType); ok || type_.IsAllType() {
		return type_
	} else {
		if type_ == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index 0"))
		}
		b := r.acquireUnionBuilder()
		b.AddAlternate(type_)
		b.AddAlternate(r.GetNativeType(JSTypeNative_VOID_TYPE))
		res := b.Build()
		r.releaseUnionBuilder(b)
		return res
	}
}

// CreateNullableType creates a type representing nullable values of the
// given type. Returns the union of the type and the Null type
//
// Java: JSTypeRegistry.createNullableType(JSType)
func (r *JSTypeRegistry) CreateNullableType(type_ JSType) JSType {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("at index 0"))
	}
	b := r.acquireUnionBuilder()
	b.AddAlternate(type_)
	b.AddAlternate(r.GetNativeType(JSTypeNative_NULL_TYPE))
	res := b.Build()
	r.releaseUnionBuilder(b)
	return res
}

// CreateOptionalNullableType creates a nullable and undefine-able value of
// the given type. Returns the union of the type and null and undefined.
//
// Java: JSTypeRegistry.createOptionalNullableType(JSType)
func (r *JSTypeRegistry) CreateOptionalNullableType(type_ JSType) JSType {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("at index 0"))
	}
	b := r.acquireUnionBuilder()
	b.AddAlternate(type_)
	b.AddAlternate(r.GetNativeType(JSTypeNative_VOID_TYPE))
	b.AddAlternate(r.GetNativeType(JSTypeNative_NULL_TYPE))
	res := b.Build()
	r.releaseUnionBuilder(b)
	return res
}

// CreateUnionType creates a union type whose variants are the arguments.
//
// Java: JSTypeRegistry.createUnionType(JSType...)
func (r *JSTypeRegistry) CreateUnionType(variants ...JSType) JSType {
	for i, e := range variants {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
	}
	return r.CreateUnionTypeOfList(variants)
}

// Java: JSTypeRegistry.createUnionType(List<? extends JSType>)
func (r *JSTypeRegistry) CreateUnionTypeOfList(variants []JSType) JSType {
	b := r.acquireUnionBuilder()
	res := b.AddAlternates(variants).Build()
	r.releaseUnionBuilder(b)
	return res
}

// CreateUnionTypeOfNatives creates a union type whose variants are the
// built-in types specified by the arguments.
//
// Java: JSTypeRegistry.createUnionType(JSTypeNative...)
func (r *JSTypeRegistry) CreateUnionTypeOfNatives(variants ...JSTypeNative) JSType {
	b := r.acquireUnionBuilder()
	for _, type_ := range variants {
		b.AddAlternate(r.GetNativeType(type_))
	}
	res := b.Build()
	r.releaseUnionBuilder(b)
	return res
}

// CreateFunctionType creates a function type.
//
// returnType: the function's return type. parameterTypes: the parameters'
// types.
//
// Java: JSTypeRegistry.createFunctionType(JSType, JSType...)
func (r *JSTypeRegistry) CreateFunctionType(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateFunctionTypeWithParameters(returnType, r.CreateParameters(parameterTypes...))
}

// CreateFunctionTypeWithParameters creates a function type.
//
// parameters: the function's parameters or null to indicate that the
// parameter types are unknown. returnType: the function's return type or
// null to indicate that the return type is unknown.
//
// Java: JSTypeRegistry.createFunctionType(JSType, List<FunctionType.Parameter>)
func (r *JSTypeRegistry) CreateFunctionTypeWithParameters(returnType JSType, parameters []*FunctionType_Parameter) FunctionType {
	return FunctionType_Builder_(r).WithParametersOfList(parameters).WithReturnType(returnType).Build()
}

// CreateFunctionTypeWithVarArgs creates a function type. The last parameter
// type of the function is considered a variable length argument.
//
// Java: JSTypeRegistry.createFunctionTypeWithVarArgs(JSType, JSType...)
func (r *JSTypeRegistry) CreateFunctionTypeWithVarArgs(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateFunctionTypeWithParameters(returnType, r.CreateParametersWithVarArgs(parameterTypes...))
}

// CreateNativeFunctionTypeWithVarArgs creates a function type. The last
// parameter type of the function is considered a variable length argument.
//
// Java: JSTypeRegistry.createNativeFunctionTypeWithVarArgs(JSType, JSType...)
func (r *JSTypeRegistry) CreateNativeFunctionTypeWithVarArgs(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateNativeFunctionType(returnType, r.CreateParametersWithVarArgs(parameterTypes...))
}

// CreateFunctionTypeWithInstanceType creates a function type in which this
// refers to an object instance.
//
// instanceType: the type of this. returnType: the function's return type.
// parameterTypes: the parameters' types.
//
// Java: JSTypeRegistry.createFunctionTypeWithInstanceType(ObjectType, JSType, List<JSType>)
func (r *JSTypeRegistry) CreateFunctionTypeWithInstanceType(
	instanceType ObjectType, returnType JSType, parameterTypes []JSType) JSType {
	paramsNode := r.CreateParameters(parameterTypes...)
	return FunctionType_Builder_(r).
		WithParametersOfList(paramsNode).
		WithReturnType(returnType).
		WithTypeOfThis(instanceType).
		Build()
}

// CreateParameters creates a tree hierarchy representing a typed argument
// list.
//
// Java: JSTypeRegistry.createParameters(JSType...)
func (r *JSTypeRegistry) CreateParameters(parameterTypes ...JSType) []*FunctionType_Parameter {
	return r.CreateParametersWithLastVarArgs(false, parameterTypes...)
}

// CreateParametersWithLastVarArgs creates a tree hierarchy representing a
// typed argument list.
//
// lastVarArgs: whether the last type should considered as a variable length
// argument. parameterTypes: the parameter types. The last element of this
// array is considered a variable length argument is lastVarArgs is true.
//
// Java: JSTypeRegistry.createParameters(boolean, JSType...)
func (r *JSTypeRegistry) CreateParametersWithLastVarArgs(lastVarArgs bool, parameterTypes ...JSType) []*FunctionType_Parameter {
	builder := NewFunctionParamBuilderWithCapacity(r, len(parameterTypes))
	max := len(parameterTypes) - 1
	for i := 0; i <= max; i++ {
		if lastVarArgs && i == max {
			builder.AddVarArgs(parameterTypes[i])
		} else {
			builder.AddRequiredParam(parameterTypes[i])
		}
	}
	return builder.Build()
}

// CreateParametersWithVarArgs creates a tree hierarchy representing a typed
// argument list. The last parameter type is considered a variable length
// argument.
//
// Java: JSTypeRegistry.createParametersWithVarArgs(JSType...)
func (r *JSTypeRegistry) CreateParametersWithVarArgs(parameterTypes ...JSType) []*FunctionType_Parameter {
	return r.CreateParametersWithLastVarArgs(true, parameterTypes...)
}

// CreateOptionalParameters creates a tree hierarchy representing a typed
// parameter list in which every parameter is optional.
//
// Java: JSTypeRegistry.createOptionalParameters(JSType...)
func (r *JSTypeRegistry) CreateOptionalParameters(parameterTypes ...JSType) []*FunctionType_Parameter {
	builder := NewFunctionParamBuilderWithCapacity(r, len(parameterTypes))
	builder.AddOptionalParams(parameterTypes...)
	return builder.Build()
}

// CreateFunctionTypeWithNewReturnType creates a new function type based on
// an existing function type but with a new return type.
//
// Java: JSTypeRegistry.createFunctionTypeWithNewReturnType(FunctionType, JSType)
func (r *JSTypeRegistry) CreateFunctionTypeWithNewReturnType(existingFunctionType FunctionType, returnType JSType) FunctionType {
	return existingFunctionType.ToBuilder().WithReturnType(returnType).Build()
}

// Java: JSTypeRegistry.createNativeFunctionType(JSType, List<FunctionType.Parameter>)
func (r *JSTypeRegistry) CreateNativeFunctionType(returnType JSType, parameters []*FunctionType_Parameter) FunctionType {
	return FunctionType_Builder_(r).
		WithParametersOfList(parameters).
		WithReturnType(returnType).
		ForNativeType().
		Build()
}

// Java: JSTypeRegistry.buildRecordTypeFromObject(ObjectType)
func (r *JSTypeRegistry) BuildRecordTypeFromObject(objType ObjectType) JSType {
	recType := objType.ToMaybeRecordType()
	// If it can be casted to a record type then return
	if recType != nil {
		return recType
	}
	// TODO(lpino): Handle inherited properties
	propNames := objType.GetOwnPropertyNames()
	// If the type has no properties then return Object
	if len(propNames) == 0 {
		return r.GetNativeType(JSTypeNative_OBJECT_TYPE)
	}
	var builder RecordTypeBuilder
	builder.init(r)
	// Otherwise collect the properties and build a record type
	for _, propName := range propNames {
		builder.AddProperty(propName, objType.GetPropertyType(propName), NilNode)
	}
	return builder.Build()
}

// Java: JSTypeRegistry.createRecordType(Map<String, ? extends JSType>)
func (r *JSTypeRegistry) CreateRecordType(props *jutil.LinkedHashMap[string, JSType]) JSType {
	var builder RecordTypeBuilder
	builder.init(r)
	for k, v := range props.All() {
		builder.AddProperty(k, v, NilNode)
	}
	return builder.Build()
}

// CreateObjectType creates an object type.
//
// Java: JSTypeRegistry.createObjectType(String, ObjectType)
func (r *JSTypeRegistry) CreateObjectType(name *string, implicitPrototype ObjectType) ObjectType {
	builder := PrototypeObjectType_Builder{
		registry:          r,
		classNameID:       rhinoStringIDFromPtr(name),
		implicitPrototype: implicitPrototype,
		templateTypeMap:   r.GetEmptyTemplateTypeMap(),
	}
	return NewPrototypeObjectType(&builder)
}

// CreateAnonymousObjectType creates an anonymous object type.
//
// info: used to mark object literals as structs; can be null
//
// Java: JSTypeRegistry.createAnonymousObjectType(JSDocInfo)
func (r *JSTypeRegistry) CreateAnonymousObjectType(info *JSDocInfo) ObjectType {
	builder := PrototypeObjectType_Builder{
		registry:        r,
		anonymousType:   true,
		templateTypeMap: r.GetEmptyTemplateTypeMap(),
	}
	type_ := NewPrototypeObjectType(&builder)
	type_.SetPrettyPrint(true)
	type_.SetJSDocInfo(info)
	return type_
}

// ResetImplicitPrototype sets the implicit prototype if it's possible to do
// so. There are a few different reasons why this could be a no-op: for
// example, numbers can't be implicit prototypes, and we don't want to change
// the implicit prototype if other classes have already subclassed this one.
//
// Java: JSTypeRegistry.resetImplicitPrototype(JSType, ObjectType)
func (r *JSTypeRegistry) ResetImplicitPrototype(type_ JSType, newImplicitProto ObjectType) {
	if poType, ok := type_.(PrototypeObjectType); ok {
		poType.ClearCachedValues()
		poType.SetImplicitPrototype(newImplicitProto)
	}
}

// CreateConstructorType creates a constructor function type.
//
// name: the function's name or null to indicate that the function is
// anonymous. source: the node defining this function. Its type
// (Node#getToken()) must be Token#FUNCTION. parameters: the function's
// parameters or null to indicate that the parameter types are unknown.
// returnType: the function's return type or null to indicate that the
// return type is unknown. templateKeys: the templatized types for the class.
// isAbstract: whether the function type represents an abstract class
//
// Java: JSTypeRegistry.createConstructorType(String, Node, List<FunctionType.Parameter>, JSType, ImmutableList<TemplateType>, boolean)
func (r *JSTypeRegistry) CreateConstructorType(
	name *string,
	source Node,
	parameters []*FunctionType_Parameter,
	returnType JSType,
	templateKeys []*TemplateType,
	isAbstract bool) FunctionType {
	jutil.CheckArgument(source.IsNil() || source.IsFunction() || source.IsClass())
	return FunctionType_Builder_(r).
		ForConstructor().
		WithName(name).
		WithSourceNode(source).
		WithParametersOfList(parameters).
		WithReturnType(returnType).
		WithTemplateKeysOfList(templateKeys).
		WithIsAbstract(isAbstract).
		Build()
}

// Java: JSTypeRegistry.createTemplateType(String)
func (r *JSTypeRegistry) CreateTemplateType(name string) *TemplateType {
	return NewTemplateType(r, name)
}

// Java: JSTypeRegistry.createTemplateType(String, JSType)
func (r *JSTypeRegistry) CreateTemplateTypeWithBound(name string, bound JSType) *TemplateType {
	return NewTemplateTypeWithBound(r, name, bound)
}

// Java: JSTypeRegistry.createTemplateTypeWithTransformation(String, Node)
func (r *JSTypeRegistry) CreateTemplateTypeWithTransformation(name string, expr Node) *TemplateType {
	return NewTemplateTypeWithTransformation(r, name, expr)
}

// Java: JSTypeRegistry.getEmptyTemplateTypeMap()
func (r *JSTypeRegistry) GetEmptyTemplateTypeMap() *TemplateTypeMap {
	return r.emptyTemplateTypeMap
}

// CreateTemplatizedTypeOfList creates a templatized instance of the
// specified type. Only ObjectTypes can currently be templatized; extend the
// logic in this function when more types can be templatized.
//
// baseType: the type to be templatized. templatizedTypes: a list of the
// template JSTypes. Will be matched by list order to the template keys on
// the base type.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, ImmutableList<JSType>)
func (r *JSTypeRegistry) CreateTemplatizedTypeOfList(baseType ObjectType, templatizedTypes []JSType) *TemplatizedType {
	jutil.CheckNotNull(baseType)
	// Only ObjectTypes can currently be templatized; extend this logic when
	// more types can be templatized.
	return NewTemplatizedType(r, baseType, templatizedTypes)
}

// CreateTemplatizedTypeOfMap creates a templatized instance of the
// specified type. Only ObjectTypes can currently be templatized; extend the
// logic in this function when more types can be templatized.
//
// baseType: the type to be templatized. templatizedTypes: a map from
// TemplateType to corresponding JSType value. Any unfilled TemplateTypes on
// the baseType that are *not* contained in this map will have UNKNOWN_TYPE
// used as their value.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, Map<TemplateType, JSType>)
func (r *JSTypeRegistry) CreateTemplatizedTypeOfMap(
	baseType ObjectType, templatizedTypes *jutil.LinkedHashMap[*TemplateType, JSType]) *TemplatizedType {
	unknownType := r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	return r.CreateTemplatizedTypeOfList(
		baseType,
		JSTypeIterations_MapTypes(
			func(key *TemplateType) JSType { return templatizedTypes.GetOrDefault(key, unknownType) },
			baseType.GetTypeParameters()))
}

// CreateTemplatizedType creates a templatized instance of the specified
// type. Only ObjectTypes can currently be templatized; extend the logic in
// this function when more types can be templatized.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, JSType...)
func (r *JSTypeRegistry) CreateTemplatizedType(baseType ObjectType, templatizedTypes ...JSType) *TemplatizedType {
	for i, e := range templatizedTypes {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
	}
	return r.CreateTemplatizedTypeOfList(baseType, templatizedTypes)
}

// CreateNamedType creates a named type.
//
// Java: JSTypeRegistry.createNamedType(StaticTypedScope, String, String, int, int)
func (r *JSTypeRegistry) CreateNamedType(
	scope StaticTypedScope, reference string, sourceName *string, lineno int, charno int) *NamedType {
	return r.createNamedTypeByID(scope, reference, rhinoStringIDFromPtr(sourceName), lineno, charno)
}

func (r *JSTypeRegistry) createNamedTypeByID(
	scope StaticTypedScope, reference string, sourceNameID rhinoStringID, lineno int, charno int) *NamedType {
	b := NamedType_Builder{
		registry:          r,
		referenceName:     reference,
		hasReferenceName:  true,
		scope:             scope,
		sourceNameID:      sourceNameID,
		lineno:            lineno,
		charno:            charno,
		referencedType:    r.GetNativeType(JSTypeNative_UNKNOWN_TYPE),
		templateTypes:     emptyJSTypeSlice,
		resolutionKind:    NamedType_ResolutionKind_TYPE_NAME,
		hasResolutionKind: true,
	}
	return NewNamedType(&b)
}

// IdentifyNonNullableName identifies the name of a typedef or enum before
// we actually declare it.
//
// Java: JSTypeRegistry.identifyNonNullableName(StaticScope, String)
func (r *JSTypeRegistry) IdentifyNonNullableName(scope StaticScope, name string) {
	lookupScope := JSTypeRegistry_GetLookupScope(scope, name)
	rootNode := r.GetRootNodeForScope(lookupScope)
	names := r.nonNullableTypeNames.Get(rootNode)
	if names == nil {
		names = jutil.NewLinkedHashSet[string]()
		r.nonNullableTypeNames.Put(rootNode, names)
	}
	names.Add(name)
}

// IsNonNullableName identifies the name of a typedef or enum before we
// actually declare it.
//
// Java: JSTypeRegistry.isNonNullableName(StaticScope, String)
func (r *JSTypeRegistry) IsNonNullableName(scope StaticScope, name string) bool {
	scope = JSTypeRegistry_GetLookupScope(scope, name)
	names := r.nonNullableTypeNames.Get(r.GetRootNodeForScope(scope))
	return names != nil && names.Contains(name)
}

// Java: JSTypeRegistry.evaluateTypeExpression(JSTypeExpression, StaticTypedScope)
func (r *JSTypeRegistry) EvaluateTypeExpression(expr *JSTypeExpression, scope StaticTypedScope) JSType {
	return r.CreateTypeFromCommentNodeWithScope(expr.GetRoot(), expr.GetSourceName(), scope)
}

// Java: JSTypeRegistry.createTypeFromCommentNode(Node)
func (r *JSTypeRegistry) CreateTypeFromCommentNode(n Node) JSType {
	return r.CreateTypeFromCommentNodeWithScope(n, "[internal]", nil)
}

// CreateTypeFromCommentNodeWithScope creates a JSType from the nodes
// representing a type.
//
// n: the node with type info. sourceName: the source file name. scope: a
// scope for doing type name lookups.
//
// Java: JSTypeRegistry.createTypeFromCommentNode(Node, String, StaticTypedScope)
func (r *JSTypeRegistry) CreateTypeFromCommentNodeWithScope(n Node, sourceName string, scope StaticTypedScope) JSType {
	switch n.GetToken() {
	case Token_LC:
		// Record type.
		return r.CreateRecordTypeFromNodes(n.GetFirstChild(), sourceName, scope)
	case Token_BANG:
		// Not nullable
		child := r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope)
		if namedType, ok := child.(*NamedType); ok {
			return namedType.GetBangType()
		}
		return child.RestrictByNotNullOrUndefined()
	case Token_QMARK:
		// Nullable or unknown
		firstChild := n.GetFirstChild()
		if firstChild.IsNil() {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return r.CreateNullableType(r.CreateTypeFromCommentNodeWithScope(firstChild, sourceName, scope))
	case Token_EQUALS:
		// Optional
		// TODO(b/117162687): stop automatically converting {string=} to
		// {(string|undefined)]}
		return r.CreateOptionalType(r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope))
	case Token_ITER_REST:
		// Var args
		return r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope)
	case Token_STAR:
		// The AllType
		return r.GetNativeType(JSTypeNative_ALL_TYPE)
	case Token_PIPE:
		// Union type
		var buf [8]JSType
		builder := buf[:0]
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			builder = append(builder, jutil.CheckNotNull(r.CreateTypeFromCommentNodeWithScope(child, sourceName, scope)))
		}
		return r.CreateUnionTypeOfList(builder)
	case Token_EMPTY:
		// When the return value of a function is not specified
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	case Token_VOID:
		// Only allowed in the return value of a function.
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	case Token_TYPEOF:
		name := n.GetFirstChild().GetString()
		// TODO(sdh): require var to be const?
		qname := QualifiedName_Of(name)
		root := qname.GetRoot()
		declarationScope := scope.GetTopmostScopeOfEventualDeclaration(root)
		rootSlot := scope.GetSlot(root)
		type_ := scope.LookupQualifiedName(qname)
		if type_ == nil || type_.IsUnknownType() || rootSlot.GetScope() != declarationScope {
			// Create a NamedType via getType so that it will be added to the list
			// of types to eventually resolve if necessary.
			return NamedType_Builder_(r, "typeof "+name).
				SetScope(scope).
				SetResolutionKind(NamedType_ResolutionKind_TYPEOF).
				SetErrorReportingLocationFrom(n).
				Build()
		}
		if type_.IsLiteralObject() {
			scopeType := type_
			type_ = NamedType_Builder_(r, "typeof "+name).
				SetResolutionKind(NamedType_ResolutionKind_NONE).
				SetReferencedType(scopeType).
				Build()
		}
		return type_
	case Token_STRINGLIT:
		nominalType := r.getTypeWithLocationString(scope, n.GetString(), sourceName, n.GetLineno(), n.GetCharno())
		templateArgs := r.ParseTemplateArgs(nominalType, n, sourceName, scope)

		// Handle forward declared types
		if nominalType.IsNamedType() && !nominalType.IsResolved() {
			if templateArgs != nil {
				nominalType = nominalType.ToMaybeNamedType().ToBuilder().SetTemplateTypes(templateArgs).Build()
			}
			return r.AddNullabilityBasedOnParseContext(n, nominalType, scope)
		}

		objectType, ok := nominalType.(ObjectType)
		if !ok {
			return nominalType
		}

		if templateArgs == nil || !nominalType.IsRawTypeOfTemplatizedType() {
			// TODO(nickreid): This case leaves template parameters unbound if users
			// fail to specify any arguments, allowing raw types to leak into
			// programs.
			return r.AddNullabilityBasedOnParseContext(n, nominalType, scope)
		}

		return r.AddNullabilityBasedOnParseContext(
			n, r.CreateTemplatizedTypeOfList(objectType, templateArgs), scope)
	case Token_FUNCTION:
		var thisType JSType
		isConstructor := false
		current := n.GetFirstChild()
		if current.IsThis() || current.IsNew() {
			contextNode := current.GetFirstChild()

			candidateThisType := r.CreateTypeFromCommentNodeWithScope(contextNode, sourceName, scope)

			// Allow null/undefined 'this' types to indicate that the function is
			// not called in a deliberate context, and 'this' access should raise
			// warnings.
			if candidateThisType.IsNullType() || candidateThisType.IsVoidType() {
				thisType = candidateThisType
			} else if current.IsThis() {
				thisType = candidateThisType.RestrictByNotNullOrUndefined()
			} else if current.IsNew() {
				thisType = jsTypeOrNil(ObjectType_Cast(candidateThisType.RestrictByNotNullOrUndefined()))
				if thisType == nil {
					r.reporter.Warning(
						Msg_JSDOC_FUNCTION_NEWNOTOBJECT.Format(),
						sourceName,
						contextNode.GetLineno(),
						contextNode.GetCharno())
				}
			}

			isConstructor = current.GetToken() == Token_NEW
			current = current.GetNext()
		}

		paramBuilder := NewFunctionParamBuilder(r)

		if current.GetToken() == Token_PARAM_LIST {
			for arg := current.GetFirstChild(); !arg.IsNil(); arg = arg.GetNext() {
				if arg.GetToken() == Token_ITER_REST {
					if !arg.HasChildren() {
						paramBuilder.AddVarArgs(r.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
					} else {
						paramBuilder.AddVarArgs(
							r.CreateTypeFromCommentNodeWithScope(arg.GetFirstChild(), sourceName, scope))
					}
				} else {
					type_ := r.CreateTypeFromCommentNodeWithScope(arg, sourceName, scope)
					if arg.GetToken() == Token_EQUALS {
						addSuccess := paramBuilder.AddOptionalParam(type_)
						if !addSuccess {
							r.reporter.Warning(
								Msg_JSDOC_FUNCTION_VARARGS.Format(),
								sourceName,
								arg.GetLineno(),
								arg.GetCharno())
						}
					} else {
						paramBuilder.AddRequiredParam(type_)
					}
				}
			}
			current = current.GetNext()
		}

		returnType := r.CreateTypeFromCommentNodeWithScope(current, sourceName, scope)

		kind := FunctionType_Kind_ORDINARY
		if isConstructor {
			kind = FunctionType_Kind_CONSTRUCTOR
		}
		return FunctionType_Builder_(r).
			WithParams(paramBuilder).
			WithReturnType(returnType).
			WithTypeOfThis(thisType).
			WithKind(kind).
			Build()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node in type expression: " + n.String()))
}

// jsTypeStaticScopeOrNil converts a StaticTypedScope to a StaticScope,
// keeping Java null.
func jsTypeStaticScopeOrNil(scope StaticTypedScope) StaticScope {
	if scope == nil {
		return nil
	}
	return scope
}

// Java: JSTypeRegistry.addNullabilityBasedOnParseContext(Node, JSType, StaticScope)
func (r *JSTypeRegistry) AddNullabilityBasedOnParseContext(n Node, type_ JSType, scope StaticTypedScope) JSType {
	// Other node types may be appropriate in the future.
	jutil.CheckStateMsg1(n.IsName() || n.IsStringLit(), "%s", n)
	jutil.CheckNotNull(type_)

	if r.IsNonNullableName(jsTypeStaticScopeOrNil(scope), n.GetString()) {
		return type_
	} else if type_.IsTemplateType() {
		// Template types represent the substituted type exactly and should not
		// be wrapped.
		return type_
	} else if n.HasParent() && n.GetParent().GetToken() == Token_BANG {
		// Names parsed from beneath a BANG never need nullability added.
		return type_
	} else {
		return r.CreateNullableType(type_)
	}
}

// Java: JSTypeRegistry.parseTemplateArgs(JSType, Node, String, StaticTypedScope)
func (r *JSTypeRegistry) ParseTemplateArgs(
	nominalType JSType, typeNode Node, sourceName string, scope StaticTypedScope) []JSType {
	typeList := typeNode.GetFirstChild()
	if typeList.IsNil() {
		return nil
	}

	count := 0
	for templateNode := typeList.GetFirstChild(); !templateNode.IsNil(); templateNode = templateNode.GetNext() {
		count++
	}
	// TODO(b/138617950): Eliminate the special case for `Object`.
	name := typeNode.GetString()
	isObject := name == "Object" || name == "window.Object"
	capCount := count
	if isObject && count == 1 {
		capCount = 2
	}
	templateArgs := make([]JSType, 0, capCount)
	if isObject && count == 1 {
		templateArgs = append(templateArgs, r.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
	}
	for templateNode := typeList.GetFirstChild(); !templateNode.IsNil(); templateNode = templateNode.GetNext() {
		templateArgs = append(templateArgs, jutil.CheckNotNull(r.CreateTypeFromCommentNodeWithScope(templateNode, sourceName, scope)))
	}

	if nominalType.IsNamedType() && !nominalType.IsResolved() {
		// The required number of template args will not be known until
		// resolution, so just return all the args.
		return templateArgs
	}

	requiredTemplateArgCount := nominalType.GetTemplateParamCount()
	if len(templateArgs) <= requiredTemplateArgCount {
		return templateArgs
	}

	if !nominalType.IsUnknownType() &&
		// TODO(b/287880204): delete the following case after cleaning up the
		// codebase
		!r.IsNonNullableName(jsTypeStaticScopeOrNil(scope), name) {
		firstExtraTemplateParam := typeNode.GetFirstChild().GetChildAtIndex(requiredTemplateArgCount)
		message := "Too many template parameters\nFound " +
			jutil.ToJavaString(len(templateArgs)) +
			", required at most " +
			jutil.ToJavaString(requiredTemplateArgCount)
		r.reporter.Warning(
			message,
			sourceName,
			firstExtraTemplateParam.GetLineno(),
			firstExtraTemplateParam.GetCharno())
	}
	clear(templateArgs[requiredTemplateArgCount:])
	return templateArgs[:requiredTemplateArgCount:requiredTemplateArgCount]
}

// CreateRecordTypeFromNodes creates a RecordType from the nodes representing
// said record type.
//
// n: the node with type info. sourceName: the source file name. scope: a
// scope for doing type name lookups.
//
// Java: JSTypeRegistry.createRecordTypeFromNodes(Node, String, StaticTypedScope)
func (r *JSTypeRegistry) CreateRecordTypeFromNodes(n Node, sourceName string, scope StaticTypedScope) JSType {
	var builder RecordTypeBuilder
	builder.init(r)

	// For each of the fields in the record type.
	for fieldTypeNode := n.GetFirstChild(); !fieldTypeNode.IsNil(); fieldTypeNode = fieldTypeNode.GetNext() {

		// Get the property's name.
		fieldNameNode := fieldTypeNode
		hasType := false

		if fieldTypeNode.GetToken() == Token_COLON {
			fieldNameNode = fieldTypeNode.GetFirstChild()
			hasType = true
		}

		fieldName := fieldNameNode.GetString()

		// TODO(user): Move this into the lexer/parser. Remove the string literal
		// characters around a field name, if any.
		if jutil.StringStartsWith(fieldName, "'") || jutil.StringStartsWith(fieldName, "\"") {
			fieldName = jutil.StringSubstring(fieldName, 1, jutil.StringLength(fieldName)-1)
		}

		// Get the property's type.
		var fieldType JSType

		if hasType {
			// We have a declared type.
			fieldType = r.CreateTypeFromCommentNodeWithScope(fieldTypeNode.GetLastChild(), sourceName, scope)
		} else {
			// Otherwise, the type is UNKNOWN.
			fieldType = r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}

		builder.AddProperty(fieldName, fieldType, fieldNameNode)
	}

	return builder.Build()
}

// RegisterTemplateTypeNamesInScope registers template types on the given
// scope root. This takes a Node rather than a StaticScope because at the time
// it is called, the scope has not yet been created.
//
// Java: JSTypeRegistry.registerTemplateTypeNamesInScope(Iterable<TemplateType>, Node)
func (r *JSTypeRegistry) RegisterTemplateTypeNamesInScope(keys []*TemplateType, scopeRoot Node) {
	for _, key := range keys {
		r.scopedNameTablePut(scopeRoot, *key.GetReferenceName(), key)
	}
}

// CreateScopeWithTemplates returns a new scope that includes the given
// template names for type resolution purposes.
//
// Java: JSTypeRegistry.createScopeWithTemplates(StaticTypedScope, Iterable<TemplateType>)
func (r *JSTypeRegistry) CreateScopeWithTemplates(scope StaticTypedScope, templates []*TemplateType) StaticTypedScope {
	return NewJSTypeRegistry_SyntheticTemplateScope(scope, templates)
}

// JSTypeRegistry_SyntheticTemplateScope is a synthetic scope that includes
// template names. This is necessary for resolving template names outside the
// body of templated functions (e.g. when evaluating JSDoc on things assigned
// to a prototype, or the parameter or return types of an annotated
// function), since there is not yet (and may never be) any real scope to
// attach the types to.
type JSTypeRegistry_SyntheticTemplateScope struct {
	delegate StaticTypedScope
	types    *HamtPMap[string, *TemplateType]
}

// Java: JSTypeRegistry.SyntheticTemplateScope(StaticTypedScope, Iterable<TemplateType>)
func NewJSTypeRegistry_SyntheticTemplateScope(delegate StaticTypedScope, templates []*TemplateType) *JSTypeRegistry_SyntheticTemplateScope {
	s := &JSTypeRegistry_SyntheticTemplateScope{}
	s.delegate = delegate
	var types *HamtPMap[string, *TemplateType]
	if syntheticTemplateScope, ok := delegate.(*JSTypeRegistry_SyntheticTemplateScope); ok {
		types = syntheticTemplateScope.types
	}
	for _, key := range templates {
		types = types.Plus(*key.GetReferenceName(), key)
	}
	s.types = types
	return s
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getRootNode()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetRootNode() Node {
	return s.delegate.GetRootNode()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getParentScope()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetParentScope() StaticScope {
	return s.delegate.GetParentScope()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getSlot(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetSlot(name string) StaticSlot {
	return s.delegate.GetSlot(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getOwnSlot(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetOwnSlot(name string) StaticSlot {
	return s.delegate.GetOwnSlot(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTypeOfThis()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTypeOfThis() JSType {
	return s.delegate.GetTypeOfThis()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTemplateType(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTemplateType(name string) *TemplateType {
	return s.types.Get(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTopmostScopeOfEventualDeclaration(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTopmostScopeOfEventualDeclaration(name string) StaticScope {
	if s.types.Get(name) != nil {
		return s
	}
	return s.delegate.GetTopmostScopeOfEventualDeclaration(name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *JSTypeRegistry_SyntheticTemplateScope) LookupQualifiedName(qname QualifiedName) JSType {
	return StaticTypedScope_LookupQualifiedName(s, qname)
}

// Port of com/google/javascript/jscomp/FunctionTypeBuilder.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	emptyTemplateTypeSlice       = []*rhino.TemplateType{}
	emptyFunctionParametersSlice = []*rhino.FunctionType_Parameter{}
)

// FunctionTypeBuilder is a builder for FunctionTypes, because FunctionTypes
// are so ridiculously complex. All methods return this for ease of use.
//
// Right now, this mostly uses JSDocInfo to infer type information about
// functions. In the long term, developers should extend it to use other
// signals by overloading the various "inferXXX" methods. For example, we
// might want to use goog.inherits calls as a signal for inheritance, or
// return statements as a signal for return type.
//
// NOTE(nicksantos): Organizationally, this feels like it should be in Rhino.
// But it depends on some coding convention stuff that's really part of
// JSCompiler.
type FunctionTypeBuilder struct {
	fnName           string
	compiler         AbstractCompiler
	codingConvention CodingConvention
	typeRegistry     *rhino.JSTypeRegistry
	errorRoot        rhino.Node

	sourceNode         rhino.Node
	hasNonEmptyReturns bool

	syntacticFnName              string
	hasSyntacticFnName           bool
	returnType                   rhino.JSType
	returnTypeInferred           bool
	implementedInterfaces        []rhino.ObjectType
	extendedInterfaces           []rhino.ObjectType
	baseType                     rhino.ObjectType
	thisType                     rhino.JSType
	isClass                      bool
	isConstructor                bool
	makesStructs                 bool
	makesUnrestricted            bool
	makesDicts                   bool
	isInterface                  bool
	isRecord                     bool
	isAbstract                   bool
	isKnownAmbiguous             bool
	parameters                   []*rhino.FunctionType_Parameter
	closurePrimitiveId           *rhino.ClosurePrimitive
	templateTypeNames            []*rhino.TemplateType
	constructorTemplateTypeNames []*rhino.TemplateType
	declarationScope             *TypedScope
	templateScope                rhino.StaticTypedScope
}

var FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF = DiagnosticType_Warning(
	"JSC_EXTENDS_WITHOUT_TYPEDEF",
	"@extends used without @constructor or @interface for {0}")

var FunctionTypeBuilder_EXTENDS_NON_OBJECT = DiagnosticType_Warning("JSC_EXTENDS_NON_OBJECT", "{0} @extends non-object type {1}")

var FunctionTypeBuilder_RESOLVED_TAG_EMPTY = DiagnosticType_Warning("JSC_RESOLVED_TAG_EMPTY", "Could not resolve type in {0} tag of {1}")

var FunctionTypeBuilder_CONSTRUCTOR_REQUIRED = DiagnosticType_Warning("JSC_CONSTRUCTOR_REQUIRED", "{0} used without @constructor for {1}")

var FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST = DiagnosticType_Warning("JSC_VAR_ARGS_MUST_BE_LAST", "variable length argument must be last")

var FunctionTypeBuilder_OPTIONAL_ARG_AT_END = DiagnosticType_Warning("JSC_OPTIONAL_ARG_AT_END", "optional arguments must be at the end")

var FunctionTypeBuilder_INEXISTENT_PARAM = DiagnosticType_Warning(
	"JSC_INEXISTENT_PARAM", "parameter {0} does not appear in {1}''s parameter list")

var FunctionTypeBuilder_TYPE_REDEFINITION = DiagnosticType_Warning(
	"JSC_TYPE_REDEFINITION",
	"attempted re-definition of type {0}\n"+
		"found   : {1}\n"+
		"expected: {2}")

var FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS = DiagnosticType_Warning(
	"JSC_TEMPLATE_TRANSFORMATION_ON_CLASS",
	"Template type transformation {0} not allowed on classes or interfaces")

var FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND = DiagnosticType_Error(
	"JSC_TEMPLATE_TYPE_ILLEGAL_BOUND",
	"Illegal upper bound ''{0}'' on template type parameter {1}")

var FunctionTypeBuilder_ALL_DIAGNOSTICS = NewDiagnosticGroup(
	FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF,
	FunctionTypeBuilder_EXTENDS_NON_OBJECT,
	FunctionTypeBuilder_RESOLVED_TAG_EMPTY,
	FunctionTypeBuilder_CONSTRUCTOR_REQUIRED,
	FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST,
	FunctionTypeBuilder_OPTIONAL_ARG_AT_END,
	FunctionTypeBuilder_INEXISTENT_PARAM,
	FunctionTypeBuilder_TYPE_REDEFINITION,
	FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS,
	FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND,
	TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS)

// CreateExtendedTypeValidator returns the Predicate<JSType> that validates
// @extends types without retaining the surrounding FunctionTypeBuilder.
//
// Java: FunctionTypeBuilder.createExtendedTypeValidator()
func (f *FunctionTypeBuilder) CreateExtendedTypeValidator() func(rhino.JSType) bool {
	errorRoot := f.errorRoot
	compiler := f.compiler
	formattedFnName := f.FormatFnName()
	return func(type_ rhino.JSType) bool {
		objectType := rhino.ObjectType_Cast(type_)
		if objectType == nil {
			compiler.Report(JSError_MakeWithNode(errorRoot, FunctionTypeBuilder_EXTENDS_NON_OBJECT, formattedFnName, type_.String()))
			return false
		}
		if objectType.IsEmptyType() {
			compiler.Report(JSError_MakeWithNode(errorRoot, FunctionTypeBuilder_RESOLVED_TAG_EMPTY, "@extends", formattedFnName))
			return false
		}
		if objectType.IsUnknownType() {
			if FunctionTypeBuilder_HasMoreTagsToResolve(objectType) || type_.IsTemplateType() {
				return true
			}
			compiler.Report(JSError_MakeWithNode(errorRoot, FunctionTypeBuilder_RESOLVED_TAG_EMPTY, "@extends", formattedFnName))
			return false
		}
		return true
	}
}

// CreateImplementedTypeValidator returns the Predicate<JSType> that validates
// @implements types without retaining the surrounding FunctionTypeBuilder.
//
// Java: FunctionTypeBuilder.createImplementedTypeValidator()
func (f *FunctionTypeBuilder) CreateImplementedTypeValidator() func(rhino.JSType) bool {
	errorRoot := f.errorRoot
	compiler := f.compiler
	formattedFnName := f.FormatFnName()
	return func(type_ rhino.JSType) bool {
		objectType := rhino.ObjectType_Cast(type_)
		if objectType == nil {
			compiler.Report(JSError_MakeWithNode(errorRoot, TypeCheck_BAD_IMPLEMENTED_TYPE, formattedFnName))
			return false
		}
		if objectType.IsEmptyType() {
			compiler.Report(JSError_MakeWithNode(errorRoot, FunctionTypeBuilder_RESOLVED_TAG_EMPTY, "@implements", formattedFnName))
			return false
		}
		if objectType.IsUnknownType() {
			if FunctionTypeBuilder_HasMoreTagsToResolve(objectType) {
				return true
			}
			compiler.Report(JSError_MakeWithNode(errorRoot, FunctionTypeBuilder_RESOLVED_TAG_EMPTY, "@implements", formattedFnName))
			return false
		}
		return true
	}
}

// InitFunctionTypeBuilder initializes a stack-allocated or reused FunctionTypeBuilder.
func (f *FunctionTypeBuilder) InitFunctionTypeBuilder(fnName *string, compiler AbstractCompiler, errorRoot rhino.Node, scope *TypedScope) {
	rhino.CheckNotNull(errorRoot)
	*f = FunctionTypeBuilder{
		templateTypeNames:            emptyTemplateTypeSlice,
		constructorTemplateTypeNames: emptyTemplateTypeSlice,
		codingConvention:             compiler.GetCodingConvention(),
		typeRegistry:                 compiler.GetTypeRegistry(),
		errorRoot:                    errorRoot,
		compiler:                     compiler,
	}
	if fnName != nil {
		f.fnName = *fnName
	}
	// Go-only: language-adapter — a nil *TypedScope must stay a nil interface.
	if scope != nil {
		f.templateScope = scope
	}
}

// NewFunctionTypeBuilder creates a builder.
//
// fnName is the function name to be used in error messages; compiler the
// compiler; errorRoot the node to associate with any warning generated by
// this builder; scope the syntactic scope.
//
// Java: FunctionTypeBuilder(String, AbstractCompiler, Node, TypedScope)
func NewFunctionTypeBuilder(fnName *string, compiler AbstractCompiler, errorRoot rhino.Node, scope *TypedScope) *FunctionTypeBuilder {
	f := &FunctionTypeBuilder{}
	f.InitFunctionTypeBuilder(fnName, compiler, errorRoot, scope)
	return f
}

// functionTypeBuilder_staticScope passes a possibly-null TypedScope where
// Java passes it as a StaticScope.
//
// Go-only: language-adapter — a nil *TypedScope must become a nil
// interface, not a typed nil.
func functionTypeBuilder_staticScope(scope *TypedScope) rhino.StaticScope {
	if scope == nil {
		return nil
	}
	return scope
}

// FormatFnName formats the function name for use in warnings.
//
// Java: FunctionTypeBuilder.formatFnName()
func (f *FunctionTypeBuilder) FormatFnName() string {
	if f.fnName == "" {
		return "<anonymous>"
	}
	return f.fnName
}

// SetSyntacticFunctionName sets the name with which this new type will be
// declared in the type registry.
//
// Java: FunctionTypeBuilder.setSyntacticFunctionName(String)
func (f *FunctionTypeBuilder) SetSyntacticFunctionName(syntacticFnName *string) *FunctionTypeBuilder {
	if syntacticFnName != nil {
		f.syntacticFnName = *syntacticFnName
	} else {
		f.syntacticFnName = ""
	}
	f.hasSyntacticFnName = true
	return f
}

// SetAstContents sets the AST contents of this function without allocating a separate struct.
func (f *FunctionTypeBuilder) SetAstContents(n rhino.Node, hasNonEmptyReturns bool) *FunctionTypeBuilder {
	if !n.IsNil() {
		f.sourceNode = n
		f.hasNonEmptyReturns = hasNonEmptyReturns
	}
	return f
}

// SetDeclarationScope sets a declaration scope explicitly. This is important
// with block scopes because a function declared in an inner scope with 'var'
// needs to use the inner scope to resolve names, but needs to be declared in
// the outer scope.
//
// Java: FunctionTypeBuilder.setDeclarationScope(TypedScope)
func (f *FunctionTypeBuilder) SetDeclarationScope(declarationScope *TypedScope) *FunctionTypeBuilder {
	f.declarationScope = declarationScope
	return f
}

// InferReturnType infers the return type from JSDocInfo.
//
// fromInlineDoc indicates whether return type is inferred from inline doc
// attached to function name.
//
// Java: FunctionTypeBuilder.inferReturnType(JSDocInfo, boolean)
func (f *FunctionTypeBuilder) InferReturnType(info *rhino.JSDocInfo, fromInlineDoc bool) *FunctionTypeBuilder {
	if info != nil {
		var returnTypeExpr *rhino.JSTypeExpression
		if fromInlineDoc {
			returnTypeExpr = info.GetType()
		} else {
			returnTypeExpr = info.GetReturnType()
		}
		if returnTypeExpr != nil {
			f.returnType = returnTypeExpr.Evaluate(f.templateScope, f.typeRegistry)
			f.returnTypeInferred = false
		}
	}

	return f
}

// Java: FunctionTypeBuilder.usingClassSyntax()
func (f *FunctionTypeBuilder) UsingClassSyntax() *FunctionTypeBuilder {
	f.isClass = true
	return f
}

// InferKind infers whether the function is a normal function, a
// constructor, or an interface.
//
// Java: FunctionTypeBuilder.inferKind(JSDocInfo)
func (f *FunctionTypeBuilder) InferKind(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if info != nil {
		if !NodeUtil_IsMethodDeclaration(f.errorRoot) {
			f.isConstructor = info.IsConstructor()
			f.isInterface = info.IsInterface()
			f.isRecord = info.UsesImplicitMatch()
			f.makesStructs = info.MakesStructs()
			f.makesUnrestricted = info.MakesUnrestricted()
			f.makesDicts = info.MakesDicts()
		}
		f.isAbstract = info.IsAbstract()
	}
	if f.isClass {
		// If a CLASS literal has not been explicitly declared an interface, it's a constructor.
		// If it's not expicitly @dict or @unrestricted then it's @struct.
		f.isConstructor = !f.isInterface
		f.makesStructs = info == nil || (!f.makesDicts && !info.MakesUnrestricted())
	}

	if f.makesStructs && !(f.isConstructor || f.isInterface) {
		f.ReportWarning(FunctionTypeBuilder_CONSTRUCTOR_REQUIRED, "@struct", f.FormatFnName())
	} else if f.makesDicts && !f.isConstructor {
		f.ReportWarning(FunctionTypeBuilder_CONSTRUCTOR_REQUIRED, "@dict", f.FormatFnName())
	}
	return f
}

func jsDocTypeTransformations(info *rhino.JSDocInfo) *jutil.LinkedHashMap[string, rhino.Node] {
	if info == nil {
		return nil
	}
	if m, ok := rhino.JSDocInfo_TYPE_TRANSFORMATIONS.Get(info).(*jutil.LinkedHashMap[string, rhino.Node]); ok {
		return m
	}
	return nil
}

// MaybeUseNativeClassTemplateNames clobbers the templateTypeNames from the
// JSDoc with builtin ones for native types.
//
// Java: FunctionTypeBuilder.maybeUseNativeClassTemplateNames(JSDocInfo)
func (f *FunctionTypeBuilder) MaybeUseNativeClassTemplateNames(info *rhino.JSDocInfo) bool {
	if f.fnName == "" {
		return false
	}
	infoTemplateCount := len(info.GetTemplateTypeEntries())
	if infoTemplateCount == 0 && len(f.templateTypeNames) == 0 && jsDocTypeTransformations(info) == nil {
		return false
	}
	nativeKeys := f.typeRegistry.MaybeGetTemplateTypesOfBuiltin(
		functionTypeBuilder_staticScope(f.declarationScope), f.fnName)
	// TODO(b/73386087): Make infoTemplateTypeNames.size() == nativeKeys.size() a
	// Preconditions check. It currently fails for "var symbol" in the externs.
	if nativeKeys != nil && infoTemplateCount == len(nativeKeys) {
		f.templateTypeNames = nativeKeys
		return true
	}
	return false
}

// InferInheritance infers any supertypes from the JSDocInfo or the passed-in
// base type.
//
// info is the JSDoc info that is attached to the type declaration, if any;
// classExtendsType the type of the extends clause in
// `class C extends SuperClass {}`, if present.
//
// Java: FunctionTypeBuilder.inferInheritance(JSDocInfo, ObjectType)
func (f *FunctionTypeBuilder) InferInheritance(info *rhino.JSDocInfo, classExtendsType rhino.ObjectType) *FunctionTypeBuilder {

	if info != nil && info.HasBaseType() {
		if f.isConstructor || f.isInterface {
			infoBaseType := info.GetBaseType().Evaluate(f.templateScope, f.typeRegistry)
			if !FunctionTypeBuilder_AreCompatibleExtendsTypes(infoBaseType.ToMaybeObjectType(), classExtendsType) {
				f.isKnownAmbiguous = true
			}
			if infoBaseType.SetValidator(f.CreateExtendedTypeValidator()) {
				f.baseType = infoBaseType.ToObjectType()
			} else {
				f.isKnownAmbiguous = true
			}
		} else {
			f.ReportWarning(FunctionTypeBuilder_EXTENDS_WITHOUT_TYPEDEF, f.FormatFnName())
			f.isKnownAmbiguous = true
		}
	} else if classExtendsType != nil && (f.isConstructor || f.isInterface) {
		// This case is:
		// // no JSDoc here
		// class extends astBaseType {...}
		//
		// It may well be that astBaseType is something dynamically created, like a value passed into
		// a function. A common pattern is:
		//
		// function mixinX(superClass) {
		//   return class extends superClass {
		//     ...
		//   };
		// }
		// The ExtendedTypeValidator() used in the JSDocInfo case above will report errors for these
		// cases, and we don't want that.
		// Since astBaseType is an actual value in code rather than an annotation, we can
		// rely on validation elsewhere to ensure it is actually defined.
		f.baseType = classExtendsType
	}

	// Implemented interfaces (for constructors only).
	if info != nil && info.GetImplementedInterfaceCount() > 0 {
		if f.isConstructor {
			implExprs := info.GetImplementedInterfaces()
			f.implementedInterfaces = make([]rhino.ObjectType, 0, len(implExprs))
			for _, t := range implExprs {
				maybeInterType := t.Evaluate(f.templateScope, f.typeRegistry)

				if maybeInterType.SetValidator(f.CreateImplementedTypeValidator()) {
					f.implementedInterfaces = append(f.implementedInterfaces, maybeInterType.(rhino.ObjectType))
				}
			}
		} else if f.isInterface {
			f.ReportWarning(TypeCheck_CONFLICTING_IMPLEMENTED_TYPE, f.FormatFnName())
		} else {
			f.ReportWarning(FunctionTypeBuilder_CONSTRUCTOR_REQUIRED, "@implements", f.FormatFnName())
		}
	}

	// extended interfaces (for interfaces only)
	// We've already emitted a warning if this is not an interface.
	if f.isInterface {
		extCap := 0
		if info != nil {
			extCap += info.GetExtendedInterfacesCount()
		}
		if classExtendsType != nil {
			extCap++
		}
		if extCap > 0 {
			f.extendedInterfaces = make([]rhino.ObjectType, 0, extCap)
		} else {
			f.extendedInterfaces = []rhino.ObjectType{}
		}
		if info != nil {
			for _, t := range info.GetExtendedInterfaces() {
				maybeInterfaceType := t.Evaluate(f.templateScope, f.typeRegistry)
				if maybeInterfaceType != nil {
					// setValidator runs validation and returns whether validation was successful (except
					// for not-yet resolved named types, where validation is delayed).
					// This code must run even for non-object types (which we know are invalid) to generate
					// and record the user error message.
					isValid := maybeInterfaceType.SetValidator(f.CreateExtendedTypeValidator())
					// ExtendedTypeValidator guarantees that maybeInterfaceType is an object type, but
					// setValidator might not (e.g. due to delayed execution).
					if isValid && maybeInterfaceType.ToMaybeObjectType() != nil {
						f.extendedInterfaces = append(f.extendedInterfaces, maybeInterfaceType.ToMaybeObjectType())
					}
				}
				// de-dupe baseType (from extends keyword) if it's also in @extends jsdoc.
				if classExtendsType != nil && maybeInterfaceType.IsSubtypeOf(classExtendsType) {
					classExtendsType = nil
				}
			}
		}
		if classExtendsType != nil &&
			classExtendsType.SetValidator(f.CreateExtendedTypeValidator()) {
			// case is:
			// /**
			//  * @interface
			//  * @extends {OtherInterface}
			//  */
			// class SomeInterface extends astBaseType {}
			// Add the explicit extends type to the extended interfaces listed in JSDoc.
			f.extendedInterfaces = append(f.extendedInterfaces, classExtendsType)
		}
	}

	return f
}

// FunctionTypeBuilder_AreCompatibleExtendsTypes decides if the types in the
// @extends annotation and the extends clause of a class are a "sensible"
// combination.
//
// Sensible is vague here, depending on how dynamic we allow types to be.
// It's not just a supertype/subtype check. Generally, we want to trust user
// annotations as declarations of intent, but we also want to protect users
// from dangerous lies. UNKNOWN as one of the types is a particularly
// uncertain case.
//
// Java: FunctionTypeBuilder.areCompatibleExtendsTypes(ObjectType, ObjectType)
func FunctionTypeBuilder_AreCompatibleExtendsTypes(annotated rhino.ObjectType, extendsClause rhino.ObjectType) bool {
	if extendsClause == nil || annotated == extendsClause {
		return true
	}

	// Allow `/** @extends {Foo<T>} */ class Bar extends Foo`
	return annotated.IsTemplatizedType() &&
		annotated.ToMaybeTemplatizedType().GetReferencedType() == extendsClause
}

// InferThisTypeWithFallback infers the type of this; type is the type of
// this if the info is missing.
//
// Java: FunctionTypeBuilder.inferThisType(JSDocInfo, JSType)
func (f *FunctionTypeBuilder) InferThisTypeWithFallback(info *rhino.JSDocInfo, type_ rhino.JSType) *FunctionTypeBuilder {
	// Look at the @this annotation first.
	f.InferThisType(info)

	if f.thisType == nil {
		objType := rhino.ObjectType_Cast(type_)
		if objType != nil && (info == nil || !info.HasType()) {
			f.thisType = objType
		}
	}

	return f
}

// InferThisType infers the type of this from the JSDocInfo for this
// function.
//
// Java: FunctionTypeBuilder.inferThisType(JSDocInfo)
func (f *FunctionTypeBuilder) InferThisType(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if info != nil && info.HasThisType() {
		// TODO(johnlenz): In ES5 strict mode a function can have a null or
		// undefined "this" value, but all the existing "@this" annotations
		// don't declare restricted types.
		maybeThisType :=
			info.GetThisType().Evaluate(f.templateScope, f.typeRegistry).RestrictByNotNullOrUndefined()
		if maybeThisType != nil {
			f.thisType = maybeThisType
		}
	}

	return f
}

// Java: FunctionTypeBuilder.inferClosurePrimitive(JSDocInfo)
func (f *FunctionTypeBuilder) InferClosurePrimitive(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if info != nil && info.HasClosurePrimitiveId() {
		f.closurePrimitiveId = rhino.ClosurePrimitive_FromStringId(info.GetClosurePrimitiveId())
	}
	return f
}

// ProvideDefaultReturnType sets the returnType for this function using very
// basic type inference.
//
// Java: FunctionTypeBuilder.provideDefaultReturnType()
func (f *FunctionTypeBuilder) ProvideDefaultReturnType() {
	if !f.sourceNode.IsNil() {
		if f.sourceNode.IsAsyncGeneratorFunction() {
			// Set the return type of a generator function to:
			//   @return {!AsyncGenerator<?>}
			generatorType := f.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ASYNC_GENERATOR_TYPE)
			f.returnType = f.typeRegistry.CreateTemplatizedType(
				generatorType, f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
			return
		}
		if f.sourceNode.IsGeneratorFunction() {
			// Set the return type of a generator function to:
			//   @return {!Generator<?>}
			generatorType := f.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GENERATOR_TYPE)
			f.returnType = f.typeRegistry.CreateTemplatizedType(
				generatorType, f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
			return
		}
	}

	inferredReturnType := f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if !f.sourceNode.IsNil() &&
		!f.hasNonEmptyReturns &&
		!functionTypeBuilder_hasSingleThrow(f.sourceNode) &&
		!f.sourceNode.IsFromExterns() {
		// Infer return types for non-generator functions.
		// We need to be extremely conservative about this, because of two
		// competing needs.
		// 1) If we infer the return type of f too widely, then we won't be able
		//    to assign f to other functions.
		// 2) If we infer the return type of f too narrowly, then we won't be
		//    able to override f in subclasses.
		// So we only infer in cases where the user doesn't expect to write
		// @return annotations--when it's very obvious that the function returns
		// nothing.
		inferredReturnType = f.typeRegistry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		f.returnTypeInferred = true
	}

	if !f.sourceNode.IsNil() && f.sourceNode.IsAsyncFunction() {
		// Set the return type of an async function:
		//   @return {!Promise<?>} or @return {!Promise<undefined>}
		promiseType := f.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE)
		f.returnType = f.typeRegistry.CreateTemplatizedType(promiseType, inferredReturnType)
	} else {
		f.returnType = inferredReturnType
	}
}

// BuildAndRegister builds the function type, and puts it in the registry.
//
// Java: FunctionTypeBuilder.buildAndRegister()
func (f *FunctionTypeBuilder) BuildAndRegister() rhino.FunctionType {
	if f.returnType == nil {
		f.ProvideDefaultReturnType()
		jutil.CheckNotNull(f.returnType)
	}

	if f.parameters == nil {
		panic(jutil.NewIllegalStateExceptionMsg("All Function types must have params and a return type"))
	}

	var fnType rhino.FunctionType
	if f.isConstructor {
		fnType = f.GetOrCreateConstructor()
	} else if f.isInterface {
		fnType = f.GetOrCreateInterface()
	} else {
		fnType = f.CreateDefaultBuilder().
			WithParametersOfList(f.parameters).
			WithReturnTypeWithInferred(f.returnType, f.returnTypeInferred).
			WithTypeOfThis(f.thisType).
			WithIsAbstract(f.isAbstract).
			WithClosurePrimitiveId(f.closurePrimitiveId).
			Build()
		f.MaybeSetBaseType(fnType)
	}

	if f.implementedInterfaces != nil && fnType.IsConstructor() {
		fnType.SetImplementedInterfaces(f.implementedInterfaces)
	}

	if f.extendedInterfaces != nil {
		fnType.SetExtendedInterfaces(f.extendedInterfaces)
	}

	if f.isRecord {
		fnType.SetImplicitMatch(true)
	}

	return fnType
}

// Java: FunctionTypeBuilder.maybeSetBaseType(FunctionType)
func (f *FunctionTypeBuilder) MaybeSetBaseType(fnType rhino.FunctionType) {
	if !fnType.HasInstanceType() || f.baseType == nil {
		return
	}

	fnType.SetPrototypeBasedOn(f.baseType)
	fnType.GetInstanceType().MergeSupertypeTemplateTypes(f.baseType)
}

// Java: FunctionTypeBuilder.createDefaultBuilder()
func (f *FunctionTypeBuilder) CreateDefaultBuilder() *rhino.FunctionType_Builder {
	return rhino.FunctionType_Builder_(f.typeRegistry).
		WithNameString(f.fnName).
		WithSourceNode(f.sourceNode).
		WithTemplateKeysOfList(f.templateTypeNames).
		SetIsKnownAmbiguous(f.isKnownAmbiguous).
		SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(f.declarationScope))
}

// GetOrCreateConstructor returns a constructor function either by returning
// it from the registry if it exists or creating and registering a new type.
// If there is already a type, then warn if the existing type is different
// than the one we are creating, though still return the existing function
// if possible. The primary purpose of this is that registering a
// constructor will fail for all built-in types that are initialized in
// JSTypeRegistry. We a) want to make sure that the type information
// specified in the externs file matches what is in the registry and b)
// annotate the externs with the JSType from the registry so that there are
// not two separate JSType objects for one type.
//
// Java: FunctionTypeBuilder.getOrCreateConstructor()
func (f *FunctionTypeBuilder) GetOrCreateConstructor() rhino.FunctionType {
	fnType := f.CreateDefaultBuilder().
		ForConstructor().
		WithParametersOfList(f.parameters).
		WithReturnType(f.returnType).
		WithConstructorTemplateKeys(f.constructorTemplateTypeNames).
		WithIsAbstract(f.isAbstract).
		Build()

	if f.makesStructs {
		fnType.SetStruct()
	} else if f.makesDicts {
		fnType.SetDict()
	} else if f.makesUnrestricted {
		fnType.SetExplicitUnrestricted()
	}

	// There are two cases where this type already exists in the current scope:
	//   1. The type is a built-in that we initialized in JSTypeRegistry and is also defined in
	//  externs.
	//   2. Cases like "class C {} C = class {}"
	// See https://github.com/google/closure-compiler/issues/2928 for some related bugs.
	existingType := f.typeRegistry.GetType(functionTypeBuilder_staticScope(f.declarationScope), f.fnName)
	if existingType != nil {
		isInstanceObject := existingType.IsInstanceType()
		if isInstanceObject || f.fnName == "Function" {
			var existingFn rhino.FunctionType
			if isInstanceObject {
				existingFn = existingType.ToObjectType().GetConstructor()
			} else {
				existingFn = f.typeRegistry.GetNativeFunctionType(rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE)
			}

			if existingFn.GetSource().IsNil() {
				existingFn.SetSource(f.sourceNode)
			}

			if !existingFn.HasEqualCallType(fnType) {
				f.ReportWarning(
					FunctionTypeBuilder_TYPE_REDEFINITION, f.FormatFnName(), fnType.String(), existingFn.String())
			}

			return existingFn
		} else {
			// We fall through and return the created type, even though it will fail
			// to register. We have no choice as we have to return a function. We
			// issue an error elsewhere though, so the user should fix it.
		}
	}

	f.MaybeSetBaseType(fnType)

	// TODO(johnlenz): determine what we are supposed to do for:
	//   @constructor
	//   this.Foo = ...
	//
	if !f.hasSyntacticFnName {
		panic(jutil.NewNullPointerException())
	}
	if f.syntacticFnName != "" && !strings.HasPrefix(f.syntacticFnName, "this.") {
		f.typeRegistry.DeclareTypeForExactScope(
			functionTypeBuilder_staticScope(f.declarationScope), f.syntacticFnName, fnType.GetInstanceType())
	}
	return fnType
}

// Java: FunctionTypeBuilder.getOrCreateInterface()
func (f *FunctionTypeBuilder) GetOrCreateInterface() rhino.FunctionType {
	if !f.hasSyntacticFnName {
		panic(jutil.NewNullPointerException())
	}
	var fnType rhino.FunctionType

	type_ := f.typeRegistry.GetType(functionTypeBuilder_staticScope(f.declarationScope), f.syntacticFnName)
	if type_ != nil && type_.IsInstanceType() {
		ctor := type_.ToMaybeObjectType().GetConstructor()
		if ctor.IsInterface() {
			fnType = ctor
			fnType.SetSource(f.sourceNode)
		}
	}

	if fnType == nil {
		fnType = f.CreateDefaultBuilder().ForInterface().WithParameters().Build()
		if f.makesStructs {
			fnType.SetStruct()
		}
		if f.fnName != "" {
			f.typeRegistry.DeclareTypeForExactScope(
				functionTypeBuilder_staticScope(f.declarationScope), f.syntacticFnName, fnType.GetInstanceType())
		}
		f.MaybeSetBaseType(fnType)
	}
	return fnType
}

// Java: FunctionTypeBuilder.reportWarning(DiagnosticType, String...)
func (f *FunctionTypeBuilder) ReportWarning(warning *DiagnosticType, args ...string) {
	f.compiler.Report(JSError_MakeWithNode(f.errorRoot, warning, args...))
}

// Java: FunctionTypeBuilder.reportError(DiagnosticType, String...)
func (f *FunctionTypeBuilder) ReportError(error_ *DiagnosticType, args ...string) {
	f.compiler.Report(JSError_MakeWithNode(f.errorRoot, error_, args...))
}

// FunctionTypeBuilder_IsFunctionTypeDeclaration determines whether the
// given JsDoc info declares a function type.
//
// Java: FunctionTypeBuilder.isFunctionTypeDeclaration(JSDocInfo)
func FunctionTypeBuilder_IsFunctionTypeDeclaration(info *rhino.JSDocInfo) bool {
	return info.GetParameterCount() > 0 ||
		info.HasReturnType() ||
		info.HasThisType() ||
		info.IsConstructor() ||
		info.IsInterface() ||
		info.IsAbstract()
}

// FunctionTypeBuilder_HasMoreTagsToResolve checks whether a type is
// resolvable in the future. If this has a supertype that hasn't been
// resolved yet, then we can assume this type will be OK once the super type
// resolves.
//
// Returns true if objectType is resolvable in the future.
//
// Java: FunctionTypeBuilder.hasMoreTagsToResolve(ObjectType)
func FunctionTypeBuilder_HasMoreTagsToResolve(objectType rhino.ObjectType) bool {
	jutil.CheckArgument(objectType.IsUnknownType())
	ctor := objectType.GetConstructor()
	if ctor != nil {
		// interface extends interfaces
		for _, interfaceType := range ctor.GetExtendedInterfaces() {
			if !interfaceType.IsResolved() {
				return true
			}
		}
	}
	if objectType.GetImplicitPrototype() != nil {
		// constructor extends class
		return !objectType.GetImplicitPrototype().IsResolved()
	}
	return false
}

func functionTypeBuilder_hasSingleThrow(n rhino.Node) bool {
	block := n.GetLastChild()
	return block.HasOneChild() && block.GetFirstChild().IsThrow()
}

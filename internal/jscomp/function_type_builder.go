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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	emptyTemplateTypeSlice       = []*rhino.TemplateType{}
	emptyFunctionParametersSlice = []*rhino.FunctionType_Parameter{}
	emptyFunctionNameString      = ""
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
	fnNamePtr        *string
	compiler         AbstractCompiler
	codingConvention CodingConvention
	typeRegistry     *rhino.JSTypeRegistry
	errorRoot        rhino.Node

	contents        FunctionTypeBuilder_FunctionContents
	defaultContents FunctionTypeBuilder_AstFunctionContents

	syntacticFnName              *string
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

// NewFunctionTypeBuilder creates a builder.
//
// fnName is the function name to be used in error messages; compiler the
// compiler; errorRoot the node to associate with any warning generated by
// this builder; scope the syntactic scope.
//
// Java: FunctionTypeBuilder(String, AbstractCompiler, Node, TypedScope)
func NewFunctionTypeBuilder(fnName *string, compiler AbstractCompiler, errorRoot rhino.Node, scope *TypedScope) *FunctionTypeBuilder {
	rhino.CheckNotNull(errorRoot)
	f := &FunctionTypeBuilder{
		contents:                     functionTypeBuilder_UnknownFunctionContents_singleton,
		templateTypeNames:            emptyTemplateTypeSlice,
		constructorTemplateTypeNames: emptyTemplateTypeSlice,
		codingConvention:             compiler.GetCodingConvention(),
		typeRegistry:                 compiler.GetTypeRegistry(),
		errorRoot:                    errorRoot,
		compiler:                     compiler,
	}
	if fnName == nil || *fnName == "" {
		f.fnName = ""
		f.fnNamePtr = &emptyFunctionNameString
	} else {
		f.fnName = *fnName
		f.fnNamePtr = fnName
	}
	// Go-only: language-adapter — a nil *TypedScope must stay a nil interface.
	if scope != nil {
		f.templateScope = scope
	}
	return f
}

// functionTypeBuilder_nullToEmpty is Guava's Strings.nullToEmpty.
//
// Go-only: language-adapter — Java null strings are nil *string.
func functionTypeBuilder_nullToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
	if syntacticFnName == nil || *syntacticFnName == "" {
		f.syntacticFnName = &emptyFunctionNameString
	} else {
		f.syntacticFnName = syntacticFnName
	}
	return f
}

// SetContents sets the contents of this function.
//
// Java: FunctionTypeBuilder.setContents(FunctionContents)
func (f *FunctionTypeBuilder) SetContents(contents FunctionTypeBuilder_FunctionContents) *FunctionTypeBuilder {
	if contents != nil {
		f.contents = contents
	}
	return f
}

// SetAstContents sets the AST contents of this function without allocating a separate struct.
func (f *FunctionTypeBuilder) SetAstContents(n rhino.Node, hasNonEmptyReturns bool) *FunctionTypeBuilder {
	if !n.IsNil() {
		f.defaultContents.n = n
		f.defaultContents.hasNonEmptyReturns = hasNonEmptyReturns
		f.contents = &f.defaultContents
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

// InferFromOverriddenFunction infers the parameter and return types of a
// function from the parameter and return types of the function it is
// overriding.
//
// oldType is the function being overridden. Does nothing if this is null.
// paramsParent is the PARAM_LIST node of the function that we're assigning
// to. If null, that just means we're not initializing this to a function
// literal.
//
// Java: FunctionTypeBuilder.inferFromOverriddenFunction(FunctionType, Node)
func (f *FunctionTypeBuilder) InferFromOverriddenFunction(oldType rhino.FunctionType, paramsParent rhino.Node) *FunctionTypeBuilder {
	if oldType == nil {
		return f
	}

	// Propagate the template types, if they exist.
	f.templateTypeNames = oldType.GetTemplateTypeMap().GetTemplateKeys()

	f.returnType = oldType.GetReturnType()
	f.returnTypeInferred = oldType.IsReturnTypeInferred()
	if paramsParent.IsNil() {
		// Not a function literal.
		f.parameters = oldType.GetParameters()
		if f.parameters == nil {
			f.parameters = emptyFunctionParametersSlice
		}
	} else {
		// We're overriding with a function literal. Apply type information
		// to each parameter of the literal.
		paramBuilder := rhino.NewFunctionParamBuilder(f.typeRegistry)
		oldParams := oldType.GetParameters()
		oldParamsIndex := 0
		warnedAboutArgList := false
		oldParamsListHitOptArgs := false
		for currentParam := paramsParent.GetFirstChild(); !currentParam.IsNil(); currentParam = currentParam.GetNext() {
			if oldParamsIndex < len(oldParams) {
				oldParam := oldParams[oldParamsIndex]
				oldParamsIndex++

				oldParamsListHitOptArgs =
					oldParamsListHitOptArgs || oldParam.IsVariadic() || oldParam.IsOptional()

				// The subclass method might write its var_args as individual arguments.

				isOptionalArg := oldParam.IsOptional()
				isVarArgs := oldParam.IsVariadic()
				if !currentParam.GetNext().IsNil() && isVarArgs {
					isVarArgs = false
					isOptionalArg = true
				}
				// The subclass method might also make a required parameter into an optional parameter
				// with a default value
				if currentParam.IsDefaultValue() {
					isOptionalArg = true
				}
				paramBuilder.NewParameterFrom(
					rhino.FunctionType_Parameter_Create(oldParam.GetJSType(), isOptionalArg, isVarArgs))
			} else {
				warnedAboutArgList = f.AddParameter(
					paramBuilder,
					f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE),
					warnedAboutArgList,
					f.codingConvention.IsOptionalParameter(currentParam) ||
						oldParamsListHitOptArgs ||
						currentParam.IsDefaultValue(),
					f.codingConvention.IsVarArgsParameter(currentParam)) || warnedAboutArgList
			}
		}

		// Clone any remaining params that aren't in the function literal,
		// but make them optional.
		for oldParamsIndex < len(oldParams) {
			paramBuilder.NewOptionalParameterFrom(oldParams[oldParamsIndex])
			oldParamsIndex++
		}

		f.parameters = paramBuilder.Build()
	}
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

func jsDocTemplateTypes(info *rhino.JSDocInfo) *jutil.LinkedHashMap[string, *rhino.JSTypeExpression] {
	if info == nil {
		return nil
	}
	if m, ok := rhino.JSDocInfo_TEMPLATE_TYPE_NAMES.Get(info).(*jutil.LinkedHashMap[string, *rhino.JSTypeExpression]); ok {
		return m
	}
	return nil
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
	infoTemplateTypes := jsDocTemplateTypes(info)
	var infoTemplateCount int
	if infoTemplateTypes != nil {
		infoTemplateCount = infoTemplateTypes.Size()
	}
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
			f.implementedInterfaces = []rhino.ObjectType{}
			for _, t := range info.GetImplementedInterfaces() {
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
		f.extendedInterfaces = []rhino.ObjectType{}
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

// InferParameterTypes infers the parameter types from the doc info
// alone.
//
// Java: FunctionTypeBuilder.inferParameterTypes(JSDocInfo)
func (f *FunctionTypeBuilder) InferParameterTypes(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	// Create a fake args parent.
	lp := rhino.IR_ParamList()
	paramCount := info.GetParameterCount()
	for i := 0; i < paramCount; i++ {
		name, _ := info.GetParameterNameValueAt(i)
		lp.AddChildToBack(rhino.IR_Name(name))
	}

	return f.InferParameterTypesWithParamsParent(lp, info)
}

// InferParameterTypesWithParamsParent infers the parameter types from the list of parameter
// names and the JSDoc info.
//
// Java: FunctionTypeBuilder.inferParameterTypes(Node, JSDocInfo)
func (f *FunctionTypeBuilder) InferParameterTypesWithParamsParent(paramsParent rhino.Node, info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if paramsParent.IsNil() {
		if info == nil {
			return f
		} else {
			return f.InferParameterTypes(info)
		}
	}

	// arguments
	// Go-only: language-adapter — the Iterator<Parameter> over parameters is
	// an index into oldParameters.
	var oldParameters []*rhino.FunctionType_Parameter
	oldParametersIndex := 0
	var oldParameterType *rhino.FunctionType_Parameter
	if f.parameters != nil {
		oldParameters = f.parameters
		if oldParametersIndex < len(oldParameters) {
			oldParameterType = oldParameters[oldParametersIndex]
			oldParametersIndex++
		}
	}

	builder := rhino.NewFunctionParamBuilder(f.typeRegistry)
	warnedAboutArgList := false
	hasJsDocParams := info != nil && info.GetParameterCount() > 0
	var seen []string
	isVarArgs := false
	paramIndex := 0
	for param := paramsParent.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
		isOptionalParam := false
		var paramLhs rhino.Node

		if param.IsRest() {
			isVarArgs = true
			paramLhs = param.GetOnlyChild()
		} else if param.IsDefaultValue() {
			// The first child is the actual positional parameter
			paramLhs = rhino.CheckNotNullMsg1(param.GetFirstChild(), "%s", param)
			isOptionalParam = true
		} else {
			isVarArgs = f.IsVarArgsParameterByConvention(param)
			isOptionalParam = f.IsOptionalParameterByConvention(param)
			paramLhs = param
		}

		var paramName string
		var hasParamName bool
		if paramLhs.IsName() {
			paramName = paramLhs.GetString()
			hasParamName = true
		} else {
			jutil.CheckState(paramLhs.IsDestructuringPattern())
			// Right now, the only way to match a JSDoc param to a destructuring parameter is through
			// ordering the JSDoc parameters. So the third formal parameter will correspond to the
			// third JSDoc parameter.
			if info != nil {
				if p, ok := info.GetParameterNameValueAt(paramIndex); ok {
					paramName = p
					hasParamName = true
				}
			}
		}
		if hasParamName && hasJsDocParams {
			seen = append(seen, paramName)
		}

		// type from JSDocInfo
		var parameterType rhino.JSType
		var parameterTypeExpression *rhino.JSTypeExpression
		if info != nil && hasParamName {
			parameterTypeExpression = info.GetParameterType(paramName)
		}
		if parameterTypeExpression != nil {
			parameterType = parameterTypeExpression.Evaluate(f.templateScope, f.typeRegistry)
			isOptionalParam = isOptionalParam || parameterTypeExpression.IsOptionalArg()
			isVarArgs = isVarArgs || parameterTypeExpression.IsVarArgs()
		} else if paramLhs.GetJSDocInfo() != nil && paramLhs.GetJSDocInfo().HasType() {
			parameterTypeExpression := paramLhs.GetJSDocInfo().GetType()
			parameterType = parameterTypeExpression.Evaluate(f.templateScope, f.typeRegistry)
			isOptionalParam = parameterTypeExpression.IsOptionalArg()
			isVarArgs = parameterTypeExpression.IsVarArgs()
		} else if oldParameterType != nil && oldParameterType.GetJSType() != nil {
			parameterType = oldParameterType.GetJSType()
			isOptionalParam = oldParameterType.IsOptional()
			isVarArgs = oldParameterType.IsVariadic()
		} else {
			parameterType = f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}

		warnedAboutArgList = f.AddParameter(builder, parameterType, warnedAboutArgList, isOptionalParam, isVarArgs) ||
			warnedAboutArgList

		oldParameterType = nil
		if oldParametersIndex < len(oldParameters) {
			oldParameterType = oldParameters[oldParametersIndex]
			oldParametersIndex++
		}
		paramIndex++
	}
	// Copy over any old parameters that aren't in the param list.
	if !isVarArgs {
		for oldParameterType != nil && !isVarArgs {
			builder.NewParameterFrom(oldParameterType)
			oldParameterType = nil
			if oldParametersIndex < len(oldParameters) {
				oldParameterType = oldParameters[oldParametersIndex]
				oldParametersIndex++
			}
		}
	}

	if hasJsDocParams {
		paramCount := info.GetParameterCount()
		for i := 0; i < paramCount; i++ {
			jsDocParam, _ := info.GetParameterNameValueAt(i)
			if !slices.Contains(seen, jsDocParam) {
				f.ReportWarning(FunctionTypeBuilder_INEXISTENT_PARAM, jsDocParam, f.FormatFnName())
			}
		}
	}

	f.parameters = builder.Build()
	return f
}

// RegisterTemplates registers the template keys in a template scope and on
// the function node.
//
// Java: FunctionTypeBuilder.registerTemplates(Iterable, Node)
func (f *FunctionTypeBuilder) RegisterTemplates(templates []*rhino.TemplateType, scopeRoot rhino.Node) {
	if len(templates) != 0 {
		// Add any templates from JSDoc into our template scope.
		f.templateScope = f.typeRegistry.CreateScopeWithTemplates(f.templateScope, templates)
		// Register the template types on the scope root node, if there is one.
		if !scopeRoot.IsNil() {
			f.typeRegistry.RegisterTemplateTypeNamesInScope(templates, scopeRoot)
		}
	}
}

// InferConstructorParameters infers parameters from the params list and
// info. Also maybe adds extra templates.
//
// Java: FunctionTypeBuilder.inferConstructorParameters(Node, JSDocInfo)
func (f *FunctionTypeBuilder) InferConstructorParameters(argsParent rhino.Node, info *rhino.JSDocInfo) *FunctionTypeBuilder {
	// Look for template parameters in 'info': these will be added to anything from the class.
	if info != nil {
		f.SetConstructorTemplateTypeNames(
			f.BuildTemplateTypesFromJSDocInfo(info, true), argsParent.GetParent())
	}

	f.InferParameterTypesWithParamsParent(argsParent, info)

	return f
}

// Java: FunctionTypeBuilder.inferImplicitConstructorParameters(ImmutableList)
func (f *FunctionTypeBuilder) InferImplicitConstructorParameters(parameters []*rhino.FunctionType_Parameter) *FunctionTypeBuilder {
	f.parameters = parameters
	return f
}

// Java: FunctionTypeBuilder.inferClosurePrimitive(JSDocInfo)
func (f *FunctionTypeBuilder) InferClosurePrimitive(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if info != nil && info.HasClosurePrimitiveId() {
		f.closurePrimitiveId = rhino.ClosurePrimitive_FromStringId(info.GetClosurePrimitiveId())
	}
	return f
}

// Java: FunctionTypeBuilder.setConstructorTemplateTypeNames(List, Node)
func (f *FunctionTypeBuilder) SetConstructorTemplateTypeNames(templates []*rhino.TemplateType, ctor rhino.Node) {
	if len(templates) != 0 {
		f.constructorTemplateTypeNames = append([]*rhino.TemplateType{}, templates...)
		if len(f.templateTypeNames) == 0 {
			f.templateTypeNames = append([]*rhino.TemplateType{}, templates...)
		} else {
			names := make([]*rhino.TemplateType, 0, len(f.templateTypeNames)+len(f.constructorTemplateTypeNames))
			names = append(names, f.templateTypeNames...)
			names = append(names, f.constructorTemplateTypeNames...)
			f.templateTypeNames = names
		}
		f.RegisterTemplates(templates, ctor)
	}
}

// IsOptionalParameterByConvention returns whether the given param is an
// optional param.
//
// Java: FunctionTypeBuilder.isOptionalParameterByConvention(Node)
func (f *FunctionTypeBuilder) IsOptionalParameterByConvention(param rhino.Node) bool {
	if param.IsDestructuringPattern() {
		return false
	}
	return f.codingConvention.IsOptionalParameter(param)
}

// IsVarArgsParameterByConvention determines whether this is a var args
// parameter.
//
// Java: FunctionTypeBuilder.isVarArgsParameterByConvention(Node)
func (f *FunctionTypeBuilder) IsVarArgsParameterByConvention(param rhino.Node) bool {
	if param.IsDestructuringPattern() {
		return false
	}

	return f.codingConvention.IsVarArgsParameter(param)
}

// Java: FunctionTypeBuilder.buildTemplateTypesFromJSDocInfo(JSDocInfo, boolean)
func (f *FunctionTypeBuilder) BuildTemplateTypesFromJSDocInfo(info *rhino.JSDocInfo, allowTypeTransformations bool) []*rhino.TemplateType {
	infoTypeKeys := jsDocTemplateTypes(info)
	infoTypeTransformations := jsDocTypeTransformations(info)
	hasKeys := infoTypeKeys != nil && !infoTypeKeys.IsEmpty()
	hasTransformations := infoTypeTransformations != nil && !infoTypeTransformations.IsEmpty()
	if !hasKeys && !hasTransformations {
		return emptyTemplateTypeSlice
	}

	// Temporarily bootstrap the template environment with unbound (unknown bound) template types
	var unboundedTemplates []*rhino.TemplateType
	if hasKeys {
		unboundedTemplates = make([]*rhino.TemplateType, 0, infoTypeKeys.Size())
		for templateKey := range infoTypeKeys.Keys() {
			unboundedTemplates = append(unboundedTemplates, f.typeRegistry.CreateTemplateType(templateKey))
		}
	}
	f.templateScope = f.typeRegistry.CreateScopeWithTemplates(f.templateScope, unboundedTemplates)

	// Evaluate template type bounds with bootstrapped environment and reroute the bounds to these
	templates := []*rhino.TemplateType{}
	if hasKeys {
		templatesToBounds := jutil.NewLinkedHashMap[*rhino.TemplateType, rhino.JSType]()
		for key, expr := range infoTypeKeys.All() {
			typeBound := f.typeRegistry.EvaluateTypeExpression(expr, f.templateScope)
			// It's an error to mark a template bound explicitly {?}. Unbounded templates have an implicit
			// unknown bound. Allowing explicit unknowns would make it more difficult to stricten their
			// treatment in the future, since "unknown" is currently used as a proxy for "implicit".
			if expr.IsExplicitUnknownTemplateBound() {
				f.ReportError(FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND, jutil.ToJavaString(typeBound), key)
			}
			template := f.typeRegistry.GetType(f.templateScope, key).ToMaybeTemplateType()
			if template != nil {
				templatesToBounds.Put(template, typeBound)
			} else {
				templatesToBounds.Put(
					f.typeRegistry.CreateTemplateTypeWithBound(key, typeBound), typeBound)
			}
		}

		for template, bound := range templatesToBounds.All() {
			template.SetBound(bound)
			templates = append(templates, template)
		}
	}

	if hasTransformations {
		for key, value := range infoTypeTransformations.All() {
			if allowTypeTransformations {
				templates = append(templates,
					f.typeRegistry.CreateTemplateTypeWithTransformation(key, value))
			} else {
				f.ReportWarning(FunctionTypeBuilder_TEMPLATE_TRANSFORMATION_ON_CLASS, key)
			}
		}
	}

	builtTemplates := templates
	for _, template := range builtTemplates {
		if template.ContainsCycle() {
			f.ReportError(
				RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR,
				"Cycle detected in inheritance chain of type "+*template.GetReferenceName())
		}
	}

	return builtTemplates
}

// InferTemplateTypeName infers the template type from the doc info.
//
// Java: FunctionTypeBuilder.inferTemplateTypeName(JSDocInfo, JSType)
func (f *FunctionTypeBuilder) InferTemplateTypeName(info *rhino.JSDocInfo, ownerType rhino.JSType) *FunctionTypeBuilder {
	// NOTE: these template type names may override a list
	// of inherited ones from an overridden function.

	if info != nil && !f.MaybeUseNativeClassTemplateNames(info) {
		templates := f.BuildTemplateTypesFromJSDocInfo(info, !(f.isConstructor || f.isInterface))
		if len(templates) != 0 {
			f.templateTypeNames = templates
		}
	}

	var ownerTypeKeys []*rhino.TemplateType
	if ownerType != nil {
		ownerTypeKeys = ownerType.GetTemplateTypeMap().GetTemplateKeys()
	} else {
		ownerTypeKeys = emptyTemplateTypeSlice
	}

	if len(f.templateTypeNames) != 0 || len(ownerTypeKeys) != 0 {
		// TODO(sdh): The order of these should be switched to avoid class templates shadowing
		// method templates, but this currently loosens type checking of arrays more than we'd like.
		// See http://github.com/google/closure-compiler/issues/2973
		if len(ownerTypeKeys) == 0 {
			f.RegisterTemplates(f.templateTypeNames, f.contents.GetSourceNode())
		} else if len(f.templateTypeNames) == 0 {
			f.RegisterTemplates(ownerTypeKeys, f.contents.GetSourceNode())
		} else {
			concat := make([]*rhino.TemplateType, 0, len(f.templateTypeNames)+len(ownerTypeKeys))
			concat = append(concat, f.templateTypeNames...)
			concat = append(concat, ownerTypeKeys...)
			f.RegisterTemplates(concat, f.contents.GetSourceNode())
		}
	}

	return f
}

// AddParameter adds a parameter to the param list.
//
// builder is a builder; paramType the parameter type; warnedAboutArgList
// whether we've already warned about arg ordering issues (like if optional
// args appeared before required ones); isOptional whether this is an
// optional parameter; isVarArgs whether this is a var args parameter.
// Returns whether a warning was emitted.
//
// Java: FunctionTypeBuilder.addParameter(FunctionParamBuilder, JSType, boolean, boolean, boolean)
func (f *FunctionTypeBuilder) AddParameter(
	builder *rhino.FunctionParamBuilder,
	paramType rhino.JSType,
	warnedAboutArgList bool,
	isOptional bool,
	isVarArgs bool) bool {
	emittedWarning := false
	if isOptional {
		// Remembering that an optional parameter has been encountered
		// so that if a non optional param is encountered later, an
		// error can be reported.
		if !builder.AddOptionalParams(paramType) && !warnedAboutArgList {
			f.ReportWarning(FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
			emittedWarning = true
		}
	} else if isVarArgs {
		if !builder.AddVarArgs(paramType) && !warnedAboutArgList {
			f.ReportWarning(FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
			emittedWarning = true
		}
	} else {
		if !builder.AddRequiredParams(paramType) && !warnedAboutArgList {
			// An optional parameter was seen and this argument is not an optional
			// or var arg so it is an error.
			if builder.HasVarArgs() {
				f.ReportWarning(FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
			} else {
				f.ReportWarning(FunctionTypeBuilder_OPTIONAL_ARG_AT_END)
			}
			emittedWarning = true
		}
	}
	return emittedWarning
}

// ProvideDefaultReturnType sets the returnType for this function using very
// basic type inference.
//
// Java: FunctionTypeBuilder.provideDefaultReturnType()
func (f *FunctionTypeBuilder) ProvideDefaultReturnType() {
	if !f.contents.GetSourceNode().IsNil() && f.contents.GetSourceNode().IsAsyncGeneratorFunction() {
		// Set the return type of a generator function to:
		//   @return {!AsyncGenerator<?>}
		generatorType := f.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ASYNC_GENERATOR_TYPE)
		f.returnType = f.typeRegistry.CreateTemplatizedType(
			generatorType, f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
		return
	} else if !f.contents.GetSourceNode().IsNil() && f.contents.GetSourceNode().IsGeneratorFunction() {
		// Set the return type of a generator function to:
		//   @return {!Generator<?>}
		generatorType := f.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GENERATOR_TYPE)
		f.returnType = f.typeRegistry.CreateTemplatizedType(
			generatorType, f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
		return
	}

	inferredReturnType := f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if !f.contents.MayHaveNonEmptyReturns() &&
		!f.contents.MayHaveSingleThrow() &&
		!f.contents.MayBeFromExterns() {
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

	if !f.contents.GetSourceNode().IsNil() && f.contents.GetSourceNode().IsAsyncFunction() {
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
	namePtr := f.fnNamePtr
	if namePtr == nil {
		namePtr = jutil.StringPtr(f.fnName)
	}
	return rhino.FunctionType_Builder_(f.typeRegistry).
		WithName(namePtr).
		WithSourceNode(f.contents.GetSourceNode()).
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
				existingFn.SetSource(f.contents.GetSourceNode())
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
	if *f.syntacticFnName != "" && !jutil.StringStartsWith(*f.syntacticFnName, "this.") {
		f.typeRegistry.DeclareTypeForExactScope(
			functionTypeBuilder_staticScope(f.declarationScope), *f.syntacticFnName, fnType.GetInstanceType())
	}
	return fnType
}

// Java: FunctionTypeBuilder.getOrCreateInterface()
func (f *FunctionTypeBuilder) GetOrCreateInterface() rhino.FunctionType {
	var fnType rhino.FunctionType

	type_ := f.typeRegistry.GetType(functionTypeBuilder_staticScope(f.declarationScope), *f.syntacticFnName)
	if type_ != nil && type_.IsInstanceType() {
		ctor := type_.ToMaybeObjectType().GetConstructor()
		if ctor.IsInterface() {
			fnType = ctor
			fnType.SetSource(f.contents.GetSourceNode())
		}
	}

	if fnType == nil {
		fnType = f.CreateDefaultBuilder().ForInterface().WithParameters().Build()
		if f.makesStructs {
			fnType.SetStruct()
		}
		if f.fnName != "" {
			f.typeRegistry.DeclareTypeForExactScope(
				functionTypeBuilder_staticScope(f.declarationScope), *f.syntacticFnName, fnType.GetInstanceType())
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

// FunctionTypeBuilder_FunctionContents holds data dynamically inferred about
// functions.
type FunctionTypeBuilder_FunctionContents interface {
	// GetSourceNode returns the source node of this function. May be null.
	//
	// Java: FunctionTypeBuilder.FunctionContents.getSourceNode()
	GetSourceNode() rhino.Node

	// MayBeFromExterns returns if the function may be in externs.
	//
	// Java: FunctionTypeBuilder.FunctionContents.mayBeFromExterns()
	MayBeFromExterns() bool

	// MayHaveNonEmptyReturns returns if a return of a real value (not
	// undefined) appears.
	//
	// Java: FunctionTypeBuilder.FunctionContents.mayHaveNonEmptyReturns()
	MayHaveNonEmptyReturns() bool

	// MayHaveSingleThrow returns if this consists of a single throw.
	//
	// Java: FunctionTypeBuilder.FunctionContents.mayHaveSingleThrow()
	MayHaveSingleThrow() bool
}

// FunctionTypeBuilder_UnknownFunctionContents is the FunctionContents of a
// function whose body is unknown.
type FunctionTypeBuilder_UnknownFunctionContents struct{}

var functionTypeBuilder_UnknownFunctionContents_singleton = &FunctionTypeBuilder_UnknownFunctionContents{}

// Java: FunctionTypeBuilder.UnknownFunctionContents.get()
func FunctionTypeBuilder_UnknownFunctionContents_Get() FunctionTypeBuilder_FunctionContents {
	return functionTypeBuilder_UnknownFunctionContents_singleton
}

// Java: FunctionTypeBuilder.UnknownFunctionContents.getSourceNode()
func (c *FunctionTypeBuilder_UnknownFunctionContents) GetSourceNode() rhino.Node {
	return rhino.NilNode
}

// Java: FunctionTypeBuilder.UnknownFunctionContents.mayBeFromExterns()
func (c *FunctionTypeBuilder_UnknownFunctionContents) MayBeFromExterns() bool {
	return true
}

// Java: FunctionTypeBuilder.UnknownFunctionContents.mayHaveNonEmptyReturns()
func (c *FunctionTypeBuilder_UnknownFunctionContents) MayHaveNonEmptyReturns() bool {
	return true
}

// Java: FunctionTypeBuilder.UnknownFunctionContents.mayHaveSingleThrow()
func (c *FunctionTypeBuilder_UnknownFunctionContents) MayHaveSingleThrow() bool {
	return true
}

// FunctionTypeBuilder_AstFunctionContents is the FunctionContents of a
// function literal in the AST.
type FunctionTypeBuilder_AstFunctionContents struct {
	n                  rhino.Node
	hasNonEmptyReturns bool
}

// Java: FunctionTypeBuilder.AstFunctionContents(Node)
func NewFunctionTypeBuilder_AstFunctionContents(n rhino.Node) *FunctionTypeBuilder_AstFunctionContents {
	c := &FunctionTypeBuilder_AstFunctionContents{}
	c.n = n
	return c
}

// Java: FunctionTypeBuilder.AstFunctionContents.getSourceNode()
func (c *FunctionTypeBuilder_AstFunctionContents) GetSourceNode() rhino.Node {
	return c.n
}

// Java: FunctionTypeBuilder.AstFunctionContents.mayBeFromExterns()
func (c *FunctionTypeBuilder_AstFunctionContents) MayBeFromExterns() bool {
	return c.n.IsFromExterns()
}

// Java: FunctionTypeBuilder.AstFunctionContents.mayHaveNonEmptyReturns()
func (c *FunctionTypeBuilder_AstFunctionContents) MayHaveNonEmptyReturns() bool {
	return c.hasNonEmptyReturns
}

// Java: FunctionTypeBuilder.AstFunctionContents.recordNonEmptyReturn()
func (c *FunctionTypeBuilder_AstFunctionContents) RecordNonEmptyReturn() {
	c.hasNonEmptyReturns = true
}

// Java: FunctionTypeBuilder.AstFunctionContents.mayHaveSingleThrow()
func (c *FunctionTypeBuilder_AstFunctionContents) MayHaveSingleThrow() bool {
	block := c.n.GetLastChild()
	return block.HasOneChild() && block.GetFirstChild().IsThrow()
}

var (
	_ FunctionTypeBuilder_FunctionContents = (*FunctionTypeBuilder_UnknownFunctionContents)(nil)
	_ FunctionTypeBuilder_FunctionContents = (*FunctionTypeBuilder_AstFunctionContents)(nil)
)

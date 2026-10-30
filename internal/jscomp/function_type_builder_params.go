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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
		var paramBuilder rhino.FunctionParamBuilder
		paramBuilder.Init(f.typeRegistry, paramsParent.GetChildCount())
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
					&paramBuilder,
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

func (f *FunctionTypeBuilder) paramNameConventions(name string) (isOptional, isVarArgs bool) {
	for conv := f.codingConvention; conv != nil; {
		switch c := conv.(type) {
		case *GoogleCodingConvention:
			if strings.HasPrefix(name, googleOptionalArgPrefix) {
				isOptional = true
			}
			if name == googleVarArgsName {
				isVarArgs = true
			}
			conv = c.CodingConvention
		case *ClosureCodingConvention:
			conv = c.CodingConvention
		case *ChromeCodingConvention:
			conv = c.CodingConvention
		case *CodingConventions_DefaultCodingConvention:
			return isOptional, isVarArgs
		default:
			tmp := rhino.IR_Name(name)
			return isOptional || conv.IsOptionalParameter(tmp), isVarArgs || conv.IsVarArgsParameter(tmp)
		}
	}
	return isOptional, isVarArgs
}

// InferParameterTypes infers the parameter types from the doc info
// alone.
//
// Java: FunctionTypeBuilder.inferParameterTypes(JSDocInfo)
func (f *FunctionTypeBuilder) InferParameterTypes(info *rhino.JSDocInfo) *FunctionTypeBuilder {
	paramCount := info.GetParameterCount()
	capHint := paramCount
	if len(f.parameters) > capHint {
		capHint = len(f.parameters)
	}
	var builder rhino.FunctionParamBuilder
	builder.Init(f.typeRegistry, capHint)
	warnedAboutArgList := false
	oldIdx := 0
	isVarArgs := false
	for i := 0; i < paramCount; i++ {
		paramName, _ := info.GetParameterNameValueAt(i)
		isOptionalParam, varArgsByConv := f.paramNameConventions(paramName)
		isVarArgs = varArgsByConv

		var oldParameterType rhino.FunctionType_Parameter
		if oldIdx < len(f.parameters) {
			oldParameterType = f.parameters[oldIdx]
			oldIdx++
		}

		var parameterType rhino.JSType
		if parameterTypeExpression := info.GetParameterType(paramName); parameterTypeExpression != nil {
			parameterType = parameterTypeExpression.Evaluate(f.templateScope, f.typeRegistry)
			isOptionalParam = isOptionalParam || parameterTypeExpression.IsOptionalArg()
			isVarArgs = isVarArgs || parameterTypeExpression.IsVarArgs()
		} else if oldParameterType.GetJSType() != nil {
			parameterType = oldParameterType.GetJSType()
			isOptionalParam = oldParameterType.IsOptional()
			isVarArgs = oldParameterType.IsVariadic()
		} else {
			parameterType = f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}

		warnedAboutArgList = f.AddParameter(&builder, parameterType, warnedAboutArgList, isOptionalParam, isVarArgs) ||
			warnedAboutArgList
	}
	if !isVarArgs {
		for ; oldIdx < len(f.parameters); oldIdx++ {
			builder.NewParameterFrom(f.parameters[oldIdx])
		}
	}
	f.parameters = builder.Build()
	return f
}

// InferParameterTypesWithParamsParent infers the parameter types from the list of parameter
// names and the JSDoc info.
//
// Java: FunctionTypeBuilder.inferParameterTypes(Node, JSDocInfo)
func (f *FunctionTypeBuilder) InferParameterTypesWithParamsParent(paramsParent rhino.Node, info *rhino.JSDocInfo) *FunctionTypeBuilder {
	if paramsParent.IsNil() {
		if info == nil {
			return f
		}
		return f.InferParameterTypes(info)
	}

	var builder rhino.FunctionParamBuilder
	builder.Init(f.typeRegistry, paramsParent.GetChildCount())
	warnedAboutArgList := false
	hasJsDocParams := info != nil && info.GetParameterCount() > 0
	var seenBuf [8]string
	seen := seenBuf[:0]
	oldIdx := 0
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

		var oldParameterType rhino.FunctionType_Parameter
		if oldIdx < len(f.parameters) {
			oldParameterType = f.parameters[oldIdx]
			oldIdx++
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
		} else if oldParameterType.GetJSType() != nil {
			parameterType = oldParameterType.GetJSType()
			isOptionalParam = oldParameterType.IsOptional()
			isVarArgs = oldParameterType.IsVariadic()
		} else {
			parameterType = f.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}

		warnedAboutArgList = f.AddParameter(&builder, parameterType, warnedAboutArgList, isOptionalParam, isVarArgs) ||
			warnedAboutArgList
		paramIndex++
	}
	// Copy over any old parameters that aren't in the param list.
	if !isVarArgs {
		for ; oldIdx < len(f.parameters); oldIdx++ {
			builder.NewParameterFrom(f.parameters[oldIdx])
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
func (f *FunctionTypeBuilder) InferImplicitConstructorParameters(parameters []rhino.FunctionType_Parameter) *FunctionTypeBuilder {
	f.parameters = parameters
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

type templateBoundEntry struct {
	template *rhino.TemplateType
	bound    rhino.JSType
}

// Java: FunctionTypeBuilder.buildTemplateTypesFromJSDocInfo(JSDocInfo, boolean)
func (f *FunctionTypeBuilder) BuildTemplateTypesFromJSDocInfo(info *rhino.JSDocInfo, allowTypeTransformations bool) []*rhino.TemplateType {
	infoTypeKeys := info.GetTemplateTypeEntries()
	infoTypeTransformations := jsDocTypeTransformations(info)
	hasKeys := len(infoTypeKeys) > 0
	hasTransformations := infoTypeTransformations != nil && !infoTypeTransformations.IsEmpty()
	if !hasKeys && !hasTransformations {
		return emptyTemplateTypeSlice
	}

	// Temporarily bootstrap the template environment with unbound (unknown bound) template types
	var unboundedBuf [4]*rhino.TemplateType
	unboundedTemplates := unboundedBuf[:0]
	for _, entry := range infoTypeKeys {
		unboundedTemplates = append(unboundedTemplates, f.typeRegistry.CreateTemplateType(entry.Name))
	}
	f.templateScope = f.typeRegistry.CreateScopeWithTemplates(f.templateScope, unboundedTemplates)

	// Evaluate template type bounds with bootstrapped environment and reroute the bounds to these
	totalCap := len(infoTypeKeys)
	if hasTransformations && allowTypeTransformations {
		totalCap += infoTypeTransformations.Size()
	}
	var templates []*rhino.TemplateType
	if totalCap > 0 {
		templates = make([]*rhino.TemplateType, 0, totalCap)
	} else {
		templates = emptyTemplateTypeSlice
	}
	if hasKeys {
		var boundBuf [4]templateBoundEntry
		bounds := boundBuf[:0]
		for _, entry := range infoTypeKeys {
			key, expr := entry.Name, entry.Type
			typeBound := f.typeRegistry.EvaluateTypeExpression(expr, f.templateScope)
			// It's an error to mark a template bound explicitly {?}. Unbounded templates have an implicit
			// unknown bound. Allowing explicit unknowns would make it more difficult to stricten their
			// treatment in the future, since "unknown" is currently used as a proxy for "implicit".
			if expr.IsExplicitUnknownTemplateBound() {
				f.ReportError(FunctionTypeBuilder_TEMPLATE_TYPE_ILLEGAL_BOUND, jutil.ToJavaString(typeBound), key)
			}
			template := f.typeRegistry.GetType(f.templateScope, key).ToMaybeTemplateType()
			if template == nil {
				template = f.typeRegistry.CreateTemplateTypeWithBound(key, typeBound)
			}
			updated := false
			for i := range bounds {
				if bounds[i].template == template {
					bounds[i].bound = typeBound
					updated = true
					break
				}
			}
			if !updated {
				bounds = append(bounds, templateBoundEntry{template: template, bound: typeBound})
			}
		}

		for _, b := range bounds {
			b.template.SetBound(b.bound)
			templates = append(templates, b.template)
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

	for _, template := range templates {
		if template.ContainsCycle() {
			f.ReportError(
				RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR,
				"Cycle detected in inheritance chain of type "+*template.GetReferenceName())
		}
	}

	return templates
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
			f.RegisterTemplates(f.templateTypeNames, f.sourceNode)
		} else if len(f.templateTypeNames) == 0 {
			f.RegisterTemplates(ownerTypeKeys, f.sourceNode)
		} else {
			var concatBuf [8]*rhino.TemplateType
			concat := append(concatBuf[:0], f.templateTypeNames...)
			concat = append(concat, ownerTypeKeys...)
			f.RegisterTemplates(concat, f.sourceNode)
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
		if !builder.AddOptionalParam(paramType) && !warnedAboutArgList {
			f.ReportWarning(FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
			emittedWarning = true
		}
	} else if isVarArgs {
		if !builder.AddVarArgs(paramType) && !warnedAboutArgList {
			f.ReportWarning(FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
			emittedWarning = true
		}
	} else {
		if !builder.AddRequiredParam(paramType) && !warnedAboutArgList {
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

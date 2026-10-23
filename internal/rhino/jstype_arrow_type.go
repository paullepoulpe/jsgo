// Port of com/google/javascript/rhino/jstype/ArrowType.java @ closure-compiler v20261015

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

// ArrowType models a "bare" function type: from some parameter types to a
// return type.
//
// JavaScript functions include more things like properties, type of THIS,
// etc, and are modeled by FunctionType.
type ArrowType struct {
	parameterList []*FunctionType_Parameter
	// Go-only: performance — the type's base (8 bytes) rather than a JSType
	// interface (16); GetReturnType returns its self().
	returnType *JSTypeBase

	// Whether the return type is inferred.
	returnTypeInferred bool
	resolved           bool

	// FunctionTypeBase 1-byte fields packed into ArrowType's tail padding.
	fnConstructorAmbiguity  FunctionType_ConstructorAmbiguity
	fnKind                  FunctionType_Kind
	fnPropAccess            uint8
	fnIsStructuralInterface bool
	fnIsAbstract            bool
}

// Java: ArrowType.TYPE_CLASS
const ArrowType_TYPE_CLASS = JSTypeClass_ARROW

// Java: ArrowType(JSTypeRegistry, List<Parameter>, JSType, boolean)
func (t *ArrowType) init(
	registry *JSTypeRegistry,
	parameters []*FunctionType_Parameter,
	returnType JSType,
	returnTypeInferred bool) {
	if parameters == nil {
		t.parameterList = registry.CreateParametersWithVarArgs(registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
	} else {
		t.parameterList = arrowTypeCopyParameters(parameters)
	}
	if returnType == nil {
		t.returnType = jsTypeBase(registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
	} else {
		t.returnType = jsTypeBase(returnType)
	}
	t.returnTypeInferred = returnTypeInferred
}

// arrowTypeCopyParameters is ImmutableList.copyOf(parameters), which
// rejects null elements.
func arrowTypeCopyParameters(parameters []*FunctionType_Parameter) []*FunctionType_Parameter {
	if len(parameters) == 0 {
		return emptyFunctionTypeParameterSlice
	}
	out := make([]*FunctionType_Parameter, len(parameters))
	for i, p := range parameters {
		out[i] = jutil.CheckNotNull(p)
	}
	return out
}

// Java: ArrowType.getTypeClass()
func (t *ArrowType) GetTypeClass() JSTypeClass {
	return ArrowType_TYPE_CLASS
}

// Java: ArrowType.recursionUnsafeHashCode()
func (t *ArrowType) RecursionUnsafeHashCode() int32 {
	hashCode := jutil.ObjectHashCode(jsTypeOf(t.returnType))
	for _, param := range t.parameterList {
		hashCode = hashCode*31 + jutil.ObjectHashCode(param.GetJSType())
	}
	return hashCode
}

func (t *ArrowType) HashCode() int32 {
	return t.RecursionUnsafeHashCode()
}

// Java: ArrowType.getReturnType()
func (t *ArrowType) GetReturnType() JSType {
	return jsTypeOf(t.returnType)
}

// Java: ArrowType.getParameterList()
func (t *ArrowType) GetParameterList() []*FunctionType_Parameter {
	return t.parameterList
}

// IsReturnTypeInferred reads the package-private field returnTypeInferred.
//
// Go-only: language-adapter — field read by FunctionType.
func (t *ArrowType) IsReturnTypeInferred() bool {
	return t.returnTypeInferred
}

// Java: ArrowType.getPossibleToBooleanOutcomes()
func (t *ArrowType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// Java: ArrowType.resolveInternal(ErrorReporter)
func (t *ArrowType) resolveInternal(reporter ErrorReporter) {
	if t.resolved {
		return
	}
	t.resolved = true
	t.returnType = jsTypeBase(JSType_SafeResolve(jsTypeOf(t.returnType), reporter))
	for _, param := range t.parameterList {
		param.GetJSType().Resolve(reporter)
	}
}

// Java: ArrowType.hasUnknownParamsOrReturn()
func (t *ArrowType) HasUnknownParamsOrReturn() bool {
	for _, param := range t.parameterList {
		type_ := param.GetJSType()
		if type_.IsUnknownType() {
			return true
		}
	}

	return t.returnType == nil || t.returnType.self().IsUnknownType()
}

// Java: ArrowType.appendTo(TypeStringBuilder)
func (t *ArrowType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("[ArrowType]")
}

// Java: ArrowType.hasAnyTemplateTypesInternal()
func (t *ArrowType) HasAnyTemplateTypesInternal() bool {
	return t.returnType.self().HasAnyTemplateTypes() || t.HasTemplatedParameterType()
}

func (t *ArrowType) HasAnyTemplateTypes() bool {
	return t.HasAnyTemplateTypesInternal()
}

// Java: ArrowType.HasTemplatedParameterType()
func (t *ArrowType) HasTemplatedParameterType() bool {
	for _, param := range t.parameterList {
		type_ := param.GetJSType()
		if type_.HasAnyTemplateTypes() {
			return true
		}
	}
	return false
}

// String is a string representation of this type, suitable for printing in
// warnings.
//
// Java: ArrowType.toString()
func (t *ArrowType) String() string {
	typeStringBuilder := NewTypeStringBuilder(false)
	typeStringBuilder.Append("(")
	elements := make([]any, len(t.parameterList))
	for i, p := range t.parameterList {
		elements[i] = p
	}
	typeStringBuilder.AppendAll(elements, ",")
	typeStringBuilder.Append(") -> ")
	typeStringBuilder.AppendType(jsTypeOf(t.returnType))
	return typeStringBuilder.Build()
}

// Port of com/google/javascript/rhino/jstype/FunctionParamBuilder.java @ closure-compiler v20261021

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

import "slices"

// FunctionParamBuilder is a builder for the list representing FunctionType
// Parameters.
type FunctionParamBuilder struct {
	registry   *JSTypeRegistry
	parameters []*FunctionType_Parameter
	consumed   bool
}

var emptyFunctionTypeParams = []*FunctionType_Parameter{}

// Init initializes a stack-allocated FunctionParamBuilder with an optional capacity hint.
func (b *FunctionParamBuilder) Init(registry *JSTypeRegistry, initialParameterCapacity int) {
	b.registry = registry
	b.consumed = false
	if initialParameterCapacity > 0 {
		b.parameters = make([]*FunctionType_Parameter, 0, initialParameterCapacity)
	} else {
		b.parameters = nil
	}
}

func (b *FunctionParamBuilder) ensureMutable() {
	if b.consumed {
		b.parameters = slices.Clone(b.parameters)
		b.consumed = false
	}
}

// Java: FunctionParamBuilder(JSTypeRegistry)
func NewFunctionParamBuilder(registry *JSTypeRegistry) *FunctionParamBuilder {
	return &FunctionParamBuilder{registry: registry}
}

// Java: FunctionParamBuilder(JSTypeRegistry, int)
func NewFunctionParamBuilderWithCapacity(registry *JSTypeRegistry, initialParameterCapacity int) *FunctionParamBuilder {
	b := &FunctionParamBuilder{registry: registry}
	if initialParameterCapacity > 0 {
		b.parameters = make([]*FunctionType_Parameter, 0, initialParameterCapacity)
	}
	return b
}

// AddRequiredParam adds a single required parameter without variadic slice allocation.
func (b *FunctionParamBuilder) AddRequiredParam(type_ JSType) bool {
	if b.HasOptionalOrVarArgs() {
		return false
	}
	b.NewParameter(type_ /* isOptional= */, false /* isVariadic= */, false)
	return true
}

// AddRequiredParams adds parameters of the given type to the end of the
// param list. Returns false if this is called after optional params are
// added.
//
// Java: FunctionParamBuilder.addRequiredParams(JSType...)
func (b *FunctionParamBuilder) AddRequiredParams(types ...JSType) bool {
	if b.HasOptionalOrVarArgs() {
		return false
	}
	for _, type_ := range types {
		b.NewParameter(type_ /* isOptional= */, false /* isVariadic= */, false)
	}
	return true
}

// AddOptionalParam adds a single optional parameter without variadic slice allocation.
func (b *FunctionParamBuilder) AddOptionalParam(type_ JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	b.NewParameter(
		b.registry.CreateOptionalType(type_) /* isOptional= */, true /* isVariadic= */, false)
	return true
}

// AddOptionalParams adds optional parameters of the given type to the end
// of the param list. The builder will make them undefine-able. Returns false
// if this is called after var args are added.
//
// Java: FunctionParamBuilder.addOptionalParams(JSType...)
func (b *FunctionParamBuilder) AddOptionalParams(types ...JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	for _, type_ := range types {
		b.NewParameter(
			b.registry.CreateOptionalType(type_) /* isOptional= */, true /* isVariadic= */, false)
	}
	return true
}

// AddVarArgs adds variable arguments to the end of the parameter list.
// Returns false if this is called after var args are added.
//
// Java: FunctionParamBuilder.addVarArgs(JSType)
func (b *FunctionParamBuilder) AddVarArgs(type_ JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	b.NewParameter(type_ /* isOptional= */, false /* isVariadic= */, true)
	return true
}

// NewParameterFrom copies the existing parameter into this builder.
//
// Java: FunctionParamBuilder.newParameterFrom(Parameter)
func (b *FunctionParamBuilder) NewParameterFrom(n *FunctionType_Parameter) {
	b.ensureMutable()
	b.parameters = append(b.parameters, n)
}

// NewOptionalParameterFrom copies the parameter specification from the given
// parameter, but makes sure it's optional.
//
// Java: FunctionParamBuilder.newOptionalParameterFrom(Parameter)
func (b *FunctionParamBuilder) NewOptionalParameterFrom(p *FunctionType_Parameter) {
	isOptional := p.IsOptional() || !p.IsVariadic()
	if isOptional != p.IsOptional() {
		b.NewParameter(p.GetJSType(), isOptional, p.IsVariadic())
	} else {
		b.NewParameterFrom(p)
	}
}

// newParameter adds a parameter with the given type.
//
// Java: FunctionParamBuilder.NewParameter(JSType, boolean, boolean)
func (b *FunctionParamBuilder) NewParameter(type_ JSType, isOptional bool, isVariadic bool) {
	b.ensureMutable()
	b.parameters = append(b.parameters, FunctionType_Parameter_Create(type_, isOptional, isVariadic))
}

// Java: FunctionParamBuilder.HasOptionalOrVarArgs()
func (b *FunctionParamBuilder) HasOptionalOrVarArgs() bool {
	if len(b.parameters) == 0 {
		return false
	}
	lastParam := b.parameters[len(b.parameters)-1]
	return lastParam.IsOptional() || lastParam.IsVariadic()
}

// Java: FunctionParamBuilder.hasVarArgs()
func (b *FunctionParamBuilder) HasVarArgs() bool {
	return len(b.parameters) != 0 && b.parameters[len(b.parameters)-1].IsVariadic()
}

// Java: FunctionParamBuilder.build()
func (b *FunctionParamBuilder) Build() []*FunctionType_Parameter {
	if len(b.parameters) == 0 {
		return emptyFunctionTypeParams
	}
	b.consumed = true
	return b.parameters
}

// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20261021

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSDocParam is an ordered (Name, Type) entry for @param or @template.
type JSDocParam struct {
	Name string
	Type *JSTypeExpression
}

// JSDocTypeTransform is an ordered (Name, TypeNode) entry for @template T := ...
// type transformations.
type JSDocTypeTransform struct {
	Name     string
	TypeNode Node
}

// JSDocInfo_TypeTransform transforms a type expression.
//
// Go-only: language-adapter — a functional interface, so a Go func type.
//
// Java: JSDocInfo.TypeTransform.apply(JSTypeExpression)
type JSDocInfo_TypeTransform func(arg *JSTypeExpression) *JSTypeExpression

func formatJsDocParams(s []JSDocParam) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, p := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.Name)
		b.WriteByte('=')
		if p.Type == nil {
			b.WriteString("null")
		} else {
			b.WriteString(p.Type.String())
		}
	}
	b.WriteByte('}')
	return b.String()
}

func jsDocParamsToMap(s []JSDocParam) *jutil.LinkedHashMap[string, *JSTypeExpression] {
	m := jutil.NewLinkedHashMapWithCapacity[string, *JSTypeExpression](len(s))
	for _, p := range s {
		m.Put(p.Name, p.Type)
	}
	return m
}

func formatJsDocTypeTransforms(s []JSDocTypeTransform) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, tr := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(tr.Name)
		b.WriteByte('=')
		b.WriteString(jutil.ToJavaString(tr.TypeNode))
	}
	b.WriteByte('}')
	return b.String()
}

func jsDocTypeTransformsToMap(s []JSDocTypeTransform) *jutil.LinkedHashMap[string, Node] {
	m := jutil.NewLinkedHashMapWithCapacity[string, Node](len(s))
	for _, tr := range s {
		m.Put(tr.Name, tr.TypeNode)
	}
	return m
}

// Java: JSDocInfo.TypeMapProperty.areEquivalent(JSTypeExpression, JSTypeExpression)
func jsDocInfo_TypeMapProperty_AreEquivalent(expr1 *JSTypeExpression, expr2 *JSTypeExpression) bool {
	return expr1 == expr2 || (expr1 != nil && expr1.IsEquivalentTo(expr2))
}

// IsConstructor returns whether the `@constructor` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isConstructor()
func (i *JSDocInfo) IsConstructor() bool {
	return i.CheckBit(JSDocInfo_Bit_CONSTRUCTOR)
}

// IsAbstract returns whether the `@abstract` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isAbstract()
func (i *JSDocInfo) IsAbstract() bool {
	return i.CheckBit(JSDocInfo_Bit_ABSTRACT)
}

// UsesImplicitMatch returns whether the `@record` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.usesImplicitMatch()
func (i *JSDocInfo) UsesImplicitMatch() bool {
	return i.CheckBit(JSDocInfo_Bit_RECORD)
}

// MakesUnrestricted returns whether the `@unrestricted` annotation is present
// on this JSDocInfo.
//
// Java: JSDocInfo.makesUnrestricted()
func (i *JSDocInfo) MakesUnrestricted() bool {
	return i.CheckBit(JSDocInfo_Bit_UNRESTRICTED)
}

// MakesStructs returns whether the `@struct` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.makesStructs()
func (i *JSDocInfo) MakesStructs() bool {
	return i.CheckBit(JSDocInfo_Bit_STRUCT)
}

// MakesDicts returns whether the `@dict` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.makesDicts()
func (i *JSDocInfo) MakesDicts() bool {
	return i.CheckBit(JSDocInfo_Bit_DICT)
}

// IsInterface returns whether the `@interface` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isInterface()
func (i *JSDocInfo) IsInterface() bool {
	return (i.propertyBits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_RECORD.Mask())) != 0
}

// Java: JSDocInfo.isConstructorOrInterface()
func (i *JSDocInfo) IsConstructorOrInterface() bool {
	return (i.propertyBits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_RECORD.Mask())) != 0
}

// ContainsDeclarationExcludingTypelessConst returns whether there is a
// declaration present on this JSDocInfo.
//
// Does not consider `@const` (without a following type) to indicate a
// declaration. Whether you want this method, or the`containsDeclaration`
// that includes const, depends on whether you want to consider
// `/** @const * / a.b.c = 0` a declaration or not.
//
// Java: JSDocInfo.containsDeclarationExcludingTypelessConst()
func (i *JSDocInfo) ContainsDeclarationExcludingTypelessConst() bool {
	return (i.propertyBits&
		(JSDocInfo_Bit_CONSTRUCTOR.Mask()|
			JSDocInfo_Bit_DEFINE.Mask()|
			JSDocInfo_Bit_OVERRIDE.Mask()|
			JSDocInfo_Bit_EXPORT.Mask()|
			JSDocInfo_Bit_DEPRECATED.Mask()|
			JSDocInfo_Bit_INTERFACE.Mask()|
			JSDocInfo_Bit_IMPLICITCAST.Mask()|
			JSDocInfo_Bit_NOSIDEEFFECTS.Mask()|
			JSDocInfo_Bit_RECORD.Mask())) != 0 ||
		(i.propertyKeysBitset&
			(JSDocInfo_TYPE.mask|
				JSDocInfo_RETURN_TYPE.mask|
				JSDocInfo_ENUM_PARAMETER_TYPE.mask|
				JSDocInfo_TYPEDEF_TYPE.mask|
				JSDocInfo_THIS_TYPE.mask|
				JSDocInfo_PARAMETERS.mask|
				JSDocInfo_IMPLEMENTED_INTERFACES.mask|
				JSDocInfo_BASE_TYPE.mask|
				JSDocInfo_VISIBILITY.mask)) != 0
}

// ContainsTypeDeclaration returns whether this JSDoc contains a type
// declaration such as `/** @type {string}`.
//
// Java: JSDocInfo.containsTypeDeclaration()
func (i *JSDocInfo) ContainsTypeDeclaration() bool {
	return (i.propertyKeysBitset &
		(JSDocInfo_TYPE.mask |
			JSDocInfo_RETURN_TYPE.mask |
			JSDocInfo_ENUM_PARAMETER_TYPE.mask |
			JSDocInfo_TYPEDEF_TYPE.mask |
			JSDocInfo_THIS_TYPE.mask |
			JSDocInfo_PARAMETERS.mask |
			JSDocInfo_BASE_TYPE.mask)) != 0
}

// ContainsDeclaration returns whether there is a declaration present on this
// JSDocInfo, including a typeless @const like `/** @const * / a.b.c = 0`.
//
// Java: JSDocInfo.containsDeclaration()
func (i *JSDocInfo) ContainsDeclaration() bool {
	return i.ContainsDeclarationExcludingTypelessConst() || i.CheckBit(JSDocInfo_Bit_CONST)
}

// ContainsFunctionDeclaration returns whether there is a declaration of a
// callable type.
//
// Deprecated: This method is quite heuristic, looking for @type annotations
// that start with "function". Other methods like containsDeclaration() and
// containsTypeDefinition are generally preferred.
//
// Java: JSDocInfo.containsFunctionDeclaration()
func (i *JSDocInfo) ContainsFunctionDeclaration() bool {
	if i.CheckBit(JSDocInfo_Bit_CONSTRUCTOR) ||
		(i.propertyKeysBitset&(JSDocInfo_RETURN_TYPE.mask|JSDocInfo_THIS_TYPE.mask|JSDocInfo_PARAMETERS.mask)) != 0 {
		return true
	}
	type_ := jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.getRaw(i))
	if type_ != nil {
		return type_.GetRoot().IsFunction()
	}
	return i.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// ContainsTypeDefinition is for jsdocs that create new types. Not to be
// confused with jsdocs that declare the type of a variable or property.
//
// Java: JSDocInfo.containsTypeDefinition()
func (i *JSDocInfo) ContainsTypeDefinition() bool {
	return (i.propertyBits&(JSDocInfo_Bit_CONSTRUCTOR.Mask()|JSDocInfo_Bit_INTERFACE.Mask())) != 0 ||
		(i.propertyKeysBitset&(JSDocInfo_ENUM_PARAMETER_TYPE.mask|JSDocInfo_TYPEDEF_TYPE.mask)) != 0
}

// HasEnumParameterType returns whether an enum parameter type, specified
// using the `@enum` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasEnumParameterType()
func (i *JSDocInfo) HasEnumParameterType() bool {
	return (i.propertyKeysBitset & JSDocInfo_ENUM_PARAMETER_TYPE.mask) != 0
}

// HasTypedefType returns whether a typedef parameter type, specified using
// the `@typedef` annotation, is present on this JSDoc.
//
// Java: JSDocInfo.hasTypedefType()
func (i *JSDocInfo) HasTypedefType() bool {
	return (i.propertyKeysBitset & JSDocInfo_TYPEDEF_TYPE.mask) != 0
}

// HasReturnType returns whether this JSDocInfo contains a type for `@return`
// annotation.
//
// Java: JSDocInfo.hasReturnType()
func (i *JSDocInfo) HasReturnType() bool {
	return (i.propertyKeysBitset & JSDocInfo_RETURN_TYPE.mask) != 0
}

// HasType returns whether a type, specified using the `@type` annotation, is
// present on this JSDoc.
//
// Java: JSDocInfo.hasType()
func (i *JSDocInfo) HasType() bool {
	return (i.propertyKeysBitset & JSDocInfo_TYPE.mask) != 0
}

// Java: JSDocInfo.hasTypeInformation()
func (i *JSDocInfo) HasTypeInformation() bool {
	return (i.propertyKeysBitset &
		(JSDocInfo_TYPE.mask | JSDocInfo_TYPEDEF_TYPE.mask | JSDocInfo_ENUM_PARAMETER_TYPE.mask | JSDocInfo_RETURN_TYPE.mask)) != 0
}

// IsInlineType returns whether the type annotation was inlined.
//
// Java: JSDocInfo.isInlineType()
func (i *JSDocInfo) IsInlineType() bool {
	return i.CheckBit(JSDocInfo_Bit_INLINE_TYPE)
}

// GetReturnType gets the return type specified by the `@return` annotation.
//
// Java: JSDocInfo.getReturnType()
func (i *JSDocInfo) GetReturnType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_RETURN_TYPE.getRaw(i))
}

// GetEnumParameterType gets the enum parameter type specified by the `@enum`
// annotation.
//
// Java: JSDocInfo.getEnumParameterType()
func (i *JSDocInfo) GetEnumParameterType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_ENUM_PARAMETER_TYPE.getRaw(i))
}

// GetTypedefType gets the typedef type specified by the `@type` annotation.
//
// Java: JSDocInfo.getTypedefType()
func (i *JSDocInfo) GetTypedefType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPEDEF_TYPE.getRaw(i))
}

// GetType gets the type specified by the `@type` annotation.
//
// Java: JSDocInfo.getType()
func (i *JSDocInfo) GetType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_TYPE.getRaw(i))
}

// GetThisType gets the type specified by the `@this` annotation.
//
// Java: JSDocInfo.getThisType()
func (i *JSDocInfo) GetThisType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_THIS_TYPE.getRaw(i))
}

// HasThisType returns whether this JSDocInfo contains a type for `@this`
// annotation.
//
// Java: JSDocInfo.hasThisType()
func (i *JSDocInfo) HasThisType() bool {
	return (i.propertyKeysBitset & JSDocInfo_THIS_TYPE.mask) != 0
}

// GetBaseType gets the base type specified by the `@extends` annotation.
//
// Java: JSDocInfo.getBaseType()
func (i *JSDocInfo) GetBaseType() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_BASE_TYPE.getRaw(i))
}

// HasBaseType returns whether this JSDocInfo contains a type for `@extends`
// annotation.
//
// Java: JSDocInfo.hasBaseType()
func (i *JSDocInfo) HasBaseType() bool {
	return (i.propertyKeysBitset & JSDocInfo_BASE_TYPE.mask) != 0
}

// GetLendsName gets the name we're lending to in a `@lends` annotation.
//
// In many reflection APIs, you pass an anonymous object to a function, and
// that function mixes the anonymous object into another object. The `@lends`
// annotation allows the type system to track those property assignments.
//
// Java: JSDocInfo.getLendsName()
func (i *JSDocInfo) GetLendsName() *JSTypeExpression {
	return jsdocGet[*JSTypeExpression](JSDocInfo_LENDS_NAME.getRaw(i))
}

// Java: JSDocInfo.hasLendsName()
func (i *JSDocInfo) HasLendsName() bool {
	return (i.propertyKeysBitset & JSDocInfo_LENDS_NAME.mask) != 0
}

// GetTsTypes gets the ts type declarations specified by the `@tsType`
// annotations.
//
// Java: JSDocInfo.getTsTypes()
func (i *JSDocInfo) GetTsTypes() []string {
	if types := jsdocGet[[]string](JSDocInfo_TS_TYPES.getRaw(i)); types != nil {
		return slices.Clone(types)
	}
	return []string{}
}

// GetImplementedInterfaces returns the types specified by the `@implements`
// annotation: an immutable list of JSTypeExpression objects that can be
// resolved to types.
//
// Java: JSDocInfo.getImplementedInterfaces()
func (i *JSDocInfo) GetImplementedInterfaces() []*JSTypeExpression {
	return JSDocInfo_IMPLEMENTED_INTERFACES.GetUnmodifiable(i)
}

// GetImplementedInterfaceCount gets the number of interfaces specified by
// the `@implements` annotation.
//
// Java: JSDocInfo.getImplementedInterfaceCount()
func (i *JSDocInfo) GetImplementedInterfaceCount() int {
	return len(jsdocGet[[]*JSTypeExpression](JSDocInfo_IMPLEMENTED_INTERFACES.getRaw(i)))
}

// GetExtendedInterfaces returns the interfaces extended by an interface: an
// immutable list of JSTypeExpression objects that can be resolved to types.
//
// Java: JSDocInfo.getExtendedInterfaces()
func (i *JSDocInfo) GetExtendedInterfaces() []*JSTypeExpression {
	return JSDocInfo_EXTENDED_INTERFACES.GetUnmodifiable(i)
}

// GetExtendedInterfacesCount gets the number of extended interfaces
// specified.
//
// Java: JSDocInfo.getExtendedInterfacesCount()
func (i *JSDocInfo) GetExtendedInterfacesCount() int {
	return len(jsdocGet[[]*JSTypeExpression](JSDocInfo_EXTENDED_INTERFACES.getRaw(i)))
}

// GetParameterEntries returns the ordered @param (Name, Type) entries without allocating.
func (i *JSDocInfo) GetParameterEntries() []JSDocParam {
	return jsdocGet[[]JSDocParam](JSDocInfo_PARAMETERS.getRaw(i))
}

// GetParameterType gets the type of a given named parameter; null if this
// parameter is not defined or has a null type.
//
// Java: JSDocInfo.getParameterType(String)
func (i *JSDocInfo) GetParameterType(parameter string) *JSTypeExpression {
	for _, p := range i.GetParameterEntries() {
		if p.Name == parameter {
			return p.Type
		}
	}
	return nil
}

// HasParameter returns whether the parameter is defined.
//
// Java: JSDocInfo.hasParameter(String)
func (i *JSDocInfo) HasParameter(parameter string) bool {
	for _, p := range i.GetParameterEntries() {
		if p.Name == parameter {
			return true
		}
	}
	return false
}

// HasParameterType returns whether the parameter has an attached type:
// true if the parameter has an attached type, false if the parameter has no
// attached type or does not exist.
//
// Java: JSDocInfo.hasParameterType(String)
func (i *JSDocInfo) HasParameterType(parameter string) bool {
	return i.GetParameterType(parameter) != nil
}

// GetParameterNames returns the set of names of the defined parameters. The
// iteration order of the returned set is the order in which parameters are
// defined in the JSDoc, rather than the order in which the function declares
// them. The returned set is immutable.
//
// Java: JSDocInfo.getParameterNames()
func (i *JSDocInfo) GetParameterNames() *jutil.LinkedHashSet[string] {
	params := i.GetParameterEntries()
	if len(params) > 0 {
		set := jutil.NewLinkedHashSet[string]()
		for _, p := range params {
			set.Add(p.Name)
		}
		return set
	}
	return jutil.NewLinkedHashSet[string]()
}

// GetParameterNameAt returns the nth name in the defined parameters. The
// iteration order is the order in which parameters are defined in the
// JSDoc, rather than the order in which the function declares them.
//
// Java: JSDocInfo.getParameterNameAt(int)
func (i *JSDocInfo) GetParameterNameAt(index int) *string {
	name, ok := i.GetParameterNameValueAt(index)
	if !ok {
		return nil
	}
	return &name
}

// GetParameterNameValueAt returns the nth defined parameter name without
// allocating a *string pointer.
func (i *JSDocInfo) GetParameterNameValueAt(index int) (string, bool) {
	params := i.GetParameterEntries()
	if index < 0 || index >= len(params) {
		return "", false
	}
	return params[index].Name, true
}

// GetParameterTypeAt returns the nth defined parameter type expression in JSDoc
// definition order without allocating an iterator closure.
func (i *JSDocInfo) GetParameterTypeAt(index int) *JSTypeExpression {
	params := i.GetParameterEntries()
	if index < 0 || index >= len(params) {
		return nil
	}
	return params[index].Type
}

// GetParameterCount gets the number of parameters defined.
//
// Java: JSDocInfo.getParameterCount()
func (i *JSDocInfo) GetParameterCount() int {
	return len(i.GetParameterEntries())
}

// GetTemplateTypeEntries returns the ordered @template (Name, Type) entries without allocating.
func (i *JSDocInfo) GetTemplateTypeEntries() []JSDocParam {
	if i == nil {
		return nil
	}
	return jsdocGet[[]JSDocParam](JSDocInfo_TEMPLATE_TYPE_NAMES.getRaw(i))
}

// GetTemplateTypeNames gets the @template type names.
//
// Excludes @template types from TTL; get those with
// getTypeTransformations().
//
// Java: JSDocInfo.getTemplateTypeNames()
func (i *JSDocInfo) GetTemplateTypeNames() []string {
	params := i.GetTemplateTypeEntries()
	if len(params) == 0 {
		return []string{}
	}
	out := make([]string, len(params))
	for idx, p := range params {
		out[idx] = p.Name
	}
	return out
}

// HasTemplateTypeName reports whether name is declared as a @template type name.
func (i *JSDocInfo) HasTemplateTypeName(name string) bool {
	for _, p := range i.GetTemplateTypeEntries() {
		if p.Name == name {
			return true
		}
	}
	return false
}

// Java: JSDocInfo.getTemplateTypes()
func (i *JSDocInfo) GetTemplateTypes() *jutil.LinkedHashMap[string, *JSTypeExpression] {
	params := i.GetTemplateTypeEntries()
	if len(params) != 0 {
		return jsDocParamsToMap(params)
	}
	return jutil.NewLinkedHashMap[string, *JSTypeExpression]()
}

// GetTypeTransformationEntries returns the ordered @template T := ... entries
// without allocating a LinkedHashMap.
func (i *JSDocInfo) GetTypeTransformationEntries() []JSDocTypeTransform {
	if i == nil {
		return nil
	}
	return jsdocGet[[]JSDocTypeTransform](JSDocInfo_TYPE_TRANSFORMATIONS.getRaw(i))
}

// HasTypeTransformations reports whether any @template T := ... transformations
// are present without allocating.
func (i *JSDocInfo) HasTypeTransformations() bool {
	return i != nil && (i.propertyKeysBitset&JSDocInfo_TYPE_TRANSFORMATIONS.mask) != 0
}

// GetTypeTransformations gets the type transformations.
//
// Java: JSDocInfo.getTypeTransformations()
func (i *JSDocInfo) GetTypeTransformations() *jutil.LinkedHashMap[string, Node] {
	entries := i.GetTypeTransformationEntries()
	if len(entries) != 0 {
		return jsDocTypeTransformsToMap(entries)
	}
	return jutil.NewLinkedHashMap[string, Node]()
}

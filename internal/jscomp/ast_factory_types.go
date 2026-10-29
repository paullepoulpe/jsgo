// Port of com/google/javascript/jscomp/AstFactory.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstFactory.bigintNumberStringColor (Suppliers.memoize)
var AstFactory_bigintNumberStringColor = sync.OnceValue(func() *colors.Color {
	return colors.Color_CreateUnionSlice(
		[]*colors.Color{
			colors.StandardColors_BIGINT, colors.StandardColors_STRING, colors.StandardColors_NUMBER})
})

// AstFactory_TypeMode is AstFactory.TypeMode.
type AstFactory_TypeMode int32

const (
	AstFactory_TypeMode_JSTYPE AstFactory_TypeMode = iota
	AstFactory_TypeMode_COLOR
	AstFactory_TypeMode_NONE
)

var astFactory_TypeMode_names = [...]string{"JSTYPE", "COLOR", "NONE"}

// Java: AstFactory.TypeMode.name()
func (m AstFactory_TypeMode) String() string { return astFactory_TypeMode_names[m] }

// Java: AstFactory.TypeMode.ordinal()
func (m AstFactory_TypeMode) Ordinal() int { return int(m) }

// IsAddingTypes answers: does this class instance add types to the nodes it
// creates?
//
// Java: AstFactory.isAddingTypes()
func (a *AstFactory) IsAddingTypes() bool {
	return AstFactory_TypeMode_JSTYPE == a.typeMode
}

// IsAddingColors answers: does this class instance add optimization colors
// to the nodes it creates?
//
// Java: AstFactory.isAddingColors()
func (a *AstFactory) IsAddingColors() bool {
	return AstFactory_TypeMode_COLOR == a.typeMode
}

// Java: AstFactory.assertNotAddingColors()
func (a *AstFactory) AssertNotAddingColors() {
	jutil.CheckStateMsg(!a.IsAddingColors(), "method not supported for colors")
}

// Java: AstFactory.getTypeOfThisForFunctionNode(Node)
func (a *AstFactory) GetTypeOfThisForFunctionNode(functionNode rhino.Node) rhino.JSType {
	a.AssertNotAddingColors()
	if a.IsAddingTypes() {
		functionType := a.GetFunctionType(functionNode)
		return jutil.CheckNotNullMsg1(functionType.GetTypeOfThis(), "%s", functionType)
	} else {
		return nil // not adding type information
	}
}

// Java: AstFactory.getTypeOfThisForEs6Class(Node)
func (a *AstFactory) GetTypeOfThisForEs6Class(functionNode rhino.Node) AstFactory_Type {
	// Go-only: language-adapter — the JVM's "helpful" NullPointerException
	// message (jutil.NewHelpfulNullPointerException). An async method of an
	// object literal that calls super has no enclosing class, so
	// RewriteAsyncFunctions passes null here and Java dereferences it; the
	// INTERNAL COMPILER ERROR report prints the message the JVM derives from
	// this call. A Go nil dereference has none.
	if functionNode.IsNil() {
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "com.google.javascript.rhino.Node.isClass()" because "functionNode" is null`))
	}
	jutil.CheckArgumentMsg1(functionNode.IsClass(), "%s", functionNode)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		return AstFactory_Type_OfJSType(a.GetTypeOfThisForFunctionNode(functionNode))
	case AstFactory_TypeMode_COLOR:
		return AstFactory_Type_OfColor(a.GetInstanceOfColor(functionNode.GetColor()))
	case AstFactory_TypeMode_NONE:
		return AstFactory_NoTypeInformation()
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: AstFactory.getFunctionType(Node)
func (a *AstFactory) GetFunctionType(functionNode rhino.Node) rhino.FunctionType {
	jutil.CheckStateMsg1(
		functionNode.IsFunction() || functionNode.IsClass(),
		"not a function or class: %s",
		functionNode)
	a.AssertNotAddingColors()
	// If the function declaration was cast to a different type, we want the original type
	// from before the cast.
	typeBeforeCast := functionNode.GetJSTypeBeforeCast()
	var functionType rhino.FunctionType
	if typeBeforeCast != nil {
		functionType = typeBeforeCast.AssertFunctionType()
	} else {
		functionType = functionNode.GetJSTypeRequired().AssertFunctionType()
	}
	return functionType
}

// astFactoryObjectTypeAsJSType widens an ObjectType to a JSType.
//
// Go-only: language-adapter — keeps a nil ObjectType a nil JSType.
func astFactoryObjectTypeAsJSType(t rhino.ObjectType) rhino.JSType {
	if t == nil {
		return nil
	}
	return t
}

// Java: AstFactory.getInstanceOfColor(Color)
func (a *AstFactory) GetInstanceOfColor(color *colors.Color) *colors.Color {
	possibleInstanceColors := color.GetInstanceColors()
	if len(possibleInstanceColors) == 0 {
		return colors.StandardColors_UNKNOWN
	}
	return colors.Color_CreateUnionSlice(possibleInstanceColors)
}

// Java: AstFactory.replaceTemplate(JSType, ImmutableList<JSType>)
func (a *AstFactory) ReplaceTemplate(templatedType rhino.JSType, templateTypes []rhino.JSType) rhino.JSType {
	typeMap :=
		a.registry.
			GetEmptyTemplateTypeMap().
			CopyWithExtensionOfKeys(templatedType.GetTemplateTypeMap().GetTemplateKeys(), templateTypes)
	replacer := rhino.TemplateTypeReplacer_MakeForPartialReplacement(a.registry, typeMap)
	return replacer.VisitJSType(templatedType)
}

// Java: AstFactory.getNativeType(JSTypeNative)
func (a *AstFactory) GetNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	jutil.CheckNotNullMsg(a.registry, "registry is null")
	return jutil.CheckNotNullMsg1(
		a.registry.GetNativeType(nativeType), "native type not found: %s", nativeType)
}

// Java: AstFactory.getJsTypeForProperty(Node, String)
func (a *AstFactory) GetJsTypeForProperty(receiver rhino.Node, propertyName string) rhino.JSType {
	// NOTE: we use both findPropertyType and getPropertyType because they are subtly
	// different: findPropertyType works on JSType, autoboxing scalars and joining unions,
	// but it returns null if the type is not found and does not handle dynamic types of
	// Function.prototype.call and .apply; whereas getPropertyType does not autobox nor
	// iterate over unions, but it does synthesize the function properties correctly, and
	// it returns unknown instead of null if the property is missing.
	var getpropType rhino.JSType = nil
	receiverJSType := receiver.GetJSType()
	if receiverJSType != nil {
		getpropType = receiverJSType.FindPropertyType(propertyName)
		if getpropType == nil {
			receiverObjectType := rhino.ObjectType_Cast(receiverJSType.Autobox())
			if receiverObjectType == nil {
				getpropType = a.unknownType
			} else {
				getpropType = receiverObjectType.GetPropertyType(propertyName)
			}
		}
	}
	if getpropType == nil {
		getpropType = a.unknownType
	}
	// TODO(bradfordcsmith): Special case $jscomp.global until we annotate its type correctly.
	if getpropType.IsUnknownType() &&
		propertyName == "global" &&
		receiver.MatchesName("$jscomp") {
		getpropType = a.GetNativeType(rhino.JSTypeNative_GLOBAL_THIS)
	}
	return getpropType
}

// Java: AstFactory.setJSTypeOrColor(AstFactory.Type, Node)
func (a *AstFactory) SetJSTypeOrColor(type_ AstFactory_Type, result rhino.Node) {
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(type_.GetJSType(a.registry))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(type_.GetColor(a.colorRegistry))
	case AstFactory_TypeMode_NONE:
	}
}

type astFactoryTypeKind uint8

const (
	astFactoryTypeZero     astFactoryTypeKind = 0
	astFactoryTypeFromNode astFactoryTypeKind = 1 << iota
	astFactoryTypeExplicit
	astFactoryTypeHasJSTypeNative
	astFactoryTypeHasColorId
)

// AstFactory_Type is AstFactory.Type, represented as a zero-allocation value.
type AstFactory_Type struct {
	jstype       rhino.JSType
	color        *colors.Color
	colorId      colors.ColorId
	node         rhino.Node
	jstypeNative rhino.JSTypeNative
	kind         astFactoryTypeKind
}

// IsZero reports whether t is the zero AstFactory_Type value.
func (t AstFactory_Type) IsZero() bool {
	return t.kind == astFactoryTypeZero
}

// Java: AstFactory.Type.getJSType(JSTypeRegistry)
func (t AstFactory_Type) GetJSType(registry *rhino.JSTypeRegistry) rhino.JSType {
	if t.kind == astFactoryTypeFromNode {
		jstype := t.node.GetJSType()
		// TODO(b/149843534): crash instead of defaulting to unknown
		if jstype != nil {
			return jstype
		}
		return registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	if t.jstype != nil {
		return t.jstype
	}
	if t.kind&astFactoryTypeHasJSTypeNative == 0 {
		// Go-only: language-adapter — registry.getNativeType(null) throws a
		// NullPointerException (typeId.ordinal()).
		panic(jutil.NewNullPointerException())
	}
	return registry.GetNativeType(t.jstypeNative)
}

// Java: AstFactory.Type.getColor(ColorRegistry)
func (t AstFactory_Type) GetColor(registry *colors.ColorRegistry) *colors.Color {
	if t.kind == astFactoryTypeFromNode {
		color := t.node.GetColor()
		// TODO(b/149843534): crash instead of defaulting to unknown
		if color != nil {
			return color
		}
		return colors.StandardColors_UNKNOWN
	}
	if t.color != nil {
		return t.color
	}
	if t.kind&astFactoryTypeHasColorId == 0 {
		panic(jutil.NewNullPointerException())
	}
	return registry.Get(t.colorId)
}

// AstFactory_Type_ uses the JSType or Color of the given node as a template
// for adding type information
//
// Java: AstFactory.type(Node)
func AstFactory_Type_(node rhino.Node) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeFromNode, node: node}
}

// Java: AstFactory.type(JSType)
func AstFactory_Type_OfJSType(type_ rhino.JSType) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, jstype: type_}
}

// Java: AstFactory.type(JSTypeNative)
func AstFactory_Type_OfJSTypeNative(type_ rhino.JSTypeNative) AstFactory_Type {
	return AstFactory_Type{
		kind:         astFactoryTypeExplicit | astFactoryTypeHasJSTypeNative,
		jstypeNative: type_,
	}
}

// Java: AstFactory.type(Color)
func AstFactory_Type_OfColor(type_ *colors.Color) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, color: type_}
}

// Java: AstFactory.type(ColorId)
func AstFactory_Type_OfColorId(type_ colors.ColorId) AstFactory_Type {
	return AstFactory_Type{
		kind:    astFactoryTypeExplicit | astFactoryTypeHasColorId,
		colorId: type_,
	}
}

// Java: AstFactory.type(JSType, Color)
func AstFactory_Type_OfJSTypeAndColor(type_ rhino.JSType, color *colors.Color) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, jstype: type_, color: color}
}

// Java: AstFactory.type(JSTypeNative, Color)
func AstFactory_Type_OfJSTypeNativeAndColor(type_ rhino.JSTypeNative, color *colors.Color) AstFactory_Type {
	return AstFactory_Type{
		kind:         astFactoryTypeExplicit | astFactoryTypeHasJSTypeNative,
		jstypeNative: type_,
		color:        color,
	}
}

// Java: AstFactory.noTypeInformation()
func AstFactory_NoTypeInformation() AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit}
}

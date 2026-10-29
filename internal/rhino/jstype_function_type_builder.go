// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20261021

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

// Java: FunctionType.builder(JSTypeRegistry)
func FunctionType_Builder_(registry *JSTypeRegistry) *FunctionType_Builder {
	return NewFunctionType_Builder(registry)
}

func (f *FunctionTypeBase) initBuilder(b *FunctionType_Builder) {
	b.registry = f.GetRegistry()
	b.implicitPrototype = jutil.CheckNotNull(f.GetRegistry().GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE))
	b.googModuleID = rhinoStringIDFromPtr(f.GetGoogModuleId())
	b.classNameID = rhinoStringIDFromPtr(f.GetReferenceName())
	b.nativeType = f.IsNativeObjectType()
	b.templateTypeMap = f.GetTemplateTypeMap()
	b.canonicalRepresentation = f.GetCanonicalRepresentation()
	b.isKnownAmbiguous = f.call.fnConstructorAmbiguity == FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	if p := f.GetClosurePrimitive(); p != nil {
		b.primitivePlus1 = uint8(*p) + 1
	} else {
		b.primitivePlus1 = 0
	}
	b.isAbstract = f.IsAbstract()
	b.kind = f.GetKind()
	b.parameters = f.GetParameters()
	b.returnType = f.GetReturnType()
	b.returnTypeIsInferred = f.IsReturnTypeInferred()
	b.sourceNode = f.GetSource()
	b.typeOfThis = f.GetTypeOfThis()
	b.templateParamCount = f.GetTemplateParamCount()

	if f.IsConstructor() {
		b.WithConstructorTemplateKeys(f.GetConstructorOnlyTemplateParameters())
	}
}

// ToBuilder copies all the information from another function type.
//
// Java: FunctionType.toBuilder()
func (f *FunctionTypeBase) ToBuilder() *FunctionType_Builder {
	b := f.GetRegistry().acquireFunctionBuilder()
	f.initBuilder(b)
	return b
}

// FunctionType_Builder is a builder class for function and arrow types.
type FunctionType_Builder struct {
	PrototypeObjectType_Builder

	// Go: nil is Java null (unknown parameters).
	parameters              []*FunctionType_Parameter
	returnType              JSType
	typeOfThis              JSType
	setPrototypeBasedOn     ObjectType
	constructorOnlyKeys     []*TemplateType
	canonicalRepresentation FunctionType
	sourceNode              Node
	googModuleID            rhinoStringID
	kind                    FunctionType_Kind
	primitivePlus1          uint8
	isAbstract              bool
	isKnownAmbiguous        bool
	returnTypeIsInferred    bool
	returnsOwnInstanceType  bool
	pooled                  bool
}

// Java: FunctionType.Builder(JSTypeRegistry)
func NewFunctionType_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	b := registry.acquireFunctionBuilder()
	b.InitPrototypeObjectType_Builder(registry)
	b.kind = FunctionType_Kind_ORDINARY
	b.implicitPrototype = jutil.CheckNotNull(registry.GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE))
	return b
}

// Java: PrototypeObjectType.Builder.setName(String)
func (b *FunctionType_Builder) SetName(x *string) *FunctionType_Builder {
	b.classNameID = rhinoStringIDFromPtr(x)
	return b
}

// Java: PrototypeObjectType.Builder.setImplicitPrototype(ObjectType)
func (b *FunctionType_Builder) SetImplicitPrototype(x ObjectType) *FunctionType_Builder {
	b.implicitPrototype = x
	return b
}

// Java: PrototypeObjectType.Builder.setNative(boolean)
func (b *FunctionType_Builder) SetNative(x bool) *FunctionType_Builder {
	b.nativeType = x
	return b
}

// Java: PrototypeObjectType.Builder.setAnonymous(boolean)
func (b *FunctionType_Builder) SetAnonymous(x bool) *FunctionType_Builder {
	b.anonymousType = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateTypeMap(TemplateTypeMap)
func (b *FunctionType_Builder) SetTemplateTypeMap(x *TemplateTypeMap) *FunctionType_Builder {
	b.templateTypeMap = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateParamCount(int)
func (b *FunctionType_Builder) SetTemplateParamCount(x int) *FunctionType_Builder {
	b.templateParamCount = x
	return b
}

// WithName sets the name of the function type.
//
// Java: FunctionType.Builder.withName(String)
func (b *FunctionType_Builder) WithName(name *string) *FunctionType_Builder {
	b.classNameID = rhinoStringIDFromPtr(name)
	return b
}

// WithNameString sets the non-null name of the function type without allocating a *string.
func (b *FunctionType_Builder) WithNameString(name string) *FunctionType_Builder {
	b.classNameID = rhinoStringPoolIntern(name)
	return b
}

// WithSourceNode sets the source node of the function type.
//
// Java: FunctionType.Builder.withSourceNode(Node)
func (b *FunctionType_Builder) WithSourceNode(sourceNode Node) *FunctionType_Builder {
	b.sourceNode = sourceNode
	return b
}

// SetGoogModuleId sets the ID of the goog.module in which this enum was
// declared.
//
// Java: FunctionType.Builder.setGoogModuleId(String)
func (b *FunctionType_Builder) SetGoogModuleId(x *string) *FunctionType_Builder {
	b.googModuleID = rhinoStringIDFromPtr(x)
	return b
}

// WithParametersOfList sets the parameters of the function type.
//
// Java: FunctionType.Builder.withParameters(List<Parameter>)
func (b *FunctionType_Builder) WithParametersOfList(parameters []*FunctionType_Parameter) *FunctionType_Builder {
	b.parameters = parameters
	return b
}

// WithParams sets the parameters of the function type from a FunctionParamBuilder.
func (b *FunctionType_Builder) WithParams(paramBuilder *FunctionParamBuilder) *FunctionType_Builder {
	b.parameters = paramBuilder.Build()
	return b
}

// WithParameters sets the function to take zero parameters.
//
// Java: FunctionType.Builder.withParameters()
func (b *FunctionType_Builder) WithParameters() *FunctionType_Builder {
	b.parameters = emptyFunctionTypeParameterSlice
	return b
}

// WithReturnType sets the return type.
//
// Java: FunctionType.Builder.withReturnType(JSType)
func (b *FunctionType_Builder) WithReturnType(returnType JSType) *FunctionType_Builder {
	b.returnType = returnType
	return b
}

// WithReturnTypeWithInferred sets the return type and whether it's
// inferred.
//
// Java: FunctionType.Builder.withReturnType(JSType, boolean)
func (b *FunctionType_Builder) WithReturnTypeWithInferred(returnType JSType, inferred bool) *FunctionType_Builder {
	b.returnType = returnType
	b.returnTypeIsInferred = inferred
	return b
}

// WithReturnsOwnInstanceType sets the return type to be a constructor's own
// instance type.
//
// Java: FunctionType.Builder.withReturnsOwnInstanceType()
func (b *FunctionType_Builder) WithReturnsOwnInstanceType() *FunctionType_Builder {
	b.returnsOwnInstanceType = true
	return b
}

// WithInferredReturnType sets an inferred return type.
//
// Java: FunctionType.Builder.withInferredReturnType(JSType)
func (b *FunctionType_Builder) WithInferredReturnType(returnType JSType) *FunctionType_Builder {
	b.returnType = returnType
	b.returnTypeIsInferred = true
	return b
}

// WithTypeOfThis sets the "this" type.
//
// Java: FunctionType.Builder.withTypeOfThis(JSType)
func (b *FunctionType_Builder) WithTypeOfThis(typeOfThis JSType) *FunctionType_Builder {
	b.typeOfThis = typeOfThis
	return b
}

// WithTemplateKeysOfList sets the template name.
//
// Java: FunctionType.Builder.withTemplateKeys(ImmutableList<TemplateType>)
func (b *FunctionType_Builder) WithTemplateKeysOfList(templateKeys []*TemplateType) *FunctionType_Builder {
	if templateKeys == nil {
		templateKeys = []*TemplateType{}
	}

	b.templateTypeMap = b.registry.
		GetEmptyTemplateTypeMap().
		CopyWithExtensionOfKeys(templateKeys, []JSType{})
	// TODO(nickreid): This value should only consider ctor only keys.
	b.templateParamCount = len(templateKeys)
	return b
}

// WithTemplateKeys sets the template name.
//
// Java: FunctionType.Builder.withTemplateKeys(TemplateType...)
func (b *FunctionType_Builder) WithTemplateKeys(templateKeys ...*TemplateType) *FunctionType_Builder {
	return b.WithTemplateKeysOfList(jsTypeCopyOfTemplateTypes(templateKeys))
}

// jsTypeCopyOfTemplateTypes is ImmutableList.copyOf(array), which rejects
// null elements.
func jsTypeCopyOfTemplateTypes(keys []*TemplateType) []*TemplateType {
	out := make([]*TemplateType, len(keys))
	for i, k := range keys {
		out[i] = jutil.CheckNotNull(k)
	}
	return out
}

// WithConstructorTemplateKeys specifies a subset of the template keys that
// only apply to the constructor, and should be removed from the instance
// type. These keys must still be passed to #withTemplateKeys.
//
// Java: FunctionType.Builder.withConstructorTemplateKeys(Iterable<TemplateType>)
func (b *FunctionType_Builder) WithConstructorTemplateKeys(constructorOnlyKeys []*TemplateType) *FunctionType_Builder {
	if len(constructorOnlyKeys) == 0 {
		b.constructorOnlyKeys = nil
		return b
	}
	out := make([]*TemplateType, 0, len(constructorOnlyKeys))
	for _, k := range constructorOnlyKeys {
		jutil.CheckNotNull(k)
		dup := false
		for _, existing := range out {
			if existing == k || existing.Equals(k) {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, k)
		}
	}
	b.constructorOnlyKeys = out
	return b
}

// WithKind sets the function kind.
//
// Java: FunctionType.Builder.withKind(Kind)
func (b *FunctionType_Builder) WithKind(kind FunctionType_Kind) *FunctionType_Builder {
	b.kind = kind
	return b
}

// ForConstructor makes this a constructor.
//
// Java: FunctionType.Builder.forConstructor()
func (b *FunctionType_Builder) ForConstructor() *FunctionType_Builder {
	b.kind = FunctionType_Kind_CONSTRUCTOR
	return b
}

// ForInterface makes this an interface.
//
// Java: FunctionType.Builder.forInterface()
func (b *FunctionType_Builder) ForInterface() *FunctionType_Builder {
	b.kind = FunctionType_Kind_INTERFACE
	b.parameters = []*FunctionType_Parameter{}
	return b
}

// ForNativeType makes this a native type.
//
// Java: FunctionType.Builder.forNativeType()
func (b *FunctionType_Builder) ForNativeType() *FunctionType_Builder {
	b.nativeType = true
	return b
}

// WithIsAbstract marks abstract method.
//
// Java: FunctionType.Builder.withIsAbstract(boolean)
func (b *FunctionType_Builder) WithIsAbstract(isAbstract bool) *FunctionType_Builder {
	b.isAbstract = isAbstract
	return b
}

// Java: FunctionType.Builder.setIsKnownAmbiguous(boolean)
func (b *FunctionType_Builder) SetIsKnownAmbiguous(x bool) *FunctionType_Builder {
	b.isKnownAmbiguous = x
	return b
}

// WithPrototypeBasedOn sets the prototype property of a constructor.
//
// Java: FunctionType.Builder.withPrototypeBasedOn(ObjectType)
func (b *FunctionType_Builder) WithPrototypeBasedOn(setPrototypeBasedOn ObjectType) *FunctionType_Builder {
	b.setPrototypeBasedOn = setPrototypeBasedOn
	return b
}

// WithClosurePrimitiveId sets the ClosurePrimitive corresponding to this
// function
//
// Java: FunctionType.Builder.withClosurePrimitiveId(ClosurePrimitive)
func (b *FunctionType_Builder) WithClosurePrimitiveId(id *ClosurePrimitive) *FunctionType_Builder {
	if id == nil {
		b.primitivePlus1 = 0
	} else {
		b.primitivePlus1 = uint8(*id) + 1
	}
	return b
}

// withCanonicalRepresentation sets the canonical representation of a
// constructor, if any
//
// Java: FunctionType.Builder.WithCanonicalRepresentation(FunctionType)
func (b *FunctionType_Builder) WithCanonicalRepresentation(representation FunctionType) *FunctionType_Builder {
	b.canonicalRepresentation = representation
	return b
}

// Build constructs a new function type.
//
// Java: FunctionType.Builder.build()
func (b *FunctionType_Builder) Build() FunctionType {
	// Verify that the builder is an a sensible state before instantiating a
	// function.

	switch b.kind {
	case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
		// These kinds have no implication on whether `returnsOwnInstanceType`
		// is reasonable. This configuration may be intended to synthesize an
		// instance type. The return type and instance type are independent.
	case FunctionType_Kind_NONE:
		jutil.CheckState(b.returnsOwnInstanceType)
	case FunctionType_Kind_ORDINARY:
		jutil.CheckState(!b.returnsOwnInstanceType)
	}

	if b.returnsOwnInstanceType {
		// If the return type or instance type was available to set, there's no
		// need to use `returnsOwnInstanceType`.
		jutil.CheckState(b.returnType == nil)
		jutil.CheckState(b.typeOfThis == nil)
	}

	switch b.kind {
	case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
	case FunctionType_Kind_NONE, FunctionType_Kind_ORDINARY:
		jutil.CheckState(len(b.constructorOnlyKeys) == 0)
	}

	res := NewFunctionType(b)
	b.registry.releaseFunctionBuilder(b)
	return res
}

// Java: FunctionType.Builder.buildAndResolve()
func (b *FunctionType_Builder) BuildAndResolve() FunctionType {
	return jutil.CheckNotNull(b.Build().ToMaybeFunctionType())
}

// FunctionType_Parameter models a single JavaScript parameter.
//
// This parameter has a type; optionality; and may be var_args (variadic).
//
// Java: abstract class FunctionType.Parameter, with its generated
// AutoValue_FunctionType_Parameter implementation (fields, equals, hashCode,
// toString) folded in.
type FunctionType_Parameter struct {
	// Go-only: performance — the type's base (8 bytes) rather than a JSType
	// interface (16); GetJSType returns its self().
	jSType   *JSTypeBase
	optional bool
	variadic bool
}

// Java: FunctionType.Parameter.create(JSType, boolean, boolean)
func FunctionType_Parameter_Create(type_ JSType, isOptional bool, isVariadic bool) *FunctionType_Parameter {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null JSType"))
	}
	return &FunctionType_Parameter{jSType: type_.GetJSTypeBase(), optional: isOptional, variadic: isVariadic}
}

// Java: FunctionType.Parameter.getJSType()
func (p *FunctionType_Parameter) GetJSType() JSType {
	return p.jSType.self()
}

// Java: FunctionType.Parameter.isOptional()
func (p *FunctionType_Parameter) IsOptional() bool {
	return p.optional
}

// Java: FunctionType.Parameter.isVariadic()
func (p *FunctionType_Parameter) IsVariadic() bool {
	return p.variadic
}

// Java: AutoValue_FunctionType_Parameter.toString()
func (p *FunctionType_Parameter) String() string {
	return "Parameter{" +
		"JSType=" + p.jSType.self().String() + ", " +
		"optional=" + jutil.ToJavaString(p.optional) + ", " +
		"variadic=" + jutil.ToJavaString(p.variadic) +
		"}"
}

// Java: AutoValue_FunctionType_Parameter.equals(Object)
func (p *FunctionType_Parameter) Equals(o any) bool {
	if any(p) == o {
		return true
	}
	if that, ok := o.(*FunctionType_Parameter); ok && that != nil {
		return p.jSType.self().Equals(that.GetJSType()) &&
			p.optional == that.IsOptional() &&
			p.variadic == that.IsVariadic()
	}
	return false
}

// Java: AutoValue_FunctionType_Parameter.hashCode()
func (p *FunctionType_Parameter) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= p.jSType.self().HashCode()
	h *= 1000003
	h ^= jutil.ObjectHashCode(p.optional)
	h *= 1000003
	h ^= jutil.ObjectHashCode(p.variadic)
	return h
}

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FunctionType is the derived type that provides extended information about
// a function, including its return type and argument types.
//
// Note: the parameters list is the PARAM_LIST node that is the parent of the
// actual NAME node containing the parsed argument list (annotated with
// JSDOC_TYPE_PROP's for the compile-time type of each argument.
//
// Java: class FunctionType extends PrototypeObjectType implements
// JSType.WithSourceRef (subclassed by NoObjectType).
type FunctionType = *FunctionTypeBase

type functionTypeExt struct {
	// The prototype property. This field is lazily initialized by
	// #getPrototype().
	prototypeSlot *Property

	// The interfaces directly implemented by this function (for
	// constructors) It is only relevant for constructors.
	implementedInterfaces []ObjectType

	// The interfaces directly extended by this function (for interfaces) It
	// is only relevant for constructors.
	extendedInterfaces []ObjectType

	// If non-null, the original canonical variant of this function; only used
	// for constructors.
	canonicalRepresentation FunctionType

	// The ID of the goog.module in which this function was declared.
	googModuleID rhinoStringID

	// 1 + ClosurePrimitive ordinal, or 0 if none.
	closurePrimitivePlus1 uint8
}

// FunctionTypeBase holds FunctionType's fields and methods; plain function
// types as well as NoObjectType and NoType are *FunctionTypeBase values
// distinguished by objFlags & objKindMask.
type FunctionTypeBase struct {
	PrototypeObjectTypeBase

	// [[Call]] property (also packs 1-byte FunctionTypeBase flags in its tail padding).
	call ArrowType

	// The type of this in the scope of this function.
	//
	// Go-only: performance — the type's base (8 bytes) rather than a JSType
	// interface (16).
	typeOfThis *JSTypeBase

	// The function node which this type represents. It may be null.
	source Node

	ext *functionTypeExt
}

func (f *FunctionTypeBase) isNo() bool {
	return f.objFlags&objKindMask == objKindNo
}

func (f *FunctionTypeBase) isNoObjectOrNo() bool {
	k := f.objFlags & objKindMask
	return k == objKindNoObject || k == objKindNo
}

func (f *FunctionTypeBase) ensureExt() *functionTypeExt {
	if f.ext == nil {
		f.ext = &functionTypeExt{}
	}
	return f.ext
}

func (f *FunctionTypeBase) prototypeSlot() *Property {
	if f.ext != nil {
		return f.ext.prototypeSlot
	}
	return nil
}

func (f *FunctionTypeBase) setPrototypeSlot(slot *Property) {
	if slot != nil || f.ext != nil {
		f.ensureExt().prototypeSlot = slot
	}
}

// Java: FunctionType.TYPE_CLASS
const FunctionType_TYPE_CLASS = JSTypeClass_FUNCTION

// FunctionType_Kind is FunctionType.Kind.
type FunctionType_Kind uint8

const (
	FunctionType_Kind_ORDINARY FunctionType_Kind = iota
	FunctionType_Kind_CONSTRUCTOR
	FunctionType_Kind_INTERFACE
	FunctionType_Kind_NONE
)

var functionType_Kind_names = [...]string{"ORDINARY", "CONSTRUCTOR", "INTERFACE", "NONE"}

// String is Enum.name().
func (k FunctionType_Kind) String() string { return functionType_Kind_names[k] }

// FunctionType_PropAccess is relevant only for constructors.
type FunctionType_PropAccess uint8

const (
	// An implicit any behavior (default of ES5 classes)
	FunctionType_PropAccess_ANY FunctionType_PropAccess = iota
	// An explicit @unrestricted tag
	FunctionType_PropAccess_ANY_EXPLICIT
	// An explicit @struct or the implicit behavior of ES6 classes
	FunctionType_PropAccess_STRUCT
	// An explicit @dict or the implicit behavior of ES5 classes
	FunctionType_PropAccess_DICT
)

// FunctionType_ConstructorAmbiguity is FunctionType.ConstructorAmbiguity.
type FunctionType_ConstructorAmbiguity uint8

const (
	FunctionType_ConstructorAmbiguity_UNKNOWN FunctionType_ConstructorAmbiguity = iota
	FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
)

var functionType_ConstructorAmbiguity_names = [...]string{
	"UNKNOWN", "IS_AMBIGUOUS_CONSTRUCTOR", "IS_UNAMBIGUOUS_CONSTRUCTOR",
}

// String is Enum.name().
func (a FunctionType_ConstructorAmbiguity) String() string {
	return functionType_ConstructorAmbiguity_names[a]
}

var (
	emptyObjectTypeSlice            = []ObjectType{}
	emptyFunctionTypeParameterSlice = []*FunctionType_Parameter{}
)

// NewFunctionType creates an instance for a function that might be a
// constructor.
//
// Non-subclasses must go through Builder to create a new FunctionType.
//
// Java: FunctionType(Builder)
func NewFunctionType(builder *FunctionType_Builder) *FunctionTypeBase {
	t := &FunctionTypeBase{}
	t.objFlags = objKindFunction
	t.InitFunctionTypeBase(t, builder)
	return t
}

// InitFunctionTypeBase is the Go form of the constructor for subclasses.
//
// Java: FunctionType(Builder)
func (f *FunctionTypeBase) InitFunctionTypeBase(self FunctionType, builder *FunctionType_Builder) {
	f.InitPrototypeObjectTypeBase(self, &builder.PrototypeObjectType_Builder)
	self.SetPrettyPrint(true)

	source := builder.sourceNode
	jutil.CheckArgument(source.IsNil() || source.IsFunction() || source.IsClass())
	f.source = source
	if builder.googModuleID != 0 {
		f.ensureExt().googModuleID = builder.googModuleID
	}
	f.call.fnKind = builder.kind
	if builder.isKnownAmbiguous {
		f.call.fnConstructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	} else {
		f.call.fnConstructorAmbiguity = FunctionType_ConstructorAmbiguity_UNKNOWN
	}

	if builder.typeOfThis != nil {
		f.typeOfThis = jsTypeBase(builder.typeOfThis)
	} else {
		switch f.call.fnKind {
		case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
			var typeOfThisBuilder InstanceObjectType_Builder
			typeOfThisBuilder.registry = f.GetRegistry()
			typeOfThisBuilder.classNameID = rhinoStringIDFromPtr(self.GetReferenceName())
			typeOfThisBuilder.nativeType = self.IsNativeObjectType()
			typeOfThisBuilder.anonymousType = self.IsAnonymous()
			typeOfThisBuilder.templateTypeMap = self.GetTemplateTypeMap()
			typeOfThisBuilder.templateParamCount = self.GetTemplateParamCount()
			typeOfThisBuilder.constructor = self

			ctorKeys := builder.constructorOnlyKeys
			if len(ctorKeys) > 0 {
				typeOfThisBuilder.templateTypeMap = f.templateTypeMap.copyWithoutKeySlice(ctorKeys)
				typeOfThisBuilder.templateParamCount = f.GetTemplateParamCount() - len(ctorKeys)
			}

			f.typeOfThis = NewInstanceObjectType(&typeOfThisBuilder).GetJSTypeBase()
		case FunctionType_Kind_ORDINARY:
			f.typeOfThis = f.GetRegistry().GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE).GetJSTypeBase()
		case FunctionType_Kind_NONE:
			f.typeOfThis = self.GetJSTypeBase()
		}
	}

	if f.call.fnKind == FunctionType_Kind_CONSTRUCTOR {
		f.call.fnPropAccess = uint8(FunctionType_PropAccess_ANY) + 1
	}

	var returnType JSType
	if builder.returnsOwnInstanceType {
		returnType = jsTypeOf(f.typeOfThis)
	} else {
		returnType = builder.returnType
	}
	f.call.init(f.GetRegistry(), builder.parameters, returnType, builder.returnTypeIsInferred)
	if f.GetRegistry().GetResolver().state != JSTypeResolver_State_OPEN {
		f.call.resolveInternal(f.GetRegistry().GetErrorReporter())
	}

	if builder.primitivePlus1 != 0 {
		f.ensureExt().closurePrimitivePlus1 = builder.primitivePlus1
	}
	f.call.fnIsStructuralInterface = false
	f.call.fnIsAbstract = builder.isAbstract
	canonicalRepresentation := builder.canonicalRepresentation
	jutil.CheckArgumentMsg(
		canonicalRepresentation == nil || f.call.fnKind == FunctionType_Kind_CONSTRUCTOR,
		"Only constructors should have canonical representations")
	if canonicalRepresentation != nil {
		f.ensureExt().canonicalRepresentation = canonicalRepresentation
	}

	if builder.setPrototypeBasedOn != nil {
		f.SetPrototypeBasedOn(builder.setPrototypeBasedOn)
	}

	f.GetRegistry().GetResolver().ResolveIfClosed(self, FunctionType_TYPE_CLASS)
}

// Java: FunctionType.getTypeClass() / NoObjectType.getTypeClass() / NoType.getTypeClass()
func (f *FunctionTypeBase) GetTypeClass() JSTypeClass {
	switch f.objFlags & objKindMask {
	case objKindNo:
		return JSTypeClass_NO
	case objKindNoObject:
		return JSTypeClass_NO_OBJECT
	default:
		return FunctionType_TYPE_CLASS
	}
}

// Java: FunctionType.getConstructor() / NoObjectType.getConstructor()
func (f *FunctionTypeBase) GetConstructor() FunctionType {
	if f.isNoObjectOrNo() {
		return nil
	}
	// Every function type, including `Function`, is constructed by
	// `(typeof Function)`.
	return jutil.CheckNotNull(f.GetRegistry().GetNativeFunctionType(JSTypeNative_FUNCTION_FUNCTION_TYPE))
}

// Java: FunctionType.isInstanceType()
func (f *FunctionTypeBase) IsInstanceType() bool {
	// Only `Function` is both a function type and the intance type of a
	// nominal constructor.
	return JSType(f) == f.GetRegistry().GetNativeType(JSTypeNative_FUNCTION_TYPE)
}

// Java: FunctionType.isConstructor()
func (f *FunctionTypeBase) IsConstructor() bool {
	return f.call.fnKind == FunctionType_Kind_CONSTRUCTOR
}

// Java: FunctionType.isInterface()
func (f *FunctionTypeBase) IsInterface() bool {
	return f.call.fnKind == FunctionType_Kind_INTERFACE
}

// Java: FunctionType.isOrdinaryFunction()
func (f *FunctionTypeBase) IsOrdinaryFunction() bool {
	return f.call.fnKind == FunctionType_Kind_ORDINARY
}

// Java: FunctionType.getKind()
func (f *FunctionTypeBase) GetKind() FunctionType_Kind {
	return f.call.fnKind
}

// MakesStructs: when a class B inherits from A and A is annotated as a
// struct, then B automatically gets the annotation, if B's constructor is
// not explicitly annotated.
//
// Java: FunctionType.makesStructs()
func (f *FunctionTypeBase) MakesStructs() bool {
	if !f.HasInstanceType() || f.call.fnPropAccess == 0 {
		return false
	}
	switch FunctionType_PropAccess(f.call.fnPropAccess - 1) {
	case FunctionType_PropAccess_STRUCT:
		return true
	case FunctionType_PropAccess_DICT:
		return false
	case FunctionType_PropAccess_ANY_EXPLICIT:
		// For anything EXPLICITLY marked as @unresticted do not look to the
		// super type.
		return false
	case FunctionType_PropAccess_ANY:
		superc := f.GetSuperClassConstructor()
		if superc != nil && superc.MakesStructs() {
			f.SetStruct()
			return true
		}
		return false
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// MakesDicts: when a class B inherits from A and A is annotated as a dict,
// then B automatically gets the annotation, if B's constructor is not
// explicitly annotated.
//
// Java: FunctionType.makesDicts()
func (f *FunctionTypeBase) MakesDicts() bool {
	if !f.IsConstructor() || f.call.fnPropAccess == 0 {
		return false
	}
	switch FunctionType_PropAccess(f.call.fnPropAccess - 1) {
	case FunctionType_PropAccess_DICT:
		return true
	case FunctionType_PropAccess_STRUCT:
		return false
	case FunctionType_PropAccess_ANY_EXPLICIT:
		// For anything EXPLICITLY marked as @unresticted do not look to the
		// super type.
		return false
	case FunctionType_PropAccess_ANY:
		superc := f.GetSuperClassConstructor()
		if superc != nil && superc.MakesDicts() {
			f.SetDict()
			return true
		}
		return false
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: FunctionType.setStruct()
func (f *FunctionTypeBase) SetStruct() {
	f.call.fnPropAccess = uint8(FunctionType_PropAccess_STRUCT) + 1
}

// Java: FunctionType.setDict()
func (f *FunctionTypeBase) SetDict() {
	f.call.fnPropAccess = uint8(FunctionType_PropAccess_DICT) + 1
}

// Java: FunctionType.setExplicitUnrestricted()
func (f *FunctionTypeBase) SetExplicitUnrestricted() {
	f.call.fnPropAccess = uint8(FunctionType_PropAccess_ANY_EXPLICIT) + 1
}

// Java: JSType.isFunctionType()
func (f *FunctionTypeBase) IsFunctionType() bool {
	return !f.isNoObjectOrNo()
}

// Java: FunctionType.toMaybeFunctionType() / NoObjectType.toMaybeFunctionType()
func (f *FunctionTypeBase) ToMaybeFunctionType() FunctionType {
	if f.isNoObjectOrNo() {
		return nil
	}
	return f
}

// Java: NoObjectType.isNoObjectType() / NoType.isNoObjectType()
func (f *FunctionTypeBase) IsNoObjectType() bool {
	return f.objFlags&objKindMask == objKindNoObject
}

// Java: NoType.isNoType()
func (f *FunctionTypeBase) IsNoType() bool {
	return f.isNo()
}

// Java: NoType.isNullable()
func (f *FunctionTypeBase) IsNullable() bool {
	return f.isNo()
}

// Java: NoType.isVoidable()
func (f *FunctionTypeBase) IsVoidable() bool {
	return f.isNo()
}

// Java: ObjectType.getPossibleToBooleanOutcomes() / NoType.getPossibleToBooleanOutcomes()
func (f *FunctionTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	if f.isNo() {
		return BooleanLiteralSet_BOTH
	}
	return BooleanLiteralSet_TRUE
}

// Java: PrototypeObjectType.getImplicitPrototype() / NoObjectType.getImplicitPrototype()
func (f *FunctionTypeBase) GetImplicitPrototype() ObjectType {
	if f.isNoObjectOrNo() {
		return nil
	}
	return f.PrototypeObjectTypeBase.GetImplicitPrototype()
}

// Java: ObjectType.hasReferenceName()
func (f *FunctionTypeBase) HasReferenceName() bool {
	if f.isNoObjectOrNo() {
		return false
	}
	return f.PrototypeObjectTypeBase.HasReferenceName()
}

// Java: PrototypeObjectType.getReferenceName() / NoObjectType.getReferenceName()
func (f *FunctionTypeBase) GetReferenceName() *string {
	if f.isNoObjectOrNo() {
		return nil
	}
	return f.PrototypeObjectTypeBase.GetReferenceName()
}

// Java: PrototypeObjectType.matchesNumberContext() / NoObjectType.matchesNumberContext()
func (f *FunctionTypeBase) MatchesNumberContext() bool {
	if f.isNoObjectOrNo() {
		return true
	}
	return f.PrototypeObjectTypeBase.MatchesNumberContext()
}

// Java: PrototypeObjectType.matchesStringContext() / NoObjectType.matchesStringContext()
func (f *FunctionTypeBase) MatchesStringContext() bool {
	if f.isNoObjectOrNo() {
		return true
	}
	return f.PrototypeObjectTypeBase.MatchesStringContext()
}

// Java: PrototypeObjectType.matchesSymbolContext() / NoObjectType.matchesSymbolContext()
func (f *FunctionTypeBase) MatchesSymbolContext() bool {
	if f.isNoObjectOrNo() {
		return true
	}
	return f.PrototypeObjectTypeBase.MatchesSymbolContext()
}

// Java: PrototypeObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo) / NoObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (f *FunctionTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	if f.isNoObjectOrNo() {
		// Do nothing, specific properties do not have JSDocInfo.
		return
	}
	f.PrototypeObjectTypeBase.SetPropertyJSDocInfoByKey(propertyName, info)
}

// Java: FunctionType.canBeCalled()
func (f *FunctionTypeBase) CanBeCalled() bool {
	return true
}

// Java: FunctionType.getParameters()
func (f *FunctionTypeBase) GetParameters() []*FunctionType_Parameter {
	return f.GetInternalArrowType().GetParameterList()
}

// GetMinArity gets the minimum number of arguments that this function
// requires.
//
// Java: FunctionType.getMinArity()
func (f *FunctionTypeBase) GetMinArity() int {
	// NOTE(nicksantos): There are some native functions that have optional
	// parameters before required parameters. This algorithm finds the
	// position of the last required parameter.
	min := 0
	for i, parameter := range f.GetParameters() {
		if !parameter.IsOptional() && !parameter.IsVariadic() {
			min = i + 1
		}
	}
	return min
}

// GetMaxArity gets the maximum number of arguments that this function
// requires, or Integer.MAX_VALUE if this is a variable argument function.
//
// Java: FunctionType.getMaxArity()
func (f *FunctionTypeBase) GetMaxArity() int {
	params := f.GetParameters()
	if len(params) == 0 {
		return 0
	}
	lastParam := params[len(params)-1]
	if !lastParam.IsVariadic() {
		return len(params)
	}

	return math.MaxInt32
}

// Java: FunctionType.getReturnType()
func (f *FunctionTypeBase) GetReturnType() JSType {
	return f.call.GetReturnType()
}

// Java: FunctionType.isReturnTypeInferred()
func (f *FunctionTypeBase) IsReturnTypeInferred() bool {
	return f.call.returnTypeInferred
}

// GetInternalArrowType gets the internal arrow type. For use by subclasses
// only.
//
// Java: FunctionType.getInternalArrowType()
func (f *FunctionTypeBase) GetInternalArrowType() *ArrowType {
	return &f.call
}

// Java: FunctionType.getSlot(Property.Key)
func (f *FunctionTypeBase) GetSlotByKey(name Property_Key) *Property {
	isProto := false
	if !rhinoStringPoolWasCleared() {
		isProto = name.symbol == nil && name.stringID == propKeyIDPrototype
	} else {
		isProto = name.Matches("prototype")
	}
	if isProto {
		// Lazy initialization of the prototype field.
		f.GetPrototype()
		return f.prototypeSlot()
	}
	return f.PrototypeObjectTypeBase.GetSlotByKey(name)
}

// GetOwnPropertyNames includes the prototype iff someone has created it. We
// do not want to expose the prototype for ordinary functions.
//
// Java: FunctionType.getOwnPropertyNames()
func (f *FunctionTypeBase) GetOwnPropertyNames() []string {
	if f.prototypeSlot() == nil {
		return f.PrototypeObjectTypeBase.GetOwnPropertyNames()
	}
	pm := &f.properties
	pm.ensureSorted()
	entries := pm.entries()
	names := make([]string, 0, len(entries)+1)
	names = append(names, "prototype")
	for _, prop := range entries {
		name := rhinoStringFromID(prop.nameID)
		if name != "prototype" {
			names = append(names, name)
		}
	}
	return names
}

// ForEachOwnPropertyName calls fn for each own property name in sorted order (with "prototype" first when present).
func (f *FunctionTypeBase) ForEachOwnPropertyName(fn func(string)) {
	if f.prototypeSlot() == nil {
		f.PrototypeObjectTypeBase.ForEachOwnPropertyName(fn)
		return
	}
	fn("prototype")
	f.PrototypeObjectTypeBase.ForEachOwnPropertyName(func(name string) {
		if name != "prototype" {
			fn(name)
		}
	})
}

// Java: FunctionType.getPropertyType(Property.Key)
func (f *FunctionTypeBase) GetPropertyTypeByKey(name Property_Key) JSType {
	if !f.HasOwnPropertyByKey(name) {
		// Define the "call", "apply", and "bind" functions lazily.
		isCall := name.Matches("call")
		isBind := name.Matches("bind")
		if isCall || isBind {
			f.DefineDeclaredPropertyByKey(name, f.GetCallOrBindSignature(isCall), f.source)
		} else if name.Matches("apply") {
			// Define the "apply" function lazily.
			var builder FunctionParamBuilder
			builder.Init(f.GetRegistry(), 2)

			// ECMA-262 says that apply's second argument must be an Array or an
			// arguments object. We don't model the arguments object, so let's
			// just be forgiving for now.
			// TODO(nicksantos): Model the Arguments object.
			builder.AddOptionalParam(f.GetRegistry().CreateNullableType(f.GetTypeOfThis()))
			builder.AddOptionalParam(
				f.GetRegistry().CreateNullableType(f.GetRegistry().GetNativeType(JSTypeNative_OBJECT_TYPE)))

			f.DefineDeclaredPropertyByKey(
				name,
				FunctionType_Builder_(f.GetRegistry()).
					WithParametersOfList(builder.Build()).
					WithReturnType(f.GetReturnType()).
					WithTemplateKeysOfList(f.GetTemplateTypeMap().GetTemplateKeys()).
					Build(),
				f.source)
		}
	}

	return f.PrototypeObjectTypeBase.GetPropertyTypeByKey(name)
}

// GetBindReturnType gets the return value of calling "bind" on this function
// with the specified number of arguments.
//
// If -1 is passed, then we will return a result that accepts any
// parameters.
//
// Java: FunctionType.getBindReturnType(int)
func (f *FunctionTypeBase) GetBindReturnType(argsToBind int) FunctionType {
	builder := FunctionType_Builder_(f.GetRegistry()).
		WithReturnType(f.GetReturnType()).
		WithTemplateKeysOfList(f.GetTemplateTypeMap().GetTemplateKeys())
	if argsToBind < 0 {
		return builder.Build()
	}
	origParams := f.call.GetParameterList()
	params := append([]*FunctionType_Parameter{}, origParams...)
	for i := 1; i < argsToBind && len(params) != 0; i++ {
		if params[0].IsVariadic() {
			break
		}
		params = params[1:]
	}
	builder.WithParametersOfList(params)

	return builder.Build()
}

// getCallOrBindSignature: notice that "call" and "bind" have the same
// argument signature, except that all the arguments of "bind" (except the
// first) are optional.
//
// Java: FunctionType.GetCallOrBindSignature(boolean)
func (f *FunctionTypeBase) GetCallOrBindSignature(isCall bool) FunctionType {
	isBind := !isCall
	var returnType JSType
	if isCall {
		returnType = f.GetReturnType()
	} else {
		returnType = f.GetBindReturnType(-1)
	}
	builder := FunctionType_Builder_(f.GetRegistry()).
		WithReturnType(returnType).
		WithTemplateKeysOfList(f.GetTemplateTypeMap().GetTemplateKeys())

	origParams := f.GetInternalArrowType().GetParameterList()

	params := append([]*FunctionType_Parameter{}, origParams...)

	thisType := FunctionType_Parameter_Create(
		f.GetRegistry().CreateOptionalNullableType(f.GetTypeOfThis()),
		/* isOptional= */ false,
		/* isVariadic= */ false)
	params = append([]*FunctionType_Parameter{thisType}, params...)

	if isBind {
		// The arguments of bind() are unique in that they are all optional but
		// not undefinable.
		for i := 1; i < len(params); i++ {
			current := params[i]
			optionalCopy := FunctionType_Parameter_Create(current.GetJSType() /* isOptional= */, true, current.IsVariadic())
			params[i] = optionalCopy
		}
	} else if isCall {
		// The first argument of call() is optional iff all the arguments are
		// optional. It's sufficient to check the first argument.
		var firstArg *FunctionType_Parameter
		if len(params) > 1 {
			firstArg = params[1]
		}
		if firstArg == nil || firstArg.IsOptional() || firstArg.IsVariadic() {
			optionalThisType := FunctionType_Parameter_Create(thisType.GetJSType() /* isOptional= */, true /* isVariadic= */, false)
			params[0] = optionalThisType
		}
	}

	builder.WithParametersOfList(params)

	return builder.Build()
}

// Java: FunctionType.recursionUnsafeHashCode() / NoObjectType.recursionUnsafeHashCode()
func (f *FunctionTypeBase) RecursionUnsafeHashCode(hs *JSTypeStack) int32 {
	if f.isNoObjectOrNo() {
		return jutil.IdentityHashCode(f)
	}
	// Go-only: language-adapter — Enum.hashCode() is an identity hash in
	// Java; the ordinal is a stable stand-in.
	hc := int32(f.call.fnKind)
	switch f.call.fnKind {
	// constructors use identity semantics
	case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
		return 31*hc + jutil.IdentityHashCode(f)
	case FunctionType_Kind_ORDINARY:
		hc = 31*hc + f.typeOfThis.hashCodeIn(hs)
		hc = 31*hc + f.call.RecursionUnsafeHashCode(hs)
		hc = 31*hc + functionTypeClosurePrimitiveHashCode(f.GetClosurePrimitive())
		return hc
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// functionTypeClosurePrimitiveHashCode is Objects.hashCode on a nullable
// ClosurePrimitive.
func functionTypeClosurePrimitiveHashCode(p *ClosurePrimitive) int32 {
	if p == nil {
		return 0
	}
	return int32(*p) + 1
}

// Java: FunctionType.hasEqualCallType(FunctionType)
func (f *FunctionTypeBase) HasEqualCallType(that FunctionType) bool {
	return NewEqualityChecker().SetEqMethod(EqualityChecker_EqMethod_IDENTITY).
		CheckArrows(&f.call, &that.call)
}

// AppendTo: informally, a function is represented by function (params):
// returnType where the params is a comma separated list of types, the first
// one being a special this:T if the function expects a known type for this.
//
// Java: FunctionType.appendTo(TypeStringBuilder) / NoObjectType.appendTo(TypeStringBuilder) / NoType.appendTo(TypeStringBuilder)
func (f *FunctionTypeBase) AppendTo(sb *TypeStringBuilder) {
	if f.isNoObjectOrNo() {
		if sb.IsForAnnotations() {
			sb.Append("?")
		} else if f.isNo() {
			sb.Append("None")
		} else {
			sb.Append("NoObject")
		}
		return
	}
	if !f.IsPrettyPrint() || JSType(f) == f.GetRegistry().GetNativeType(JSTypeNative_FUNCTION_TYPE) {
		if sb.IsForAnnotations() {
			sb.Append("!Function")
		} else {
			sb.Append("Function")
		}
		return
	}

	if f.HasInstanceType() && !f.GetSource().IsNil() {
		// Render function types known to be type definitions as "(typeof
		// Foo)". This includes types defined like "/** @constructor */ function
		// Foo() { }" but not to those defined like "@param
		// {function(new:Foo)}". Only the former will have a source node.
		sb.Append("(typeof ")
		sb.AppendType(f.GetInstanceType())
		sb.Append(")")
		return
	}

	f.SetPrettyPrint(false)

	sb.Append("function(")
	paramNum := len(f.call.GetParameterList())
	_, isUnknown := jsTypeOf(f.typeOfThis).(*UnknownType)
	hasKnownTypeOfThis := !isUnknown
	if hasKnownTypeOfThis {
		if f.IsConstructor() {
			sb.Append("new:")
		} else {
			sb.Append("this:")
		}
		sb.AppendType(jsTypeOf(f.typeOfThis))
	}
	if paramNum > 0 {
		if hasKnownTypeOfThis {
			sb.Append(", ")
		}
		p := f.call.GetParameterList()[0]
		f.AppendArgString(sb, p)

		for i := 1; i < paramNum; i++ {
			p = f.call.GetParameterList()[i]
			sb.Append(", ")
			f.AppendArgString(sb, p)
		}
	}
	sb.Append("): ")
	sb.AppendNonNull(f.call.GetReturnType())

	f.SetPrettyPrint(true)
}

// Java: FunctionType.AppendArgString(TypeStringBuilder, Parameter)
func (f *FunctionTypeBase) AppendArgString(sb *TypeStringBuilder, p *FunctionType_Parameter) {
	if p.IsVariadic() {
		f.AppendVarArgsString(sb, p.GetJSType())
	} else if p.IsOptional() {
		f.AppendOptionalArgString(sb, p.GetJSType())
	} else {
		sb.AppendNonNull(p.GetJSType())
	}
}

// appendVarArgsString gets the string representation of a var args param.
//
// Java: FunctionType.AppendVarArgsString(TypeStringBuilder, JSType)
func (f *FunctionTypeBase) AppendVarArgsString(sb *TypeStringBuilder, paramType JSType) {
	sb.Append("...")
	sb.AppendNonNull(paramType)
}

// appendOptionalArgString gets the string representation of an optional
// param.
//
// Java: FunctionType.AppendOptionalArgString(TypeStringBuilder, JSType)
func (f *FunctionTypeBase) AppendOptionalArgString(sb *TypeStringBuilder, paramType JSType) {
	if paramType.IsUnionType() {
		// Remove the optionality from the var arg.
		paramType = paramType.
			ToMaybeUnionType().
			GetRestrictedUnion(f.GetRegistry().GetNativeType(JSTypeNative_VOID_TYPE))
	}
	sb.AppendNonNull(paramType)
	sb.Append("=")
}

// GetTypeOfThis gets the type of this in this function.
//
// Java: FunctionType.getTypeOfThis()
func (f *FunctionTypeBase) GetTypeOfThis() JSType {
	if f.typeOfThis.self().IsEmptyType() {
		return f.GetRegistry().GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)
	}
	return f.typeOfThis.self()
}

// GetSource gets the source node or null if this is an unknown function.
//
// Java: FunctionType.getSource()
func (f *FunctionTypeBase) GetSource() Node {
	return f.source
}

// Java: FunctionType.getGoogModuleId()
func (f *FunctionTypeBase) GetGoogModuleId() *string {
	if f.ext != nil {
		return rhinoStringPtrFromID(f.ext.googModuleID)
	}
	return nil
}

// SetSource sets the source node.
//
// Java: FunctionType.setSource(Node)
func (f *FunctionTypeBase) SetSource(source Node) {
	if slot := f.prototypeSlot(); slot != nil {
		// NOTE(bashir): On one hand when source is null we want to drop any
		// references to old nodes retained in prototypeSlot. On the other hand
		// we cannot simply drop prototypeSlot, so we retain all information
		// except the propertyNode for which we use an approximation! These
		// details mostly matter in hot-swap passes.
		if source.IsNil() || slot.GetNode().IsNil() {
			f.setPrototypeSlot(NewProperty(
				slot.GetName(),
				slot.GetType(),
				slot.IsTypeInferred(),
				source))
		}
	}
	f.source = source
}

// Java: FunctionType.hasAnyTemplateTypesInternal()
func (f *FunctionTypeBase) HasAnyTemplateTypesInternal(st *JSTypeStack) bool {
	return f.GetTemplateParamCount() > 0 ||
		f.typeOfThis.hasAnyTemplateTypesIn(st) ||
		f.call.HasAnyTemplateTypesInternal(st)
}

// Java: FunctionType.hasProperties()
func (f *FunctionTypeBase) HasProperties() bool {
	return f.GetPropertyMap().HasAnyOwnProperties()
}

// Java: FunctionType.acceptsArguments(List<? extends JSType>)
func (f *FunctionTypeBase) AcceptsArguments(argumentTypes []JSType) bool {
	// NOTE(aravindpg): This code is essentially lifted from
	// TypeCheck::visitParameterList, but what small differences there are
	// make it very painful to refactor out the shared code.
	parameters := f.GetParameters()
	ai, pi := 0, 0
	var parameter *FunctionType_Parameter
	var argument JSType
	for ai < len(argumentTypes) &&
		(pi < len(parameters) || (parameter != nil && parameter.IsVariadic())) {
		// If there are no parameters left in the list, then the while loop
		// above implies that this must be a var_args function.
		if pi < len(parameters) {
			parameter = parameters[pi]
			pi++
		}
		argument = argumentTypes[ai]
		ai++

		if !argument.IsSubtypeOf(parameter.GetJSType()) {
			return false
		}
	}

	numArgs := len(argumentTypes)
	return f.GetMinArity() <= numArgs && numArgs <= f.GetMaxArity()
}

// ForgetParameterAndReturnTypes creates a new constructor with the
// parameters and return type stripped.
//
// Java: FunctionType.forgetParameterAndReturnTypes()
func (f *FunctionTypeBase) ForgetParameterAndReturnTypes() FunctionType {
	result := FunctionType_Builder_(f.GetRegistry()).
		WithName(f.GetReferenceName()).
		WithSourceNode(f.source).
		WithTypeOfThis(f.GetInstanceType()).
		WithKind(f.call.fnKind).
		WithCanonicalRepresentation(f).
		Build()
	result.SetPrototypeBasedOn(f.GetInstanceType())
	return result
}

// GetConstructorOnlyTemplateParameters returns a list of template types
// present on the constructor but not on the instance.
//
// Java: FunctionType.getConstructorOnlyTemplateParameters()
func (f *FunctionTypeBase) GetConstructorOnlyTemplateParameters() []*TemplateType {
	jutil.CheckStateMsg1(f.IsConstructor(), "%s", f)

	// Structural constructor types (e.g `function(new:Foo<T>)`) cannot have
	// template params of their own.
	if f.source.IsNil() {
		return []*TemplateType{}
	}

	// Within the `TemplateTypeMap` of a ctor type, the ctor only keys always
	// appear after the instance type keys.
	m := f.GetTemplateTypeMap()
	ctorOnlyKeyCount := f.GetTemplateParamCount() - f.GetInstanceType().GetTemplateParamCount()
	return jsTypeSubList(m.GetTemplateKeys(), m.Size()-ctorOnlyKeyCount, m.Size())
}

// GetClosurePrimitive returns the @closurePrimitive identifier associated
// with this function
//
// Java: FunctionType.getClosurePrimitive()
func (f *FunctionTypeBase) GetClosurePrimitive() *ClosurePrimitive {
	if f.ext != nil && f.ext.closurePrimitivePlus1 != 0 {
		return &closurePrimitive_values[f.ext.closurePrimitivePlus1-1]
	}
	return nil
}

// Java: FunctionType.getCanonicalRepresentation()
func (f *FunctionTypeBase) GetCanonicalRepresentation() FunctionType {
	if f.ext != nil {
		return f.ext.canonicalRepresentation
	}
	return nil
}

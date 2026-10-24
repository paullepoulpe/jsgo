// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20261015

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
type FunctionType interface {
	PrototypeObjectType
	JSType_WithSourceRef

	GetKind() FunctionType_Kind
	MakesStructs() bool
	MakesDicts() bool
	SetStruct()
	SetDict()
	SetExplicitUnrestricted()
	HasImplementedInterfaces() bool
	GetParameters() []*FunctionType_Parameter
	GetMinArity() int
	GetMaxArity() int
	GetReturnType() JSType
	IsReturnTypeInferred() bool
	GetInternalArrowType() *ArrowType
	GetPrototypeProperty() ObjectType
	GetPrototype() ObjectType
	SetPrototypeBasedOn(baseType ObjectType)
	SetPrototype(prototype ObjectType, propertyNode Node) bool
	GetAllImplementedInterfaces() []ObjectType
	GetAncestorInterfaces() []ObjectType
	GetImplementedInterfaces() []ObjectType
	GetOwnImplementedInterfaces() []ObjectType
	SetImplementedInterfaces(implementedInterfaces []ObjectType)
	GetExtendedInterfaces() []ObjectType
	GetExtendedInterfacesCount() int
	SetExtendedInterfaces(extendedInterfaces []ObjectType)
	GetBindReturnType(argsToBind int) FunctionType
	SupAndInfHelper(that FunctionType, leastSuper bool) FunctionType
	HasEqualCallType(that FunctionType) bool
	GetInstanceType() ObjectType
	HasInstanceType() bool
	SetSource(source Node)
	HasProperties() bool
	SetImplicitMatch(flag bool)
	IsAbstract() bool
	CheckExtendsLoop() []FunctionType
	AcceptsArguments(argumentTypes []JSType) bool
	ForgetParameterAndReturnTypes() FunctionType
	GetConstructorOnlyTemplateParameters() []*TemplateType
	IsAmbiguousConstructor() bool
	GetClosurePrimitive() *ClosurePrimitive
	ToBuilder() *FunctionType_Builder
	GetCanonicalRepresentation() FunctionType

	// GetFunctionTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetFunctionTypeBase() *FunctionTypeBase
}

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

	// The primitive id associated with this FunctionType, or null if none.
	closurePrimitive *ClosurePrimitive

	// The ID of the goog.module in which this function was declared.
	googModuleId *string
}

// FunctionTypeBase holds FunctionType's fields and methods; plain function
// types are *FunctionTypeBase values.
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

func (f *FunctionTypeBase) selfFunc() FunctionType {
	if f.objFlags&objKindMask == objKindFunction {
		return f
	}
	return f.self().(FunctionType)
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

// FunctionType_Kind_values is Kind.values().
func FunctionType_Kind_values() []FunctionType_Kind {
	return []FunctionType_Kind{
		FunctionType_Kind_ORDINARY, FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE, FunctionType_Kind_NONE,
	}
}

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
	f.InitPrototypeObjectTypeBase(self, builder.GetBuilderFields())
	self.SetPrettyPrint(true)

	source := builder.sourceNode
	jutil.CheckArgument(source.IsNil() || source.IsFunction() || source.IsClass())
	f.source = source
	if builder.googModuleId != nil {
		f.ensureExt().googModuleId = builder.googModuleId
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
			typeOfThisBuilder.className = self.GetReferenceName()
			typeOfThisBuilder.nativeType = self.IsNativeObjectType()
			typeOfThisBuilder.anonymousType = self.IsAnonymous()
			typeOfThisBuilder.templateTypeMap = self.GetTemplateTypeMap()
			typeOfThisBuilder.templateParamCount = self.GetTemplateParamCount()
			typeOfThisBuilder.constructor = self

			ctorKeys := builder.constructorOnlyKeys
			if ctorKeys != nil && !ctorKeys.IsEmpty() {
				typeOfThisBuilder.templateTypeMap = f.templateTypeMap.CopyWithoutKeys(ctorKeys)
				typeOfThisBuilder.templateParamCount = f.selfFunc().GetTemplateParamCount() - ctorKeys.Size()
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

	if builder.primitiveId != nil {
		f.ensureExt().closurePrimitive = builder.primitiveId
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
		f.selfFunc().SetPrototypeBasedOn(builder.setPrototypeBasedOn)
	}

	f.GetRegistry().GetResolver().ResolveIfClosed(self, FunctionType_TYPE_CLASS)
}

// GetFunctionTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (f *FunctionTypeBase) GetFunctionTypeBase() *FunctionTypeBase { return f }

// Java: FunctionType.getTypeClass()
func (f *FunctionTypeBase) GetTypeClass() JSTypeClass {
	return FunctionType_TYPE_CLASS
}

// Java: FunctionType.getConstructor()
func (f *FunctionTypeBase) GetConstructor() FunctionType {
	// Every function type, including `Function`, is constructed by
	// `(typeof Function)`.
	return jutil.CheckNotNull(f.GetRegistry().GetNativeFunctionType(JSTypeNative_FUNCTION_FUNCTION_TYPE))
}

// Java: FunctionType.isInstanceType()
func (f *FunctionTypeBase) IsInstanceType() bool {
	// Only `Function` is both a function type and the intance type of a
	// nominal constructor.
	return JSType(f.selfFunc()) == f.GetRegistry().GetNativeType(JSTypeNative_FUNCTION_TYPE)
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
	if !f.selfFunc().HasInstanceType() || f.call.fnPropAccess == 0 {
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
		superc := f.selfFunc().GetSuperClassConstructor()
		if superc != nil && superc.MakesStructs() {
			f.selfFunc().SetStruct()
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
	if !f.selfFunc().IsConstructor() || f.call.fnPropAccess == 0 {
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
		superc := f.selfFunc().GetSuperClassConstructor()
		if superc != nil && superc.MakesDicts() {
			f.selfFunc().SetDict()
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
	return true
}

// Java: FunctionType.toMaybeFunctionType()
func (f *FunctionTypeBase) ToMaybeFunctionType() FunctionType {
	return f.selfFunc()
}

// Java: FunctionType.canBeCalled()
func (f *FunctionTypeBase) CanBeCalled() bool {
	return true
}

// Java: FunctionType.hasImplementedInterfaces()
func (f *FunctionTypeBase) HasImplementedInterfaces() bool {
	if f.ext != nil && len(f.ext.implementedInterfaces) != 0 {
		return true
	}
	var superCtor FunctionType
	if f.selfFunc().IsConstructor() {
		superCtor = f.selfFunc().GetSuperClassConstructor()
	}
	if superCtor != nil {
		return superCtor.HasImplementedInterfaces()
	}
	return false
}

// Java: FunctionType.getParameters()
func (f *FunctionTypeBase) GetParameters() []*FunctionType_Parameter {
	return f.selfFunc().GetInternalArrowType().GetParameterList()
}

// GetMinArity gets the minimum number of arguments that this function
// requires.
//
// Java: FunctionType.getMinArity()
func (f *FunctionTypeBase) GetMinArity() int {
	// NOTE(nicksantos): There are some native functions that have optional
	// parameters before required parameters. This algorithm finds the
	// position of the last required parameter.
	i := 0
	min := 0
	for _, parameter := range f.selfFunc().GetParameters() {
		i++
		if !parameter.IsOptional() && !parameter.IsVariadic() {
			min = i
		}
	}
	return min
}

// GetMaxArity gets the maximum number of arguments that this function
// requires, or Integer.MAX_VALUE if this is a variable argument function.
//
// Java: FunctionType.getMaxArity()
func (f *FunctionTypeBase) GetMaxArity() int {
	params := f.selfFunc().GetParameters()
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
	if name.Matches("prototype") {
		// Lazy initialization of the prototype field.
		f.selfFunc().GetPrototype()
		return f.prototypeSlot()
	} else {
		return f.PrototypeObjectTypeBase.GetSlotByKey(name)
	}
}

// GetOwnPropertyNames includes the prototype iff someone has created it. We
// do not want to expose the prototype for ordinary functions.
//
// Java: FunctionType.getOwnPropertyNames()
func (f *FunctionTypeBase) GetOwnPropertyNames() []string {
	if f.prototypeSlot() == nil {
		return f.PrototypeObjectTypeBase.GetOwnPropertyNames()
	} else {
		names := jutil.NewLinkedHashSet[string]()
		names.Add("prototype")
		f.PrototypeObjectTypeBase.ForEachOwnPropertyName(func(name string) {
			names.Add(name)
		})
		return names.ToSlice()
	}
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

// Java: FunctionType.getPrototypeProperty()
func (f *FunctionTypeBase) GetPrototypeProperty() ObjectType {
	return f.selfFunc().GetPrototype()
}

// GetPrototype gets the prototype property of this function type. This is
// equivalent to (ObjectType) getPropertyType("prototype").
//
// Java: FunctionType.getPrototype()
func (f *FunctionTypeBase) GetPrototype() ObjectType {
	// lazy initialization of the prototype field
	if f.prototypeSlot() == nil {
		refName := f.selfFunc().GetReferenceName()
		if refName == nil {
			// Someone is trying to access the prototype of a structural function.
			// We don't want to give real properties to this prototype, because
			// then it would propagate to all structural functions.
			f.SetPrototypeNoCheck(f.GetRegistry().GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE), NilNode)
		} else {
			name := jutil.ToJavaString(refName) + ".prototype"
			fields := PrototypeObjectType_BuilderFields{
				registry:          f.GetRegistry(),
				className:         &name,
				implicitPrototype: f.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE),
				nativeType:        f.selfFunc().IsNativeObjectType(),
				templateTypeMap:   f.GetRegistry().GetEmptyTemplateTypeMap(),
			}
			f.selfFunc().SetPrototype(
				NewPrototypeObjectType(&fields),
				NilNode)
		}
	}
	return f.prototypeSlot().GetType().(ObjectType)
}

// SetPrototypeBasedOn sets the prototype, creating the prototype object from
// the given base type.
//
// Java: FunctionType.setPrototypeBasedOn(ObjectType)
func (f *FunctionTypeBase) SetPrototypeBasedOn(baseType ObjectType) {
	f.SetPrototypeBasedOnWithPropertyNode(baseType, NilNode)
}

// Java: FunctionType.setPrototypeBasedOn(ObjectType, Node)
func (f *FunctionTypeBase) SetPrototypeBasedOnWithPropertyNode(baseType ObjectType, propertyNode Node) {
	if f.selfFunc().IsNativeObjectType() {
		return
	}
	// First handle class-side inheritance for ES6 classes, before reassigning
	// baseType.
	if !f.source.IsNil() && f.source.IsClass() {
		superCtor := baseType.GetConstructor()
		if superCtor != nil {
			f.selfFunc().SetImplicitPrototype(superCtor)
		}
		f.selfFunc().MaybeLoosenTypecheckingDueToForwardReferencedSupertype(baseType)
	}
	// This is a bit weird. We need to successfully handle these two cases:
	// Foo.prototype = new Bar();
	// and
	// Foo.prototype = {baz: 3};
	// In the first case, we do not want new properties to get added to Bar.
	// In the second case, we do want new properties to get added to the type
	// of the anonymous object.
	//
	// We handle this by breaking it into two cases:
	//
	// In the first case, we create a new PrototypeObjectType and set its
	// implicitPrototype to the type being assigned. This ensures that Bar
	// will not get any properties of Foo.prototype, but properties later
	// assigned to Bar will get inherited properly.
	//
	// In the second case, we just use the anonymous object as the prototype.
	if baseType.HasReferenceName() || f.selfFunc().IsNativeObjectType() || baseType.IsFunctionPrototypeType() {
		if f.prototypeSlot() != nil && f.selfFunc().HasInstanceType() && baseType.Equals(f.selfFunc().GetInstanceType()) {
			// Bail out for cases like Foo.prototype = new Foo();
			return
		}
		name := jutil.ToJavaString(f.selfFunc().GetReferenceName()) + ".prototype"
		fields := PrototypeObjectType_BuilderFields{
			registry:          f.GetRegistry(),
			className:         &name,
			implicitPrototype: baseType,
			templateTypeMap:   f.GetRegistry().GetEmptyTemplateTypeMap(),
		}
		baseType = NewPrototypeObjectType(&fields)
	}
	f.selfFunc().SetPrototype(baseType, propertyNode)
}

// SetPrototype sets the prototype. prototype: the prototype. If this value
// is null it will silently be discarded.
//
// Java: FunctionType.setPrototype(ObjectType, Node)
func (f *FunctionTypeBase) SetPrototype(prototype ObjectType, propertyNode Node) bool {
	if prototype == nil {
		return false
	}
	// getInstanceType fails if the function is not a constructor
	if f.selfFunc().IsConstructor() && prototype == f.selfFunc().GetInstanceType() {
		return false
	}
	return f.SetPrototypeNoCheck(prototype, propertyNode)
}

// setPrototypeNoCheck sets the prototype without doing any sanity checks.
//
// Java: FunctionType.SetPrototypeNoCheck(ObjectType, Node)
func (f *FunctionTypeBase) SetPrototypeNoCheck(prototype ObjectType, propertyNode Node) bool {
	var oldPrototype ObjectType
	if slot := f.prototypeSlot(); slot != nil {
		oldPrototype = slot.GetType().(ObjectType)
	}
	replacedPrototype := oldPrototype != nil

	slotNode := propertyNode
	if propertyNode.IsNil() {
		slotNode = f.source
	}
	f.setPrototypeSlot(NewProperty("prototype", prototype, true, slotNode))
	prototype.SetOwnerFunction(f.selfFunc())

	if oldPrototype != nil {
		// Disassociating the old prototype makes this easier to debug-- we
		// don't have to worry about two prototypes running around.
		oldPrototype.SetOwnerFunction(nil)
	}

	if replacedPrototype {
		f.selfFunc().ClearCachedValues()
	}
	if f.GetRegistry() != nil {
		f.GetRegistry().leastSupertypeGen++
	}

	return true
}

// GetAllImplementedInterfaces returns all interfaces implemented by a class
// or its superclass and any superclasses for any of those interfaces. If
// this is called before all types are resolved, it may return an incomplete
// set.
//
// Java: FunctionType.getAllImplementedInterfaces()
func (f *FunctionTypeBase) GetAllImplementedInterfaces() []ObjectType {
	// Store them in a linked hash set, so that the compile job is
	// deterministic.
	interfaces := jutil.NewObjectLinkedHashSet[ObjectType]()

	for _, type_ := range f.selfFunc().GetImplementedInterfaces() {
		f.AddRelatedInterfaces(type_, interfaces)
	}
	return interfaces.ToSlice()
}

// Java: FunctionType.AddRelatedInterfaces(ObjectType, Set<ObjectType>)
func (f *FunctionTypeBase) AddRelatedInterfaces(instance ObjectType, set *jutil.ObjectLinkedHashSet[ObjectType]) {
	constructor := instance.GetConstructor()
	if constructor != nil {
		if !constructor.IsInterface() {
			return
		}

		if !set.Add(instance) {
			return
		}

		for _, interfaceType := range instance.GetCtorExtendedInterfaces() {
			f.AddRelatedInterfaces(interfaceType, set)
		}
	}
}

// Java: FunctionType.getAncestorInterfaces()
func (f *FunctionTypeBase) GetAncestorInterfaces() []ObjectType {
	result := jutil.NewObjectLinkedHashSet[ObjectType]()
	if f.selfFunc().IsConstructor() {
		for _, t := range f.selfFunc().GetImplementedInterfaces() {
			result.Add(t)
		}
	} else {
		for _, t := range f.selfFunc().GetExtendedInterfaces() {
			result.Add(t)
		}
	}
	return result.ToSlice()
}

// GetImplementedInterfaces returns interfaces implemented directly by a
// class or its superclass.
//
// Java: FunctionType.getImplementedInterfaces()
func (f *FunctionTypeBase) GetImplementedInterfaces() []ObjectType {
	var superCtor FunctionType
	if f.selfFunc().IsConstructor() {
		superCtor = f.selfFunc().GetSuperClassConstructor()
	}
	own := f.GetOwnImplementedInterfaces()
	if superCtor == nil {
		return own
	}
	builder := append([]ObjectType{}, own...)
	for superCtor != nil {
		builder = append(builder, superCtor.GetFunctionTypeBase().GetOwnImplementedInterfaces()...)
		superCtor = superCtor.GetSuperClassConstructor()
	}
	return builder
}

// GetOwnImplementedInterfaces returns interfaces directly implemented by the
// class.
//
// Java: FunctionType.getOwnImplementedInterfaces()
func (f *FunctionTypeBase) GetOwnImplementedInterfaces() []ObjectType {
	if f.ext != nil && f.ext.implementedInterfaces != nil {
		return f.ext.implementedInterfaces
	}
	return emptyObjectTypeSlice
}

// Java: FunctionType.setImplementedInterfaces(List<ObjectType>)
func (f *FunctionTypeBase) SetImplementedInterfaces(implementedInterfaces []ObjectType) {
	jutil.CheckState(f.selfFunc().IsConstructor())

	copied := functionTypeCopyOfObjectTypes(implementedInterfaces)
	if len(copied) > 0 || f.ext != nil {
		f.ensureExt().implementedInterfaces = copied
	}
	for _, type_ := range implementedInterfaces {
		f.typeOfThis.self().MergeSupertypeTemplateTypes(type_)
	}
}

// functionTypeCopyOfObjectTypes is ImmutableList.copyOf(list), which
// rejects null elements.
func functionTypeCopyOfObjectTypes(list []ObjectType) []ObjectType {
	if len(list) == 0 {
		return emptyObjectTypeSlice
	}
	out := make([]ObjectType, len(list))
	for i, e := range list {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
		out[i] = e
	}
	return out
}

// GetExtendedInterfaces returns interfaces directly extended by an
// interface.
//
// Java: FunctionType.getExtendedInterfaces()
func (f *FunctionTypeBase) GetExtendedInterfaces() []ObjectType {
	if f.ext != nil && f.ext.extendedInterfaces != nil {
		return f.ext.extendedInterfaces
	}
	return emptyObjectTypeSlice
}

// GetExtendedInterfacesCount returns the number of interfaces directly
// extended by an interface.
//
// Java: FunctionType.getExtendedInterfacesCount()
func (f *FunctionTypeBase) GetExtendedInterfacesCount() int {
	if f.ext != nil {
		return len(f.ext.extendedInterfaces)
	}
	return 0
}

// Java: FunctionType.setExtendedInterfaces(List<ObjectType>)
func (f *FunctionTypeBase) SetExtendedInterfaces(extendedInterfaces []ObjectType) {
	jutil.CheckState(f.selfFunc().IsInterface())

	copied := functionTypeCopyOfObjectTypes(extendedInterfaces)
	if len(copied) > 0 || f.ext != nil {
		f.ensureExt().extendedInterfaces = copied
	}
	for _, extendedInterface := range extendedInterfaces {
		f.typeOfThis.self().MergeSupertypeTemplateTypes(extendedInterface)
	}
}

// Java: FunctionType.getPropertyType(Property.Key)
func (f *FunctionTypeBase) GetPropertyTypeByKey(name Property_Key) JSType {
	if !f.selfFunc().HasOwnPropertyByKey(name) {
		// Define the "call", "apply", and "bind" functions lazily.
		isCall := name.Matches("call")
		isBind := name.Matches("bind")
		if isCall || isBind {
			f.selfFunc().DefineDeclaredPropertyByKey(name, f.GetCallOrBindSignature(isCall), f.source)
		} else if name.Matches("apply") {
			// Define the "apply" function lazily.
			builder := NewFunctionParamBuilder(f.GetRegistry())

			// ECMA-262 says that apply's second argument must be an Array or an
			// arguments object. We don't model the arguments object, so let's
			// just be forgiving for now.
			// TODO(nicksantos): Model the Arguments object.
			builder.AddOptionalParams(
				f.GetRegistry().CreateNullableType(f.selfFunc().GetTypeOfThis()),
				f.GetRegistry().CreateNullableType(f.GetRegistry().GetNativeType(JSTypeNative_OBJECT_TYPE)))

			f.selfFunc().DefineDeclaredPropertyByKey(
				name,
				FunctionType_Builder_(f.GetRegistry()).
					WithParametersOfList(builder.Build()).
					WithReturnType(f.selfFunc().GetReturnType()).
					WithTemplateKeysOfList(f.selfFunc().GetTemplateTypeMap().GetTemplateKeys()).
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
		WithReturnType(f.selfFunc().GetReturnType()).
		WithTemplateKeysOfList(f.selfFunc().GetTemplateTypeMap().GetTemplateKeys())
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
		returnType = f.selfFunc().GetReturnType()
	} else {
		returnType = f.selfFunc().GetBindReturnType(-1)
	}
	builder := FunctionType_Builder_(f.GetRegistry()).
		WithReturnType(returnType).
		WithTemplateKeysOfList(f.selfFunc().GetTemplateTypeMap().GetTemplateKeys())

	origParams := f.selfFunc().GetInternalArrowType().GetParameterList()

	params := append([]*FunctionType_Parameter{}, origParams...)

	thisType := FunctionType_Parameter_Create(
		f.GetRegistry().CreateOptionalNullableType(f.selfFunc().GetTypeOfThis()),
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

// Java: FunctionType.defineProperty(Property.Key, JSType, boolean, Node)
func (f *FunctionTypeBase) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	if name.Matches("prototype") {
		objType := type_.ToObjectType()
		if objType != nil {
			if slot := f.prototypeSlot(); slot != nil && objType.Equals(slot.GetType()) {
				return true
			}
			f.SetPrototypeBasedOnWithPropertyNode(objType, propertyNode)
			return true
		} else {
			return false
		}
	}
	return f.PrototypeObjectTypeBase.DefinePropertyByKey(name, type_, inferred, propertyNode)
}

// SupAndInfHelper computes the supremum or infimum of two functions. Because
// sup() and inf() share a lot of logic for functions, we use a single
// helper.
//
// leastSuper: if true, compute the supremum of this with that. Otherwise,
// compute the infimum. Returns the least supertype or greatest subtype.
//
// Java: FunctionType.supAndInfHelper(FunctionType, boolean)
func (f *FunctionTypeBase) SupAndInfHelper(that FunctionType, leastSuper bool) FunctionType {
	// NOTE(nicksantos): When we remove the unknown type, the function types
	// form a lattice with the universal constructor at the top of the
	// lattice, and the LEAST_FUNCTION_TYPE type at the bottom of the lattice.
	//
	// When we introduce the unknown type, it's much more difficult to make
	// heads or tails of the partial ordering of types, because there's no
	// clear hierarchy between the different components (parameter types and
	// return types) in the ArrowType.
	//
	// Rather than make the situation more complicated by introducing new
	// types (like unions of functions), we just fallback on the simpler
	// approach of getting things right at the top and the bottom of the
	// lattice.
	//
	// If there are unknown parameters or return types making things
	// ambiguous, then sup(A, B) is always the top function type, and inf(A,
	// B) is always the bottom function type.
	jutil.CheckNotNull(that)

	if f.selfFunc().Equals(that) {
		return f.selfFunc()
	}

	// If these are ordinary functions, then merge them. Don't do this if any
	// of the params/return values are unknown, because then there will be
	// cycles in their local lattice and they will merge in weird ways.
	if f.selfFunc().IsOrdinaryFunction() &&
		that.IsOrdinaryFunction() &&
		!f.call.HasUnknownParamsOrReturn() &&
		!that.GetFunctionTypeBase().call.HasUnknownParamsOrReturn() {

		// Check for the degenerate case, but double check that there's not a
		// cycle.
		isSubtypeOfThat := f.selfFunc().IsSubtype(that)
		isSubtypeOfThis := that.IsSubtype(f.selfFunc())
		if isSubtypeOfThat && !isSubtypeOfThis {
			if leastSuper {
				return that
			}
			return f.selfFunc()
		} else if isSubtypeOfThis && !isSubtypeOfThat {
			if leastSuper {
				return f.selfFunc()
			}
			return that
		}

		// Merge the two functions component-wise.
		merged := f.TryMergeFunctionPiecewise(that, leastSuper)
		if merged != nil {
			return merged
		}
	}

	// The function instance type is a special case that lives above the rest
	// of the lattice.
	functionInstance := f.GetRegistry().GetNativeType(JSTypeNative_FUNCTION_TYPE)
	if functionInstance.Equals(that) {
		if leastSuper {
			return that
		}
		return f.selfFunc()
	} else if functionInstance.Equals(f.selfFunc()) {
		if leastSuper {
			return f.selfFunc()
		}
		return that
	}

	// In theory, we should be using the GREATEST_FUNCTION_TYPE as the
	// greatest function. In practice, we don't because it's way too broad.
	// The greatest function takes var_args None parameters, which means that
	// all parameters register a type warning.
	//
	// Instead, we use the U2U ctor type, which has unknown type args.
	greatestFn := f.GetRegistry().GetNativeFunctionType(JSTypeNative_FUNCTION_TYPE)
	leastFn := f.GetRegistry().GetNativeFunctionType(JSTypeNative_LEAST_FUNCTION_TYPE)
	if leastSuper {
		return greatestFn
	}
	return leastFn
}

// tryMergeFunctionPiecewise tries to get the sup/inf of two functions by
// looking at the piecewise components.
//
// Java: FunctionType.TryMergeFunctionPiecewise(FunctionType, boolean)
func (f *FunctionTypeBase) TryMergeFunctionPiecewise(other FunctionType, leastSuper bool) FunctionType {
	otherBase := other.GetFunctionTypeBase()
	var newParamsNode []*FunctionType_Parameter
	if NewEqualityChecker().
		SetEqMethod(EqualityChecker_EqMethod_IDENTITY).
		CheckParameters(&f.call, &otherBase.call) {
		newParamsNode = f.call.GetParameterList()
	} else {
		// If the parameters are not equal, don't try to merge them. Someday, we
		// should try to merge the individual params.
		return nil
	}

	var newReturnType JSType
	if leastSuper {
		newReturnType = f.call.GetReturnType().GetLeastSupertype(otherBase.call.GetReturnType())
	} else {
		newReturnType = f.call.GetReturnType().GetGreatestSubtype(otherBase.call.GetReturnType())
	}

	var newTypeOfThis JSType
	if jutil.ObjectsEquals(jsTypeOf(f.typeOfThis), jsTypeOf(otherBase.typeOfThis)) {
		newTypeOfThis = jsTypeOf(f.typeOfThis)
	} else {
		var maybeNewTypeOfThis JSType
		if leastSuper {
			maybeNewTypeOfThis = f.typeOfThis.self().GetLeastSupertype(jsTypeOf(otherBase.typeOfThis))
		} else {
			maybeNewTypeOfThis = f.typeOfThis.self().GetGreatestSubtype(jsTypeOf(otherBase.typeOfThis))
		}
		newTypeOfThis = maybeNewTypeOfThis
	}

	newReturnTypeInferred := f.call.returnTypeInferred || otherBase.call.returnTypeInferred

	return FunctionType_Builder_(f.GetRegistry()).
		WithParametersOfList(newParamsNode).
		WithReturnTypeWithInferred(newReturnType, newReturnTypeInferred).
		WithTypeOfThis(newTypeOfThis).
		Build()
}

// GetSuperClassConstructor: given a constructor or an interface type, get
// its superclass constructor or null if none exists.
//
// Java: FunctionType.getSuperClassConstructor()
func (f *FunctionTypeBase) GetSuperClassConstructor() FunctionType {
	jutil.CheckArgument(f.selfFunc().IsConstructor() || f.selfFunc().IsInterface())
	maybeSuperInstanceType := f.selfFunc().GetPrototype().GetImplicitPrototype()
	if maybeSuperInstanceType == nil {
		return nil
	}
	return maybeSuperInstanceType.GetConstructor()
}

// Java: FunctionType.recursionUnsafeHashCode()
func (f *FunctionTypeBase) RecursionUnsafeHashCode() int32 {
	// Go-only: language-adapter — Enum.hashCode() is an identity hash in
	// Java; the ordinal is a stable stand-in.
	hc := int32(f.call.fnKind)
	switch f.call.fnKind {
	// constructors use identity semantics
	case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
		return 31*hc + jutil.IdentityHashCode(f.selfFunc())
	case FunctionType_Kind_ORDINARY:
		hc = 31*hc + f.typeOfThis.self().HashCode()
		hc = 31*hc + f.call.HashCode()
		hc = 31*hc + functionTypeClosurePrimitiveHashCode(f.selfFunc().GetClosurePrimitive())
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
		CheckArrows(&f.call, &that.GetFunctionTypeBase().call)
}

// AppendTo: informally, a function is represented by function (params):
// returnType where the params is a comma separated list of types, the first
// one being a special this:T if the function expects a known type for this.
//
// Java: FunctionType.appendTo(TypeStringBuilder)
func (f *FunctionTypeBase) AppendTo(sb *TypeStringBuilder) {
	if !f.selfFunc().IsPrettyPrint() || JSType(f.selfFunc()) == f.GetRegistry().GetNativeType(JSTypeNative_FUNCTION_TYPE) {
		if sb.IsForAnnotations() {
			sb.Append("!Function")
		} else {
			sb.Append("Function")
		}
		return
	}

	if f.selfFunc().HasInstanceType() && !f.selfFunc().GetSource().IsNil() {
		// Render function types known to be type definitions as "(typeof
		// Foo)". This includes types defined like "/** @constructor */ function
		// Foo() { }" but not to those defined like "@param
		// {function(new:Foo)}". Only the former will have a source node.
		sb.Append("(typeof ").AppendType(f.selfFunc().GetInstanceType()).Append(")")
		return
	}

	f.selfFunc().SetPrettyPrint(false)

	sb.Append("function(")
	paramNum := len(f.call.GetParameterList())
	_, isUnknown := jsTypeOf(f.typeOfThis).(*UnknownType)
	hasKnownTypeOfThis := !isUnknown
	if hasKnownTypeOfThis {
		if f.selfFunc().IsConstructor() {
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

	f.selfFunc().SetPrettyPrint(true)
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
	sb.Append("...").AppendNonNull(paramType)
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
	sb.AppendNonNull(paramType).Append("=")
}

// Java: FunctionType.visit(Visitor<T>)
func (f *FunctionTypeBase) Visit(visitor Visitor) any {
	return visitor.CaseFunctionType(f.selfFunc())
}

// Java: FunctionType.visit(RelationshipVisitor<T>, JSType)
func (f *FunctionTypeBase) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseFunctionType(f.selfFunc(), that)
}

// GetInstanceType gets the type of instance of this function. May return
// null if the `this` type can not be converted to "ObjectType" (see
// JSType#toObjectType).
//
// Throws IllegalStateException if this function is not a constructor (see
// #isConstructor()).
//
// Java: FunctionType.getInstanceType()
func (f *FunctionTypeBase) GetInstanceType() ObjectType {
	jutil.CheckState(f.selfFunc().HasInstanceType())

	return f.typeOfThis.self().ToObjectType()
}

// HasInstanceType returns whether this function type has an instance type.
//
// Java: FunctionType.hasInstanceType()
func (f *FunctionTypeBase) HasInstanceType() bool {
	return f.selfFunc().IsConstructor() || f.selfFunc().IsInterface()
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
		return f.ext.googModuleId
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

// ClearCachedValues: NOTE(sdh): One might assume that immediately after
// calling this, hasCachedValues() should always return false. This is not
// the case, since hasCachedValues() will return true if prototypeSlot is
// non-null, and this override does nothing to change that state.
//
// Java: FunctionType.clearCachedValues()
func (f *FunctionTypeBase) ClearCachedValues() {
	f.PrototypeObjectTypeBase.ClearCachedValues()

	if !f.selfFunc().IsNativeObjectType() {
		if f.selfFunc().HasInstanceType() {
			f.selfFunc().GetInstanceType().ClearCachedValues()
		}

		if slot := f.prototypeSlot(); slot != nil {
			slot.GetType().(ObjectType).ClearCachedValues()
		}
	}
}

// Java: FunctionType.hasCachedValues()
func (f *FunctionTypeBase) HasCachedValues() bool {
	return f.prototypeSlot() != nil || f.PrototypeObjectTypeBase.HasCachedValues()
}

// Java: FunctionType.resolveInternal(ErrorReporter)
func (f *FunctionTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	if f.call.parameterList != nil && !f.call.resolved {
		f.call.resolveInternal(reporter)
	}
	if slot := f.prototypeSlot(); slot != nil {
		slot.SetType(JSType_SafeResolve(slot.GetType(), reporter))
	}

	// Warning about typeOfThis if it doesn't resolve to an ObjectType is
	// handled further upstream.
	//
	// TODO(nicksantos): Handle this correctly if we have a UnionType.
	maybeTypeOfThis := JSType_SafeResolve(jsTypeOf(f.typeOfThis), reporter)
	if maybeTypeOfThis != nil {
		if maybeTypeOfThis.IsNullType() || maybeTypeOfThis.IsVoidType() {
			f.typeOfThis = maybeTypeOfThis.GetJSTypeBase()
		} else {
			maybeTypeOfThis = jsTypeOrNil(ObjectType_Cast(maybeTypeOfThis.RestrictByNotNullOrUndefined()))
			if maybeTypeOfThis != nil {
				f.typeOfThis = maybeTypeOfThis.GetJSTypeBase()
			}
		}
	}

	if f.ext != nil {
		if len(f.ext.implementedInterfaces) > 0 {
			resolvedImplemented := f.ResolveTypeListHelper(f.ext.implementedInterfaces, reporter)
			if resolvedImplemented != nil {
				f.ext.implementedInterfaces = resolvedImplemented
			}
		}
		if len(f.ext.extendedInterfaces) > 0 {
			resolvedExtended := f.ResolveTypeListHelper(f.ext.extendedInterfaces, reporter)
			if resolvedExtended != nil {
				f.ext.extendedInterfaces = resolvedExtended
			}
		}
	}

	return f.PrototypeObjectTypeBase.ResolveInternal(reporter)
}

// jsTypeOrNil converts an ObjectType to a JSType, keeping Java null.
func jsTypeOrNil(t ObjectType) JSType {
	if t == nil {
		return nil
	}
	return t
}

// resolveTypeListHelper resolves each item in the list, and returns a new
// list if any references changed. Otherwise, returns null.
//
// Java: FunctionType.ResolveTypeListHelper(ImmutableList<ObjectType>, ErrorReporter)
func (f *FunctionTypeBase) ResolveTypeListHelper(list []ObjectType, reporter ErrorReporter) []ObjectType {
	changed := false
	resolvedList := []ObjectType{}
	for _, type_ := range list {
		rt := type_.Resolve(reporter)
		if !rt.IsObjectType() {
			reporter.Warning(
				"not an object type: "+rt.String()+" (at "+f.selfFunc().String()+")",
				jsTypeDerefString(f.source.GetSourceFileName()),
				f.source.GetLineno(),
				f.source.GetCharno())
			continue
		}
		resolved := rt.ToObjectType()
		resolvedList = append(resolvedList, resolved)
		changed = changed || resolved != type_
	}
	if changed {
		return resolvedList
	}
	return nil
}

// jsTypeDerefString passes a nullable String to an API that takes a Go
// string; null becomes "".
//
// Go-only: language-adapter — ErrorReporter takes a Go string for the
// nullable source name.
func jsTypeDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Java: FunctionType.hasAnyTemplateTypesInternal()
func (f *FunctionTypeBase) HasAnyTemplateTypesInternal() bool {
	return f.selfFunc().GetTemplateParamCount() > 0 ||
		f.typeOfThis.self().HasAnyTemplateTypes() ||
		f.call.HasAnyTemplateTypes()
}

// Java: FunctionType.hasProperties()
func (f *FunctionTypeBase) HasProperties() bool {
	return f.GetPropertyMap().HasAnyOwnProperties()
}

// SetImplicitMatch sets the current interface type to support structural
// interface matching (abbr. SMI). flag indicates whether or not it should
// support SMI.
//
// Java: FunctionType.setImplicitMatch(boolean)
func (f *FunctionTypeBase) SetImplicitMatch(flag bool) {
	jutil.CheckState(f.selfFunc().IsInterface())
	f.call.fnIsStructuralInterface = flag
}

// Java: FunctionType.isStructuralInterface()
func (f *FunctionTypeBase) IsStructuralInterface() bool {
	return f.selfFunc().IsInterface() && f.call.fnIsStructuralInterface
}

// Java: FunctionType.isAbstract()
func (f *FunctionTypeBase) IsAbstract() bool {
	return f.call.fnIsAbstract
}

// GetPropertyTypeMap gets the map of properties to types covered in a
// function type. Returns a Map that maps the property's name to the
// property's type.
//
// Java: FunctionType.getPropertyTypeMap()
func (f *FunctionTypeBase) GetPropertyTypeMap() *jutil.LinkedHashMap[string, JSType] {
	propTypeMap := jutil.NewLinkedHashMap[string, JSType]()
	FunctionType_UpdatePropertyTypeMap(f.selfFunc(), propTypeMap, jutil.NewObjectLinkedHashSet[FunctionType]())
	return propTypeMap
}

// FunctionType_UpdatePropertyTypeMap: cache is added to prevent infinite
// recursion when retrieving the super type: see testInterfaceExtendsLoop in
// TypeCheckTest.java
//
// Java: FunctionType.updatePropertyTypeMap(FunctionType, Map<String, JSType>, LinkedHashSet<FunctionType>)
func FunctionType_UpdatePropertyTypeMap(
	type_ FunctionType, propTypeMap *jutil.LinkedHashMap[string, JSType], cache *jutil.ObjectLinkedHashSet[FunctionType]) {
	if type_ == nil {
		return
	}
	// retrieve all property types on the prototype of this class
	prototype := type_.GetPrototype()
	if prototype != nil {
		propNames := prototype.GetOwnPropertyNames()
		for _, name := range propNames {
			if !propTypeMap.ContainsKey(name) {
				propType := prototype.GetPropertyType(name)
				propTypeMap.Put(name, propType)
			}
		}
	}
	// retrieve all property types from its super class
	iterable := type_.GetExtendedInterfaces()
	if iterable != nil {
		for _, interfaceType := range iterable {
			superConstructor := interfaceType.GetConstructor()
			if superConstructor == nil || cache.Contains(superConstructor) {
				continue
			}
			cache.Add(superConstructor)
			FunctionType_UpdatePropertyTypeMap(superConstructor, propTypeMap, cache)
			cache.Remove(superConstructor)
		}
	}
}

// CheckExtendsLoop checks if there is a loop in the type extends chain.
// Returns an array of all functions in the loop chain if a loop exists,
// otherwise returns null.
//
// Java: FunctionType.checkExtendsLoop()
func (f *FunctionTypeBase) CheckExtendsLoop() []FunctionType {
	return f.CheckExtendsLoopWithCacheAndPath(jutil.NewObjectLinkedHashSet[FunctionType](), jutil.NewArrayList[FunctionType]())
}

// Java: FunctionType.checkExtendsLoop(Set<FunctionType>, List<FunctionType>)
func (f *FunctionTypeBase) CheckExtendsLoopWithCacheAndPath(
	cache *jutil.ObjectLinkedHashSet[FunctionType], path *jutil.ArrayList[FunctionType]) []FunctionType {
	iterable := f.selfFunc().GetExtendedInterfaces()
	if iterable != nil {
		for _, interfaceType := range iterable {
			superConstructor := interfaceType.GetConstructor()
			if superConstructor == nil {
				continue
			}
			if cache.Contains(superConstructor) {
				// after detecting a loop, prune and return the path, e.g.,: A -> B
				// -> C -> D -> C, will be pruned into: c -> D -> C
				path.Add(superConstructor)
				for path.Get(0) != superConstructor {
					path.RemoveAt(0)
				}
				return path.Slice()
			}
			cache.Add(superConstructor)
			path.Add(superConstructor)
			result := superConstructor.GetFunctionTypeBase().CheckExtendsLoopWithCacheAndPath(cache, path)
			if result != nil {
				return result
			}
			cache.Remove(superConstructor)
			path.RemoveAt(path.Size() - 1)
		}
	}
	return nil
}

// Java: FunctionType.acceptsArguments(List<? extends JSType>)
func (f *FunctionTypeBase) AcceptsArguments(argumentTypes []JSType) bool {
	// NOTE(aravindpg): This code is essentially lifted from
	// TypeCheck::visitParameterList, but what small differences there are
	// make it very painful to refactor out the shared code.
	parameters := f.selfFunc().GetParameters()
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
	return f.selfFunc().GetMinArity() <= numArgs && numArgs <= f.selfFunc().GetMaxArity()
}

// ForgetParameterAndReturnTypes creates a new constructor with the
// parameters and return type stripped.
//
// Java: FunctionType.forgetParameterAndReturnTypes()
func (f *FunctionTypeBase) ForgetParameterAndReturnTypes() FunctionType {
	result := FunctionType_Builder_(f.GetRegistry()).
		WithName(f.selfFunc().GetReferenceName()).
		WithSourceNode(f.source).
		WithTypeOfThis(f.selfFunc().GetInstanceType()).
		WithKind(f.call.fnKind).
		WithCanonicalRepresentation(f.selfFunc()).
		Build()
	result.SetPrototypeBasedOn(f.selfFunc().GetInstanceType())
	return result
}

// GetConstructorOnlyTemplateParameters returns a list of template types
// present on the constructor but not on the instance.
//
// Java: FunctionType.getConstructorOnlyTemplateParameters()
func (f *FunctionTypeBase) GetConstructorOnlyTemplateParameters() []*TemplateType {
	jutil.CheckStateMsg1(f.selfFunc().IsConstructor(), "%s", f.selfFunc())

	// Structural constructor types (e.g `function(new:Foo<T>)`) cannot have
	// template params of their own.
	if f.source.IsNil() {
		return []*TemplateType{}
	}

	// Within the `TemplateTypeMap` of a ctor type, the ctor only keys always
	// appear after the instance type keys.
	m := f.selfFunc().GetTemplateTypeMap()
	ctorOnlyKeyCount := f.selfFunc().GetTemplateParamCount() - f.selfFunc().GetInstanceType().GetTemplateParamCount()
	return jsTypeSubList(m.GetTemplateKeys(), m.Size()-ctorOnlyKeyCount, m.Size())
}

// IsAmbiguousConstructor returns true if the constructor does not come from
// a literal class or function in the AST, or if it extends such an
// ambiguous constructor
//
// Java: FunctionType.isAmbiguousConstructor()
func (f *FunctionTypeBase) IsAmbiguousConstructor() bool {
	if f.call.fnConstructorAmbiguity == FunctionType_ConstructorAmbiguity_UNKNOWN {
		f.call.fnConstructorAmbiguity = f.CalculateConstructorAmbiguity()
	}
	return f.call.fnConstructorAmbiguity == FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
}

// Java: FunctionType.CalculateConstructorAmbiguity()
func (f *FunctionTypeBase) CalculateConstructorAmbiguity() FunctionType_ConstructorAmbiguity {
	var constructorAmbiguity FunctionType_ConstructorAmbiguity
	if f.selfFunc().IsUnknownType() {
		constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	} else if f.selfFunc().IsNativeObjectType() {
		// native types other than unknown are never ambiguous
		constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
	} else {
		superConstructor := f.selfFunc().GetSuperClassConstructor()
		if superConstructor == nil {
			// TODO(bradfordcsmith): Why is superConstructor ever null here?
			constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
		} else if superConstructor.IsAmbiguousConstructor() {
			// Subclasses of ambiguous objects are also ambiguous
			constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
		} else if !f.source.IsNil() {
			// We can see the definition of the class, so we know all properties it
			// directly declares or references. The same is true for its superclass
			// (previous condition).
			constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
		} else if f.IsDelegateProxy() {
			// Type was created by the compiler as a proxy that inherits from the
			// real type that was in the code. Since we've made it this far, we
			// know the real type creates unambiguous objects. Therefore, the
			// proxy does, too.
			constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
		} else {
			// Type was created directly from JSDoc without a function or class
			// literal. e.g.
			// /**
			//  * @constructor
			//  * @param {string} x
			//  * @implements {SomeInterface}
			//  */
			// const MyImpl = createMyImpl();
			// The actual properties on this class are hidden from us, so we must
			// consider it ambiguous.
			constructorAmbiguity = FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
		}
	}
	return constructorAmbiguity
}

// See also TypedScopeCreator.DELEGATE_PROXY_SUFFIX Unfortunately we cannot
// use that constant here.
//
// Java: FunctionType.DELEGATE_SUFFIX
var FunctionType_DELEGATE_SUFFIX = ObjectType_CreateDelegateSuffix("Proxy")

// Java: FunctionType.IsDelegateProxy()
func (f *FunctionTypeBase) IsDelegateProxy() bool {
	// TODO(bradfordcsmith): There should be a better way to determine that we
	// have a proxy type.
	return f.selfFunc().HasReferenceName() && jutil.StringEndsWith(*f.selfFunc().GetReferenceName(), FunctionType_DELEGATE_SUFFIX)
}

// GetClosurePrimitive returns the @closurePrimitive identifier associated
// with this function
//
// Java: FunctionType.getClosurePrimitive()
func (f *FunctionTypeBase) GetClosurePrimitive() *ClosurePrimitive {
	if f.ext != nil {
		return f.ext.closurePrimitive
	}
	return nil
}

// Java: FunctionType.builder(JSTypeRegistry)
func FunctionType_Builder_(registry *JSTypeRegistry) *FunctionType_Builder {
	return NewFunctionType_Builder(registry)
}

func (f *FunctionTypeBase) initBuilder(b *FunctionType_Builder) {
	b.registry = f.GetRegistry()
	b.implicitPrototype = jutil.CheckNotNull(f.GetRegistry().GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE))
	b.googModuleId = f.selfFunc().GetGoogModuleId()
	b.className = f.selfFunc().GetReferenceName()
	b.nativeType = f.selfFunc().IsNativeObjectType()
	b.templateTypeMap = f.selfFunc().GetTemplateTypeMap()
	b.canonicalRepresentation = f.selfFunc().GetCanonicalRepresentation()
	b.isKnownAmbiguous = f.call.fnConstructorAmbiguity == FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	b.primitiveId = f.selfFunc().GetClosurePrimitive()
	b.isAbstract = f.selfFunc().IsAbstract()
	b.kind = f.selfFunc().GetKind()
	b.parameters = f.selfFunc().GetParameters()
	b.returnType = f.selfFunc().GetReturnType()
	b.returnTypeIsInferred = f.selfFunc().IsReturnTypeInferred()
	b.sourceNode = f.selfFunc().GetSource()
	b.typeOfThis = f.selfFunc().GetTypeOfThis()
	b.templateParamCount = f.selfFunc().GetTemplateParamCount()

	if f.selfFunc().IsConstructor() {
		b.WithConstructorTemplateKeys(f.selfFunc().GetConstructorOnlyTemplateParameters())
	}
}

// ToBuilder copies all the information from another function type.
//
// Java: FunctionType.toBuilder()
func (f *FunctionTypeBase) ToBuilder() *FunctionType_Builder {
	b := f.GetRegistry().acquireFunctionBuilder()
	b.self = b
	f.initBuilder(b)
	return b
}

// FunctionType_Builder is a builder class for function and arrow types.
type FunctionType_Builder struct {
	PrototypeObjectType_Builder[*FunctionType_Builder]

	sourceNode   Node
	googModuleId *string
	// Go: nil is Java null (unknown parameters).
	parameters              []*FunctionType_Parameter
	returnType              JSType
	typeOfThis              JSType
	setPrototypeBasedOn     ObjectType
	constructorOnlyKeys     *jutil.LinkedHashSet[*TemplateType]
	kind                    FunctionType_Kind
	isAbstract              bool
	isKnownAmbiguous        bool
	returnTypeIsInferred    bool
	returnsOwnInstanceType  bool
	pooled                  bool
	primitiveId             *ClosurePrimitive
	canonicalRepresentation FunctionType
}

// Java: FunctionType.Builder(JSTypeRegistry)
func NewFunctionType_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	b := registry.acquireFunctionBuilder()
	b.InitPrototypeObjectType_Builder(b, registry)
	b.kind = FunctionType_Kind_ORDINARY

	b.SetImplicitPrototype(
		jutil.CheckNotNull(registry.GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE)))
	return b
}

// WithName sets the name of the function type.
//
// Java: FunctionType.Builder.withName(String)
func (b *FunctionType_Builder) WithName(name *string) *FunctionType_Builder {
	return b.SetName(name)
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
	b.googModuleId = x
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
	b.parameters = paramBuilder.parameters
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

	return b.SetTemplateTypeMap(
		b.registry.
			GetEmptyTemplateTypeMap().
			CopyWithExtensionOfKeys(templateKeys, []JSType{})).
		// TODO(nickreid): This value should only consider ctor only keys.
		SetTemplateParamCount(len(templateKeys))
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
	b.constructorOnlyKeys = jutil.NewLinkedHashSet(jsTypeCopyOfTemplateTypes(constructorOnlyKeys)...)
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
	return b.SetNative(true)
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
	b.primitiveId = id
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
		jutil.CheckState(b.constructorOnlyKeys == nil || b.constructorOnlyKeys.IsEmpty())
	}

	res := NewFunctionType(b)
	b.registry.releaseFunctionBuilder(b)
	return res
}

// Java: FunctionType.Builder.buildAndResolve()
func (b *FunctionType_Builder) BuildAndResolve() FunctionType {
	return jutil.CheckNotNull(b.Build().ToMaybeFunctionType())
}

// Java: FunctionType.getCanonicalRepresentation()
func (f *FunctionTypeBase) GetCanonicalRepresentation() FunctionType {
	if f.ext != nil {
		return f.ext.canonicalRepresentation
	}
	return nil
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

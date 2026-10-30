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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GetInstanceType gets the type of instance of this function. May return
// null if the `this` type can not be converted to "ObjectType" (see
// JSType#toObjectType).
//
// Throws IllegalStateException if this function is not a constructor (see
// #isConstructor()).
//
// Java: FunctionType.getInstanceType()
func (f *FunctionTypeBase) GetInstanceType() ObjectType {
	jutil.CheckState(f.HasInstanceType())

	return f.typeOfThis.self().ToObjectType()
}

// HasInstanceType returns whether this function type has an instance type.
//
// Java: FunctionType.hasInstanceType()
func (f *FunctionTypeBase) HasInstanceType() bool {
	return f.IsConstructor() || f.IsInterface()
}

// Java: FunctionType.getPrototypeProperty()
func (f *FunctionTypeBase) GetPrototypeProperty() ObjectType {
	return f.GetPrototype()
}

// GetPrototype gets the prototype property of this function type. This is
// equivalent to (ObjectType) getPropertyType("prototype").
//
// Java: FunctionType.getPrototype()
func (f *FunctionTypeBase) GetPrototype() ObjectType {
	// lazy initialization of the prototype field
	if f.prototypeSlot() == nil {
		refName := f.GetReferenceName()
		if refName == nil {
			// Someone is trying to access the prototype of a structural function.
			// We don't want to give real properties to this prototype, because
			// then it would propagate to all structural functions.
			f.SetPrototypeNoCheck(f.GetRegistry().GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE), NilNode)
		} else {
			fields := PrototypeObjectType_Builder{
				registry:          f.GetRegistry(),
				classNameID:       rhinoStringPoolIntern(*refName + ".prototype"),
				implicitPrototype: f.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE),
				nativeType:        f.IsNativeObjectType(),
				templateTypeMap:   f.GetRegistry().GetEmptyTemplateTypeMap(),
			}
			f.SetPrototype(
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
	if f.IsNativeObjectType() {
		return
	}
	// First handle class-side inheritance for ES6 classes, before reassigning
	// baseType.
	if !f.source.IsNil() && f.source.IsClass() {
		superCtor := baseType.GetConstructor()
		if superCtor != nil {
			f.SetImplicitPrototype(superCtor)
		}
		f.MaybeLoosenTypecheckingDueToForwardReferencedSupertype(baseType)
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
	if baseType.HasReferenceName() || f.IsNativeObjectType() || baseType.IsFunctionPrototypeType() {
		if f.prototypeSlot() != nil && f.HasInstanceType() && baseType.Equals(f.GetInstanceType()) {
			// Bail out for cases like Foo.prototype = new Foo();
			return
		}
		fields := PrototypeObjectType_Builder{
			registry:          f.GetRegistry(),
			classNameID:       rhinoStringPoolIntern(jutil.ToJavaString(f.GetReferenceName()) + ".prototype"),
			implicitPrototype: baseType,
			templateTypeMap:   f.GetRegistry().GetEmptyTemplateTypeMap(),
		}
		baseType = NewPrototypeObjectType(&fields)
	}
	f.SetPrototype(baseType, propertyNode)
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
	if f.IsConstructor() && prototype == f.GetInstanceType() {
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
	prototype.SetOwnerFunction(f)

	if oldPrototype != nil {
		// Disassociating the old prototype makes this easier to debug-- we
		// don't have to worry about two prototypes running around.
		oldPrototype.SetOwnerFunction(nil)
	}

	if replacedPrototype {
		f.ClearCachedValues()
	}
	if reg := f.GetRegistry(); reg != nil && !reg.parallelReadOnly {
		reg.leastSupertypeGen++ // see PrototypeObjectTypeBase.SetImplicitPrototype
	}

	return true
}

// Java: FunctionType.defineProperty(Property.Key, JSType, boolean, Node) / NoObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (f *FunctionTypeBase) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	if f.isNoObjectOrNo() {
		// nothing, all properties are defined
		return true
	}
	if name.Matches("prototype") {
		objType := type_.ToObjectType()
		if objType != nil {
			if slot := f.prototypeSlot(); slot != nil && objType.Equals(slot.GetType()) {
				return true
			}
			f.SetPrototypeBasedOnWithPropertyNode(objType, propertyNode)
			return true
		}
		return false
	}
	return f.PrototypeObjectTypeBase.DefinePropertyByKey(name, type_, inferred, propertyNode)
}

// GetSuperClassConstructor: given a constructor or an interface type, get
// its superclass constructor or null if none exists.
//
// Java: FunctionType.getSuperClassConstructor()
func (f *FunctionTypeBase) GetSuperClassConstructor() FunctionType {
	jutil.CheckArgument(f.IsConstructor() || f.IsInterface())
	maybeSuperInstanceType := f.GetPrototype().GetImplicitPrototype()
	if maybeSuperInstanceType == nil {
		return nil
	}
	return maybeSuperInstanceType.GetConstructor()
}

// Java: FunctionType.hasImplementedInterfaces()
func (f *FunctionTypeBase) HasImplementedInterfaces() bool {
	if f.ext != nil && len(f.ext.implementedInterfaces) != 0 {
		return true
	}
	var superCtor FunctionType
	if f.IsConstructor() {
		superCtor = f.GetSuperClassConstructor()
	}
	if superCtor != nil {
		return superCtor.HasImplementedInterfaces()
	}
	return false
}

func appendUniqueObjectType(slice []ObjectType, elem ObjectType) []ObjectType {
	for _, existing := range slice {
		if existing == elem || (existing != nil && elem != nil && existing.Equals(elem)) {
			return slice
		}
	}
	return append(slice, elem)
}

func addRelatedInterfacesSlice(instance ObjectType, out *[]ObjectType) {
	constructor := instance.GetConstructor()
	if constructor == nil || !constructor.IsInterface() {
		return
	}
	prevLen := len(*out)
	*out = appendUniqueObjectType(*out, instance)
	if len(*out) == prevLen {
		return
	}
	for _, interfaceType := range instance.GetCtorExtendedInterfaces() {
		addRelatedInterfacesSlice(interfaceType, out)
	}
}

// GetAllImplementedInterfaces returns all interfaces implemented by a class
// or its superclass and any superclasses for any of those interfaces. If
// this is called before all types are resolved, it may return an incomplete
// set.
//
// Java: FunctionType.getAllImplementedInterfaces()
func (f *FunctionTypeBase) GetAllImplementedInterfaces() []ObjectType {
	direct := f.GetImplementedInterfaces()
	if len(direct) == 0 {
		return emptyObjectTypeSlice
	}
	interfaces := make([]ObjectType, 0, len(direct))
	for _, type_ := range direct {
		addRelatedInterfacesSlice(type_, &interfaces)
	}
	return interfaces
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
	var src []ObjectType
	if f.IsConstructor() {
		src = f.GetImplementedInterfaces()
	} else {
		src = f.GetExtendedInterfaces()
	}
	if len(src) == 0 {
		return emptyObjectTypeSlice
	}
	result := make([]ObjectType, 0, len(src))
	for _, t := range src {
		result = appendUniqueObjectType(result, t)
	}
	return result
}

// GetImplementedInterfaces returns interfaces implemented directly by a
// class or its superclass.
//
// Java: FunctionType.getImplementedInterfaces()
func (f *FunctionTypeBase) GetImplementedInterfaces() []ObjectType {
	var superCtor FunctionType
	if f.IsConstructor() {
		superCtor = f.GetSuperClassConstructor()
	}
	own := f.GetOwnImplementedInterfaces()
	if superCtor == nil {
		return own
	}
	builder := append([]ObjectType{}, own...)
	for superCtor != nil {
		builder = append(builder, superCtor.GetOwnImplementedInterfaces()...)
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
	jutil.CheckState(f.IsConstructor())

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
			panic(jutil.NewNullPointerExceptionMsg("at index " + strconv.Itoa(i)))
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
	jutil.CheckState(f.IsInterface())

	copied := functionTypeCopyOfObjectTypes(extendedInterfaces)
	if len(copied) > 0 || f.ext != nil {
		f.ensureExt().extendedInterfaces = copied
	}
	for _, extendedInterface := range extendedInterfaces {
		f.typeOfThis.self().MergeSupertypeTemplateTypes(extendedInterface)
	}
}

// SetImplicitMatch sets the current interface type to support structural
// interface matching (abbr. SMI). flag indicates whether or not it should
// support SMI.
//
// Java: FunctionType.setImplicitMatch(boolean)
func (f *FunctionTypeBase) SetImplicitMatch(flag bool) {
	jutil.CheckState(f.IsInterface())
	f.call.fnIsStructuralInterface = flag
}

// Java: FunctionType.isStructuralInterface()
func (f *FunctionTypeBase) IsStructuralInterface() bool {
	return f.IsInterface() && f.call.fnIsStructuralInterface
}

// Java: FunctionType.isAbstract()
func (f *FunctionTypeBase) IsAbstract() bool {
	return f.call.fnIsAbstract
}

func fnSliceContains(list []FunctionType, elem FunctionType) bool {
	for _, existing := range list {
		if existing == elem || (existing != nil && elem != nil && existing.Equals(elem)) {
			return true
		}
	}
	return false
}

// CheckExtendsLoop checks if there is a loop in the type extends chain.
// Returns an array of all functions in the loop chain if a loop exists,
// otherwise returns null.
//
// Java: FunctionType.checkExtendsLoop()
func (f *FunctionTypeBase) CheckExtendsLoop() []FunctionType {
	if len(f.GetExtendedInterfaces()) == 0 {
		return nil
	}
	var pathBuf [8]FunctionType
	path := pathBuf[:0]
	return f.checkExtendsLoopSlice(&path)
}

func (f *FunctionTypeBase) checkExtendsLoopSlice(path *[]FunctionType) []FunctionType {
	for _, interfaceType := range f.GetExtendedInterfaces() {
		superConstructor := interfaceType.GetConstructor()
		if superConstructor == nil {
			continue
		}
		if fnSliceContains(*path, superConstructor) {
			i := 0
			for (*path)[i] != superConstructor && !(*path)[i].Equals(superConstructor) {
				i++
			}
			loopLen := len(*path) - i + 1
			out := make([]FunctionType, loopLen)
			copy(out, (*path)[i:])
			out[loopLen-1] = superConstructor
			return out
		}
		*path = append(*path, superConstructor)
		result := superConstructor.checkExtendsLoopSlice(path)
		if result != nil {
			return result
		}
		(*path)[len(*path)-1] = nil
		*path = (*path)[:len(*path)-1]
	}
	return nil
}

// Java: FunctionType.checkExtendsLoop(Set<FunctionType>, List<FunctionType>)
func (f *FunctionTypeBase) CheckExtendsLoopWithCacheAndPath(
	cache *jutil.ObjectLinkedHashSet[FunctionType], path *[]FunctionType) []FunctionType {
	for _, interfaceType := range f.GetExtendedInterfaces() {
		superConstructor := interfaceType.GetConstructor()
		if superConstructor == nil {
			continue
		}
		if cache.Contains(superConstructor) {
			// after detecting a loop, prune and return the path, e.g.,: A -> B
			// -> C -> D -> C, will be pruned into: c -> D -> C
			*path = append(*path, superConstructor)
			i := 0
			for (*path)[i] != superConstructor {
				i++
			}
			return (*path)[i:]
		}
		cache.Add(superConstructor)
		*path = append(*path, superConstructor)
		result := superConstructor.CheckExtendsLoopWithCacheAndPath(cache, path)
		if result != nil {
			return result
		}
		cache.Remove(superConstructor)
		(*path)[len(*path)-1] = nil
		*path = (*path)[:len(*path)-1]
	}
	return nil
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

	if f.Equals(that) {
		return f
	}

	// If these are ordinary functions, then merge them. Don't do this if any
	// of the params/return values are unknown, because then there will be
	// cycles in their local lattice and they will merge in weird ways.
	if f.IsOrdinaryFunction() &&
		that.IsOrdinaryFunction() &&
		!f.call.HasUnknownParamsOrReturn() &&
		!that.call.HasUnknownParamsOrReturn() {

		// Check for the degenerate case, but double check that there's not a
		// cycle.
		isSubtypeOfThat := f.IsSubtype(that)
		isSubtypeOfThis := that.IsSubtype(f)
		if isSubtypeOfThat && !isSubtypeOfThis {
			if leastSuper {
				return that
			}
			return f
		} else if isSubtypeOfThis && !isSubtypeOfThat {
			if leastSuper {
				return f
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
		return f
	} else if functionInstance.Equals(f) {
		if leastSuper {
			return f
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
	var newParamsNode []*FunctionType_Parameter
	if NewEqualityChecker().
		SetEqMethod(EqualityChecker_EqMethod_IDENTITY).
		CheckParameters(&f.call, &other.call) {
		newParamsNode = f.call.GetParameterList()
	} else {
		// If the parameters are not equal, don't try to merge them. Someday, we
		// should try to merge the individual params.
		return nil
	}

	var newReturnType JSType
	if leastSuper {
		newReturnType = f.call.GetReturnType().GetLeastSupertype(other.call.GetReturnType())
	} else {
		newReturnType = f.call.GetReturnType().GetGreatestSubtype(other.call.GetReturnType())
	}

	var newTypeOfThis JSType
	if f.typeOfThis == other.typeOfThis ||
		(f.typeOfThis != nil && other.typeOfThis != nil && f.typeOfThis.self().Equals(other.typeOfThis.self())) {
		newTypeOfThis = jsTypeOf(f.typeOfThis)
	} else if leastSuper {
		newTypeOfThis = f.typeOfThis.self().GetLeastSupertype(jsTypeOf(other.typeOfThis))
	} else {
		newTypeOfThis = f.typeOfThis.self().GetGreatestSubtype(jsTypeOf(other.typeOfThis))
	}

	newReturnTypeInferred := f.call.returnTypeInferred || other.call.returnTypeInferred

	return FunctionType_Builder_(f.GetRegistry()).
		WithParametersOfList(newParamsNode).
		WithReturnTypeWithInferred(newReturnType, newReturnTypeInferred).
		WithTypeOfThis(newTypeOfThis).
		Build()
}

// ClearCachedValues: NOTE(sdh): One might assume that immediately after
// calling this, hasCachedValues() should always return false. This is not
// the case, since hasCachedValues() will return true if prototypeSlot is
// non-null, and this override does nothing to change that state.
//
// Java: FunctionType.clearCachedValues()
func (f *FunctionTypeBase) ClearCachedValues() {
	f.PrototypeObjectTypeBase.ClearCachedValues()

	if !f.IsNativeObjectType() {
		if f.HasInstanceType() {
			f.GetInstanceType().ClearCachedValues()
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

// Java: FunctionType.resolveInternal(ErrorReporter) / NoObjectType.resolveInternal(ErrorReporter)
func (f *FunctionTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	if f.isNoObjectOrNo() {
		panic(jutil.NewAssertionErrorNoMessage())
	}
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
				"not an object type: "+rt.String()+" (at "+f.String()+")",
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
	if f.IsUnknownType() {
		return FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	}
	if f.IsNativeObjectType() {
		// native types other than unknown are never ambiguous
		return FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
	}
	superConstructor := f.GetSuperClassConstructor()
	if superConstructor == nil {
		// TODO(bradfordcsmith): Why is superConstructor ever null here?
		return FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	}
	if superConstructor.IsAmbiguousConstructor() {
		// Subclasses of ambiguous objects are also ambiguous
		return FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
	}
	if !f.source.IsNil() {
		// We can see the definition of the class, so we know all properties it
		// directly declares or references. The same is true for its superclass
		// (previous condition).
		return FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
	}
	if f.IsDelegateProxy() {
		// Type was created by the compiler as a proxy that inherits from the
		// real type that was in the code. Since we've made it this far, we
		// know the real type creates unambiguous objects. Therefore, the
		// proxy does, too.
		return FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
	}
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
	return FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
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
	return f.HasReferenceName() && strings.HasSuffix(*f.GetReferenceName(), FunctionType_DELEGATE_SUFFIX)
}

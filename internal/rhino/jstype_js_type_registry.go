// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20261015

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSTypeRegistry is the type registry used to resolve named types.
//
// This class is not thread-safe.
type JSTypeRegistry struct {
	// The template variable in IObject<IOBJECT_KEY, IOBJECT_VALUE> (plus the
	// builtin Javascript Object).
	iObjectIndexTemplateKey *TemplateType

	// The template variable in IObject<IOBJECT_KEY, IOBJECT_VALUE> (plus the
	// builtin Javascript Object).
	iObjectElementTemplateKey *TemplateType

	// The template variable corresponding to the T type in Iterable<T,
	// TReturn, TNext>
	iterableValueTemplate *TemplateType

	// The template variable corresponding to the T type in Iterable<T,
	// TReturn, TNext>
	iterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in IteratorIterable<T,
	// TReturn, TNext>
	iteratorIterableValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// IteratorIterable<T, TReturn, TNext>
	iteratorIterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in IteratorLike<T,
	// TReturn, TNext>.
	iteratorLikeValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// IteratorLike<T, TReturn, TNext>.
	iteratorLikeReturnTemplate *TemplateType

	// The template variable corresponding to the T type in Iterator<T,
	// TReturn, TNext>.
	iteratorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in Iterator<T,
	// TReturn, TNext>.
	iteratorReturnTemplate *TemplateType

	// The template variable corresponding to the TYield type in
	// IIterableResult<TYield>
	iiterableResultValueTemplate *TemplateType

	// The template variable corresponding to the T type in AsyncIterable<T,
	// TReturn, TNext>
	asyncIterableValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// AsyncIterable<T, TReturn, TNext>
	asyncIterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in AsyncIterator<T,
	// TReturn, TNext>
	asyncIteratorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// AsyncIterator<T, TReturn, TNext>
	asyncIteratorReturnTemplate *TemplateType

	// The template variable corresponding to the T type in Generator<T,
	// TReturn, TNext>
	generatorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in Generator<T,
	// TReturn, TNext>
	generatorReturnTemplate *TemplateType

	// The template variable corresponding to the TYPE type in IThenable<TYPE>
	iThenableTemplateKey *TemplateType

	// The template variable corresponding to the TYPE in Promise<TYPE>
	promiseTemplateKey *TemplateType

	// The template variable in Array<T>.
	arrayElementTemplateKey *TemplateType

	// The template variable in ReadonlyArray<T>.
	readonlyArrayElementTemplateKey *TemplateType

	// The key variable in ReadonlyMap<K, V>.
	readonlyMapKeyTemplateKey *TemplateType

	// The value variable in ReadonlyMap<K, V>.
	readonlyMapValueTemplateKey *TemplateType

	// The key variable in Map<K, V>.
	mapKeyTemplateKey *TemplateType

	// The value variable in Map<K, V>.
	mapValueTemplateKey *TemplateType

	// TODO(user): An instance of this class should be used during
	// compilation. We also want to make all types' constructors package
	// private and force usage of this registry instead. This will allow us to
	// evolve the types without being tied by an open API.

	reporter ErrorReporter

	// We use an Array instead of an immutable list because this lookup needs
	// to be very fast. When it was an immutable list, we were spending 5% of
	// CPU time on bounds checking inside get().
	nativeTypes []JSType

	// Table<Node, String, JSType> (a HashBasedTable).
	scopedNameTable *jutil.LinkedHashMap[*Node, *jutil.LinkedHashMap[string, JSType]]

	// Only needed for type resolution at the moment
	closureNamespaces *jutil.LinkedHashMap[string, *JSTypeRegistry_ClosureNamespace]

	// NOTE: This would normally be "static final" but that causes unit test
	// failures when serializing and deserializing compiler state for
	// multistage builds.
	nameTableGlobalRoot *Node

	// NOTE(nicksantos): This is a terrible terrible hack. When type
	// expressions are evaluated, we need to be able to decide whether that
	// type name resolves to a nullable type or a non-nullable type. Object
	// types are nullable, but enum types and typedefs are not.
	//
	// Notice that it's not good enough to just declare enum types and
	// typedefs sooner. For example, if we have
	//   /** @enum {MyObject} */ var MyEnum = ...;
	// we won't be to declare "MyEnum" without evaluating the expression
	// {MyObject}, and following those dependencies starts to lead us into
	// undecidable territory. Instead, we "pre-declare" enum types and
	// typedefs, so that the expression resolver can decide whether a given
	// name is nullable or not. Also, note that this solution is buggy and the
	// type resolution still gets default nullability wrong in some cases. See
	// b/116853368. Probably NamedType should be responsible for applying
	// nullability to forward references instead of the type expression
	// evaluation.
	//
	// Multimap<Node, String> with hash keys and hash set values.
	nonNullableTypeNames *jutil.LinkedHashMap[*Node, *jutil.LinkedHashSet[string]]

	// Types that have been "forward-declared." If these types are not
	// declared anywhere in the binary, we shouldn't try to type-check them at
	// all.
	forwardDeclaredTypes jutil.Set[string]

	// A map of properties to the types on which those properties have been
	// declared. "Reference" types are excluded because those already exist in
	// eachRefTypeIndexedByProperty to avoid blowing up the size of this map.
	nonRefTypesIndexedByProperty map[string]*jutil.ObjectLinkedHashSet[JSType]

	sentinelObjectLiteral JSType

	// To avoid blowing up the size of nonRefTypesIndexedByProperty, we use the
	// sentinel object literal instead of registering arbitrarily many types.
	// But because of the way unions are constructed, some properties of
	// record types in unions are getting dropped and cause spurious
	// "non-existent property" warnings. The next two fields avoid the
	// warnings. The first field contains property names of records that
	// participate in unions, and have caused properties to be dropped. The
	// second field contains the names of the dropped properties. When
	// checking canPropertyBeDefined, if the type has a property in
	// propertiesOfSupertypesInUnions, we consider it to possibly have any
	// property in droppedPropertiesOfUnions. This is a loose check, but we
	// restrict it to records that may be present in unions, and it allows us
	// to keep nonRefTypesIndexedByProperty small.
	propertiesOfSupertypesInUnions *jutil.LinkedHashSet[string]
	droppedPropertiesOfUnions      *jutil.LinkedHashSet[string]

	// A map of properties to each reference type on which those properties
	// have been declared.
	eachRefTypeIndexedByProperty map[string]*jutil.ObjectLinkedHashSet[ObjectType]

	// A single empty TemplateTypeMap, which can be safely reused in cases
	// where there are no template types.
	emptyTemplateTypeMap *TemplateTypeMap

	resolver *JSTypeResolver
}

// Java: JSTypeRegistry.DOT_SPLITTER (Splitter.on('.'))
func JSTypeRegistry_DOT_SPLITTER_Split(s string) []string {
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			n++
		}
	}
	out := make([]string, 0, n)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// Java: JSTypeRegistry.I_OBJECT_ELEMENT_TEMPLATE
const JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE = "IOBJECT_VALUE"

// Deprecated: Java: JSTypeRegistry.OBJECT_ELEMENT_TEMPLATE
const JSTypeRegistry_OBJECT_ELEMENT_TEMPLATE = JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE

// Java: JSTypeRegistry(ErrorReporter)
func NewJSTypeRegistry(reporter ErrorReporter) *JSTypeRegistry {
	return NewJSTypeRegistryWithForwardDeclaredTypes(reporter, jutil.NewLinkedHashSet[string]())
}

// NewJSTypeRegistryWithForwardDeclaredTypes constructs a new type registry
// populated with the built-in types.
//
// Java: JSTypeRegistry(ErrorReporter, Set<String>)
func NewJSTypeRegistryWithForwardDeclaredTypes(reporter ErrorReporter, forwardDeclaredTypes jutil.Set[string]) *JSTypeRegistry {
	r := &JSTypeRegistry{
		scopedNameTable:                jutil.NewLinkedHashMap[*Node, *jutil.LinkedHashMap[string, JSType]](),
		closureNamespaces:              jutil.NewLinkedHashMap[string, *JSTypeRegistry_ClosureNamespace](),
		nameTableGlobalRoot:            NewNode(Token_ROOT),
		nonNullableTypeNames:           jutil.NewLinkedHashMap[*Node, *jutil.LinkedHashSet[string]](),
		nonRefTypesIndexedByProperty:   make(map[string]*jutil.ObjectLinkedHashSet[JSType]),
		propertiesOfSupertypesInUnions: jutil.NewLinkedHashSet[string](),
		droppedPropertiesOfUnions:      jutil.NewLinkedHashSet[string](),
		eachRefTypeIndexedByProperty:   make(map[string]*jutil.ObjectLinkedHashSet[ObjectType]),
	}
	r.reporter = reporter
	r.forwardDeclaredTypes = forwardDeclaredTypes
	r.emptyTemplateTypeMap = TemplateTypeMap_CreateEmpty(r)
	r.resolver = JSTypeResolver_Create(r)
	r.nativeTypes = make([]JSType, len(JSTypeNative_values()))

	func() {
		closer := r.resolver.OpenForDefinition()
		defer closer.Close()
		r.InitializeBuiltInTypes()
		r.InitializeRegistry()
	}()
	return r
}

// Java: JSTypeRegistry.getSentinelObjectLiteral()
func (r *JSTypeRegistry) GetSentinelObjectLiteral() JSType {
	if r.sentinelObjectLiteral == nil {
		r.sentinelObjectLiteral = r.CreateAnonymousObjectType(nil)
	}
	return r.sentinelObjectLiteral
}

// GetArrayElementKey returns the template variable for the element type of
// Arrays.
//
// Java: JSTypeRegistry.getArrayElementKey()
func (r *JSTypeRegistry) GetArrayElementKey() *TemplateType {
	return r.arrayElementTemplateKey
}

// GetReadonlyArrayElementKey returns the template variable for the element
// type of ReadonlyArrays.
//
// Java: JSTypeRegistry.getReadonlyArrayElementKey()
func (r *JSTypeRegistry) GetReadonlyArrayElementKey() *TemplateType {
	return r.readonlyArrayElementTemplateKey
}

// GetReadonlyMapKey returns the template variable for the key type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapKey()
func (r *JSTypeRegistry) GetReadonlyMapKey() *TemplateType {
	return r.readonlyMapKeyTemplateKey
}

// GetReadonlyMapValue returns the template variable for the value type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapValue()
func (r *JSTypeRegistry) GetReadonlyMapValue() *TemplateType {
	return r.readonlyMapValueTemplateKey
}

// GetObjectElementKey returns the template variable corresponding to the
// property value type for Javascript Objects and Arrays.
//
// Java: JSTypeRegistry.getObjectElementKey()
func (r *JSTypeRegistry) GetObjectElementKey() *TemplateType {
	return r.iObjectElementTemplateKey
}

// GetObjectIndexKey returns the template variable corresponding to the
// property key type of the built-in Javascript object.
//
// Java: JSTypeRegistry.getObjectIndexKey()
func (r *JSTypeRegistry) GetObjectIndexKey() *TemplateType {
	jutil.CheckNotNull(r.iObjectIndexTemplateKey)
	return r.iObjectIndexTemplateKey
}

// GetIterableValueTemplate returns the value template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableValueTemplate()
func (r *JSTypeRegistry) GetIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableValueTemplate)
}

// GetIterableReturnTemplate returns the return template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableReturnTemplate()
func (r *JSTypeRegistry) GetIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableReturnTemplate)
}

// GetIteratorIterableValueTemplate returns the value template variable for
// the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableValueTemplate()
func (r *JSTypeRegistry) GetIteratorIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableValueTemplate)
}

// GetIteratorIterableReturnTemplate returns the return template variable
// for the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableReturnTemplate()
func (r *JSTypeRegistry) GetIteratorIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableReturnTemplate)
}

// GetIIterableResultValueTemplate returns the value template variable for
// the IIterableResult interface.
//
// Java: JSTypeRegistry.getIIterableResultValueTemplate()
func (r *JSTypeRegistry) GetIIterableResultValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iiterableResultValueTemplate)
}

// GetIteratorLikeValueTemplate returns the value template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeValueTemplate()
func (r *JSTypeRegistry) GetIteratorLikeValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeValueTemplate)
}

// GetIteratorLikeReturnTemplate returns the return template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeReturnTemplate()
func (r *JSTypeRegistry) GetIteratorLikeReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeReturnTemplate)
}

// GetIteratorValueTemplate returns the value template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorValueTemplate()
func (r *JSTypeRegistry) GetIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorValueTemplate)
}

// GetIteratorReturnTemplate returns the return template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorReturnTemplate()
func (r *JSTypeRegistry) GetIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorReturnTemplate)
}

// GetGeneratorValueTemplate returns the value template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorValueTemplate()
func (r *JSTypeRegistry) GetGeneratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorValueTemplate)
}

// GetGeneratorReturnTemplate returns the return template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorReturnTemplate()
func (r *JSTypeRegistry) GetGeneratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorReturnTemplate)
}

// GetAsyncIterableValueTemplate returns the template variable for the
// AsyncIterable interface.
//
// Java: JSTypeRegistry.getAsyncIterableValueTemplate()
func (r *JSTypeRegistry) GetAsyncIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableValueTemplate)
}

// GetAsyncIterableReturnTemplate returns the template variable for the
// AsyncIterableR interface.
//
// Java: JSTypeRegistry.getAsyncIterableReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableReturnTemplate)
}

// GetAsyncIteratorValueTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorValueTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorValueTemplate)
}

// GetAsyncIteratorReturnTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorReturnTemplate)
}

// GetIThenableTemplate returns the template variable for the IThenable
// interface.
//
// Java: JSTypeRegistry.getIThenableTemplate()
func (r *JSTypeRegistry) GetIThenableTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iThenableTemplateKey)
}

// MaybeGetTemplateTypesOfBuiltin returns an immutable list of template types
// of the given builtin.
//
// Java: JSTypeRegistry.maybeGetTemplateTypesOfBuiltin(StaticScope, String)
func (r *JSTypeRegistry) MaybeGetTemplateTypesOfBuiltin(scope StaticScope, fnName string) []*TemplateType {
	type_ := r.GetType(scope, fnName)
	var objType ObjectType
	if type_ != nil {
		objType = type_.ToObjectType()
	}
	if objType != nil && objType.IsNativeObjectType() {
		return objType.GetTypeParameters()
	}
	return nil
}

// Java: JSTypeRegistry.getErrorReporter()
func (r *JSTypeRegistry) GetErrorReporter() ErrorReporter {
	return r.reporter
}

// Java: JSTypeRegistry.initializeBuiltInTypes()
func (r *JSTypeRegistry) InitializeBuiltInTypes() {
	// These locals shouldn't be all caps.
	booleanType := NewBooleanType(r)
	r.RegisterNativeType(JSTypeNative_BOOLEAN_TYPE, booleanType)

	nullType := NewNullType(r)
	r.RegisterNativeType(JSTypeNative_NULL_TYPE, nullType)

	bigIntType := NewBigIntType(r)
	r.RegisterNativeType(JSTypeNative_BIGINT_TYPE, bigIntType)

	numberType := NewNumberType(r)
	r.RegisterNativeType(JSTypeNative_NUMBER_TYPE, numberType)

	stringType := NewStringType(r)
	r.RegisterNativeType(JSTypeNative_STRING_TYPE, stringType)

	symbolType := NewSymbolType(r)
	r.RegisterNativeType(JSTypeNative_SYMBOL_TYPE, symbolType)

	unknownType := NewUnknownType(r, false)
	r.RegisterNativeType(JSTypeNative_UNKNOWN_TYPE, unknownType)
	checkedUnknownType := NewUnknownType(r, true)
	r.RegisterNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE, checkedUnknownType)

	voidType := NewVoidType(r)
	r.RegisterNativeType(JSTypeNative_VOID_TYPE, voidType)

	allType := NewAllType(r)
	r.RegisterNativeType(JSTypeNative_ALL_TYPE, allType)

	// Template Types
	// These should match the template type name in externs files.
	r.iObjectIndexTemplateKey = NewTemplateType(r, "IOBJECT_KEY")
	r.iObjectElementTemplateKey = NewTemplateType(r, JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE)

	// Array-related template types.
	iArrayLikeTemplate := NewTemplateType(r, "VALUE2")
	r.arrayElementTemplateKey = NewTemplateType(r, "T")
	r.readonlyArrayElementTemplateKey = NewTemplateType(r, "T")

	// Map-related template types.
	r.readonlyMapKeyTemplateKey = NewTemplateType(r, "K")
	r.readonlyMapValueTemplateKey = NewTemplateType(r, "V")
	r.mapKeyTemplateKey = NewTemplateType(r, "K")
	r.mapValueTemplateKey = NewTemplateType(r, "V")

	// Promise-related template types.
	r.iThenableTemplateKey = NewTemplateType(r, "TYPE")
	r.promiseTemplateKey = NewTemplateType(r, "TYPE")

	// The default implicit prototype of all functions, as well as the
	// ".prototype" field of `(typeof Function)`.
	//
	// This is the interesting `Function` protoype, as all its properties are
	// inhertied by other functions.
	//
	// Must be created and registered early so `FunctionType.Builder` has
	// access to it.
	functionPrototype := PrototypeObjectType_Builder_(r).
		SetName(jutil.StringPtr("Function.prototype")).
		SetNative(true).
		// Defer until `Object` is defined. .setImplicitPrototype(objectType)
		Build()
	r.RegisterNativeType(JSTypeNative_FUNCTION_PROTOTYPE, functionPrototype)

	// The ".prototype" property of the type `Function`.
	//
	// So named because `Function` constructs `?`s. This type is not
	// particularly interesting, since its properties aren't inherited by any
	// other type.
	functionInstancePrototype := PrototypeObjectType_Builder_(r).
		SetName(jutil.StringPtr("?.prototype")).
		SetNative(true).
		// Defer until `Object` is defined. .setImplicitPrototype(objectType)
		Build()
	r.RegisterNativeType(JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE, functionInstancePrototype)

	// `Function`
	//
	// The default implicit prototype of all `FunctionType`s is
	// `functionPrototype`.
	functionType := FunctionType_Builder_(r).
		WithName(jutil.StringPtr("Function")).
		ForConstructor().
		ForNativeType().
		WithParametersOfList(r.CreateParametersWithVarArgs(unknownType)).
		WithTypeOfThis(unknownType).
		WithReturnType(unknownType).
		Build()
	functionType.SetPrototype(functionInstancePrototype, nil)
	r.RegisterNativeType(JSTypeNative_FUNCTION_TYPE, functionType)

	// `(typeof Function)`
	//
	// The default implict prototype of all `FunctionType`s is
	// `functionPrototype`.
	functionFunctionType := FunctionType_Builder_(r).
		WithName(jutil.StringPtr("Function")).
		ForConstructor().
		ForNativeType().
		WithParametersOfList(r.CreateParametersWithVarArgs(allType)).
		WithTypeOfThis(functionType).
		// TODO(nickreid): .withReturnsOwnInstanceType()
		Build()
	functionFunctionType.SetPrototype(functionPrototype, nil)
	r.RegisterNativeType(JSTypeNative_FUNCTION_FUNCTION_TYPE, functionFunctionType)

	// `Object.prototype`
	objectPrototype := PrototypeObjectType_Builder_(r).
		SetName(jutil.StringPtr("Object.prototype")).
		SetNative(true).
		SetImplicitPrototype(nil).
		Build()
	r.RegisterNativeType(JSTypeNative_OBJECT_PROTOTYPE, objectPrototype)

	// Object
	objectFunctionType := r.NativeConstructorBuilder("Object").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnsOwnInstanceType().
		WithTemplateKeys(r.iObjectIndexTemplateKey, r.iObjectElementTemplateKey).
		Build()
	objectFunctionType.SetPrototype(objectPrototype, nil)
	r.RegisterNativeType(JSTypeNative_OBJECT_FUNCTION_TYPE, objectFunctionType)

	objectType := objectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_OBJECT_TYPE, objectType)

	functionPrototype.ClearCachedValues()
	functionPrototype.SetImplicitPrototype(objectType)
	functionInstancePrototype.ClearCachedValues()
	functionInstancePrototype.SetImplicitPrototype(objectType)

	// IObject
	iObjectFunctionType := r.NativeInterface("IObject", r.iObjectIndexTemplateKey, r.iObjectElementTemplateKey)
	r.RegisterNativeType(JSTypeNative_I_OBJECT_FUNCTION_TYPE, iObjectFunctionType)
	iObjectType := iObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_I_OBJECT_TYPE, iObjectType)

	// Bottom Types
	noType := NewNoType(r)
	r.RegisterNativeType(JSTypeNative_NO_TYPE, noType)

	noObjectType := NewNoObjectType(r)
	r.RegisterNativeType(JSTypeNative_NO_OBJECT_TYPE, noObjectType)

	// TODO(b/142881197): start using unused iterator (and related type)
	// template params https://github.com/google/closure-compiler/issues/3489

	r.iterableValueTemplate = NewTemplateType(r, "T")
	r.iterableReturnTemplate = NewTemplateType(r, "TReturn")
	iterableNextTemplate := NewTemplateType(r, "TNext")
	iterableFunctionType := r.NativeInterface(
		"Iterable", r.iterableValueTemplate, r.iterableReturnTemplate, iterableNextTemplate)
	r.RegisterNativeType(JSTypeNative_ITERABLE_FUNCTION_TYPE, iterableFunctionType)
	iterableType := iterableFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_ITERABLE_TYPE, iterableType)

	r.iteratorLikeValueTemplate = NewTemplateType(r, "T")
	r.iteratorLikeReturnTemplate = NewTemplateType(r, "TReturn")
	iteratorLikeNextTemplate := NewTemplateType(r, "TNext")
	iteratorLikeFunctionType := r.NativeInterface(
		"IteratorLike",
		r.iteratorLikeValueTemplate,
		r.iteratorLikeReturnTemplate,
		iteratorLikeNextTemplate)
	r.RegisterNativeType(JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE, iteratorLikeFunctionType)
	iteratorLikeType := iteratorLikeFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_ITERATOR_LIKE_TYPE, iteratorLikeType)

	r.iteratorValueTemplate = NewTemplateType(r, "T")
	r.iteratorReturnTemplate = NewTemplateType(r, "TReturn")
	iteratorNextTemplate := NewTemplateType(r, "TNext")

	r.iteratorIterableValueTemplate = NewTemplateType(r, "T")
	r.iteratorIterableReturnTemplate = NewTemplateType(r, "TReturn")
	iteratorIterableNextTemplate := NewTemplateType(r, "TNext")
	iteratorIterableFunctionType := r.NativeInterface(
		"IteratorIterable",
		r.iteratorIterableValueTemplate,
		r.iteratorIterableReturnTemplate,
		iteratorIterableNextTemplate)
	r.RegisterNativeType(JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE, iteratorIterableFunctionType)
	iteratorIterableFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				iterableType,
				r.iteratorIterableValueTemplate,
				r.iteratorIterableReturnTemplate,
				iteratorIterableNextTemplate),
			r.CreateTemplatizedType(
				iteratorLikeType,
				r.iteratorIterableValueTemplate,
				r.iteratorIterableReturnTemplate,
				iteratorIterableNextTemplate),
		})
	iteratorIterableType := iteratorIterableFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_ITERATOR_ITERABLE_TYPE, iteratorIterableType)
	iteratorFunctionType := r.NativeConstructorBuilder("Iterator").
		WithParameters().
		WithTemplateKeys(
			r.iteratorValueTemplate, r.iteratorReturnTemplate, iteratorNextTemplate).
		Build()
	r.RegisterNativeType(JSTypeNative_ITERATOR_FUNCTION_TYPE, iteratorFunctionType)
	iteratorFunctionType.GetPrototype() // Force initialization
	iteratorFunctionType.SetImplementedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				iteratorIterableType,
				r.iteratorValueTemplate,
				r.iteratorReturnTemplate,
				iteratorNextTemplate),
		})
	iteratorType := iteratorFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_ITERATOR_TYPE, iteratorType)

	r.iiterableResultValueTemplate = NewTemplateType(r, "TYield")
	iiterableResultFunctionType := r.NativeRecord("IIterableResult", r.iiterableResultValueTemplate)
	r.RegisterNativeType(JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE, iiterableResultFunctionType)
	iiterableResultType := iiterableResultFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_I_ITERABLE_RESULT_TYPE, iiterableResultType)

	// IArrayLike.
	iArrayLikeFunctionType := r.NativeRecord("IArrayLike", iArrayLikeTemplate)
	iArrayLikeFunctionType.SetExtendedInterfaces(
		[]ObjectType{r.CreateTemplatizedType(iObjectType, numberType, iArrayLikeTemplate)})
	r.RegisterNativeType(JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE, iArrayLikeFunctionType)
	iArrayLikeType := iArrayLikeFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_I_ARRAY_LIKE_TYPE, iArrayLikeType)

	// ReadonlyArray
	readonlyArrayFunctionType := r.NativeRecord("ReadonlyArray", r.readonlyArrayElementTemplateKey)
	r.RegisterNativeType(JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE, readonlyArrayFunctionType)
	readonlyArrayFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(iArrayLikeType, r.readonlyArrayElementTemplateKey),
			r.CreateTemplatizedType(iterableType, r.readonlyArrayElementTemplateKey),
		})
	readonlyArrayType := readonlyArrayFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_READONLY_ARRAY_TYPE, readonlyArrayType)

	// Array
	arrayFunctionType := r.NativeConstructorBuilder("Array").
		WithParametersOfList(r.CreateParametersWithVarArgs(allType)).
		WithReturnsOwnInstanceType().
		WithTemplateKeys(r.arrayElementTemplateKey).
		Build()
	arrayFunctionType.GetPrototype() // Force initialization
	arrayFunctionType.SetImplementedInterfaces(
		[]ObjectType{r.CreateTemplatizedType(readonlyArrayType, r.arrayElementTemplateKey)})
	r.RegisterNativeType(JSTypeNative_ARRAY_FUNCTION_TYPE, arrayFunctionType)

	arrayType := arrayFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_ARRAY_TYPE, arrayType)

	// ITemplateArray extends !Array<string>
	arrayOfString := r.CreateTemplatizedType(arrayType, stringType)
	iTemplateArrayFunctionType := r.NativeConstructorBuilder("ITemplateArray").WithParameters().Build()
	iTemplateArrayPrototype := PrototypeObjectType_Builder_(r).
		SetName(jutil.StringPtr("ITemplateArray.prototype")).
		SetNative(true).
		SetImplicitPrototype(arrayOfString).
		Build()
	iTemplateArrayFunctionType.SetPrototype(iTemplateArrayPrototype, nil)
	iTemplateArrayFunctionType.GetInstanceType().MergeSupertypeTemplateTypes(arrayOfString)
	r.RegisterNativeType(
		JSTypeNative_I_TEMPLATE_ARRAY_TYPE, iTemplateArrayFunctionType.GetInstanceType())

	r.generatorValueTemplate = NewTemplateType(r, "T")
	r.generatorReturnTemplate = NewTemplateType(r, "TReturn")
	generatorNextTemplate := NewTemplateType(r, "TNext")
	generatorFunctionType := r.NativeInterface(
		"Generator", r.generatorValueTemplate, r.generatorReturnTemplate, generatorNextTemplate)
	generatorFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				iteratorIterableType,
				r.generatorValueTemplate,
				r.generatorReturnTemplate,
				generatorNextTemplate),
		})
	r.RegisterNativeType(JSTypeNative_GENERATOR_FUNCTION_TYPE, generatorFunctionType)
	r.RegisterNativeType(JSTypeNative_GENERATOR_TYPE, generatorFunctionType.GetInstanceType())

	r.asyncIteratorValueTemplate = NewTemplateType(r, "T")
	r.asyncIteratorReturnTemplate = NewTemplateType(r, "TReturn")
	asyncIteratorNextTemplate := NewTemplateType(r, "TNext")
	asyncIteratorFunctionType := r.NativeInterface(
		"AsyncIterator",
		r.asyncIteratorValueTemplate,
		r.asyncIteratorReturnTemplate,
		asyncIteratorNextTemplate)
	r.RegisterNativeType(JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE, asyncIteratorFunctionType)
	r.RegisterNativeType(
		JSTypeNative_ASYNC_ITERATOR_TYPE, asyncIteratorFunctionType.GetInstanceType())

	r.asyncIterableValueTemplate = NewTemplateType(r, "T")
	r.asyncIterableReturnTemplate = NewTemplateType(r, "TReturn")
	asyncIterableNextTemplate := NewTemplateType(r, "TNext")
	asyncIterableFunctionType := r.NativeInterface(
		"AsyncIterable",
		r.asyncIterableValueTemplate,
		r.asyncIterableReturnTemplate,
		asyncIterableNextTemplate)
	r.RegisterNativeType(JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE, asyncIterableFunctionType)
	r.RegisterNativeType(
		JSTypeNative_ASYNC_ITERABLE_TYPE, asyncIterableFunctionType.GetInstanceType())

	asyncIteratorIterableTemplate := NewTemplateType(r, "T")
	asyncIteratorIterableReturnTemplate := NewTemplateType(r, "TReturn")
	asyncIteratorIterableNextTemplate := NewTemplateType(r, "TNext")
	asyncIteratorIterableFunctionType := r.NativeInterface(
		"AsyncIteratorIterable",
		asyncIteratorIterableTemplate,
		asyncIteratorIterableReturnTemplate,
		asyncIteratorIterableNextTemplate)
	asyncIteratorIterableFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				asyncIteratorFunctionType.GetInstanceType(),
				asyncIteratorIterableTemplate,
				asyncIteratorIterableReturnTemplate,
				asyncIteratorIterableNextTemplate),
			r.CreateTemplatizedType(
				asyncIterableFunctionType.GetInstanceType(),
				asyncIteratorIterableTemplate,
				asyncIteratorIterableReturnTemplate,
				asyncIteratorIterableNextTemplate),
		})

	r.RegisterNativeType(
		JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE, asyncIteratorIterableFunctionType)
	r.RegisterNativeType(
		JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE,
		asyncIteratorIterableFunctionType.GetInstanceType())

	asyncGeneratorValueTemplate := NewTemplateType(r, "T")
	asyncGeneratorReturnTemplate := NewTemplateType(r, "TReturn")
	asyncGeneratorNextTemplate := NewTemplateType(r, "TNext")
	asyncGeneratorFunctionType := r.NativeInterface(
		"AsyncGenerator",
		asyncGeneratorValueTemplate,
		asyncGeneratorReturnTemplate,
		asyncGeneratorNextTemplate)
	r.RegisterNativeType(JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE, asyncGeneratorFunctionType)
	asyncGeneratorFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				asyncIteratorIterableFunctionType.GetInstanceType(),
				asyncGeneratorValueTemplate,
				asyncGeneratorReturnTemplate,
				asyncGeneratorNextTemplate),
		})
	r.RegisterNativeType(
		JSTypeNative_ASYNC_GENERATOR_TYPE, asyncGeneratorFunctionType.GetInstanceType())

	ithenableFunctionType := r.NativeInterface("IThenable", r.iThenableTemplateKey)
	ithenableFunctionType.SetStruct()
	r.RegisterNativeType(JSTypeNative_I_THENABLE_FUNCTION_TYPE, ithenableFunctionType)
	ithenableType := ithenableFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_I_THENABLE_TYPE, ithenableType)

	// Thenable is an @typedef
	thenProps := jutil.NewLinkedHashMap[string, JSType]()
	thenProps.Put("then", unknownType)
	thenableType := r.CreateRecordType(thenProps)
	r.IdentifyNonNullableName(nil, "Thenable")
	r.RegisterNativeType(JSTypeNative_THENABLE_TYPE, thenableType)

	// Create built-in Promise type, whose constructor takes one parameter.
	// @param {function(
	//             function((TYPE|IThenable<TYPE>|Thenable|null)=),
	//             function(*=))} resolver
	promiseParameterType := r.CreateFunctionType(
		/* returnType= */ unknownType,
		/* parameterTypes...= */ r.CreateFunctionTypeWithParameters(
			unknownType,
			r.CreateOptionalParameters(
				r.CreateUnionType(
					r.promiseTemplateKey,
					r.CreateTemplatizedType(ithenableType, r.promiseTemplateKey),
					thenableType,
					nullType))),
		r.CreateFunctionTypeWithParameters(unknownType, r.CreateOptionalParameters(allType)))

	promiseFunctionType := r.NativeConstructorBuilder("Promise").
		WithParametersOfList(r.CreateParameters(promiseParameterType)).
		WithTemplateKeys(r.promiseTemplateKey).
		Build()
	promiseFunctionType.SetImplementedInterfaces(
		[]ObjectType{r.CreateTemplatizedType(ithenableType, r.promiseTemplateKey)})

	r.RegisterNativeType(JSTypeNative_PROMISE_FUNCTION_TYPE, promiseFunctionType)
	r.RegisterNativeType(JSTypeNative_PROMISE_TYPE, promiseFunctionType.GetInstanceType())

	// Arguments
	argumentsFunctionType := r.NativeConstructorBuilder("Arguments").WithParameters().Build()
	argumentsFunctionType.SetImplementedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(iArrayLikeType, unknownType),
			r.CreateTemplatizedType(iterableType, unknownType),
		})
	r.RegisterNativeType(JSTypeNative_ARGUMENTS_FUNCTION_TYPE, argumentsFunctionType)
	r.RegisterNativeType(JSTypeNative_ARGUMENTS_TYPE, argumentsFunctionType.GetInstanceType())

	// (bigint,number,string)
	bigintNumberString := r.CreateUnionType(bigIntType, numberType, stringType)
	r.RegisterNativeType(JSTypeNative_BIGINT_NUMBER_STRING, bigintNumberString)

	// BigInt
	bigIntObjectFunctionType := r.NativeConstructorBuilder("BigInt").
		WithParametersOfList(r.CreateParameters(bigintNumberString)).
		WithReturnType(bigIntType).
		Build()
	bigIntObjectFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE, bigIntObjectFunctionType)

	bigIntObjectType := bigIntObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_BIGINT_OBJECT_TYPE, bigIntObjectType)

	// Boolean
	booleanObjectFunctionType := r.NativeConstructorBuilder("Boolean").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnType(booleanType).
		Build()
	booleanObjectFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE, booleanObjectFunctionType)

	booleanObjectType := booleanObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_BOOLEAN_OBJECT_TYPE, booleanObjectType)

	// Date
	dateFunctionType := r.NativeConstructorBuilder("Date").
		WithParametersOfList(
			r.CreateOptionalParameters(
				unknownType,
				unknownType,
				unknownType,
				unknownType,
				unknownType,
				unknownType,
				unknownType)).
		WithReturnType(stringType).
		Build()
	dateFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_DATE_FUNCTION_TYPE, dateFunctionType)

	dateType := dateFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_DATE_TYPE, dateType)

	// Number
	numberObjectFunctionType := r.NativeConstructorBuilder("Number").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnType(numberType).
		Build()
	numberObjectFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE, numberObjectFunctionType)

	numberObjectType := numberObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_NUMBER_OBJECT_TYPE, numberObjectType)

	// RegExp
	regexpFunctionType := r.NativeConstructorBuilder("RegExp").
		WithParametersOfList(r.CreateOptionalParameters(allType, allType)).
		WithReturnsOwnInstanceType().
		Build()
	regexpFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_REGEXP_FUNCTION_TYPE, regexpFunctionType)

	regexpType := regexpFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_REGEXP_TYPE, regexpType)

	// String
	stringObjectFunctionType := r.NativeConstructorBuilder("String").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnType(stringType).
		Build()
	stringObjectFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_STRING_OBJECT_FUNCTION_TYPE, stringObjectFunctionType)

	stringObjectType := stringObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_STRING_OBJECT_TYPE, stringObjectType)

	// Symbol
	// NOTE: While "Symbol" is a class, with an instance type and prototype it
	// is illegal to call "new Symbol". This is checked in the type checker.
	symbolObjectFunctionType := r.NativeConstructorBuilder("Symbol").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnType(symbolType).
		Build()
	symbolObjectFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE, symbolObjectFunctionType)

	symbolObjectType := symbolObjectFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE, symbolObjectType)

	gbigintFunctionType := r.NativeInterface("gbigint")
	gbigintFunctionType.GetPrototype() // Force initialization
	r.RegisterNativeType(JSTypeNative_GBIGINT_TYPE, gbigintFunctionType.GetInstanceType())

	// (null|void)
	nullVoid := r.CreateUnionType(nullType, voidType)
	r.RegisterNativeType(JSTypeNative_NULL_VOID, nullVoid)

	// (string|number|boolean)
	numberStringBoolean := r.CreateUnionType(numberType, stringType, booleanType)
	r.RegisterNativeType(JSTypeNative_NUMBER_STRING_BOOLEAN, numberStringBoolean)

	// (string|number|boolean|symbol)
	valueTypes := r.CreateUnionType(numberType, stringType, booleanType, symbolType)
	r.RegisterNativeType(JSTypeNative_VALUE_TYPES, valueTypes)

	// (number|symbol)
	numberSymbol := r.CreateUnionType(numberType, symbolType)
	r.RegisterNativeType(JSTypeNative_NUMBER_SYMBOL, numberSymbol)

	// (string|symbol)
	stringSymbol := r.CreateUnionType(stringType, symbolType)
	r.RegisterNativeType(JSTypeNative_STRING_SYMBOL, stringSymbol)

	// (string,number)
	numberString := r.CreateUnionType(numberType, stringType)
	r.RegisterNativeType(JSTypeNative_NUMBER_STRING, numberString)

	// (bigint,number)
	bigintNumber := r.CreateUnionType(bigIntType, numberType)
	r.RegisterNativeType(JSTypeNative_BIGINT_NUMBER, bigintNumber)

	// (BigInt,Number)
	bigintNumberObject := r.CreateUnionType(bigIntObjectType, numberObjectType)
	r.RegisterNativeType(JSTypeNative_BIGINT_NUMBER_OBJECT, bigintNumberObject)

	// (Bigint,Number,String)
	bigintNumberStringObject := r.CreateUnionType(bigIntObjectType, numberObjectType, stringObjectType)
	r.RegisterNativeType(JSTypeNative_BIGINT_NUMBER_STRING_OBJECT, bigintNumberStringObject)

	// (number,string,symbol)
	numberStringSymbol := r.CreateUnionType(numberType, stringType, symbolType)
	r.RegisterNativeType(JSTypeNative_NUMBER_STRING_SYMBOL, numberStringSymbol)

	// (boolean,Boolean,number,Number,null,undefined)
	numberAdditionSupertype := r.CreateUnionTypeOfNatives(
		JSTypeNative_VOID_TYPE,
		JSTypeNative_NULL_TYPE,
		JSTypeNative_NUMBER_TYPE,
		JSTypeNative_NUMBER_OBJECT_TYPE,
		JSTypeNative_BOOLEAN_TYPE,
		JSTypeNative_BOOLEAN_OBJECT_TYPE)
	r.RegisterNativeType(JSTypeNative_NUMBER_ADDITION_SUPERTYPE, numberAdditionSupertype)

	// Native object properties are filled in by externs...

	// least function type, i.e. (All...) -> NoType
	leastFunctionType := r.CreateNativeFunctionTypeWithVarArgs(noType, allType)
	r.RegisterNativeType(JSTypeNative_LEAST_FUNCTION_TYPE, leastFunctionType)

	// the 'this' object in the global scope
	globalThisCtor := r.NativeConstructorBuilder("global this").
		WithParametersOfList(r.CreateParameters(allType)).
		WithReturnType(numberType).
		Build()
	globalThis := globalThisCtor.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_GLOBAL_THIS, globalThis)

	// greatest function type, i.e. (NoType...) -> All
	greatestFunctionType := r.CreateNativeFunctionTypeWithVarArgs(allType, noType)
	r.RegisterNativeType(JSTypeNative_GREATEST_FUNCTION_TYPE, greatestFunctionType)

	// Register the prototype property. See the comments below in
	// registerPropertyOnType about the bootstrapping process.
	r.RegisterPropertyOnType("prototype", objectFunctionType)

	// ReadonlyMap
	readonlyMapFunctionType := r.NativeInterface("ReadonlyMap", r.readonlyMapKeyTemplateKey, r.readonlyMapValueTemplateKey)
	r.RegisterNativeType(JSTypeNative_READONLY_MAP_FUNCTION_TYPE, readonlyMapFunctionType)
	readonlyMapFunctionType.SetExtendedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(
				iterableType,
				r.CreateTemplatizedType(
					arrayType,
					r.CreateUnionType(r.readonlyMapKeyTemplateKey, r.readonlyMapValueTemplateKey))),
		})
	readonlyMapType := readonlyMapFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_READONLY_MAP_TYPE, readonlyMapType)

	// Map
	mapFunctionType := r.NativeConstructorBuilder("Map").
		WithParametersOfList(
			r.CreateOptionalParameters(
				r.CreateUnionType(
					nullType,
					r.CreateTemplatizedType(
						iterableType,
						r.CreateTemplatizedType(
							arrayType,
							r.CreateUnionType(r.mapKeyTemplateKey, r.mapValueTemplateKey))),
					r.CreateTemplatizedType(
						arrayType,
						r.CreateTemplatizedType(
							arrayType,
							r.CreateUnionType(r.mapKeyTemplateKey, r.mapValueTemplateKey)))))).
		WithTemplateKeys(r.mapKeyTemplateKey, r.mapValueTemplateKey).
		Build()
	r.RegisterNativeType(JSTypeNative_MAP_FUNCTION_TYPE, mapFunctionType)
	mapFunctionType.GetPrototype() // Force initialization
	mapFunctionType.SetImplementedInterfaces(
		[]ObjectType{
			r.CreateTemplatizedType(readonlyMapType, r.mapKeyTemplateKey, r.mapValueTemplateKey),
		})
	mapType := mapFunctionType.GetInstanceType()
	r.RegisterNativeType(JSTypeNative_MAP_TYPE, mapType)
}

// Java: JSTypeRegistry.initializeRegistry()
func (r *JSTypeRegistry) InitializeRegistry() {
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ARGUMENTS_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ARRAY_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_READONLY_ARRAY_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_MAP_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_READONLY_MAP_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ASYNC_ITERABLE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ASYNC_ITERATOR_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ASYNC_GENERATOR_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_BIGINT_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_BIGINT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_GBIGINT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_BOOLEAN_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_BOOLEAN_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_I_ARRAY_LIKE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ITERABLE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ITERATOR_LIKE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ITERATOR_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_ITERATOR_ITERABLE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_GENERATOR_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_DATE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_I_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_I_ITERABLE_RESULT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_I_TEMPLATE_ARRAY_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_I_THENABLE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_NULL_TYPE))
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_NULL_TYPE), "Null")
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_NUMBER_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_NUMBER_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_PROMISE_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_REGEXP_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_STRING_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_STRING_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE))
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_SYMBOL_TYPE))
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_THENABLE_TYPE), "Thenable")
	r.RegisterGlobalType(r.GetNativeType(JSTypeNative_VOID_TYPE))
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_VOID_TYPE), "Undefined")
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_VOID_TYPE), "void")
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_FUNCTION_TYPE), "Function")
	r.RegisterGlobalTypeWithName(r.GetNativeType(JSTypeNative_GLOBAL_THIS), "Global")
}

// Java: JSTypeRegistry.getRootElementOfName(String)
func JSTypeRegistry_GetRootElementOfName(name string) string {
	index := strings.IndexByte(name, '.')
	if index != -1 {
		return name[:index]
	}
	return name
}

// JSTypeRegistry_GetLookupScope returns which scope in the provided scope
// chain the provided name is declared in, or else null.
//
// This assumed that the Scope construction is complete. It can not be used
// during scope construction to determine if a name is already defined as a
// shadowed name from a parent scope would be returned.
//
// Java: JSTypeRegistry.getLookupScope(StaticScope, String)
func JSTypeRegistry_GetLookupScope(scope StaticScope, name string) StaticScope {
	if scope != nil && scope.GetParentScope() != nil {
		return scope.GetTopmostScopeOfEventualDeclaration(JSTypeRegistry_GetRootElementOfName(name))
	}
	return scope
}

// Java: JSTypeRegistry.getRootNodeForScope(StaticScope)
func (r *JSTypeRegistry) GetRootNodeForScope(scope StaticScope) *Node {
	var root *Node
	if scope != nil {
		root = scope.GetRootNode()
	}
	if root == nil || root.IsRoot() || root.IsScript() {
		return r.nameTableGlobalRoot
	}
	return root
}

// Java: JSTypeRegistry.isDeclaredForScope(StaticScope, String)
func (r *JSTypeRegistry) IsDeclaredForScope(scope StaticScope, name string) bool {
	return r.GetTypeInternal(scope, name) != nil
}

// Java: JSTypeRegistry.checkTypeName(String)
func JSTypeRegistry_CheckTypeName(typeName string) {
	if strings.IndexByte(typeName, '<') >= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Type names cannot contain template annotations."))
	}
}

// Java: JSTypeRegistry.getTypeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeInternal(scope StaticScope, name string) JSType {
	JSTypeRegistry_CheckTypeName(name)
	if syntheticTemplateScope, ok := scope.(*JSTypeRegistry_SyntheticTemplateScope); ok {
		type_ := syntheticTemplateScope.GetTemplateType(name)
		if type_ != nil {
			return type_
		}
	}
	declarationScope := JSTypeRegistry_GetLookupScope(scope, name)
	resolvedViaTable := r.GetTypeForScopeInternal(declarationScope, name)
	if resolvedViaTable != nil {
		return resolvedViaTable
	}
	bestScope := scope
	if declarationScope != nil {
		bestScope = declarationScope
	}
	return r.ResolveViaComponents(bestScope, name)
}

// Java: JSTypeRegistry.resolveViaComponents(StaticScope, String)
func (r *JSTypeRegistry) ResolveViaComponents(scope StaticScope, qualifiedName string) JSType {
	resolutionScope, ok := scope.(StaticTypedScope)
	if qualifiedName == "" || !ok {
		return nil
	}
	// Skip closure namespace resolution of types whose root component is
	// defined in a local scope (not a global scope). Those will follow the
	// normal resolution scheme. (For legacy compatibility reasons we don't
	// check for global names that are the same as the module root).
	if !JSTypeRegistry_IsNameDefinedLocally(resolutionScope, JSTypeRegistry_GetRootElementOfName(qualifiedName)) {
		resolvedViaClosureNamespace := r.ResolveViaClosureNamespace(qualifiedName)
		if resolvedViaClosureNamespace != nil {
			return resolvedViaClosureNamespace
		}
	}
	return r.ResolveViaProperties(resolutionScope, qualifiedName)
}

// Java: JSTypeRegistry.isNameDefinedLocally(StaticTypedScope, String)
func JSTypeRegistry_IsNameDefinedLocally(resolutionScope StaticTypedScope, reference string) bool {
	slot := resolutionScope.GetSlot(reference)
	return slot != nil && slot.GetScope() != nil && slot.GetScope().GetParentScope() != nil
}

// ResolveViaClosureNamespace resolves a named type by checking for the
// longest prefix that matches some Closure namespace, if any, then
// attempting to resolve via properties based on the type of the `exports`
// object in that namespace.
//
// Java: JSTypeRegistry.resolveViaClosureNamespace(String)
func (r *JSTypeRegistry) ResolveViaClosureNamespace(reference string) JSType {
	if r.closureNamespaces.IsEmpty() {
		return nil
	}
	// Find the `exports` type of the longest prefix match of this namespace,
	// if any. Then resolve it via property.
	prefix := reference
	for {
		namespace := r.closureNamespaces.Get(prefix)
		if namespace != nil {
			if namespace.IsLegacy() {
				// Try to resolve this name via registry or properties.
				return nil
			}
			// Always stop resolution here whether successful or not, instead of
			// continuing with resolution via registry or via properties, to
			// match legacy behavior.
			var components []string
			if len(prefix) < len(reference) {
				components = JSTypeRegistry_DOT_SPLITTER_Split(reference[len(prefix)+1:])
			}
			return r.ResolveViaPropertyGivenSlot(
				namespace.Type(), namespace.DefinitionNode(), components)
		}

		lastDot := strings.LastIndexByte(prefix, '.')
		if lastDot < 0 {
			return nil
		}
		prefix = prefix[:lastDot]
	}
}

// Java: JSTypeRegistry.resolveViaProperties(StaticTypedScope, String)
func (r *JSTypeRegistry) ResolveViaProperties(declarationScope StaticTypedScope, qualifiedName string) JSType {
	jutil.CheckArgument(qualifiedName != "")
	dotIdx := strings.IndexByte(qualifiedName, '.')
	rootName := qualifiedName
	if dotIdx != -1 {
		rootName = qualifiedName[:dotIdx]
	}
	slot, _ := declarationScope.GetOwnSlot(rootName).(StaticTypedSlot)
	if slot == nil {
		return nil
	}
	var componentNames []string
	if dotIdx != -1 {
		componentNames = JSTypeRegistry_DOT_SPLITTER_Split(qualifiedName[dotIdx+1:])
	}
	slotType := slot.GetType()
	return r.ResolveViaPropertyGivenSlot(slotType, nil, componentNames)
}

// ResolveViaPropertyGivenSlot resolves a type using a given StaticTypedSlot
// and list of properties on that type.
//
// slotType: the JSType of the slot, possibly null. definitionNode: if known,
// the Node representing the type definition.
//
// Java: JSTypeRegistry.resolveViaPropertyGivenSlot(JSType, Node, List<String>)
func (r *JSTypeRegistry) ResolveViaPropertyGivenSlot(slotType JSType, definitionNode *Node, componentNames []string) JSType {
	if len(componentNames) == 0 {
		typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(definitionNode)
		if typedefType != nil {
			return typedefType
		}
	}

	// If the first component has a type of 'Unknown', then any type names
	// using it should be regarded as silently 'Unknown' rather than be noisy
	// about it.
	if slotType == nil || slotType.IsAllType() || slotType.IsNoType() {
		return nil
	}

	// resolving component by component
	for i := 0; i < len(componentNames); i++ {
		component := componentNames[i]
		parentObj := ObjectType_Cast(slotType)
		if parentObj == nil || component == "" || !parentObj.HasOwnProperty(component) {
			return nil
		}
		if i == len(componentNames)-1 {
			// Look for a typedefTypeProp on the definition node of the last
			// component.
			def := parentObj.GetPropertyDefSite(component)
			typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(def)
			if typedefType != nil {
				return typedefType
			}
		}
		slotType = parentObj.GetPropertyType(component)
	}

	// Translate "constructor" types to "instance" types.
	if slotType == nil {
		return nil
	} else if slotType.IsFunctionType() && (slotType.IsConstructor() || slotType.IsInterface()) {
		return jsTypeOrNil(slotType.ToMaybeFunctionType().GetInstanceType())
	} else if slotType.IsNoObjectType() {
		return r.GetNativeObjectType(JSTypeNative_NO_OBJECT_TYPE)
	} else if enumType, ok := slotType.(*EnumType); ok {
		return enumType.GetElementsType()
	} else {
		return nil
	}
}

// JSTypeRegistry_ResolveTypeFromNodeIfTypedef checks the given Node for a
// typedef annotation, resolving to that type if existent.
//
// Java: JSTypeRegistry.resolveTypeFromNodeIfTypedef(Node)
func JSTypeRegistry_ResolveTypeFromNodeIfTypedef(node *Node) JSType {
	if node == nil {
		return nil
	}
	return node.GetTypedefTypeProp()
}

// Java: JSTypeRegistry.getTypeForScopeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeForScopeInternal(scope StaticScope, name string) JSType {
	rootNode := r.GetRootNodeForScope(scope)
	row := r.scopedNameTable.Get(rootNode)
	if row == nil {
		return nil
	}
	return row.Get(name)
}

// Java: JSTypeRegistry.registerGlobalType(JSType)
func (r *JSTypeRegistry) RegisterGlobalType(type_ JSType) {
	r.Register(nil, type_, type_.String())
}

// Java: JSTypeRegistry.registerGlobalType(JSType, String)
func (r *JSTypeRegistry) RegisterGlobalTypeWithName(type_ JSType, name string) {
	r.Register(nil, type_, name)
}

// Java: JSTypeRegistry.register(StaticScope, JSType, String)
func (r *JSTypeRegistry) Register(scope StaticScope, type_ JSType, name string) {
	JSTypeRegistry_CheckTypeName(name)
	r.RegisterForScope(JSTypeRegistry_GetLookupScope(scope, name), type_, name)
}

// Java: JSTypeRegistry.registerForScope(StaticScope, JSType, String)
func (r *JSTypeRegistry) RegisterForScope(scope StaticScope, type_ JSType, name string) {
	r.scopedNameTablePut(r.GetRootNodeForScope(scope), name, type_)
}

// scopedNameTablePut is Table.put(row, column, value), which rejects nulls.
func (r *JSTypeRegistry) scopedNameTablePut(rowKey *Node, columnKey string, value JSType) {
	jutil.CheckNotNull(rowKey)
	jutil.CheckNotNull(value)
	row := r.scopedNameTable.Get(rowKey)
	if row == nil {
		row = jutil.NewLinkedHashMap[string, JSType]()
		r.scopedNameTable.Put(rowKey, row)
	}
	row.Put(columnKey, value)
}

// Java: JSTypeRegistry.registerNativeType(JSTypeNative, JSType)
func (r *JSTypeRegistry) RegisterNativeType(typeId JSTypeNative, type_ JSType) {
	r.nativeTypes[typeId.Ordinal()] = type_
}

// JSTypeRegistry_IsObjectLiteralThatCanBeSkipped: when t is an object that
// is not the prototype of some class, and its nominal type is Object, and it
// has some properties, we don't need to store these properties in the
// propertyIndex separately.
//
// Java: JSTypeRegistry.isObjectLiteralThatCanBeSkipped(JSType)
func JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(t JSType) bool {
	t = t.RestrictByNotNullOrUndefined()
	return t.IsRecordType() || t.IsLiteralObject()
}

// Java: JSTypeRegistry.registerDroppedPropertiesInUnion(RecordType, RecordType)
func (r *JSTypeRegistry) RegisterDroppedPropertiesInUnion(subtype *RecordType, supertype *RecordType) {
	foundDroppedProperty := false
	subMap := subtype.GetPropertyMap()
	subMap.ensureSorted()
	for i := 0; i < len(subMap.entries); i++ {
		pname := subMap.entries[i].key
		if !supertype.HasProperty(pname) {
			foundDroppedProperty = true
			r.droppedPropertiesOfUnions.Add(pname)
		}
	}
	if foundDroppedProperty {
		superMap := supertype.GetPropertyMap()
		superMap.ensureSorted()
		for i := 0; i < len(superMap.entries); i++ {
			r.propertiesOfSupertypesInUnions.Add(superMap.entries[i].key)
		}
	}
}

// RegisterPropertyOnType tells the type system that owner may have a
// property named propertyName. This allows the registry to keep track of
// what types a property is defined upon.
//
// This is NOT the same as saying that owner must have a property named type.
// ObjectType#hasProperty attempts to minimize false positives ("if we're not
// sure, then don't type check this property"). The type registry, on the
// other hand, should attempt to minimize false negatives ("if this property
// is assigned anywhere in the program, it must show up in the type
// registry").
//
// Only string property keys are registered, not symbol property keys. This
// is largely for historical reasons - JSCompiler has allowed accessing
// properties on a struct via any symbol type, instead of reporting missing
// properties; and even if we do implement missing properties for symbol prop
// accesses, this method is not needed for "STRICT_MISSING_PROPERTIES" errors
// which we might default all symbol props to.
//
// Java: JSTypeRegistry.registerPropertyOnType(String, JSType)
func (r *JSTypeRegistry) RegisterPropertyOnType(propertyName string, type_ JSType) {
	if type_.IsUnionType() {
		for _, alternate := range type_.ToMaybeUnionType().GetAlternates() {
			r.RegisterPropertyOnType(propertyName, alternate)
		}
		return
	}

	if JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(type_) {
		type_ = r.GetSentinelObjectLiteral()
	}

	if objType, ok := type_.(ObjectType); ok && objType.HasReferenceName() {
		set := r.eachRefTypeIndexedByProperty[propertyName]
		if set == nil {
			set = jutil.NewObjectLinkedHashSet[ObjectType]()
			r.eachRefTypeIndexedByProperty[propertyName] = set
		}
		set.Add(objType)
	} else {
		set := r.nonRefTypesIndexedByProperty[propertyName]
		if set == nil {
			set = jutil.NewObjectLinkedHashSet[JSType]()
			r.nonRefTypesIndexedByProperty[propertyName] = set
		}
		set.Add(type_)
	}
}

// JSTypeRegistry_PropDefinitionKind is a tristate value returned from
// canPropertyBeDefined.
type JSTypeRegistry_PropDefinitionKind int32

const (
	JSTypeRegistry_PropDefinitionKind_UNKNOWN     JSTypeRegistry_PropDefinitionKind = iota // The property is not known to be part of this type
	JSTypeRegistry_PropDefinitionKind_KNOWN                                                // The properties is known to be defined on a type or its super types
	JSTypeRegistry_PropDefinitionKind_LOOSE                                                // The property is loosely associated with a type, typically one of its subtypes
	JSTypeRegistry_PropDefinitionKind_LOOSE_UNION                                          // The property is loosely associated with a union type
)

var jsTypeRegistry_PropDefinitionKind_names = [...]string{"UNKNOWN", "KNOWN", "LOOSE", "LOOSE_UNION"}

// JSTypeRegistry_PropDefinitionKind_values is PropDefinitionKind.values().
func JSTypeRegistry_PropDefinitionKind_values() []JSTypeRegistry_PropDefinitionKind {
	return []JSTypeRegistry_PropDefinitionKind{
		JSTypeRegistry_PropDefinitionKind_UNKNOWN, JSTypeRegistry_PropDefinitionKind_KNOWN,
		JSTypeRegistry_PropDefinitionKind_LOOSE, JSTypeRegistry_PropDefinitionKind_LOOSE_UNION,
	}
}

// String is Enum.name().
func (k JSTypeRegistry_PropDefinitionKind) String() string {
	return jsTypeRegistry_PropDefinitionKind_names[k]
}

// CanPropertyBeDefined returns whether the given property can possibly be
// set on the given type.
//
// Java: JSTypeRegistry.canPropertyBeDefined(JSType, String)
func (r *JSTypeRegistry) CanPropertyBeDefined(type_ JSType, propertyName string) JSTypeRegistry_PropDefinitionKind {
	if type_.IsStruct() {
		// We are stricter about "struct" types and only allow access to
		// properties that to the best of our knowledge are available at
		// creation time and specifically not properties only defined on
		// subtypes.

		switch type_.GetPropertyKind(propertyName) {
		case JSType_HasPropertyKind_KNOWN_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_KNOWN
		case JSType_HasPropertyKind_MAYBE_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
		}
		return JSTypeRegistry_PropDefinitionKind_UNKNOWN // ABSENT
	} else {
		if !type_.IsEmptyType() && !type_.IsUnknownType() {
			switch type_.GetPropertyKind(propertyName) {
			case JSType_HasPropertyKind_KNOWN_PRESENT:
				return JSTypeRegistry_PropDefinitionKind_KNOWN
			case JSType_HasPropertyKind_MAYBE_PRESENT:
				return JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
			case JSType_HasPropertyKind_ABSENT:
				// check for loose properties below.
			}
		}

		if set := r.nonRefTypesIndexedByProperty[propertyName]; set != nil {
			for alternative := range set.All() {
				greatestSubtype := alternative.GetGreatestSubtype(type_)
				if !greatestSubtype.IsEmptyType() {
					// We've found a type with this property. Now we just have to make
					// sure it's not a type used for internal bookkeeping.
					maybeRecordType := greatestSubtype.ToMaybeRecordType()
					if maybeRecordType != nil && maybeRecordType.IsSynthetic() {
						continue
					}

					return JSTypeRegistry_PropDefinitionKind_LOOSE
				}
			}
		}
		if set := r.eachRefTypeIndexedByProperty[propertyName]; set != nil {
			for alternative := range set.All() {
				greatestSubtype := alternative.GetGreatestSubtype(type_)
				if !greatestSubtype.IsEmptyType() {
					// We've found a type with this property. Now we just have to make
					// sure it's not a type used for internal bookkeeping.
					maybeRecordType := greatestSubtype.ToMaybeRecordType()
					if maybeRecordType != nil && maybeRecordType.IsSynthetic() {
						continue
					}

					return JSTypeRegistry_PropDefinitionKind_LOOSE
				}
			}
		}

		if rec := type_.ToMaybeRecordType(); rec != nil && r.droppedPropertiesOfUnions.Contains(propertyName) {
			propMap := rec.GetPropertyMap()
			for i := 0; i < len(propMap.entries); i++ {
				if r.propertiesOfSupertypesInUnions.Contains(propMap.entries[i].key) {
					return JSTypeRegistry_PropDefinitionKind_LOOSE
				}
			}
		}
	}
	return JSTypeRegistry_PropDefinitionKind_UNKNOWN
}

// GetEachReferenceTypeWithProperty returns each reference type that has a
// property propertyName defined on it.
//
// Unlike most types in our type system, the collection of types returned
// will not be collapsed. This means that if a type is defined on Object and
// on Array, this method must return [Object, Array]. It would not be correct
// to collapse them to [Object].
//
// Java: JSTypeRegistry.getEachReferenceTypeWithProperty(String)
func (r *JSTypeRegistry) GetEachReferenceTypeWithProperty(propertyName string) []ObjectType {
	if set := r.eachRefTypeIndexedByProperty[propertyName]; set != nil {
		return set.ToSlice()
	}
	return []ObjectType{}
}

// FindCommonSuperObject finds the common supertype of the two given object
// types.
//
// Java: JSTypeRegistry.findCommonSuperObject(ObjectType, ObjectType)
func (r *JSTypeRegistry) FindCommonSuperObject(a ObjectType, b ObjectType) ObjectType {
	var bufA, bufB [16]ObjectType
	stackA := bufA[:0]
	for current := a; current != nil; current = current.GetImplicitPrototype() {
		stackA = append(stackA, current)
	}
	stackB := bufB[:0]
	for current := b; current != nil; current = current.GetImplicitPrototype() {
		stackB = append(stackB, current)
	}

	result := r.GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
	for len(stackA) != 0 && len(stackB) != 0 {
		currentA := stackA[len(stackA)-1]
		stackA = stackA[:len(stackA)-1]
		currentB := stackB[len(stackB)-1]
		stackB = stackB[:len(stackB)-1]
		if currentA.Equals(currentB) {
			result = currentA
		} else {
			return result
		}
	}
	return result
}

// Java: JSTypeRegistry.getSuperStack(ObjectType)
func JSTypeRegistry_GetSuperStack(a ObjectType) []ObjectType {
	stack := make([]ObjectType, 0, 5)
	for current := a; current != nil; current = current.GetImplicitPrototype() {
		stack = append(stack, current)
	}
	return stack
}

// DeclareType records declared global type names. This makes resolution
// faster and more robust in the common case.
//
// name: the name of the type to be recorded. type: the actual type being
// associated with the name. Returns true if this name is not already
// defined, false otherwise.
//
// Java: JSTypeRegistry.declareType(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareType(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(JSTypeRegistry_GetLookupScope(scope, name), name) != nil {
		return false
	}

	r.Register(scope, type_, name)
	return true
}

// DeclareTypeForExactScope records declared global type names. This makes
// resolution faster and more robust in the common case.
//
// Java: JSTypeRegistry.declareTypeForExactScope(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareTypeForExactScope(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(scope, name) != nil {
		return false
	}

	r.RegisterForScope(scope, type_, name)
	return true
}

// OverwriteDeclaredType overrides a declared global type name. Throws an
// exception if this type name hasn't been declared yet.
//
// Java: JSTypeRegistry.overwriteDeclaredType(StaticScope, String, JSType)
func (r *JSTypeRegistry) OverwriteDeclaredType(scope StaticScope, name string, type_ JSType) {
	jutil.CheckStateMsg1(r.IsDeclaredForScope(scope, name), "missing name %s", name)
	r.Register(scope, type_, name)
}

// IsForwardDeclaredType reports whether this is a forward-declared type
// name.
//
// Java: JSTypeRegistry.isForwardDeclaredType(String)
func (r *JSTypeRegistry) IsForwardDeclaredType(name string) bool {
	return r.forwardDeclaredTypes.Contains(name)
}

// GetReadableTypeName first dereferences the JSType to remove
// null/undefined then returns a human-readable type name
//
// Java: JSTypeRegistry.getReadableTypeName(Node)
func (r *JSTypeRegistry) GetReadableTypeName(n *Node) string {
	return r.GetReadableJSTypeName(n, true)
}

// Java: JSTypeRegistry.getReadableTypeNameNoDeref(Node)
func (r *JSTypeRegistry) GetReadableTypeNameNoDeref(n *Node) string {
	return r.GetReadableJSTypeName(n, false)
}

// Java: JSTypeRegistry.getSimpleReadableJSTypeName(JSType)
func JSTypeRegistry_GetSimpleReadableJSTypeName(type_ JSType) *string {
	if _, ok := type_.(*AllType); ok {
		s := type_.String()
		return &s
	} else if _, ok := type_.(ValueType); ok {
		s := type_.String()
		return &s
	} else if type_.IsFunctionPrototypeType() {
		s := type_.String()
		return &s
	} else if _, ok := type_.(ObjectType); ok {
		if !JSTypeRegistry_IsReadableObjectType(type_.ToMaybeObjectType()) {
			return nil
		}

		source := type_.ToObjectType().GetConstructor().GetSource()
		if source != nil {
			jutil.CheckStateMsg1(source.IsFunction() || source.IsClass(), "%s", source)
			readable := source.GetFirstChild().GetOriginalName()
			if readable != nil {
				return readable
			}
		}
		s := type_.String()
		return &s

	} else if _, ok := type_.(*UnionType); ok {
		unionType := type_.ToMaybeUnionType()
		var union *string
		for _, alternate := range unionType.GetAlternates() {
			name := JSTypeRegistry_GetSimpleReadableJSTypeName(alternate)
			if name == nil {
				return nil
			}
			if union == nil {
				u := "(" + *name
				union = &u
			} else {
				u := *union + "|" + *name
				union = &u
			}
		}
		u := jutil.ToJavaString(union) + ")"
		return &u
	}
	return nil
}

// GetReadableJSTypeName: given a node, get a human-readable name for the
// type of that node so that will be easy for the programmer to find the
// original declaration.
//
// For example, if SubFoo's property "bar" might have the human-readable name
// "Foo.prototype.bar".
//
// n: the node. dereference: if true, the type of the node will be
// dereferenced to an Object type, if possible. Prefer to call
// #getReadableTypeName(String) or #getReadableTypeNameNoDeref(String)
// instead of passing this as an argument.
//
// Java: JSTypeRegistry.getReadableJSTypeName(Node, boolean)
func (r *JSTypeRegistry) GetReadableJSTypeName(n *Node, dereference bool) string {
	type_ := r.GetJSTypeOrUnknown(n)
	if dereference {
		autoboxed := type_.Autobox()
		if !autoboxed.IsNoType() {
			type_ = autoboxed
		}
	}

	name := JSTypeRegistry_GetSimpleReadableJSTypeName(type_)
	if name != nil {
		return *name
	}

	// If we're analyzing a GETPROP, the property may be inherited by the
	// prototype chain. So climb the prototype chain and find out where the
	// property was originally defined.
	if n.IsGetProp() {
		objectType := r.GetJSTypeOrUnknown(n.GetFirstChild()).Dereference()
		if objectType != nil {
			propName := n.GetString()
			objectType = objectType.GetClosestDefiningType(propName)

			// Don't show complex function names or anonymous types. Instead, try to
			// get a human-readable type name.
			if JSTypeRegistry_IsReadableObjectType(objectType) {
				return objectType.String() + "." + propName
			}
		}
	}

	if n.IsQualifiedName() {
		return *n.GetQualifiedName()
	} else if type_.IsFunctionType() {
		// Don't show complex function names.
		return "function"
	} else {
		return type_.String()
	}
}

// Java: JSTypeRegistry.isReadableObjectType(ObjectType)
func JSTypeRegistry_IsReadableObjectType(type_ ObjectType) bool {
	if type_ == nil {
		return false
	} else if type_.IsInstanceType() {
		return true
	} else if type_.IsFunctionPrototypeType() {
		return true
	} else if type_.IsEnumElementType() {
		// TODO(b/147236174): This case should be deleted as impossible.
		return type_.ToMaybeEnumElementType().GetConstructor() != nil
	}
	return false
}

// Java: JSTypeRegistry.getJSTypeOrUnknown(Node)
func (r *JSTypeRegistry) GetJSTypeOrUnknown(n *Node) JSType {
	jsType := n.GetJSType()
	if jsType == nil {
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	} else {
		return jsType
	}
}

// Java: JSTypeRegistry.getGlobalType(String)
func (r *JSTypeRegistry) GetGlobalType(jsTypeName string) JSType {
	return r.GetType(nil, jsTypeName)
}

// GetType looks up a native type by name. Returns the corresponding JSType
// object or null it cannot be found.
//
// Java: JSTypeRegistry.getType(StaticScope, String)
func (r *JSTypeRegistry) GetType(scope StaticScope, jsTypeName string) JSType {
	return r.GetTypeInternal(scope, jsTypeName)
}

// GetTypeWithLocation looks up a type by name. To allow for forward
// references to types, an unrecognized string has to be bound to a
// NamedType object that will be resolved later.
//
// scope: a scope for doing type name resolution. jsTypeName: the name
// string. sourceName: the name of the source file where this reference
// appears. lineno: the line number of the reference. Returns a NamedType if
// the string argument is not one of the known types, otherwise the
// corresponding JSType object.
//
// Java: JSTypeRegistry.getType(StaticTypedScope, String, String, int, int)
func (r *JSTypeRegistry) GetTypeWithLocation(
	scope StaticTypedScope, jsTypeName string, sourceName *string, lineno int, charno int) JSType {
	switch jsTypeName {
	case "boolean":
		return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
	case "number":
		return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
	case "bigint":
		return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
	case "string":
		return r.GetNativeType(JSTypeNative_STRING_TYPE)
	case "undefined", "void":
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	}
	// Resolve template type names
	var type_ JSType
	var thisType JSType
	if scope != nil && scope.GetTypeOfThis() != nil {
		thisType = jsTypeOrNil(scope.GetTypeOfThis().ToObjectType())
	}
	if thisType != nil {
		// get the 'last' template type key with this name because if a
		// superclass template key shares the same name as a subclass key, the
		// subclass key is guaranteed to be last.
		key := thisType.GetTemplateTypeMap().GetLastTemplateTypeKeyByName(jsTypeName)
		if key != nil {
			type_ = key
			jutil.CheckStateMsg1(type_.IsTemplateType(), "expected:%s", type_)
			return type_
		}
	}

	// TODO(sdh): The use of "getType" here is incorrect. This currently will
	// pick up a type in an outer scope if it will be shadowed by a local
	// type. But creating a unique NamedType object for every name referenced
	// (even if interned) in every scope would be expensive.
	//
	// Instead perhaps a standard (untyped) scope object might be used to pick
	// the right scope during type construction.
	var untypedScope StaticScope
	if scope != nil {
		untypedScope = scope
	}
	type_ = r.GetType(untypedScope, jsTypeName)
	if type_ == nil {
		// TODO(user): Each instance should support named type creation using
		// interning.
		return r.CreateNamedType(scope, jsTypeName, sourceName, lineno, charno)
	}
	return type_
}

// Java: JSTypeRegistry.getNativeType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeType(typeId JSTypeNative) JSType {
	return r.nativeTypes[typeId.Ordinal()]
}

// Java: JSTypeRegistry.getNativeObjectType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeObjectType(typeId JSTypeNative) ObjectType {
	t := r.GetNativeType(typeId)
	if t == nil {
		return nil
	}
	return t.(ObjectType)
}

// Java: JSTypeRegistry.getNativeFunctionType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeFunctionType(typeId JSTypeNative) FunctionType {
	t := r.GetNativeType(typeId)
	if t == nil {
		return nil
	}
	return t.(FunctionType)
}

// Java: JSTypeRegistry.getResolver()
func (r *JSTypeRegistry) GetResolver() *JSTypeResolver {
	return r.resolver
}

// Java: JSTypeRegistry.evaluateTypeExpressionInGlobalScope(JSTypeExpression)
func (r *JSTypeRegistry) EvaluateTypeExpressionInGlobalScope(expr *JSTypeExpression) JSType {
	return expr.Evaluate(nil, r)
}

// CreateOptionalType creates a type representing optional values of the
// given type. Returns the union of the type and the void type
//
// Java: JSTypeRegistry.createOptionalType(JSType)
func (r *JSTypeRegistry) CreateOptionalType(type_ JSType) JSType {
	if _, ok := type_.(*UnknownType); ok || type_.IsAllType() {
		return type_
	} else {
		if type_ == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index 0"))
		}
		var b UnionType_Builder
		initUnionType_Builder(&b, r)
		b.AddAlternate(type_)
		b.AddAlternate(r.GetNativeType(JSTypeNative_VOID_TYPE))
		return b.Build()
	}
}

// CreateNullableType creates a type representing nullable values of the
// given type. Returns the union of the type and the Null type
//
// Java: JSTypeRegistry.createNullableType(JSType)
func (r *JSTypeRegistry) CreateNullableType(type_ JSType) JSType {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("at index 0"))
	}
	var b UnionType_Builder
	initUnionType_Builder(&b, r)
	b.AddAlternate(type_)
	b.AddAlternate(r.GetNativeType(JSTypeNative_NULL_TYPE))
	return b.Build()
}

// CreateOptionalNullableType creates a nullable and undefine-able value of
// the given type. Returns the union of the type and null and undefined.
//
// Java: JSTypeRegistry.createOptionalNullableType(JSType)
func (r *JSTypeRegistry) CreateOptionalNullableType(type_ JSType) JSType {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("at index 0"))
	}
	var b UnionType_Builder
	initUnionType_Builder(&b, r)
	b.AddAlternate(type_)
	b.AddAlternate(r.GetNativeType(JSTypeNative_VOID_TYPE))
	b.AddAlternate(r.GetNativeType(JSTypeNative_NULL_TYPE))
	return b.Build()
}

// CreateUnionType creates a union type whose variants are the arguments.
//
// Java: JSTypeRegistry.createUnionType(JSType...)
func (r *JSTypeRegistry) CreateUnionType(variants ...JSType) JSType {
	for i, e := range variants {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
	}
	return r.CreateUnionTypeOfList(variants)
}

// Java: JSTypeRegistry.createUnionType(List<? extends JSType>)
func (r *JSTypeRegistry) CreateUnionTypeOfList(variants []JSType) JSType {
	var builder UnionType_Builder
	initUnionType_Builder(&builder, r)
	return builder.AddAlternates(variants).Build()
}

// CreateUnionTypeOfNatives creates a union type whose variants are the
// built-in types specified by the arguments.
//
// Java: JSTypeRegistry.createUnionType(JSTypeNative...)
func (r *JSTypeRegistry) CreateUnionTypeOfNatives(variants ...JSTypeNative) JSType {
	var builder UnionType_Builder
	initUnionType_Builder(&builder, r)
	for _, type_ := range variants {
		builder.AddAlternate(r.GetNativeType(type_))
	}
	return builder.Build()
}

// CreateFunctionType creates a function type.
//
// returnType: the function's return type. parameterTypes: the parameters'
// types.
//
// Java: JSTypeRegistry.createFunctionType(JSType, JSType...)
func (r *JSTypeRegistry) CreateFunctionType(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateFunctionTypeWithParameters(returnType, r.CreateParameters(parameterTypes...))
}

// CreateFunctionTypeWithParameters creates a function type.
//
// parameters: the function's parameters or null to indicate that the
// parameter types are unknown. returnType: the function's return type or
// null to indicate that the return type is unknown.
//
// Java: JSTypeRegistry.createFunctionType(JSType, List<FunctionType.Parameter>)
func (r *JSTypeRegistry) CreateFunctionTypeWithParameters(returnType JSType, parameters []*FunctionType_Parameter) FunctionType {
	return FunctionType_Builder_(r).WithParametersOfList(parameters).WithReturnType(returnType).Build()
}

// CreateFunctionTypeWithVarArgs creates a function type. The last parameter
// type of the function is considered a variable length argument.
//
// Java: JSTypeRegistry.createFunctionTypeWithVarArgs(JSType, JSType...)
func (r *JSTypeRegistry) CreateFunctionTypeWithVarArgs(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateFunctionTypeWithParameters(returnType, r.CreateParametersWithVarArgs(parameterTypes...))
}

// CreateNativeFunctionTypeWithVarArgs creates a function type. The last
// parameter type of the function is considered a variable length argument.
//
// Java: JSTypeRegistry.createNativeFunctionTypeWithVarArgs(JSType, JSType...)
func (r *JSTypeRegistry) CreateNativeFunctionTypeWithVarArgs(returnType JSType, parameterTypes ...JSType) FunctionType {
	return r.CreateNativeFunctionType(returnType, r.CreateParametersWithVarArgs(parameterTypes...))
}

// CreateFunctionTypeWithInstanceType creates a function type in which this
// refers to an object instance.
//
// instanceType: the type of this. returnType: the function's return type.
// parameterTypes: the parameters' types.
//
// Java: JSTypeRegistry.createFunctionTypeWithInstanceType(ObjectType, JSType, List<JSType>)
func (r *JSTypeRegistry) CreateFunctionTypeWithInstanceType(
	instanceType ObjectType, returnType JSType, parameterTypes []JSType) JSType {
	paramsNode := r.CreateParameters(parameterTypes...)
	return FunctionType_Builder_(r).
		WithParametersOfList(paramsNode).
		WithReturnType(returnType).
		WithTypeOfThis(instanceType).
		Build()
}

// CreateParameters creates a tree hierarchy representing a typed argument
// list.
//
// Java: JSTypeRegistry.createParameters(JSType...)
func (r *JSTypeRegistry) CreateParameters(parameterTypes ...JSType) []*FunctionType_Parameter {
	return r.CreateParametersWithLastVarArgs(false, parameterTypes...)
}

// CreateParametersWithLastVarArgs creates a tree hierarchy representing a
// typed argument list.
//
// lastVarArgs: whether the last type should considered as a variable length
// argument. parameterTypes: the parameter types. The last element of this
// array is considered a variable length argument is lastVarArgs is true.
//
// Java: JSTypeRegistry.createParameters(boolean, JSType...)
func (r *JSTypeRegistry) CreateParametersWithLastVarArgs(lastVarArgs bool, parameterTypes ...JSType) []*FunctionType_Parameter {
	builder := NewFunctionParamBuilder(r)
	max := len(parameterTypes) - 1
	for i := 0; i <= max; i++ {
		if lastVarArgs && i == max {
			builder.AddVarArgs(parameterTypes[i])
		} else {
			builder.AddRequiredParams(parameterTypes[i])
		}
	}
	return builder.Build()
}

// CreateParametersWithVarArgs creates a tree hierarchy representing a typed
// argument list. The last parameter type is considered a variable length
// argument.
//
// Java: JSTypeRegistry.createParametersWithVarArgs(JSType...)
func (r *JSTypeRegistry) CreateParametersWithVarArgs(parameterTypes ...JSType) []*FunctionType_Parameter {
	return r.CreateParametersWithLastVarArgs(true, parameterTypes...)
}

// CreateOptionalParameters creates a tree hierarchy representing a typed
// parameter list in which every parameter is optional.
//
// Java: JSTypeRegistry.createOptionalParameters(JSType...)
func (r *JSTypeRegistry) CreateOptionalParameters(parameterTypes ...JSType) []*FunctionType_Parameter {
	builder := NewFunctionParamBuilder(r)
	builder.AddOptionalParams(parameterTypes...)
	return builder.Build()
}

// CreateFunctionTypeWithNewReturnType creates a new function type based on
// an existing function type but with a new return type.
//
// Java: JSTypeRegistry.createFunctionTypeWithNewReturnType(FunctionType, JSType)
func (r *JSTypeRegistry) CreateFunctionTypeWithNewReturnType(existingFunctionType FunctionType, returnType JSType) FunctionType {
	return existingFunctionType.ToBuilder().WithReturnType(returnType).Build()
}

// Java: JSTypeRegistry.createNativeFunctionType(JSType, List<FunctionType.Parameter>)
func (r *JSTypeRegistry) CreateNativeFunctionType(returnType JSType, parameters []*FunctionType_Parameter) FunctionType {
	return FunctionType_Builder_(r).
		WithParametersOfList(parameters).
		WithReturnType(returnType).
		ForNativeType().
		Build()
}

// Java: JSTypeRegistry.buildRecordTypeFromObject(ObjectType)
func (r *JSTypeRegistry) BuildRecordTypeFromObject(objType ObjectType) JSType {
	recType := objType.ToMaybeRecordType()
	// If it can be casted to a record type then return
	if recType != nil {
		return recType
	}
	// TODO(lpino): Handle inherited properties
	propNames := objType.GetOwnPropertyNames()
	// If the type has no properties then return Object
	if len(propNames) == 0 {
		return r.GetNativeType(JSTypeNative_OBJECT_TYPE)
	}
	props := jutil.NewLinkedHashMap[string, JSType]()
	// Otherwise collect the properties and build a record type
	for _, propName := range propNames {
		props.Put(propName, objType.GetPropertyType(propName))
	}
	return r.CreateRecordType(props)
}

// Java: JSTypeRegistry.createRecordType(Map<String, ? extends JSType>)
func (r *JSTypeRegistry) CreateRecordType(props *jutil.LinkedHashMap[string, JSType]) JSType {
	builder := NewRecordTypeBuilder(r)
	for k, v := range props.All() {
		builder.AddProperty(k, v, nil)
	}
	return builder.Build()
}

// CreateObjectType creates an object type.
//
// Java: JSTypeRegistry.createObjectType(String, ObjectType)
func (r *JSTypeRegistry) CreateObjectType(name *string, implicitPrototype ObjectType) ObjectType {
	fields := PrototypeObjectType_BuilderFields{
		registry:          r,
		className:         name,
		implicitPrototype: implicitPrototype,
		templateTypeMap:   r.GetEmptyTemplateTypeMap(),
	}
	return NewPrototypeObjectType(&fields)
}

// CreateAnonymousObjectType creates an anonymous object type.
//
// info: used to mark object literals as structs; can be null
//
// Java: JSTypeRegistry.createAnonymousObjectType(JSDocInfo)
func (r *JSTypeRegistry) CreateAnonymousObjectType(info *JSDocInfo) ObjectType {
	fields := PrototypeObjectType_BuilderFields{
		registry:        r,
		anonymousType:   true,
		templateTypeMap: r.GetEmptyTemplateTypeMap(),
	}
	type_ := NewPrototypeObjectType(&fields)
	type_.SetPrettyPrint(true)
	type_.SetJSDocInfo(info)
	return type_
}

// ResetImplicitPrototype sets the implicit prototype if it's possible to do
// so. There are a few different reasons why this could be a no-op: for
// example, numbers can't be implicit prototypes, and we don't want to change
// the implicit prototype if other classes have already subclassed this one.
//
// Java: JSTypeRegistry.resetImplicitPrototype(JSType, ObjectType)
func (r *JSTypeRegistry) ResetImplicitPrototype(type_ JSType, newImplicitProto ObjectType) {
	if poType, ok := type_.(PrototypeObjectType); ok {
		poType.ClearCachedValues()
		poType.SetImplicitPrototype(newImplicitProto)
	}
}

// CreateConstructorType creates a constructor function type.
//
// name: the function's name or null to indicate that the function is
// anonymous. source: the node defining this function. Its type
// (Node#getToken()) must be Token#FUNCTION. parameters: the function's
// parameters or null to indicate that the parameter types are unknown.
// returnType: the function's return type or null to indicate that the
// return type is unknown. templateKeys: the templatized types for the class.
// isAbstract: whether the function type represents an abstract class
//
// Java: JSTypeRegistry.createConstructorType(String, Node, List<FunctionType.Parameter>, JSType, ImmutableList<TemplateType>, boolean)
func (r *JSTypeRegistry) CreateConstructorType(
	name *string,
	source *Node,
	parameters []*FunctionType_Parameter,
	returnType JSType,
	templateKeys []*TemplateType,
	isAbstract bool) FunctionType {
	jutil.CheckArgument(source == nil || source.IsFunction() || source.IsClass())
	return FunctionType_Builder_(r).
		ForConstructor().
		WithName(name).
		WithSourceNode(source).
		WithParametersOfList(parameters).
		WithReturnType(returnType).
		WithTemplateKeysOfList(templateKeys).
		WithIsAbstract(isAbstract).
		Build()
}

// Java: JSTypeRegistry.createTemplateType(String)
func (r *JSTypeRegistry) CreateTemplateType(name string) *TemplateType {
	return NewTemplateType(r, name)
}

// Java: JSTypeRegistry.createTemplateType(String, JSType)
func (r *JSTypeRegistry) CreateTemplateTypeWithBound(name string, bound JSType) *TemplateType {
	return NewTemplateTypeWithBound(r, name, bound)
}

// Java: JSTypeRegistry.createTemplateTypeWithTransformation(String, Node)
func (r *JSTypeRegistry) CreateTemplateTypeWithTransformation(name string, expr *Node) *TemplateType {
	return NewTemplateTypeWithTransformation(r, name, expr)
}

// Java: JSTypeRegistry.getEmptyTemplateTypeMap()
func (r *JSTypeRegistry) GetEmptyTemplateTypeMap() *TemplateTypeMap {
	return r.emptyTemplateTypeMap
}

// CreateTemplatizedTypeOfList creates a templatized instance of the
// specified type. Only ObjectTypes can currently be templatized; extend the
// logic in this function when more types can be templatized.
//
// baseType: the type to be templatized. templatizedTypes: a list of the
// template JSTypes. Will be matched by list order to the template keys on
// the base type.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, ImmutableList<JSType>)
func (r *JSTypeRegistry) CreateTemplatizedTypeOfList(baseType ObjectType, templatizedTypes []JSType) *TemplatizedType {
	jutil.CheckNotNull(baseType)
	// Only ObjectTypes can currently be templatized; extend this logic when
	// more types can be templatized.
	return NewTemplatizedType(r, baseType, templatizedTypes)
}

// CreateTemplatizedTypeOfMap creates a templatized instance of the
// specified type. Only ObjectTypes can currently be templatized; extend the
// logic in this function when more types can be templatized.
//
// baseType: the type to be templatized. templatizedTypes: a map from
// TemplateType to corresponding JSType value. Any unfilled TemplateTypes on
// the baseType that are *not* contained in this map will have UNKNOWN_TYPE
// used as their value.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, Map<TemplateType, JSType>)
func (r *JSTypeRegistry) CreateTemplatizedTypeOfMap(
	baseType ObjectType, templatizedTypes *jutil.LinkedHashMap[*TemplateType, JSType]) *TemplatizedType {
	unknownType := r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	return r.CreateTemplatizedTypeOfList(
		baseType,
		JSTypeIterations_MapTypes(
			func(key *TemplateType) JSType { return templatizedTypes.GetOrDefault(key, unknownType) },
			baseType.GetTypeParameters()))
}

// CreateTemplatizedType creates a templatized instance of the specified
// type. Only ObjectTypes can currently be templatized; extend the logic in
// this function when more types can be templatized.
//
// Java: JSTypeRegistry.createTemplatizedType(ObjectType, JSType...)
func (r *JSTypeRegistry) CreateTemplatizedType(baseType ObjectType, templatizedTypes ...JSType) *TemplatizedType {
	return r.CreateTemplatizedTypeOfList(baseType, templateTypeMapCopyOfTypes(templatizedTypes))
}

// CreateNamedType creates a named type.
//
// Java: JSTypeRegistry.createNamedType(StaticTypedScope, String, String, int, int)
func (r *JSTypeRegistry) CreateNamedType(
	scope StaticTypedScope, reference string, sourceName *string, lineno int, charno int) *NamedType {
	return NamedType_Builder_(r, reference).
		SetScope(scope).
		SetResolutionKind(NamedType_ResolutionKind_TYPE_NAME).
		SetErrorReportingLocation(sourceName, lineno, charno).
		Build()
}

// IdentifyNonNullableName identifies the name of a typedef or enum before
// we actually declare it.
//
// Java: JSTypeRegistry.identifyNonNullableName(StaticScope, String)
func (r *JSTypeRegistry) IdentifyNonNullableName(scope StaticScope, name string) {
	lookupScope := JSTypeRegistry_GetLookupScope(scope, name)
	rootNode := r.GetRootNodeForScope(lookupScope)
	names := r.nonNullableTypeNames.Get(rootNode)
	if names == nil {
		names = jutil.NewLinkedHashSet[string]()
		r.nonNullableTypeNames.Put(rootNode, names)
	}
	names.Add(name)
}

// IsNonNullableName identifies the name of a typedef or enum before we
// actually declare it.
//
// Java: JSTypeRegistry.isNonNullableName(StaticScope, String)
func (r *JSTypeRegistry) IsNonNullableName(scope StaticScope, name string) bool {
	scope = JSTypeRegistry_GetLookupScope(scope, name)
	names := r.nonNullableTypeNames.Get(r.GetRootNodeForScope(scope))
	return names != nil && names.Contains(name)
}

// Java: JSTypeRegistry.evaluateTypeExpression(JSTypeExpression, StaticTypedScope)
func (r *JSTypeRegistry) EvaluateTypeExpression(expr *JSTypeExpression, scope StaticTypedScope) JSType {
	return r.CreateTypeFromCommentNodeWithScope(expr.GetRoot(), expr.GetSourceName(), scope)
}

// Java: JSTypeRegistry.createTypeFromCommentNode(Node)
func (r *JSTypeRegistry) CreateTypeFromCommentNode(n *Node) JSType {
	return r.CreateTypeFromCommentNodeWithScope(n, "[internal]", nil)
}

// CreateTypeFromCommentNodeWithScope creates a JSType from the nodes
// representing a type.
//
// n: the node with type info. sourceName: the source file name. scope: a
// scope for doing type name lookups.
//
// Java: JSTypeRegistry.createTypeFromCommentNode(Node, String, StaticTypedScope)
func (r *JSTypeRegistry) CreateTypeFromCommentNodeWithScope(n *Node, sourceName string, scope StaticTypedScope) JSType {
	switch n.GetToken() {
	case Token_LC:
		// Record type.
		return r.CreateRecordTypeFromNodes(n.GetFirstChild(), sourceName, scope)
	case Token_BANG:
		// Not nullable
		child := r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope)
		if namedType, ok := child.(*NamedType); ok {
			return namedType.GetBangType()
		}
		return child.RestrictByNotNullOrUndefined()
	case Token_QMARK:
		// Nullable or unknown
		firstChild := n.GetFirstChild()
		if firstChild == nil {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return r.CreateNullableType(r.CreateTypeFromCommentNodeWithScope(firstChild, sourceName, scope))
	case Token_EQUALS:
		// Optional
		// TODO(b/117162687): stop automatically converting {string=} to
		// {(string|undefined)]}
		return r.CreateOptionalType(r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope))
	case Token_ITER_REST:
		// Var args
		return r.CreateTypeFromCommentNodeWithScope(n.GetFirstChild(), sourceName, scope)
	case Token_STAR:
		// The AllType
		return r.GetNativeType(JSTypeNative_ALL_TYPE)
	case Token_PIPE:
		// Union type
		var buf [8]JSType
		builder := buf[:0]
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			builder = append(builder, jutil.CheckNotNull(r.CreateTypeFromCommentNodeWithScope(child, sourceName, scope)))
		}
		return r.CreateUnionTypeOfList(builder)
	case Token_EMPTY:
		// When the return value of a function is not specified
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	case Token_VOID:
		// Only allowed in the return value of a function.
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	case Token_TYPEOF:
		name := n.GetFirstChild().GetString()
		// TODO(sdh): require var to be const?
		qname := QualifiedName_Of(name)
		root := qname.GetRoot()
		declarationScope := scope.GetTopmostScopeOfEventualDeclaration(root)
		rootSlot := scope.GetSlot(root)
		type_ := scope.LookupQualifiedName(qname)
		if type_ == nil || type_.IsUnknownType() || rootSlot.GetScope() != declarationScope {
			// Create a NamedType via getType so that it will be added to the list
			// of types to eventually resolve if necessary.
			return NamedType_Builder_(r, "typeof "+name).
				SetScope(scope).
				SetResolutionKind(NamedType_ResolutionKind_TYPEOF).
				SetErrorReportingLocationFrom(n).
				Build()
		}
		if type_.IsLiteralObject() {
			scopeType := type_
			type_ = NamedType_Builder_(r, "typeof "+name).
				SetResolutionKind(NamedType_ResolutionKind_NONE).
				SetReferencedType(scopeType).
				Build()
		}
		return type_
	case Token_STRINGLIT:
		nominalType := r.GetTypeWithLocation(scope, n.GetString(), &sourceName, n.GetLineno(), n.GetCharno())
		templateArgs := r.ParseTemplateArgs(nominalType, n, sourceName, scope)

		// Handle forward declared types
		if nominalType.IsNamedType() && !nominalType.IsResolved() {
			if templateArgs != nil {
				nominalType = nominalType.ToMaybeNamedType().ToBuilder().SetTemplateTypes(templateArgs).Build()
			}
			return r.AddNullabilityBasedOnParseContext(n, nominalType, scope)
		}

		objectType, ok := nominalType.(ObjectType)
		if !ok {
			return nominalType
		}

		if templateArgs == nil || !nominalType.IsRawTypeOfTemplatizedType() {
			// TODO(nickreid): This case leaves template parameters unbound if users
			// fail to specify any arguments, allowing raw types to leak into
			// programs.
			return r.AddNullabilityBasedOnParseContext(n, nominalType, scope)
		}

		return r.AddNullabilityBasedOnParseContext(
			n, r.CreateTemplatizedTypeOfList(objectType, templateArgs), scope)
	case Token_FUNCTION:
		var thisType JSType
		isConstructor := false
		current := n.GetFirstChild()
		if current.IsThis() || current.IsNew() {
			contextNode := current.GetFirstChild()

			candidateThisType := r.CreateTypeFromCommentNodeWithScope(contextNode, sourceName, scope)

			// Allow null/undefined 'this' types to indicate that the function is
			// not called in a deliberate context, and 'this' access should raise
			// warnings.
			if candidateThisType.IsNullType() || candidateThisType.IsVoidType() {
				thisType = candidateThisType
			} else if current.IsThis() {
				thisType = candidateThisType.RestrictByNotNullOrUndefined()
			} else if current.IsNew() {
				thisType = jsTypeOrNil(ObjectType_Cast(candidateThisType.RestrictByNotNullOrUndefined()))
				if thisType == nil {
					r.reporter.Warning(
						Msg_JSDOC_FUNCTION_NEWNOTOBJECT.Format(),
						sourceName,
						contextNode.GetLineno(),
						contextNode.GetCharno())
				}
			}

			isConstructor = current.GetToken() == Token_NEW
			current = current.GetNext()
		}

		paramBuilder := NewFunctionParamBuilder(r)

		if current.GetToken() == Token_PARAM_LIST {
			for arg := current.GetFirstChild(); arg != nil; arg = arg.GetNext() {
				if arg.GetToken() == Token_ITER_REST {
					if !arg.HasChildren() {
						paramBuilder.AddVarArgs(r.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
					} else {
						paramBuilder.AddVarArgs(
							r.CreateTypeFromCommentNodeWithScope(arg.GetFirstChild(), sourceName, scope))
					}
				} else {
					type_ := r.CreateTypeFromCommentNodeWithScope(arg, sourceName, scope)
					if arg.GetToken() == Token_EQUALS {
						addSuccess := paramBuilder.AddOptionalParams(type_)
						if !addSuccess {
							r.reporter.Warning(
								Msg_JSDOC_FUNCTION_VARARGS.Format(),
								sourceName,
								arg.GetLineno(),
								arg.GetCharno())
						}
					} else {
						paramBuilder.AddRequiredParams(type_)
					}
				}
			}
			current = current.GetNext()
		}

		returnType := r.CreateTypeFromCommentNodeWithScope(current, sourceName, scope)

		kind := FunctionType_Kind_ORDINARY
		if isConstructor {
			kind = FunctionType_Kind_CONSTRUCTOR
		}
		return FunctionType_Builder_(r).
			WithParametersOfList(paramBuilder.Build()).
			WithReturnType(returnType).
			WithTypeOfThis(thisType).
			WithKind(kind).
			Build()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node in type expression: " + n.String()))
}

// jsTypeStaticScopeOrNil converts a StaticTypedScope to a StaticScope,
// keeping Java null.
func jsTypeStaticScopeOrNil(scope StaticTypedScope) StaticScope {
	if scope == nil {
		return nil
	}
	return scope
}

// Java: JSTypeRegistry.addNullabilityBasedOnParseContext(Node, JSType, StaticScope)
func (r *JSTypeRegistry) AddNullabilityBasedOnParseContext(n *Node, type_ JSType, scope StaticTypedScope) JSType {
	// Other node types may be appropriate in the future.
	jutil.CheckStateMsg1(n.IsName() || n.IsStringLit(), "%s", n)
	jutil.CheckNotNull(type_)

	if r.IsNonNullableName(jsTypeStaticScopeOrNil(scope), n.GetString()) {
		return type_
	} else if type_.IsTemplateType() {
		// Template types represent the substituted type exactly and should not
		// be wrapped.
		return type_
	} else if n.HasParent() && n.GetParent().GetToken() == Token_BANG {
		// Names parsed from beneath a BANG never need nullability added.
		return type_
	} else {
		return r.CreateNullableType(type_)
	}
}

// Java: JSTypeRegistry.parseTemplateArgs(JSType, Node, String, StaticTypedScope)
func (r *JSTypeRegistry) ParseTemplateArgs(
	nominalType JSType, typeNode *Node, sourceName string, scope StaticTypedScope) []JSType {
	typeList := typeNode.GetFirstChild()
	if typeList == nil {
		return nil
	}

	count := 0
	for templateNode := typeList.GetFirstChild(); templateNode != nil; templateNode = templateNode.GetNext() {
		count++
	}
	// TODO(b/138617950): Eliminate the special case for `Object`.
	name := typeNode.GetString()
	isObject := name == "Object" || name == "window.Object"
	capCount := count
	if isObject && count == 1 {
		capCount = 2
	}
	templateArgs := make([]JSType, 0, capCount)
	if isObject && count == 1 {
		templateArgs = append(templateArgs, r.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
	}
	for templateNode := typeList.GetFirstChild(); templateNode != nil; templateNode = templateNode.GetNext() {
		templateArgs = append(templateArgs, jutil.CheckNotNull(r.CreateTypeFromCommentNodeWithScope(templateNode, sourceName, scope)))
	}

	if nominalType.IsNamedType() && !nominalType.IsResolved() {
		// The required number of template args will not be known until
		// resolution, so just return all the args.
		return templateArgs
	}

	requiredTemplateArgCount := nominalType.GetTemplateParamCount()
	if len(templateArgs) <= requiredTemplateArgCount {
		return templateArgs
	}

	if !nominalType.IsUnknownType() &&
		// TODO(b/287880204): delete the following case after cleaning up the
		// codebase
		!r.IsNonNullableName(jsTypeStaticScopeOrNil(scope), name) {
		firstExtraTemplateParam := typeNode.GetFirstChild().GetChildAtIndex(requiredTemplateArgCount)
		message := "Too many template parameters\nFound " +
			jutil.ToJavaString(len(templateArgs)) +
			", required at most " +
			jutil.ToJavaString(requiredTemplateArgCount)
		r.reporter.Warning(
			message,
			sourceName,
			firstExtraTemplateParam.GetLineno(),
			firstExtraTemplateParam.GetCharno())
	}
	clear(templateArgs[requiredTemplateArgCount:])
	return templateArgs[:requiredTemplateArgCount:requiredTemplateArgCount]
}

// CreateRecordTypeFromNodes creates a RecordType from the nodes representing
// said record type.
//
// n: the node with type info. sourceName: the source file name. scope: a
// scope for doing type name lookups.
//
// Java: JSTypeRegistry.createRecordTypeFromNodes(Node, String, StaticTypedScope)
func (r *JSTypeRegistry) CreateRecordTypeFromNodes(n *Node, sourceName string, scope StaticTypedScope) JSType {

	builder := NewRecordTypeBuilder(r)

	// For each of the fields in the record type.
	for fieldTypeNode := n.GetFirstChild(); fieldTypeNode != nil; fieldTypeNode = fieldTypeNode.GetNext() {

		// Get the property's name.
		fieldNameNode := fieldTypeNode
		hasType := false

		if fieldTypeNode.GetToken() == Token_COLON {
			fieldNameNode = fieldTypeNode.GetFirstChild()
			hasType = true
		}

		fieldName := fieldNameNode.GetString()

		// TODO(user): Move this into the lexer/parser. Remove the string literal
		// characters around a field name, if any.
		if jutil.StringStartsWith(fieldName, "'") || jutil.StringStartsWith(fieldName, "\"") {
			fieldName = jutil.StringSubstring(fieldName, 1, jutil.StringLength(fieldName)-1)
		}

		// Get the property's type.
		var fieldType JSType

		if hasType {
			// We have a declared type.
			fieldType = r.CreateTypeFromCommentNodeWithScope(fieldTypeNode.GetLastChild(), sourceName, scope)
		} else {
			// Otherwise, the type is UNKNOWN.
			fieldType = r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}

		builder.AddProperty(fieldName, fieldType, fieldNameNode)
	}

	return builder.Build()
}

// RegisterTemplateTypeNamesInScope registers template types on the given
// scope root. This takes a Node rather than a StaticScope because at the time
// it is called, the scope has not yet been created.
//
// Java: JSTypeRegistry.registerTemplateTypeNamesInScope(Iterable<TemplateType>, Node)
func (r *JSTypeRegistry) RegisterTemplateTypeNamesInScope(keys []*TemplateType, scopeRoot *Node) {
	for _, key := range keys {
		r.scopedNameTablePut(scopeRoot, *key.GetReferenceName(), key)
	}
}

// CreateScopeWithTemplates returns a new scope that includes the given
// template names for type resolution purposes.
//
// Java: JSTypeRegistry.createScopeWithTemplates(StaticTypedScope, Iterable<TemplateType>)
func (r *JSTypeRegistry) CreateScopeWithTemplates(scope StaticTypedScope, templates []*TemplateType) StaticTypedScope {
	return NewJSTypeRegistry_SyntheticTemplateScope(scope, templates)
}

// JSTypeRegistry_SyntheticTemplateScope is a synthetic scope that includes
// template names. This is necessary for resolving template names outside the
// body of templated functions (e.g. when evaluating JSDoc on things assigned
// to a prototype, or the parameter or return types of an annotated
// function), since there is not yet (and may never be) any real scope to
// attach the types to.
type JSTypeRegistry_SyntheticTemplateScope struct {
	delegate StaticTypedScope
	types    PMap[string, *TemplateType]
}

// Java: JSTypeRegistry.SyntheticTemplateScope(StaticTypedScope, Iterable<TemplateType>)
func NewJSTypeRegistry_SyntheticTemplateScope(delegate StaticTypedScope, templates []*TemplateType) *JSTypeRegistry_SyntheticTemplateScope {
	s := &JSTypeRegistry_SyntheticTemplateScope{}
	s.delegate = delegate
	var types PMap[string, *TemplateType]
	if syntheticTemplateScope, ok := delegate.(*JSTypeRegistry_SyntheticTemplateScope); ok {
		types = syntheticTemplateScope.types
	} else {
		types = HamtPMap_Empty[string, *TemplateType]()
	}
	for _, key := range templates {
		types = types.Plus(*key.GetReferenceName(), key)
	}
	s.types = types
	return s
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getRootNode()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetRootNode() *Node {
	return s.delegate.GetRootNode()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getParentScope()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetParentScope() StaticScope {
	return s.delegate.GetParentScope()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getSlot(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetSlot(name string) StaticSlot {
	return s.delegate.GetSlot(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getOwnSlot(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetOwnSlot(name string) StaticSlot {
	return s.delegate.GetOwnSlot(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTypeOfThis()
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTypeOfThis() JSType {
	return s.delegate.GetTypeOfThis()
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTemplateType(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTemplateType(name string) *TemplateType {
	return s.types.Get(name)
}

// Java: JSTypeRegistry.SyntheticTemplateScope.getTopmostScopeOfEventualDeclaration(String)
func (s *JSTypeRegistry_SyntheticTemplateScope) GetTopmostScopeOfEventualDeclaration(name string) StaticScope {
	if s.types.Get(name) != nil {
		return s
	}
	return s.delegate.GetTopmostScopeOfEventualDeclaration(name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *JSTypeRegistry_SyntheticTemplateScope) LookupQualifiedName(qname QualifiedName) JSType {
	return StaticTypedScope_LookupQualifiedName(s, qname)
}

// Java: JSTypeRegistry.nativeConstructorBuilder(String)
func (r *JSTypeRegistry) NativeConstructorBuilder(name string) *FunctionType_Builder {
	return FunctionType_Builder_(r).ForNativeType().ForConstructor().WithName(&name)
}

// Java: JSTypeRegistry.nativeInterface(String, TemplateType...)
func (r *JSTypeRegistry) NativeInterface(name string, templateKeys ...*TemplateType) FunctionType {
	builder := FunctionType_Builder_(r).ForNativeType().ForInterface().WithName(&name)
	if len(templateKeys) > 0 {
		builder.WithTemplateKeys(templateKeys...)
	}
	return builder.Build()
}

// Java: JSTypeRegistry.nativeRecord(String, TemplateType...)
func (r *JSTypeRegistry) NativeRecord(name string, templateKeys ...*TemplateType) FunctionType {
	type_ := r.NativeInterface(name, templateKeys...)
	type_.SetImplicitMatch(true)
	return type_
}

// RegisterNonLegacyClosureNamespace registers a goog.module namespace (that
// does not have goog.module.declareLegacyNamespace)
//
// This allows JSTypeRegistry to resolve types that refer to goog.modules by
// namespace. These have unique handling because they exist only in the type
// space and do not have a corresponding value space value.
//
// Java: JSTypeRegistry.registerNonLegacyClosureNamespace(String, Node, JSType)
func (r *JSTypeRegistry) RegisterNonLegacyClosureNamespace(moduleName string, definitionNode *Node, type_ JSType) {
	r.closureNamespaces.Put(
		moduleName, JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ false, definitionNode, type_))
}

// RegisterLegacyClosureNamespace registers a goog.provide or legacy
// goog.module namespace with the type registry
//
// Java: JSTypeRegistry.registerLegacyClosureNamespace(String)
func (r *JSTypeRegistry) RegisterLegacyClosureNamespace(moduleName string) {
	r.closureNamespaces.Put(moduleName, JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ true, nil, nil))
}

// JSTypeRegistry_ClosureNamespace stores information about a Closure
// namespace.
//
// Java: abstract class JSTypeRegistry.ClosureNamespace, with its generated
// AutoValue_JSTypeRegistry_ClosureNamespace implementation folded in.
type JSTypeRegistry_ClosureNamespace struct {
	isLegacy       bool
	definitionNode *Node
	type_          JSType
}

// IsLegacy returns true if this is a goog.provide'd namespace or a
// goog.module namespace followed by `goog.module.declareLegacyNamespace()
//
// Java: JSTypeRegistry.ClosureNamespace.isLegacy()
func (c *JSTypeRegistry_ClosureNamespace) IsLegacy() bool { return c.isLegacy }

// Java: JSTypeRegistry.ClosureNamespace.definitionNode()
func (c *JSTypeRegistry_ClosureNamespace) DefinitionNode() *Node { return c.definitionNode }

// Java: JSTypeRegistry.ClosureNamespace.type()
func (c *JSTypeRegistry_ClosureNamespace) Type() JSType { return c.type_ }

// Java: JSTypeRegistry.ClosureNamespace.create(boolean, Node, JSType)
func JSTypeRegistry_ClosureNamespace_Create(isLegacy bool, definitionNode *Node, type_ JSType) *JSTypeRegistry_ClosureNamespace {
	return &JSTypeRegistry_ClosureNamespace{isLegacy: isLegacy, definitionNode: definitionNode, type_: type_}
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.toString()
func (c *JSTypeRegistry_ClosureNamespace) String() string {
	return "ClosureNamespace{" +
		"isLegacy=" + jutil.ToJavaString(c.isLegacy) + ", " +
		"definitionNode=" + jutil.ToJavaString(c.definitionNode) + ", " +
		"type=" + jutil.ToJavaString(c.type_) +
		"}"
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.equals(Object)
func (c *JSTypeRegistry_ClosureNamespace) Equals(o any) bool {
	if any(c) == o {
		return true
	}
	if that, ok := o.(*JSTypeRegistry_ClosureNamespace); ok && that != nil {
		return c.isLegacy == that.IsLegacy() &&
			jutil.ObjectsEquals(c.definitionNode, that.DefinitionNode()) &&
			jutil.ObjectsEquals(c.type_, that.Type())
	}
	return false
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.hashCode()
func (c *JSTypeRegistry_ClosureNamespace) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.isLegacy)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.definitionNode)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.type_)
	return h
}

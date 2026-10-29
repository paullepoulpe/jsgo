// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20261021

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
	functionPrototype := NewPrototypeObjectType(&PrototypeObjectType_Builder{
		registry:        r,
		classNameID:     rhinoStringPoolIntern("Function.prototype"),
		nativeType:      true,
		templateTypeMap: r.GetEmptyTemplateTypeMap(),
	})
	r.RegisterNativeType(JSTypeNative_FUNCTION_PROTOTYPE, functionPrototype)

	// The ".prototype" property of the type `Function`.
	//
	// So named because `Function` constructs `?`s. This type is not
	// particularly interesting, since its properties aren't inherited by any
	// other type.
	functionInstancePrototype := NewPrototypeObjectType(&PrototypeObjectType_Builder{
		registry:        r,
		classNameID:     rhinoStringPoolIntern("?.prototype"),
		nativeType:      true,
		templateTypeMap: r.GetEmptyTemplateTypeMap(),
	})
	r.RegisterNativeType(JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE, functionInstancePrototype)

	// `Function`
	//
	// The default implicit prototype of all `FunctionType`s is
	// `functionPrototype`.
	functionType := FunctionType_Builder_(r).
		WithNameString("Function").
		ForConstructor().
		ForNativeType().
		WithParametersOfList(r.CreateParametersWithVarArgs(unknownType)).
		WithTypeOfThis(unknownType).
		WithReturnType(unknownType).
		Build()
	functionType.SetPrototype(functionInstancePrototype, NilNode)
	r.RegisterNativeType(JSTypeNative_FUNCTION_TYPE, functionType)

	// `(typeof Function)`
	//
	// The default implict prototype of all `FunctionType`s is
	// `functionPrototype`.
	functionFunctionType := FunctionType_Builder_(r).
		WithNameString("Function").
		ForConstructor().
		ForNativeType().
		WithParametersOfList(r.CreateParametersWithVarArgs(allType)).
		WithTypeOfThis(functionType).
		// TODO(nickreid): .withReturnsOwnInstanceType()
		Build()
	functionFunctionType.SetPrototype(functionPrototype, NilNode)
	r.RegisterNativeType(JSTypeNative_FUNCTION_FUNCTION_TYPE, functionFunctionType)

	// `Object.prototype`
	objectPrototype := NewPrototypeObjectType(&PrototypeObjectType_Builder{
		registry:        r,
		classNameID:     rhinoStringPoolIntern("Object.prototype"),
		nativeType:      true,
		templateTypeMap: r.GetEmptyTemplateTypeMap(),
	})
	r.RegisterNativeType(JSTypeNative_OBJECT_PROTOTYPE, objectPrototype)

	// Object
	objectFunctionType := r.NativeConstructorBuilder("Object").
		WithParametersOfList(r.CreateOptionalParameters(allType)).
		WithReturnsOwnInstanceType().
		WithTemplateKeys(r.iObjectIndexTemplateKey, r.iObjectElementTemplateKey).
		Build()
	objectFunctionType.SetPrototype(objectPrototype, NilNode)
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
	iTemplateArrayFunctionType.SetPrototype(iTemplateArrayPrototype, NilNode)
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

// Java: JSTypeRegistry.nativeConstructorBuilder(String)
func (r *JSTypeRegistry) NativeConstructorBuilder(name string) *FunctionType_Builder {
	return FunctionType_Builder_(r).ForNativeType().ForConstructor().WithNameString(name)
}

// Java: JSTypeRegistry.nativeInterface(String, TemplateType...)
func (r *JSTypeRegistry) NativeInterface(name string, templateKeys ...*TemplateType) FunctionType {
	builder := FunctionType_Builder_(r).ForNativeType().ForInterface().WithNameString(name)
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

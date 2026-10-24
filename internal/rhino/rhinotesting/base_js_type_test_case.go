// Port of com/google/javascript/rhino/testing/BaseJSTypeTestCase.java @ closure-compiler v20261021

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

package rhinotesting

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// BaseJSTypeTestCase is the base of the JSType unit tests: a registry, a
// test error reporter and the native types.
//
// Go-only: language-adapter — a JUnit test class becomes a struct built per
// test function; @After runs as a t.Cleanup.
type BaseJSTypeTestCase struct {
	Truth

	ErrorReporter *TestErrorReporter
	Registry      *rhino.JSTypeRegistry

	ALL_TYPE                     rhino.JSType
	NO_OBJECT_TYPE               rhino.ObjectType
	NO_TYPE                      rhino.ObjectType
	ARRAY_FUNCTION_TYPE          rhino.FunctionType
	ARRAY_TYPE                   rhino.ObjectType
	BIGINT_NUMBER                rhino.JSType
	BIGINT_NUMBER_OBJECT         rhino.JSType
	BIGINT_NUMBER_STRING         rhino.JSType
	BIGINT_NUMBER_STRING_OBJECT  rhino.JSType
	BIGINT_OBJECT_TYPE           rhino.ObjectType
	BIGINT_TYPE                  rhino.JSType
	BOOLEAN_OBJECT_FUNCTION_TYPE rhino.JSType
	BOOLEAN_OBJECT_TYPE          rhino.ObjectType
	BOOLEAN_TYPE                 rhino.JSType
	CHECKED_UNKNOWN_TYPE         rhino.ObjectType
	DATE_FUNCTION_TYPE           rhino.JSType
	DATE_TYPE                    rhino.ObjectType
	FUNCTION_FUNCTION_TYPE       rhino.FunctionType
	FUNCTION_PROTOTYPE           rhino.ObjectType
	FUNCTION_TYPE                rhino.FunctionType
	GREATEST_FUNCTION_TYPE       rhino.JSType
	LEAST_FUNCTION_TYPE          rhino.JSType
	NULL_TYPE                    rhino.JSType
	NUMBER_OBJECT_FUNCTION_TYPE  rhino.JSType
	NUMBER_OBJECT_TYPE           rhino.ObjectType
	NUMBER_STRING                rhino.JSType
	NUMBER_STRING_BOOLEAN        rhino.JSType
	NUMBER_TYPE                  rhino.JSType
	OBJECT_FUNCTION_TYPE         rhino.FunctionType
	NULL_VOID                    rhino.JSType
	OBJECT_NUMBER_STRING         rhino.JSType
	OBJECT_NUMBER_STRING_BOOLEAN rhino.JSType
	OBJECT_PROTOTYPE             rhino.JSType
	OBJECT_TYPE                  rhino.ObjectType
	READONLY_ARRAY_TYPE          rhino.ObjectType
	REGEXP_FUNCTION_TYPE         rhino.JSType
	REGEXP_TYPE                  rhino.ObjectType
	STRING_OBJECT_FUNCTION_TYPE  rhino.JSType
	STRING_OBJECT_TYPE           rhino.ObjectType
	STRING_TYPE                  rhino.JSType
	SYMBOL_OBJECT_TYPE           rhino.ObjectType
	SYMBOL_TYPE                  rhino.JSType
	UNKNOWN_TYPE                 rhino.ObjectType
	VOID_TYPE                    rhino.JSType

	NATIVE_PROPERTIES_COUNT int
}

// Java: BaseJSTypeTestCase.LINE_JOINER
func baseJSTypeTestCase_LINE_JOINER(lines []string) string {
	return strings.Join(lines, "\n")
}

// NewBaseJSTypeTestCase returns a test case with a fresh registry.
//
// Java: BaseJSTypeTestCase()
func NewBaseJSTypeTestCase(t testing.TB) *BaseJSTypeTestCase {
	return NewBaseJSTypeTestCaseWithRegistry(t, nil)
}

// Java: BaseJSTypeTestCase(JSTypeRegistry)
func NewBaseJSTypeTestCaseWithRegistry(t testing.TB, registry *rhino.JSTypeRegistry) *BaseJSTypeTestCase {
	b := &BaseJSTypeTestCase{Truth: Truth{T: t}, ErrorReporter: NewTestErrorReporter()}
	if registry == nil {
		b.Registry = rhino.NewJSTypeRegistryWithForwardDeclaredTypes(b.ErrorReporter, jutil.NewLinkedHashSet[string]())
	} else {
		b.Registry = registry
	}
	b.InitTypes()
	t.Cleanup(b.ValidateWarningsAndErrors)
	return b
}

// Java: BaseJSTypeTestCase.resetRegistryWithForwardDeclaredName(String)
func (b *BaseJSTypeTestCase) ResetRegistryWithForwardDeclaredName(name string) {
	b.Registry = rhino.NewJSTypeRegistryWithForwardDeclaredTypes(b.ErrorReporter, jutil.NewLinkedHashSet(name))
	b.InitTypes()
}

// Java: BaseJSTypeTestCase.validateWarningsAndErrors()
func (b *BaseJSTypeTestCase) ValidateWarningsAndErrors() {
	b.T.Helper()
	b.ErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(b.T)
}

// Java: BaseJSTypeTestCase.initTypes()
func (b *BaseJSTypeTestCase) InitTypes() {
	b.ALL_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	b.NO_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	b.NO_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_NO_TYPE)
	b.ARRAY_FUNCTION_TYPE = b.Registry.GetNativeFunctionType(rhino.JSTypeNative_ARRAY_FUNCTION_TYPE)
	b.ARRAY_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	b.BIGINT_NUMBER = b.Registry.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER)
	b.BIGINT_NUMBER_OBJECT = b.Registry.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_OBJECT)
	b.BIGINT_NUMBER_STRING = b.Registry.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING)
	b.BIGINT_NUMBER_STRING_OBJECT = b.Registry.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING_OBJECT)
	b.BIGINT_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	b.BIGINT_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	b.BOOLEAN_OBJECT_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE)
	b.BOOLEAN_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE)
	b.BOOLEAN_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	b.CHECKED_UNKNOWN_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	b.DATE_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_DATE_FUNCTION_TYPE)
	b.DATE_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_DATE_TYPE)
	b.FUNCTION_FUNCTION_TYPE = b.Registry.GetNativeFunctionType(rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE)
	b.FUNCTION_PROTOTYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_FUNCTION_PROTOTYPE)
	b.FUNCTION_TYPE = b.Registry.GetNativeFunctionType(rhino.JSTypeNative_FUNCTION_TYPE)
	b.GREATEST_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_GREATEST_FUNCTION_TYPE)
	b.LEAST_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_LEAST_FUNCTION_TYPE)
	b.NULL_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	b.NUMBER_OBJECT_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE)
	b.NUMBER_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_NUMBER_OBJECT_TYPE)
	b.NUMBER_STRING = b.Registry.GetNativeType(rhino.JSTypeNative_NUMBER_STRING)
	b.NUMBER_STRING_BOOLEAN = b.Registry.GetNativeType(rhino.JSTypeNative_NUMBER_STRING_BOOLEAN)
	b.NUMBER_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	b.OBJECT_FUNCTION_TYPE = b.Registry.GetNativeFunctionType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)
	b.NULL_VOID = b.Registry.GetNativeType(rhino.JSTypeNative_NULL_VOID)
	b.OBJECT_PROTOTYPE = b.Registry.GetNativeType(rhino.JSTypeNative_OBJECT_PROTOTYPE)
	b.OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE)
	b.READONLY_ARRAY_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_READONLY_ARRAY_TYPE)
	b.REGEXP_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_REGEXP_FUNCTION_TYPE)
	b.REGEXP_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_REGEXP_TYPE)
	b.STRING_OBJECT_FUNCTION_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE)
	b.STRING_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_STRING_OBJECT_TYPE)
	b.STRING_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	b.SYMBOL_OBJECT_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_SYMBOL_OBJECT_TYPE)
	b.SYMBOL_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
	b.UNKNOWN_TYPE = b.Registry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
	b.VOID_TYPE = b.Registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)

	b.OBJECT_NUMBER_STRING = b.Registry.CreateUnionType(b.OBJECT_TYPE, b.NUMBER_TYPE, b.STRING_TYPE)
	b.OBJECT_NUMBER_STRING_BOOLEAN = b.Registry.CreateUnionType(b.OBJECT_TYPE, b.NUMBER_TYPE, b.STRING_TYPE, b.BOOLEAN_TYPE)

	BaseJSTypeTestCase_AddNativeProperties(b.Registry)

	b.NATIVE_PROPERTIES_COUNT = b.OBJECT_TYPE.GetPropertiesCount()
}

// AddNativeProperties adds a basic set of properties to the native types.
//
// Java: BaseJSTypeTestCase.addNativeProperties(JSTypeRegistry)
func BaseJSTypeTestCase_AddNativeProperties(registry *rhino.JSTypeRegistry) {
	booleanType := registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	numberType := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	stringType := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	unknownType := registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)

	objectType := registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE)
	arrayType := registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	dateType := registry.GetNativeObjectType(rhino.JSTypeNative_DATE_TYPE)
	regexpType := registry.GetNativeObjectType(rhino.JSTypeNative_REGEXP_TYPE)
	booleanObjectType := registry.GetNativeObjectType(rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE)
	numberObjectType := registry.GetNativeObjectType(rhino.JSTypeNative_NUMBER_OBJECT_TYPE)
	stringObjectType := registry.GetNativeObjectType(rhino.JSTypeNative_STRING_OBJECT_TYPE)

	objectPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "constructor", objectType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "toLocaleString", stringType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "valueOf", unknownType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "hasOwnProperty", booleanType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "isPrototypeOf", booleanType)
	baseJSTypeTestCase_addMethod(registry, objectPrototype, "propertyIsEnumerable", booleanType)

	arrayPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_ARRAY_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "constructor", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "toLocaleString", stringType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "concat", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "join", stringType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "pop", unknownType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "push", numberType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "reverse", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "shift", unknownType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "slice", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "sort", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "splice", arrayType)
	baseJSTypeTestCase_addMethod(registry, arrayPrototype, "unshift", numberType)
	arrayType.DefineDeclaredProperty("length", numberType, rhino.NilNode)

	booleanPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, booleanPrototype, "constructor", booleanObjectType)
	baseJSTypeTestCase_addMethod(registry, booleanPrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, booleanPrototype, "valueOf", booleanType)

	datePrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_DATE_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, datePrototype, "constructor", dateType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toDateString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toTimeString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toLocaleString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toLocaleDateString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toLocaleTimeString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "valueOf", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getTime", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getFullYear", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCFullYear", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getMonth", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCMonth", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getDate", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCDate", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getDay", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCDay", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getHours", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCHours", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getMinutes", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCMinutes", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getSeconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCSeconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getMilliseconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getUTCMilliseconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "getTimezoneOffset", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setTime", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setMilliseconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCMilliseconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setSeconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCSeconds", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setMinutes", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCMinutes", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setHours", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCHours", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setDate", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCDate", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setMonth", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCMonth", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setFullYear", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "setUTCFullYear", numberType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toUTCString", stringType)
	baseJSTypeTestCase_addMethod(registry, datePrototype, "toGMTString", stringType)

	numberPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "constructor", numberObjectType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "toLocaleString", stringType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "valueOf", numberType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "toFixed", stringType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "toExponential", stringType)
	baseJSTypeTestCase_addMethod(registry, numberPrototype, "toPrecision", stringType)

	regexpPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_REGEXP_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, regexpPrototype, "constructor", regexpType)
	baseJSTypeTestCase_addMethod(registry, regexpPrototype, "exec", registry.CreateNullableType(arrayType))
	baseJSTypeTestCase_addMethod(registry, regexpPrototype, "test", booleanType)
	baseJSTypeTestCase_addMethod(registry, regexpPrototype, "toString", stringType)
	regexpType.DefineDeclaredProperty("source", stringType, rhino.NilNode)
	regexpType.DefineDeclaredProperty("global", booleanType, rhino.NilNode)
	regexpType.DefineDeclaredProperty("ignoreCase", booleanType, rhino.NilNode)
	regexpType.DefineDeclaredProperty("multiline", booleanType, rhino.NilNode)
	regexpType.DefineDeclaredProperty("lastIndex", numberType, rhino.NilNode)

	stringPrototype := registry.GetNativeFunctionType(rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE).GetPrototype()
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "constructor", stringObjectType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "toString", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "valueOf", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "charAt", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "charCodeAt", numberType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "concat", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "indexOf", numberType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "lastIndexOf", numberType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "localeCompare", numberType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "match", registry.CreateNullableType(arrayType))
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "replace", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "search", numberType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "slice", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "split", arrayType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "substr", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "substring", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "toLowerCase", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "toLocaleLowerCase", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "toUpperCase", stringType)
	baseJSTypeTestCase_addMethod(registry, stringPrototype, "toLocaleUpperCase", stringType)
	stringObjectType.DefineDeclaredProperty("length", numberType, rhino.NilNode)
}

// Java: BaseJSTypeTestCase.addMethod(JSTypeRegistry, ObjectType, String, JSType)
func baseJSTypeTestCase_addMethod(
	registry *rhino.JSTypeRegistry, receivingType rhino.ObjectType, methodName string, returnType rhino.JSType) {
	receivingType.DefineDeclaredProperty(
		methodName, rhino.FunctionType_Builder_(registry).WithReturnType(returnType).Build(), rhino.NilNode)
}

// Java: BaseJSTypeTestCase.createUnionType(JSType...)
func (b *BaseJSTypeTestCase) CreateUnionType(variants ...rhino.JSType) rhino.JSType {
	return b.Registry.CreateUnionType(variants...)
}

// Java: BaseJSTypeTestCase.createRecordTypeBuilder()
func (b *BaseJSTypeTestCase) CreateRecordTypeBuilder() *rhino.RecordTypeBuilder {
	return rhino.NewRecordTypeBuilder(b.Registry)
}

// Java: BaseJSTypeTestCase.createNullableType(JSType)
func (b *BaseJSTypeTestCase) CreateNullableType(type_ rhino.JSType) rhino.JSType {
	return b.Registry.CreateNullableType(type_)
}

// Java: BaseJSTypeTestCase.createOptionalType(JSType)
func (b *BaseJSTypeTestCase) CreateOptionalType(type_ rhino.JSType) rhino.JSType {
	return b.Registry.CreateOptionalType(type_)
}

// Java: BaseJSTypeTestCase.createTemplatizedType(ObjectType, ImmutableList<JSType>)
func (b *BaseJSTypeTestCase) CreateTemplatizedTypeOfList(baseType rhino.ObjectType, templatizedTypes []rhino.JSType) *rhino.TemplatizedType {
	return b.Registry.CreateTemplatizedTypeOfList(baseType, templatizedTypes)
}

// Java: BaseJSTypeTestCase.createTemplatizedType(ObjectType, JSType...)
func (b *BaseJSTypeTestCase) CreateTemplatizedType(baseType rhino.ObjectType, templatizedType ...rhino.JSType) *rhino.TemplatizedType {
	return b.CreateTemplatizedTypeOfList(baseType, append([]rhino.JSType{}, templatizedType...))
}

// Asserts that a Node representing a type expression resolves to the
// correct JSType.
//
// Java: BaseJSTypeTestCase.assertTypeEquals(JSType, Node)
func (b *BaseJSTypeTestCase) AssertTypeEqualsNode(expected rhino.JSType, actual rhino.Node) {
	b.T.Helper()
	b.AssertTypeEqualsExpr(expected, rhino.NewJSTypeExpression(actual, "<BaseJSTypeTestCase.java>"))
}

// Asserts that a a type expression resolves to the correct JSType.
//
// Java: BaseJSTypeTestCase.assertTypeEquals(JSType, JSTypeExpression)
func (b *BaseJSTypeTestCase) AssertTypeEqualsExpr(expected rhino.JSType, actual *rhino.JSTypeExpression) {
	b.T.Helper()
	b.AssertTypeEquals(expected, b.Resolve(actual))
}

// Java: BaseJSTypeTestCase.assertTypeEquals(JSType, JSType)
func (b *BaseJSTypeTestCase) AssertTypeEquals(a rhino.JSType, bType rhino.JSType) {
	b.T.Helper()
	b.AssertType(bType).IsEqualTo(a)
}

// Java: BaseJSTypeTestCase.assertTypeEquals(String, JSType, JSType)
func (b *BaseJSTypeTestCase) AssertTypeEqualsMsg(msg string, a rhino.JSType, bType rhino.JSType) {
	b.T.Helper()
	b.AssertWithMessage("%s", msg).AboutTypes().That(bType).IsEqualTo(a)
}

// Resolves a type expression, expecting the given warnings.
//
// Java: BaseJSTypeTestCase.resolve(JSTypeExpression, String...)
func (b *BaseJSTypeTestCase) Resolve(n *rhino.JSTypeExpression, warnings ...string) rhino.JSType {
	b.ErrorReporter.ExpectAllWarnings(warnings...)
	return n.Evaluate(nil, b.Registry)
}

// Java: BaseJSTypeTestCase.ALL_NATIVE_EXTERN_TYPES
var BaseJSTypeTestCase_ALL_NATIVE_EXTERN_TYPES = BaseJSTypeTestCase_Lines(
	"/**",
	" * @constructor",
	" * @param {*=} opt_value",
	" * @return {!Object}",
	" */",
	"function Object(opt_value) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {...*} var_args",
	" */",
	"",
	"function Function(var_args) {}",
	"/**",
	" * @constructor",
	" * @param {...*} var_args",
	" * @return {!Array.<?>}",
	" * @template T",
	" */",
	"function Array(var_args) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {*=} opt_value",
	" * @return {boolean}",
	" */",
	"function Boolean(opt_value) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {*=} opt_value",
	" * @return {number}",
	" */",
	"function Number(opt_value) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {?=} opt_yr_num",
	" * @param {?=} opt_mo_num",
	" * @param {?=} opt_day_num",
	" * @param {?=} opt_hr_num",
	" * @param {?=} opt_min_num",
	" * @param {?=} opt_sec_num",
	" * @param {?=} opt_ms_num",
	" * @return {string}",
	" */",
	"function Date(opt_yr_num, opt_mo_num, opt_day_num, opt_hr_num,",
	"    opt_min_num, opt_sec_num, opt_ms_num) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {*=} opt_str",
	" * @return {string}",
	" */",
	"function String(opt_str) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {*=} opt_pattern",
	" * @param {*=} opt_flags",
	" * @return {!RegExp}",
	" */",
	"function RegExp(opt_pattern, opt_flags) {}",
	"",
	"/**",
	" * @constructor",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!Error}",
	" */",
	"function Error(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!EvalError}",
	" */",
	"function EvalError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!RangeError}",
	" */",
	"function RangeError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!ReferenceError}",
	" */",
	"function ReferenceError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!SyntaxError}",
	" */",
	"function SyntaxError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!TypeError}",
	" */",
	"function TypeError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @constructor",
	" * @extends {Error}",
	" * @param {*=} opt_message",
	" * @param {*=} opt_file",
	" * @param {*=} opt_line",
	" * @return {!URIError}",
	" */",
	"function URIError(opt_message, opt_file, opt_line) {}",
	"",
	"/**",
	" * @param {string} progId",
	" * @param {string=} opt_location",
	" * @constructor",
	" */",
	"function ActiveXObject(progId, opt_location) {}")

// Java: BaseJSTypeTestCase.assertTypeNotEquals(JSType, JSType)
func (b *BaseJSTypeTestCase) AssertTypeNotEquals(a rhino.JSType, bType rhino.JSType) {
	b.T.Helper()
	b.AssertType(bType).IsNotEqualTo(a)
}

// Java: BaseJSTypeTestCase.assertCanTestForEqualityWith(JSType, JSType)
func (b *BaseJSTypeTestCase) AssertCanTestForEqualityWith(t1 rhino.JSType, t2 rhino.JSType) {
	b.T.Helper()
	b.AssertThat(t1.CanTestForEqualityWith(t2)).IsTrue()
	b.AssertThat(t2.CanTestForEqualityWith(t1)).IsTrue()
}

// Java: BaseJSTypeTestCase.assertCannotTestForEqualityWith(JSType, JSType)
func (b *BaseJSTypeTestCase) AssertCannotTestForEqualityWith(t1 rhino.JSType, t2 rhino.JSType) {
	b.T.Helper()
	b.AssertThat(t1.CanTestForEqualityWith(t2)).IsFalse()
	b.AssertThat(t2.CanTestForEqualityWith(t1)).IsFalse()
}

// Java: BaseJSTypeTestCase.lines(String...)
func BaseJSTypeTestCase_Lines(lines ...string) string {
	return baseJSTypeTestCase_LINE_JOINER(lines)
}

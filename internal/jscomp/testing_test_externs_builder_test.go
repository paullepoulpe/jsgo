// Port of com/google/javascript/jscomp/testing/TestExternsBuilder.java @ closure-compiler v20260415

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// Go-only: placement — com.google.javascript.jscomp.testing is test-support
// code; it lives in package jscomp's test files (jscomp.testing is part of
// the jscomp SCC, PORTING.md §2).

import "strings"

// Closure definitions that are not technically externs but serve a similar purpose.
const testExternsBuilder_CLOSURE_EXTERNS = `/** @const */ var goog = {};
goog.module = function(ns) {};
goog.module.declareLegacyNamespace = function() {};
goog.module.preventModuleExportSealing = function() {};
/** @return {?} */
goog.module.get = function(ns) {};
goog.provide = function(ns) {};
/** @return {?} */
goog.require = function(ns) {};
/** @return {?} */
goog.requireType = function(ns) {};
goog.requireDynamic = function(ns) {};
goog.loadModule = function(ns) {};
/** @return {?} */
goog.forwardDeclare = function(ns) {};
goog.setTestOnly = function() {};
goog.scope = function(fn) {};
goog.declareModuleId = function(ns) {};
/**
 * @param {string} name
 * @param {T} defaultValue
 * @return {T}
 * @template T
 */
goog.define = function(name, defaultValue) {};
goog.exportSymbol = function(publicName, symbol) {};
goog.inherits = function(childCtor, parentCtor) {
  childCtor.superClass_ = parentCtor.prototype;
};
goog.getMsg = function(str) {};
/**
 * @param {T} symbol
 * @return {T|undefined}
 * @template T
 * @noinline
 */
goog.weakUsage = function(symbol) {};
`

// testExternsBuilder_GetClosureExternsAsSource: there are some rare cases
// where we want to use the closure externs defined above as if they were
// sources.
//
// Using this method avoids the automatic addition of the `@fileoverview`
// comment containing `@externs`, which forces the compiler to treat the code
// as externs definitions.
//
// Go-only: language-adapter — lower-case t, because `go test` takes every
// exported func named Test* in a _test.go file for a test function.
//
// Java: TestExternsBuilder.getClosureExternsAsSource()
func testExternsBuilder_GetClosureExternsAsSource() string {
	return testExternsBuilder_CLOSURE_EXTERNS
}

// Runtime library stubs that are not technically externs but serve a similar purpose.
const testExternsBuilder_JSCOMP_LIBRARIES = `/** @const */
var $jscomp = {};
/**
 * @param {function(new: ?)} subclass
 * @param {function(new: ?)} superclass
 */
$jscomp.inherits = function(subclass, superclass) {};
/**
 * @param {!Iterable<T>} iterable
 * @return {!Array<T>}
 * @template T
 */
$jscomp.arrayFromIterable = function(iterable) {};
/**
 * @param {string|!Iterable<T>|!Iterator<T>|!Arguments} iterable
 * @return {!Iterator<T>}
 * @template T
 */
$jscomp.makeIterator = function(iterable) {};
$jscomp.makeAsyncIterator = function(asyncIterable) {};
/**
 * @param {!Iterator<T>} iterator
 * @return {!Array<T>}
 * @template T
 */
$jscomp.arrayFromIterator = function(iterator) {};
/**
 * @param {function(): !Generator<?>} generatorFunction
 * @return {!Promise<?>}
 */
$jscomp.asyncExecutePromiseGeneratorFunction = function(generatorFunction) {};
/**
* @type {!Global}
*/
var globalThis;
/** @type {!Global} */
$jscomp.global = globalThis;
/** @const {typeof Reflect.construct} */
$jscomp.construct = Reflect.construct;
/** @constructor */
$jscomp.AsyncGeneratorWrapper = function(generator) {};
/** @constructor */
$jscomp.AsyncGeneratorWrapper$ActionRecord = function(action, value) {};
/** @enum {number} */
$jscomp.AsyncGeneratorWrapper$ActionEnum = {
  YIELD_VALUE: 0,
  YIELD_STAR: 1,
  AWAIT_VALUE: 2,
};`

const testExternsBuilder_BIGINT_EXTERNS = `/**
 * @constructor
 * @param {bigint|number|string} arg
 * @return {bigint}
 */
function BigInt(arg) {}

/**
 * @this {BigInt|bigint}
 * @param {number} width
 * @param {bigint} bigint
 * @return {bigint}
 */
BigInt.asIntN = function(width, bigint) {};

/**
 * @this {BigInt|bigint}
 * @param {number} width
 * @param {bigint} bigint
 * @return {bigint}
 */
BigInt.asUintN = function(width, bigint) {};

/**
 * @param {string|Array<string>=} locales
 * @param {?Object=} options
 * @return {string}
 */
BigInt.prototype.toLocaleString = function(locales, options) {};

/**
 * @this {BigInt|bigint}
 * @param {number=} radix
 * @return {string}
 */
BigInt.prototype.toString = function(radix) {};

/**
 * @return {bigint}
 */
BigInt.prototype.valueOf = function() {};
`

const testExternsBuilder_ITERABLE_EXTERNS = `// Symbol is needed for Symbol.iterator
/**
 * @constructor
 * @param {*=} opt_description
 * @return {symbol}
 * @nosideeffects
 */
function Symbol(opt_description) {}

/** @const {!symbol} */ Symbol.iterator;

/**
 * @record
 * @template TYield
 */
function IIterableResult() {};
/** @type {boolean} */
IIterableResult.prototype.done;
/** @type {TYield} */
IIterableResult.prototype.value;

/**
 * @interface
 * @template T, TReturn, TNext
 */
function Iterator() {}
/**
 * @param {T=} value
 * @return {!IIterableResult<T>}
 */
Iterator.prototype.next;

/**
 * @interface
 * @template T, TReturn, TNext
 */
function Iterable() {}

/**
 * @return {!Iterator<T, ?, *>}}
 * @suppress {externsValidation}
 */
Iterable.prototype[Symbol.iterator] = function() {};

/**
 * @interface
 * @extends {Iterator<T, ?, *>}}
 * @extends {Iterable<T, ?, *>}}
 * @template T, TReturn, TNext
 */
function IteratorIterable() {}

/**
 * @interface
 * @extends {IteratorIterable<T, ?, *>}}
 * @template T, TReturn, TNext
 */
function Generator() {}
/**
 * @param {?=} opt_value
 * @return {!IIterableResult<T>}
 * @override
 */
Generator.prototype.next = function(opt_value) {};
/**
 * @param {T} value
 * @return {!IIterableResult<T>}
 */
Generator.prototype.return = function(value) {};
/**
 * @param {?} exception
 * @return {!IIterableResult<T>}
 */
Generator.prototype.throw = function(exception) {};
`

const testExternsBuilder_STRING_EXTERNS = `/**
 * @constructor
 * @implements {Iterable<string>}
 * @param {*=} arg
 * @return {string}
 */
function String(arg) {}
/** @override @return {!Iterator<string>} */
String.prototype[Symbol.iterator] = function() {};
/** @type {number} */
String.prototype.length;
/** @param {number} sliceArg */
String.prototype.slice = function(sliceArg) {};
/**
 * @this {string|!String}
 * @param {*=} opt_separator
 * @param {number=} opt_limit
 * @return {!Array<string>}
 */
String.prototype.split = function(opt_separator, opt_limit) {};
/**
 * @this {string|!String}
 * @param {string} search_string
 * @param {number=} opt_position
 * @return {boolean}
 * @nosideeffects
 */
String.prototype.startsWith = function(search_string, opt_position) {};
/**
 * @this {!String|string}
 * @param {?} regex
 * @param {?} str
 * @param {string=} opt_flags
 * @return {string}
 */
String.prototype.replace = function(regex, str, opt_flags) {};
/**
 * @this {!String|string}
 * @param {number} index
 * @return {string}
 */
String.prototype.charAt = function(index) {};
/**
 * @this {!String|string}
 * @param {*} regexp
 * @return {!Array<string>}
 */
String.prototype.match = function(regexp) {};
/**
 * @this {!String|string}
 * @return {string}
 */
String.prototype.toLowerCase = function() {};

/**
 * @param {number} count
 * @this {!String|string}
 * @return {string}
 * @nosideeffects
 */
String.prototype.repeat = function(count) {};

/**
 * @param {string} searchString
 * @param {number=} position
 * @return {boolean}
 * @nosideeffects
 */
String.prototype.includes = function(searchString, position) {};
`

const testExternsBuilder_FUNCTION_EXTERNS = `/**
 * @constructor
 * @param {...*} var_args
 */
function Function(var_args) {}
/** @type {!Function} */ Function.prototype.apply;
/** @type {!Function} */ Function.prototype.bind;
/** @type {!Function} */ Function.prototype.call;
/** @type {number} */
Function.prototype.length;
/** @type {string} */
Function.prototype.name;
`

const testExternsBuilder_OBJECT_EXTERNS = `/**
 * @record
 * @template THIS
 */
function ObjectPropertyDescriptor() {}
/** @type {(*|undefined)} */
ObjectPropertyDescriptor.prototype.value;
/** @type {(function(this: THIS):?)|undefined} */
ObjectPropertyDescriptor.prototype.get;

/** @type {(function(this: THIS, ?):void)|undefined} */
ObjectPropertyDescriptor.prototype.set;

/** @type {boolean|undefined} */
ObjectPropertyDescriptor.prototype.writable;

/** @type {boolean|undefined} */
ObjectPropertyDescriptor.prototype.enumerable;

/** @type {boolean|undefined} */
ObjectPropertyDescriptor.prototype.configurable;

/**
 * @constructor
 * @param {*=} opt_value
 * @return {!Object}
 * @template K,V
 */
function Object(opt_value) {}

/** @type {?Object} */ Object.prototype.__proto__;
/** @return {string} */
Object.prototype.toString = function() {};
/**
 * @param {*} propertyName
 * @return {boolean}
 */
Object.prototype.hasOwnProperty = function(propertyName) {};
/** @type {?Function} */ Object.prototype.constructor;
/** @return {*} */
Object.prototype.valueOf = function() {};
/**
 * @param {!Object} obj
 * @param {string} prop
 * @return {!ObjectPropertyDescriptor|undefined}
 * @nosideeffects
 */
Object.getOwnPropertyDescriptor = function(obj, prop) {};
/**
 * @param {!Object} obj
 * @param {string | symbol} prop
 * @param {!ObjectPropertyDescriptor} descriptor
 * @return {!Object}
 */
Object.defineProperty = function(obj, prop, descriptor) {};

/**
 * @template T
 * @param {T} obj
 * @param {!Object<string|symbol, !ObjectPropertyDescriptor<T>>} props
 * @return {T}
 */
Object.defineProperties = function(obj, props) {};

/**
 * @param {?Object} proto
 * @param {?Object=} opt_properties
 * @return {!Object}
 */
Object.create = function(proto, opt_properties) {};
/**
 * @param {!Object} obj
 * @param {?} proto
 * @return {!Object}
 */
Object.setPrototypeOf = function(obj, proto) {};
/**
 * @param {!Object} obj
 * @return {?Object}
 * @nosideeffects
 */
Object.getPrototypeOf = function(obj) {};

/**
 * @param {!Object} target
 * @param {...(Object|null|undefined)} var_args
 * @return {!Object}
 */
Object.assign = function(target, var_args) {};

/**
 * @param {T} obj
 * @return {T}
 * @template T
 */
Object.seal = function(obj) {}
`

const testExternsBuilder_REFLECT_EXTERNS = `/** @const */
var Reflect = {}

/**
 * @param {function(new: ?, ...?)} targetConstructorFn
 * @param {!Array<?>} argList
 * @param {function(new: TARGET, ...?)=} opt_newTargetConstructorFn
 * @return {TARGET}
 * @template TARGET
 * @nosideeffects
 */
Reflect.construct = function(
    targetConstructorFn, argList, opt_newTargetConstructorFn) {};

/**
 * @param {!Object} target
 * @param {?Object} proto
 * @return {boolean}
 */
Reflect.setPrototypeOf = function(target, proto) {};
`

const testExternsBuilder_ARRAY_EXTERNS = `/**
 * @interface
 * @template KEY1, VALUE1
 */
function IObject() {};

/**
 * @record
 * @extends IObject<number, VALUE2>
 * @template VALUE2
 */
function IArrayLike() {};

/** @type {number} */
IArrayLike.prototype.length;

/**
 * @template T
 * @record
 * @extends {IArrayLike<T>}
 * @extends {Iterable<T>}
 */
function ReadonlyArray(var_args) {}
/** @type {number} */ ReadonlyArray.prototype.length;
/**
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @this {!IArrayLike<T>|string}
 * @template T,S
 * @return {undefined}
 */
ReadonlyArray.prototype.forEach;
/**
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T,S
 */
ReadonlyArray.prototype.filter;
/**
 * @param {...*} var_args
 * @return {!Array<?>}
 * @this {*}
 */
ReadonlyArray.prototype.concat;
/**
 * @param {?number=} begin Zero-based index at which to begin extraction.
 * @param {?number=} end Zero-based index at which to end extraction.  slice
 * extracts up to but not including end.
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
ReadonlyArray.prototype.slice;

/**
 * @return {!IteratorIterable<T>}
 */
ReadonlyArray.prototype.values;

/**
 * @param {T} searchElement
 * @param {number=} fromIndex
 * @return {boolean}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
ReadonlyArray.prototype.includes;
/**
 * @template T
 * @constructor
 * @implements {ReadonlyArray<T>}
 * @implements {IArrayLike<T>}
 * @implements {Iterable<T>}
 * @param {...*} var_args
 * @return {!Array<?>}
 */
function Array(var_args) {}
/** @override */
Array.prototype[Symbol.iterator] = function() {};
/** @override @type {number} */ Array.prototype.length;
/**
 * @param {*} arr
 * @return {boolean}
 */
Array.isArray = function(arr) {};

/**
 * @param {string|!IArrayLike<T>|!Iterable<T>} arrayLike
 * @param {function(this:S, (string|T), number): R=} mapFn
 * @param {S=} thisObj
 * @return {!Array<R>}
 * @template T,S,R
 */
Array.from = function(arrayLike, mapFn, thisObj) {}

/**
 * @param {...T} var_args
 * @return {number} The new length of the array.
 * @this {!IArrayLike<T>}
 * @template T
 * @modifies {this}
 */
Array.prototype.push = function(var_args) {};
/**
 * @this {!IArrayLike<T>}
 * @return {T}
 * @template T
 */
Array.prototype.shift = function() {};
/**
 * @override
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @this {!IArrayLike<T>|string}
 * @template T,S
 * @return {undefined}
 */
Array.prototype.forEach = function(callback, opt_thisobj) {};
/**
 * @override
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T,S
 */
Array.prototype.filter = function(callback, opt_thisobj) {};
/**
 * @override
 * @param {...*} var_args
 * @return {!Array<?>}
 * @this {*}
 */
Array.prototype.concat = function(var_args) {};
/**
 * @override
 * @param {?number=} begin Zero-based index at which to begin extraction.
 * @param {?number=} end Zero-based index at which to end extraction.  slice
 * extracts up to but not including end.
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
Array.prototype.slice = function(begin, end) {};

/**
 * @override
 * @return {!IteratorIterable<T>}
 */
Array.prototype.values;

/**
 * @override
 * @param {T} searchElement
 * @param {number=} fromIndex
 * @return {boolean}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
Array.prototype.includes = function(searchElement, fromIndex) {};
`

const testExternsBuilder_MAP_EXTERNS = `/**
 * @interface
 * @extends {Iterable<K|V>}
 * @template K, V
 */
function ReadonlyMap() {}
/**
 * @return {!IteratorIterable<K|V>}
 */
ReadonlyMap.prototype.entries = function() {};
/**
 * @constructor @struct
 * @param {?Iterable<!Array<K|V>>|!Array<!Array<K|V>>=} opt_iterable
 * @implements {ReadonlyMap<K, V>}
 * @template K, V
 */
function Map(opt_iterable) {}
/** @override */
Map.prototype[Symbol.iterator] = function() {};
/**
 * @override
 * @return {!IteratorIterable<K|V>}
 */
Map.prototype.entries = function() {};`

const testExternsBuilder_ARGUMENTS_EXTERNS = `/**
 * @constructor
 * @implements {IArrayLike<?>}
 * @implements {Iterable<?>}
 */
function Arguments() {}

/** @type {number} */
Arguments.prototype.length;
/** @override */
Arguments.prototype[Symbol.iterator] = function() {};

/** @type {!Arguments} */
var arguments;
`

const testExternsBuilder_CONSOLE_EXTERNS = `/** @constructor */
function Console() {};

/**
 * @param {...*} var_args
 * @return {undefined}
 */
Console.prototype.log = function(var_args) {};

/** @const {!Console} */
var console;
`

const testExternsBuilder_ALERT_EXTERNS = `/**
 * @param {*} message
 * @return {undefined}
 */
function alert(message) {}
`

const testExternsBuilder_UNDEFINED_EXTERNS = `/** @type {undefined} */
var undefined;
`

const testExternsBuilder_INFINITY_EXTERNS = `/**
 * @type {number}
 * @see http://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Infinity
 * @const
 */
var Infinity;
`

const testExternsBuilder_NAN_EXTERNS = `/**
 * @type {number}
 * @see http://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/NaN
 * @const
 */
var NaN;
`

const testExternsBuilder_PROMISE_EXTERNS = `
/**
 * @typedef {{then: ?}}
 */
var Thenable;


/**
 * @interface
 * @template TYPE
 */
function IThenable() {}


/**
 * @param {?(function(TYPE):VALUE)=} opt_onFulfilled
 * @param {?(function(*): *)=} opt_onRejected
 * @return {RESULT}
 * @template VALUE
 *
 * @template RESULT := type('IThenable',
 * cond(isUnknown(VALUE), unknown(),
 * mapunion(VALUE, (V) =>
 * cond(isTemplatized(V) && sub(rawTypeOf(V), 'IThenable'),
 * templateTypeOf(V, 0),
 * cond(sub(V, 'Thenable'),
 * unknown(),
 * V)))))
 * =:
 */
IThenable.prototype.then = function(opt_onFulfilled, opt_onRejected) {};


/**
 * @param {function(
 * function((TYPE|IThenable<TYPE>|Thenable|null)=),
 * function(*=))} resolver
 * @constructor
 * @implements {IThenable<TYPE>}
 * @template TYPE
 */
function Promise(resolver) {}


/**
 * @param {VALUE=} opt_value
 * @return {RESULT}
 * @template VALUE
 * @template RESULT := type('Promise',
 * cond(isUnknown(VALUE), unknown(),
 * mapunion(VALUE, (V) =>
 * cond(isTemplatized(V) && sub(rawTypeOf(V), 'IThenable'),
 * templateTypeOf(V, 0),
 * cond(sub(V, 'Thenable'),
 * unknown(),
 * V)))))
 * =:
 */
Promise.resolve = function(opt_value) {};


/**
 * @param {*=} opt_error
 * @return {!Promise<?>}
 */
Promise.reject = function(opt_error) {};


/**
 * @param {!Iterable<VALUE>} iterable
 * @return {!Promise<!Array<RESULT>>}
 * @template VALUE
 * @template RESULT := mapunion(VALUE, (V) =>
 * cond(isUnknown(V),
 * unknown(),
 * cond(isTemplatized(V) && sub(rawTypeOf(V), 'IThenable'),
 * templateTypeOf(V, 0),
 * cond(sub(V, 'Thenable'), unknown(), V))))
 * =:
 */
Promise.all = function(iterable) {};


/**
 * @param {!Iterable<VALUE>} iterable
 * @return {!Promise<RESULT>}
 * @template VALUE
 * @template RESULT := mapunion(VALUE, (V) =>
 * cond(isUnknown(V),
 * unknown(),
 * cond(isTemplatized(V) && sub(rawTypeOf(V), 'IThenable'),
 * templateTypeOf(V, 0),
 * cond(sub(V, 'Thenable'), unknown(), V))))
 * =:
 */
Promise.race = function(iterable) {};


/**
 * @param {?(function(this:void, TYPE):VALUE)=} opt_onFulfilled
 * @param {?(function(this:void, *): *)=} opt_onRejected
 * @return {RESULT}
 * @template VALUE
 *
 * @template RESULT := type('Promise',
 * cond(isUnknown(VALUE), unknown(),
 * mapunion(VALUE, (V) =>
 * cond(isTemplatized(V) && sub(rawTypeOf(V), 'IThenable'),
 * templateTypeOf(V, 0),
 * cond(sub(V, 'Thenable'),
 * unknown(),
 * V)))))
 * =:
 * @override
 */
Promise.prototype.then = function(opt_onFulfilled, opt_onRejected) {};


/**
 * @param {function(*): RESULT} onRejected
 * @return {!Promise<RESULT>}
 * @template RESULT
 */
Promise.prototype.catch = function(onRejected) {};


/**
 * @param {function()} callback
 * @return {!Promise<TYPE>}
 */
Promise.prototype.finally = function(callback) {};
`

const testExternsBuilder_I_TEMPLATE_ARRAY_EXTERNS = `/**
 * @constructor
 * @extends {Array<string>}
 */
function ITemplateArray() {}`

const testExternsBuilder_ASYNC_ITERABLE_EXTERNS = `/**
 * @const {symbol}
 */
Symbol.asyncIterator;
/**
 * @interface
 * @template T, TReturn, TNext
 */
function AsyncIterator() {}
/**
 * @param {?=} opt_value
 * @return {!Promise<!IIterableResult<T>>}
 */
AsyncIterator.prototype.next;
/**
 * @interface
 * @template T, TReturn, TNext
 */
function AsyncIterable() {}
/**
 * @return {!AsyncIterator<T,?,*>}
 */
AsyncIterable.prototype[Symbol.asyncIterator] = function() {};
/**
 * @interface
 * @extends {AsyncIterator<T,?,*>}
 * @extends {AsyncIterable<T,?,*>}
 * @template T, TReturn, TNext
 */
function AsyncIteratorIterable() {}
/**
 * @interface
 * @extends {AsyncIteratorIterable<T,?,*>}
 * @template T, TReturn, TNext
 */
function AsyncGenerator() {}
/**
 * @param {?=} opt_value
 * @return {!Promise<!IIterableResult<T>>}
 * @override
 */
AsyncGenerator.prototype.next = function(opt_value) {};
/**
 * @param {T} value
 * @return {!Promise<!IIterableResult<T>>}
 */
AsyncGenerator.prototype.return = function(value) {};
/**
 * @param {?} exception
 * @return {!Promise<!IIterableResult<T>>}
 */
AsyncGenerator.prototype.throw = function(exception) {};`

const testExternsBuilder_ES6_CLASS_TRANSPILATION_EXTERNS = `var $jscomp = {};

/**
 * @param {?} subClass
 * @param {?} superClass
 * @return {?} newClass
 */
$jscomp.inherits = function(subClass, superClass) {};
`

const testExternsBuilder_MATH_EXTERNS = `
/** @const */
var Math = {};

/**
 * @return {number}
 * @nosideeffects
 */
Math.random = function() {};

/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 * @nosideeffects
 */
Math.max = function(a, b) {};

/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 * @nosideeffects
 */
Math.min = function(a, b) {};

/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 * @nosideeffects
 */
Math.pow = function(a, b) {};
`

const testExternsBuilder_REG_EXP_EXTERNS = `/**
 * @constructor
 * @param {*=} opt_pattern
 * @param {*=} opt_flags
 * @return {!RegExp}
 * @throws {SyntaxError} if opt_pattern is an invalid pattern.
 */
function RegExp(opt_pattern, opt_flags) {}
/**
 * @param {*} str The string to search.
 * @return {?RegExpResult}
 */
RegExp.prototype.exec = function(str) {};
/**
 * @constructor
 * @extends {Array<string>}
 */
var RegExpResult = function() {};
/** @type {string} */
RegExp.$1;`

const testExternsBuilder_HTML_ELEMENT_EXTERNS = `/** @constructor */
function Node() {}
/** @constructor @extends {Node} */
function Element() {}
/** @constructor @extends {Element} */
function HTMLElement() {}`

const testExternsBuilder_POLYMER_EXTERNS = `
/**
 * @param {!Object} init
 * @return {!function(new:HTMLElement)}
 */
function Polymer(init) {}

Polymer.ElementMixin = function(mixin) {}

/** @typedef {!Object} */
var PolymerElementProperties;

/** @interface */
function Polymer_ElementMixin() {}
/** @type {string} */
Polymer_ElementMixin.prototype._importPath;

/**
* @interface
* @extends {Polymer_ElementMixin}
*/
function Polymer_LegacyElementMixin(){}
/** @type {boolean} */
Polymer_LegacyElementMixin.prototype.isAttached;
/**
 * @constructor
 * @extends {HTMLElement}
 * @implements {Polymer_LegacyElementMixin}
 */
var PolymerElement = function() {};  // Polymer 1

/** @constructor @extends {HTMLElement} */
Polymer.Element = function() {};  // Polymer 2
`

// TestExternsBuilder builds a string containing standard externs definitions
// for use in testing.
//
// The externs in this file are intentionally far from complete. Please use
// the actual externs definitions (es3.js, etc.) as a model when you need to
// add something for a test case.
type TestExternsBuilder struct {
	includeBigIntExterns                bool
	includeIterableExterns              bool
	includeStringExterns                bool
	includeFunctionExterns              bool
	includeObjectExterns                bool
	includeArrayExterns                 bool
	includeMapExterns                   bool
	includeArgumentsExterns             bool
	includeUndefinedExterns             bool
	includeInfinityExterns              bool
	includeNaNExterns                   bool
	includeITemplateArrayExterns        bool
	includeConsoleExterns               bool
	includeAlertExterns                 bool
	includePromiseExterns               bool
	includeAsyncIterableExterns         bool
	includeEs6ClassTranspilationExterns bool
	includeReflectExterns               bool
	includeClosureExterns               bool
	includeJSCompLibraries              bool
	includeMathExterns                  bool
	includeRegExpExterns                bool
	includeHtmlElementExterns           bool
	includePolymerExterns               bool
	extraExterns                        []string
}

// NewTestExternsBuilder is new TestExternsBuilder().
//
// Java: TestExternsBuilder()
func NewTestExternsBuilder() *TestExternsBuilder {
	return &TestExternsBuilder{extraExterns: []string{}}
}

// Java: TestExternsBuilder.addBigInt()
func (b *TestExternsBuilder) AddBigInt() *TestExternsBuilder {
	b.includeBigIntExterns = true
	return b
}

// Java: TestExternsBuilder.addIterable()
func (b *TestExternsBuilder) AddIterable() *TestExternsBuilder {
	b.includeIterableExterns = true
	return b
}

// Java: TestExternsBuilder.addString()
func (b *TestExternsBuilder) AddString() *TestExternsBuilder {
	b.includeStringExterns = true
	b.AddIterable() // String implements Iterable<string>
	return b
}

// Java: TestExternsBuilder.addFunction()
func (b *TestExternsBuilder) AddFunction() *TestExternsBuilder {
	b.includeFunctionExterns = true
	return b
}

// Java: TestExternsBuilder.addObject()
func (b *TestExternsBuilder) AddObject() *TestExternsBuilder {
	b.includeObjectExterns = true
	b.AddFunction() // Object.prototype.constructor has type {?Function}
	return b
}

// Java: TestExternsBuilder.addArray()
func (b *TestExternsBuilder) AddArray() *TestExternsBuilder {
	b.includeArrayExterns = true
	b.AddIterable() // Array implements Iterable
	return b
}

// Java: TestExternsBuilder.addMap()
func (b *TestExternsBuilder) AddMap() *TestExternsBuilder {
	b.includeMapExterns = true
	b.AddArray() // Map requires array for its ctor arg.
	return b
}

// Java: TestExternsBuilder.addArguments()
func (b *TestExternsBuilder) AddArguments() *TestExternsBuilder {
	b.includeArgumentsExterns = true
	b.AddArray()    // Arguments implements IArrayLike
	b.AddIterable() // Arguments implements Iterable
	return b
}

// Java: TestExternsBuilder.addUndefined()
func (b *TestExternsBuilder) AddUndefined() *TestExternsBuilder {
	b.includeUndefinedExterns = true
	return b
}

// Java: TestExternsBuilder.addInfinity()
func (b *TestExternsBuilder) AddInfinity() *TestExternsBuilder {
	b.includeInfinityExterns = true
	return b
}

// Java: TestExternsBuilder.addNaN()
func (b *TestExternsBuilder) AddNaN() *TestExternsBuilder {
	b.includeNaNExterns = true
	return b
}

// Java: TestExternsBuilder.addITemplateArray()
func (b *TestExternsBuilder) AddITemplateArray() *TestExternsBuilder {
	b.includeITemplateArrayExterns = true
	b.AddArray() // ITemplateArray extends Array
	return b
}

// Java: TestExternsBuilder.addPromise()
func (b *TestExternsBuilder) AddPromise() *TestExternsBuilder {
	b.includePromiseExterns = true
	b.AddIterable() // Promise.all() and Promise.race() need Iterable
	return b
}

// AddConsole adds declaration of `console.log()`
//
// Java: TestExternsBuilder.addConsole()
func (b *TestExternsBuilder) AddConsole() *TestExternsBuilder {
	b.includeConsoleExterns = true
	return b
}

// AddAlert adds declaration of `alert(message)`
//
// Java: TestExternsBuilder.addAlert()
func (b *TestExternsBuilder) AddAlert() *TestExternsBuilder {
	b.includeAlertExterns = true
	return b
}

// Java: TestExternsBuilder.addAsyncIterable()
func (b *TestExternsBuilder) AddAsyncIterable() *TestExternsBuilder {
	b.includeAsyncIterableExterns = true
	b.AddIterable() // IIterableResult + Symbol
	b.AddPromise()  // Promise
	return b
}

// Java: TestExternsBuilder.addReflect()
func (b *TestExternsBuilder) AddReflect() *TestExternsBuilder {
	b.includeReflectExterns = true
	b.AddObject() // Reflect shares many things in common with Object
	return b
}

// AddEs6ClassTranspilationExterns adds the externs needed for successful
// transpilation of ES6 classes without injecting the runtime code.
//
// ES6 class transpilation depends on some runtime code that we often don't
// want to actually generate in test cases, so we use a non-injecting
// compiler and include these externs definitions to keep the type checker
// happy.
//
// Java: TestExternsBuilder.addEs6ClassTranspilationExterns()
func (b *TestExternsBuilder) AddEs6ClassTranspilationExterns() *TestExternsBuilder {
	b.includeEs6ClassTranspilationExterns = true
	b.AddFunction() // need definition of Function.prototype.apply
	return b
}

// Java: TestExternsBuilder.addClosureExterns()
func (b *TestExternsBuilder) AddClosureExterns() *TestExternsBuilder {
	b.includeClosureExterns = true
	return b
}

// Java: TestExternsBuilder.addJSCompLibraries()
func (b *TestExternsBuilder) AddJSCompLibraries() *TestExternsBuilder {
	b.includeJSCompLibraries = true
	b.AddArguments() // need definition of Arguments
	b.AddIterable()  // need definition of Iterable and Array
	b.AddReflect()   // need Reflect.construct
	return b
}

// Java: TestExternsBuilder.addMath()
func (b *TestExternsBuilder) AddMath() *TestExternsBuilder {
	b.includeMathExterns = true
	return b
}

// Java: TestExternsBuilder.addRegExp()
func (b *TestExternsBuilder) AddRegExp() *TestExternsBuilder {
	b.AddArray() // RegExpResult needs definition of Array
	b.includeRegExpExterns = true
	return b
}

// Java: TestExternsBuilder.addHtmlElement()
func (b *TestExternsBuilder) AddHtmlElement() *TestExternsBuilder {
	b.includeHtmlElementExterns = true
	return b
}

// Java: TestExternsBuilder.addPolymer()
func (b *TestExternsBuilder) AddPolymer() *TestExternsBuilder {
	b.AddHtmlElement()
	b.includePolymerExterns = true
	return b
}

// Java: TestExternsBuilder.addExtra(String...)
func (b *TestExternsBuilder) AddExtra(lines ...string) *TestExternsBuilder {
	b.extraExterns = append(b.extraExterns, lines...)
	return b
}

// Java: TestExternsBuilder.build()
func (b *TestExternsBuilder) Build() string {
	externSections := []string{}
	externSections = append(externSections,
		`/**
 * @fileoverview
 * @externs
 */
`)
	if b.includeBigIntExterns {
		externSections = append(externSections, testExternsBuilder_BIGINT_EXTERNS)
	}
	if b.includeIterableExterns {
		externSections = append(externSections, testExternsBuilder_ITERABLE_EXTERNS)
	}
	if b.includeStringExterns {
		externSections = append(externSections, testExternsBuilder_STRING_EXTERNS)
	}
	if b.includeFunctionExterns {
		externSections = append(externSections, testExternsBuilder_FUNCTION_EXTERNS)
	}
	if b.includeObjectExterns {
		externSections = append(externSections, testExternsBuilder_OBJECT_EXTERNS)
	}
	if b.includeArrayExterns {
		externSections = append(externSections, testExternsBuilder_ARRAY_EXTERNS)
	}
	if b.includeMapExterns {
		externSections = append(externSections, testExternsBuilder_MAP_EXTERNS)
	}
	if b.includeReflectExterns {
		externSections = append(externSections, testExternsBuilder_REFLECT_EXTERNS)
	}
	if b.includeArgumentsExterns {
		externSections = append(externSections, testExternsBuilder_ARGUMENTS_EXTERNS)
	}
	if b.includeUndefinedExterns {
		externSections = append(externSections, testExternsBuilder_UNDEFINED_EXTERNS)
	}
	if b.includeInfinityExterns {
		externSections = append(externSections, testExternsBuilder_INFINITY_EXTERNS)
	}
	if b.includeNaNExterns {
		externSections = append(externSections, testExternsBuilder_NAN_EXTERNS)
	}
	if b.includeITemplateArrayExterns {
		externSections = append(externSections, testExternsBuilder_I_TEMPLATE_ARRAY_EXTERNS)
	}
	if b.includeConsoleExterns {
		externSections = append(externSections, testExternsBuilder_CONSOLE_EXTERNS)
	}
	if b.includeAlertExterns {
		externSections = append(externSections, testExternsBuilder_ALERT_EXTERNS)
	}
	if b.includePromiseExterns {
		externSections = append(externSections, testExternsBuilder_PROMISE_EXTERNS)
	}
	if b.includeAsyncIterableExterns {
		externSections = append(externSections, testExternsBuilder_ASYNC_ITERABLE_EXTERNS)
	}
	if b.includeMathExterns {
		externSections = append(externSections, testExternsBuilder_MATH_EXTERNS)
	}
	if b.includeRegExpExterns {
		externSections = append(externSections, testExternsBuilder_REG_EXP_EXTERNS)
	}
	if b.includeEs6ClassTranspilationExterns {
		externSections = append(externSections, testExternsBuilder_ES6_CLASS_TRANSPILATION_EXTERNS)
	}
	if b.includeClosureExterns {
		externSections = append(externSections, testExternsBuilder_CLOSURE_EXTERNS)
	}
	if b.includeJSCompLibraries {
		externSections = append(externSections, testExternsBuilder_JSCOMP_LIBRARIES)
	}
	if b.includeHtmlElementExterns {
		externSections = append(externSections, testExternsBuilder_HTML_ELEMENT_EXTERNS)
	}
	if b.includePolymerExterns {
		externSections = append(externSections, testExternsBuilder_POLYMER_EXTERNS)
	}
	externSections = append(externSections, b.extraExterns...)
	return strings.Join(externSections, "\n")
}

// Java: TestExternsBuilder.buildExternsFile(String)
func (b *TestExternsBuilder) BuildExternsFile(filePath string) *SourceFile {
	externsString := b.Build()
	return SourceFile_FromCode(filePath, externsString)
}

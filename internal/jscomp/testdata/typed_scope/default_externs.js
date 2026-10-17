/**
 * @fileoverview
 * @externs
 */

// Symbol is needed for Symbol.iterator
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
 * @record
 * @template T, TReturn, TNext
 */
function IteratorLike() {}
/**
 * @param {T=} value
 * @return {!IIterableResult<T>}
 */
IteratorLike.prototype.next;
/**
 * @type {((function(T=): !IIterableResult<T>)|undefined)}
 */
IteratorLike.prototype.return;
/**
 * @type {((function(?=): !IIterableResult<T>)|undefined)}
 */
IteratorLike.prototype.throw;

/**
 * @interface
 * @template T, TReturn, TNext
 */
function Iterable() {}

/**
 * @return {!IteratorLike<T, ?, *>}
 * @suppress {externsValidation}
 */
Iterable.prototype[Symbol.iterator] = function() {};

/**
 * @interface
 * @extends {IteratorLike<T, ?, *>}
 * @extends {Iterable<T, ?, *>}
 * @template T, TReturn, TNext
 */
function IteratorIterable() {}

/**
 * @constructor
 * @abstract
 * @implements {IteratorIterable<T, TReturn, TNext>}
 * @template T, TReturn, TNext
 * @return {?}
 */
function Iterator() {}
/**
 * @param {string|!Iterable<T>|!IteratorLike<T>} iterable
 * @return {!Iterator<T>}
 * @template T
 */
Iterator.from = function(iterable) {};
/**
 * @override
 * @param {?=} value
 * @return {!IIterableResult<T>}
 */
Iterator.prototype.next;
/**
 * @type {((function(T=): !IIterableResult<T>)|undefined)}
 */
Iterator.prototype.return;
/**
 * @type {((function(?=): !IIterableResult<T>)|undefined)}
 */
Iterator.prototype.throw;
/**
 * @override
 * @return {!Iterator<T, TReturn, TNext>}
 * @suppress {externsValidation}
 */
Iterator.prototype[Symbol.iterator] = function() {};
/**
 * @template U
 * @param {function(T, number): U} mapperFn
 * @return {!Iterator<U>}
 */
Iterator.prototype.map;
/**
 * @param {function(T, number): *} filtererFn
 * @return {!Iterator<T>}
 */
Iterator.prototype.filter;
/**
 * @param {number} limit
 * @return {!Iterator<T>}
 */
Iterator.prototype.take;
/**
 * @param {number} limit
 * @return {!Iterator<T>}
 */
Iterator.prototype.drop;
/**
 * @template U
 * @param {function(T, number): (!Iterable<U>|!IteratorLike<U>|!Iterator<U>)} mapperFn
 * @return {!Iterator<U>}
 */
Iterator.prototype.flatMap;
/**
 * @template U
 * @param {function(U, T, number): U} reducerFn
 * @param {U=} initialValue
 * @return {U}
 */
Iterator.prototype.reduce;
/**
 * @return {!Array<T>}
 */
Iterator.prototype.toArray;
/**
 * @param {function(T, number): *} callbackFn
 * @return {void}
 */
Iterator.prototype.forEach;
/**
 * @param {function(T, number): *} predicate
 * @return {boolean}
 */
Iterator.prototype.some;
/**
 * @param {function(T, number): *} predicate
 * @return {boolean}
 */
Iterator.prototype.every;
/**
 * @param {function(T, number): *} predicate
 * @return {T|undefined}
 */
Iterator.prototype.find;

/**
 * @interface
 * @extends {IteratorIterable<T, ?, *>}
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
 * @param {T=} value
 * @return {!IIterableResult<T>}
 */
Generator.prototype.return = function(value) {};
/**
 * @param {?=} exception
 * @return {!IIterableResult<T>}
 */
Generator.prototype.throw = function(exception) {};

/**
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

/**
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

/**
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

/**
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

/** @type {undefined} */
var undefined;

/**
 * @constructor
 * @extends {Array<string>}
 */
function ITemplateArray() {}

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

/** @const */ var goog = {};
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

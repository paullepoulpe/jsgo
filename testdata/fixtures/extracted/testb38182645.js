/**
 * @interface
 * @template VALUE
 */
function MyI() {}


/**
 * @constructor
 * @implements {MyI<K|V>}
 * @template K, V
 */
function MyMap() {}


/**
 * @param {!MyMap<string, T>} map
 * @return {T}
 * @template T
 */
function getValueFromNameAndMap(map) {
  return /** @type {?} */ (123);
}
var m = /** @type {!MyMap<string,number>} */ (new MyMap());
var /** null */ n = getValueFromNameAndMap(m);

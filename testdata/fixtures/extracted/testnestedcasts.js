/** @constructor */var T = function() {};
/** @constructor */var V = function() {};
/**
* @param {boolean} b
* @return {T|V}
*/
function f(b) { return b ? new T() : new V(); }
/**
* @param {boolean} b
* @return {boolean|undefined}
*/
function g(b) { return b ? true : undefined; }
/** @return {T} */
function h() {
return /** @type {T} */ (f(/** @type {boolean} */ (g(true))));
}

/** @constructor
 * @extends {Base} */ function A() {}
/** @constructor
 * @extends {Base} */ function B() {}
/** @constructor */ function Base() {}
/** @param {A|B} x
 * @return {B|A} */
function foo(x) { return x; }

/** @constructor
 * @extends {Base} */ function A() {}
/** @constructor
 * @extends {Base} */ function B() {}
/** @param {A|B} x
 * @return {B|A} */
function foo(x) { return x; }
/** @constructor */ function Base() {}

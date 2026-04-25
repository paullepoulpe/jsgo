/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

/** @type {Generic<!Foo>} */ var x;
/** @type {Generic<!Bar>} */ var y;

x = y;
/** @type {null} */ var z1 = y;

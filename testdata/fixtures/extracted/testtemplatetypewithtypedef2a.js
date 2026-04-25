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

function f(/** Generic<!Bar> */ x) {}
/** @type {Generic<!Foo>} */ var x;
f(x);

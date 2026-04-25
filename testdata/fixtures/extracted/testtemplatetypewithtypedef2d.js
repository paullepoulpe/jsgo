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

function f(/** Generic<!Foo> */ x) {}
/** @type {Generic<!Bar>} */ var x;
f(x);

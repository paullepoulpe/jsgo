/**
 * @constructor
 * @struct
 */
function Foo() {}
/** @param {!Foo} foo */
function getIterator(foo) { return foo[Symbol.iterator](); }

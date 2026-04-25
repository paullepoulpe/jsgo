/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;
/** @interface @template U */ function Bar() {};
/** @type {U} */
Bar.prototype.p;

/**
 * @param {!Foo<number>|!Bar<number>} x
 */
var f = function(x) { x.p = 'not a number'; };

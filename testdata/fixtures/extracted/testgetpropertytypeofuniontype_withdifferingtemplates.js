/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;
/** @interface @template U */ function Bar() {};
/** @type {U} */
Bar.prototype.p;

/**
 * @param {!Foo<number>|!Bar<string>} x
 * @return {string}
 */
var f = function(x) { return x.p; };

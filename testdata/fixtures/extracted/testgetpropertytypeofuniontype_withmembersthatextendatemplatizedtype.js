/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;

/** @interface @extends {Foo<number>} */ function Bar() {};
/** @interface @extends {Foo<number>} */ function Baz() {}

/**
 * @param {!Bar|!Baz} x
 * @return {string}
 */
var f = function(x) { return x.p; };

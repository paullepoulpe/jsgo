/** @type {Foo} */ var a;
/** @type {?Foo} */ var b;
/** @type {!Foo} */ var c;
/** @type {!number} */ var d;
/** @type {!function()} */ var e;
/** @type {!{a: number}} */ var f;
/** @param {Array<Foo>} x @param {number=} y @return {Object} */
function g(x, y) {}
/** @this {number} */ function h() {}
/** @this {!number} */ function h2() {}
/** @typedef {Foo|Bar} */ var TD;
/** @enum {!string} */ var EN = {A: 'a'};
/** @template T @param {T} t @return {T} */ function id(t) { return t; }
/** @type {function(new:Foo, this:Bar, typeof Baz)} */ var k;
/** @type {Foo|null} */ var nn = null;
/** @type {Foo} */ x.y = null;
/** @type {Boolean} */ var bo;
/** @param {Number|String} n */ function prim(n) { return new Number(n) + new String(n) + new Boolean(n) + new Object(); }
/** @const {Map<string, Foo>} */ var m = new Map();
/** @type {Foo} */ let lnull = null;

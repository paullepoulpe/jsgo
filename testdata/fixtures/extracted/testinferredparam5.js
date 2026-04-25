/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {...number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @param {number=} x
 * @param {...number} y  */
SubFoo.prototype.bar =
    function(x, y) { f(x); }; (new SubFoo()).bar();

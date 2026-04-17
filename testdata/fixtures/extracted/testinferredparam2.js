/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @return {void} */ SubFoo.prototype.bar =
    function(x) { f(x); }

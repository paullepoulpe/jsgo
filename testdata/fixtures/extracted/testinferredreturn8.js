/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @param {number} x */ SubFoo.prototype.bar =
    function(x) { return 3; }

/** @constructor */ function Foo() {}
/** @return {Object} */ Foo.prototype.bar =
    function() { return {}; };
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {SubFoo}
 * @override */ SubFoo.prototype.bar =
    function() { return new Foo(); }

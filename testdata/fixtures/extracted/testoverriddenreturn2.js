/** @constructor */ function Foo() {}
/** @return {SubFoo} */ Foo.prototype.bar =
    function() { return new SubFoo(); };
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {Foo} x
 * @override */ SubFoo.prototype.bar =
    function() { return new SubFoo(); }

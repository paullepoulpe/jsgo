/** @constructor
 * @template T */ function Foo() {}
/** @return {T} */ Foo.prototype.bar =
    function() { return null; };
/** @constructor
 * @extends {Foo<string>} */ function SubFoo() {}
/** @override */ SubFoo.prototype.bar =
    function() { return 3; }

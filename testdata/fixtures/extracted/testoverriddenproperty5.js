/** @constructor */ function Foo() {}
Foo.prototype.bar = null;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/** @override */ SubFoo.prototype.bar = 3;

/** @interface  @template T */ function Int() {};
/** @return {T} */ Int.prototype.foo = function() {};

/** @constructor @implements {Int<string>} */ function Foo() {};
/** @return {number}  @override */ Foo.prototype.foo = function() {};

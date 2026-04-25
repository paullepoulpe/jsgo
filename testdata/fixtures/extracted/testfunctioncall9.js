/** @constructor
 * @template T
 **/ function Foo() {}
/** @param {T} x */ Foo.prototype.bar = function(x) {}
var foo = /** @type {Foo<string>} */ (new Foo());
foo.bar(3);

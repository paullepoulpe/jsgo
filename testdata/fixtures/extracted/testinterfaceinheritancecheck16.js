/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @return {T} */A.prototype.bar = function() {};
/** @constructor
 @implements {A<string>} */function B() {};
/** @override */B.prototype.foo = function() { return 'string'};
/** @override */B.prototype.bar = function() { return 3 };

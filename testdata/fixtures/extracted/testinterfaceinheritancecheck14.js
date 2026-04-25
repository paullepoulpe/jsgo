/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @interface
 @template U
 @extends {A<U>} */function B() {};
/** @return {U} */B.prototype.bar = function() {};
/** @constructor
 @implements {B<string>} */function C() {};
/** @return {string}
 @override */C.prototype.foo = function() {};
/** @return {string}
 @override */C.prototype.bar = function() {};

/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @interface
 @template U
 @extends {A<U>} */function B() {};
/** @return {U} */B.prototype.bar = function() {};
/** @constructor
 @template V
 @implements {B<V>}
 */function C() {};
/** @return {V}
 @override */C.prototype.foo = function() {};
/** @return {V}
 @override */C.prototype.bar = function() {};

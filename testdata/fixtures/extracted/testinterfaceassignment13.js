/** @interface */var I = function() {};
/** @constructor
 * @implements {I}*/var T = function() {};
/** @constructor */function Super() {};
/** @return {I} */Super.prototype.foo = function() { return new T(); };
/** @constructor
 * @extends {Super} */function Sub() {}
/** @override
 * @return {T} */Sub.prototype.foo = function() { return new T(); };

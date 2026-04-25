/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor */var T = function() {};
/** @return {!I1|!I2|!T} */function f() { return new T(); }
/** @type {!I1} */var i1 = f();

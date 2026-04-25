/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I2} */var T = function() {};
/** @return {!I1|!I2} */function f() { return new T(); }
/** @type {!I1} */var i1 = f();

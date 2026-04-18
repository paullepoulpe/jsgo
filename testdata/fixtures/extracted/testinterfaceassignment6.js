/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I1} */var T = function() {};
/** @type {!I1} */var i1 = new T();
/** @type {!I2} */var i2 = i1;

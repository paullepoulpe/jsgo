/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I1} */var T = function() {};
var t = new T();
/** @type {I1|I2} */var i = t;

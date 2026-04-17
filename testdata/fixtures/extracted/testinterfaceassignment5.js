/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
 * @implements {I1}
 * @implements {I2}*/
var T = function() {};
var t = new T();
/** @type {I1} */var i1 = t;
/** @type {I2} */var i2 = t;

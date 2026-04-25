/** @interface */
var I1 = function() {};
/** @interface */
var I2 = function() {}
/** @interface @extends {I1} @extends {I2} */
var I3 = function() {};
/** @constructor @implements {I3} */
var T = function() {};
var t = new T();
/**
 * @param {I1} x
 * @param {I2} y
 * @param {I3} z
 */
function foo(x,y,z){};
foo(t,t,t)

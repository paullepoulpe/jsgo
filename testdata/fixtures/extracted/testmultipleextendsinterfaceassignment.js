/** @interface */var I1 = function() {};
/** @interface */ var I2 = function() {}
/** @interface
@extends {I1}
@extends {I2}*/
var I3 = function() {};
/** @constructor
@implements {I3}*/var T = function() {};
var t = new T();
/** @type {I1} */var i1 = t;
/** @type {I2} */var i2 = t;
/** @type {I3} */var i3 = t;
i1 = i3;
i2 = i3;

/** @interface */var I = function() {};
/** @constructor
@implements{I}*/var T1 = function() {};
/** @constructor
@extends {T1}*/var T2 = function() {};
/** @return {I} */function f() { return new T2(); }

/** @constructor */ u.A = function() {};
/**
* @constructor
* @extends {u.A}
*/
u.A.A = function() {};
/**
* @constructor
* @extends {u.A}
*/
u.A.B = function() {};
var ab = new u.A.B();
/** @type {!u.A} */ var a = ab;
/** @type {!u.A.A} */ var aa = ab;

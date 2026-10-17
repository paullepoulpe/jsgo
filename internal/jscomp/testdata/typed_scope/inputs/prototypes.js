/** @constructor */ function A() {}
A.prototype = { m: function() {}, n: 3 };
/** @constructor @extends {A} */ function B() {}
B.prototype = new A();
/** @constructor */ function C() {}
/** @type {number} */ C.prototype.x = 1;
C.prototype.y = function() {};
/** @param {number} a */ C.prototype.z;
function D() {}
D.prototype.foo = 1;
/** @const */ var cAlias = C;
/** @constructor */ var CAlias2 = C;
var protoAlias = C.prototype;
/** @type {!Object} */ var api = {};
api.foo;

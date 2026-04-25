/** @const */ var ns1 = {};
/** @const */ ns1.ns2 = {};
/** @constructor */ ns1.ns2.Outer = function C() {};
/** @return {!ns1.ns2.Outer.Inner} */
ns1.ns2.Outer.prototype.method = function() {
  return new ns1.ns2.Outer.Inner();
};
/** @constructor */ ns1.ns2.Outer.Inner = function () {};

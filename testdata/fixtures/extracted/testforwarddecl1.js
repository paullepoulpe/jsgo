/** @const */ var ns = {};
/** @constructor */ ns.Outer = function C() {};
/** @return {!ns.Outer.Inner} */
ns.Outer.prototype.method = function() {
  return new ns.Outer.Inner();
};
/** @constructor */ ns.Outer.Inner = function () {};

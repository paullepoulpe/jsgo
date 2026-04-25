/** @param {number} x */ function f(x) {}
/** @constructor */ function Ctor() {}
Ctor.prototype = {
  method: function() {
    f(this);
  }
};

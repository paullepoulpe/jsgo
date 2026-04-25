/** @const */ var namespace = {};
(function(ns) {
  /**
   * @constructor
   */
   ns.Ctor = function() {};
   /** @type {boolean} */ ns.Ctor.prototype.bar = true;
})(namespace);
/** @param {namespace.Ctor} x
  * @return {number} */ function f(x) { return x.bar; }

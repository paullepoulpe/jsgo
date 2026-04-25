var ns = {};
(function() {
  /** @constructor
  * @param {number} x */
  ns.Foo = function(x) {};
  /** @param {ns.Foo} x */ function f(x) {}
  f(new ns.Foo(true));
})();

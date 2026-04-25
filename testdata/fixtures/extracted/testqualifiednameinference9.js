var ns = {};
ns.ns2 = {};
(function() {
  /** @constructor
  * @param {number} x */
  ns.ns2.Foo = function(x) {};
  /** @param {ns.ns2.Foo} x */ function f(x) {}
  f(new ns.ns2.Foo(true));
})();

var ns = {};
ns.ns2 = {};
(function() {
  /** @interface */
  ns.ns2.Foo = function() {};
  /** @constructor
  * @implements {ns.ns2.Foo} */
  function F() {}
  (new F());
})();

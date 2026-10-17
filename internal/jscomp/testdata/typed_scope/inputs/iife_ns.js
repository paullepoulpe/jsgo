var ns = {};
(function() {
  /** @return {number} */
  ns.method = function() { return 1; };
  ns.inferred = function() { return 2; };
  /** @constructor */
  ns.Klass = function() {};
})();
ns.method();
(function(global) { global.foo = 1; })(this);

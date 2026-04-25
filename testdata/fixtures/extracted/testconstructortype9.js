var ns = {};
ns.create = function() { return function() {}; };
ns.extend = function(x) { return x; };
/** @constructor */ ns.Foo = ns.create();
ns.Foo.prototype = ns.extend({x: 0, y: 0});
/**
 * @param {ns.Foo} foo
 * @return {string}
 */
function f(foo) {
  return foo.x;
}

var ns = {};
ns.create = function() { return function() {}; };
/** @constructor */ ns.Foo = ns.create();
ns.Foo.prototype = {x: 0, y: 0};
/**
 * @param {ns.Foo} foo
 * @return {string}
 */
function f(foo) {
  return foo.x;
}
